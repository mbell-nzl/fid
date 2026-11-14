package main

import (
	"fmt"

	"github.com/mbell-nzl/fid/internal/cli"
	"github.com/mbell-nzl/fid/internal/handoffrecord"
	"github.com/mbell-nzl/fid/internal/holdwork"
	"github.com/spf13/cobra"
)

func init() { rootFlagCapabilities["handoff"] = rootFlagCapability{Project: true, JSON: true} }

func newHandoffCommand() *cobra.Command {
	var hold, actor, pr, key string
	cmd := &cobra.Command{Use: "handoff", Short: "Hand the held task to review with its recorded checks.", Args: cobra.NoArgs, RunE: func(cmd *cobra.Command, _ []string) error {
		for _, flag := range []string{"hold", "actor", "pr", "idempotency-key"} {
			if err := requireCanonicalFlag(cmd, "task.handoff", flag); err != nil {
				return err
			}
		}
		if !handoffrecord.ValidPR(pr) {
			return writeGrammarFailure(cmd, "task.handoff", cli.CodeInvalidArgument, "--pr must be a GitHub pull request URL", map[string]any{"argument": "--pr", "value": pr, "allowed": []string{"https://github.com/<owner>/<repo>/pull/<n>"}})
		}
		root, err := rebuildWorkspaceRoot(cmd)
		if err != nil {
			return renderCheckOperationFailure(cmd, "task.handoff", err)
		}
		result, err := holdwork.Handoff(holdwork.Options{WorkspaceRoot: root, MigrationsDir: issueMigrationsDir(), HoldID: hold, Actor: actor, PR: pr, IdempotencyKey: key, Clock: commandDisplayEnvironment(cmd).now})
		if err != nil {
			return renderCheckOperationFailure(cmd, "task.handoff", err)
		}
		jsonMode, err := readJSONMode(cmd)
		if err != nil {
			return err
		}
		if jsonMode {
			return cli.WriteSuccess(cmd.OutOrStdout(), "task.handoff", result.Fields, &cli.Event{ID: result.EventID, OccurredAt: result.OccurredAt})
		}
		names := result.Fields["required_checks"].([]any)
		checks := "no checks required"
		if len(names) > 0 {
			checks = fmt.Sprintf("%d checks passed", len(names))
		}
		_, err = fmt.Fprintf(cmd.OutOrStdout(), "✓ Handed %s to review · %s · PR #%s\n", cli.EscapeControlText(result.Fields["task_id"].(string)), checks, handoffrecord.PRNumber(result.Fields["pr"].(string)))
		return err
	}}
	cmd.Flags().StringVar(&hold, "hold", "", "The hold ID for this session")
	cmd.Flags().StringVar(&actor, "actor", "", "Actor attribution")
	cmd.Flags().StringVar(&pr, "pr", "", "GitHub pull request URL")
	cmd.Flags().StringVar(&key, "idempotency-key", "", "Retry key")
	return cmd
}
