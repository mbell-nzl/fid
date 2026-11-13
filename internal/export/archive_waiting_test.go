package export

import (
	"testing"
	"time"

	"github.com/mbell-nzl/fid/internal/holds"
	"github.com/mbell-nzl/fid/internal/journal"
)

func waitingArchiveEvent() journal.Event {
	now := time.Date(2026, 11, 10, 9, 5, 0, 0, time.UTC)
	effect := holds.End("HLD-1", "waiting", now)
	effect.Values["kept"] = true
	event := journal.Event{EventSchema: journal.EventSchemaV2, EventID: "EVT-20261110-000003", EventType: "hold.waiting", OccurredAt: now, ActorV2: journal.EventActor{ActorID: "builder", ActorType: "agent", Producer: "fid-standalone", ProvenanceRef: "project-local-attribution"}, Aggregate: journal.Aggregate{ID: "HLD-1", Type: "hold", Sequence: 2}, CommandName: "hold.park", Effects: []journal.Effect{effect, {Op: "insert", Table: "questions", PK: map[string]any{"id": "QST-1"}, Values: map[string]any{"task_id": "FID-a", "hold_id": "HLD-1", "text": "Include timeouts?", "asker": "agent:builder", "asked_at": holds.Timestamp(now), "resolved_at": nil, "resolution": nil, "decision_id": nil}}}}
	takeoverArchiveRetry(&event)
	event.Idempotency.Scope = "hold.park:project-test"
	event.Effects[2].PK["operation_scope"] = event.Idempotency.Scope
	event.Effects[2].Values["aggregate_sequence"] = 2
	return event
}

func answerArchiveEvent() journal.Event {
	now := time.Date(2026, 11, 10, 9, 6, 0, 0, time.UTC)
	return journal.Event{EventSchema: journal.EventSchemaV2, EventID: "EVT-20261110-000004", EventType: "task.answered", OccurredAt: now, ActorV2: journal.EventActor{ActorID: "mark", ActorType: "human", Producer: "fid-standalone", ProvenanceRef: "project-local-attribution"}, Aggregate: journal.Aggregate{ID: "FID-a", Type: "task", Sequence: 1}, CommandName: "task.answer", Effects: []journal.Effect{{Op: "insert", Table: "person_decisions", PK: map[string]any{"id": "PDN-1"}, Values: map[string]any{"kind": "answer", "actor": "human:mark", "subject_type": "task", "subject_id": "FID-a", "reason": "Yes", "binding_id": "QST-1", "decided_at": holds.Timestamp(now)}}, {Op: "update", Table: "questions", PK: map[string]any{"id": "QST-1"}, Values: map[string]any{"resolved_at": holds.Timestamp(now), "resolution": "answered", "decision_id": "PDN-1"}, Preconditions: map[string]any{"resolved_at": nil}}}}
}

func TestArchiveWaitingAndAnswerBindings(t *testing.T) {
	t.Parallel()
	for _, test := range []struct {
		name    string
		waiting bool
		valid   bool
		mutate  func(*journal.Event)
	}{
		{name: "complete", valid: true},
		{name: "question task differs", waiting: true, mutate: func(e *journal.Event) { e.Effects[1].Values["task_id"] = "FID-b" }},
		{name: "question asker differs", waiting: true, mutate: func(e *journal.Event) { e.Effects[1].Values["asker"] = "agent:other" }},
		{name: "empty question", waiting: true, mutate: func(e *journal.Event) { e.Effects[1].Values["text"] = " " }},
		{name: "reservation not kept", waiting: true, mutate: func(e *journal.Event) { e.Effects[0].Values["kept"] = false }},
		{name: "missing question", waiting: true, mutate: func(e *journal.Event) { e.Effects = append(e.Effects[:1], e.Effects[2:]...) }},
		{name: "wrong decision binding", mutate: func(e *journal.Event) { e.Effects[0].Values["binding_id"] = "FID-a" }},
		{name: "wrong subject", mutate: func(e *journal.Event) { e.Effects[0].Values["subject_id"] = "FID-b" }},
		{name: "wrong decision link", mutate: func(e *journal.Event) { e.Effects[1].Values["decision_id"] = "PDN-other" }},
		{name: "rewritten question", mutate: func(e *journal.Event) { e.Effects[1].Values["text"] = "New question" }},
		{name: "agent answer", mutate: func(e *journal.Event) { e.ActorV2.ActorType = "agent"; e.Effects[0].Values["actor"] = "agent:mark" }},
		{name: "missing resolution", mutate: func(e *journal.Event) { e.Effects = e.Effects[:1] }},
	} {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			validator, err := newSemanticValidator(testMigrationsDir(t), 23)
			if err != nil {
				t.Fatal(err)
			}
			defer func() { _ = validator.close() }()
			for _, id := range []string{"FID-a", "FID-b"} {
				if _, err := validator.conn.Exec(`INSERT INTO issues (id,kind,title,workflow_id,status,owner_role,created_at,updated_at) VALUES (?,'task','Work','task-default','todo','developer','2026-11-10T00:00:00Z','2026-11-10T00:00:00Z')`, id); err != nil {
					t.Fatal(err)
				}
			}
			accept := func(e journal.Event) error {
				raw, err := journal.MarshalCanonicalJSON(e)
				if err != nil {
					t.Fatal(err)
				}
				return validator.accept(raw)
			}
			if err := accept(archiveHoldEvent(func(*journal.Event) {})); err != nil {
				t.Fatal(err)
			}
			waiting := waitingArchiveEvent()
			if test.waiting {
				test.mutate(&waiting)
			}
			err = accept(waiting)
			if test.waiting {
				if err == nil {
					t.Fatal("invalid waiting event accepted")
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			answer := answerArchiveEvent()
			if test.mutate != nil {
				test.mutate(&answer)
			}
			err = accept(answer)
			if (err == nil) != test.valid {
				t.Fatalf("answer error=%v valid=%v", err, test.valid)
			}
			if test.valid {
				answer.EventID = "EVT-20261110-000005"
				answer.Aggregate.Sequence++
				if err := accept(answer); err == nil {
					t.Fatal("already answered question accepted")
				}
			}
		})
	}
}
