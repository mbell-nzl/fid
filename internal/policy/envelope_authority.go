package policy

import (
	"fmt"
	"strings"
)

type EnvelopeGrant struct {
	ActorID              string
	RequiredCapabilities []string
	Actions              []string
	IssueID              string
}

// EnvelopeGrantEvaluation contains the Activation-resolved values required for one authority decision.
type EnvelopeGrantEvaluation struct {
	Active            bool
	ActorID           string
	ActorCapabilities []string
	Action            string
	Capability        string
	IssueID           string
	EnvelopeID        string
	EnvelopeStatus    string
	EnvelopeRevision  int
	Revision          int
	Grants            []EnvelopeGrant
}

// EnvelopeAuthorityDeniedError reports a fail-closed authority decision.
type EnvelopeAuthorityDeniedError struct {
	ActorID    string
	Action     string
	Capability string
}

func (err *EnvelopeAuthorityDeniedError) Error() string {
	return fmt.Sprintf("envelope authority denied for actor %q action %q capability %q", err.ActorID, err.Action, err.Capability)
}

// EvaluateEnvelopeGrant permits only an active Actor and an exact active Envelope Grant match.
func EvaluateEnvelopeGrant(evaluation EnvelopeGrantEvaluation) error {
	if !validEnvelopeGrantEvaluation(evaluation) ||
		evaluation.EnvelopeStatus != "active" ||
		evaluation.EnvelopeRevision != evaluation.Revision ||
		!contains(evaluation.ActorCapabilities, evaluation.Capability) {
		return envelopeAuthorityDenied(evaluation)
	}
	for _, grant := range evaluation.Grants {
		if grant.ActorID == evaluation.ActorID &&
			grant.IssueID == evaluation.IssueID &&
			contains(grant.Actions, evaluation.Action) &&
			contains(grant.RequiredCapabilities, evaluation.Capability) {
			return nil
		}
	}
	return envelopeAuthorityDenied(evaluation)
}

func validEnvelopeGrantEvaluation(evaluation EnvelopeGrantEvaluation) bool {
	return evaluation.Active &&
		strings.TrimSpace(evaluation.ActorID) != "" &&
		strings.TrimSpace(evaluation.Action) != "" &&
		strings.TrimSpace(evaluation.Capability) != "" &&
		strings.TrimSpace(evaluation.IssueID) != "" &&
		strings.TrimSpace(evaluation.EnvelopeID) != "" &&
		evaluation.Revision > 0
}

func envelopeAuthorityDenied(evaluation EnvelopeGrantEvaluation) error {
	return &EnvelopeAuthorityDeniedError{ActorID: evaluation.ActorID, Action: evaluation.Action, Capability: evaluation.Capability}
}

func contains(values []string, wanted string) bool {
	for _, value := range values {
		if value == wanted {
			return true
		}
	}
	return false
}
