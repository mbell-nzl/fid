package requirements

// SyncOptions defines the inputs for syncing requirement pointers from a source
// markdown document into the workspace requirements table.
type SyncOptions struct {
	WorkspaceRoot string
	MigrationsDir string
	SourceDoc     string
	ActorRole     string
	TypeFilter    string
	RequirementID string
}

// SyncResult reports the persisted changes from a requirement sync run.
type SyncResult struct {
	SourceDoc      string
	Created        int
	Updated        int
	Unchanged      int
	RequirementIDs []string
}

type requirementRow struct {
	ID             string
	Type           string
	Summary        string
	SourceDoc      string
	SourceSection  *string
	LastSyncedHash string
	Status         string
}
