package main

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/mbell-nzl/fid/internal/backup"
	"github.com/mbell-nzl/fid/internal/cli"
	"github.com/mbell-nzl/fid/internal/project"
	"github.com/spf13/cobra"
)

func init() {
	rootFlagCapabilities["backup"] = rootFlagCapability{Project: true, JSON: true}
}

// backupEnvironmentKey carries a backupEnvironment in the root command's
// context, so tests can give a backup its own HOME, config and cache.
type backupEnvironmentKey struct{}

type backupEnvironment struct {
	environ  []string
	lockWait time.Duration
}

func commandBackupEnvironment(cmd *cobra.Command) backupEnvironment {
	if environment, ok := cmd.Root().Context().Value(backupEnvironmentKey{}).(backupEnvironment); ok {
		return environment
	}
	return backupEnvironment{environ: os.Environ(), lockWait: backup.DefaultLockWait}
}

func newBackupCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "backup",
		Short: "Push the ledger's export to the backup repository.",
		Long: "Export the ledger into the project's folder in the backup repository, commit and push.\n\n" +
			"The target is [backup] target in $XDG_CONFIG_HOME/fid/config.toml; the folder is the project ID in .fid.toml.\n" +
			"fid works in its own clone under $XDG_CACHE_HOME/fid/backup and records every attempt in the ledger's backups.jsonl.",
		RunE: func(cmd *cobra.Command, _ []string) error {
			root, err := issueMutationWorkspaceRoot(cmd)
			if err != nil {
				return renderBackupFailure(cmd, err)
			}
			environment := commandBackupEnvironment(cmd)
			record, recordPath, err := backup.Run(backup.Options{WorkspaceRoot: root, Environment: environment.environ,
				LockWait: environment.lockWait, MigrationsDir: issueMigrationsDir()})
			if err != nil {
				return renderBackupFailure(cmd, err)
			}
			return renderBackupSuccess(cmd, record, recordPath)
		},
	}
}

func renderBackupSuccess(cmd *cobra.Command, record backup.Record, recordPath string) error {
	jsonMode, err := readJSONMode(cmd)
	if err != nil {
		return err
	}
	if jsonMode {
		return cli.WriteSuccess(cmd.OutOrStdout(), "project.backup", map[string]any{"backup": record, "record_path": recordPath}, nil)
	}
	line := fmt.Sprintf("backup pushed: %s/ at commit %s\n", *record.ProjectID, *record.Commit)
	if record.Outcome == backup.OutcomeUnchanged {
		line = fmt.Sprintf("backup unchanged: %s/ is current at commit %s\n", *record.ProjectID, *record.Commit)
	}
	if _, err := fmt.Fprint(cmd.OutOrStdout(), line); err != nil {
		return fmt.Errorf("write backup output: %w", err)
	}
	return nil
}

func renderBackupFailure(cmd *cobra.Command, cause error) error {
	if refused := renderRepositoryRefusal(cmd, cause); refused != nil {
		return refused
	}
	jsonMode, err := readJSONMode(cmd)
	if err != nil {
		return err
	}
	failure, err := backupFailure(cmd, cause)
	if err != nil {
		return err
	}
	if jsonMode {
		return cli.WriteFailureAndReturn(cmd.ErrOrStderr(), "project.backup", failure, cause)
	}
	return writeHumanFailure(cmd, "project.backup", failure, cause)
}

// backupFailure maps a failed backup onto existing failure codes. The
// record keeps the backup's own kind; these keep the CLI contract closed.
func backupFailure(cmd *cobra.Command, cause error) (cli.Failure, error) {
	if errors.Is(cause, project.ErrNotFound) {
		return issueMutationFailure(cmd, "project.backup", "", cause)
	}
	var failed *backup.Error
	if !errors.As(cause, &failed) {
		return cli.NewFailure(cli.CodeInternalError, cli.BoundedMessage("backup failed: "+cause.Error()), nil, map[string]any{})
	}
	invalid := func(allowed string) (cli.Failure, error) {
		return cli.NewFailure(cli.CodeInvalidArgument, cli.BoundedMessage(failed.Error()), nil,
			map[string]any{"argument": failed.Path, "value": nil, "allowed": []string{allowed}})
	}
	unavailable := func(reason string) (cli.Failure, error) {
		return cli.NewFailure(cli.CodeRepositoryUnavailable, cli.BoundedMessage(failed.Message), nil,
			map[string]any{"reason": reason, "detail": failed.Detail})
	}
	switch failed.Kind {
	case backup.KindNoTarget:
		return invalid("[backup] target: a git URL or local path")
	case backup.KindUserConfig:
		return invalid("valid project.backup configuration")
	case backup.KindProjectIdentity:
		return invalid("[project] id matching the ledger's project ID")
	case backup.KindLockUnavailable:
		return cli.NewFailure(cli.CodeLockUnavailable, cli.BoundedMessage(failed.Error()), nil, map[string]any{"lock_path": failed.Path})
	case backup.KindTargetUnreachable:
		return unavailable("not-a-repository")
	case backup.KindPushRejected:
		return unavailable("push-rejected")
	}
	return cli.NewFailure(cli.CodeInternalError, cli.BoundedMessage(failed.Error()), nil, map[string]any{})
}
