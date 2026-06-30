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

type BlockOptions struct {
	WorkspaceRoot string
	MigrationsDir string
	IssueID       string
	ActorRole     string
	ReasonMDPath  string
}

type BlockResult struct {
	IssueID    string
	FromState  string
	ToState    string
	Completion issues.CompleteResult
}

func Block(options BlockOptions) (BlockResult, error) {
	if strings.TrimSpace(options.WorkspaceRoot) == "" {
		return BlockResult{}, fmt.Errorf("dev block: workspace root is required")
	}
	if strings.TrimSpace(options.MigrationsDir) == "" {
		return BlockResult{}, fmt.Errorf("dev block: migrations dir is required")
	}
	issueID := strings.TrimSpace(options.IssueID)
	if issueID == "" {
		return BlockResult{}, fmt.Errorf("dev block: issue id is required")
	}
	reasonMDPath := strings.TrimSpace(options.ReasonMDPath)
	if reasonMDPath == "" {
		return BlockResult{}, fmt.Errorf("dev block: reason path is required")
	}
	if err := policy.Authorize(options.WorkspaceRoot, options.ActorRole, actionDevBlock); err != nil {
		return BlockResult{}, err
	}

	commandArgv := blockArgv(issueID, options.ActorRole, reasonMDPath)
	var result BlockResult
	err := withDevConn(options.WorkspaceRoot, options.MigrationsDir, commandArgv, options.ActorRole, func(lock *journal.WriteLock, conn *sql.DB) error {
		current, err := loadIssueState(conn, issueID)
		if err != nil {
			return err
		}
		if err := requireTaskOrSubtask(current.ID, current.Kind); err != nil {
			return err
		}
		if _, allowed := workflow.RulesForEdge(current.WorkflowID, current.Status, "failed_gate"); !allowed {
			return &issues.InvalidTransitionError{
				IssueID:        current.ID,
				CurrentState:   current.Status,
				RequestedState: "failed_gate",
			}
		}

		completionEffects, completionResult, err := issues.BuildCompletionEffects(conn, issues.CompletionBuildOptions{
			IssueID:     issueID,
			ActorRole:   options.ActorRole,
			Status:      "failed_gate",
			NotesMDPath: &reasonMDPath,
		})
		if err != nil {
			return err
		}
		releaseEffect, err := claims.BuildIssueReleaseEffect(conn, issueID, options.ActorRole)
		if err != nil {
			return err
		}
		activityID, err := journal.GenerateID("ACT", func(candidate string) (bool, error) {
			return activityIDExists(conn, candidate)
		})
		if err != nil {
			return err
		}

		event, err := journal.PrepareEvent(lock, journal.Event{
			EventSchema: journal.CurrentEventSchema,
			EventType:   "dev.blocked",
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
		event.Effects = append(event.Effects, completionEffects...)
		event.Effects = append(event.Effects, buildIssueTransitionEffects(current, "failed_gate", options.ActorRole, event.EventID, activityID, timestamp)...)
		event.Effects = append(event.Effects, releaseEffect)

		if err := appendAndProjectEvent(options.WorkspaceRoot, conn, lock, event); err != nil {
			return err
		}

		result = BlockResult{
			IssueID:    issueID,
			FromState:  current.Status,
			ToState:    "failed_gate",
			Completion: completionResult,
		}
		return nil
	})
	if err != nil {
		return BlockResult{}, err
	}

	return result, nil
}

func blockArgv(issueID, actorRole, reasonMDPath string) []string {
	return []string{"fid", "dev", "block", issueID, "--reason-md", reasonMDPath, "--by", actorRole}
}
