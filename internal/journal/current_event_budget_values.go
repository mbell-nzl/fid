package journal

import (
	"database/sql"
	"fmt"
	"time"
)

func validateCurrentBudgetSampleChain(tx *sql.Tx, effect Effect, policyRunID any) error {
	policy, err := currentPublicRow(tx, "run_budget_bindings", map[string]any{"run_id": policyRunID, "name": effect.Values["budget_name"]})
	if err != nil {
		return err
	}
	budget := policy["budget"].(map[string]any)
	owned := currentEnum(budget["unit"], "attempts", "failures", "milliseconds", "tool-failures", "interventions", "review-cycles")
	if (effect.Values["source"] == "fid") != owned {
		return fmt.Errorf("current budget sample: source does not own unit")
	}
	var latest sql.NullString
	err = tx.QueryRow(`SELECT id FROM budget_samples WHERE envelope_id=? AND envelope_revision=? AND budget_name=? ORDER BY observed_at DESC,id DESC LIMIT 1`, effect.Values["envelope_id"], effect.Values["envelope_revision"], effect.Values["budget_name"]).Scan(&latest)
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	var expected any
	if latest.Valid {
		expected = latest.String
	}
	if !currentValuesEqual(effect.Values["previous_sample_id"], expected) {
		return fmt.Errorf("current budget sample: predecessor is not latest recorded sample")
	}
	return nil
}

func validateCurrentBudgetFacts(tx *sql.Tx, event Event) error {
	runID := any(event.Aggregate.ID)
	switch event.EventType {
	case "run.budget-observed", "run.execution-blocked":
	case "run.recovered":
		runID = currentInserted(event, "runs")["id"]
	case "gate.execution-started":
		runID = currentInserted(event, "execution_attempts")["run_id"]
	case "gate.receipt-recorded":
		attempt, err := currentReference(tx, "execution_attempts", map[string]any{"id": currentInserted(event, "execution_receipts")["attempt_id"]})
		if err != nil {
			return err
		}
		runID = attempt["run_id"]
	default:
		return nil
	}
	names, err := currentQueryValues(tx, "SELECT name FROM run_budget_bindings WHERE run_id=? ORDER BY name", runID)
	if err != nil {
		return err
	}
	required := map[string]map[string]any{}
	for _, name := range names {
		before, err := currentProjectedBudget(tx, runID, name, event.EventID)
		if err != nil {
			return err
		}
		after, err := currentProjectedBudget(tx, runID, name)
		if err != nil {
			return err
		}
		if err := validateCurrentAutomaticCharge(event, name, before["unit"]); err != nil {
			return err
		}
		changed := event.EventType == "run.execution-blocked"
		for _, effect := range event.Effects {
			if effect.Table == "budget_samples" && effect.Values["budget_name"] == name {
				changed = true
			}
		}
		if !changed || (event.EventType != "run.execution-blocked" && before["state"] == "exhausted") {
			continue
		}
		code := ""
		if after["state"] == "unknown" {
			code = "budget-data-unavailable"
		}
		if after["state"] == "exhausted" {
			code = "budget-exhausted"
		}
		if code != "" {
			required[name.(string)+"/"+code] = after
		}
	}
	return validateCurrentBudgetStops(event, required, runID)
}

func validateCurrentAutomaticCharge(event Event, name, unit any) error {
	if event.EventType != "gate.execution-started" && event.EventType != "gate.receipt-recorded" && event.EventType != "run.recovered" {
		return nil
	}
	var charge *int64
	var sourceID, observedAt any
	switch event.EventType {
	case "run.recovered":
		recovery := currentInserted(event, "run_recoveries")
		authority, err := decodeCurrentJSON(recovery["authority_json"].(string))
		if err != nil {
			return err
		}
		sourceID, observedAt = recovery["id"], currentInserted(event, "runs")["started_at"]
		if authority.(map[string]any)["mode"] == "force" && unit == "interventions" {
			v := int64(1)
			charge = &v
		}
	case "gate.execution-started":
		attempt := currentInserted(event, "execution_attempts")
		sourceID, observedAt = attempt["id"], attempt["started_at"]
		if unit == "attempts" {
			v := int64(1)
			charge = &v
		}
	default:
		receipt := currentInserted(event, "execution_receipts")
		sourceID, observedAt = receipt["id"], receipt["ended_at"]
		switch unit {
		case "failures":
			if receipt["terminal_status"] != "succeeded" {
				v := int64(1)
				charge = &v
			}
		case "tool-failures":
			if receipt["terminal_status"] == "spawn-failed" {
				v := int64(1)
				charge = &v
			}
		case "milliseconds":
			start, a := time.Parse(time.RFC3339Nano, receipt["started_at"].(string))
			end, b := time.Parse(time.RFC3339Nano, receipt["ended_at"].(string))
			if a != nil || b != nil || end.Before(start) {
				return fmt.Errorf("current charge: invalid Receipt timing")
			}
			v := end.Sub(start).Milliseconds()
			charge = &v
		}
	}
	count := 0
	for _, effect := range event.Effects {
		if effect.Table != "budget_samples" || effect.Values["budget_name"] != name {
			continue
		}
		count++
		if event.EventType == "run.recovered" && (effect.Values["actor_id"] != event.ActorV2.ActorID || effect.Values["measurement_ref_json"] != nil) {
			return fmt.Errorf("current intervention: invalid Actor or measurement reference")
		}
		if charge == nil || !currentValuesEqual(effect.Values["value"], *charge) || effect.Values["source"] != "fid" || effect.Values["sample_kind"] != "delta" || effect.Values["source_sample_id"] != sourceID || effect.Values["observed_at"] != observedAt {
			return fmt.Errorf("current charge: sample differs from captured execution facts")
		}
	}
	if (charge != nil && count != 1) || (charge == nil && count != 0) {
		return fmt.Errorf("current charge: incomplete automatic samples")
	}
	return nil
}

func validateCurrentBudgetStops(event Event, required map[string]map[string]any, runID any) error {
	for _, effect := range event.Effects {
		if effect.Table != "stops" {
			continue
		}
		if effect.Values["subject_type"] != "run" || effect.Values["subject_id"] != runID {
			return fmt.Errorf("current budget boundary: Stop Run binding mismatch")
		}
		value, err := decodeCurrentJSON(effect.Values["details_json"].(string))
		if err != nil {
			return err
		}
		details := value.(map[string]any)
		name, ok := details["budget_name"].(string)
		if !ok {
			return fmt.Errorf("current budget boundary: unrelated Stop")
		}
		key := name + "/" + effect.Values["code"].(string)
		budget, found := required[key]
		if !found || !currentMatching(details, budget, "unit", "maximum", "consumed", "remaining") {
			return fmt.Errorf("current budget boundary: Stop differs from derived boundary")
		}
		delete(required, key)
	}
	if len(required) != 0 {
		return fmt.Errorf("current budget boundary: missing required Stop")
	}
	return nil
}
