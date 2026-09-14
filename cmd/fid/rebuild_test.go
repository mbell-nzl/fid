package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	fiddb "github.com/mbell-nzl/fid/internal/db"
	"github.com/mbell-nzl/fid/internal/journal"
)

func TestRebuildJSONSucceedsWithoutActorAuthority(t *testing.T) {
	workspace := t.TempDir()
	runInitCommand(t, workspace, "init")
	journalPath := mustJournalPath(t, workspace)
	beforeJournal, err := os.ReadFile(journalPath)
	if err != nil {
		t.Fatalf("ReadFile(%q) error = %v", journalPath, err)
	}
	if err := os.Remove(filepath.Join(workspace, ".fid", "fid.sqlite")); err != nil {
		t.Fatalf("Remove(projection) error = %v", err)
	}

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	cmd := newRootCommandAt(workspace)
	cmd.SetOut(stdout)
	cmd.SetErr(stderr)
	cmd.SetArgs([]string{"--json", "rebuild"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v, want role-free rebuild success", err)
	}
	if stderr.Len() != 0 {
		t.Fatalf("stderr = %q, want empty", stderr.String())
	}
	var success struct {
		OK        bool   `json:"ok"`
		Operation string `json:"operation"`
		Event     *any   `json:"event"`
	}
	if err := json.Unmarshal(stdout.Bytes(), &success); err != nil {
		t.Fatalf("Unmarshal(stdout) error = %v", err)
	}
	if !success.OK || success.Operation != "project.rebuild" || success.Event != nil {
		t.Fatalf("success = %+v, want role-free project.rebuild result", success)
	}
	afterJournal, err := os.ReadFile(journalPath)
	if err != nil {
		t.Fatalf("ReadFile(%q) after rebuild error = %v", journalPath, err)
	}
	if !bytes.Equal(afterJournal, beforeJournal) {
		t.Fatal("role-free rebuild appended journal data")
	}
}

func TestRebuildRejectsRemovedByFlagWithoutMutation(t *testing.T) {
	workspace := t.TempDir()
	runInitCommand(t, workspace, "init")

	journalPath := mustJournalPath(t, workspace)
	databasePath := filepath.Join(workspace, ".fid", "fid.sqlite")
	beforeJournal, err := os.ReadFile(journalPath)
	if err != nil {
		t.Fatalf("ReadFile(%q) error = %v", journalPath, err)
	}
	beforeDatabase, err := os.ReadFile(databasePath)
	if err != nil {
		t.Fatalf("ReadFile(%q) error = %v", databasePath, err)
	}

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	cmd := newRootCommandAt(workspace)
	cmd.SetOut(stdout)
	cmd.SetErr(stderr)
	cmd.SetArgs([]string{"--json", "rebuild", "--by", "workflow_runner"})

	err = cmd.Execute()
	if err == nil || err.Error() != "unknown flag: --by" {
		t.Fatalf("Execute() error = %v, want removed --by rejection", err)
	}

	afterJournal, err := os.ReadFile(journalPath)
	if err != nil {
		t.Fatalf("ReadFile(%q) after rejection error = %v", journalPath, err)
	}
	afterDatabase, err := os.ReadFile(databasePath)
	if err != nil {
		t.Fatalf("ReadFile(%q) after rejection error = %v", databasePath, err)
	}
	if !bytes.Equal(afterJournal, beforeJournal) || !bytes.Equal(afterDatabase, beforeDatabase) {
		t.Fatal("removed-flag rebuild changed journal or projection")
	}
}

func TestRebuildRejectsEveryFormerByValueWithoutMutation(t *testing.T) {
	workspace := t.TempDir()
	runInitCommand(t, workspace, "init")
	journalPath := mustJournalPath(t, workspace)
	beforeJournal, err := os.ReadFile(journalPath)
	if err != nil {
		t.Fatalf("ReadFile(%q) error = %v", journalPath, err)
	}

	for _, testCase := range []struct {
		name  string
		actor string
	}{
		{name: "unknown", actor: "intern"},
		{name: "denied", actor: "developer"},
	} {
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
			stdout := &bytes.Buffer{}
			stderr := &bytes.Buffer{}
			cmd := newRootCommandAt(workspace)
			cmd.SetOut(stdout)
			cmd.SetErr(stderr)
			cmd.SetArgs([]string{"--json", "rebuild", "--by", testCase.actor})
			err := cmd.Execute()
			if err == nil || err.Error() != "unknown flag: --by" {
				t.Fatalf("Execute() error = %v, want removed --by rejection", err)
			}
		})
	}

	afterJournal, err := os.ReadFile(journalPath)
	if err != nil {
		t.Fatalf("ReadFile(%q) after rejection error = %v", journalPath, err)
	}
	if !bytes.Equal(afterJournal, beforeJournal) {
		t.Fatal("removed-flag rebuild appended journal data")
	}
}

func TestRebuildJSONRejectsClosedRequestsBeforeProjectionWork(t *testing.T) {
	workspace := t.TempDir()
	runInitCommand(t, workspace, "init")
	journalPath := mustJournalPath(t, workspace)
	beforeJournal, err := os.ReadFile(journalPath)
	if err != nil {
		t.Fatalf("ReadFile(%q) error = %v", journalPath, err)
	}

	for _, testCase := range []struct {
		name         string
		args         []string
		wantArgument string
		wantValue    string
	}{
		{name: "positional argument", args: []string{"--json", "rebuild", "unexpected"}, wantArgument: "arguments", wantValue: "unexpected"},
		{name: "verbose JSON", args: []string{"--json", "rebuild", "--verbose"}, wantArgument: "verbose", wantValue: "true"},
	} {
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
			stdout := &bytes.Buffer{}
			stderr := &bytes.Buffer{}
			cmd := newRootCommandAt(workspace)
			cmd.SetOut(stdout)
			cmd.SetErr(stderr)
			cmd.SetArgs(testCase.args)
			if err := cmd.Execute(); err == nil {
				t.Fatal("Execute() error = nil, want invalid request rejection")
			}
			if stdout.Len() != 0 {
				t.Fatalf("stdout = %q, want empty", stdout.String())
			}
			var failure struct {
				OK        bool   `json:"ok"`
				Operation string `json:"operation"`
				Error     struct {
					Code      string `json:"code"`
					Category  string `json:"category"`
					Retryable bool   `json:"retryable"`
					Subject   *any   `json:"subject"`
					Details   struct {
						Argument string   `json:"argument"`
						Value    string   `json:"value"`
						Allowed  []string `json:"allowed"`
					} `json:"details"`
				} `json:"error"`
			}
			if err := json.Unmarshal(stderr.Bytes(), &failure); err != nil {
				t.Fatalf("Unmarshal(stderr) error = %v\nstderr=%q", err, stderr.String())
			}
			if failure.OK || failure.Operation != "project.rebuild" || failure.Error.Code != "invalid-argument" ||
				failure.Error.Category != "invalid-input" || failure.Error.Retryable || failure.Error.Subject != nil ||
				failure.Error.Details.Argument != testCase.wantArgument || failure.Error.Details.Value != testCase.wantValue ||
				len(failure.Error.Details.Allowed) == 0 {
				t.Fatalf("request failure = %+v, want typed invalid %s", failure, testCase.wantArgument)
			}
		})
	}

	afterJournal, err := os.ReadFile(journalPath)
	if err != nil {
		t.Fatalf("ReadFile(%q) after rejection error = %v", journalPath, err)
	}
	if !bytes.Equal(afterJournal, beforeJournal) {
		t.Fatal("invalid rebuild request appended journal data")
	}
}

func TestRebuildJSONMapsExplicitMissingProjectToSelectedPath(t *testing.T) {
	workingDirectory := t.TempDir()
	missingProject := filepath.Join(t.TempDir(), "missing-project")
	if err := os.Mkdir(missingProject, 0o755); err != nil {
		t.Fatalf("Mkdir(%q) error = %v", missingProject, err)
	}

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	cmd := newRootCommandAt(workingDirectory)
	cmd.SetOut(stdout)
	cmd.SetErr(stderr)
	cmd.SetArgs([]string{"--project", missingProject, "--json", "rebuild"})
	if err := cmd.Execute(); err == nil {
		t.Fatal("Execute() error = nil, want selected-project rejection")
	}
	if stdout.Len() != 0 {
		t.Fatalf("stdout = %q, want empty", stdout.String())
	}
	var failure struct {
		OK        bool   `json:"ok"`
		Operation string `json:"operation"`
		Error     struct {
			Code    string `json:"code"`
			Subject struct {
				Type string `json:"type"`
				ID   string `json:"id"`
			} `json:"subject"`
			Details struct {
				Path string `json:"path"`
			} `json:"details"`
		} `json:"error"`
	}
	if err := json.Unmarshal(stderr.Bytes(), &failure); err != nil {
		t.Fatalf("Unmarshal(stderr) error = %v", err)
	}
	if failure.OK || failure.Operation != "project.rebuild" || failure.Error.Code != "project-not-found" ||
		failure.Error.Subject.Type != "project" || failure.Error.Subject.ID != missingProject ||
		failure.Error.Details.Path != missingProject {
		t.Fatalf("project failure = %+v, want selected missing project %q", failure, missingProject)
	}
}

func TestRebuildJSONReportsRetryableHeldLockAndSucceedsAfterRelease(t *testing.T) {
	workspace := t.TempDir()
	runInitCommand(t, workspace, "init")
	journalPath := mustJournalPath(t, workspace)
	beforeJournal, err := os.ReadFile(journalPath)
	if err != nil {
		t.Fatalf("ReadFile(%q) error = %v", journalPath, err)
	}

	lockedStdout := &bytes.Buffer{}
	lockedStderr := &bytes.Buffer{}
	var lockedErr error
	err = journal.WithLock(workspace, journal.LockOwner{Command: "test held lock"}, func(_ *journal.WriteLock) error {
		cmd := newRootCommandAt(workspace)
		cmd.SetOut(lockedStdout)
		cmd.SetErr(lockedStderr)
		cmd.SetArgs([]string{"--json", "rebuild"})
		lockedErr = cmd.Execute()
		return nil
	})
	if err != nil {
		t.Fatalf("WithLock() error = %v", err)
	}
	if lockedErr == nil {
		t.Fatal("Execute() under held lock error = nil, want lock rejection")
	}
	if lockedStdout.Len() != 0 {
		t.Fatalf("held-lock stdout = %q, want empty", lockedStdout.String())
	}
	assertEnvelopeKeySet(t, lockedStderr.Bytes(), []string{"contract_version", "error", "ok", "operation"})
	assertNestedEnvelopeKeySet(t, lockedStderr.Bytes(), "error.details", []string{"lock_path"})
	var lockedFailure struct {
		OK        bool   `json:"ok"`
		Operation string `json:"operation"`
		Error     struct {
			Code      string `json:"code"`
			Category  string `json:"category"`
			Retryable bool   `json:"retryable"`
			Subject   *any   `json:"subject"`
			Details   struct {
				LockPath string `json:"lock_path"`
			} `json:"details"`
		} `json:"error"`
	}
	if err := json.Unmarshal(lockedStderr.Bytes(), &lockedFailure); err != nil {
		t.Fatalf("Unmarshal(held-lock stderr) error = %v", err)
	}
	if lockedFailure.OK || lockedFailure.Operation != "project.rebuild" ||
		lockedFailure.Error.Code != "lock-unavailable" || lockedFailure.Error.Category != "conflict" ||
		!lockedFailure.Error.Retryable || lockedFailure.Error.Subject != nil ||
		lockedFailure.Error.Details.LockPath != ".fid/write.lock" {
		t.Fatalf("held-lock failure = %+v, want retryable typed lock-unavailable", lockedFailure)
	}
	afterHeldJournal, err := os.ReadFile(journalPath)
	if err != nil {
		t.Fatalf("ReadFile(%q) after held-lock rejection error = %v", journalPath, err)
	}
	if !bytes.Equal(afterHeldJournal, beforeJournal) {
		t.Fatal("held-lock rebuild appended journal data")
	}

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	cmd := newRootCommandAt(workspace)
	cmd.SetOut(stdout)
	cmd.SetErr(stderr)
	cmd.SetArgs([]string{"--json", "rebuild"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() after lock release error = %v", err)
	}
	if stderr.Len() != 0 {
		t.Fatalf("retry stderr = %q, want empty", stderr.String())
	}
	assertEnvelopeKeySet(t, stdout.Bytes(), []string{"contract_version", "event", "ok", "operation", "result"})
	assertNestedEnvelopeKeySet(t, stdout.Bytes(), "result", []string{"events_replayed", "files_replayed", "schema_version", "upgraded_events"})
	var success struct {
		OK        bool   `json:"ok"`
		Operation string `json:"operation"`
		Event     *any   `json:"event"`
	}
	if err := json.Unmarshal(stdout.Bytes(), &success); err != nil {
		t.Fatalf("Unmarshal(retry stdout) error = %v", err)
	}
	if !success.OK || success.Operation != "project.rebuild" || success.Event != nil {
		t.Fatalf("retry success = %+v, want project.rebuild with event null", success)
	}
}

func TestRebuildJSONMapsReplayStopWithClosedDetails(t *testing.T) {
	workspace := t.TempDir()
	runInitCommand(t, workspace, "init")
	writeCommandFixture(t, workspace, "corrupt-tail.jsonl", "2026-06-05.jsonl")

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	cmd := newRootCommandAt(workspace)
	cmd.SetOut(stdout)
	cmd.SetErr(stderr)
	cmd.SetArgs([]string{"--json", "rebuild"})
	if err := cmd.Execute(); err == nil {
		t.Fatal("Execute() error = nil, want replay-stop failure")
	}
	if stdout.Len() != 0 {
		t.Fatalf("stdout = %q, want empty", stdout.String())
	}
	assertEnvelopeKeySet(t, stderr.Bytes(), []string{"contract_version", "error", "ok", "operation"})
	assertNestedEnvelopeKeySet(t, stderr.Bytes(), "error.details", []string{"event_id", "journal_file", "line"})
	var failure struct {
		OK        bool   `json:"ok"`
		Operation string `json:"operation"`
		Error     struct {
			Code      string `json:"code"`
			Category  string `json:"category"`
			Retryable bool   `json:"retryable"`
			Subject   *any   `json:"subject"`
			Details   struct {
				JournalFile string  `json:"journal_file"`
				Line        *int    `json:"line"`
				EventID     *string `json:"event_id"`
			} `json:"details"`
		} `json:"error"`
	}
	if err := json.Unmarshal(stderr.Bytes(), &failure); err != nil {
		t.Fatalf("Unmarshal(stderr) error = %v", err)
	}
	if failure.OK || failure.Operation != "project.rebuild" || failure.Error.Code != "journal-replay-stopped" ||
		failure.Error.Category != "integrity" || failure.Error.Retryable || failure.Error.Subject != nil ||
		failure.Error.Details.JournalFile != ".fid/journal/2026-06-05.jsonl" || failure.Error.Details.Line == nil ||
		*failure.Error.Details.Line != 2 || failure.Error.Details.EventID != nil {
		t.Fatalf("replay failure = %+v, want closed journal-replay-stopped details", failure)
	}
}

func TestRebuildCommandReplaysJournalAndReportsCounts(t *testing.T) {
	workspace := newRebuildWorkspace(t)
	writeCommandFixture(t, workspace, "replay-v1.golden.jsonl", "2026-06-05.jsonl")

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	cmd := newRootCommandAt(workspace)
	cmd.SetOut(stdout)
	cmd.SetErr(stderr)
	cmd.SetArgs([]string{"rebuild"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	if !strings.Contains(stdout.String(), "rebuild complete") {
		t.Fatalf("stdout = %q, want rebuild completion", stdout.String())
	}
	if stderr.Len() != 0 {
		t.Fatalf("stderr = %q, want empty", stderr.String())
	}
	if got := stdout.String(); !strings.Contains(got, "replayed 2 events from 1 journal files") {
		t.Fatalf("stdout = %q, want replay count line", got)
	}
	if got := stdout.String(); !strings.Contains(got, "schema version: "+fmt.Sprint(fiddb.ActiveMigrationVersion)) {
		t.Fatalf("stdout = %q, want active schema version %d", got, fiddb.ActiveMigrationVersion)
	}

	conn, err := fiddb.Open(filepath.Join(workspace, ".fid", "fid.sqlite"))
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	defer closeCommandDB(t, conn)

	var title string
	if err := conn.QueryRow(`SELECT title FROM issues WHERE id = ?`, "FID-a1b2c3").Scan(&title); err != nil {
		t.Fatalf("query rebuilt issue: %v", err)
	}
	if title != "Implement replay" {
		t.Fatalf("rebuilt issue title = %q, want %q", title, "Implement replay")
	}
}

func TestRebuildCommandReportsCorruptJournal(t *testing.T) {
	workspace := newRebuildWorkspace(t)
	writeCommandFixture(t, workspace, "corrupt-tail.jsonl", "2026-06-05.jsonl")

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	cmd := newRootCommandAt(workspace)
	cmd.SetOut(stdout)
	cmd.SetErr(stderr)
	cmd.SetArgs([]string{"rebuild"})

	err := cmd.Execute()
	if err == nil {
		t.Fatal("Execute() error = nil, want corrupt journal error")
	}
	if got := stderr.String(); !strings.Contains(got, "malformed journal entry at .fid/journal/2026-06-05.jsonl:2") {
		t.Fatalf("stderr = %q, want corrupt journal line", got)
	}
}

func TestRebuildCommandReportsPreconditionConflict(t *testing.T) {
	workspace := t.TempDir()
	runInitCommand(t, workspace, "init")
	writeCommandFixture(t, workspace, "conflicting-precondition.jsonl", "2026-06-05.jsonl")

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	cmd := newRootCommandAt(workspace)
	cmd.SetOut(stdout)
	cmd.SetErr(stderr)
	cmd.SetArgs([]string{"rebuild"})

	err := cmd.Execute()
	if err == nil {
		t.Fatal("Execute() error = nil, want precondition conflict")
	}

	if stdout.Len() != 0 {
		t.Fatalf("stdout = %q, want empty", stdout.String())
	}

	got := stderr.String()
	for _, want := range []string{
		"error: replay precondition conflict at .fid/journal/2026-06-05.jsonl:2",
		"event: EVT-20260605-000002 (issue.transitioned)",
		"entity: issues pk={id:FID-a1b2c3}",
		"expected: {status:in_progress}",
		"actual: {status:todo}",
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("stderr = %q, want substring %q", got, want)
		}
	}
}

func newRebuildWorkspace(t *testing.T) string {
	t.Helper()

	workspace := t.TempDir()
	if err := os.MkdirAll(filepath.Join(workspace, ".fid", "journal"), 0o755); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}
	return workspace
}

func commandFixturePath(t *testing.T, name string) string {
	t.Helper()

	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller(0) failed")
	}

	return filepath.Clean(filepath.Join(filepath.Dir(currentFile), "..", "..", "testdata", "journal", name))
}

func writeCommandFixture(t *testing.T, workspaceRoot, fixtureName, journalFileName string) {
	t.Helper()

	payload, err := os.ReadFile(commandFixturePath(t, fixtureName))
	if err != nil {
		t.Fatalf("ReadFile(%q) error = %v", fixtureName, err)
	}

	journalPath := filepath.Join(workspaceRoot, ".fid", "journal", journalFileName)
	if err := os.WriteFile(journalPath, payload, 0o644); err != nil {
		t.Fatalf("WriteFile(%q) error = %v", journalPath, err)
	}
}

func closeCommandDB(t *testing.T, conn *sql.DB) {
	t.Helper()

	if err := conn.Close(); err != nil {
		t.Fatalf("Close() error = %v", err)
	}
}
