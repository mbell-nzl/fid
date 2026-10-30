package maintenance

import (
	"database/sql"
	"encoding/json"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/mbell-nzl/fid/internal/journal"
)

func TestCheckErrorsReportsAllSeededStructuralViolationFamilies(t *testing.T) {
	t.Parallel()

	fixture := newStructuralFixture(t)
	beforeJournal := fixture.journalBytes(t)
	beforeChanges := fixture.changeCount(t)

	got, err := Check(CheckOptions{
		DB:         fixture.db,
		JournalDir: fixture.journalDir,
	})
	if err != nil {
		t.Fatalf("Check() error = %v", err)
	}

	wantChecks := []string{
		"activity_journal_events",
		"activity_journal_events",
		"activity_refs",
		"activity_refs",
		"activity_refs",
		"dependent_row_kind",
		"dependent_row_kind",
		"dependent_row_kind",
		"foreign_key_integrity",
		"foreign_key_integrity",
		"foreign_key_integrity",
		"foreign_key_integrity",
		"foreign_key_integrity",
		"kind_hierarchy",
		"kind_hierarchy",
		"kind_hierarchy",
		"link_refs",
		"link_relation_compatibility",
		"link_scope",
		"link_scope",
		"sequence_scheduling",
		"sequence_scheduling",
		"sequence_scheduling",
		"sequence_scheduling",
		"workflow_state",
		"workflow_state",
	}
	if gotChecks := findingChecks(got.Errors); !reflect.DeepEqual(gotChecks, wantChecks) {
		t.Fatalf("Check().Errors checks = %#v, want %#v", gotChecks, wantChecks)
	}
	assertFindingsOrdered(t, got.Errors)
	assertActionableFindings(t, got.Errors)

	again, err := Check(CheckOptions{DB: fixture.db, JournalDir: fixture.journalDir})
	if err != nil {
		t.Fatalf("second Check() error = %v", err)
	}
	if !reflect.DeepEqual(again, got) {
		t.Fatalf("second Check() = %#v, want %#v", again, got)
	}
	if afterJournal := fixture.journalBytes(t); !reflect.DeepEqual(afterJournal, beforeJournal) {
		t.Fatalf("journal bytes changed after Check()")
	}
	if afterChanges := fixture.changeCount(t); afterChanges != beforeChanges {
		t.Fatalf("database changes after Check() = %d, want %d", afterChanges, beforeChanges)
	}
}

func assertActionableFindings(t *testing.T, findings []Finding) {
	t.Helper()

	foreignKeys := 0
	issueLabelForeignKeys := 0
	issueWorkflowForeignKeys := 0
	for _, finding := range findings {
		if finding.EntityID == "" || finding.Message == "" {
			t.Fatalf("finding = %#v, want entity ID and message", finding)
		}
		if finding.Check == "foreign_key_integrity" {
			foreignKeys++
			switch {
			case strings.HasPrefix(finding.EntityID, "issue_labels:") && strings.Contains(finding.Message, "references missing"):
				issueLabelForeignKeys++
			case strings.HasPrefix(finding.EntityID, "issues:") && strings.Contains(finding.Message, "references missing workflow"):
				issueWorkflowForeignKeys++
			default:
				t.Fatalf("foreign key finding = %#v, want seeded issue-label or workflow violation", finding)
			}
		}
	}
	if foreignKeys != 5 || issueLabelForeignKeys != 2 || issueWorkflowForeignKeys != 3 {
		t.Fatalf("foreign key findings = total %d, issue_labels %d, issues workflow %d; want 5, 2, 3", foreignKeys, issueLabelForeignKeys, issueWorkflowForeignKeys)
	}
}

func TestCheckAcceptsPopulatedGenericReferencesWithoutStructuralErrors(t *testing.T) {
	t.Parallel()
	fixture := newValidStructuralFixture(t)
	got, err := Check(CheckOptions{DB: fixture.db, JournalDir: fixture.journalDir})
	if err != nil {
		t.Fatalf("Check() error = %v", err)
	}
	if len(got.Errors) != 0 {
		t.Fatalf("Check().Errors = %#v, want no errors for valid generic references", got.Errors)
	}
}

func TestCheckWorkflowStatesReadsProjectedWorkflowTables(t *testing.T) {
	t.Parallel()

	fixture := newMaintenanceFixture(t)
	seedJournalEvent(t, fixture.journalPath, "EVT-20260715-000001")
	seedIssue(t, fixture.db, "FID-registry-audit", "task", nil, nil, "task-default", "todo")

	setFixturePragmas(t, fixture.db, false, false)
	if _, err := fixture.db.Exec(`
		DELETE FROM workflow_states
		WHERE workflow_id = ? AND state = ?
	`, "task-default", "todo"); err != nil {
		t.Fatalf("delete projected workflow state: %v", err)
	}
	setFixturePragmas(t, fixture.db, true, false)

	result, err := Check(CheckOptions{DB: fixture.db, JournalDir: fixture.journalDir})
	if err != nil {
		t.Fatalf("Check() error = %v", err)
	}
	if !hasFinding(result.Errors, "workflow_state", "FID-registry-audit") {
		t.Fatalf("Check().Errors = %#v, want workflow_state finding for projected-table mismatch", result.Errors)
	}
}

func TestMaintenanceReportsBrokenRecoveryLinkWithoutMutation(t *testing.T) {
	t.Parallel()
	fixture := newMaintenanceFixture(t)
	seedIssue(t, fixture.db, "FID-recovery-link", "task", nil, nil, "task-default", "todo")
	seedIssue(t, fixture.db, "FID-final-schema", "task", nil, nil, "task-default", "todo")
	seedTask21CorruptJournal(t, fixture.journalPath)
	seedTask21CorruptProjection(t, fixture.db)
	if _, err := fixture.db.Exec(`
		INSERT INTO run_observations (id, run_id, kind, actor_id, observed_at, observation_json, created_event_id)
		VALUES ('OBS-recovery-link', 'RUN-1', 'liveness', 'actor-maintenance', '2026-08-22T00:00:00Z', '{}', 'EVT-observation')
	`); err != nil {
		t.Fatalf("seed recovery observation: %v", err)
	}
	setFixturePragmas(t, fixture.db, false, false)
	if _, err := fixture.db.Exec(`
		INSERT INTO run_recoveries (id, prior_run_id, successor_run_id, prior_claim_id, successor_claim_id, observation_id, reconciliation_json, adopted_paths_json, authority_json, reason, created_event_id)
		VALUES ('REC-broken-link', 'RUN-1', 'RUN-missing', 'CLM-1', 'CLM-missing', 'OBS-recovery-link', '{}', '[]', '{}', 'test', 'EVT-recovery')
	`); err != nil {
		t.Fatalf("seed broken recovery link: %v", err)
	}
	setFixturePragmas(t, fixture.db, true, false)
	beforeJournal := fixture.journalBytes(t)
	beforeChanges := fixture.changeCount(t)

	result, err := Check(CheckOptions{DB: fixture.db, JournalDir: fixture.journalDir})
	if err != nil {
		t.Fatalf("Check() error = %v", err)
	}
	if !hasFinding(result.Errors, "recovery-link-integrity", "REC-broken-link") {
		t.Fatalf("Check().Errors = %#v, want broken recovery link finding", result.Errors)
	}
	if afterJournal := fixture.journalBytes(t); !reflect.DeepEqual(afterJournal, beforeJournal) {
		t.Fatal("maintenance Check changed journal bytes")
	}
	if afterChanges := fixture.changeCount(t); afterChanges != beforeChanges {
		t.Fatalf("maintenance Check changed projection rows = %d, want %d", afterChanges, beforeChanges)
	}
}

func TestMaintenanceRecoveryIntegrityMatrixDoesNotMutate(t *testing.T) {
	t.Parallel()
	t.Run("healthy controls", func(t *testing.T) {
		fixture := newMaintenanceFixture(t)
		beforeChanges := fixture.changeCount(t)
		result, err := Check(CheckOptions{DB: fixture.db, JournalDir: fixture.journalDir, Now: time.Date(2026, time.September, 8, 0, 0, 0, 0, time.UTC)})
		if err != nil {
			t.Fatalf("Check() error = %v", err)
		}
		if len(result.Errors) != 0 || len(result.Warnings) != 0 {
			t.Fatalf("Check() = %#v, want clean healthy controls", result)
		}
		if afterChanges := fixture.changeCount(t); afterChanges != beforeChanges {
			t.Fatalf("maintenance Check changed healthy projection rows = %d, want %d", afterChanges, beforeChanges)
		}
	})

	t.Run("completed complete-forward recovery later superseded", func(t *testing.T) {
		fixture := newMaintenanceFixture(t)
		if _, err := fixture.db.Exec(`INSERT INTO kit_activations (id, kit_identity, backend_identity, phase, status, point_of_no_return, direction, created_at, updated_at) VALUES ('ACT-superseded', 'kit', 'fid', 22, 'superseded', 1, 'activate', '2026-09-08T00:00:00Z', '2026-09-08T00:00:00Z')`); err != nil {
			t.Fatalf("seed superseded activation: %v", err)
		}
		if _, err := fixture.db.Exec(`INSERT INTO activation_recoveries (activation_id, direction, actor_attribution, reason, selected_event_id, completed_event_id) VALUES ('ACT-superseded', 'complete-forward', 'actor-maintenance', 'test', 'EVT-selected', 'EVT-completed')`); err != nil {
			t.Fatalf("seed completed activation recovery: %v", err)
		}
		beforeChanges := fixture.changeCount(t)
		result, err := Check(CheckOptions{DB: fixture.db, JournalDir: fixture.journalDir})
		if err != nil {
			t.Fatalf("Check() error = %v", err)
		}
		if hasFinding(result.Errors, "activation-recovery-integrity", "ACT-superseded") {
			t.Fatalf("Check().Errors = %#v, want later supersession accepted", result.Errors)
		}
		if afterChanges := fixture.changeCount(t); afterChanges != beforeChanges {
			t.Fatalf("maintenance Check changed superseded activation rows = %d, want %d", afterChanges, beforeChanges)
		}
	})

	t.Run("recovery corruption classes", func(t *testing.T) {
		fixture := newMaintenanceFixture(t)
		seedIssue(t, fixture.db, "FID-final-schema", "task", nil, nil, "task-default", "todo")
		seedTask21CorruptJournal(t, fixture.journalPath)
		seedTask21CorruptProjection(t, fixture.db)
		if _, err := fixture.db.Exec(`INSERT INTO envelope_budgets (envelope_id, revision, name, budget_json) VALUES ('ENV-1', 1, 'attempts', '{}')`); err != nil {
			t.Fatalf("seed budget definition: %v", err)
		}
		if _, err := fixture.db.Exec(`INSERT INTO budget_samples (id, envelope_id, envelope_revision, budget_name, run_id, actor_id, source, source_sample_id, sample_kind, previous_sample_id, value, measurement_ref_json, observed_at, created_event_id) VALUES ('BUD-missing-binding', 'ENV-1', 1, 'attempts', 'RUN-1', 'actor-maintenance', 'fid', 'sample-1', 'delta', NULL, 1, NULL, '2026-08-22T00:00:00Z', 'EVT-budget')`); err != nil {
			t.Fatalf("seed budget sample: %v", err)
		}
		if _, err := fixture.db.Exec(`INSERT INTO activation_recoveries (activation_id, direction, actor_attribution, reason, selected_event_id, completed_event_id) VALUES ('ACT-1', 'complete-forward', 'actor-maintenance', 'test', 'EVT-activation', NULL)`); err != nil {
			t.Fatalf("seed incomplete activation recovery: %v", err)
		}
		if _, err := fixture.db.Exec(`INSERT INTO run_observations (id, run_id, kind, actor_id, observed_at, observation_json, created_event_id) VALUES ('OBS-matrix', 'RUN-1', 'liveness', 'actor-maintenance', '2026-08-22T00:00:00Z', '{}', 'EVT-observation-matrix')`); err != nil {
			t.Fatalf("seed recovery observation: %v", err)
		}
		setFixturePragmas(t, fixture.db, false, false)
		if _, err := fixture.db.Exec(`INSERT INTO run_recoveries (id, prior_run_id, successor_run_id, prior_claim_id, successor_claim_id, observation_id, reconciliation_json, adopted_paths_json, authority_json, reason, created_event_id) VALUES ('REC-matrix', 'RUN-1', 'RUN-missing', 'CLM-1', 'CLM-missing', 'OBS-matrix', '{}', '[]', '{}', 'test', 'EVT-recovery-matrix')`); err != nil {
			t.Fatalf("seed broken recovery link: %v", err)
		}
		setFixturePragmas(t, fixture.db, true, false)
		beforeJournal := fixture.journalBytes(t)
		beforeChanges := fixture.changeCount(t)

		result, err := Check(CheckOptions{DB: fixture.db, JournalDir: fixture.journalDir, Now: time.Date(2026, time.September, 8, 0, 0, 0, 0, time.UTC)})
		if err != nil {
			t.Fatalf("Check() error = %v", err)
		}
		for _, expected := range []struct{ check, entityID string }{
			{"recovery-link-integrity", "REC-matrix"},
			{"attempt-receipt-integrity", "ATT-1"},
			{"budget-binding-integrity", "BUD-missing-binding"},
			{"activation-recovery-integrity", "ACT-1"},
			{"artifact-integrity", "ART-1"},
		} {
			if !hasFinding(result.Errors, expected.check, expected.entityID) {
				t.Fatalf("Check().Errors = %#v, want %s for %s", result.Errors, expected.check, expected.entityID)
			}
		}
		if !hasFinding(result.Warnings, "stale_claims", "FID-final-schema") {
			t.Fatalf("Check().Warnings = %#v, want stale Claims warning", result.Warnings)
		}
		if afterJournal := fixture.journalBytes(t); !reflect.DeepEqual(afterJournal, beforeJournal) {
			t.Fatal("maintenance Check changed journal bytes")
		}
		if afterChanges := fixture.changeCount(t); afterChanges != beforeChanges {
			t.Fatalf("maintenance Check changed projection rows = %d, want %d", afterChanges, beforeChanges)
		}
	})
}

func TestCheckReportsFinalSchemaIntegrityClassesWithoutRepair(t *testing.T) {
	t.Parallel()
	fixture := newMaintenanceFixture(t)
	seedIssue(t, fixture.db, "FID-final-schema", "task", nil, nil, "task-default", "todo")
	seedTask21CorruptJournal(t, fixture.journalPath)
	seedTask21CorruptProjection(t, fixture.db)
	beforeJournal := fixture.journalBytes(t)
	beforeChanges := fixture.changeCount(t)

	result, err := Check(CheckOptions{DB: fixture.db, JournalDir: fixture.journalDir})
	if err != nil {
		t.Fatalf("Check() error = %v", err)
	}

	want := []Finding{
		{Check: "aggregate-sequence", EntityID: "label:LABEL-maintenance", Message: "aggregate sequence 3 follows 1, want 2"},
		{Check: "artifact-integrity", EntityID: "ART-1", Message: "Artifact is unavailable"},
		{Check: "attempt-receipt-integrity", EntityID: "ATT-1", Message: "Attempt has no terminal Receipt"},
		{Check: "event-binding", EntityID: "task21:binding", Message: "idempotency result does not bind its journal Event"},
		{Check: "event-sequence", EntityID: "EVT-20260822-000001", Message: "Event ID occurs more than once"},
		{Check: "removed-state-residue", EntityID: "roles", Message: "removed table is present"},
	}
	if !reflect.DeepEqual(result.Errors, want) {
		t.Fatalf("Check().Errors = %#v, want %#v", result.Errors, want)
	}
	if afterJournal := fixture.journalBytes(t); !reflect.DeepEqual(afterJournal, beforeJournal) {
		t.Fatal("maintenance Check changed journal bytes")
	}
	if afterChanges := fixture.changeCount(t); afterChanges != beforeChanges {
		t.Fatalf("maintenance Check changed projection rows = %d, want %d", afterChanges, beforeChanges)
	}
}

func seedTask21CorruptJournal(t *testing.T, path string) {
	t.Helper()
	first := task21MaintenanceEvent("EVT-20260822-000001", 1)
	second := task21MaintenanceEvent("EVT-20260822-000001", 3)
	payload := make([]byte, 0)
	for _, event := range []journal.Event{first, second} {
		line, err := json.Marshal(event)
		if err != nil {
			t.Fatalf("marshal corrupt Event: %v", err)
		}
		payload = append(payload, line...)
		payload = append(payload, '\n')
	}
	if err := os.WriteFile(path, payload, 0o644); err != nil {
		t.Fatalf("write corrupt journal: %v", err)
	}
}

func task21MaintenanceEvent(eventID string, sequence int) journal.Event {
	return journal.Event{
		EventID:     eventID,
		EventSchema: journal.EventSchemaV2,
		EventType:   "label.added",
		OccurredAt:  time.Date(2026, time.August, 22, 0, 0, 0, 0, time.UTC),
		ActorV2: journal.EventActor{
			ActorID:       "actor-maintenance",
			ActorType:     "service",
			Producer:      "fid",
			ProvenanceRef: "test:maintenance",
		},
		Aggregate:   journal.Aggregate{ID: "LABEL-maintenance", Type: "label", Sequence: sequence},
		CommandName: "label.add",
		Effects: []journal.Effect{{
			Op: "insert", Table: "labels", PK: map[string]any{"name": "maintenance"},
			Values: map[string]any{"colour": "#010203", "description": "maintenance"},
		}},
	}
}

func seedTask21CorruptProjection(t *testing.T, db *sql.DB) {
	t.Helper()
	statements := []string{
		`INSERT INTO idempotency_results (operation_scope, idempotency_key, request_digest, result_entity_id, result_event_id, aggregate_sequence, result_json, created_at, request_field_digests_json) VALUES ('task21', 'binding', 'sha256:test', 'LABEL-maintenance', 'EVT-missing', 9, '{}', '2026-08-22T00:00:00Z', '{}')`,
		`INSERT INTO kit_activations (id, kit_identity, backend_identity, phase, status, point_of_no_return, direction, created_at, updated_at) VALUES ('ACT-1', 'kit-1', 'fid-1', 21, 'finalized', 1, 'activate', '2026-08-22T00:00:00Z', '2026-08-22T00:00:00Z')`,
		`INSERT INTO activation_gate_definitions (activation_id, definition_id, definition_version, definition_digest, definition_json) VALUES ('ACT-1', 'DEF-1', 1, 'sha256:def', '{}')`,
		`INSERT INTO gate_registrations (id, version, activation_id, definition_id, definition_version, definition_digest, binding_digest, project_id, worktree_id, scope_json, policy_json, authorization_provenance, created_event_id, created_at) VALUES ('REG-1', 1, 'ACT-1', 'DEF-1', 1, 'sha256:def', 'sha256:binding', 'project-1', 'worktree-1', '{}', '{}', 'test', 'EVT-20260822-000001', '2026-08-22T00:00:00Z')`,
		`INSERT INTO work_envelopes (id, current_revision, status, created_at) VALUES ('ENV-1', 1, 'active', '2026-08-22T00:00:00Z')`,
		`INSERT INTO envelope_revisions (envelope_id, revision, prior_revision, outcome, non_goals_json, risk_result_json, repository_scope_json, status, created_actor_id, created_at) VALUES ('ENV-1', 1, NULL, 'test', '[]', '{}', '{}', 'active', 'actor-maintenance', '2026-08-22T00:00:00Z')`,
		`INSERT INTO envelope_obligations (id, envelope_id, revision, obligation_type, registration_id, registration_version, obligation_json) VALUES ('OBL-1', 'ENV-1', 1, 'gate-evidence', 'REG-1', 1, '{}')`,
		`INSERT INTO claims (id, issue_id, envelope_id, envelope_revision, actor_id, worktree_id, scope_json, mode, version, status, acquired_at, expires_at) VALUES ('CLM-1', 'FID-final-schema', 'ENV-1', 1, 'actor-maintenance', 'worktree-1', '{}', 'mutation', 1, 'active', '2026-08-22T00:00:00Z', '2026-08-22T01:00:00Z')`,
		`INSERT INTO runs (id, issue_id, envelope_id, envelope_revision, claim_id, claim_version, actor_id, start_repository_json, status, started_at) VALUES ('RUN-1', 'FID-final-schema', 'ENV-1', 1, 'CLM-1', 1, 'actor-maintenance', '{}', 'active', '2026-08-22T00:00:00Z')`,
		`INSERT INTO execution_attempts (id, registration_id, registration_version, obligation_id, run_id, actor_id, start_repository_json, resolved_policy_digest, producer, started_at) VALUES ('ATT-1', 'REG-1', 1, 'OBL-1', 'RUN-1', 'actor-maintenance', '{}', 'sha256:policy', 'fid', '2026-08-22T00:00:00Z')`,
		`INSERT INTO referenced_artifacts (id, digest, media_type, size_bytes, retention_class, availability, location_policy, created_at) VALUES ('ART-1', 'sha256:missing', 'text/plain', 0, 'evidence', 'missing', 'content-addressed', '2026-08-22T00:00:00Z')`,
		`CREATE TABLE roles (name TEXT PRIMARY KEY)`,
	}
	for _, statement := range statements {
		if _, err := db.Exec(statement); err != nil {
			t.Fatalf("seed task-21 corrupt projection: %v", err)
		}
	}
}

func hasFinding(findings []Finding, check, entityID string) bool {
	for _, finding := range findings {
		if finding.Check == check && finding.EntityID == entityID {
			return true
		}
	}
	return false
}
