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
	"github.com/mbell-nzl/fid/internal/policy"
	"github.com/mbell-nzl/fid/internal/review"
	"github.com/spf13/cobra"
)

func newConcernCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "concern",
		Short: "Manage concerns.",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				return fmt.Errorf("unknown command %q for %q", args[0], cmd.CommandPath())
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, _ []string) error { return cmd.Help() },
	}

	cmd.AddCommand(newConcernDisposeCommand())
	cmd.AddCommand(newConcernListCommand())
	cmd.AddCommand(newConcernRecordCommand())
	return cmd
}

func init() {
	rootFlagCapabilities["concern list"] = rootFlagCapability{Project: true, JSON: true}
	rootFlagCapabilities["concern record"] = rootFlagCapability{Project: true, JSON: true}
	rootFlagCapabilities["concern dispose"] = rootFlagCapability{Project: true, JSON: true}
}

func newConcernRecordCommand() *cobra.Command {
	return newConcernMutationCommand("record", func(options review.JudgmentOptions, inputPath string, cmd *cobra.Command) (review.ConcernMutationResult, error) {
		var request review.ConcernRecordRequest
		fields := []string{"id", "subject", "scope", "severity", "summary", "reasoning", "references"}
		if err := decodeClosedJudgmentInput(cmd.InOrStdin(), inputPath, fields, &request); err != nil {
			return review.ConcernMutationResult{}, err
		}
		return review.RecordConcern(options, request)
	})
}

func newConcernDisposeCommand() *cobra.Command {
	return newConcernMutationCommand("dispose", func(options review.JudgmentOptions, inputPath string, cmd *cobra.Command) (review.ConcernMutationResult, error) {
		var request review.ConcernDisposeRequest
		fields := []string{"concern_id", "action", "issue_id", "supersedes_concern_id", "supersedes_finding_id", "reason"}
		if err := decodeClosedJudgmentInput(cmd.InOrStdin(), inputPath, fields, &request); err != nil {
			return review.ConcernMutationResult{}, err
		}
		return review.DisposeConcern(options, request)
	})
}

type concernMutation func(review.JudgmentOptions, string, *cobra.Command) (review.ConcernMutationResult, error)

func newConcernMutationCommand(name string, mutate concernMutation) *cobra.Command {
	operation := "concern." + name
	var inputPath, actor, key string
	cmd := &cobra.Command{Use: name, Args: canonicalNoArgs(operation), RunE: func(cmd *cobra.Command, _ []string) error {
		for _, flag := range []string{"input", "actor", "idempotency-key"} {
			if err := requireCanonicalFlag(cmd, operation, flag); err != nil {
				return err
			}
		}
		workspace, err := issueMutationWorkspaceRoot(cmd)
		if err != nil {
			return renderConcernError(cmd, operation, err)
		}
		result, err := mutate(review.JudgmentOptions{WorkspaceRoot: workspace, MigrationsDir: issueMigrationsDir(), ActorID: actor, IdempotencyKey: key}, inputPath, cmd)
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
	return cmd
}

func renderConcernError(cmd *cobra.Command, operation string, cause error) error {
	jsonMode, err := readJSONMode(cmd)
	if err != nil {
		return err
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
		return writeClaimFailure(cmd, jsonMode, operation, cause, cli.CodeIdempotencyConflict, "idempotency key conflicts with a different request", nil, map[string]any{"operation_scope": conflict.OperationScope, "idempotency_key": conflict.IdempotencyKey, "stored_request_digest": conflict.StoredRequestDigest, "supplied_request_digest": conflict.SuppliedRequestDigest})
	}
	var resultLimit *review.ResultLimitExceededError
	if errors.As(cause, &resultLimit) {
		return writeClaimFailure(cmd, jsonMode, operation, cause, cli.CodeResultLimitExceeded, "Concern result limit exceeded", nil, map[string]any{"collection": "items", "limit": resultLimit.Limit, "observed": resultLimit.Observed})
	}
	return writeClaimFailure(cmd, jsonMode, operation, cause, cli.CodeInvalidArgument, "invalid Concern request", nil, map[string]any{"argument": "input", "value": nil, "allowed": []string{"accepted " + operation + " input"}})
}

func decodeClosedJudgmentInput(stdin io.Reader, path string, required []string, target any) error {
	payload, err := os.ReadFile(path)
	if path == "-" {
		payload, err = io.ReadAll(io.LimitReader(stdin, 1_048_577))
	}
	if err != nil || !utf8.Valid(payload) {
		return fmt.Errorf("read judgment input: invalid input")
	}
	if err := rejectDuplicateJSONKeys(payload); err != nil {
		return fmt.Errorf("decode judgment input: %w", err)
	}
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(payload, &fields); err != nil {
		return fmt.Errorf("decode judgment input: %w", err)
	}
	if len(fields) != len(required) {
		return fmt.Errorf("decode judgment input: wrong field set")
	}
	for _, name := range required {
		if _, found := fields[name]; !found {
			return fmt.Errorf("decode judgment input: missing %s", name)
		}
	}
	decoder := json.NewDecoder(strings.NewReader(string(payload)))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(target); err != nil {
		return fmt.Errorf("decode judgment input: %w", err)
	}
	var extra any
	if err := decoder.Decode(&extra); err != io.EOF {
		return fmt.Errorf("judgment input must contain one object")
	}
	return nil
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
