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
	rootFlagCapabilities["activation prepare"] = rootFlagCapability{Project: true, JSON: true}
	rootFlagCapabilities["activation component"] = rootFlagCapability{Project: true, JSON: true}
	rootFlagCapabilities["activation commit"] = rootFlagCapability{Project: true, JSON: true}
	rootFlagCapabilities["activation abort"] = rootFlagCapability{Project: true, JSON: true}
	rootFlagCapabilities["activation finalize"] = rootFlagCapability{Project: true, JSON: true}
	rootFlagCapabilities["activation status"] = rootFlagCapability{Project: true, JSON: true}
}

func newActivationCommand() *cobra.Command {
	cmd := &cobra.Command{Use: "activation", Short: "Manage kit activation trust."}
	cmd.AddCommand(newActivationPrepareCommand())
	cmd.AddCommand(newActivationComponentCommand())
	cmd.AddCommand(newActivationTransitionCommand("commit"))
	cmd.AddCommand(newActivationTransitionCommand("abort"))
	cmd.AddCommand(newActivationFinalizeCommand())
	cmd.AddCommand(newActivationStatusCommand())
	return cmd
}

func newActivationFinalizeCommand() *cobra.Command {
	var inputPath, key, bootstrapActor string
	cmd := &cobra.Command{Use: "finalize", Short: "Finalize an activation.", Args: canonicalNoArgs("activation.finalize"), RunE: func(cmd *cobra.Command, _ []string) error {
		if err := requireCanonicalFlag(cmd, "activation.finalize", "input"); err != nil {
			return err
		}
		if err := requireCanonicalFlag(cmd, "activation.finalize", "idempotency-key"); err != nil {
			return err
		}
		var request workflow.ActivationFinalizeRequest
		if err := decodeClosedActivationInput(cmd.InOrStdin(), inputPath, &request); err != nil {
			return renderActivationError(cmd, "activation.finalize", err)
		}
		workspaceRoot, err := issueMutationWorkspaceRoot(cmd)
		if err != nil {
			return renderActivationError(cmd, "activation.finalize", err)
		}
		result, err := workflow.FinalizeActivation(workflow.ActivationTransitionOptions{WorkspaceRoot: workspaceRoot, MigrationsDir: issueMigrationsDir(), BootstrapActor: bootstrapActor, IdempotencyKey: key}, request)
		if err != nil {
			return renderActivationError(cmd, "activation.finalize", err)
		}
		jsonMode, err := readJSONMode(cmd)
		if err != nil {
			return err
		}
		if jsonMode {
			return cli.WriteSuccess(cmd.OutOrStdout(), "activation.finalize", map[string]any{"activation": result.Activation, "idempotency_key": result.IdempotencyKey, "replayed": result.Replayed}, &cli.Event{ID: result.EventID, OccurredAt: result.Activation.UpdatedAt})
		}
		_, err = fmt.Fprintf(cmd.OutOrStdout(), "finalized activation %s\n", result.Activation.ID)
		return err
	}}
	cmd.Flags().StringVar(&inputPath, "input", "", "Activation JSON input")
	cmd.Flags().StringVar(&key, "idempotency-key", "", "Idempotency key")
	cmd.Flags().StringVar(&bootstrapActor, "actor", "", "Untrusted bootstrap attribution")
	cmd.SetFlagErrorFunc(canonicalFlagError("activation.finalize"))
	return cmd
}

func newActivationStatusCommand() *cobra.Command {
	cmd := &cobra.Command{Use: "status <activation-id>", Short: "Read an activation status.", Args: canonicalExactArgs("activation.status", 1), RunE: func(cmd *cobra.Command, args []string) error {
		workspaceRoot, err := issueMutationWorkspaceRoot(cmd)
		if err != nil {
			return renderActivationError(cmd, "activation.status", err)
		}
		activation, err := workflow.ActivationStatus(workspaceRoot, issueMigrationsDir(), args[0])
		if err != nil {
			return renderActivationError(cmd, "activation.status", err)
		}
		jsonMode, err := readJSONMode(cmd)
		if err != nil {
			return err
		}
		if jsonMode {
			return cli.WriteSuccess(cmd.OutOrStdout(), "activation.status", map[string]any{"activation": activation}, nil)
		}
		_, err = fmt.Fprintf(cmd.OutOrStdout(), "%s\t%s\n", activation.ID, activation.Status)
		return err
	}}
	cmd.SetFlagErrorFunc(canonicalFlagError("activation.status"))
	return cmd
}

func newActivationTransitionCommand(name string) *cobra.Command {
	var inputPath, key, bootstrapActor string
	operation := "activation." + name
	cmd := &cobra.Command{Use: name, Short: name + " an activation.", Args: canonicalNoArgs(operation), RunE: func(cmd *cobra.Command, _ []string) error {
		if err := requireCanonicalFlag(cmd, operation, "input"); err != nil {
			return err
		}
		if err := requireCanonicalFlag(cmd, operation, "idempotency-key"); err != nil {
			return err
		}
		workspaceRoot, err := issueMutationWorkspaceRoot(cmd)
		if err != nil {
			return renderActivationError(cmd, operation, err)
		}
		options := workflow.ActivationTransitionOptions{WorkspaceRoot: workspaceRoot, MigrationsDir: issueMigrationsDir(), BootstrapActor: bootstrapActor, IdempotencyKey: key}
		var result workflow.ActivationPrepareResult
		if name == "commit" {
			var request workflow.ActivationCommitRequest
			if err := decodeClosedActivationInput(cmd.InOrStdin(), inputPath, &request); err != nil {
				return renderActivationError(cmd, operation, err)
			}
			result, err = workflow.CommitActivation(options, request)
		} else {
			var request workflow.ActivationAbortRequest
			if err := decodeClosedActivationInput(cmd.InOrStdin(), inputPath, &request); err != nil {
				return renderActivationError(cmd, operation, err)
			}
			result, err = workflow.AbortActivation(options, request)
		}
		if err != nil {
			return renderActivationError(cmd, operation, err)
		}
		jsonMode, err := readJSONMode(cmd)
		if err != nil {
			return err
		}
		if jsonMode {
			return cli.WriteSuccess(cmd.OutOrStdout(), operation, map[string]any{"activation": result.Activation, "idempotency_key": result.IdempotencyKey, "replayed": result.Replayed}, &cli.Event{ID: result.EventID, OccurredAt: result.Activation.UpdatedAt})
		}
		_, err = fmt.Fprintf(cmd.OutOrStdout(), "%sed activation %s\n", name, result.Activation.ID)
		return err
	}}
	cmd.Flags().StringVar(&inputPath, "input", "", "Activation JSON input")
	cmd.Flags().StringVar(&key, "idempotency-key", "", "Idempotency key")
	cmd.Flags().StringVar(&bootstrapActor, "actor", "", "Untrusted bootstrap attribution")
	cmd.SetFlagErrorFunc(canonicalFlagError(operation))
	return cmd
}

func newActivationComponentCommand() *cobra.Command {
	var inputPath, key, bootstrapActor string
	cmd := &cobra.Command{
		Use: "component", Short: "Stage an activation component.", Args: canonicalNoArgs("activation.component"),
		RunE: func(cmd *cobra.Command, _ []string) error {
			if err := requireCanonicalFlag(cmd, "activation.component", "input"); err != nil {
				return err
			}
			if err := requireCanonicalFlag(cmd, "activation.component", "idempotency-key"); err != nil {
				return err
			}
			var request workflow.ActivationComponentRequest
			if err := decodeClosedActivationInput(cmd.InOrStdin(), inputPath, &request); err != nil {
				return renderActivationError(cmd, "activation.component", err)
			}
			workspaceRoot, err := issueMutationWorkspaceRoot(cmd)
			if err != nil {
				return renderActivationError(cmd, "activation.component", err)
			}
			result, err := workflow.StageActivationComponent(workflow.ActivationComponentOptions{WorkspaceRoot: workspaceRoot, MigrationsDir: issueMigrationsDir(), BootstrapActor: bootstrapActor, IdempotencyKey: key, Request: request})
			if err != nil {
				return renderActivationError(cmd, "activation.component", err)
			}
			jsonMode, err := readJSONMode(cmd)
			if err != nil {
				return err
			}
			if jsonMode {
				return cli.WriteSuccess(cmd.OutOrStdout(), "activation.component", map[string]any{"activation": result.Activation, "component": result.Component, "idempotency_key": result.IdempotencyKey, "replayed": result.Replayed}, &cli.Event{ID: result.EventID, OccurredAt: result.Activation.UpdatedAt})
			}
			_, err = fmt.Fprintf(cmd.OutOrStdout(), "staged activation component %s\n", result.Component.ID)
			return err
		},
	}
	cmd.Flags().StringVar(&inputPath, "input", "", "Activation JSON input")
	cmd.Flags().StringVar(&key, "idempotency-key", "", "Idempotency key")
	cmd.Flags().StringVar(&bootstrapActor, "actor", "", "Untrusted bootstrap attribution")
	cmd.SetFlagErrorFunc(canonicalFlagError("activation.component"))
	return cmd
}

func newActivationPrepareCommand() *cobra.Command {
	var inputPath, key, bootstrapActor string
	cmd := &cobra.Command{
		Use:   "prepare",
		Short: "Prepare an untrusted activation.",
		Args:  canonicalNoArgs("activation.prepare"),
		RunE: func(cmd *cobra.Command, _ []string) error {
			if err := requireCanonicalFlag(cmd, "activation.prepare", "input"); err != nil {
				return err
			}
			if err := requireCanonicalFlag(cmd, "activation.prepare", "idempotency-key"); err != nil {
				return err
			}
			request, err := readActivationPrepareInput(cmd, inputPath)
			if err != nil {
				return renderActivationError(cmd, "activation.prepare", err)
			}
			workspaceRoot, err := issueMutationWorkspaceRoot(cmd)
			if err != nil {
				return renderActivationError(cmd, "activation.prepare", err)
			}
			result, err := workflow.PrepareActivation(workflow.ActivationPrepareOptions{
				WorkspaceRoot: workspaceRoot, MigrationsDir: issueMigrationsDir(), BootstrapActor: bootstrapActor,
				IdempotencyKey: key, Request: request,
			})
			if err != nil {
				return renderActivationError(cmd, "activation.prepare", err)
			}
			jsonMode, err := readJSONMode(cmd)
			if err != nil {
				return err
			}
			if jsonMode {
				return cli.WriteSuccess(cmd.OutOrStdout(), "activation.prepare", map[string]any{
					"activation": result.Activation, "idempotency_key": result.IdempotencyKey, "replayed": result.Replayed,
				}, &cli.Event{ID: result.EventID, OccurredAt: result.Activation.UpdatedAt})
			}
			_, err = fmt.Fprintf(cmd.OutOrStdout(), "prepared activation %s; bootstrap attribution is untrusted\n", result.Activation.ID)
			return err
		},
	}
	cmd.Flags().StringVar(&inputPath, "input", "", "Activation JSON input")
	cmd.Flags().StringVar(&key, "idempotency-key", "", "Idempotency key")
	cmd.Flags().StringVar(&bootstrapActor, "actor", "", "Untrusted bootstrap attribution")
	cmd.SetFlagErrorFunc(canonicalFlagError("activation.prepare"))
	return cmd
}

func readActivationPrepareInput(cmd *cobra.Command, path string) (workflow.ActivationPrepareRequest, error) {
	var request workflow.ActivationPrepareRequest
	if err := decodeClosedActivationInput(cmd.InOrStdin(), path, &request); err != nil {
		return workflow.ActivationPrepareRequest{}, err
	}
	return request, nil
}

func decodeClosedActivationInput(stdin io.Reader, path string, target any) error {
	var payload []byte
	var err error
	if path == "-" {
		payload, err = io.ReadAll(io.LimitReader(stdin, 1_048_577))
	} else {
		payload, err = os.ReadFile(path)
	}
	if err != nil {
		return fmt.Errorf("read activation input: %w", err)
	}
	if len(payload) > 1_048_576 {
		return fmt.Errorf("activation input exceeds 1048576 bytes")
	}
	decoder := json.NewDecoder(strings.NewReader(string(payload)))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(target); err != nil {
		return fmt.Errorf("decode activation input: %w", err)
	}
	if decoder.More() {
		return fmt.Errorf("activation input must contain one object")
	}
	var extra any
	if err := decoder.Decode(&extra); err != io.EOF {
		return fmt.Errorf("activation input must contain one object")
	}
	return nil
}

func renderActivationError(cmd *cobra.Command, operation string, cause error) error {
	jsonMode, err := readJSONMode(cmd)
	if err != nil {
		return err
	}
	if !jsonMode {
		if err := renderIssueError(cmd.ErrOrStderr(), cause); err != nil {
			return err
		}
		return cause
	}
	var conflict *workflow.ActivationConflictError
	if errors.As(cause, &conflict) {
		failure, err := cli.NewFailure(cli.CodeActivationConflict, "another activation blocks this request", &cli.Subject{Type: cli.SubjectActivation, ID: ""}, map[string]any{"active_id": conflict.ActiveID, "nonterminal_id": conflict.NonterminalID, "requested_direction": "activate"})
		if err == nil {
			err = cli.WriteFailure(cmd.ErrOrStderr(), operation, failure)
		}
		if err != nil {
			return err
		}
		return cause
	}
	var recovery *workflow.ActivationRecoveryUnavailableError
	if errors.As(cause, &recovery) {
		failure, err := cli.NewFailure(cli.CodeActivationRecoveryUnavailable, "activation recovery is unavailable", &cli.Subject{Type: cli.SubjectActivation, ID: ""}, map[string]any{"status": recovery.Status, "point_of_no_return": true, "owner_phase": 22})
		if err == nil {
			err = cli.WriteFailure(cmd.ErrOrStderr(), operation, failure)
		}
		if err != nil {
			return err
		}
		return cause
	}
	failure, err := cli.NewFailure(cli.CodeInvalidArgument, "invalid activation input", nil, map[string]any{"argument": "input", "value": nil, "allowed": []string{"accepted activation input"}})
	if err != nil {
		return err
	}
	if err := cli.WriteFailure(cmd.ErrOrStderr(), operation, failure); err != nil {
		return err
	}
	return cause
}
