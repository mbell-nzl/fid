package journal

import (
	"errors"
	"reflect"
	"testing"
)

func TestUpgradeToCurrentLeavesV1EventUnchanged(t *testing.T) {
	t.Parallel()

	event := testV1Event()

	got, err := UpgradeToCurrent(event)
	if err != nil {
		t.Fatalf("UpgradeToCurrent() error = %v", err)
	}

	if !reflect.DeepEqual(got, event) {
		t.Fatalf("UpgradeToCurrent() changed event = %#v, want %#v", got, event)
	}
}

func TestUpgradeChainDispatchesByEventSchema(t *testing.T) {
	t.Parallel()

	called := false
	chain, err := NewUpgradeChain(2,
		SchemaUpgrader{
			FromSchema: 1,
			ToSchema:   2,
			Upgrade: func(event Event) (Event, error) {
				called = true
				event.EventSchema = 2
				event.EventType = "issue.transitioned.upgraded"
				return event, nil
			},
		},
		SchemaUpgrader{
			FromSchema: 2,
			ToSchema:   2,
			Upgrade: func(event Event) (Event, error) {
				return event, nil
			},
		},
	)
	if err != nil {
		t.Fatalf("NewUpgradeChain() error = %v", err)
	}

	event := testV1Event()
	got, err := chain.Upgrade(event)
	if err != nil {
		t.Fatalf("Upgrade() error = %v", err)
	}

	if !called {
		t.Fatal("Upgrade() did not dispatch the v1 upgrader")
	}
	if got.EventSchema != 2 {
		t.Fatalf("Upgrade() schema = %d, want 2", got.EventSchema)
	}
	if got.EventType != "issue.transitioned.upgraded" {
		t.Fatalf("Upgrade() type = %q, want %q", got.EventType, "issue.transitioned.upgraded")
	}
}

func TestNewUpgradeChainRejectsInvalidDefinitions(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		current int
		input   []SchemaUpgrader
		wantErr string
	}{
		{
			name:    "current schema must be positive",
			current: 0,
			wantErr: "new upgrade chain: current schema must be positive",
		},
		{
			name:    "from schema must be positive",
			current: 1,
			input:   []SchemaUpgrader{{FromSchema: 0, ToSchema: 1, Upgrade: func(event Event) (Event, error) { return event, nil }}},
			wantErr: "new upgrade chain: from schema must be positive",
		},
		{
			name:    "to schema must be positive",
			current: 1,
			input:   []SchemaUpgrader{{FromSchema: 1, ToSchema: 0, Upgrade: func(event Event) (Event, error) { return event, nil }}},
			wantErr: "new upgrade chain: to schema must be positive",
		},
		{
			name:    "upgrade func required",
			current: 1,
			input:   []SchemaUpgrader{{FromSchema: 1, ToSchema: 1}},
			wantErr: "new upgrade chain: upgrader for schema 1 is nil",
		},
		{
			name:    "duplicate from schema rejected",
			current: 2,
			input: []SchemaUpgrader{
				{FromSchema: 1, ToSchema: 1, Upgrade: func(event Event) (Event, error) { return event, nil }},
				{FromSchema: 1, ToSchema: 2, Upgrade: func(event Event) (Event, error) { return event, nil }},
			},
			wantErr: "new upgrade chain: duplicate upgrader for schema 1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewUpgradeChain(tt.current, tt.input...)
			if err == nil {
				t.Fatalf("NewUpgradeChain() error = nil, want %q", tt.wantErr)
			}
			if err.Error() != tt.wantErr {
				t.Fatalf("NewUpgradeChain() error = %q, want %q", err.Error(), tt.wantErr)
			}
		})
	}
}

func TestUpgradeRejectsInvalidSchemaStates(t *testing.T) {
	t.Parallel()

	identity := SchemaUpgrader{
		FromSchema: 1,
		ToSchema:   1,
		Upgrade: func(event Event) (Event, error) {
			return event, nil
		},
	}

	tests := []struct {
		name    string
		chain   UpgradeChain
		event   Event
		wantErr string
	}{
		{
			name:    "event schema must be positive",
			chain:   mustUpgradeChain(t, 1, identity),
			event:   Event{EventID: "EVT-1"},
			wantErr: `upgrade event "EVT-1": event schema must be positive`,
		},
		{
			name:    "event schema newer than current",
			chain:   mustUpgradeChain(t, 1, identity),
			event:   Event{EventID: "EVT-1", EventSchema: 2},
			wantErr: `upgrade event "EVT-1": schema 2 is newer than current 1`,
		},
		{
			name:    "missing upgrader",
			chain:   mustUpgradeChain(t, 2, SchemaUpgrader{FromSchema: 2, ToSchema: 2, Upgrade: func(event Event) (Event, error) { return event, nil }}),
			event:   Event{EventID: "EVT-1", EventSchema: 1},
			wantErr: `upgrade event "EVT-1": no upgrader registered for schema 1`,
		},
		{
			name: "upgrader returns error",
			chain: mustUpgradeChain(t, 1, SchemaUpgrader{
				FromSchema: 1,
				ToSchema:   1,
				Upgrade: func(event Event) (Event, error) {
					return Event{}, errors.New("boom")
				},
			}),
			event:   Event{EventID: "EVT-1", EventSchema: 1},
			wantErr: `upgrade event "EVT-1" from schema 1: boom`,
		},
		{
			name: "identity upgrader on non-current schema rejected",
			chain: mustUpgradeChain(t, 2,
				SchemaUpgrader{
					FromSchema: 1,
					ToSchema:   1,
					Upgrade: func(event Event) (Event, error) {
						return event, nil
					},
				},
				SchemaUpgrader{
					FromSchema: 2,
					ToSchema:   2,
					Upgrade: func(event Event) (Event, error) {
						return event, nil
					},
				},
			),
			event:   Event{EventID: "EVT-1", EventSchema: 1},
			wantErr: `upgrade event "EVT-1": identity upgrader for non-current schema 1`,
		},
		{
			name: "identity upgrader must preserve schema",
			chain: mustUpgradeChain(t, 1, SchemaUpgrader{
				FromSchema: 1,
				ToSchema:   1,
				Upgrade: func(event Event) (Event, error) {
					event.EventSchema = 2
					return event, nil
				},
			}),
			event:   Event{EventID: "EVT-1", EventSchema: 1},
			wantErr: `upgrade event "EVT-1": identity upgrader returned schema 2`,
		},
		{
			name: "upgrader must return declared to schema",
			chain: mustUpgradeChain(t, 2,
				SchemaUpgrader{
					FromSchema: 1,
					ToSchema:   2,
					Upgrade: func(event Event) (Event, error) {
						event.EventSchema = 3
						return event, nil
					},
				},
				SchemaUpgrader{
					FromSchema: 2,
					ToSchema:   2,
					Upgrade: func(event Event) (Event, error) {
						return event, nil
					},
				},
			),
			event:   Event{EventID: "EVT-1", EventSchema: 1},
			wantErr: `upgrade event "EVT-1": upgrader for schema 1 returned schema 3, want 2`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.chain.Upgrade(tt.event)
			if err == nil {
				t.Fatalf("Upgrade() error = nil, want %q", tt.wantErr)
			}
			if err.Error() != tt.wantErr {
				t.Fatalf("Upgrade() error = %q, want %q", err.Error(), tt.wantErr)
			}
		})
	}
}

func mustUpgradeChain(t *testing.T, current int, upgraders ...SchemaUpgrader) UpgradeChain {
	t.Helper()

	chain, err := NewUpgradeChain(current, upgraders...)
	if err != nil {
		t.Fatalf("NewUpgradeChain() error = %v", err)
	}

	return chain
}
