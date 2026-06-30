package dev

import (
	"database/sql"
	"errors"
	"path/filepath"
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

	conn := openDevDB(t, filepath.Join(workspace, ".fid", "fid.sqlite"))
	defer closeDevDB(t, conn)

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
