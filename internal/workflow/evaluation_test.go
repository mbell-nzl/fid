package workflow

import (
	"reflect"
	"testing"

	"github.com/mbell-nzl/fid/internal/journal"
)

func TestEvaluateObligationsReturnsIdenticalUnsatisfiedReasonsForBothGuards(t *testing.T) {
	obligation := Obligation{ObligationInput: ObligationInput{ID: "OBL-1", Type: "gate-evidence", Applicability: "required"}}
	for _, action := range []string{"issue.transition-terminal-success", "envelope.complete"} {
		got := EvaluateObligations(ObligationEvaluationInput{
			Action:      action,
			Obligations: []Obligation{obligation},
			Evidence:    map[string]bool{"OBL-1": false},
		})
		want := []EvaluationReason{{Code: "obligation-evidence-absent", Detail: "qualifying Evidence is absent", SubjectID: "OBL-1", SubjectType: "obligation"}}
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("EvaluateObligations(%q) = %#v, want %#v", action, got, want)
		}
	}
}

func TestTerminalSuccessEvaluationIncludesEnvelopeStop(t *testing.T) {
	options, _ := recoveryFixture(t)
	conn, err := openWorkflowDB(options.WorkspaceRoot)
	if err != nil {
		t.Fatalf("open workflow database: %v", err)
	}
	defer closeWorkflowDB(conn)
	event := stopRaisedEvent("", "STOP-envelope", "ATT-envelope", true)
	event.Effects[0].Values["subject_type"] = "envelope"
	event.Effects[0].Values["subject_id"] = "ENV-recovery"
	event.Effects[0].Values["blocked_actions_json"] = `["envelope.complete","issue.transition-terminal-success"]`
	event.Effects[0].Values["resolution_condition_json"] = `{"type":"work-ended","subject":{"type":"envelope","id":"ENV-recovery"},"budget_name":null,"capability_id":null}`
	event.Effects[1].Values["subject_type"] = "envelope"
	event.Effects[1].Values["subject_id"] = "ENV-recovery"
	event = refreshStopFixtureResult(event)
	if err := journal.WithLock(options.WorkspaceRoot, journal.LockOwner{Command: "test envelope Stop"}, func(lock *journal.WriteLock) error {
		prepared, err := journal.PrepareEvent(lock, event)
		if err != nil {
			return err
		}
		_, err = journal.CommitEvent(lock, conn, eventWithID(event, prepared.EventID))
		return err
	}); err != nil {
		t.Fatalf("commit envelope Stop: %v", err)
	}
	version := 1
	transition, err := EvaluateExecutionTransition(options.WorkspaceRoot, options.MigrationsDir, ExecutionTransitionRequest{IssueID: "ISS-recovery", ActorID: options.ActorID, ExpectedIssueRevision: &version, ActualIssueRevision: version, WorkflowEdgeValid: true, TerminalSuccess: true})
	if err != nil || !evaluationHasStop(transition.Reasons, "STOP-envelope") {
		t.Fatalf("terminal Issue evaluation = %#v, %v; want STOP-envelope", transition, err)
	}
}

func eventWithID(event journal.Event, eventID string) journal.Event {
	event.EventID = eventID
	for index := range event.Effects {
		if event.Effects[index].Table == "idempotency_results" {
			event.Effects[index].Values["result_event_id"] = eventID
		}
		if event.Effects[index].Values["created_event_id"] == "" {
			event.Effects[index].Values["created_event_id"] = eventID
		}
	}
	return event
}

func evaluationHasStop(reasons []EvaluationReason, stopID string) bool {
	for _, reason := range reasons {
		if reason.SubjectType == "stop" && reason.SubjectID == stopID {
			return true
		}
	}
	return false
}
