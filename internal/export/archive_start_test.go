package export

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/mbell-nzl/fid/internal/journal"
)

func archiveStartEvent(t *testing.T) journal.Event {
	t.Helper()
	now := time.Date(2026, 11, 10, 9, 5, 0, 0, time.UTC)
	event := journal.Event{EventSchema: journal.EventSchemaV2, EventID: "EVT-20261110-000003", EventType: "task.started", CommandName: "task.start", OccurredAt: now, ActorV2: journal.EventActor{ActorID: "builder", ActorType: "agent", Producer: "fid-standalone", ProvenanceRef: "project-local-attribution"}, Aggregate: journal.Aggregate{ID: "FID-a", Type: "task", Sequence: 1}}
	for _, id := range []string{"FID-a", "FID-epic"} {
		event.Effects = append(event.Effects, journal.Effect{Op: "update", Table: "issues", PK: map[string]any{"id": id}, Values: map[string]any{"status": "in_progress", "revision": 2, "updated_at": now.Format(time.RFC3339)}, Preconditions: map[string]any{"status": "todo", "revision": 1}})
	}
	takeoverArchiveRetry(&event)
	event.Idempotency.Scope = "task.start:project-test"
	last := &event.Effects[len(event.Effects)-1]
	last.PK["operation_scope"] = event.Idempotency.Scope
	raw, err := journal.MarshalCanonicalJSON(map[string]any{"fields": map[string]any{"task_id": "FID-a", "hold_id": "HLD-1", "status": "in_progress", "lane": "standard"}, "occurred_at": now.Format(time.RFC3339), "replayed": false})
	if err != nil {
		t.Fatal(err)
	}
	last.Values["result_json"] = string(raw)
	return event
}

func startArchiveValidator(t *testing.T) *semanticValidator {
	t.Helper()
	validator, err := newSemanticValidator(testMigrationsDir(t), 24)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = validator.close() })
	for _, query := range []string{
		`INSERT INTO issues(id,kind,title,workflow_id,status,owner_role,created_at,updated_at) VALUES('FID-epic','epic','Epic','epic-default','todo','tl','2026-11-10T00:00:00Z','2026-11-10T00:00:00Z')`,
		`INSERT INTO epic_metadata(issue_id,scheduling) VALUES('FID-epic','prioritised')`,
		`INSERT INTO issues(id,kind,parent_id,title,workflow_id,status,lane,owner_role,created_at,updated_at) VALUES('FID-a','task','FID-epic','Task','task-default','todo','standard','developer','2026-11-10T00:00:00Z','2026-11-10T00:00:00Z')`,
	} {
		if _, err := validator.conn.Exec(query); err != nil {
			t.Fatal(err)
		}
	}
	raw, err := journal.MarshalCanonicalJSON(archiveHoldEvent(func(*journal.Event) {}))
	if err != nil {
		t.Fatal(err)
	}
	if err := validator.accept(raw); err != nil {
		t.Fatal(err)
	}
	return validator
}

func TestArchiveStartBindsHoldLaneTaskAndEpicAtomically(t *testing.T) {
	t.Parallel()
	for _, row := range []struct {
		name   string
		valid  bool
		sql    string
		mutate func(*journal.Event)
	}{
		{name: "complete start", valid: true},
		{name: "no epic effect", mutate: func(e *journal.Event) { e.Effects = append(e.Effects[:1], e.Effects[2:]...) }},
		{name: "another actor", mutate: func(e *journal.Event) { e.ActorV2.ActorID = "other" }},
		{name: "service actor", mutate: func(e *journal.Event) { e.ActorV2.ActorType = "service" }},
		{name: "another epic", mutate: func(e *journal.Event) { e.Effects[1].PK["id"] = "FID-a" }},
		{name: "unrelated scalar", mutate: func(e *journal.Event) { e.Effects[0].Values["priority"] = 1 }},
		{name: "wrong revision", mutate: func(e *journal.Event) { e.Effects[0].Values["revision"] = 3 }},
		{name: "wrong command", mutate: func(e *journal.Event) { e.CommandName = "issue.transition" }},
		{name: "missing lane", sql: `UPDATE issues SET lane=NULL WHERE id='FID-a'`},
		{name: "already in review", sql: `UPDATE issues SET status='in_review' WHERE id='FID-a'`},
		{name: "expired hold", sql: `UPDATE holds SET expires_at='2026-11-10T09:05:00Z' WHERE id='HLD-1'`},
		{name: "ended hold", sql: `UPDATE holds SET ended_at='2026-11-10T09:01:00Z',end_reason='takeover' WHERE id='HLD-1'`},
		{name: "retry result differs", mutate: func(e *journal.Event) {
			e.Effects[2].Values["result_json"] = `{"fields":{"task_id":"FID-a","hold_id":"HLD-1","status":"in_progress","lane":"assured"},"occurred_at":"2026-11-10T09:05:00Z","replayed":false}`
		}},
	} {
		t.Run(row.name, func(t *testing.T) {
			t.Parallel()
			validator := startArchiveValidator(t)
			if row.sql != "" {
				if _, err := validator.conn.Exec(row.sql); err != nil {
					t.Fatal(err)
				}
			}
			var before string
			if err := validator.conn.QueryRow(`SELECT status FROM issues WHERE id='FID-a'`).Scan(&before); err != nil {
				t.Fatal(err)
			}
			event := archiveStartEvent(t)
			if row.mutate != nil {
				row.mutate(&event)
			}
			raw, err := journal.MarshalCanonicalJSON(event)
			if err != nil {
				t.Fatal(err)
			}
			err = validator.accept(raw)
			if (err == nil) != row.valid {
				t.Fatalf("accept %s = %v, valid=%v", row.name, err, row.valid)
			}
			var task, epic string
			var bindings int
			if err := validator.conn.QueryRow(`SELECT status FROM issues WHERE id='FID-a'`).Scan(&task); err != nil {
				t.Fatal(err)
			}
			if err := validator.conn.QueryRow(`SELECT status FROM issues WHERE id='FID-epic'`).Scan(&epic); err != nil {
				t.Fatal(err)
			}
			if err := validator.conn.QueryRow(`SELECT count(*) FROM idempotency_results`).Scan(&bindings); err != nil {
				t.Fatal(err)
			}
			if row.valid {
				if task != "in_progress" || epic != "in_progress" || bindings != 1 {
					t.Fatalf("start state = %s/%s/%d", task, epic, bindings)
				}
			} else if task != before || epic != "todo" || bindings != 0 {
				t.Fatalf("rejected start left partial state %s/%s/%d", task, epic, bindings)
			}
		})
	}
}

func TestStartReplayRejectsMalformedRetryResultWithoutMutation(t *testing.T) {
	t.Parallel()
	for _, row := range []struct {
		name  string
		value any
	}{
		{"missing", nil},
		{"number", 42},
		{"null", nil},
	} {
		t.Run(row.name, func(t *testing.T) {
			t.Parallel()
			validator := startArchiveValidator(t)
			event := archiveStartEvent(t)
			result := event.Effects[len(event.Effects)-1].Values
			result["result_json"] = row.value
			if row.name == "missing" {
				delete(result, "result_json")
			}
			err := journal.ProjectEvent(validator.conn, journal.ProjectionInput{Event: event, EventHash: "sha256:malformed-result", File: "journal", Line: 2, Offset: 2})
			if err == nil || !strings.Contains(err.Error(), "invalid original idempotency result") {
				t.Fatalf("replay malformed retry result = %v", err)
			}
			var task, epic string
			var bindings int
			if err := validator.conn.QueryRow(`SELECT status FROM issues WHERE id='FID-a'`).Scan(&task); err != nil {
				t.Fatal(err)
			}
			if err := validator.conn.QueryRow(`SELECT status FROM issues WHERE id='FID-epic'`).Scan(&epic); err != nil {
				t.Fatal(err)
			}
			if err := validator.conn.QueryRow(`SELECT count(*) FROM idempotency_results`).Scan(&bindings); err != nil {
				t.Fatal(err)
			}
			if task != "todo" || epic != "todo" || bindings != 0 {
				t.Fatalf("rejected replay left partial state %s/%s/%d", task, epic, bindings)
			}
		})
	}
}

func TestArchiveLaneLoweringKeepsThePlanningAuthorityBoundary(t *testing.T) {
	t.Parallel()
	for _, row := range []struct {
		actor, status, lane string
		valid               bool
	}{
		{"agent", "todo", "assured", true}, {"agent", "todo", "routine", false}, {"human", "todo", "routine", true}, {"human", "in_progress", "routine", false}, {"service", "todo", "assured", false},
	} {
		t.Run(row.actor+row.status+row.lane, func(t *testing.T) {
			t.Parallel()
			validator := startArchiveValidator(t)
			if _, err := validator.conn.Exec(`UPDATE issues SET status=? WHERE id='FID-a'`, row.status); err != nil {
				t.Fatal(err)
			}
			event := journal.Event{EventSchema: journal.EventSchemaV2, EventID: "EVT-20261110-000003", EventType: "issue.updated", CommandName: "issue.update", OccurredAt: time.Date(2026, 11, 10, 9, 5, 0, 0, time.UTC), ActorV2: journal.EventActor{ActorID: "builder", ActorType: row.actor, Producer: "fid-standalone", ProvenanceRef: "project-local-attribution"}, Aggregate: journal.Aggregate{ID: "FID-a", Type: "issue", Sequence: 1}, Effects: []journal.Effect{
				{Op: "update", Table: "issues", PK: map[string]any{"id": "FID-a"}, Values: map[string]any{"lane": row.lane}, Preconditions: map[string]any{"lane": "standard"}},
				{Op: "activity", Table: "activity_log", PK: map[string]any{"id": "lane"}, Values: map[string]any{"journal_event_id": "EVT-20261110-000003", "ref_type": "issue", "ref_id": "FID-a", "event": "issue-updated"}},
			}}
			raw, err := json.Marshal(event)
			if err != nil {
				t.Fatal(err)
			}
			err = validator.accept(raw)
			if (err == nil) != row.valid {
				t.Fatalf("lane import = %v, valid=%v", err, row.valid)
			}
			var lane string
			if err := validator.conn.QueryRow(`SELECT lane FROM issues WHERE id='FID-a'`).Scan(&lane); err != nil {
				t.Fatal(err)
			}
			want := "standard"
			if row.valid {
				want = row.lane
			}
			if lane != want {
				t.Fatalf("lane after import = %s, want %s", lane, want)
			}
		})
	}
}
