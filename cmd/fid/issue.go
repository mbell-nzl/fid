package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"strings"

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
	cmd.AddCommand(newIssueUpdateCommand())
	cmd.AddCommand(newIssueGetCommand())
	cmd.AddCommand(newIssueListCommand())
	return cmd
}

func newIssueCreateCommand() *cobra.Command {
	var kind string
	var title string
	var actorRole string
	var epicID string
	var parentID string
	var priority int
	var sequence int
	var specMDPath string
	var estTokens int
	var dueDate string
	var labels []string
	var inbox bool
	var ownerRole string
	var ownerActor string
	var scheduling string

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create an issue.",
		RunE: func(cmd *cobra.Command, args []string) error {
			workspaceRoot, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("resolve workspace root: %w", err)
			}

			issueParentID, err := issueCreateParentID(kind, epicID, parentID)
			if err != nil {
				if renderErr := renderIssueError(cmd.ErrOrStderr(), err); renderErr != nil {
					return renderErr
				}
				return err
			}

			result, err := issues.Create(issues.CreateOptions{
				WorkspaceRoot: workspaceRoot,
				MigrationsDir: issueMigrationsDir(),
				Kind:          kind,
				Title:         title,
				ActorRole:     actorRole,
				ParentID:      issueParentID,
				Priority:      issueCreateOptionalInt(cmd, "priority", priority),
				Sequence:      issueCreateOptionalInt(cmd, "sequence", sequence),
				SpecMDPath:    issueCreateOptionalString(cmd, "spec-md", specMDPath),
				EstTokens:     issueCreateOptionalInt(cmd, "est-tokens", estTokens),
				DueDate:       issueCreateOptionalString(cmd, "due-date", dueDate),
				Labels:        labels,
				InInbox:       inbox,
				OwnerRole:     issueCreateOptionalString(cmd, "owner-role", ownerRole),
				OwnerActor:    issueCreateOptionalString(cmd, "owner-actor", ownerActor),
				Scheduling:    scheduling,
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
	cmd.Flags().StringVar(&epicID, "epic", "", "Parent epic ID")
	cmd.Flags().StringVar(&parentID, "parent", "", "Parent task ID")
	cmd.Flags().IntVar(&priority, "priority", 0, "Issue priority")
	cmd.Flags().IntVar(&sequence, "sequence", 0, "Issue sequence")
	cmd.Flags().StringVar(&specMDPath, "spec-md", "", "Spec markdown path")
	cmd.Flags().IntVar(&estTokens, "est-tokens", 0, "Estimated tokens")
	cmd.Flags().StringVar(&dueDate, "due-date", "", "Due date")
	cmd.Flags().StringSliceVar(&labels, "label", nil, "Issue label")
	cmd.Flags().BoolVar(&inbox, "inbox", false, "Add issue to inbox")
	cmd.Flags().StringVar(&ownerRole, "owner-role", "", "Owner role")
	cmd.Flags().StringVar(&ownerActor, "owner-actor", "", "Owner actor")
	cmd.Flags().StringVar(&actorRole, "by", "", "Acting role")
	cmd.Flags().StringVar(&scheduling, "scheduling", "", "Epic scheduling")
	return cmd
}

func newIssueGetCommand() *cobra.Command {
	var formatName string
	var sectionName string

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

			outputValue := any(issue)
			if strings.TrimSpace(sectionName) != "" {
				outputValue, err = issues.SectionValue(issue, sectionName)
				if err != nil {
					if renderErr := renderIssueError(cmd.ErrOrStderr(), err); renderErr != nil {
						return renderErr
					}
					return err
				}
			}

			if err := fidformat.Write(cmd.OutOrStdout(), formatName, outputValue); err != nil {
				if renderErr := renderIssueError(cmd.ErrOrStderr(), err); renderErr != nil {
					return renderErr
				}
				return err
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&formatName, "format", fidformat.JSON, "Output format")
	cmd.Flags().StringVar(&sectionName, "section", "", "Issue detail section")
	return cmd
}

func newIssueUpdateCommand() *cobra.Command {
	var (
		title      string
		priority   int
		labels     []string
		ownerRole  string
		ownerActor string
		specMDPath string
		estTokens  int
		dueDate    string
		actorRole  string
	)

	cmd := &cobra.Command{
		Use:   "update <id>",
		Short: "Update an issue.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			workspaceRoot, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("resolve workspace root: %w", err)
			}

			result, err := issues.Update(issues.UpdateOptions{
				WorkspaceRoot: workspaceRoot,
				MigrationsDir: issueMigrationsDir(),
				ID:            args[0],
				ActorRole:     actorRole,
				Title:         issueCreateOptionalString(cmd, "title", title),
				Priority:      issueCreateOptionalInt(cmd, "priority", priority),
				Labels:        issueOptionalStringSlice(cmd, "labels", labels),
				SpecMDPath:    issueCreateOptionalString(cmd, "spec-md", specMDPath),
				EstTokens:     issueCreateOptionalInt(cmd, "est-tokens", estTokens),
				DueDate:       issueCreateOptionalString(cmd, "due-date", dueDate),
				OwnerRole:     issueCreateOptionalString(cmd, "owner-role", ownerRole),
				OwnerActor:    issueCreateOptionalString(cmd, "owner-actor", ownerActor),
			})
			if err != nil {
				if renderErr := renderIssueError(cmd.ErrOrStderr(), err); renderErr != nil {
					return renderErr
				}
				return err
			}

			return writeIssueLine(cmd.ErrOrStderr(), "updated %s: %s\n", result.ID, issueUpdateSummary(result.Changes))
		},
	}

	cmd.Flags().StringVar(&title, "title", "", "New title")
	cmd.Flags().IntVar(&priority, "priority", 0, "New priority")
	cmd.Flags().StringSliceVar(&labels, "labels", nil, "Replacement issue labels")
	cmd.Flags().StringVar(&ownerRole, "owner-role", "", "New owner role")
	cmd.Flags().StringVar(&ownerActor, "owner-actor", "", "New owner actor")
	cmd.Flags().StringVar(&specMDPath, "spec-md", "", "New spec markdown path")
	cmd.Flags().IntVar(&estTokens, "est-tokens", 0, "New estimated tokens")
	cmd.Flags().StringVar(&dueDate, "due-date", "", "New due date")
	cmd.Flags().StringVar(&actorRole, "by", "", "Acting role")
	return cmd
}

func newIssueListCommand() *cobra.Command {
	var (
		kind           string
		status         string
		statusCategory string
		epicID         string
		label          string
		priorityFilter string
		inboxOnly      bool
		hasConcerns    bool
		order          string
		formatName     string
	)

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List issues.",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			workspaceRoot, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("resolve workspace root: %w", err)
			}

			priorities, err := issues.ParsePriorityFilter(priorityFilter)
			if err != nil {
				if renderErr := renderIssueError(cmd.ErrOrStderr(), err); renderErr != nil {
					return renderErr
				}
				return err
			}

			listOptions := issues.ListOptions{
				WorkspaceRoot:  workspaceRoot,
				MigrationsDir:  issueMigrationsDir(),
				Kind:           kind,
				Status:         status,
				StatusCategory: statusCategory,
				EpicID:         epicID,
				Label:          label,
				Priorities:     priorities,
				Order:          order,
			}
			if cmd.Flags().Changed("inbox") {
				listOptions.InInbox = &inboxOnly
			}
			if cmd.Flags().Changed("has-concerns") {
				listOptions.HasOpenConcerns = &hasConcerns
			}

			items, err := issues.List(listOptions)
			if err != nil {
				if renderErr := renderIssueError(cmd.ErrOrStderr(), err); renderErr != nil {
					return renderErr
				}
				return err
			}

			if err := fidformat.Write(cmd.OutOrStdout(), formatName, items); err != nil {
				if renderErr := renderIssueError(cmd.ErrOrStderr(), err); renderErr != nil {
					return renderErr
				}
				return err
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&kind, "kind", "", "Filter by kind")
	cmd.Flags().StringVar(&status, "status", "", "Filter by exact status")
	cmd.Flags().StringVar(&statusCategory, "status-category", "", "Filter by status category")
	cmd.Flags().StringVar(&epicID, "epic", "", "Filter by parent epic ID")
	cmd.Flags().StringVar(&label, "label", "", "Filter by label")
	cmd.Flags().StringVar(&priorityFilter, "priority", "", "Filter by priority list")
	cmd.Flags().BoolVar(&inboxOnly, "inbox", false, "Show only inbox items")
	cmd.Flags().BoolVar(&hasConcerns, "has-concerns", false, "Show only issues with open concerns")
	cmd.Flags().StringVar(&order, "order", "", "Ordering")
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

	handled, renderErr := renderJournalReplayError(func(format string, args ...any) error {
		return writeIssueLine(writer, format, args...)
	}, err, "read stopped - manual repair required")
	if handled {
		return renderErr
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
	return strings.Join(parts, ", ")
}

func issueUpdateFieldLabel(field string) string {
	switch field {
	case "spec_md_path":
		return "spec-md"
	case "est_tokens":
		return "est-tokens"
	case "due_date":
		return "due-date"
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
