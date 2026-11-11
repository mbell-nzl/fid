package export

import (
	"testing"
	"time"

	"github.com/mbell-nzl/fid/internal/journal"
)

func archiveHoldEvent(mutate func(*journal.Event)) journal.Event {
	event := journal.Event{EventSchema: journal.EventSchemaV2, EventID: "EVT-20261110-000002", EventType: "hold.taken",
		OccurredAt: time.Date(2026, 11, 10, 9, 0, 0, 0, time.UTC),
		ActorV2:    journal.EventActor{ActorID: "builder", ActorType: "agent", Producer: "fid-standalone", ProvenanceRef: "project-local-attribution"},
		Aggregate:  journal.Aggregate{ID: "HLD-1", Sequence: 1, Type: "hold"}, CommandName: "next.claim",
		Effects: []journal.Effect{{Op: "insert", Table: "holds", PK: map[string]any{"id": "HLD-1"}, Values: map[string]any{
			"task_id": "FID-a", "worktree_id": "sha256:one", "worktree_root": "/w/one", "actor": "agent:builder",
			"since": "2026-11-10T09:00:00Z", "expires_at": "2026-11-10T10:00:00Z", "replaces_hold_id": nil, "ended_at": nil, "end_reason": nil}}}}
	mutate(&event)
	return event
}

// An archive's hold Events are checked against the hold family: archives
// from before migration 20 hold none, and a hold Event carries only its
// holds, in order, with the new hold as its aggregate.
func TestSemanticValidatorChecksHoldEvents(t *testing.T) {
	t.Parallel()
	for _, test := range []struct {
		name    string
		profile int
		mutate  func(*journal.Event)
		check   string // the refusal's check code; empty when accepted
	}{
		{"nothing changed", 20, func(*journal.Event) {}, ""},
		{"an end before the new hold", 20, func(event *journal.Event) {
			event.Effects = append([]journal.Effect{{Op: "update", Table: "holds", PK: map[string]any{"id": "HLD-0"},
				Values: map[string]any{"ended_at": "2026-11-10T09:00:00Z", "end_reason": "replaced"}, Preconditions: map[string]any{"ended_at": nil}}}, event.Effects...)
		}, ""},
		{"before migration 20", 19, func(*journal.Event) {}, "effect-invalid"},
		{"an issue effect", 20, func(event *journal.Event) {
			event.Effects = append(event.Effects, journal.Effect{Op: "update", Table: "issues", PK: map[string]any{"id": "FID-a"}, Values: map[string]any{"status": "done"}})
		}, "effect-invalid"},
		{"an end after the new hold", 20, func(event *journal.Event) {
			event.Effects = append(event.Effects, journal.Effect{Op: "update", Table: "holds", PK: map[string]any{"id": "HLD-0"},
				Values: map[string]any{"ended_at": "2026-11-10T09:00:00Z", "end_reason": "replaced"}, Preconditions: map[string]any{"ended_at": nil}})
		}, "effect-invalid"},
		{"a changed task", 20, func(event *journal.Event) {
			event.Effects = append([]journal.Effect{{Op: "update", Table: "holds", PK: map[string]any{"id": "HLD-0"},
				Values: map[string]any{"task_id": "FID-b"}}}, event.Effects...)
		}, "effect-invalid"},
		{"a new hold other than its aggregate", 20, func(event *journal.Event) { event.Aggregate.ID = "HLD-2" }, "reference-invalid"},
		{"another aggregate", 20, func(event *journal.Event) { event.Aggregate.Type = "claim" }, "effect-invalid"},
	} {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			validator, err := newSemanticValidator(testMigrationsDir(t), test.profile)
			if err != nil {
				t.Fatalf("new semantic validator: %v", err)
			}
			defer func() { _ = validator.close() }()
			for _, id := range []string{"FID-a", "FID-b"} {
				if _, err := validator.conn.Exec(`INSERT INTO issues (id, kind, title, workflow_id, status, owner_role, created_at, updated_at)
					VALUES (?, 'task', 'Task', 'task-default', 'todo', 'developer', '2026-11-10T00:00:00Z', '2026-11-10T00:00:00Z')`, id); err != nil {
					t.Fatal(err)
				}
			}
			if test.profile == 20 {
				if _, err := validator.conn.Exec(`INSERT INTO holds (id, task_id, worktree_id, worktree_root, actor, since, expires_at)
					VALUES ('HLD-0', 'FID-b', 'sha256:two', '/w/two', 'agent:two', '2026-11-10T08:00:00Z', '2026-11-10T09:00:00Z')`); err != nil {
					t.Fatal(err)
				}
			}
			raw, err := journal.MarshalCanonicalJSON(archiveHoldEvent(test.mutate))
			if err != nil {
				t.Fatal(err)
			}
			err = validator.accept(raw)
			if test.check == "" && err != nil {
				t.Fatalf("accept hold Event with %s = %v", test.name, err)
			}
			if code := semanticCheckCode(err); test.check != "" && code != test.check {
				t.Fatalf("accept hold Event with %s = %v (check %q), want %s", test.name, err, code, test.check)
			}
		})
	}
}
