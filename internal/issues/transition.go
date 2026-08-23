package issues

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/mbell-nzl/fid/internal/claims"
	"github.com/mbell-nzl/fid/internal/config"
	fiddb "github.com/mbell-nzl/fid/internal/db"
	"github.com/mbell-nzl/fid/internal/journal"
	"github.com/mbell-nzl/fid/internal/workflow"
)

type TransitionOptions struct {
	WorkspaceRoot    string
	MigrationsDir    string
	ID               string
	ToState          string
	Actor            string
	IdempotencyKey   string
	ActorRole        string
	ExpectedRevision *int
	now              func() time.Time
}

type TransitionResult struct {
	ID             string
	Revision       int
	FromState      string
	ToState        string
	EventID        string
	OccurredAt     string
	IdempotencyKey string
	Replayed       bool
}

type InvalidTransitionError struct {
	IssueID        string
	CurrentState   string
	RequestedState string
}

func (err *InvalidTransitionError) Error() string {
	if err == nil {
		return "<nil>"
	}
	return fmt.Sprintf(
		"transition issue %q from %q to %q is not allowed",
		err.IssueID,
		err.CurrentState,
		err.RequestedState,
	)
}

type BlockedTransitionError struct {
	IssueID        string
	CurrentState   string
	RequestedState string
	Unmet          []workflow.UnmetCondition
	Cause          error
}

func (err *BlockedTransitionError) Unwrap() error {
	if err == nil {
		return nil
	}
	return err.Cause
}

func (err *BlockedTransitionError) Error() string {
	if err == nil {
		return "<nil>"
	}

	var builder strings.Builder
	fmt.Fprintf(
		&builder,
		"transition issue %q from %q to %q blocked by workflow rules:",
		err.IssueID,
		err.CurrentState,
		err.RequestedState,
	)
	for _, condition := range err.Unmet {
		builder.WriteString("\n  - ")
		builder.WriteString(condition.Rule)
		details := make([]string, 0, 2)
		if len(condition.EntityIDs) > 0 {
			details = append(details, "entities "+strings.Join(condition.EntityIDs, ", "))
		}
		if len(condition.MissingRoles) > 0 {
			details = append(details, "missing roles "+strings.Join(condition.MissingRoles, ", "))
		}
		if len(details) > 0 {
			builder.WriteString(": ")
			builder.WriteString(strings.Join(details, "; "))
		}
	}
	return builder.String()
}

type issueTransitionState struct {
	ID         string
	Kind       string
	WorkflowID string
	Status     string
	UpdatedAt  string
	Revision   int
}

func Transition(options TransitionOptions) (TransitionResult, error) {
	if strings.TrimSpace(options.WorkspaceRoot) == "" {
		return TransitionResult{}, fmt.Errorf("transition issue: workspace root is required")
	}
	if strings.TrimSpace(options.MigrationsDir) == "" {
		return TransitionResult{}, fmt.Errorf("transition issue: migrations dir is required")
	}
	if strings.TrimSpace(options.ID) == "" {
		return TransitionResult{}, fmt.Errorf("transition issue: issue id is required")
	}
	actor, err := standaloneActorForOptions(options.Actor, options.ActorRole)
	if err != nil {
		return TransitionResult{}, err
	}

	toState := strings.TrimSpace(options.ToState)
	if toState == "" {
		return TransitionResult{}, fmt.Errorf("transition issue: target state is required")
	}
	if _, err := journal.EnsureProjectionFresh(options.WorkspaceRoot, options.MigrationsDir); err != nil {
		return TransitionResult{}, err
	}
	projectConfig, err := config.ReadProjectConfig(filepath.Join(options.WorkspaceRoot, config.WorkspaceDirName, config.ConfigFileName))
	if err != nil {
		return TransitionResult{}, err
	}
	requestDigest, err := issueRequestDigest("issue.transition", projectConfig.ProjectID, actor, transitionRequestValue(toState, options.ExpectedRevision))
	if err != nil {
		return TransitionResult{}, err
	}

	commandArgv := issueTransitionArgv(options.ID, toState, options.Actor, options.IdempotencyKey, options.ExpectedRevision)
	var result TransitionResult
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
		idempotencyRequest := journal.IdempotencyRequest{OperationScope: "issue.transition:" + projectConfig.ProjectID, IdempotencyKey: options.IdempotencyKey, RequestDigest: requestDigest}
		if outcome, found, err := journal.ReplayIdempotent(conn, idempotencyRequest); err != nil {
			return err
		} else if found {
			result, err = transitionResultFromJSON(outcome.ResultJSON, outcome.EventID, true)
			return err
		}

		current, err := loadIssueTransitionState(conn, options.ID)
		if err != nil {
			return err
		}
		if options.ExpectedRevision != nil && *options.ExpectedRevision != current.Revision {
			return &StateConflictError{IssueID: current.ID, ExpectedRevision: *options.ExpectedRevision, ActualRevision: current.Revision}
		}
		registry, err := workflow.Load(conn)
		if err != nil {
			return err
		}
		if _, known := registry.CategoryForState(current.WorkflowID, toState); !known {
			return &InvalidTransitionError{
				IssueID:        current.ID,
				CurrentState:   current.Status,
				RequestedState: toState,
			}
		}
		rules, allowed := registry.RulesForEdge(current.WorkflowID, current.Status, toState)
		if !allowed {
			return &InvalidTransitionError{
				IssueID:        current.ID,
				CurrentState:   current.Status,
				RequestedState: toState,
			}
		}

		unmet, err := evaluateTransitionRules(conn, current, rules)
		if err != nil {
			return err
		}
		if len(unmet) > 0 {
			return &BlockedTransitionError{
				IssueID:        current.ID,
				CurrentState:   current.Status,
				RequestedState: toState,
				Unmet:          unmet,
			}
		}
		if entersStartedCategory(registry, current.WorkflowID, current.Status, toState) {
			if err := claims.RequireLiveIssueClaimHeld(conn, current.ID, actor.id); err != nil {
				return &BlockedTransitionError{
					IssueID:        current.ID,
					CurrentState:   current.Status,
					RequestedState: toState,
					Unmet:          []workflow.UnmetCondition{{Rule: "live_claim_held"}},
					Cause:          err,
				}
			}
		}

		activityID, err := nextIssueActivityID(conn)
		if err != nil {
			return err
		}

		event := journal.Event{
			EventSchema: journal.EventSchemaV2,
			EventType:   "issue.transitioned",
			OccurredAt:  issueMutationNow(options.now),
			ActorV2:     actor.eventActor(),
			Aggregate:   journal.Aggregate{ID: current.ID, Sequence: current.Revision + 1, Type: "issue"},
			CommandName: "issue.transition",
		}
		timestamp := event.OccurredAt.UTC().Format(time.RFC3339)
		effects := []journal.Effect{
			{
				Op:    "update",
				Table: "issues",
				PK:    map[string]any{"id": current.ID},
				Values: map[string]any{
					"status":     toState,
					"updated_at": timestamp,
				},
				Prior: map[string]any{
					"status": current.Status,
				},
				Preconditions: map[string]any{
					"status":      current.Status,
					"workflow_id": current.WorkflowID,
				},
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
		if leavesStartedCategory(registry, current.WorkflowID, current.Status, toState) {
			releaseEffect, err := claims.MaybeBuildIssueReleaseEffect(conn, current.ID, actor.id)
			if err != nil {
				return err
			}
			if releaseEffect != nil {
				effects = append(effects, *releaseEffect)
			}
		}
		event.Effects = effects
		prepared, err := journal.PrepareEvent(lock, event)
		if err != nil {
			return err
		}
		event.Effects = append(effects, issueActivityEffect(activityID, prepared.EventID, actor.id, timestamp, current.ID, "status-changed", map[string]any{"from": current.Status, "to": toState}))
		resultJSON, err := transitionResultJSON(TransitionResult{ID: current.ID, Revision: current.Revision + 1, FromState: current.Status, ToState: toState, IdempotencyKey: options.IdempotencyKey})
		if err != nil {
			return err
		}
		outcome, err := journal.CommitIdempotent(lock, conn, idempotencyRequest, event, func(journal.Event) (journal.IdempotencyMutation, error) {
			return journal.IdempotencyMutation{EntityID: current.ID, ResultJSON: resultJSON}, nil
		})
		if err != nil {
			return err
		}
		result, err = transitionResultFromJSON(outcome.ResultJSON, outcome.EventID, outcome.Replayed)
		return err
	})
	if err != nil {
		return TransitionResult{}, err
	}

	return result, nil
}

func transitionRequestValue(toState string, expectedRevision *int) map[string]any {
	return map[string]any{"expected_revision": expectedRevision, "to": toState}
}

func entersStartedCategory(registry *workflow.Registry, workflowID, fromState, toState string) bool {
	fromCategory, toCategory, ok := transitionCategories(registry, workflowID, fromState, toState)
	if !ok {
		return false
	}
	return fromCategory != "started" && toCategory == "started"
}

func leavesStartedCategory(registry *workflow.Registry, workflowID, fromState, toState string) bool {
	fromCategory, toCategory, ok := transitionCategories(registry, workflowID, fromState, toState)
	if !ok {
		return false
	}
	return fromCategory == "started" && toCategory != "started"
}

func transitionCategories(registry *workflow.Registry, workflowID, fromState, toState string) (string, string, bool) {
	fromCategory, ok := registry.CategoryForState(workflowID, fromState)
	if !ok {
		return "", "", false
	}
	toCategory, ok := registry.CategoryForState(workflowID, toState)
	if !ok {
		return "", "", false
	}
	return fromCategory, toCategory, true
}
