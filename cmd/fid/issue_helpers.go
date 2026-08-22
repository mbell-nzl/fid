package main

import (
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
	"time"

	"github.com/mbell-nzl/fid/internal/claims"
	"github.com/mbell-nzl/fid/internal/issues"
	"github.com/mbell-nzl/fid/internal/policy"
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

	var permissionDeniedErr *policy.PermissionDeniedError
	if errors.As(err, &permissionDeniedErr) {
		if err := writeIssueLine(writer, "error: %s\n", permissionDeniedErr.Error()); err != nil {
			return err
		}
		return writeIssueLine(writer, "  check `fid policy show` for permitted role/action combinations\n")
	}

	var claimNotHeldErr *claims.ClaimNotHeldError
	if errors.As(err, &claimNotHeldErr) {
		return writeIssueLine(writer, "error: %s\n", claimNotHeldErr.Error())
	}

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

func issueClaimTTLSeconds(cmd *cobra.Command, ttl string) (*int, error) {
	if !cmd.Flags().Changed("ttl") {
		return nil, nil
	}

	duration, err := time.ParseDuration(ttl)
	if err != nil {
		return nil, fmt.Errorf("parse --ttl: %w", err)
	}
	if duration <= 0 {
		return nil, fmt.Errorf("--ttl must be greater than 0s")
	}
	if duration%time.Second != 0 {
		return nil, fmt.Errorf("--ttl must resolve to whole seconds")
	}

	ttlSeconds := int(duration / time.Second)
	return &ttlSeconds, nil
}

func issueSplitCommaSeparated(value string) []string {
	if strings.TrimSpace(value) == "" {
		return nil
	}
	parts := strings.Split(value, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed == "" {
			continue
		}
		result = append(result, trimmed)
	}
	return result
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
