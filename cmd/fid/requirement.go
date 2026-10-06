package main

import (
	"errors"

	"github.com/mbell-nzl/fid/internal/cli"
	"github.com/mbell-nzl/fid/internal/review"
	"github.com/spf13/cobra"
)

func init() {
	for _, leaf := range []string{"record-reference", "link", "get", "list"} {
		rootFlagCapabilities["requirement "+leaf] = rootFlagCapability{Project: true, JSON: true}
	}
}

func newRequirementCommand() *cobra.Command {
	command := &cobra.Command{Use: "requirement", Short: "Record and inspect Requirement References."}
	command.AddCommand(newRequirementMutationCommand("record-reference"), newRequirementMutationCommand("link"), newRequirementGetCommand(), newRequirementListCommand())
	return command
}

func newRequirementMutationCommand(leaf string) *cobra.Command {
	operation := "requirement." + leaf
	var inputPath, actor, key string
	command := &cobra.Command{Use: leaf, Short: "Record a Requirement Reference.", Args: canonicalNoArgs(operation), RunE: func(cmd *cobra.Command, _ []string) error {
		for _, flag := range []string{"input", "actor", "idempotency-key"} {
			if err := requireCanonicalFlag(cmd, operation, flag); err != nil {
				return err
			}
		}
		var recordRequest review.RequirementReferenceRequest
		var linkRequest review.RequirementLinkRequest
		if leaf == "record-reference" {
			if err := readDescriptorInput(cmd.InOrStdin(), inputPath, inputDescriptorRegistry[rootCommandKey(cmd)], &recordRequest); err != nil {
				return renderRequirementError(cmd, operation, err)
			}
		} else {
			if err := readDescriptorInput(cmd.InOrStdin(), inputPath, inputDescriptorRegistry[rootCommandKey(cmd)], &linkRequest); err != nil {
				return renderRequirementError(cmd, operation, err)
			}
		}
		options, err := requirementMutationOptions(cmd, actor, key)
		if err != nil {
			return renderRequirementError(cmd, operation, err)
		}
		var result review.RequirementMutationResult
		if leaf == "record-reference" {
			result, err = review.RecordRequirementReference(options, recordRequest)
		} else {
			result, err = review.LinkRequirement(options, linkRequest)
		}
		if err != nil {
			return renderRequirementError(cmd, operation, err)
		}
		var name string
		var value map[string]any
		if leaf == "record-reference" {
			name, value = "requirement", result.Requirement
		} else {
			name, value = "link", result.Link
		}
		jsonMode, err := readJSONMode(cmd)
		if err != nil {
			return err
		}
		if jsonMode {
			return cli.WriteSuccess(cmd.OutOrStdout(), operation, map[string]any{name: value, "idempotency_key": result.IdempotencyKey, "replayed": result.Replayed}, &cli.Event{ID: result.EventID, OccurredAt: value["created_at"].(string)})
		}
		return writeP23Human(cmd, operation, value)
	}}
	command.Flags().StringVar(&inputPath, "input", "", "JSON input")
	command.Flags().StringVar(&actor, "actor", "", "Finalized activation actor ID")
	command.Flags().StringVar(&key, "idempotency-key", "", "Idempotency key")
	command.SetFlagErrorFunc(canonicalFlagError(operation))
	installInputDescriptor(command, inputDescriptorRegistry["requirement "+leaf])
	return command
}

func requirementMutationOptions(cmd *cobra.Command, actor, key string) (review.JudgmentOptions, error) {
	workspace, err := issueMutationWorkspaceRoot(cmd)
	if err != nil {
		return review.JudgmentOptions{}, err
	}
	return review.JudgmentOptions{WorkspaceRoot: workspace, MigrationsDir: issueMigrationsDir(), ActorID: actor, IdempotencyKey: key}, nil
}

func newRequirementGetCommand() *cobra.Command {
	return &cobra.Command{Use: "get <requirement-id>", Short: "Read a Requirement Reference.", Args: cobra.ExactArgs(1), RunE: func(cmd *cobra.Command, args []string) error {
		workspace, err := issueMutationWorkspaceRoot(cmd)
		if err != nil {
			return renderRequirementError(cmd, "requirement.get", err)
		}
		value, links, err := review.GetRequirement(review.RequirementListOptions{WorkspaceRoot: workspace, MigrationsDir: issueMigrationsDir()}, args[0])
		if err != nil {
			return renderRequirementError(cmd, "requirement.get", err)
		}
		jsonMode, err := readJSONMode(cmd)
		if err != nil {
			return err
		}
		if jsonMode {
			return cli.WriteSuccess(cmd.OutOrStdout(), "requirement.get", map[string]any{"requirement": value, "links": links}, nil)
		}
		return writeP23Human(cmd, "requirement.get", append([]map[string]any{value}, links...)...)
	}}
}

func newRequirementListCommand() *cobra.Command {
	var status, issueID string
	var limit int
	command := &cobra.Command{Use: "list", Short: "List Requirement References.", Args: canonicalNoArgs("requirement.list"), RunE: func(cmd *cobra.Command, _ []string) error {
		workspace, err := issueMutationWorkspaceRoot(cmd)
		if err != nil {
			return renderRequirementError(cmd, "requirement.list", err)
		}
		items, err := review.ListRequirements(review.RequirementListOptions{WorkspaceRoot: workspace, MigrationsDir: issueMigrationsDir(), Status: status, IssueID: issueID, Limit: limit})
		if err != nil {
			return renderRequirementError(cmd, "requirement.list", err)
		}
		jsonMode, err := readJSONMode(cmd)
		if err != nil {
			return err
		}
		if jsonMode {
			return cli.WriteSuccess(cmd.OutOrStdout(), "requirement.list", map[string]any{"items": items}, nil)
		}
		return writeP23Human(cmd, "requirement.list", items...)
	}}
	command.Flags().StringVar(&status, "status", "", "Filter by status")
	command.Flags().StringVar(&issueID, "issue", "", "Filter by linked Issue")
	command.Flags().IntVar(&limit, "limit", 100, "Maximum references to return")
	return command
}

func renderRequirementError(cmd *cobra.Command, operation string, cause error) error {
	var notFound *review.RequirementNotFoundError
	if errors.As(cause, &notFound) {
		jsonMode, err := readJSONMode(cmd)
		if err != nil {
			return err
		}
		return writeClaimFailure(cmd, jsonMode, operation, cause, cli.CodeRequirementNotFound, "Requirement Reference was not found", &cli.Subject{Type: cli.SubjectRequirement, ID: notFound.ID}, map[string]any{"id": notFound.ID})
	}
	return renderDecisionError(cmd, operation, cause)
}
