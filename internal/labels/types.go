package labels

import "fmt"

const ReadCollectionLimit = 10000

type AddOptions struct {
	WorkspaceRoot string
	MigrationsDir string
	Name          string
	Colour        *string
	Description   *string
	ActorRole     string
}

type AddResult struct {
	Label      Label
	EventID    string
	OccurredAt string
}

type ListOptions struct {
	WorkspaceRoot string
	MigrationsDir string
}

type Label struct {
	Name        string  `json:"name"`
	Colour      *string `json:"colour"`
	Description *string `json:"description"`
}

type ResultLimitExceededError struct {
	Collection string
	Limit      int
	Observed   int
}

func (err *ResultLimitExceededError) Error() string {
	return fmt.Sprintf("%s contains %d items, exceeding limit %d", err.Collection, err.Observed, err.Limit)
}
