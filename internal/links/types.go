package links

type AddOptions struct {
	WorkspaceRoot string
	MigrationsDir string
	FromID        string
	Relation      string
	ToID          string
	ActorRole     string
}

type AddResult struct {
	Declared Link
	Inverse  Link
}

type RemoveOptions struct {
	WorkspaceRoot string
	MigrationsDir string
	FromID        string
	Relation      string
	ToID          string
	ActorRole     string
}

type RemoveResult struct {
	Declared Link
	Inverse  Link
}

type ListOptions struct {
	WorkspaceRoot string
	MigrationsDir string
	EntityID      string
}

type ListItem struct {
	Direction   string  `json:"direction"`
	Relation    string  `json:"relation"`
	EntityType  string  `json:"entity_type"`
	EntityID    string  `json:"entity_id"`
	EntityTitle *string `json:"entity_title"`
}

type Link struct {
	FromType string
	FromID   string
	Relation string
	ToType   string
	ToID     string
}
