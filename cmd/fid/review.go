package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"unicode/utf8"

	"github.com/mbell-nzl/fid/internal/cli"
	"github.com/mbell-nzl/fid/internal/journal"
	"github.com/mbell-nzl/fid/internal/review"
	"github.com/spf13/cobra"
)

func init() {
	rootFlagCapabilities["review record"] = rootFlagCapability{Project: true, JSON: true}
	rootFlagCapabilities["review get"] = rootFlagCapability{Project: true, JSON: true}
}

func newReviewCommand() *cobra.Command {
	cmd := &cobra.Command{Use: "review", Short: "Record immutable Reviews."}
	cmd.AddCommand(newReviewRecordCommand(), newReviewGetCommand())
	return cmd
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
		_, err = fmt.Fprintf(cmd.OutOrStdout(), "%s\t%s\n", value["id"], value["verdict"])
		return err
	}}
	cmd.SetFlagErrorFunc(canonicalFlagError("review.get"))
	return cmd
}

func newReviewRecordCommand() *cobra.Command {
	var inputPath, actor, key string
	cmd := &cobra.Command{Use: "record", Args: canonicalNoArgs("review.record"), RunE: func(cmd *cobra.Command, _ []string) error {
		for _, flag := range []string{"input", "actor", "idempotency-key"} {
			if err := requireCanonicalFlag(cmd, "review.record", flag); err != nil {
				return err
			}
		}
		var request review.RecordRequest
		if err := decodeClosedReviewInput(cmd.InOrStdin(), inputPath, &request); err != nil {
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
		_, err = fmt.Fprintf(cmd.OutOrStdout(), "recorded review %s\n", result.Review["id"])
		return err
	}}
	cmd.Flags().StringVar(&inputPath, "input", "", "Review JSON input")
	cmd.Flags().StringVar(&actor, "actor", "", "Finalized activation actor ID")
	cmd.Flags().StringVar(&key, "idempotency-key", "", "Idempotency key")
	cmd.SetFlagErrorFunc(canonicalFlagError("review.record"))
	return cmd
}

func decodeClosedReviewInput(stdin io.Reader, path string, target any) error {
	payload, err := os.ReadFile(path)
	if path == "-" {
		payload, err = io.ReadAll(io.LimitReader(stdin, 1_048_577))
	}
	if err != nil {
		return fmt.Errorf("read review input: %w", err)
	}
	if !utf8.Valid(payload) {
		return fmt.Errorf("decode review input: input is not valid UTF-8")
	}
	if err := rejectDuplicateJSONKeys(payload); err != nil {
		return fmt.Errorf("decode review input: %w", err)
	}
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(payload, &fields); err == nil {
		if raw, found := fields["inapplicable_findings"]; found && string(raw) == "null" {
			return fmt.Errorf("decode review input: inapplicable_findings must be an array")
		}
	}
	decoder := json.NewDecoder(strings.NewReader(string(payload)))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(target); err != nil {
		return fmt.Errorf("decode review input: %w", err)
	}
	var extra any
	if err := decoder.Decode(&extra); err != io.EOF {
		return fmt.Errorf("review input must contain one object")
	}
	return nil
}

func rejectDuplicateJSONKeys(payload []byte) error {
	decoder := json.NewDecoder(strings.NewReader(string(payload)))
	var walk func() error
	walk = func() error {
		token, err := decoder.Token()
		if err != nil {
			return err
		}
		delimiter, ok := token.(json.Delim)
		if !ok {
			return nil
		}
		switch delimiter {
		case '{':
			seen := map[string]bool{}
			for decoder.More() {
				keyToken, err := decoder.Token()
				if err != nil {
					return err
				}
				key := keyToken.(string)
				if seen[key] {
					return fmt.Errorf("duplicate key %q", key)
				}
				seen[key] = true
				if err := walk(); err != nil {
					return err
				}
			}
			_, err = decoder.Token()
			return err
		case '[':
			for decoder.More() {
				if err := walk(); err != nil {
					return err
				}
			}
			_, err = decoder.Token()
			return err
		default:
			return fmt.Errorf("unexpected JSON delimiter %q", delimiter)
		}
	}
	return walk()
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
		return writeClaimFailure(cmd, jsonMode, "review.record", cause, cli.CodeRepositoryMismatch, "repository identity does not match", nil, map[string]any{"expected": repositoryDetails(mismatch.Expected), "observed": repositoryDetails(mismatch.Observed)})
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
		return writeClaimFailure(cmd, jsonMode, "review.record", cause, cli.CodeIdempotencyConflict, "idempotency key conflicts with a different request", nil, map[string]any{"operation_scope": conflict.OperationScope, "idempotency_key": conflict.IdempotencyKey, "stored_request_digest": conflict.StoredRequestDigest, "supplied_request_digest": conflict.SuppliedRequestDigest})
	}
	return writeClaimFailure(cmd, jsonMode, "review.record", cause, cli.CodeInvalidArgument, "invalid review request", nil, map[string]any{"argument": "input", "value": nil, "allowed": []string{"accepted review.record input"}})
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
