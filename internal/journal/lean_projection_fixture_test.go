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

// Admission refuses Event-v1 writes to holds, so the fixture seeds its rows
// directly, as the archive fixtures do.
func seedLeanStartProjection(t *testing.T, conn *sql.DB) {
	t.Helper()
	for _, statement := range []string{
		`INSERT INTO issues (id, kind, title, workflow_id, status, lane, created_at, updated_at) VALUES ('FID-start', 'task', 'Fixture task', 'task-default', 'todo', 'standard', '2026-09-06T00:00:00Z', '2026-09-06T00:00:00Z')`,
		`INSERT INTO holds (id, task_id, worktree_id, worktree_root, actor, since, expires_at, last_heartbeat_at) VALUES ('HLD-start', 'FID-start', 'WT-start', '/fixture', 'agent:builder', '2026-09-06T00:00:00Z', '2026-09-06T01:00:00Z', '2026-09-06T00:00:00Z')`,
	} {
		if _, err := conn.Exec(statement); err != nil {
			t.Fatal(err)
		}
	}
}
