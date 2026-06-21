package policy

import (
	"errors"
	"testing"
)

func TestAuthorizePermitsAllowedCombinationAndRejectsDeniedCombination(t *testing.T) {
	t.Parallel()

	workspace := newPolicyWorkspace(t)

	if err := Authorize(workspace, "developer", "issue.create.task"); err != nil {
		t.Fatalf("Authorize(allowed) error = %v", err)
	}

	err := Authorize(workspace, "developer", "issue.create.epic")
	if err == nil {
		t.Fatal("Authorize(denied) error = nil, want permission denied error")
	}

	var deniedErr *PermissionDeniedError
	if !errors.As(err, &deniedErr) {
		t.Fatalf("Authorize(denied) error = %T, want *PermissionDeniedError", err)
	}
	if deniedErr.Role != "developer" {
		t.Fatalf("PermissionDeniedError.Role = %q, want %q", deniedErr.Role, "developer")
	}
	if deniedErr.Action != "issue.create.epic" {
		t.Fatalf("PermissionDeniedError.Action = %q, want %q", deniedErr.Action, "issue.create.epic")
	}
	if deniedErr.Error() != `role "developer" not permitted for action "issue.create.epic"` {
		t.Fatalf("PermissionDeniedError.Error() = %q, want %q", deniedErr.Error(), `role "developer" not permitted for action "issue.create.epic"`)
	}
}

func TestAuthorizePreservesMissingAndUnknownRoleErrors(t *testing.T) {
	t.Parallel()

	workspace := newPolicyWorkspace(t)

	err := Authorize(workspace, "  ", "issue.create.task")
	if err == nil {
		t.Fatal("Authorize(missing role) error = nil, want missing role error")
	}

	var missingErr *MissingRoleError
	if !errors.As(err, &missingErr) {
		t.Fatalf("Authorize(missing role) error = %T, want *MissingRoleError", err)
	}

	err = Authorize(workspace, "intern", "issue.create.task")
	if err == nil {
		t.Fatal("Authorize(unknown role) error = nil, want unknown role error")
	}

	var unknownErr *UnknownRoleError
	if !errors.As(err, &unknownErr) {
		t.Fatalf("Authorize(unknown role) error = %T, want *UnknownRoleError", err)
	}
	if unknownErr.Role != "intern" {
		t.Fatalf("UnknownRoleError.Role = %q, want %q", unknownErr.Role, "intern")
	}
}
