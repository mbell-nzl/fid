package maintenance

import (
	"database/sql"
	"fmt"
)

func checkKindHierarchy(db *sql.DB, _ string) ([]Finding, error) {
	rows, err := db.Query(`
		SELECT child.id, child.kind, child.parent_id, parent.kind
		FROM issues AS child
		LEFT JOIN issues AS parent ON parent.id = child.parent_id
		ORDER BY child.id
	`)
	if err != nil {
		return nil, fmt.Errorf("check kind hierarchy: query issues: %w", err)
	}

	findings := make([]Finding, 0)
	for rows.Next() {
		var id, kind string
		var parentID, parentKind sql.NullString
		if err := rows.Scan(&id, &kind, &parentID, &parentKind); err != nil {
			_ = rows.Close()
			return nil, fmt.Errorf("check kind hierarchy: scan issue: %w", err)
		}
		if finding, invalid := hierarchyFinding(id, kind, parentID, parentKind); invalid {
			findings = append(findings, finding)
		}
	}
	if err := rows.Err(); err != nil {
		_ = rows.Close()
		return nil, fmt.Errorf("check kind hierarchy: iterate issues: %w", err)
	}
	if err := rows.Close(); err != nil {
		return nil, fmt.Errorf("check kind hierarchy: close issues: %w", err)
	}
	return findings, nil
}

func hierarchyFinding(id, kind string, parentID, parentKind sql.NullString) (Finding, bool) {
	const check = "kind_hierarchy"
	if kind == "epic" && parentID.Valid {
		return Finding{Check: check, EntityID: id, Message: fmt.Sprintf("epic %s has parent %s; epics must not have parents", id, parentID.String)}, true
	}
	if kind == "task" && parentID.Valid && (!parentKind.Valid || parentKind.String != "epic") {
		return Finding{Check: check, EntityID: id, Message: fmt.Sprintf("task %s parent %s must be an epic", id, parentID.String)}, true
	}
	if kind == "subtask" && (!parentID.Valid || !parentKind.Valid || parentKind.String != "task") {
		parent := "none"
		if parentID.Valid {
			parent = parentID.String
		}
		return Finding{Check: check, EntityID: id, Message: fmt.Sprintf("subtask %s parent %s must be a task", id, parent)}, true
	}
	return Finding{}, false
}

func checkDependentRowKinds(db *sql.DB, _ string) ([]Finding, error) {
	checks := []dependentRowCheck{
		{table: "acceptance_criteria", idColumn: "id", issueColumn: "issue_id", allowedKinds: []string{"task", "subtask"}},
		{table: "epic_metadata", idColumn: "issue_id", issueColumn: "issue_id", allowedKinds: []string{"epic"}},
		{table: "triage_events", idColumn: "id", issueColumn: "epic_id", allowedKinds: []string{"epic"}},
	}
	findings := make([]Finding, 0)
	for _, check := range checks {
		rows, err := db.Query(check.query())
		if err != nil {
			return nil, fmt.Errorf("check dependent row kinds: query %s: %w", check.table, err)
		}
		for rows.Next() {
			var id string
			var issueID, kind sql.NullString
			if err := rows.Scan(&id, &issueID, &kind); err != nil {
				_ = rows.Close()
				return nil, fmt.Errorf("check dependent row kinds: scan %s: %w", check.table, err)
			}
			if !kindAllowed(kind, check.allowedKinds) {
				findings = append(findings, dependentRowFinding(check, id, issueID, kind))
			}
		}
		if err := rows.Err(); err != nil {
			_ = rows.Close()
			return nil, fmt.Errorf("check dependent row kinds: iterate %s: %w", check.table, err)
		}
		if err := rows.Close(); err != nil {
			return nil, fmt.Errorf("check dependent row kinds: close %s: %w", check.table, err)
		}
	}
	return findings, nil
}

type dependentRowCheck struct {
	table, idColumn, issueColumn string
	allowedKinds                 []string
	optional                     bool
}

func (check dependentRowCheck) query() string {
	where := ""
	if check.optional {
		where = " WHERE dependent." + check.issueColumn + " IS NOT NULL"
	}
	return "SELECT dependent." + check.idColumn + ", dependent." + check.issueColumn + ", issue.kind FROM " + check.table + " AS dependent LEFT JOIN issues AS issue ON issue.id = dependent." + check.issueColumn + where + " ORDER BY dependent." + check.idColumn
}

func kindAllowed(kind sql.NullString, allowed []string) bool {
	for _, candidate := range allowed {
		if kind.Valid && kind.String == candidate {
			return true
		}
	}
	return false
}

func dependentRowFinding(check dependentRowCheck, id string, issueID, kind sql.NullString) Finding {
	issue := "missing"
	if issueID.Valid {
		issue = issueID.String
	}
	return Finding{
		Check:    "dependent_row_kind",
		EntityID: id,
		Message:  fmt.Sprintf("%s %s references %s, which must be %s", check.table, id, issue, joinKinds(check.allowedKinds)),
	}
}

func joinKinds(kinds []string) string {
	if len(kinds) == 1 {
		return "an " + kinds[0]
	}
	return "a " + kinds[0] + " or " + kinds[1]
}

func checkSequenceScheduling(db *sql.DB, _ string) ([]Finding, error) {
	rows, err := db.Query(`
		SELECT task.id, task.parent_id, task.sequence, parent.kind, metadata.scheduling
		FROM issues AS task
		LEFT JOIN issues AS parent ON parent.id = task.parent_id
		LEFT JOIN epic_metadata AS metadata ON metadata.issue_id = task.parent_id
		WHERE task.kind = 'task'
		ORDER BY task.id
	`)
	if err != nil {
		return nil, fmt.Errorf("check sequence scheduling: query tasks: %w", err)
	}

	findings := make([]Finding, 0)
	for rows.Next() {
		var id string
		var parentID, parentKind, scheduling sql.NullString
		var sequence sql.NullInt64
		if err := rows.Scan(&id, &parentID, &sequence, &parentKind, &scheduling); err != nil {
			_ = rows.Close()
			return nil, fmt.Errorf("check sequence scheduling: scan task: %w", err)
		}
		if finding, invalid := sequenceFinding(id, parentID, sequence, parentKind, scheduling); invalid {
			findings = append(findings, finding)
		}
	}
	if err := rows.Err(); err != nil {
		_ = rows.Close()
		return nil, fmt.Errorf("check sequence scheduling: iterate tasks: %w", err)
	}
	if err := rows.Close(); err != nil {
		return nil, fmt.Errorf("check sequence scheduling: close tasks: %w", err)
	}
	return findings, nil
}

func sequenceFinding(id string, parentID sql.NullString, sequence sql.NullInt64, parentKind, scheduling sql.NullString) (Finding, bool) {
	const check = "sequence_scheduling"
	if !parentID.Valid && sequence.Valid {
		return Finding{Check: check, EntityID: id, Message: fmt.Sprintf("unattached task %s must not have sequence", id)}, true
	}
	if !parentID.Valid || !parentKind.Valid || parentKind.String != "epic" {
		return Finding{}, false
	}
	if !scheduling.Valid {
		return Finding{Check: check, EntityID: id, Message: fmt.Sprintf("task %s parent epic %s must have epic_metadata", id, parentID.String)}, true
	}
	if scheduling.String == "sequential" && !sequence.Valid {
		return Finding{Check: check, EntityID: id, Message: fmt.Sprintf("task %s under sequential epic %s requires sequence", id, parentID.String)}, true
	}
	if scheduling.String == "prioritised" && sequence.Valid {
		return Finding{Check: check, EntityID: id, Message: fmt.Sprintf("task %s under prioritised epic %s must not have sequence", id, parentID.String)}, true
	}
	return Finding{}, false
}

func checkLinkReferences(db *sql.DB, _ string) ([]Finding, error) {
	rows, err := db.Query(`SELECT from_type, from_id, relation, to_type, to_id FROM links ORDER BY from_type, from_id, relation, to_type, to_id`)
	if err != nil {
		return nil, fmt.Errorf("check link references: query links: %w", err)
	}

	links := make([]linkRow, 0)
	for rows.Next() {
		var link linkRow
		if err := rows.Scan(&link.fromType, &link.fromID, &link.relation, &link.toType, &link.toID); err != nil {
			_ = rows.Close()
			return nil, fmt.Errorf("check link references: scan link: %w", err)
		}
		links = append(links, link)
	}
	if err := rows.Err(); err != nil {
		_ = rows.Close()
		return nil, fmt.Errorf("check link references: iterate links: %w", err)
	}
	if err := rows.Close(); err != nil {
		return nil, fmt.Errorf("check link references: close links: %w", err)
	}

	findings := make([]Finding, 0)
	for _, link := range links {
		rowFindings, err := linkFindings(db, link.fromType, link.fromID, link.relation, link.toType, link.toID)
		if err != nil {
			return nil, fmt.Errorf("check link references: query link %s: %w", link.entityID(), err)
		}
		findings = append(findings, rowFindings...)
	}
	return findings, nil
}

type linkRow struct{ fromType, fromID, relation, toType, toID string }

func (link linkRow) entityID() string {
	return linkEntityID(link.fromType, link.fromID, link.relation, link.toType, link.toID)
}

func linkFindings(db *sql.DB, fromType, fromID, relation, toType, toID string) ([]Finding, error) {
	entityID := linkEntityID(fromType, fromID, relation, toType, toID)
	if !validLinkType(fromType) || !validLinkType(toType) {
		return []Finding{{Check: "link_scope", EntityID: entityID, Message: fmt.Sprintf("link %s must use issue or requirement endpoints in v0.1", entityID)}}, nil
	}
	if !linkRelationAllowed(fromType, relation, toType) {
		return []Finding{{Check: "link_relation_compatibility", EntityID: entityID, Message: fmt.Sprintf("link %s relation %s is not allowed for %s to %s", entityID, relation, fromType, toType)}}, nil
	}
	findings := make([]Finding, 0, 2)
	if exists, err := referencedEntityExists(db, fromType, fromID); err != nil {
		return nil, err
	} else if !exists {
		findings = append(findings, Finding{Check: "link_refs", EntityID: entityID, Message: fmt.Sprintf("link %s from endpoint %s:%s is missing", entityID, fromType, fromID)})
	}
	if exists, err := referencedEntityExists(db, toType, toID); err != nil {
		return nil, err
	} else if !exists {
		findings = append(findings, Finding{Check: "link_refs", EntityID: entityID, Message: fmt.Sprintf("link %s to endpoint %s:%s is missing", entityID, toType, toID)})
	}
	return findings, nil
}

func validLinkType(refType string) bool { return refType == "issue" }

func linkRelationAllowed(fromType, relation, toType string) bool {
	if relation == "relates_to" {
		return true
	}
	return fromType == "issue" && toType == "issue" && map[string]bool{
		"blocks": true, "blocked_by": true, "duplicates": true, "duplicated_by": true,
		"supersedes": true, "superseded_by": true, "caused_by": true, "causes": true,
		"discovered_from": true, "discovered": true,
	}[relation]
}

func linkEntityID(fromType, fromID, relation, toType, toID string) string {
	return fromType + ":" + fromID + ":" + relation + ":" + toType + ":" + toID
}

func referencedEntityExists(db *sql.DB, refType, refID string) (bool, error) {
	table, ok := referenceTable(refType)
	if !ok {
		return false, nil
	}
	var exists bool
	column := "id"
	if refType == "requirement" {
		column = "requirement_id"
	}
	if err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM "+table+" WHERE "+column+" = ?)", refID).Scan(&exists); err != nil {
		return false, err
	}
	return exists, nil
}

func referenceTable(refType string) (string, bool) {
	switch refType {
	case "issue":
		return "issues", true
	case "requirement":
		return "requirement_refs", true
	case "decision":
		return "decision_refs", true
	default:
		return "", false
	}
}
