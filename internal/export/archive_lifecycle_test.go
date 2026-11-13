package export

import (
	"strings"
	"testing"
	"time"

	"github.com/mbell-nzl/fid/internal/holds"
	"github.com/mbell-nzl/fid/internal/journal"
	"github.com/mbell-nzl/fid/internal/project"
)

func takeoverArchiveEvent() journal.Event {
	now := time.Date(2026, 11, 10, 9, 5, 0, 0, time.UTC)
	old := "HLD-1"
	hold := holds.New("HLD-2", "FID-a", project.WorktreeID("/w/new"), "/w/new", "human:mark", now, &old)
	return journal.Event{EventSchema: journal.EventSchemaV2, EventID: "EVT-20261110-000003", EventType: "task.taken-over", OccurredAt: now,
		ActorV2: journal.EventActor{ActorID: "mark", ActorType: "human", Producer: "fid-standalone", ProvenanceRef: "project-local-attribution"}, Aggregate: journal.Aggregate{ID: "FID-a", Type: "task", Sequence: 1}, CommandName: "task.takeover",
		Effects: []journal.Effect{holds.End(old, "takeover", now), {Op: "insert", Table: "person_decisions", PK: map[string]any{"id": "PDN-1"}, Values: map[string]any{"kind": "takeover", "actor": "human:mark", "subject_type": "task", "subject_id": "FID-a", "binding_id": "FID-a", "reason": "Agent stuck", "decided_at": holds.Timestamp(now)}}, hold.Insert()}}
}

func TestArchiveChecksTakeoverAndPersonDecisionBindings(t *testing.T) {
	t.Parallel()
	for _, test := range []struct {
		name   string
		mutate func(*journal.Event)
		closed bool
		valid  bool
	}{
		{name: "complete transaction", valid: true},
		{name: "keyed transaction", valid: true, mutate: func(e *journal.Event) { takeoverArchiveRetry(e) }},
		{name: "retry result bound to another event", mutate: func(e *journal.Event) {
			takeoverArchiveRetry(e)
			e.Effects[3].Values["result_event_id"] = "EVT-20261110-000002"
		}},
		{name: "agent decision", mutate: func(e *journal.Event) { e.ActorV2.ActorType = "agent"; e.Effects[1].Values["actor"] = "agent:mark" }},
		{name: "service decision", mutate: func(e *journal.Event) { e.ActorV2.ActorType = "service"; e.Effects[1].Values["actor"] = "service:mark" }},
		{name: "unregistered decision kind", mutate: func(e *journal.Event) { e.Effects[1].Values["kind"] = "answer" }},
		{name: "another binding", mutate: func(e *journal.Event) { e.Effects[1].Values["binding_id"] = "FID-b" }},
		{name: "missing decision ID", mutate: func(e *journal.Event) { e.Effects[1].PK["id"] = nil }},
		{name: "missing decision", mutate: func(e *journal.Event) { e.Effects = append(e.Effects[:1], e.Effects[2:]...) }},
		{name: "new hold has another actor", mutate: func(e *journal.Event) { e.Effects[2].Values["actor"] = "agent:builder" }},
		{name: "missing predecessor link", mutate: func(e *journal.Event) { e.Effects[2].Values["replaces_hold_id"] = nil }},
		{name: "closed task", closed: true},
		{name: "wrong event command", mutate: func(e *journal.Event) { e.CommandName = "next.claim" }},
	} {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			validator, err := newSemanticValidator(testMigrationsDir(t), 22)
			if err != nil {
				t.Fatal(err)
			}
			defer func() { _ = validator.close() }()
			for _, id := range []string{"FID-a", "FID-b"} {
				if _, err := validator.conn.Exec(`INSERT INTO issues (id,kind,title,workflow_id,status,owner_role,created_at,updated_at) VALUES (?,'task','Work','task-default','todo','developer','2026-11-10T00:00:00Z','2026-11-10T00:00:00Z')`, id); err != nil {
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
			if test.closed {
				if _, err := validator.conn.Exec(`UPDATE issues SET status='done' WHERE id='FID-a'`); err != nil {
					t.Fatal(err)
				}
			}
			event := takeoverArchiveEvent()
			if test.mutate != nil {
				test.mutate(&event)
			}
			raw, err = journal.MarshalCanonicalJSON(event)
			if err != nil {
				t.Fatal(err)
			}
			err = validator.accept(raw)
			if (err == nil) != test.valid {
				t.Fatalf("accept %s = %v, valid=%v", test.name, err, test.valid)
			}
			if !test.valid {
				var decisions int
				if err := validator.conn.QueryRow(`SELECT count(*) FROM person_decisions`).Scan(&decisions); err != nil || decisions != 0 {
					t.Fatalf("failed takeover left %d decisions: %v", decisions, err)
				}
			}
		})
	}
}

func TestArchiveChecksHeartbeatSessionAndExpiry(t *testing.T) {
	t.Parallel()
	for _, test := range []struct {
		name   string
		at     time.Time
		actor  string
		mutate func(*journal.Event)
		valid  bool
	}{
		{name: "renewed", at: time.Date(2026, 11, 10, 9, 5, 0, 0, time.UTC), actor: "builder", valid: true},
		{name: "at expiry", at: time.Date(2026, 11, 10, 10, 0, 0, 0, time.UTC), actor: "builder"},
		{name: "another actor", at: time.Date(2026, 11, 10, 9, 5, 0, 0, time.UTC), actor: "other"},
		{name: "rewrites since", at: time.Date(2026, 11, 10, 9, 5, 0, 0, time.UTC), actor: "builder", mutate: func(e *journal.Event) { e.Effects[0].Values["since"] = "2026-11-10T09:05:00Z" }},
	} {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			validator, err := newSemanticValidator(testMigrationsDir(t), 22)
			if err != nil {
				t.Fatal(err)
			}
			defer func() { _ = validator.close() }()
			if _, err := validator.conn.Exec(`INSERT INTO issues (id,kind,title,workflow_id,status,owner_role,created_at,updated_at) VALUES ('FID-a','task','Work','task-default','todo','developer','2026-11-10T00:00:00Z','2026-11-10T00:00:00Z')`); err != nil {
				t.Fatal(err)
			}
			raw, err := journal.MarshalCanonicalJSON(archiveHoldEvent(func(*journal.Event) {}))
			if err != nil {
				t.Fatal(err)
			}
			if err := validator.accept(raw); err != nil {
				t.Fatal(err)
			}
			event := journal.Event{EventSchema: journal.EventSchemaV2, EventID: "EVT-20261110-000003", EventType: "hold.heartbeat", OccurredAt: test.at, ActorV2: journal.EventActor{ActorID: test.actor, ActorType: "agent", Producer: "fid-standalone", ProvenanceRef: "project-local-attribution"}, Aggregate: journal.Aggregate{ID: "HLD-1", Type: "hold", Sequence: 2}, CommandName: "hold.heartbeat", Effects: []journal.Effect{{Op: "update", Table: "holds", PK: map[string]any{"id": "HLD-1"}, Values: map[string]any{"last_heartbeat_at": holds.Timestamp(test.at), "expires_at": holds.Timestamp(test.at.Add(time.Hour))}, Preconditions: map[string]any{"ended_at": nil, "expires_at": "2026-11-10T10:00:00Z"}}}}
			if test.mutate != nil {
				test.mutate(&event)
			}
			raw, err = journal.MarshalCanonicalJSON(event)
			if err != nil {
				t.Fatal(err)
			}
			err = validator.accept(raw)
			if (err == nil) != test.valid {
				t.Fatalf("accept %s = %v, valid=%v", test.name, err, test.valid)
			}
		})
	}
}

func takeoverArchiveRetry(event *journal.Event) {
	digest := "sha256:" + strings.Repeat("a", 64)
	event.Idempotency = &journal.EventIdempotency{Key: "takeover-1", RequestDigest: digest, Scope: "task.takeover:project-test"}
	event.Effects = append(event.Effects, journal.Effect{Op: "insert", Table: "idempotency_results", PK: map[string]any{"operation_scope": event.Idempotency.Scope, "idempotency_key": event.Idempotency.Key}, Values: map[string]any{"request_digest": digest, "request_field_digests_json": `{"--actor":"` + digest + `"}`, "result_entity_id": event.Aggregate.ID, "result_event_id": event.EventID, "aggregate_sequence": 1, "result_json": `{"fields":{},"occurred_at":"2026-11-10T09:05:00Z","replayed":false}`, "created_at": "2026-11-10T09:05:00Z"}})
}
