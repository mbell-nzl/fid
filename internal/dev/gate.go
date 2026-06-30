package dev

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/mbell-nzl/fid/internal/claims"
	"github.com/mbell-nzl/fid/internal/journal"
	"github.com/mbell-nzl/fid/internal/policy"
)

type GateOptions struct {
	WorkspaceRoot string
	MigrationsDir string
	IssueID       string
	ActorRole     string
	Command       string
	ExitCode      int
	OutputExcerpt *string
}

type GateResult struct {
	GateID        string
	IssueID       string
	AttemptNumber int
	AgentRunID    *string
	Command       string
	ExitCode      int
	GatesPassed   int
	GatesFailed   int
}

func Gate(options GateOptions) (GateResult, error) {
	if strings.TrimSpace(options.WorkspaceRoot) == "" {
		return GateResult{}, fmt.Errorf("dev gate: workspace root is required")
	}
	if strings.TrimSpace(options.MigrationsDir) == "" {
		return GateResult{}, fmt.Errorf("dev gate: migrations dir is required")
	}
	issueID := strings.TrimSpace(options.IssueID)
	if issueID == "" {
		return GateResult{}, fmt.Errorf("dev gate: issue id is required")
	}
	command := strings.TrimSpace(options.Command)
	if command == "" {
		return GateResult{}, fmt.Errorf("dev gate: command is required")
	}
	if err := policy.Authorize(options.WorkspaceRoot, options.ActorRole, actionDevGate); err != nil {
		return GateResult{}, err
	}

	commandArgv := gateArgv(issueID, options.ActorRole, command, options.ExitCode, options.OutputExcerpt)
	var result GateResult
	err := withDevConn(options.WorkspaceRoot, options.MigrationsDir, commandArgv, options.ActorRole, func(lock *journal.WriteLock, conn *sql.DB) error {
		current, err := loadIssueState(conn, issueID)
		if err != nil {
			return err
		}
		if err := requireTaskOrSubtask(current.ID, current.Kind); err != nil {
			return err
		}

		category, err := issueStatusCategory(current.WorkflowID, current.Status)
		if err != nil {
			return err
		}
		if category != "started" {
			return fmt.Errorf(
				"dev gate: issue %q is in state %q; gate recording requires a started-category state",
				issueID,
				current.Status,
			)
		}
		if err := claims.RequireLiveIssueClaimHeld(conn, issueID, options.ActorRole); err != nil {
			return err
		}

		attemptNumber, err := activeAttemptNumber(conn, issueID)
		if err != nil {
			return err
		}
		gatesPassed, gatesFailed, err := gateCounts(conn, issueID, attemptNumber)
		if err != nil {
			return err
		}
		agentRunID, err := openRunID(conn, issueID)
		if err != nil {
			return err
		}
		gateID, err := journal.GenerateID("GATE", func(candidate string) (bool, error) {
			return gateIDExists(conn, candidate)
		})
		if err != nil {
			return err
		}

		event, err := journal.PrepareEvent(lock, journal.Event{
			EventSchema: journal.CurrentEventSchema,
			EventType:   "dev.gate.recorded",
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

		values := map[string]any{
			"issue_id":         issueID,
			"attempt_number":   attemptNumber,
			"recorded_by_role": options.ActorRole,
			"command":          command,
			"exit_code":        options.ExitCode,
			"ran_at":           event.OccurredAt.UTC().Format(time.RFC3339),
		}
		if agentRunID != nil {
			values["agent_run_id"] = *agentRunID
		}
		if outputExcerpt := trimmedOptionalString(options.OutputExcerpt); outputExcerpt != nil {
			values["output_excerpt"] = *outputExcerpt
		}
		event.Effects = []journal.Effect{{
			Op:     "insert",
			Table:  "gate_runs",
			PK:     map[string]any{"id": gateID},
			Values: values,
		}}

		if err := appendAndProjectEvent(options.WorkspaceRoot, conn, lock, event); err != nil {
			return err
		}

		if options.ExitCode == 0 {
			gatesPassed++
		} else {
			gatesFailed++
		}
		result = GateResult{
			GateID:        gateID,
			IssueID:       issueID,
			AttemptNumber: attemptNumber,
			AgentRunID:    agentRunID,
			Command:       command,
			ExitCode:      options.ExitCode,
			GatesPassed:   gatesPassed,
			GatesFailed:   gatesFailed,
		}
		return nil
	})
	if err != nil {
		return GateResult{}, err
	}

	return result, nil
}

func gateArgv(issueID, actorRole, command string, exitCode int, outputExcerpt *string) []string {
	argv := []string{"fid", "dev", "gate", issueID, "--cmd", command, "--exit-code", fmt.Sprintf("%d", exitCode)}
	if value := trimmedOptionalString(outputExcerpt); value != nil {
		argv = append(argv, "--output-excerpt", *value)
	}
	return append(argv, "--by", actorRole)
}

func trimmedOptionalString(value *string) *string {
	if value == nil {
		return nil
	}
	trimmed := strings.TrimSpace(*value)
	if trimmed == "" {
		return nil
	}
	return &trimmed
}
