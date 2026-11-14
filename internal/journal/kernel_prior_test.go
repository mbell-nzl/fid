package journal

// These assertions cover retiring kernel/historical behavior; retain until its owning cut.
import (
	"path/filepath"
	"reflect"
	"testing"
)

func TestKernelEventRequiresExactCapturedPrior(t *testing.T) {
	t.Parallel()
	for _, kind := range []string{"valid", "missing", "wrong", "extra"} {
		t.Run(kind, func(t *testing.T) {
			event := checkpointedProjectionEvent()
			event.Effects[1].Prior = map[string]any{"version": 1}
			switch kind {
			case "missing":
				event.Effects[1].Prior = nil
			case "wrong":
				event.Effects[1].Prior["version"] = 99
			case "extra":
				event.Effects[1].Prior["executor_id"] = "executor-checkpoint"
			}
			conn := openReplayDB(t, filepath.Join(t.TempDir(), "prior.sqlite"))
			defer closeProjectDB(t, conn)
			seedCheckpointProjectionRun(t, conn)
			before := currentProjectionSnapshot(t, conn)
			hash, err := HashEvent(event)
			if err != nil {
				t.Fatal(err)
			}
			err = ProjectEvent(conn, ProjectionInput{Event: event, EventHash: hash})
			if kind == "valid" {
				if err != nil {
					t.Fatal(err)
				}
				return
			}
			if err == nil {
				t.Fatal("invalid captured prior accepted")
			}
			if !reflect.DeepEqual(before, currentProjectionSnapshot(t, conn)) {
				t.Fatal("invalid prior partially projected")
			}
		})
	}
}
