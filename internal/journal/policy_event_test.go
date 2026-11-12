package journal

import (
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/mbell-nzl/fid/internal/config"
	"github.com/mbell-nzl/fid/internal/policy"
)

func snapshotEvent(t *testing.T, sequence int) Event {
	t.Helper()
	file, err := config.DecodeProject(nil)
	if err != nil {
		t.Fatal(err)
	}
	effective := policy.Resolve(file)
	encoded, err := effective.Canonical()
	if err != nil {
		t.Fatal(err)
	}
	return Event{EventSchema: EventSchemaV2, EventID: "EVT-20261112-000001", EventType: "policy.snapshotted",
		OccurredAt: time.Date(2026, 11, 12, 0, 0, 0, 0, time.UTC),
		ActorV2:    EventActor{ActorID: "fid", ActorType: "service", Producer: "fid-standalone", ProvenanceRef: "built-in-policy"},
		Aggregate:  Aggregate{ID: "policy", Sequence: sequence, Type: "policy"}, CommandName: "policy.snapshot",
		Effects: []Effect{{Op: "insert", Table: "policy_snapshots", PK: map[string]any{"sequence": sequence},
			Values: map[string]any{"source_ref": "built-in defaults", "source_commit": nil, "digest": effective.Digest(), "effective_json": string(encoded)}}}}
}

func TestPolicySnapshotEvent(t *testing.T) {
	t.Parallel()
	for _, test := range []struct {
		name   string
		mutate func(*Event)
		valid  bool
	}{
		{"valid", func(*Event) {}, true},
		{"wrong digest", func(event *Event) { event.Effects[0].Values["digest"] = "sha256:" + strings.Repeat("0", 64) }, false},
		{"incomplete payload", func(event *Event) { event.Effects[0].Values["effective_json"] = "{}" }, false},
		{"wrong source", func(event *Event) { event.Effects[0].Values["source_ref"] = "main" }, false},
		{"missing origin commit", func(event *Event) { event.Effects[0].Values["source_ref"] = "origin/main" }, false},
		{"origin", func(event *Event) {
			event.Effects[0].Values["source_ref"] = "origin/trunk"
			event.Effects[0].Values["source_commit"] = strings.Repeat("a", 40)
		}, true},
		{"wrong sequence", func(event *Event) { event.Effects[0].PK["sequence"] = 2 }, false},
		{"wrong aggregate", func(event *Event) { event.Aggregate.ID = "other" }, false},
		{"wrong writer", func(event *Event) { event.ActorV2.ActorType = "agent" }, false},
		{"unrelated effect", func(event *Event) {
			event.Effects = append(event.Effects, Effect{Op: "delete", Table: "labels", PK: map[string]any{"name": "x"}})
		}, false},
		{"duplicate effect", func(event *Event) { event.Effects = append(event.Effects, event.Effects[0]) }, false},
		{"another family", func(event *Event) { event.EventType = "hold.taken"; event.Aggregate.Type = "hold" }, false},
	} {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			conn := openReplayDB(t, filepath.Join(t.TempDir(), "projection.sqlite"))
			defer closeProjectDB(t, conn)
			event := snapshotEvent(t, 1)
			test.mutate(&event)
			err := DryRunProjectEvent(conn, event)
			if (err == nil) != test.valid {
				t.Fatalf("snapshot validation = %v, valid = %v", err, test.valid)
			}
		})
	}
}
