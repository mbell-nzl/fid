package main

import (
	"fmt"

	"github.com/mbell-nzl/fid/internal/cli"
	"github.com/mbell-nzl/fid/internal/holdwork"
	"github.com/spf13/cobra"
)

func init() { rootFlagCapabilities["start"] = rootFlagCapability{Project: true, JSON: true} }

func newStartCommand() *cobra.Command {
	var hold, actor, key string
	cmd := &cobra.Command{Use: "start", Short: "Start the task held by this session.", RunE: func(cmd *cobra.Command, _ []string) error {
		for _, flag := range []string{"hold", "actor", "idempotency-key"} {
			if err := requireCanonicalFlag(cmd, "task.start", flag); err != nil {
				return err
			}
		}
		root, err := rebuildWorkspaceRoot(cmd)
		if err != nil {
			return renderHoldFailure(cmd, "task.start", err)
		}
		result, err := holdwork.Start(holdwork.Options{WorkspaceRoot: root, MigrationsDir: issueMigrationsDir(), HoldID: hold, Actor: actor, IdempotencyKey: key, Clock: commandDisplayEnvironment(cmd).now})
		if err != nil {
			return renderHoldFailure(cmd, "task.start", err)
		}
		jsonMode, err := readJSONMode(cmd)
		if err != nil {
			return err
		}
		if jsonMode {
			var event *cli.Event
			if result.EventID != "" {
				event = &cli.Event{ID: result.EventID, OccurredAt: result.OccurredAt}
			}
			return cli.WriteSuccess(cmd.OutOrStdout(), "task.start", result.Fields, event)
		}
		message := fmt.Sprintf("Started %s", result.Fields["task_id"])
		if result.EventID == "" {
			message = fmt.Sprintf("%s is already started", result.Fields["task_id"])
		}
		_, err = fmt.Fprintf(cmd.OutOrStdout(), "✓ %s · %s\n", cli.EscapeControlText(message), laneLabel(result.Fields["lane"].(string)))
		return err
	}}
	cmd.Flags().StringVar(&hold, "hold", "", "The hold ID for this session")
	cmd.Flags().StringVar(&actor, "actor", "", "Actor attribution")
	cmd.Flags().StringVar(&key, "idempotency-key", "", "Retry key")
	return cmd
}
