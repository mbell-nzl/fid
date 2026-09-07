package workflow

import (
	"database/sql"
	"sort"
	"testing"
	"time"

	"github.com/mbell-nzl/fid/internal/journal"
)

func TestAttentionResolutionRequiresExactDispositionEventsAndKeepsStopsImmutable(t *testing.T) {
	workspace, migrationsDir := initializedWorkflowWorkspace(t)
	conn, err := openWorkflowDB(workspace)
	if err != nil {
		t.Fatalf("open workflow database: %v", err)
	}
	defer closeWorkflowDB(conn)
	installStopTestActor(t, conn)
	request := StopRaiseRequest{Subject: Subject{Type: "activation", ID: "ACT-stop-test"}, Code: "operator-stop", BlockedActions: []string{"run.recover"}, Details: []byte(`{"reason":"operator intervention","reference":"DEC-188"}`), ResolutionCondition: Condition{Type: "activation-recovered", Subject: Subject{Type: "activation", ID: "ACT-stop-test"}}, ResolverCapability: "stop.resolve"}
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
	factID := commitActivationRecoveryFact(t, workspace, conn)
	options.IdempotencyKey = testIdempotencyKey()
	firstResolution, err := ResolveStop(options, StopResolveRequest{StopID: first.Stop.ID, FactEventID: factID})
	if err != nil {
		t.Fatalf("ResolveStop(first) error = %v", err)
	}
	replayedStop, err := ResolveStop(options, StopResolveRequest{StopID: first.Stop.ID, FactEventID: factID})
	if err != nil || !replayedStop.Replayed || replayedStop.EventID != firstResolution.EventID || replayedStop.Stop.ID != firstResolution.Stop.ID {
		t.Fatalf("ResolveStop(replay) = (%#v, %v), want original resolution", replayedStop, err)
	}
	options.IdempotencyKey = testIdempotencyKey()
	secondResolution, err := ResolveStop(options, StopResolveRequest{StopID: second.Stop.ID, FactEventID: factID})
	if err != nil {
		t.Fatalf("ResolveStop(second) error = %v", err)
	}
	before := readWorkflowTestFile(t, workflowJournalPath(t, workspace))
	options.IdempotencyKey = testIdempotencyKey()
	if _, err := ResolveAttention(options, AttentionResolveRequest{AttentionID: first.Attention.ID, ExpectedVersion: 1, FactEventIDs: []string{firstResolution.Stop.Resolution.EventID}}); err == nil {
		t.Fatal("ResolveAttention(missing disposition) error = nil")
	}
	if after := readWorkflowTestFile(t, workflowJournalPath(t, workspace)); string(after) != string(before) {
		t.Fatal("rejected Attention resolution changed journal bytes")
	}
	ids := []string{firstResolution.Stop.Resolution.EventID, secondResolution.Stop.Resolution.EventID}
	sort.Strings(ids)
	options.IdempotencyKey = testIdempotencyKey()
	resolved, err := ResolveAttention(options, AttentionResolveRequest{AttentionID: first.Attention.ID, ExpectedVersion: 1, FactEventIDs: ids})
	if err != nil {
		t.Fatalf("ResolveAttention(exact dispositions) error = %v", err)
	}
	if resolved.Attention.Status != "resolved" || resolved.Attention.Version != 2 {
		t.Fatalf("resolved Attention = %#v, want resolved version 2", resolved.Attention)
	}
	replayedAttention, err := ResolveAttention(options, AttentionResolveRequest{AttentionID: first.Attention.ID, ExpectedVersion: 1, FactEventIDs: ids})
	if err != nil || !replayedAttention.Replayed || replayedAttention.EventID != resolved.EventID || replayedAttention.Attention.ID != resolved.Attention.ID {
		t.Fatalf("ResolveAttention(replay) = (%#v, %v), want original resolution", replayedAttention, err)
	}
	for _, id := range []string{first.Stop.ID, second.Stop.ID} {
		stop, err := GetStop(workspace, migrationsDir, id)
		if err != nil || stop.Resolution == nil || stop.Code != "operator-stop" || string(stop.Details) != `{"reason":"operator intervention","reference":"DEC-188"}` {
			t.Fatalf("immutable resolved Stop %q = %#v, %v", id, stop, err)
		}
	}
}

func commitActivationRecoveryFact(t *testing.T, workspace string, db *sql.DB) string {
	t.Helper()
	var eventID string
	err := journal.WithLock(workspace, journal.LockOwner{Command: "test activation recovery"}, func(lock *journal.WriteLock) error {
		event := journal.Event{EventSchema: journal.EventSchemaV2, EventType: "activation.finalized", OccurredAt: time.Now().UTC(), ActorV2: journal.EventActor{ActorID: "agent:stop-test", ActorType: "agent", Producer: "test", ProvenanceRef: "test:attention"}, Aggregate: journal.Aggregate{ID: "ACT-stop-test", Type: "activation", Sequence: 1}, CommandName: "activation.finalize", Effects: []journal.Effect{{Op: "update", Table: "kit_activations", PK: map[string]any{"id": "ACT-stop-test"}, Values: map[string]any{"updated_at": time.Now().UTC().Format(time.RFC3339Nano)}}}}
		appended, err := journal.CommitEvent(lock, db, event)
		eventID = appended.EventID
		return err
	})
	if err != nil {
		t.Fatalf("commit activation recovery fact: %v", err)
	}
	return eventID
}
