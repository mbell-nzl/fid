package main

import (
	"fmt"

	"github.com/mbell-nzl/fid/internal/cli"
	"github.com/mbell-nzl/fid/internal/holdwork"
	"github.com/spf13/cobra"
)

func init() { rootFlagCapabilities["accept"] = rootFlagCapability{Project: true, JSON: true} }

func newAcceptCommand() *cobra.Command {
	var actor, reason string
	cmd := &cobra.Command{Use: "accept <issue>...", Short: "Accept discovered work as a person.", Args: cobra.MinimumNArgs(1), RunE: func(cmd *cobra.Command, args []string) error {
		const operation = "task.accept"
		if err := requireCanonicalFlag(cmd, operation, "actor"); err != nil {
			return err
		}
		root, err := rebuildWorkspaceRoot(cmd)
		if err != nil {
			return renderHoldFailure(cmd, operation, err)
		}
		decisions := []map[string]any{}
		for _, id := range args {
			result, err := holdwork.Accept(holdwork.Options{WorkspaceRoot: root, MigrationsDir: issueMigrationsDir(), TaskID: id, Actor: actor, Reason: reason, Clock: commandDisplayEnvironment(cmd).now})
			if err != nil {
				return renderHoldFailure(cmd, operation, err)
			}
			decisions = append(decisions, result.Fields)
		}
		jsonMode, err := readJSONMode(cmd)
		if err != nil {
			return err
		}
		if jsonMode {
			return cli.WriteSuccess(cmd.OutOrStdout(), operation, map[string]any{"decisions": decisions}, nil)
		}
		for _, d := range decisions {
			if _, err := fmt.Fprintf(cmd.OutOrStdout(), "✓ Accepted %s · decision %s\n", cli.EscapeControlText(fmt.Sprint(d["issue_id"])), cli.EscapeControlText(fmt.Sprint(d["decision_id"]))); err != nil {
				return err
			}
		}
		return nil
	}}
	cmd.Flags().StringVar(&actor, "actor", "", "Person accepting work")
	cmd.Flags().StringVar(&reason, "reason", "", "Why the work is accepted")
	return cmd
}
