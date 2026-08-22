package maintenance

import (
	"errors"
	"fmt"
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

func TestCheckStartedWithoutClaimGraceBoundary(t *testing.T) {
	t.Parallel()

	startedAt := time.Date(2040, time.January, 2, 3, 4, 5, 0, time.UTC)
	tests := []struct {
		name        string
		now         time.Time
		wantFinding bool
	}{
		{name: "just before", now: startedAt.Add(300*time.Second - time.Nanosecond)},
		{name: "exact boundary", now: startedAt.Add(300 * time.Second), wantFinding: true},
		{name: "just after", now: startedAt.Add(300*time.Second + time.Nanosecond), wantFinding: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			fixture := newMaintenanceFixture(t)
			seedJournalEvent(t, fixture.journalPath, "EVT-20400102-000001")
			seedIssueAt(t, fixture.db, "FID-boundary", "task", "in_progress", startedAt)
			result, err := Check(CheckOptions{
				DB:                              fixture.db,
				JournalDir:                      fixture.journalDir,
				Now:                             test.now,
				StartedWithoutClaimGraceSeconds: 300,
			})
			if err != nil {
				t.Fatalf("Check() error = %v", err)
			}
			if got := hasFinding(result.Warnings, "started_without_claim", "FID-boundary"); got != test.wantFinding {
				t.Fatalf("started_without_claim finding = %t, want %t", got, test.wantFinding)
			}
		})
	}
}

func TestCheckRejectsMoreThanOneThousandWarningsWithoutTruncating(t *testing.T) {
	t.Parallel()

	now := time.Date(2040, time.January, 2, 3, 10, 5, 0, time.UTC)
	fixture := newMaintenanceFixture(t)
	for index := 0; index < 1001; index++ {
		seedIssueAt(t, fixture.db, fmt.Sprintf("FID-warning-%04d", index), "task", "in_progress", now.Add(-300*time.Second))
	}

	result, err := Check(CheckOptions{
		DB:                              fixture.db,
		JournalDir:                      fixture.journalDir,
		Now:                             now,
		StartedWithoutClaimGraceSeconds: 300,
	})
	if err == nil {
		t.Fatalf("Check() result = %#v, want an over-limit error for 1001 warnings", result)
	}
	var limit *ResultLimitExceededError
	if !errors.As(err, &limit) || limit.Collection != "maintenance.warnings" || limit.Limit != 1000 || limit.Observed != 1001 {
		t.Fatalf("Check() error = %v, want maintenance.warnings limit 1000 observed 1001", err)
	}
}
