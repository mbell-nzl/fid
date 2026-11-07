package main

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/mbell-nzl/fid/internal/cli"
	"github.com/mbell-nzl/fid/internal/config"
	"github.com/mbell-nzl/fid/internal/gittest"
	"golang.org/x/sys/unix"
)

type backupFailureForTest struct {
	Code      string         `json:"code"`
	Category  string         `json:"category"`
	Retryable bool           `json:"retryable"`
	Message   string         `json:"message"`
	Subject   any            `json:"subject"`
	Details   map[string]any `json:"details"`
}

func decodeBackupFailure(t *testing.T, result commandExecutionResult, exit int) backupFailureForTest {
	t.Helper()
	if cli.ExitCode(result.err) != exit || result.stdout != "" {
		t.Fatalf("backup = %#v (exit %d), want exit %d and no stdout", result, cli.ExitCode(result.err), exit)
	}
	var envelope struct {
		OK        bool                 `json:"ok"`
		Operation string               `json:"operation"`
		Error     backupFailureForTest `json:"error"`
	}
	if err := json.Unmarshal([]byte(result.stderr), &envelope); err != nil {
		t.Fatalf("decode failure %q: %v", result.stderr, err)
	}
	if envelope.OK || envelope.Operation != "project.backup" {
		t.Fatalf("failure envelope = %s", result.stderr)
	}
	return envelope.Error
}

func TestBackupFailuresAreRecordedAndTyped(t *testing.T) {
	t.Parallel()
	for _, test := range []struct {
		name    string
		setup   func(*backupFixture, string)
		exit    int
		code    string
		kind    string
		message string // in the recorded message, where kinds share a cause
		details func(*backupFixture, string) map[string]any
	}{
		{name: "no target", exit: 2, code: "invalid-argument", kind: "no-target",
			setup: func(*backupFixture, string) {},
			details: func(fixture *backupFixture, _ string) map[string]any {
				return map[string]any{"argument": fixture.userConfigPath(), "value": nil, "allowed": []any{"[backup] target: a git URL or local path"}}
			}},
		{name: "malformed user config", exit: 2, code: "invalid-argument", kind: "user-config",
			setup: func(fixture *backupFixture, _ string) {
				writeBackupTestFile(fixture.t, fixture.userConfigPath(), "[backup]\ntarget = broken\n")
			},
			details: func(fixture *backupFixture, _ string) map[string]any {
				return map[string]any{"argument": fixture.userConfigPath(), "value": nil, "allowed": []any{"valid project.backup configuration"}}
			}},
		{name: "unreachable target", exit: 4, code: "repository-unavailable", kind: "target-unreachable",
			setup: func(fixture *backupFixture, _ string) { fixture.setTarget(filepath.Join(fixture.base, "missing.git")) }},
		{name: "missing project identity", exit: 2, code: "invalid-argument", kind: "project-identity",
			setup: func(fixture *backupFixture, workspace string) {
				fixture.setTarget(fixture.target)
				if err := os.Remove(filepath.Join(workspace, config.ProjectIdentityFileName)); err != nil {
					fixture.t.Fatal(err)
				}
			},
			details: backupIdentityDetails},
		{name: "another project's identity", exit: 2, code: "invalid-argument", kind: "project-identity",
			setup: func(fixture *backupFixture, workspace string) {
				fixture.setTarget(fixture.target)
				writeBackupTestFile(fixture.t, filepath.Join(workspace, config.ProjectIdentityFileName), "[project]\nid = \"project-"+strings.Repeat("a", 32)+"\"\n")
			},
			details: backupIdentityDetails},
		{name: "export failure", exit: 7, code: "internal-error", kind: "export-failed", message: "export the ledger: ",
			setup: func(fixture *backupFixture, workspace string) {
				fixture.setTarget(fixture.target)
				journal, err := os.OpenFile(mustJournalPath(fixture.t, workspace), os.O_APPEND|os.O_WRONLY, 0)
				if err != nil {
					fixture.t.Fatal(err)
				}
				defer func() { _ = journal.Close() }()
				if _, err := journal.WriteString("not an event\n"); err != nil {
					fixture.t.Fatal(err)
				}
			},
			details: func(*backupFixture, string) map[string]any { return map[string]any{} }},
		// The archive is created, but validation rejects it: the identity in
		// .fid.toml and the ledger config agree, but differ from the journal's.
		{name: "export rejected by validation", exit: 7, code: "internal-error", kind: "export-failed", message: "validate the ledger's export: ",
			setup: func(fixture *backupFixture, workspace string) {
				fixture.setTarget(fixture.target)
				original, err := config.ReadProjectIdentity(workspace)
				if err != nil {
					fixture.t.Fatal(err)
				}
				replacement := "project-" + strings.Repeat("b", 32)
				for _, path := range []string{filepath.Join(workspace, config.ProjectIdentityFileName), filepath.Join(config.LedgerDir(workspace), config.ConfigFileName)} {
					payload, err := os.ReadFile(path)
					if err != nil || !strings.Contains(string(payload), original) {
						fixture.t.Fatalf("%s holds no %s: %v", path, original, err)
					}
					writeBackupTestFile(fixture.t, path, strings.ReplaceAll(string(payload), original, replacement))
				}
			},
			details: func(*backupFixture, string) map[string]any { return map[string]any{} }},
	} {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			fixture := newBackupFixture(t)
			workspace, _ := fixture.project()
			test.setup(fixture, workspace)

			text := fixture.backup(workspace)
			if cli.ExitCode(text.err) != test.exit || text.stdout != "" ||
				!strings.HasPrefix(text.stderr, "error: ") || !strings.Contains(text.stderr, "\nreason: "+test.code) {
				t.Fatalf("text backup = %#v, want exit %d and a %s refusal", text, test.exit, test.code)
			}
			failure := decodeBackupFailure(t, fixture.backup(workspace, "--json"), test.exit)
			if failure.Code != test.code || failure.Message == "" || failure.Subject != nil {
				t.Fatalf("failure = %+v, want %s", failure, test.code)
			}
			if test.details != nil {
				if got, want := mustJSON(t, failure.Details), mustJSON(t, test.details(fixture, workspace)); got != want {
					t.Fatalf("details = %s, want %s", got, want)
				}
			}
			if fixture.targetMain() != "" {
				t.Fatal("a failed backup pushed to the target")
			}
			records := backupRecords(t, workspace)
			if len(records) != 2 {
				t.Fatalf("records = %+v, want one per failed backup", records)
			}
			for _, record := range records {
				if record.Outcome != "failed" || record.Failure == nil || record.Failure.Kind != test.kind || record.Failure.Message == "" ||
					record.Commit != nil || record.ArchiveDigest != nil || record.EventCount != nil || record.TerminalEventHash != nil {
					t.Fatalf("record = %+v, want a %s failure", record, test.kind)
				}
				if (record.ProjectID == nil) != (test.kind == "project-identity") {
					t.Fatalf("record project_id = %v for %s", record.ProjectID, test.kind)
				}
			}
			if !strings.Contains(records[1].Failure.Message, failure.Message) || !strings.Contains(records[1].Failure.Message, test.message) {
				t.Fatalf("recorded message %q does not carry the reported %q and %q", records[1].Failure.Message, failure.Message, test.message)
			}
		})
	}
}

func backupIdentityDetails(fixture *backupFixture, workspace string) map[string]any {
	return map[string]any{"argument": filepath.Join(gittest.ResolvedPath(fixture.t, workspace), config.ProjectIdentityFileName),
		"value": nil, "allowed": []any{"[project] id matching the ledger's project ID"}}
}

func TestBackupUnreachableTargetNamesGitsReason(t *testing.T) {
	t.Parallel()
	fixture := newBackupFixture(t)
	workspace, _ := fixture.project()
	missing := filepath.Join(fixture.base, "missing.git")
	fixture.setTarget(missing)

	failure := decodeBackupFailure(t, fixture.backup(workspace, "--json"), 4)
	if failure.Code != "repository-unavailable" || failure.Category != "precondition" || !failure.Retryable ||
		failure.Details["reason"] != "not-a-repository" || !strings.Contains(failure.Details["detail"].(string), "does not appear to be a git repository") {
		t.Fatalf("failure = %+v, want a retryable unreachable target with git's reason", failure)
	}
}

func TestBackupWithoutALedgerRefusesAndRecordsNothing(t *testing.T) {
	t.Parallel()
	fixture := newBackupFixture(t)
	fixture.setTarget(fixture.target)
	workspace := gittest.NewWorkspace(t)

	failure := decodeBackupFailure(t, fixture.backup(workspace, "--json"), 3)
	if failure.Code != "project-not-found" {
		t.Fatalf("failure = %+v, want project-not-found", failure)
	}
	if _, err := os.Lstat(config.LedgerDir(workspace)); !os.IsNotExist(err) {
		t.Fatalf("refused backup created a ledger: %v", err)
	}
}

func TestBackupRejectedPushIsRecordedAndTheNextRunRecovers(t *testing.T) {
	t.Parallel()
	fixture := newBackupFixture(t)
	workspace, id := fixture.project()
	fixture.setTarget(fixture.target)
	hook := filepath.Join(fixture.target, "hooks", "pre-receive")
	writeBackupTestFile(t, hook, "#!/bin/sh\necho 'denied by test' >&2\nexit 1\n")
	if err := os.Chmod(hook, 0o755); err != nil {
		t.Fatal(err)
	}

	failure := decodeBackupFailure(t, fixture.backup(workspace, "--json"), 4)
	if failure.Code != "repository-unavailable" || failure.Details["reason"] != "push-rejected" ||
		!strings.Contains(failure.Details["detail"].(string), "pre-receive hook declined") {
		t.Fatalf("failure = %+v, want a rejected push", failure)
	}
	if fixture.targetMain() != "" {
		t.Fatal("rejected push reached the target")
	}

	if err := os.Remove(hook); err != nil {
		t.Fatal(err)
	}
	record, _ := decodeBackupSuccess(t, fixture.backup(workspace, "--json"))
	if record.Outcome != "pushed" || *record.Commit != fixture.targetMain() {
		t.Fatalf("retry = %+v, want a push", record)
	}
	if count := fixture.targetGit("rev-list", "--count", "main"); count != "1" {
		t.Fatalf("target commits = %s, want only the recovered backup", count)
	}
	records := backupRecords(t, workspace)
	if got := backupOutcomes(records); strings.Join(got, ",") != "failed,pushed" || records[0].Failure.Kind != "push-rejected" || *records[0].ProjectID != id {
		t.Fatalf("records = %+v, want the rejection then the recovery", records)
	}
}

func TestBackupWaitsForTheLockAndRecordsLockUnavailable(t *testing.T) {
	t.Parallel()
	fixture := newBackupFixture(t)
	workspace, _ := fixture.project()
	fixture.setTarget(fixture.target)
	decodeBackupSuccess(t, fixture.backup(workspace, "--json"))
	locks, err := filepath.Glob(filepath.Join(fixture.cacheHome, "fid", "backup", "*", "lock"))
	if err != nil || len(locks) != 1 {
		t.Fatalf("backup locks = %v, %v; want one beside fid's clone", locks, err)
	}
	holder, err := os.Open(locks[0])
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = holder.Close() }()
	if err := unix.Flock(int(holder.Fd()), unix.LOCK_EX|unix.LOCK_NB); err != nil {
		t.Fatal(err)
	}

	failure := decodeBackupFailure(t, fixture.backupWaiting(workspace, 200*time.Millisecond, "--json"), 4)
	if failure.Code != "lock-unavailable" || failure.Details["lock_path"] != locks[0] {
		t.Fatalf("failure = %+v, want lock-unavailable at %s", failure, locks[0])
	}
	if err := unix.Flock(int(holder.Fd()), unix.LOCK_UN); err != nil {
		t.Fatal(err)
	}
	decodeBackupSuccess(t, fixture.backup(workspace, "--json"))
	records := backupRecords(t, workspace)
	if got := backupOutcomes(records); strings.Join(got, ",") != "pushed,failed,unchanged" || records[1].Failure.Kind != "lock-unavailable" {
		t.Fatalf("records = %+v, want the lock wait recorded as a failure", records)
	}
}

func TestConcurrentBackupsShareFidsCloneSafely(t *testing.T) {
	t.Parallel()
	fixture := newBackupFixture(t)
	first, firstID := fixture.project()
	second, secondID := fixture.project()
	fixture.setTarget(fixture.target)

	start := make(chan struct{})
	results := make(chan commandExecutionResult, 3)
	for _, workspace := range []string{first, first, second} {
		go func() {
			<-start
			results <- fixture.backup(workspace, "--json")
		}()
	}
	close(start)
	for range 3 {
		if result := <-results; result.err != nil {
			t.Fatalf("concurrent backup = %#v", result)
		}
	}
	firstOutcomes := backupOutcomes(backupRecords(t, first))
	sort.Strings(firstOutcomes)
	if strings.Join(firstOutcomes, ",") != "pushed,unchanged" {
		t.Fatalf("first project's outcomes = %v, want one push and one unchanged", firstOutcomes)
	}
	if got := backupOutcomes(backupRecords(t, second)); strings.Join(got, ",") != "pushed" {
		t.Fatalf("second project's outcomes = %v", got)
	}
	if count := fixture.targetGit("rev-list", "--count", "main"); count != "2" {
		t.Fatalf("target commits = %s, want one per changed project", count)
	}
	for _, id := range []string{firstID, secondID} {
		fixture.targetGit("cat-file", "-e", "main:"+id+"/export.zip")
	}
	if record, _ := decodeBackupSuccess(t, fixture.backup(first, "--json")); record.Outcome != "unchanged" {
		t.Fatalf("backup after the race = %+v, want unchanged", record)
	}
}

// A record that can't be written is itself a failure fid reports, after a
// push as after a failed backup, so no outcome goes unrecorded silently.
func TestBackupReportsARecordItCannotWrite(t *testing.T) {
	t.Parallel()
	fixture := newBackupFixture(t)
	workspace, _ := fixture.project()
	if err := os.Mkdir(filepath.Join(config.LedgerDir(workspace), "backups.jsonl"), 0o700); err != nil {
		t.Fatal(err)
	}

	failure := decodeBackupFailure(t, fixture.backup(workspace, "--json"), 2)
	if failure.Code != "invalid-argument" || !strings.Contains(failure.Message, "no backup target is set") ||
		!strings.Contains(failure.Message, "recording the failure failed") {
		t.Fatalf("failure = %+v, want the missing target and the unwritable record", failure)
	}
	fixture.setTarget(fixture.target)
	failure = decodeBackupFailure(t, fixture.backup(workspace, "--json"), 7)
	if failure.Code != "internal-error" || !strings.Contains(failure.Message, "backup pushed, but recording it failed") || fixture.targetMain() == "" {
		t.Fatalf("failure = %+v, want a push reported with its unwritable record", failure)
	}
}

// A write cut short, by a full disk or a killed process, leaves the record
// without a final newline. The next backup must not append onto it: its own
// line must parse, and every complete line before it must survive.
func TestBackupRepairsATornRecordBeforeAppending(t *testing.T) {
	t.Parallel()
	for _, test := range []struct {
		name string
		tear func(t *testing.T, path string, payload []byte) []byte // returns the complete lines to keep
		want string
	}{
		{name: "fragment of a record", want: "pushed,unchanged", tear: func(t *testing.T, path string, payload []byte) []byte {
			appendBackupTestBytes(t, path, payload[:48])
			return payload
		}},
		{name: "whole record without its newline", want: "pushed,pushed,unchanged", tear: func(t *testing.T, path string, payload []byte) []byte {
			appendBackupTestBytes(t, path, bytes.TrimSuffix(payload, []byte("\n")))
			return append(append([]byte{}, payload...), payload...)
		}},
	} {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			fixture := newBackupFixture(t)
			workspace, _ := fixture.project()
			fixture.setTarget(fixture.target)
			decodeBackupSuccess(t, fixture.backup(workspace, "--json"))
			path := filepath.Join(config.LedgerDir(workspace), "backups.jsonl")
			payload, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			kept := test.tear(t, path, payload)

			record, _ := decodeBackupSuccess(t, fixture.backup(workspace, "--json"))
			after, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			if !bytes.HasPrefix(after, kept) {
				t.Fatalf("record = %q, want the complete lines %q kept", after, kept)
			}
			records := backupRecords(t, workspace)
			if got := backupOutcomes(records); strings.Join(got, ",") != test.want || mustJSON(t, records[len(records)-1]) != mustJSON(t, record) {
				t.Fatalf("records = %v, want %s ending with the new attempt", got, test.want)
			}
		})
	}
}

func appendBackupTestBytes(t *testing.T, path string, payload []byte) {
	t.Helper()
	file, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = file.Close() }()
	if _, err := file.Write(payload); err != nil {
		t.Fatal(err)
	}
}

func writeBackupTestFile(t *testing.T, path, contents string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(contents), 0o600); err != nil {
		t.Fatal(err)
	}
}

func mustJSON(t *testing.T, value any) string {
	t.Helper()
	payload, err := json.Marshal(value)
	if err != nil {
		t.Fatal(err)
	}
	return string(payload)
}
