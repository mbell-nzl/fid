package main

import (
	"fmt"

	"github.com/mbell-nzl/fid/internal/cli"
	"github.com/mbell-nzl/fid/internal/requirements"
	"github.com/spf13/cobra"
)

func newRequirementCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "requirement",
		Short: "Manage requirement pointers.",
	}

	cmd.AddCommand(newRequirementSyncCommand())
	return cmd
}

func newRequirementSyncCommand() *cobra.Command {
	var (
		sourceDoc     string
		typeFilter    string
		requirementID string
		actor         string
		key           string
	)

	cmd := &cobra.Command{
		Use:   "sync",
		Short: "Sync requirement pointers from a source markdown document.",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			for _, flag := range []string{"source-doc", "actor", "idempotency-key"} {
				if err := requireCanonicalFlag(cmd, "requirement.sync", flag); err != nil {
					return err
				}
			}
			workspaceRoot, err := issueMutationWorkspaceRoot(cmd)
			if err != nil {
				return fmt.Errorf("resolve workspace root: %w", err)
			}

			result, err := requirements.Sync(requirements.SyncOptions{
				WorkspaceRoot:  workspaceRoot,
				MigrationsDir:  issueMigrationsDir(),
				SourceDoc:      sourceDoc,
				Actor:          actor,
				IdempotencyKey: key,
				TypeFilter:     typeFilter,
				RequirementID:  requirementID,
			})
			if err != nil {
				if renderErr := renderIssueError(cmd.ErrOrStderr(), err); renderErr != nil {
					return renderErr
				}
				return err
			}

			jsonMode, err := readJSONMode(cmd)
			if err != nil {
				return err
			}
			if jsonMode {
				return cli.WriteSuccess(cmd.OutOrStdout(), "requirement.sync", map[string]any{"source_doc": result.SourceDoc, "created": result.Created, "updated": result.Updated, "unchanged": result.Unchanged, "requirement_ids": result.RequirementIDs, "idempotency_key": result.IdempotencyKey, "replayed": result.Replayed}, &cli.Event{ID: result.EventID, OccurredAt: result.OccurredAt})
			}
			if _, err := fmt.Fprintf(cmd.OutOrStdout(),
				"synced requirements from %s\n  created: %d\n  updated: %d\n  unchanged: %d\n",
				result.SourceDoc,
				result.Created,
				result.Updated,
				result.Unchanged,
			); err != nil {
				return err
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&sourceDoc, "source-doc", "", "Path to the source markdown document")
	cmd.Flags().StringVar(&typeFilter, "type", "", "Optional requirement type filter: PRD or TRD")
	cmd.Flags().StringVar(&requirementID, "id", "", "Optional single requirement ID to sync")
	cmd.Flags().StringVar(&actor, "actor", "", "Standalone attribution")
	cmd.Flags().StringVar(&key, "idempotency-key", "", "Idempotency key")
	cmd.SetFlagErrorFunc(canonicalFlagError("requirement.sync"))
	return cmd
}

func init() {
	rootFlagCapabilities["requirement sync"] = rootFlagCapability{Project: true, JSON: true}
}
