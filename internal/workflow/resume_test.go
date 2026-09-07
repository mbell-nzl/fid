package workflow

import "testing"

func TestResumeKeepsOpenAttentionAfterAllLinkedStopsAreResolved(t *testing.T) {
	workspace, migrationsDir := initializedWorkflowWorkspace(t)
	conn, err := openWorkflowDB(workspace)
	if err != nil {
		t.Fatalf("open workflow database: %v", err)
	}
	defer closeWorkflowDB(conn)
	installStopTestActor(t, conn)
	request := StopRaiseRequest{Subject: Subject{Type: "activation", ID: "ACT-stop-test"}, Code: "activation-incomplete", BlockedActions: []string{"gate.run"}, Details: []byte(`{"activation_id":"ACT-stop-test","point_of_no_return":false,"components":["backend"]}`), ResolutionCondition: Condition{Type: "activation-recovered", Subject: Subject{Type: "activation", ID: "ACT-stop-test"}}, ResolverCapability: "stop.resolve"}
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
	if _, err := ResolveStop(options, StopResolveRequest{StopID: first.Stop.ID, FactEventID: factID}); err != nil {
		t.Fatalf("ResolveStop() error = %v", err)
	}
	options.IdempotencyKey = testIdempotencyKey()
	if _, err := ResolveStop(options, StopResolveRequest{StopID: second.Stop.ID, FactEventID: factID}); err != nil {
		t.Fatalf("ResolveStop(second) error = %v", err)
	}
	stops, attention, err := resumeStopsAndAttention(conn, "ENV-none", nil)
	if err != nil {
		t.Fatalf("resumeStopsAndAttention() error = %v", err)
	}
	if len(stops) != 0 || len(attention) != 1 || attention[0].ID != first.Attention.ID || attention[0].Status != "open" {
		t.Fatalf("resume Stops/Attention = %#v / %#v, want no unresolved Stops and independently open Attention", stops, attention)
	}
}
