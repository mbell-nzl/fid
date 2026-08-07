package maintenance

import (
	"database/sql"
	"fmt"
)

func checkWorkflowAssignments(db *sql.DB, _ string) ([]Finding, error) {
	rows, err := db.Query(`
		SELECT assignment.scope, assignment.kind, assignment.workflow_id, workflow.applies_to
		FROM workflow_assignments AS assignment
		JOIN workflows AS workflow ON workflow.id = assignment.workflow_id
		ORDER BY assignment.scope, assignment.kind, assignment.workflow_id
	`)
	if err != nil {
		return nil, fmt.Errorf("check workflow assignments: query assignments: %w", err)
	}

	findings := make([]Finding, 0)
	for rows.Next() {
		var scope, workflowID string
		var kind, appliesTo sql.NullString
		if err := rows.Scan(&scope, &kind, &workflowID, &appliesTo); err != nil {
			_ = rows.Close()
			return nil, fmt.Errorf("check workflow assignments: scan assignment: %w", err)
		}
		if finding, incompatible := workflowAssignmentFinding(scope, kind, workflowID, appliesTo); incompatible {
			findings = append(findings, finding)
		}
	}
	if err := rows.Err(); err != nil {
		_ = rows.Close()
		return nil, fmt.Errorf("check workflow assignments: iterate assignments: %w", err)
	}
	if err := rows.Close(); err != nil {
		return nil, fmt.Errorf("check workflow assignments: close assignments: %w", err)
	}
	return findings, nil
}

func workflowAssignmentFinding(scope string, kind sql.NullString, workflowID string, appliesTo sql.NullString) (Finding, bool) {
	if scope == "kind" && (!appliesTo.Valid || appliesTo.String == kind.String) {
		return Finding{}, false
	}
	if scope == "project" && !appliesTo.Valid {
		return Finding{}, false
	}

	entityID := scope
	if kind.Valid {
		entityID += ":" + kind.String
	}
	appliesToValue := "any"
	if appliesTo.Valid {
		appliesToValue = appliesTo.String
	}
	return Finding{
		Check:    "workflow_assignment",
		EntityID: entityID,
		Message:  fmt.Sprintf("%s assignment %s uses workflow %s with incompatible applies_to %s", scope, entityID, workflowID, appliesToValue),
	}, true
}
