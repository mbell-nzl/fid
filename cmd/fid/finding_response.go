package main

import (
	"fmt"

	"github.com/mbell-nzl/fid/internal/cli"
	"github.com/mbell-nzl/fid/internal/reviewwork"
	"github.com/spf13/cobra"
)

func init() {
	for _, leaf := range []string{"finding answer", "finding defer"} {
		rootFlagCapabilities[leaf] = rootFlagCapability{Project: true, JSON: true}
	}
}

func newFindingResponseCommand(deferFinding bool) *cobra.Command {
	o := reviewwork.ResponseOptions{Defer: deferFinding}
	verb, flag := "answer", "text"
	if deferFinding {
		verb, flag = "defer", "title"
	}
	operation := "finding." + verb
	cmd := &cobra.Command{Use: verb + " <finding>", Short: "Record a builder's or person's finding response.", Args: canonicalExactArgs(operation, 1), RunE: func(cmd *cobra.Command, args []string) error {
		for _, required := range []string{"actor", flag} {
			if err := requireCanonicalFlag(cmd, operation, required); err != nil {
				return err
			}
		}
		root, err := rebuildWorkspaceRoot(cmd)
		if err != nil {
			return reviewFailure(cmd, operation, err, "")
		}
		o.WorkspaceRoot, o.MigrationsDir, o.FindingID, o.Clock = root, issueMigrationsDir(), args[0], commandDisplayEnvironment(cmd).now
		result, err := reviewwork.Respond(o)
		if err != nil {
			return reviewFailure(cmd, operation, err, "")
		}
		postFindingResponse(cmd, root, args[0])
		jsonMode, err := readJSONMode(cmd)
		if err != nil {
			return err
		}
		if jsonMode {
			return cli.WriteSuccess(cmd.OutOrStdout(), operation, map[string]any{"response": result.Response}, &cli.Event{ID: result.EventID, OccurredAt: result.Response.At})
		}
		message := "Answered " + args[0] + " · " + result.Response.Text
		if d := result.Response.DeferredTask; d != nil {
			message = "Deferred " + args[0] + " → " + d.ID + " · " + d.Title
		}
		_, err = fmt.Fprintln(cmd.OutOrStdout(), "✓ "+cli.EscapeControlText(message))
		return err
	}}
	cmd.Flags().StringVar(&o.Actor, "actor", "", "Builder or person recording the response")
	if deferFinding {
		cmd.Flags().StringVar(&o.Title, "title", "", "Title of the discovered task")
		cmd.Flags().StringVar(&o.EpicID, "epic", "", "Epic for the discovered task (defaults to the reviewed task's epic)")
		cmd.Flags().StringVar(&o.Lane, "lane", "", "Lane for the discovered task (at least Standard)")
	} else {
		cmd.Flags().StringVar(&o.Text, "text", "", "Answer to the finding")
	}
	return cmd
}
