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
		"--actor", "agent:developer",
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

	addACStdout, addACStderr := runFIDCommand(
		t,
		workspace,
		"issue", "ac", "add", finishIssueID,
		"--text", "Activity lifecycle finish path is ready",
		"--ordinal", "1",
		"--by", "developer",
	)
	assertNoNestedLockFailure(t, addACStderr)
	finishACID := parseAddedAcceptanceCriterionID(t, addACStdout)

	assertIssueACVerifyUnavailable(t, workspace, finishIssueID, finishACID, "agent:developer")
	seedHistoricalACVerification(t, workspace, finishACID, "developer")

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
		"--actor", "agent:developer",
	)
	assertNoNestedLockFailure(t, peerCreateStderr)
	peerIssueID := parseCreatedIssueID(t, peerCreateStderr)

	_, linkStderr := runFIDCommand(
		t,
		workspace,
		"link", "add", finishIssueID, "blocks", peerIssueID, "--from-type", "issue", "--to-type", "issue",
		"--by", "developer",
	)
	assertNoNestedLockFailure(t, linkStderr)

	_, unclaimCreateStderr := runFIDCommand(
		t,
		workspace,
		"issue", "create",
		"--kind", "task",
		"--title", "Activity lifecycle unclaim path",
		"--actor", "agent:developer",
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
		"--actor", "agent:developer",
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

	concernID := "CON-v01e-activity"
	assertConcernMutationUnavailable(t, workspace, []string{"concern", "raise", "--on", concernIssueID, "--actor", "agent:developer", "--in-completion", firstCompletionID, "--summary", "Follow-up required for lifecycle parity"})

	_, secondCompleteStderr := runFIDCommand(
		t,
		workspace,
		"issue", "complete",
		concernIssueID,
		"--by", "developer",
	)
	assertNoNestedLockFailure(t, secondCompleteStderr)
	secondCompletionID := parseCompletionID(t, secondCompleteStderr)

	assertConcernMutationUnavailable(t, workspace, []string{"concern", "address", concernID, "--in-completion", secondCompletionID, "--actor", "human:reviewer"})

	activityStdout, activityStderr := runFIDCommand(t, workspace, "--json", "activity", "list")
	if activityStderr != "" {
		t.Fatalf("activity list stderr = %q, want empty", activityStderr)
	}

	rows := decodeActivityLifecycleRows(t, activityStdout)

	want := []activityLifecycleExpectation{
		{RefID: finishIssueID, Event: "issue-created"},
		{RefID: finishIssueID, Event: "status-changed"},
		{RefID: finishIssueID, Event: "ac-added"},
		{RefID: finishIssueID, Event: "gate-recorded"},
		{RefID: finishIssueID, Event: "status-changed"},
		{RefID: finishIssueID, Event: "task-completed"},
		{RefID: peerIssueID, Event: "issue-created"},
		{RefID: finishIssueID, Event: "link-added"},
		{RefID: unclaimIssueID, Event: "issue-created"},
		{RefID: unclaimIssueID, Event: "status-changed"},
		{RefID: unclaimIssueID, Event: "claim-released"},
		{RefID: concernIssueID, Event: "issue-created"},
		{RefID: concernIssueID, Event: "task-completed"},
		{RefID: concernIssueID, Event: "task-completed"},
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
		"issue-created",
		"status-changed",
		"ac-added",
		"gate-recorded",
		"status-changed",
		"task-completed",
		"link-added",
	})
	assertActivityEventsForIssue(t, workspace, peerIssueID, []string{
		"issue-created",
	})
	assertActivityEventsForIssue(t, workspace, unclaimIssueID, []string{
		"issue-created",
		"status-changed",
		"claim-released",
	})
	assertActivityEventsForIssue(t, workspace, concernIssueID, []string{
		"issue-created",
		"task-completed",
		"task-completed",
	})
	if !strings.Contains(activityStdout, `"event":"gate-recorded"`) {
		t.Fatalf("activity list stdout = %q, want gate-recorded", activityStdout)
	}
	for _, event := range []string{"ac-verified", "concern-raised", "concern-addressed"} {
		if strings.Contains(activityStdout, `"event":"`+event+`"`) {
			t.Fatalf("activity list stdout = %q, want no %s event after fail-closed mutation", activityStdout, event)
		}
	}
	if !strings.Contains(activityStdout, `"event":"claim-released"`) {
		t.Fatalf("activity list stdout = %q, want claim-released", activityStdout)
	}
	if !strings.Contains(activityStdout, `"event":"link-added"`) {
		t.Fatalf("activity list stdout = %q, want link-added", activityStdout)
	}
	if !strings.Contains(activityStdout, `"gate_run_id":"`) {
		t.Fatalf("activity list stdout = %q, want gate_run_id details", activityStdout)
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

	rebuildStdout, rebuildStderr := runFIDCommand(t, workspace, "rebuild")
	assertNoNestedLockFailure(t, rebuildStderr)
	if !strings.Contains(rebuildStdout, "rebuild complete") || rebuildStderr != "" {
		t.Fatalf("rebuild stdout/stderr = %q/%q, want stdout completion and empty stderr", rebuildStdout, rebuildStderr)
	}

	afterStdout, afterStderr := runFIDCommand(t, workspace, "--json", "activity", "list")
	if afterStderr != "" {
		t.Fatalf("activity list stderr after rebuild = %q, want empty", afterStderr)
	}
	if got := snapshotActivityLogRows(t, workspace); !reflect.DeepEqual(got, beforeRebuild) {
		t.Fatalf("activity_log after rebuild = %#v, want %#v", got, beforeRebuild)
	}

	afterRows := decodeActivityLifecycleRows(t, afterStdout)
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

	stdout, stderr := runFIDCommand(t, workspace, "--json", "activity", "list", "--on", "issue/"+issueID)
	if stderr != "" {
		t.Fatalf("activity list stderr for %s = %q, want empty", issueID, stderr)
	}

	rows := decodeActivityLifecycleRows(t, stdout)

	got := make([]string, 0, len(rows))
	for _, row := range rows {
		got = append(got, row.Event)
	}
	if !reflect.DeepEqual(countLifecycleEventNames(got), countLifecycleEventNames(want)) {
		t.Fatalf("activity events for %s = %#v, want %#v", issueID, got, want)
	}
}

func decodeActivityLifecycleRows(t *testing.T, output string) []activityLifecycleRow {
	t.Helper()

	envelope := decodeReadEnvelope(t, output, "activity.list")
	payload, err := json.Marshal(envelope.Result["items"])
	if err != nil {
		t.Fatalf("Marshal(activity list items) error = %v", err)
	}
	var rows []activityLifecycleRow
	if err := json.Unmarshal(payload, &rows); err != nil {
		t.Fatalf("Unmarshal(activity list items) error = %v", err)
	}
	return rows
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
