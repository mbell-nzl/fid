package issues

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mbell-nzl/fid/internal/journal"
)

func TestTransitionAllowsWorkflowEdgeProjectsActivityAndRebuilds(t *testing.T) {
	t.Parallel()

	workspace := newIssueWorkspace(t)

	created, err := Create(CreateOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
		Kind:          "task",
		Title:         "Transition me",
		ActorRole:     "developer",
	})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	journalPath := journalPathForWorkspace(t, workspace)
	linesBefore := countJournalLines(journalPath)

	result, err := Transition(TransitionOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
		ID:            created.ID,
		ToState:       "in_progress",
		ActorRole:     "developer",
	})
	if err != nil {
		t.Fatalf("Transition() error = %v", err)
	}

	if result.ID != created.ID {
		t.Fatalf("TransitionResult.ID = %q, want %q", result.ID, created.ID)
	}
	if result.FromState != "todo" {
		t.Fatalf("TransitionResult.FromState = %q, want todo", result.FromState)
	}
	if result.ToState != "in_progress" {
		t.Fatalf("TransitionResult.ToState = %q, want in_progress", result.ToState)
	}

	issue, err := Get(GetOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
		ID:            created.ID,
	})
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}
	if issue.Status != "in_progress" {
		t.Fatalf("Get().Status = %q, want in_progress", issue.Status)
	}
	if issue.StatusCategory != "started" {
		t.Fatalf("Get().StatusCategory = %q, want started", issue.StatusCategory)
	}

	if got := countJournalLines(journalPath); got != linesBefore+1 {
		t.Fatalf("journal lines after transition = %d, want %d", got, linesBefore+1)
	}

	stored := readIssueStoredEvents(t, journalPath)
	last := stored[len(stored)-1]
	if last.EventType != "issue.transitioned" {
		t.Fatalf("last event type = %q, want issue.transitioned", last.EventType)
	}

	var (
		updateEffect   *journal.Effect
		activityEffect *journal.Effect
	)
	for index := range last.Effects {
		effect := &last.Effects[index]
		switch {
		case effect.Table == "issues" && effect.Op == "update":
			updateEffect = effect
		case effect.Table == "activity_log" && effect.Op == "activity":
			activityEffect = effect
		}
	}
	if updateEffect == nil {
		t.Fatal("missing issues update effect")
	}
	if activityEffect == nil {
		t.Fatal("missing activity effect")
	}
	if got := updateEffect.Preconditions["status"]; got != "todo" {
		t.Fatalf("update precondition status = %#v, want todo", got)
	}
	if got := updateEffect.Preconditions["workflow_id"]; got != "task-default" {
		t.Fatalf("update precondition workflow_id = %#v, want task-default", got)
	}
	if got := updateEffect.Values["status"]; got != "in_progress" {
		t.Fatalf("update value status = %#v, want in_progress", got)
	}
	if got := activityEffect.Values["event"]; got != "status_changed" {
		t.Fatalf("activity event = %#v, want status_changed", got)
	}
	details, ok := activityEffect.Values["details"].(map[string]any)
	if !ok {
		t.Fatalf("activity details = %#v, want map", activityEffect.Values["details"])
	}
	if got := details["from"]; got != "todo" {
		t.Fatalf("activity details from = %#v, want todo", got)
	}
	if got := details["to"]; got != "in_progress" {
		t.Fatalf("activity details to = %#v, want in_progress", got)
	}

	conn := openIssueDB(t, filepath.Join(workspace, ".fid", "fid.sqlite"))
	defer closeIssueTestDB(t, conn)

	var (
		activityCount int
		actor         string
		eventName     string
	)
	if err := conn.QueryRow(`
		SELECT COUNT(*), MIN(actor), MIN(event)
		FROM activity_log
		WHERE ref_type = 'issue' AND ref_id = ?
	`, created.ID).Scan(&activityCount, &actor, &eventName); err != nil {
		t.Fatalf("query activity_log row: %v", err)
	}
	if activityCount != 1 {
		t.Fatalf("activity_log grouped row count = %d, want 1", activityCount)
	}
	if actor != "developer" {
		t.Fatalf("activity_log actor = %q, want developer", actor)
	}
	if eventName != "status_changed" {
		t.Fatalf("activity_log event = %q, want status_changed", eventName)
	}

	beforeRebuild := issue
	dbPath := filepath.Join(workspace, ".fid", "fid.sqlite")
	if err := os.Remove(dbPath); err != nil {
		t.Fatalf("Remove(%q) error = %v", dbPath, err)
	}

	afterRebuild, err := Get(GetOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
		ID:            created.ID,
	})
	if err != nil {
		t.Fatalf("Get(after rebuild) error = %v", err)
	}
	assertIssuesEqual(t, afterRebuild, beforeRebuild)
}

func TestTransitionRejectsInvalidEdgeWithoutAppending(t *testing.T) {
	t.Parallel()

	workspace := newIssueWorkspace(t)

	created, err := Create(CreateOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
		Kind:          "task",
		Title:         "Reject invalid transition",
		ActorRole:     "developer",
	})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	before, err := Get(GetOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
		ID:            created.ID,
	})
	if err != nil {
		t.Fatalf("Get(before invalid transition) error = %v", err)
	}

	journalPath := journalPathForWorkspace(t, workspace)
	linesBefore := countJournalLines(journalPath)

	_, err = Transition(TransitionOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
		ID:            created.ID,
		ToState:       "done",
		ActorRole:     "developer",
	})
	if err == nil {
		t.Fatal("Transition() error = nil, want rejection")
	}

	var invalidErr *InvalidTransitionError
	if !errors.As(err, &invalidErr) {
		t.Fatalf("Transition() error = %T, want *InvalidTransitionError", err)
	}
	if invalidErr.CurrentState != "todo" {
		t.Fatalf("InvalidTransitionError.CurrentState = %q, want todo", invalidErr.CurrentState)
	}
	if invalidErr.RequestedState != "done" {
		t.Fatalf("InvalidTransitionError.RequestedState = %q, want done", invalidErr.RequestedState)
	}
	for _, want := range []string{`"todo"`, `"done"`} {
		if !strings.Contains(err.Error(), want) {
			t.Fatalf("Transition() error = %q, want substring %q", err.Error(), want)
		}
	}

	after, err := Get(GetOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
		ID:            created.ID,
	})
	if err != nil {
		t.Fatalf("Get(after invalid transition) error = %v", err)
	}
	assertIssuesEqual(t, after, before)

	if got := countJournalLines(journalPath); got != linesBefore {
		t.Fatalf("journal lines after rejected transition = %d, want unchanged %d", got, linesBefore)
	}
}

func TestTransitionRequiresDeveloperVerificationsAndCompletionReport(t *testing.T) {
	t.Parallel()

	workspace := newIssueWorkspace(t)

	created, err := Create(CreateOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
		Kind:          "task",
		Title:         "Developer gate",
		ActorRole:     "developer",
	})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	seedTransitionAcceptanceCriteria(t, workspace, created.ID,
		acceptanceCriterionSeed{ID: "AC-dev-001", Ordinal: 1, Text: "first AC"},
		acceptanceCriterionSeed{ID: "AC-dev-002", Ordinal: 2, Text: "second AC"},
	)

	if _, err := Transition(TransitionOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
		ID:            created.ID,
		ToState:       "in_progress",
		ActorRole:     "developer",
	}); err != nil {
		t.Fatalf("Transition(todo->in_progress) error = %v", err)
	}

	journalPath := journalPathForWorkspace(t, workspace)
	linesBeforeReject := countJournalLines(journalPath)

	_, err = Transition(TransitionOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
		ID:            created.ID,
		ToState:       "dev_complete",
		ActorRole:     "developer",
	})
	if err == nil {
		t.Fatal("Transition(in_progress->dev_complete) error = nil, want rejection")
	}

	var blockedErr *BlockedTransitionError
	if !errors.As(err, &blockedErr) {
		t.Fatalf("Transition() error = %T, want *BlockedTransitionError", err)
	}
	for _, want := range []string{
		"all_acs_verified",
		"has_completion_report",
		"AC-dev-001",
		"AC-dev-002",
	} {
		if !strings.Contains(err.Error(), want) {
			t.Fatalf("Transition() error = %q, want substring %q", err.Error(), want)
		}
	}
	if got := countJournalLines(journalPath); got != linesBeforeReject {
		t.Fatalf("journal lines after rejected rule transition = %d, want unchanged %d", got, linesBeforeReject)
	}

	seedTransitionACVerifications(t, workspace,
		acVerificationSeed{ID: "ACV-dev-001", ACID: "AC-dev-001", Role: "developer", Status: "pass"},
		acVerificationSeed{ID: "ACV-dev-002", ACID: "AC-dev-002", Role: "developer", Status: "pass"},
	)
	seedTransitionCompletionReport(t, workspace, created.ID, "TC-dev-001")

	result, err := Transition(TransitionOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
		ID:            created.ID,
		ToState:       "dev_complete",
		ActorRole:     "developer",
	})
	if err != nil {
		t.Fatalf("Transition(dev_complete success) error = %v", err)
	}
	if result.FromState != "in_progress" || result.ToState != "dev_complete" {
		t.Fatalf("TransitionResult = %+v, want in_progress -> dev_complete", result)
	}

	after, err := Get(GetOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
		ID:            created.ID,
	})
	if err != nil {
		t.Fatalf("Get(after dev_complete) error = %v", err)
	}
	if after.Status != "dev_complete" {
		t.Fatalf("Get().Status = %q, want dev_complete", after.Status)
	}
	if after.StatusCategory != "started" {
		t.Fatalf("Get().StatusCategory = %q, want started", after.StatusCategory)
	}
}

func TestTransitionRequiresQALeadVerificationForQAPass(t *testing.T) {
	t.Parallel()

	workspace := newIssueWorkspace(t)

	created, err := Create(CreateOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
		Kind:          "task",
		Title:         "QA gate",
		ActorRole:     "developer",
	})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	seedTransitionAcceptanceCriteria(t, workspace, created.ID,
		acceptanceCriterionSeed{ID: "AC-qa-001", Ordinal: 1, Text: "qa AC"},
	)
	seedTransitionACVerifications(t, workspace,
		acVerificationSeed{ID: "ACV-qa-dev-001", ACID: "AC-qa-001", Role: "developer", Status: "pass"},
	)
	seedTransitionCompletionReport(t, workspace, created.ID, "TC-qa-001")

	for _, toState := range []string{"in_progress", "dev_complete", "qa_review"} {
		if _, err := Transition(TransitionOptions{
			WorkspaceRoot: workspace,
			MigrationsDir: testIssueMigrationsDir(t),
			ID:            created.ID,
			ToState:       toState,
			ActorRole:     "developer",
		}); err != nil {
			t.Fatalf("Transition(%s) error = %v", toState, err)
		}
	}

	journalPath := journalPathForWorkspace(t, workspace)
	linesBeforeReject := countJournalLines(journalPath)

	_, err = Transition(TransitionOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
		ID:            created.ID,
		ToState:       "qa_pass",
		ActorRole:     "qa_lead",
	})
	if err == nil {
		t.Fatal("Transition(qa_review->qa_pass) error = nil, want rejection")
	}
	for _, want := range []string{
		"all_acs_verified",
		"missing roles qa_lead",
		"AC-qa-001",
	} {
		if !strings.Contains(err.Error(), want) {
			t.Fatalf("Transition() error = %q, want substring %q", err.Error(), want)
		}
	}
	if got := countJournalLines(journalPath); got != linesBeforeReject {
		t.Fatalf("journal lines after rejected qa transition = %d, want unchanged %d", got, linesBeforeReject)
	}

	seedTransitionACVerifications(t, workspace,
		acVerificationSeed{ID: "ACV-qa-lead-001", ACID: "AC-qa-001", Role: "qa_lead", Status: "pass"},
	)

	if _, err := Transition(TransitionOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
		ID:            created.ID,
		ToState:       "qa_pass",
		ActorRole:     "qa_lead",
	}); err != nil {
		t.Fatalf("Transition(qa_pass success) error = %v", err)
	}

	after, err := Get(GetOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
		ID:            created.ID,
	})
	if err != nil {
		t.Fatalf("Get(after qa_pass) error = %v", err)
	}
	if after.Status != "qa_pass" {
		t.Fatalf("Get().Status = %q, want qa_pass", after.Status)
	}
	if after.StatusCategory != "completed" {
		t.Fatalf("Get().StatusCategory = %q, want completed", after.StatusCategory)
	}
}

func TestTransitionRequiresDoneOrCancelledChildrenToWrapEpic(t *testing.T) {
	t.Parallel()

	workspace := newIssueWorkspace(t)

	epic, err := Create(CreateOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
		Kind:          "epic",
		Title:         "Wrap epic",
		ActorRole:     "tl",
		Scheduling:    "sequential",
	})
	if err != nil {
		t.Fatalf("Create(epic) error = %v", err)
	}

	childOne, err := Create(CreateOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
		Kind:          "task",
		Title:         "Child one",
		ActorRole:     "developer",
		ParentID:      stringPointer(epic.ID),
		Sequence:      intPointer(1),
	})
	if err != nil {
		t.Fatalf("Create(child one) error = %v", err)
	}
	childTwo, err := Create(CreateOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
		Kind:          "task",
		Title:         "Child two",
		ActorRole:     "developer",
		ParentID:      stringPointer(epic.ID),
		Sequence:      intPointer(2),
	})
	if err != nil {
		t.Fatalf("Create(child two) error = %v", err)
	}

	if _, err := Transition(TransitionOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
		ID:            epic.ID,
		ToState:       "in_progress",
		ActorRole:     "tl",
	}); err != nil {
		t.Fatalf("Transition(epic todo->in_progress) error = %v", err)
	}

	journalPath := journalPathForWorkspace(t, workspace)
	linesBeforeReject := countJournalLines(journalPath)

	_, err = Transition(TransitionOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
		ID:            epic.ID,
		ToState:       "wrapped",
		ActorRole:     "tl",
	})
	if err == nil {
		t.Fatal("Transition(epic in_progress->wrapped) error = nil, want rejection")
	}
	for _, want := range []string{
		"all_child_tasks_done_or_cancelled",
		childOne.ID,
		childTwo.ID,
	} {
		if !strings.Contains(err.Error(), want) {
			t.Fatalf("Transition() error = %q, want substring %q", err.Error(), want)
		}
	}
	if got := countJournalLines(journalPath); got != linesBeforeReject {
		t.Fatalf("journal lines after rejected epic wrap = %d, want unchanged %d", got, linesBeforeReject)
	}

	seedTransitionStatusUpdate(t, workspace, childOne.ID, "task-default", "todo", "done")
	seedTransitionStatusUpdate(t, workspace, childTwo.ID, "task-default", "todo", "cancelled")

	if _, err := Transition(TransitionOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
		ID:            epic.ID,
		ToState:       "wrapped",
		ActorRole:     "tl",
	}); err != nil {
		t.Fatalf("Transition(epic wrapped success) error = %v", err)
	}

	after, err := Get(GetOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
		ID:            epic.ID,
	})
	if err != nil {
		t.Fatalf("Get(after epic wrap) error = %v", err)
	}
	if after.Status != "wrapped" {
		t.Fatalf("Get().Status = %q, want wrapped", after.Status)
	}
	if after.StatusCategory != "started" {
		t.Fatalf("Get().StatusCategory = %q, want started", after.StatusCategory)
	}
}

func TestTransitionSupportsSubtaskWorkflowAndCancellationRules(t *testing.T) {
	t.Parallel()

	workspace := newIssueWorkspace(t)

	parent, err := Create(CreateOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
		Kind:          "task",
		Title:         "Parent task",
		ActorRole:     "developer",
	})
	if err != nil {
		t.Fatalf("Create(parent) error = %v", err)
	}
	subtask, err := Create(CreateOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
		Kind:          "subtask",
		Title:         "Child subtask",
		ActorRole:     "developer",
		ParentID:      stringPointer(parent.ID),
	})
	if err != nil {
		t.Fatalf("Create(subtask) error = %v", err)
	}
	task, err := Create(CreateOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
		Kind:          "task",
		Title:         "Cancel from todo",
		ActorRole:     "developer",
	})
	if err != nil {
		t.Fatalf("Create(task) error = %v", err)
	}

	if _, err := Transition(TransitionOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
		ID:            subtask.ID,
		ToState:       "in_progress",
		ActorRole:     "developer",
	}); err != nil {
		t.Fatalf("Transition(subtask todo->in_progress) error = %v", err)
	}

	subtaskAfterStart, err := Get(GetOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
		ID:            subtask.ID,
	})
	if err != nil {
		t.Fatalf("Get(subtask after start) error = %v", err)
	}
	if subtaskAfterStart.Status != "in_progress" {
		t.Fatalf("subtask status = %q, want in_progress", subtaskAfterStart.Status)
	}
	if subtaskAfterStart.StatusCategory != "started" {
		t.Fatalf("subtask status_category = %q, want started", subtaskAfterStart.StatusCategory)
	}

	if _, err := Transition(TransitionOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
		ID:            task.ID,
		ToState:       "cancelled",
		ActorRole:     "developer",
	}); err != nil {
		t.Fatalf("Transition(task todo->cancelled) error = %v", err)
	}
	if _, err := Transition(TransitionOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
		ID:            subtask.ID,
		ToState:       "cancelled",
		ActorRole:     "developer",
	}); err != nil {
		t.Fatalf("Transition(subtask in_progress->cancelled) error = %v", err)
	}

	journalPath := journalPathForWorkspace(t, workspace)
	linesBeforeReject := countJournalLines(journalPath)

	_, err = Transition(TransitionOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
		ID:            subtask.ID,
		ToState:       "cancelled",
		ActorRole:     "developer",
	})
	if err == nil {
		t.Fatal("Transition(cancelled->cancelled) error = nil, want rejection")
	}
	for _, want := range []string{`"cancelled"`, `"cancelled"`} {
		if !strings.Contains(err.Error(), want) {
			t.Fatalf("Transition() error = %q, want substring %q", err.Error(), want)
		}
	}
	if got := countJournalLines(journalPath); got != linesBeforeReject {
		t.Fatalf("journal lines after rejected cancelled transition = %d, want unchanged %d", got, linesBeforeReject)
	}
}

type acceptanceCriterionSeed struct {
	ID      string
	Ordinal int
	Text    string
}

type acVerificationSeed struct {
	ID     string
	ACID   string
	Role   string
	Status string
}

func seedTransitionAcceptanceCriteria(t *testing.T, workspace, issueID string, criteria ...acceptanceCriterionSeed) {
	t.Helper()

	effects := make([]journal.Effect, 0, len(criteria))
	for _, criterion := range criteria {
		effects = append(effects, journal.Effect{
			Op:    "insert",
			Table: "acceptance_criteria",
			PK:    map[string]any{"id": criterion.ID},
			Values: map[string]any{
				"issue_id": issueID,
				"ordinal":  criterion.Ordinal,
				"text":     criterion.Text,
			},
		})
	}
	appendIssueJournalEvent(t, workspace, journal.Event{
		EventSchema: journal.EventSchemaV1,
		EventType:   "acceptance_criteria.seeded",
		Actor: journal.Actor{
			Role: "qa_lead",
		},
		Command: journal.Command{
			Argv:       []string{"fid", "maintenance", "seed"},
			CWD:        ".",
			FIDVersion: fidVersion,
		},
		Effects: effects,
	})
}

func seedTransitionACVerifications(t *testing.T, workspace string, verifications ...acVerificationSeed) {
	t.Helper()

	effects := make([]journal.Effect, 0, len(verifications))
	for _, verification := range verifications {
		effects = append(effects, journal.Effect{
			Op:    "insert",
			Table: "ac_verifications",
			PK:    map[string]any{"id": verification.ID},
			Values: map[string]any{
				"ac_id":         verification.ACID,
				"verifier_role": verification.Role,
				"status":        verification.Status,
				"verified_at":   "2026-06-13T10:00:00Z",
			},
		})
	}
	appendIssueJournalEvent(t, workspace, journal.Event{
		EventSchema: journal.EventSchemaV1,
		EventType:   "ac_verifications.seeded",
		Actor: journal.Actor{
			Role: "qa_lead",
		},
		Command: journal.Command{
			Argv:       []string{"fid", "maintenance", "seed"},
			CWD:        ".",
			FIDVersion: fidVersion,
		},
		Effects: effects,
	})
}

func seedTransitionCompletionReport(t *testing.T, workspace, issueID, completionID string) {
	t.Helper()

	appendIssueJournalEvent(t, workspace, journal.Event{
		EventSchema: journal.EventSchemaV1,
		EventType:   "task_completion.seeded",
		Actor: journal.Actor{
			Role: "developer",
		},
		Command: journal.Command{
			Argv:       []string{"fid", "maintenance", "seed"},
			CWD:        ".",
			FIDVersion: fidVersion,
		},
		Effects: []journal.Effect{{
			Op:    "insert",
			Table: "task_completions",
			PK:    map[string]any{"id": completionID},
			Values: map[string]any{
				"issue_id":       issueID,
				"attempt_number": 1,
				"status":         "dev_complete",
				"completed_at":   "2026-06-13T10:05:00Z",
			},
		}},
	})
}

func seedTransitionStatusUpdate(t *testing.T, workspace, issueID, workflowID, fromStatus, toStatus string) {
	t.Helper()

	appendIssueJournalEvent(t, workspace, journal.Event{
		EventSchema: journal.EventSchemaV1,
		EventType:   "issue.seed_status_changed",
		Actor: journal.Actor{
			Role: "developer",
		},
		Command: journal.Command{
			Argv:       []string{"fid", "maintenance", "seed"},
			CWD:        ".",
			FIDVersion: fidVersion,
		},
		Effects: []journal.Effect{{
			Op:    "update",
			Table: "issues",
			PK:    map[string]any{"id": issueID},
			Values: map[string]any{
				"status":     toStatus,
				"updated_at": "2026-06-13T10:10:00Z",
			},
			Prior: map[string]any{
				"status": fromStatus,
			},
			Preconditions: map[string]any{
				"status":      fromStatus,
				"workflow_id": workflowID,
			},
		}},
	})
}
