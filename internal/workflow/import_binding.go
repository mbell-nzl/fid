package workflow

import (
	"database/sql"
	"path/filepath"

	"github.com/mbell-nzl/fid/internal/config"
	"github.com/mbell-nzl/fid/internal/project"
)

func reconcileRunRepository(workspaceRoot string, start Repository, current project.Observation, checkpointID *string, paths []string) (Reconciliation, error) {
	if importedDestinationBinding(workspaceRoot, start, current) {
		return ReconcileImportedDestination(start, current, checkpointID, paths)
	}
	return ReconcileRepository(start, current, checkpointID, paths)
}

func importedDestinationBinding(workspaceRoot string, start Repository, current project.Observation) bool {
	binding, err := config.ReadImportBinding(workspaceRoot)
	if err != nil || binding.ProjectID == "" || binding.ProjectID != start.ProjectID || current.WorktreeID == "" {
		return false
	}
	root, err := filepath.Abs(workspaceRoot)
	return err == nil && binding.TargetRoot == root
}

// SameWorktree accepts a local observation or a destination bound to imported
// history for this project. It deliberately ignores HEAD and dirty state.
func SameWorktree(workspaceRoot string, recorded, observed Repository) bool {
	if recorded.ProjectID == "" || recorded.ProjectID != observed.ProjectID || observed.WorktreeID == "" {
		return false
	}
	if recorded.WorktreeID == observed.WorktreeID {
		return true
	}
	binding, err := config.ReadImportBinding(workspaceRoot)
	if err != nil || binding.ProjectID != observed.ProjectID {
		return false
	}
	root, err := filepath.Abs(workspaceRoot)
	return err == nil && binding.TargetRoot == root
}

func importedObservation(workspaceRoot string, conn *sql.DB, observationID string) bool {
	binding, err := config.ReadImportBinding(workspaceRoot)
	if err != nil {
		return false
	}
	var eventID string
	if err := conn.QueryRow(`SELECT created_event_id FROM run_observations WHERE id = ?`, observationID).Scan(&eventID); err != nil {
		return false
	}
	for _, imported := range binding.ImportedEventIDs {
		if eventID == imported {
			return true
		}
	}
	return false
}
