package workflow

import (
	"reflect"
	"testing"
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
