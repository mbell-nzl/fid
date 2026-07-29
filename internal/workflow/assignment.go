package workflow

import (
	"database/sql"
	"fmt"
)

// AssignmentRequest is a parsed request to bind one workflow to one scope.
type AssignmentRequest struct {
	Name    string
	Kind    string
	Project bool
}

// ValidateAssignmentRequest validates assignment scope and workflow applicability.
func ValidateAssignmentRequest(conn *sql.DB, request AssignmentRequest) error {
	if request.Kind != "" && request.Project {
		return fmt.Errorf("workflow assignment specifies both kind and project scopes")
	}
	if request.Kind == "" && !request.Project {
		return fmt.Errorf("workflow assignment requires a kind or project scope")
	}
	appliesTo, err := workflowAppliesTo(conn, request.Name)
	if err != nil {
		return err
	}
	if request.Project {
		return validateProjectAssignment(request.Name, appliesTo)
	}
	return validateKindAssignment(request.Name, request.Kind, appliesTo)
}

func workflowAppliesTo(conn *sql.DB, name string) (sql.NullString, error) {
	if conn == nil {
		return sql.NullString{}, fmt.Errorf("check workflow %q: database connection is required", name)
	}
	var appliesTo sql.NullString
	err := conn.QueryRow(`SELECT applies_to FROM workflows WHERE id = ? LIMIT 1`, name).Scan(&appliesTo)
	if err == sql.ErrNoRows {
		return sql.NullString{}, fmt.Errorf("workflow %q does not exist", name)
	}
	if err != nil {
		return sql.NullString{}, fmt.Errorf("load workflow %q for assignment: %w", name, err)
	}
	return appliesTo, nil
}

func validateProjectAssignment(name string, appliesTo sql.NullString) error {
	if !appliesTo.Valid {
		return nil
	}
	return fmt.Errorf("workflow %q applies_to %q cannot be assigned to project", name, appliesTo.String)
}

func validateKindAssignment(name, kind string, appliesTo sql.NullString) error {
	if !appliesTo.Valid || appliesTo.String == kind {
		return nil
	}
	return fmt.Errorf("workflow %q applies_to %q is incompatible with kind %q", name, appliesTo.String, kind)
}
