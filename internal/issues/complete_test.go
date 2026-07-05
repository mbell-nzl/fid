package issues

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mbell-nzl/fid/internal/journal"
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

	stored := readIssueStoredEvents(t, journalPathForWorkspace(t, workspace))
	last := stored[len(stored)-1]
	if last.EventType != "issue.completed" {
		t.Fatalf("last event type = %q, want issue.completed", last.EventType)
	}

	var activityEffect *journal.Effect
	for index := range last.Effects {
		effect := &last.Effects[index]
		if effect.Table == "activity_log" && effect.Op == "activity" {
			activityEffect = effect
			break
		}
	}
	if activityEffect == nil {
		t.Fatal("missing task_completed activity effect")
	}
	if got := activityEffect.Values["event"]; got != "task_completed" {
		t.Fatalf("activity event = %#v, want task_completed", got)
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

	conn := openIssueDB(t, filepath.Join(workspace, ".fid", "fid.sqlite"))
	defer closeIssueTestDB(t, conn)
	assertIssueLabelActivity(t, conn, last.EventID, "task_completed", "developer", result.ID, `"status":"dev_complete"`, `"attempt_number":1`)
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

func TestCompleteStoresHumanInterventionStateAndProjectsActivity(t *testing.T) {
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

	rows, err := conn.Query(`
		SELECT event, COALESCE(details, '')
		FROM activity_log
		WHERE ref_type = 'issue' AND ref_id IN (?, ?)
		ORDER BY occurred_at, id
	`, created.ID, unflaggedCreated.ID)
	if err != nil {
		t.Fatalf("query completion activity rows: %v", err)
	}
	defer closeIssueRows(rows)

	type completionActivity struct {
		event   string
		details string
	}
	var completionActivities []completionActivity
	for rows.Next() {
		var activity completionActivity
		if err := rows.Scan(&activity.event, &activity.details); err != nil {
			t.Fatalf("scan completion activity row: %v", err)
		}
		if activity.event == "task_completed" {
			completionActivities = append(completionActivities, activity)
		}
	}
	if err := rows.Err(); err != nil {
		t.Fatalf("iterate completion activity rows: %v", err)
	}
	if len(completionActivities) != 2 {
		t.Fatalf("task_completed activity rows = %d, want 2", len(completionActivities))
	}
	if !strings.Contains(completionActivities[0].details, `"attempt_number":1`) {
		t.Fatalf("flagged completion activity = %#v, want attempt 1", completionActivities[0])
	}
	if !strings.Contains(completionActivities[1].details, `"attempt_number":1`) {
		t.Fatalf("unflagged completion activity = %#v, want attempt 1", completionActivities[1])
	}
}

func TestCompleteSnapshotsMatchingGateRunsAndPreservesSourceRows(t *testing.T) {
	t.Parallel()

	workspace := newIssueWorkspace(t)

	created, err := Create(CreateOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
		Kind:          "task",
		Title:         "Snapshot gate runs",
		ActorRole:     "developer",
	})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	seedCompletionGateRuns(t, workspace, created.ID, []completionGateRunSeed{
		{
			ID:             "GATE-000101",
			AttemptNumber:  1,
			RecordedByRole: "developer",
			Command:        "go test ./...",
			ExitCode:       0,
			OutputExcerpt:  stringPointer("ok\tgithub.com/mbell-nzl/fid/internal/issues\t0.321s"),
			RanAt:          "2026-06-20T09:00:00Z",
		},
		{
			ID:             "GATE-000102",
			AttemptNumber:  1,
			RecordedByRole: "developer",
			Command:        "golangci-lint run ./...",
			ExitCode:       1,
			OutputExcerpt:  stringPointer("internal/issues/complete.go:128: line is 61 chars"),
			RanAt:          "2026-06-20T09:01:00Z",
		},
		{
			ID:             "GATE-000103",
			AttemptNumber:  2,
			RecordedByRole: "developer",
			Command:        "go test ./... -run TestFutureAttempt",
			ExitCode:       0,
			OutputExcerpt:  stringPointer("ok"),
			RanAt:          "2026-06-20T09:02:00Z",
		},
	})

	result, err := Complete(CompleteOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
		IssueID:       created.ID,
		ActorRole:     "developer",
	})
	if err != nil {
		t.Fatalf("Complete() error = %v", err)
	}

	conn := openIssueDB(t, filepath.Join(workspace, ".fid", "fid.sqlite"))
	defer closeIssueTestDB(t, conn)

	var gatesPassed, gatesFailed int
	var gateResultsJSON string
	if err := conn.QueryRow(`
		SELECT gates_passed, gates_failed, gate_results
		FROM task_completions
		WHERE id = ?
	`, result.ID).Scan(&gatesPassed, &gatesFailed, &gateResultsJSON); err != nil {
		t.Fatalf("query completion gate snapshot: %v", err)
	}
	if gatesPassed != 1 {
		t.Fatalf("gates_passed = %d, want 1", gatesPassed)
	}
	if gatesFailed != 1 {
		t.Fatalf("gates_failed = %d, want 1", gatesFailed)
	}

	var gateResults []map[string]any
	if err := json.Unmarshal([]byte(gateResultsJSON), &gateResults); err != nil {
		t.Fatalf("Unmarshal(gate_results) error = %v", err)
	}
	wantGateResults := []map[string]any{
		{
			"cmd":            "go test ./...",
			"exit_code":      float64(0),
			"ran_at":         "2026-06-20T09:00:00Z",
			"output_excerpt": "ok\tgithub.com/mbell-nzl/fid/internal/issues\t0.321s",
		},
		{
			"cmd":            "golangci-lint run ./...",
			"exit_code":      float64(1),
			"ran_at":         "2026-06-20T09:01:00Z",
			"output_excerpt": "internal/issues/complete.go:128: line is 61 chars",
		},
	}
	if !equalGateResultMaps(gateResults, wantGateResults) {
		t.Fatalf("gate_results = %#v, want %#v", gateResults, wantGateResults)
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
	if got.Completions[0].GatesPassed != 1 {
		t.Fatalf("Get().Completions[0].GatesPassed = %d, want 1", got.Completions[0].GatesPassed)
	}
	if got.Completions[0].GatesFailed != 1 {
		t.Fatalf("Get().Completions[0].GatesFailed = %d, want 1", got.Completions[0].GatesFailed)
	}
	if len(got.Completions[0].GateResults) != 2 {
		t.Fatalf("Get().Completions[0].GateResults len = %d, want 2", len(got.Completions[0].GateResults))
	}
	if got.Completions[0].GateResults[0].Cmd != "go test ./..." || got.Completions[0].GateResults[0].ExitCode != 0 {
		t.Fatalf("Get().Completions[0].GateResults[0] = %#v, want go test pass", got.Completions[0].GateResults[0])
	}
	if got.Completions[0].GateResults[1].Cmd != "golangci-lint run ./..." || got.Completions[0].GateResults[1].ExitCode != 1 {
		t.Fatalf("Get().Completions[0].GateResults[1] = %#v, want lint fail", got.Completions[0].GateResults[1])
	}

	rows, err := conn.Query(`
		SELECT id, attempt_number, command, exit_code, output_excerpt, ran_at
		FROM gate_runs
		WHERE issue_id = ?
		ORDER BY attempt_number, ran_at, id
	`, created.ID)
	if err != nil {
		t.Fatalf("query gate_runs: %v", err)
	}
	defer closeIssueRows(rows)

	var preserved []string
	for rows.Next() {
		var id string
		var attemptNumber int
		var command string
		var exitCode int
		var outputExcerpt sql.NullString
		var ranAt string
		if err := rows.Scan(&id, &attemptNumber, &command, &exitCode, &outputExcerpt, &ranAt); err != nil {
			t.Fatalf("scan gate_runs: %v", err)
		}
		excerpt := "<nil>"
		if outputExcerpt.Valid {
			excerpt = outputExcerpt.String
		}
		preserved = append(preserved, fmt.Sprintf("%s|%s|%s|%s:%d:%d", id, command, ranAt, excerpt, attemptNumber, exitCode))
	}
	if err := rows.Err(); err != nil {
		t.Fatalf("iterate gate_runs: %v", err)
	}
	wantPreserved := []string{
		"GATE-000101|go test ./...|2026-06-20T09:00:00Z|ok\tgithub.com/mbell-nzl/fid/internal/issues\t0.321s:1:0",
		"GATE-000102|golangci-lint run ./...|2026-06-20T09:01:00Z|internal/issues/complete.go:128: line is 61 chars:1:1",
		"GATE-000103|go test ./... -run TestFutureAttempt|2026-06-20T09:02:00Z|ok:2:0",
	}
	if !equalStrings(preserved, wantPreserved) {
		t.Fatalf("preserved gate_runs = %#v, want %#v", preserved, wantPreserved)
	}
}

func TestCompleteLeavesGateSnapshotEmptyWhenNoMatchingGateRunsExist(t *testing.T) {
	t.Parallel()

	workspace := newIssueWorkspace(t)

	created, err := Create(CreateOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
		Kind:          "task",
		Title:         "No matching gates",
		ActorRole:     "developer",
	})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	seedCompletionGateRuns(t, workspace, created.ID, []completionGateRunSeed{
		{
			ID:             "GATE-000201",
			AttemptNumber:  2,
			RecordedByRole: "developer",
			Command:        "go test ./... -run TestAttemptTwo",
			ExitCode:       0,
			RanAt:          "2026-06-20T10:00:00Z",
		},
	})

	if _, err := Complete(CompleteOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
		IssueID:       created.ID,
		ActorRole:     "developer",
	}); err != nil {
		t.Fatalf("Complete() error = %v", err)
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
	if len(got.Completions[0].GateResults) != 0 {
		t.Fatalf("Get().Completions[0].GateResults = %#v, want empty", got.Completions[0].GateResults)
	}
	if got.Completions[0].GatesPassed != 0 {
		t.Fatalf("Get().Completions[0].GatesPassed = %d, want 0", got.Completions[0].GatesPassed)
	}
	if got.Completions[0].GatesFailed != 0 {
		t.Fatalf("Get().Completions[0].GatesFailed = %d, want 0", got.Completions[0].GatesFailed)
	}
}

type completionGateRunSeed struct {
	ID             string
	AttemptNumber  int
	RecordedByRole string
	Command        string
	ExitCode       int
	OutputExcerpt  *string
	RanAt          string
}

func seedCompletionGateRuns(t *testing.T, workspace, issueID string, runs []completionGateRunSeed) {
	t.Helper()

	effects := make([]journal.Effect, 0, len(runs))
	for _, run := range runs {
		values := map[string]any{
			"issue_id":         issueID,
			"attempt_number":   run.AttemptNumber,
			"recorded_by_role": run.RecordedByRole,
			"command":          run.Command,
			"exit_code":        run.ExitCode,
			"ran_at":           run.RanAt,
		}
		if run.OutputExcerpt != nil {
			values["output_excerpt"] = *run.OutputExcerpt
		}
		effects = append(effects, journal.Effect{
			Op:    "insert",
			Table: "gate_runs",
			PK: map[string]any{
				"id": run.ID,
			},
			Values: values,
		})
	}

	appendIssueJournalEvent(t, workspace, journal.Event{
		EventSchema: journal.EventSchemaV1,
		EventType:   "gate.seeded",
		Actor: journal.Actor{
			Role: "developer",
		},
		Command: journal.Command{
			Argv:       []string{"fid", "maintenance", "seed"},
			CWD:        ".",
			FIDVersion: "0.1.0",
		},
		Effects: effects,
	})
}

func equalGateResultMaps(left, right []map[string]any) bool {
	if len(left) != len(right) {
		return false
	}
	for index := range left {
		if len(left[index]) != len(right[index]) {
			return false
		}
		for key, rightValue := range right[index] {
			if left[index][key] != rightValue {
				return false
			}
		}
	}
	return true
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
