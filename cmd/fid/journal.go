package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"path/filepath"

	"github.com/mbell-nzl/fid/internal/config"
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
	cmd := &cobra.Command{
		Use:   "inspect <event-id>",
		Short: "Inspect one journal event by id.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			workspaceRoot, err := discoverWorkspaceRoot(cmd)
			if refused := renderRepositoryRefusal(cmd, err); refused != nil {
				return refused
			} else if err != nil {
				return fmt.Errorf("resolve workspace root: %w", err)
			}

			found, err := journal.FindEvent(filepath.Join(config.LedgerDir(workspaceRoot), config.JournalDirName), args[0])
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
			}
			encoder := json.NewEncoder(cmd.OutOrStdout())
			encoder.SetEscapeHTML(false)
			if err := encoder.Encode(output); err != nil {
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
