package export

import (
	"fmt"
	"strings"

	"github.com/mbell-nzl/fid/internal/journal"
)

func (v *semanticValidator) validateReferences(event journal.Event) error {
	for _, effect := range event.Effects {
		for column, value := range effect.Values {
			if value == nil || !strings.HasSuffix(column, "event_id") {
				continue
			}
			id, ok := value.(string)
			_, prior := v.eventIDs[id]
			if !ok || (id != event.EventID && !prior) {
				return fmt.Errorf("%s.%s references absent Event %q", effect.Table, column, id)
			}
		}
		if event.EventType == "hold.taken" && effect.Table == "holds" && effect.Op == "insert" && effect.PK["id"] != event.Aggregate.ID {
			return fmt.Errorf("hold identity differs from Event aggregate")
		}

	}
	return nil
}
