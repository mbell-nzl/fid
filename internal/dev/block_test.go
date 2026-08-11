package dev

import (
	"database/sql"
	"errors"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mbell-nzl/fid/internal/issues"
)

func TestBlockRecordsFailedGateCompletionTransitionsAndReleasesClaim(t *testing.T) {
	t.Parallel()

	workspace := newDevWorkspace(t)
	issueID, runID, _ := seedFinishReadyIssue(t, workspace)
	if _, err := Gate(GateOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testDevMigrationsDir(t),
		IssueID:       issueID,
		ActorRole:     "developer",
		Command:       "golangci-lint run ./...",
		ExitCode:      1,
		OutputExcerpt: stringPointer("internal/dev/finish.go: line too long"),
	}); err != nil {
		t.Fatalf("Gate() error = %v", err)
	}

	journalPath := journalPathForDevWorkspace(t, workspace)
	linesBefore := countDevJournalLines(journalPath)

	result, err := Block(BlockOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testDevMigrationsDir(t),
		IssueID:       issueID,
		ActorRole:     "developer",
		ReasonMDPath:  "project-docs/docs/phases/phase-15/tasks/task-3-done.md",
	})
	if err != nil {
		t.Fatalf("Block() error = %v", err)
	}

	if result.IssueID != issueID {
		t.Fatalf("Block().IssueID = %q, want %q", result.IssueID, issueID)
	}
	if result.FromState != "in_progress" {
		t.Fatalf("Block().FromState = %q, want in_progress", result.FromState)
	}
	if result.ToState != "failed_gate" {
		t.Fatalf("Block().ToState = %q, want failed_gate", result.ToState)
	}
	if result.Completion.Status != "failed_gate" {
		t.Fatalf("Block().Completion.Status = %q, want failed_gate", result.Completion.Status)
	}

	if got := countDevJournalLines(journalPath); got != linesBefore+1 {
		t.Fatalf("journal lines after Block() = %d, want %d", got, linesBefore+1)
	}

	stored := readDevStoredEvents(t, journalPath)
	last := stored[len(stored)-1]
	if last.EventType != "dev.blocked" {
		t.Fatalf("last event type = %q, want dev.blocked", last.EventType)
	}
	assertDevReplayableRevisionEffect(t, last, issueID, 2, 3)

	var (
		statusChangedEffects int
		taskCompletedEffects int
	)
	for _, effect := range last.Effects {
		if effect.Table != "activity_log" || effect.Op != "activity" {
			continue
		}
		switch effect.Values["event"] {
		case "status_changed":
			statusChangedEffects++
		case "task_completed":
			taskCompletedEffects++
		}
	}
	if statusChangedEffects != 1 || taskCompletedEffects != 1 {
		t.Fatalf("block activity effects = status_changed:%d task_completed:%d, want 1 each", statusChangedEffects, taskCompletedEffects)
	}

	conn := openDevDB(t, filepath.Join(workspace, ".fid", "fid.sqlite"))
	defer closeDevDB(t, conn)

	rows, err := conn.Query(`
		SELECT event, COALESCE(details, '')
		FROM activity_log
		WHERE journal_event_id = ?
		ORDER BY event
	`, last.EventID)
	if err != nil {
		t.Fatalf("query block activity rows: %v", err)
	}
	defer closeDevRows(rows)

	activityDetails := map[string]string{}
	for rows.Next() {
		var eventName string
		var details string
		if err := rows.Scan(&eventName, &details); err != nil {
			t.Fatalf("scan block activity row: %v", err)
		}
		activityDetails[eventName] = details
	}
	if err := rows.Err(); err != nil {
		t.Fatalf("iterate block activity rows: %v", err)
	}
	if len(activityDetails) != 2 {
		t.Fatalf("block activity row count = %d, want 2", len(activityDetails))
	}
	if !strings.Contains(activityDetails["task_completed"], `"status":"failed_gate"`) {
		t.Fatalf("task_completed details = %q, want failed_gate status", activityDetails["task_completed"])
	}

	var status string
	if err := conn.QueryRow(`SELECT status FROM issues WHERE id = ?`, issueID).Scan(&status); err != nil {
		t.Fatalf("query issues row: %v", err)
	}
	if status != "failed_gate" {
		t.Fatalf("issues.status = %q, want failed_gate", status)
	}

	var claimedBy sql.NullString
	if err := conn.QueryRow(`
		SELECT claimed_by
		FROM claims
		WHERE ref_type = 'issue' AND ref_id = ?
	`, issueID).Scan(&claimedBy); !errors.Is(err, sql.ErrNoRows) {
		t.Fatalf("query claims row error = %v, want sql.ErrNoRows", err)
	}

	var endedAt sql.NullString
	if err := conn.QueryRow(`SELECT ended_at FROM agent_runs WHERE id = ?`, runID).Scan(&endedAt); err != nil {
		t.Fatalf("query agent_runs row: %v", err)
	}
	if endedAt.Valid {
		t.Fatalf("agent_runs.ended_at = %#v, want open run unchanged", endedAt)
	}

	var notesPath string
	if err := conn.QueryRow(`
		SELECT notes_md_path
		FROM task_completions
		WHERE id = ?
	`, result.Completion.ID).Scan(&notesPath); err != nil {
		t.Fatalf("query task_completions row: %v", err)
	}
	if notesPath != "project-docs/docs/phases/phase-15/tasks/task-3-done.md" {
		t.Fatalf("task_completions.notes_md_path = %q, want reason path", notesPath)
	}
}

func TestBlockSupersedesPriorCompletions(t *testing.T) {
	t.Parallel()

	workspace := newDevWorkspace(t)
	issueID, _, _ := seedFinishReadyIssue(t, workspace)
	first, err := issues.Complete(issues.CompleteOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testDevMigrationsDir(t),
		IssueID:       issueID,
		ActorRole:     "developer",
	})
	if err != nil {
		t.Fatalf("Complete() error = %v", err)
	}

	result, err := Block(BlockOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testDevMigrationsDir(t),
		IssueID:       issueID,
		ActorRole:     "developer",
		ReasonMDPath:  "project-docs/docs/phases/phase-15/tasks/task-3-done.md",
	})
	if err != nil {
		t.Fatalf("Block() error = %v", err)
	}
	if result.Completion.AttemptNumber != 2 {
		t.Fatalf("Block().Completion.AttemptNumber = %d, want 2", result.Completion.AttemptNumber)
	}

	conn := openDevDB(t, filepath.Join(workspace, ".fid", "fid.sqlite"))
	defer closeDevDB(t, conn)

	var firstStatus string
	if err := conn.QueryRow(`SELECT status FROM task_completions WHERE id = ?`, first.ID).Scan(&firstStatus); err != nil {
		t.Fatalf("query prior task_completion: %v", err)
	}
	if firstStatus != "superseded" {
		t.Fatalf("prior task_completion status = %q, want superseded", firstStatus)
	}
}

func TestBlockRejectsEdgeMissingFromProjectedWorkflowRegistry(t *testing.T) {
	workspace := newDevWorkspace(t)
	issueID, _, _ := seedFinishReadyIssue(t, workspace)

	conn := openDevDB(t, filepath.Join(workspace, ".fid", "fid.sqlite"))
	if _, err := conn.Exec(`
		DELETE FROM workflow_transitions
		WHERE workflow_id = ? AND from_state = ? AND to_state = ?
	`, "task-default", "in_progress", "failed_gate"); err != nil {
		closeDevDB(t, conn)
		t.Fatalf("delete projected transition: %v", err)
	}
	closeDevDB(t, conn)

	_, err := Block(BlockOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testDevMigrationsDir(t),
		IssueID:       issueID,
		ActorRole:     "developer",
		ReasonMDPath:  "project-docs/docs/block.md",
	})
	var invalidErr *issues.InvalidTransitionError
	if !errors.As(err, &invalidErr) {
		t.Fatalf("Block() error = %v, want projected-registry invalid transition", err)
	}
}
