package main

import (
	"errors"
	"fmt"
	"strings"

	"github.com/mbell-nzl/fid/internal/claims"
	"github.com/mbell-nzl/fid/internal/cli"
	"github.com/mbell-nzl/fid/internal/journal"
	"github.com/mbell-nzl/fid/internal/policy"
	"github.com/mbell-nzl/fid/internal/workflow"
	"github.com/spf13/cobra"
)

func init() {
	for _, operation := range []string{"claim acquire", "claim expire", "claim renew", "claim release", "claim list"} {
		rootFlagCapabilities[operation] = rootFlagCapability{Project: true, JSON: true}
	}
}

func newClaimCommand() *cobra.Command {
	cmd := &cobra.Command{Use: "claim", Short: "Manage generic mutation claims."}
	cmd.AddCommand(newClaimAcquireCommand(), newClaimExpireCommand(), newClaimRenewCommand(), newClaimReleaseCommand(), newClaimListCommand())
	return cmd
}

func newClaimExpireCommand() *cobra.Command {
	var inputPath, actor, key string
	cmd := &cobra.Command{Use: "expire", Short: "Record one mutation claim expiry.", Args: canonicalNoArgs("claim.expire"), RunE: func(cmd *cobra.Command, _ []string) error {
		var request claims.ExpireRequest
		if err := claimRequest(cmd, "claim.expire", inputPath, actor, key, &request); err != nil {
			return err
		}
		workspace, err := issueMutationWorkspaceRoot(cmd)
		if err != nil {
			return renderClaimError(cmd, "claim.expire", err)
		}
		result, err := claims.ExpireGeneric(claims.MutationOptions{WorkspaceRoot: workspace, MigrationsDir: issueMigrationsDir(), ActorID: actor, IdempotencyKey: key}, request)
		if err != nil {
			return renderClaimError(cmd, "claim.expire", err)
		}
		return writeClaimSuccess(cmd, "claim.expire", result)
	}}
	claimMutationFlags(cmd, &inputPath, &actor, &key)
	return cmd
}

func newClaimAcquireCommand() *cobra.Command {
	var inputPath, actor, key string
	cmd := &cobra.Command{Use: "acquire", Short: "Acquire one Claim.", Args: canonicalNoArgs("claim.acquire"), RunE: func(cmd *cobra.Command, _ []string) error {
		var request claims.AcquireRequest
		if err := claimRequest(cmd, "claim.acquire", inputPath, actor, key, &request); err != nil {
			return err
		}
		workspace, err := issueMutationWorkspaceRoot(cmd)
		if err != nil {
			return renderClaimError(cmd, "claim.acquire", err)
		}
		result, err := claims.AcquireGeneric(claims.MutationOptions{WorkspaceRoot: workspace, MigrationsDir: issueMigrationsDir(), ActorID: actor, IdempotencyKey: key}, request)
		if err != nil {
			return renderClaimError(cmd, "claim.acquire", err)
		}
		return writeClaimSuccess(cmd, "claim.acquire", result)
	}}
	claimMutationFlags(cmd, &inputPath, &actor, &key)
	return cmd
}

func newClaimRenewCommand() *cobra.Command {
	var inputPath, actor, key string
	cmd := &cobra.Command{Use: "renew", Short: "Renew one active mutation claim.", Args: canonicalNoArgs("claim.renew"), RunE: func(cmd *cobra.Command, _ []string) error {
		var request claims.RenewRequest
		if err := claimRequest(cmd, "claim.renew", inputPath, actor, key, &request); err != nil {
			return err
		}
		workspace, err := issueMutationWorkspaceRoot(cmd)
		if err != nil {
			return renderClaimError(cmd, "claim.renew", err)
		}
		result, err := claims.RenewGeneric(claims.MutationOptions{WorkspaceRoot: workspace, MigrationsDir: issueMigrationsDir(), ActorID: actor, IdempotencyKey: key}, request)
		if err != nil {
			return renderClaimError(cmd, "claim.renew", err)
		}
		return writeClaimSuccess(cmd, "claim.renew", result)
	}}
	claimMutationFlags(cmd, &inputPath, &actor, &key)
	return cmd
}

func newClaimReleaseCommand() *cobra.Command {
	var inputPath, actor, key string
	cmd := &cobra.Command{Use: "release", Short: "Release one mutation claim.", Args: canonicalNoArgs("claim.release"), RunE: func(cmd *cobra.Command, _ []string) error {
		var request claims.ReleaseRequest
		if err := claimRequest(cmd, "claim.release", inputPath, actor, key, &request); err != nil {
			return err
		}
		workspace, err := issueMutationWorkspaceRoot(cmd)
		if err != nil {
			return renderClaimError(cmd, "claim.release", err)
		}
		result, err := claims.ReleaseGeneric(claims.MutationOptions{WorkspaceRoot: workspace, MigrationsDir: issueMigrationsDir(), ActorID: actor, IdempotencyKey: key}, request)
		if err != nil {
			return renderClaimError(cmd, "claim.release", err)
		}
		return writeClaimSuccess(cmd, "claim.release", result)
	}}
	claimMutationFlags(cmd, &inputPath, &actor, &key)
	return cmd
}

func claimMutationFlags(cmd *cobra.Command, inputPath, actor, key *string) {
	cmd.Flags().StringVar(inputPath, "input", "", "Claim JSON input")
	cmd.Flags().StringVar(actor, "actor", "", "Finalized activation actor ID")
	cmd.Flags().StringVar(key, "idempotency-key", "", "Idempotency key")
	cmd.SetFlagErrorFunc(canonicalFlagError(strings.ReplaceAll(rootCommandKey(cmd), " ", ".")))
	installInputDescriptor(cmd, inputDescriptorRegistry["claim "+cmd.Name()])
}

func claimRequest(cmd *cobra.Command, operation, inputPath, actor, key string, target any) error {
	for _, flag := range []string{"input", "actor", "idempotency-key"} {
		if err := requireCanonicalFlag(cmd, operation, flag); err != nil {
			return err
		}
	}
	if err := readDescriptorInput(cmd.InOrStdin(), inputPath, inputDescriptorRegistry[strings.ReplaceAll(operation, ".", " ")], target); err != nil {
		return renderClaimError(cmd, operation, err)
	}
	return nil
}

func newClaimListCommand() *cobra.Command {
	var issueID, envelopeID, status string
	var limit int
	cmd := &cobra.Command{Use: "list", Short: "List generic mutation claims.", Args: canonicalNoArgs("claim.list"), RunE: func(cmd *cobra.Command, _ []string) error {
		workspace, err := issueMutationWorkspaceRoot(cmd)
		if err != nil {
			return renderClaimError(cmd, "claim.list", err)
		}
		items, err := claims.ListGeneric(claims.GenericListOptions{WorkspaceRoot: workspace, MigrationsDir: issueMigrationsDir(), IssueID: issueID, EnvelopeID: envelopeID, Status: status, Limit: limit})
		if err != nil {
			return renderClaimError(cmd, "claim.list", err)
		}
		if jsonMode, err := readJSONMode(cmd); err != nil {
			return err
		} else if jsonMode {
			return cli.WriteSuccess(cmd.OutOrStdout(), "claim.list", map[string]any{"items": items}, nil)
		}
		for _, item := range items {
			if _, err := fmt.Fprintf(cmd.OutOrStdout(), "%s\t%s\n", item.ID, item.Status); err != nil {
				return err
			}
		}
		return nil
	}}
	cmd.Flags().StringVar(&issueID, "issue", "", "Issue ID")
	cmd.Flags().StringVar(&envelopeID, "envelope", "", "Envelope ID")
	cmd.Flags().StringVar(&status, "status", "", "Claim status")
	cmd.Flags().IntVar(&limit, "limit", 100, "Maximum claims")
	cmd.SetFlagErrorFunc(canonicalFlagError("claim.list"))
	return cmd
}

func writeClaimSuccess(cmd *cobra.Command, operation string, result claims.MutationResult) error {
	if jsonMode, err := readJSONMode(cmd); err != nil {
		return err
	} else if jsonMode {
		return cli.WriteSuccess(cmd.OutOrStdout(), operation, map[string]any{"claim": result.Claim, "idempotency_key": result.IdempotencyKey, "replayed": result.Replayed}, &cli.Event{ID: result.EventID, OccurredAt: result.Claim.AcquiredAt})
	}
	_, err := fmt.Fprintf(cmd.OutOrStdout(), "%s claim %s version %d\n", result.Claim.Status, result.Claim.ID, result.Claim.Version)
	return err
}

func renderClaimError(cmd *cobra.Command, operation string, cause error) error {
	var stopped *workflow.StopBlockedError
	if errors.As(cause, &stopped) {
		return renderStopError(cmd, operation, cause)
	}
	jsonMode, err := readJSONMode(cmd)
	if err != nil {
		return err
	}
	if failure, handled := renderJudgmentObservationError(cmd, jsonMode, operation, cause); handled {
		return failure
	}
	var idempotency *journal.IdempotencyConflictError
	if errors.As(cause, &idempotency) {
		return writeClaimFailure(cmd, jsonMode, operation, cause, cli.CodeIdempotencyConflict, "idempotency key conflicts with a different request", nil, idempotencyConflictDetails(idempotency))
	}
	var notFound *claims.ClaimNotFoundError
	if errors.As(cause, &notFound) {
		return writeClaimFailure(cmd, jsonMode, operation, cause, cli.CodeClaimNotFound, "claim was not found", &cli.Subject{Type: cli.SubjectClaim, ID: notFound.ID}, map[string]any{"id": notFound.ID})
	}
	var conflict *claims.ClaimConflictError
	var expired *claims.ClaimExpiredError
	var stale *claims.ClaimStaleError
	switch {
	case errors.As(cause, &conflict):
		return writeClaimConflict(cmd, jsonMode, operation, cause, cli.CodeClaimConflict, (*claims.ClaimConflictError)(conflict))
	case errors.As(cause, &expired):
		return writeClaimConflict(cmd, jsonMode, operation, cause, cli.CodeClaimExpired, (*claims.ClaimConflictError)(expired))
	case errors.As(cause, &stale):
		return writeClaimConflict(cmd, jsonMode, operation, cause, cli.CodeClaimStale, (*claims.ClaimConflictError)(stale))
	}
	var denied *policy.EnvelopeAuthorityDeniedError
	if errors.As(cause, &denied) {
		return writeClaimFailure(cmd, jsonMode, operation, cause, cli.CodeAuthorityDenied, "active envelope authority is denied", nil, map[string]any{"operation": operation, "actor_id": denied.ActorID, "required_actions": []string{denied.Action}, "required_capabilities": []string{denied.Capability}})
	}
	return writeClaimFailure(cmd, jsonMode, operation, cause, cli.CodeInvalidArgument, "invalid claim request", nil, map[string]any{"argument": "input", "value": nil, "allowed": []string{"accepted " + operation + " input"}})
}

func writeClaimConflict(cmd *cobra.Command, jsonMode bool, operation string, original error, code cli.Code, cause *claims.ClaimConflictError) error {
	return writeClaimFailure(cmd, jsonMode, operation, original, code, string(code), &cli.Subject{Type: cli.SubjectClaim, ID: claimErrorID(cause)}, map[string]any{"issue_id": cause.IssueID, "claim_id": claimErrorString(cause.ClaimID), "expected_version": claimErrorVersion(cause.ExpectedVersion), "actual_version": claimErrorVersion(cause.ActualVersion), "expires_at": claimErrorString(cause.ExpiresAt)})
}

func claimErrorID(cause *claims.ClaimConflictError) string {
	if cause.ClaimID == nil {
		return ""
	}
	return *cause.ClaimID
}

func claimErrorString(value *string) any {
	if value == nil {
		return nil
	}
	return *value
}

func claimErrorVersion(value *int) any {
	if value == nil {
		return nil
	}
	return *value
}

func writeClaimFailure(cmd *cobra.Command, jsonMode bool, operation string, cause error, code cli.Code, message string, subject *cli.Subject, details map[string]any) error {
	failure, err := cli.NewFailure(code, message, subject, details)
	if err != nil {
		return err
	}
	if !jsonMode {
		_ = renderIssueError(cmd.ErrOrStderr(), cause)
		return cli.NewFailureError(failure, cause)
	}
	return cli.WriteFailureAndReturn(cmd.ErrOrStderr(), operation, failure, fmt.Errorf("%s failed", operation))
}
