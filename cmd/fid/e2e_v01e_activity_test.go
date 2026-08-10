package main

import (
	"database/sql"
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	fiddb "github.com/mbell-nzl/fid/internal/db"
)

func TestV01EActivityLifecycleCoverageAndRebuildParity(t *testing.T) {
	workspace := t.TempDir()
	runInitCommand(t, workspace, "init")

	finishReportPath := writeV01EMarkdown(t, workspace, "reports/activity-finish.md", "# finish\n")

	_, finishCreateStderr := runFIDCommand(
		t,
		workspace,
		"issue", "create",
		"--kind", "task",
		"--title", "Activity lifecycle finish path",
		"--by", "developer",
	)
	assertNoNestedLockFailure(t, finishCreateStderr)
	finishIssueID := parseCreatedIssueID(t, finishCreateStderr)

	_, startFinishStderr := runFIDCommand(
		t,
		workspace,
		"dev", "start", finishIssueID,
		"--model", "claude-sonnet-4-6",
		"--by", "developer",
	)
	assertNoNestedLockFailure(t, startFinishStderr)

	_, addACStderr := runFIDCommand(
		t,
		workspace,
		"issue", "ac", "add", finishIssueID,
		"--text", "Activity lifecycle finish path is ready",
		"--ordinal", "1",
		"--by", "developer",
	)
	assertNoNestedLockFailure(t, addACStderr)
	finishACID := parseAddedAcceptanceCriterionID(t, addACStderr)

	_, verifyACStderr := runFIDCommand(
		t,
		workspace,
		"issue", "ac", "verify", finishIssueID, finishACID,
		"--status", "pass",
		"--by", "developer",
	)
	assertNoNestedLockFailure(t, verifyACStderr)

	_, gateStderr := runFIDCommand(
		t,
		workspace,
		"dev", "gate", finishIssueID,
		"--cmd", "go test ./...",
		"--exit-code", "0",
		"--output-excerpt", "ok",
		"--by", "developer",
	)
	assertNoNestedLockFailure(t, gateStderr)

	_, finishStderr := runFIDCommand(
		t,
		workspace,
		"dev", "finish", finishIssueID,
		"--completion-md", finishReportPath,
		"--tokens", "14823",
		"--commits", "abc123",
		"--files", "cmd/fid/e2e_v01e_activity_test.go",
		"--by", "developer",
	)
	assertNoNestedLockFailure(t, finishStderr)
	finishCompletionID := strings.Fields(parseCompletionID(t, finishStderr))[0]

	_, peerCreateStderr := runFIDCommand(
		t,
		workspace,
		"issue", "create",
		"--kind", "task",
		"--title", "Activity lifecycle peer issue",
		"--by", "developer",
	)
	assertNoNestedLockFailure(t, peerCreateStderr)
	peerIssueID := parseCreatedIssueID(t, peerCreateStderr)

	_, linkStderr := runFIDCommand(
		t,
		workspace,
		"link", "add", finishIssueID, "blocks", peerIssueID,
		"--by", "developer",
	)
	assertNoNestedLockFailure(t, linkStderr)

	_, unclaimCreateStderr := runFIDCommand(
		t,
		workspace,
		"issue", "create",
		"--kind", "task",
		"--title", "Activity lifecycle unclaim path",
		"--by", "developer",
	)
	assertNoNestedLockFailure(t, unclaimCreateStderr)
	unclaimIssueID := parseCreatedIssueID(t, unclaimCreateStderr)

	_, startUnclaimStderr := runFIDCommand(
		t,
		workspace,
		"dev", "start", unclaimIssueID,
		"--model", "claude-sonnet-4-6",
		"--by", "developer",
	)
	assertNoNestedLockFailure(t, startUnclaimStderr)

	_, unclaimStderr := runFIDCommand(
		t,
		workspace,
		"issue", "unclaim", unclaimIssueID,
		"--by", "developer",
	)
	assertNoNestedLockFailure(t, unclaimStderr)

	_, concernCreateStderr := runFIDCommand(
		t,
		workspace,
		"issue", "create",
		"--kind", "task",
		"--title", "Activity lifecycle concern path",
		"--by", "developer",
	)
	assertNoNestedLockFailure(t, concernCreateStderr)
	concernIssueID := parseCreatedIssueID(t, concernCreateStderr)

	_, firstCompleteStderr := runFIDCommand(
		t,
		workspace,
		"issue", "complete",
		concernIssueID,
		"--by", "developer",
	)
	assertNoNestedLockFailure(t, firstCompleteStderr)
	firstCompletionID := parseCompletionID(t, firstCompleteStderr)

	_, raiseConcernStderr := runFIDCommand(
		t,
		workspace,
		"concern", "raise",
		"--on", concernIssueID,
		"--by", "developer",
		"--in-completion", firstCompletionID,
		"--summary", "Follow-up required for lifecycle parity",
	)
	assertNoNestedLockFailure(t, raiseConcernStderr)
	concernID := parseConcernID(t, raiseConcernStderr)

	_, secondCompleteStderr := runFIDCommand(
		t,
		workspace,
		"issue", "complete",
		concernIssueID,
		"--by", "developer",
	)
	assertNoNestedLockFailure(t, secondCompleteStderr)
	secondCompletionID := parseCompletionID(t, secondCompleteStderr)

	_, addressConcernStderr := runFIDCommand(
		t,
		workspace,
		"concern", "address", concernID,
		"--in-completion", secondCompletionID,
		"--by", "developer",
	)
	assertNoNestedLockFailure(t, addressConcernStderr)

	activityStdout, activityStderr := runFIDCommand(t, workspace, "activity", "list", "--format", "json")
	if activityStderr != "" {
		t.Fatalf("activity list stderr = %q, want empty", activityStderr)
	}

	var rows []activityLifecycleRow
	if err := json.Unmarshal([]byte(activityStdout), &rows); err != nil {
		t.Fatalf("Unmarshal(activity list) error = %v", err)
	}

	want := []activityLifecycleExpectation{
		{RefID: finishIssueID, Event: "issue_created"},
		{RefID: finishIssueID, Event: "status_changed"},
		{RefID: finishIssueID, Event: "ac_added"},
		{RefID: finishIssueID, Event: "ac_verified"},
		{RefID: finishIssueID, Event: "gate_recorded"},
		{RefID: finishIssueID, Event: "status_changed"},
		{RefID: finishIssueID, Event: "task_completed"},
		{RefID: peerIssueID, Event: "issue_created"},
		{RefID: finishIssueID, Event: "link_added"},
		{RefID: unclaimIssueID, Event: "issue_created"},
		{RefID: unclaimIssueID, Event: "status_changed"},
		{RefID: unclaimIssueID, Event: "claim_released"},
		{RefID: concernIssueID, Event: "issue_created"},
		{RefID: concernIssueID, Event: "task_completed"},
		{RefID: concernIssueID, Event: "concern_raised"},
		{RefID: concernIssueID, Event: "task_completed"},
		{RefID: concernIssueID, Event: "concern_addressed"},
	}
	if len(rows) != len(want) {
		t.Fatalf("activity row count = %d, want %d", len(rows), len(want))
	}
	if got, wantCounts := countLifecycleEvents(rows), countLifecycleExpectations(want); !reflect.DeepEqual(got, wantCounts) {
		t.Fatalf("activity event counts = %#v, want %#v", got, wantCounts)
	}
	for index := range rows {
		if rows[index].RefType != "issue" {
			t.Fatalf("activity row %d ref_type = %q, want issue", index, rows[index].RefType)
		}
		if rows[index].JournalEventID == "" || !strings.HasPrefix(rows[index].JournalEventID, "EVT-") {
			t.Fatalf("activity row %d journal_event_id = %q, want EVT- prefix", index, rows[index].JournalEventID)
		}
	}
	assertActivityEventsForIssue(t, workspace, finishIssueID, []string{
		"issue_created",
		"status_changed",
		"ac_added",
		"ac_verified",
		"gate_recorded",
		"status_changed",
		"task_completed",
		"link_added",
	})
	assertActivityEventsForIssue(t, workspace, peerIssueID, []string{
		"issue_created",
	})
	assertActivityEventsForIssue(t, workspace, unclaimIssueID, []string{
		"issue_created",
		"status_changed",
		"claim_released",
	})
	assertActivityEventsForIssue(t, workspace, concernIssueID, []string{
		"issue_created",
		"task_completed",
		"concern_raised",
		"task_completed",
		"concern_addressed",
	})
	if !strings.Contains(activityStdout, `"event":"gate_recorded"`) {
		t.Fatalf("activity list stdout = %q, want gate_recorded", activityStdout)
	}
	if !strings.Contains(activityStdout, `"event":"concern_addressed"`) {
		t.Fatalf("activity list stdout = %q, want concern_addressed", activityStdout)
	}
	if !strings.Contains(activityStdout, `"event":"claim_released"`) {
		t.Fatalf("activity list stdout = %q, want claim_released", activityStdout)
	}
	if !strings.Contains(activityStdout, `"event":"link_added"`) {
		t.Fatalf("activity list stdout = %q, want link_added", activityStdout)
	}
	if !strings.Contains(activityStdout, `"gate_run_id":"`) {
		t.Fatalf("activity list stdout = %q, want gate_run_id details", activityStdout)
	}
	if !strings.Contains(activityStdout, `"concern_id":"`+concernID+`"`) {
		t.Fatalf("activity list stdout = %q, want concern_id %q", activityStdout, concernID)
	}
	if !strings.Contains(activityStdout, `"attempt_number":1`) {
		t.Fatalf("activity list stdout = %q, want attempt_number details", activityStdout)
	}
	if !strings.Contains(activityStdout, `"attempt_number":2`) {
		t.Fatalf("activity list stdout = %q, want second completion attempt details", activityStdout)
	}
	if !strings.Contains(activityStdout, `"completion_id":"`+finishCompletionID+`"`) {
		t.Fatalf("activity list stdout = %q, want finish completion id %q", activityStdout, finishCompletionID)
	}

	beforeRebuild := snapshotActivityLogRows(t, workspace)

	dbPath := filepath.Join(workspace, ".fid", "fid.sqlite")
	if err := os.Remove(dbPath); err != nil {
		t.Fatalf("Remove(%q) error = %v", dbPath, err)
	}

	_, rebuildStderr := runFIDCommand(t, workspace, "rebuild", "--by", "workflow_runner")
	assertNoNestedLockFailure(t, rebuildStderr)
	if !strings.Contains(rebuildStderr, "rebuild complete") {
		t.Fatalf("rebuild stderr = %q, want completion output", rebuildStderr)
	}

	afterStdout, afterStderr := runFIDCommand(t, workspace, "activity", "list", "--format", "json")
	if afterStderr != "" {
		t.Fatalf("activity list stderr after rebuild = %q, want empty", afterStderr)
	}
	if got := snapshotActivityLogRows(t, workspace); !reflect.DeepEqual(got, beforeRebuild) {
		t.Fatalf("activity_log after rebuild = %#v, want %#v", got, beforeRebuild)
	}

	var afterRows []activityLifecycleRow
	if err := json.Unmarshal([]byte(afterStdout), &afterRows); err != nil {
		t.Fatalf("Unmarshal(activity list after rebuild) error = %v", err)
	}
	if !reflect.DeepEqual(afterRows, rows) {
		t.Fatalf("activity list after rebuild = %#v, want %#v", afterRows, rows)
	}
}

type activityLifecycleExpectation struct {
	RefID string
	Event string
}

type activityLifecycleRow struct {
	ID             string          `json:"id"`
	JournalEventID string          `json:"journal_event_id"`
	RefType        string          `json:"ref_type"`
	RefID          string          `json:"ref_id"`
	Event          string          `json:"event"`
	Actor          string          `json:"actor"`
	OccurredAt     string          `json:"occurred_at"`
	Details        json.RawMessage `json:"details"`
}

type activityLogSnapshotRow struct {
	ID             string
	JournalEventID string
	RefType        string
	RefID          string
	Event          string
	Actor          string
	OccurredAt     string
	Details        string
}

func snapshotActivityLogRows(t *testing.T, workspace string) []activityLogSnapshotRow {
	t.Helper()

	conn, err := fiddb.Open(filepath.Join(workspace, ".fid", "fid.sqlite"))
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	defer closeCommandDB(t, conn)

	rows, err := conn.Query(`
		SELECT id, journal_event_id, ref_type, ref_id, event, COALESCE(actor, ''), COALESCE(occurred_at, ''), COALESCE(details, '')
		FROM activity_log
		ORDER BY occurred_at ASC, id ASC
	`)
	if err != nil {
		t.Fatalf("query activity_log: %v", err)
	}
	defer closeActivityLifecycleRows(t, rows)

	snapshot := []activityLogSnapshotRow{}
	for rows.Next() {
		var row activityLogSnapshotRow
		if err := rows.Scan(
			&row.ID,
			&row.JournalEventID,
			&row.RefType,
			&row.RefID,
			&row.Event,
			&row.Actor,
			&row.OccurredAt,
			&row.Details,
		); err != nil {
			t.Fatalf("scan activity_log: %v", err)
		}
		snapshot = append(snapshot, row)
	}
	if err := rows.Err(); err != nil {
		t.Fatalf("iterate activity_log: %v", err)
	}
	return snapshot
}

func closeActivityLifecycleRows(t *testing.T, rows *sql.Rows) {
	t.Helper()
	if err := rows.Close(); err != nil {
		t.Fatalf("Close() rows error = %v", err)
	}
}

func assertActivityEventsForIssue(t *testing.T, workspace, issueID string, want []string) {
	t.Helper()

	stdout, stderr := runFIDCommand(t, workspace, "activity", "list", "--on", "issue/"+issueID, "--format", "json")
	if stderr != "" {
		t.Fatalf("activity list stderr for %s = %q, want empty", issueID, stderr)
	}

	var rows []activityLifecycleRow
	if err := json.Unmarshal([]byte(stdout), &rows); err != nil {
		t.Fatalf("Unmarshal(activity list for %s) error = %v", issueID, err)
	}

	got := make([]string, 0, len(rows))
	for _, row := range rows {
		got = append(got, row.Event)
	}
	if !reflect.DeepEqual(countLifecycleEventNames(got), countLifecycleEventNames(want)) {
		t.Fatalf("activity events for %s = %#v, want %#v", issueID, got, want)
	}
}

func countLifecycleEvents(rows []activityLifecycleRow) map[string]int {
	counts := map[string]int{}
	for _, row := range rows {
		counts[row.RefID+"|"+row.Event]++
	}
	return counts
}

func countLifecycleExpectations(rows []activityLifecycleExpectation) map[string]int {
	counts := map[string]int{}
	for _, row := range rows {
		counts[row.RefID+"|"+row.Event]++
	}
	return counts
}

func countLifecycleEventNames(events []string) map[string]int {
	counts := map[string]int{}
	for _, event := range events {
		counts[event]++
	}
	return counts
}
