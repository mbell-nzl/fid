package main

import (
	"errors"
	"fmt"

	"github.com/mbell-nzl/fid/internal/cli"
	"github.com/mbell-nzl/fid/internal/journal"
	"github.com/mbell-nzl/fid/internal/policy"
	"github.com/mbell-nzl/fid/internal/review"
	"github.com/spf13/cobra"
)

func newConcernCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "concern",
		Short: "Manage concerns.",
		RunE:  func(cmd *cobra.Command, _ []string) error { return cmd.Help() },
	}

	cmd.AddCommand(newConcernDisposeCommand())
	cmd.AddCommand(newConcernListCommand(), newConcernGetCommand())
	cmd.AddCommand(newConcernRecordCommand())
	return cmd
}

func init() {
	rootFlagCapabilities["concern get"] = rootFlagCapability{Project: true, JSON: true}
	rootFlagCapabilities["concern list"] = rootFlagCapability{Project: true, JSON: true}
	rootFlagCapabilities["concern record"] = rootFlagCapability{Project: true, JSON: true}
	rootFlagCapabilities["concern dispose"] = rootFlagCapability{Project: true, JSON: true}
}

func newConcernRecordCommand() *cobra.Command {
	return newConcernMutationCommand("record", func(cmd *cobra.Command, inputPath string) (concernMutation, error) {
		var request review.ConcernRecordRequest
		if err := readDescriptorInput(cmd.InOrStdin(), inputPath, inputDescriptorRegistry[rootCommandKey(cmd)], &request); err != nil {
			return nil, err
		}
		return func(options review.JudgmentOptions) (review.ConcernMutationResult, error) {
			return review.RecordConcern(options, request)
		}, nil
	})
}

func newConcernDisposeCommand() *cobra.Command {
	return newConcernMutationCommand("dispose", func(cmd *cobra.Command, inputPath string) (concernMutation, error) {
		var request review.ConcernDisposeRequest
		if err := readDescriptorInput(cmd.InOrStdin(), inputPath, inputDescriptorRegistry[rootCommandKey(cmd)], &request); err != nil {
			return nil, err
		}
		return func(options review.JudgmentOptions) (review.ConcernMutationResult, error) {
			return review.DisposeConcern(options, request)
		}, nil
	})
}

type concernMutation func(review.JudgmentOptions) (review.ConcernMutationResult, error)
type concernDecoder func(*cobra.Command, string) (concernMutation, error)

func newConcernMutationCommand(name string, decode concernDecoder) *cobra.Command {
	operation := "concern." + name
	var inputPath, actor, key string
	cmd := &cobra.Command{Use: name, Short: "Record or dispose a Concern.", Args: canonicalNoArgs(operation), RunE: func(cmd *cobra.Command, _ []string) error {
		for _, flag := range []string{"input", "actor", "idempotency-key"} {
			if err := requireCanonicalFlag(cmd, operation, flag); err != nil {
				return err
			}
		}
		mutate, err := decode(cmd, inputPath)
		if err != nil {
			return renderConcernError(cmd, operation, err)
		}
		workspace, err := issueMutationWorkspaceRoot(cmd)
		if err != nil {
			return renderConcernError(cmd, operation, err)
		}
		result, err := mutate(review.JudgmentOptions{WorkspaceRoot: workspace, MigrationsDir: issueMigrationsDir(), ActorID: actor, IdempotencyKey: key})
		if err != nil {
			return renderConcernError(cmd, operation, err)
		}
		jsonMode, err := readJSONMode(cmd)
		if err != nil {
			return err
		}
		if jsonMode {
			occurredAt := result.Concern["created_at"].(string)
			if disposition, ok := result.Concern["disposition"].(map[string]any); ok {
				occurredAt = disposition["disposed_at"].(string)
			}
			return cli.WriteSuccess(cmd.OutOrStdout(), operation, map[string]any{"concern": result.Concern, "idempotency_key": result.IdempotencyKey, "replayed": result.Replayed}, &cli.Event{ID: result.EventID, OccurredAt: occurredAt})
		}
		return writeP23Human(cmd, operation, result.Concern)
	}}
	cmd.Flags().StringVar(&inputPath, "input", "", "Concern JSON input")
	cmd.Flags().StringVar(&actor, "actor", "", "Finalized activation actor ID")
	cmd.Flags().StringVar(&key, "idempotency-key", "", "Idempotency key")
	cmd.SetFlagErrorFunc(canonicalFlagError(operation))
	installInputDescriptor(cmd, inputDescriptorRegistry["concern "+name])
	return cmd
}

func renderConcernError(cmd *cobra.Command, operation string, cause error) error {
	jsonMode, err := readJSONMode(cmd)
	if err != nil {
		return err
	}
	if rendered, ok := renderJudgmentObservationError(cmd, jsonMode, operation, cause); ok {
		return rendered
	}
	var missing *review.ConcernNotFoundError
	if errors.As(cause, &missing) {
		return writeClaimFailure(cmd, jsonMode, operation, cause, cli.CodeConcernNotFound, "Concern was not found", &cli.Subject{Type: cli.SubjectConcern, ID: missing.ID}, map[string]any{"id": missing.ID})
	}
	var issueMissing *review.ConcernIssueNotFoundError
	if errors.As(cause, &issueMissing) {
		return writeClaimFailure(cmd, jsonMode, operation, cause, cli.CodeIssueNotFound, "Issue was not found", &cli.Subject{Type: cli.SubjectIssue, ID: issueMissing.ID}, map[string]any{"id": issueMissing.ID})
	}
	var findingMissing *review.FindingNotFoundError
	if errors.As(cause, &findingMissing) {
		return writeClaimFailure(cmd, jsonMode, operation, cause, cli.CodeFindingNotFound, "Finding was not found", &cli.Subject{Type: cli.SubjectFinding, ID: findingMissing.ID}, map[string]any{"id": findingMissing.ID})
	}
	var disposed *review.ConcernAlreadyDisposedError
	if errors.As(cause, &disposed) {
		return writeClaimFailure(cmd, jsonMode, operation, cause, cli.CodeConcernAlreadyDisposed, "Concern is already disposed", &cli.Subject{Type: cli.SubjectConcern, ID: disposed.ID}, map[string]any{"concern_id": disposed.ID, "disposed_at": disposed.DisposedAt})
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
		return writeClaimFailure(cmd, jsonMode, operation, cause, cli.CodeResultLimitExceeded, "Concern result limit exceeded", nil, map[string]any{"collection": "items", "limit": resultLimit.Limit, "observed": resultLimit.Observed})
	}
	return writeGenericFailure(cmd, jsonMode, operation, cause, cli.CodeInvalidArgument, "invalid Concern request", nil, map[string]any{"argument": "input", "value": nil, "allowed": []string{"accepted " + operation + " input"}})
}

func newConcernListCommand() *cobra.Command {
	var (
		subject, status, severity string
		limit                     int
	)

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List concerns.",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			workspaceRoot, err := commandWorkingDirectory(cmd)
			if err != nil {
				return fmt.Errorf("resolve workspace root: %w", err)
			}

			items, err := review.ListConcerns(review.ConcernListOptions{WorkspaceRoot: workspaceRoot, MigrationsDir: issueMigrationsDir(), Subject: subject, Status: status, Severity: severity, Limit: limit})
			if err != nil {
				return renderConcernError(cmd, "concern.list", err)
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
				return cli.WriteSuccess(cmd.OutOrStdout(), "concern.list", map[string]any{"items": values}, nil)
			}
			return writeP23Human(cmd, "concern.list", items...)
		},
	}

	cmd.Flags().StringVar(&subject, "subject", "", "Filter by type:id subject")
	cmd.Flags().StringVar(&status, "status", "", "Filter by concern status")
	cmd.Flags().StringVar(&severity, "severity", "", "Filter by severity")
	cmd.Flags().IntVar(&limit, "limit", 100, "Maximum Concerns to return")
	cmd.SetFlagErrorFunc(canonicalFlagError("concern.list"))
	return cmd
}

func newConcernGetCommand() *cobra.Command {
	cmd := &cobra.Command{Use: "get <concern-id>", Short: "Read one Concern.", Args: canonicalExactArgs("concern.get", 1), RunE: func(cmd *cobra.Command, args []string) error {
		workspace, err := issueMutationWorkspaceRoot(cmd)
		if err != nil {
			return renderConcernError(cmd, "concern.get", err)
		}
		value, err := review.GetConcern(workspace, issueMigrationsDir(), args[0])
		if err != nil {
			return renderConcernError(cmd, "concern.get", err)
		}
		jsonMode, err := readJSONMode(cmd)
		if err != nil {
			return err
		}
		if jsonMode {
			return cli.WriteSuccess(cmd.OutOrStdout(), "concern.get", map[string]any{"concern": value}, nil)
		}
		return writeP23Human(cmd, "concern.get", value)
	}}
	cmd.SetFlagErrorFunc(canonicalFlagError("concern.get"))
	return cmd
}
