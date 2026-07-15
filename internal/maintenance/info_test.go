package maintenance

import (
	"reflect"
	"testing"
	"time"
)

func TestCheckInfoCountsAndResultAreDeterministic(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, time.July, 15, 12, 0, 0, 0, time.UTC)
	fixture := newInfoFixture(t, now)
	options := CheckOptions{DB: fixture.db, JournalDir: fixture.journalDir, Now: now, StartedWithoutClaimGraceSeconds: 300}

	got, err := Check(options)
	if err != nil {
		t.Fatalf("Check() error = %v", err)
	}
	if want := (Info{IssueCount: 3, ClaimCountLive: 1, RecentActivity: 2}); got.Info != want {
		t.Fatalf("Check().Info = %#v, want %#v", got.Info, want)
	}

	again, err := Check(options)
	if err != nil {
		t.Fatalf("second Check() error = %v", err)
	}
	if !reflect.DeepEqual(again, got) {
		t.Fatalf("second Check() = %#v, want %#v", again, got)
	}
}
