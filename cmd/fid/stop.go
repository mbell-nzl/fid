package main

import (
	"errors"
	"fmt"

	"github.com/mbell-nzl/fid/internal/cli"
	"github.com/mbell-nzl/fid/internal/journal"
	"github.com/mbell-nzl/fid/internal/workflow"
	"github.com/spf13/cobra"
)

func init() {
	for _, operation := range []string{"stop raise", "stop resolve", "stop list"} {
		rootFlagCapabilities[operation] = rootFlagCapability{Project: true, JSON: true}
	}
}

func newStopCommand() *cobra.Command {
	cmd := &cobra.Command{Use: "stop", Short: "Manage durable Stops."}
	cmd.AddCommand(newStopRaiseCommand(), newStopResolveCommand(), newStopListCommand())
	return cmd
}

func newStopRaiseCommand() *cobra.Command {
	var input, actor, key string
	cmd := &cobra.Command{Use: "raise", Short: "Raise a durable Stop.", Args: canonicalNoArgs("stop.raise"), RunE: func(cmd *cobra.Command, _ []string) error {
		var request workflow.StopRaiseRequest
		if err := stopRequest(cmd, "stop.raise", input, actor, key, &request); err != nil {
			return err
		}
		workspace, err := issueMutationWorkspaceRoot(cmd)
		if err != nil {
			return renderStopError(cmd, "stop.raise", err)
		}
		result, err := workflow.RaiseStop(workflow.StopOptions{WorkspaceRoot: workspace, MigrationsDir: issueMigrationsDir(), ActorID: actor, IdempotencyKey: key}, request)
		if err != nil {
			return renderStopError(cmd, "stop.raise", err)
		}
		return writeStopSuccess(cmd, "stop.raise", map[string]any{"stop": result.Stop, "attention": result.Attention, "idempotency_key": result.IdempotencyKey, "replayed": result.Replayed}, result.EventID, result.Stop.RaisedAt)
	}}
	runMutationFlags(cmd, &input, &actor, &key)
	installInputDescriptor(cmd, inputDescriptorRegistry["stop raise"])
	return cmd
}

func newStopResolveCommand() *cobra.Command {
	var input, actor, key string
	cmd := &cobra.Command{Use: "resolve", Short: "Resolve a Stop with a qualifying fact.", Args: canonicalNoArgs("stop.resolve"), RunE: func(cmd *cobra.Command, _ []string) error {
		var request workflow.StopResolveRequest
		if err := stopRequest(cmd, "stop.resolve", input, actor, key, &request); err != nil {
			return err
		}
		workspace, err := issueMutationWorkspaceRoot(cmd)
		if err != nil {
			return renderStopError(cmd, "stop.resolve", err)
		}
		result, err := workflow.ResolveStop(workflow.StopOptions{WorkspaceRoot: workspace, MigrationsDir: issueMigrationsDir(), ActorID: actor, IdempotencyKey: key}, request)
		if err != nil {
			return renderStopError(cmd, "stop.resolve", err)
		}
		return writeStopSuccess(cmd, "stop.resolve", map[string]any{"stop": result.Stop, "idempotency_key": result.IdempotencyKey, "replayed": result.Replayed}, result.EventID, result.Stop.Resolution.ResolvedAt)
	}}
	runMutationFlags(cmd, &input, &actor, &key)
	installInputDescriptor(cmd, inputDescriptorRegistry["stop resolve"])
	return cmd
}

func newStopListCommand() *cobra.Command {
	var subject, status string
	var limit int
	cmd := &cobra.Command{Use: "list", Short: "List Stops.", Args: canonicalNoArgs("stop.list"), RunE: func(cmd *cobra.Command, _ []string) error {
		workspace, err := issueMutationWorkspaceRoot(cmd)
		if err != nil {
			return renderStopError(cmd, "stop.list", err)
		}
		items, err := workflow.ListStops(workflow.StopListOptions{WorkspaceRoot: workspace, MigrationsDir: issueMigrationsDir(), Subject: subject, Status: status, Limit: limit})
		if err != nil {
			return renderStopError(cmd, "stop.list", err)
		}
		return writeStopSuccess(cmd, "stop.list", map[string]any{"items": items}, "", "")
	}}
	cmd.Flags().StringVar(&subject, "subject", "", "Subject type and ID")
	cmd.Flags().StringVar(&status, "status", "", "Stop status")
	cmd.Flags().IntVar(&limit, "limit", 100, "Maximum Stops")
	cmd.SetFlagErrorFunc(canonicalFlagError("stop.list"))
	return cmd
}

func stopRequest(cmd *cobra.Command, operation, input, actor, key string, target any) error {
	for _, flag := range []string{"input", "actor", "idempotency-key"} {
		if err := requireCanonicalFlag(cmd, operation, flag); err != nil {
			return err
		}
	}
	descriptor := inputDescriptorRegistry[rootCommandKey(cmd)]
	if err := readDescriptorInput(cmd.InOrStdin(), input, descriptor, target); err != nil {
		var structural *cli.InputError
		if errors.As(err, &structural) {
			jsonMode, modeErr := readJSONMode(cmd)
			if modeErr != nil {
				return modeErr
			}
			return writeStopFailure(cmd, jsonMode, operation, err, cli.CodeInvalidArgument, err.Error(), nil, structural.Details())
		}
		return renderStopError(cmd, operation, err)
	}
	return nil
}

func writeStopSuccess(cmd *cobra.Command, operation string, result map[string]any, eventID, occurredAt string) error {
	jsonMode, err := readJSONMode(cmd)
	if err != nil {
		return err
	}
	if jsonMode {
		var event *cli.Event
		if eventID != "" {
			event = &cli.Event{ID: eventID, OccurredAt: occurredAt}
		}
		return cli.WriteSuccess(cmd.OutOrStdout(), operation, result, event)
	}
	return writeStopHuman(cmd, operation, result)
}

func renderStopError(cmd *cobra.Command, operation string, cause error) error {
	jsonMode, err := readJSONMode(cmd)
	if err != nil {
		return err
	}
	var idempotency *journal.IdempotencyConflictError
	if errors.As(cause, &idempotency) {
		return writeStopFailure(cmd, jsonMode, operation, cause, cli.CodeIdempotencyConflict, "idempotency key is bound to another request", nil, idempotencyConflictDetails(idempotency))
	}
	var stopNotFound *workflow.StopNotFoundError
	if errors.As(cause, &stopNotFound) {
		return writeStopFailure(cmd, jsonMode, operation, cause, cli.CodeStopNotFound, "Stop was not found", &cli.Subject{Type: cli.SubjectStop, ID: stopNotFound.ID}, map[string]any{"id": stopNotFound.ID})
	}
	var attentionNotFound *workflow.AttentionNotFoundError
	if errors.As(cause, &attentionNotFound) {
		return writeStopFailure(cmd, jsonMode, operation, cause, cli.CodeAttentionNotFound, "Attention was not found", &cli.Subject{Type: cli.SubjectAttention, ID: attentionNotFound.ID}, map[string]any{"id": attentionNotFound.ID})
	}
	var attentionVersion *workflow.AttentionVersionConflictError
	if errors.As(cause, &attentionVersion) {
		return writeStopFailure(cmd, jsonMode, operation, cause, cli.CodeAttentionVersionConflict, "Attention version conflict", &cli.Subject{Type: cli.SubjectAttention, ID: attentionVersion.AttentionID}, map[string]any{"attention_id": attentionVersion.AttentionID, "expected_version": attentionVersion.ExpectedVersion, "actual_version": attentionVersion.ActualVersion})
	}
	var unresolved *workflow.StopBlockedError
	if errors.As(cause, &unresolved) {
		if len(unresolved.Stops) > 0 && unresolved.Stops[0].Code == "review-rework-exhausted" && unresolved.Budget != nil {
			ids := make([]string, 0, len(unresolved.Stops))
			for _, stop := range unresolved.Stops {
				ids = append(ids, stop.ID)
			}
			budget := map[string]any{
				"name": unresolved.Budget.Name, "scope": unresolved.Budget.Scope, "unit": unresolved.Budget.Unit,
				"maximum": unresolved.Budget.Maximum, "exhaustion_action": unresolved.Budget.ExhaustionAction,
				"warning_threshold": optionalBudgetValue(unresolved.Budget.WarningThreshold), "consumed": optionalBudgetValue(unresolved.Budget.Consumed),
				"remaining": optionalBudgetValue(unresolved.Budget.Remaining), "state": unresolved.Budget.State,
			}
			return writeStopFailure(cmd, jsonMode, operation, cause, cli.CodeReviewReworkExhausted, "review rework budget is exhausted", &cli.Subject{Type: cli.SubjectEnvelope, ID: unresolved.Stops[0].Subject.ID}, map[string]any{"budget": budget, "stop_ids": ids, "attention_ids": unresolved.AttentionIDs, "event": nil, "replayed": false})
		}
		ids := make([]string, 0, len(unresolved.Stops))
		reasons := make([]map[string]any, 0, len(unresolved.Stops))
		for _, stop := range unresolved.Stops {
			ids = append(ids, stop.ID)
			reasons = append(reasons, map[string]any{"code": stop.Code, "detail": unresolved.Action, "subject_id": stop.ID, "subject_type": "stop"})
		}
		return writeStopFailure(cmd, jsonMode, operation, cause, cli.CodeStopUnresolved, "Stop remains unresolved", &cli.Subject{Type: cli.SubjectStop, ID: nil}, map[string]any{"stop_ids": ids, "reasons": reasons})
	}
	var resolutionInvalid *workflow.ResolutionInvalidError
	if errors.As(cause, &resolutionInvalid) {
		kind := cli.SubjectStop
		if operation == "attention.resolve" {
			kind = cli.SubjectAttention
		}
		return writeStopFailure(cmd, jsonMode, operation, cause, cli.CodeResolutionInvalid, "Stop resolution is invalid", &cli.Subject{Type: kind, ID: nil}, map[string]any{"subject": map[string]any{"type": resolutionInvalid.Subject.Type, "id": resolutionInvalid.Subject.ID}, "fact_event_ids": resolutionInvalid.FactEventIDs, "reasons": cliEvaluationReasons(resolutionInvalid.Reasons)})
	}
	return writeStopFailure(cmd, jsonMode, operation, cause, cli.CodeInvalidArgument, "invalid Stop request", nil, map[string]any{"argument": "input", "value": nil, "allowed": []string{"accepted " + operation + " input"}})
}

func writeStopFailure(cmd *cobra.Command, jsonMode bool, operation string, cause error, code cli.Code, message string, subject *cli.Subject, details map[string]any) error {
	failure, err := cli.NewFailure(code, message, subject, details)
	if err != nil {
		return err
	}
	if jsonMode {
		return cli.WriteFailureAndReturn(cmd.ErrOrStderr(), operation, failure, cause)
	}
	if _, err := fmt.Fprintln(cmd.ErrOrStderr(), "error: "+failure.Message); err != nil {
		return err
	}
	if code == cli.CodeReviewReworkExhausted && subject != nil {
		if _, err := fmt.Fprintln(cmd.ErrOrStderr(), "next: an activated Human records an approved Human Outcome of class review-rework-extension for this Envelope's Issue; see fid stop list --subject envelope:"+fmt.Sprint(subject.ID)); err != nil {
			return err
		}
	}
	return cli.NewFailureError(failure, cause)
}

func cliEvaluationReasons(reasons []workflow.EvaluationReason) []map[string]any {
	result := make([]map[string]any, 0, len(reasons))
	for _, reason := range reasons {
		result = append(result, map[string]any{"code": reason.Code, "detail": reason.Detail, "subject_id": reason.SubjectID, "subject_type": reason.SubjectType})
	}
	return result
}
