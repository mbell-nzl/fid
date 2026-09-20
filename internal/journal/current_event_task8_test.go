package journal

import "testing"

func TestConcernDisposedRequiresExactFieldsAndOpenGuard(t *testing.T) {
	effect := Effect{Op: "update", Table: "judgment_concerns", PK: map[string]any{"id": "CON-1"}, Preconditions: map[string]any{"status": "open"}, Values: map[string]any{
		"status": "disposed", "disposition_action": "rejected", "disposition_issue_id": nil,
		"supersedes_concern_id": nil, "supersedes_finding_id": nil, "disposition_reason": "not applicable",
		"disposed_actor_id": "actor-builder", "disposed_at": "2026-09-21T00:00:00Z", "disposed_event_id": "EVT-1",
	}}
	event := Event{EventType: "concern.disposed", EventID: "EVT-1", ActorV2: EventActor{ActorID: "actor-builder"}}
	if err := validateCurrentUpdate(event, effect); err != nil {
		t.Fatalf("validateCurrentUpdate() error = %v", err)
	}
	missingGuard := effect
	missingGuard.Preconditions = map[string]any{}
	if err := validateCurrentUpdate(event, missingGuard); err == nil {
		t.Fatal("validateCurrentUpdate() accepted disposal without open-state guard")
	}
	wrongTransition := effect
	wrongTransition.Preconditions = map[string]any{"status": "disposed"}
	if err := validateCurrentUpdate(event, wrongTransition); err == nil {
		t.Fatal("validateCurrentUpdate() accepted disposed-to-disposed transition")
	}
}

func TestTask8ValueShapesReuseJudgmentSubject(t *testing.T) {
	for _, kind := range []string{"Concern", "ConcernDisposition", "Deviation"} {
		if currentValueShapes[kind] == "" {
			t.Fatalf("currentValueShapes[%q] is missing", kind)
		}
	}
	if got := currentValueShapes["JudgmentSubject"]; got != "type:=issue,envelope,run,review,finding,requirement id:ID" {
		t.Fatalf("JudgmentSubject shape = %q", got)
	}
	if got := currentValueShapes["Subject"]; got != "type:=run,envelope,activation id:ID" {
		t.Fatalf("Subject shape changed to %q", got)
	}
}
