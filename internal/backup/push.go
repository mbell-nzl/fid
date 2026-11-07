package backup

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/mbell-nzl/fid/internal/export"
	"github.com/mbell-nzl/fid/internal/journal"
)

const (
	// Branch is the backup repository's branch that holds every folder.
	Branch = "main"
	// ArchiveName and ManifestName are the files in a project's folder.
	ArchiveName  = "export.zip"
	ManifestName = "backup.json"
	// ManifestSchema versions the folder's manifest.
	ManifestSchema = "fid.backup.v1"

	remoteRef   = "refs/heads/" + Branch
	trackingRef = "refs/remotes/origin/" + Branch
)

// manifest describes the archive in a project's folder, so a restore can
// check it before importing. It holds nothing that changes unless the
// export does, so an unchanged ledger makes an unchanged folder.
type manifest struct {
	Schema            string `json:"schema"`
	ProjectID         string `json:"project_id"`
	Archive           string `json:"archive"`
	ArchiveDigest     string `json:"archive_digest"`
	ManifestSchema    string `json:"manifest_schema"`
	EventCount        int    `json:"event_count"`
	TerminalEventHash string `json:"terminal_event_hash"`
	ProjectionSchema  int    `json:"projection_schema"`
	EntryCount        int    `json:"entry_count"`
	GenesisEventHash  string `json:"genesis_event_hash"`
}

// push brings fid's clone up to date with the target, writes the project's
// folder over the target's tip and pushes it, unless nothing changed. A
// failed push leaves only unreferenced objects; the next backup starts
// again from the target's tip.
func (run *attempt) push() error {
	repository := clone{gitDir: filepath.Join(run.directory, "repo.git"), workingDir: run.directory, environment: gitEnvironment(run.options.Environment)}
	if _, err := os.Stat(filepath.Join(repository.gitDir, "HEAD")); err != nil {
		if _, err := runGit(run.directory, repository.environment, nil, "init", "--bare", "-q", "--template=", repository.gitDir); err != nil {
			return failure(KindInternal, err, "create fid's clone of the backup repository")
		}
	}
	parent, err := run.fetch(repository)
	if err != nil {
		return err
	}
	if run.genesis, err = genesisHash(run.options.WorkspaceRoot); err != nil {
		return failure(KindExportFailed, err, "find the ledger's genesis Event")
	}
	if err := run.checkFolder(repository, parent); err != nil {
		return err
	}
	scratch := filepath.Join(run.directory, "export")
	defer func() { _ = os.RemoveAll(scratch) }()
	summary, archive, err := run.export(scratch)
	if err != nil {
		return err
	}
	tree, err := run.tree(repository, parent, summary, archive)
	if err != nil {
		return err
	}
	run.record.ArchiveDigest, run.record.EventCount, run.record.TerminalEventHash = &summary.ArchiveDigest, &summary.EventCount, &summary.TerminalEventHash
	if parent != "" {
		parentTree, err := repository.text("rev-parse", parent+"^{tree}")
		if err != nil {
			return failure(KindInternal, err, "read the backup repository's tree")
		}
		if parentTree == tree {
			run.record.Outcome, run.record.Commit = OutcomeUnchanged, &parent
			return nil
		}
	}
	commit, err := run.commit(repository, parent, tree, summary)
	if err != nil {
		return err
	}
	if _, err := repository.run(nil, "push", "--porcelain", run.target, commit+":"+remoteRef); err != nil {
		if rejected := rejectedPush(err); rejected != "" {
			return &Error{Kind: KindPushRejected, Message: "the backup repository rejected the push", Detail: rejected, Err: err}
		}
		return failure(KindTargetUnreachable, err, "the backup repository is unreachable")
	}
	// The push succeeded, so the backup did. Tracking the new tip only saves
	// the next backup a fetch and keeps its objects from being pruned.
	_, _ = repository.run(nil, "update-ref", trackingRef, commit)
	run.record.Outcome, run.record.Commit = OutcomePushed, &commit
	return nil
}

// fetch returns the target branch's tip, fetching it into fid's clone when
// the clone doesn't already track it, or "" when the branch doesn't exist.
func (run *attempt) fetch(repository clone) (string, error) {
	remote, err := remoteTip(repository, run.target)
	if err != nil || remote == "" {
		return "", err
	}
	if tracked, err := repository.text("rev-parse", "--verify", "-q", trackingRef); err == nil && tracked == remote {
		return remote, nil
	}
	return fetchBranch(repository, run.target)
}

// remoteTip asks target for its branch and returns the tip, or "" when the
// branch doesn't exist.
func remoteTip(repository clone, target string) (string, error) {
	listed, err := repository.run(nil, "ls-remote", target, remoteRef)
	if err != nil {
		return "", failure(KindTargetUnreachable, err, "the backup repository is unreachable")
	}
	for _, line := range strings.Split(listed, "\n") {
		if oid, ref, ok := strings.Cut(line, "\t"); ok && ref == remoteRef {
			return oid, nil
		}
	}
	return "", nil
}

// fetchBranch fetches target's branch into the tracking ref, with any extra
// fetch options, and returns the fetched tip.
func fetchBranch(repository clone, target string, options ...string) (string, error) {
	arguments := append(append([]string{"fetch", "-q", "--no-tags"}, options...), target, "+"+remoteRef+":"+trackingRef)
	if _, err := repository.run(nil, arguments...); err != nil {
		return "", failure(KindTargetUnreachable, err, "the backup repository is unreachable")
	}
	tip, err := repository.text("rev-parse", "--verify", trackingRef+"^{commit}")
	if err != nil {
		return "", failure(KindInternal, err, "read the fetched backup branch")
	}
	return tip, nil
}

// export writes and validates the project's export in scratch, emptied
// first in case a crashed backup left it behind.
func (run *attempt) export(scratch string) (export.Summary, string, error) {
	if err := os.RemoveAll(scratch); err != nil {
		return export.Summary{}, "", failure(KindInternal, err, "clear the export directory")
	}
	if err := os.Mkdir(scratch, 0o700); err != nil {
		return export.Summary{}, "", failure(KindInternal, err, "create the export directory")
	}
	archive := filepath.Join(scratch, ArchiveName)
	summary, err := export.Create(export.CreateOptions{WorkspaceRoot: run.options.WorkspaceRoot, ArchivePath: archive})
	if err != nil {
		return export.Summary{}, "", failure(KindExportFailed, err, "export the ledger")
	}
	if _, err := export.Validate(export.ValidateOptions{ArchivePath: archive, MigrationsDir: run.options.MigrationsDir}); err != nil {
		var invalid *export.InvalidArchiveError
		if errors.As(err, &invalid) {
			codes := make([]string, 0, len(invalid.Checks))
			for _, check := range invalid.Checks {
				codes = append(codes, check.Code+": "+check.Detail)
			}
			err = fmt.Errorf("%w: %s", err, strings.Join(codes, "; "))
		}
		return export.Summary{}, "", failure(KindExportFailed, err, "validate the ledger's export")
	}
	return summary, archive, nil
}

// tree returns the backup repository's tree with the project's folder
// replaced by the new archive and manifest. Other folders are untouched.
func (run *attempt) tree(repository clone, parent string, summary export.Summary, archive string) (string, error) {
	encoded, err := journal.MarshalCanonicalJSON(manifest{Schema: ManifestSchema, ProjectID: run.projectID, Archive: ArchiveName,
		ArchiveDigest: summary.ArchiveDigest, ManifestSchema: summary.ManifestSchema, EventCount: summary.EventCount,
		TerminalEventHash: summary.TerminalEventHash, ProjectionSchema: summary.ProjectionSchema, EntryCount: summary.EntryCount,
		GenesisEventHash: run.genesis})
	if err != nil {
		return "", failure(KindInternal, err, "encode the folder manifest")
	}
	// --no-filters: attributes in the user's git config never touch the
	// archive. Content on stdin has no path, so git applies no filters to it.
	archiveBlob, err := repository.text("hash-object", "-w", "--no-filters", archive)
	if err != nil {
		return "", failure(KindInternal, err, "store the export")
	}
	manifestBlob, err := repository.output(strings.NewReader(string(encoded)+"\n"), "hash-object", "-w", "--stdin")
	if err != nil {
		return "", failure(KindInternal, err, "store the folder manifest")
	}
	folder, err := repository.output(strings.NewReader(fmt.Sprintf("100644 blob %s\t%s\n100644 blob %s\t%s\n", archiveBlob, ArchiveName, manifestBlob, ManifestName)), "mktree")
	if err != nil {
		return "", failure(KindInternal, err, "build the project's folder")
	}
	entries := ""
	if parent != "" {
		listed, err := repository.run(nil, "ls-tree", "-z", parent)
		if err != nil {
			return "", failure(KindInternal, err, "read the backup repository's tree")
		}
		for _, entry := range strings.Split(listed, "\x00") {
			if _, name, ok := strings.Cut(entry, "\t"); ok && name != run.projectID {
				entries += entry + "\x00"
			}
		}
	}
	entries += "040000 tree " + folder + "\t" + run.projectID + "\x00"
	root, err := repository.output(strings.NewReader(entries), "mktree", "-z")
	if err != nil {
		return "", failure(KindInternal, err, "build the backup repository's tree")
	}
	return root, nil
}

func (run *attempt) commit(repository clone, parent, tree string, summary export.Summary) (string, error) {
	message := fmt.Sprintf("fid backup %s\n\nevent_count: %d\nterminal_event_hash: %s\narchive_digest: %s\n",
		run.projectID, summary.EventCount, summary.TerminalEventHash, summary.ArchiveDigest)
	arguments := []string{"commit-tree", "--no-gpg-sign", "-F", "-"}
	if parent != "" {
		arguments = append(arguments, "-p", parent)
	}
	commit, err := repository.output(strings.NewReader(message), append(arguments, tree)...)
	if err != nil {
		return "", failure(KindInternal, err, "commit the backup")
	}
	return commit, nil
}
