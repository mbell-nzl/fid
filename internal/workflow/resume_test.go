package workflow

import "testing"

func TestResumeKeepsOpenAttentionAfterAllLinkedStopsAreResolved(t *testing.T) {
	t.Parallel()
	options, activation, factID := finalizedStopConditionActivation(t)
	workspace := options.WorkspaceRoot
	conn, err := openWorkflowDB(workspace)
	if err != nil {
		t.Fatalf("open workflow database: %v", err)
	}
	defer closeWorkflowDB(conn)
	request := StopRaiseRequest{Subject: Subject{Type: "activation", ID: activation.ID}, Code: "activation-incomplete", BlockedActions: []string{"gate.run"}, Details: mustRunJSON(map[string]any{"activation_id": activation.ID, "point_of_no_return": false, "components": []string{"backend"}}), ResolutionCondition: Condition{Type: "activation-recovered", Subject: Subject{Type: "activation", ID: activation.ID}}, ResolverCapability: "stop.resolve"}
	first, err := RaiseStop(options, request)
	if err != nil {
		t.Fatalf("RaiseStop(first) error = %v", err)
	}
	options.IdempotencyKey = testIdempotencyKey()
	second, err := RaiseStop(options, request)
	if err != nil {
		t.Fatalf("RaiseStop(second) error = %v", err)
	}
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
