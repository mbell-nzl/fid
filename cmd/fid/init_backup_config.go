package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/mbell-nzl/fid/internal/cli"
	"github.com/mbell-nzl/fid/internal/config"
	"github.com/spf13/cobra"
)

// initBackupConfigPath returns a path only if init should ask for a target.
// Non-interactive and JSON runs do not even open the user config.
func initBackupConfigPath(cmd *cobra.Command, isTerminal func(uintptr) bool) (string, error) {
	jsonMode, err := cmd.Root().PersistentFlags().GetBool("json")
	if err != nil {
		return "", fmt.Errorf("read json flag: %w", err)
	}
	if jsonMode || !initIsTerminal(cmd.InOrStdin(), isTerminal) || !initIsTerminal(cmd.OutOrStdout(), isTerminal) {
		return "", nil
	}
	path, err := config.UserConfigPath()
	if err != nil {
		return "", err
	}
	user, err := config.ReadUserConfig(path)
	if err != nil {
		return "", err
	}
	if strings.TrimSpace(user.BackupTarget) != "" {
		return "", nil
	}
	return path, nil
}

func initIsTerminal(stream any, isTerminal func(uintptr) bool) bool {
	file, ok := stream.(interface{ Fd() uintptr })
	return ok && isTerminal(file.Fd())
}

func promptInitBackupTarget(cmd *cobra.Command, path string) error {
	if path == "" {
		return nil
	}
	if err := writeInitLine(cmd.OutOrStdout(), "Backup repository (git URL or local path; empty to skip): "); err != nil {
		return err
	}
	answer, err := bufio.NewReader(cmd.InOrStdin()).ReadString('\n')
	if err != nil && !errors.Is(err, io.EOF) {
		return fmt.Errorf("read backup target: %w", err)
	}
	answer = strings.TrimSuffix(strings.TrimSuffix(answer, "\n"), "\r")
	if strings.TrimSpace(answer) == "" {
		return nil
	}
	if err := config.SetBackupTarget(path, answer); err != nil {
		return err
	}
	return writeInitLine(cmd.OutOrStdout(), "Backup target saved in %s\n", path)
}

// Only Human init can ask for a target. The success lines have already
// reported the identity, so this error describes the failed setup alone.
func renderInitBackupFailure(cmd *cobra.Command, cause error) error {
	failure, err := cli.NewFailure(cli.CodeInternalError,
		cli.BoundedMessage(fmt.Sprintf("ledger initialized, but backup target setup failed: %v", cause)), nil, map[string]any{})
	if err != nil {
		return err
	}
	return writeHumanFailure(cmd, "project.init", failure, cause)
}
