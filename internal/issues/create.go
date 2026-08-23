package issues

import (
	"database/sql"
	"errors"
	"fmt"
)

func Create(options CreateOptions) (CreateResult, error) {
	options.IdempotencyKey = idempotencyKeyForOptions(options.IdempotencyKey, options.ActorRole)
	request, err := normalizeCreateRequest(options)
	if err != nil {
		return CreateResult{}, err
	}
	return commitCreate(options, request)
}

func parentIDValue(parentID *string) any {
	if parentID == nil {
		return nil
	}
	return *parentID
}

func issueIDExists(conn *sql.DB, issueID string) (bool, error) {
	var existing string
	err := conn.QueryRow(`SELECT id FROM issues WHERE id = ? LIMIT 1`, issueID).Scan(&existing)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	return false, fmt.Errorf("query issue id %q: %w", issueID, err)
}

func supportedIssueKind(kind string) bool {
	switch kind {
	case "epic", "task", "subtask":
		return true
	default:
		return false
	}
}

func defaultIssueOwnerRole(kind string) string {
	if kind == "epic" {
		return "tl"
	}
	return "developer"
}

func validateCreateOptions(kind string, parentID *string, priority, sequence *int, scheduling string) error {
	switch kind {
	case "epic":
		if scheduling == "" {
			return fmt.Errorf("create issue: epic scheduling is required")
		}
	case "task":
		if scheduling != "" {
			return fmt.Errorf("create issue: scheduling is only valid for epics")
		}
	case "subtask":
		if scheduling != "" {
			return fmt.Errorf("create issue: scheduling is only valid for epics")
		}
	}

	if priority != nil && (*priority < 0 || *priority > 5) {
		return fmt.Errorf("create issue: priority must be between 0 and 5")
	}
	if sequence != nil && *sequence <= 0 {
		return fmt.Errorf("create issue: sequence must be positive")
	}
	if scheduling != "" && scheduling != "sequential" && scheduling != "prioritised" {
		return fmt.Errorf("create issue: scheduling must be sequential or prioritised")
	}
	return nil
}

func validateParentConstraint(conn *sql.DB, kind string, parentID *string) error {
	if parentID == nil {
		if kind == "subtask" {
			return hierarchyInvalidError(kind, nil, nil)
		}
		return nil
	}

	parentKind, found, err := issueKindForID(conn, *parentID)
	if err != nil {
		return err
	}
	if !found {
		return &NotFoundError{IssueID: *parentID}
	}
	switch kind {
	case "epic":
		return hierarchyInvalidError(kind, parentID, &parentKind)
	case "task":
		if parentKind != "epic" {
			return hierarchyInvalidError(kind, parentID, &parentKind)
		}
	case "subtask":
		if parentKind != "task" {
			return hierarchyInvalidError(kind, parentID, &parentKind)
		}
	}
	return nil
}

func hierarchyInvalidError(kind string, parentID, parentKind *string) *HierarchyInvalidError {
	rule := ""
	switch kind {
	case "epic":
		rule = "epics must not have a parent"
	case "task":
		rule = "tasks must have an epic parent or no parent"
	case "subtask":
		rule = "subtasks must have a task parent"
	}
	return &HierarchyInvalidError{
		Kind:       kind,
		ParentID:   parentID,
		ParentKind: parentKind,
		Rule:       rule,
	}
}

func requireExistingLabels(conn *sql.DB, labels []string) error {
	for _, label := range labels {
		var existing string
		err := conn.QueryRow(`SELECT name FROM labels WHERE name = ? LIMIT 1`, label).Scan(&existing)
		if err == nil {
			continue
		}
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("create issue: label %q does not exist", label)
		}
		return fmt.Errorf("query label %q: %w", label, err)
	}
	return nil
}

func issueCreateArgv(kind, title, actor, idempotencyKey string, parentID *string, priority, sequence *int, specMDPath *string, labels []string, scheduling string) []string {
	argv := []string{"fid", "issue", "create", "--kind", kind, "--title", title}
	switch kind {
	case "task":
		if parentID != nil {
			argv = append(argv, "--epic", *parentID)
		}
	case "subtask":
		if parentID != nil {
			argv = append(argv, "--parent", *parentID)
		}
	}
	if priority != nil {
		argv = append(argv, "--priority", fmt.Sprintf("%d", *priority))
	}
	if sequence != nil {
		argv = append(argv, "--sequence", fmt.Sprintf("%d", *sequence))
	}
	if specMDPath != nil {
		argv = append(argv, "--spec-md", *specMDPath)
	}
	for _, label := range labels {
		argv = append(argv, "--label", label)
	}
	if scheduling != "" {
		argv = append(argv, "--scheduling", scheduling)
	}
	argv = append(argv, "--actor", actor, "--idempotency-key", idempotencyKey)
	return argv
}
