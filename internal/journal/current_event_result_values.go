package journal

import (
	"database/sql"
	"fmt"
)

func currentFillRun(tx *sql.Tx, row map[string]any) error {
	runtime, err := currentReference(tx, "run_runtime", map[string]any{"run_id": row["id"]})
	if err != nil {
		return err
	}
	row["version"], row["obligation_id"] = runtime["version"], runtime["obligation_id"]
	row["executor"] = map[string]any{"id": runtime["executor_id"], "source": "activated-harness"}
	row["predecessor_run_id"], row["successor_run_id"] = nil, nil
	for key, query := range map[string]string{"predecessor_run_id": "SELECT prior_run_id FROM run_recoveries WHERE successor_run_id=?", "successor_run_id": "SELECT successor_run_id FROM run_recoveries WHERE prior_run_id=?"} {
		values, err := currentQueryValues(tx, query, row["id"])
		if err != nil {
			return err
		}
		if len(values) > 1 {
			return fmt.Errorf("ambiguous Run link")
		}
		if len(values) == 1 {
			row[key] = values[0]
		}
	}
	names, err := currentQueryValues(tx, "SELECT name FROM run_budget_bindings WHERE run_id=? ORDER BY name", row["id"])
	if err != nil {
		return err
	}
	budgets := []any{}
	for _, name := range names {
		budget, err := currentProjectedBudget(tx, row["id"], name)
		if err != nil {
			return err
		}
		budgets = append(budgets, budget)
	}
	row["budgets"] = budgets
	return nil
}
func currentFillActivation(tx *sql.Tx, row map[string]any) error {
	digests, err := currentQueryValues(tx, "SELECT content_digest FROM activation_components WHERE activation_id=? ORDER BY content_digest", row["id"])
	if err != nil {
		return err
	}
	row["component_digests"] = digests
	resources, err := currentActivationResources(tx, row["id"])
	if err != nil {
		return err
	}
	row["resources"] = resources
	row["recovery"] = nil
	var count int
	if err := tx.QueryRow("SELECT COUNT(*) FROM activation_recoveries WHERE activation_id=?", row["id"]).Scan(&count); err != nil {
		return err
	}
	if count != 0 {
		v, err := currentProjectedValue(tx, "ActivationRecovery", "activation_recoveries", map[string]any{"activation_id": row["id"]})
		if err != nil {
			return err
		}
		row["recovery"] = v
	}
	return nil
}
func currentFillStopAttention(tx *sql.Tx, kind string, row map[string]any) error {
	row["subject"] = map[string]any{"type": row["subject_type"], "id": row["subject_id"]}
	if kind == "Attention" {
		ids, err := currentQueryValues(tx, "SELECT stop_id FROM attention_stop_links WHERE attention_id=? ORDER BY stop_id", row["id"])
		row["stop_ids"] = ids
		return err
	}
	row["resolution"] = nil
	var count int
	if err := tx.QueryRow("SELECT COUNT(*) FROM stop_resolutions WHERE stop_id=?", row["id"]).Scan(&count); err != nil {
		return err
	}
	if count != 0 {
		v, err := currentProjectedValue(tx, "StopResolution", "stop_resolutions", map[string]any{"stop_id": row["id"]})
		if err != nil {
			return err
		}
		row["resolution"] = v
	}
	return nil
}
func currentFillReceipt(tx *sql.Tx, row map[string]any) error {
	attempt, err := currentReference(tx, "execution_attempts", map[string]any{"id": row["attempt_id"]})
	if err != nil {
		return err
	}
	for _, field := range []string{"registration_id", "registration_version", "obligation_id", "run_id", "actor_id"} {
		row[field] = attempt[field]
	}
	artifacts, err := currentCollection(tx, "Artifact", "referenced_artifacts", "id", "SELECT artifact_id FROM receipt_artifacts WHERE receipt_id=? ORDER BY artifact_id", row["id"])
	row["artifacts"] = artifacts
	return err
}

func currentFillReview(tx *sql.Tx, row map[string]any) error {
	findings, err := currentCollection(tx, "Finding", "review_findings", "id", "SELECT id FROM review_findings WHERE review_id=? ORDER BY id", row["id"])
	if err != nil {
		return err
	}
	row["findings"] = findings
	resolved, err := currentQueryValues(tx, "SELECT finding_id FROM finding_dispositions WHERE qualifying_review_id=? AND action='resolved' ORDER BY finding_id", row["id"])
	if err != nil {
		return err
	}
	row["resolved_finding_ids"] = resolved
	return nil
}

func currentFillFinding(tx *sql.Tx, row map[string]any) error {
	dispositions, err := currentCollection(tx, "FindingDisposition", "finding_dispositions", "id", "SELECT id FROM finding_dispositions WHERE finding_id=? ORDER BY disposed_at, id", row["id"])
	if err != nil {
		return err
	}
	row["subject"] = map[string]any{"type": row["subject_type"], "id": row["subject_id"]}
	row["dispositions"] = dispositions
	row["state"] = "open"
	if len(dispositions) != 0 {
		action, ok := dispositions[len(dispositions)-1].(map[string]any)["action"].(string)
		if !ok {
			return fmt.Errorf("invalid projected Finding disposition")
		}
		row["state"] = action
	}
	return nil
}

// Read the complete resource snapshot in one ordered query. Result validation
// still compares every field, without preparing one query per resource.
func currentActivationResources(tx *sql.Tx, activationID any) ([]any, error) {
	rows, err := tx.Query(`SELECT component, version, old_identity, new_identity, observed_identity, state, evidence_ref, error_code FROM activation_resources WHERE activation_id=? ORDER BY component`, activationID)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()
	resources := []any{}
	for rows.Next() {
		var component, version, oldIdentity, newIdentity, observed, state, evidence, code any
		if err := rows.Scan(&component, &version, &oldIdentity, &newIdentity, &observed, &state, &evidence, &code); err != nil {
			return nil, err
		}
		resources = append(resources, map[string]any{"component": component, "version": version, "old_identity": oldIdentity, "new_identity": newIdentity, "observed_identity": observed, "state": state, "evidence_ref": evidence, "error_code": code})
	}
	return resources, rows.Err()
}
