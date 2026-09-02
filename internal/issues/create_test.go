package issues

import (
	"errors"
	"reflect"
	"testing"
)

func TestCreateAllowsLegalHierarchy(t *testing.T) {
	workspace := newIssueWorkspace(t)

	epic := mustCreateHierarchyIssue(t, workspace, CreateOptions{
		Kind:  "epic",
		Title: "Root epic",
		Actor: testActorForRole("tl"), IdempotencyKey: testIdempotencyKey(),
		Scheduling: "sequential",
	})
	mustCreateHierarchyIssue(t, workspace, CreateOptions{
		Kind:  "task",
		Title: "Root task",
		Actor: testActorForRole("developer"), IdempotencyKey: testIdempotencyKey(),
	})
	task := mustCreateHierarchyIssue(t, workspace, CreateOptions{
		Kind:  "task",
		Title: "Epic task",
		Actor: testActorForRole("developer"), IdempotencyKey: testIdempotencyKey(),
		ParentID: stringPointer(epic.ID),
		Sequence: intPointer(1),
	})
	mustCreateHierarchyIssue(t, workspace, CreateOptions{
		Kind:  "subtask",
		Title: "Task subtask",
		Actor: testActorForRole("developer"), IdempotencyKey: testIdempotencyKey(),
		ParentID: stringPointer(task.ID),
	})
}

func TestCreateRejectsIllegalHierarchyBeforeAppend(t *testing.T) {
	workspace := newIssueWorkspace(t)

	epic := mustCreateHierarchyIssue(t, workspace, CreateOptions{
		Kind:  "epic",
		Title: "Hierarchy epic",
		Actor: testActorForRole("tl"), IdempotencyKey: testIdempotencyKey(),
		Scheduling: "prioritised",
	})
	rootTask := mustCreateHierarchyIssue(t, workspace, CreateOptions{
		Kind:  "task",
		Title: "Root task",
		Actor: testActorForRole("developer"), IdempotencyKey: testIdempotencyKey(),
	})
	childTask := mustCreateHierarchyIssue(t, workspace, CreateOptions{
		Kind:  "task",
		Title: "Epic child task",
		Actor: testActorForRole("developer"), IdempotencyKey: testIdempotencyKey(),
		ParentID: stringPointer(epic.ID),
	})
	childSubtask := mustCreateHierarchyIssue(t, workspace, CreateOptions{
		Kind:  "subtask",
		Title: "Task child subtask",
		Actor: testActorForRole("developer"), IdempotencyKey: testIdempotencyKey(),
		ParentID: stringPointer(childTask.ID),
	})

	journalPath := journalPathForWorkspace(t, workspace)
	linesBefore := countJournalLines(journalPath)

	tests := []struct {
		name          string
		options       CreateOptions
		wantHierarchy *HierarchyInvalidError
		wantNotFound  string
	}{
		{
			name: "epic with parent",
			options: CreateOptions{
				Kind:  "epic",
				Title: "Nested epic",
				Actor: testActorForRole("tl"), IdempotencyKey: testIdempotencyKey(),
				ParentID:   stringPointer(rootTask.ID),
				Scheduling: "sequential",
			},
			wantHierarchy: &HierarchyInvalidError{
				Kind:       "epic",
				ParentID:   stringPointer(rootTask.ID),
				ParentKind: stringPointer("task"),
				Rule:       "epics must not have a parent",
			},
		},
		{
			name: "subtask under epic",
			options: CreateOptions{
				Kind:  "subtask",
				Title: "Wrong subtask parent",
				Actor: testActorForRole("developer"), IdempotencyKey: testIdempotencyKey(),
				ParentID: stringPointer(epic.ID),
			},
			wantHierarchy: &HierarchyInvalidError{
				Kind:       "subtask",
				ParentID:   stringPointer(epic.ID),
				ParentKind: stringPointer("epic"),
				Rule:       "subtasks must have a task parent",
			},
		},
		{
			name: "subtask without parent",
			options: CreateOptions{
				Kind:  "subtask",
				Title: "Orphan subtask",
				Actor: testActorForRole("developer"), IdempotencyKey: testIdempotencyKey(),
			},
			wantHierarchy: &HierarchyInvalidError{
				Kind: "subtask",
				Rule: "subtasks must have a task parent",
			},
		},
		{
			name: "task under subtask",
			options: CreateOptions{
				Kind:  "task",
				Title: "Too-deep task",
				Actor: testActorForRole("developer"), IdempotencyKey: testIdempotencyKey(),
				ParentID: stringPointer(childSubtask.ID),
			},
			wantHierarchy: &HierarchyInvalidError{
				Kind:       "task",
				ParentID:   stringPointer(childSubtask.ID),
				ParentKind: stringPointer("subtask"),
				Rule:       "tasks must have an epic parent or no parent",
			},
		},
		{
			name: "missing parent",
			options: CreateOptions{
				Kind:  "subtask",
				Title: "Missing parent subtask",
				Actor: testActorForRole("developer"), IdempotencyKey: testIdempotencyKey(),
				ParentID: stringPointer("FID-missing-parent"),
			},
			wantNotFound: "FID-missing-parent",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Create(withHierarchyWorkspace(t, workspace, tt.options))
			if err == nil {
				t.Fatal("Create() error = nil, want hierarchy rejection")
			}

			if tt.wantNotFound != "" {
				var notFound *NotFoundError
				if !errors.As(err, &notFound) {
					t.Fatalf("Create() error = %T, want *NotFoundError", err)
				}
				if notFound.IssueID != tt.wantNotFound {
					t.Fatalf("NotFoundError.IssueID = %q, want %q", notFound.IssueID, tt.wantNotFound)
				}
			} else {
				var hierarchy *HierarchyInvalidError
				if !errors.As(err, &hierarchy) {
					t.Fatalf("Create() error = %T, want *HierarchyInvalidError", err)
				}
				if !reflect.DeepEqual(hierarchy, tt.wantHierarchy) {
					t.Fatalf("HierarchyInvalidError = %#v, want %#v", hierarchy, tt.wantHierarchy)
				}
			}

			if got := countJournalLines(journalPath); got != linesBefore {
				t.Fatalf("journal lines after rejected create = %d, want unchanged %d", got, linesBefore)
			}
		})
	}
}

func TestIssueMutationRequestsCannotChangeParent(t *testing.T) {
	for _, request := range []any{UpdateOptions{}, TransitionOptions{}} {
		requestType := reflect.TypeOf(request)
		if _, exists := requestType.FieldByName("ParentID"); exists {
			t.Fatalf("%s exposes ParentID after Issue creation", requestType.Name())
		}
	}
}

func mustCreateHierarchyIssue(t *testing.T, workspace string, options CreateOptions) CreateResult {
	t.Helper()
	created, err := Create(withHierarchyWorkspace(t, workspace, options))
	if err != nil {
		t.Fatalf("Create(%q) error = %v", options.Title, err)
	}
	return created
}

func withHierarchyWorkspace(t *testing.T, workspace string, options CreateOptions) CreateOptions {
	t.Helper()
	options.WorkspaceRoot = workspace
	options.MigrationsDir = testIssueMigrationsDir(t)
	return options
}
