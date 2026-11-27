package main

import (
	"errors"
	"fmt"

	"github.com/mbell-nzl/fid/internal/cli"
	"github.com/mbell-nzl/fid/internal/usagework"
	fidmigrations "github.com/mbell-nzl/fid/migrations"
	"github.com/spf13/cobra"
)

func init() { rootFlagCapabilities["usage record"] = rootFlagCapability{Project: true, JSON: true} }
func newUsageCommand() *cobra.Command {
	var o usagework.Options
	cmd := &cobra.Command{Use: "usage", Short: "Record agent session token usage."}
	record := &cobra.Command{Use: "record", Short: "Read tokens from a session's agent output.", RunE: func(cmd *cobra.Command, _ []string) error {
		for _, flag := range []string{"hold", "actor", "format", "input", "idempotency-key"} {
			if err := requireCanonicalFlag(cmd, "usage.record", flag); err != nil {
				return err
			}
		}
		root, err := rebuildWorkspaceRoot(cmd)
		if err != nil {
			return renderHoldFailure(cmd, "usage.record", err)
		}
		o.WorkspaceRoot, o.MigrationsFS, o.Clock = root, fidmigrations.FS, commandDisplayEnvironment(cmd).now
		result, event, err := usagework.Record(o)
		if err != nil {
			var input *usagework.InputError
			if errors.As(err, &input) {
				return writeGrammarFailure(cmd, "usage.record", cli.CodeInvalidArgument, input.Message, map[string]any{"argument": input.Argument, "value": nil, "allowed": []string{}})
			}
			return renderHoldFailure(cmd, "usage.record", err)
		}
		jsonMode, err := readJSONMode(cmd)
		if err != nil {
			return err
		}
		if jsonMode {
			return cli.WriteSuccess(cmd.OutOrStdout(), "usage.record", map[string]any{"usage": result.Usage}, &cli.Event{ID: event, OccurredAt: result.Usage.RecordedAt})
		}
		_, err = fmt.Fprintf(cmd.OutOrStdout(), "✓ Recorded %d tokens for %s (%s)\n", result.Usage.Counts.Total, cli.EscapeControlText(result.Usage.TaskID), cli.EscapeControlText(result.Usage.Role))
		return err
	}}
	record.Flags().StringVar(&o.HoldID, "hold", "", "The build or review session ID")
	record.Flags().StringVar(&o.Actor, "actor", "", "The session's actor")
	record.Flags().StringVar(&o.Format, "format", "", "Agent output format: codex or claude")
	record.Flags().StringVar(&o.Input, "input", "", "Agent stdout file")
	record.Flags().StringVar(&o.ToolVersion, "tool-version", "", "Caller-reported tool version when output has none")
	record.Flags().StringVar(&o.IdempotencyKey, "idempotency-key", "", "Retry key")
	cmd.AddCommand(record)
	return cmd
}
