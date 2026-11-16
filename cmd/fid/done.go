package main

import (
	"errors"
	"fmt"

	"github.com/mbell-nzl/fid/internal/cli"
	"github.com/mbell-nzl/fid/internal/donework"
	"github.com/mbell-nzl/fid/internal/github"
	"github.com/mbell-nzl/fid/internal/persondecision"
	"github.com/spf13/cobra"
)

func init() { rootFlagCapabilities["done"] = rootFlagCapability{Project: true, JSON: true} }
func newDoneCommand(run github.Runner) *cobra.Command {
	var actor, key, reason string
	cmd := &cobra.Command{Use: "done [issue]", Short: "Complete a reviewed task, close an epic, or reconcile tasks in review.", Args: doneArgs, RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireCanonicalFlag(cmd, "task.done", "actor"); err != nil {
			return err
		}
		if len(args) == 0 && cmd.Flags().Changed("idempotency-key") {
			return writeGrammarFailure(cmd, "task.done", cli.CodeInvalidArgument, "--idempotency-key requires a task", map[string]any{"argument": "--idempotency-key", "value": key, "allowed": []string{"supply a task"}})
		}
		root, err := rebuildWorkspaceRoot(cmd)
		if err != nil {
			return renderCheckOperationFailure(cmd, "task.done", err)
		}
		if len(args) == 0 && cmd.Flags().Changed("reason") {
			return doneReasonUsage(cmd, reason)
		}
		if len(args) == 0 {
			return reconcileDone(cmd, donework.Options{WorkspaceRoot: root, MigrationsDir: issueMigrationsDir(), Actor: actor, Clock: commandDisplayEnvironment(cmd).now, Runner: run})
		}
		options := donework.Options{WorkspaceRoot: root, MigrationsDir: issueMigrationsDir(), TaskID: args[0], Actor: actor, IdempotencyKey: key, Reason: reason, Clock: commandDisplayEnvironment(cmd).now, Runner: run}
		epic, err := donework.EpicTarget(options)
		if err != nil {
			return renderCheckOperationFailure(cmd, "task.done", err)
		}
		if epic {
			return finishEpic(cmd, options)
		}
		if cmd.Flags().Changed("reason") {
			return doneReasonUsage(cmd, reason)
		}
		result, err := donework.Finish(cmd.Context(), options)
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
	cmd.Flags().StringVar(&reason, "reason", "", "Reason for closing an epic")
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

func doneArgs(cmd *cobra.Command, args []string) error {
	if len(args) > 1 {
		return canonicalExactArgs("task.done", 1)(cmd, args)
	}
	return nil
}

func reconcileDone(cmd *cobra.Command, options donework.Options) error {
	outcomes, err := donework.Reconcile(cmd.Context(), options)
	if err != nil {
		return renderCheckOperationFailure(cmd, "task.done", err)
	}
	jsonMode, err := readJSONMode(cmd)
	if err != nil {
		return err
	}
	if jsonMode {
		return cli.WriteSuccess(cmd.OutOrStdout(), "task.done", map[string]any{"outcomes": outcomes}, nil)
	}
	if len(outcomes) == 0 {
		_, err = fmt.Fprintln(cmd.OutOrStdout(), "No tasks in review to reconcile.")
		return err
	}
	for _, o := range outcomes {
		label, message := "Completed", "merged with CI green"
		if o.Refusal != nil {
			message = o.Refusal.Message
			label = map[string]string{"waiting": "Waiting", "refused": "Refused", "github-unavailable": "GitHub unavailable"}[o.Status]
			if gate, ok := o.Refusal.Details["gate"].(string); ok {
				message += " (" + gate + ")"
			}
			if o.Refusal.Next != "" {
				message += "; " + o.Refusal.Next
			}
		}
		if _, err := fmt.Fprintf(cmd.OutOrStdout(), "%s %s · %s\n", label, cli.EscapeControlText(o.TaskID), cli.EscapeControlText(message)); err != nil {
			return err
		}
	}
	return nil
}

func finishEpic(cmd *cobra.Command, o donework.Options) error {
	result, err := donework.CloseEpic(o)
	if err != nil {
		var required *persondecision.ReasonRequiredError
		if errors.As(err, &required) {
			return renderHoldFailure(cmd, "task.done", err)
		}
		return renderDoneFailure(cmd, o.TaskID, err)
	}
	jsonMode, err := readJSONMode(cmd)
	if err != nil {
		return err
	}
	if jsonMode {
		return cli.WriteSuccess(cmd.OutOrStdout(), "task.done", map[string]any{"status": result.Status, "decision": result.Decision, "revision": result.Revision}, &cli.Event{ID: result.EventID, OccurredAt: result.Decision.At})
	}
	_, err = fmt.Fprintf(cmd.OutOrStdout(), "✓ Closed %s · %s\n", cli.EscapeControlText(o.TaskID), cli.EscapeControlText(o.Reason))
	return err
}

func doneReasonUsage(cmd *cobra.Command, reason string) error {
	return writeGrammarFailure(cmd, "task.done", cli.CodeInvalidArgument, "--reason requires an epic", map[string]any{"argument": "--reason", "value": reason, "allowed": []string{"supply an epic"}})
}
