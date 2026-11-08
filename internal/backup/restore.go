package backup

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/mbell-nzl/fid/internal/config"
	"github.com/mbell-nzl/fid/internal/export"
)

// Check codes for a damaged folder. An archive-invalid failure carries them
// beside export validation's own codes.
const (
	CheckManifestInvalid = "backup-manifest-invalid"
	CheckArchiveMissing  = "archive-missing"
	CheckArchiveDigest   = "archive-digest-mismatch"
)

// RestoreOptions identifies the project and the environment a restore runs
// in.
type RestoreOptions struct {
	// WorkspaceRoot is the worktree's top level, as project.Discover returns
	// it. The ledger is the one every worktree of the repository shares.
	WorkspaceRoot string
	// Environment is the process environment. It locates user config and the
	// git credentials, and git runs with it as it does for a backup.
	Environment   []string
	MigrationsDir string
}

// Restored is the result fid restore reports.
type Restored struct {
	ProjectID         string `json:"project_id"`
	Commit            string `json:"commit"`
	ArchiveDigest     string `json:"archive_digest"`
	EventCount        int    `json:"event_count"`
	TerminalEventHash string `json:"terminal_event_hash"`
	GenesisEventHash  string `json:"genesis_event_hash"`
}

// folderHeader is what a restore reads from backup.json before it knows
// the folder's format.
type folderHeader struct {
	Schema    string `json:"schema"`
	ProjectID string `json:"project_id"`
}

// Restore imports the project's folder, from the tip of the backup
// repository's branch, into a repository that has no ledger, and returns
// the ledger's path. It fetches into a private temporary repository that it
// deletes afterwards, never fid's backup clone, so it takes no lock. It
// writes no backup record: a restored ledger starts without one. A failed
// restore returns an *Error.
func Restore(options RestoreOptions) (Restored, string, error) {
	root := options.WorkspaceRoot
	ledgerPath := config.LedgerDir(root)
	if err := noLedger(ledgerPath); err != nil {
		return Restored{}, ledgerPath, err
	}
	id, err := projectIdentity(root)
	if err != nil {
		return Restored{}, ledgerPath, err
	}
	getenv := lookup(options.Environment)
	target, err := configuredTarget(getenv)
	if err != nil {
		return Restored{}, ledgerPath, err
	}
	directory, err := os.MkdirTemp(filepath.Dir(ledgerPath), "fid-restore-")
	if err != nil {
		return Restored{}, ledgerPath, failure(KindInternal, err, "create the restore's temporary directory")
	}
	defer func() { _ = os.RemoveAll(directory) }()
	repository := clone{gitDir: filepath.Join(directory, "repo.git"), workingDir: directory, environment: gitEnvironment(options.Environment)}
	if _, err := runGit(directory, repository.environment, nil, "init", "--bare", "-q", "--template=", repository.gitDir); err != nil {
		return Restored{}, ledgerPath, failure(KindInternal, err, "create the restore's temporary repository")
	}
	commit, err := fetchFolder(repository, target, id)
	if err != nil {
		return Restored{}, ledgerPath, err
	}
	header, payload, err := readManifest(repository, commit, id)
	if err != nil {
		return Restored{}, ledgerPath, err
	}
	archive, err := folderArchive(repository, commit, id, header.Schema, payload, directory)
	if err != nil {
		return Restored{}, ledgerPath, err
	}
	imported, err := importArchive(options, archive, commit, id, ledgerPath)
	if err != nil {
		return Restored{}, ledgerPath, err
	}
	genesis, err := genesisHash(root)
	if err != nil {
		return Restored{}, ledgerPath, failure(KindInternal, err, "restored %s/, but reading its genesis Event failed", ledgerPath)
	}
	return Restored{ProjectID: id, Commit: commit, ArchiveDigest: imported.ArchiveDigest, EventCount: imported.EventCount,
		TerminalEventHash: imported.TerminalEventHash, GenesisEventHash: genesis}, ledgerPath, nil
}

// noLedger refuses when anything is at the ledger's path. Restore only fills
// a repository that has no ledger; it never replaces or merges one.
func noLedger(ledgerPath string) error {
	if _, err := os.Lstat(ledgerPath); err == nil {
		return ledgerExists(ledgerPath, nil)
	} else if !errors.Is(err, os.ErrNotExist) {
		return failure(KindInternal, err, "look for a ledger at %s", ledgerPath)
	}
	return nil
}

func ledgerExists(ledgerPath string, cause error) *Error {
	return &Error{Kind: KindLedgerExists, Path: ledgerPath, Err: cause,
		Message: fmt.Sprintf("a ledger already exists at %s/, so fid won't restore over it", ledgerPath)}
}

// fetchFolder fetches only the tip of the target's branch, so a restore's
// time doesn't grow with the backup history, and returns it once it is
// known to hold the project's folder.
func fetchFolder(repository clone, target, id string) (string, error) {
	remote, err := remoteTip(repository, target)
	if err != nil {
		return "", err
	}
	if remote == "" {
		return "", &Error{Kind: KindNoBackup, Detail: remoteRef + " does not exist",
			Message: fmt.Sprintf("the backup repository has no %s branch, so it holds no backup of %s", Branch, id)}
	}
	commit, err := fetchBranch(repository, target, "--depth=1")
	if err != nil {
		return "", err
	}
	listed, err := repository.text("ls-tree", "-z", commit, "--", id)
	if err != nil {
		return "", failure(KindInternal, err, "read the backup repository's tree")
	}
	if listed == "" {
		return "", &Error{Kind: KindNoBackup, Detail: fmt.Sprintf("%s at %s has no %s/", Branch, commit, id),
			Message: fmt.Sprintf("the backup repository holds no backup of %s: its %s has no %s/ folder", id, Branch, id)}
	}
	return commit, nil
}

// readManifest reads the folder's backup.json and checks what doesn't depend
// on the folder's format: that it names this project.
func readManifest(repository clone, commit, id string) (folderHeader, []byte, error) {
	path := id + "/" + ManifestName
	var payload bytes.Buffer
	found, err := folderBlob(repository, commit, path, &payload)
	if err != nil {
		return folderHeader{}, nil, err
	}
	if !found {
		return folderHeader{}, nil, damaged(id, folderCheck(CheckManifestInvalid, "%s is missing", path))
	}
	var header folderHeader
	if err := json.Unmarshal(payload.Bytes(), &header); err != nil {
		return folderHeader{}, nil, damaged(id, folderCheck(CheckManifestInvalid, "%s is not a manifest: %v", path, err))
	}
	if header.ProjectID != id {
		return folderHeader{}, nil, damaged(id, folderCheck(CheckManifestInvalid, "%s names project %q, not %s", path, header.ProjectID, id))
	}
	return header, payload.Bytes(), nil
}

// folderArchive turns the project's fetched folder into the archive to
// import, by the format the manifest's schema names, and returns the
// archive's path in scratch. It is the one place that knows how a folder
// holds the export.
func folderArchive(repository clone, commit, id, schema string, payload []byte, scratch string) (string, error) {
	if schema != ManifestSchema {
		return "", damaged(id, folderCheck(CheckManifestInvalid, "%s/%s has schema %q; fid restores %s", id, ManifestName, schema, ManifestSchema))
	}
	var written manifest
	if err := json.Unmarshal(payload, &written); err != nil {
		return "", damaged(id, folderCheck(CheckManifestInvalid, "%s/%s is not a %s manifest: %v", id, ManifestName, schema, err))
	}
	if err := checkExportIdentity(repository, commit, id); err != nil {
		return "", err
	}
	return zipFolder(repository, commit, id, written.ArchiveDigest, scratch)
}

// folderBlob writes the file at path in commit to output and reports whether
// it is there.
func folderBlob(repository clone, commit, path string, output io.Writer) (bool, error) {
	listed, err := repository.text("ls-tree", "-l", "-z", commit, "--", path)
	if err != nil {
		return false, failure(KindInternal, err, "read the backup repository's tree")
	}
	fields := strings.Fields(strings.SplitN(listed, "\t", 2)[0])
	if len(fields) != 4 || fields[1] != "blob" {
		return false, nil
	}
	id, _, _ := strings.Cut(path, "/")
	size, err := strconv.ParseInt(fields[3], 10, 64)
	limit := export.ArchiveResourceLimits().NonArtifactBytes
	if err != nil || size < 0 || size > limit {
		return false, folderResourceLimit(id, path)
	}
	bounded := &folderBlobWriter{output: output, remaining: limit}
	if err := repository.stream(bounded, "cat-file", "blob", fields[2]); err != nil {
		if bounded.exceeded {
			return false, folderResourceLimit(id, path)
		}
		return false, failure(KindInternal, err, "read %s from the backup repository", path)
	}
	return true, nil
}

// importArchive installs the archive as the ledger. A ledger that appeared
// while the restore fetched, including one another restore of the same
// backup installed, is refused like one that was there to begin with.
func importArchive(options RestoreOptions, archive, commit, id, ledgerPath string) (export.ImportResult, error) {
	imported, err := export.Import(export.ImportOptions{ArchivePath: archive, RepositoryRoot: options.WorkspaceRoot,
		MigrationsDir: options.MigrationsDir, IdempotencyKey: "fid-restore:" + commit})
	var conflict *export.ImportTargetConflictError
	var invalid *export.InvalidArchiveError
	switch {
	case errors.As(err, &conflict):
		return export.ImportResult{}, ledgerExists(ledgerPath, err)
	case errors.As(err, &invalid):
		return export.ImportResult{}, damaged(id, invalid.Checks...)
	case err != nil:
		return export.ImportResult{}, failure(KindInternal, err, "import the backup")
	case imported.Replayed:
		return export.ImportResult{}, ledgerExists(ledgerPath, nil)
	}
	return imported, nil
}

// damaged is the failure for a folder that can't be restored. It names the
// folder's archive, not the restore's temporary copy, which is deleted.
func damaged(id string, checks ...export.Check) *Error {
	archive := id + "/" + ArchiveName
	details := make([]string, 0, len(checks))
	for _, check := range checks {
		details = append(details, check.Code+": "+check.Detail)
	}
	return &Error{Kind: KindArchiveInvalid, Path: archive, Detail: strings.Join(details, "; "),
		Message: fmt.Sprintf("the backup in %s/ is damaged, so fid won't restore it", id),
		Err:     &export.InvalidArchiveError{ArchivePath: archive, Checks: checks}}
}

func folderCheck(code, format string, arguments ...any) export.Check {
	return export.Check{Code: code, Detail: fmt.Sprintf(format, arguments...)}
}
