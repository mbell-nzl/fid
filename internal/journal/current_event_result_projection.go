package journal

import (
	"database/sql"
	"fmt"
	"strings"
)

// Read only the projection inside the Event transaction, after its Effects but
// before commit. A mismatching cached response rolls the complete Event back.
func validateCurrentAppliedResults(tx *sql.Tx, event Event) error {
	if strings.HasPrefix(event.EventType, "check.") {
		return validateCheckResponse(tx, event)
	}
	if _, current := currentEventRules[event.EventType]; !current || !event.IsV2() {
		return nil
	}
	var schema int
	if err := tx.QueryRow("SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name='run_runtime'").Scan(&schema); err != nil {
		return err
	}
	if schema == 0 {
		return nil
	}
	if err := validateCurrentBudgetFacts(tx, event); err != nil {
		return err
	}
	var result map[string]any
	for _, effect := range event.Effects {
		if effect.Table == "idempotency_results" {
			v, err := decodeCurrentJSON(effect.Values["result_json"].(string))
			if err != nil {
				return err
			}
			result = v.(map[string]any)
		}
	}
	if result == nil {
		return fmt.Errorf("current Event missing result")
	}
	expected, err := currentExpectedResult(tx, event)
	if err != nil {
		return err
	}
	expected["idempotency_key"] = event.Idempotency.Key
	expected["replayed"] = false
	if !currentValuesEqual(result, expected) {
		for key, want := range expected {
			if !currentValuesEqual(result[key], want) {
				return fmt.Errorf("current Event %s: stored result %s disagrees with recorded facts", event.EventType, currentResultDifference(key, result[key], want))
			}
		}
		return fmt.Errorf("current Event: stored result fields mismatch")
	}
	return nil
}

func currentResultDifference(prefix string, actual, expected any) string {
	if want, ok := expected.(map[string]any); ok {
		got, ok := actual.(map[string]any)
		if !ok {
			return prefix
		}
		for _, key := range sortedMapKeys(want) {
			if !currentValuesEqual(got[key], want[key]) {
				return currentResultDifference(prefix+"."+key, got[key], want[key])
			}
		}
	}
	return prefix
}
func currentExpectedResult(tx *sql.Tx, event Event) (map[string]any, error) {
	result := map[string]any{}
	var failure error
	put := func(key, kind, table string, pk map[string]any) {
		if failure != nil {
			return
		}
		var v map[string]any
		v, failure = currentProjectedValue(tx, kind, table, pk)
		result[key] = v
	}
	id := func(value any) map[string]any { return map[string]any{"id": value} }
	switch event.EventType {
	case "review.recorded":
		put("review", "Review", "review_records", id(event.Aggregate.ID))
		if failure == nil {
			var budgetResult map[string]any
			budgetResult, failure = currentReviewBudgetResult(tx, event)
			if failure == nil {
				result["budget"], result["stops"], result["attention"] = budgetResult["budget"], budgetResult["stops"], budgetResult["attention"]
			}
		}
	case "finding.addressed":
		put("finding", "Finding", "review_findings", id(event.Aggregate.ID))
	case "concern.recorded", "concern.disposed":
		put("concern", "Concern", "concerns", id(event.Aggregate.ID))
	case "deviation.recorded":
		put("deviation", "Deviation", "deviations", id(event.Aggregate.ID))
	case "decision.referenced":
		put("decision", "DecisionRef", "decision_refs", id(event.Aggregate.ID))
		if failure == nil {
			result["protected_decisions"], failure = currentProtectedDecisions(tx, result["decision"].(map[string]any)["scope"].(map[string]any)["issue_id"].(string))
		}
		if failure == nil {
			result["stops"], result["attention"], failure = currentDecisionResolutionValues(tx, event)
		}
	case "human-outcome.recorded":
		put("outcome", "HumanOutcome", "human_outcomes", id(event.Aggregate.ID))
		if failure == nil {
			result["protected_decisions"], failure = currentProtectedDecisions(tx, result["outcome"].(map[string]any)["scope"].(map[string]any)["issue_id"].(string))
		}
		if failure == nil {
			result["stops"], result["attention"], failure = currentDecisionResolutionValues(tx, event)
		}
	case "requirement.referenced":
		put("requirement", "RequirementRef", "requirement_refs", id(event.Aggregate.ID))
	case "requirement.linked":
		for _, effect := range event.Effects {
			if effect.Table == "requirement_issue_links" && effect.Op == "insert" {
				put("link", "RequirementLink", "requirement_issue_links", effect.PK)
			}
		}
	case "run.started", "run.closed", "run.observed", "run.checkpointed", "run.abandoned":
		put("run", "Run", "runs", id(event.Aggregate.ID))
		if event.EventType == "run.abandoned" && failure == nil {
			put("claim", "Claim", "claims", id(result["run"].(map[string]any)["claim_id"]))
		}
		if event.EventType == "run.observed" || event.EventType == "run.checkpointed" {
			kind, key := "Observation", "observation"
			if event.EventType == "run.checkpointed" {
				kind, key = "Checkpoint", "checkpoint"
			}
			put(key, kind, "run_observations", id(currentInserted(event, "run_observations")["id"]))
		}
	case "run.recovered":
		recovery := currentInserted(event, "run_recoveries")
		put("recovery", "Recovery", "run_recoveries", id(recovery["id"]))
		put("run", "Run", "runs", id(recovery["successor_run_id"]))
		put("claim", "Claim", "claims", id(recovery["successor_claim_id"]))
	case "claim.expired":
		put("claim", "Claim", "claims", id(event.Aggregate.ID))
	case "stop.raised", "stop.resolved":
		put("stop", "Stop", "stops", id(event.Aggregate.ID))
		if event.EventType == "stop.raised" {
			for _, e := range event.Effects {
				if e.Table == "attention_stop_links" {
					put("attention", "Attention", "attention_items", id(e.PK["attention_id"]))
				}
			}
		}
	case "attention.resolved":
		put("attention", "Attention", "attention_items", id(event.Aggregate.ID))
	case "gate.execution-started":
		put("attempt", "Attempt", "execution_attempts", id(event.Aggregate.ID))
		result["receipt"] = nil
	case "gate.receipt-recorded":
		put("receipt", "Receipt", "execution_receipts", id(event.Aggregate.ID))
	case "run.budget-observed", "run.execution-blocked":
		return currentBudgetEventResult(tx, event)
	default:
		if !strings.HasPrefix(event.EventType, "activation.") {
			return nil, fmt.Errorf("unknown current result")
		}
		put("activation", "Activation", "kit_activations", id(event.Aggregate.ID))
		if event.EventType == "activation.component-staged" {
			for _, e := range event.Effects {
				if e.Table == "activation_components" && e.Op == "insert" {
					row, err := currentReference(tx, e.Table, e.PK)
					if err != nil {
						return nil, err
					}
					result["component"] = map[string]any{"type": row["component_type"], "id": row["component_id"], "version": row["component_version"], "digest": row["content_digest"]}
				}
			}
		}
		if event.EventType == "activation.resource-observed" {
			for _, e := range event.Effects {
				if e.Table == "activation_resources" {
					put("resource", "ActivationResource", e.Table, e.PK)
				}
			}
		}
	}
	return result, failure
}
func currentPublicRow(tx *sql.Tx, table string, pk map[string]any) (map[string]any, error) {
	row, err := currentReference(tx, table, pk)
	if err != nil {
		return nil, err
	}
	result := map[string]any{}
	for column, value := range row {
		name := column
		if strings.HasSuffix(column, "_json") {
			name = strings.TrimSuffix(column, "_json")
			if value != nil {
				value, err = decodeCurrentJSON(value.(string))
				if err != nil {
					return nil, err
				}
			}
		}
		if column == "created_event_id" {
			name = "event_id"
		}
		if column == "terminal_status" {
			name = "outcome"
		}
		if currentEventTables[table].columns[column] == "b" {
			value = currentValuesEqual(value, 1) || value == true
		}
		result[name] = value
	}
	return result, nil
}
func currentProjectedValue(tx *sql.Tx, kind, table string, pk map[string]any) (map[string]any, error) {
	row, err := currentPublicRow(tx, table, pk)
	if err != nil {
		return nil, err
	}
	switch kind {
	case "Observation", "Checkpoint":
		return row["observation"].(map[string]any), nil
	case "Run":
		err = currentFillRun(tx, row)
	case "Activation":
		err = currentFillActivation(tx, row)
	case "Stop", "Attention":
		err = currentFillStopAttention(tx, kind, row)
	case "Receipt":
		err = currentFillReceipt(tx, row)
	case "Review":
		err = currentFillReview(tx, row)
	case "Finding":
		err = currentFillFinding(tx, row)
	case "Concern":
		err = currentFillConcern(row)
	case "Deviation":
		err = currentFillDeviation(row)
	case "DecisionRef":
		err = currentFillDecisionRef(row)
	case "HumanOutcome":
		err = currentFillHumanOutcome(row)
	}
	if err != nil {
		return nil, err
	}
	result := map[string]any{}
	for _, field := range strings.Fields(currentValueShapes[kind]) {
		name, _, _ := strings.Cut(field, ":")
		v, ok := row[name]
		if !ok {
			return nil, fmt.Errorf("missing projected %s.%s", kind, name)
		}
		result[name] = v
	}
	return result, nil
}

func currentFillConcern(row map[string]any) error {
	row["subject"] = map[string]any{"type": row["subject_type"], "id": row["subject_id"]}
	row["disposition"] = nil
	if row["disposition_action"] != nil {
		row["disposition"] = map[string]any{
			"action": row["disposition_action"], "issue_id": row["disposition_issue_id"],
			"supersedes_concern_id": row["supersedes_concern_id"], "supersedes_finding_id": row["supersedes_finding_id"],
			"reason": row["disposition_reason"], "actor_id": row["disposed_actor_id"],
			"disposed_at": row["disposed_at"], "event_id": row["disposed_event_id"],
		}
	}
	return nil
}

func currentFillDeviation(row map[string]any) error {
	row["subject"] = map[string]any{"type": row["subject_type"], "id": row["subject_id"]}
	row["authority"] = map[string]any{"actor_id": row["authority_actor_id"], "reference": row["authority_reference"]}
	return nil
}

func currentFillDecisionRef(row map[string]any) error {
	row["authority"] = map[string]any{"actor_id": row["authority_actor_id"], "required_capabilities": row["authority_capabilities"]}
	row["scope"] = map[string]any{"issue_id": row["scope_issue_id"]}
	return nil
}

func currentFillHumanOutcome(row map[string]any) error {
	row["subject"] = map[string]any{"type": row["subject_type"], "id": row["subject_id"]}
	row["scope"] = map[string]any{"issue_id": row["scope_issue_id"]}
	row["provenance"] = map[string]any{"activation_id": row["activation_id"], "activation_version": row["activation_version"], "profile_ref": row["profile_ref"]}
	return nil
}

func currentProtectedDecisions(tx *sql.Tx, issueID string) ([]any, error) {
	rows, err := tx.Query(`
		SELECT p.ordinal,p.decision_class,p.required_capabilities_json,p.affected_actions_json
		FROM envelope_protected_decisions p
		JOIN work_envelopes e ON e.id=p.envelope_id AND e.current_revision=p.revision
		JOIN envelope_issue_selection s ON s.envelope_id=p.envelope_id AND s.revision=p.revision
		WHERE s.issue_id=? AND e.status IN ('active','draft')
		  AND p.envelope_id=(
		    SELECT selected.id
		    FROM work_envelopes selected
		    JOIN envelope_issue_selection selected_scope ON selected_scope.envelope_id=selected.id AND selected_scope.revision=selected.current_revision
		    WHERE selected_scope.issue_id=? AND selected.status IN ('active','draft')
		    ORDER BY CASE selected.status WHEN 'active' THEN 0 ELSE 1 END,selected.id
		    LIMIT 1
		  )
		ORDER BY p.ordinal`, issueID, issueID)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()
	result := []any{}
	for rows.Next() {
		var ordinal int
		var class, capabilitiesJSON, actionsJSON string
		if err := rows.Scan(&ordinal, &class, &capabilitiesJSON, &actionsJSON); err != nil {
			return nil, err
		}
		capabilities, err := decodeCurrentJSON(capabilitiesJSON)
		if err != nil {
			return nil, err
		}
		actions, err := decodeCurrentJSON(actionsJSON)
		if err != nil {
			return nil, err
		}
		decisionRefID, err := currentOptionalID(tx, `SELECT id FROM decision_refs WHERE class=? AND scope_issue_id=? AND outcome='approved' AND authority_capabilities_json=? ORDER BY created_event_id DESC LIMIT 1`, class, issueID, capabilitiesJSON)
		if err != nil {
			return nil, err
		}
		outcomeID, err := currentOptionalID(tx, `SELECT id FROM human_outcomes WHERE class=? AND scope_issue_id=? AND result='approved' ORDER BY created_event_id DESC LIMIT 1`, class, issueID)
		if err != nil {
			return nil, err
		}
		status := "unsatisfied"
		if decisionRefID != nil && outcomeID != nil {
			status = "satisfied"
		}
		result = append(result, map[string]any{"ordinal": ordinal, "decision_class": class, "required_capabilities": capabilities, "affected_actions": actions, "status": status, "decision_ref_id": decisionRefID, "human_outcome_id": outcomeID})
	}
	return result, rows.Err()
}

func currentOptionalID(tx *sql.Tx, query string, args ...any) (any, error) {
	var id string
	if err := tx.QueryRow(query, args...).Scan(&id); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return id, nil
}

func currentDecisionResolutionValues(tx *sql.Tx, event Event) ([]any, []any, error) {
	stops := []any{}
	attention := []any{}
	seen := map[any]bool{}
	for _, effect := range event.Effects {
		if effect.Table == "stop_resolutions" && effect.Op == "insert" {
			value, err := currentProjectedValue(tx, "Stop", "stops", map[string]any{"id": effect.PK["stop_id"]})
			if err != nil {
				return nil, nil, err
			}
			stops = append(stops, value)
		}
		if effect.Table == "attention_items" && effect.Op == "update" && !seen[effect.PK["id"]] {
			value, err := currentProjectedValue(tx, "Attention", "attention_items", effect.PK)
			if err != nil {
				return nil, nil, err
			}
			attention = append(attention, value)
			seen[effect.PK["id"]] = true
		}
	}
	return stops, attention, nil
}

func currentQueryValues(tx *sql.Tx, query string, args ...any) ([]any, error) {
	rows, err := tx.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()
	result := []any{}
	for rows.Next() {
		var value any
		if err := rows.Scan(&value); err != nil {
			return nil, err
		}
		result = append(result, value)
	}
	return result, rows.Err()
}
func currentCollection(tx *sql.Tx, kind, table, pkColumn, query string, args ...any) ([]any, error) {
	ids, err := currentQueryValues(tx, query, args...)
	if err != nil {
		return nil, err
	}
	result := []any{}
	for _, id := range ids {
		value, err := currentProjectedValue(tx, kind, table, map[string]any{pkColumn: id})
		if err != nil {
			return nil, err
		}
		result = append(result, value)
	}
	return result, nil
}
