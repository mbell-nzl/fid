package main

import (
	"fmt"

	"github.com/mbell-nzl/fid/internal/cli"
	"github.com/mbell-nzl/fid/internal/holdwork"
	"github.com/spf13/cobra"
)

func init() { rootFlagCapabilities["finding close"] = rootFlagCapability{Project: true, JSON: true} }

func newFindingCommand() *cobra.Command {
	parent := &cobra.Command{Use: "finding", Short: "Record a person's finding decision."}
	var actor, reason, key string
	cmd := &cobra.Command{Use: "close <finding>", Short: "Close a finding as a person, with a reason.", Args: canonicalExactArgs("finding.close", 1), RunE: func(cmd *cobra.Command, args []string) error {
		for _, flag := range []string{"actor", "reason"} {
			if err := requireCanonicalFlag(cmd, "finding.close", flag); err != nil {
				return err
			}
		}
		root, err := rebuildWorkspaceRoot(cmd)
		if err != nil {
			return renderHoldFailure(cmd, "finding.close", err)
		}
		result, err := holdwork.CloseFinding(holdwork.Options{WorkspaceRoot: root, MigrationsDir: issueMigrationsDir(), FindingID: args[0], Actor: actor, Reason: reason, IdempotencyKey: key, Clock: commandDisplayEnvironment(cmd).now})
		if err != nil {
			return renderHoldFailure(cmd, "finding.close", err)
		}
		jsonMode, err := readJSONMode(cmd)
		if err != nil {
			return err
		}
		if jsonMode {
			return cli.WriteSuccess(cmd.OutOrStdout(), "finding.close", result.Fields, &cli.Event{ID: result.EventID, OccurredAt: result.OccurredAt})
		}
		_, err = fmt.Fprintf(cmd.OutOrStdout(), "✓ Closed %s · decision %s · %s\n", cli.EscapeControlText(args[0]), cli.EscapeControlText(fmt.Sprint(result.Fields["decision_id"])), cli.EscapeControlText(reason))
		return err
	}}
	cmd.Flags().StringVar(&actor, "actor", "", "Person making the decision")
	cmd.Flags().StringVar(&reason, "reason", "", "Why the finding is closed")
	cmd.Flags().StringVar(&key, "idempotency-key", "", "Retry key")
	parent.AddCommand(cmd)
	return parent
}
