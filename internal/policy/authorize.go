package policy

import (
	"fmt"
	"strings"
)

type PermissionDeniedError struct {
	Role   string
	Action string
}

func (err *PermissionDeniedError) Error() string {
	if err == nil {
		return "<nil>"
	}
	return fmt.Sprintf("role %q not permitted for action %q", err.Role, err.Action)
}

func Authorize(workspaceRoot, role, action string) error {
	if err := RequireKnownRole(workspaceRoot, role); err != nil {
		return err
	}

	row, ok := MatrixEntry(action)
	if !ok {
		return fmt.Errorf("authorize action %q: no permissions matrix entry", strings.TrimSpace(action))
	}

	trimmedRole := strings.TrimSpace(role)
	allowed, err := row.allows(trimmedRole)
	if err != nil {
		return err
	}
	if allowed {
		return nil
	}

	return &PermissionDeniedError{
		Role:   trimmedRole,
		Action: row.Action,
	}
}
