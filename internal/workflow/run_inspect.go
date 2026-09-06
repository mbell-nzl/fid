package workflow

import (
	"strings"

	"github.com/mbell-nzl/fid/internal/project"
)

type RunInspection struct {
	Run            Run            `json:"run"`
	Reconciliation Reconciliation `json:"reconciliation"`
	Liveness       *Liveness      `json:"liveness"`
}

// InspectRun composes the current repository state without appending an Event.
func InspectRun(workspaceRoot, migrationsDir, runID, repositoryRoot string) (RunInspection, error) {
	if strings.TrimSpace(workspaceRoot) == "" || strings.TrimSpace(migrationsDir) == "" || strings.TrimSpace(runID) == "" || strings.TrimSpace(repositoryRoot) == "" {
		return RunInspection{}, invalidRun("run inspect input is invalid")
	}
	conn, err := openWorkflowDB(workspaceRoot)
	if err != nil {
		return RunInspection{}, err
	}
	defer closeWorkflowDB(conn)
	run, err := loadRun(conn, runID)
	if err != nil {
		return RunInspection{}, err
	}
	checkpoint, err := latestRunCheckpoint(conn, run.ID)
	if err != nil {
		return RunInspection{}, err
	}
	base, checkpointID := run.StartRepository, (*string)(nil)
	if checkpoint != nil {
		base, checkpointID = checkpoint.Reconciliation.Repository, &checkpoint.ID
	}
	paths, err := runClaimPaths(conn, run.ClaimID)
	if err != nil {
		return RunInspection{}, err
	}
	current, err := project.ObserveRepository(repositoryRoot)
	if err != nil {
		return RunInspection{}, err
	}
	reconciliation, err := ReconcileRepository(base, current, checkpointID, paths)
	if err != nil {
		return RunInspection{}, err
	}
	liveness, err := latestLiveness(conn, run.ID)
	if err != nil {
		return RunInspection{}, err
	}
	return RunInspection{Run: run, Reconciliation: reconciliation, Liveness: liveness}, nil
}
