package journal

import (
	"errors"
	"strings"
	"testing"
)

// These immutable fields come from the retained DDL, independently of the
// admission allowlist. Pin the shared update refusal even if a later family
// validator or a SQL constraint would also refuse the malformed write.
func TestImmutableLeanColumnsAdmission(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct{ table, columns string }{
		{"check_requests", "id task_id hold_id actor names_json"},
		{"referenced_artifacts", "id digest media_type size_bytes retention_class location_policy created_at expires_at"},
		{"idempotency_results", "operation_scope idempotency_key request_digest created_at request_field_digests_json"},
		{"questions", "id task_id hold_id text asker asked_at"},
		{"holds", "id task_id worktree_id worktree_root actor model since replaces_hold_id"},
	} {
		for _, column := range strings.Fields(tc.columns) {
			t.Run(tc.table+"/"+column, func(t *testing.T) {
				t.Parallel()
				e := admissionLaneRaise()
				e.Effects = []Effect{{Op: "update", Table: tc.table, PK: map[string]any{"id": "row"}, Values: map[string]any{column: "changed"}}}
				err := ValidateEventUpdates(e)
				var refused *RowIdentityError
				if !errors.As(err, &refused) || refused.Table != tc.table || refused.Reason != column+" is immutable in this Event family" {
					t.Fatalf("malformed update = %v, want %s.%s immutable refusal", err, tc.table, column)
				}
			})
		}
	}
}

func TestMutableLeanColumnsAdmission(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct{ table, columns string }{
		{"check_requests", "response_json"},
		{"referenced_artifacts", "availability"},
		{"idempotency_results", "result_entity_id result_event_id aggregate_sequence result_json"},
		{"questions", "resolved_at resolution decision_id"},
		{"holds", "ended_at end_reason last_heartbeat_at expires_at kept park_reason rework_at"},
		{"issues", "title"},
	} {
		for _, column := range strings.Fields(tc.columns) {
			t.Run(tc.table+"/"+column, func(t *testing.T) {
				t.Parallel()
				e := admissionLaneRaise()
				e.Effects = []Effect{{Op: "update", Table: tc.table, Values: map[string]any{column: nil}}}
				if err := ValidateEventUpdates(e); err != nil {
					t.Fatalf("mutable column refused: %v", err)
				}
			})
		}
	}
}

func TestImmutableFactsMayBeInserted(t *testing.T) {
	t.Parallel()
	for _, table := range []string{"holds", "questions", "check_requests", "idempotency_results", "referenced_artifacts"} {
		t.Run(table, func(t *testing.T) {
			t.Parallel()
			e := admissionLaneRaise()
			e.Effects = []Effect{{Op: "insert", Table: table, Values: map[string]any{"id": "new"}}}
			if err := ValidateEventUpdates(e); err != nil {
				t.Fatalf("insert refused by update rule: %v", err)
			}
		})
	}
}
