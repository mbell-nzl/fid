package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"sort"

	"github.com/mbell-nzl/fid/internal/cli"
	"github.com/mbell-nzl/fid/internal/journal"
	"github.com/mbell-nzl/fid/internal/policy"
	"github.com/mbell-nzl/fid/internal/project"
	"github.com/mbell-nzl/fid/internal/review"
	"github.com/mbell-nzl/fid/internal/workflow"
	"github.com/spf13/cobra"
)

func init() {
	rootFlagCapabilities["decision record-reference"] = rootFlagCapability{Project: true, JSON: true}
	rootFlagCapabilities["decision list"] = rootFlagCapability{Project: true, JSON: true}
}

func protectedDecisionFailure(blocked *workflow.StopBlockedError) (cli.Failure, error) {
	if blocked == nil || len(blocked.Stops) != 1 || blocked.Stops[0].Code != "protected-decision-required" {
		return cli.Failure{}, fmt.Errorf("protected Decision denial is invalid")
	}
	stop := blocked.Stops[0]
	var details map[string]any
	if err := json.Unmarshal(stop.Details, &details); err != nil {
		return cli.Failure{}, fmt.Errorf("decode protected Decision details: %w", err)
	}
	stopIDs := []string{stop.ID}
	attentionIDs := append([]string{}, blocked.AttentionIDs...)
	sort.Strings(stopIDs)
	sort.Strings(attentionIDs)
	details["stop_ids"] = stopIDs
	details["attention_ids"] = attentionIDs
	details["event"] = nil
	details["replayed"] = false
	return cli.NewFailure(cli.CodeProtectedDecisionRequired, "protected Decision requires an approved reference and Human Outcome", &cli.Subject{Type: cli.SubjectEnvelope, ID: stop.Subject.ID}, details)
}

func newDecisionCommand() *cobra.Command {
	command := &cobra.Command{Use: "decision", Short: "Record and list Decision References."}
	command.AddCommand(newDecisionRecordReferenceCommand(), newDecisionListCommand())
	return command
}

func newDecisionRecordReferenceCommand() *cobra.Command {
	var inputPath, actor, key string
	command := &cobra.Command{Use: "record-reference", Short: "Record a Decision Reference.", Args: canonicalNoArgs("decision.record-reference"), RunE: func(cmd *cobra.Command, _ []string) error {
		for _, flag := range []string{"input", "actor", "idempotency-key"} {
			if err := requireCanonicalFlag(cmd, "decision.record-reference", flag); err != nil {
				return err
			}
		}
		var request review.DecisionReferenceRequest
		if err := readDescriptorInput(cmd.InOrStdin(), inputPath, inputDescriptorRegistry[rootCommandKey(cmd)], &request); err != nil {
			return renderDecisionError(cmd, "decision.record-reference", err)
		}
		workspace, err := issueMutationWorkspaceRoot(cmd)
		if err != nil {
			return renderDecisionError(cmd, "decision.record-reference", err)
		}
		result, err := review.RecordDecisionReference(review.JudgmentOptions{WorkspaceRoot: workspace, MigrationsDir: issueMigrationsDir(), ActorID: actor, IdempotencyKey: key}, request)
		if err != nil {
			return renderDecisionError(cmd, "decision.record-reference", err)
		}
		jsonMode, err := readJSONMode(cmd)
		if err != nil {
			return err
		}
		if jsonMode {
			value := map[string]any{"decision": result.Decision, "protected_decisions": result.ProtectedDecisions, "stops": result.Stops, "attention": result.Attention, "idempotency_key": result.IdempotencyKey, "replayed": result.Replayed}
			return cli.WriteSuccess(cmd.OutOrStdout(), "decision.record-reference", value, &cli.Event{ID: result.EventID, OccurredAt: result.Decision["created_at"].(string)})
		}
		return writeP23MutationHuman(cmd, "decision.record-reference", result.Decision, nil, result.ProtectedDecisions, result.Stops, result.Attention)
	}}
	command.Flags().StringVar(&inputPath, "input", "", "Decision Reference JSON input")
	command.Flags().StringVar(&actor, "actor", "", "Finalized activation actor ID")
	command.Flags().StringVar(&key, "idempotency-key", "", "Idempotency key")
	command.SetFlagErrorFunc(canonicalFlagError("decision.record-reference"))
	installInputDescriptor(command, inputDescriptorRegistry["decision record-reference"])
	return command
}

func newDecisionListCommand() *cobra.Command {
	var class, outcome, issueID string
	var limit int
	command := &cobra.Command{Use: "list", Short: "List Decision References.", Args: canonicalNoArgs("decision.list"), RunE: func(cmd *cobra.Command, _ []string) error {
		workspace, err := issueMutationWorkspaceRoot(cmd)
		if err != nil {
			return renderDecisionError(cmd, "decision.list", err)
		}
		items, err := review.ListDecisionReferences(review.DecisionListOptions{WorkspaceRoot: workspace, MigrationsDir: issueMigrationsDir(), Class: class, Outcome: outcome, IssueID: issueID, Limit: limit})
		if err != nil {
			return renderDecisionError(cmd, "decision.list", err)
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
			return cli.WriteSuccess(cmd.OutOrStdout(), "decision.list", map[string]any{"items": values}, nil)
		}
		return writeP23Human(cmd, "decision.list", items...)
	}}
	command.Flags().StringVar(&class, "class", "", "Filter by Decision class")
	command.Flags().StringVar(&outcome, "outcome", "", "Filter by outcome")
	command.Flags().StringVar(&issueID, "issue", "", "Filter by Issue scope")
	command.Flags().IntVar(&limit, "limit", 100, "Maximum Decision References to return")
	command.SetFlagErrorFunc(canonicalFlagError("decision.list"))
	return command
}

func renderDecisionError(cmd *cobra.Command, operation string, cause error) error {
	jsonMode, err := readJSONMode(cmd)
	if err != nil {
		return err
	}
	if rendered, ok := renderJudgmentObservationError(cmd, jsonMode, operation, cause); ok {
		return rendered
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
		return writeClaimFailure(cmd, jsonMode, operation, cause, cli.CodeResultLimitExceeded, "Decision Reference result limit exceeded", nil, map[string]any{"collection": "items", "limit": resultLimit.Limit, "observed": resultLimit.Observed})
	}
	return writeGenericFailure(cmd, jsonMode, operation, cause, cli.CodeInvalidArgument, "invalid Decision Reference request", nil, map[string]any{"argument": "input", "value": nil, "allowed": []string{"accepted " + operation + " input"}})
}

func renderJudgmentObservationError(cmd *cobra.Command, jsonMode bool, operation string, cause error) (error, bool) {
	var structural *cli.InputError
	if errors.As(cause, &structural) {
		return writeClaimFailure(cmd, jsonMode, operation, cause, cli.CodeInvalidArgument,
			structural.Error(), nil, structural.Details()), true
	}
	var mismatch *project.SourceDigestMismatchError
	if errors.As(cause, &mismatch) {
		return writeClaimFailure(cmd, jsonMode, operation, cause, cli.CodeSourceDigestMismatch,
			"Source digest does not match the observed project file", nil,
			map[string]any{"display_label": mismatch.DisplayLabel, "path": optionalString(mismatch.Path),
				"expected_digest": optionalString(mismatch.ExpectedDigest), "observed_digest": optionalString(mismatch.ObservedDigest)}), true
	}
	var unavailable *project.RepositoryUnavailableError
	if errors.As(cause, &unavailable) {
		return writeClaimFailure(cmd, jsonMode, operation, cause, cli.CodeRepositoryUnavailable,
			"repository observation is unavailable", nil,
			map[string]any{"reason": unavailable.Reason, "detail": unavailable.Detail}), true
	}
	var invalid *project.InvalidSourceInputError
	if errors.As(cause, &invalid) {
		return writeClaimFailure(cmd, jsonMode, operation, cause, cli.CodeInvalidArgument,
			invalid.Detail, nil, map[string]any{"argument": invalid.Field, "value": nil,
				"allowed": []string{"valid SourceInput " + invalid.Field}}), true
	}
	return nil, false
}
