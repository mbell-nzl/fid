package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/mbell-nzl/fid/internal/checkrecord"
	"github.com/mbell-nzl/fid/internal/checkrun"
	"github.com/mbell-nzl/fid/internal/config"
	"github.com/mbell-nzl/fid/internal/gittest"
	"github.com/mbell-nzl/fid/internal/holdwork"
)

type restoredForTest struct {
	ProjectID         string `json:"project_id"`
	Commit            string `json:"commit"`
	ArchiveDigest     string `json:"archive_digest"`
	EventCount        int    `json:"event_count"`
	TerminalEventHash string `json:"terminal_event_hash"`
	GenesisEventHash  string `json:"genesis_event_hash"`
}

// runRestoreWith runs fid restore in-process with environ as its
// environment, through the same seam as fid backup.
func runRestoreWith(environ []string, workspace string, args ...string) commandExecutionResult {
	root := newRootCommandAt(workspace)
	root.SetContext(context.WithValue(root.Context(), backupEnvironmentKey{}, backupEnvironment{environ: environ, lockWait: time.Minute}))
	stdout, stderr := &bytes.Buffer{}, &bytes.Buffer{}
	root.SetOut(stdout)
	root.SetErr(stderr)
	root.SetArgs(append(args, "restore"))
	err := root.Execute()
	return commandExecutionResult{stdout: stdout.String(), stderr: stderr.String(), err: err}
}

// restoreWith runs fid restore with environ and a temporary directory of
// its own, and checks the restore left nothing in it.
func restoreWith(t *testing.T, environ []string, workspace string, args ...string) commandExecutionResult {
	t.Helper()
	temporary := t.TempDir()
	result := runRestoreWith(append(append([]string{}, environ...), "TMPDIR="+temporary), workspace, args...)
	if left, err := os.ReadDir(temporary); err != nil || len(left) != 0 {
		t.Fatalf("restore left %v in its temporary directory: %v", left, err)
	}
	return result
}

func (fixture *backupFixture) restore(workspace string, args ...string) commandExecutionResult {
	fixture.t.Helper()
	return restoreWith(fixture.t, fixture.env, workspace, args...)
}

func decodeRestoreSuccess(t *testing.T, result commandExecutionResult) (restoredForTest, string) {
	t.Helper()
	if result.err != nil || result.stderr != "" {
		t.Fatalf("restore = %v, stderr=%q", result.err, result.stderr)
	}
	var envelope struct {
		OK        bool            `json:"ok"`
		Operation string          `json:"operation"`
		Event     json.RawMessage `json:"event"`
		Result    struct {
			Restore    json.RawMessage `json:"restore"`
			LedgerPath string          `json:"ledger_path"`
		} `json:"result"`
	}
	if err := json.Unmarshal([]byte(result.stdout), &envelope); err != nil {
		t.Fatalf("decode restore envelope %q: %v", result.stdout, err)
	}
	if !envelope.OK || envelope.Operation != "project.restore" || string(envelope.Event) != "null" {
		t.Fatalf("restore envelope = %s, want ok project.restore with no event", result.stdout)
	}
	var keys map[string]json.RawMessage
	if err := json.Unmarshal(envelope.Result.Restore, &keys); err != nil || len(keys) != 6 {
		t.Fatalf("restore result %s has %d keys, want all 6: %v", envelope.Result.Restore, len(keys), err)
	}
	var restored restoredForTest
	decoder := json.NewDecoder(bytes.NewReader(envelope.Result.Restore))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&restored); err != nil {
		t.Fatalf("decode restore result %s: %v", envelope.Result.Restore, err)
	}
	return restored, envelope.Result.LedgerPath
}

// v8IDs names what the V8 fixture's read commands take.
type v8IDs struct {
	issues            []string
	envelope, receipt string
}

// v8Ledger builds a ledger holding issues, a label, a link, and a gate run
// whose receipt keeps an Artifact, with the project's .fid.toml committed.
func v8Ledger(t *testing.T) (string, v8IDs) {
	t.Helper()
	workspace, issueID, envelopeID := runEnvelopeForTestWithDefinition(t, func(definition map[string]any) {
		definition["argv"] = []string{"/bin/true"}
		definition["artifact_paths"] = []string{"artifact.txt"}
	})
	if err := os.WriteFile(filepath.Join(workspace, "artifact.txt"), []byte("kept by the receipt\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	terminal, result := executeGateRunForTest(t, workspace, issueID, envelopeID, "v8-gate")
	if result.err != nil || terminal.Data.ReceiptID == "" {
		t.Fatalf("gate run = %#v", result)
	}
	runFIDCommand(t, workspace, "label", "add", "restore", "--by", "developer")
	_, created := runFIDCommand(t, workspace, "issue", "create", "--kind", "task", "--title", "Restored with the ledger",
		"--label", "restore", "--actor", "agent:developer")
	second := parseCreatedIssueID(t, created)
	runFIDCommand(t, workspace, "link", "add", second, "discovered_from", issueID, "--from-type", "issue", "--to-type", "issue", "--by", "developer")
	// Two holds: one its task's cancellation ended, and one still open. The
	// held tasks block the second task, so fid next takes them first.
	ids := []string{issueID, second}
	for _, priority := range []string{"0", "1"} {
		_, created := runFIDCommand(t, workspace, "issue", "create", "--kind", "task", "--title", "Held "+priority, "--lane", "routine", "--priority", priority, "--actor", "human:mark")
		ids = append(ids, parseCreatedIssueID(t, created))
		runFIDCommand(t, workspace, "link", "add", ids[len(ids)-1], "blocks", second, "--from-type", "issue", "--to-type", "issue", "--actor", "human:mark")
	}
	runFIDCommand(t, workspace, "next", "--claim", "--actor", "agent:developer")
	runFIDCommand(t, workspace, "issue", "transition", ids[2], "--to", "cancelled", "--actor", "human:mark", "--reason", "Test cancellation")
	runFIDCommand(t, workspace, "next", "--claim", "--actor", "agent:developer")
	// Keep each lifecycle family and both kinds of person's decision in V8.
	// A heartbeat, park and linked resume precede a person's live takeover.
	answer, _ := decodeNext(t, executeCommandForTest(workspace, "--json", "next", "--claim", "--actor", "agent:developer"), "next.claim")
	originalHold := answer["hold"].(map[string]any)["id"].(string)
	runFIDCommand(t, workspace, "heartbeat", "--hold", originalHold, "--actor", "agent:developer")
	runFIDCommand(t, workspace, "issue", "update", ids[3], "--lane", "standard", "--actor", "agent:developer")
	runFIDCommand(t, workspace, "start", "--hold", originalHold, "--actor", "agent:developer", "--idempotency-key", "v8-start")
	runFIDCommand(t, workspace, "park", "--hold", originalHold, "--actor", "agent:developer", "--interrupted", "--reason", "Usage limit", "--idempotency-key", "v8-park-resume")
	resumed, _ := decodeNext(t, executeCommandForTest(workspace, "--json", "next", "--claim", "--actor", "agent:developer"), "next.claim")
	if resumed["hold"].(map[string]any)["replaces"] != originalHold {
		t.Fatalf("V8 resume = %v", resumed)
	}
	runFIDCommand(t, workspace, "takeover", ids[3], "--actor", "human:mark", "--reason", "Agent stuck")
	// Another worktree keeps a parked task, so a restore must preserve an
	// ended hold's reservation and the worktree path it belongs to.
	_, epicCreated := runFIDCommand(t, workspace, "issue", "create", "--kind", "epic", "--title", "Started through its task", "--scheduling", "prioritised", "--actor", "human:mark")
	epicID := parseCreatedIssueID(t, epicCreated)
	ids = append(ids, epicID)
	_, created = runFIDCommand(t, workspace, "issue", "create", "--kind", "task", "--epic", epicID, "--lane", "assured", "--title", "Kept during backup", "--actor", "human:mark")
	keptTask := parseCreatedIssueID(t, created)
	ids = append(ids, keptTask)
	other := gittest.AddWorktree(t, workspace)
	took, _ := decodeNext(t, executeCommandForTest(other, "--json", "takeover", keptTask, "--actor", "human:mark", "--reason", "Keep this worktree"), "task.takeover")
	runFIDCommand(t, other, "start", "--hold", took["hold_id"].(string), "--actor", "human:mark", "--idempotency-key", "v8-epic-start")
	runFIDCommand(t, other, "park", "--hold", took["hold_id"].(string), "--actor", "human:mark", "--interrupted", "--reason", "Usage limit", "--idempotency-key", "v8-kept")
	// A waiting task retains the first answer and the next open question.
	resumedWaiting, _ := decodeNext(t, executeCommandForTest(other, "--json", "next", "--claim", "--actor", "human:mark"), "next.claim")
	waitingHold := resumedWaiting["hold"].(map[string]any)["id"].(string)
	runFIDCommand(t, other, "park", "--hold", waitingHold, "--actor", "human:mark", "--waiting", "Include authentication failures?", "--idempotency-key", "v8-question-1")
	runFIDCommand(t, workspace, "answer", keptTask, "Exclude them.", "--actor", "human:mark", "--idempotency-key", "v8-answer")
	resumedWaiting, _ = decodeNext(t, executeCommandForTest(other, "--json", "next", "--claim", "--actor", "human:mark"), "next.claim")
	if resumedWaiting["answer"] == nil {
		t.Fatal("V8 resume lost answer")
	}
	waitingHold = resumedWaiting["hold"].(map[string]any)["id"].(string)
	runFIDCommand(t, other, "park", "--hold", waitingHold, "--actor", "human:mark", "--waiting", "Include timeouts?", "--idempotency-key", "v8-question-2")
	// Record missing and explicitly empty test paths, so backup, export,
	// replay and restore must preserve their different digests and setup state.
	appendPolicy(t, workspace, "\n[review]\ntest_adequacy = []\n[checks.restore]\nrun = ['/bin/sleep', '0.15']\ntimeout = '15m'\nartifacts = ['lean-artifact.txt']\n[checks.missing]\nrun = ['/bin/true']\ntimeout = '1m'\nartifacts = ['missing-artifact.txt']\n[lanes.standard]\nchecks = []\n")
	policyOrigin(t, workspace)
	runFIDCommand(t, workspace, "next", "--claim", "--actor", "agent:developer")
	writeBackupTestFile(t, filepath.Join(workspace, "lean-artifact.txt"), "kept by the lean check\n")
	identityGit(t, workspace, "add", ".")
	identityGit(t, workspace, "-c", "user.name=fid test", "-c", "user.email=fid@example.test", "commit", "-q", "-m", "Lean check files")
	current, _ := decodeNext(t, executeCommandForTest(workspace, "--json", "next", "--claim", "--actor", "human:mark"), "next.claim")
	_, checkErr := checkrun.Run(context.Background(), checkrun.Options{WorkspaceRoot: workspace, MigrationsDir: issueMigrationsDir(), HoldID: current["hold"].(map[string]any)["id"].(string), Actor: "human:mark", IdempotencyKey: "v8-lean-check", Names: []string{"restore", "missing"}, Version: checkrun.Version(), RenewInterval: 20 * time.Millisecond})
	var refusal *checkrecord.Refusal
	if !errors.As(checkErr, &refusal) || refusal.Code != "check-not-qualified" {
		t.Fatalf("V8 lean check = %v", checkErr)
	}
	if _, err := checkrun.Run(context.Background(), checkrun.Options{WorkspaceRoot: workspace, MigrationsDir: issueMigrationsDir(), HoldID: current["hold"].(map[string]any)["id"].(string), Actor: "human:mark", IdempotencyKey: "v8-empty-check"}); err != nil {
		t.Fatalf("V8 empty check = %v", err)
	}
	_, stoppedErr := checkrun.Run(context.Background(), checkrun.Options{WorkspaceRoot: workspace, MigrationsDir: issueMigrationsDir(), HoldID: current["hold"].(map[string]any)["id"].(string), Actor: "human:mark", IdempotencyKey: "v8-stopped-check", Names: []string{"missing", "restore"}, Version: checkrun.Version(), After: func(checkrecord.Record) error {
		_, err := holdwork.Takeover(holdwork.Options{WorkspaceRoot: workspace, MigrationsDir: issueMigrationsDir(), TaskID: ids[3], Actor: "human:operator", Reason: "Stop between checks"})
		return err
	}})
	if !errors.As(stoppedErr, &refusal) || refusal.Code != "hold-not-current" {
		t.Fatalf("V8 stopped check = %v", stoppedErr)
	}

	return workspace, v8IDs{issues: ids, envelope: envelopeID, receipt: terminal.Data.ReceiptID}
}

// v8State is everything V8 compares in a ledger, apart from its projection.
type v8State struct {
	journalNames []string
	events       []byte
	files        map[string][]byte // every other ledger file, by path in the ledger
	reads        map[string]string // each read command's JSON output
	genesis      string
}

func captureV8State(t *testing.T, workspace string, ids v8IDs) v8State {
	t.Helper()
	ledger := config.LedgerDir(workspace)
	state := v8State{files: map[string][]byte{}, reads: map[string]string{}, genesis: ledgerGenesisHash(t, workspace)}
	err := filepath.WalkDir(ledger, func(path string, entry fs.DirEntry, err error) error {
		if err != nil || entry.IsDir() {
			return err
		}
		relative, err := filepath.Rel(ledger, path)
		if err != nil {
			return err
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		if directory, name := filepath.Split(relative); directory == config.JournalDirName+"/" {
			state.journalNames = append(state.journalNames, name)
		} else if !strings.HasPrefix(relative, config.DatabaseFileName) {
			state.files[relative] = data
		}
		return nil
	})
	if err != nil {
		t.Fatalf("read the ledger: %v", err)
	}
	// Replay reads the journal's files in name order, so their contents in
	// that order are the Events in order, however the files are cut.
	sort.Strings(state.journalNames)
	for _, name := range state.journalNames {
		data, err := os.ReadFile(filepath.Join(ledger, config.JournalDirName, name))
		if err != nil {
			t.Fatal(err)
		}
		state.events = append(state.events, data...)
	}
	reads := [][]string{{"issue", "list"}, {"label", "list"}, {"ready"}, {"workflow", "list"}, {"activity", "list"},
		{"claim", "list"}, {"run", "list"}, {"receipt", "list"}, {"receipt", "get", ids.receipt}, {"envelope", "get", ids.envelope},
		{"stop", "list"}, {"attention", "list"}, {"concern", "list"}, {"deviation", "list"}, {"decision", "list"},
		{"human-outcome", "list"}, {"review", "list"}, {"finding", "list"}, {"requirement", "list"}, {"maintenance", "check"}}
	for _, id := range ids.issues {
		reads = append(reads, []string{"issue", "get", id}, []string{"issue", "ac", "list", id}, []string{"link", "list", id, "--type", "issue"})
	}
	for _, args := range reads {
		result := executeCommandForTest(workspace, append([]string{"--json"}, args...)...)
		if result.err != nil || result.stderr != "" {
			t.Fatalf("fid %v = %#v", args, result)
		}
		state.reads[strings.Join(args, " ")] = result.stdout
	}
	// The holds, as the projection replayed from the journal has them. No
	// read command lists them, and fid next's answer depends on the
	// worktree it runs in.
	holds := executeCommandForTest(workspace, "sql", "SELECT * FROM holds ORDER BY id")
	if holds.err != nil || holds.stderr != "" {
		t.Fatalf("fid sql holds = %#v", holds)
	}
	state.reads["holds"] = holds.stdout
	decisions := executeCommandForTest(workspace, "sql", "SELECT * FROM person_decisions ORDER BY id")
	if decisions.err != nil || decisions.stderr != "" {
		t.Fatalf("fid sql person decisions = %#v", decisions)
	}
	state.reads["person decisions"] = decisions.stdout
	questions := executeCommandForTest(workspace, "sql", "SELECT * FROM questions ORDER BY rowid")
	if questions.err != nil || questions.stderr != "" {
		t.Fatalf("fid sql questions = %#v", questions)
	}
	state.reads["questions"] = questions.stdout
	// Policy snapshots are journal-backed too. Compare all their projected
	// fields, so a restore cannot silently lose the policy fid trusted.
	snapshots := executeCommandForTest(workspace, "sql", "SELECT * FROM policy_snapshots ORDER BY sequence")
	if snapshots.err != nil || snapshots.stderr != "" {
		t.Fatalf("fid sql policy snapshots = %#v", snapshots)
	}
	state.reads["policy snapshots"] = snapshots.stdout
	for _, table := range []string{"check_requests", "check_starts", "check_processes", "check_results", "check_files"} {
		checkRows := executeCommandForTest(workspace, "sql", "SELECT * FROM "+table+" ORDER BY rowid")
		if checkRows.err != nil || checkRows.stderr != "" || strings.Contains(checkRows.stdout, `"rows":[]`) {
			t.Fatalf("V8 %s = %#v", table, checkRows)
		}
		state.reads[table] = checkRows.stdout
	}
	rows := policySnapshots(t, workspace)
	if len(rows) != 2 || rows[0].Digest == rows[1].Digest {
		t.Fatalf("V8 policy snapshots = %#v, want distinct missing/empty test paths", rows)
	}
	for index, row := range rows {
		var effective struct {
			TestAdequacy    []string `json:"test_adequacy"`
			TestAdequacySet *bool    `json:"test_adequacy_set"`
		}
		if err := json.Unmarshal([]byte(row.Effective), &effective); err != nil {
			t.Fatal(err)
		}
		if effective.TestAdequacySet == nil || *effective.TestAdequacySet != (index == 1) || effective.TestAdequacy == nil || len(effective.TestAdequacy) != 0 {
			t.Fatalf("V8 review policy %d = %#v", index, effective)
		}
	}
	return state
}

// exportForTest runs fid export create on the ledger and returns the
// archive's bytes and the export's summary.
func exportForTest(t *testing.T, workspace string) ([]byte, map[string]any) {
	t.Helper()
	path := filepath.Join(t.TempDir(), "export.zip")
	result := executeCommandForTest(workspace, "--json", "export", "create", path)
	var envelope struct {
		Result struct {
			Export map[string]any `json:"export"`
		} `json:"result"`
	}
	if result.err != nil || json.Unmarshal([]byte(result.stdout), &envelope) != nil {
		t.Fatalf("export create = %#v", result)
	}
	archive, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return archive, envelope.Result.Export
}

// V8 (check:restore): back up, wipe, restore and compare. A restore gives
// back the ledger that was backed up: the same export and backup, the same
// Events in the same order, the same config and Artifacts, and the same
// answers from read commands. It goes through fid backup and fid restore
// only, so it holds whatever the backup folder's format.
func TestV8BackupWipeRestoreYieldsIdenticalState(t *testing.T) {
	t.Parallel()
	for _, test := range []struct {
		name string
		// wipe loses the ledger, and returns where to restore it and the
		// environment to restore it with.
		wipe func(t *testing.T, fixture *backupFixture, original string) (string, []string)
	}{
		{"new machine", func(t *testing.T, fixture *backupFixture, original string) (string, []string) {
			workspace := gittest.NewWorkspace(t)
			identity, err := os.ReadFile(filepath.Join(original, config.ProjectIdentityFileName))
			if err != nil {
				t.Fatal(err)
			}
			writeBackupTestFile(t, filepath.Join(workspace, config.ProjectIdentityFileName), string(identity))
			fixture.git("-C", workspace, "add", config.ProjectIdentityFileName)
			fixture.git("-C", workspace, "-c", "user.name=human", "-c", "user.email=human@example.test", "commit", "-q", "-m", "identity")
			return workspace, backupTestEnvironment(fixture.home, fixture.configHome, filepath.Join(fixture.base, "new-machine-cache"))
		}},
		{"in place", func(t *testing.T, fixture *backupFixture, original string) (string, []string) {
			if err := os.RemoveAll(config.LedgerDir(original)); err != nil {
				t.Fatal(err)
			}
			return original, fixture.env
		}},
	} {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			fixture := newBackupFixture(t)
			original, ids := v8Ledger(t)
			id, err := config.ReadProjectIdentity(original)
			if err != nil {
				t.Fatal(err)
			}
			fixture.setTarget(fixture.target)
			ours, _ := decodeBackupSuccess(t, fixture.backup(original, "--json"))
			// Another project's backup shares the target, and is its tip.
			other, otherID := fixture.project()
			theirs, _ := decodeBackupSuccess(t, fixture.backup(other, "--json"))
			before := captureV8State(t, original, ids)
			archiveBefore, exported := exportForTest(t, original)

			workspace, environ := test.wipe(t, fixture, original)
			if _, err := os.Lstat(config.LedgerDir(workspace)); !os.IsNotExist(err) {
				t.Fatalf("the wipe left a ledger at %s: %v", config.LedgerDir(workspace), err)
			}
			restored, ledgerPath := decodeRestoreSuccess(t, restoreWith(t, environ, workspace, "--json"))

			if restored.ProjectID != id || restored.Commit != *theirs.Commit || restored.ArchiveDigest != *ours.ArchiveDigest ||
				restored.EventCount != bytes.Count(before.events, []byte("\n")) || restored.EventCount != int(exported["event_count"].(float64)) ||
				restored.TerminalEventHash != exported["terminal_event_hash"] || restored.GenesisEventHash != before.genesis {
				t.Fatalf("restore result = %+v, want %s at the target's tip %s, with the backed-up archive %s and the original's %d Events, terminal %v and genesis %s",
					restored, id, *theirs.Commit, *ours.ArchiveDigest, bytes.Count(before.events, []byte("\n")), exported["terminal_event_hash"], before.genesis)
			}
			if ledgerPath != config.LedgerDir(workspace) {
				t.Fatalf("ledger_path = %s, want %s", ledgerPath, config.LedgerDir(workspace))
			}
			if *theirs.ProjectID != otherID || fixture.targetMain() != *theirs.Commit {
				t.Fatalf("target main = %s, want the other project's backup %s untouched", fixture.targetMain(), *theirs.Commit)
			}

			after := captureV8State(t, workspace, ids)
			assertV8ExpectedDifferences(t, before, after)
			assertV8SameState(t, before, after)
			if archiveAfter, _ := exportForTest(t, workspace); !bytes.Equal(archiveAfter, archiveBefore) {
				t.Fatal("the restored ledger's export differs from the original's")
			}
			unchanged, _ := decodeBackupSuccess(t, runBackupWith(environ, workspace, time.Minute, "--json"))
			if unchanged.Outcome != "unchanged" || *unchanged.Commit != restored.Commit || *unchanged.ArchiveDigest != *ours.ArchiveDigest {
				t.Fatalf("first backup after the restore = %+v, want unchanged at %s", unchanged, restored.Commit)
			}
		})
	}
}

var journalPieceName = regexp.MustCompile(`^[0-9]{8}\.jsonl$`)

// assertV8ExpectedDifferences checks what a restore changes on purpose: the
// journal arrives as the export's numbered pieces, the ledger records the
// import, and it starts without a backup record.
func assertV8ExpectedDifferences(t *testing.T, before, after v8State) {
	t.Helper()
	for index, name := range after.journalNames {
		if !journalPieceName.MatchString(name) || name != pieceNameForTest(index+1) {
			t.Fatalf("restored journal files = %v, want the export's pieces 00000001.jsonl onwards", after.journalNames)
		}
	}
	for _, name := range before.journalNames {
		if journalPieceName.MatchString(name) {
			t.Fatalf("original journal files = %v, want dated files", before.journalNames)
		}
	}
	if _, ok := before.files["backups.jsonl"]; !ok {
		t.Fatal("the original ledger has no backup record")
	}
	if _, ok := after.files["backups.jsonl"]; ok {
		t.Fatal("the restored ledger has a backup record before its first backup")
	}
	if _, ok := after.files["import.marker.json"]; !ok {
		t.Fatal("the restored ledger has no import marker")
	}
	if _, ok := before.files["import.marker.json"]; ok {
		t.Fatal("the original ledger has an import marker")
	}
}

func pieceNameForTest(index int) string { return fmt.Sprintf("%08d.jsonl", index) }

// assertV8SameState checks everything else is identical.
func assertV8SameState(t *testing.T, before, after v8State) {
	t.Helper()
	if !bytes.Equal(after.events, before.events) {
		t.Fatal("the restored journal's Events differ from the original's")
	}
	if after.genesis != before.genesis {
		t.Fatalf("restored genesis = %s, want %s", after.genesis, before.genesis)
	}
	delete(before.files, "backups.jsonl")
	delete(after.files, "import.marker.json")
	artifacts := 0
	for name, data := range before.files {
		if !bytes.Equal(after.files[name], data) {
			t.Errorf("restored %s differs from the original's", name)
		}
		if strings.HasPrefix(name, "artifacts/") {
			artifacts++
		}
	}
	if len(after.files) != len(before.files) || artifacts == 0 || before.files[config.ConfigFileName] == nil {
		t.Fatalf("restored ledger files %v, want the original's %v, with its config and at least one Artifact", sortedFileNames(after.files), sortedFileNames(before.files))
	}
	for args, output := range before.reads {
		if after.reads[args] != output {
			t.Errorf("fid --json %s after the restore:\n%s\nbefore:\n%s", args, after.reads[args], output)
		}
	}
	if strings.Count(before.reads["holds"], `"end_reason":"task-cancelled"`) != 1 || strings.Count(before.reads["holds"], `"end_reason":null`) != 1 {
		t.Fatalf("original holds = %s, want one its task's cancellation ended and one open", before.reads["holds"])
	}
	if strings.Count(before.reads["questions"], `"resolution":"answered"`) != 1 || strings.Count(before.reads["questions"], `"resolution":null`) != 1 || strings.Count(before.reads["person decisions"], `"kind":"answer"`) != 1 {
		t.Fatalf("V8 waiting history missing: %s", before.reads["questions"])
	}
	if strings.Count(before.reads["person decisions"], `"kind":"takeover"`) != 3 || strings.Count(before.reads["person decisions"], `"kind":"cancel"`) != 1 || strings.Count(before.reads["holds"], `"kept":1`) != 1 {
		t.Fatalf("V8 lifecycle missing: holds=%s decisions=%s", before.reads["holds"], before.reads["person decisions"])
	}
	if !strings.Contains(before.reads["check_requests"], "hold-not-current") || !strings.Contains(before.reads["check_results"], "artifact-not-kept") {
		t.Fatalf("V8 is missing its stopped request or refused file: %s", before.reads["check_requests"])
	}
	if !bytes.Contains(before.events, []byte(`"event_type":"check.completed"`)) {
		t.Fatal("V8 is missing its completed empty check")
	}
	if strings.Count(before.reads["policy snapshots"], `"source_ref":"built-in defaults"`) != 1 {
		t.Fatalf("original policy snapshots = %s, want one built-in policy", before.reads["policy snapshots"])
	}
	if !strings.Contains(after.reads["maintenance check"], `"status":"pass"`) || !strings.Contains(after.reads["maintenance check"], `"errors":[]`) {
		t.Fatalf("maintenance check after the restore = %s, want a pass", after.reads["maintenance check"])
	}
}

func sortedFileNames(files map[string][]byte) []string {
	names := make([]string, 0, len(files))
	for name := range files {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

func TestRestoreTextOutputNamesFolderCommitAndLedger(t *testing.T) {
	t.Parallel()
	fixture := newBackupFixture(t)
	original, id := fixture.project()
	fixture.setTarget(fixture.target)
	pushed, _ := decodeBackupSuccess(t, fixture.backup(original, "--json"))
	workspace := fixture.restoreWorkspace(id)

	result := fixture.restore(workspace)
	if want := "restored " + id + "/ from commit " + *pushed.Commit + " into " + config.LedgerDir(workspace) + "/\n"; result.err != nil || result.stderr != "" || result.stdout != want {
		t.Fatalf("text restore = %#v, want %q", result, want)
	}
	if listed := executeCommandForTest(workspace, "--json", "issue", "list"); listed.err != nil {
		t.Fatalf("issue list after the restore = %#v", listed)
	}
}

// Restore imports the backup under the idempotency key
// fid-restore:<commit>, which the restored ledger's import marker records.
func TestRestoreImportsUnderTheBackupCommitsKey(t *testing.T) {
	t.Parallel()
	fixture := newBackupFixture(t)
	original, id := fixture.project()
	fixture.setTarget(fixture.target)
	pushed, _ := decodeBackupSuccess(t, fixture.backup(original, "--json"))
	workspace := fixture.restoreWorkspace(id)

	restored, _ := decodeRestoreSuccess(t, fixture.restore(workspace, "--json"))
	binding, err := config.ReadImportBinding(workspace)
	if err != nil {
		t.Fatalf("read the restored ledger's import marker: %v", err)
	}
	if want := "fid-restore:" + *pushed.Commit; restored.Commit != *pushed.Commit || binding.IdempotencyKey != want {
		t.Fatalf("restore of %s imported with key %q, want %q", restored.Commit, binding.IdempotencyKey, want)
	}
}

// restoreWorkspace is a fresh clone of the project with ID id: a git
// repository holding only its .fid.toml, and no ledger.
func (fixture *backupFixture) restoreWorkspace(id string) string {
	fixture.t.Helper()
	workspace := gittest.NewWorkspace(fixture.t)
	writeBackupTestFile(fixture.t, filepath.Join(workspace, config.ProjectIdentityFileName), "[project]\nid = \""+id+"\"\n")
	return workspace
}

// Every worktree shares the repository's ledger, so a restore from a linked
// worktree, found through --project, restores the shared ledger, which the
// main worktree then uses.
func TestRestoreFromALinkedWorktreeRestoresTheSharedLedger(t *testing.T) {
	t.Parallel()
	fixture := newBackupFixture(t)
	original, id := fixture.project()
	fixture.setTarget(fixture.target)
	runFIDCommand(t, original, "issue", "create", "--kind", "task", "--title", "Shared by every worktree", "--actor", "agent:developer")
	decodeBackupSuccess(t, fixture.backup(original, "--json"))
	issues := executeCommandForTest(original, "--json", "issue", "list")
	fixture.git("-C", original, "add", config.ProjectIdentityFileName)
	fixture.git("-C", original, "-c", "user.name=human", "-c", "user.email=human@example.test", "commit", "-q", "-m", "identity")
	linked := gittest.AddWorktree(t, original)
	if err := os.RemoveAll(config.LedgerDir(original)); err != nil {
		t.Fatal(err)
	}

	restored, ledgerPath := decodeRestoreSuccess(t, fixture.restore(t.TempDir(), "--project", linked, "--json"))
	if restored.ProjectID != id || ledgerPath != config.LedgerDir(original) || ledgerPath != config.LedgerDir(linked) {
		t.Fatalf("restore = %+v into %s, want %s in the shared ledger %s", restored, ledgerPath, id, config.LedgerDir(original))
	}
	if again := executeCommandForTest(original, "--json", "issue", "list"); again.err != nil || again.stdout != issues.stdout {
		t.Fatalf("issue list in the main worktree = %#v, want %s", again, issues.stdout)
	}
	refused := decodeRestoreFailure(t, fixture.restore(linked, "--json"), 4)
	if refused.Code != "project-already-initialized" {
		t.Fatalf("restore with the shared ledger present = %+v, want project-already-initialized", refused)
	}
}
