package main

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/mbell-nzl/fid/internal/cli"
	"github.com/mbell-nzl/fid/internal/journal"
	"github.com/mbell-nzl/fid/internal/workflow"
	"github.com/spf13/cobra"
)

func init() {
	rootFlagCapabilities["activation prepare"] = rootFlagCapability{Project: true, JSON: true}
	rootFlagCapabilities["activation component"] = rootFlagCapability{Project: true, JSON: true}
	rootFlagCapabilities["activation commit"] = rootFlagCapability{Project: true, JSON: true}
	rootFlagCapabilities["activation abort"] = rootFlagCapability{Project: true, JSON: true}
	rootFlagCapabilities["activation finalize"] = rootFlagCapability{Project: true, JSON: true}
	rootFlagCapabilities["activation resource"] = rootFlagCapability{Project: true, JSON: true}
	rootFlagCapabilities["activation recover"] = rootFlagCapability{Project: true, JSON: true}
	rootFlagCapabilities["activation recover-complete"] = rootFlagCapability{Project: true, JSON: true}
	rootFlagCapabilities["activation status"] = rootFlagCapability{Project: true, JSON: true}
}

func newActivationCommand() *cobra.Command {
	cmd := &cobra.Command{Use: "activation", Short: "Manage kit activation trust."}
	cmd.AddCommand(newActivationPrepareCommand())
	cmd.AddCommand(newActivationComponentCommand())
	cmd.AddCommand(newActivationTransitionCommand("commit"))
	cmd.AddCommand(newActivationTransitionCommand("abort"))
	cmd.AddCommand(newActivationFinalizeCommand())
	cmd.AddCommand(newActivationResourceCommand())
	cmd.AddCommand(newActivationRecoverCommand())
	cmd.AddCommand(newActivationRecoverCompleteCommand())
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
		if err := readDescriptorInput(cmd.InOrStdin(), inputPath, inputDescriptorRegistry["activation finalize"], &request); err != nil {
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
	installInputDescriptor(cmd, inputDescriptorRegistry["activation finalize"])
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
		return writeActivationStatusHuman(cmd, activation)
	}}
	cmd.SetFlagErrorFunc(canonicalFlagError("activation.status"))
	return cmd
}

func newActivationResourceCommand() *cobra.Command {
	return newActivationRecoveryMutationCommand("resource", func(options workflow.ActivationTransitionOptions, input any) (workflow.ActivationPrepareResult, *workflow.ActivationResourceResult, error) {
		request := input.(*workflow.ActivationResourceObservationRequest)
		result, err := workflow.ObserveActivationResource(options, *request)
		return workflow.ActivationPrepareResult{}, &result, err
	})
}

func newActivationRecoverCommand() *cobra.Command {
	return newActivationRecoveryMutationCommand("recover", func(options workflow.ActivationTransitionOptions, input any) (workflow.ActivationPrepareResult, *workflow.ActivationResourceResult, error) {
		request := input.(*workflow.ActivationRecoveryRequest)
		result, err := workflow.RecoverActivation(options, *request)
		return result, nil, err
	})
}

func newActivationRecoverCompleteCommand() *cobra.Command {
	return newActivationRecoveryMutationCommand("recover-complete", func(options workflow.ActivationTransitionOptions, input any) (workflow.ActivationPrepareResult, *workflow.ActivationResourceResult, error) {
		request := input.(*workflow.ActivationRecoveryCompleteRequest)
		result, err := workflow.CompleteActivationRecovery(options, *request)
		return result, nil, err
	})
}

type activationRecoveryMutation func(workflow.ActivationTransitionOptions, any) (workflow.ActivationPrepareResult, *workflow.ActivationResourceResult, error)

func newActivationRecoveryMutationCommand(name string, mutate activationRecoveryMutation) *cobra.Command {
	var inputPath, key, bootstrapActor string
	operation := "activation." + name
	cmd := &cobra.Command{Use: name, Short: name + " an activation.", Args: canonicalNoArgs(operation), RunE: func(cmd *cobra.Command, _ []string) error {
		if err := requireCanonicalFlag(cmd, operation, "input"); err != nil {
			return err
		}
		if err := requireCanonicalFlag(cmd, operation, "idempotency-key"); err != nil {
			return err
		}
		descriptor := inputDescriptorRegistry["activation "+name]
		request := reflect.New(descriptor.Type).Interface()
		if err := readDescriptorInput(cmd.InOrStdin(), inputPath, descriptor, request); err != nil {
			return renderActivationError(cmd, operation, err)
		}
		workspaceRoot, err := issueMutationWorkspaceRoot(cmd)
		if err != nil {
			return renderActivationError(cmd, operation, err)
		}
		options := workflow.ActivationTransitionOptions{WorkspaceRoot: workspaceRoot, MigrationsDir: issueMigrationsDir(), BootstrapActor: bootstrapActor, IdempotencyKey: key}
		activationResult, resourceResult, err := mutate(options, request)
		if err != nil {
			return renderActivationError(cmd, operation, err)
		}
		jsonMode, err := readJSONMode(cmd)
		if err != nil {
			return err
		}
		if resourceResult != nil {
			if jsonMode {
				return cli.WriteSuccess(cmd.OutOrStdout(), operation, map[string]any{"activation": resourceResult.Activation, "resource": resourceResult.Resource, "idempotency_key": resourceResult.IdempotencyKey, "replayed": resourceResult.Replayed}, &cli.Event{ID: resourceResult.EventID, OccurredAt: resourceResult.Activation.UpdatedAt})
			}
			_, err = fmt.Fprintf(cmd.OutOrStdout(), "observed activation resource %s\n", resourceResult.Resource.Component)
			return err
		}
		if jsonMode {
			return cli.WriteSuccess(cmd.OutOrStdout(), operation, map[string]any{"activation": activationResult.Activation, "idempotency_key": activationResult.IdempotencyKey, "replayed": activationResult.Replayed}, &cli.Event{ID: activationResult.EventID, OccurredAt: activationResult.Activation.UpdatedAt})
		}
		_, err = fmt.Fprintf(cmd.OutOrStdout(), "%s activation %s\n", name, activationResult.Activation.ID)
		return err
	}}
	cmd.Flags().StringVar(&inputPath, "input", "", "Activation JSON input")
	cmd.Flags().StringVar(&key, "idempotency-key", "", "Idempotency key")
	cmd.Flags().StringVar(&bootstrapActor, "actor", "", "Untrusted bootstrap attribution")
	cmd.SetFlagErrorFunc(canonicalFlagError(operation))
	installInputDescriptor(cmd, inputDescriptorRegistry["activation "+name])
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
		descriptor := inputDescriptorRegistry["activation "+name]
		request := reflect.New(descriptor.Type).Interface()
		if err := readDescriptorInput(cmd.InOrStdin(), inputPath, descriptor, request); err != nil {
			return renderActivationError(cmd, operation, err)
		}
		workspaceRoot, err := issueMutationWorkspaceRoot(cmd)
		if err != nil {
			return renderActivationError(cmd, operation, err)
		}
		options := workflow.ActivationTransitionOptions{WorkspaceRoot: workspaceRoot, MigrationsDir: issueMigrationsDir(), BootstrapActor: bootstrapActor, IdempotencyKey: key}
		var result workflow.ActivationPrepareResult
		if name == "commit" {
			result, err = workflow.CommitActivation(options, *request.(*workflow.ActivationCommitRequest))
		} else {
			result, err = workflow.AbortActivation(options, *request.(*workflow.ActivationAbortRequest))
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
	installInputDescriptor(cmd, inputDescriptorRegistry["activation "+name])
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
			if err := readDescriptorInput(cmd.InOrStdin(), inputPath, inputDescriptorRegistry["activation component"], &request); err != nil {
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
	installInputDescriptor(cmd, inputDescriptorRegistry["activation component"])
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
	installInputDescriptor(cmd, inputDescriptorRegistry["activation prepare"])
	return cmd
}

func readActivationPrepareInput(cmd *cobra.Command, path string) (workflow.ActivationPrepareRequest, error) {
	var request workflow.ActivationPrepareRequest
	if err := readDescriptorInput(cmd.InOrStdin(), path, inputDescriptorRegistry["activation prepare"], &request); err != nil {
		return workflow.ActivationPrepareRequest{}, err
	}
	return request, nil
}

func renderActivationError(cmd *cobra.Command, operation string, cause error) error {
	jsonMode, err := readJSONMode(cmd)
	if err != nil {
		return err
	}
	failure, err := activationFailure(cause)
	if err != nil {
		return err
	}
	if !jsonMode {
		return writeHumanFailure(cmd, operation, failure, cause)
	}
	return cli.WriteFailureAndReturn(cmd.ErrOrStderr(), operation, failure, cause)
}

func activationFailure(cause error) (cli.Failure, error) {
	var structural *cli.InputError
	if errors.As(cause, &structural) {
		return cli.NewFailure(cli.CodeInvalidArgument, structural.Error(), nil, structural.Details())
	}
	var idempotency *journal.IdempotencyConflictError
	if errors.As(cause, &idempotency) {
		return cli.NewFailure(cli.CodeIdempotencyConflict, "idempotency key conflicts with a different request", nil, idempotencyConflictDetails(idempotency))
	}
	var conflict *workflow.ActivationConflictError
	if errors.As(cause, &conflict) {
		return cli.NewFailure(cli.CodeActivationConflict, "another activation blocks this request", &cli.Subject{Type: cli.SubjectActivation, ID: ""}, map[string]any{"active_id": activationConflictID(conflict.ActiveID), "nonterminal_id": activationConflictID(conflict.NonterminalID), "requested_direction": "activate"})
	}
	var version *workflow.ActivationVersionConflictError
	if errors.As(cause, &version) {
		return cli.NewFailure(cli.CodeActivationVersionConflict, "activation version does not match", &cli.Subject{Type: cli.SubjectActivation, ID: version.ActivationID}, map[string]any{"activation_id": version.ActivationID, "expected_version": version.ExpectedVersion, "actual_version": version.ActualVersion})
	}
	var incomplete *workflow.ActivationIncompleteError
	if errors.As(cause, &incomplete) {
		reasons := make([]map[string]any, 0, len(incomplete.Reasons))
		for _, reason := range incomplete.Reasons {
			reasons = append(reasons, map[string]any{"code": reason.Code, "detail": reason.Detail, "subject_id": reason.SubjectID, "subject_type": reason.SubjectType})
		}
		return cli.NewFailure(cli.CodeActivationIncomplete, "activation resource agreement is incomplete", &cli.Subject{Type: cli.SubjectActivation, ID: incomplete.ActivationID}, map[string]any{"activation_id": incomplete.ActivationID, "components": incomplete.Components, "reasons": reasons})
	}
	var direction *workflow.ActivationDirectionConflictError
	if errors.As(cause, &direction) {
		return cli.NewFailure(cli.CodeActivationDirectionConflict, "activation recovery direction is not permitted", &cli.Subject{Type: cli.SubjectActivation, ID: direction.ActivationID}, map[string]any{"activation_id": direction.ActivationID, "point_of_no_return": direction.PointOfNoReturn, "required_direction": direction.RequiredDirection})
	}
	var notFound *workflow.ActivationNotFoundError
	if errors.As(cause, &notFound) {
		return cli.NewFailure(cli.CodeActivationNotFound, "activation was not found", &cli.Subject{Type: cli.SubjectActivation, ID: notFound.ID}, map[string]any{"id": notFound.ID})
	}
	var invalid *workflow.ActivationInvalidError
	if errors.As(cause, &invalid) {
		return cli.NewFailure(cli.CodeActivationInvalid, "activation state is not permitted", &cli.Subject{Type: cli.SubjectActivation, ID: ""}, map[string]any{"status": invalid.Status, "expected_statuses": invalid.ExpectedStatuses, "component_errors": []map[string]any{}})
	}
	return genericFailure(cause, cli.CodeInvalidArgument, "invalid activation input", nil, map[string]any{"argument": "input", "value": nil, "allowed": []string{"accepted activation input"}})
}

func activationConflictID(value *string) any {
	if value == nil {
		return nil
	}
	return *value
}
