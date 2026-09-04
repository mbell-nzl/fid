package issues

// VerifyAcceptanceCriterionOptions preserves the unavailable public request
// surface without retaining a Phase 21 verification writer.
type VerifyAcceptanceCriterionOptions struct {
	WorkspaceRoot  string
	MigrationsDir  string
	IssueID        string
	CriterionID    string
	Actor          string
	IdempotencyKey string
	Status         string
	EvidenceType   *string
	EvidenceRef    *string
	Notes          *string
}

type VerifyAcceptanceCriterionResult struct {
	ID           string
	IssueID      string
	CriterionID  string
	Status       string
	EvidenceType *string
	EvidenceRef  *string
	Notes        *string
	VerifiedAt   string
}

func VerifyAcceptanceCriterion(VerifyAcceptanceCriterionOptions) (VerifyAcceptanceCriterionResult, error) {
	return VerifyAcceptanceCriterionResult{}, unavailableCapability("reviews.coverage.v1", 23)
}
