package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	internalactivity "github.com/mbell-nzl/fid/internal/activity"
	fidformat "github.com/mbell-nzl/fid/internal/format"
	"github.com/spf13/cobra"
)

func newActivityCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "activity",
		Short: "Inspect the projected activity timeline.",
	}

	cmd.AddCommand(newActivityListCommand())
	return cmd
}

func newActivityListCommand() *cobra.Command {
	var (
		on         string
		since      string
		eventName  string
		formatName string
	)

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List projected activity entries.",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			workspaceRoot, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("resolve workspace root: %w", err)
			}

			items, err := internalactivity.List(internalactivity.ListOptions{
				WorkspaceRoot: workspaceRoot,
				MigrationsDir: issueMigrationsDir(),
				On:            on,
				Since:         since,
				Event:         eventName,
			})
			if err != nil {
				if renderErr := renderIssueError(cmd.ErrOrStderr(), err); renderErr != nil {
					return renderErr
				}
				return err
			}

			if err := writeActivityListOutput(cmd, formatName, items); err != nil {
				if renderErr := renderIssueError(cmd.ErrOrStderr(), err); renderErr != nil {
					return renderErr
				}
				return err
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&on, "on", "", "Filter to one entity reference, e.g. issue/FID-a1b2c3")
	cmd.Flags().StringVar(&since, "since", "", "Include entries on or after this timestamp")
	cmd.Flags().StringVar(&eventName, "event", "", "Filter by projected activity event name")
	cmd.Flags().StringVar(&formatName, "format", fidformat.JSON, "Output format")
	return cmd
}

func writeActivityListOutput(cmd *cobra.Command, formatName string, items []internalactivity.Entry) error {
	if fidformat.Normalize(formatName) == fidformat.Markdown {
		return writeActivityMarkdown(cmd.OutOrStdout(), items)
	}
	return fidformat.Write(cmd.OutOrStdout(), formatName, items)
}

func writeActivityMarkdown(writer interface{ Write([]byte) (int, error) }, items []internalactivity.Entry) error {
	var builder strings.Builder

	builder.WriteString("# Activity\n\n")
	builder.WriteString("| Occurred At | Event | Ref | Actor | ID | Details |\n")
	builder.WriteString("| --- | --- | --- | --- | --- | --- |\n")
	for _, item := range items {
		fmt.Fprintf(
			&builder,
			"| %s | %s | %s/%s | %s | %s | %s |\n",
			markdownCell(item.OccurredAt),
			markdownCell(item.Event),
			markdownCell(item.RefType),
			markdownCell(item.RefID),
			markdownCell(item.Actor),
			markdownCell(item.ID),
			markdownCell(activityDetailsMarkdown(item.Details)),
		)
	}

	if _, err := writer.Write([]byte(builder.String())); err != nil {
		return fmt.Errorf("write markdown output: %w", err)
	}
	return nil
}

func activityDetailsMarkdown(details any) string {
	if details == nil {
		return "null"
	}

	payload, err := json.Marshal(details)
	if err != nil {
		return fmt.Sprintf("%v", details)
	}
	return string(payload)
}

func markdownCell(value string) string {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return "(none)"
	}
	return strings.ReplaceAll(trimmed, "|", "\\|")
}
