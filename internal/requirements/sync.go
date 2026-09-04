package requirements

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/mbell-nzl/fid/internal/config"
	fiddb "github.com/mbell-nzl/fid/internal/db"
	"github.com/mbell-nzl/fid/internal/journal"
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
	actor, err := standaloneActorForOptions(options.Actor)
	if err != nil {
		return SyncResult{}, err
	}

	parsed, err := parseSourceDocument(options.WorkspaceRoot, sourceDoc, options.TypeFilter, options.RequirementID)
	if err != nil {
		return SyncResult{}, err
	}

	commandArgv := requirementSyncArgv(sourceDoc, options.TypeFilter, options.RequirementID, options.Actor, options.IdempotencyKey)
	result := SyncResult{
		SourceDoc:      filepath.ToSlash(filepath.Clean(sourceDoc)),
		RequirementIDs: collectRequirementIDs(parsed),
	}

	err = journal.WithLock(options.WorkspaceRoot, journal.LockOwner{
		Command: strings.Join(commandArgv, " "),
	}, func(lock *journal.WriteLock) error {
		if _, err := journal.EnsureProjectionFreshWithLock(lock, options.MigrationsDir); err != nil {
			return err
		}

		conn, err := fiddb.Open(filepath.Join(options.WorkspaceRoot, config.WorkspaceDirName, config.DatabaseFileName))
		if err != nil {
			return err
		}
		defer closeRequirementsConn(conn)
		projectConfig, err := config.ReadProjectConfig(filepath.Join(options.WorkspaceRoot, config.WorkspaceDirName, config.ConfigFileName))
		if err != nil {
			return err
		}
		digest, err := requirementRequestDigest(projectConfig.ProjectID, actor, map[string]any{"id": options.RequirementID, "source_doc": result.SourceDoc, "type": options.TypeFilter})
		if err != nil {
			return err
		}
		idempotencyRequest := journal.IdempotencyRequest{OperationScope: "requirement.sync:" + projectConfig.ProjectID, IdempotencyKey: options.IdempotencyKey, RequestDigest: digest}
		if outcome, found, err := journal.ReplayIdempotent(conn, idempotencyRequest); err != nil {
			return err
		} else if found {
			decoded, err := syncResultFromJSON(outcome.ResultJSON, outcome.EventID, true)
			if err != nil {
				return err
			}
			result = decoded
			return nil
		}

		effects, created, updated, unchanged, err := buildSyncEffects(conn, parsed)
		if err != nil {
			return err
		}
		result.Created = created
		result.Updated = updated
		result.Unchanged = unchanged

		if len(effects) == 0 {
			result.IdempotencyKey = options.IdempotencyKey
			return nil
		}
		sequence, err := journal.NextAggregateSequence(filepath.Join(options.WorkspaceRoot, config.WorkspaceDirName, config.JournalDirName), journal.Aggregate{ID: result.SourceDoc, Type: "requirement-sync"})
		if err != nil {
			return err
		}
		event := journal.Event{EventSchema: journal.EventSchemaV2, EventType: "requirement.synced", ActorV2: actor.eventActor(), Aggregate: journal.Aggregate{ID: result.SourceDoc, Sequence: sequence, Type: "requirement-sync"}, CommandName: "requirement.sync", Effects: effects}
		seed := result
		seed.IdempotencyKey = options.IdempotencyKey
		outcome, err := journal.CommitIdempotent(lock, conn, idempotencyRequest, event, func(prepared journal.Event) (journal.IdempotencyMutation, error) {
			seed.OccurredAt = prepared.OccurredAt.UTC().Format(time.RFC3339Nano)
			encoded, err := syncResultJSON(seed)
			if err != nil {
				return journal.IdempotencyMutation{}, err
			}
			return journal.IdempotencyMutation{EntityID: result.SourceDoc, ResultJSON: encoded}, nil
		})
		if err != nil {
			return err
		}
		decoded, err := syncResultFromJSON(outcome.ResultJSON, outcome.EventID, outcome.Replayed)
		if err != nil {
			return err
		}
		result = decoded
		return nil
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

func requirementSyncArgv(sourceDoc, typeFilter, requirementID, actor, key string) []string {
	argv := []string{"fid", "requirement", "sync", "--source-doc", sourceDoc}
	if trimmed := strings.TrimSpace(typeFilter); trimmed != "" {
		argv = append(argv, "--type", trimmed)
	}
	if trimmed := strings.TrimSpace(requirementID); trimmed != "" {
		argv = append(argv, "--id", trimmed)
	}
	return append(argv, "--actor", actor, "--idempotency-key", key)
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

func syncResultJSON(result SyncResult) (string, error) {
	payload, err := journal.MarshalCanonicalJSON(map[string]any{"created": result.Created, "event_occurred_at": result.OccurredAt, "idempotency_key": result.IdempotencyKey, "replayed": false, "requirement_ids": result.RequirementIDs, "source_doc": result.SourceDoc, "unchanged": result.Unchanged, "updated": result.Updated})
	return string(payload), err
}

func syncResultFromJSON(raw, eventID string, replayed bool) (SyncResult, error) {
	var value struct {
		Created        int      `json:"created"`
		OccurredAt     string   `json:"event_occurred_at"`
		IdempotencyKey string   `json:"idempotency_key"`
		RequirementIDs []string `json:"requirement_ids"`
		SourceDoc      string   `json:"source_doc"`
		Unchanged      int      `json:"unchanged"`
		Updated        int      `json:"updated"`
	}
	if err := json.Unmarshal([]byte(raw), &value); err != nil {
		return SyncResult{}, fmt.Errorf("decode stored requirement sync result: %w", err)
	}
	return SyncResult{SourceDoc: value.SourceDoc, Created: value.Created, Updated: value.Updated, Unchanged: value.Unchanged, RequirementIDs: value.RequirementIDs, EventID: eventID, OccurredAt: value.OccurredAt, IdempotencyKey: value.IdempotencyKey, Replayed: replayed}, nil
}
