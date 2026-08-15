package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"

	internalsqlquery "github.com/mbell-nzl/fid/internal/sqlquery"
	"github.com/spf13/cobra"
)

func newSQLCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sql <query>",
		Short: "Run a read-only SQL query against the projected database.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			workspaceRoot, err := commandWorkingDirectory(cmd)
			if err != nil {
				return fmt.Errorf("resolve workspace root: %w", err)
			}

			result, err := internalsqlquery.Run(internalsqlquery.Options{
				WorkspaceRoot: workspaceRoot,
				MigrationsDir: issueMigrationsDir(),
				Query:         args[0],
			})
			if err != nil {
				if renderErr := renderSQLCommandError(cmd.ErrOrStderr(), err); renderErr != nil {
					return renderErr
				}
				return err
			}

			encoder := json.NewEncoder(cmd.OutOrStdout())
			encoder.SetEscapeHTML(false)
			if err := encoder.Encode(result.Rows); err != nil {
				if renderErr := renderIssueError(cmd.ErrOrStderr(), err); renderErr != nil {
					return renderErr
				}
				return err
			}
			return nil
		},
	}

	return cmd
}

func renderSQLCommandError(writer io.Writer, err error) error {
	var rejectedErr *internalsqlquery.RejectionError
	if errors.As(err, &rejectedErr) {
		if writeErr := writeIssueLine(writer, "error: %s\n", rejectedErr.Error()); writeErr != nil {
			return writeErr
		}
		return writeIssueLine(writer, "  rejected: %s\n", rejectedErr.Statement)
	}
	return renderIssueError(writer, err)
}
