package journal

import "fmt"

type UpgradeFunc func(Event) (Event, error)

type SchemaUpgrader struct {
	FromSchema int
	ToSchema   int
	Upgrade    UpgradeFunc
}

type UpgradeChain struct {
	currentSchema       int
	upgraders           map[int]SchemaUpgrader
	allowLegacyIdentity bool
}

var defaultUpgradeChain = mustNewDefaultUpgradeChain()

func NewUpgradeChain(currentSchema int, upgraders ...SchemaUpgrader) (UpgradeChain, error) {
	if currentSchema <= 0 {
		return UpgradeChain{}, fmt.Errorf("new upgrade chain: current schema must be positive")
	}

	registry := make(map[int]SchemaUpgrader, len(upgraders))
	for _, upgrader := range upgraders {
		if upgrader.FromSchema <= 0 {
			return UpgradeChain{}, fmt.Errorf("new upgrade chain: from schema must be positive")
		}
		if upgrader.ToSchema <= 0 {
			return UpgradeChain{}, fmt.Errorf("new upgrade chain: to schema must be positive")
		}
		if upgrader.Upgrade == nil {
			return UpgradeChain{}, fmt.Errorf("new upgrade chain: upgrader for schema %d is nil", upgrader.FromSchema)
		}
		if _, exists := registry[upgrader.FromSchema]; exists {
			return UpgradeChain{}, fmt.Errorf("new upgrade chain: duplicate upgrader for schema %d", upgrader.FromSchema)
		}

		registry[upgrader.FromSchema] = upgrader
	}

	return UpgradeChain{
		currentSchema: currentSchema,
		upgraders:     registry,
	}, nil
}

func UpgradeToCurrent(event Event) (Event, error) {
	return defaultUpgradeChain.Upgrade(event)
}

func (chain UpgradeChain) Upgrade(event Event) (Event, error) {
	if event.EventSchema <= 0 {
		return Event{}, fmt.Errorf("upgrade event %q: event schema must be positive", event.EventID)
	}
	if event.EventSchema > chain.currentSchema {
		return Event{}, fmt.Errorf("upgrade event %q: schema %d is newer than current %d", event.EventID, event.EventSchema, chain.currentSchema)
	}

	current := event
	for {
		upgrader, exists := chain.upgraders[current.EventSchema]
		if !exists {
			return Event{}, fmt.Errorf("upgrade event %q: no upgrader registered for schema %d", current.EventID, current.EventSchema)
		}

		next, err := upgrader.Upgrade(current)
		if err != nil {
			return Event{}, fmt.Errorf("upgrade event %q from schema %d: %w", current.EventID, current.EventSchema, err)
		}

		if upgrader.ToSchema == current.EventSchema {
			if current.EventSchema != chain.currentSchema {
				if chain.allowLegacyIdentity {
					return next, nil
				}
				return Event{}, fmt.Errorf("upgrade event %q: identity upgrader for non-current schema %d", current.EventID, current.EventSchema)
			}
			if next.EventSchema != current.EventSchema {
				return Event{}, fmt.Errorf("upgrade event %q: identity upgrader returned schema %d", current.EventID, next.EventSchema)
			}
			return next, nil
		}

		if next.EventSchema != upgrader.ToSchema {
			return Event{}, fmt.Errorf("upgrade event %q: upgrader for schema %d returned schema %d, want %d", current.EventID, current.EventSchema, next.EventSchema, upgrader.ToSchema)
		}

		current = next
		if current.EventSchema == chain.currentSchema {
			return current, nil
		}
	}
}

func mustNewDefaultUpgradeChain() UpgradeChain {
	chain, err := NewUpgradeChain(EventSchemaV2,
		SchemaUpgrader{
			FromSchema: EventSchemaV1,
			ToSchema:   EventSchemaV1,
			Upgrade: func(event Event) (Event, error) {
				return event, nil
			},
		}, SchemaUpgrader{
			FromSchema: EventSchemaV2,
			ToSchema:   EventSchemaV2,
			Upgrade: func(event Event) (Event, error) {
				return event, nil
			},
		})
	if err != nil {
		panic(err)
	}

	chain.allowLegacyIdentity = true
	return chain
}
