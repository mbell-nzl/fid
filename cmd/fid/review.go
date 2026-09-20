package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/mbell-nzl/fid/internal/cli"
	"github.com/mbell-nzl/fid/internal/review"
	"github.com/spf13/cobra"
)

func init() { rootFlagCapabilities["review record"] = rootFlagCapability{Project: true, JSON: true} }

func newReviewCommand() *cobra.Command {
	cmd := &cobra.Command{Use: "review", Short: "Record immutable Reviews."}
	cmd.AddCommand(newReviewRecordCommand())
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
		if err := review.Record(review.RecordOptions{WorkspaceRoot: workspace, MigrationsDir: issueMigrationsDir(), ActorID: actor}, request); err != nil {
			return renderReviewError(cmd, err)
		}
		return fmt.Errorf("review record implementation is incomplete")
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

func renderReviewError(cmd *cobra.Command, cause error) error {
	jsonMode, err := readJSONMode(cmd)
	if err != nil {
		return err
	}
	var independent *review.NotIndependentError
	if errors.As(cause, &independent) {
		return writeClaimFailure(cmd, jsonMode, "review.record", cause, cli.CodeReviewNotIndependent, "review actor is not independent", nil, map[string]any{"review_actor_id": independent.ReviewActorID, "conflicting_subject": map[string]any{"type": independent.ConflictingType, "id": independent.ConflictingID}, "conflicting_actor_id": independent.ConflictingActorID})
	}
	return writeClaimFailure(cmd, jsonMode, "review.record", cause, cli.CodeInvalidArgument, "invalid review request", nil, map[string]any{"argument": "input", "value": nil, "allowed": []string{"accepted review.record input"}})
}
