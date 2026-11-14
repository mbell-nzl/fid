package maintenance

import (
	"errors"
	"fmt"
	"testing"
)

func TestCheckRejectsMoreThanOneThousandFindingsWithoutTruncating(t *testing.T) {
	t.Parallel()
	fixture := newMaintenanceFixture(t)
	for index := 0; index < 1001; index++ {
		if _, err := fixture.db.Exec(`INSERT INTO referenced_artifacts(id,digest,media_type,size_bytes,retention_class,availability,location_policy,created_at) VALUES(?, ?, 'text/plain', 0, 'evidence', 'missing', 'content-addressed', '2026-11-15T00:00:00Z')`, fmt.Sprintf("ART-%04d", index), fmt.Sprintf("sha256:missing-%04d", index)); err != nil {
			t.Fatal(err)
		}
	}
	before := fixture.changeCount(t)
	_, err := Check(CheckOptions{DB: fixture.db, JournalDir: fixture.journalDir})
	var limit *ResultLimitExceededError
	if !errors.As(err, &limit) || limit.Collection != "maintenance.errors" || limit.Limit != 1000 || limit.Observed != 1001 {
		t.Fatalf("Check error = %v, want errors limit 1000 observed 1001", err)
	}
	if after := fixture.changeCount(t); after != before {
		t.Fatalf("Check changed rows: %d, want %d", after, before)
	}
}
