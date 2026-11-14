package main

import (
	"errors"
	"fmt"
	"strings"

	"github.com/mbell-nzl/fid/internal/cli"
	"github.com/mbell-nzl/fid/internal/display"
	"github.com/mbell-nzl/fid/internal/holdwork"
	"github.com/mbell-nzl/fid/internal/persondecision"
	"github.com/mbell-nzl/fid/internal/policy"
	"github.com/spf13/cobra"
)

func init() {
	for _, name := range []string{"heartbeat", "park", "takeover"} {
		rootFlagCapabilities[name] = rootFlagCapability{Project: true, JSON: true}
	}
}

func newHoldCommand(name string) *cobra.Command {
	var hold, actor, reason, key, waiting string
	var interrupted bool
	operation := "hold.heartbeat"
	if name == "park" {
		operation = "hold.park"
	}
	if name == "takeover" {
		operation = "task.takeover"
	}
	use := name
	if name == "takeover" {
		use += " <task>"
	}
	cmd := &cobra.Command{Use: use, Short: map[string]string{"heartbeat": "Keep this session's hold alive.", "park": "Park work as interrupted or waiting.", "takeover": "Take over a task as a person."}[name], RunE: func(cmd *cobra.Command, args []string) error {
		if name == "park" && (interrupted == cmd.Flags().Changed("waiting") || (cmd.Flags().Changed("waiting") && cmd.Flags().Changed("reason"))) {
			return writeGrammarFailure(cmd, operation, cli.CodeInvalidArgument, "park requires either --interrupted --reason <text> or --waiting <question>", map[string]any{"argument": "request", "value": nil, "allowed": []string{"--interrupted --reason <text>", "--waiting <question>"}})
		}
		required := []string{"actor"}
		if name != "takeover" {
			required = append(required, "hold")
		}
		if name != "heartbeat" && !cmd.Flags().Changed("waiting") {
			required = append(required, "reason")
		}
		if name == "park" {
			required = append(required, "idempotency-key")
			if cmd.Flags().Changed("waiting") {
				required = append(required, "waiting")
			}
		}
		for _, flag := range required {
			if err := requireCanonicalFlag(cmd, operation, flag); err != nil {
				return err
			}
		}
		root, err := rebuildWorkspaceRoot(cmd)
		if err != nil {
			return renderHoldFailure(cmd, operation, err)
		}
		options := holdwork.Options{WorkspaceRoot: root, MigrationsDir: issueMigrationsDir(), HoldID: hold, Actor: actor, Reason: reason, IdempotencyKey: key, Clock: commandDisplayEnvironment(cmd).now, Question: waiting}
		action := holdwork.Heartbeat
		if name == "park" {
			action = holdwork.Park
		}
		if name == "takeover" {
			options.TaskID = args[0]
			action = holdwork.Takeover
		}
		result, err := action(options)
		if err != nil {
			return renderHoldFailure(cmd, operation, err)
		}
		jsonMode, err := readJSONMode(cmd)
		if err != nil {
			return err
		}
		if jsonMode {
			return cli.WriteSuccess(cmd.OutOrStdout(), operation, result.Fields, &cli.Event{ID: result.EventID, OccurredAt: result.OccurredAt})
		}
		message := ""
		switch name {
		case "heartbeat":
			message = fmt.Sprintf("Held %s until %s", result.Fields["task_id"], display.Clock(result.Fields["expires_at"].(string), commandDisplayEnvironment(cmd).now()))
		case "park":
			place := "back in the queue"
			if result.Fields["kept_for_worktree"] == true {
				place = "kept for this worktree"
			}
			if result.Fields["mode"] == "waiting" {
				place = "waiting for a person's answer"
			}
			message = fmt.Sprintf("Parked %s · %s", result.Fields["task_id"], place)
		case "takeover":
			message = fmt.Sprintf("Took over %s · it's yours in this worktree", result.Fields["task_id"])
		}
		_, err = fmt.Fprintln(cmd.OutOrStdout(), "✓ "+cli.EscapeControlText(message))
		return err
	}}
	if name == "takeover" {
		cmd.Args = cobra.ExactArgs(1)
	} else {
		cmd.Args = cobra.NoArgs
		cmd.Flags().StringVar(&hold, "hold", "", "The hold ID for this session")
	}
	cmd.Flags().StringVar(&actor, "actor", "", "Actor attribution")
	if name != "heartbeat" {
		cmd.Flags().StringVar(&reason, "reason", "", "Why the work is parked or taken over")
		cmd.Flags().StringVar(&key, "idempotency-key", "", "Retry key")
	}
	if name == "park" {
		cmd.Flags().StringVar(&waiting, "waiting", "", "The question for a person")
		cmd.Flags().BoolVar(&interrupted, "interrupted", false, "Park as interrupted")
	}
	return cmd
}

func personDecisionFailure(operation string, cause error) (cli.Failure, bool, error) {
	var denied *persondecision.DeniedError
	if errors.As(cause, &denied) {
		f, e := cli.NewFailure(cli.CodeAuthorityDenied, denied.Error(), nil, map[string]any{"actor_id": denied.Actor, "operation": operation, "required_actions": []string{operation}, "required_capabilities": []string{operation}})
		return f, true, e
	}
	var required *persondecision.ReasonRequiredError
	if errors.As(cause, &required) {
		f, e := cli.NewFailure(cli.CodeInvalidArgument, "--reason is required", nil, map[string]any{"argument": "--reason", "value": nil, "allowed": []string{"nonempty text"}})
		return f, true, e
	}
	return cli.Failure{}, false, nil
}

func renderHoldFailure(cmd *cobra.Command, operation string, cause error) error {
	if refused := renderRepositoryRefusal(cmd, cause); refused != nil {
		return refused
	}
	failure, handled, err := personDecisionFailure(operation, cause)
	var source *policy.SourceError
	if errors.As(cause, &source) {
		failure, err = cli.NewFailure(cli.CodePolicyUnavailable, cli.BoundedMessage(source.Error()), nil, source.Details())
		failure.NextAction = source.NextAction()
		handled = true
	}
	var refusal *holdwork.Refusal
	if !handled && errors.As(cause, &refusal) {
		var subject *cli.Subject
		if refusal.Code == "unsupported-kind" {
			failure, err = cli.NewFailure(cli.CodeUnsupportedKind, refusal.Message, nil, map[string]any{"argument": "task", "value": refusal.Details["kind"], "allowed": []string{"task"}})
		} else {
			var id any = refusal.TaskID
			if refusal.TaskID == "" {
				id = nil
			}
			subject = &cli.Subject{Type: cli.SubjectIssue, ID: id}
			failure, err = cli.NewFailure(cli.Code(refusal.Code), refusal.Message, subject, refusal.Details)
		}
		handled = true
	}
	if !handled {
		failure, err = issueMutationFailure(cmd, operation, "", cause)
	}
	if err != nil {
		return err
	}
	jsonMode, err := readJSONMode(cmd)
	if err != nil {
		return err
	}
	if jsonMode {
		return cli.WriteFailureAndReturn(cmd.ErrOrStderr(), operation, failure, cause)
	}
	// The authority envelope retains the operation; ordinary text names only
	// the actor, just as next does.
	if failure.Code == cli.CodeAuthorityDenied {
		failure.Details = map[string]any{"actor_id": failure.Details["actor_id"]}
	}
	if err := cli.WriteHumanFailure(cmd.ErrOrStderr(), operation, strings.TrimSpace(cmd.CommandPath()), failure, false); err != nil {
		return err
	}
	return cli.NewFailureError(failure, cause)
}
