package issues

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/mbell-nzl/fid/internal/claims"
	"github.com/mbell-nzl/fid/internal/config"
	fiddb "github.com/mbell-nzl/fid/internal/db"
	"github.com/mbell-nzl/fid/internal/journal"
	"github.com/mbell-nzl/fid/internal/policy"
	"github.com/mbell-nzl/fid/internal/workflow"
)

type TransitionOptions struct {
	WorkspaceRoot string
	MigrationsDir string
	ID            string
	ToState       string
	ActorRole     string
}

type TransitionResult struct {
	ID        string
	FromState string
	ToState   string
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

	toState := strings.TrimSpace(options.ToState)
	if toState == "" {
		return TransitionResult{}, fmt.Errorf("transition issue: target state is required")
	}
	if err := policy.Authorize(options.WorkspaceRoot, options.ActorRole, issueTransitionAction(toState)); err != nil {
		return TransitionResult{}, err
	}

	if _, err := journal.EnsureProjectionFresh(options.WorkspaceRoot, options.MigrationsDir); err != nil {
		return TransitionResult{}, err
	}

	commandArgv := issueTransitionArgv(options.ID, toState, options.ActorRole)
	var result TransitionResult
	err := journal.WithLock(options.WorkspaceRoot, journal.LockOwner{
		Command:   strings.Join(commandArgv, " "),
		ActorRole: options.ActorRole,
	}, func(lock *journal.WriteLock) error {
		dbPath := filepath.Join(options.WorkspaceRoot, config.WorkspaceDirName, config.DatabaseFileName)
		conn, err := fiddb.Open(dbPath)
		if err != nil {
			return err
		}
		defer closeIssueDB(conn)

		current, err := loadIssueTransitionState(conn, options.ID)
		if err != nil {
			return err
		}
		registry, err := workflow.Load(conn)
		if err != nil {
			return err
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
			if err := claims.RequireLiveIssueClaimHeld(conn, current.ID, options.ActorRole); err != nil {
				return err
			}
		}

		activityID, err := journal.GenerateID("ACT", func(candidate string) (bool, error) {
			return activityIDExists(conn, candidate)
		})
		if err != nil {
			return err
		}

		event, err := journal.PrepareEvent(lock, journal.Event{
			EventSchema: journal.CurrentEventSchema,
			EventType:   "issue.transitioned",
			OccurredAt:  nextIssueUpdateTime(current.UpdatedAt, time.Now().UTC()),
			Actor: journal.Actor{
				Role: options.ActorRole,
			},
			Command: journal.Command{
				Argv:       commandArgv,
				CWD:        ".",
				FIDVersion: fidVersion,
			},
		})
		if err != nil {
			return err
		}

		timestamp := event.OccurredAt.UTC().Format(time.RFC3339)
		event.Effects = []journal.Effect{
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
		if leavesStartedCategory(registry, current.WorkflowID, current.Status, toState) {
			releaseEffect, err := claims.MaybeBuildIssueReleaseEffect(conn, current.ID, options.ActorRole)
			if err != nil {
				return err
			}
			if releaseEffect != nil {
				event.Effects = append(event.Effects, *releaseEffect)
			}
		}
		event.Effects = append(event.Effects,
			journal.Effect{
				Op:    "activity",
				Table: "activity_log",
				PK:    map[string]any{"id": activityID},
				Values: map[string]any{
					"journal_event_id": event.EventID,
					"ref_type":         "issue",
					"ref_id":           current.ID,
					"event":            "status_changed",
					"details": map[string]any{
						"from": current.Status,
						"to":   toState,
					},
					"actor":       options.ActorRole,
					"occurred_at": timestamp,
				},
			},
		)

		if err := journal.DryRunProjectEvent(conn, event); err != nil {
			return err
		}

		appended, err := journal.AppendEvent(lock, event)
		if err != nil {
			return err
		}

		eventHash, err := journal.HashEvent(appended)
		if err != nil {
			return fmt.Errorf("hash issue event %q: %w", appended.EventID, err)
		}

		journalFile := appended.OccurredAt.UTC().Format("2006-01-02") + ".jsonl"
		journalPath := filepath.Join(options.WorkspaceRoot, config.WorkspaceDirName, config.JournalDirName, journalFile)
		info, err := os.Stat(journalPath)
		if err != nil {
			return fmt.Errorf("stat issue journal %q: %w", journalPath, err)
		}

		if err := journal.ProjectEvent(conn, journal.ProjectionInput{
			Event:     appended,
			EventHash: eventHash,
			File:      journalFile,
			Line:      countJournalLines(journalPath),
			Offset:    info.Size(),
		}); err != nil {
			return err
		}

		result = TransitionResult{
			ID:        current.ID,
			FromState: current.Status,
			ToState:   toState,
		}
		return nil
	})
	if err != nil {
		return TransitionResult{}, err
	}

	return result, nil
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
