package issues

// AddressConcernOptions preserves the unavailable public request surface
// without retaining a Phase 21 concern writer.
type AddressConcernOptions struct {
	WorkspaceRoot           string
	MigrationsDir           string
	ConcernID               string
	Actor                   string
	IdempotencyKey          string
	AddressedInCompletionID string
}

type AddressConcernResult struct {
	ID                     string
	IssueID                string
	AddressingCompletionID string
	Status                 string
	UpdatedAt              string
}

func AddressConcern(AddressConcernOptions) (AddressConcernResult, error) {
	return AddressConcernResult{}, unavailableCapability("human.outcomes.v1", 23)
}
