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

type Link struct {
	FromType string
	FromID   string
	Relation string
	ToType   string
	ToID     string
}
