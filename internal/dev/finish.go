package dev

import (
	"database/sql"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/mbell-nzl/fid/internal/claims"
	"github.com/mbell-nzl/fid/internal/issues"
	"github.com/mbell-nzl/fid/internal/journal"
	"github.com/mbell-nzl/fid/internal/policy"
	"github.com/mbell-nzl/fid/internal/workflow"
)

type FinishOptions struct {
	WorkspaceRoot         string
	MigrationsDir         string
	IssueID               string
	ActorRole             string
	CompletionMDPath      string
	ActualTokens          *int
	CommitHashes          []string
	FilesTouched          []string
	HumanIntervention     bool
	HumanInterventionNote *string
}

type FinishResult struct {
	IssueID    string
	FromState  string
	ToState    string
	RunID      *string
	Completion issues.CompleteResult
}

func Finish(options FinishOptions) (FinishResult, error) {
	if strings.TrimSpace(options.WorkspaceRoot) == "" {
		return FinishResult{}, fmt.Errorf("dev finish: workspace root is required")
	}
	if strings.TrimSpace(options.MigrationsDir) == "" {
		return FinishResult{}, fmt.Errorf("dev finish: migrations dir is required")
	}
	issueID := strings.TrimSpace(options.IssueID)
	if issueID == "" {
		return FinishResult{}, fmt.Errorf("dev finish: issue id is required")
	}
	completionMDPath := strings.TrimSpace(options.CompletionMDPath)
	if completionMDPath == "" {
		return FinishResult{}, fmt.Errorf("dev finish: completion report path is required")
	}
	if err := policy.Authorize(options.WorkspaceRoot, options.ActorRole, actionDevFinish); err != nil {
		return FinishResult{}, err
	}

	commandArgv := finishArgv(issueID, options.ActorRole, completionMDPath, options.ActualTokens, options.CommitHashes, options.FilesTouched, options.HumanIntervention, options.HumanInterventionNote)
	var result FinishResult
	err := withDevConn(options.WorkspaceRoot, options.MigrationsDir, commandArgv, options.ActorRole, func(lock *journal.WriteLock, conn *sql.DB) error {
		current, err := loadIssueState(conn, issueID)
		if err != nil {
			return err
		}
		if err := requireTaskOrSubtask(current.ID, current.Kind); err != nil {
			return err
		}

		if err := validateFinishTransition(conn, current); err != nil {
			return err
		}

		runID, err := openRunID(conn, issueID)
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
			EventType:   "dev.finished",
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
		completionEffects, completionResult, err := issues.BuildCompletionEffects(conn, issues.CompletionBuildOptions{
			IssueID:               issueID,
			ActorRole:             options.ActorRole,
			Status:                "dev_complete",
			NotesMDPath:           &completionMDPath,
			ActualTokens:          options.ActualTokens,
			CommitHashes:          options.CommitHashes,
			FilesTouched:          options.FilesTouched,
			HumanIntervention:     options.HumanIntervention,
			HumanInterventionNote: trimmedOptionalString(options.HumanInterventionNote),
		}, event.EventID, timestamp)
		if err != nil {
			return err
		}
		effects := make([]journal.Effect, 0, len(completionEffects)+4)
		if runID != nil {
			effects = append(effects, closeRunEffect(*runID, issueID, timestamp))
		}
		effects = append(effects, completionEffects...)
		transitionEffects, err := buildIssueTransitionEffects(current, "dev_complete", options.ActorRole, event.EventID, activityID, timestamp)
		if err != nil {
			return err
		}
		effects = append(effects, transitionEffects...)
		effects = append(effects, releaseEffect)
		event.Effects = effects

		if err := appendAndProjectEvent(options.WorkspaceRoot, conn, lock, event); err != nil {
			return err
		}

		result = FinishResult{
			IssueID:    issueID,
			FromState:  current.Status,
			ToState:    "dev_complete",
			RunID:      runID,
			Completion: completionResult,
		}
		return nil
	})
	if err != nil {
		return FinishResult{}, err
	}

	return result, nil
}

func validateFinishTransition(conn *sql.DB, current issueState) error {
	registry, err := workflow.Load(conn)
	if err != nil {
		return fmt.Errorf("dev finish: load workflow registry: %w", err)
	}
	rules, allowed := registry.RulesForEdge(current.WorkflowID, current.Status, "dev_complete")
	if !allowed {
		return &issues.InvalidTransitionError{
			IssueID:        current.ID,
			CurrentState:   current.Status,
			RequestedState: "dev_complete",
		}
	}
	unmet, err := finishUnmetRules(conn, current.ID, rules)
	if err != nil {
		return err
	}
	if len(unmet) == 0 {
		return nil
	}
	return &issues.BlockedTransitionError{
		IssueID:        current.ID,
		CurrentState:   current.Status,
		RequestedState: "dev_complete",
		Unmet:          unmet,
	}
}

func finishUnmetRules(conn *sql.DB, issueID string, rules []workflow.RuleReference) ([]workflow.UnmetCondition, error) {
	unmet := []workflow.UnmetCondition{}
	for _, rule := range rules {
		switch rule.Name {
		case workflow.RuleAllACsVerified:
			criteria, err := loadAcceptanceCriterionVerifications(conn, issueID)
			if err != nil {
				return nil, err
			}
			unmet = append(unmet, workflow.AllACsVerified(rule.RequiredRoles, criteria).Unmet...)
		case workflow.RuleHasCompletionReport:
			continue
		default:
			return nil, fmt.Errorf("dev finish: unsupported workflow rule %q", rule.Name)
		}
	}
	return unmet, nil
}

func loadAcceptanceCriterionVerifications(conn *sql.DB, issueID string) ([]workflow.AcceptanceCriterionVerification, error) {
	rows, err := conn.Query(`
		SELECT id
		FROM acceptance_criteria
		WHERE issue_id = ?
		ORDER BY ordinal, id
	`, issueID)
	if err != nil {
		return nil, fmt.Errorf("query acceptance criteria for %q: %w", issueID, err)
	}
	defer closeDevRows(rows)

	criteria := []workflow.AcceptanceCriterionVerification{}
	for rows.Next() {
		var criterionID string
		if err := rows.Scan(&criterionID); err != nil {
			return nil, fmt.Errorf("scan acceptance criterion for %q: %w", issueID, err)
		}
		passingRoles, err := latestPassingRoles(conn, criterionID)
		if err != nil {
			return nil, err
		}
		criteria = append(criteria, workflow.AcceptanceCriterionVerification{
			AcceptanceCriterionID: criterionID,
			PassingRoles:          passingRoles,
		})
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate acceptance criteria for %q: %w", issueID, err)
	}
	return criteria, nil
}

func latestPassingRoles(conn *sql.DB, criterionID string) ([]string, error) {
	rows, err := conn.Query(`
		SELECT verifier_role, status
		FROM ac_verifications
		WHERE ac_id = ?
		ORDER BY verifier_role, verified_at DESC, id DESC
	`, criterionID)
	if err != nil {
		return nil, fmt.Errorf("query acceptance criterion verifications for %q: %w", criterionID, err)
	}
	defer closeDevRows(rows)

	passingRoles := []string{}
	seen := map[string]bool{}
	for rows.Next() {
		var role, status string
		if err := rows.Scan(&role, &status); err != nil {
			return nil, fmt.Errorf("scan acceptance criterion verification for %q: %w", criterionID, err)
		}
		if seen[role] {
			continue
		}
		seen[role] = true
		if status == "pass" {
			passingRoles = append(passingRoles, role)
		}
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate acceptance criterion verifications for %q: %w", criterionID, err)
	}
	sort.Strings(passingRoles)
	return passingRoles, nil
}

func closeRunEffect(runID, issueID, timestamp string) journal.Effect {
	return journal.Effect{
		Op:    "update",
		Table: "agent_runs",
		PK:    map[string]any{"id": runID},
		Values: map[string]any{
			"ended_at": timestamp,
		},
		Preconditions: map[string]any{
			"issue_id": issueID,
			"ended_at": nil,
		},
	}
}

func finishArgv(
	issueID, actorRole, completionMDPath string,
	actualTokens *int,
	commitHashes, filesTouched []string,
	humanIntervention bool,
	humanInterventionNote *string,
) []string {
	argv := []string{"fid", "dev", "finish", issueID, "--completion-md", completionMDPath}
	if actualTokens != nil {
		argv = append(argv, "--tokens", fmt.Sprintf("%d", *actualTokens))
	}
	if len(compactTrimmedStrings(commitHashes)) > 0 {
		argv = append(argv, "--commits", strings.Join(compactTrimmedStrings(commitHashes), ", "))
	}
	for _, path := range compactTrimmedStrings(filesTouched) {
		argv = append(argv, "--files", path)
	}
	if humanIntervention {
		argv = append(argv, "--human-intervention")
	}
	if note := trimmedOptionalString(humanInterventionNote); note != nil {
		argv = append(argv, "--intervention-note", *note)
	}
	return append(argv, "--by", actorRole)
}

func compactTrimmedStrings(values []string) []string {
	result := make([]string, 0, len(values))
	for _, value := range values {
		trimmed := strings.TrimSpace(value)
		if trimmed == "" {
			continue
		}
		result = append(result, trimmed)
	}
	return result
}
