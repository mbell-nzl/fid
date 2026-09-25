package main

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/mbell-nzl/fid/internal/cli"
	"github.com/mbell-nzl/fid/internal/journal"
	"github.com/mbell-nzl/fid/internal/policy"
	"github.com/mbell-nzl/fid/internal/review"
	"github.com/mbell-nzl/fid/internal/workflow"
	"github.com/spf13/cobra"
)

func init() {
	rootFlagCapabilities["human-outcome record"] = rootFlagCapability{Project: true, JSON: true}
	rootFlagCapabilities["human-outcome list"] = rootFlagCapability{Project: true, JSON: true}
}

// humanOutcomeRequiredFailure renders the human-outcome-required precedence
// case, mirroring protectedDecisionFailure for its Stop-carried twin.
func humanOutcomeRequiredFailure(blocked *workflow.StopBlockedError) (cli.Failure, error) {
	if blocked == nil || len(blocked.Stops) != 1 || blocked.Stops[0].Code != "human-outcome-required" {
		return cli.Failure{}, fmt.Errorf("human-outcome-required denial is invalid")
	}
	stop := blocked.Stops[0]
	var details map[string]any
	if err := json.Unmarshal(stop.Details, &details); err != nil {
		return cli.Failure{}, fmt.Errorf("decode human-outcome-required details: %w", err)
	}
	return cli.NewFailure(cli.CodeHumanOutcomeRequired, "human-outcome Obligation is unsatisfied", &cli.Subject{Type: cli.SubjectEnvelope, ID: stop.Subject.ID}, details)
}

func newHumanOutcomeCommand() *cobra.Command {
	command := &cobra.Command{Use: "human-outcome", Short: "Record and list Human Outcomes."}
	command.AddCommand(newHumanOutcomeRecordCommand(), newHumanOutcomeListCommand())
	return command
}

func newHumanOutcomeRecordCommand() *cobra.Command {
	var inputPath, actor, key string
	command := &cobra.Command{Use: "record", Args: canonicalNoArgs("human-outcome.record"), RunE: func(cmd *cobra.Command, _ []string) error {
		for _, flag := range []string{"input", "actor", "idempotency-key"} {
			if err := requireCanonicalFlag(cmd, "human-outcome.record", flag); err != nil {
				return err
			}
		}
		var request review.HumanOutcomeRequest
		if err := decodeClosedJudgmentInput(cmd.InOrStdin(), inputPath, []string{"id", "class", "subject", "result", "scope", "reasoning", "references"}, &request); err != nil {
			return renderHumanOutcomeError(cmd, "human-outcome.record", err)
		}
		workspace, err := issueMutationWorkspaceRoot(cmd)
		if err != nil {
			return renderHumanOutcomeError(cmd, "human-outcome.record", err)
		}
		result, err := review.RecordHumanOutcome(review.JudgmentOptions{WorkspaceRoot: workspace, MigrationsDir: issueMigrationsDir(), ActorID: actor, IdempotencyKey: key}, request)
		if err != nil {
			return renderHumanOutcomeError(cmd, "human-outcome.record", err)
		}
		jsonMode, err := readJSONMode(cmd)
		if err != nil {
			return err
		}
		if jsonMode {
			value := map[string]any{"outcome": result.Outcome, "protected_decisions": result.ProtectedDecisions, "stops": result.Stops, "attention": result.Attention, "idempotency_key": result.IdempotencyKey, "replayed": result.Replayed}
			return cli.WriteSuccess(cmd.OutOrStdout(), "human-outcome.record", value, &cli.Event{ID: result.EventID, OccurredAt: result.Outcome["created_at"].(string)})
		}
		_, err = fmt.Fprintf(cmd.OutOrStdout(), "recorded Human Outcome %s\n", result.Outcome["id"])
		return err
	}}
	command.Flags().StringVar(&inputPath, "input", "", "Human Outcome JSON input")
	command.Flags().StringVar(&actor, "actor", "", "Finalized activation human actor ID")
	command.Flags().StringVar(&key, "idempotency-key", "", "Idempotency key")
	command.SetFlagErrorFunc(canonicalFlagError("human-outcome.record"))
	return command
}

func newHumanOutcomeListCommand() *cobra.Command {
	var class, result, issueID string
	var limit int
	command := &cobra.Command{Use: "list", Args: canonicalNoArgs("human-outcome.list"), RunE: func(cmd *cobra.Command, _ []string) error {
		workspace, err := issueMutationWorkspaceRoot(cmd)
		if err != nil {
			return renderHumanOutcomeError(cmd, "human-outcome.list", err)
		}
		items, err := review.ListHumanOutcomes(review.HumanOutcomeListOptions{WorkspaceRoot: workspace, MigrationsDir: issueMigrationsDir(), Class: class, Result: result, IssueID: issueID, Limit: limit})
		if err != nil {
			return renderHumanOutcomeError(cmd, "human-outcome.list", err)
		}
		jsonMode, err := readJSONMode(cmd)
		if err != nil {
			return err
		}
		if jsonMode {
			values := make([]any, len(items))
			for index := range items {
				values[index] = items[index]
			}
			return cli.WriteSuccess(cmd.OutOrStdout(), "human-outcome.list", map[string]any{"items": values}, nil)
		}
		for _, item := range items {
			if _, err := fmt.Fprintf(cmd.OutOrStdout(), "%s\t%s\n", item["id"], item["result"]); err != nil {
				return err
			}
		}
		return nil
	}}
	command.Flags().StringVar(&class, "class", "", "Filter by Human Outcome class")
	command.Flags().StringVar(&result, "result", "", "Filter by result")
	command.Flags().StringVar(&issueID, "issue", "", "Filter by Issue scope")
	command.Flags().IntVar(&limit, "limit", 100, "Maximum Human Outcomes to return")
	command.SetFlagErrorFunc(canonicalFlagError("human-outcome.list"))
	return command
}

func renderHumanOutcomeError(cmd *cobra.Command, operation string, cause error) error {
	jsonMode, err := readJSONMode(cmd)
	if err != nil {
		return err
	}
	var authorityMissing *review.AuthorityMissingError
	if errors.As(cause, &authorityMissing) {
		return writeClaimFailure(cmd, jsonMode, operation, cause, cli.CodeAuthorityMissing, "active Kit Activation authority is missing", nil, map[string]any{"operation": authorityMissing.Operation, "actor": authorityMissing.ActorID})
	}
	var denied *policy.EnvelopeAuthorityDeniedError
	if errors.As(cause, &denied) {
		return writeClaimFailure(cmd, jsonMode, operation, cause, cli.CodeAuthorityDenied, "active envelope authority is denied", nil, map[string]any{"operation": operation, "actor": denied.ActorID})
	}
	var conflict *journal.IdempotencyConflictError
	if errors.As(cause, &conflict) {
		return writeClaimFailure(cmd, jsonMode, operation, cause, cli.CodeIdempotencyConflict, "idempotency key conflicts with a different request", nil, map[string]any{"operation_scope": conflict.OperationScope, "idempotency_key": conflict.IdempotencyKey, "stored_request_digest": conflict.StoredRequestDigest, "supplied_request_digest": conflict.SuppliedRequestDigest})
	}
	var resultLimit *review.ResultLimitExceededError
	if errors.As(cause, &resultLimit) {
		return writeClaimFailure(cmd, jsonMode, operation, cause, cli.CodeResultLimitExceeded, "Human Outcome result limit exceeded", nil, map[string]any{"collection": "items", "limit": resultLimit.Limit, "observed": resultLimit.Observed})
	}
	return writeClaimFailure(cmd, jsonMode, operation, cause, cli.CodeInvalidArgument, "invalid Human Outcome request", nil, map[string]any{"argument": "input", "value": nil, "allowed": []string{"accepted " + operation + " input"}})
}
