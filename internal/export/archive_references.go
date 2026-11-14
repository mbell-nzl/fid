package export

import (
	"fmt"
	"strings"

	"github.com/mbell-nzl/fid/internal/journal"
)

func (v *semanticValidator) validateReferences(event journal.Event) error {
	if err := validateAggregateEffectIdentity(event); err != nil {
		return err
	}
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

var aggregateIdentityTables = map[string]map[string]struct{}{
	"issue":   {"issues": {}},
	"label":   {"labels": {}},
	"link":    {"links": {}},
	"project": {"projects": {}},
}

func validateAggregateEffectIdentity(event journal.Event) error {
	tables := aggregateIdentityTables[event.Aggregate.Type]
	for _, effect := range event.Effects {
		if _, relevant := tables[effect.Table]; !relevant {
			continue
		}
		id, hasID := effect.PK["id"].(string)
		if hasID && id != event.Aggregate.ID {
			return fmt.Errorf("%s identity differs from Event aggregate", effect.Table)
		}
	}
	return nil
}
