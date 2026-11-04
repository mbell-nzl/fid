package main

import (
	"errors"
	"strings"

	"github.com/mbell-nzl/fid/internal/cli"
	"github.com/mbell-nzl/fid/internal/issues"
	"github.com/mbell-nzl/fid/internal/project"
	"github.com/spf13/cobra"
)

func init() {
	rootFlagCapabilities["issue get"] = rootFlagCapability{Project: true, JSON: true}
	rootFlagCapabilities["issue list"] = rootFlagCapability{Project: true, JSON: true}
}

func newIssueGetCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "get <id>",
		Short: "Get an issue by ID.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			jsonMode, err := readJSONMode(cmd)
			if err != nil {
				return err
			}
			workspaceRoot, err := rebuildWorkspaceRoot(cmd)
			if err != nil {
				return renderIssueReadError(cmd, jsonMode, "issue.get", err)
			}
			issue, err := issues.ReadIssueDetail(issues.GetOptions{WorkspaceRoot: workspaceRoot, MigrationsDir: issueMigrationsDir(), ID: args[0]})
			if err != nil {
				return renderIssueReadError(cmd, jsonMode, "issue.get", err)
			}
			if jsonMode {
				return cli.WriteSuccess(cmd.OutOrStdout(), "issue.get", map[string]any{"issue": issue}, nil)
			}
			return writeIssueLine(cmd.OutOrStdout(), "%s | %s | %s\n", issue.ID, issue.Status, issue.Title)
		},
	}
}

func newIssueListCommand() *cobra.Command {
	var kind, status, statusCategory, epicID, label, priorityFilter, order string
	var hasConcerns bool
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List issues.",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			jsonMode, err := readJSONMode(cmd)
			if err != nil {
				return err
			}
			workspaceRoot, err := rebuildWorkspaceRoot(cmd)
			if err != nil {
				return renderIssueReadError(cmd, jsonMode, "issue.list", err)
			}
			priorities, err := issues.ParsePriorityFilter(priorityFilter)
			if err != nil {
				return renderIssueReadError(cmd, jsonMode, "issue.list", err)
			}
			options := issues.ListOptions{WorkspaceRoot: workspaceRoot, MigrationsDir: issueMigrationsDir(), Kind: kind, Status: status, StatusCategory: statusCategory, EpicID: epicID, Label: label, Priorities: priorities, Order: order}
			if cmd.Flags().Changed("has-concerns") {
				options.HasOpenConcerns = &hasConcerns
			}
			items, err := issues.ListIssueSummaries(options)
			if err != nil {
				return renderIssueReadError(cmd, jsonMode, "issue.list", err)
			}
			if jsonMode {
				return cli.WriteSuccess(cmd.OutOrStdout(), "issue.list", map[string]any{"issues": items, "count": len(items)}, nil)
			}
			for _, item := range items {
				if err := writeIssueLine(cmd.OutOrStdout(), "%s | %s | %s\n", item.ID, item.Status, item.Title); err != nil {
					return err
				}
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
	cmd.Flags().BoolVar(&hasConcerns, "has-concerns", false, "Show only issues with open concerns")
	cmd.Flags().StringVar(&order, "order", "", "Ordering")
	return cmd
}

func readJSONMode(cmd *cobra.Command) (bool, error) {
	return cmd.Root().PersistentFlags().GetBool("json")
}

func renderIssueReadError(cmd *cobra.Command, jsonMode bool, operation string, readErr error) error {
	if refused := renderRepositoryRefusal(cmd, readErr); refused != nil {
		return refused
	}
	failure, err := issueReadFailure(operation, readErr)
	if err != nil {
		return err
	}
	if !jsonMode {
		return writeHumanFailure(cmd, operation, failure, readErr)
	}
	return cli.WriteFailureAndReturn(cmd.ErrOrStderr(), operation, failure, readErr)
}

func issueReadFailure(operation string, readErr error) (cli.Failure, error) {
	if errors.Is(readErr, project.ErrNotFound) {
		return cli.NewFailure(cli.CodeProjectNotFound, "fid project was not found", &cli.Subject{Type: cli.SubjectProject, ID: ""}, map[string]any{"path": ""})
	}
	var notFound *issues.NotFoundError
	if strings.TrimSpace(operation) != "" && errors.As(readErr, &notFound) {
		return cli.NewFailure(cli.CodeIssueNotFound, "issue was not found", &cli.Subject{Type: cli.SubjectIssue, ID: notFound.IssueID}, map[string]any{"id": notFound.IssueID})
	}
	var limit *issues.ResultLimitExceededError
	if errors.As(readErr, &limit) {
		return cli.NewFailure(cli.CodeResultLimitExceeded, "read result exceeds its collection limit", nil, map[string]any{"collection": limit.Collection, "limit": limit.Limit, "observed": limit.Observed})
	}
	if failure, handled, err := journalReplayStoppedFailure(readErr); err != nil {
		return cli.Failure{}, err
	} else if handled {
		return failure, nil
	}
	return genericFailure(readErr, cli.CodeInvalidArgument, "invalid issue read argument", nil, map[string]any{"argument": "request", "value": "", "allowed": []string{"accepted issue read input"}})
}
