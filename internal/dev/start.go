package dev

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/mbell-nzl/fid/internal/claims"
	"github.com/mbell-nzl/fid/internal/issues"
	"github.com/mbell-nzl/fid/internal/journal"
	"github.com/mbell-nzl/fid/internal/policy"
	"github.com/mbell-nzl/fid/internal/workflow"
)

type StartOptions struct {
	WorkspaceRoot string
	MigrationsDir string
	IssueID       string
	ActorRole     string
	Model         string
}

type StartResult struct {
	IssueID   string
	ClaimedBy string
	ExpiresAt string
	FromState string
	ToState   string
	Model     string
	RunID     string
}

func Start(options StartOptions) (StartResult, error) {
	if strings.TrimSpace(options.WorkspaceRoot) == "" {
		return StartResult{}, fmt.Errorf("dev start: workspace root is required")
	}
	if strings.TrimSpace(options.MigrationsDir) == "" {
		return StartResult{}, fmt.Errorf("dev start: migrations dir is required")
	}
	issueID := strings.TrimSpace(options.IssueID)
	if issueID == "" {
		return StartResult{}, fmt.Errorf("dev start: issue id is required")
	}
	model := strings.TrimSpace(options.Model)
	if model == "" {
		return StartResult{}, fmt.Errorf("dev start: model is required")
	}
	if err := policy.Authorize(options.WorkspaceRoot, options.ActorRole, actionDevStart); err != nil {
		return StartResult{}, err
	}

	commandArgv := startArgv(issueID, options.ActorRole, model)
	var result StartResult
	err := withDevConn(options.WorkspaceRoot, options.MigrationsDir, commandArgv, options.ActorRole, func(lock *journal.WriteLock, conn *sql.DB) error {
		current, err := loadIssueState(conn, issueID)
		if err != nil {
			return err
		}
		if err := requireTaskOrSubtask(current.ID, current.Kind); err != nil {
			return err
		}

		registry, err := workflow.Load(conn)
		if err != nil {
			return fmt.Errorf("dev start: load workflow registry: %w", err)
		}
		if _, allowed := registry.RulesForEdge(current.WorkflowID, current.Status, "in_progress"); !allowed {
			return &issues.InvalidTransitionError{
				IssueID:        current.ID,
				CurrentState:   current.Status,
				RequestedState: "in_progress",
			}
		}

		now := time.Now().UTC()
		currentClaim, err := loadCurrentClaim(conn, issueID)
		if err != nil {
			return err
		}
		if currentClaim != nil && currentClaim.isLiveAt(now) {
			return &claims.LiveClaimConflictError{
				RefID:     issueID,
				ClaimedBy: currentClaim.ClaimedBy,
				ExpiresAt: currentClaim.expiresAt().Format(time.RFC3339),
			}
		}

		runID, err := journal.GenerateID("RUN", func(candidate string) (bool, error) {
			return runIDExists(conn, candidate)
		})
		if err != nil {
			return err
		}
		activityID, err := journal.GenerateID("ACT", func(candidate string) (bool, error) {
			return activityIDExists(conn, candidate)
		})
		if err != nil {
			return err
		}

		occurredAt := nextIssueUpdateTime(current.UpdatedAt, now)
		event, err := journal.PrepareEvent(lock, journal.Event{
			EventSchema: journal.CurrentEventSchema,
			EventType:   "dev.started",
			OccurredAt:  occurredAt,
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
				Op:    "insert",
				Table: "claims",
				PK: map[string]any{
					"ref_type": "issue",
					"ref_id":   issueID,
				},
				Values: map[string]any{
					"claimed_by":  options.ActorRole,
					"claimed_at":  timestamp,
					"ttl_seconds": defaultTTLSeconds,
				},
			},
			{
				Op:    "update",
				Table: "issues",
				PK:    map[string]any{"id": current.ID},
				Values: map[string]any{
					"status":     "in_progress",
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
			{
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
						"to":   "in_progress",
					},
					"actor":       options.ActorRole,
					"occurred_at": timestamp,
				},
			},
			{
				Op:    "insert",
				Table: "agent_runs",
				PK:    map[string]any{"id": runID},
				Values: map[string]any{
					"issue_id":   issueID,
					"role":       options.ActorRole,
					"model":      model,
					"started_at": timestamp,
				},
			},
		}

		if err := appendAndProjectEvent(options.WorkspaceRoot, conn, lock, event); err != nil {
			return err
		}

		result = StartResult{
			IssueID:   issueID,
			ClaimedBy: options.ActorRole,
			ExpiresAt: event.OccurredAt.UTC().Add(time.Duration(defaultTTLSeconds) * time.Second).Format(time.RFC3339),
			FromState: current.Status,
			ToState:   "in_progress",
			Model:     model,
			RunID:     runID,
		}
		return nil
	})
	if err != nil {
		return StartResult{}, err
	}

	return result, nil
}

func startArgv(issueID, actorRole, model string) []string {
	return []string{"fid", "dev", "start", issueID, "--model", model, "--by", actorRole}
}
