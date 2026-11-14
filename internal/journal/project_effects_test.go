package journal

import (
	"database/sql"
	"fmt"
	"path/filepath"
	"strings"
	"testing"
	"time"

	fiddb "github.com/mbell-nzl/fid/internal/db"
)

func TestProjectEventFenceRejectsMalformedPhase23EventsWithoutProjectionChange(t *testing.T) {
	t.Parallel()
	dbPath := filepath.Join(t.TempDir(), "fid.sqlite")
	if err := fiddb.MigrateUp(dbPath, testProjectMigrationsDir(t)); err != nil {
		t.Fatalf("MigrateUp() error = %v", err)
	}
	conn, err := fiddb.Open(dbPath)
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	defer closeProjectDB(t, conn)

	for _, eventType := range []string{
		"review.recorded", "finding.addressed", "concern.recorded", "concern.disposed",
		"deviation.recorded", "decision.referenced", "human-outcome.recorded",
		"requirement.referenced", "requirement.linked", "unimplemented.phase23-event",
	} {
		t.Run(eventType, func(t *testing.T) {
			before := fencedProjectionState(t, conn)
			event := Event{
				EventID:     "EVT-fence-" + strings.ReplaceAll(eventType, ".", "-"),
				EventSchema: EventSchemaV2,
				EventType:   eventType,
				OccurredAt:  time.Date(2026, time.September, 20, 0, 0, 0, 0, time.UTC),
				Effects:     []Effect{{Op: "insert", Table: "labels", PK: map[string]any{"name": "fence"}, Values: map[string]any{"colour": "#000000"}}},
			}
			err := ProjectEvent(conn, ProjectionInput{Event: event, EventHash: "sha256:fence"})
			if err == nil {
				t.Fatal("malformed Phase 23 Event applied")
			}
			if after := fencedProjectionState(t, conn); after != before {
				t.Fatalf("fenced event changed projection: before %q, after %q", before, after)
			}
		})
	}
}

func fencedProjectionState(t *testing.T, conn interface{ QueryRow(string, ...any) *sql.Row }) string {
	t.Helper()
	var labelCount, metaCount int
	if err := conn.QueryRow(`SELECT COUNT(*) FROM labels`).Scan(&labelCount); err != nil {
		t.Fatalf("count labels: %v", err)
	}
	if err := conn.QueryRow(`SELECT COUNT(*) FROM projection_meta`).Scan(&metaCount); err != nil {
		t.Fatalf("count projection meta: %v", err)
	}
	return fmt.Sprintf("labels=%d;projection_meta=%d", labelCount, metaCount)
}
