package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/mbell-nzl/fid/internal/config"
	"github.com/mbell-nzl/fid/internal/gittest"
)

// Clones of one project each make their own ledger with the project's ID, as
// does fid init on a new machine. Their backups name the same folder, so a
// backup must never replace a folder another ledger wrote: it may be that
// ledger's only copy.
func TestBackupRefusesAFolderAnotherLedgerWrote(t *testing.T) {
	t.Parallel()
	fixture := newBackupFixture(t)
	first, id := fixture.project()
	fixture.setTarget(fixture.target)
	written, _ := decodeBackupSuccess(t, fixture.backup(first, "--json"))
	second := fixture.projectClone(id)
	firstGenesis, secondGenesis := ledgerGenesisHash(t, first), ledgerGenesisHash(t, second)
	if firstGenesis == secondGenesis {
		t.Fatalf("both ledgers have genesis %s, want one each", firstGenesis)
	}

	text := fixture.backup(second)
	if backupExitCode(text.err) != 4 || text.stdout != "" || !strings.HasPrefix(text.stderr, "error: ") ||
		!strings.Contains(text.stderr, "\nreason: repository-unavailable") || !strings.Contains(text.stderr, "remove "+id+"/") {
		t.Fatalf("text backup = %#v, want exit 4 and a refusal saying how to clear it", text)
	}
	failure := decodeBackupFailure(t, fixture.backup(second, "--json"), 4)
	detail, _ := failure.Details["detail"].(string)
	if failure.Code != "repository-unavailable" || failure.Category != "precondition" || failure.Subject != nil ||
		failure.Details["reason"] != "different-ledger" || !strings.Contains(failure.Message, id+"/") ||
		!strings.Contains(detail, firstGenesis) || !strings.Contains(detail, secondGenesis) {
		t.Fatalf("failure = %+v, want different-ledger naming both ledgers' genesis Events", failure)
	}
	if main := fixture.targetMain(); main != *written.Commit {
		t.Fatalf("target main = %s, want the first ledger's %s untouched", main, *written.Commit)
	}
	assertDifferentLedgerRecords(t, second, id, 2)

	// The folder is still the first ledger's to update, and still refused to
	// the second.
	runFIDCommand(t, first, "issue", "create", "--kind", "task", "--title", "Backed up again", "--actor", "agent:developer")
	if record, _ := decodeBackupSuccess(t, fixture.backup(first, "--json")); record.Outcome != "pushed" {
		t.Fatalf("first ledger's backup = %+v, want a push", record)
	}
	decodeBackupFailure(t, fixture.backup(second, "--json"), 4)
	if manifest := targetFolderManifest(fixture, id); manifest.GenesisEventHash != firstGenesis {
		t.Fatalf("folder manifest genesis = %s, want the first ledger's %s", manifest.GenesisEventHash, firstGenesis)
	}
}

// A folder that doesn't say which ledger wrote it can't be shown to be this
// ledger's, so fid refuses it too rather than guess.
func TestBackupRefusesAFolderThatNamesNoLedger(t *testing.T) {
	t.Parallel()
	for _, test := range []struct {
		name  string
		files func(id string) map[string]string
	}{
		{"no manifest", func(id string) map[string]string { return map[string]string{id + "/notes.txt": "mine\n"} }},
		{"manifest without a genesis", func(id string) map[string]string {
			return map[string]string{id + "/backup.json": `{"schema":"fid.backup.v1","project_id":"` + id + `"}` + "\n"}
		}},
		{"manifest that is not JSON", func(id string) map[string]string { return map[string]string{id + "/backup.json": "not json\n"} }},
		{"a file, not a folder", func(id string) map[string]string { return map[string]string{id: "mine\n"} }},
	} {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			fixture := newBackupFixture(t)
			workspace, id := fixture.project()
			fixture.setTarget(fixture.target)
			planted := fixture.commitToTarget(test.files(id))

			failure := decodeBackupFailure(t, fixture.backup(workspace, "--json"), 4)
			if failure.Code != "repository-unavailable" || failure.Details["reason"] != "different-ledger" {
				t.Fatalf("failure = %+v, want different-ledger", failure)
			}
			if main := fixture.targetMain(); main != planted {
				t.Fatalf("target main = %s, want the planted %s untouched", main, planted)
			}
			assertDifferentLedgerRecords(t, workspace, id, 1)
		})
	}
}

// Restoring imports the folder's archive, which keeps the journal and so
// the genesis Event: the restored ledger is the same ledger, and its backups
// update the folder.
func TestBackupFromARestoredLedgerUpdatesItsFolder(t *testing.T) {
	t.Parallel()
	fixture := newBackupFixture(t)
	original, id := fixture.project()
	fixture.setTarget(fixture.target)
	decodeBackupSuccess(t, fixture.backup(original, "--json"))
	archive, _ := exportForTest(t, original)
	archivePath := filepath.Join(t.TempDir(), "export.zip")
	if err := os.WriteFile(archivePath, archive, 0o600); err != nil {
		t.Fatal(err)
	}
	restored := gittest.NewWorkspace(t)
	if imported := executeCommandForTest(restored, "export", "import", archivePath, "--repository", restored, "--idempotency-key", "restore", "--json"); imported.err != nil {
		t.Fatalf("export import = %#v", imported)
	}
	writeBackupTestFile(t, filepath.Join(restored, config.ProjectIdentityFileName), "[project]\nid = \""+id+"\"\n")
	if got, want := ledgerGenesisHash(t, restored), ledgerGenesisHash(t, original); got != want {
		t.Fatalf("restored genesis = %s, want the original's %s", got, want)
	}

	decodeBackupSuccess(t, fixture.backup(restored, "--json"))
	runFIDCommand(t, restored, "issue", "create", "--kind", "task", "--title", "After the restore", "--actor", "agent:developer")
	record, _ := decodeBackupSuccess(t, fixture.backup(restored, "--json"))
	if record.Outcome != "pushed" || *record.Commit != fixture.targetMain() || targetFolderManifest(fixture, id).EventCount != *record.EventCount {
		t.Fatalf("restored ledger's backup = %+v, want it to update %s/", record, id)
	}
}

// projectClone is another clone of the project with ID id: its .fid.toml is
// committed with that ID, and fid init makes the clone its own ledger.
func (fixture *backupFixture) projectClone(id string) string {
	fixture.t.Helper()
	workspace := gittest.NewWorkspace(fixture.t)
	writeBackupTestFile(fixture.t, filepath.Join(workspace, config.ProjectIdentityFileName), "[project]\nid = \""+id+"\"\n")
	runInitCommand(fixture.t, workspace, "init", "--json")
	if got, err := config.ReadProjectIdentity(workspace); err != nil || got != id {
		fixture.t.Fatalf("clone's project ID = %q, %v; want %s", got, err, id)
	}
	return workspace
}

// commitToTarget commits files to the target's main as a person would, and
// returns the new tip.
func (fixture *backupFixture) commitToTarget(files map[string]string) string {
	fixture.t.Helper()
	human := filepath.Join(fixture.t.TempDir(), "human")
	fixture.git("clone", "-q", fixture.target, human)
	names := make([]string, 0, len(files))
	for name, contents := range files {
		writeBackupTestFile(fixture.t, filepath.Join(human, name), contents)
		names = append(names, name)
	}
	sort.Strings(names)
	fixture.git(append([]string{"-C", human, "add", "--"}, names...)...)
	fixture.git("-C", human, "-c", "user.name=human", "-c", "user.email=human@example.test", "commit", "-q", "-m", "by hand")
	fixture.git("-C", human, "push", "-q", "origin", "HEAD:refs/heads/main")
	return fixture.targetMain()
}

// ledgerGenesisHash is the hash fid reports for the ledger's genesis Event,
// the first line of its first journal file after fid init.
func ledgerGenesisHash(t *testing.T, workspace string) string {
	t.Helper()
	journalDir := filepath.Join(config.LedgerDir(workspace), config.JournalDirName)
	files, err := filepath.Glob(filepath.Join(journalDir, "*.jsonl"))
	if err != nil || len(files) == 0 {
		t.Fatalf("journal files = %v, %v", files, err)
	}
	sort.Strings(files)
	payload, err := os.ReadFile(files[0])
	if err != nil {
		t.Fatal(err)
	}
	var genesis struct {
		EventID   string `json:"event_id"`
		EventType string `json:"event_type"`
	}
	first, _, _ := strings.Cut(string(payload), "\n")
	if err := json.Unmarshal([]byte(first), &genesis); err != nil || genesis.EventType != "project.initialized" {
		t.Fatalf("first journal line = %q, %v; want the genesis Event", first, err)
	}
	inspected := executeCommandForTest(workspace, "journal", "inspect", genesis.EventID)
	var event struct {
		Hash string `json:"hash"`
	}
	if inspected.err != nil || json.Unmarshal([]byte(inspected.stdout), &event) != nil || !strings.HasPrefix(event.Hash, "sha256:") {
		t.Fatalf("journal inspect %s = %#v", genesis.EventID, inspected)
	}
	return event.Hash
}

func targetFolderManifest(fixture *backupFixture, id string) backupFolderManifestForTest {
	fixture.t.Helper()
	var manifest backupFolderManifestForTest
	decoder := json.NewDecoder(strings.NewReader(fixture.targetGit("cat-file", "blob", "main:"+id+"/backup.json")))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&manifest); err != nil {
		fixture.t.Fatal(err)
	}
	return manifest
}

func assertDifferentLedgerRecords(t *testing.T, workspace, id string, count int) {
	t.Helper()
	records := backupRecords(t, workspace)
	if len(records) != count {
		t.Fatalf("records = %+v, want %d", records, count)
	}
	for _, record := range records {
		if record.Outcome != "failed" || record.Failure == nil || record.Failure.Kind != "different-ledger" ||
			record.ProjectID == nil || *record.ProjectID != id || record.Commit != nil || record.ArchiveDigest != nil {
			t.Fatalf("record = %+v, want a different-ledger failure for %s", record, id)
		}
	}
}

// fid checks the folder against the target's tip as it read it, so another
// ledger's backup can land between that check and the push. fid never forces
// the push, so git must then reject it rather than replace that backup. A
// remote helper stands in for the target's transport and publishes the
// other ledger's backup, made earlier on another machine, as fid's push
// connects.
func TestBackupNeverReplacesAFolderPublishedAfterItsCheck(t *testing.T) {
	t.Parallel()
	fixture := newBackupFixture(t)
	first, id := fixture.project()
	second := fixture.projectClone(id)
	other := newBackupFixture(t)
	other.setTarget(other.target)
	published, _ := decodeBackupSuccess(t, other.backup(second, "--json"))

	bin, marker := filepath.Join(fixture.base, "bin"), filepath.Join(fixture.base, "published")
	writeBackupTestFile(t, filepath.Join(bin, "git-remote-fidrace"), strings.Join([]string{
		"#!/bin/sh",
		"while read -r command; do",
		"  case \"$command\" in",
		"  capabilities) printf 'connect\\n\\n' ;;",
		"  'connect git-receive-pack')",
		"    if [ ! -e '" + marker + "' ]; then",
		"      touch '" + marker + "'",
		"      git --git-dir='" + other.target + "' push -q \"$2\" main:refs/heads/main >&2 || exit 1",
		"    fi",
		"    printf '\\n'; exec git receive-pack \"$2\" ;;",
		"  'connect git-upload-pack') printf '\\n'; exec git upload-pack \"$2\" ;;",
		"  *) exit 1 ;;",
		"  esac",
		"done",
		"",
	}, "\n"))
	if err := os.Chmod(filepath.Join(bin, "git-remote-fidrace"), 0o755); err != nil {
		t.Fatal(err)
	}
	fixture.setTarget("fidrace::" + fixture.target)
	environ := append(append([]string{}, fixture.env...), "PATH="+bin+string(os.PathListSeparator)+os.Getenv("PATH"))

	failure := decodeBackupFailure(t, runBackupWith(environ, first, time.Minute, "--json"), 4)
	if failure.Code != "repository-unavailable" || failure.Details["reason"] != "push-rejected" {
		t.Fatalf("failure = %+v, want the push rejected", failure)
	}
	if _, err := os.Lstat(marker); err != nil {
		t.Fatalf("the other ledger's backup was never published: %v", err)
	}
	if main := fixture.targetMain(); main != *published.Commit {
		t.Fatalf("target main = %s, want the other ledger's %s kept", main, *published.Commit)
	}
	if got := backupOutcomes(backupRecords(t, first)); strings.Join(got, ",") != "failed" || backupRecords(t, first)[0].Failure.Kind != "push-rejected" {
		t.Fatalf("records = %v, want the rejected push", got)
	}

	// The next backup sees the other ledger's folder, and refuses it.
	failure = decodeBackupFailure(t, runBackupWith(environ, first, time.Minute, "--json"), 4)
	if failure.Details["reason"] != "different-ledger" || fixture.targetMain() != *published.Commit ||
		targetFolderManifest(fixture, id).GenesisEventHash != ledgerGenesisHash(t, second) {
		t.Fatalf("retry = %+v, want different-ledger with the other ledger's folder kept", failure)
	}
}
