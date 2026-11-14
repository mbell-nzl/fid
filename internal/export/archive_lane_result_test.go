package export

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/mbell-nzl/fid/internal/journal"
)

func TestLaneUpdateReplayAndArchiveRejectMalformedRetryResult(t *testing.T) {
	t.Parallel()
	for _, row := range []struct {
		name  string
		value any
		valid bool
		from  bool
	}{
		{"valid", "assured", true, false}, {"null", nil, false, false}, {"number", 42, false, false},
		{"empty", "", false, false}, {"unknown", "urgent", false, false}, {"missing", nil, false, false},
		{"from-number", 42, false, true}, {"from-empty", "", false, true}, {"from-unknown", "urgent", false, true}, {"from-unset", nil, true, true},
	} {
		for _, mode := range []string{"replay", "archive"} {
			t.Run(row.name+"/"+mode, func(t *testing.T) {
				t.Parallel()
				validator := startArchiveValidator(t)
				event := journal.Event{EventSchema: journal.EventSchemaV2, EventID: "EVT-20261110-000003", EventType: "issue.updated", CommandName: "issue.update", OccurredAt: time.Date(2026, 11, 10, 9, 5, 0, 0, time.UTC), ActorV2: journal.EventActor{ActorID: "builder", ActorType: "agent", Producer: "fid-standalone", ProvenanceRef: "project-local-attribution"}, Aggregate: journal.Aggregate{ID: "FID-a", Type: "issue", Sequence: 1}, Effects: []journal.Effect{
					{Op: "update", Table: "issues", PK: map[string]any{"id": "FID-a"}, Values: map[string]any{"lane": "assured", "revision": 2}, Preconditions: map[string]any{"lane": "standard", "revision": 1}},
					{Op: "activity", Table: "activity_log", PK: map[string]any{"id": "lane"}, Values: map[string]any{"journal_event_id": "EVT-20261110-000003", "ref_type": "issue", "ref_id": "FID-a", "event": "issue-updated"}},
				}}
				takeoverArchiveRetry(&event)
				event.Idempotency.Scope = "issue.update:project-test"
				last := &event.Effects[len(event.Effects)-1]
				last.PK["operation_scope"] = event.Idempotency.Scope
				change := map[string]any{"field": "lane", "from": "standard", "to": row.value}
				if row.from {
					change["from"], change["to"] = row.value, "assured"
				}
				if row.name == "missing" {
					delete(change, "to")
				}
				result, err := json.Marshal(map[string]any{"issue": map[string]any{"id": "FID-a", "revision": 2}, "changes": []any{change}, "replayed": false})
				if err != nil {
					t.Fatal(err)
				}
				last.Values["result_json"] = string(result)
				if mode == "replay" {
					err = journal.ProjectEvent(validator.conn, journal.ProjectionInput{Event: event, EventHash: "sha256:lane-result", File: "journal", Line: 2, Offset: 2})
				} else {
					var raw []byte
					raw, err = json.Marshal(event)
					if err != nil {
						t.Fatal(err)
					}
					err = validator.accept(raw)
				}
				if (err == nil) != row.valid || (!row.valid && !strings.Contains(err.Error(), "invalid lane change result")) {
					t.Fatalf("%s lane retry result = %v, valid=%v", mode, err, row.valid)
				}
				var lane string
				var revision, activities, bindings int
				if err := validator.conn.QueryRow(`SELECT lane,revision FROM issues WHERE id='FID-a'`).Scan(&lane, &revision); err != nil {
					t.Fatal(err)
				}
				if err := validator.conn.QueryRow(`SELECT (SELECT count(*) FROM activity_log), (SELECT count(*) FROM idempotency_results)`).Scan(&activities, &bindings); err != nil {
					t.Fatal(err)
				}
				if row.valid {
					if lane != "assured" || revision != 2 || activities != 1 || bindings != 1 {
						t.Fatalf("valid update = %s/%d/%d/%d", lane, revision, activities, bindings)
					}
				} else if lane != "standard" || revision != 1 || activities != 0 || bindings != 0 {
					t.Fatalf("rejected update left partial state = %s/%d/%d/%d", lane, revision, activities, bindings)
				}
			})
		}
	}
}
