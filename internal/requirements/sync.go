package requirements

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/mbell-nzl/fid/internal/config"
	fiddb "github.com/mbell-nzl/fid/internal/db"
	"github.com/mbell-nzl/fid/internal/journal"
	"github.com/mbell-nzl/fid/internal/policy"
)

const (
	actionRequirementSync = "requirement.sync"
	fidVersion            = "0.1.0"
)

// Sync parses requirement pointers from a markdown source document and upserts
// them into the requirements table through one journal event.
func Sync(options SyncOptions) (SyncResult, error) {
	if strings.TrimSpace(options.WorkspaceRoot) == "" {
		return SyncResult{}, fmt.Errorf("requirement sync: workspace root is required")
	}
	if strings.TrimSpace(options.MigrationsDir) == "" {
		return SyncResult{}, fmt.Errorf("requirement sync: migrations dir is required")
	}
	sourceDoc := strings.TrimSpace(options.SourceDoc)
	if sourceDoc == "" {
		return SyncResult{}, fmt.Errorf("requirement sync: --source-doc is required")
	}
	if err := policy.Authorize(options.WorkspaceRoot, options.ActorRole, actionRequirementSync); err != nil {
		return SyncResult{}, err
	}

	parsed, err := parseSourceDocument(options.WorkspaceRoot, sourceDoc, options.TypeFilter, options.RequirementID)
	if err != nil {
		return SyncResult{}, err
	}

	commandArgv := requirementSyncArgv(sourceDoc, options.TypeFilter, options.RequirementID, options.ActorRole)
	result := SyncResult{
		SourceDoc:      filepath.ToSlash(filepath.Clean(sourceDoc)),
		RequirementIDs: collectRequirementIDs(parsed),
	}

	err = journal.WithLock(options.WorkspaceRoot, journal.LockOwner{
		Command:   strings.Join(commandArgv, " "),
		ActorRole: options.ActorRole,
	}, func(lock *journal.WriteLock) error {
		if _, err := journal.EnsureProjectionFreshWithLock(lock, options.MigrationsDir); err != nil {
			return err
		}

		conn, err := fiddb.Open(filepath.Join(options.WorkspaceRoot, config.WorkspaceDirName, config.DatabaseFileName))
		if err != nil {
			return err
		}
		defer closeRequirementsConn(conn)

		effects, created, updated, unchanged, err := buildSyncEffects(conn, parsed)
		if err != nil {
			return err
		}
		result.Created = created
		result.Updated = updated
		result.Unchanged = unchanged

		if len(effects) == 0 {
			return nil
		}

		event, err := journal.PrepareEvent(lock, journal.Event{
			EventSchema: journal.CurrentEventSchema,
			EventType:   "requirement.synced",
			Actor: journal.Actor{
				Role: options.ActorRole,
			},
			Command: journal.Command{
				Argv:       commandArgv,
				CWD:        ".",
				FIDVersion: fidVersion,
			},
			Effects: effects,
		})
		if err != nil {
			return err
		}

		if err := journal.DryRunProjectEvent(conn, event); err != nil {
			return err
		}

		appended, err := journal.AppendEvent(lock, event)
		if err != nil {
			return err
		}
		return projectRequirementsEvent(options.WorkspaceRoot, conn, appended)
	})
	if err != nil {
		return SyncResult{}, err
	}

	return result, nil
}

func buildSyncEffects(conn *sql.DB, parsed []requirementRow) ([]journal.Effect, int, int, int, error) {
	effects := make([]journal.Effect, 0, len(parsed))
	created := 0
	updated := 0
	unchanged := 0

	for _, candidate := range parsed {
		current, found, err := loadRequirementRow(conn, candidate.ID)
		if err != nil {
			return nil, 0, 0, 0, err
		}
		if !found {
			effects = append(effects, journal.Effect{
				Op:    "insert",
				Table: "requirements",
				PK:    map[string]any{"id": candidate.ID},
				Values: map[string]any{
					"type":             candidate.Type,
					"summary":          candidate.Summary,
					"source_doc":       candidate.SourceDoc,
					"source_section":   dereferenceRequirementString(candidate.SourceSection),
					"last_synced_hash": candidate.LastSyncedHash,
					"status":           candidate.Status,
				},
			})
			created++
			continue
		}
		if requirementRowsEqual(current, candidate) {
			unchanged++
			continue
		}

		effects = append(effects, journal.Effect{
			Op:    "update",
			Table: "requirements",
			PK:    map[string]any{"id": candidate.ID},
			Values: map[string]any{
				"type":             candidate.Type,
				"summary":          candidate.Summary,
				"source_doc":       candidate.SourceDoc,
				"source_section":   dereferenceRequirementString(candidate.SourceSection),
				"last_synced_hash": candidate.LastSyncedHash,
				"status":           candidate.Status,
			},
			Prior: map[string]any{
				"type":             current.Type,
				"summary":          current.Summary,
				"source_doc":       current.SourceDoc,
				"source_section":   dereferenceRequirementString(current.SourceSection),
				"last_synced_hash": current.LastSyncedHash,
				"status":           current.Status,
			},
			Preconditions: map[string]any{
				"type":             current.Type,
				"summary":          current.Summary,
				"source_doc":       current.SourceDoc,
				"source_section":   dereferenceRequirementString(current.SourceSection),
				"last_synced_hash": current.LastSyncedHash,
				"status":           current.Status,
			},
		})
		updated++
	}

	return effects, created, updated, unchanged, nil
}

func requirementSyncArgv(sourceDoc, typeFilter, requirementID, actorRole string) []string {
	argv := []string{"fid", "requirement", "sync", "--source-doc", sourceDoc}
	if trimmed := strings.TrimSpace(typeFilter); trimmed != "" {
		argv = append(argv, "--type", trimmed)
	}
	if trimmed := strings.TrimSpace(requirementID); trimmed != "" {
		argv = append(argv, "--id", trimmed)
	}
	return append(argv, "--by", actorRole)
}

func collectRequirementIDs(rows []requirementRow) []string {
	ids := make([]string, 0, len(rows))
	for _, row := range rows {
		ids = append(ids, row.ID)
	}
	slices.Sort(ids)
	return ids
}

func loadRequirementRow(conn *sql.DB, requirementID string) (requirementRow, bool, error) {
	var (
		row           requirementRow
		sourceSection sql.NullString
		hash          sql.NullString
		status        sql.NullString
	)
	err := conn.QueryRow(`
		SELECT id, type, summary, source_doc, source_section, last_synced_hash, status
		FROM requirements
		WHERE id = ?
	`, requirementID).Scan(
		&row.ID,
		&row.Type,
		&row.Summary,
		&row.SourceDoc,
		&sourceSection,
		&hash,
		&status,
	)
	if err == sql.ErrNoRows {
		return requirementRow{}, false, nil
	}
	if err != nil {
		return requirementRow{}, false, fmt.Errorf("query requirement %q: %w", requirementID, err)
	}
	row.SourceSection = nullRequirementString(sourceSection)
	row.LastSyncedHash = strings.TrimSpace(hash.String)
	row.Status = strings.TrimSpace(status.String)
	return row, true, nil
}

func requirementRowsEqual(left, right requirementRow) bool {
	return left.ID == right.ID &&
		left.Type == right.Type &&
		left.Summary == right.Summary &&
		left.SourceDoc == right.SourceDoc &&
		dereferenceRequirementString(left.SourceSection) == dereferenceRequirementString(right.SourceSection) &&
		left.LastSyncedHash == right.LastSyncedHash &&
		left.Status == right.Status
}

func projectRequirementsEvent(workspaceRoot string, conn *sql.DB, event journal.Event) error {
	eventHash, err := journal.HashEvent(event)
	if err != nil {
		return fmt.Errorf("hash requirement event %q: %w", event.EventID, err)
	}

	journalFile := event.OccurredAt.UTC().Format("2006-01-02") + ".jsonl"
	journalPath := filepath.Join(workspaceRoot, config.WorkspaceDirName, config.JournalDirName, journalFile)
	info, err := os.Stat(journalPath)
	if err != nil {
		return fmt.Errorf("stat requirement journal %q: %w", journalPath, err)
	}

	return journal.ProjectEvent(conn, journal.ProjectionInput{
		Event:     event,
		EventHash: eventHash,
		File:      journalFile,
		Line:      countJournalLines(journalPath),
		Offset:    info.Size(),
	})
}

func countJournalLines(path string) int {
	payload, err := os.ReadFile(path)
	if err != nil {
		return 0
	}
	trimmed := strings.TrimSpace(string(payload))
	if trimmed == "" {
		return 0
	}
	return len(strings.Split(trimmed, "\n"))
}

func nullRequirementString(value sql.NullString) *string {
	if !value.Valid {
		return nil
	}
	copy := value.String
	return &copy
}

func dereferenceRequirementString(value *string) any {
	if value == nil {
		return nil
	}
	return *value
}

func closeRequirementsConn(conn *sql.DB) {
	_ = conn.Close()
}
