package main

import (
	"fmt"

	"github.com/mbell-nzl/fid/internal/cli"
	"github.com/mbell-nzl/fid/internal/holdwork"
	fidmigrations "github.com/mbell-nzl/fid/migrations"
	"github.com/spf13/cobra"
)

func init() { rootFlagCapabilities["rework allow"] = rootFlagCapability{Project: true, JSON: true} }

func newReworkCommand() *cobra.Command {
	parent := &cobra.Command{Use: "rework", Short: "Record a person's rework allowance."}
	var actor, reason string
	var rounds int64
	cmd := &cobra.Command{Use: "allow <task>", Short: "Allow more rework rounds for one task.", RunE: func(cmd *cobra.Command, args []string) error {
		const operation = "task.rework-allow"
		for _, flag := range []string{"actor", "reason", "rounds"} {
			if err := requireCanonicalFlag(cmd, operation, flag); err != nil {
				return err
			}
		}
		root, err := rebuildWorkspaceRoot(cmd)
		if err != nil {
			return renderHoldFailure(cmd, operation, err)
		}
		result, err := holdwork.AllowRework(holdwork.Options{WorkspaceRoot: root, MigrationsFS: fidmigrations.FS, TaskID: args[0], Actor: actor, Reason: reason, Rounds: rounds, Clock: commandDisplayEnvironment(cmd).now})
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
		_, err = fmt.Fprintf(cmd.OutOrStdout(), "✓ Allowed %d more rework rounds for %s · %s\n", rounds, cli.EscapeControlText(args[0]), cli.EscapeControlText(reason))
		return err
	}}
	cmd.Flags().StringVar(&actor, "actor", "", "Person allowing more rounds")
	cmd.Flags().StringVar(&reason, "reason", "", "Guidance for the builder")
	cmd.Flags().Int64Var(&rounds, "rounds", 0, "Additional rounds (at least 1)")
	parent.AddCommand(cmd)
	return parent
}
