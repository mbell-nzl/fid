package maintenance

// These assertions cover retiring kernel/historical behavior; retain until its owning cut.
import (
	"testing"
	"time"
)

func TestKernelCheckInfoLiveClaimCount(t *testing.T) {
	t.Parallel()
	now := time.Date(2026, time.July, 15, 12, 0, 0, 0, time.UTC)
	fixture := newInfoFixture(t, now)
	seedClaimAt(t, fixture.db, "FID-info-live", now.Add(-time.Second), 2)
	seedClaimAt(t, fixture.db, "FID-info-stale", now.Add(-time.Second), 1)
	got, err := Check(CheckOptions{DB: fixture.db, JournalDir: fixture.journalDir, Now: now, StartedWithoutClaimGraceSeconds: 300})
	if err != nil {
		t.Fatal(err)
	}
	if want := (Info{IssueCount: 3, ClaimCountLive: 1, RecentActivity: 2}); got.Info != want {
		t.Fatalf("Info=%#v, want %#v", got.Info, want)
	}
}
