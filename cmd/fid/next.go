package main

import (
	"errors"
	"fmt"
	"regexp"
	"time"

	"github.com/mbell-nzl/fid/internal/cli"
	"github.com/mbell-nzl/fid/internal/display"
	"github.com/mbell-nzl/fid/internal/next"
	"github.com/mbell-nzl/fid/internal/project"
	"github.com/spf13/cobra"
)

func init() {
	rootFlagCapabilities["next"] = rootFlagCapability{Project: true, JSON: true}
}

// nextActorPattern is a well-formed actor, as every other --actor takes it.
// fid next checks it itself, once it knows whether --claim was given.
const nextActorPattern = `^(agent|human|service):[A-Za-z0-9][A-Za-z0-9._-]{0,127}$`

var nextActor = regexp.MustCompile(nextActorPattern)

// nextActorDomain lets any --actor through the parser, so fid next can
// refuse it under the operation the request names.
var nextActorDomain = flagDomain{[]string{"string"}, "RESET 1.4a next.claim; P21 §2 standalone Actor", func(string) bool { return true }}

func newNextCommand() *cobra.Command {
	var (
		claim, explain bool
		actor          string
	)
	cmd := &cobra.Command{
		Use:   "next",
		Short: "Show the task to work on next in this worktree, or take it.",
		Long: "Show the task fid would hand this worktree next: the one it already holds, then work others wait on, " +
			"then unstarted tasks in fid status's order, passing over held, blocked and out-of-sequence tasks.\n\n" +
			"--claim --actor <actor> takes the task as a hold for this worktree, in one Event; run again, it returns the " +
			"same hold. --explain lists the whole order with a verdict for each task.",
		RunE: func(cmd *cobra.Command, _ []string) error {
			operation := leafOperation(cmd)
			if refused := refuseNextRequest(cmd, claim, actor); refused != nil {
				return refused
			}
			root, err := rebuildWorkspaceRoot(cmd)
			if err != nil {
				return renderNextFailure(cmd, operation, err)
			}
			environment := commandDisplayEnvironment(cmd)
			options := next.Options{WorkspaceRoot: root, MigrationsDir: issueMigrationsDir(), Clock: environment.now, Actor: actor}
			read := next.Read
			if claim {
				read = next.Claim
			}
			result, err := read(options)
			if err != nil {
				return renderNextFailure(cmd, operation, err)
			}
			jsonMode, err := readJSONMode(cmd)
			if err != nil {
				return err
			}
			if jsonMode {
				var event *cli.Event
				if result.Event != nil {
					event = &cli.Event{ID: result.Event.ID, OccurredAt: result.Event.OccurredAt.UTC().Format(time.RFC3339Nano)}
				}
				return cli.WriteSuccess(cmd.OutOrStdout(), operation, result.Fields(), event)
			}
			renderer := display.NewRenderer(display.IsTerminal(cmd.OutOrStdout()), environmentValue(environment.environ, "NO_COLOR"))
			return next.WriteText(cmd.OutOrStdout(), result, renderer, explain)
		},
	}
	cmd.Flags().BoolVar(&claim, "claim", false, "Take the task as a hold for this worktree")
	cmd.Flags().StringVar(&actor, "actor", "", "Who takes the task, agent:<id> or human:<id>; only with --claim")
	cmd.Flags().BoolVar(&explain, "explain", false, "List the whole order with a verdict for each task")
	return cmd
}

// refuseNextRequest refuses, before fid looks for a project, a request
// whose flags don't go together or whose actor can't take a task. A flag
// is never accepted and ignored.
func refuseNextRequest(cmd *cobra.Command, claim bool, actor string) error {
	switch {
	case !claim && cmd.Flags().Changed("actor"):
		return writeGrammarFailure(cmd, "next.get", cli.CodeInvalidArgument, "--actor is accepted only when taking the task",
			map[string]any{"argument": flagCombinationArgument, "value": nil, "allowed": []string{"accepted next.get input"}})
	case !claim:
		return nil
	case !cmd.Flags().Changed("actor"):
		return requireCanonicalFlag(cmd, "next.claim", "actor")
	case !nextActor.MatchString(actor):
		return writeGrammarFailure(cmd, "next.claim", cli.CodeInvalidArgument, fmt.Sprintf("invalid argument %q for %q flag", actor, "--actor"),
			map[string]any{"argument": "--actor", "value": actor, "allowed": []string{nextActorPattern}})
	}
	if err := next.CheckActor(actor); err != nil {
		return refuseNextActor(cmd, actor)
	}
	return nil
}

// refuseNextActor refuses an actor whose type takes no tasks. JSON keeps
// authority-denied's details; the human text names the actor alone, as the
// operation's name holds a word fid's text never shows.
func refuseNextActor(cmd *cobra.Command, actor string) error {
	const message = "service actors can't take tasks"
	failure, err := cli.NewFailure(cli.CodeAuthorityDenied, message, nil,
		map[string]any{"actor_id": actor, "operation": "next.claim", "required_actions": []string{"next.claim"}, "required_capabilities": []string{"next.claim"}})
	if err != nil {
		return err
	}
	jsonMode, err := readJSONMode(cmd)
	if err != nil {
		return err
	}
	cause := errors.New(message)
	if jsonMode {
		return cli.WriteFailureAndReturn(cmd.ErrOrStderr(), "next.claim", failure, cause)
	}
	human := failure
	human.Details = map[string]any{"actor_id": actor}
	if err := cli.WriteHumanFailure(cmd.ErrOrStderr(), "next.claim", cmd.CommandPath(), human, false); err != nil {
		return err
	}
	return cli.NewFailureError(failure, cause)
}

func renderNextFailure(cmd *cobra.Command, operation string, cause error) error {
	if refused := renderRepositoryRefusal(cmd, cause); refused != nil {
		return refused
	}
	jsonMode, err := readJSONMode(cmd)
	if err != nil {
		return err
	}
	failure, err := nextFailure(cmd, operation, cause)
	if err != nil {
		return err
	}
	if jsonMode {
		return cli.WriteFailureAndReturn(cmd.ErrOrStderr(), operation, failure, cause)
	}
	return writeHumanFailure(cmd, operation, failure, cause)
}

func nextFailure(cmd *cobra.Command, operation string, cause error) (cli.Failure, error) {
	if errors.Is(cause, project.ErrNotFound) {
		return issueMutationFailure(cmd, operation, "", cause)
	}
	if failure, handled, err := lockUnavailableFailure(cause); err != nil || handled {
		return failure, err
	}
	if failure, handled, err := journalReplayStoppedFailure(cause); err != nil || handled {
		return failure, err
	}
	return cli.NewFailure(cli.CodeInternalError, cli.BoundedMessage("next failed: "+cause.Error()), nil, map[string]any{})
}
