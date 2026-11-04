package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"sort"

	"github.com/mbell-nzl/fid/internal/cli"
	"github.com/mbell-nzl/fid/internal/journal"
	"github.com/mbell-nzl/fid/internal/policy"
	"github.com/mbell-nzl/fid/internal/review"
	"github.com/mbell-nzl/fid/internal/workflow"
	"github.com/spf13/cobra"
)

func init() {
	rootFlagCapabilities["human-outcome record"] = rootFlagCapability{Project: true, JSON: true}
	rootFlagCapabilities["human-outcome get"] = rootFlagCapability{Project: true, JSON: true}
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
	stopIDs := []string{stop.ID}
	attentionIDs := append([]string{}, blocked.AttentionIDs...)
	sort.Strings(stopIDs)
	sort.Strings(attentionIDs)
	details["stop_ids"] = stopIDs
	details["attention_ids"] = attentionIDs
	details["event"] = nil
	details["replayed"] = false
	return cli.NewFailure(cli.CodeHumanOutcomeRequired, "human-outcome Obligation is unsatisfied", &cli.Subject{Type: cli.SubjectEnvelope, ID: stop.Subject.ID}, details)
}

func newHumanOutcomeCommand() *cobra.Command {
	command := &cobra.Command{Use: "human-outcome", Short: "Record and list Human Outcomes."}
	command.AddCommand(newHumanOutcomeRecordCommand(), newHumanOutcomeListCommand(), newHumanOutcomeGetCommand())
	return command
}

func newHumanOutcomeRecordCommand() *cobra.Command {
	var inputPath, actor, key string
	command := &cobra.Command{Use: "record", Short: "Record a Human Outcome.", Args: canonicalNoArgs("human-outcome.record"), RunE: func(cmd *cobra.Command, _ []string) error {
		for _, flag := range []string{"input", "actor", "idempotency-key"} {
			if err := requireCanonicalFlag(cmd, "human-outcome.record", flag); err != nil {
				return err
			}
		}
		var request review.HumanOutcomeRequest
		if err := readDescriptorInput(cmd.InOrStdin(), inputPath, inputDescriptorRegistry[rootCommandKey(cmd)], &request); err != nil {
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
		return writeP23MutationHuman(cmd, "human-outcome.record", result.Outcome, nil, result.ProtectedDecisions, result.Stops, result.Attention)
	}}
	command.Flags().StringVar(&inputPath, "input", "", "Human Outcome JSON input")
	command.Flags().StringVar(&actor, "actor", "", "Finalized activation human actor ID")
	command.Flags().StringVar(&key, "idempotency-key", "", "Idempotency key")
	command.SetFlagErrorFunc(canonicalFlagError("human-outcome.record"))
	installInputDescriptor(command, inputDescriptorRegistry["human-outcome record"])
	return command
}

func newHumanOutcomeListCommand() *cobra.Command {
	var class, result, issueID string
	var limit int
	command := &cobra.Command{Use: "list", Short: "List Human Outcomes.", Args: canonicalNoArgs("human-outcome.list"), RunE: func(cmd *cobra.Command, _ []string) error {
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
		return writeP23Human(cmd, "human-outcome.list", items...)
	}}
	command.Flags().StringVar(&class, "class", "", "Filter by Human Outcome class")
	command.Flags().StringVar(&result, "result", "", "Filter by result")
	command.Flags().StringVar(&issueID, "issue", "", "Filter by Issue scope")
	command.Flags().IntVar(&limit, "limit", 100, "Maximum Human Outcomes to return")
	command.SetFlagErrorFunc(canonicalFlagError("human-outcome.list"))
	return command
}

func renderHumanOutcomeError(cmd *cobra.Command, operation string, cause error) error {
	if refused := renderRepositoryRefusal(cmd, cause); refused != nil {
		return refused
	}
	jsonMode, err := readJSONMode(cmd)
	if err != nil {
		return err
	}
	if rendered, ok := renderJudgmentObservationError(cmd, jsonMode, operation, cause); ok {
		return rendered
	}
	var missing *review.HumanOutcomeNotFoundError
	if errors.As(cause, &missing) {
		return writeClaimFailure(cmd, jsonMode, operation, cause, cli.CodeOutcomeNotFound, "Human Outcome was not found", &cli.Subject{Type: cli.SubjectHumanOutcome, ID: missing.ID}, map[string]any{"id": missing.ID})
	}
	var authorityMissing *review.AuthorityMissingError
	if errors.As(cause, &authorityMissing) {
		return writeClaimFailure(cmd, jsonMode, operation, cause, cli.CodeAuthorityMissing, "active Kit Activation authority is missing", nil, map[string]any{"operation": authorityMissing.Operation, "actor_id": authorityMissing.ActorID, "required_actions": []string{authorityMissing.Operation}, "required_capabilities": []string{authorityMissing.Operation}})
	}
	var denied *policy.EnvelopeAuthorityDeniedError
	if errors.As(cause, &denied) {
		return writeClaimFailure(cmd, jsonMode, operation, cause, cli.CodeAuthorityDenied, "active envelope authority is denied", nil, map[string]any{"operation": operation, "actor_id": denied.ActorID, "required_actions": []string{denied.Action}, "required_capabilities": []string{denied.Capability}})
	}
	var conflict *journal.IdempotencyConflictError
	if errors.As(cause, &conflict) {
		return writeClaimFailure(cmd, jsonMode, operation, cause, cli.CodeIdempotencyConflict, "idempotency key conflicts with a different request", nil, idempotencyConflictDetails(conflict))
	}
	var resultLimit *review.ResultLimitExceededError
	if errors.As(cause, &resultLimit) {
		return writeClaimFailure(cmd, jsonMode, operation, cause, cli.CodeResultLimitExceeded, "Human Outcome result limit exceeded", nil, map[string]any{"collection": "items", "limit": resultLimit.Limit, "observed": resultLimit.Observed})
	}
	return writeGenericFailure(cmd, jsonMode, operation, cause, cli.CodeInvalidArgument, "invalid Human Outcome request", nil, map[string]any{"argument": "input", "value": nil, "allowed": []string{"accepted " + operation + " input"}})
}

func newHumanOutcomeGetCommand() *cobra.Command {
	cmd := &cobra.Command{Use: "get <outcome-id>", Short: "Read one Human Outcome.", Args: canonicalExactArgs("human-outcome.get", 1), RunE: func(cmd *cobra.Command, args []string) error {
		workspace, err := issueMutationWorkspaceRoot(cmd)
		if err != nil {
			return renderHumanOutcomeError(cmd, "human-outcome.get", err)
		}
		value, err := review.GetHumanOutcome(workspace, issueMigrationsDir(), args[0])
		if err != nil {
			return renderHumanOutcomeError(cmd, "human-outcome.get", err)
		}
		jsonMode, err := readJSONMode(cmd)
		if err != nil {
			return err
		}
		if jsonMode {
			return cli.WriteSuccess(cmd.OutOrStdout(), "human-outcome.get", map[string]any{"outcome": value}, nil)
		}
		return writeP23Human(cmd, "human-outcome.get", value)
	}}
	cmd.SetFlagErrorFunc(canonicalFlagError("human-outcome.get"))
	return cmd
}
