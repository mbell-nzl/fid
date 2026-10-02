package maintenance

import "testing"

func TestCheckWorkflowAssignmentsFlagsIncompatibleAppliesTo(t *testing.T) {
	t.Parallel()
	fixture := newMaintenanceFixture(t)
	seedJournalEvent(t, fixture.journalPath, "EVT-20260715-000001")
	if _, err := fixture.db.Exec(`
		INSERT INTO workflow_assignments (scope, kind, workflow_id)
		VALUES ('kind', 'epic', 'task-default')
	`); err != nil {
		t.Fatalf("seed incompatible assignment: %v", err)
	}

	result, err := Check(CheckOptions{DB: fixture.db, JournalDir: fixture.journalDir})
	if err != nil {
		t.Fatalf("Check() error = %v", err)
	}
	if !hasFinding(result.Errors, "workflow_assignment", "kind:epic") {
		t.Fatalf("Check().Errors = %#v, want workflow_assignment finding for kind:epic", result.Errors)
	}
}

func TestCheckWorkflowAssignmentsAcceptsCompatibleKindAndProjectScopes(t *testing.T) {
	t.Parallel()
	fixture := newMaintenanceFixture(t)
	seedJournalEvent(t, fixture.journalPath, "EVT-20260715-000001")
	if _, err := fixture.db.Exec(`
		INSERT INTO workflows (id, name, applies_to, active)
		VALUES ('any-workflow', 'any-workflow', NULL, 1)
	`); err != nil {
		t.Fatalf("seed any workflow: %v", err)
	}
	if _, err := fixture.db.Exec(`
		INSERT INTO workflow_assignments (scope, kind, workflow_id)
		VALUES ('kind', 'task', 'task-default'), ('project', NULL, 'any-workflow')
	`); err != nil {
		t.Fatalf("seed compatible assignments: %v", err)
	}

	result, err := Check(CheckOptions{DB: fixture.db, JournalDir: fixture.journalDir})
	if err != nil {
		t.Fatalf("Check() error = %v", err)
	}
	if hasFinding(result.Errors, "workflow_assignment", "kind:task") || hasFinding(result.Errors, "workflow_assignment", "project") {
		t.Fatalf("Check().Errors = %#v, want no workflow_assignment findings", result.Errors)
	}
}
