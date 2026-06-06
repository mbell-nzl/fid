package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/mbell-nzl/fid/internal/journal"
	"github.com/mbell-nzl/fid/internal/policy"
	"github.com/spf13/cobra"
)

func newRebuildCommand() *cobra.Command {
	var actorRole string
	var verbose bool

	cmd := &cobra.Command{
		Use:   "rebuild",
		Short: "Rebuild the SQLite projection from the journal.",
		RunE: func(cmd *cobra.Command, args []string) error {
			workspaceRoot, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("resolve workspace root: %w", err)
			}

			if err := policy.RequireKnownRole(workspaceRoot, actorRole); err != nil {
				if renderErr := renderIssueError(cmd.ErrOrStderr(), err); renderErr != nil {
					return renderErr
				}
				return err
			}

			startedAt := time.Now()
			if err := writeRebuildLine(cmd.ErrOrStderr(), "rebuilding .fid/fid.sqlite from journal...\n"); err != nil {
				return err
			}

			stats, err := journal.RebuildWorkspace(workspaceRoot, rebuildMigrationsDir(), actorRole)
			if err != nil {
				if renderErr := renderRebuildError(cmd, err); renderErr != nil {
					return renderErr
				}
				return err
			}

			if err := writeRebuildLine(cmd.ErrOrStderr(), "  replayed %d events from %d journal files\n", stats.EventsReplayed, stats.FilesReplayed); err != nil {
				return err
			}
			if err := writeRebuildLine(cmd.ErrOrStderr(), "  schema version: %d\n", 1); err != nil {
				return err
			}
			if verbose || stats.UpgradedEvents > 0 {
				if err := writeRebuildLine(cmd.ErrOrStderr(), "  upgraded %d events\n", stats.UpgradedEvents); err != nil {
					return err
				}
			}
			if err := writeRebuildLine(cmd.ErrOrStderr(), "  elapsed: %s\n", time.Since(startedAt).Round(time.Millisecond)); err != nil {
				return err
			}
			if err := writeRebuildLine(cmd.ErrOrStderr(), "rebuild complete\n"); err != nil {
				return err
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&actorRole, "by", "", "Acting role")
	cmd.Flags().BoolVar(&verbose, "verbose", false, "Print per-file replay details")

	return cmd
}

func rebuildMigrationsDir() string {
	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		return "migrations"
	}

	return filepath.Clean(filepath.Join(filepath.Dir(currentFile), "..", "..", "migrations"))
}

func renderRebuildError(cmd *cobra.Command, err error) error {
	var corruptErr *journal.CorruptJournalError
	if errors.As(err, &corruptErr) {
		if writeErr := writeRebuildLine(cmd.ErrOrStderr(), "error: malformed journal entry at %s:%d\n", corruptErr.File, corruptErr.Line); writeErr != nil {
			return writeErr
		}
		if writeErr := writeRebuildLine(cmd.ErrOrStderr(), "  parse error: %v\n", corruptErr.Err); writeErr != nil {
			return writeErr
		}
		return writeRebuildLine(cmd.ErrOrStderr(), "  rebuild stopped - manual repair required\n")
	}

	return writeRebuildLine(cmd.ErrOrStderr(), "error: %v\n", err)
}

func writeRebuildLine(writer io.Writer, format string, args ...any) error {
	if _, err := fmt.Fprintf(writer, format, args...); err != nil {
		return fmt.Errorf("write rebuild output: %w", err)
	}
	return nil
}
