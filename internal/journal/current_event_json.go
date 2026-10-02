package journal

import (
	"fmt"
	"strings"
)

func validateCurrentJSON(event Event, effect Effect, column, raw string) error {
	value, err := decodeCurrentJSON(raw)
	if err != nil {
		return err
	}
	kind := ""
	switch effect.Table + "." + column {
	case "activity_log.details":
		if _, ok := value.(map[string]any); !ok {
			return fmt.Errorf("activity details must be a JSON object")
		}
		return nil
	case "run_observations.observation_json":
		return validateCurrentObservation(value, effect)
	case "run_recoveries.reconciliation_json":
		kind = "Reconciliation"
	case "run_recoveries.authority_json":
		kind = "Authority"
	case "run_recoveries.adopted_paths_json":
		kind = "AdoptedPaths"
	case "claims.scope_json":
		kind = "Scope"
	case "stops.blocked_actions_json":
		kind = "BlockedActions"
	case "stops.resolution_condition_json":
		if err := validateCurrentValue("Condition", value); err != nil {
			return err
		}
		subject := map[string]any{"type": effect.Values["subject_type"], "id": effect.Values["subject_id"]}
		if !currentValuesEqual(value.(map[string]any)["subject"], subject) {
			return fmt.Errorf("stop condition subject differs from row")
		}
		return nil
	case "stops.details_json":
		code, _ := effect.Values["code"].(string)
		shape, ok := currentStopDetailShapes[code]
		if !ok {
			return fmt.Errorf("unknown Stop code")
		}
		_, err := validateCurrentShape(value, shape)
		return err
	case "run_budget_bindings.budget_json":
		kind = "ImmutableBudget"
	case "budget_samples.measurement_ref_json":
		kind = "?Source"
	case "activation_actor_profiles.capabilities_json":
		kind = "[]ID"
	case "activation_gate_definitions.definition_json":
		kind = "GateDefinition"
	case "activation_components.evidence_json":
		row, mergeErr := mergeEffectValues(effect.PK, effect.Values)
		if mergeErr != nil {
			return mergeErr
		}
		switch row["component_type"] {
		case "authority-profile":
			kind = "[]ActorProfile"
		case "gate-definitions":
			kind = "[]GateDefinition"
		default:
			return fmt.Errorf("invalid component type")
		}
	case "execution_receipts.toolchain_json":
		kind = "Toolchain"
	case "execution_receipts.rejected_artifacts_json":
		kind = "[]ArtifactRejection"
	case "review_records.repository_json":
		kind = "Repository"
	case "review_records.coverage_json":
		kind = "Coverage"
	case "review_records.evidence_refs_json":
		kind = "[]Source"
	case "review_records.receipt_ids_json":
		kind = "[]ID"
	case "review_findings.scope_json":
		kind = "Scope"
	case "concerns.scope_json":
		kind = "Scope"
	case "concerns.references_json", "deviations.references_json":
		kind = "[]Source"
	case "decision_refs.source_json":
		kind = "Source"
	case "requirement_refs.source_json":
		kind = "Source"
	case "decision_refs.authority_capabilities_json":
		kind = "[]ID"
	case "human_outcomes.references_json":
		kind = "[]Source"
	case "idempotency_results.result_json":
		return validateCurrentResultJSON(event, value)
	default:
		if column == "start_repository_json" || column == "terminal_repository_json" {
			kind = "Repository"
		} else {
			return fmt.Errorf("unregistered current JSON column %s.%s", effect.Table, column)
		}
	}
	return validateCurrentValue(kind, value)
}
func validateCurrentColumnDomain(column string, value any) error {
	if value == nil {
		return nil
	}
	if column == "request_digest" {
		// The existing Event idempotency grammar admits both hash spellings.
		return validateCurrentValue("Hash", value)
	}
	if strings.HasSuffix(column, "_at") {
		return validateCurrentValue("Timestamp", value)
	}
	if column == "id" || strings.HasSuffix(column, "_id") || column == "idempotency_key" {
		return validateCurrentValue("ID", value)
	}
	if column == "digest" || strings.HasSuffix(column, "_digest") || column == "semantic_key" {
		return validateCurrentValue("Digest", value)
	}
	if column == "declared_path" {
		return validateCurrentValue("Path", value)
	}
	return nil
}
func validateCurrentObservation(value any, effect Effect) error {
	kind := "Observation"
	if effect.Values["kind"] == "checkpoint" {
		kind = "Checkpoint"
	}
	if err := validateCurrentValue(kind, value); err != nil {
		return err
	}
	o := value.(map[string]any)
	if o["id"] != effect.PK["id"] || !currentMatching(o, effect.Values, "run_id", "kind", "actor_id", "observed_at") {
		return fmt.Errorf("observation row binding mismatch")
	}
	return nil
}
func validateCurrentBlockedActions(value any) error {
	if err := validateCurrentValue("[]=claim.acquire,run.start,run.recover,gate.run,issue.transition-terminal-success,envelope.complete", value); err != nil {
		return err
	}
	return currentSortedValues(value, true)
}
func validateCurrentAdoptedPaths(value any) error {
	if err := validateCurrentValue("[]Path", value); err != nil {
		return err
	}
	for _, v := range value.([]any) {
		p := v.(string)
		if p == ".fid" || p == ".git" || strings.HasPrefix(p, ".fid/") || strings.HasPrefix(p, ".git/") {
			return fmt.Errorf("protected adopted path")
		}
	}
	return currentSortedValues(value, false)
}
func validateCurrentResultJSON(event Event, value any) error {
	fields := map[string]string{
		"review.recorded":        "review:Review budget:?Budget stops:[]Stop attention:[]Attention",
		"finding.addressed":      "finding:Finding",
		"concern.recorded":       "concern:Concern",
		"concern.disposed":       "concern:Concern",
		"deviation.recorded":     "deviation:Deviation",
		"decision.referenced":    "decision:DecisionRef protected_decisions:[]ProtectedDecision stops:[]Stop attention:[]Attention",
		"human-outcome.recorded": "outcome:HumanOutcome protected_decisions:[]ProtectedDecision stops:[]Stop attention:[]Attention",
		"requirement.referenced": "requirement:RequirementRef",
		"requirement.linked":     "link:RequirementLink",
		"run.started":            "run:Run", "run.observed": "observation:Observation run:Run", "run.checkpointed": "checkpoint:Checkpoint run:Run",
		"run.recovered": "recovery:Recovery run:Run claim:Claim", "run.abandoned": "run:Run claim:Claim", "run.closed": "run:Run", "claim.expired": "claim:Claim",
		"run.budget-observed": "sample:BudgetSample budget:Budget stops:[]Stop attention:[]Attention", "run.execution-blocked": "denial:Denial",
		"gate.execution-started": "attempt:Attempt receipt:Null", "gate.receipt-recorded": "receipt:Receipt",
		"stop.raised": "stop:Stop attention:Attention", "stop.resolved": "stop:Stop", "attention.resolved": "attention:Attention",
	}[event.EventType]
	if strings.HasPrefix(event.EventType, "activation.") {
		fields = "activation:Activation"
		if event.EventType == "activation.component-staged" {
			fields += " component:Component"
		}
		if event.EventType == "activation.resource-observed" {
			fields += " resource:ActivationResource"
		}
	}
	if fields == "" {
		return fmt.Errorf("unregistered current result family")
	}
	o, err := validateCurrentShape(value, fields+" idempotency_key:ID replayed:False")
	if err != nil {
		return err
	}
	if event.Idempotency != nil && o["idempotency_key"] != event.Idempotency.Key {
		return fmt.Errorf("result idempotency key mismatch")
	}
	return nil
}
