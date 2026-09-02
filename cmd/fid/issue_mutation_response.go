package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/mbell-nzl/fid/internal/cli"
	"github.com/mbell-nzl/fid/internal/issues"
	"github.com/mbell-nzl/fid/internal/project"
	"github.com/spf13/cobra"
)

func issueMutationWorkspaceRoot(cmd *cobra.Command) (string, error) {
	workingDirectory, err := commandWorkingDirectory(cmd)
	if err != nil {
		return "", err
	}
	projectPath, err := rebuildProjectCandidate(cmd, workingDirectory)
	if err != nil {
		return "", err
	}
	return project.Resolve(workingDirectory, projectPath)
}

func issueExpectedRevision(cmd *cobra.Command, raw string) (*int, error) {
	if !cmd.Flags().Changed("expected-revision") {
		return nil, nil
	}
	if raw == "" || strings.Trim(raw, "0123456789") != "" {
		return nil, fmt.Errorf("invalid expected revision %q", raw)
	}
	value, err := strconv.ParseInt(raw, 10, 64)
	if err != nil || value < 1 || int64(int(value)) != value {
		return nil, fmt.Errorf("invalid expected revision %q", raw)
	}
	result := int(value)
	return &result, nil
}

func renderIssueCreateSuccess(cmd *cobra.Command, result issues.CreateResult) error {
	jsonMode, err := readJSONMode(cmd)
	if err != nil {
		return err
	}
	if jsonMode {
		return cli.WriteSuccess(cmd.OutOrStdout(), "issue.create", map[string]any{
			"issue": map[string]any{
				"id": result.ID, "kind": result.Kind, "title": result.Title,
				"status": result.Status, "revision": result.Revision,
			},
			"idempotency_key": result.IdempotencyKey,
			"replayed":        result.Replayed,
		}, &cli.Event{ID: result.EventID, OccurredAt: result.OccurredAt})
	}
	return writeIssueLine(cmd.OutOrStdout(), "created issue %s (%s) %q\n", result.ID, result.Kind, result.Title)
}

func renderIssueUpdateSuccess(cmd *cobra.Command, result issues.UpdateResult) error {
	jsonMode, err := readJSONMode(cmd)
	if err != nil {
		return err
	}
	if jsonMode {
		changes := make([]map[string]any, 0, len(result.Changes))
		for _, change := range result.Changes {
			changes = append(changes, map[string]any{"field": change.Field, "from": change.From, "to": change.To})
		}
		return cli.WriteSuccess(cmd.OutOrStdout(), "issue.update", map[string]any{
			"issue":           map[string]any{"id": result.ID, "revision": result.Revision},
			"changes":         changes,
			"idempotency_key": result.IdempotencyKey,
			"replayed":        result.Replayed,
		}, &cli.Event{ID: result.EventID, OccurredAt: result.OccurredAt})
	}
	return writeIssueLine(cmd.OutOrStdout(), "updated %s:\n%s\n", result.ID, issueUpdateSummary(result.Changes))
}

func renderIssueTransitionSuccess(cmd *cobra.Command, result issues.TransitionResult) error {
	jsonMode, err := readJSONMode(cmd)
	if err != nil {
		return err
	}
	if jsonMode {
		return cli.WriteSuccess(cmd.OutOrStdout(), "issue.transition", map[string]any{
			"issue":           map[string]any{"id": result.ID, "revision": result.Revision},
			"from_state":      result.FromState,
			"to_state":        result.ToState,
			"idempotency_key": result.IdempotencyKey,
			"replayed":        result.Replayed,
		}, &cli.Event{ID: result.EventID, OccurredAt: result.OccurredAt})
	}
	return writeIssueLine(cmd.OutOrStdout(), "transitioned %s: %s → %s\n", result.ID, result.FromState, result.ToState)
}

func renderIssueMutationError(cmd *cobra.Command, operation, issueID string, mutationErr error) error {
	jsonMode, err := readJSONMode(cmd)
	if err != nil {
		return err
	}
	failure, err := issueMutationFailure(cmd, operation, issueID, mutationErr)
	if err != nil {
		return err
	}
	if !jsonMode {
		if err := renderIssueError(cmd.ErrOrStderr(), mutationErr); err != nil {
			return err
		}
		return cli.NewFailureError(failure, mutationErr)
	}
	return cli.WriteFailureAndReturn(cmd.ErrOrStderr(), operation, failure, mutationErr)
}

func issueMutationFlagError(operation string) func(*cobra.Command, error) error {
	return func(cmd *cobra.Command, flagErr error) error {
		jsonMode, err := readJSONMode(cmd)
		if err != nil {
			return err
		}
		failure, err := cli.NewFailure(cli.CodeInvalidArgument, "invalid Issue mutation request", nil, map[string]any{
			"argument": "request",
			"value":    nil,
			"allowed":  []string{"accepted " + operation + " input"},
		})
		if err != nil {
			return err
		}
		if !jsonMode {
			if err := writeIssueLine(cmd.ErrOrStderr(), "error: invalid Issue mutation request\n"); err != nil {
				return err
			}
			return cli.NewFailureError(failure, flagErr)
		}
		return cli.WriteFailureAndReturn(cmd.ErrOrStderr(), operation, failure, flagErr)
	}
}

func issueMutationFailure(cmd *cobra.Command, operation, issueID string, mutationErr error) (cli.Failure, error) {
	if errors.Is(mutationErr, project.ErrNotFound) {
		workingDirectory, err := commandWorkingDirectory(cmd)
		if err != nil {
			return cli.Failure{}, err
		}
		path, err := rebuildProjectCandidate(cmd, workingDirectory)
		if err != nil {
			return cli.Failure{}, err
		}
		return cli.NewFailure(cli.CodeProjectNotFound, "fid project was not found", &cli.Subject{Type: cli.SubjectProject, ID: path}, map[string]any{"path": path})
	}
	if failure, handled, err := lockUnavailableFailure(mutationErr); err != nil {
		return cli.Failure{}, err
	} else if handled {
		return failure, nil
	}
	if failure, handled, err := journalReplayStoppedFailure(mutationErr); err != nil {
		return cli.Failure{}, err
	} else if handled {
		return failure, nil
	}

	var stateConflict *issues.StateConflictError
	if errors.As(mutationErr, &stateConflict) {
		return cli.NewFailure(cli.CodeStateConflict, "issue revision does not match the expected revision", &cli.Subject{Type: cli.SubjectIssue, ID: stateConflict.IssueID}, map[string]any{"expected_revision": stateConflict.ExpectedRevision, "actual_revision": stateConflict.ActualRevision})
	}
	var notFound *issues.NotFoundError
	if errors.As(mutationErr, &notFound) {
		return cli.NewFailure(cli.CodeIssueNotFound, "issue was not found", &cli.Subject{Type: cli.SubjectIssue, ID: notFound.IssueID}, map[string]any{"id": notFound.IssueID})
	}
	var labelNotFound *issues.LabelNotFoundError
	if errors.As(mutationErr, &labelNotFound) {
		return cli.NewFailure(cli.CodeLabelNotFound, "label was not found", &cli.Subject{Type: cli.SubjectLabel, ID: labelNotFound.Label}, map[string]any{"id": labelNotFound.Label})
	}
	var duplicateLabel *issues.DuplicateIssueLabelError
	if errors.As(mutationErr, &duplicateLabel) {
		return cli.NewFailure(cli.CodeDuplicateEntity, "entity already exists", &cli.Subject{Type: cli.SubjectIssue, ID: duplicateLabel.IssueID}, map[string]any{"entity_type": "issue_label", "identity": map[string]any{"fields": []string{"issue_id", "label"}, "values": []string{duplicateLabel.IssueID, duplicateLabel.Label}}})
	}
	var duplicateAC *issues.DuplicateAcceptanceCriterionError
	if errors.As(mutationErr, &duplicateAC) {
		return cli.NewFailure(cli.CodeDuplicateEntity, "entity already exists", &cli.Subject{Type: cli.SubjectIssue, ID: duplicateAC.IssueID}, map[string]any{"entity_type": "acceptance_criterion", "identity": map[string]any{"fields": []string{"issue_id", "ordinal"}, "values": []string{duplicateAC.IssueID, fmt.Sprintf("%d", duplicateAC.Ordinal)}}})
	}
	var limit *issues.ResultLimitExceededError
	if errors.As(mutationErr, &limit) {
		return cli.NewFailure(cli.CodeResultLimitExceeded, "result exceeds its collection limit", nil, map[string]any{"collection": limit.Collection, "limit": limit.Limit, "observed": limit.Observed})
	}
	var unsupportedKind *issues.UnsupportedKindError
	if errors.As(mutationErr, &unsupportedKind) {
		return cli.NewFailure(cli.CodeUnsupportedKind, "issue kind is unsupported", nil, map[string]any{"argument": "kind", "value": unsupportedKind.Kind, "allowed": []string{"epic", "task", "subtask"}})
	}
	var hierarchy *issues.HierarchyInvalidError
	if errors.As(mutationErr, &hierarchy) {
		return cli.NewFailure(cli.CodeHierarchyInvalid, "issue hierarchy is invalid", &cli.Subject{Type: cli.SubjectIssue, ID: issueID}, map[string]any{"kind": hierarchy.Kind, "parent_id": optionalString(hierarchy.ParentID), "parent_kind": optionalString(hierarchy.ParentKind), "rule": hierarchy.Rule})
	}
	var invalidTransition *issues.InvalidTransitionError
	if errors.As(mutationErr, &invalidTransition) {
		return transitionNotPermittedFailure(invalidTransition.IssueID, invalidTransition.CurrentState, invalidTransition.RequestedState, nil)
	}
	var executionTransition *issues.ExecutionBoundTransitionError
	if errors.As(mutationErr, &executionTransition) {
		unsatisfied := make([]string, 0, len(executionTransition.Reasons))
		for _, reason := range executionTransition.Reasons {
			unsatisfied = append(unsatisfied, reason.Code)
		}
		return transitionNotPermittedFailure(executionTransition.IssueID, executionTransition.CurrentState, executionTransition.RequestedState, unsatisfied)
	}
	var blockedTransition *issues.BlockedTransitionError
	if errors.As(mutationErr, &blockedTransition) {
		unsatisfied := make([]string, 0, len(blockedTransition.Unmet))
		for _, unmet := range blockedTransition.Unmet {
			unsatisfied = append(unsatisfied, unmet.Rule)
		}
		return transitionNotPermittedFailure(blockedTransition.IssueID, blockedTransition.CurrentState, blockedTransition.RequestedState, unsatisfied)
	}
	return cli.NewFailure(cli.CodeInvalidArgument, "invalid Issue mutation request", nil, map[string]any{"argument": "request", "value": nil, "allowed": []string{"accepted " + operation + " input"}})
}

func transitionNotPermittedFailure(issueID, fromState, toState string, unsatisfied []string) (cli.Failure, error) {
	if unsatisfied == nil {
		unsatisfied = []string{}
	}
	return cli.NewFailure(cli.CodeTransitionNotPermitted, "Issue transition is not permitted", &cli.Subject{Type: cli.SubjectIssue, ID: issueID}, map[string]any{"from_state": fromState, "to_state": toState, "unsatisfied": unsatisfied})
}

func optionalString(value *string) any {
	if value == nil {
		return nil
	}
	return *value
}
