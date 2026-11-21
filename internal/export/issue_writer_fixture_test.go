package export

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/mbell-nzl/fid/internal/journal"
)

// These fixtures are private in-memory projections without repositories for
// commands to open. Their revision row comes from the captured command corpus;
// scalar changes keep the fixture's original arguments and guards.
func writerUpdateFixture(t *testing.T, event journal.Event, observed int) journal.Event {
	t.Helper()
	raw, err := os.ReadFile("../journal/testdata/row-writer-templates.json")
	if err != nil {
		t.Fatal(err)
	}
	var corpus struct{ Templates map[string]journal.Event }
	if err := json.Unmarshal(raw, &corpus); err != nil {
		t.Fatal(err)
	}
	revision := corpus.Templates["issue.ac-added"].Effects[1]
	revision.PK = map[string]any{"id": event.Aggregate.ID}
	revision.Values = map[string]any{"revision": observed + 1}
	revision.Prior = map[string]any{"revision": observed}
	revision.Preconditions = map[string]any{"revision": observed}
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
