package claims

import (
	"errors"
	"path/filepath"
	"testing"
	"time"
)

func TestOperationTimeUsesSuppliedUTCAndDefaultsToCurrentUTC(t *testing.T) {
	t.Parallel()

	zone := time.FixedZone("test", 12*60*60)
	supplied := time.Date(2040, time.February, 3, 4, 5, 6, 7, zone)
	if got := operationTime(supplied); !got.Equal(supplied.UTC()) || got.Location() != time.UTC {
		t.Fatalf("operationTime(supplied) = %s in %s, want %s in UTC", got, got.Location(), supplied.UTC())
	}

	before := time.Now().UTC()
	got := operationTime(time.Time{})
	after := time.Now().UTC()
	if got.Before(before) || got.After(after) {
		t.Fatalf("operationTime(zero) = %s, want between %s and %s", got, before, after)
	}
	if got.Location() != time.UTC {
		t.Fatalf("operationTime(zero) location = %s, want UTC", got.Location())
	}
}

func TestRenewUsesSuppliedTimeForValidationAndTimestamp(t *testing.T) {
	t.Parallel()

	claimedAt := time.Date(2041, time.March, 4, 5, 6, 7, 0, time.UTC)
	tests := []struct {
		name        string
		now         time.Time
		wantNotHeld bool
	}{
		{name: "before expiry", now: claimedAt.Add(30 * time.Second)},
		{name: "at expiry", now: claimedAt.Add(60 * time.Second), wantNotHeld: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			workspace, issueID := acquireFixedClaim(t, claimedAt, 60)
			result, err := Renew(RenewOptions{
				WorkspaceRoot: workspace,
				MigrationsDir: testClaimsMigrationsDir(t),
				IssueID:       issueID,
				ActorRole:     "developer",
				Now:           test.now,
			})
			if test.wantNotHeld {
				var notHeldErr *ClaimNotHeldError
				if !errors.As(err, &notHeldErr) {
					t.Fatalf("Renew() error = %v, want *ClaimNotHeldError", err)
				}
				return
			}
			if err != nil {
				t.Fatalf("Renew() error = %v", err)
			}
			if result.ClaimedAt != test.now.Format(time.RFC3339) {
				t.Fatalf("Renew().ClaimedAt = %q, want %q", result.ClaimedAt, test.now.Format(time.RFC3339))
			}
		})
	}
}

func TestReleaseUsesSuppliedTimeForValidationAndTimestamp(t *testing.T) {
	t.Parallel()

	claimedAt := time.Date(2042, time.April, 5, 6, 7, 8, 0, time.UTC)
	tests := []struct {
		name        string
		now         time.Time
		wantNotHeld bool
	}{
		{name: "before expiry", now: claimedAt.Add(60*time.Second - time.Nanosecond)},
		{name: "at expiry", now: claimedAt.Add(60 * time.Second), wantNotHeld: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			workspace, issueID := acquireFixedClaim(t, claimedAt, 60)
			_, err := Release(ReleaseOptions{
				WorkspaceRoot: workspace,
				MigrationsDir: testClaimsMigrationsDir(t),
				IssueID:       issueID,
				ActorRole:     "developer",
				Now:           test.now,
			})
			if test.wantNotHeld {
				var notHeldErr *ClaimNotHeldError
				if !errors.As(err, &notHeldErr) {
					t.Fatalf("Release() error = %v, want *ClaimNotHeldError", err)
				}
				return
			}
			if err != nil {
				t.Fatalf("Release() error = %v", err)
			}

			journalPath := filepath.Join(workspace, ".fid", "journal", test.now.Format("2006-01-02")+".jsonl")
			events := readClaimStoredEvents(t, journalPath)
			last := events[len(events)-1]
			if !last.OccurredAt.Equal(test.now) {
				t.Fatalf("release occurred_at = %s, want %s", last.OccurredAt, test.now)
			}
		})
	}
}

func TestListAndLiveValidationUseSuppliedTimeAtExpiryBoundary(t *testing.T) {
	t.Parallel()

	claimedAt := time.Date(2043, time.May, 6, 7, 8, 9, 0, time.UTC)
	workspace, issueID := acquireFixedClaim(t, claimedAt, 60)
	beforeExpiry := claimedAt.Add(60*time.Second - time.Nanosecond)
	atExpiry := claimedAt.Add(60 * time.Second)

	live, err := List(ListOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testClaimsMigrationsDir(t),
		Now:           beforeExpiry,
	})
	if err != nil {
		t.Fatalf("List(before expiry) error = %v", err)
	}
	if len(live) != 1 || live[0].RefID != issueID || live[0].IsStale {
		t.Fatalf("List(before expiry) = %#v, want one live claim on %q", live, issueID)
	}

	stale, err := List(ListOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testClaimsMigrationsDir(t),
		StaleOnly:     true,
		Now:           atExpiry,
	})
	if err != nil {
		t.Fatalf("List(at expiry) error = %v", err)
	}
	if len(stale) != 1 || stale[0].RefID != issueID || !stale[0].IsStale {
		t.Fatalf("List(at expiry) = %#v, want one stale claim on %q", stale, issueID)
	}

	conn := openClaimsDB(t, filepath.Join(workspace, ".fid", "fid.sqlite"))
	defer closeClaimsDB(t, conn)
	if _, err := requireLiveIssueClaimHeldAt(conn, issueID, "developer", beforeExpiry); err != nil {
		t.Fatalf("requireLiveIssueClaimHeldAt(before expiry) error = %v", err)
	}
	if _, err := requireLiveIssueClaimHeldAt(conn, issueID, "developer", atExpiry); err == nil {
		t.Fatal("requireLiveIssueClaimHeldAt(at expiry) error = nil, want rejection")
	}
}

func acquireFixedClaim(t *testing.T, claimedAt time.Time, ttlSeconds int) (string, string) {
	t.Helper()

	workspace := newClaimsWorkspace(t)
	created := createClaimIssue(t, workspace, "Deterministic claim time")
	if _, err := Acquire(AcquireOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testClaimsMigrationsDir(t),
		IssueID:       created.ID,
		ActorRole:     "developer",
		TTLSeconds:    &ttlSeconds,
		Now:           claimedAt,
	}); err != nil {
		t.Fatalf("Acquire(fixed) error = %v", err)
	}
	return workspace, created.ID
}
