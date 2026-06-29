package main

import (
	"fmt"
	"os"

	"github.com/mbell-nzl/fid/internal/claims"
	"github.com/spf13/cobra"
)

func newIssueClaimCommand() *cobra.Command {
	var (
		actorRole string
		ttl       string
		force     bool
		renew     bool
		note      string
	)

	cmd := &cobra.Command{
		Use:   "claim <issue-id>",
		Short: "Acquire or renew an issue claim.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !cmd.Flags().Changed("by") {
				err := fmt.Errorf("--by is required")
				if renderErr := renderIssueError(cmd.ErrOrStderr(), err); renderErr != nil {
					return renderErr
				}
				return err
			}

			workspaceRoot, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("resolve workspace root: %w", err)
			}

			ttlSeconds, err := issueClaimTTLSeconds(cmd, ttl)
			if err != nil {
				if renderErr := renderIssueError(cmd.ErrOrStderr(), err); renderErr != nil {
					return renderErr
				}
				return err
			}
			if force && renew {
				err := fmt.Errorf("--force and --renew are mutually exclusive")
				if renderErr := renderIssueError(cmd.ErrOrStderr(), err); renderErr != nil {
					return renderErr
				}
				return err
			}

			if renew {
				result, err := claims.Renew(claims.RenewOptions{
					WorkspaceRoot: workspaceRoot,
					MigrationsDir: issueMigrationsDir(),
					IssueID:       args[0],
					ActorRole:     actorRole,
					TTLSeconds:    ttlSeconds,
				})
				if err != nil {
					if renderErr := renderIssueError(cmd.ErrOrStderr(), err); renderErr != nil {
						return renderErr
					}
					return err
				}

				return writeIssueLine(cmd.ErrOrStderr(), "renewed claim on %s by %s (expires %s)\n", result.RefID, result.ClaimedBy, result.ExpiresAt)
			}

			result, err := claims.Acquire(claims.AcquireOptions{
				WorkspaceRoot: workspaceRoot,
				MigrationsDir: issueMigrationsDir(),
				IssueID:       args[0],
				ActorRole:     actorRole,
				TTLSeconds:    ttlSeconds,
				Note:          issueCreateOptionalString(cmd, "note", note),
				Force:         force,
			})
			if err != nil {
				if renderErr := renderIssueError(cmd.ErrOrStderr(), err); renderErr != nil {
					return renderErr
				}
				return err
			}

			if force && result.Previous != nil {
				return writeIssueLine(
					cmd.ErrOrStderr(),
					"force-acquired claim on %s by %s (previously held by %s until %s)\n",
					result.RefID,
					result.ClaimedBy,
					result.Previous.ClaimedBy,
					result.Previous.ExpiresAt,
				)
			}
			return writeIssueLine(cmd.ErrOrStderr(), "claimed %s by %s (expires %s)\n", result.RefID, result.ClaimedBy, result.ExpiresAt)
		},
	}

	cmd.Flags().StringVar(&actorRole, "by", "", "Claiming role")
	cmd.Flags().StringVar(&ttl, "ttl", "30m", "Lease duration")
	cmd.Flags().BoolVar(&force, "force", false, "Force-acquire a live claim")
	cmd.Flags().BoolVar(&renew, "renew", false, "Renew the current live claim")
	cmd.Flags().StringVar(&note, "note", "", "Optional claim context")
	return cmd
}

func newIssueUnclaimCommand() *cobra.Command {
	var actorRole string

	cmd := &cobra.Command{
		Use:   "unclaim <issue-id>",
		Short: "Release the caller's current issue claim.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !cmd.Flags().Changed("by") {
				err := fmt.Errorf("--by is required")
				if renderErr := renderIssueError(cmd.ErrOrStderr(), err); renderErr != nil {
					return renderErr
				}
				return err
			}

			workspaceRoot, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("resolve workspace root: %w", err)
			}

			result, err := claims.Release(claims.ReleaseOptions{
				WorkspaceRoot: workspaceRoot,
				MigrationsDir: issueMigrationsDir(),
				IssueID:       args[0],
				ActorRole:     actorRole,
			})
			if err != nil {
				if renderErr := renderIssueError(cmd.ErrOrStderr(), err); renderErr != nil {
					return renderErr
				}
				return err
			}

			return writeIssueLine(cmd.ErrOrStderr(), "released claim on %s by %s\n", result.RefID, result.ReleasedBy)
		},
	}

	cmd.Flags().StringVar(&actorRole, "by", "", "Acting role releasing its claim")
	return cmd
}
