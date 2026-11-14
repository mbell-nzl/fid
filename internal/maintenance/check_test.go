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

func TestCheckReportsFinalSchemaIntegrityClassesWithoutRepair(t *testing.T) {
	t.Parallel()
	fixture := newMaintenanceFixture(t)
	seedIssue(t, fixture.db, "FID-final-schema", "task", nil, nil, "task-default", "todo")
	seedTask21CorruptJournal(t, fixture.journalPath)
	seedRetainedCorruptProjection(t, fixture.db)
	beforeJournal := fixture.journalBytes(t)
	beforeChanges := fixture.changeCount(t)

	result, err := Check(CheckOptions{DB: fixture.db, JournalDir: fixture.journalDir})
	if err != nil {
		t.Fatalf("Check() error = %v", err)
	}

	want := []Finding{
		{Check: "aggregate-sequence", EntityID: "label:LABEL-maintenance", Message: "aggregate sequence 3 follows 1, want 2"},
		{Check: "artifact-integrity", EntityID: "ART-1", Message: "Artifact is unavailable"},
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

func hasFinding(findings []Finding, check, entityID string) bool {
	for _, finding := range findings {
		if finding.Check == check && finding.EntityID == entityID {
			return true
		}
	}
	return false
}

func seedRetainedCorruptProjection(t *testing.T, db *sql.DB) {
	t.Helper()
	for _, statement := range []string{
		`INSERT INTO idempotency_results (operation_scope, idempotency_key, request_digest, result_entity_id, result_event_id, aggregate_sequence, result_json, created_at, request_field_digests_json) VALUES ('task21', 'binding', 'sha256:test', 'LABEL-maintenance', 'EVT-missing', 9, '{}', '2026-08-22T00:00:00Z', '{}')`,
		`INSERT INTO referenced_artifacts (id, digest, media_type, size_bytes, retention_class, availability, location_policy, created_at) VALUES ('ART-1', 'sha256:missing', 'text/plain', 0, 'evidence', 'missing', 'content-addressed', '2026-08-22T00:00:00Z')`,
		`CREATE TABLE roles (name TEXT PRIMARY KEY)`,
	} {
		if _, err := db.Exec(statement); err != nil {
			t.Fatal(err)
		}
	}
}
