package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/mbell-nzl/fid/internal/cli"
	"github.com/mbell-nzl/fid/internal/policy"
	"github.com/mbell-nzl/fid/internal/workflow"
	"github.com/spf13/cobra"
)

func init() {
	for _, operation := range []string{"run start", "run finish", "run observe", "run checkpoint", "run recover", "run inspect", "run get", "run list"} {
		rootFlagCapabilities[operation] = rootFlagCapability{Project: true, JSON: true}
	}
}

func newRunCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "run",
		Short: "Manage generic Runs.",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return nil
			}
			return fmt.Errorf("unknown command %q for %q", args[0], cmd.CommandPath())
		},
		RunE: func(cmd *cobra.Command, _ []string) error {
			return cmd.Help()
		},
	}
	cmd.AddCommand(newRunStartCommand(), newRunFinishCommand(), newRunObserveCommand(), newRunCheckpointCommand(), newRunInspectCommand(), newRunRecoverCommand(), newRunGetCommand(), newRunListCommand())
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

func newRunFinishCommand() *cobra.Command {
	var inputPath, actor, key string
	cmd := &cobra.Command{Use: "finish", Short: "Terminalize one active Run.", Args: canonicalNoArgs("run.finish"), RunE: func(cmd *cobra.Command, _ []string) error {
		var request workflow.RunFinishRequest
		if err := runRequest(cmd, "run.finish", inputPath, actor, key, &request); err != nil {
			return err
		}
		workspace, err := issueMutationWorkspaceRoot(cmd)
		if err != nil {
			return renderRunError(cmd, "run.finish", err)
		}
		result, err := workflow.FinishRun(workflow.RunOptions{WorkspaceRoot: workspace, MigrationsDir: issueMigrationsDir(), ActorID: actor, IdempotencyKey: key}, request)
		if err != nil {
			return renderRunError(cmd, "run.finish", err)
		}
		return writeRunSuccess(cmd, "run.finish", result)
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
	cmd.Flags().StringVar(inputPath, "input", "", "Run JSON input")
	cmd.Flags().StringVar(actor, "actor", "", "Finalized activation actor ID")
	cmd.Flags().StringVar(key, "idempotency-key", "", "Idempotency key")
	cmd.SetFlagErrorFunc(canonicalFlagError(strings.ReplaceAll(rootCommandKey(cmd), " ", ".")))
}

func runRequest(cmd *cobra.Command, operation, inputPath, actor, key string, target any) error {
	for _, flag := range []string{"input", "actor", "idempotency-key"} {
		if err := requireCanonicalFlag(cmd, operation, flag); err != nil {
			return err
		}
	}
	if err := decodeClosedRunInput(cmd.InOrStdin(), inputPath, target); err != nil {
		return renderRunError(cmd, operation, err)
	}
	return nil
}

func decodeClosedRunInput(stdin io.Reader, path string, target any) error {
	var payload []byte
	var err error
	if path == "-" {
		payload, err = io.ReadAll(io.LimitReader(stdin, 1_048_577))
	} else {
		payload, err = os.ReadFile(path)
	}
	if err != nil {
		return fmt.Errorf("read run input: %w", err)
	}
	if len(payload) > 1_048_576 {
		return fmt.Errorf("run input exceeds 1048576 bytes")
	}
	decoder := json.NewDecoder(strings.NewReader(string(payload)))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(target); err != nil {
		return fmt.Errorf("decode run input: %w", err)
	}
	var extra any
	if err := decoder.Decode(&extra); err != io.EOF {
		return fmt.Errorf("run input must contain one object")
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

func runEventTime(run workflow.Run) string {
	if run.EndedAt != nil {
		return *run.EndedAt
	}
	return run.StartedAt
}

func renderRunError(cmd *cobra.Command, operation string, cause error) error {
	jsonMode, err := readJSONMode(cmd)
	if err != nil {
		return err
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
	var denied *policy.EnvelopeAuthorityDeniedError
	if errors.As(cause, &denied) {
		return writeRunFailure(cmd, jsonMode, operation, cause, cli.CodeAuthorityDenied, "active envelope authority is denied", nil, map[string]any{"operation": operation, "actor": denied.ActorID})
	}
	return writeRunFailure(cmd, jsonMode, operation, cause, cli.CodeInvalidArgument, "invalid Run request", nil, map[string]any{"argument": "input", "value": nil, "allowed": []string{"accepted " + operation + " input"}})
}

func writeRunFailure(cmd *cobra.Command, jsonMode bool, operation string, cause error, code cli.Code, message string, subject *cli.Subject, details map[string]any) error {
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
