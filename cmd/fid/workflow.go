package main

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/mbell-nzl/fid/internal/cli"
	"github.com/mbell-nzl/fid/internal/workflow"
	"github.com/spf13/cobra"
)

func newWorkflowCommand() *cobra.Command {
	cmd := &cobra.Command{Use: "workflow", Short: "Define, assign, and inspect workflows."}
	cmd.AddCommand(newWorkflowDefineCommand(), newWorkflowAssignCommand(), newWorkflowListCommand())
	return cmd
}

func newWorkflowDefineCommand() *cobra.Command {
	var name, appliesTo, states, categories, transitions, actor, key string
	cmd := &cobra.Command{
		Use:   "define",
		Short: "Define a workflow.",
		Args:  canonicalNoArgs("workflow.define"),
		RunE: func(cmd *cobra.Command, _ []string) error {
			for _, flag := range []string{"name", "applies-to", "states", "categories", "transitions"} {
				if err := requireCanonicalFlag(cmd, "workflow.define", flag); err != nil {
					return err
				}
			}
			for _, flag := range []string{"actor", "idempotency-key"} {
				if err := requireCanonicalFlag(cmd, "workflow.define", flag); err != nil {
					return err
				}
			}
			workspaceRoot, err := issueMutationWorkspaceRoot(cmd)
			if err != nil {
				return renderWorkflowError(cmd, "workflow.define", err)
			}
			request, err := parseWorkflowDefinition(name, appliesTo, states, categories, transitions)
			if err != nil {
				return renderWorkflowError(cmd, "workflow.define", err)
			}
			result, err := workflow.Define(workflow.DefineOptions{
				WorkspaceRoot: workspaceRoot, MigrationsDir: issueMigrationsDir(), Actor: actor, IdempotencyKey: key, Request: request,
				RequestInput: map[string]any{"--name": name, "--applies-to": appliesTo, "--states": states, "--categories": categories, "--transitions": transitions},
			})
			if err != nil {
				return renderWorkflowError(cmd, "workflow.define", err)
			}
			jsonMode, err := readJSONMode(cmd)
			if err != nil {
				return err
			}
			if jsonMode {
				return cli.WriteSuccess(cmd.OutOrStdout(), "workflow.define", map[string]any{"workflow": map[string]any{
					"id": result.Name, "name": result.Name, "applies_to": result.AppliesTo,
					"state_count": result.States, "transition_count": result.Transitions,
				}, "idempotency_key": result.IdempotencyKey, "replayed": result.Replayed}, &cli.Event{ID: result.EventID, OccurredAt: result.OccurredAt.Format(time.RFC3339Nano)})
			}
			_, err = fmt.Fprintf(cmd.OutOrStdout(), "defined workflow %s\n", result.Name)
			return err
		},
	}
	cmd.Flags().StringVar(&name, "name", "", "Workflow slug")
	cmd.Flags().StringVar(&appliesTo, "applies-to", "", "Workflow applicability")
	cmd.Flags().StringVar(&states, "states", "", "Comma-separated workflow states")
	cmd.Flags().StringVar(&categories, "categories", "", "Comma-separated state categories")
	cmd.Flags().StringVar(&transitions, "transitions", "", "Comma-separated workflow transitions")
	cmd.Flags().StringVar(&actor, "actor", "", "Standalone attribution")
	cmd.Flags().StringVar(&key, "idempotency-key", "", "Idempotency key")
	cmd.SetFlagErrorFunc(canonicalFlagError("workflow.define"))
	return cmd
}

func newWorkflowAssignCommand() *cobra.Command {
	var name, kind, actor, key string
	var projectWide bool
	cmd := &cobra.Command{
		Use:   "assign",
		Short: "Assign a workflow.",
		Args:  canonicalNoArgs("workflow.assign"),
		RunE: func(cmd *cobra.Command, _ []string) error {
			if err := requireCanonicalFlag(cmd, "workflow.assign", "name"); err != nil {
				return err
			}
			for _, flag := range []string{"actor", "idempotency-key"} {
				if err := requireCanonicalFlag(cmd, "workflow.assign", flag); err != nil {
					return err
				}
			}
			if cmd.Flags().Changed("kind") == projectWide {
				return renderCanonicalRequestError(cmd, "workflow.assign", fmt.Errorf("exactly one of --kind and --project-wide is required"))
			}
			if cmd.Flags().Changed("kind") && !isWorkflowKind(kind) {
				return renderWorkflowError(cmd, "workflow.assign", fmt.Errorf("invalid workflow kind %q", strings.TrimSpace(kind)))
			}
			workspaceRoot, err := issueMutationWorkspaceRoot(cmd)
			if err != nil {
				return renderWorkflowError(cmd, "workflow.assign", err)
			}
			requestInput := map[string]any{"--name": name}
			if projectWide {
				requestInput["--project-wide"] = true
			} else {
				requestInput["--kind"] = kind
			}
			result, err := workflow.Assign(workflow.AssignOptions{
				WorkspaceRoot: workspaceRoot, MigrationsDir: issueMigrationsDir(), Actor: actor, IdempotencyKey: key,
				Request:      workflow.AssignmentRequest{Name: strings.TrimSpace(name), Kind: strings.TrimSpace(kind), Project: projectWide},
				RequestInput: requestInput,
			})
			if err != nil {
				return renderWorkflowError(cmd, "workflow.assign", err)
			}
			jsonMode, err := readJSONMode(cmd)
			if err != nil {
				return err
			}
			scopeValue := any(result.Kind)
			if result.Project {
				scopeValue = nil
			}
			if jsonMode {
				return cli.WriteSuccess(cmd.OutOrStdout(), "workflow.assign", map[string]any{"assignment": map[string]any{
					"scope_type": map[bool]string{true: "project", false: "kind"}[result.Project], "scope_value": scopeValue, "workflow_id": result.Name,
				}, "idempotency_key": result.IdempotencyKey, "replayed": result.Replayed}, &cli.Event{ID: result.EventID, OccurredAt: result.OccurredAt.Format(time.RFC3339Nano)})
			}
			_, err = fmt.Fprintf(cmd.OutOrStdout(), "assigned workflow %s\n", result.Name)
			return err
		},
	}
	cmd.Flags().StringVar(&name, "name", "", "Workflow slug")
	cmd.Flags().StringVar(&kind, "kind", "", "Issue kind")
	cmd.Flags().BoolVar(&projectWide, "project-wide", false, "Assign as project default")
	cmd.Flags().StringVar(&actor, "actor", "", "Standalone attribution")
	cmd.Flags().StringVar(&key, "idempotency-key", "", "Idempotency key")
	cmd.SetFlagErrorFunc(canonicalFlagError("workflow.assign"))
	return cmd
}

func newWorkflowListCommand() *cobra.Command {
	var kind string
	var active bool
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List defined workflows.",
		Args:  canonicalNoArgs("workflow.list"),
		RunE: func(cmd *cobra.Command, _ []string) error {
			if cmd.Flags().Changed("kind") && !isWorkflowKind(kind) {
				return renderWorkflowError(cmd, "workflow.list", fmt.Errorf("invalid workflow kind %q", strings.TrimSpace(kind)))
			}
			workspaceRoot, err := issueMutationWorkspaceRoot(cmd)
			if err != nil {
				return renderWorkflowError(cmd, "workflow.list", err)
			}
			items, err := workflow.List(workflow.ListOptions{WorkspaceRoot: workspaceRoot, MigrationsDir: issueMigrationsDir(), Kind: strings.TrimSpace(kind), Active: active})
			if err != nil {
				return renderWorkflowError(cmd, "workflow.list", err)
			}
			jsonMode, err := readJSONMode(cmd)
			if err != nil {
				return err
			}
			if jsonMode {
				return cli.WriteSuccess(cmd.OutOrStdout(), "workflow.list", map[string]any{"workflows": workflowListResult(items), "count": len(items)}, nil)
			}
			for _, item := range items {
				if _, err := fmt.Fprintln(cmd.OutOrStdout(), item.Name); err != nil {
					return err
				}
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&kind, "kind", "", "Only workflows applicable to this kind")
	cmd.Flags().BoolVar(&active, "active", false, "Only active workflows")
	cmd.SetFlagErrorFunc(canonicalFlagError("workflow.list"))
	return cmd
}

func parseWorkflowDefinition(name, appliesTo, states, categories, transitions string) (workflow.DefinitionRequest, error) {
	request := workflow.DefinitionRequest{Name: strings.TrimSpace(name), AppliesTo: strings.TrimSpace(appliesTo)}
	if !isWorkflowAppliesTo(request.AppliesTo) {
		return workflow.DefinitionRequest{}, fmt.Errorf("invalid workflow applies-to %q", request.AppliesTo)
	}
	var err error
	if request.States, err = parseWorkflowCSV(states); err != nil {
		return workflow.DefinitionRequest{}, err
	}
	if request.Categories, err = parseWorkflowCSV(categories); err != nil {
		return workflow.DefinitionRequest{}, err
	}
	if len(request.States) > 64 {
		return workflow.DefinitionRequest{}, fmt.Errorf("workflow has %d states, maximum is 64", len(request.States))
	}
	request.Transitions, err = parseWorkflowTransitions(transitions)
	if err != nil {
		return workflow.DefinitionRequest{}, err
	}
	if len(request.Transitions) > 4096 {
		return workflow.DefinitionRequest{}, fmt.Errorf("workflow has %d transitions, maximum is 4096", len(request.Transitions))
	}
	return request, nil
}

func parseWorkflowCSV(value string) ([]string, error) {
	parts := strings.Split(value, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			return nil, fmt.Errorf("workflow list contains a blank item")
		}
		result = append(result, part)
	}
	return result, nil
}

func isWorkflowAppliesTo(value string) bool { return value == "any" || isWorkflowKind(value) }

func isWorkflowKind(value string) bool {
	switch strings.TrimSpace(value) {
	case "epic", "task", "subtask":
		return true
	default:
		return false
	}
}

func parseWorkflowTransitions(value string) ([]workflow.Transition, error) {
	parts, err := parseWorkflowCSV(value)
	if err != nil {
		return nil, err
	}
	transitions := make([]workflow.Transition, 0, len(parts))
	for _, part := range parts {
		if strings.Count(part, "->") != 1 {
			return nil, fmt.Errorf("invalid workflow transition %q", part)
		}
		endpoints := strings.Split(part, "->")
		from, to := strings.TrimSpace(endpoints[0]), strings.TrimSpace(endpoints[1])
		if from == "" || to == "" {
			return nil, fmt.Errorf("invalid workflow transition %q", part)
		}
		transitions = append(transitions, workflow.Transition{From: from, To: to})
	}
	return transitions, nil
}

func renderWorkflowError(cmd *cobra.Command, operation string, cause error) error {
	if refused := renderRepositoryRefusal(cmd, cause); refused != nil {
		return refused
	}
	jsonMode, err := readJSONMode(cmd)
	if err != nil {
		return err
	}
	var notFound *workflow.NotFoundError
	if errors.As(cause, &notFound) {
		failure, err := cli.NewFailure(cli.CodeWorkflowNotFound, "workflow was not found", &cli.Subject{Type: cli.SubjectWorkflow, ID: notFound.ID}, map[string]any{"id": notFound.ID})
		if err != nil {
			return err
		}
		if !jsonMode {
			return writeHumanFailure(cmd, operation, failure, cause)
		}
		return cli.WriteFailureAndReturn(cmd.ErrOrStderr(), operation, failure, cause)
	}
	var duplicate *workflow.DuplicateError
	if errors.As(cause, &duplicate) {
		failure, err := cli.NewFailure(cli.CodeDuplicateEntity, "entity already exists", &cli.Subject{Type: cli.SubjectWorkflow, ID: duplicate.ID}, map[string]any{
			"entity_type": "workflow", "identity": map[string]any{"fields": []string{"id"}, "values": []string{duplicate.ID}},
		})
		if err != nil {
			return err
		}
		if !jsonMode {
			return writeHumanFailure(cmd, operation, failure, cause)
		}
		return cli.WriteFailureAndReturn(cmd.ErrOrStderr(), operation, failure, cause)
	}
	var limit *workflow.ResultLimitExceededError
	if errors.As(cause, &limit) {
		failure, err := cli.NewFailure(cli.CodeResultLimitExceeded, "read result exceeds its collection limit", nil, map[string]any{"collection": limit.Collection, "limit": limit.Limit, "observed": limit.Observed})
		if err != nil {
			return err
		}
		if !jsonMode {
			return writeHumanFailure(cmd, operation, failure, cause)
		}
		return cli.WriteFailureAndReturn(cmd.ErrOrStderr(), operation, failure, cause)
	}
	return renderIssueMutationError(cmd, operation, "", cause)
}

func workflowListResult(items []workflow.ListItem) []map[string]any {
	result := make([]map[string]any, 0, len(items))
	for _, item := range items {
		assignments := make([]map[string]any, 0, len(item.AssignedTo))
		for _, assignment := range item.AssignedTo {
			scopeValue := any(assignment.Kind)
			if assignment.Scope == "project" {
				scopeValue = nil
			}
			assignments = append(assignments, map[string]any{"scope_type": assignment.Scope, "scope_value": scopeValue})
		}
		result = append(result, map[string]any{
			"id": item.ID, "name": item.Name, "applies_to": item.AppliesTo, "active": item.Active,
			"states": item.States, "transitions": item.Transitions, "assigned_to": assignments,
		})
	}
	return result
}

func init() {
	rootFlagCapabilities["workflow define"] = rootFlagCapability{Project: true, JSON: true}
	rootFlagCapabilities["workflow assign"] = rootFlagCapability{Project: true, JSON: true}
	rootFlagCapabilities["workflow list"] = rootFlagCapability{Project: true, JSON: true}
}
