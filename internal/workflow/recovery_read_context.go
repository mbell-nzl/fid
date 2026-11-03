package workflow

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"slices"

	"github.com/mbell-nzl/fid/internal/config"
)

// RecoveryReadContext explains local delivery and observation provenance.
// It is a read-only supplement, never an admission decision or a wire shape.
type RecoveryReadContext struct {
	Import             *ImportedHistoryContext
	Liveness           *Liveness
	LivenessHistorical bool
}

type ImportedHistoryContext struct {
	ProjectID     string
	ArchiveDigest string
	Destination   string
	EventCount    int
}

// ReadRecoveryContext verifies delivery metadata before labeling preserved
// history, and distinguishes imported observations from subsequent local ones.
func ReadRecoveryContext(workspaceRoot, runID string) (RecoveryReadContext, error) {
	var context RecoveryReadContext
	var binding config.ImportBinding
	_, err := os.Lstat(filepath.Join(config.LedgerDir(workspaceRoot), config.ImportBindingFileName))
	if err == nil {
		binding, err = config.ReadImportBinding(workspaceRoot)
		if err != nil {
			return context, fmt.Errorf("read recovery provenance: %w", err)
		}
		context.Import = &ImportedHistoryContext{ProjectID: binding.ProjectID, ArchiveDigest: binding.ArchiveDigest,
			Destination: binding.TargetRoot, EventCount: len(binding.ImportedEventIDs)}
	} else if !os.IsNotExist(err) {
		return context, fmt.Errorf("read recovery provenance: %w", err)
	}
	if runID == "" {
		return context, nil
	}
	conn, err := openWorkflowDB(workspaceRoot)
	if err != nil {
		return context, err
	}
	defer closeWorkflowDB(conn)
	var id, eventID string
	err = conn.QueryRow(`SELECT id, created_event_id FROM run_observations WHERE run_id = ? AND kind = 'liveness' ORDER BY observed_at DESC, id DESC LIMIT 1`, runID).Scan(&id, &eventID)
	if err == sql.ErrNoRows {
		return context, nil
	}
	if err != nil {
		return context, fmt.Errorf("read recovery liveness provenance: %w", err)
	}
	observation, err := loadRunObservation(conn, id)
	if err != nil {
		return context, err
	}
	context.Liveness = &observation.Liveness
	context.LivenessHistorical = slices.Contains(binding.ImportedEventIDs, eventID)
	return context, nil
}
