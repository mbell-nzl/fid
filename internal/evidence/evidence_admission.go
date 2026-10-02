package evidence

import "fmt"

type evidenceAdmission struct {
	ObligationID           string
	ObligationType         string
	AcceptedEvidenceTypes  []string
	AcceptedTrustSources   []string
	EvidenceType           string
	TrustSource            string
	HumanOutcomeID         string
	HumanOutcomeResult     string
	HumanOutcomeClass      string
	ObligationOutcomeClass string
}

type EvidenceTypeNotAcceptedError struct {
	ObligationID          string
	EvidenceType          string
	AcceptedEvidenceTypes []string
}

func (err *EvidenceTypeNotAcceptedError) Error() string {
	return fmt.Sprintf("Evidence type %q is not accepted by Obligation %q", err.EvidenceType, err.ObligationID)
}

type OutcomeNotPermittedError struct {
	ObligationID   string
	ObligationType string
	OutcomeID      string
	Result         string
}

func (err *OutcomeNotPermittedError) Error() string {
	return fmt.Sprintf("Human Outcome %q is not permitted for Obligation %q", err.OutcomeID, err.ObligationID)
}

type ReviewChangesRequiredError struct {
	ReviewID           string
	ObligationID       string
	BlockingFindingIDs []string
}

func (err *ReviewChangesRequiredError) Error() string {
	return fmt.Sprintf("Review %q cannot satisfy Obligation %q", err.ReviewID, err.ObligationID)
}

type ReviewNotIndependentError struct {
	ReviewActorID      string
	ConflictingType    string
	ConflictingID      string
	ConflictingActorID string
}

func (err *ReviewNotIndependentError) Error() string {
	return "review actor is not independent at Evidence link time"
}

type ReviewNotFoundError struct{ ID string }

func (err *ReviewNotFoundError) Error() string { return fmt.Sprintf("Review %q was not found", err.ID) }

type OutcomeNotFoundError struct{ ID string }

func (err *OutcomeNotFoundError) Error() string {
	return fmt.Sprintf("Human Outcome %q was not found", err.ID)
}

func classifyEvidenceAdmission(input evidenceAdmission) error {
	if input.EvidenceType == "human-outcome" &&
		(input.ObligationType != "human-outcome" || input.HumanOutcomeResult != "approved" || input.HumanOutcomeClass != input.ObligationOutcomeClass) {
		return &OutcomeNotPermittedError{
			ObligationID:   input.ObligationID,
			ObligationType: input.ObligationType,
			OutcomeID:      input.HumanOutcomeID,
			Result:         input.HumanOutcomeResult,
		}
	}
	if !contains(input.AcceptedEvidenceTypes, input.EvidenceType) ||
		!contains(input.AcceptedTrustSources, input.TrustSource) {
		return &EvidenceTypeNotAcceptedError{
			ObligationID:          input.ObligationID,
			EvidenceType:          input.EvidenceType,
			AcceptedEvidenceTypes: append([]string{}, input.AcceptedEvidenceTypes...),
		}
	}
	return nil
}

func validateReviewEvidenceAdmission(obligationID, reviewID, verdict string, blockingFindingIDs []string) error {
	if verdict == "pass" {
		return nil
	}
	return &ReviewChangesRequiredError{
		ReviewID:           reviewID,
		ObligationID:       obligationID,
		BlockingFindingIDs: append([]string{}, blockingFindingIDs...),
	}
}

func contains(values []string, target string) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}
	return false
}
