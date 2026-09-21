package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/mbell-nzl/fid/internal/cli"
	"github.com/mbell-nzl/fid/internal/workflow"
	"github.com/spf13/cobra"
)

func init() {
	rootFlagCapabilities["envelope create"] = rootFlagCapability{Project: true, JSON: true}
	rootFlagCapabilities["envelope get"] = rootFlagCapability{Project: true, JSON: true}
	rootFlagCapabilities["envelope activate"] = rootFlagCapability{Project: true, JSON: true}
	rootFlagCapabilities["envelope complete"] = rootFlagCapability{Project: true, JSON: true}
}

func newEnvelopeCommand() *cobra.Command {
	cmd := &cobra.Command{Use: "envelope", Short: "Manage the bounded work envelope."}
	cmd.AddCommand(newEnvelopeCreateCommand())
	cmd.AddCommand(newEnvelopeGetCommand())
	cmd.AddCommand(newEnvelopeActivateCommand())
	cmd.AddCommand(newEnvelopeCompleteCommand())
	return cmd
}

func newEnvelopeCreateCommand() *cobra.Command {
	var inputPath, actor, key string
	cmd := &cobra.Command{Use: "create", Short: "Create one immutable envelope revision.", Args: canonicalNoArgs("envelope.create"), RunE: func(cmd *cobra.Command, _ []string) error {
		if err := requireCanonicalFlag(cmd, "envelope.create", "input"); err != nil {
			return err
		}
		if err := requireCanonicalFlag(cmd, "envelope.create", "actor"); err != nil {
			return err
		}
		if err := requireCanonicalFlag(cmd, "envelope.create", "idempotency-key"); err != nil {
			return err
		}
		var request workflow.EnvelopeCreateRequest
		if err := decodeClosedEnvelopeInput(cmd.InOrStdin(), inputPath, &request); err != nil {
			return renderEnvelopeError(cmd, "envelope.create", err, actor)
		}
		workspaceRoot, err := issueMutationWorkspaceRoot(cmd)
		if err != nil {
			return renderEnvelopeError(cmd, "envelope.create", err, actor)
		}
		result, err := workflow.CreateEnvelope(workflow.EnvelopeOptions{WorkspaceRoot: workspaceRoot, MigrationsDir: issueMigrationsDir(), ActorID: actor, IdempotencyKey: key}, request)
		if err != nil {
			return renderEnvelopeError(cmd, "envelope.create", err, actor)
		}
		return writeEnvelopeSuccess(cmd, "envelope.create", result, false)
	}}
	cmd.Flags().StringVar(&inputPath, "input", "", "Envelope JSON input")
	cmd.Flags().StringVar(&actor, "actor", "", "Finalized activation actor ID")
	cmd.Flags().StringVar(&key, "idempotency-key", "", "Idempotency key")
	cmd.SetFlagErrorFunc(canonicalFlagError("envelope.create"))
	return cmd
}

func newEnvelopeGetCommand() *cobra.Command {
	cmd := &cobra.Command{Use: "get <envelope-id>", Short: "Read an envelope.", Args: canonicalExactArgs("envelope.get", 1), RunE: func(cmd *cobra.Command, args []string) error {
		workspaceRoot, err := issueMutationWorkspaceRoot(cmd)
		if err != nil {
			return renderEnvelopeError(cmd, "envelope.get", err, "")
		}
		envelope, err := workflow.GetEnvelope(workspaceRoot, issueMigrationsDir(), args[0])
		if err != nil {
			return renderEnvelopeError(cmd, "envelope.get", err, "")
		}
		jsonMode, err := readJSONMode(cmd)
		if err != nil {
			return err
		}
		if jsonMode {
			return cli.WriteSuccess(cmd.OutOrStdout(), "envelope.get", map[string]any{"envelope": envelope}, nil)
		}
		_, err = fmt.Fprintf(cmd.OutOrStdout(), "%s\t%s\n", envelope.ID, envelope.Status)
		return err
	}}
	cmd.SetFlagErrorFunc(canonicalFlagError("envelope.get"))
	return cmd
}

func newEnvelopeActivateCommand() *cobra.Command {
	var inputPath, actor, key string
	cmd := &cobra.Command{Use: "activate", Short: "Activate an envelope revision.", Args: canonicalNoArgs("envelope.activate"), RunE: func(cmd *cobra.Command, _ []string) error {
		if err := requireCanonicalFlag(cmd, "envelope.activate", "input"); err != nil {
			return err
		}
		if err := requireCanonicalFlag(cmd, "envelope.activate", "actor"); err != nil {
			return err
		}
		if err := requireCanonicalFlag(cmd, "envelope.activate", "idempotency-key"); err != nil {
			return err
		}
		var request workflow.EnvelopeActivateRequest
		if err := decodeClosedEnvelopeInput(cmd.InOrStdin(), inputPath, &request); err != nil {
			return renderEnvelopeError(cmd, "envelope.activate", err, actor)
		}
		workspaceRoot, err := issueMutationWorkspaceRoot(cmd)
		if err != nil {
			return renderEnvelopeError(cmd, "envelope.activate", err, actor)
		}
		result, err := workflow.ActivateEnvelope(workflow.EnvelopeOptions{WorkspaceRoot: workspaceRoot, MigrationsDir: issueMigrationsDir(), ActorID: actor, IdempotencyKey: key}, request)
		if err != nil {
			return renderEnvelopeError(cmd, "envelope.activate", err, actor)
		}
		return writeEnvelopeSuccess(cmd, "envelope.activate", result, true)
	}}
	cmd.Flags().StringVar(&inputPath, "input", "", "Envelope activation JSON input")
	cmd.Flags().StringVar(&actor, "actor", "", "Finalized activation actor ID")
	cmd.Flags().StringVar(&key, "idempotency-key", "", "Idempotency key")
	cmd.SetFlagErrorFunc(canonicalFlagError("envelope.activate"))
	return cmd
}

func newEnvelopeCompleteCommand() *cobra.Command {
	var inputPath, actor, key string
	cmd := &cobra.Command{Use: "complete", Short: "Complete a terminal envelope.", Args: canonicalNoArgs("envelope.complete"), RunE: func(cmd *cobra.Command, _ []string) error {
		for _, flag := range []string{"input", "actor", "idempotency-key"} {
			if err := requireCanonicalFlag(cmd, "envelope.complete", flag); err != nil {
				return err
			}
		}
		var request workflow.EnvelopeCompleteRequest
		if err := decodeClosedEnvelopeInput(cmd.InOrStdin(), inputPath, &request); err != nil {
			return renderEnvelopeError(cmd, "envelope.complete", err, actor)
		}
		workspaceRoot, err := issueMutationWorkspaceRoot(cmd)
		if err != nil {
			return renderEnvelopeError(cmd, "envelope.complete", err, actor)
		}
		result, err := workflow.CompleteEnvelope(workflow.EnvelopeOptions{WorkspaceRoot: workspaceRoot, MigrationsDir: issueMigrationsDir(), ActorID: actor, IdempotencyKey: key}, request)
		if err != nil {
			return renderEnvelopeError(cmd, "envelope.complete", err, actor)
		}
		return writeEnvelopeSuccess(cmd, "envelope.complete", result, true)
	}}
	cmd.Flags().StringVar(&inputPath, "input", "", "Envelope completion JSON input")
	cmd.Flags().StringVar(&actor, "actor", "", "Finalized activation actor ID")
	cmd.Flags().StringVar(&key, "idempotency-key", "", "Idempotency key")
	cmd.SetFlagErrorFunc(canonicalFlagError("envelope.complete"))
	return cmd
}

func writeEnvelopeSuccess(cmd *cobra.Command, operation string, result workflow.EnvelopeMutationResult, withEvaluation bool) error {
	jsonMode, err := readJSONMode(cmd)
	if err != nil {
		return err
	}
	if jsonMode {
		payload := map[string]any{"envelope": result.Envelope, "idempotency_key": result.IdempotencyKey, "replayed": result.Replayed}
		if withEvaluation {
			payload["evaluation"] = result.Evaluation
		}
		occurredAt := result.Envelope.CreatedAt
		if withEvaluation && result.Envelope.ActivatedAt != nil {
			occurredAt = *result.Envelope.ActivatedAt
		}
		if operation == "envelope.complete" && result.Envelope.CompletedAt != nil {
			occurredAt = *result.Envelope.CompletedAt
		}
		return cli.WriteSuccess(cmd.OutOrStdout(), operation, payload, &cli.Event{ID: result.EventID, OccurredAt: occurredAt})
	}
	_, err = fmt.Fprintf(cmd.OutOrStdout(), "%s envelope %s\n", result.Envelope.Status, result.Envelope.ID)
	return err
}

func decodeClosedEnvelopeInput(stdin io.Reader, path string, target any) error {
	var payload []byte
	var err error
	if path == "-" {
		payload, err = io.ReadAll(io.LimitReader(stdin, 1_048_577))
	} else {
		payload, err = os.ReadFile(path)
	}
	if err != nil {
		return fmt.Errorf("read envelope input: %w", err)
	}
	if len(payload) > 1_048_576 {
		return fmt.Errorf("envelope input exceeds 1048576 bytes")
	}
	if err := rejectAuthoredBudgetState(payload, target); err != nil {
		return err
	}
	decoder := json.NewDecoder(strings.NewReader(string(payload)))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(target); err != nil {
		return fmt.Errorf("decode envelope input: %w", err)
	}
	var extra any
	if err := decoder.Decode(&extra); err != io.EOF {
		return fmt.Errorf("envelope input must contain one object")
	}
	return nil
}

func rejectAuthoredBudgetState(payload []byte, target any) error {
	if _, ok := target.(*workflow.EnvelopeCreateRequest); !ok {
		return nil
	}
	var input struct {
		Budgets []map[string]json.RawMessage `json:"budgets"`
	}
	if err := json.Unmarshal(payload, &input); err != nil {
		return nil
	}
	for _, budget := range input.Budgets {
		if _, found := budget["state"]; found {
			return fmt.Errorf("envelope budget state is output-only")
		}
	}
	return nil
}

func renderEnvelopeError(cmd *cobra.Command, operation string, cause error, actor string) error {
	jsonMode, err := readJSONMode(cmd)
	if err != nil {
		return err
	}
	var notFound *workflow.EnvelopeNotFoundError
	var authority *workflow.EnvelopeAuthorityError
	var completion *workflow.EnvelopeCompletionError
	var protectedDecision *workflow.StopBlockedError
	if errors.As(cause, &notFound) {
		return writeEnvelopeFailure(cmd, jsonMode, operation, cause, cli.CodeEnvelopeNotFound, "envelope was not found", &cli.Subject{Type: cli.SubjectEnvelope, ID: notFound.ID}, map[string]any{"id": notFound.ID})
	}
	if errors.As(cause, &authority) {
		return writeEnvelopeFailure(cmd, jsonMode, operation, cause, cli.CodeAuthorityUnavailable, "finalized envelope authority is unavailable", nil, map[string]any{"operation": operation, "actor_id": actor, "required_actions": []string{authority.Capability}, "required_capabilities": []string{authority.Capability}})
	}
	if errors.As(cause, &protectedDecision) {
		failure, err := protectedDecisionFailure(protectedDecision)
		if err != nil {
			return err
		}
		if !jsonMode {
			_ = renderIssueError(cmd.ErrOrStderr(), cause)
			return cli.NewFailureError(failure, cause)
		}
		return cli.WriteFailureAndReturn(cmd.ErrOrStderr(), operation, failure, fmt.Errorf("%s failed", operation))
	}
	if errors.As(cause, &completion) {
		return writeEnvelopeFailure(cmd, jsonMode, operation, cause, cli.CodeObligationUnsatisfied, "Envelope completion obligations are unsatisfied", &cli.Subject{Type: cli.SubjectEnvelope, ID: completion.EnvelopeID}, map[string]any{"reasons": envelopeReasonValues(completion.Reasons)})
	}
	return writeEnvelopeFailure(cmd, jsonMode, operation, cause, cli.CodeInvalidArgument, "invalid envelope input", nil, map[string]any{"argument": "input", "value": nil, "allowed": []string{"accepted " + operation + " input"}})
}

func envelopeReasonValues(reasons []workflow.EvaluationReason) []any {
	values := make([]any, 0, len(reasons))
	for _, reason := range reasons {
		values = append(values, map[string]any{"code": reason.Code, "detail": reason.Detail, "subject_id": reason.SubjectID, "subject_type": reason.SubjectType})
	}
	return values
}

func writeEnvelopeFailure(cmd *cobra.Command, jsonMode bool, operation string, cause error, code cli.Code, message string, subject *cli.Subject, details map[string]any) error {
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
