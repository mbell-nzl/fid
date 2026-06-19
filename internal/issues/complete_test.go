package issues

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestCompleteCreatesAttemptOneAndPersistsCompletionFields(t *testing.T) {
	t.Parallel()

	workspace := newIssueWorkspace(t)

	created, err := Create(CreateOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
		Kind:          "task",
		Title:         "Ship completion write path",
		ActorRole:     "developer",
	})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	result, err := Complete(CompleteOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
		IssueID:       created.ID,
		ActorRole:     "developer",
		NotesMDPath:   stringPointer("project-docs/docs/phases/phase-8/tasks/task-1-done.md"),
		ActualTokens:  intPointer(14823),
		CommitHashes:  []string{"abc123", "def456"},
		FilesTouched:  []string{"internal/issues/complete.go", "cmd/fid/issue.go"},
	})
	if err != nil {
		t.Fatalf("Complete() error = %v", err)
	}

	if !strings.HasPrefix(result.ID, "TC-") {
		t.Fatalf("Complete().ID = %q, want TC- prefix", result.ID)
	}
	if result.IssueID != created.ID {
		t.Fatalf("Complete().IssueID = %q, want %q", result.IssueID, created.ID)
	}
	if result.AttemptNumber != 1 {
		t.Fatalf("Complete().AttemptNumber = %d, want 1", result.AttemptNumber)
	}
	if result.Status != "dev_complete" {
		t.Fatalf("Complete().Status = %q, want dev_complete", result.Status)
	}

	got, err := Get(GetOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
		ID:            created.ID,
	})
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}

	if len(got.Completions) != 1 {
		t.Fatalf("Get().Completions len = %d, want 1", len(got.Completions))
	}

	completion := got.Completions[0]
	if completion.ID != result.ID {
		t.Fatalf("Get().Completions[0].ID = %q, want %q", completion.ID, result.ID)
	}
	if completion.AttemptNumber != 1 {
		t.Fatalf("Get().Completions[0].AttemptNumber = %d, want 1", completion.AttemptNumber)
	}
	if completion.DeveloperRole == nil || *completion.DeveloperRole != "developer" {
		t.Fatalf("Get().Completions[0].DeveloperRole = %#v, want developer", completion.DeveloperRole)
	}
	if completion.Status == nil || *completion.Status != "dev_complete" {
		t.Fatalf("Get().Completions[0].Status = %#v, want dev_complete", completion.Status)
	}
	if completion.NotesMDPath == nil || *completion.NotesMDPath != "project-docs/docs/phases/phase-8/tasks/task-1-done.md" {
		t.Fatalf("Get().Completions[0].NotesMDPath = %#v, want completion report path", completion.NotesMDPath)
	}
	if completion.ActualTokens == nil || *completion.ActualTokens != 14823 {
		t.Fatalf("Get().Completions[0].ActualTokens = %#v, want 14823", completion.ActualTokens)
	}
	if !equalStrings(completion.CommitHashes, []string{"abc123", "def456"}) {
		t.Fatalf("Get().Completions[0].CommitHashes = %#v, want commits", completion.CommitHashes)
	}
	if !equalStrings(completion.FilesTouched, []string{"internal/issues/complete.go", "cmd/fid/issue.go"}) {
		t.Fatalf("Get().Completions[0].FilesTouched = %#v, want files", completion.FilesTouched)
	}
	if completion.HumanIntervention {
		t.Fatal("Get().Completions[0].HumanIntervention = true, want false")
	}
	if completion.HumanInterventionNote != nil {
		t.Fatalf("Get().Completions[0].HumanInterventionNote = %#v, want nil", completion.HumanInterventionNote)
	}
}

func TestCompleteSupersedesPriorAttemptAndRetainsHistory(t *testing.T) {
	t.Parallel()

	workspace := newIssueWorkspace(t)

	created, err := Create(CreateOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
		Kind:          "task",
		Title:         "Record multiple completion attempts",
		ActorRole:     "developer",
	})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	first, err := Complete(CompleteOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
		IssueID:       created.ID,
		ActorRole:     "developer",
		ActualTokens:  intPointer(4000),
	})
	if err != nil {
		t.Fatalf("Complete(first) error = %v", err)
	}
	second, err := Complete(CompleteOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
		IssueID:       created.ID,
		ActorRole:     "developer",
		ActualTokens:  intPointer(5000),
	})
	if err != nil {
		t.Fatalf("Complete(second) error = %v", err)
	}

	if first.AttemptNumber != 1 {
		t.Fatalf("first attempt_number = %d, want 1", first.AttemptNumber)
	}
	if second.AttemptNumber != 2 {
		t.Fatalf("second attempt_number = %d, want 2", second.AttemptNumber)
	}

	got, err := Get(GetOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
		ID:            created.ID,
	})
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}
	if len(got.Completions) != 2 {
		t.Fatalf("Get().Completions len = %d, want 2", len(got.Completions))
	}
	if got.Completions[0].ID != second.ID || got.Completions[0].AttemptNumber != 2 {
		t.Fatalf("latest completion = %#v, want second attempt", got.Completions[0])
	}
	if got.Completions[0].Status == nil || *got.Completions[0].Status != "dev_complete" {
		t.Fatalf("latest completion status = %#v, want dev_complete", got.Completions[0].Status)
	}
	if got.Completions[1].ID != first.ID || got.Completions[1].AttemptNumber != 1 {
		t.Fatalf("prior completion = %#v, want first attempt", got.Completions[1])
	}
	if got.Completions[1].Status == nil || *got.Completions[1].Status != "superseded" {
		t.Fatalf("prior completion status = %#v, want superseded", got.Completions[1].Status)
	}
}

func TestCompleteStoresHumanInterventionStateWithoutActivityLogRows(t *testing.T) {
	t.Parallel()

	workspace := newIssueWorkspace(t)

	created, err := Create(CreateOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
		Kind:          "task",
		Title:         "Capture intervention details",
		ActorRole:     "developer",
	})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	conn := openIssueDB(t, filepath.Join(workspace, ".fid", "fid.sqlite"))
	defer closeIssueTestDB(t, conn)

	var activityBefore int
	if err := conn.QueryRow(`SELECT COUNT(*) FROM activity_log`).Scan(&activityBefore); err != nil {
		t.Fatalf("count activity_log before: %v", err)
	}

	if _, err := Complete(CompleteOptions{
		WorkspaceRoot:         workspace,
		MigrationsDir:         testIssueMigrationsDir(t),
		IssueID:               created.ID,
		ActorRole:             "developer",
		HumanIntervention:     true,
		HumanInterventionNote: stringPointer("paired on schema fix"),
	}); err != nil {
		t.Fatalf("Complete(flagged) error = %v", err)
	}

	got, err := Get(GetOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
		ID:            created.ID,
	})
	if err != nil {
		t.Fatalf("Get() after flagged completion error = %v", err)
	}
	if len(got.Completions) != 1 {
		t.Fatalf("Get().Completions len after flagged = %d, want 1", len(got.Completions))
	}
	if !got.Completions[0].HumanIntervention {
		t.Fatal("flagged completion human_intervention = false, want true")
	}
	if got.Completions[0].HumanInterventionNote == nil || *got.Completions[0].HumanInterventionNote != "paired on schema fix" {
		t.Fatalf("flagged completion note = %#v, want paired on schema fix", got.Completions[0].HumanInterventionNote)
	}

	unflaggedCreated, err := Create(CreateOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
		Kind:          "task",
		Title:         "No intervention",
		ActorRole:     "developer",
	})
	if err != nil {
		t.Fatalf("Create(unflagged) error = %v", err)
	}
	if _, err := Complete(CompleteOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
		IssueID:       unflaggedCreated.ID,
		ActorRole:     "developer",
	}); err != nil {
		t.Fatalf("Complete(unflagged) error = %v", err)
	}

	unflagged, err := Get(GetOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
		ID:            unflaggedCreated.ID,
	})
	if err != nil {
		t.Fatalf("Get() after unflagged completion error = %v", err)
	}
	if len(unflagged.Completions) != 1 {
		t.Fatalf("unflagged completions len = %d, want 1", len(unflagged.Completions))
	}
	if unflagged.Completions[0].HumanIntervention {
		t.Fatal("unflagged completion human_intervention = true, want false")
	}
	if unflagged.Completions[0].HumanInterventionNote != nil {
		t.Fatalf("unflagged completion note = %#v, want nil", unflagged.Completions[0].HumanInterventionNote)
	}

	var activityAfter int
	if err := conn.QueryRow(`SELECT COUNT(*) FROM activity_log`).Scan(&activityAfter); err != nil {
		t.Fatalf("count activity_log after: %v", err)
	}
	if activityAfter != activityBefore {
		t.Fatalf("activity_log rows = %d after completions, want unchanged %d", activityAfter, activityBefore)
	}
}

func TestCompleteRejectsEpicBeforeAppend(t *testing.T) {
	t.Parallel()

	workspace := newIssueWorkspace(t)

	created, err := Create(CreateOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
		Kind:          "epic",
		Title:         "Epic completion reject",
		ActorRole:     "tl",
		Scheduling:    "sequential",
	})
	if err != nil {
		t.Fatalf("Create(epic) error = %v", err)
	}

	journalPath := journalPathForWorkspace(t, workspace)
	linesBefore := countJournalLines(journalPath)

	_, err = Complete(CompleteOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
		IssueID:       created.ID,
		ActorRole:     "developer",
	})
	if err == nil {
		t.Fatal("Complete(epic) error = nil, want rejection")
	}
	if !strings.Contains(err.Error(), `is kind "epic"; task completions require task or subtask`) {
		t.Fatalf("Complete(epic) error = %q, want kind rejection", err.Error())
	}
	if got := countJournalLines(journalPath); got != linesBefore {
		t.Fatalf("journal line count after epic rejection = %d, want unchanged %d", got, linesBefore)
	}

	conn := openIssueDB(t, filepath.Join(workspace, ".fid", "fid.sqlite"))
	defer closeIssueTestDB(t, conn)

	var count int
	if err := conn.QueryRow(`SELECT COUNT(*) FROM task_completions`).Scan(&count); err != nil {
		t.Fatalf("count task_completions: %v", err)
	}
	if count != 0 {
		t.Fatalf("task_completions rows = %d, want 0", count)
	}
}
