package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/mbell-nzl/fid/internal/cli"
	"github.com/mbell-nzl/fid/internal/evidence"
	"github.com/mbell-nzl/fid/internal/workflow"
	"github.com/spf13/cobra"
)

func init() {
	rootFlagCapabilities["gate register"] = rootFlagCapability{Project: true, JSON: true}
	rootFlagCapabilities["gate resolve"] = rootFlagCapability{Project: true, JSON: true}
	rootFlagCapabilities["gate run"] = rootFlagCapability{Project: true, JSON: true}
}

func newGateCommand() *cobra.Command {
	cmd := &cobra.Command{Use: "gate", Short: "Register immutable gate definitions."}
	cmd.AddCommand(newGateRegisterCommand())
	cmd.AddCommand(newGateResolveCommand())
	cmd.AddCommand(newGateRunCommand())
	return cmd
}

func newGateRunCommand() *cobra.Command {
	var inputPath, actor, key string
	cmd := &cobra.Command{Use: "run", Short: "Execute one registered gate.", Args: canonicalNoArgs("gate.run"), RunE: func(cmd *cobra.Command, _ []string) error {
		for _, flag := range []string{"input", "actor", "idempotency-key"} {
			if err := requireCanonicalFlag(cmd, "gate.run", flag); err != nil {
				return err
			}
		}
		var request evidence.GateRunRequest
		if err := decodeClosedGateInput(cmd.InOrStdin(), inputPath, &request); err != nil {
			return renderGateError(cmd, "gate.run", err)
		}
		workspaceRoot, err := issueMutationWorkspaceRoot(cmd)
		if err != nil {
			return renderGateError(cmd, "gate.run", err)
		}
		result, err := evidence.RunGate(context.Background(), evidence.GateRunOptions{WorkspaceRoot: workspaceRoot, MigrationsDir: issueMigrationsDir(), ActorID: actor, IdempotencyKey: key, Request: request})
		if err != nil {
			var denied *evidence.BudgetAdmissionDeniedError
			if errors.As(err, &denied) {
				return writeCommittedBudgetDenial(cmd, request.RunID, result, denied)
			}
			if len(result.Live.Frames) == 0 {
				return renderGateError(cmd, "gate.run", err)
			}
			response := gateInternalFailureEnvelope()
			if appendErr := result.Live.AppendFidError(result.Attempt.ID, response); appendErr != nil {
				return appendErr
			}
			return writeGateLiveResult(cmd, result.Live, "unsatisfied", true)
		}
		if result.Receipt == nil {
			return fmt.Errorf("gate run: terminal execution has no receipt")
		}
		response, failure, err := gateTerminalResponse(result)
		if err != nil {
			return err
		}
		if err := result.Live.AppendTerminal(evidence.LiveTerminal{AttemptID: result.Attempt.ID, ReceiptID: &result.Receipt.ID, Outcome: result.Receipt.Outcome, ExitCode: result.Receipt.ExitCode, Signal: result.Receipt.Signal, TimedOut: result.Receipt.TimedOut, ArtifactIDs: []string{}, Response: response, Replayed: result.Replayed}); err != nil {
			return err
		}
		if err := writeGateLiveResult(cmd, result.Live, "unsatisfied", false); err != nil {
			return err
		}
		if failure != nil {
			return cli.NewFailureError(*failure, fmt.Errorf("gate run observed %s", result.Receipt.Outcome))
		}
		return nil
	}}
	cmd.Flags().StringVar(&inputPath, "input", "", "Gate run JSON input")
	cmd.Flags().StringVar(&actor, "actor", "", "Finalized activation actor ID")
	cmd.Flags().StringVar(&key, "idempotency-key", "", "Idempotency key")
	cmd.SetFlagErrorFunc(canonicalFlagError("gate.run"))
	return cmd
}

func writeCommittedBudgetDenial(cmd *cobra.Command, runID string, result evidence.GateRunResult, denied *evidence.BudgetAdmissionDeniedError) error {
	jsonMode, err := readJSONMode(cmd)
	if err != nil {
		return err
	}
	code := cli.CodeBudgetDataUnavailable
	if denied.Code == "budget-exhausted" {
		code = cli.CodeBudgetExhausted
	}
	denial := denied.Denial
	budget := map[string]any{
		"name": denial.Budget.Name, "scope": denial.Budget.Scope, "unit": denial.Budget.Unit,
		"maximum": denial.Budget.Maximum, "exhaustion_action": denial.Budget.ExhaustionAction,
	}
	budget["warning_threshold"] = optionalBudgetValue(denial.Budget.WarningThreshold)
	budget["consumed"] = optionalBudgetValue(denial.Budget.Consumed)
	budget["remaining"] = optionalBudgetValue(denial.Budget.Remaining)
	failure, err := cli.NewFailure(code, "gate execution admission was denied", &cli.Subject{Type: cli.SubjectRun, ID: runID}, map[string]any{
		"budget": budget, "stop_ids": denial.StopIDs, "attention_ids": denial.AttentionIDs,
		"event": map[string]any{"id": denial.EventID, "occurred_at": denial.EventOccurredAt}, "replayed": result.Replayed,
	})
	if err != nil {
		return err
	}
	if !jsonMode {
		return cli.NewFailureError(failure, fmt.Errorf("gate run admission denied"))
	}
	response := failureEnvelope("gate.run", failure)
	if err := result.Live.AppendCommittedFidError(runID, response, result.Replayed); err != nil {
		return err
	}
	if err := writeGateLiveResult(cmd, result.Live, "unsatisfied", false); err != nil {
		return err
	}
	return cli.NewFailureError(failure, fmt.Errorf("gate run admission denied"))
}

func optionalBudgetValue(value *int64) any {
	if value == nil {
		return nil
	}
	return *value
}

func gateTerminalResponse(result evidence.GateRunResult) (map[string]any, *cli.Failure, error) {
	if result.Receipt.Outcome == "succeeded" {
		return map[string]any{"contract_version": cli.ContractVersion, "event": map[string]any{"id": result.EventID, "occurred_at": result.EventOccurredAt}, "ok": true, "operation": "gate.run", "result": map[string]any{"receipt": result.Receipt, "idempotency_key": result.IdempotencyKey, "replayed": result.Replayed}}, nil, nil
	}
	failure, err := cli.NewFailure(cli.CodeGateUnsuccessful, "gate execution was unsuccessful", &cli.Subject{Type: cli.SubjectReceipt, ID: result.Receipt.ID}, map[string]any{"attempt_id": result.Attempt.ID, "receipt_id": result.Receipt.ID, "outcome": result.Receipt.Outcome, "exit_code": optionalExitCode(result.Receipt.ExitCode), "signal": optionalSignal(result.Receipt.Signal), "timed_out": result.Receipt.TimedOut})
	if err != nil {
		return nil, nil, err
	}
	return failureEnvelope("gate.run", failure), &failure, nil
}

func optionalExitCode(value *int) any {
	if value == nil {
		return nil
	}
	return *value
}

func optionalSignal(value *string) any {
	if value == nil {
		return nil
	}
	return *value
}

func gateInternalFailureEnvelope() map[string]any {
	failure, err := cli.NewFailure(cli.CodeInternalError, "gate execution infrastructure failed", nil, map[string]any{})
	if err != nil {
		panic(err)
	}
	return failureEnvelope("gate.run", failure)
}

func failureEnvelope(operation string, failure cli.Failure) map[string]any {
	var output bytes.Buffer
	if err := cli.WriteFailure(&output, operation, failure); err != nil {
		panic(err)
	}
	var envelope map[string]any
	if err := json.Unmarshal(output.Bytes(), &envelope); err != nil {
		panic(err)
	}
	return envelope
}

func writeGateLiveResult(cmd *cobra.Command, result evidence.LiveResult, obligationStatus string, fidError bool) error {
	jsonMode, err := readJSONMode(cmd)
	if err != nil {
		return err
	}
	if jsonMode {
		for _, line := range result.JSONLines() {
			if _, err := fmt.Fprintln(cmd.OutOrStdout(), line); err != nil {
				return err
			}
		}
		if fidError {
			failure, _ := cli.NewFailure(cli.CodeInternalError, "gate execution infrastructure failed", nil, map[string]any{})
			return cli.NewFailureError(failure, fmt.Errorf("gate execution infrastructure failed"))
		}
		return nil
	}
	_, err = fmt.Fprint(cmd.OutOrStdout(), renderLiveHuman(result, obligationStatus))
	return err
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
	var stopped *workflow.StopBlockedError
	if errors.As(cause, &stopped) {
		return renderStopError(cmd, operation, cause)
	}
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
	requiredCapability := gateOperationCapability(operation)
	switch {
	case errors.As(cause, &unavailable):
		return writeGateFailure(cmd, jsonMode, operation, cause, cli.CodeAuthorityUnavailable, "finalized activation authority is unavailable", nil, map[string]any{"operation": operation, "actor_id": nil, "required_actions": []string{requiredCapability}, "required_capabilities": []string{requiredCapability}})
	case errors.As(cause, &policyErr):
		return writeGateFailure(cmd, jsonMode, operation, cause, cli.CodeExecutorPolicyUnsupported, "executor policy is unsupported", &cli.Subject{Type: cli.SubjectGateDefinition, ID: ""}, map[string]any{"fields": policyErr.Fields, "supported": map[string]any{"sandbox": "none", "network": "allowed", "write": "allowed", "shell": false}})
	case errors.As(cause, &definitionErr):
		return writeGateFailure(cmd, jsonMode, operation, cause, cli.CodeRegistrationNotFound, "activated gate definition was not found", &cli.Subject{Type: cli.SubjectGateDefinition, ID: definitionErr.DefinitionID}, map[string]any{"id": definitionErr.DefinitionID})
	case errors.As(cause, &registrationErr):
		return writeGateFailure(cmd, jsonMode, operation, cause, cli.CodeRegistrationNotFound, "gate registration was not found", &cli.Subject{Type: cli.SubjectGateRegistration, ID: registrationErr.ID}, map[string]any{"id": registrationErr.ID})
	case errors.As(cause, &actorErr):
		return writeGateFailure(cmd, jsonMode, operation, cause, cli.CodeAuthorityUnknownActor, "finalized activation actor is unknown", nil, map[string]any{"operation": operation, "actor_id": actorErr.ActorID, "required_actions": []string{requiredCapability}, "required_capabilities": []string{requiredCapability}})
	case errors.As(cause, &deniedErr):
		return writeGateFailure(cmd, jsonMode, operation, cause, cli.CodeAuthorityDenied, "finalized activation actor is denied", nil, map[string]any{"operation": operation, "actor_id": deniedErr.ActorID, "required_actions": []string{requiredCapability}, "required_capabilities": []string{requiredCapability}})
	default:
		return writeGateFailure(cmd, jsonMode, operation, cause, cli.CodeInvalidArgument, "invalid gate request", nil, map[string]any{"argument": "input", "value": nil, "allowed": []string{"accepted gate registration input"}})
	}
}

func gateOperationCapability(operation string) string {
	switch operation {
	case "artifact.verify":
		return "artifact.verify"
	case "gate.run":
		return "gate.run"
	default:
		return "gate.register"
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
