package journal

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
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
		{"unknown lane", changeSnapshotPolicy(t, func(effective *policy.Effective) {
			effective.Lanes["typo"] = effective.Lanes["assured"]
			delete(effective.Lanes, "assured")
		}), false},
		{"empty argv", changeSnapshotPolicy(t, func(effective *policy.Effective) {
			effective.Checks["test"] = config.Check{Run: []string{}, Timeout: time.Minute, DefinedBy: []string{}, WhenChanged: []string{}}
		}), false},
		{"empty executable", changeSnapshotPolicy(t, func(effective *policy.Effective) {
			effective.Checks["test"] = config.Check{Run: []string{""}, Timeout: time.Minute, DefinedBy: []string{}, WhenChanged: []string{}}
		}), false},
		{"negative timeout", changeSnapshotPolicy(t, func(effective *policy.Effective) {
			effective.Checks["test"] = config.Check{Run: []string{"true"}, Timeout: -time.Minute, DefinedBy: []string{}, WhenChanged: []string{}}
		}), false},
		{"zero timeout", changeSnapshotPolicy(t, func(effective *policy.Effective) {
			effective.Checks["test"] = config.Check{Run: []string{"true"}, DefinedBy: []string{}, WhenChanged: []string{}}
		}), false},
		{"bad defining glob", changeSnapshotPolicy(t, func(effective *policy.Effective) {
			effective.Checks["test"] = config.Check{Run: []string{"true"}, Timeout: time.Minute, DefinedBy: []string{"../Makefile"}, WhenChanged: []string{}}
		}), false},
		{"bad conditional glob", changeSnapshotPolicy(t, func(effective *policy.Effective) {
			effective.Checks["test"] = config.Check{Run: []string{"true"}, Timeout: time.Minute, DefinedBy: []string{}, WhenChanged: []string{"/src/**"}}
		}), false},
		{"invalid lane values", changeSnapshotPolicy(t, func(effective *policy.Effective) {
			effective.Lanes["routine"] = config.Lane{Paths: []string{}, Checks: []string{}, ReworkCap: -1, TokenBudget: -1}
		}), false},
		{"undefined check", changeSnapshotPolicy(t, func(effective *policy.Effective) {
			effective.Lanes["routine"] = config.Lane{Paths: []string{}, Checks: []string{"missing"}, ChecksSet: true}
		}), false},
		{"unset nonempty checks", changeSnapshotPolicy(t, func(effective *policy.Effective) {
			effective.Checks["test"] = config.Check{Run: []string{"true"}, Timeout: time.Minute, DefinedBy: []string{}, WhenChanged: []string{}}
			effective.Lanes["routine"] = config.Lane{Paths: []string{}, Checks: []string{"test"}, ChecksSet: false}
		}), false},
		{"standard paths", changeSnapshotPolicy(t, func(effective *policy.Effective) {
			effective.Lanes["standard"] = config.Lane{Paths: []string{"src/**"}, Checks: []string{}}
		}), false},
		{"assured auto merge", changeSnapshotPolicy(t, func(effective *policy.Effective) {
			effective.Lanes["assured"] = config.Lane{Paths: []string{}, Checks: []string{}, AutoMerge: true}
		}), false},
		{"invalid lane glob", changeSnapshotPolicy(t, func(effective *policy.Effective) {
			effective.Lanes["routine"] = config.Lane{Paths: []string{"../src"}, Checks: []string{}}
		}), false},
		{"missing test paths marker", func(event *Event) {
			values := event.Effects[0].Values
			raw := strings.Replace(values["effective_json"].(string), `,"test_adequacy_set":false`, "", 1)
			values["effective_json"] = raw
			values["digest"] = fmt.Sprintf("sha256:%x", sha256.Sum256([]byte(raw)))
		}, false},
		{"configured test paths", changeSnapshotPolicy(t, func(effective *policy.Effective) {
			effective.TestAdequacy = []string{"**/*_test.go"}
			effective.TestAdequacySet = true
		}), true},
		{"explicit empty test paths", changeSnapshotPolicy(t, func(effective *policy.Effective) {
			effective.TestAdequacySet = true
		}), true},
		{"unset nonempty test paths", changeSnapshotPolicy(t, func(effective *policy.Effective) {
			effective.TestAdequacy = []string{"**/*_test.go"}
		}), false},
		{"invalid review glob", changeSnapshotPolicy(t, func(effective *policy.Effective) {
			effective.TestAdequacy = []string{"[test]"}
			effective.TestAdequacySet = true
		}), false},
		{"invalid hard floor glob", changeSnapshotPolicy(t, func(effective *policy.Effective) {
			effective.HardAssured = []string{"../policy"}
		}), false},
		{"missing actor type", changeSnapshotPolicy(t, func(effective *policy.Effective) {
			delete(effective.ActorRules, "human")
		}), false},
		{"unknown actor type", changeSnapshotPolicy(t, func(effective *policy.Effective) {
			effective.ActorRules["robot"] = effective.ActorRules["human"]
			delete(effective.ActorRules, "human")
		}), false},
		{"undefined action set", changeSnapshotPolicy(t, func(effective *policy.Effective) {
			effective.ActorRules["human"] = []policy.ActionSet{"missing"}
		}), false},
		{"missing action set", changeSnapshotPolicy(t, func(effective *policy.Effective) {
			delete(effective.ActionSets, policy.PersonDecisions)
		}), false},
		{"empty action", changeSnapshotPolicy(t, func(effective *policy.Effective) {
			effective.ActionSets[policy.WorkSet] = []string{""}
		}), false},
		{"historical built-in values", changeSnapshotPolicy(t, func(effective *policy.Effective) {
			effective.RulesVersion++
			effective.HardAssured = []string{"legacy-policy/**"}
			effective.ActionSets[policy.WorkSet] = []string{"historical-action"}
			effective.ActionSets[policy.PersonDecisions] = []string{"historical-decision"}
			effective.ActorRules["human"] = []policy.ActionSet{policy.PersonDecisions}
		}), true},
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

// Recompute both fields so malformed-value cases exercise schema validation,
// rather than merely detecting an obsolete digest.
func changeSnapshotPolicy(t *testing.T, mutate func(*policy.Effective)) func(*Event) {
	t.Helper()
	return func(event *Event) {
		t.Helper()
		var effective policy.Effective
		if err := json.Unmarshal([]byte(event.Effects[0].Values["effective_json"].(string)), &effective); err != nil {
			t.Fatal(err)
		}
		mutate(&effective)
		encoded, err := effective.Canonical()
		if err != nil {
			t.Fatal(err)
		}
		event.Effects[0].Values["effective_json"] = string(encoded)
		event.Effects[0].Values["digest"] = effective.Digest()
	}
}
