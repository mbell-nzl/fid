package main

import (
	"errors"

	"github.com/mbell-nzl/fid/internal/cli"
	"github.com/mbell-nzl/fid/internal/journal"
	"github.com/mbell-nzl/fid/internal/policy"
	"github.com/mbell-nzl/fid/internal/review"
	"github.com/spf13/cobra"
)

func init() {
	rootFlagCapabilities["deviation record"] = rootFlagCapability{Project: true, JSON: true}
	rootFlagCapabilities["deviation list"] = rootFlagCapability{Project: true, JSON: true}
}

func newDeviationCommand() *cobra.Command {
	cmd := &cobra.Command{Use: "deviation", Short: "Record and list append-only Deviations."}
	cmd.AddCommand(newDeviationRecordCommand(), newDeviationListCommand())
	return cmd
}

func newDeviationRecordCommand() *cobra.Command {
	var inputPath, actor, key string
	cmd := &cobra.Command{Use: "record", Short: "Record a Deviation.", Args: canonicalNoArgs("deviation.record"), RunE: func(cmd *cobra.Command, _ []string) error {
		for _, flag := range []string{"input", "actor", "idempotency-key"} {
			if err := requireCanonicalFlag(cmd, "deviation.record", flag); err != nil {
				return err
			}
		}
		var request review.DeviationRecordRequest
		if err := readDescriptorInput(cmd.InOrStdin(), inputPath, inputDescriptorRegistry[rootCommandKey(cmd)], &request); err != nil {
			return renderDeviationError(cmd, "deviation.record", err)
		}
		workspace, err := issueMutationWorkspaceRoot(cmd)
		if err != nil {
			return renderDeviationError(cmd, "deviation.record", err)
		}
		result, err := review.RecordDeviation(review.JudgmentOptions{WorkspaceRoot: workspace, MigrationsDir: issueMigrationsDir(), ActorID: actor, IdempotencyKey: key}, request)
		if err != nil {
			return renderDeviationError(cmd, "deviation.record", err)
		}
		jsonMode, err := readJSONMode(cmd)
		if err != nil {
			return err
		}
		if jsonMode {
			return cli.WriteSuccess(cmd.OutOrStdout(), "deviation.record", map[string]any{"deviation": result.Deviation, "idempotency_key": result.IdempotencyKey, "replayed": result.Replayed}, &cli.Event{ID: result.EventID, OccurredAt: result.Deviation["created_at"].(string)})
		}
		return writeP23Human(cmd, "deviation.record", result.Deviation)
	}}
	cmd.Flags().StringVar(&inputPath, "input", "", "Deviation JSON input")
	cmd.Flags().StringVar(&actor, "actor", "", "Finalized activation actor ID")
	cmd.Flags().StringVar(&key, "idempotency-key", "", "Idempotency key")
	cmd.SetFlagErrorFunc(canonicalFlagError("deviation.record"))
	installInputDescriptor(cmd, inputDescriptorRegistry["deviation record"])
	return cmd
}

func newDeviationListCommand() *cobra.Command {
	var subject, disposition string
	var limit int
	cmd := &cobra.Command{Use: "list", Short: "List Deviations.", Args: canonicalNoArgs("deviation.list"), RunE: func(cmd *cobra.Command, _ []string) error {
		workspace, err := issueMutationWorkspaceRoot(cmd)
		if err != nil {
			return renderDeviationError(cmd, "deviation.list", err)
		}
		items, err := review.ListDeviations(review.DeviationListOptions{WorkspaceRoot: workspace, MigrationsDir: issueMigrationsDir(), Subject: subject, Disposition: disposition, Limit: limit})
		if err != nil {
			return renderDeviationError(cmd, "deviation.list", err)
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
			return cli.WriteSuccess(cmd.OutOrStdout(), "deviation.list", map[string]any{"items": values}, nil)
		}
		return writeP23Human(cmd, "deviation.list", items...)
	}}
	cmd.Flags().StringVar(&subject, "subject", "", "Filter by type:id subject")
	cmd.Flags().StringVar(&disposition, "disposition", "", "Filter by disposition")
	cmd.Flags().IntVar(&limit, "limit", 100, "Maximum Deviations to return")
	cmd.SetFlagErrorFunc(canonicalFlagError("deviation.list"))
	return cmd
}

func renderDeviationError(cmd *cobra.Command, operation string, cause error) error {
	jsonMode, err := readJSONMode(cmd)
	if err != nil {
		return err
	}
	if rendered, ok := renderJudgmentObservationError(cmd, jsonMode, operation, cause); ok {
		return rendered
	}
	var missing *review.DeviationNotFoundError
	if errors.As(cause, &missing) {
		return writeClaimFailure(cmd, jsonMode, operation, cause, cli.CodeDeviationNotFound, "Deviation was not found", &cli.Subject{Type: cli.SubjectDeviation, ID: missing.ID}, map[string]any{"id": missing.ID})
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
		return writeClaimFailure(cmd, jsonMode, operation, cause, cli.CodeResultLimitExceeded, "Deviation result limit exceeded", nil, map[string]any{"collection": "items", "limit": resultLimit.Limit, "observed": resultLimit.Observed})
	}
	return writeGenericFailure(cmd, jsonMode, operation, cause, cli.CodeInvalidArgument, "invalid Deviation request", nil, map[string]any{"argument": "input", "value": nil, "allowed": []string{"accepted " + operation + " input"}})
}
