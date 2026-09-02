package issues

import (
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/mbell-nzl/fid/internal/journal"
)

func TestTask8IssueEventsUseFixedClockAndAggregateOrder(t *testing.T) {
	workspace := newIssueWorkspace(t)
	fixed := time.Now().UTC()
	clock := func() time.Time { return fixed }

	created, err := Create(CreateOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
		Kind:          "task", Title: "Fixed clock Issue", Actor: "agent:task-8", IdempotencyKey: "fixed-create",
		now: clock,
	})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if _, err := Update(UpdateOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
		ID:            created.ID, Actor: "agent:task-8", IdempotencyKey: "fixed-update", Title: stringPointer("Fixed clock update"), ExpectedRevision: intPointer(1),
		now: clock,
	}); err != nil {
		t.Fatalf("Update() error = %v", err)
	}
	if _, err := Transition(TransitionOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
		ID:            created.ID, ToState: "cancelled", Actor: "agent:task-8", IdempotencyKey: "fixed-transition", ExpectedRevision: intPointer(2),
		now: clock,
	}); err != nil {
		t.Fatalf("Transition() error = %v", err)
	}

	var issueEvents []journal.Event
	for _, event := range readIssueStoredEvents(t, journalPathForWorkspace(t, workspace)) {
		if strings.HasPrefix(event.EventType, "issue.") {
			issueEvents = append(issueEvents, event)
		}
	}
	if len(issueEvents) != 3 {
		t.Fatalf("Issue Event count = %d, want 3", len(issueEvents))
	}
	for index, event := range issueEvents {
		if event.EventSchema != journal.EventSchemaV2 || event.Aggregate.ID != created.ID || event.Aggregate.Type != "issue" || event.Aggregate.Sequence != index+1 {
			t.Fatalf("Issue Event %d = %#v, want Event v2 aggregate sequence %d", index, event, index+1)
		}
		if event.OccurredAt.After(fixed) || !event.OccurredAt.Equal(fixed) {
			t.Fatalf("Issue Event %d occurred_at = %s, want fixed nonfuture %s", index, event.OccurredAt, fixed)
		}
	}
}

func TestTask8LegacyInProcessIssueMutationsReceiveIdempotencyKeys(t *testing.T) {
	workspace := newIssueWorkspace(t)
	created, err := Create(CreateOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
		Kind:          "task",
		Title:         "Legacy in-process Issue",
		Actor:         testActorForRole("developer"), IdempotencyKey: testIdempotencyKey(),
	})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	updated, err := Update(UpdateOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
		ID:            created.ID,
		Actor:         testActorForRole("developer"), IdempotencyKey: testIdempotencyKey(),
		Title:            stringPointer("Legacy in-process update"),
		ExpectedRevision: intPointer(1),
	})
	if err != nil {
		t.Fatalf("Update() error = %v", err)
	}
	transitioned, err := Transition(TransitionOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
		ID:            created.ID,
		ToState:       "cancelled",
		Actor:         testActorForRole("developer"), IdempotencyKey: testIdempotencyKey(),
		ExpectedRevision: intPointer(2),
	})
	if err != nil {
		t.Fatalf("Transition() error = %v", err)
	}
	keys := []string{created.IdempotencyKey, updated.IdempotencyKey, transitioned.IdempotencyKey}
	for index, key := range keys {
		if key == "" {
			t.Fatalf("mutation %d idempotency key is empty", index)
		}
	}
	if keys[0] == keys[1] || keys[1] == keys[2] || keys[0] == keys[2] {
		t.Fatalf("legacy idempotency keys = %v, want distinct values", keys)
	}
}

func TestTask8IssueCreationSourceBudget(t *testing.T) {
	_, sourceFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller(0) failed")
	}
	path := filepath.Join(filepath.Dir(sourceFile), "create.go")
	fileSet := token.NewFileSet()
	file, err := parser.ParseFile(fileSet, path, nil, 0)
	if err != nil {
		t.Fatalf("parse %s: %v", path, err)
	}
	if lineCount := fileSet.File(file.Pos()).LineCount(); lineCount > 400 {
		t.Fatalf("create.go has %d lines, want <= 400", lineCount)
	}
	for _, declaration := range file.Decls {
		function, ok := declaration.(*ast.FuncDecl)
		if !ok || function.Body == nil {
			continue
		}
		lines := fileSet.Position(function.Body.Rbrace).Line - fileSet.Position(function.Pos()).Line + 1
		if lines > 60 {
			t.Fatalf("%s has %d lines, want <= 60", function.Name.Name, lines)
		}
	}
}
