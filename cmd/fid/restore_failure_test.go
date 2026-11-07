package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"strings"
	"testing"

	"github.com/mbell-nzl/fid/internal/cli"
	"github.com/mbell-nzl/fid/internal/config"
	"github.com/mbell-nzl/fid/internal/gittest"
)

func decodeRestoreFailure(t *testing.T, result commandExecutionResult, exit int) backupFailureForTest {
	t.Helper()
	if cli.ExitCode(result.err) != exit || result.stdout != "" {
		t.Fatalf("restore = %#v (exit %d), want exit %d and no stdout", result, cli.ExitCode(result.err), exit)
	}
	var envelope struct {
		OK        bool                 `json:"ok"`
		Operation string               `json:"operation"`
		Error     backupFailureForTest `json:"error"`
	}
	if err := json.Unmarshal([]byte(result.stderr), &envelope); err != nil {
		t.Fatalf("decode failure %q: %v", result.stderr, err)
	}
	if envelope.OK || envelope.Operation != "project.restore" {
		t.Fatalf("failure envelope = %s", result.stderr)
	}
	return envelope.Error
}

// restoreCase is a restore that must fail: setup returns the workspace to
// restore into.
type restoreCase struct {
	name    string
	setup   func(fixture *backupFixture) string
	exit    int
	code    string
	reason  string // repository-unavailable's, where details vary
	details func(fixture *backupFixture, workspace string) map[string]any
}

func TestRestoreFailuresAreTypedAndChangeNothing(t *testing.T) {
	t.Parallel()
	identityDetails := func(_ *backupFixture, workspace string) map[string]any {
		return map[string]any{"argument": filepath.Join(gittest.ResolvedPath(t, workspace), config.ProjectIdentityFileName),
			"value": nil, "allowed": []any{"[project] id naming the backup folder"}}
	}
	someID := "project-" + strings.Repeat("c", 32)
	for _, test := range []restoreCase{
		// Checked first: a broken user config and a missing identity don't
		// hide the ledger that is already there.
		{name: "a ledger already exists", exit: 4, code: "project-already-initialized",
			setup: func(fixture *backupFixture) string {
				workspace, _ := fixture.project()
				writeBackupTestFile(fixture.t, fixture.userConfigPath(), "[backup]\ntarget = broken\n")
				if err := os.Remove(filepath.Join(workspace, config.ProjectIdentityFileName)); err != nil {
					fixture.t.Fatal(err)
				}
				return workspace
			},
			details: func(_ *backupFixture, workspace string) map[string]any {
				return map[string]any{"path": gittest.ResolvedPath(t, workspace)}
			}},
		{name: "an empty ledger directory exists", exit: 4, code: "project-already-initialized",
			setup: func(fixture *backupFixture) string {
				fixture.setTarget(fixture.target)
				workspace := fixture.restoreWorkspace(someID)
				if err := os.Mkdir(config.LedgerDir(workspace), 0o700); err != nil {
					fixture.t.Fatal(err)
				}
				return workspace
			},
			details: func(_ *backupFixture, workspace string) map[string]any {
				return map[string]any{"path": gittest.ResolvedPath(t, workspace)}
			}},
		{name: "missing project identity", exit: 2, code: "invalid-argument", details: identityDetails,
			setup: func(fixture *backupFixture) string {
				fixture.setTarget(fixture.target)
				return gittest.NewWorkspace(fixture.t)
			}},
		{name: "invalid project identity", exit: 2, code: "invalid-argument", details: identityDetails,
			setup: func(fixture *backupFixture) string {
				fixture.setTarget(fixture.target)
				workspace := gittest.NewWorkspace(fixture.t)
				writeBackupTestFile(fixture.t, filepath.Join(workspace, config.ProjectIdentityFileName), "[project]\nid = \"FID\"\n")
				return workspace
			}},
		{name: "no target", exit: 2, code: "invalid-argument",
			setup: func(fixture *backupFixture) string { return fixture.restoreWorkspace(someID) },
			details: func(fixture *backupFixture, _ string) map[string]any {
				return map[string]any{"argument": fixture.userConfigPath(), "value": nil, "allowed": []any{"[backup] target: a git URL or local path"}}
			}},
		{name: "malformed user config", exit: 2, code: "invalid-argument",
			setup: func(fixture *backupFixture) string {
				writeBackupTestFile(fixture.t, fixture.userConfigPath(), "[backup]\ntarget = broken\n")
				return fixture.restoreWorkspace(someID)
			},
			details: func(fixture *backupFixture, _ string) map[string]any {
				return map[string]any{"argument": fixture.userConfigPath(), "value": nil, "allowed": []any{"valid project.restore configuration"}}
			}},
		{name: "unreachable target", exit: 4, code: "repository-unavailable", reason: "not-a-repository",
			setup: func(fixture *backupFixture) string {
				fixture.setTarget(filepath.Join(fixture.base, "missing.git"))
				return fixture.restoreWorkspace(someID)
			}},
		{name: "target that is not a repository", exit: 4, code: "repository-unavailable", reason: "not-a-repository",
			setup: func(fixture *backupFixture) string {
				directory := filepath.Join(fixture.base, "plain")
				if err := os.Mkdir(directory, 0o700); err != nil {
					fixture.t.Fatal(err)
				}
				fixture.setTarget(directory)
				return fixture.restoreWorkspace(someID)
			}},
		{name: "target without main", exit: 4, code: "repository-unavailable",
			details: func(*backupFixture, string) map[string]any {
				return map[string]any{"reason": "no-backup", "detail": "refs/heads/main does not exist"}
			},
			setup: func(fixture *backupFixture) string {
				fixture.setTarget(fixture.target)
				return fixture.restoreWorkspace(someID)
			}},
		{name: "main without the project's folder", exit: 4, code: "repository-unavailable",
			details: func(fixture *backupFixture, _ string) map[string]any {
				return map[string]any{"reason": "no-backup", "detail": "main at " + fixture.targetMain() + " has no " + someID + "/"}
			},
			setup: func(fixture *backupFixture) string {
				other, _ := fixture.project()
				fixture.setTarget(fixture.target)
				decodeBackupSuccess(fixture.t, fixture.backup(other, "--json"))
				return fixture.restoreWorkspace(someID)
			}},
	} {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			fixture := newBackupFixture(t)
			assertRestoreRefused(t, fixture, test)
		})
	}
}

// assertRestoreRefused runs a failing restore in text and JSON, and checks
// it changed neither the workspace's repository nor the target, and wrote no
// ledger where there was none.
func assertRestoreRefused(t *testing.T, fixture *backupFixture, test restoreCase) backupFailureForTest {
	t.Helper()
	workspace := test.setup(fixture)
	_, ledgerErr := os.Lstat(config.LedgerDir(workspace))
	repositoryBefore, targetBefore := snapshotWorkspace(t, filepath.Dir(config.LedgerDir(workspace))), fixture.targetMain()

	text := fixture.restore(workspace)
	if cli.ExitCode(text.err) != test.exit || text.stdout != "" ||
		!strings.HasPrefix(text.stderr, "error: ") || !strings.Contains(text.stderr, "\nreason: "+test.code) {
		t.Fatalf("text restore = %#v, want exit %d and a %s refusal", text, test.exit, test.code)
	}
	failure := decodeRestoreFailure(t, fixture.restore(workspace, "--json"), test.exit)
	if failure.Code != test.code || failure.Message == "" || test.reason != "" && failure.Details["reason"] != test.reason {
		t.Fatalf("failure = %+v, want %s %s", failure, test.code, test.reason)
	}
	if test.details != nil {
		if got, want := mustJSON(t, failure.Details), mustJSON(t, test.details(fixture, workspace)); got != want {
			t.Fatalf("details = %s, want %s", got, want)
		}
	}
	if !reflect.DeepEqual(repositoryBefore, snapshotWorkspace(t, filepath.Dir(config.LedgerDir(workspace)))) || fixture.targetMain() != targetBefore {
		t.Fatal("a refused restore changed the repository or the target")
	}
	if _, err := os.Lstat(config.LedgerDir(workspace)); os.IsNotExist(ledgerErr) != os.IsNotExist(err) {
		t.Fatalf("a refused restore changed whether a ledger exists: %v, then %v", ledgerErr, err)
	}
	return failure
}

func TestRestoreRefusesADamagedFolder(t *testing.T) {
	t.Parallel()
	source := newBackupFixture(t)
	original, id := source.project()
	source.setTarget(source.target)
	decodeBackupSuccess(t, source.backup(original, "--json"))
	archive, err := source.gitBytes("--git-dir="+source.target, "cat-file", "blob", "main:"+id+"/export.zip")
	if err != nil {
		t.Fatal(err)
	}
	manifest := source.targetGit("cat-file", "blob", "main:"+id+"/backup.json") + "\n"
	notAnArchive := "not a zip archive\n"
	digest := sha256.Sum256([]byte(notAnArchive))
	manifestFor := func(replace ...string) string { return strings.NewReplacer(replace...).Replace(manifest) }
	var written struct {
		ArchiveDigest string `json:"archive_digest"`
	}
	if err := json.Unmarshal([]byte(manifest), &written); err != nil {
		t.Fatal(err)
	}
	folder := func(files ...string) map[string]string {
		planted := map[string]string{}
		for index := 0; index < len(files); index += 2 {
			planted[id+"/"+files[index]] = files[index+1]
		}
		return planted
	}
	for _, test := range []struct {
		name  string
		files map[string]string
		code  string
	}{
		{"backup.json missing", folder("export.zip", string(archive)), "backup-manifest-invalid"},
		{"backup.json not JSON", folder("export.zip", string(archive), "backup.json", "not json\n"), "backup-manifest-invalid"},
		{"backup.json with an unknown schema", folder("export.zip", string(archive), "backup.json", manifestFor(`"fid.backup.v1"`, `"fid.backup.v9"`)), "backup-manifest-invalid"},
		{"backup.json naming another project", folder("export.zip", string(archive), "backup.json", manifestFor(`"project_id":"`+id, `"project_id":"project-`+strings.Repeat("d", 32))), "backup-manifest-invalid"},
		{"export.zip missing", folder("backup.json", manifest), "archive-missing"},
		{"export.zip that doesn't match its digest", folder("export.zip", string(archive)+"x", "backup.json", manifest), "archive-digest-mismatch"},
		{"a file where the folder should be", map[string]string{id: "not a folder\n"}, "backup-manifest-invalid"},
		// The digest matches, so only import's own validation finds it.
		{"export.zip that is not an export", folder("export.zip", notAnArchive, "backup.json",
			manifestFor(written.ArchiveDigest, "sha256:"+hex.EncodeToString(digest[:]))), ""},
	} {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			fixture := newBackupFixture(t)
			failure := assertRestoreRefused(t, fixture, restoreCase{exit: 6, code: "archive-invalid", setup: func(fixture *backupFixture) string {
				fixture.setTarget(fixture.target)
				fixture.commitToTarget(test.files)
				return fixture.restoreWorkspace(id)
			}})
			subject, _ := json.Marshal(failure.Subject)
			checks, _ := failure.Details["checks"].([]any)
			if failure.Category != "integrity" || failure.Retryable || string(subject) != `{"id":"`+id+`/export.zip","type":"export"}` ||
				failure.Details["archive_path"] != id+"/export.zip" || len(checks) == 0 || !strings.Contains(failure.Message, id+"/") {
				t.Fatalf("failure = %+v, want archive-invalid for %s/export.zip", failure, id)
			}
			codes := []string{}
			for _, check := range checks {
				codes = append(codes, check.(map[string]any)["code"].(string))
			}
			ours := []string{"backup-manifest-invalid", "archive-missing", "archive-digest-mismatch"}
			if test.code != "" && (len(codes) != 1 || codes[0] != test.code || checks[0].(map[string]any)["entry"] != nil) ||
				test.code == "" && slices.ContainsFunc(codes, func(code string) bool { return slices.Contains(ours, code) }) {
				t.Fatalf("checks = %v, want %q", checks, test.code)
			}
		})
	}
}
