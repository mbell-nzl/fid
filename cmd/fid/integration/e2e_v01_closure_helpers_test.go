//go:build integration

package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

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
	return closureRun(t, workspace, binary, args...)
}

func closureRunFIDExpectError(t *testing.T, binary, workspace string, args ...string) (string, string, error) {
	t.Helper()
	return closureRunCommand(workspace, binary, args...)
}

func closureRunCommand(workspace, command string, args ...string) (string, string, error) {
	cmd := exec.Command(command, args...)
	cmd.Dir = workspace
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
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

	payload := closureReadFile(t, filepath.Join(workspace, ".fid", ".gitignore"))
	entries := strings.Fields(string(payload))
	if len(entries) != 1 || entries[0] != "fid.sqlite" {
		t.Fatalf(".fid/.gitignore = %q, want only fid.sqlite", payload)
	}
	_, _ = closureRun(t, workspace, "git", "check-ignore", "-q", ".fid/fid.sqlite")
}

func closureJournalPath(t *testing.T, workspace string) string {
	t.Helper()

	matches, err := filepath.Glob(filepath.Join(workspace, ".fid", "journal", "*.jsonl"))
	if err != nil {
		t.Fatalf("find journal: %v", err)
	}
	if len(matches) != 1 {
		t.Fatalf("journal files = %v, want one", matches)
	}
	return matches[0]
}

func assertClosureJournalTracked(t *testing.T, workspace string) {
	t.Helper()

	journalPath := closureJournalPath(t, workspace)
	relative, err := filepath.Rel(workspace, journalPath)
	if err != nil {
		t.Fatalf("relative journal path: %v", err)
	}
	stdout, _ := closureRun(t, workspace, "git", "ls-files", "--error-unmatch", "--", filepath.ToSlash(relative))
	if strings.TrimSpace(stdout) != filepath.ToSlash(relative) {
		t.Fatalf("tracked journal = %q, want %q", stdout, filepath.ToSlash(relative))
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
		args = append(args, "--epic", *epicID)
	}
	if sequence != nil {
		args = append(args, "--sequence", fmt.Sprintf("%d", *sequence))
	}
	args = append(args, extras...)
	actor := "developer"
	if kind == "epic" {
		actor = "tl"
	}
	args = append(args, "--by", actor)
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

func closureRaiseAndAddressConcern(t *testing.T, binary, workspace, issueID, completionID string) (string, string) {
	t.Helper()
	_, stderr := closureRunFID(t, binary, workspace, "concern", "raise", "--on", issueID, "--by", "developer", "--in-completion", completionID, "--summary", "Closure concern")
	concernID := closureIDAfterWord(t, stderr, "concern")
	addressCompletionID := closureCompleteIssue(t, binary, workspace, issueID)
	closureRunFID(t, binary, workspace, "concern", "address", concernID, "--in-completion", addressCompletionID, "--by", "developer")
	return concernID, addressCompletionID
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

func assertClosureLinkReadback(t *testing.T, binary, workspace, mainIssueID, peerIssueID string) {
	t.Helper()

	mainLinks, _ := closureRunFID(t, binary, workspace, "--json", "link", "list", mainIssueID, "--type", "issue")
	for _, want := range []string{"PRD-LINK-002", peerIssueID, "implements", "discovered_from"} {
		if !strings.Contains(mainLinks, want) {
			t.Fatalf("main link list = %q, want %q", mainLinks, want)
		}
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

func assertClosureTelemetry(t *testing.T, binary, workspace, issueID, completionID string) {
	t.Helper()

	stdout, _ := closureRunFID(t, binary, workspace, "sql", "SELECT tokens_in, tokens_out, price_in_per_mtok, price_out_per_mtok, cost_usd, completion_id FROM agent_runs WHERE issue_id = '"+issueID+"' ORDER BY id")
	var rows []struct {
		TokensIn   int     `json:"tokens_in"`
		TokensOut  int     `json:"tokens_out"`
		PriceIn    float64 `json:"price_in_per_mtok"`
		PriceOut   float64 `json:"price_out_per_mtok"`
		CostUSD    float64 `json:"cost_usd"`
		Completion string  `json:"completion_id"`
	}
	if err := json.Unmarshal([]byte(stdout), &rows); err != nil {
		t.Fatalf("decode run snapshot: %v", err)
	}
	if len(rows) != 2 {
		t.Fatalf("run snapshots = %#v, want two completion-linked runs", rows)
	}
	var primary *struct {
		TokensIn   int     `json:"tokens_in"`
		TokensOut  int     `json:"tokens_out"`
		PriceIn    float64 `json:"price_in_per_mtok"`
		PriceOut   float64 `json:"price_out_per_mtok"`
		CostUSD    float64 `json:"cost_usd"`
		Completion string  `json:"completion_id"`
	}
	for index := range rows {
		if rows[index].Completion == completionID {
			primary = &rows[index]
		}
	}
	if primary == nil {
		t.Fatalf("run snapshots = %#v, want one linked to %q", rows, completionID)
	}
	wantCost := (float64(primary.TokensIn)*primary.PriceIn + float64(primary.TokensOut)*primary.PriceOut) / 1_000_000
	if primary.PriceIn <= 0 || primary.PriceOut <= 0 || math.Abs(primary.CostUSD-wantCost) > 1e-12 {
		t.Fatalf("run snapshot = %#v, want priced cost %v", primary, wantCost)
	}
	cost, _ := closureRunFID(t, binary, workspace, "stats", "cost", "--by", "role")
	tokens, _ := closureRunFID(t, binary, workspace, "stats", "tokens", "--by", "role")
	if !strings.Contains(cost, "developer") || !strings.Contains(cost, fmt.Sprintf("%g", primary.CostUSD)) {
		t.Fatalf("cost stats = %q, want developer cost %g", cost, primary.CostUSD)
	}
	if !strings.Contains(tokens, `"total_tokens_in":1200`) || !strings.Contains(tokens, `"total_tokens_out":3400`) {
		t.Fatalf("token stats = %q, want 1200/3400", tokens)
	}
}

func assertClosureMaintenanceAndReadCommands(t *testing.T, binary, workspace, issueID string, journalBefore []byte) {
	t.Helper()

	check, _ := closureRunFID(t, binary, workspace, "--json", "maintenance", "check")
	if !strings.Contains(check, `"status":"pass"`) {
		t.Fatalf("maintenance check = %q, want pass", check)
	}
	assertClosureJournalBytesEqual(t, closureJournalPath(t, workspace), journalBefore, "maintenance check")
	closureRunFID(t, binary, workspace, "maintenance", "vacuum", "--by", "tl")
	assertClosureJournalBytesEqual(t, closureJournalPath(t, workspace), journalBefore, "maintenance vacuum")
	stdout, stderr, err := closureRunFIDExpectError(t, binary, workspace, "decision", "get")
	if err == nil || stdout != "" || stderr != "" {
		t.Fatalf("decision command = (%v, %q, %q), want unknown command with empty output", err, stdout, stderr)
	}
	assertClosureJournalBytesEqual(t, closureJournalPath(t, workspace), journalBefore, "unknown decision command")
	for _, args := range [][]string{
		{"--json", "issue", "list"},
		{"--json", "ready"},
		{"--json", "link", "list", issueID, "--type", "issue"},
		{"stats", "cost", "--by", "role"},
		{"stats", "tokens", "--by", "role"},
		{"policy", "show"},
	} {
		closureRunFID(t, binary, workspace, args...)
		assertClosureJournalBytesEqual(t, closureJournalPath(t, workspace), journalBefore, strings.Join(args, " "))
	}
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

	conn, err := fiddb.Open(filepath.Join(workspace, ".fid", "fid.sqlite"))
	if err != nil {
		t.Fatalf("open projection for maintenance check: %v", err)
	}
	defer func() { _ = conn.Close() }()

	result, err := maintenance.Check(maintenance.CheckOptions{
		DB:                              conn,
		JournalDir:                      filepath.Join(workspace, ".fid", "journal"),
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
		name:  "ac_verifications",
		query: "SELECT * FROM ac_verifications ORDER BY id",
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
		query: "SELECT * FROM claims ORDER BY ref_type, ref_id, claimed_by",
	},
	{
		name:  "concerns",
		query: "SELECT * FROM concerns ORDER BY id",
	},
	{
		name:  "decisions",
		query: "SELECT * FROM decisions ORDER BY id",
	},
	{
		name:  "epic_metadata",
		query: "SELECT * FROM epic_metadata ORDER BY issue_id",
	},
	{
		name:  "escalations",
		query: "SELECT * FROM escalations ORDER BY id",
	},
	{
		name:  "gate_runs",
		query: "SELECT * FROM gate_runs ORDER BY id",
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
		name:  "requirements",
		query: "SELECT * FROM requirements ORDER BY id",
	},
	{
		name:  "roles",
		query: "SELECT * FROM roles ORDER BY name",
	},
	{
		name:  "task_completions",
		query: "SELECT * FROM task_completions ORDER BY id",
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
	paths, err := filepath.Glob(filepath.Join(workspace, ".fid", "journal", "*.jsonl"))
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
