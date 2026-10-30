package main

import (
	"errors"

	"github.com/mbell-nzl/fid/internal/cli"
	"github.com/mbell-nzl/fid/internal/journal"
	"github.com/mbell-nzl/fid/internal/review"
	"github.com/spf13/cobra"
)

func init() {
	rootFlagCapabilities["review record"] = rootFlagCapability{Project: true, JSON: true}
	rootFlagCapabilities["review get"] = rootFlagCapability{Project: true, JSON: true}
	rootFlagCapabilities["review list"] = rootFlagCapability{Project: true, JSON: true}
}

func newReviewCommand() *cobra.Command {
	cmd := &cobra.Command{Use: "review", Short: "Record immutable Reviews."}
	cmd.AddCommand(newReviewRecordCommand(), newReviewGetCommand(), newReviewListCommand())
	return cmd
}

func newReviewListCommand() *cobra.Command {
	var issueID, envelopeID, verdict string
	var limit int
	cmd := &cobra.Command{Use: "list", Short: "List Reviews.", Args: canonicalNoArgs("review.list"), RunE: func(cmd *cobra.Command, _ []string) error {
		workspace, err := issueMutationWorkspaceRoot(cmd)
		if err != nil {
			return renderReviewListError(cmd, err)
		}
		items, err := review.ListReviews(review.ReviewListOptions{WorkspaceRoot: workspace, MigrationsDir: issueMigrationsDir(), IssueID: issueID, EnvelopeID: envelopeID, Verdict: verdict, Limit: limit})
		if err != nil {
			return renderReviewListError(cmd, err)
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
			return cli.WriteSuccess(cmd.OutOrStdout(), "review.list", map[string]any{"items": values}, nil)
		}
		return writeP23Human(cmd, "review.list", items...)
	}}
	cmd.Flags().StringVar(&issueID, "issue", "", "Filter by Issue")
	cmd.Flags().StringVar(&envelopeID, "envelope", "", "Filter by Envelope")
	cmd.Flags().StringVar(&verdict, "verdict", "", "Filter by verdict")
	cmd.Flags().IntVar(&limit, "limit", 100, "Maximum Reviews to return")
	cmd.SetFlagErrorFunc(canonicalFlagError("review.list"))
	return cmd
}

func renderReviewListError(cmd *cobra.Command, cause error) error {
	jsonMode, err := readJSONMode(cmd)
	if err != nil {
		return err
	}
	var limit *review.ResultLimitExceededError
	if errors.As(cause, &limit) {
		return writeClaimFailure(cmd, jsonMode, "review.list", cause, cli.CodeResultLimitExceeded, "review result limit exceeded", nil, map[string]any{"collection": "items", "limit": limit.Limit, "observed": limit.Observed})
	}
	return writeGenericFailure(cmd, jsonMode, "review.list", cause, cli.CodeInvalidArgument, "invalid review list request", nil, map[string]any{"argument": "filter", "value": nil, "allowed": []string{"accepted review.list filters"}})
}

func newReviewGetCommand() *cobra.Command {
	cmd := &cobra.Command{Use: "get <review-id>", Short: "Read one Review.", Args: canonicalExactArgs("review.get", 1), RunE: func(cmd *cobra.Command, args []string) error {
		workspace, err := issueMutationWorkspaceRoot(cmd)
		if err != nil {
			return renderReviewError(cmd, err)
		}
		value, err := review.Get(workspace, issueMigrationsDir(), args[0])
		if err != nil {
			return renderReviewError(cmd, err)
		}
		jsonMode, err := readJSONMode(cmd)
		if err != nil {
			return err
		}
		if jsonMode {
			return cli.WriteSuccess(cmd.OutOrStdout(), "review.get", map[string]any{"review": value}, nil)
		}
		return writeP23Human(cmd, "review.get", value)
	}}
	cmd.SetFlagErrorFunc(canonicalFlagError("review.get"))
	return cmd
}

func newReviewRecordCommand() *cobra.Command {
	var inputPath, actor, key string
	cmd := &cobra.Command{Use: "record", Short: "Record a Review.", Args: canonicalNoArgs("review.record"), RunE: func(cmd *cobra.Command, _ []string) error {
		for _, flag := range []string{"input", "actor", "idempotency-key"} {
			if err := requireCanonicalFlag(cmd, "review.record", flag); err != nil {
				return err
			}
		}
		var request review.RecordRequest
		if err := readDescriptorInput(cmd.InOrStdin(), inputPath, inputDescriptorRegistry["review record"], &request); err != nil {
			return renderReviewError(cmd, err)
		}
		workspace, err := issueMutationWorkspaceRoot(cmd)
		if err != nil {
			return renderReviewError(cmd, err)
		}
		result, err := review.Record(review.RecordOptions{WorkspaceRoot: workspace, MigrationsDir: issueMigrationsDir(), ActorID: actor, IdempotencyKey: key}, request)
		if err != nil {
			return renderReviewError(cmd, err)
		}
		jsonMode, err := readJSONMode(cmd)
		if err != nil {
			return err
		}
		if jsonMode {
			return cli.WriteSuccess(cmd.OutOrStdout(), "review.record", map[string]any{"review": result.Review, "budget": result.Budget, "stops": result.Stops, "attention": result.Attention, "idempotency_key": result.IdempotencyKey, "replayed": result.Replayed}, &cli.Event{ID: result.EventID, OccurredAt: result.Review["recorded_at"].(string)})
		}
		return writeP23MutationHuman(cmd, "review.record", result.Review, result.Budget, nil, result.Stops, result.Attention)
	}}
	cmd.Flags().StringVar(&inputPath, "input", "", "Review JSON input")
	cmd.Flags().StringVar(&actor, "actor", "", "Finalized activation actor ID")
	cmd.Flags().StringVar(&key, "idempotency-key", "", "Idempotency key")
	cmd.SetFlagErrorFunc(canonicalFlagError("review.record"))
	installInputDescriptor(cmd, inputDescriptorRegistry["review record"])
	return cmd
}

func renderReviewError(cmd *cobra.Command, cause error) error {
	jsonMode, err := readJSONMode(cmd)
	if err != nil {
		return err
	}
	var independent *review.NotIndependentError
	if errors.As(cause, &independent) {
		return writeClaimFailure(cmd, jsonMode, "review.record", cause, cli.CodeReviewNotIndependent, "review actor is not independent", nil, map[string]any{"review_actor_id": independent.ReviewActorID, "conflicting_subject": map[string]any{"type": independent.ConflictingType, "id": independent.ConflictingID}, "conflicting_actor_id": independent.ConflictingActorID})
	}
	var verdict *review.VerdictInvalidError
	if errors.As(cause, &verdict) {
		return writeClaimFailure(cmd, jsonMode, "review.record", cause, cli.CodeReviewVerdictInvalid, "review verdict is invalid", nil, map[string]any{"verdict": verdict.Verdict, "reason": verdict.Reason})
	}
	var mismatch *review.RepositoryMismatchError
	if errors.As(cause, &mismatch) {
		return writeClaimFailure(cmd, jsonMode, "review.record", cause, cli.CodeRepositoryMismatch, "repository identity does not match", nil, map[string]any{"expected": map[string]any{"head_oid": mismatch.Expected.HeadOID, "dirty": mismatch.Expected.Dirty}, "observed": repositoryDetails(mismatch.Observed)})
	}
	if rendered, ok := renderJudgmentObservationError(cmd, jsonMode, "review.record", cause); ok {
		return rendered
	}
	var reviewConflict *review.ReviewConflictError
	if errors.As(cause, &reviewConflict) {
		return writeClaimFailure(cmd, jsonMode, "review.record", cause, cli.CodeDuplicateEntity, "Review already exists for Run", &cli.Subject{Type: cli.SubjectRun, ID: reviewConflict.RunID}, map[string]any{"entity_type": "review", "identity": map[string]any{"fields": []string{"run_id"}, "values": []string{reviewConflict.RunID}}})
	}
	var findingMissing *review.FindingNotFoundError
	if errors.As(cause, &findingMissing) {
		return writeClaimFailure(cmd, jsonMode, "review.record", cause, cli.CodeFindingNotFound, "Finding was not found", &cli.Subject{Type: cli.SubjectFinding, ID: findingMissing.ID}, map[string]any{"id": findingMissing.ID})
	}
	var disposed *review.FindingAlreadyDisposedError
	if errors.As(cause, &disposed) {
		return writeClaimFailure(cmd, jsonMode, "review.record", cause, cli.CodeFindingAlreadyDisposed, "Finding is already terminal", &cli.Subject{Type: cli.SubjectFinding, ID: disposed.ID}, map[string]any{"finding_id": disposed.ID, "state": disposed.State})
	}
	var invalidTarget *review.FindingTargetInvalidError
	if errors.As(cause, &invalidTarget) {
		return writeClaimFailure(cmd, jsonMode, "review.record", cause, cli.CodeInvalidArgument, "invalid Finding target", nil, map[string]any{"argument": "finding_id", "value": invalidTarget.ID, "allowed": []string{"earlier Finding on the same Issue and Envelope revision"}})
	}
	var missing *review.ReviewNotFoundError
	if errors.As(cause, &missing) {
		return writeClaimFailure(cmd, jsonMode, "review.get", cause, cli.CodeReviewNotFound, "review was not found", &cli.Subject{Type: cli.SubjectReview, ID: missing.ID}, map[string]any{"id": missing.ID})
	}
	var conflict *journal.IdempotencyConflictError
	if errors.As(cause, &conflict) {
		return writeClaimFailure(cmd, jsonMode, "review.record", cause, cli.CodeIdempotencyConflict, "idempotency key conflicts with a different request", nil, idempotencyConflictDetails(conflict))
	}
	return writeGenericFailure(cmd, jsonMode, "review.record", cause, cli.CodeInvalidArgument, "invalid review request", nil, map[string]any{"argument": "input", "value": nil, "allowed": []string{"accepted review.record input"}})
}

func repositoryDetails(repository review.Repository) map[string]any {
	var baseOID, dirtyDigest any
	if repository.BaseOID != nil {
		baseOID = *repository.BaseOID
	}
	if repository.DirtyDigest != nil {
		dirtyDigest = *repository.DirtyDigest
	}
	return map[string]any{"base_oid": baseOID, "capture_source": repository.CaptureSource, "dirty": repository.Dirty, "dirty_digest": dirtyDigest, "head_oid": repository.HeadOID, "project_id": repository.ProjectID, "tree_oid": repository.TreeOID, "worktree_id": repository.WorktreeID}
}
