package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/mbell-nzl/fid/internal/config"
	fiddb "github.com/mbell-nzl/fid/internal/db"
	fidformat "github.com/mbell-nzl/fid/internal/format"
	"github.com/mbell-nzl/fid/internal/journal"
	"github.com/mbell-nzl/fid/internal/maintenance"
	"github.com/spf13/cobra"
)

var errMaintenanceCheckFailed = errors.New("maintenance check failed")

const legacyStartedWithoutClaimGraceSeconds = 300

func newMaintenanceCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "maintenance",
		Short: "Inspect and maintain the derived projection.",
	}
	cmd.AddCommand(newMaintenanceCheckCommand())
	cmd.AddCommand(newMaintenanceVacuumCommand())
	return cmd
}

func newMaintenanceVacuumCommand() *cobra.Command {
	var actorRole string

	cmd := &cobra.Command{
		Use:   "vacuum",
		Short: "Compact the derived SQLite projection without changing the journal.",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			workspaceRoot, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("resolve workspace root: %w", err)
			}

			result, err := maintenance.Vacuum(workspaceRoot, issueMigrationsDir(), actorRole)
			if err != nil {
				return renderMaintenanceVacuumError(cmd, err)
			}
			if err := writeIssueLine(cmd.ErrOrStderr(), "vacuum complete — .fid/fid.sqlite: %d B → %d B\n", result.BeforeBytes, result.AfterBytes); err != nil {
				return err
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&actorRole, "by", "", "Acting role")
	return cmd
}

func renderMaintenanceVacuumError(cmd *cobra.Command, err error) error {
	if renderErr := renderIssueError(cmd.ErrOrStderr(), err); renderErr != nil {
		return renderErr
	}
	return err
}

func newMaintenanceCheckCommand() *cobra.Command {
	var (
		formatName  string
		warnAsError bool
	)

	cmd := &cobra.Command{
		Use:   "check",
		Short: "Audit projected workspace integrity without changing the journal.",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			workspaceRoot, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("resolve workspace root: %w", err)
			}
			if _, err := journal.EnsureProjectionFresh(workspaceRoot, issueMigrationsDir()); err != nil {
				return renderMaintenanceCheckError(cmd, err)
			}

			paths := config.WorkspacePathsForRoot(workspaceRoot)
			conn, err := fiddb.Open(paths.DatabasePath)
			if err != nil {
				return renderMaintenanceCheckError(cmd, err)
			}
			defer func() { _ = conn.Close() }()

			result, err := maintenance.Check(maintenance.CheckOptions{
				DB:                              conn,
				JournalDir:                      paths.JournalDir,
				StartedWithoutClaimGraceSeconds: legacyStartedWithoutClaimGraceSeconds,
			})
			if err != nil {
				return renderMaintenanceCheckError(cmd, err)
			}

			view := fidformat.MaintenanceCheckView{
				Status: maintenanceCheckStatus(result, warnAsError),
				Result: result,
			}
			if err := fidformat.Write(cmd.OutOrStdout(), formatName, view); err != nil {
				return renderMaintenanceCheckError(cmd, err)
			}
			if view.Status == "fail" {
				return errMaintenanceCheckFailed
			}
			return nil
		},
	}

	cmd.Flags().BoolVar(&warnAsError, "warn-as-error", false, "Promote warning findings to a failing outcome")
	cmd.Flags().StringVar(&formatName, "format", fidformat.JSON, "Output format")
	return cmd
}

func maintenanceCheckStatus(result maintenance.Result, warnAsError bool) string {
	if len(result.Errors) > 0 || (warnAsError && len(result.Warnings) > 0) {
		return "fail"
	}
	if len(result.Warnings) > 0 {
		return "warn"
	}
	return "pass"
}

func renderMaintenanceCheckError(cmd *cobra.Command, err error) error {
	if renderErr := renderIssueError(cmd.ErrOrStderr(), err); renderErr != nil {
		return renderErr
	}
	return err
}
