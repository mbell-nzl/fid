package main

import (
	"encoding/json"
	"fmt"

	"github.com/mbell-nzl/fid/internal/issues"
	"github.com/spf13/cobra"
)

func newConcernCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "concern",
		Short: "Manage concerns.",
	}

	cmd.AddCommand(newConcernAddressCommand())
	cmd.AddCommand(newConcernListCommand())
	cmd.AddCommand(newConcernRaiseCommand())
	return cmd
}

func init() {
	rootFlagCapabilities["concern raise"] = rootFlagCapability{Project: true, JSON: true}
	rootFlagCapabilities["concern address"] = rootFlagCapability{Project: true, JSON: true}
}

func newConcernListCommand() *cobra.Command {
	var (
		issueID string
		status  string
	)

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List concerns.",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			workspaceRoot, err := commandWorkingDirectory(cmd)
			if err != nil {
				return fmt.Errorf("resolve workspace root: %w", err)
			}

			items, err := issues.ListConcerns(issues.ConcernListOptions{
				WorkspaceRoot: workspaceRoot,
				MigrationsDir: issueMigrationsDir(),
				IssueID:       issueID,
				Status:        status,
			})
			if err != nil {
				if renderErr := renderIssueError(cmd.ErrOrStderr(), err); renderErr != nil {
					return renderErr
				}
				return err
			}

			encoder := json.NewEncoder(cmd.OutOrStdout())
			encoder.SetEscapeHTML(false)
			if err := encoder.Encode(items); err != nil {
				if renderErr := renderIssueError(cmd.ErrOrStderr(), err); renderErr != nil {
					return renderErr
				}
				return err
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&issueID, "issue", "", "Filter by issue ID")
	cmd.Flags().StringVar(&status, "status", "", "Filter by concern status")
	return cmd
}

func newConcernAddressCommand() *cobra.Command {
	var actor, completionID string

	cmd := &cobra.Command{
		Use:   "address <concern-id>",
		Short: "Address an open concern from a later completion.",
		Args:  canonicalExactArgs("concern.address", 1),
		RunE: func(cmd *cobra.Command, _ []string) error {
			return writeCapabilityUnavailable(cmd, "concern.address", "human.outcomes.v1", 23)
		},
	}

	cmd.Flags().StringVar(&actor, "actor", "", "Standalone attribution")
	cmd.Flags().StringVar(&completionID, "in-completion", "", "Completion ID that addressed the concern")
	cmd.SetFlagErrorFunc(canonicalFlagError("concern.address"))
	return cmd
}

func newConcernRaiseCommand() *cobra.Command {
	var actor, issueID, completionID, summary, bodyMDPath string

	cmd := &cobra.Command{
		Use:   "raise",
		Short: "Raise a concern against a completed issue.",
		Args:  canonicalNoArgs("concern.raise"),
		RunE: func(cmd *cobra.Command, _ []string) error {
			return writeCapabilityUnavailable(cmd, "concern.raise", "human.outcomes.v1", 23)
		},
	}

	cmd.Flags().StringVar(&actor, "actor", "", "Standalone attribution")
	cmd.Flags().StringVar(&issueID, "on", "", "Issue ID")
	cmd.Flags().StringVar(&completionID, "in-completion", "", "Completion ID this concern was raised in")
	cmd.Flags().StringVar(&summary, "summary", "", "One-line concern summary")
	cmd.Flags().StringVar(&bodyMDPath, "body-md", "", "Path to detailed concern body")
	cmd.SetFlagErrorFunc(canonicalFlagError("concern.raise"))
	return cmd
}
