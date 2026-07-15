package maintenance

import (
	"reflect"
	"testing"
	"time"
)

func TestCheckWarningsClassifiesQualifyingAndNonQualifyingStates(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, time.July, 15, 12, 0, 0, 0, time.UTC)
	fixture := newWarningFixture(t, now)

	got, err := Check(CheckOptions{
		DB:                              fixture.db,
		JournalDir:                      fixture.journalDir,
		Now:                             now,
		StartedWithoutClaimGraceSeconds: 300,
	})
	if err != nil {
		t.Fatalf("Check() error = %v", err)
	}

	wantChecks := []string{
		"stale_claims",
		"stale_verifications",
		"started_without_claim",
		"started_without_claim",
		"started_without_claim",
		"started_without_claim",
		"telemetry_gaps",
		"telemetry_gaps",
	}
	if gotChecks := findingChecks(got.Warnings); !reflect.DeepEqual(gotChecks, wantChecks) {
		t.Fatalf("Check().Warnings checks = %#v, want %#v", gotChecks, wantChecks)
	}
	assertFindingsOrdered(t, got.Warnings)

	wantEntityIDs := []string{
		"FID-stale-claim",
		"AV-stale",
		"FID-start-after",
		"FID-start-at",
		"FID-start-fallback",
		"FID-start-non-status-update",
		"TC-issue-only-run",
		"TC-missing-run",
	}
	if gotEntityIDs := findingEntityIDs(got.Warnings); !reflect.DeepEqual(gotEntityIDs, wantEntityIDs) {
		t.Fatalf("Check().Warnings entity IDs = %#v, want %#v", gotEntityIDs, wantEntityIDs)
	}
}

func TestCheckWarningsIncludesStartedIssueAtZeroGrace(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, time.July, 15, 12, 0, 0, 0, time.UTC)
	fixture := newMaintenanceFixture(t)
	seedIssueAt(t, fixture.db, "FID-start-zero", "task", "in_progress", now)

	got, err := Check(CheckOptions{DB: fixture.db, JournalDir: fixture.journalDir, Now: now})
	if err != nil {
		t.Fatalf("Check() error = %v", err)
	}
	if want := []string{"FID-start-zero"}; !reflect.DeepEqual(findingEntityIDs(got.Warnings), want) {
		t.Fatalf("Check().Warnings entity IDs = %#v, want %#v", findingEntityIDs(got.Warnings), want)
	}
}
