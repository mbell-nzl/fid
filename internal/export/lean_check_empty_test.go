package export

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/mbell-nzl/fid/internal/checkrecord"
	"github.com/mbell-nzl/fid/internal/journal"
	"github.com/mbell-nzl/fid/internal/policy"
)

func TestLeanEmptyCheckReplayAndImportRejectCoherentForgeries(t *testing.T) {
	t.Parallel()
	fixture := newExportSourceFixture(t, addLeanExportFacts)
	lines := leanJournalLines(t, fixture.workspace)
	for _, row := range []struct {
		name, want string
	}{
		{"valid", ""},
		{"lane", "empty check requires a started task, lane and current hold"},
		{"unstarted task", "empty check requires a started task, lane and current hold"},
		{"hold actor", "empty check requires a started task, lane and current hold"},
		{"expired hold", "empty check requires a started task, lane and current hold"},
		{"checks not configured", "empty check requires configured lane checks"},
	} {
		t.Run(row.name, func(t *testing.T) {
			t.Parallel()
			for _, mode := range []string{"replay", "import"} {
				t.Run(mode, func(t *testing.T) {
					t.Parallel()
					validator, err := newSemanticValidator(testMigrationsDir(t), 25)
					if err != nil {
						t.Fatal(err)
					}
					defer func() { _ = validator.close() }()
					var completed journal.Event
					for _, raw := range lines {
						var event journal.Event
						if err := json.Unmarshal(raw, &event); err != nil {
							t.Fatal(err)
						}
						if event.EventType == "check.completed" {
							completed = event
							break
						}
						if row.name == "unstarted task" && event.EventType == "task.started" {
							// The authentic hold and policy already exist. Present the
							// empty completion before this task has ever started.
							continue
						}
						if row.name == "unstarted task" && strings.HasPrefix(event.EventType, "check.") {
							continue
						}
						if err := applyLeanEmptyEvent(validator, mode, event); err != nil {
							t.Fatalf("fixture %s: %v", event.EventType, err)
						}
					}
					if completed.EventID == "" {
						t.Fatal("fixture is missing an empty completion")
					}
					forgeLeanEmptyCheck(t, validator, mode, &completed, row.name)
					err = applyLeanEmptyEvent(validator, mode, completed)
					if row.want == "" {
						if err != nil {
							t.Fatalf("valid empty completion: %v", err)
						}
					} else if err == nil || !strings.Contains(err.Error(), row.want) {
						t.Fatalf("%s accepted forgery or rejected it for another reason: %v; want %q", mode, err, row.want)
					}
				})
			}
		})
	}
}

func applyLeanEmptyEvent(validator *semanticValidator, mode string, event journal.Event) error {
	if mode == "import" {
		raw, err := json.Marshal(event)
		if err != nil {
			return err
		}
		return validator.accept(raw)
	}
	hash, err := journal.HashEvent(event)
	if err != nil {
		return err
	}
	return journal.ProjectEvent(validator.conn, journal.ProjectionInput{Event: event, EventHash: hash, File: "journal", Line: 1, Offset: 1})
}

func forgeLeanEmptyCheck(t *testing.T, validator *semanticValidator, mode string, event *journal.Event, forgery string) {
	t.Helper()
	var response checkrecord.Response
	if err := json.Unmarshal([]byte(event.Effects[0].Values["response_json"].(string)), &response); err != nil {
		t.Fatal(err)
	}
	switch forgery {
	case "lane":
		response.Lane = "assured"
	case "hold actor":
		event.ActorV2.ActorID = "another-person"
		event.Effects[0].Values["actor"] = "human:another-person"
	case "expired hold":
		var expires string
		if err := validator.conn.QueryRow(`SELECT expires_at FROM holds WHERE id=?`, event.Effects[0].Values["hold_id"]).Scan(&expires); err != nil {
			t.Fatal(err)
		}
		var err error
		event.OccurredAt, err = time.Parse(time.RFC3339Nano, expires)
		if err != nil {
			t.Fatal(err)
		}
		response.OccurredAt = event.OccurredAt.UTC().Format(time.RFC3339Nano)
	case "checks not configured":
		unsetLeanEmptyChecks(t, validator, mode, *event)
	}
	// Keep both responses identical, including the timestamp in the expiry
	// case, so response-equality validation cannot mask the intended guard.
	raw := encodedLean(response)
	event.Effects[0].Values["response_json"] = raw
	event.Effects[1].Values["result_json"] = raw
}

func unsetLeanEmptyChecks(t *testing.T, validator *semanticValidator, mode string, completed journal.Event) {
	t.Helper()
	var sequence int
	var raw string
	if err := validator.conn.QueryRow(`SELECT sequence,effective_json FROM policy_snapshots ORDER BY sequence DESC LIMIT 1`).Scan(&sequence, &raw); err != nil {
		t.Fatal(err)
	}
	var effective policy.Effective
	if err := json.Unmarshal([]byte(raw), &effective); err != nil {
		t.Fatal(err)
	}
	lane := effective.Lanes["standard"]
	lane.ChecksSet, lane.Checks = false, []string{}
	effective.Lanes["standard"] = lane
	encoded, err := effective.Canonical()
	if err != nil {
		t.Fatal(err)
	}
	// Add a valid policy change, rather than corrupting an earlier snapshot or
	// its digest. Only the following empty completion should be rejected.
	snapshot := journal.Event{EventSchema: journal.EventSchemaV2, EventID: "EVT-20261114-999999", EventType: "policy.snapshotted",
		OccurredAt: completed.OccurredAt, CommandName: "policy.snapshot",
		ActorV2:   journal.EventActor{ActorID: "fid", ActorType: "service", Producer: "fid-standalone", ProvenanceRef: "built-in-policy"},
		Aggregate: journal.Aggregate{ID: "policy", Sequence: sequence + 1, Type: "policy"},
		Effects: []journal.Effect{{Op: "insert", Table: "policy_snapshots", PK: map[string]any{"sequence": sequence + 1},
			Values: map[string]any{"source_ref": "built-in defaults", "source_commit": nil, "digest": effective.Digest(), "effective_json": string(encoded)}}}}
	if err := applyLeanEmptyEvent(validator, mode, snapshot); err != nil {
		t.Fatalf("fixture policy change: %v", err)
	}
}
