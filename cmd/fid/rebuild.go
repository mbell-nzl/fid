package main

import (
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/mbell-nzl/fid/internal/cli"
	"github.com/mbell-nzl/fid/internal/journal"
	"github.com/mbell-nzl/fid/internal/policy"
	"github.com/mbell-nzl/fid/internal/project"
	"github.com/spf13/cobra"
)

func init() {
	rootFlagCapabilities["rebuild"] = rootFlagCapability{Project: true, JSON: true}
}

func newRebuildCommand() *cobra.Command {
	var actorRole string
	var verbose bool

	cmd := &cobra.Command{
		Use:   "rebuild",
		Short: "Rebuild the SQLite projection from the journal.",
		Args:  validateRebuildArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			jsonMode, err := cmd.Root().PersistentFlags().GetBool("json")
			if err != nil {
				return fmt.Errorf("read json flag: %w", err)
			}
			if jsonMode && verbose {
				return renderRebuildFailure(cmd, jsonMode, cli.CodeInvalidArgument, "verbose is unavailable with JSON output", nil, map[string]any{
					"argument": "verbose",
					"value":    "true",
					"allowed":  []string{"human output only"},
				})
			}

			workspaceRoot, err := rebuildWorkspaceRoot(cmd)
			if err != nil {
				return renderRebuildProjectError(cmd, jsonMode, err)
			}

			if err := policy.Authorize(workspaceRoot, actorRole, "maintenance.rebuild"); err != nil {
				return renderRebuildAuthorityError(cmd, jsonMode, actorRole, err)
			}

			startedAt := time.Now()
			stats, err := journal.RebuildWorkspace(workspaceRoot, rebuildMigrationsDir(), actorRole)
			if err != nil {
				return renderRebuildOperationError(cmd, jsonMode, err)
			}

			if jsonMode {
				return cli.WriteSuccess(cmd.OutOrStdout(), "project.rebuild", map[string]any{
					"events_replayed": stats.EventsReplayed,
					"files_replayed":  stats.FilesReplayed,
					"upgraded_events": stats.UpgradedEvents,
					"schema_version":  1,
				}, nil)
			}

			if err := writeRebuildLine(cmd.OutOrStdout(), "rebuilding .fid/fid.sqlite from journal...\n"); err != nil {
				return err
			}
			if err := writeRebuildLine(cmd.OutOrStdout(), "  replayed %d events from %d journal files\n", stats.EventsReplayed, stats.FilesReplayed); err != nil {
				return err
			}
			if err := writeRebuildLine(cmd.OutOrStdout(), "  schema version: %d\n", 1); err != nil {
				return err
			}
			if verbose || stats.UpgradedEvents > 0 {
				if err := writeRebuildLine(cmd.OutOrStdout(), "  upgraded %d events\n", stats.UpgradedEvents); err != nil {
					return err
				}
			}
			if err := writeRebuildLine(cmd.OutOrStdout(), "  elapsed: %s\n", time.Since(startedAt).Round(time.Millisecond)); err != nil {
				return err
			}
			if err := writeRebuildLine(cmd.OutOrStdout(), "rebuild complete\n"); err != nil {
				return err
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&actorRole, "by", "", "Acting role")
	cmd.Flags().BoolVar(&verbose, "verbose", false, "Print per-file replay details")

	return cmd
}

func validateRebuildArgs(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return nil
	}
	jsonMode, err := cmd.Root().PersistentFlags().GetBool("json")
	if err != nil {
		return fmt.Errorf("read json flag: %w", err)
	}
	return renderRebuildFailure(cmd, jsonMode, cli.CodeInvalidArgument, "rebuild accepts no positional arguments", nil, map[string]any{
		"argument": "arguments",
		"value":    strings.Join(args, " "),
		"allowed":  []string{"no positional arguments"},
	})
}

func rebuildWorkspaceRoot(cmd *cobra.Command) (string, error) {
	workingDirectory, err := commandWorkingDirectory(cmd)
	if err != nil {
		return "", err
	}
	projectPath, err := rebuildProjectCandidate(cmd, workingDirectory)
	if err != nil {
		return "", err
	}
	return project.Resolve(workingDirectory, projectPath)
}

func rebuildProjectCandidate(cmd *cobra.Command, workingDirectory string) (string, error) {
	projectPath, err := cmd.Root().PersistentFlags().GetString("project")
	if err != nil {
		return "", fmt.Errorf("read project flag: %w", err)
	}
	if strings.TrimSpace(projectPath) == "" {
		return workingDirectory, nil
	}
	if !filepath.IsAbs(projectPath) {
		projectPath = filepath.Join(workingDirectory, projectPath)
	}
	root, err := filepath.Abs(projectPath)
	if err != nil {
		return "", fmt.Errorf("resolve project path: %w", err)
	}
	return filepath.Clean(root), nil
}

func renderRebuildProjectError(cmd *cobra.Command, jsonMode bool, rebuildErr error) error {
	if !errors.Is(rebuildErr, project.ErrNotFound) {
		return renderRebuildOperationError(cmd, jsonMode, rebuildErr)
	}
	workingDirectory, err := commandWorkingDirectory(cmd)
	if err != nil {
		return err
	}
	projectPath, err := rebuildProjectCandidate(cmd, workingDirectory)
	if err != nil {
		return err
	}
	return renderRebuildFailure(cmd, jsonMode, cli.CodeProjectNotFound, "fid project was not found", &cli.Subject{
		Type: cli.SubjectProject,
		ID:   projectPath,
	}, map[string]any{"path": projectPath})
}

func renderRebuildAuthorityError(cmd *cobra.Command, jsonMode bool, actorRole string, rebuildErr error) error {
	if !jsonMode {
		if err := renderIssueError(cmd.ErrOrStderr(), rebuildErr); err != nil {
			return err
		}
		return rebuildErr
	}
	failure, handled, err := authorityFailure("project.rebuild", actorRole, rebuildErr)
	if err != nil {
		return err
	}
	if handled {
		return renderRebuildFailure(cmd, jsonMode, failure.Code, failure.Message, failure.Subject, failure.Details)
	}
	return renderRebuildOperationError(cmd, jsonMode, rebuildErr)
}

func renderRebuildOperationError(cmd *cobra.Command, jsonMode bool, rebuildErr error) error {
	if !jsonMode {
		if renderErr := renderRebuildError(cmd, rebuildErr); renderErr != nil {
			return renderErr
		}
		return rebuildErr
	}
	if failure, handled, err := lockUnavailableFailure(rebuildErr); err != nil {
		return err
	} else if handled {
		return renderRebuildFailure(cmd, jsonMode, failure.Code, failure.Message, failure.Subject, failure.Details)
	}
	if failure, handled, err := journalReplayStoppedFailure(rebuildErr); err != nil {
		return err
	} else if handled {
		return renderRebuildFailure(cmd, jsonMode, failure.Code, failure.Message, failure.Subject, failure.Details)
	}
	failure, err := cli.NewFailure(cli.CodeInternalError, "projection rebuild failed", nil, map[string]any{})
	if err != nil {
		return err
	}
	if err := cli.WriteFailure(cmd.ErrOrStderr(), "project.rebuild", failure); err != nil {
		return err
	}
	return rebuildErr
}

func renderRebuildFailure(cmd *cobra.Command, jsonMode bool, code cli.Code, message string, subject *cli.Subject, details map[string]any) error {
	failure, err := cli.NewFailure(code, message, subject, details)
	if err != nil {
		return err
	}
	if jsonMode {
		if err := cli.WriteFailure(cmd.ErrOrStderr(), "project.rebuild", failure); err != nil {
			return err
		}
		return fmt.Errorf("%s", code)
	}
	if err := writeRebuildLine(cmd.ErrOrStderr(), "error: %s\n", message); err != nil {
		return err
	}
	return fmt.Errorf("%s", code)
}

func rebuildMigrationsDir() string {
	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		return "migrations"
	}

	return filepath.Clean(filepath.Join(filepath.Dir(currentFile), "..", "..", "migrations"))
}

func renderRebuildError(cmd *cobra.Command, err error) error {
	handled, renderErr := renderJournalReplayError(func(format string, args ...any) error {
		return writeRebuildLine(cmd.ErrOrStderr(), format, args...)
	}, err, "rebuild stopped - manual repair required")
	if handled {
		return renderErr
	}

	return writeRebuildLine(cmd.ErrOrStderr(), "error: %v\n", err)
}

func writeRebuildLine(writer io.Writer, format string, args ...any) error {
	if _, err := fmt.Fprintf(writer, format, args...); err != nil {
		return fmt.Errorf("write rebuild output: %w", err)
	}
	return nil
}
