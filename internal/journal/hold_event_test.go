package journal

import (
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// holdEvent is a hold.taken Event for hold id on task: it ends the holds
// in ended, then records the new one.
func holdEvent(eventID, id, task, worktree string, ended ...string) Event {
	effects := []Effect{}
	for _, previous := range ended {
		effects = append(effects, Effect{Op: "update", Table: "holds", PK: map[string]any{"id": previous},
			Values: map[string]any{"ended_at": "2026-11-10T10:30:00Z", "end_reason": "replaced"}, Preconditions: map[string]any{"ended_at": nil}})
	}
	effects = append(effects, Effect{Op: "insert", Table: "holds", PK: map[string]any{"id": id}, Values: map[string]any{
		"task_id": task, "worktree_id": "sha256:" + worktree, "worktree_root": "/w/" + worktree, "actor": "agent:" + worktree,
		"since": "2026-11-10T10:30:00Z", "expires_at": "2026-11-10T11:30:00Z", "replaces_hold_id": nil, "ended_at": nil, "end_reason": nil}})
	return Event{EventSchema: EventSchemaV2, EventID: eventID, EventType: "hold.taken", OccurredAt: time.Date(2026, 11, 10, 10, 30, 0, 0, time.UTC),
		ActorV2:   EventActor{ActorID: worktree, ActorType: "agent", Producer: "fid-standalone", ProvenanceRef: "project-local-attribution"},
		Aggregate: Aggregate{ID: id, Sequence: 1, Type: "hold"}, CommandName: "next.claim", Effects: effects}
}

// The hold family projects its holds and nothing else, and no other family
// may write a hold except issue.transitioned, which ends one.
func TestHoldEventsProjectOnlyHolds(t *testing.T) {
	t.Parallel()
	conn := openReplayDB(t, filepath.Join(t.TempDir(), "projection.sqlite"))
	defer closeProjectDB(t, conn)
	for _, id := range []string{"FID-a", "FID-b"} {
		if _, err := conn.Exec(`INSERT INTO issues (id, kind, title, workflow_id, status, owner_role, created_at, updated_at, revision)
			VALUES (?, 'task', 'Task', 'task-default', 'todo', 'developer', '2026-11-10T00:00:00Z', '2026-11-10T00:00:00Z', 1)`, id); err != nil {
			t.Fatal(err)
		}
	}
	project := func(event Event) error {
		return ProjectEvent(conn, ProjectionInput{Event: event, EventHash: "sha256:" + strings.Repeat("0", 64), File: "2026-11-10.jsonl", Line: 1})
	}
	if err := project(holdEvent("EVT-20261110-000001", "HLD-1", "FID-a", "one")); err != nil {
		t.Fatalf("first hold: %v", err)
	}
	if err := project(holdEvent("EVT-20261110-000002", "HLD-2", "FID-a", "two")); err == nil {
		t.Fatal("a second open hold on FID-a projected")
	}
	if err := project(holdEvent("EVT-20261110-000003", "HLD-2", "FID-a", "two", "HLD-1")); err != nil {
		t.Fatalf("a hold replacing HLD-1: %v", err)
	}

	smuggled := holdEvent("EVT-20261110-000004", "HLD-3", "FID-b", "three")
	smuggled.Effects = append([]Effect{{Op: "update", Table: "issues", PK: map[string]any{"id": "FID-b"}, Values: map[string]any{"status": "done"}}}, smuggled.Effects...)
	wrongAggregate := holdEvent("EVT-20261110-000005", "HLD-3", "FID-b", "three")
	wrongAggregate.Aggregate.Type = "claim"
	reordered := holdEvent("EVT-20261110-000006", "HLD-3", "FID-b", "three", "HLD-2")
	reordered.Effects[0], reordered.Effects[1] = reordered.Effects[1], reordered.Effects[0]
	labelled := holdEvent("EVT-20261110-000007", "HLD-3", "FID-b", "three")
	labelled.EventType, labelled.Aggregate.Type = "label.added", "label"
	for name, event := range map[string]Event{"an issue effect": smuggled, "another aggregate": wrongAggregate,
		"an end after the insert": reordered, "another family writing a hold": labelled} {
		if err := project(event); err == nil {
			t.Errorf("hold Event with %s projected", name)
		}
	}
	var open int
	if err := conn.QueryRow(`SELECT count(*) FROM holds WHERE ended_at IS NULL`).Scan(&open); err != nil || open != 1 {
		t.Fatalf("open holds = %d, %v; want only HLD-2", open, err)
	}
}
