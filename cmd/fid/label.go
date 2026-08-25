package main

import (
	"fmt"

	"github.com/mbell-nzl/fid/internal/cli"
	"github.com/mbell-nzl/fid/internal/labels"
	"github.com/spf13/cobra"
)

func newLabelCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "label",
		Short: "Manage labels.",
	}

	cmd.AddCommand(newLabelAddCommand())
	cmd.AddCommand(newLabelListCommand())
	return cmd
}

func newLabelAddCommand() *cobra.Command {
	var (
		actor       string
		key         string
		colour      string
		description string
	)

	cmd := &cobra.Command{
		Use:   "add <name>",
		Short: "Create a label.",
		Args:  canonicalExactArgs("label.add", 1),
		RunE: func(cmd *cobra.Command, args []string) error {
			for _, flag := range []string{"actor", "idempotency-key"} {
				if err := requireCanonicalFlag(cmd, "label.add", flag); err != nil {
					return err
				}
			}
			workspaceRoot, err := issueMutationWorkspaceRoot(cmd)
			if err != nil {
				return renderLabelError(cmd, "label.add", err)
			}

			result, err := labels.Add(labels.AddOptions{
				WorkspaceRoot:  workspaceRoot,
				MigrationsDir:  issueMigrationsDir(),
				Name:           args[0],
				Colour:         optionalStringFlag(cmd, "colour", colour),
				Description:    optionalStringFlag(cmd, "description", description),
				Actor:          actor,
				IdempotencyKey: key,
			})
			if err != nil {
				return renderLabelError(cmd, "label.add", err)
			}
			jsonMode, err := readJSONMode(cmd)
			if err != nil {
				return err
			}
			if jsonMode {
				return cli.WriteSuccess(cmd.OutOrStdout(), "label.add", map[string]any{"label": result.Label, "idempotency_key": result.IdempotencyKey, "replayed": result.Replayed}, &cli.Event{ID: result.EventID, OccurredAt: result.OccurredAt})
			}
			if _, err := fmt.Fprintf(cmd.OutOrStdout(), "created label %q\n", result.Label.Name); err != nil {
				return err
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&colour, "colour", "", "Label colour")
	cmd.Flags().StringVar(&description, "description", "", "Label description")
	cmd.Flags().StringVar(&actor, "actor", "", "Standalone attribution")
	cmd.Flags().StringVar(&key, "idempotency-key", "", "Idempotency key")
	cmd.SetFlagErrorFunc(canonicalFlagError("label.add"))
	return cmd
}

func optionalStringFlag(cmd *cobra.Command, name, value string) *string {
	if !cmd.Flags().Changed(name) {
		return nil
	}
	return &value
}

func newLabelListCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List labels.",
		Args:  canonicalNoArgs("label.list"),
		RunE: func(cmd *cobra.Command, args []string) error {
			workspaceRoot, err := issueMutationWorkspaceRoot(cmd)
			if err != nil {
				return renderLabelError(cmd, "label.list", err)
			}

			items, err := labels.List(labels.ListOptions{
				WorkspaceRoot: workspaceRoot,
				MigrationsDir: issueMigrationsDir(),
			})
			if err != nil {
				return renderLabelError(cmd, "label.list", err)
			}
			jsonMode, err := readJSONMode(cmd)
			if err != nil {
				return err
			}
			if jsonMode {
				return cli.WriteSuccess(cmd.OutOrStdout(), "label.list", map[string]any{"labels": items, "count": len(items)}, nil)
			}
			for _, item := range items {
				if _, err := fmt.Fprintf(cmd.OutOrStdout(), "%s\n", item.Name); err != nil {
					return err
				}
			}
			return nil
		},
	}
	cmd.SetFlagErrorFunc(canonicalFlagError("label.list"))

	return cmd
}

func init() {
	rootFlagCapabilities["label add"] = rootFlagCapability{Project: true, JSON: true}
	rootFlagCapabilities["label list"] = rootFlagCapability{Project: true, JSON: true}
}
