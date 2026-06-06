package journal

import (
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
