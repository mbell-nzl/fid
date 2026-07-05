package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	fidformat "github.com/mbell-nzl/fid/internal/format"
	internalsqlquery "github.com/mbell-nzl/fid/internal/sqlquery"
	"github.com/spf13/cobra"
)

func newSQLCommand() *cobra.Command {
	var formatName string

	cmd := &cobra.Command{
		Use:   "sql <query>",
		Short: "Run a read-only SQL query against the projected database.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			workspaceRoot, err := os.Getwd()
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

			if err := writeSQLCommandOutput(cmd.OutOrStdout(), formatName, result); err != nil {
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

func writeSQLCommandOutput(writer io.Writer, formatName string, result internalsqlquery.Result) error {
	if fidformat.Normalize(formatName) == fidformat.Markdown {
		return writeSQLMarkdown(writer, result)
	}
	return fidformat.Write(writer, formatName, result.Rows)
}

func writeSQLMarkdown(writer io.Writer, result internalsqlquery.Result) error {
	var builder strings.Builder

	builder.WriteString("# SQL Query Results\n\n")
	builder.WriteString("| ")
	builder.WriteString(strings.Join(result.Columns, " | "))
	builder.WriteString(" |\n")
	builder.WriteString("| ")
	builder.WriteString(strings.TrimRight(strings.Repeat("--- | ", len(result.Columns)), " "))
	builder.WriteString("\n")
	for _, row := range result.Rows {
		builder.WriteString("| ")
		for index, column := range result.Columns {
			if index > 0 {
				builder.WriteString(" | ")
			}
			builder.WriteString(sqlMarkdownCell(row[column]))
		}
		builder.WriteString(" |\n")
	}

	if _, err := io.WriteString(writer, builder.String()); err != nil {
		return fmt.Errorf("write markdown output: %w", err)
	}
	return nil
}

func sqlMarkdownCell(value any) string {
	switch typed := value.(type) {
	case nil:
		return "null"
	case string:
		return strings.ReplaceAll(typed, "|", "\\|")
	default:
		payload, err := json.Marshal(typed)
		if err != nil {
			return strings.ReplaceAll(fmt.Sprint(typed), "|", "\\|")
		}
		text := string(payload)
		text = strings.TrimPrefix(text, "\"")
		text = strings.TrimSuffix(text, "\"")
		return strings.ReplaceAll(text, "|", "\\|")
	}
}
