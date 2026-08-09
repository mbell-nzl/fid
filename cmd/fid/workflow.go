package main

import (
	"fmt"
	"os"
	"strings"

	fidformat "github.com/mbell-nzl/fid/internal/format"
	"github.com/mbell-nzl/fid/internal/workflow"
	"github.com/spf13/cobra"
)

func newWorkflowCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "workflow",
		Short: "Define and inspect workflows.",
	}
	cmd.AddCommand(newWorkflowDefineCommand(), newWorkflowAssignCommand(), newWorkflowListCommand())
	return cmd
}

func newWorkflowListCommand() *cobra.Command {
	var formatName, kind string
	var active bool
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List defined workflows.",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			workspaceRoot, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("resolve workspace root: %w", err)
			}
			items, err := workflow.List(workflow.ListOptions{
				WorkspaceRoot: workspaceRoot,
				MigrationsDir: issueMigrationsDir(),
				Kind:          strings.TrimSpace(kind),
				Active:        active,
			})
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
	cmd.Flags().StringVar(&kind, "kind", "", "Only workflows applicable to this kind")
	cmd.Flags().BoolVar(&active, "active", false, "Only active workflows")
	cmd.Flags().StringVar(&formatName, "format", fidformat.JSON, "Output format")
	return cmd
}

func newWorkflowAssignCommand() *cobra.Command {
	var name, kind, actorRole string
	var projectWide bool
	cmd := &cobra.Command{
		Use:   "assign",
		Short: "Assign a workflow.",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			workspaceRoot, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("resolve workspace root: %w", err)
			}
			if err := workflow.AuthorizeAssignment(workspaceRoot, actorRole); err != nil {
				if renderErr := renderIssueError(cmd.ErrOrStderr(), err); renderErr != nil {
					return renderErr
				}
				return err
			}
			result, err := workflow.Assign(workflow.AssignOptions{
				WorkspaceRoot: workspaceRoot,
				MigrationsDir: issueMigrationsDir(),
				ActorRole:     actorRole,
				Request: workflow.AssignmentRequest{
					Name:    strings.TrimSpace(name),
					Kind:    strings.TrimSpace(kind),
					Project: projectWide,
				},
			})
			if err != nil {
				if renderErr := renderIssueError(cmd.ErrOrStderr(), err); renderErr != nil {
					return renderErr
				}
				return err
			}
			if result.Project {
				return writeIssueLine(cmd.ErrOrStderr(), "assigned workflow %s to project\n", result.Name)
			}
			return writeIssueLine(cmd.ErrOrStderr(), "assigned workflow %s to kind %s\n", result.Name, result.Kind)
		},
	}
	cmd.Flags().StringVar(&name, "name", "", "Workflow slug")
	cmd.Flags().StringVar(&kind, "kind", "", "Issue kind")
	cmd.Flags().BoolVar(&projectWide, "project-wide", false, "Assign as project default")
	cmd.Flags().StringVar(&actorRole, "by", "", "Acting role")
	for _, flagName := range []string{"name", "by"} {
		if err := cmd.MarkFlagRequired(flagName); err != nil {
			panic(err)
		}
	}
	return cmd
}

func newWorkflowDefineCommand() *cobra.Command {
	var name, appliesTo, states, categories, transitions, actorRole string
	var requirements []string
	cmd := &cobra.Command{
		Use:   "define",
		Short: "Define a workflow.",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			workspaceRoot, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("resolve workspace root: %w", err)
			}
			if err := workflow.AuthorizeDefinition(workspaceRoot, actorRole); err != nil {
				if renderErr := renderIssueError(cmd.ErrOrStderr(), err); renderErr != nil {
					return renderErr
				}
				return err
			}
			request, err := parseWorkflowDefinition(name, appliesTo, states, categories, transitions, requirements)
			if err != nil {
				if renderErr := renderIssueError(cmd.ErrOrStderr(), err); renderErr != nil {
					return renderErr
				}
				return err
			}
			result, err := workflow.Define(workflow.DefineOptions{
				WorkspaceRoot: workspaceRoot,
				MigrationsDir: issueMigrationsDir(),
				ActorRole:     actorRole,
				Request:       request,
			})
			if err != nil {
				if renderErr := renderIssueError(cmd.ErrOrStderr(), err); renderErr != nil {
					return renderErr
				}
				return err
			}
			return writeIssueLine(cmd.ErrOrStderr(), "defined workflow %s for %s\n  states: %d\n  transitions: %d\n  requirements: %d\n", result.Name, result.AppliesTo, result.States, result.Transitions, result.Requirements)
		},
	}
	cmd.Flags().StringVar(&name, "name", "", "Workflow slug")
	cmd.Flags().StringVar(&appliesTo, "applies-to", "", "Workflow applicability")
	cmd.Flags().StringVar(&states, "states", "", "Ordered comma-separated state names")
	cmd.Flags().StringVar(&categories, "categories", "", "Comma-separated state categories")
	cmd.Flags().StringVar(&transitions, "transitions", "", "Comma-separated workflow transitions")
	cmd.Flags().StringArrayVar(&requirements, "require", nil, "Transition rule requirement")
	cmd.Flags().StringVar(&actorRole, "by", "", "Acting role")
	for _, flagName := range []string{"name", "applies-to", "states", "categories", "transitions", "by"} {
		if err := cmd.MarkFlagRequired(flagName); err != nil {
			panic(err)
		}
	}
	return cmd
}

func parseWorkflowDefinition(name, appliesTo, states, categories, transitions string, requirements []string) (workflow.DefinitionRequest, error) {
	request := workflow.DefinitionRequest{
		Name:       strings.TrimSpace(name),
		AppliesTo:  strings.TrimSpace(appliesTo),
		States:     issueSplitCommaSeparated(states),
		Categories: issueSplitCommaSeparated(categories),
	}
	if !isWorkflowAppliesTo(request.AppliesTo) {
		return workflow.DefinitionRequest{}, fmt.Errorf("invalid workflow applies-to %q", request.AppliesTo)
	}
	parsedTransitions, err := parseWorkflowTransitions(transitions)
	if err != nil {
		return workflow.DefinitionRequest{}, err
	}
	request.Transitions = parsedTransitions
	parsedRequirements, err := parseWorkflowRequirements(requirements)
	if err != nil {
		return workflow.DefinitionRequest{}, err
	}
	request.Requirements = parsedRequirements
	return request, nil
}

func isWorkflowAppliesTo(value string) bool {
	switch value {
	case "epic", "task", "subtask", "qa_review", "escalation", "any":
		return true
	default:
		return false
	}
}

func parseWorkflowTransitions(value string) ([]workflow.Transition, error) {
	parts := issueSplitCommaSeparated(value)
	transitions := make([]workflow.Transition, 0, len(parts))
	for _, part := range parts {
		transition, err := parseWorkflowTransition(part)
		if err != nil {
			return nil, err
		}
		transitions = append(transitions, transition)
	}
	return transitions, nil
}

func parseWorkflowTransition(value string) (workflow.Transition, error) {
	normalized := normalizeWorkflowArrow(value)
	parts := strings.Split(normalized, "->")
	if len(parts) != 2 || strings.TrimSpace(parts[0]) == "" || strings.TrimSpace(parts[1]) == "" {
		return workflow.Transition{}, fmt.Errorf("invalid workflow transition %q", value)
	}
	return workflow.Transition{From: strings.TrimSpace(parts[0]), To: strings.TrimSpace(parts[1])}, nil
}

func parseWorkflowRequirements(values []string) ([]workflow.Requirement, error) {
	requirements := make([]workflow.Requirement, 0, len(values))
	for _, value := range values {
		normalized := normalizeWorkflowArrow(value)
		arrow := strings.Index(normalized, "->")
		if arrow < 0 {
			return nil, fmt.Errorf("invalid workflow requirement %q", value)
		}
		colonOffset := strings.Index(normalized[arrow+2:], ":")
		if colonOffset < 0 {
			return nil, fmt.Errorf("invalid workflow requirement %q", value)
		}
		colon := arrow + 2 + colonOffset
		transition, err := parseWorkflowTransition(normalized[:colon])
		if err != nil {
			return nil, fmt.Errorf("invalid workflow requirement %q", value)
		}
		ruleAndParams := strings.TrimSpace(normalized[colon+1:])
		if ruleAndParams == "" {
			return nil, fmt.Errorf("invalid workflow requirement %q", value)
		}
		rule, paramsJSON := splitWorkflowRuleParams(ruleAndParams)
		requirements = append(requirements, workflow.Requirement{From: transition.From, To: transition.To, Rule: rule, ParamsJSON: paramsJSON})
	}
	return requirements, nil
}

func normalizeWorkflowArrow(value string) string {
	normalized := strings.ReplaceAll(value, "-->", "->")
	normalized = strings.ReplaceAll(normalized, "→", "->")
	return strings.ReplaceAll(normalized, "⟶", "->")
}

func splitWorkflowRuleParams(value string) (string, string) {
	index := strings.IndexAny(value, "[{")
	if index < 0 {
		return strings.TrimSpace(value), ""
	}
	return strings.TrimSpace(value[:index]), strings.TrimSpace(value[index:])
}
