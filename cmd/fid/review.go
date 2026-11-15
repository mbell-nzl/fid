package main

import (
	"errors"
	"fmt"

	"github.com/mbell-nzl/fid/internal/cli"
	"github.com/mbell-nzl/fid/internal/display"
	"github.com/mbell-nzl/fid/internal/journal"
	"github.com/mbell-nzl/fid/internal/policy"
	"github.com/mbell-nzl/fid/internal/project"
	"github.com/mbell-nzl/fid/internal/reviewrecord"
	"github.com/mbell-nzl/fid/internal/reviewwork"
	"github.com/spf13/cobra"
)

func init() {
	for _, leaf := range []string{"review record", "review list"} {
		rootFlagCapabilities[leaf] = rootFlagCapability{Project: true, JSON: true}
	}
}

func newReviewCommand() *cobra.Command {
	cmd := &cobra.Command{Use: "review", Short: "Record and read attributed reviews."}
	cmd.AddCommand(newReviewRecordCommand(), newReviewListCommand())
	return cmd
}

func newReviewRecordCommand() *cobra.Command {
	var input, actor, key string
	cmd := &cobra.Command{Use: "record <task>", Short: "Record a review session and its findings.", Args: canonicalExactArgs("review.record", 1), RunE: func(cmd *cobra.Command, args []string) error {
		for _, flag := range []string{"input", "actor", "idempotency-key"} {
			if err := requireCanonicalFlag(cmd, "review.record", flag); err != nil {
				return err
			}
		}
		if !policy.Allows(actor, policy.WorkSet) {
			return reviewFailure(cmd, "review.record", fmt.Errorf("service actors cannot record reviews"), actor)
		}
		var in reviewrecord.Input
		if err := readDescriptorInput(cmd.InOrStdin(), input, inputDescriptorRegistry["review record"], &in); err != nil {
			return reviewFailure(cmd, "review.record", err, "")
		}
		root, err := rebuildWorkspaceRoot(cmd)
		if err != nil {
			return reviewFailure(cmd, "review.record", err, "")
		}
		result, err := reviewwork.Record(reviewwork.Options{WorkspaceRoot: root, MigrationsDir: issueMigrationsDir(), TaskID: args[0], Actor: actor, IdempotencyKey: key, Input: in, Clock: commandDisplayEnvironment(cmd).now})
		if err != nil {
			return reviewFailure(cmd, "review.record", err, "")
		}
		jsonMode, err := readJSONMode(cmd)
		if err != nil {
			return err
		}
		if jsonMode {
			return cli.WriteSuccess(cmd.OutOrStdout(), "review.record", map[string]any{"round": result.Round}, &cli.Event{ID: result.EventID, OccurredAt: result.At})
		}
		_, err = fmt.Fprintf(cmd.OutOrStdout(), "✓ Recorded review %s · %s · %s\n", cli.EscapeControlText(result.Round.ID), result.Round.Stream, result.Round.Verdict)
		return err
	}}
	cmd.Flags().StringVar(&input, "input", "", "Review JSON input")
	cmd.Flags().StringVar(&actor, "actor", "", "Who records the attributed review")
	cmd.Flags().StringVar(&key, "idempotency-key", "", "Retry key")
	installInputDescriptor(cmd, inputDescriptorRegistry["review record"])
	return cmd
}

func newReviewListCommand() *cobra.Command {
	return &cobra.Command{Use: "list <task>", Short: "Read recorded reviews and their findings.", RunE: func(cmd *cobra.Command, args []string) error {
		root, err := rebuildWorkspaceRoot(cmd)
		if err != nil {
			return reviewFailure(cmd, "review.list", err, "")
		}
		if _, err := journal.EnsureProjectionFresh(root, issueMigrationsDir()); err != nil {
			return reviewFailure(cmd, "review.list", err, "")
		}
		views, err := reviewwork.List(root, args[0])
		if err != nil {
			return reviewFailure(cmd, "review.list", err, "")
		}
		jsonMode, err := readJSONMode(cmd)
		if err != nil {
			return err
		}
		if jsonMode {
			return cli.WriteSuccess(cmd.OutOrStdout(), "review.list", map[string]any{"task_id": args[0], "rounds": views}, nil)
		}
		environment := commandDisplayEnvironment(cmd)
		renderer := display.NewRenderer(display.IsTerminal(cmd.OutOrStdout()), environmentValue(environment.environ, "NO_COLOR"))
		return reviewwork.WriteText(cmd.OutOrStdout(), args[0], views, renderer, environment.now())
	}}
}

func reviewFailure(cmd *cobra.Command, operation string, cause error, actor string) error {
	if refused := renderRepositoryRefusal(cmd, cause); refused != nil {
		return refused
	}
	jsonMode, err := readJSONMode(cmd)
	if err != nil {
		return err
	}
	var common *cli.Failure
	if errors.Is(cause, project.ErrNotFound) {
		failure, err := issueMutationFailure(cmd, operation, "", cause)
		if err != nil {
			return err
		}
		common = &failure
	}
	for _, classify := range []func(error) (cli.Failure, bool, error){lockUnavailableFailure, journalReplayStoppedFailure} {
		failure, handled, err := classify(cause)
		if err != nil {
			return err
		}
		if handled {
			common = &failure
			break
		}
	}
	if common != nil {
		if jsonMode {
			return cli.WriteFailureAndReturn(cmd.ErrOrStderr(), operation, *common, cause)
		}
		return writeHumanFailure(cmd, operation, *common, cause)
	}
	if actor != "" {
		return writeFailure(cmd, jsonMode, operation, cause, cli.CodeAuthorityDenied, cause.Error(), nil, map[string]any{"actor_id": actor, "operation": operation, "required_actions": []string{"review"}, "required_capabilities": []string{"review"}})
	}
	var conflict *journal.IdempotencyConflictError
	if errors.As(cause, &conflict) {
		return writeFailure(cmd, jsonMode, operation, cause, cli.CodeIdempotencyConflict, "idempotency key is bound to another request", nil, idempotencyConflictDetails(conflict))
	}
	var input *cli.InputError
	if errors.As(cause, &input) {
		return writeFailure(cmd, jsonMode, operation, cause, cli.CodeInvalidArgument, input.Error(), nil, input.Details())
	}
	return writeFailure(cmd, jsonMode, operation, cause, cli.CodeInvalidArgument, cli.BoundedMessage(cause.Error()), nil, map[string]any{"argument": "input", "value": nil, "allowed": []string{"an independent review of the latest hand-off"}})
}
