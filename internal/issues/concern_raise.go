package issues

// RaiseConcernOptions preserves the unavailable public request surface without
// retaining a Phase 21 concern writer.
type RaiseConcernOptions struct {
	WorkspaceRoot        string
	MigrationsDir        string
	IssueID              string
	Actor                string
	IdempotencyKey       string
	RaisedInCompletionID string
	Summary              string
	BodyMDPath           *string
}

type RaiseConcernResult struct {
	ID                   string
	IssueID              string
	RaisedInCompletionID string
	Summary              string
	BodyMDPath           *string
	Status               string
	CreatedAt            string
	UpdatedAt            string
}

func RaiseConcern(RaiseConcernOptions) (RaiseConcernResult, error) {
	return RaiseConcernResult{}, unavailableCapability("human.outcomes.v1", 23)
}
