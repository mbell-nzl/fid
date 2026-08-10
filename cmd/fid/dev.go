package main

import (
	"fmt"

	internaldev "github.com/mbell-nzl/fid/internal/dev"
	"github.com/spf13/cobra"
)

func newDevCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "dev",
		Short: "Run developer composite workflows.",
	}

	cmd.AddCommand(newDevStartCommand())
	cmd.AddCommand(newDevGateCommand())
	cmd.AddCommand(newDevFinishCommand())
	cmd.AddCommand(newDevBlockCommand())
	return cmd
}

func newDevStartCommand() *cobra.Command {
	var (
		actorRole string
		model     string
	)

	cmd := &cobra.Command{
		Use:   "start <issue-id>",
		Short: "Claim a task, transition it to in_progress, and open a run.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			requiredFlags := []string{"model", "by"}
			for _, name := range requiredFlags {
				if cmd.Flags().Changed(name) {
					continue
				}
				err := fmt.Errorf("--%s is required", name)
				if renderErr := renderIssueError(cmd.ErrOrStderr(), err); renderErr != nil {
					return renderErr
				}
				return err
			}

			workspaceRoot, err := commandWorkingDirectory(cmd)
			if err != nil {
				return fmt.Errorf("resolve workspace root: %w", err)
			}

			result, err := internaldev.Start(internaldev.StartOptions{
				WorkspaceRoot: workspaceRoot,
				MigrationsDir: issueMigrationsDir(),
				IssueID:       args[0],
				ActorRole:     actorRole,
				Model:         model,
			})
			if err != nil {
				if renderErr := renderIssueError(cmd.ErrOrStderr(), err); renderErr != nil {
					return renderErr
				}
				return err
			}

			if err := writeIssueLine(cmd.ErrOrStderr(), "dev start %s\n", result.IssueID); err != nil {
				return err
			}
			if err := writeIssueLine(cmd.ErrOrStderr(), "  claimed by %s (expires %s)\n", result.ClaimedBy, result.ExpiresAt); err != nil {
				return err
			}
			if err := writeIssueLine(cmd.ErrOrStderr(), "  transitioned: %s → %s\n", result.FromState, result.ToState); err != nil {
				return err
			}
			return writeIssueLine(cmd.ErrOrStderr(), "  agent run: %s opened\n", result.RunID)
		},
	}

	cmd.Flags().StringVar(&model, "model", "", "Model used for the opened run")
	cmd.Flags().StringVar(&actorRole, "by", "", "Acting role")
	return cmd
}

func newDevGateCommand() *cobra.Command {
	var (
		actorRole     string
		command       string
		exitCode      int
		outputExcerpt string
	)

	cmd := &cobra.Command{
		Use:   "gate <issue-id>",
		Short: "Record one reported gate result on the active attempt.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			requiredFlags := []string{"cmd", "exit-code", "by"}
			for _, name := range requiredFlags {
				if cmd.Flags().Changed(name) {
					continue
				}
				err := fmt.Errorf("--%s is required", name)
				if renderErr := renderIssueError(cmd.ErrOrStderr(), err); renderErr != nil {
					return renderErr
				}
				return err
			}

			workspaceRoot, err := commandWorkingDirectory(cmd)
			if err != nil {
				return fmt.Errorf("resolve workspace root: %w", err)
			}

			result, err := internaldev.Gate(internaldev.GateOptions{
				WorkspaceRoot: workspaceRoot,
				MigrationsDir: issueMigrationsDir(),
				IssueID:       args[0],
				ActorRole:     actorRole,
				Command:       command,
				ExitCode:      exitCode,
				OutputExcerpt: changedOptionalString(cmd, "output-excerpt", outputExcerpt),
			})
			if err != nil {
				if renderErr := renderIssueError(cmd.ErrOrStderr(), err); renderErr != nil {
					return renderErr
				}
				return err
			}

			if err := writeIssueLine(cmd.ErrOrStderr(), "gate recorded on %s (attempt %d)\n", result.IssueID, result.AttemptNumber); err != nil {
				return err
			}
			if err := writeIssueLine(cmd.ErrOrStderr(), "  cmd: %s\n", result.Command); err != nil {
				return err
			}
			if err := writeIssueLine(cmd.ErrOrStderr(), "  exit_code: %d\n", result.ExitCode); err != nil {
				return err
			}
			return writeIssueLine(cmd.ErrOrStderr(), "  gates_passed: %d, gates_failed: %d\n", result.GatesPassed, result.GatesFailed)
		},
	}

	cmd.Flags().StringVar(&command, "cmd", "", "Reported gate command")
	cmd.Flags().IntVar(&exitCode, "exit-code", 0, "Reported gate exit code")
	cmd.Flags().StringVar(&outputExcerpt, "output-excerpt", "", "Optional reported gate output excerpt")
	cmd.Flags().StringVar(&actorRole, "by", "", "Acting role")
	return cmd
}

func newDevFinishCommand() *cobra.Command {
	var (
		actorRole         string
		completionMDPath  string
		actualTokens      int
		commits           string
		filesTouched      []string
		humanIntervention bool
		interventionNote  string
	)

	cmd := &cobra.Command{
		Use:   "finish <issue-id>",
		Short: "Close a run, complete the task, transition it, and release the claim.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			requiredFlags := []string{"completion-md", "by"}
			for _, name := range requiredFlags {
				if cmd.Flags().Changed(name) {
					continue
				}
				err := fmt.Errorf("--%s is required", name)
				if renderErr := renderIssueError(cmd.ErrOrStderr(), err); renderErr != nil {
					return renderErr
				}
				return err
			}

			workspaceRoot, err := commandWorkingDirectory(cmd)
			if err != nil {
				return fmt.Errorf("resolve workspace root: %w", err)
			}

			result, err := internaldev.Finish(internaldev.FinishOptions{
				WorkspaceRoot:         workspaceRoot,
				MigrationsDir:         issueMigrationsDir(),
				IssueID:               args[0],
				ActorRole:             actorRole,
				CompletionMDPath:      completionMDPath,
				ActualTokens:          changedOptionalInt(cmd, "tokens", actualTokens),
				CommitHashes:          issueSplitCommaSeparated(commits),
				FilesTouched:          filesTouched,
				HumanIntervention:     humanIntervention,
				HumanInterventionNote: changedOptionalString(cmd, "intervention-note", interventionNote),
			})
			if err != nil {
				if renderErr := renderIssueError(cmd.ErrOrStderr(), err); renderErr != nil {
					return renderErr
				}
				return err
			}

			if err := writeIssueLine(cmd.ErrOrStderr(), "dev finish %s\n", result.IssueID); err != nil {
				return err
			}
			if result.RunID != nil {
				runLine := "  agent run: %s closed\n"
				if result.Completion.ActualTokens != nil {
					runLine = "  agent run: %s closed (%d tokens)\n"
					if err := writeIssueLine(cmd.ErrOrStderr(), runLine, *result.RunID, *result.Completion.ActualTokens); err != nil {
						return err
					}
				} else if err := writeIssueLine(cmd.ErrOrStderr(), runLine, *result.RunID); err != nil {
					return err
				}
			}
			if err := writeIssueLine(cmd.ErrOrStderr(), "  completion: %s (attempt %d)\n", result.Completion.ID, result.Completion.AttemptNumber); err != nil {
				return err
			}
			if err := writeIssueLine(cmd.ErrOrStderr(), "  transitioned: %s → %s\n", result.FromState, result.ToState); err != nil {
				return err
			}
			return writeIssueLine(cmd.ErrOrStderr(), "  claim released\n")
		},
	}

	cmd.Flags().StringVar(&completionMDPath, "completion-md", "", "Path to completion report")
	cmd.Flags().IntVar(&actualTokens, "tokens", 0, "Actual tokens used")
	cmd.Flags().StringVar(&commits, "commits", "", "Comma-separated commit hashes")
	cmd.Flags().StringSliceVar(&filesTouched, "files", nil, "Files touched")
	cmd.Flags().BoolVar(&humanIntervention, "human-intervention", false, "Mark that a human intervened during the attempt")
	cmd.Flags().StringVar(&interventionNote, "intervention-note", "", "Short note stored when --human-intervention is set")
	cmd.Flags().StringVar(&actorRole, "by", "", "Acting role")
	return cmd
}

func newDevBlockCommand() *cobra.Command {
	var (
		actorRole    string
		reasonMDPath string
	)

	cmd := &cobra.Command{
		Use:   "block <issue-id>",
		Short: "Record a failed-gate completion, transition the task, and release the claim.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			requiredFlags := []string{"reason-md", "by"}
			for _, name := range requiredFlags {
				if cmd.Flags().Changed(name) {
					continue
				}
				err := fmt.Errorf("--%s is required", name)
				if renderErr := renderIssueError(cmd.ErrOrStderr(), err); renderErr != nil {
					return renderErr
				}
				return err
			}

			workspaceRoot, err := commandWorkingDirectory(cmd)
			if err != nil {
				return fmt.Errorf("resolve workspace root: %w", err)
			}

			result, err := internaldev.Block(internaldev.BlockOptions{
				WorkspaceRoot: workspaceRoot,
				MigrationsDir: issueMigrationsDir(),
				IssueID:       args[0],
				ActorRole:     actorRole,
				ReasonMDPath:  reasonMDPath,
			})
			if err != nil {
				if renderErr := renderIssueError(cmd.ErrOrStderr(), err); renderErr != nil {
					return renderErr
				}
				return err
			}

			if err := writeIssueLine(cmd.ErrOrStderr(), "dev block %s\n", result.IssueID); err != nil {
				return err
			}
			if err := writeIssueLine(cmd.ErrOrStderr(), "  transitioned: %s → %s\n", result.FromState, result.ToState); err != nil {
				return err
			}
			if err := writeIssueLine(cmd.ErrOrStderr(), "  completion: %s (attempt %d, status: %s)\n", result.Completion.ID, result.Completion.AttemptNumber, result.Completion.Status); err != nil {
				return err
			}
			return writeIssueLine(cmd.ErrOrStderr(), "  claim released\n")
		},
	}

	cmd.Flags().StringVar(&reasonMDPath, "reason-md", "", "Path to the failed-gate reason report")
	cmd.Flags().StringVar(&actorRole, "by", "", "Acting role")
	return cmd
}
