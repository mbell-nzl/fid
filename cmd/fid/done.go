package main

import (
	"errors"
	"fmt"

	"github.com/mbell-nzl/fid/internal/cli"
	"github.com/mbell-nzl/fid/internal/donework"
	"github.com/mbell-nzl/fid/internal/github"
	"github.com/spf13/cobra"
)

func init() { rootFlagCapabilities["done"] = rootFlagCapability{Project: true, JSON: true} }
func newDoneCommand(run github.Runner) *cobra.Command {
	var actor, key string
	cmd := &cobra.Command{Use: "done <task>", Short: "Complete a reviewed task after its merge and CI pass.", Args: canonicalExactArgs("task.done", 1), RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireCanonicalFlag(cmd, "task.done", "actor"); err != nil {
			return err
		}
		root, err := rebuildWorkspaceRoot(cmd)
		if err != nil {
			return renderCheckOperationFailure(cmd, "task.done", err)
		}
		result, err := donework.Finish(cmd.Context(), donework.Options{WorkspaceRoot: root, MigrationsDir: issueMigrationsDir(), TaskID: args[0], Actor: actor, IdempotencyKey: key, Clock: commandDisplayEnvironment(cmd).now, Runner: run})
		if err != nil {
			return renderDoneFailure(cmd, args[0], err)
		}
		jsonMode, err := readJSONMode(cmd)
		if err != nil {
			return err
		}
		if jsonMode {
			return cli.WriteSuccess(cmd.OutOrStdout(), "task.done", map[string]any{"status": "done", "completion": result.Completion}, &cli.Event{ID: result.EventID, OccurredAt: result.Completion.At})
		}
		_, err = fmt.Fprintf(cmd.OutOrStdout(), "✓ Completed %s · merged with CI green\n", cli.EscapeControlText(result.Completion.TaskID))
		return err
	}}
	cmd.Flags().StringVar(&actor, "actor", "", "Actor attribution")
	cmd.Flags().StringVar(&key, "idempotency-key", "", "Retry key")
	return cmd
}
func renderDoneFailure(cmd *cobra.Command, task string, cause error) error {
	var r *donework.Refusal
	if !errors.As(cause, &r) {
		return renderCheckOperationFailure(cmd, "task.done", cause)
	}
	var subject *cli.Subject
	if r.Code != "invalid-argument" {
		subject = &cli.Subject{Type: cli.SubjectIssue, ID: task}
	}
	f, err := cli.NewFailure(cli.Code(r.Code), cli.BoundedMessage(r.Message), subject, r.Details)
	if err != nil {
		return err
	}
	f.NextAction = r.Next
	jsonMode, err := readJSONMode(cmd)
	if err != nil {
		return err
	}
	if jsonMode {
		return cli.WriteFailureAndReturn(cmd.ErrOrStderr(), "task.done", f, cause)
	}
	if err := cli.WriteHumanFailure(cmd.ErrOrStderr(), "task.done", cmd.CommandPath(), f, false); err != nil {
		return err
	}
	return cli.NewFailureError(f, cause)
}
