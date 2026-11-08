package main

import (
	"errors"
	"fmt"

	"github.com/mbell-nzl/fid/internal/backup"
	"github.com/mbell-nzl/fid/internal/cli"
	"github.com/mbell-nzl/fid/internal/export"
	"github.com/spf13/cobra"
)

func init() {
	rootFlagCapabilities["restore"] = rootFlagCapability{Project: true, JSON: true}
}

func newRestoreCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "restore",
		Short: "Restore the ledger from the backup repository.",
		Long: "Fetch the project's folder from the backup repository and import it as this repository's ledger.\n\n" +
			"The target is [backup] target in $XDG_CONFIG_HOME/fid/config.toml; the folder is the project ID in .fid.toml.\n" +
			"fid restores only into a repository with no ledger, and never replaces one.",
		RunE: func(cmd *cobra.Command, _ []string) error {
			root, err := discoverWorkspaceRoot(cmd)
			if err != nil {
				return renderRestoreFailure(cmd, root, err)
			}
			environment := commandBackupEnvironment(cmd)
			restored, ledgerPath, err := backup.Restore(backup.RestoreOptions{WorkspaceRoot: root, Environment: environment.environ,
				MigrationsDir: issueMigrationsDir()})
			if err != nil {
				return renderRestoreFailure(cmd, root, err)
			}
			return renderRestoreSuccess(cmd, restored, ledgerPath)
		},
	}
}

func renderRestoreSuccess(cmd *cobra.Command, restored backup.Restored, ledgerPath string) error {
	jsonMode, err := readJSONMode(cmd)
	if err != nil {
		return err
	}
	if jsonMode {
		return cli.WriteSuccess(cmd.OutOrStdout(), "project.restore", map[string]any{"restore": restored, "ledger_path": ledgerPath}, nil)
	}
	if _, err := fmt.Fprintf(cmd.OutOrStdout(), "restored %s/ from commit %s into %s/\n", restored.ProjectID, restored.Commit, ledgerPath); err != nil {
		return fmt.Errorf("write restore output: %w", err)
	}
	return nil
}

func renderRestoreFailure(cmd *cobra.Command, root string, cause error) error {
	if refused := renderRepositoryRefusal(cmd, cause); refused != nil {
		return refused
	}
	jsonMode, err := readJSONMode(cmd)
	if err != nil {
		return err
	}
	failure, err := restoreFailure(root, cause)
	if err != nil {
		return err
	}
	if jsonMode {
		return cli.WriteFailureAndReturn(cmd.ErrOrStderr(), "project.restore", failure, cause)
	}
	return writeHumanFailure(cmd, "project.restore", failure, cause)
}

// restoreFailure maps a failed restore onto failure codes: backup's, plus
// project-already-initialized, as init gives it, and archive-invalid.
func restoreFailure(root string, cause error) (cli.Failure, error) {
	var failed *backup.Error
	if !errors.As(cause, &failed) {
		return cli.NewFailure(cli.CodeInternalError, cli.BoundedMessage("restore failed: "+cause.Error()), nil, map[string]any{})
	}
	switch failed.Kind {
	case backup.KindLedgerExists:
		return cli.NewFailure(cli.CodeProjectAlreadyInit, cli.BoundedMessage(failed.Message), &cli.Subject{Type: cli.SubjectProject, ID: root},
			map[string]any{"path": root})
	case backup.KindProjectIdentity:
		return backupInvalidArgument(failed, "[project] id naming the backup folder")
	case backup.KindNoTarget:
		return backupInvalidArgument(failed, "[backup] target: a git URL or local path")
	case backup.KindUserConfig:
		return backupInvalidArgument(failed, "valid project.restore configuration")
	case backup.KindTargetUnreachable:
		return backupRepositoryUnavailable(failed, "not-a-repository")
	case backup.KindNoBackup:
		return backupRepositoryUnavailable(failed, "no-backup")
	case backup.KindArchiveInvalid:
		var invalid *export.InvalidArchiveError
		if errors.As(failed, &invalid) {
			return cli.NewFailure(cli.CodeArchiveInvalid, cli.BoundedMessage(failed.Error()), &cli.Subject{Type: cli.SubjectExport, ID: failed.Path},
				map[string]any{"archive_path": failed.Path, "checks": archiveChecks(invalid)})
		}
	}
	return cli.NewFailure(cli.CodeInternalError, cli.BoundedMessage(failed.Error()), nil, map[string]any{})
}
