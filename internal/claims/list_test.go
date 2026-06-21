package claims

import (
	"path/filepath"
	"testing"
	"time"
)

func TestListFiltersClaimsByRoleAndStalenessAndOrdersByRowID(t *testing.T) {
	t.Parallel()

	workspace := newClaimsWorkspace(t)
	first := createClaimIssue(t, workspace, "First live claim")
	second := createClaimIssue(t, workspace, "Stale claim")
	third := createClaimIssue(t, workspace, "Human live claim")

	conn := openClaimsDB(t, filepath.Join(workspace, ".fid", "fid.sqlite"))
	defer closeClaimsDB(t, conn)

	now := time.Now().UTC().Truncate(time.Second)
	statements := []struct {
		refID      string
		claimedBy  string
		claimedAt  string
		ttlSeconds int
		note       any
	}{
		{
			refID:      first.ID,
			claimedBy:  "developer",
			claimedAt:  now.Add(-5 * time.Minute).Format(time.RFC3339),
			ttlSeconds: 1800,
			note:       "live developer claim",
		},
		{
			refID:      second.ID,
			claimedBy:  "qa_reviewer",
			claimedAt:  now.Add(-2 * time.Hour).Format(time.RFC3339),
			ttlSeconds: 1800,
			note:       nil,
		},
		{
			refID:      third.ID,
			claimedBy:  "human",
			claimedAt:  now.Add(-10 * time.Minute).Format(time.RFC3339),
			ttlSeconds: 1800,
			note:       "human held claim",
		},
	}

	for _, statement := range statements {
		if _, err := conn.Exec(`
			INSERT INTO claims (ref_type, ref_id, claimed_by, claimed_at, ttl_seconds, note)
			VALUES (?, ?, ?, ?, ?, ?)
		`, "issue", statement.refID, statement.claimedBy, statement.claimedAt, statement.ttlSeconds, statement.note); err != nil {
			t.Fatalf("insert claim for %s: %v", statement.refID, err)
		}
	}

	liveItems, err := List(ListOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testClaimsMigrationsDir(t),
	})
	if err != nil {
		t.Fatalf("List(live) error = %v", err)
	}
	if len(liveItems) != 2 {
		t.Fatalf("List(live) len = %d, want 2", len(liveItems))
	}
	if liveItems[0].RefID != first.ID || liveItems[1].RefID != third.ID {
		t.Fatalf("List(live) ids = %#v, want [%q, %q]", []string{liveItems[0].RefID, liveItems[1].RefID}, first.ID, third.ID)
	}
	if liveItems[0].ClaimedBy != "developer" {
		t.Fatalf("first live claimed_by = %q, want developer", liveItems[0].ClaimedBy)
	}
	if liveItems[1].ClaimedBy != "human" {
		t.Fatalf("second live claimed_by = %q, want human", liveItems[1].ClaimedBy)
	}
	if liveItems[0].IsStale || liveItems[1].IsStale {
		t.Fatalf("live items stale flags = %#v, want all false", []bool{liveItems[0].IsStale, liveItems[1].IsStale})
	}
	if liveItems[0].ExpiresAt == "" || liveItems[1].ExpiresAt == "" {
		t.Fatalf("live expires_at values = [%q, %q], want populated", liveItems[0].ExpiresAt, liveItems[1].ExpiresAt)
	}

	humanItems, err := List(ListOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testClaimsMigrationsDir(t),
		ClaimedBy:     "human",
	})
	if err != nil {
		t.Fatalf("List(human) error = %v", err)
	}
	if len(humanItems) != 1 {
		t.Fatalf("List(human) len = %d, want 1", len(humanItems))
	}
	if humanItems[0].RefID != third.ID || humanItems[0].ClaimedBy != "human" {
		t.Fatalf("List(human) item = %#v, want human claim on %q", humanItems[0], third.ID)
	}

	staleItems, err := List(ListOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testClaimsMigrationsDir(t),
		StaleOnly:     true,
	})
	if err != nil {
		t.Fatalf("List(stale) error = %v", err)
	}
	if len(staleItems) != 1 {
		t.Fatalf("List(stale) len = %d, want 1", len(staleItems))
	}
	if staleItems[0].RefID != second.ID {
		t.Fatalf("List(stale) ref_id = %q, want %q", staleItems[0].RefID, second.ID)
	}
	if !staleItems[0].IsStale {
		t.Fatal("List(stale) is_stale = false, want true")
	}
}
