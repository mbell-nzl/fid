package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"

	fidformat "github.com/mbell-nzl/fid/internal/format"
	"github.com/mbell-nzl/fid/internal/issues"
	"github.com/mbell-nzl/fid/internal/policy"
	"github.com/spf13/cobra"
)

func newIssueCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "issue",
		Short: "Manage issues.",
	}

	cmd.AddCommand(newIssueCreateCommand())
	cmd.AddCommand(newIssueGetCommand())
	return cmd
}

func newIssueCreateCommand() *cobra.Command {
	var kind string
	var title string
	var actorRole string

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create an issue.",
		RunE: func(cmd *cobra.Command, args []string) error {
			workspaceRoot, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("resolve workspace root: %w", err)
			}

			result, err := issues.Create(issues.CreateOptions{
				WorkspaceRoot: workspaceRoot,
				MigrationsDir: issueMigrationsDir(),
				Kind:          kind,
				Title:         title,
				ActorRole:     actorRole,
			})
			if err != nil {
				if renderErr := renderIssueError(cmd.ErrOrStderr(), err); renderErr != nil {
					return renderErr
				}
				return err
			}

			return writeIssueLine(cmd.ErrOrStderr(), "created issue %s (%s) %q\n", result.ID, result.Kind, result.Title)
		},
	}

	cmd.Flags().StringVar(&kind, "kind", "", "Issue kind")
	cmd.Flags().StringVar(&title, "title", "", "Issue title")
	cmd.Flags().StringVar(&actorRole, "by", "", "Acting role")
	return cmd
}

func newIssueGetCommand() *cobra.Command {
	var formatName string

	cmd := &cobra.Command{
		Use:   "get <id>",
		Short: "Get an issue by ID.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			workspaceRoot, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("resolve workspace root: %w", err)
			}

			issue, err := issues.Get(issues.GetOptions{
				WorkspaceRoot: workspaceRoot,
				MigrationsDir: issueMigrationsDir(),
				ID:            args[0],
			})
			if err != nil {
				if renderErr := renderIssueError(cmd.ErrOrStderr(), err); renderErr != nil {
					return renderErr
				}
				return err
			}

			if err := fidformat.Write(cmd.OutOrStdout(), formatName, issue); err != nil {
				if renderErr := renderIssueError(cmd.ErrOrStderr(), err); renderErr != nil {
					return renderErr
				}
				return err
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&formatName, "format", fidformat.JSON, "Output format")
	return cmd
}

func issueMigrationsDir() string {
	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		return "migrations"
	}

	return filepath.Clean(filepath.Join(filepath.Dir(currentFile), "..", "..", "migrations"))
}

func renderIssueError(writer io.Writer, err error) error {
	var missingRoleErr *policy.MissingRoleError
	if errors.As(err, &missingRoleErr) {
		return writeIssueLine(writer, "error: %s\n", missingRoleErr.Error())
	}

	var unknownRoleErr *policy.UnknownRoleError
	if errors.As(err, &unknownRoleErr) {
		if err := writeIssueLine(writer, "error: %s\n", unknownRoleErr.Error()); err != nil {
			return err
		}
		return writeIssueLine(writer, "  known roles: %s\n", joinCommaSeparated(unknownRoleErr.KnownRoles))
	}

	var unsupportedFormatErr *fidformat.UnsupportedFormatError
	if errors.As(err, &unsupportedFormatErr) {
		return writeIssueLine(writer, "error: %s\n", unsupportedFormatErr.Error())
	}

	var unsupportedKindErr *issues.UnsupportedKindError
	if errors.As(err, &unsupportedKindErr) {
		return writeIssueLine(writer, "error: %s\n", unsupportedKindErr.Error())
	}

	var notFoundErr *issues.NotFoundError
	if errors.As(err, &notFoundErr) {
		return writeIssueLine(writer, "error: %s\n", notFoundErr.Error())
	}

	return writeIssueLine(writer, "error: %v\n", err)
}

func writeIssueLine(writer io.Writer, format string, args ...any) error {
	if _, err := fmt.Fprintf(writer, format, args...); err != nil {
		return fmt.Errorf("write issue output: %w", err)
	}
	return nil
}

func joinCommaSeparated(values []string) string {
	if len(values) == 0 {
		return ""
	}
	result := values[0]
	for _, value := range values[1:] {
		result += ", " + value
	}
	return result
}
