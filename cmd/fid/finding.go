package main

import (
	"errors"

	"github.com/mbell-nzl/fid/internal/claims"
	"github.com/mbell-nzl/fid/internal/cli"
	"github.com/mbell-nzl/fid/internal/journal"
	"github.com/mbell-nzl/fid/internal/policy"
	"github.com/mbell-nzl/fid/internal/review"
	"github.com/mbell-nzl/fid/internal/workflow"
	"github.com/spf13/cobra"
)

func init() {
	rootFlagCapabilities["finding address"] = rootFlagCapability{Project: true, JSON: true}
	rootFlagCapabilities["finding list"] = rootFlagCapability{Project: true, JSON: true}
}

func newFindingCommand() *cobra.Command {
	cmd := &cobra.Command{Use: "finding", Short: "Address and list Review Findings."}
	cmd.AddCommand(newFindingAddressCommand(), newFindingListCommand())
	return cmd
}

func newFindingAddressCommand() *cobra.Command {
	var inputPath, actor, key string
	cmd := &cobra.Command{Use: "address", Short: "Address a Review Finding.", Args: canonicalNoArgs("finding.address"), RunE: func(cmd *cobra.Command, _ []string) error {
		for _, flag := range []string{"input", "actor", "idempotency-key"} {
			if err := requireCanonicalFlag(cmd, "finding.address", flag); err != nil {
				return err
			}
		}
		var request review.AddressRequest
		if err := decodeClosedReviewInput(cmd.InOrStdin(), inputPath, &request); err != nil {
			return renderFindingError(cmd, "finding.address", err)
		}
		workspace, err := issueMutationWorkspaceRoot(cmd)
		if err != nil {
			return renderFindingError(cmd, "finding.address", err)
		}
		result, err := review.Address(review.AddressOptions{WorkspaceRoot: workspace, MigrationsDir: issueMigrationsDir(), ActorID: actor, IdempotencyKey: key}, request)
		if err != nil {
			return renderFindingError(cmd, "finding.address", err)
		}
		jsonMode, err := readJSONMode(cmd)
		if err != nil {
			return err
		}
		if jsonMode {
			dispositions := result.Finding["dispositions"].([]any)
			occurredAt := dispositions[len(dispositions)-1].(map[string]any)["disposed_at"].(string)
			return cli.WriteSuccess(cmd.OutOrStdout(), "finding.address", map[string]any{"finding": result.Finding, "idempotency_key": result.IdempotencyKey, "replayed": result.Replayed}, &cli.Event{ID: result.EventID, OccurredAt: occurredAt})
		}
		return writeP23Human(cmd, "finding.address", result.Finding)
	}}
	cmd.Flags().StringVar(&inputPath, "input", "", "Finding JSON input")
	cmd.Flags().StringVar(&actor, "actor", "", "Finalized activation actor ID")
	cmd.Flags().StringVar(&key, "idempotency-key", "", "Idempotency key")
	cmd.SetFlagErrorFunc(canonicalFlagError("finding.address"))
	return cmd
}

func newFindingListCommand() *cobra.Command {
	var issueID, reviewID, state, severity string
	var limit int
	cmd := &cobra.Command{Use: "list", Short: "List Review Findings.", Args: canonicalNoArgs("finding.list"), RunE: func(cmd *cobra.Command, _ []string) error {
		workspace, err := issueMutationWorkspaceRoot(cmd)
		if err != nil {
			return renderFindingError(cmd, "finding.list", err)
		}
		items, err := review.List(review.ListOptions{WorkspaceRoot: workspace, MigrationsDir: issueMigrationsDir(), IssueID: issueID, ReviewID: reviewID, State: state, Severity: severity, Limit: limit})
		if err != nil {
			return renderFindingError(cmd, "finding.list", err)
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
			return cli.WriteSuccess(cmd.OutOrStdout(), "finding.list", map[string]any{"items": values}, nil)
		}
		return writeP23Human(cmd, "finding.list", items...)
	}}
	cmd.Flags().StringVar(&issueID, "issue", "", "Filter by Issue ID")
	cmd.Flags().StringVar(&reviewID, "review", "", "Filter by Review ID")
	cmd.Flags().StringVar(&state, "state", "", "Filter by derived Finding state")
	cmd.Flags().StringVar(&severity, "severity", "", "Filter by severity")
	cmd.Flags().IntVar(&limit, "limit", 100, "Maximum Findings to return")
	cmd.SetFlagErrorFunc(canonicalFlagError("finding.list"))
	return cmd
}

func renderFindingError(cmd *cobra.Command, operation string, cause error) error {
	jsonMode, err := readJSONMode(cmd)
	if err != nil {
		return err
	}
	var missing *review.FindingNotFoundError
	if errors.As(cause, &missing) {
		return writeClaimFailure(cmd, jsonMode, operation, cause, cli.CodeFindingNotFound, "Finding was not found", &cli.Subject{Type: cli.SubjectFinding, ID: missing.ID}, map[string]any{"id": missing.ID})
	}
	var disposed *review.FindingAlreadyDisposedError
	if errors.As(cause, &disposed) {
		return writeClaimFailure(cmd, jsonMode, operation, cause, cli.CodeFindingAlreadyDisposed, "Finding is already terminal", &cli.Subject{Type: cli.SubjectFinding, ID: disposed.ID}, map[string]any{"finding_id": disposed.ID, "state": disposed.State})
	}
	var invalidMode *workflow.ClaimModeInvalidError
	if errors.As(cause, &invalidMode) {
		return writeClaimFailure(cmd, jsonMode, operation, cause, cli.CodeClaimModeInvalid, "claim mode is invalid", &cli.Subject{Type: cli.SubjectClaim, ID: invalidMode.ClaimID}, map[string]any{"claim_id": invalidMode.ClaimID, "mode": invalidMode.Mode, "required_mode": invalidMode.RequiredMode})
	}
	var claimMissing *claims.ClaimNotFoundError
	if errors.As(cause, &claimMissing) {
		return writeClaimFailure(cmd, jsonMode, operation, cause, cli.CodeClaimNotFound, "claim was not found", &cli.Subject{Type: cli.SubjectClaim, ID: claimMissing.ID}, map[string]any{"id": claimMissing.ID})
	}
	var stale *claims.ClaimStaleError
	if errors.As(cause, &stale) {
		value := (*claims.ClaimConflictError)(stale)
		return writeClaimConflict(cmd, jsonMode, operation, cause, cli.CodeClaimStale, value)
	}
	var denied *policy.EnvelopeAuthorityDeniedError
	if errors.As(cause, &denied) {
		return writeClaimFailure(cmd, jsonMode, operation, cause, cli.CodeAuthorityDenied, "active envelope authority is denied", nil, map[string]any{"operation": operation, "actor_id": denied.ActorID, "required_actions": []string{denied.Action}, "required_capabilities": []string{denied.Capability}})
	}
	var conflict *journal.IdempotencyConflictError
	if errors.As(cause, &conflict) {
		return writeClaimFailure(cmd, jsonMode, operation, cause, cli.CodeIdempotencyConflict, "idempotency key conflicts with a different request", nil, idempotencyConflictDetails(conflict))
	}
	var limit *review.ResultLimitExceededError
	if errors.As(cause, &limit) {
		return writeClaimFailure(cmd, jsonMode, operation, cause, cli.CodeResultLimitExceeded, "Finding result limit exceeded", nil, map[string]any{"collection": "items", "limit": limit.Limit, "observed": limit.Observed})
	}
	var target *review.FindingTargetInvalidError
	if errors.As(cause, &target) {
		return writeClaimFailure(cmd, jsonMode, operation, cause, cli.CodeInvalidArgument, "invalid Finding target", nil, map[string]any{"argument": "finding_id", "value": target.ID, "allowed": []string{"earlier Finding on the same Issue and Envelope revision"}})
	}
	return writeClaimFailure(cmd, jsonMode, operation, cause, cli.CodeInvalidArgument, "invalid Finding request", nil, map[string]any{"argument": "input", "value": nil, "allowed": []string{"accepted " + operation + " input"}})
}
