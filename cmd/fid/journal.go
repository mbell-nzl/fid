package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/mbell-nzl/fid/internal/config"
	fidformat "github.com/mbell-nzl/fid/internal/format"
	"github.com/mbell-nzl/fid/internal/journal"
	"github.com/spf13/cobra"
)

type journalInspectOutput struct {
	EventID     string           `json:"event_id"`
	EventSchema int              `json:"event_schema"`
	EventType   string           `json:"event_type"`
	JournalFile string           `json:"journal_file"`
	Line        int              `json:"line"`
	Hash        string           `json:"hash"`
	Actor       journal.Actor    `json:"actor"`
	Command     journal.Command  `json:"command"`
	Effects     []journal.Effect `json:"effects"`
	Revert      journal.Verdict  `json:"revert"`
}

func newJournalCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "journal",
		Short: "Inspect journal events.",
	}

	cmd.AddCommand(newJournalInspectCommand())
	return cmd
}

func newJournalInspectCommand() *cobra.Command {
	var formatName string

	cmd := &cobra.Command{
		Use:   "inspect <event-id>",
		Short: "Inspect one journal event by id.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			workspaceRoot, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("resolve workspace root: %w", err)
			}

			found, err := journal.FindEvent(filepath.Join(workspaceRoot, config.WorkspaceDirName, config.JournalDirName), args[0])
			if err != nil {
				if renderErr := renderJournalInspectError(cmd.ErrOrStderr(), err); renderErr != nil {
					return renderErr
				}
				return err
			}

			output := journalInspectOutput{
				EventID:     found.Event.EventID,
				EventSchema: found.Event.EventSchema,
				EventType:   found.Event.EventType,
				JournalFile: found.JournalFile,
				Line:        found.Line,
				Hash:        found.Hash,
				Actor:       found.Event.Actor,
				Command:     found.Event.Command,
				Effects:     found.Event.Effects,
				Revert:      journal.RevertVerdict(found.Event),
			}
			if err := writeJournalInspectOutput(cmd.OutOrStdout(), formatName, output); err != nil {
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

func renderJournalInspectError(writer io.Writer, err error) error {
	var notFoundErr *journal.EventNotFoundError
	if errors.As(err, &notFoundErr) {
		if writeErr := writeIssueLine(writer, "error: %s\n", notFoundErr.Error()); writeErr != nil {
			return writeErr
		}
		return writeIssueLine(writer, "  check the event id and confirm it exists under .fid/journal/\n")
	}
	return renderIssueError(writer, err)
}

func writeJournalInspectOutput(writer io.Writer, formatName string, output journalInspectOutput) error {
	if fidformat.Normalize(formatName) == fidformat.Markdown {
		return writeJournalInspectMarkdown(writer, output)
	}
	return fidformat.Write(writer, formatName, output)
}

func writeJournalInspectMarkdown(writer io.Writer, output journalInspectOutput) error {
	var builder strings.Builder

	builder.WriteString("# Journal Event ")
	builder.WriteString(output.EventID)
	builder.WriteString("\n\n")
	fmt.Fprintf(&builder, "**Schema:** %d\n", output.EventSchema)
	fmt.Fprintf(&builder, "**Type:** %s\n", output.EventType)
	fmt.Fprintf(&builder, "**Location:** %s:%d\n", output.JournalFile, output.Line)
	fmt.Fprintf(&builder, "**Hash:** %s\n", output.Hash)
	builder.WriteString("\n## Actor\n")
	fmt.Fprintf(&builder, "- Role: %s\n", output.Actor.Role)
	fmt.Fprintf(&builder, "- Actor: %s\n", markdownOptionalString(output.Actor.Actor))
	builder.WriteString("\n## Command\n")
	fmt.Fprintf(&builder, "- Argv: %s\n", strings.Join(output.Command.Argv, " "))
	fmt.Fprintf(&builder, "- CWD: %s\n", output.Command.CWD)
	fmt.Fprintf(&builder, "- fid_version: %s\n", output.Command.FIDVersion)
	builder.WriteString("\n## Effects\n")
	for index, effect := range output.Effects {
		fmt.Fprintf(&builder, "- Effect %d: %s\n", index+1, mustMarshalInspectMarkdown(effect))
	}
	builder.WriteString("\n## Revert\n")
	fmt.Fprintf(&builder, "- Eligible: %t\n", output.Revert.Eligible)
	if strings.TrimSpace(output.Revert.Mode) != "" {
		fmt.Fprintf(&builder, "- Mode: %s\n", output.Revert.Mode)
	}
	if len(output.Revert.Restores) > 0 {
		fmt.Fprintf(&builder, "- Restores: %s\n", mustMarshalInspectMarkdown(output.Revert.Restores))
	}
	if len(output.Revert.Precondition) > 0 {
		fmt.Fprintf(&builder, "- Precondition: %s\n", mustMarshalInspectMarkdown(output.Revert.Precondition))
	}
	if strings.TrimSpace(output.Revert.Correction) != "" {
		fmt.Fprintf(&builder, "- Correction: %s\n", output.Revert.Correction)
	}
	if strings.TrimSpace(output.Revert.Note) != "" {
		fmt.Fprintf(&builder, "- Note: %s\n", output.Revert.Note)
	}

	if _, err := io.WriteString(writer, builder.String()); err != nil {
		return fmt.Errorf("write markdown output: %w", err)
	}
	return nil
}

func mustMarshalInspectMarkdown(value any) string {
	payload, err := json.Marshal(value)
	if err != nil {
		return fmt.Sprintf("%v", value)
	}
	return string(payload)
}

func markdownOptionalString(value *string) string {
	if value == nil || strings.TrimSpace(*value) == "" {
		return "(none)"
	}
	return *value
}
