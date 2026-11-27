package main

import (
	"fmt"
	"strings"

	"github.com/mbell-nzl/fid/internal/persondecision"
	"github.com/mbell-nzl/fid/internal/policy"

	"github.com/mbell-nzl/fid/internal/cli"
	"github.com/mbell-nzl/fid/internal/mergehold"
	fidmigrations "github.com/mbell-nzl/fid/migrations"
	"github.com/spf13/cobra"
)

func init() { rootFlagCapabilities["hold"] = rootFlagCapability{Project: true, JSON: true} }
func newMergeHoldCommand() *cobra.Command {
	var actor, reason, key string
	cmd := &cobra.Command{Use: "hold <task>", Short: "Record a person's HOLD call on a hand-off.", Args: canonicalExactArgs("task.hold", 1), RunE: func(cmd *cobra.Command, args []string) error {
		for _, flag := range []string{"actor", "reason"} {
			if err := requireCanonicalFlag(cmd, "task.hold", flag); err != nil {
				return err
			}
		}
		root, err := rebuildWorkspaceRoot(cmd)
		if err != nil {
			return renderHoldFailure(cmd, "task.hold", err)
		}
		if !policy.Allows(actor, policy.PersonDecisions) {
			return renderHoldFailure(cmd, "task.hold", &persondecision.DeniedError{Actor: actor})
		}
		if strings.TrimSpace(reason) == "" {
			return renderHoldFailure(cmd, "task.hold", &persondecision.ReasonRequiredError{})
		}
		result, err := mergehold.Record(mergehold.Options{WorkspaceRoot: root, MigrationsFS: fidmigrations.FS, TaskID: args[0], Actor: actor, Reason: reason, IdempotencyKey: key, Clock: commandDisplayEnvironment(cmd).now})
		if err != nil {
			return reviewFailure(cmd, "task.hold", err, "")
		}
		postRevealedJudges(cmd, root, args[0])
		jsonMode, err := readJSONMode(cmd)
		if err != nil {
			return err
		}
		if jsonMode {
			return cli.WriteSuccess(cmd.OutOrStdout(), "task.hold", result.Fields, &cli.Event{ID: result.EventID, OccurredAt: result.Fields["at"].(string)})
		}
		_, err = fmt.Fprintf(cmd.OutOrStdout(), "✓ Held %s · %s\n", cli.EscapeControlText(args[0]), cli.EscapeControlText(reason))
		return err
	}}
	cmd.Flags().StringVar(&actor, "actor", "", "Person making the call")
	cmd.Flags().StringVar(&reason, "reason", "", "Why the task should wait")
	cmd.Flags().StringVar(&key, "idempotency-key", "", "Retry key")
	return cmd
}
