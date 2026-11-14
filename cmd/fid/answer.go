package main

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/mbell-nzl/fid/internal/cli"
	"github.com/mbell-nzl/fid/internal/holdwork"
	"github.com/spf13/cobra"
)

func init() { rootFlagCapabilities["answer"] = rootFlagCapability{Project: true, JSON: true} }

func newAnswerCommand() *cobra.Command {
	var actor, key string
	cmd := &cobra.Command{Use: "answer <task> <text>", Short: "Answer a task's open question as a person.", Args: cobra.ExactArgs(2), RunE: func(cmd *cobra.Command, args []string) error {
		const operation = "task.answer"
		if err := requireCanonicalFlag(cmd, operation, "actor"); err != nil {
			return err
		}
		if strings.TrimSpace(args[1]) == "" || !utf8.ValidString(args[1]) || len(args[1]) > 4096 {
			return writeGrammarFailure(cmd, operation, cli.CodeInvalidArgument, "answer requires nonempty text of at most 4096 bytes", map[string]any{"argument": "text", "value": nil, "allowed": []string{"nonempty text of at most 4096 bytes"}})
		}
		root, err := rebuildWorkspaceRoot(cmd)
		if err != nil {
			return renderHoldFailure(cmd, operation, err)
		}
		result, err := holdwork.Answer(holdwork.Options{WorkspaceRoot: root, MigrationsDir: issueMigrationsDir(), TaskID: args[0], Actor: actor, Answer: args[1], IdempotencyKey: key, Clock: commandDisplayEnvironment(cmd).now})
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

		_, err = fmt.Fprintln(cmd.OutOrStdout(), "✓ "+cli.EscapeControlText(fmt.Sprintf("Answered %s · it resumes in %s", args[0], result.Worktree)))
		return err
	}}
	cmd.Flags().StringVar(&actor, "actor", "", "Actor attribution")
	cmd.Flags().StringVar(&key, "idempotency-key", "", "Retry key")
	return cmd
}
