package maintenance

import (
	"reflect"
	"strings"
	"testing"
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
		"dependent_row_kind",
		"dependent_row_kind",
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
		"link_refs",
		"link_relation_compatibility",
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

func TestCheckErrorsExcludesEveryResolvedV01Rule(t *testing.T) {
	t.Parallel()

	fixture := newValidStructuralFixture(t)
	got, err := Check(CheckOptions{DB: fixture.db, JournalDir: fixture.journalDir})
	if err != nil {
		t.Fatalf("Check() error = %v", err)
	}
	if len(got.Errors) != 0 {
		t.Fatalf("Check().Errors = %#v, want no errors", got.Errors)
	}
}

func TestCheckWorkflowStatesReadsProjectedWorkflowTables(t *testing.T) {

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

func hasFinding(findings []Finding, check, entityID string) bool {
	for _, finding := range findings {
		if finding.Check == check && finding.EntityID == entityID {
			return true
		}
	}
	return false
}
