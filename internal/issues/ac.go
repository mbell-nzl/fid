package issues

import (
	"database/sql"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/mbell-nzl/fid/internal/config"
	fiddb "github.com/mbell-nzl/fid/internal/db"
	"github.com/mbell-nzl/fid/internal/journal"
	"github.com/mbell-nzl/fid/internal/validation"
)

const AcceptanceCriteriaLimit = 1000

type AddAcceptanceCriterionOptions struct {
	WorkspaceRoot  string
	MigrationsDir  string
	IssueID        string
	Text           string
	Ordinal        int
	Actor          string
	IdempotencyKey string
	ActorRole      string
}

type AddAcceptanceCriterionResult struct {
	ID             string
	IssueID        string
	Text           string
	Ordinal        int
	IssueRevision  int
	EventID        string
	OccurredAt     string
	IdempotencyKey string
	Replayed       bool
}

func AddAcceptanceCriterion(options AddAcceptanceCriterionOptions) (AddAcceptanceCriterionResult, error) {
	options.IdempotencyKey = idempotencyKeyForOptions(options.IdempotencyKey, options.ActorRole)
	if strings.TrimSpace(options.WorkspaceRoot) == "" {
		return AddAcceptanceCriterionResult{}, fmt.Errorf("add acceptance criterion: workspace root is required")
	}
	if strings.TrimSpace(options.MigrationsDir) == "" {
		return AddAcceptanceCriterionResult{}, fmt.Errorf("add acceptance criterion: migrations dir is required")
	}
	issueID := strings.TrimSpace(options.IssueID)
	if issueID == "" {
		return AddAcceptanceCriterionResult{}, fmt.Errorf("add acceptance criterion: issue id is required")
	}
	if strings.TrimSpace(options.Text) == "" {
		return AddAcceptanceCriterionResult{}, fmt.Errorf("add acceptance criterion: text is required")
	}
	text, err := validation.NormalizeText("acceptance criterion text", options.Text, 4096)
	if err != nil {
		return AddAcceptanceCriterionResult{}, err
	}
	if options.Ordinal <= 0 {
		return AddAcceptanceCriterionResult{}, fmt.Errorf("add acceptance criterion: ordinal must be a positive integer")
	}
	if options.Ordinal > 2147483647 {
		return AddAcceptanceCriterionResult{}, fmt.Errorf("add acceptance criterion: ordinal must be between 1 and 2147483647")
	}
	actor, err := standaloneActorForOptions(options.Actor, options.ActorRole)
	if err != nil {
		return AddAcceptanceCriterionResult{}, err
	}

	if _, err := journal.EnsureProjectionFresh(options.WorkspaceRoot, options.MigrationsDir); err != nil {
		return AddAcceptanceCriterionResult{}, err
	}

	projectConfig, err := config.ReadProjectConfig(filepath.Join(options.WorkspaceRoot, config.WorkspaceDirName, config.ConfigFileName))
	if err != nil {
		return AddAcceptanceCriterionResult{}, err
	}
	requestDigest, err := issueRequestDigest("issue.ac.add", projectConfig.ProjectID, actor, map[string]any{"issue_id": issueID, "ordinal": options.Ordinal, "text": text})
	if err != nil {
		return AddAcceptanceCriterionResult{}, err
	}
	commandArgv := issueACAddArgv(issueID, text, options.Ordinal, options.Actor, options.IdempotencyKey)
	var result AddAcceptanceCriterionResult
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
		idempotencyRequest := journal.IdempotencyRequest{OperationScope: "issue.ac.add:" + projectConfig.ProjectID, IdempotencyKey: options.IdempotencyKey, RequestDigest: requestDigest}
		if outcome, found, err := journal.ReplayIdempotent(conn, idempotencyRequest); err != nil {
			return err
		} else if found {
			result, err = acAddResultFromJSON(outcome.ResultJSON, outcome.EventID, true)
			return err
		}

		issueRevision, err := acceptanceCriteriaIssueRevision(conn, issueID)
		if err != nil {
			return err
		}
		var existing string
		err = conn.QueryRow(`SELECT id FROM acceptance_criteria WHERE issue_id = ? AND ordinal = ?`, issueID, options.Ordinal).Scan(&existing)
		if err == nil {
			return &DuplicateAcceptanceCriterionError{IssueID: issueID, Ordinal: options.Ordinal}
		}
		if err != sql.ErrNoRows {
			return fmt.Errorf("query acceptance criterion ordinal: %w", err)
		}
		var criterionCount int
		if err := conn.QueryRow(`SELECT COUNT(*) FROM acceptance_criteria WHERE issue_id = ?`, issueID).Scan(&criterionCount); err != nil {
			return fmt.Errorf("count acceptance criteria for issue %q: %w", issueID, err)
		}
		if criterionCount >= AcceptanceCriteriaLimit {
			return &ResultLimitExceededError{Collection: "acceptance_criteria", Limit: AcceptanceCriteriaLimit, Observed: criterionCount + 1}
		}

		criterionID, err := journal.GenerateID("AC", func(candidate string) (bool, error) {
			return acceptanceCriterionIDExists(conn, candidate)
		})
		if err != nil {
			return err
		}
		activityID, err := nextIssueActivityID(conn)
		if err != nil {
			return err
		}

		aggregateSequence, err := journal.NextAggregateSequence(filepath.Join(options.WorkspaceRoot, config.WorkspaceDirName, config.JournalDirName), journal.Aggregate{ID: issueID, Type: "issue"})
		if err != nil {
			return err
		}
		event := journal.Event{EventSchema: journal.EventSchemaV2, EventType: "issue.ac-added", OccurredAt: issueMutationNow(nil), ActorV2: actor.eventActor(), Aggregate: journal.Aggregate{ID: issueID, Sequence: aggregateSequence, Type: "issue"}, CommandName: "issue.ac.add"}
		timestamp := event.OccurredAt.UTC().Format(time.RFC3339)
		effects := []journal.Effect{
			{
				Op:    "insert",
				Table: "acceptance_criteria",
				PK:    map[string]any{"id": criterionID},
				Values: map[string]any{
					"issue_id": issueID,
					"ordinal":  options.Ordinal,
					"text":     text,
				},
			},
		}
		revisionEffects, err := BuildIssueRevisionEffects([]IssueRevision{{IssueID: issueID, Observed: issueRevision}})
		if err != nil {
			return err
		}
		effects = append(effects, revisionEffects...)
		event.Effects = effects
		prepared, err := journal.PrepareEvent(lock, event)
		if err != nil {
			return err
		}
		event.EventID, event.OccurredAt = prepared.EventID, prepared.OccurredAt
		event.Effects = append(effects, issueActivityEffect(activityID, event.EventID, actor.id, timestamp, issueID, "ac-added", map[string]any{"ac_id": criterionID, "ordinal": options.Ordinal}))
		seed := AddAcceptanceCriterionResult{ID: criterionID, IssueID: issueID, Text: text, Ordinal: options.Ordinal, IssueRevision: issueRevision + 1, OccurredAt: event.OccurredAt.UTC().Format(time.RFC3339Nano), IdempotencyKey: options.IdempotencyKey}
		resultJSON, err := acAddResultJSON(seed)
		if err != nil {
			return err
		}
		outcome, err := journal.CommitIdempotent(lock, conn, idempotencyRequest, event, func(journal.Event) (journal.IdempotencyMutation, error) {
			return journal.IdempotencyMutation{EntityID: criterionID, ResultJSON: resultJSON}, nil
		})
		if err != nil {
			return err
		}
		result, err = acAddResultFromJSON(outcome.ResultJSON, outcome.EventID, outcome.Replayed)
		return err
	})
	if err != nil {
		return AddAcceptanceCriterionResult{}, err
	}

	return result, nil
}

func acceptanceCriteriaIssueRevision(conn *sql.DB, issueID string) (int, error) {
	var kind string
	var revision int
	err := conn.QueryRow(`SELECT kind, revision FROM issues WHERE id = ?`, issueID).Scan(&kind, &revision)
	if err == sql.ErrNoRows {
		return 0, &NotFoundError{IssueID: issueID}
	}
	if err != nil {
		return 0, fmt.Errorf("query issue %q: %w", issueID, err)
	}
	if kind == "task" || kind == "subtask" {
		return revision, nil
	}
	return 0, fmt.Errorf(
		"add acceptance criterion: issue %q is kind %q; acceptance criteria require task or subtask",
		issueID,
		kind,
	)
}

func acceptanceCriterionIDExists(conn *sql.DB, criterionID string) (bool, error) {
	var existing string
	err := conn.QueryRow(`SELECT id FROM acceptance_criteria WHERE id = ? LIMIT 1`, criterionID).Scan(&existing)
	if err == nil {
		return true, nil
	}
	if err == sql.ErrNoRows {
		return false, nil
	}
	return false, fmt.Errorf("query acceptance criterion id %q: %w", criterionID, err)
}

func issueACAddArgv(issueID, text string, ordinal int, actor, key string) []string {
	return []string{
		"fid", "issue", "ac", "add", issueID,
		"--text", text,
		"--ordinal", fmt.Sprintf("%d", ordinal),
		"--actor", actor,
		"--idempotency-key", key,
	}
}
