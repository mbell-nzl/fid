package dev

import (
	"database/sql"
	"encoding/json"
	"errors"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mbell-nzl/fid/internal/issues"
)

func TestFinishClosesRunCompletesTransitionsAndReleasesClaimInOneEvent(t *testing.T) {
	t.Parallel()

	workspace := newDevWorkspace(t)
	issueID, runID, criterionID := seedFinishReadyIssue(t, workspace)
	recordFinishDeveloperGate(t, workspace, issueID)
	verifyDeveloperAC(t, workspace, issueID, criterionID)

	journalPath := journalPathForDevWorkspace(t, workspace)
	linesBefore := countDevJournalLines(journalPath)

	result, err := Finish(FinishOptions{
		WorkspaceRoot:         workspace,
		MigrationsDir:         testDevMigrationsDir(t),
		IssueID:               issueID,
		ActorRole:             "developer",
		CompletionMDPath:      "project-docs/docs/phases/phase-15/tasks/task-3-done.md",
		ActualTokens:          intPointer(14823),
		CommitHashes:          []string{"abc123", "def456"},
		FilesTouched:          []string{"internal/dev/finish.go", "cmd/fid/dev.go"},
		HumanIntervention:     true,
		HumanInterventionNote: stringPointer("paired on final review"),
	})
	if err != nil {
		t.Fatalf("Finish() error = %v", err)
	}

	if result.IssueID != issueID {
		t.Fatalf("Finish().IssueID = %q, want %q", result.IssueID, issueID)
	}
	if result.FromState != "in_progress" {
		t.Fatalf("Finish().FromState = %q, want in_progress", result.FromState)
	}
	if result.ToState != "dev_complete" {
		t.Fatalf("Finish().ToState = %q, want dev_complete", result.ToState)
	}
	if result.RunID == nil || *result.RunID != runID {
		t.Fatalf("Finish().RunID = %#v, want %q", result.RunID, runID)
	}
	if result.Completion.AttemptNumber != 1 {
		t.Fatalf("Finish().Completion.AttemptNumber = %d, want 1", result.Completion.AttemptNumber)
	}
	if result.Completion.Status != "dev_complete" {
		t.Fatalf("Finish().Completion.Status = %q, want dev_complete", result.Completion.Status)
	}

	if got := countDevJournalLines(journalPath); got != linesBefore+1 {
		t.Fatalf("journal lines after Finish() = %d, want %d", got, linesBefore+1)
	}

	stored := readDevStoredEvents(t, journalPath)
	last := stored[len(stored)-1]
	if last.EventType != "dev.finished" {
		t.Fatalf("last event type = %q, want dev.finished", last.EventType)
	}

	conn := openDevDB(t, filepath.Join(workspace, ".fid", "fid.sqlite"))
	defer closeDevDB(t, conn)

	var (
		status     string
		claimedBy  sql.NullString
		endedAt    sql.NullString
		gatesPass  int
		gatesFail  int
		notesPath  sql.NullString
		humanNote  sql.NullString
		gateResult string
	)
	if err := conn.QueryRow(`SELECT status FROM issues WHERE id = ?`, issueID).Scan(&status); err != nil {
		t.Fatalf("query issues row: %v", err)
	}
	if status != "dev_complete" {
		t.Fatalf("issues.status = %q, want dev_complete", status)
	}
	if err := conn.QueryRow(`
		SELECT claimed_by
		FROM claims
		WHERE ref_type = 'issue' AND ref_id = ?
	`, issueID).Scan(&claimedBy); !errors.Is(err, sql.ErrNoRows) {
		t.Fatalf("query claims row error = %v, want sql.ErrNoRows", err)
	}
	if err := conn.QueryRow(`SELECT ended_at FROM agent_runs WHERE id = ?`, runID).Scan(&endedAt); err != nil {
		t.Fatalf("query agent_runs row: %v", err)
	}
	if !endedAt.Valid || strings.TrimSpace(endedAt.String) == "" {
		t.Fatalf("agent_runs.ended_at = %#v, want closed run timestamp", endedAt)
	}
	if err := conn.QueryRow(`
		SELECT gates_passed, gates_failed, gate_results, notes_md_path, human_intervention_note
		FROM task_completions
		WHERE id = ?
	`, result.Completion.ID).Scan(&gatesPass, &gatesFail, &gateResult, &notesPath, &humanNote); err != nil {
		t.Fatalf("query task_completions row: %v", err)
	}
	if gatesPass != 1 || gatesFail != 0 {
		t.Fatalf("completion gate counts = (%d, %d), want (1, 0)", gatesPass, gatesFail)
	}
	if !notesPath.Valid || notesPath.String != "project-docs/docs/phases/phase-15/tasks/task-3-done.md" {
		t.Fatalf("task_completions.notes_md_path = %#v, want completion report path", notesPath)
	}
	if !humanNote.Valid || humanNote.String != "paired on final review" {
		t.Fatalf("task_completions.human_intervention_note = %#v, want paired on final review", humanNote)
	}

	var snapshots []map[string]any
	if err := json.Unmarshal([]byte(gateResult), &snapshots); err != nil {
		t.Fatalf("Unmarshal(gate_results) error = %v", err)
	}
	if len(snapshots) != 1 {
		t.Fatalf("gate_results len = %d, want 1", len(snapshots))
	}
}

func TestFinishRejectsWhenDeveloperACsUnverifiedBeforeAppend(t *testing.T) {
	t.Parallel()

	workspace := newDevWorkspace(t)
	issueID, _, _ := seedFinishReadyIssue(t, workspace)

	journalPath := journalPathForDevWorkspace(t, workspace)
	linesBefore := countDevJournalLines(journalPath)

	_, err := Finish(FinishOptions{
		WorkspaceRoot:    workspace,
		MigrationsDir:    testDevMigrationsDir(t),
		IssueID:          issueID,
		ActorRole:        "developer",
		CompletionMDPath: "project-docs/docs/phases/phase-15/tasks/task-3-done.md",
	})
	if err == nil {
		t.Fatal("Finish() error = nil, want blocked transition")
	}

	var blockedErr *issues.BlockedTransitionError
	if !errors.As(err, &blockedErr) {
		t.Fatalf("Finish() error = %T, want *issues.BlockedTransitionError", err)
	}
	if !strings.Contains(err.Error(), "all_acs_verified") {
		t.Fatalf("Finish() error = %q, want all_acs_verified", err.Error())
	}
	if got := countDevJournalLines(journalPath); got != linesBefore {
		t.Fatalf("journal lines after rejected Finish() = %d, want unchanged %d", got, linesBefore)
	}
}

func TestFinishSucceedsWithoutPreexistingCompletionOnceDeveloperACsPass(t *testing.T) {
	t.Parallel()

	workspace := newDevWorkspace(t)
	issueID, _, criterionID := seedFinishReadyIssue(t, workspace)
	verifyDeveloperAC(t, workspace, issueID, criterionID)

	result, err := Finish(FinishOptions{
		WorkspaceRoot:    workspace,
		MigrationsDir:    testDevMigrationsDir(t),
		IssueID:          issueID,
		ActorRole:        "developer",
		CompletionMDPath: "project-docs/docs/phases/phase-15/tasks/task-3-done.md",
	})
	if err != nil {
		t.Fatalf("Finish() error = %v", err)
	}
	if result.Completion.AttemptNumber != 1 {
		t.Fatalf("Finish().Completion.AttemptNumber = %d, want 1", result.Completion.AttemptNumber)
	}

	conn := openDevDB(t, filepath.Join(workspace, ".fid", "fid.sqlite"))
	defer closeDevDB(t, conn)

	var count int
	if err := conn.QueryRow(`SELECT COUNT(*) FROM task_completions WHERE issue_id = ?`, issueID).Scan(&count); err != nil {
		t.Fatalf("count task_completions: %v", err)
	}
	if count != 1 {
		t.Fatalf("task_completions count = %d, want 1", count)
	}
}

func seedFinishReadyIssue(t *testing.T, workspace string) (string, string, string) {
	t.Helper()

	created, err := issues.Create(issues.CreateOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testDevMigrationsDir(t),
		Kind:          "task",
		Title:         "Finish composite target",
		ActorRole:     "developer",
	})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	criterion, err := issues.AddAcceptanceCriterion(issues.AddAcceptanceCriterionOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testDevMigrationsDir(t),
		IssueID:       created.ID,
		Text:          "Developer verified acceptance criterion",
		Ordinal:       1,
		ActorRole:     "developer",
	})
	if err != nil {
		t.Fatalf("AddAcceptanceCriterion() error = %v", err)
	}
	started, err := Start(StartOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testDevMigrationsDir(t),
		IssueID:       created.ID,
		ActorRole:     "developer",
		Model:         "claude-sonnet-4-6",
	})
	if err != nil {
		t.Fatalf("Start() error = %v", err)
	}
	return created.ID, started.RunID, criterion.ID
}

func recordFinishDeveloperGate(t *testing.T, workspace, issueID string) {
	t.Helper()

	if _, err := Gate(GateOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testDevMigrationsDir(t),
		IssueID:       issueID,
		ActorRole:     "developer",
		Command:       "go test ./...",
		ExitCode:      0,
		OutputExcerpt: stringPointer("ok"),
	}); err != nil {
		t.Fatalf("Gate() error = %v", err)
	}
}

func verifyDeveloperAC(t *testing.T, workspace, issueID, criterionID string) {
	t.Helper()

	if _, err := issues.VerifyAcceptanceCriterion(issues.VerifyAcceptanceCriterionOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testDevMigrationsDir(t),
		IssueID:       issueID,
		CriterionID:   criterionID,
		ActorRole:     "developer",
		Status:        "pass",
	}); err != nil {
		t.Fatalf("VerifyAcceptanceCriterion() error = %v", err)
	}
}

func intPointer(value int) *int {
	return &value
}
