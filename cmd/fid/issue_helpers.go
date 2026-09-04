package main

import (
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"runtime"
	"slices"
	"strings"

	"github.com/mbell-nzl/fid/internal/issues"
	"github.com/spf13/cobra"
)

func issueMigrationsDir() string {
	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		return "migrations"
	}

	return filepath.Clean(filepath.Join(filepath.Dir(currentFile), "..", "..", "migrations"))
}

func renderIssueError(writer io.Writer, err error) error {
	handled, renderErr := renderJournalReplayError(func(format string, args ...any) error {
		return writeIssueLine(writer, format, args...)
	}, err, "read stopped - manual repair required")
	if handled {
		return renderErr
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

func issueCreateParentID(kind, epicID, parentID string) (*string, error) {
	trimmedEpicID := strings.TrimSpace(epicID)
	trimmedParentID := strings.TrimSpace(parentID)

	if trimmedEpicID != "" && trimmedParentID != "" {
		return nil, fmt.Errorf("create issue: --epic and --parent are mutually exclusive")
	}

	switch strings.TrimSpace(kind) {
	case "task":
		if trimmedParentID != "" {
			return nil, fmt.Errorf("create issue: --parent is only valid for subtasks")
		}
		if trimmedEpicID == "" {
			return nil, nil
		}
		return &trimmedEpicID, nil
	case "subtask":
		if trimmedEpicID != "" {
			return nil, fmt.Errorf("create issue: --epic is only valid for tasks")
		}
		if trimmedParentID == "" {
			return nil, nil
		}
		return &trimmedParentID, nil
	default:
		if trimmedEpicID != "" {
			return nil, fmt.Errorf("create issue: --epic is only valid for tasks")
		}
		if trimmedParentID != "" {
			return nil, fmt.Errorf("create issue: --parent is only valid for subtasks")
		}
		return nil, nil
	}
}

func issueCreateOptionalInt(cmd *cobra.Command, flagName string, value int) *int {
	if !cmd.Flags().Changed(flagName) {
		return nil
	}
	result := value
	return &result
}

func issueCreateOptionalString(cmd *cobra.Command, flagName, value string) *string {
	if !cmd.Flags().Changed(flagName) {
		return nil
	}
	return &value
}

func issueOptionalStringSlice(cmd *cobra.Command, flagName string, values []string) *[]string {
	if !cmd.Flags().Changed(flagName) {
		return nil
	}
	result := slices.Clone(values)
	return &result
}

func issueUpdateSummary(changes []issues.FieldChange) string {
	if len(changes) == 0 {
		return "no changes"
	}

	parts := make([]string, 0, len(changes))
	for _, change := range changes {
		parts = append(parts, fmt.Sprintf("%s %s -> %s", issueUpdateFieldLabel(change.Field), issueUpdateValue(change.From), issueUpdateValue(change.To)))
	}
	return strings.Join(parts, "\n")
}

func issueUpdateFieldLabel(field string) string {
	switch field {
	case "spec_md_path":
		return "spec-md"
	case "owner_role":
		return "owner-role"
	case "owner_actor":
		return "owner-actor"
	default:
		return field
	}
}

func issueUpdateValue(value any) string {
	switch typed := value.(type) {
	case nil:
		return "(none)"
	case string:
		return fmt.Sprintf("%q", typed)
	case []string:
		return "[" + strings.Join(typed, ", ") + "]"
	default:
		return fmt.Sprintf("%v", typed)
	}
}
