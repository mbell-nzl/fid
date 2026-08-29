package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/mbell-nzl/fid/internal/cli"
	"github.com/mbell-nzl/fid/internal/evidence"
	"github.com/spf13/cobra"
)

func init() {
	rootFlagCapabilities["gate register"] = rootFlagCapability{Project: true, JSON: true}
	rootFlagCapabilities["gate resolve"] = rootFlagCapability{Project: true, JSON: true}
}

func newGateCommand() *cobra.Command {
	cmd := &cobra.Command{Use: "gate", Short: "Register immutable gate definitions."}
	cmd.AddCommand(newGateRegisterCommand())
	cmd.AddCommand(newGateResolveCommand())
	return cmd
}

func newGateRegisterCommand() *cobra.Command {
	var inputPath, actor, key string
	cmd := &cobra.Command{Use: "register", Short: "Register an activated gate definition.", Args: canonicalNoArgs("gate.register"), RunE: func(cmd *cobra.Command, _ []string) error {
		if err := requireCanonicalFlag(cmd, "gate.register", "input"); err != nil {
			return err
		}
		if err := requireCanonicalFlag(cmd, "gate.register", "actor"); err != nil {
			return err
		}
		if err := requireCanonicalFlag(cmd, "gate.register", "idempotency-key"); err != nil {
			return err
		}
		var request evidence.RegisterRequest
		if err := decodeClosedGateInput(cmd.InOrStdin(), inputPath, &request); err != nil {
			return renderGateError(cmd, "gate.register", err)
		}
		workspaceRoot, err := issueMutationWorkspaceRoot(cmd)
		if err != nil {
			return renderGateError(cmd, "gate.register", err)
		}
		result, err := evidence.Register(evidence.RegisterOptions{WorkspaceRoot: workspaceRoot, MigrationsDir: issueMigrationsDir(), ActorID: actor, IdempotencyKey: key, Request: request})
		if err != nil {
			return renderGateError(cmd, "gate.register", err)
		}
		jsonMode, err := readJSONMode(cmd)
		if err != nil {
			return err
		}
		if jsonMode {
			return cli.WriteSuccess(cmd.OutOrStdout(), "gate.register", map[string]any{"registration": result.Registration, "idempotency_key": result.IdempotencyKey, "replayed": result.Replayed}, &cli.Event{ID: result.EventID, OccurredAt: result.Registration.CreatedAt})
		}
		_, err = fmt.Fprintf(cmd.OutOrStdout(), "registered gate %s version %d\n", result.Registration.ID, result.Registration.Version)
		return err
	}}
	cmd.Flags().StringVar(&inputPath, "input", "", "Gate registration JSON input")
	cmd.Flags().StringVar(&actor, "actor", "", "Finalized activation actor ID")
	cmd.Flags().StringVar(&key, "idempotency-key", "", "Idempotency key")
	cmd.SetFlagErrorFunc(canonicalFlagError("gate.register"))
	return cmd
}

func newGateResolveCommand() *cobra.Command {
	var version int
	cmd := &cobra.Command{Use: "resolve <registration-id>", Short: "Resolve an immutable gate registration.", Args: canonicalExactArgs("gate.resolve", 1), RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireCanonicalFlag(cmd, "gate.resolve", "version"); err != nil {
			return err
		}
		workspaceRoot, err := issueMutationWorkspaceRoot(cmd)
		if err != nil {
			return renderGateError(cmd, "gate.resolve", err)
		}
		registration, err := evidence.Resolve(workspaceRoot, issueMigrationsDir(), args[0], version)
		if err != nil {
			return renderGateError(cmd, "gate.resolve", err)
		}
		jsonMode, err := readJSONMode(cmd)
		if err != nil {
			return err
		}
		if jsonMode {
			return cli.WriteSuccess(cmd.OutOrStdout(), "gate.resolve", map[string]any{"registration": registration}, nil)
		}
		_, err = fmt.Fprintf(cmd.OutOrStdout(), "%s\t%d\n", registration.ID, registration.Version)
		return err
	}}
	cmd.Flags().IntVar(&version, "version", 0, "Registration version")
	cmd.SetFlagErrorFunc(canonicalFlagError("gate.resolve"))
	return cmd
}

func decodeClosedGateInput(stdin io.Reader, path string, target any) error {
	var payload []byte
	var err error
	if path == "-" {
		payload, err = io.ReadAll(io.LimitReader(stdin, 1_048_577))
	} else {
		payload, err = os.ReadFile(path)
	}
	if err != nil {
		return fmt.Errorf("read gate input: %w", err)
	}
	if len(payload) > 1_048_576 {
		return fmt.Errorf("gate input exceeds 1048576 bytes")
	}
	decoder := json.NewDecoder(strings.NewReader(string(payload)))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(target); err != nil {
		return fmt.Errorf("decode gate input: %w", err)
	}
	if decoder.More() {
		return fmt.Errorf("gate input must contain one object")
	}
	var extra any
	if err := decoder.Decode(&extra); err != io.EOF {
		return fmt.Errorf("gate input must contain one object")
	}
	return nil
}

func renderGateError(cmd *cobra.Command, operation string, cause error) error {
	jsonMode, err := readJSONMode(cmd)
	if err != nil {
		return err
	}
	var unavailable *evidence.ActiveActivationUnavailableError
	var policyErr *evidence.ExecutorPolicyUnsupportedError
	var definitionErr *evidence.GateDefinitionNotFoundError
	var registrationErr *evidence.RegistrationNotFoundError
	var actorErr *evidence.TrustedActorUnavailableError
	var deniedErr *evidence.AuthorityDeniedError
	switch {
	case errors.As(cause, &unavailable):
		return writeGateFailure(cmd, jsonMode, operation, cause, cli.CodeAuthorityUnavailable, "finalized activation authority is unavailable", nil, map[string]any{"operation": operation, "actor_id": nil, "required_actions": []string{"gate.register"}, "required_capabilities": []string{"gate.register"}})
	case errors.As(cause, &policyErr):
		return writeGateFailure(cmd, jsonMode, operation, cause, cli.CodeExecutorPolicyUnsupported, "executor policy is unsupported", &cli.Subject{Type: cli.SubjectGateDefinition, ID: ""}, map[string]any{"fields": policyErr.Fields, "supported": map[string]any{"sandbox": "none", "network": "allowed", "write": "allowed", "shell": false}})
	case errors.As(cause, &definitionErr):
		return writeGateFailure(cmd, jsonMode, operation, cause, cli.CodeRegistrationNotFound, "activated gate definition was not found", &cli.Subject{Type: cli.SubjectGateDefinition, ID: definitionErr.DefinitionID}, map[string]any{"id": definitionErr.DefinitionID})
	case errors.As(cause, &registrationErr):
		return writeGateFailure(cmd, jsonMode, operation, cause, cli.CodeRegistrationNotFound, "gate registration was not found", &cli.Subject{Type: cli.SubjectGateRegistration, ID: registrationErr.ID}, map[string]any{"id": registrationErr.ID})
	case errors.As(cause, &actorErr):
		return writeGateFailure(cmd, jsonMode, operation, cause, cli.CodeAuthorityUnknownActor, "finalized activation actor is unknown", nil, map[string]any{"operation": operation, "actor_id": actorErr.ActorID, "required_actions": []string{"gate.register"}, "required_capabilities": []string{"gate.register"}})
	case errors.As(cause, &deniedErr):
		return writeGateFailure(cmd, jsonMode, operation, cause, cli.CodeAuthorityDenied, "finalized activation actor is denied", nil, map[string]any{"operation": operation, "actor_id": deniedErr.ActorID, "required_actions": []string{"gate.register"}, "required_capabilities": []string{"gate.register"}})
	default:
		return writeGateFailure(cmd, jsonMode, operation, cause, cli.CodeInvalidArgument, "invalid gate request", nil, map[string]any{"argument": "input", "value": nil, "allowed": []string{"accepted gate registration input"}})
	}
}

func writeGateFailure(cmd *cobra.Command, jsonMode bool, operation string, cause error, code cli.Code, message string, subject *cli.Subject, details map[string]any) error {
	failure, err := cli.NewFailure(code, message, subject, details)
	if err != nil {
		return err
	}
	if !jsonMode {
		if err := renderIssueError(cmd.ErrOrStderr(), cause); err != nil {
			return err
		}
		return cli.NewFailureError(failure, cause)
	}
	return cli.WriteFailureAndReturn(cmd.ErrOrStderr(), operation, failure, fmt.Errorf("%s failed", operation))
}
