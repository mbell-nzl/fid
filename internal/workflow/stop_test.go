package workflow

import (
	"database/sql"
	"testing"
	"time"

	"github.com/mbell-nzl/fid/internal/journal"
)

func TestEquivalentStopsShareAttentionButKeepDistinctFacts(t *testing.T) {
	workspace, _ := initializedWorkflowWorkspace(t)
	conn, err := openWorkflowDB(workspace)
	if err != nil {
		t.Fatalf("open workflow database: %v", err)
	}
	defer closeWorkflowDB(conn)

	first := stopRaisedEvent("EVT-20260905-000001", "STOP-1", "ATT-1", true)
	second := stopRaisedEvent("EVT-20260905-000002", "STOP-2", "ATT-1", false)
	if err := journal.WithLock(workspace, journal.LockOwner{Command: "test stop raise"}, func(lock *journal.WriteLock) error {
		if _, err := journal.CommitEvent(lock, conn, first); err != nil {
			return err
		}
		_, err := journal.CommitEvent(lock, conn, second)
		return err
	}); err != nil {
		t.Fatalf("commit equivalent stop events: %v", err)
	}

	var stops, attention, links int
	if err := conn.QueryRow(`SELECT COUNT(*) FROM stops`).Scan(&stops); err != nil {
		t.Fatalf("count stops: %v", err)
	}
	if err := conn.QueryRow(`SELECT COUNT(*) FROM attention_items WHERE status = 'open'`).Scan(&attention); err != nil {
		t.Fatalf("count open Attention: %v", err)
	}
	if err := conn.QueryRow(`SELECT COUNT(*) FROM attention_stop_links WHERE attention_id = 'ATT-1'`).Scan(&links); err != nil {
		t.Fatalf("count Attention links: %v", err)
	}
	if stops != 2 || attention != 1 || links != 2 {
		t.Fatalf("stops/Attention/links = %d/%d/%d, want 2/1/2", stops, attention, links)
	}
}

func TestRaiseStopUsesOneOpenAttentionAndRejectsRoutineReadsWithoutEffects(t *testing.T) {
	workspace, migrationsDir := initializedWorkflowWorkspace(t)
	conn, err := openWorkflowDB(workspace)
	if err != nil {
		t.Fatalf("open workflow database: %v", err)
	}
	defer closeWorkflowDB(conn)
	installStopTestActor(t, conn)
	request := StopRaiseRequest{
		Subject: Subject{Type: "activation", ID: "ACT-stop-test"}, Code: "operator-stop",
		BlockedActions: []string{"run.recover"}, Details: []byte(`{"reason":"operator intervention","reference":"DEC-188"}`),
		ResolutionCondition: Condition{Type: "activation-recovered", Subject: Subject{Type: "activation", ID: "ACT-stop-test"}}, ResolverCapability: "stop.resolve",
	}
	options := StopOptions{WorkspaceRoot: workspace, MigrationsDir: migrationsDir, ActorID: "agent:stop-test", IdempotencyKey: testIdempotencyKey()}
	first, err := RaiseStop(options, request)
	if err != nil {
		t.Fatalf("RaiseStop(first) error = %v", err)
	}
	options.IdempotencyKey = testIdempotencyKey()
	second, err := RaiseStop(options, request)
	if err != nil {
		t.Fatalf("RaiseStop(second) error = %v", err)
	}
	if first.Stop.ID == second.Stop.ID || first.Attention.ID != second.Attention.ID || len(second.Attention.StopIDs) != 2 {
		t.Fatalf("RaiseStop identities = %#v / %#v, want distinct Stops linked to one Attention", first, second)
	}
	before := readWorkflowTestFile(t, workflowJournalPath(t, workspace))
	if _, err := GetStop(workspace, migrationsDir, "missing"); err == nil {
		t.Fatal("GetStop(missing) error = nil, want rejected read")
	}
	after := readWorkflowTestFile(t, workflowJournalPath(t, workspace))
	if string(before) != string(after) {
		t.Fatal("rejected Stop read produced a Stop or Attention Event")
	}
}

func installStopTestActor(t *testing.T, conn *sql.DB) {
	t.Helper()
	if _, err := conn.Exec(`INSERT INTO kit_activations (id, kit_identity, backend_identity, phase, status, point_of_no_return, direction, created_at, updated_at, version) VALUES ('ACT-stop-test', 'kit', 'backend', 22, 'finalized', 0, 'activate', '2026-09-05T00:00:00Z', '2026-09-05T00:00:00Z', 1)`); err != nil {
		t.Fatalf("seed activation: %v", err)
	}
	if _, err := conn.Exec(`INSERT INTO activation_actor_profiles (activation_id, actor_id, actor_type, producer, provenance_ref, capabilities_json) VALUES ('ACT-stop-test', 'agent:stop-test', 'agent', 'test', 'test:stop', '["attention.resolve","stop.raise","stop.resolve"]')`); err != nil {
		t.Fatalf("seed actor profile: %v", err)
	}
}

func stopRaisedEvent(eventID, stopID, attentionID string, createAttention bool) journal.Event {
	effects := []journal.Effect{{
		Op: "insert", Table: "stops", PK: map[string]any{"id": stopID},
		Values: map[string]any{
			"subject_type": "run", "subject_id": "RUN-1", "code": "operator-stop",
			"blocked_actions_json": `["run.recover"]`, "details_json": `{"reason":"operator intervention","reference":"DEC-188"}`,
			"resolution_condition_json": `{"type":"work-ended","subject":{"type":"run","id":"RUN-1"}}`,
			"resolver_capability":       "stop.resolve", "actor_id": "agent:tl", "raised_at": "2026-09-05T00:00:00Z", "created_event_id": eventID,
		},
	}}
	if createAttention {
		effects = append(effects, journal.Effect{Op: "insert", Table: "attention_items", PK: map[string]any{"id": attentionID}, Values: map[string]any{
			"semantic_key": "same-semantic-key", "subject_type": "run", "subject_id": "RUN-1", "code": "operator-stop", "status": "open", "version": 1, "raised_at": "2026-09-05T00:00:00Z", "resolution_event_id": nil,
		}})
	}
	effects = append(effects, journal.Effect{Op: "insert", Table: "attention_stop_links", PK: map[string]any{"attention_id": attentionID, "stop_id": stopID}})
	return journal.Event{
		EventID: eventID, EventSchema: journal.EventSchemaV2, EventType: "stop.raised", OccurredAt: time.Date(2026, 9, 5, 0, 0, 0, 0, time.UTC),
		ActorV2:   journal.EventActor{ActorID: "agent:tl", ActorType: "agent", Producer: "test", ProvenanceRef: "test"},
		Aggregate: journal.Aggregate{ID: stopID, Type: "stop", Sequence: 1}, CommandName: "stop.raise", Effects: effects,
	}
}
