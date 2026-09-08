package main

import (
	"errors"
	"fmt"

	"github.com/mbell-nzl/fid/internal/cli"
	"github.com/mbell-nzl/fid/internal/workflow"
	"github.com/spf13/cobra"
)

func newRunCheckpointCommand() *cobra.Command {
	var inputPath, actor, key string
	cmd := &cobra.Command{Use: "checkpoint", Short: "Record a trusted Run checkpoint.", Args: canonicalNoArgs("run.checkpoint"), RunE: func(cmd *cobra.Command, _ []string) error {
		var request workflow.RunCheckpointRequest
		if err := runRequest(cmd, "run.checkpoint", inputPath, actor, key, &request); err != nil {
			return err
		}
		workspace, err := issueMutationWorkspaceRoot(cmd)
		if err != nil {
			return renderRunRecoveryError(cmd, "run.checkpoint", err, request.RunID, request.ReconciliationDigest)
		}
		result, err := workflow.RunCheckpoint(workflow.RunOptions{WorkspaceRoot: workspace, MigrationsDir: issueMigrationsDir(), ActorID: actor, IdempotencyKey: key}, request)
		if err != nil {
			return renderRunRecoveryError(cmd, "run.checkpoint", err, request.RunID, request.ReconciliationDigest)
		}
		if jsonMode, err := readJSONMode(cmd); err != nil {
			return err
		} else if jsonMode {
			return cli.WriteSuccess(cmd.OutOrStdout(), "run.checkpoint", map[string]any{"checkpoint": result.Checkpoint, "run": result.Run, "idempotency_key": result.IdempotencyKey, "replayed": result.Replayed}, &cli.Event{ID: result.EventID, OccurredAt: result.Checkpoint.ObservedAt})
		}
		_, err = fmt.Fprintf(cmd.OutOrStdout(), "checkpointed Run %s\n", result.Run.ID)
		return err
	}}
	runMutationFlags(cmd, &inputPath, &actor, &key)
	return cmd
}

func newRunInspectCommand() *cobra.Command {
	var repository string
	cmd := &cobra.Command{Use: "inspect <run-id>", Short: "Inspect one Run's current reconciliation.", Args: canonicalExactArgs("run.inspect", 1), RunE: func(cmd *cobra.Command, args []string) error {
		workspace, err := issueMutationWorkspaceRoot(cmd)
		if err != nil {
			return renderRunRecoveryError(cmd, "run.inspect", err, args[0], "")
		}
		if repository == "" {
			repository, err = commandWorkingDirectory(cmd)
			if err != nil {
				return renderRunError(cmd, "run.inspect", err)
			}
		}
		result, err := workflow.InspectRun(workspace, issueMigrationsDir(), args[0], repository)
		if err != nil {
			return renderRunRecoveryError(cmd, "run.inspect", err, args[0], "")
		}
		if jsonMode, err := readJSONMode(cmd); err != nil {
			return err
		} else if jsonMode {
			return cli.WriteSuccess(cmd.OutOrStdout(), "run.inspect", map[string]any{"run": result.Run, "reconciliation": result.Reconciliation, "liveness": result.Liveness}, nil)
		}
		_, err = fmt.Fprintf(cmd.OutOrStdout(), "%s\t%s\n", result.Run.ID, result.Reconciliation.Digest)
		return err
	}}
	cmd.Flags().StringVar(&repository, "repository", "", "Repository worktree path")
	cmd.SetFlagErrorFunc(canonicalFlagError("run.inspect"))
	return cmd
}

func newRunRecoverCommand() *cobra.Command {
	var inputPath, actor, key string
	cmd := &cobra.Command{Use: "recover", Short: "Create one linked successor Run.", Args: canonicalNoArgs("run.recover"), RunE: func(cmd *cobra.Command, _ []string) error {
		var request workflow.RunRecoveryRequest
		if err := runRequest(cmd, "run.recover", inputPath, actor, key, &request); err != nil {
			return err
		}
		workspace, err := issueMutationWorkspaceRoot(cmd)
		if err != nil {
			return renderRunRecoveryError(cmd, "run.recover", err, request.PriorRunID, request.ReconciliationDigest)
		}
		result, err := workflow.RecoverRun(workflow.RunOptions{WorkspaceRoot: workspace, MigrationsDir: issueMigrationsDir(), ActorID: actor, IdempotencyKey: key}, request)
		if err != nil {
			return renderRunRecoveryError(cmd, "run.recover", err, request.PriorRunID, request.ReconciliationDigest)
		}
		if jsonMode, err := readJSONMode(cmd); err != nil {
			return err
		} else if jsonMode {
			return cli.WriteSuccess(cmd.OutOrStdout(), "run.recover", map[string]any{"recovery": result.Recovery, "run": result.Run, "claim": result.Claim, "idempotency_key": result.IdempotencyKey, "replayed": result.Replayed}, &cli.Event{ID: result.EventID, OccurredAt: result.Run.StartedAt})
		}
		return writeRunRecoveryHuman(cmd, result)
	}}
	runMutationFlags(cmd, &inputPath, &actor, &key)
	return cmd
}

func renderRunRecoveryError(cmd *cobra.Command, operation string, cause error, runID, digest string) error {
	jsonMode, err := readJSONMode(cmd)
	if err != nil {
		return err
	}
	var duplicate *workflow.DuplicateRunExecutorError
	if errors.As(cause, &duplicate) {
		return writeRunFailure(cmd, jsonMode, operation, cause, cli.CodeDuplicateEntity, "executor identity already exists", &cli.Subject{Type: cli.SubjectRun, ID: runID}, map[string]any{"entity_type": "run", "identity": map[string]any{"fields": []string{"executor_id"}, "values": []string{duplicate.ExecutorID}}})
	}
	var conflict *workflow.RecoveryStateConflictError
	if errors.As(cause, &conflict) {
		return writeRunFailure(cmd, jsonMode, operation, cause, cli.CodeRecoveryStateConflict, "recovery state conflict", &cli.Subject{Type: cli.SubjectRun, ID: runID}, map[string]any{"run_id": runID, "observation_digest": nullableRunDigest(digest), "paths": append([]string{}, conflict.Paths...), "reasons": conflict.Reasons})
	}
	var unknown *workflow.ExecutorLivenessUnknownError
	if errors.As(cause, &unknown) {
		return writeRunFailure(cmd, jsonMode, operation, cause, cli.CodeExecutorLivenessUnknown, "executor liveness is unknown", &cli.Subject{Type: cli.SubjectRun, ID: unknown.RunID}, map[string]any{"run_id": unknown.RunID, "executor_id": unknown.ExecutorID, "observation_id": nullableRunDigest(unknown.ObservationID)})
	}
	var live *workflow.ExecutorStillLiveError
	if errors.As(cause, &live) {
		return writeRunFailure(cmd, jsonMode, operation, cause, cli.CodeExecutorStillLive, "executor is still live", &cli.Subject{Type: cli.SubjectRun, ID: live.RunID}, map[string]any{"run_id": live.RunID, "executor_id": live.ExecutorID, "observation_id": live.ObservationID})
	}
	var missing *workflow.ObservationNotFoundError
	if errors.As(cause, &missing) {
		return writeRunFailure(cmd, jsonMode, operation, cause, cli.CodeObservationNotFound, "observation was not found", &cli.Subject{Type: cli.SubjectRun, ID: runID}, map[string]any{"id": missing.ID})
	}
	var linked *workflow.RecoveryAlreadyLinkedError
	if errors.As(cause, &linked) {
		return writeRunFailure(cmd, jsonMode, operation, cause, cli.CodeRecoveryAlreadyLinked, "recovery already linked", &cli.Subject{Type: cli.SubjectRun, ID: linked.RunID}, map[string]any{"run_id": linked.RunID, "successor_run_id": linked.SuccessorRunID})
	}
	return renderRunError(cmd, operation, cause)
}

func nullableRunDigest(value string) any {
	if value == "" {
		return nil
	}
	return value
}
