package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/mbell-nzl/fid/internal/bootstrap"
	"github.com/mbell-nzl/fid/internal/cli"
	"github.com/mbell-nzl/fid/internal/config"
	"github.com/spf13/cobra"
)

func init() {
	rootFlagCapabilities["init"] = rootFlagCapability{Project: true, JSON: true}
}

func newInitCommand() *cobra.Command {
	var prefix string

	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize a fid workspace.",
		RunE: func(cmd *cobra.Command, args []string) error {
			workspaceRoot, err := initWorkspaceRoot(cmd)
			if err != nil {
				return renderInitFailure(cmd, workspaceRoot, err)
			}
			if strings.TrimSpace(prefix) == "" {
				prefix = bootstrap.DefaultPrefix(workspaceRoot)
				if err := config.ValidateIssuePrefix(prefix); err != nil {
					return renderInitFailure(cmd, workspaceRoot, &initArgumentError{argument: "prefix", value: prefix, err: err})
				}
			} else if err := config.ValidateIssuePrefix(prefix); err != nil {
				return renderInitFailure(cmd, workspaceRoot, &initArgumentError{argument: "prefix", value: prefix, err: err})
			}

			result, err := bootstrap.InitializeWorkspace(bootstrap.InitOptions{
				WorkspaceRoot: workspaceRoot,
				MigrationsDir: initMigrationsDir(),
				Prefix:        prefix,
			})
			if err != nil {
				return renderInitFailure(cmd, workspaceRoot, err)
			}

			return renderInitSuccess(cmd, workspaceRoot, result)
		},
	}

	cmd.Flags().StringVar(&prefix, "prefix", "", "Issue ID prefix")
	return cmd
}

type initArgumentError struct {
	argument string
	value    string
	err      error
}

func (err *initArgumentError) Error() string {
	if err == nil || err.err == nil {
		return "invalid init argument"
	}
	return fmt.Sprintf("invalid init argument %q: %v", err.argument, err.err)
}

func initWorkspaceRoot(cmd *cobra.Command) (string, error) {
	project, err := cmd.Root().PersistentFlags().GetString("project")
	if err != nil {
		return "", fmt.Errorf("read project flag: %w", err)
	}
	workingDirectory, err := commandWorkingDirectory(cmd)
	if err != nil {
		return "", err
	}
	if strings.TrimSpace(project) == "" {
		project = workingDirectory
	} else if !filepath.IsAbs(project) {
		project = filepath.Join(workingDirectory, project)
	}
	root, err := filepath.Abs(project)
	if err != nil {
		return "", &initArgumentError{argument: "project", value: project, err: err}
	}
	info, err := os.Stat(root)
	if err != nil || !info.IsDir() {
		if err == nil {
			err = fmt.Errorf("not a directory")
		}
		return root, &initArgumentError{argument: "project", value: project, err: err}
	}
	return root, nil
}

func initMigrationsDir() string {
	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		return "migrations"
	}

	return filepath.Clean(filepath.Join(filepath.Dir(currentFile), "..", "..", "migrations"))
}

func renderInitSuccess(cmd *cobra.Command, workspaceRoot string, result bootstrap.InitResult) error {
	jsonMode, err := cmd.Root().PersistentFlags().GetBool("json")
	if err != nil {
		return fmt.Errorf("read json flag: %w", err)
	}
	if jsonMode {
		return cli.WriteSuccess(cmd.OutOrStdout(), "project.init", map[string]any{
			"project": map[string]any{
				"id":           result.ProjectID,
				"issue_prefix": result.Prefix,
				"root":         workspaceRoot,
			},
			"journal_path": result.JournalPath,
		}, &cli.Event{ID: result.EventID, OccurredAt: result.OccurredAt.UTC().Format(time.RFC3339Nano)})
	}
	return writeInitLine(cmd.OutOrStdout(), "initialized fid project in .fid/\n")
}

func renderInitFailure(cmd *cobra.Command, workspaceRoot string, initErr error) error {
	jsonMode, err := cmd.Root().PersistentFlags().GetBool("json")
	if err != nil {
		return fmt.Errorf("read json flag: %w", err)
	}
	failure, err := initFailure(workspaceRoot, initErr)
	if err != nil {
		return err
	}
	if jsonMode {
		if err := cli.WriteFailure(cmd.ErrOrStderr(), "project.init", failure); err != nil {
			return err
		}
		return initErr
	}
	if err := renderInitError(cmd, initErr); err != nil {
		return err
	}
	return initErr
}

func initFailure(workspaceRoot string, initErr error) (cli.Failure, error) {
	var alreadyInitialized *bootstrap.AlreadyInitializedError
	if errors.As(initErr, &alreadyInitialized) {
		return cli.NewFailure(cli.CodeProjectAlreadyInit, "project is already initialized", &cli.Subject{
			Type: cli.SubjectProject,
			ID:   workspaceRoot,
		}, map[string]any{"path": workspaceRoot})
	}
	var argumentErr *initArgumentError
	if errors.As(initErr, &argumentErr) {
		return cli.NewFailure(cli.CodeInvalidArgument, "invalid project initialization argument", nil, map[string]any{
			"argument": argumentErr.argument,
			"value":    argumentErr.value,
			"allowed":  []string{"valid project.init input"},
		})
	}
	return cli.NewFailure(cli.CodeInternalError, "project initialization failed", nil, map[string]any{})
}

func renderInitError(cmd *cobra.Command, initErr error) error {
	var alreadyInitialized *bootstrap.AlreadyInitializedError
	if errors.As(initErr, &alreadyInitialized) {
		if err := writeInitLine(cmd.ErrOrStderr(), "error: project is already initialized\n"); err != nil {
			return err
		}
		return writeInitLine(cmd.ErrOrStderr(), "  choose a new directory or use the existing project\n")
	}
	var argumentErr *initArgumentError
	if errors.As(initErr, &argumentErr) {
		return writeInitLine(cmd.ErrOrStderr(), "error: invalid project initialization argument\n")
	}
	return writeInitLine(cmd.ErrOrStderr(), "error: project initialization failed\n")
}

func writeInitLine(writer io.Writer, format string, args ...any) error {
	if _, err := fmt.Fprintf(writer, format, args...); err != nil {
		return fmt.Errorf("write init output: %w", err)
	}
	return nil
}
