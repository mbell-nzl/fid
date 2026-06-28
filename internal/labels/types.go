package labels

type AddOptions struct {
	WorkspaceRoot string
	MigrationsDir string
	Name          string
	Colour        *string
	Description   *string
	ActorRole     string
}

type AddResult struct {
	Label Label
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
