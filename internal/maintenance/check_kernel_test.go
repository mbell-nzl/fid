package maintenance

// These assertions cover retiring kernel/historical behavior; retain until its owning cut.
import (
	"reflect"
	"testing"
)

func TestKernelCheckReportsAttemptReceiptIntegrityWithoutRepair(t *testing.T) {
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

func TestKernelCheckAcceptsReferenceActivities(t *testing.T) {
	t.Parallel()
	fixture := newValidStructuralFixture(t)
	seedRequirementRef(t, fixture.db, "REQ-valid")
	seedDecisionRef(t, fixture.db, "DEC-valid", "FID-task-sequential")
	seedActivity(t, fixture.db, "ACT-requirement", "EVT-20260715-000001", "requirement", "REQ-valid")
	seedActivity(t, fixture.db, "ACT-decision", "EVT-20260715-000001", "decision", "DEC-valid")
	got, err := Check(CheckOptions{DB: fixture.db, JournalDir: fixture.journalDir})
	if err != nil || len(got.Errors) != 0 {
		t.Fatalf("Check=%#v/%v", got, err)
	}
}
