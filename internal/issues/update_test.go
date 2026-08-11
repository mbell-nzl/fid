package issues

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mbell-nzl/fid/internal/journal"
)

func TestUpdateChangesTitleAndPriorityAndPreservesUntouchedFields(t *testing.T) {
	t.Parallel()

	workspace := newIssueWorkspace(t)
	seedIssueLabels(t, workspace, "core")

	created, err := Create(CreateOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
		Kind:          "task",
		Title:         "Implement update command",
		ActorRole:     "developer",
		Priority:      intPointer(2),
		SpecMDPath:    stringPointer("project-docs/specs/06_DESIGN_PLAN_CLI.md"),
		Labels:        []string{"core"},
		OwnerActor:    stringPointer("agent:before"),
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
		t.Fatalf("Get(before update) error = %v", err)
	}

	_, err = Update(UpdateOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
		ID:            created.ID,
		ActorRole:     "developer",
		Title:         stringPointer("Implement issue update"),
		Priority:      intPointer(1),
	})
	if err != nil {
		t.Fatalf("Update() error = %v", err)
	}

	after, err := Get(GetOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
		ID:            created.ID,
	})
	if err != nil {
		t.Fatalf("Get(after update) error = %v", err)
	}

	if after.Title != "Implement issue update" {
		t.Fatalf("Get(after update) title = %q, want %q", after.Title, "Implement issue update")
	}
	if after.Priority == nil || *after.Priority != 1 {
		t.Fatalf("Get(after update) priority = %#v, want 1", after.Priority)
	}
	if !equalStringPointers(after.SpecMDPath, before.SpecMDPath) {
		t.Fatalf("Get(after update) spec_md_path = %#v, want %#v", after.SpecMDPath, before.SpecMDPath)
	}
	if !equalStrings(after.Labels, before.Labels) {
		t.Fatalf("Get(after update) labels = %#v, want %#v", after.Labels, before.Labels)
	}
	if !equalStringPointers(after.OwnerActor, before.OwnerActor) {
		t.Fatalf("Get(after update) owner_actor = %#v, want %#v", after.OwnerActor, before.OwnerActor)
	}
}

func TestUpdateRejectsInvalidChangesWithoutAppending(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		mutate  func(createdID string) UpdateOptions
		wantErr string
	}{
		{
			name: "priority out of range",
			mutate: func(createdID string) UpdateOptions {
				return UpdateOptions{
					ID:        createdID,
					ActorRole: "developer",
					Priority:  intPointer(7),
				}
			},
			wantErr: "update issue: priority must be between 0 and 5",
		},
		{
			name: "blank title",
			mutate: func(createdID string) UpdateOptions {
				return UpdateOptions{
					ID:        createdID,
					ActorRole: "developer",
					Title:     stringPointer("   "),
				}
			},
			wantErr: "update issue: title must not be blank",
		},
		{
			name: "unknown label",
			mutate: func(createdID string) UpdateOptions {
				labels := []string{"core", "missing"}
				return UpdateOptions{
					ID:        createdID,
					ActorRole: "developer",
					Labels:    &labels,
				}
			},
			wantErr: `update issue: label "missing" does not exist`,
		},
		{
			name: "unknown owner role",
			mutate: func(createdID string) UpdateOptions {
				return UpdateOptions{
					ID:        createdID,
					ActorRole: "developer",
					OwnerRole: stringPointer("intern"),
				}
			},
			wantErr: `role "intern" not in project profile`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			workspace := newIssueWorkspace(t)
			seedIssueLabels(t, workspace, "core")

			created, err := Create(CreateOptions{
				WorkspaceRoot: workspace,
				MigrationsDir: testIssueMigrationsDir(t),
				Kind:          "task",
				Title:         "Reject invalid update",
				ActorRole:     "developer",
				Priority:      intPointer(2),
				Labels:        []string{"core"},
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
				t.Fatalf("Get(before invalid update) error = %v", err)
			}

			journalPath := journalPathForWorkspace(t, workspace)
			linesBefore := countJournalLines(journalPath)

			options := tt.mutate(created.ID)
			options.WorkspaceRoot = workspace
			options.MigrationsDir = testIssueMigrationsDir(t)

			_, err = Update(options)
			if err == nil {
				t.Fatal("Update() error = nil, want rejection")
			}
			if err.Error() != tt.wantErr {
				t.Fatalf("Update() error = %q, want %q", err.Error(), tt.wantErr)
			}

			after, err := Get(GetOptions{
				WorkspaceRoot: workspace,
				MigrationsDir: testIssueMigrationsDir(t),
				ID:            created.ID,
			})
			if err != nil {
				t.Fatalf("Get(after invalid update) error = %v", err)
			}
			assertIssuesEqual(t, after, before)

			if got := countJournalLines(journalPath); got != linesBefore {
				t.Fatalf("journal lines after rejected update = %d, want unchanged %d", got, linesBefore)
			}
		})
	}
}

func TestUpdateAppendsOneEventProjectsActivityAndRebuilds(t *testing.T) {
	t.Parallel()

	workspace := newIssueWorkspace(t)
	seedIssueLabels(t, workspace, "core", "api", "legacy")

	created, err := Create(CreateOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
		Kind:          "task",
		Title:         "Initial update target",
		ActorRole:     "developer",
		Priority:      intPointer(2),
		SpecMDPath:    stringPointer("project-docs/specs/02_TRD.md"),
		Labels:        []string{"core", "legacy"},
		OwnerActor:    stringPointer("agent:before"),
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
		t.Fatalf("Get(before update) error = %v", err)
	}

	journalPath := journalPathForWorkspace(t, workspace)
	linesBefore := countJournalLines(journalPath)
	labels := []string{"api", "core"}
	result, err := Update(UpdateOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
		ID:            created.ID,
		ActorRole:     "tl",
		Title:         stringPointer("Updated issue target"),
		Priority:      intPointer(1),
		Labels:        &labels,
		SpecMDPath:    stringPointer("project-docs/specs/06_DESIGN_PLAN_CLI.md"),
		OwnerRole:     stringPointer("qa_lead"),
		OwnerActor:    stringPointer("human:mark"),
	})
	if err != nil {
		t.Fatalf("Update() error = %v", err)
	}

	after, err := Get(GetOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
		ID:            created.ID,
	})
	if err != nil {
		t.Fatalf("Get(after update) error = %v", err)
	}

	if result.ID != created.ID {
		t.Fatalf("UpdateResult.ID = %q, want %q", result.ID, created.ID)
	}
	if after.Title != "Updated issue target" {
		t.Fatalf("Get(after update) title = %q, want Updated issue target", after.Title)
	}
	if after.Priority == nil || *after.Priority != 1 {
		t.Fatalf("Get(after update) priority = %#v, want 1", after.Priority)
	}
	if after.SpecMDPath == nil || *after.SpecMDPath != "project-docs/specs/06_DESIGN_PLAN_CLI.md" {
		t.Fatalf("Get(after update) spec_md_path = %#v", after.SpecMDPath)
	}
	if after.OwnerRole == nil || *after.OwnerRole != "qa_lead" {
		t.Fatalf("Get(after update) owner_role = %#v, want qa_lead", after.OwnerRole)
	}
	if after.OwnerActor == nil || *after.OwnerActor != "human:mark" {
		t.Fatalf("Get(after update) owner_actor = %#v, want human:mark", after.OwnerActor)
	}
	if want := []string{"api", "core"}; !equalStrings(after.Labels, want) {
		t.Fatalf("Get(after update) labels = %#v, want %#v", after.Labels, want)
	}
	if after.UpdatedAt == nil || before.UpdatedAt == nil || *after.UpdatedAt == *before.UpdatedAt {
		t.Fatalf("updated_at before/after = %#v / %#v, want advance", before.UpdatedAt, after.UpdatedAt)
	}

	if got := countJournalLines(journalPath); got != linesBefore+1 {
		t.Fatalf("journal lines after update = %d, want %d", got, linesBefore+1)
	}

	stored := readIssueStoredEvents(t, journalPath)
	last := stored[len(stored)-1]
	if last.EventType != "issue.updated" {
		t.Fatalf("last event type = %q, want issue.updated", last.EventType)
	}

	var (
		updateEffect   *journal.Effect
		activityEffect *journal.Effect
		labelInserts   int
		labelDeletes   int
	)
	for index := range last.Effects {
		effect := &last.Effects[index]
		switch {
		case effect.Table == "issues" && effect.Op == "update":
			updateEffect = effect
		case effect.Table == "activity_log" && effect.Op == "activity":
			activityEffect = effect
		case effect.Table == "issue_labels" && effect.Op == "insert":
			labelInserts++
		case effect.Table == "issue_labels" && effect.Op == "delete":
			labelDeletes++
		}
	}
	if updateEffect == nil {
		t.Fatal("missing issues update effect")
	}
	if activityEffect == nil {
		t.Fatal("missing activity effect")
	}
	if labelInserts != 1 || labelDeletes != 1 {
		t.Fatalf("label delta effects = %d inserts / %d deletes, want 1 / 1", labelInserts, labelDeletes)
	}
	if updateEffect.Preconditions["title"] != "Initial update target" {
		t.Fatalf("title precondition = %#v, want Initial update target", updateEffect.Preconditions["title"])
	}
	if updateEffect.Preconditions["priority"] != float64(2) {
		t.Fatalf("priority precondition = %#v, want 2", updateEffect.Preconditions["priority"])
	}
	if updateEffect.Preconditions["owner_role"] != "developer" {
		t.Fatalf("owner_role precondition = %#v, want developer", updateEffect.Preconditions["owner_role"])
	}
	if updateEffect.Preconditions["owner_actor"] != "agent:before" {
		t.Fatalf("owner_actor precondition = %#v, want agent:before", updateEffect.Preconditions["owner_actor"])
	}
	if updateEffect.Preconditions["updated_at"] != *before.UpdatedAt {
		t.Fatalf("updated_at precondition = %#v, want %#v", updateEffect.Preconditions["updated_at"], before.UpdatedAt)
	}
	if updateEffect.Values["updated_at"] != *after.UpdatedAt {
		t.Fatalf("updated_at value = %#v, want %#v", updateEffect.Values["updated_at"], after.UpdatedAt)
	}
	if activityEffect.Values["journal_event_id"] != last.EventID {
		t.Fatalf("activity journal_event_id = %#v, want %q", activityEffect.Values["journal_event_id"], last.EventID)
	}

	conn := openIssueDB(t, filepath.Join(workspace, ".fid", "fid.sqlite"))
	defer closeIssueTestDB(t, conn)

	var (
		activityCount int
		detailsRaw    string
		actor         string
		eventName     string
	)
	if err := conn.QueryRow(`
		SELECT COUNT(*), COALESCE(details, ''), COALESCE(actor, ''), event
		FROM activity_log
		WHERE journal_event_id = ?
		GROUP BY details, actor, event
	`, last.EventID).Scan(&activityCount, &detailsRaw, &actor, &eventName); err != nil {
		t.Fatalf("query activity_log row: %v", err)
	}
	if activityCount != 1 {
		t.Fatalf("activity_log grouped row count = %d, want 1", activityCount)
	}
	if actor != "tl" {
		t.Fatalf("activity_log actor = %q, want tl", actor)
	}
	if eventName != "issue_updated" {
		t.Fatalf("activity_log event = %q, want issue_updated", eventName)
	}
	for _, want := range []string{`"field":"owner_role"`, `"from":"developer"`, `"to":"qa_lead"`, `"field":"labels"`} {
		if !strings.Contains(detailsRaw, want) {
			t.Fatalf("activity_log details = %q, want substring %q", detailsRaw, want)
		}
	}

	beforeRebuild := after
	dbPath := filepath.Join(workspace, ".fid", "fid.sqlite")
	if err := os.Remove(dbPath); err != nil {
		t.Fatalf("Remove(%q) error = %v", dbPath, err)
	}

	if _, err := journal.RebuildWorkspace(workspace, testIssueMigrationsDir(t), "developer"); err != nil {
		t.Fatalf("RebuildWorkspace() error = %v", err)
	}

	rebuilt, err := Get(GetOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
		ID:            created.ID,
	})
	if err != nil {
		t.Fatalf("Get(after rebuild) error = %v", err)
	}
	assertIssuesEqual(t, rebuilt, beforeRebuild)
}

func TestUpdateRejectsNoOpWithoutAppending(t *testing.T) {
	t.Parallel()

	workspace := newIssueWorkspace(t)

	created, err := Create(CreateOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
		Kind:          "task",
		Title:         "No-op target",
		ActorRole:     "developer",
	})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	journalPath := journalPathForWorkspace(t, workspace)
	linesBefore := countJournalLines(journalPath)

	_, err = Update(UpdateOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
		ID:            created.ID,
		ActorRole:     "developer",
		Title:         stringPointer("No-op target"),
	})
	if err == nil {
		t.Fatal("Update() error = nil, want no-op rejection")
	}
	if err.Error() != "update issue: at least one field must change" {
		t.Fatalf("Update() error = %q, want no-op rejection", err.Error())
	}
	if got := countJournalLines(journalPath); got != linesBefore {
		t.Fatalf("journal lines after no-op update = %d, want %d", got, linesBefore)
	}
}
