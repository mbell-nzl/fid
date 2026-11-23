package export

import (
	"testing"

	"github.com/mbell-nzl/fid/internal/journal"
)

// These fixtures are private in-memory projections without repositories for
// commands to open. Their revision row explicitly models the guarded bump;
// scalar changes keep the fixture's original arguments and guards.
func writerUpdateFixture(t *testing.T, event journal.Event, observed int) journal.Event {
	t.Helper()
	revision := journal.Effect{Op: "update", Table: "issues", PK: map[string]any{"id": event.Aggregate.ID}, Values: map[string]any{"revision": observed + 1}, Prior: map[string]any{"revision": observed}, Preconditions: map[string]any{"revision": observed}}
	effects := make([]journal.Effect, 0, len(event.Effects)+1)
	inserted := false
	for _, effect := range event.Effects {
		if effect.Table != "issues" {
			effects = append(effects, effect)
			continue
		}
		delete(effect.Values, "revision")
		delete(effect.Preconditions, "revision")
		delete(effect.Prior, "revision")
		if len(effect.Values) == 0 {
			continue
		}
		effect.Values["updated_at"] = event.OccurredAt.UTC().Format("2006-01-02T15:04:05Z07:00")
		effects = append(effects, effect)
		if !inserted {
			effects = append(effects, revision)
			inserted = true
		}
	}
	event.Effects = effects
	return event
}
