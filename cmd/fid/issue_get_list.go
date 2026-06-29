package main

import (
	"fmt"
	"os"
	"strings"

	fidformat "github.com/mbell-nzl/fid/internal/format"
	"github.com/mbell-nzl/fid/internal/issues"
	"github.com/spf13/cobra"
)

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
