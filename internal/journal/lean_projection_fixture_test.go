package journal

import (
	"database/sql"
	"testing"
	"time"
)

// The generic prior/idempotency probes start a fixed-state task under a lean
// hold. No kernel schema or configurable workflow is needed.
func leanStartProjectionEvent() Event {
	at := time.Date(2026, 9, 6, 0, 1, 0, 0, time.UTC)
	request := testIdempotencyRequest()
	request.OperationScope = "task.start:project-1"
	event := Event{EventSchema: EventSchemaV2, EventID: "EVT-20260906-000001", EventType: "task.started", OccurredAt: at,
		ActorV2:   EventActor{ActorID: "builder", ActorType: "agent", Producer: "fid-standalone", ProvenanceRef: "project-local-attribution"},
		Aggregate: Aggregate{ID: "FID-start", Type: "task", Sequence: 1}, CommandName: "task.start",
		Idempotency: &EventIdempotency{Key: request.IdempotencyKey, Scope: request.OperationScope, RequestDigest: request.RequestDigest},
		Effects: []Effect{{Op: "update", Table: "issues", PK: map[string]any{"id": "FID-start"},
			Values: map[string]any{"status": "in_progress", "revision": 2, "updated_at": at.Format(time.RFC3339)}, Preconditions: map[string]any{"status": "todo", "revision": 1}}}}
	raw, _ := MarshalCanonicalJSON(map[string]any{"fields": map[string]any{"task_id": "FID-start", "hold_id": "HLD-start", "status": "in_progress", "lane": "standard"}, "occurred_at": at.Format(time.RFC3339), "idempotency_key": request.IdempotencyKey, "replayed": false})
	event.Effects = append(event.Effects, idempotencyResultEffect(request, event, "FID-start", string(raw)))
	return event
}

func seedLeanStartProjection(t *testing.T, conn *sql.DB) {
	t.Helper()
	event := retainedProjectionEvent()
	event.Effects = []Effect{
		{Op: "insert", Table: "issues", PK: map[string]any{"id": "FID-start"}, Values: map[string]any{"kind": "task", "title": "Fixture task", "workflow_id": "task-default", "status": "todo", "lane": "standard", "created_at": "2026-09-06T00:00:00Z", "updated_at": "2026-09-06T00:00:00Z"}},
		{Op: "insert", Table: "holds", PK: map[string]any{"id": "HLD-start"}, Values: map[string]any{"task_id": "FID-start", "worktree_id": "WT-start", "worktree_root": "/fixture", "actor": "agent:builder", "since": "2026-09-06T00:00:00Z", "expires_at": "2026-09-06T01:00:00Z", "last_heartbeat_at": "2026-09-06T00:00:00Z"}},
	}
	hash, err := HashEvent(event)
	if err != nil {
		t.Fatal(err)
	}
	if err := ProjectEvent(conn, ProjectionInput{Event: event, EventHash: hash}); err != nil {
		t.Fatal(err)
	}
}
