package links

type AddOptions struct {
	WorkspaceRoot string
	MigrationsDir string
	FromType      string
	FromID        string
	Relation      string
	ToType        string
	ToID          string
	ActorRole     string
}

type AddResult struct {
	Declared       Link
	Inverse        Link
	IssueRevisions []IssueRevision
	EventID        string
	OccurredAt     string
}

type RemoveOptions struct {
	WorkspaceRoot string
	MigrationsDir string
	FromType      string
	FromID        string
	Relation      string
	ToType        string
	ToID          string
	ActorRole     string
}

type RemoveResult struct {
	Declared       Link
	Inverse        Link
	IssueRevisions []IssueRevision
	EventID        string
	OccurredAt     string
}

type ListOptions struct {
	WorkspaceRoot string
	MigrationsDir string
	EntityType    string
	EntityID      string
}

type IssueRevision struct {
	ID       string `json:"id"`
	Revision int    `json:"revision"`
}

type ListItem struct {
	Direction   string  `json:"direction"`
	Relation    string  `json:"relation"`
	EntityType  string  `json:"entity_type"`
	EntityID    string  `json:"entity_id"`
	EntityTitle *string `json:"entity_title"`
}

type Link struct {
	FromType string `json:"from_type"`
	FromID   string `json:"from_id"`
	Relation string `json:"relation"`
	ToType   string `json:"to_type"`
	ToID     string `json:"to_id"`
}
