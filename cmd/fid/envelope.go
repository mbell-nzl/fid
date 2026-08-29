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
}

func newEnvelopeCommand() *cobra.Command {
	cmd := &cobra.Command{Use: "envelope", Short: "Manage the bounded work envelope."}
	cmd.AddCommand(newEnvelopeCreateCommand())
	cmd.AddCommand(newEnvelopeGetCommand())
	cmd.AddCommand(newEnvelopeActivateCommand())
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

func renderEnvelopeError(cmd *cobra.Command, operation string, cause error, actor string) error {
	jsonMode, err := readJSONMode(cmd)
	if err != nil {
		return err
	}
	var notFound *workflow.EnvelopeNotFoundError
	var authority *workflow.EnvelopeAuthorityError
	if errors.As(cause, &notFound) {
		return writeEnvelopeFailure(cmd, jsonMode, operation, cause, cli.CodeEnvelopeNotFound, "envelope was not found", &cli.Subject{Type: cli.SubjectEnvelope, ID: notFound.ID}, map[string]any{"id": notFound.ID})
	}
	if errors.As(cause, &authority) {
		return writeEnvelopeFailure(cmd, jsonMode, operation, cause, cli.CodeAuthorityUnavailable, "finalized envelope authority is unavailable", nil, map[string]any{"operation": operation, "actor_id": actor, "required_actions": []string{authority.Capability}, "required_capabilities": []string{authority.Capability}})
	}
	return writeEnvelopeFailure(cmd, jsonMode, operation, cause, cli.CodeInvalidArgument, "invalid envelope input", nil, map[string]any{"argument": "input", "value": nil, "allowed": []string{"accepted " + operation + " input"}})
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
