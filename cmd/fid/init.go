package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/mbell-nzl/fid/internal/bootstrap"
	"github.com/mbell-nzl/fid/internal/config"
	"github.com/spf13/cobra"
)

func newInitCommand() *cobra.Command {
	var prefix string

	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize a fid workspace in the current directory.",
		RunE: func(cmd *cobra.Command, args []string) error {
			workspaceRoot, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("resolve workspace root: %w", err)
			}
			if strings.TrimSpace(prefix) == "" {
				prefix = bootstrap.DefaultPrefix(workspaceRoot)
				if err := config.ValidateIssuePrefix(prefix); err != nil {
					validationErr := fmt.Errorf("infer default prefix: %w", err)
					if renderErr := renderInitError(cmd, validationErr); renderErr != nil {
						return renderErr
					}
					return validationErr
				}
			}

			result, err := bootstrap.InitializeWorkspace(bootstrap.InitOptions{
				WorkspaceRoot: workspaceRoot,
				MigrationsDir: initMigrationsDir(),
				Prefix:        prefix,
			})
			if err != nil {
				if renderErr := renderInitError(cmd, err); renderErr != nil {
					return renderErr
				}
				return err
			}

			if err := writeInitLine(cmd.ErrOrStderr(), "initialized fid project in .fid/\n"); err != nil {
				return err
			}
			if err := writeInitLine(cmd.ErrOrStderr(), "  prefix: %s\n", result.Prefix); err != nil {
				return err
			}
			if err := writeInitLine(cmd.ErrOrStderr(), "  config: .fid/config.toml\n"); err != nil {
				return err
			}
			if err := writeInitLine(cmd.ErrOrStderr(), "  pricing: .fid/pricing.toml\n"); err != nil {
				return err
			}
			if err := writeInitLine(cmd.ErrOrStderr(), "  journal: %s\n", result.JournalPath); err != nil {
				return err
			}
			if err := writeInitLine(cmd.ErrOrStderr(), "  event: %s project.initialized\n", result.EventID); err != nil {
				return err
			}
			return writeInitLine(cmd.ErrOrStderr(), "  roles seeded: %s\n", strings.Join(result.SeededRoles, ", "))
		},
	}

	cmd.Flags().StringVar(&prefix, "prefix", "", "Issue ID prefix")
	return cmd
}

func initMigrationsDir() string {
	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		return "migrations"
	}

	return filepath.Clean(filepath.Join(filepath.Dir(currentFile), "..", "..", "migrations"))
}

func renderInitError(cmd *cobra.Command, err error) error {
	var alreadyInitialized *bootstrap.AlreadyInitializedError
	if errors.As(err, &alreadyInitialized) {
		if writeErr := writeInitLine(cmd.ErrOrStderr(), "error: %s\n", alreadyInitialized.Error()); writeErr != nil {
			return writeErr
		}
		return writeInitLine(cmd.ErrOrStderr(), "  to reinitialize, remove .fid/ and run fid init again\n")
	}

	return writeInitLine(cmd.ErrOrStderr(), "error: %v\n", err)
}

func writeInitLine(writer io.Writer, format string, args ...any) error {
	if _, err := fmt.Fprintf(writer, format, args...); err != nil {
		return fmt.Errorf("write init output: %w", err)
	}
	return nil
}
