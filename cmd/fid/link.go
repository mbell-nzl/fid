package main

import (
	"fmt"

	"github.com/mbell-nzl/fid/internal/cli"
	"github.com/mbell-nzl/fid/internal/links"
	"github.com/spf13/cobra"
)

func newLinkCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "link",
		Short: "Manage typed links.",
	}

	cmd.AddCommand(newLinkAddCommand())
	cmd.AddCommand(newLinkListCommand())
	cmd.AddCommand(newLinkRemoveCommand())
	return cmd
}

func newLinkAddCommand() *cobra.Command {
	var actor, key, fromType, toType string

	cmd := &cobra.Command{
		Use:   "add <from-id> <relation> <to-id>",
		Short: "Add a typed link and its inverse.",
		Args:  canonicalExactArgs("link.add", 3),
		RunE: func(cmd *cobra.Command, args []string) error {
			for _, flag := range []string{"actor", "idempotency-key"} {
				if err := requireCanonicalFlag(cmd, "link.add", flag); err != nil {
					return err
				}
			}
			if err := requireCanonicalFlag(cmd, "link.add", "from-type"); err != nil {
				return err
			}
			if err := requireCanonicalFlag(cmd, "link.add", "to-type"); err != nil {
				return err
			}
			workspaceRoot, err := issueMutationWorkspaceRoot(cmd)
			if err != nil {
				return renderLinkError(cmd, "link.add", err)
			}

			result, err := links.Add(links.AddOptions{
				WorkspaceRoot:  workspaceRoot,
				MigrationsDir:  issueMigrationsDir(),
				FromType:       fromType,
				FromID:         args[0],
				Relation:       args[1],
				ToType:         toType,
				ToID:           args[2],
				Actor:          actor,
				IdempotencyKey: key,
			})
			if err != nil {
				return renderLinkError(cmd, "link.add", err)
			}

			jsonMode, err := readJSONMode(cmd)
			if err != nil {
				return err
			}
			if jsonMode {
				return cli.WriteSuccess(cmd.OutOrStdout(), "link.add", map[string]any{
					"link":            result.Declared,
					"issue_revisions": result.IssueRevisions,
					"idempotency_key": result.IdempotencyKey,
					"replayed":        result.Replayed,
				}, &cli.Event{ID: result.EventID, OccurredAt: result.OccurredAt})
			}
			if _, err := fmt.Fprintf(
				cmd.OutOrStdout(),
				"linked %s %s %s\n  inverse: %s %s %s\n",
				result.Declared.FromID,
				result.Declared.Relation,
				result.Declared.ToID,
				result.Inverse.FromID,
				result.Inverse.Relation,
				result.Inverse.ToID,
			); err != nil {
				return err
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&actor, "actor", "", "Standalone attribution")
	cmd.Flags().StringVar(&key, "idempotency-key", "", "Idempotency key")
	cmd.Flags().StringVar(&fromType, "from-type", "", "Type of the source endpoint")
	cmd.Flags().StringVar(&toType, "to-type", "", "Type of the destination endpoint")
	cmd.SetFlagErrorFunc(canonicalFlagError("link.add"))
	return cmd
}

func newLinkRemoveCommand() *cobra.Command {
	var actor, key, fromType, toType string

	cmd := &cobra.Command{
		Use:   "remove <from-id> <relation> <to-id>",
		Short: "Remove a typed link and its inverse.",
		Args:  canonicalExactArgs("link.remove", 3),
		RunE: func(cmd *cobra.Command, args []string) error {
			for _, flag := range []string{"actor", "idempotency-key"} {
				if err := requireCanonicalFlag(cmd, "link.remove", flag); err != nil {
					return err
				}
			}
			if err := requireCanonicalFlag(cmd, "link.remove", "from-type"); err != nil {
				return err
			}
			if err := requireCanonicalFlag(cmd, "link.remove", "to-type"); err != nil {
				return err
			}
			workspaceRoot, err := issueMutationWorkspaceRoot(cmd)
			if err != nil {
				return renderLinkError(cmd, "link.remove", err)
			}

			result, err := links.Remove(links.RemoveOptions{
				WorkspaceRoot:  workspaceRoot,
				MigrationsDir:  issueMigrationsDir(),
				FromType:       fromType,
				FromID:         args[0],
				Relation:       args[1],
				ToType:         toType,
				ToID:           args[2],
				Actor:          actor,
				IdempotencyKey: key,
			})
			if err != nil {
				return renderLinkError(cmd, "link.remove", err)
			}

			jsonMode, err := readJSONMode(cmd)
			if err != nil {
				return err
			}
			if jsonMode {
				return cli.WriteSuccess(cmd.OutOrStdout(), "link.remove", map[string]any{
					"link":            result.Declared,
					"issue_revisions": result.IssueRevisions,
					"idempotency_key": result.IdempotencyKey,
					"replayed":        result.Replayed,
				}, &cli.Event{ID: result.EventID, OccurredAt: result.OccurredAt})
			}
			if _, err := fmt.Fprintf(
				cmd.OutOrStdout(),
				"removed link %s %s %s\n  inverse removed: %s %s %s\n",
				result.Declared.FromID,
				result.Declared.Relation,
				result.Declared.ToID,
				result.Inverse.FromID,
				result.Inverse.Relation,
				result.Inverse.ToID,
			); err != nil {
				return err
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&actor, "actor", "", "Standalone attribution")
	cmd.Flags().StringVar(&key, "idempotency-key", "", "Idempotency key")
	cmd.Flags().StringVar(&fromType, "from-type", "", "Type of the source endpoint")
	cmd.Flags().StringVar(&toType, "to-type", "", "Type of the destination endpoint")
	cmd.SetFlagErrorFunc(canonicalFlagError("link.remove"))
	return cmd
}

func newLinkListCommand() *cobra.Command {
	var entityType string

	cmd := &cobra.Command{
		Use:   "list <id>",
		Short: "List outgoing and incoming links for an entity.",
		Args:  canonicalExactArgs("link.list", 1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := requireCanonicalFlag(cmd, "link.list", "type"); err != nil {
				return err
			}
			workspaceRoot, err := issueMutationWorkspaceRoot(cmd)
			if err != nil {
				return renderLinkError(cmd, "link.list", err)
			}

			items, err := links.List(links.ListOptions{
				WorkspaceRoot: workspaceRoot,
				MigrationsDir: issueMigrationsDir(),
				EntityType:    entityType,
				EntityID:      args[0],
			})
			if err != nil {
				return renderLinkError(cmd, "link.list", err)
			}

			jsonMode, err := readJSONMode(cmd)
			if err != nil {
				return err
			}
			if jsonMode {
				return cli.WriteSuccess(cmd.OutOrStdout(), "link.list", map[string]any{
					"entity_type": entityType,
					"entity_id":   args[0],
					"links":       items,
					"count":       len(items),
				}, nil)
			}
			for _, item := range items {
				if _, err := fmt.Fprintf(cmd.OutOrStdout(), "%s %s %s\n", item.Direction, item.Relation, item.EntityID); err != nil {
					return err
				}
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&entityType, "type", "", "Type of the linked entity")
	cmd.SetFlagErrorFunc(canonicalFlagError("link.list"))
	return cmd
}

func init() {
	rootFlagCapabilities["link add"] = rootFlagCapability{Project: true, JSON: true}
	rootFlagCapabilities["link remove"] = rootFlagCapability{Project: true, JSON: true}
	rootFlagCapabilities["link list"] = rootFlagCapability{Project: true, JSON: true}
}
