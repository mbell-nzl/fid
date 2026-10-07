package main

import (
	"errors"
	"fmt"
	"strings"

	"github.com/mbell-nzl/fid/internal/cli"
	"github.com/mbell-nzl/fid/internal/journal"
	"github.com/mbell-nzl/fid/internal/policy"
	"github.com/mbell-nzl/fid/internal/review"
	"github.com/mbell-nzl/fid/internal/workflow"
	"github.com/spf13/cobra"
)

func init() {
	for _, operation := range []string{"run start", "run close", "run abandon", "run observe", "run observe-budget", "run checkpoint", "run recover", "run inspect", "run get", "run list"} {
		rootFlagCapabilities[operation] = rootFlagCapability{Project: true, JSON: true}
	}
}

func newRunCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "run",
		Short: "Manage generic Runs.",
		RunE: func(cmd *cobra.Command, _ []string) error {
			return cmd.Help()
		},
	}
	cmd.AddCommand(newRunStartCommand(), newRunCloseCommand(), newRunAbandonCommand(), newRunObserveCommand(), newRunObserveBudgetCommand(), newRunCheckpointCommand(), newRunInspectCommand(), newRunRecoverCommand(), newRunGetCommand(), newRunListCommand())
	return cmd
}

func newRunObserveCommand() *cobra.Command {
	var inputPath, actor, key string
	cmd := &cobra.Command{Use: "observe", Short: "Record a trusted Run liveness observation.", Args: canonicalNoArgs("run.observe"), RunE: func(cmd *cobra.Command, _ []string) error {
		var request workflow.RunObserveRequest
		if err := runRequest(cmd, "run.observe", inputPath, actor, key, &request); err != nil {
			return err
		}
		workspace, err := issueMutationWorkspaceRoot(cmd)
		if err != nil {
			return renderRunError(cmd, "run.observe", err)
		}
		result, err := workflow.ObserveRun(workflow.RunOptions{WorkspaceRoot: workspace, MigrationsDir: issueMigrationsDir(), ActorID: actor, IdempotencyKey: key}, request)
		if err != nil {
			return renderRunError(cmd, "run.observe", err)
		}
		return writeRunObservationSuccess(cmd, result)
	}}
	runMutationFlags(cmd, &inputPath, &actor, &key)
	return cmd
}

func newRunObserveBudgetCommand() *cobra.Command {
	var inputPath, actor, key string
	cmd := &cobra.Command{Use: "observe-budget", Short: "Record a trusted external budget observation.", Args: canonicalNoArgs("run.observe-budget"), RunE: func(cmd *cobra.Command, _ []string) error {
		var request workflow.RunBudgetObserveRequest
		if err := runRequest(cmd, "run.observe-budget", inputPath, actor, key, &request); err != nil {
			return err
		}
		workspace, err := issueMutationWorkspaceRoot(cmd)
		if err != nil {
			return renderRunError(cmd, "run.observe-budget", err)
		}
		result, err := workflow.ObserveRunBudget(workflow.RunOptions{WorkspaceRoot: workspace, MigrationsDir: issueMigrationsDir(), ActorID: actor, IdempotencyKey: key}, request)
		if err != nil {
			return renderBudgetSampleError(cmd, err, request.RunID)
		}
		if jsonMode, err := readJSONMode(cmd); err != nil {
			return err
		} else if jsonMode {
			return cli.WriteSuccess(cmd.OutOrStdout(), "run.observe-budget", map[string]any{"sample": result.Sample, "budget": result.Budget, "stops": result.Stops, "attention": result.Attention, "idempotency_key": result.IdempotencyKey, "replayed": result.Replayed}, &cli.Event{ID: result.EventID, OccurredAt: result.Sample.ObservedAt})
		}
		_, err = fmt.Fprintf(cmd.OutOrStdout(), "observed budget %s\n", result.Budget.Name)
		return err
	}}
	runMutationFlags(cmd, &inputPath, &actor, &key)
	return cmd
}

func newRunStartCommand() *cobra.Command {
	var inputPath, actor, key string
	cmd := &cobra.Command{Use: "start", Short: "Start one Run from an active Claim.", Args: canonicalNoArgs("run.start"), RunE: func(cmd *cobra.Command, _ []string) error {
		var request workflow.RunStartRequest
		if err := runRequest(cmd, "run.start", inputPath, actor, key, &request); err != nil {
			return err
		}
		workspace, err := issueMutationWorkspaceRoot(cmd)
		if err != nil {
			return renderRunError(cmd, "run.start", err)
		}
		result, err := workflow.StartRun(workflow.RunOptions{WorkspaceRoot: workspace, MigrationsDir: issueMigrationsDir(), ActorID: actor, IdempotencyKey: key}, request)
		if err != nil {
			return renderRunError(cmd, "run.start", err)
		}
		return writeRunSuccess(cmd, "run.start", result)
	}}
	runMutationFlags(cmd, &inputPath, &actor, &key)
	return cmd
}

func newRunCloseCommand() *cobra.Command {
	var inputPath, actor, key string
	cmd := &cobra.Command{Use: "close", Short: "Close one Run from trusted observations and Evidence.", Args: canonicalNoArgs("run.close"), RunE: func(cmd *cobra.Command, _ []string) error {
		var request workflow.RunCloseRequest
		if err := runRequest(cmd, "run.close", inputPath, actor, key, &request); err != nil {
			return err
		}
		workspace, err := issueMutationWorkspaceRoot(cmd)
		if err != nil {
			return renderRunError(cmd, "run.close", err)
		}
		result, err := workflow.CloseRun(workflow.RunOptions{WorkspaceRoot: workspace, MigrationsDir: issueMigrationsDir(), ActorID: actor, IdempotencyKey: key}, request)
		if err != nil {
			return renderRunRecoveryError(cmd, "run.close", err, request.RunID, "")
		}
		return writeRunSuccess(cmd, "run.close", result)
	}}
	runMutationFlags(cmd, &inputPath, &actor, &key)
	return cmd
}

func newRunAbandonCommand() *cobra.Command {
	var inputPath, actor, key string
	cmd := &cobra.Command{Use: "abandon", Short: "Abandon one Run with recovery authority.", Args: canonicalNoArgs("run.abandon"), RunE: func(cmd *cobra.Command, _ []string) error {
		var request workflow.RunAbandonRequest
		if err := runRequest(cmd, "run.abandon", inputPath, actor, key, &request); err != nil {
			return err
		}
		workspace, err := issueMutationWorkspaceRoot(cmd)
		if err != nil {
			return renderRunError(cmd, "run.abandon", err)
		}
		result, err := workflow.AbandonRun(workflow.RunOptions{WorkspaceRoot: workspace, MigrationsDir: issueMigrationsDir(), ActorID: actor, IdempotencyKey: key}, request)
		if err != nil {
			return renderRunRecoveryError(cmd, "run.abandon", err, request.RunID, "")
		}
		return writeRunAbandonSuccess(cmd, result)
	}}
	runMutationFlags(cmd, &inputPath, &actor, &key)
	return cmd
}

func newRunGetCommand() *cobra.Command {
	cmd := &cobra.Command{Use: "get <run-id>", Short: "Read one Run.", Args: canonicalExactArgs("run.get", 1), RunE: func(cmd *cobra.Command, args []string) error {
		workspace, err := issueMutationWorkspaceRoot(cmd)
		if err != nil {
			return renderRunError(cmd, "run.get", err)
		}
		run, err := workflow.GetRun(workspace, issueMigrationsDir(), args[0])
		if err != nil {
			return renderRunError(cmd, "run.get", err)
		}
		if jsonMode, err := readJSONMode(cmd); err != nil {
			return err
		} else if jsonMode {
			return cli.WriteSuccess(cmd.OutOrStdout(), "run.get", map[string]any{"run": run}, nil)
		}
		_, err = fmt.Fprintf(cmd.OutOrStdout(), "%s\t%s\n", run.ID, run.Status)
		return err
	}}
	cmd.SetFlagErrorFunc(canonicalFlagError("run.get"))
	return cmd
}

func newRunListCommand() *cobra.Command {
	var issueID, envelopeID, status string
	var limit int
	cmd := &cobra.Command{Use: "list", Short: "List generic Runs.", Args: canonicalNoArgs("run.list"), RunE: func(cmd *cobra.Command, _ []string) error {
		workspace, err := issueMutationWorkspaceRoot(cmd)
		if err != nil {
			return renderRunError(cmd, "run.list", err)
		}
		items, err := workflow.ListRuns(workflow.RunListOptions{WorkspaceRoot: workspace, MigrationsDir: issueMigrationsDir(), IssueID: issueID, EnvelopeID: envelopeID, Status: status, Limit: limit})
		if err != nil {
			return renderRunError(cmd, "run.list", err)
		}
		if jsonMode, err := readJSONMode(cmd); err != nil {
			return err
		} else if jsonMode {
			return cli.WriteSuccess(cmd.OutOrStdout(), "run.list", map[string]any{"items": items}, nil)
		}
		for _, run := range items {
			if _, err := fmt.Fprintf(cmd.OutOrStdout(), "%s\t%s\n", run.ID, run.Status); err != nil {
				return err
			}
		}
		return nil
	}}
	cmd.Flags().StringVar(&issueID, "issue", "", "Issue ID")
	cmd.Flags().StringVar(&envelopeID, "envelope", "", "Envelope ID")
	cmd.Flags().StringVar(&status, "status", "", "Run status")
	cmd.Flags().IntVar(&limit, "limit", 100, "Maximum Runs")
	cmd.SetFlagErrorFunc(canonicalFlagError("run.list"))
	return cmd
}

func runMutationFlags(cmd *cobra.Command, inputPath, actor, key *string) {
	cmd.Flags().StringVar(inputPath, "input", "", "JSON input file (or - for stdin)")
	cmd.Flags().StringVar(actor, "actor", "", "Finalized activation actor ID")
	cmd.Flags().StringVar(key, "idempotency-key", "", "Idempotency key")
	if descriptor, ok := inputDescriptorRegistry["run "+cmd.Name()]; ok {
		installInputDescriptor(cmd, descriptor)
	}
	cmd.SetFlagErrorFunc(canonicalFlagError(strings.ReplaceAll(rootCommandKey(cmd), " ", ".")))
}

func runRequest(cmd *cobra.Command, operation, inputPath, actor, key string, target any) error {
	for _, flag := range []string{"input", "actor", "idempotency-key"} {
		if err := requireCanonicalFlag(cmd, operation, flag); err != nil {
			return err
		}
	}
	if err := readDescriptorInput(cmd.InOrStdin(), inputPath, inputDescriptorRegistry[strings.ReplaceAll(operation, ".", " ")], target); err != nil {
		return renderRunError(cmd, operation, err)
	}
	return nil
}

func writeRunSuccess(cmd *cobra.Command, operation string, result workflow.RunMutationResult) error {
	if jsonMode, err := readJSONMode(cmd); err != nil {
		return err
	} else if jsonMode {
		return cli.WriteSuccess(cmd.OutOrStdout(), operation, map[string]any{"run": result.Run, "idempotency_key": result.IdempotencyKey, "replayed": result.Replayed}, &cli.Event{ID: result.EventID, OccurredAt: runEventTime(result.Run)})
	}
	_, err := fmt.Fprintf(cmd.OutOrStdout(), "%s Run %s\n", result.Run.Status, result.Run.ID)
	return err
}

func writeRunObservationSuccess(cmd *cobra.Command, result workflow.RunObservationResult) error {
	if jsonMode, err := readJSONMode(cmd); err != nil {
		return err
	} else if jsonMode {
		return cli.WriteSuccess(cmd.OutOrStdout(), "run.observe", map[string]any{"observation": result.Observation, "run": result.Run, "idempotency_key": result.IdempotencyKey, "replayed": result.Replayed}, &cli.Event{ID: result.EventID, OccurredAt: result.Observation.ObservedAt})
	}
	_, err := fmt.Fprintf(cmd.OutOrStdout(), "observed Run %s\n", result.Run.ID)
	return err
}

func writeRunAbandonSuccess(cmd *cobra.Command, result workflow.RunAbandonResult) error {
	if jsonMode, err := readJSONMode(cmd); err != nil {
		return err
	} else if jsonMode {
		return cli.WriteSuccess(cmd.OutOrStdout(), "run.abandon", map[string]any{"run": result.Run, "claim": result.Claim, "idempotency_key": result.IdempotencyKey, "replayed": result.Replayed}, &cli.Event{ID: result.EventID, OccurredAt: runEventTime(result.Run)})
	}
	_, err := fmt.Fprintf(cmd.OutOrStdout(), "abandoned Run %s\n", result.Run.ID)
	return err
}

func runEventTime(run workflow.Run) string {
	if run.EndedAt != nil {
		return *run.EndedAt
	}
	return run.StartedAt
}

func renderRunError(cmd *cobra.Command, operation string, cause error) error {
	var stopped *workflow.StopBlockedError
	if errors.As(cause, &stopped) {
		return renderStopError(cmd, operation, cause)
	}
	jsonMode, err := readJSONMode(cmd)
	if err != nil {
		return err
	}
	var mismatch *workflow.RunRepositoryMismatchError
	if errors.As(cause, &mismatch) {
		return writeRunFailure(cmd, jsonMode, operation, cause, cli.CodeRepositoryMismatch, "repository identity does not match", nil, map[string]any{"expected": map[string]any{"head_oid": mismatch.Expected.HeadOID, "dirty": mismatch.Expected.Dirty}, "observed": repositoryDetails(review.Repository(mismatch.Observed))})
	}
	if rendered, ok := renderJudgmentObservationError(cmd, jsonMode, operation, cause); ok {
		return rendered
	}
	var idempotency *journal.IdempotencyConflictError
	if errors.As(cause, &idempotency) {
		return writeRunFailure(cmd, jsonMode, operation, cause, cli.CodeIdempotencyConflict, "idempotency key conflicts with a different request", nil, idempotencyConflictDetails(idempotency))
	}
	var notFound *workflow.RunNotFoundError
	if errors.As(cause, &notFound) {
		return writeRunFailure(cmd, jsonMode, operation, cause, cli.CodeRunNotFound, "run was not found", &cli.Subject{Type: cli.SubjectRun, ID: notFound.ID}, map[string]any{"id": notFound.ID})
	}
	var inactive *workflow.RunInactiveError
	if errors.As(cause, &inactive) {
		return writeRunFailure(cmd, jsonMode, operation, cause, cli.CodeRunInactive, "run is not active", &cli.Subject{Type: cli.SubjectRun, ID: inactive.RunID}, map[string]any{"run_id": inactive.RunID, "status": inactive.Status})
	}
	var versionConflict *workflow.RunVersionConflictError
	if errors.As(cause, &versionConflict) {
		return writeRunFailure(cmd, jsonMode, operation, cause, cli.CodeRunVersionConflict, "run version conflict", &cli.Subject{Type: cli.SubjectRun, ID: versionConflict.RunID}, map[string]any{"run_id": versionConflict.RunID, "expected_version": versionConflict.ExpectedVersion, "actual_version": versionConflict.ActualVersion})
	}
	var duplicateExecutor *workflow.DuplicateRunExecutorError
	if errors.As(cause, &duplicateExecutor) {
		return writeRunFailure(cmd, jsonMode, operation, cause, cli.CodeDuplicateEntity, "executor identity already exists", &cli.Subject{Type: cli.SubjectRun, ID: nil}, map[string]any{"entity_type": "run", "identity": map[string]any{"fields": []string{"executor_id"}, "values": []string{duplicateExecutor.ExecutorID}}})
	}
	var invalidClaimMode *workflow.ClaimModeInvalidError
	if errors.As(cause, &invalidClaimMode) {
		return writeRunFailure(cmd, jsonMode, operation, cause, cli.CodeClaimModeInvalid, "claim mode is invalid", &cli.Subject{Type: cli.SubjectClaim, ID: invalidClaimMode.ClaimID}, map[string]any{"claim_id": invalidClaimMode.ClaimID, "mode": invalidClaimMode.Mode, "required_mode": invalidClaimMode.RequiredMode})
	}
	var denied *policy.EnvelopeAuthorityDeniedError
	if errors.As(cause, &denied) {
		return writeRunFailure(cmd, jsonMode, operation, cause, cli.CodeAuthorityDenied, "active envelope authority is denied", nil, map[string]any{"operation": operation, "actor_id": denied.ActorID, "required_actions": []string{denied.Action}, "required_capabilities": []string{denied.Capability}})
	}
	return writeGenericFailure(cmd, jsonMode, operation, cause, cli.CodeInvalidArgument, "invalid Run request", nil, map[string]any{"argument": "input", "value": nil, "allowed": []string{"accepted " + operation + " input"}})
}

func writeRunFailure(cmd *cobra.Command, jsonMode bool, operation string, cause error, code cli.Code, message string, subject *cli.Subject, details map[string]any) error {
	failure, err := cli.NewFailure(code, message, subject, details)
	if err != nil {
		return err
	}
	if !jsonMode {
		return writeHumanFailure(cmd, operation, failure, cause)
	}
	return cli.WriteFailureAndReturn(cmd.ErrOrStderr(), operation, failure, fmt.Errorf("%s failed", operation))
}

func renderBudgetSampleError(cmd *cobra.Command, cause error, runID string) error {
	var conflict *workflow.BudgetSampleConflictError
	if !errors.As(cause, &conflict) {
		return renderRunError(cmd, "run.observe-budget", cause)
	}
	jsonMode, err := readJSONMode(cmd)
	if err != nil {
		return err
	}
	return writeRunFailure(cmd, jsonMode, "run.observe-budget", cause, cli.CodeBudgetSampleConflict, "budget sample conflict", &cli.Subject{Type: cli.SubjectRun, ID: runID}, map[string]any{"budget_name": conflict.BudgetName, "source_sample_id": conflict.SourceSampleID, "expected_previous_sample_id": optionalString(conflict.ExpectedPreviousSampleID), "actual_previous_sample_id": optionalString(conflict.ActualPreviousSampleID)})
}
