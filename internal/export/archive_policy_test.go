package export

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/mbell-nzl/fid/internal/config"
	"github.com/mbell-nzl/fid/internal/journal"
	"github.com/mbell-nzl/fid/internal/policy"
)

func TestSemanticValidatorChecksPolicySnapshots(t *testing.T) {
	t.Parallel()
	file, err := config.DecodeProject(nil)
	if err != nil {
		t.Fatal(err)
	}
	effective := policy.Resolve(file)
	encoded, err := effective.Canonical()
	if err != nil {
		t.Fatal(err)
	}
	for _, test := range []struct {
		name    string
		profile int
		mutate  func(*journal.Event)
		check   string
	}{
		{"complete snapshot", 21, func(*journal.Event) {}, ""},
		{"before policy schema", 20, func(*journal.Event) {}, "effect-invalid"},
		{"other aggregate", 21, func(event *journal.Event) { event.Aggregate.Type = "issue" }, "effect-invalid"},
		{"wrong digest", 21, func(event *journal.Event) { event.Effects[0].Values["digest"] = "sha256:" + strings.Repeat("0", 64) }, "reference-invalid"},
		{"incomplete policy", 21, func(event *journal.Event) { event.Effects[0].Values["effective_json"] = "{}" }, "reference-invalid"},
		{"unknown lane", 21, changeArchivePolicy(t, func(effective *policy.Effective) {
			effective.Lanes["typo"] = effective.Lanes["assured"]
			delete(effective.Lanes, "assured")
		}), "reference-invalid"},
		{"empty argv", 21, changeArchivePolicy(t, func(effective *policy.Effective) {
			effective.Checks["test"] = config.Check{Run: []string{}, Timeout: time.Minute, DefinedBy: []string{}, WhenChanged: []string{}}
		}), "reference-invalid"},
		{"negative timeout", 21, changeArchivePolicy(t, func(effective *policy.Effective) {
			effective.Checks["test"] = config.Check{Run: []string{"true"}, Timeout: -time.Minute, DefinedBy: []string{}, WhenChanged: []string{}}
		}), "reference-invalid"},
		{"undefined check", 21, changeArchivePolicy(t, func(effective *policy.Effective) {
			effective.Lanes["standard"] = config.Lane{Paths: []string{}, Checks: []string{"missing"}, ChecksSet: true}
		}), "reference-invalid"},
		{"bad defining glob", 21, changeArchivePolicy(t, func(effective *policy.Effective) {
			effective.Checks["test"] = config.Check{Run: []string{"true"}, Timeout: time.Minute, DefinedBy: []string{"../Makefile"}, WhenChanged: []string{}}
		}), "reference-invalid"},
		{"undefined action set", 21, changeArchivePolicy(t, func(effective *policy.Effective) {
			effective.ActorRules["human"] = []policy.ActionSet{"missing"}
		}), "reference-invalid"},
		{"historical built-in values", 21, changeArchivePolicy(t, func(effective *policy.Effective) {
			effective.RulesVersion++
			effective.HardAssured = []string{"legacy-policy/**"}
			effective.ActionSets[policy.WorkSet] = []string{"historical-action"}
			effective.ActionSets[policy.PersonDecisions] = []string{"historical-decision"}
			effective.ActorRules["human"] = []policy.ActionSet{policy.PersonDecisions}
		}), ""},
		{"unrelated effect", 21, func(event *journal.Event) {
			event.Effects = append(event.Effects, journal.Effect{Op: "delete", Table: "labels", PK: map[string]any{"name": "x"}})
		}, "effect-invalid"},
	} {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			validator, err := newSemanticValidator(testMigrationsDir(t), test.profile)
			if err != nil {
				t.Fatal(err)
			}
			defer func() { _ = validator.close() }()
			event := journal.Event{EventSchema: journal.EventSchemaV2, EventID: "EVT-20261112-000001", EventType: "policy.snapshotted",
				OccurredAt: time.Date(2026, 11, 12, 0, 0, 0, 0, time.UTC),
				ActorV2:    journal.EventActor{ActorID: "fid", ActorType: "service", Producer: "fid-standalone", ProvenanceRef: "built-in-policy"},
				Aggregate:  journal.Aggregate{ID: "policy", Sequence: 1, Type: "policy"}, CommandName: "policy.snapshot",
				Effects: []journal.Effect{{Op: "insert", Table: "policy_snapshots", PK: map[string]any{"sequence": 1}, Values: map[string]any{
					"source_ref": "built-in defaults", "source_commit": nil, "digest": effective.Digest(), "effective_json": string(encoded),
				}}}}
			test.mutate(&event)
			raw, err := journal.MarshalCanonicalJSON(event)
			if err != nil {
				t.Fatal(err)
			}
			err = validator.accept(raw)
			if test.check == "" && err != nil || test.check != "" && semanticCheckCode(err) != test.check {
				t.Fatalf("accept = %v, want check %q", err, test.check)
			}
		})
	}
}

func changeArchivePolicy(t *testing.T, mutate func(*policy.Effective)) func(*journal.Event) {
	t.Helper()
	return func(event *journal.Event) {
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
