//go:build integration

package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/mbell-nzl/fid/internal/config"
	fiddb "github.com/mbell-nzl/fid/internal/db"
	"github.com/mbell-nzl/fid/internal/maintenance"
)

func closureRun(t *testing.T, workspace, command string, args ...string) (string, string) {
	t.Helper()

	stdout, stderr, err := closureRunCommand(workspace, command, args...)
	if err != nil {
		t.Fatalf("%s %s: %v\nstdout:\n%s\nstderr:\n%s", command, strings.Join(args, " "), err, stdout, stderr)
	}
	return stdout, stderr
}

func closureRunFID(t *testing.T, binary, workspace string, args ...string) (string, string) {
	t.Helper()
	args = withClosureIssueMutationKey(args)
	args = withTask9IntegrationMutationRequest(args)
	return closureRun(t, workspace, binary, args...)
}

func closureRunFIDExpectError(t *testing.T, binary, workspace string, args ...string) (string, string, error) {
	t.Helper()
	args = withClosureIssueMutationKey(args)
	args = withTask9IntegrationMutationRequest(args)
	return closureRunCommand(workspace, binary, args...)
}

var closureIssueMutationSequence uint64
var task9IntegrationMutationSequence uint64

func withClosureIssueMutationKey(args []string) []string {
	if !closureTypedIssueMutation(args) || closureContainsArgument(args, "--idempotency-key") {
		return args
	}
	return append(args, "--idempotency-key", fmt.Sprintf("task8-integration-%d", atomic.AddUint64(&closureIssueMutationSequence, 1)))
}

func closureTypedIssueMutation(args []string) bool {
	for index := 0; index+1 < len(args); index++ {
		if args[index] == "issue" && (args[index+1] == "create" || args[index+1] == "update" || args[index+1] == "transition") {
			return closureContainsArgument(args, "--actor")
		}
	}
	return false
}

func closureContainsArgument(args []string, target string) bool {
	for _, argument := range args {
		if argument == target {
			return true
		}
	}
	return false
}

func withTask9IntegrationMutationRequest(args []string) []string {
	if !isTask9IntegrationMutation(args) || closureContainsArgument(args, "--actor") {
		return args
	}
	for index := 0; index+1 < len(args); index++ {
		if args[index] != "--by" {
			continue
		}
		converted := append([]string{}, args...)
		converted[index] = "--actor"
		converted[index+1] = "agent:" + args[index+1]
		return append(converted, "--idempotency-key", fmt.Sprintf("task9-integration-%d", atomic.AddUint64(&task9IntegrationMutationSequence, 1)))
	}
	return args
}

func isTask9IntegrationMutation(args []string) bool {
	for index := range args {
		if index+2 < len(args) && args[index] == "issue" && args[index+1] == "ac" && args[index+2] == "add" {
			return true
		}
		if index+2 < len(args) && args[index] == "issue" && args[index+1] == "label" && (args[index+2] == "add" || args[index+2] == "remove") {
			return true
		}
		if index+1 < len(args) && args[index] == "label" && args[index+1] == "add" {
			return true
		}
		if index+1 < len(args) && args[index] == "link" && (args[index+1] == "add" || args[index+1] == "remove") {
			return true
		}
		if index+1 < len(args) && args[index] == "requirement" && args[index+1] == "sync" {
			return true
		}
		if index+1 < len(args) && args[index] == "workflow" && (args[index+1] == "define" || args[index+1] == "assign") {
			return true
		}
	}
	return false
}

func TestTask9IntegrationMutationRequestsUseActorAndIdempotency(t *testing.T) {
	t.Parallel()
	for _, args := range [][]string{
		{"workflow", "define", "--by", "tl"},
		{"workflow", "assign", "--by", "tl"},
		{"issue", "ac", "add", "FID-1", "--by", "developer"},
		{"issue", "label", "add", "FID-1", "core", "--by", "developer"},
		{"label", "add", "core", "--by", "developer"},
		{"link", "add", "FID-1", "relates_to", "FID-2", "--by", "developer"},
		{"requirement", "sync", "--source-doc", "spec.md", "--by", "tl"},
	} {
		converted := withTask9IntegrationMutationRequest(args)
		if closureContainsArgument(converted, "--by") || !closureContainsArgument(converted, "--actor") || !closureContainsArgument(converted, "--idempotency-key") {
			t.Fatalf("task-9 integration request %v = %v, want actor/key transport without --by", args, converted)
		}
	}

	for _, args := range [][]string{
		{"issue", "create", "--by", "developer"},
		{"issue", "claim", "FID-1", "--by", "developer"},
		{"issue", "ac", "verify", "FID-1", "AC-1", "--by", "qa_lead"},
		{"concern", "raise", "--by", "developer"},
		{"run", "record", "--by", "developer"},
	} {
		converted := withTask9IntegrationMutationRequest(args)
		if strings.Join(converted, "\x00") != strings.Join(args, "\x00") {
			t.Fatalf("non-task-9 integration request %v = %v, want unchanged", args, converted)
		}
	}
}

func closureRunCommand(workspace, command string, args ...string) (string, string, error) {
	cmd := exec.Command(command, args...)
	cmd.Dir = workspace
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	p24CaptureKeyedCommand(workspace, args, stdout.String(), stderr.String(), err)
	return stdout.String(), stderr.String(), err
}

func stringPointer(value string) *string {
	return &value
}

func intPointer(value int) *int {
	return &value
}

func assertClosureGitLayout(t *testing.T, workspace string) {
	t.Helper()

	payload := closureReadFile(t, filepath.Join(config.LedgerDir(workspace), ".gitignore"))
	entries := strings.Fields(string(payload))
	if len(entries) != 1 || entries[0] != "fid.sqlite" {
		t.Fatalf("ledger .gitignore = %q, want only fid.sqlite", payload)
	}
	_, _ = closureRun(t, workspace, "git", "check-ignore", "-q", workspaceRelative(t, workspace, filepath.Join(config.LedgerDir(workspace), "fid.sqlite")))
}

// workspaceRelative names path relative to workspace, slash-separated, for Git
// pathspecs and shell commands that run in the workspace.
func workspaceRelative(t *testing.T, workspace, path string) string {
	t.Helper()

	relative, err := filepath.Rel(workspace, path)
	if err != nil {
		t.Fatalf("relative path of %q: %v", path, err)
	}
	return filepath.ToSlash(relative)
}

func closureJournalPath(t *testing.T, workspace string) string {
	t.Helper()

	matches, err := filepath.Glob(filepath.Join(config.LedgerDir(workspace), "journal", "*.jsonl"))
	if err != nil {
		t.Fatalf("find journal: %v", err)
	}
	if len(matches) != 1 {
		t.Fatalf("journal files = %v, want one", matches)
	}
	return matches[0]
}

// assertClosureJournalNeverTracked checks that git never tracks the journal.
// The ledger is in the git common directory (ADR-0013), so adding the
// journal, or everything, stages nothing, and status never lists the ledger.
func assertClosureJournalNeverTracked(t *testing.T, workspace string) {
	t.Helper()

	relative := workspaceRelative(t, workspace, closureJournalPath(t, workspace))
	closureRun(t, workspace, "git", "add", "--", relative)
	closureRun(t, workspace, "git", "add", "--all")
	if stdout, _, err := closureRunCommand(workspace, "git", "ls-files", "--error-unmatch", "--", relative); err == nil {
		t.Fatalf("journal %q is tracked: %q", relative, stdout)
	}
	ledger := workspaceRelative(t, workspace, config.LedgerDir(workspace))
	tracked, _ := closureRun(t, workspace, "git", "ls-files", "--cached")
	status, _ := closureRun(t, workspace, "git", "status", "--porcelain", "--ignored", "--untracked-files=all")
	if strings.Contains(tracked, ledger) || strings.Contains(status, ledger) {
		t.Fatalf("git lists the ledger %q: tracked %q, status %q", ledger, tracked, status)
	}
}

func closureWriteRequirementFixture(t *testing.T, workspace string) {
	t.Helper()

	path := filepath.Join(workspace, "project-docs", "specs", "closure.md")
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("create requirement fixture directory: %v", err)
	}
	payload := "# Requirements\n\n| ID | Requirement | Priority | Notes |\n|---|---|---|---|\n| PRD-LINK-002 | Closure journey requirement link. | Must | v0.1 |\n"
	if err := os.WriteFile(path, []byte(payload), 0o644); err != nil {
		t.Fatalf("write requirement fixture: %v", err)
	}
}

func closureCreateIssue(t *testing.T, binary, workspace, kind, title string, scheduling, epicID *string, sequence *int, extras ...string) string {
	t.Helper()

	args := []string{"issue", "create", "--kind", kind, "--title", title}
	if scheduling != nil {
		args = append(args, "--scheduling", *scheduling)
	}
	if epicID != nil {
		flag := "--epic"
		if kind == "subtask" {
			flag = "--parent"
		}
		args = append(args, flag, *epicID)
	}
	if sequence != nil {
		args = append(args, "--sequence", fmt.Sprintf("%d", *sequence))
	}
	args = append(args, extras...)
	actor := "developer"
	if kind == "epic" {
		actor = "tl"
	}
	args = append(args, "--actor", "agent:"+actor)
	stdout, stderr := closureRunFID(t, binary, workspace, args...)
	if stderr != "" {
		t.Fatalf("issue create stderr = %q, want empty", stderr)
	}
	return closureIDAfterWord(t, stdout, "issue")
}

func closureAddAcceptanceCriterion(t *testing.T, binary, workspace, issueID string) string {
	t.Helper()
	stdout, stderr := closureRunFID(t, binary, workspace, "issue", "ac", "add", issueID, "--text", "Closure acceptance criterion", "--ordinal", "1", "--by", "developer")
	if stderr != "" {
		t.Fatalf("issue ac add stderr = %q, want empty", stderr)
	}
	return closureIDAfterWord(t, stdout, "AC")
}

func closureCompleteIssue(t *testing.T, binary, workspace, issueID string) string {
	t.Helper()
	_, stderr := closureRunFID(t, binary, workspace, "issue", "complete", issueID, "--by", "developer")
	for _, line := range strings.Split(stderr, "\n") {
		if strings.HasPrefix(strings.TrimSpace(line), "completion: ") {
			return strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(line), "completion: "))
		}
	}
	t.Fatalf("completion stderr = %q, want completion ID", stderr)
	return ""
}

func assertClosureConcernMutationsUnavailable(t *testing.T, binary, workspace, issueID, completionID string) {
	t.Helper()
	for _, leaf := range []string{"raise", "address"} {
		assertClosureUnknownCommandWithoutJournalAppend(t, binary, workspace, closureJournalPath(t, workspace),
			[]string{"concern", leaf})
	}
	_, _ = issueID, completionID
}

func closureIDAfterWord(t *testing.T, text, word string) string {
	t.Helper()
	fields := strings.Fields(text)
	for index := range fields {
		if fields[index] == word && index+1 < len(fields) {
			return strings.TrimSpace(fields[index+1])
		}
	}
	t.Fatalf("output = %q, want %s ID", text, word)
	return ""
}

func assertClosureRejectedWithoutJournalAppend(t *testing.T, binary, workspace, journalPath string, args []string, want string) {
	t.Helper()

	before := closureReadFile(t, journalPath)
	stdout, stderr, err := closureRunFIDExpectError(t, binary, workspace, args...)
	if err == nil {
		t.Fatalf("fid %s succeeded, want rejection", strings.Join(args, " "))
	}
	if stdout != "" {
		t.Fatalf("fid %s stdout = %q, want empty", strings.Join(args, " "), stdout)
	}
	if !strings.Contains(stderr, want) {
		t.Fatalf("fid %s stderr = %q, want %q", strings.Join(args, " "), stderr, want)
	}
	assertClosureJournalBytesEqual(t, journalPath, before, strings.Join(args, " "))
}

func assertClosureUnknownCommandWithoutJournalAppend(t *testing.T, binary, workspace, journalPath string, args []string) {
	t.Helper()

	before := closureReadFile(t, journalPath)
	stdout, stderr, err := closureRunFIDExpectError(t, binary, workspace, args...)
	if err == nil {
		t.Fatalf("fid %s succeeded, want rejection", strings.Join(args, " "))
	}
	assertRemovedCommandFailure(t, args, stdout, stderr, err)
	jsonArgs := append([]string{"--json"}, args...)
	jsonOut, jsonErr, jsonCause := closureRunFIDExpectError(t, binary, workspace, jsonArgs...)
	assertRemovedCommandFailure(t, jsonArgs, jsonOut, jsonErr, jsonCause)
	assertClosureJournalBytesEqual(t, journalPath, before, strings.Join(args, " "))
}

func assertClosureLinkReadback(t *testing.T, binary, workspace, mainIssueID, peerIssueID string) {
	t.Helper()

	mainLinks, _ := closureRunFID(t, binary, workspace, "--json", "link", "list", mainIssueID, "--type", "issue")
	for _, want := range []string{peerIssueID, "discovered_from"} {
		if !strings.Contains(mainLinks, want) {
			t.Fatalf("main link list = %q, want %q", mainLinks, want)
		}
	}
	if strings.Contains(mainLinks, "PRD-LINK-002") || strings.Contains(mainLinks, "implements") {
		t.Fatalf("removed requirement link appears in list: %q", mainLinks)
	}
	peerLinks, _ := closureRunFID(t, binary, workspace, "--json", "link", "list", peerIssueID, "--type", "issue")
	if !strings.Contains(peerLinks, mainIssueID) || !strings.Contains(peerLinks, "discovered_from") {
		t.Fatalf("peer link list = %q, want inverse to %q", peerLinks, mainIssueID)
	}
}

func assertClosureReadyTask(t *testing.T, binary, workspace, epicID, wantID string) {
	t.Helper()

	stdout, _ := closureRunFID(t, binary, workspace, "--json", "ready", "--epic", epicID)
	var payload struct {
		Result struct {
			Task struct {
				ID string `json:"id"`
			} `json:"task"`
		} `json:"result"`
	}
	if err := json.Unmarshal([]byte(stdout), &payload); err != nil {
		t.Fatalf("decode ready output: %v", err)
	}
	if payload.Result.Task.ID != wantID {
		t.Fatalf("ready task for %s = %q, want %q", epicID, payload.Result.Task.ID, wantID)
	}
}

func assertClosureMaintenanceAndReadCommands(t *testing.T, binary, workspace, issueID string, journalBefore []byte) {
	t.Helper()

	check, _ := closureRunFID(t, binary, workspace, "--json", "maintenance", "check")
	if !strings.Contains(check, `"status":"pass"`) || strings.Contains(check, `"check":"telemetry_gaps"`) {
		t.Fatalf("maintenance check = %q, want no retired telemetry warning", check)
	}
	assertClosureJournalBytesEqual(t, closureJournalPath(t, workspace), journalBefore, "maintenance check")
	stdout, stderr, err := closureRunFIDExpectError(t, binary, workspace, "maintenance", "vacuum")
	if err == nil || stdout != "" || !strings.Contains(stderr, "capability is unavailable") {
		t.Fatalf("maintenance vacuum = (%v, %q, %q), want fail-closed unavailable", err, stdout, stderr)
	}
	assertClosureJournalBytesEqual(t, closureJournalPath(t, workspace), journalBefore, "maintenance vacuum")
	decision, _ := closureRunFID(t, binary, workspace, "--json", "decision", "list")
	if !strings.Contains(decision, `"items":[]`) {
		t.Fatalf("decision list = %q, want empty read result", decision)
	}
	assertClosureJournalBytesEqual(t, closureJournalPath(t, workspace), journalBefore, "decision list")
	for _, args := range [][]string{
		{"--json", "issue", "list"},
		{"--json", "ready"},
		{"--json", "link", "list", issueID, "--type", "issue"},
		{"stats", "cost", "--group-by", "actor"},
		{"stats", "tokens", "--group-by", "actor"},
	} {
		closureRunFID(t, binary, workspace, args...)
		assertClosureJournalBytesEqual(t, closureJournalPath(t, workspace), journalBefore, strings.Join(args, " "))
	}
	assertClosureUnknownCommandWithoutJournalAppend(t, binary, workspace, closureJournalPath(t, workspace), []string{"policy", "show"})
	assertClosureJournalBytesEqual(t, closureJournalPath(t, workspace), journalBefore, "policy show")
}

func closureAgedMaintenanceTime(t *testing.T, workspace string) time.Time {
	t.Helper()

	var latest time.Time
	for path, payload := range closureJournalFileBytes(t, workspace) {
		for lineNumber, line := range bytes.Split(bytes.TrimSpace(payload), []byte{'\n'}) {
			var event struct {
				OccurredAt time.Time `json:"occurred_at"`
			}
			if err := json.Unmarshal(line, &event); err != nil {
				t.Fatalf("decode %s:%d for maintenance time: %v", path, lineNumber+1, err)
			}
			if event.OccurredAt.After(latest) {
				latest = event.OccurredAt
			}
		}
	}
	if latest.IsZero() {
		t.Fatal("latest journal event time is zero")
	}
	return latest.UTC().Add(24 * time.Hour)
}

func closureMaintenanceResultAt(t *testing.T, workspace string, now time.Time) maintenance.Result {
	t.Helper()

	conn, err := fiddb.Open(filepath.Join(config.LedgerDir(workspace), "fid.sqlite"))
	if err != nil {
		t.Fatalf("open projection for maintenance check: %v", err)
	}
	defer func() { _ = conn.Close() }()

	result, err := maintenance.Check(maintenance.CheckOptions{
		DB:                              conn,
		JournalDir:                      filepath.Join(config.LedgerDir(workspace), "journal"),
		Now:                             now,
		StartedWithoutClaimGraceSeconds: 300,
	})
	if err != nil {
		t.Fatalf("maintenance check at %s: %v", now, err)
	}
	return result
}

func assertClosureMaintenanceStateEqual(t *testing.T, got, want maintenance.Result) {
	t.Helper()

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("maintenance result after rebuild = %#v, want %#v", got, want)
	}
}

var closureProjectedTableQueries = []struct {
	name  string
	query string
}{
	{
		name:  "acceptance_criteria",
		query: "SELECT * FROM acceptance_criteria ORDER BY id",
	},
	{
		name:  "activity_log",
		query: "SELECT * FROM activity_log ORDER BY occurred_at, id",
	},
	{
		name:  "agent_runs",
		query: "SELECT * FROM agent_runs ORDER BY id",
	},
	{
		name:  "claims",
		query: "SELECT * FROM claims ORDER BY id",
	},
	{
		name:  "concerns",
		query: "SELECT * FROM concerns ORDER BY id",
	},
	{
		name:  "decision_refs",
		query: "SELECT * FROM decision_refs ORDER BY id",
	},
	{
		name:  "epic_metadata",
		query: "SELECT * FROM epic_metadata ORDER BY issue_id",
	},
	{
		name:  "issue_labels",
		query: "SELECT * FROM issue_labels ORDER BY issue_id, label",
	},
	{
		name:  "issues",
		query: "SELECT * FROM issues ORDER BY id",
	},
	{
		name:  "labels",
		query: "SELECT * FROM labels ORDER BY name",
	},
	{
		name:  "links",
		query: "SELECT * FROM links ORDER BY from_type, from_id, relation, to_type, to_id",
	},
	{
		name:  "projection_meta",
		query: "SELECT * FROM projection_meta",
	},
	{
		name:  "requirement_refs",
		query: "SELECT * FROM requirement_refs ORDER BY id",
	},
	{
		name:  "requirement_issue_links",
		query: "SELECT * FROM requirement_issue_links ORDER BY requirement_id, issue_id, relation",
	},
	{
		name:  "review_records",
		query: "SELECT * FROM review_records ORDER BY id",
	},
	{
		name:  "triage_events",
		query: "SELECT * FROM triage_events ORDER BY id",
	},
	{
		name:  "workflows",
		query: "SELECT * FROM workflows ORDER BY id",
	},
	{
		name:  "workflow_states",
		query: "SELECT * FROM workflow_states ORDER BY workflow_id, ordinal, state",
	},
	{
		name:  "workflow_transitions",
		query: "SELECT * FROM workflow_transitions ORDER BY workflow_id, from_state, to_state",
	},
	{
		name:  "workflow_transition_requires",
		query: "SELECT * FROM workflow_transition_requires ORDER BY workflow_id, from_state, to_state, rowid",
	},
	{
		name:  "workflow_assignments",
		query: "SELECT * FROM workflow_assignments ORDER BY scope, kind, workflow_id",
	},
}

func closureProjectedTableSnapshots(t *testing.T, binary, workspace string) map[string]string {
	t.Helper()

	snapshot := make(map[string]string, len(closureProjectedTableQueries))
	for _, table := range closureProjectedTableQueries {
		stdout, _ := closureRunFID(t, binary, workspace, "sql", table.query)
		snapshot[table.name] = stdout
	}
	return snapshot
}

func assertClosureProjectedTablesEqual(t *testing.T, got, want map[string]string) {
	t.Helper()

	for _, table := range closureProjectedTableQueries {
		if got[table.name] != want[table.name] {
			t.Fatalf("projected table %s after rebuild = %s, want %s", table.name, got[table.name], want[table.name])
		}
	}
}

func assertClosureWorkflowsAfterRebuild(t *testing.T, binary, workspace string) {
	t.Helper()

	stdout, _ := closureRunFID(t, binary, workspace, "sql", "SELECT id FROM workflows WHERE active = 1 ORDER BY id")
	for _, workflowID := range []string{"epic-default", "subtask-default", "task-default"} {
		if !strings.Contains(stdout, `"id":"`+workflowID+`"`) {
			t.Fatalf("workflows after rebuild = %q, want %q", stdout, workflowID)
		}
	}
}

func closureReadFile(t *testing.T, path string) []byte {
	t.Helper()
	payload, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return payload
}

func assertClosureJournalBytesEqual(t *testing.T, journalPath string, want []byte, step string) {
	t.Helper()
	if got := closureReadFile(t, journalPath); !bytes.Equal(got, want) {
		t.Fatalf("journal bytes changed after %s", step)
	}
}

func closureJournalFileBytes(t *testing.T, workspace string) map[string][]byte {
	t.Helper()
	paths, err := filepath.Glob(filepath.Join(config.LedgerDir(workspace), "journal", "*.jsonl"))
	if err != nil {
		t.Fatalf("find journal files: %v", err)
	}
	files := make(map[string][]byte, len(paths))
	for _, path := range paths {
		files[filepath.Base(path)] = closureReadFile(t, path)
	}
	return files
}

func assertClosureJournalFilesEqual(t *testing.T, got, want map[string][]byte, step string) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("journal files changed after %s", step)
	}
}

func closureSchemaSnapshot(t *testing.T, binary, workspace string) string {
	t.Helper()
	stdout, _ := closureRunFID(t, binary, workspace, "sql", `SELECT type, name, tbl_name, sql FROM sqlite_schema WHERE name NOT LIKE 'sqlite_%' ORDER BY type, name`)
	return stdout
}

func assertClosureForeignKeysClean(t *testing.T, binary, workspace string) {
	t.Helper()
	stdout, _ := closureRunFID(t, binary, workspace, "sql", "PRAGMA foreign_key_check")
	if strings.TrimSpace(stdout) != "[]" {
		t.Fatalf("foreign_key_check = %s, want []", stdout)
	}
}
