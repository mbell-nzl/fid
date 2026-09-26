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

func TestProjectEventFenceRejectsEveryUnavailablePhase23EventWithoutProjectionChange(t *testing.T) {
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

	for eventType := range unavailableJudgmentEvents {
		t.Run(eventType, func(t *testing.T) {
			before := fencedProjectionState(t, conn)
			event := Event{
				EventID:     "EVT-fence-" + strings.ReplaceAll(eventType, ".", "-"),
				EventSchema: EventSchemaV2,
				EventType:   eventType,
				OccurredAt:  time.Date(2026, time.September, 20, 0, 0, 0, 0, time.UTC),
			}
			err := ProjectEvent(conn, ProjectionInput{Event: event, EventHash: "sha256:fence"})
			want := `project event "` + event.EventID + `": Phase 23 Event "` + eventType + `" is unavailable`
			if err == nil || err.Error() != want {
				t.Fatalf("ProjectEvent() error = %v, want %q", err, want)
			}
			if after := fencedProjectionState(t, conn); after != before {
				t.Fatalf("fenced event changed projection: before %q, after %q", before, after)
			}
		})
	}
}

// TestConcernDeviationFenceConvergencePreservesEveryRemainingEvent retains
// VAL-003's convergence boundary:
// implemented judgment Events take the strict current-Event path while every remaining
// unavailable Phase 23 Event stays covered by the fence loop above.
func TestConcernDeviationFenceConvergencePreservesEveryRemainingEvent(t *testing.T) {
	t.Parallel()
	if unavailableJudgmentEvents["review.recorded"] {
		t.Fatal("review.recorded remains fenced after Task 4 implemented its strict event path")
	}
	if unavailableJudgmentEvents["finding.addressed"] {
		t.Fatal("finding.addressed remains fenced after Task 5 implemented its strict event path")
	}
	for _, eventType := range []string{"concern.recorded", "concern.disposed", "deviation.recorded"} {
		if unavailableJudgmentEvents[eventType] {
			t.Fatalf("%s remains fenced after Task 8 implemented its strict event path", eventType)
		}
	}
	if unavailableJudgmentEvents["decision.referenced"] {
		t.Fatal("decision.referenced remains fenced after Task 9 implemented its strict event path")
	}
	if unavailableJudgmentEvents["human-outcome.recorded"] {
		t.Fatal("human-outcome.recorded remains fenced after Task 10 implemented its strict event path")
	}
	want := []string{
		"requirement.linked", "requirement.referenced",
	}
	if len(unavailableJudgmentEvents) != len(want) {
		t.Fatalf("remaining unavailable Event count = %d, want %d", len(unavailableJudgmentEvents), len(want))
	}
	for _, eventType := range want {
		if !unavailableJudgmentEvents[eventType] {
			t.Errorf("remaining Event %q lost its progressive fence", eventType)
		}
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
