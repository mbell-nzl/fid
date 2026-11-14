package export

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/mbell-nzl/fid/internal/journal"
)

func recoveryEventWithMissingReference() journal.Event {
	encode := func(value any) string {
		raw, err := json.Marshal(value)
		if err != nil {
			panic(err)
		}
		return string(raw)
	}
	repository := map[string]any{"project_id": "PRJ-missing", "worktree_id": "WT-missing", "capture_source": "git", "head_oid": nil, "tree_oid": nil, "base_oid": nil, "dirty": false, "dirty_digest": nil}
	scope := map[string]any{"paths": []any{}}
	reconciliation := map[string]any{"repository": repository, "checkpoint_id": nil, "paths": []any{}, "unmerged_paths": []any{}, "operation": "none", "ancestry": "unchanged"}
	reconciliation["digest"] = fmt.Sprintf("sha256:%x", sha256.Sum256([]byte(encode(reconciliation))))
	authority := map[string]any{"mode": "ordinary", "actor_id": "agent:recovery", "reason": "semantic reference test", "reference": nil, "grant_ref": nil}
	event := journal.Event{
		EventID:     "EVT-invalid-recovery-reference",
		EventSchema: journal.EventSchemaV2,
		EventType:   "run.recovered",
		OccurredAt:  time.Date(2026, time.September, 8, 0, 0, 0, 0, time.UTC),
		ActorV2: journal.EventActor{
			ActorID:       "agent:recovery",
			ActorType:     "agent",
			Producer:      "fid",
			ProvenanceRef: "test:archive-semantic-reference",
		},
		Aggregate:   journal.Aggregate{ID: "RUN-prior-missing", Sequence: 1, Type: "run"},
		CommandName: "run.recover",
		Effects: []journal.Effect{
			{Op: "insert", Table: "claims", PK: map[string]any{"id": "CLM-successor-missing"}, Values: map[string]any{
				"issue_id": "ISS-missing", "envelope_id": "ENV-missing", "envelope_revision": 1, "actor_id": "agent:recovery", "worktree_id": "WT-missing",
				"scope_json": encode(scope), "mode": "mutation", "version": 1, "status": "active", "acquired_at": "2026-09-08T00:00:00Z", "expires_at": "2026-09-08T01:00:00Z", "released_at": nil, "release_reason": nil,
			}},
			{Op: "insert", Table: "runs", PK: map[string]any{"id": "RUN-successor-missing"}, Values: map[string]any{
				"issue_id": "ISS-missing", "envelope_id": "ENV-missing", "envelope_revision": 1, "claim_id": "CLM-successor-missing", "claim_version": 1, "actor_id": "agent:recovery",
				"start_repository_json": encode(repository), "terminal_repository_json": nil, "status": "active", "started_at": "2026-09-08T00:00:00Z", "ended_at": nil,
			}},
			{Op: "insert", Table: "run_runtime", PK: map[string]any{"run_id": "RUN-successor-missing"}, Values: map[string]any{"version": 1, "executor_id": "EXEC-successor", "obligation_id": nil}},
			{Op: "insert", Table: "run_recoveries", PK: map[string]any{"id": "REC-invalid-reference"}, Values: map[string]any{
				"prior_run_id": "RUN-prior-missing", "successor_run_id": "RUN-successor-missing", "prior_claim_id": "CLM-prior-missing", "successor_claim_id": "CLM-successor-missing", "observation_id": "OBS-missing",
				"reconciliation_json": encode(reconciliation), "adopted_paths_json": "[]", "authority_json": encode(authority), "reason": "semantic reference test", "created_event_id": "EVT-invalid-recovery-reference",
			}},
		},
	}
	// Complete nested values ensure this fixture reaches reference validation.
	claim := map[string]any{"id": "CLM-successor-missing", "issue_id": "ISS-missing", "envelope_id": "ENV-missing", "envelope_revision": 1, "actor_id": "agent:recovery", "worktree_id": "WT-missing", "scope": scope, "mode": "mutation", "version": 1, "status": "active", "acquired_at": "2026-09-08T00:00:00Z", "expires_at": "2026-09-08T01:00:00Z", "released_at": nil, "release_reason": nil}
	run := map[string]any{"id": "RUN-successor-missing", "issue_id": "ISS-missing", "envelope_id": "ENV-missing", "envelope_revision": 1, "claim_id": "CLM-successor-missing", "claim_version": 1, "actor_id": "agent:recovery", "start_repository": repository, "terminal_repository": nil, "status": "active", "budgets": []any{}, "started_at": "2026-09-08T00:00:00Z", "ended_at": nil, "version": 1, "executor": map[string]any{"id": "EXEC-successor", "source": "activated-harness"}, "obligation_id": nil, "predecessor_run_id": "RUN-prior-missing", "successor_run_id": nil}
	recovery := map[string]any{"id": "REC-invalid-reference", "prior_run_id": "RUN-prior-missing", "successor_run_id": "RUN-successor-missing", "prior_claim_id": "CLM-prior-missing", "successor_claim_id": "CLM-successor-missing", "observation_id": "OBS-missing", "reconciliation": reconciliation, "adopted_paths": []any{}, "authority": authority, "event_id": event.EventID}
	result := map[string]any{"run": run, "claim": claim, "recovery": recovery, "idempotency_key": "missing-reference", "replayed": false}
	event.Idempotency = &journal.EventIdempotency{Key: "missing-reference", Scope: "run.recover:missing-project", RequestDigest: strings.Repeat("1", 64)}
	event.Effects = append(event.Effects, journal.Effect{Op: "insert", Table: "idempotency_results", PK: map[string]any{"operation_scope": event.Idempotency.Scope, "idempotency_key": event.Idempotency.Key}, Values: map[string]any{"request_digest": event.Idempotency.RequestDigest, "request_field_digests_json": "{}", "result_entity_id": "REC-invalid-reference", "result_event_id": event.EventID, "aggregate_sequence": 1, "result_json": encode(result), "created_at": event.OccurredAt.Format(time.RFC3339Nano)}})
	return event
}

func TestSemanticValidatorRejectsMissingRecoveryReference(t *testing.T) {
	t.Parallel()
	validator, err := newSemanticValidator(testMigrationsDir(t), 19)
	if err != nil {
		t.Fatalf("new semantic validator: %v", err)
	}
	defer func() { _ = validator.close() }()

	raw, err := journal.MarshalCanonicalJSON(recoveryEventWithMissingReference())
	if err != nil {
		t.Fatalf("encode recovery Event: %v", err)
	}
	err = validator.accept(raw)
	if err == nil || !strings.Contains(err.Error(), "missing reference") {
		t.Fatalf("accept missing recovery reference error = %v, want foreign-key rejection", err)
	}
	if got := semanticCheckCode(err); got != "reference-invalid" {
		t.Fatalf("semantic check code = %q, want reference-invalid", got)
	}
}

func TestSemanticValidatorRejectsRecoveryEffectFamilyTampering(t *testing.T) {
	t.Parallel()
	for _, testCase := range []struct {
		name   string
		mutate func(*journal.Event)
	}{
		{
			name: "wrong operation",
			mutate: func(event *journal.Event) {
				event.Effects[0].Op = "update"
			},
		},
		{
			name: "reordered successor run",
			mutate: func(event *journal.Event) {
				event.Effects[0], event.Effects[1] = event.Effects[1], event.Effects[0]
			},
		},
		{
			name: "wrong aggregate",
			mutate: func(event *journal.Event) {
				event.Aggregate.Type = "issue"
			},
		},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			validator, err := newSemanticValidator(testMigrationsDir(t), 17)
			if err != nil {
				t.Fatalf("new semantic validator: %v", err)
			}
			defer func() { _ = validator.close() }()

			event := recoveryEventWithMissingReference()
			testCase.mutate(&event)
			raw, err := journal.MarshalCanonicalJSON(event)
			if err != nil {
				t.Fatalf("encode tampered recovery Event: %v", err)
			}
			err = validator.accept(raw)
			if err == nil {
				t.Fatal("accept tampered recovery Event error = nil")
			}
			if got := semanticCheckCode(err); got != "effect-invalid" {
				t.Fatalf("semantic check code = %q, want effect-invalid; error = %v", got, err)
			}
		})
	}
}
