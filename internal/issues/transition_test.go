package issues

import (
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/mbell-nzl/fid/internal/config"
	"github.com/mbell-nzl/fid/internal/issuestate"
	"github.com/mbell-nzl/fid/internal/journal"
)

func TestTransitionAllowsWorkflowEdgeProjectsActivityAndRebuilds(t *testing.T) {
	t.Parallel()

	workspace := newIssueWorkspace(t)

	created, err := createTransitionSubtask(t, CreateOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
		Kind:          "task",
		Title:         "Transition me",
		Actor:         testActorForRole("developer"), IdempotencyKey: testIdempotencyKey(),
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
		Actor:         testActorForRole("developer"), IdempotencyKey: testIdempotencyKey(),
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
		revisionEffect *journal.Effect
		activityEffect *journal.Effect
	)
	for index := range last.Effects {
		effect := &last.Effects[index]
		switch {
		case effect.Table == "issues" && effect.Op == "update" && effect.Values["status"] != nil:
			updateEffect = effect
		case effect.Table == "issues" && effect.Op == "update" && effect.Values["revision"] != nil:
			revisionEffect = effect
		case effect.Table == "activity_log" && effect.Op == "activity":
			activityEffect = effect
		}
	}
	if updateEffect == nil {
		t.Fatal("missing issues update effect")
	}
	if revisionEffect == nil {
		t.Fatal("missing Issue revision effect")
	}
	if got := revisionEffect.Values["revision"]; got != float64(2) {
		t.Fatalf("revision effect value = %#v, want 2", got)
	}
	if got := revisionEffect.Prior["revision"]; got != float64(1) {
		t.Fatalf("revision effect prior = %#v, want 1", got)
	}
	if got := revisionEffect.Preconditions["revision"]; got != float64(1) {
		t.Fatalf("revision effect precondition = %#v, want 1", got)
	}
	if activityEffect == nil {
		t.Fatal("missing activity effect")
	}
	if got := updateEffect.Preconditions["status"]; got != "todo" {
		t.Fatalf("update precondition status = %#v, want todo", got)
	}
	if got := updateEffect.Preconditions["workflow_id"]; got != "subtask-default" {
		t.Fatalf("update precondition workflow_id = %#v, want subtask-default", got)
	}
	if got := updateEffect.Values["status"]; got != "in_progress" {
		t.Fatalf("update value status = %#v, want in_progress", got)
	}
	if got := activityEffect.Values["event"]; got != "status-changed" {
		t.Fatalf("activity event = %#v, want status-changed", got)
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

	conn := openIssueDB(t, filepath.Join(config.LedgerDir(workspace), "fid.sqlite"))
	defer closeIssueTestDB(t, conn)

	var (
		activityCount int
		actor         string
		eventName     string
	)
	if err := conn.QueryRow(`
		SELECT COUNT(*), MIN(actor), MIN(event)
		FROM activity_log
		WHERE journal_event_id = ?
	`, last.EventID).Scan(&activityCount, &actor, &eventName); err != nil {
		t.Fatalf("query activity_log row: %v", err)
	}
	if activityCount != 1 {
		t.Fatalf("activity_log grouped row count = %d, want 1", activityCount)
	}
	if actor != "developer" {
		t.Fatalf("activity_log actor = %q, want developer", actor)
	}
	if eventName != "status-changed" {
		t.Fatalf("activity_log event = %q, want status-changed", eventName)
	}

	beforeRebuild := issue
	dbPath := filepath.Join(config.LedgerDir(workspace), "fid.sqlite")
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

func TestTransitionRejectsEdgeMissingFromDBRegistry(t *testing.T) {
	t.Parallel()
	workspace := newIssueWorkspace(t)

	created, err := createTransitionSubtask(t, CreateOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
		Kind:          "task",
		Title:         "Registry edge source",
		Actor:         testActorForRole("developer"), IdempotencyKey: testIdempotencyKey(),
	})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	conn := openIssueDB(t, filepath.Join(config.LedgerDir(workspace), "fid.sqlite"))
	if _, err := conn.Exec(`
		DELETE FROM workflow_transitions
		WHERE workflow_id = ? AND from_state = ? AND to_state = ?
	`, "subtask-default", "todo", "in_progress"); err != nil {
		closeIssueTestDB(t, conn)
		t.Fatalf("delete registry edge: %v", err)
	}
	closeIssueTestDB(t, conn)

	_, err = Transition(TransitionOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
		ID:            created.ID,
		ToState:       "in_progress",
		Actor:         testActorForRole("human"), IdempotencyKey: testIdempotencyKey(),
	})
	if err == nil {
		t.Fatal("Transition() error = nil, want missing registry edge rejection")
	}
	if got, want := err.Error(), `transition issue "`+created.ID+`" from "todo" to "in_progress" is not allowed`; got != want {
		t.Fatalf("Transition() error = %q, want %q", got, want)
	}
}

func TestTransitionRegistryMatchesFrozenV01OrderedStateMatrix(t *testing.T) {
	t.Parallel()
	workspace := newIssueWorkspace(t)
	conn := openIssueDB(t, filepath.Join(config.LedgerDir(workspace), "fid.sqlite"))
	defer closeIssueTestDB(t, conn)

	registry, err := issuestate.Load(conn)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	for workflowID, definition := range issuestate.FrozenDefinitions() {
		for fromState := range definition.States {
			for toState := range definition.States {
				wantRules, wantAllowed := issuestate.RulesForEdge(workflowID, fromState, toState)
				gotRules, gotAllowed := registry.RulesForEdge(workflowID, fromState, toState)
				if gotAllowed != wantAllowed {
					t.Fatalf("registry edge %s %s -> %s allowed = %t, want %t", workflowID, fromState, toState, gotAllowed, wantAllowed)
				}
				if !reflect.DeepEqual(gotRules, wantRules) {
					t.Fatalf("registry edge %s %s -> %s rules = %#v, want %#v", workflowID, fromState, toState, gotRules, wantRules)
				}
			}
		}
	}
}

func TestTransitionUsesDBRuleReferences(t *testing.T) {
	t.Parallel()
	workspace := newIssueWorkspace(t)

	created, err := createTransitionSubtask(t, CreateOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
		Kind:          "task",
		Title:         "Registry rule source",
		Actor:         testActorForRole("developer"), IdempotencyKey: testIdempotencyKey(),
	})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if _, err := Transition(TransitionOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
		ID:            created.ID,
		ToState:       "in_progress",
		Actor:         testActorForRole("human"), IdempotencyKey: testIdempotencyKey(),
	}); err != nil {
		t.Fatalf("Transition(todo->in_progress) error = %v", err)
	}

	conn := openIssueDB(t, filepath.Join(config.LedgerDir(workspace), "fid.sqlite"))
	if _, err := conn.Exec(`
		DELETE FROM workflow_transition_requires
		WHERE workflow_id = ? AND from_state = ? AND to_state = ?
	`, "subtask-default", "in_progress", "done"); err != nil {
		closeIssueTestDB(t, conn)
		t.Fatalf("delete registry rules: %v", err)
	}
	closeIssueTestDB(t, conn)

	if _, err := Transition(TransitionOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
		ID:            created.ID,
		ToState:       "done",
		Actor:         testActorForRole("human"), IdempotencyKey: testIdempotencyKey(),
	}); err != nil {
		t.Fatalf("Transition(in_progress->done) error = %v", err)
	}
}

func TestTransitionUsesDBCategoriesForClaimGating(t *testing.T) {
	t.Parallel()
	workspace := newIssueWorkspace(t)

	created, err := createTransitionSubtask(t, CreateOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
		Kind:          "task",
		Title:         "Registry category source",
		Actor:         testActorForRole("developer"), IdempotencyKey: testIdempotencyKey(),
	})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	conn := openIssueDB(t, filepath.Join(config.LedgerDir(workspace), "fid.sqlite"))
	if _, err := conn.Exec(`
		UPDATE workflow_states
		SET category = ?
		WHERE workflow_id = ? AND state = ?
	`, "unstarted", "subtask-default", "in_progress"); err != nil {
		closeIssueTestDB(t, conn)
		t.Fatalf("update registry category: %v", err)
	}
	closeIssueTestDB(t, conn)

	if _, err := Transition(TransitionOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
		ID:            created.ID,
		ToState:       "in_progress",
		Actor:         testActorForRole("human"), IdempotencyKey: testIdempotencyKey(),
	}); err != nil {
		t.Fatalf("Transition(todo->in_progress) error = %v", err)
	}
}

func TestTransitionRejectsCancelledEdgeMissingFromDBRegistry(t *testing.T) {
	t.Parallel()
	workspace := newIssueWorkspace(t)

	created, err := Create(CreateOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
		Kind:          "task",
		Title:         "Explicit cancellation edge",
		Actor:         testActorForRole("developer"), IdempotencyKey: testIdempotencyKey(),
	})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	conn := openIssueDB(t, filepath.Join(config.LedgerDir(workspace), "fid.sqlite"))
	if _, err := conn.Exec(`
		DELETE FROM workflow_transitions
		WHERE workflow_id = ? AND from_state = ? AND to_state = ?
	`, "task-default", "todo", "cancelled"); err != nil {
		closeIssueTestDB(t, conn)
		t.Fatalf("delete cancelled registry edge: %v", err)
	}
	closeIssueTestDB(t, conn)

	_, err = Transition(TransitionOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
		ID:            created.ID,
		ToState:       "cancelled", Reason: "Test cancellation",
		Actor: "human:mark", IdempotencyKey: testIdempotencyKey(),
	})
	if err == nil {
		t.Fatal("Transition() error = nil, want missing cancelled registry edge rejection")
	}
	if got, want := err.Error(), `transition issue "`+created.ID+`" from "todo" to "cancelled" is not allowed`; got != want {
		t.Fatalf("Transition() error = %q, want %q", got, want)
	}
}

func TestTransitionTreatsReservedGateRuleAsUnsatisfied(t *testing.T) {
	t.Parallel()
	workspace := newIssueWorkspace(t)
	created, err := createTransitionSubtask(t, CreateOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
		Kind:          "task",
		Title:         "Reserved gate seam",
		Actor:         testActorForRole("developer"), IdempotencyKey: testIdempotencyKey(),
	})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	conn := openIssueDB(t, filepath.Join(config.LedgerDir(workspace), "fid.sqlite"))
	if _, err := conn.Exec(`
		INSERT INTO workflow_transition_requires (workflow_id, from_state, to_state, rule, params_json)
		VALUES (?, ?, ?, ?, ?)
	`, "subtask-default", "todo", "in_progress", issuestate.RuleRequiredGatesCurrent, nil); err != nil {
		closeIssueTestDB(t, conn)
		t.Fatalf("insert reserved workflow rule: %v", err)
	}
	closeIssueTestDB(t, conn)

	journalPath := journalPathForWorkspace(t, workspace)
	before, err := os.ReadFile(journalPath)
	if err != nil {
		t.Fatalf("read journal before reserved rule: %v", err)
	}
	_, err = Transition(TransitionOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
		ID:            created.ID,
		ToState:       "in_progress",
		Actor:         testActorForRole("developer"), IdempotencyKey: testIdempotencyKey(),
	})
	if err == nil {
		t.Fatal("Transition() error = nil, want reserved-rule rejection")
	}
	var blocked *BlockedTransitionError
	if !errors.As(err, &blocked) {
		t.Fatalf("Transition() error = %T, want *BlockedTransitionError", err)
	}
	if blocked.CurrentState != "todo" || blocked.RequestedState != "in_progress" || len(blocked.Unmet) != 1 || blocked.Unmet[0].Rule != issuestate.RuleRequiredGatesCurrent {
		t.Fatalf("BlockedTransitionError = %+v, want todo/in_progress/%q", blocked, issuestate.RuleRequiredGatesCurrent)
	}
	after, readErr := os.ReadFile(journalPath)
	if readErr != nil {
		t.Fatalf("read journal after reserved rule: %v", readErr)
	}
	if !reflect.DeepEqual(after, before) {
		t.Fatal("journal changed while reserved gate rule was rejected")
	}
}

func TestTransitionIntoStartedRequiresStart(t *testing.T) {
	t.Parallel()
	workspace := newIssueWorkspace(t)
	created, err := Create(CreateOptions{WorkspaceRoot: workspace, MigrationsDir: testIssueMigrationsDir(t), Kind: "task", Title: "Start requires its command", Actor: "agent:developer", IdempotencyKey: testIdempotencyKey()})
	if err != nil {
		t.Fatal(err)
	}
	before, err := Get(GetOptions{WorkspaceRoot: workspace, MigrationsDir: testIssueMigrationsDir(t), ID: created.ID})
	if err != nil {
		t.Fatal(err)
	}
	path := journalPathForWorkspace(t, workspace)
	lines := countJournalLines(path)
	_, err = Transition(TransitionOptions{WorkspaceRoot: workspace, MigrationsDir: testIssueMigrationsDir(t), ID: created.ID, ToState: "in_progress", Actor: "agent:developer", IdempotencyKey: testIdempotencyKey()})
	var refused *CommandTransitionError
	if !errors.As(err, &refused) || refused.Rule != "use-start" {
		t.Fatalf("transition = %v", err)
	}
	after, err := Get(GetOptions{WorkspaceRoot: workspace, MigrationsDir: testIssueMigrationsDir(t), ID: created.ID})
	if err != nil {
		t.Fatal(err)
	}
	assertIssuesEqual(t, after, before)
	if countJournalLines(path) != lines {
		t.Fatal("refused start appended an Event")
	}
}

func TestTransitionRejectsInvalidEdgeWithoutAppending(t *testing.T) {
	t.Parallel()

	workspace := newIssueWorkspace(t)

	created, err := Create(CreateOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
		Kind:          "task",
		Title:         "Reject invalid transition",
		Actor:         testActorForRole("developer"), IdempotencyKey: testIdempotencyKey(),
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
		Actor:         testActorForRole("tl"), IdempotencyKey: testIdempotencyKey(),
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

func TestTransitionRejectsRetiredDeveloperState(t *testing.T) {
	t.Parallel()
	assertRetiredTaskStates(t, []string{"dev_complete", "failed_gate"})
}

func TestTransitionRejectsRetiredQAStates(t *testing.T) {
	t.Parallel()
	assertRetiredTaskStates(t, []string{"qa_review", "qa_block", "qa_pass"})
}

func assertRetiredTaskStates(t *testing.T, states []string) {
	t.Helper()
	workspace := newIssueWorkspace(t)
	created, err := Create(CreateOptions{WorkspaceRoot: workspace, MigrationsDir: testIssueMigrationsDir(t), Kind: "task", Title: "Retired task states", Actor: "agent:developer", IdempotencyKey: testIdempotencyKey()})
	if err != nil {
		t.Fatal(err)
	}
	seedTransitionAcceptanceCriteria(t, workspace, created.ID, acceptanceCriterionSeed{ID: "AC-retired", Ordinal: 1, Text: "Existing criterion"})
	seedTransitionStatusUpdate(t, workspace, created.ID, "task-default", "todo", "in_progress")
	before, err := Get(GetOptions{WorkspaceRoot: workspace, MigrationsDir: testIssueMigrationsDir(t), ID: created.ID})
	if err != nil {
		t.Fatal(err)
	}
	path := journalPathForWorkspace(t, workspace)
	lines := countJournalLines(path)
	for _, state := range states {
		_, err := Transition(TransitionOptions{WorkspaceRoot: workspace, MigrationsDir: testIssueMigrationsDir(t), ID: created.ID, ToState: state, Actor: "human:mark", IdempotencyKey: testIdempotencyKey()})
		var invalid *InvalidTransitionError
		if !errors.As(err, &invalid) || invalid.RequestedState != state {
			t.Fatalf("retired %s = %v", state, err)
		}
		after, err := Get(GetOptions{WorkspaceRoot: workspace, MigrationsDir: testIssueMigrationsDir(t), ID: created.ID})
		if err != nil {
			t.Fatal(err)
		}
		assertIssuesEqual(t, after, before)
		if countJournalLines(path) != lines {
			t.Fatal("retired state appended an Event")
		}
	}
}

func TestTransitionRejectsRetiredEpicStatesWithOpenOrFinishedChildren(t *testing.T) {
	t.Parallel()
	workspace := newIssueWorkspace(t)
	epic, err := Create(CreateOptions{WorkspaceRoot: workspace, MigrationsDir: testIssueMigrationsDir(t), Kind: "epic", Title: "Epic states are fixed", Scheduling: "prioritised", Actor: "human:mark", IdempotencyKey: testIdempotencyKey()})
	if err != nil {
		t.Fatal(err)
	}
	child, err := Create(CreateOptions{WorkspaceRoot: workspace, MigrationsDir: testIssueMigrationsDir(t), Kind: "task", ParentID: &epic.ID, Title: "Child task", Actor: "human:mark", IdempotencyKey: testIdempotencyKey()})
	if err != nil {
		t.Fatal(err)
	}
	seedTransitionStatusUpdate(t, workspace, epic.ID, "epic-default", "todo", "in_progress")
	for _, childState := range []string{"todo", "done", "cancelled"} {
		if childState != "todo" {
			from := "todo"
			if childState == "cancelled" {
				from = "done"
			}
			seedTransitionStatusUpdate(t, workspace, child.ID, "task-default", from, childState)
		}
		before, err := Get(GetOptions{WorkspaceRoot: workspace, MigrationsDir: testIssueMigrationsDir(t), ID: epic.ID})
		if err != nil {
			t.Fatal(err)
		}
		path := journalPathForWorkspace(t, workspace)
		lines := countJournalLines(path)
		for _, state := range []string{"wrapped", "closed"} {
			_, err := Transition(TransitionOptions{WorkspaceRoot: workspace, MigrationsDir: testIssueMigrationsDir(t), ID: epic.ID, ToState: state, Actor: "human:mark", IdempotencyKey: testIdempotencyKey()})
			var invalid *InvalidTransitionError
			if !errors.As(err, &invalid) {
				t.Fatalf("retired epic %s with child %s = %v", state, childState, err)
			}
		}
		after, err := Get(GetOptions{WorkspaceRoot: workspace, MigrationsDir: testIssueMigrationsDir(t), ID: epic.ID})
		if err != nil {
			t.Fatal(err)
		}
		assertIssuesEqual(t, after, before)
		if countJournalLines(path) != lines {
			t.Fatal("retired epic state appended an Event")
		}
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
		Actor:         testActorForRole("developer"), IdempotencyKey: testIdempotencyKey(),
	})
	if err != nil {
		t.Fatalf("Create(parent) error = %v", err)
	}
	subtask, err := Create(CreateOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
		Kind:          "subtask",
		Title:         "Child subtask",
		Actor:         testActorForRole("developer"), IdempotencyKey: testIdempotencyKey(),
		ParentID: stringPointer(parent.ID),
	})
	if err != nil {
		t.Fatalf("Create(subtask) error = %v", err)
	}
	task, err := Create(CreateOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
		Kind:          "task",
		Title:         "Cancel from todo",
		Actor:         testActorForRole("developer"), IdempotencyKey: testIdempotencyKey(),
	})
	if err != nil {
		t.Fatalf("Create(task) error = %v", err)
	}

	if _, err := Transition(TransitionOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
		ID:            subtask.ID,
		ToState:       "in_progress",
		Actor:         testActorForRole("developer"), IdempotencyKey: testIdempotencyKey(),
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
		ToState:       "cancelled", Reason: "Test cancellation",
		Actor: "human:mark", IdempotencyKey: testIdempotencyKey(),
	}); err != nil {
		t.Fatalf("Transition(task todo->cancelled) error = %v", err)
	}
	if _, err := Transition(TransitionOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
		ID:            subtask.ID,
		ToState:       "cancelled",
		Actor:         testActorForRole("tl"), IdempotencyKey: testIdempotencyKey(),
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
		Actor:         testActorForRole("tl"), IdempotencyKey: testIdempotencyKey(),
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

// Generic transition machinery remains observable through unchanged subtasks.
func createTransitionSubtask(t *testing.T, options CreateOptions) (CreateResult, error) {
	t.Helper()
	parent, err := Create(CreateOptions{WorkspaceRoot: options.WorkspaceRoot, MigrationsDir: options.MigrationsDir, Kind: "task", Title: "Transition parent", Actor: options.Actor, IdempotencyKey: testIdempotencyKey()})
	if err != nil {
		return CreateResult{}, err
	}
	options.Kind = "subtask"
	options.ParentID = &parent.ID
	return Create(options)
}
