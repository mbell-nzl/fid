package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/mbell-nzl/fid/internal/epicwork"
	fidmigrations "github.com/mbell-nzl/fid/migrations"
	"github.com/spf13/cobra"
)

func init() {
	for _, leaf := range []string{"epic validation verify", "epic validation pass"} {
		rootFlagCapabilities[leaf] = rootFlagCapability{Project: true, JSON: true}
	}
}

func newEpicEvidenceCommand(action string) *cobra.Command {
	var o epicwork.Options
	op := "epic.validation." + action
	use, n := action+" <epic>", 1
	if action == "pass" {
		use += " <id>"
		n = 2
	}
	cmd := &cobra.Command{Use: use, Short: "Record validation evidence on the default branch.", Args: canonicalExactArgs(op, n), RunE: func(cmd *cobra.Command, args []string) error {
		required := []string{"actor"}
		if action == "verify" {
			required = append(required, "idempotency-key")
		} else {
			required = append(required, "reason")
		}
		for _, flag := range required {
			if err := requireCanonicalFlag(cmd, op, flag); err != nil {
				return err
			}
		}
		root, err := rebuildWorkspaceRoot(cmd)
		if err != nil {
			return renderHoldFailure(cmd, op, err)
		}
		o.WorkspaceRoot, o.MigrationsFS, o.EpicID, o.Clock = root, fidmigrations.FS, args[0], commandDisplayEnvironment(cmd).now
		var r epicwork.Result
		if action == "pass" {
			o.ValidationID = args[1]
			r, err = epicwork.Pass(o)
		} else {
			ctx, cancel := signal.NotifyContext(cmd.Context(), os.Interrupt, syscall.SIGTERM)
			defer cancel()
			r, err = epicwork.Verify(ctx, o)
		}
		if err != nil {
			return epicFailure(cmd, op, o.Actor, err)
		}
		return writeEpicResult(cmd, op, r)
	}}
	cmd.Flags().StringVar(&o.Actor, "actor", "", "Who records the evidence")
	cmd.Flags().StringVar(&o.IdempotencyKey, "idempotency-key", "", "Retry key")
	if action == "pass" {
		cmd.Flags().StringVar(&o.Reason, "reason", "", "Why this validation passed")
	}
	return cmd
}
