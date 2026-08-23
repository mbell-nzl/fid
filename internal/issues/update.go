package issues

import (
	"database/sql"
	"fmt"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/mbell-nzl/fid/internal/config"
	fiddb "github.com/mbell-nzl/fid/internal/db"
	"github.com/mbell-nzl/fid/internal/journal"
)

type UpdateOptions struct {
	WorkspaceRoot  string
	MigrationsDir  string
	ID             string
	Actor          string
	IdempotencyKey string
	ActorRole      string
	Title          *string
	Priority       *int
	Labels         *[]string
	SpecMDPath     *string
	// OwnerRole and OwnerActor are retained only for source compatibility with
	// V1 in-process callers. Canonical Issue updates deliberately ignore them.
	OwnerRole        *string
	OwnerActor       *string
	ExpectedRevision *int
	now              func() time.Time
}

type UpdateResult struct {
	ID             string
	Revision       int
	Changes        []FieldChange
	EventID        string
	OccurredAt     string
	IdempotencyKey string
	Replayed       bool
}

type FieldChange struct {
	Field string
	From  any
	To    any
}

type issueUpdateState struct {
	ID         string
	Title      string
	Priority   *int
	Labels     []string
	SpecMDPath *string
	UpdatedAt  string
	Revision   int
}

func Update(options UpdateOptions) (UpdateResult, error) {
	if strings.TrimSpace(options.WorkspaceRoot) == "" {
		return UpdateResult{}, fmt.Errorf("update issue: workspace root is required")
	}
	if strings.TrimSpace(options.MigrationsDir) == "" {
		return UpdateResult{}, fmt.Errorf("update issue: migrations dir is required")
	}
	if strings.TrimSpace(options.ID) == "" {
		return UpdateResult{}, fmt.Errorf("update issue: issue id is required")
	}
	actor, err := standaloneActorForOptions(options.Actor, options.ActorRole)
	if err != nil {
		return UpdateResult{}, err
	}
	request, err := normalizeUpdateRequest(options)
	if err != nil {
		return UpdateResult{}, err
	}

	if _, err := journal.EnsureProjectionFresh(options.WorkspaceRoot, options.MigrationsDir); err != nil {
		return UpdateResult{}, err
	}
	projectConfig, err := config.ReadProjectConfig(filepath.Join(options.WorkspaceRoot, config.WorkspaceDirName, config.ConfigFileName))
	if err != nil {
		return UpdateResult{}, err
	}
	requestDigest, err := issueRequestDigest("issue.update", projectConfig.ProjectID, actor, updateRequestValue(request, options.ExpectedRevision))
	if err != nil {
		return UpdateResult{}, err
	}

	commandArgv := issueUpdateArgv(options.ID, options.Actor, options.IdempotencyKey, request, options.ExpectedRevision)
	var result UpdateResult
	err = journal.WithLock(options.WorkspaceRoot, journal.LockOwner{
		Command:   strings.Join(commandArgv, " "),
		ActorRole: actor.id,
	}, func(lock *journal.WriteLock) error {
		dbPath := filepath.Join(options.WorkspaceRoot, config.WorkspaceDirName, config.DatabaseFileName)
		conn, err := fiddb.Open(dbPath)
		if err != nil {
			return err
		}
		defer closeIssueDB(conn)
		idempotencyRequest := journal.IdempotencyRequest{OperationScope: "issue.update:" + projectConfig.ProjectID, IdempotencyKey: options.IdempotencyKey, RequestDigest: requestDigest}
		if outcome, found, err := journal.ReplayIdempotent(conn, idempotencyRequest); err != nil {
			return err
		} else if found {
			result, err = updateResultFromJSON(outcome.ResultJSON, outcome.EventID, true)
			return err
		}

		current, err := loadIssueUpdateState(conn, options.ID)
		if err != nil {
			return err
		}
		if options.ExpectedRevision != nil && *options.ExpectedRevision != current.Revision {
			return &StateConflictError{IssueID: current.ID, ExpectedRevision: *options.ExpectedRevision, ActualRevision: current.Revision}
		}
		if request.labels != nil {
			if err := requireExistingLabels(conn, *request.labels); err != nil {
				return fmt.Errorf("%s", strings.Replace(err.Error(), "create issue:", "update issue:", 1))
			}
		}

		scalarValues, scalarPreconditions, scalarPrior, changes := buildIssueUpdateChanges(current, request)
		labelEffects, labelChanged, labelChange := buildLabelEffects(current.ID, current.Labels, request.labels)
		if labelChanged {
			changes = append(changes, labelChange)
		}
		sort.Slice(changes, func(left, right int) bool {
			return changes[left].Field < changes[right].Field
		})
		if len(changes) == 0 {
			return fmt.Errorf("update issue: at least one field must change")
		}

		activityID, err := nextIssueActivityID(conn)
		if err != nil {
			return err
		}
		event := journal.Event{
			EventSchema: journal.EventSchemaV2,
			EventType:   "issue.updated",
			OccurredAt:  issueMutationNow(options.now),
			ActorV2:     actor.eventActor(),
			Aggregate:   journal.Aggregate{ID: current.ID, Sequence: current.Revision + 1, Type: "issue"},
			CommandName: "issue.update",
		}
		timestamp := event.OccurredAt.UTC().Format(time.RFC3339)
		scalarValues["updated_at"] = timestamp
		scalarPreconditions["updated_at"] = current.UpdatedAt
		scalarPrior["updated_at"] = current.UpdatedAt

		effects := []journal.Effect{
			{
				Op:            "update",
				Table:         "issues",
				PK:            map[string]any{"id": current.ID},
				Values:        scalarValues,
				Prior:         scalarPrior,
				Preconditions: scalarPreconditions,
			},
		}
		revisionEffects, err := BuildIssueRevisionEffects([]IssueRevision{{
			IssueID:  current.ID,
			Observed: current.Revision,
		}})
		if err != nil {
			return err
		}
		effects = append(effects, revisionEffects...)
		effects = append(effects, labelEffects...)
		event.Effects = effects
		prepared, err := journal.PrepareEvent(lock, event)
		if err != nil {
			return err
		}
		event.Effects = append(effects, issueActivityEffect(activityID, prepared.EventID, actor.id, timestamp, current.ID, "issue-updated", map[string]any{"changes": activityDetailChanges(changes)}))
		resultJSON, err := updateResultJSON(UpdateResult{ID: current.ID, Revision: current.Revision + 1, Changes: changes, IdempotencyKey: options.IdempotencyKey})
		if err != nil {
			return err
		}
		outcome, err := journal.CommitIdempotent(lock, conn, idempotencyRequest, event, func(journal.Event) (journal.IdempotencyMutation, error) {
			return journal.IdempotencyMutation{EntityID: current.ID, ResultJSON: resultJSON}, nil
		})
		if err != nil {
			return err
		}
		result, err = updateResultFromJSON(outcome.ResultJSON, outcome.EventID, outcome.Replayed)
		return err
	})
	if err != nil {
		return UpdateResult{}, err
	}

	return result, nil
}

func updateRequestValue(request normalizedUpdateRequest, expectedRevision *int) map[string]any {
	return map[string]any{
		"expected_revision": expectedRevision,
		"labels":            request.labels,
		"priority":          request.priority,
		"spec_md_path":      request.specMDPath,
		"title":             request.title,
	}
}

func loadIssueUpdateState(conn *sql.DB, issueID string) (issueUpdateState, error) {
	row := conn.QueryRow(`
		SELECT id, title, priority, spec_md_path, updated_at, revision
		FROM issues
		WHERE id = ?
	`, issueID)

	var (
		state      issueUpdateState
		priority   sql.NullInt64
		specMDPath sql.NullString
		updatedAt  sql.NullString
	)
	if err := row.Scan(
		&state.ID,
		&state.Title,
		&priority,
		&specMDPath,
		&updatedAt,
		&state.Revision,
	); err != nil {
		if err == sql.ErrNoRows {
			return issueUpdateState{}, &NotFoundError{IssueID: issueID}
		}
		return issueUpdateState{}, fmt.Errorf("query issue %q: %w", issueID, err)
	}
	state.Priority = nullIntPointer(priority)
	state.SpecMDPath = nullStringPointer(specMDPath)
	if !updatedAt.Valid || strings.TrimSpace(updatedAt.String) == "" {
		return issueUpdateState{}, fmt.Errorf("query issue %q: updated_at is required", issueID)
	}
	state.UpdatedAt = updatedAt.String

	labels, err := issueLabels(conn, issueID)
	if err != nil {
		return issueUpdateState{}, err
	}
	state.Labels = labels
	return state, nil
}
