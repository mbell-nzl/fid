package main

import (
	"fmt"

	"github.com/mbell-nzl/fid/internal/issues"
	"github.com/spf13/cobra"
)

func newIssueCompleteCommand() *cobra.Command {
	var (
		actorRole         string
		notesMDPath       string
		actualTokens      int
		commits           string
		filesTouched      []string
		humanIntervention bool
		interventionNote  string
	)

	cmd := &cobra.Command{
		Use:   "complete <issue-id>",
		Short: "Record a task completion attempt.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !cmd.Flags().Changed("by") {
				err := fmt.Errorf("--by is required")
				if renderErr := renderIssueError(cmd.ErrOrStderr(), err); renderErr != nil {
					return renderErr
				}
				return err
			}

			workspaceRoot, err := commandWorkingDirectory(cmd)
			if err != nil {
				return fmt.Errorf("resolve workspace root: %w", err)
			}

			result, err := issues.Complete(issues.CompleteOptions{
				WorkspaceRoot:         workspaceRoot,
				MigrationsDir:         issueMigrationsDir(),
				IssueID:               args[0],
				ActorRole:             actorRole,
				NotesMDPath:           issueCreateOptionalString(cmd, "notes-md", notesMDPath),
				ActualTokens:          issueCreateOptionalInt(cmd, "tokens", actualTokens),
				CommitHashes:          issueSplitCommaSeparated(commits),
				FilesTouched:          filesTouched,
				HumanIntervention:     humanIntervention,
				HumanInterventionNote: issueCreateOptionalString(cmd, "intervention-note", interventionNote),
			})
			if err != nil {
				if renderErr := renderIssueError(cmd.ErrOrStderr(), err); renderErr != nil {
					return renderErr
				}
				return err
			}

			if err := writeIssueLine(cmd.ErrOrStderr(), "completed %s (attempt %d)\n", result.IssueID, result.AttemptNumber); err != nil {
				return err
			}
			if err := writeIssueLine(cmd.ErrOrStderr(), "  completion: %s\n", result.ID); err != nil {
				return err
			}
			if result.ActualTokens != nil {
				if err := writeIssueLine(cmd.ErrOrStderr(), "  tokens: %d\n", *result.ActualTokens); err != nil {
					return err
				}
			}
			if len(result.CommitHashes) > 0 {
				if err := writeIssueLine(cmd.ErrOrStderr(), "  commits: %s\n", joinCommaSeparated(result.CommitHashes)); err != nil {
					return err
				}
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&actorRole, "by", "", "Acting role")
	cmd.Flags().StringVar(&notesMDPath, "notes-md", "", "Path to completion notes")
	cmd.Flags().IntVar(&actualTokens, "tokens", 0, "Actual tokens used")
	cmd.Flags().StringVar(&commits, "commits", "", "Comma-separated commit hashes")
	cmd.Flags().StringSliceVar(&filesTouched, "files", nil, "Files touched")
	cmd.Flags().BoolVar(&humanIntervention, "human-intervention", false, "Mark that a human intervened during the attempt")
	cmd.Flags().StringVar(&interventionNote, "intervention-note", "", "Short note stored when --human-intervention is set")
	return cmd
}
