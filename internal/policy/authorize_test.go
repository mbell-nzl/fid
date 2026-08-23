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

func TestAuthorizeRejectsUnknownCompatibilityRole(t *testing.T) {
	t.Parallel()

	for _, tt := range []struct {
		name string
		role string
	}{
		{name: "unknown name remains rejected after role-free bootstrap", role: "intern"},
	} {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			workspace := newPolicyJournalWorkspace(t)
			writePolicyJournalFile(t, workspace, "2026-08-23.jsonl", policyBootstrapEvent())

			err := Authorize(workspace, tt.role, "issue.create.task")
			var unknownErr *UnknownRoleError
			if !errors.As(err, &unknownErr) {
				t.Fatalf("Authorize() error = %T %v, want *UnknownRoleError", err, err)
			}
			if unknownErr.Role != tt.role {
				t.Fatalf("UnknownRoleError.Role = %q, want %q", unknownErr.Role, tt.role)
			}
		})
	}
}

func TestAuthorizeIssueLabelActionsMatchMatrixExpectations(t *testing.T) {
	t.Parallel()

	workspace := newPolicyWorkspace(t)

	for _, action := range []string{"issue.label.add", "issue.label.remove"} {
		for _, role := range []string{"developer", "tl", "workflow_runner", "human"} {
			if err := Authorize(workspace, role, action); err != nil {
				t.Fatalf("Authorize(%q, %q) error = %v, want allowed", role, action, err)
			}
		}

		err := Authorize(workspace, "qa_reviewer", action)
		if err == nil {
			t.Fatalf("Authorize(qa_reviewer, %q) error = nil, want permission denied", action)
		}

		var deniedErr *PermissionDeniedError
		if !errors.As(err, &deniedErr) {
			t.Fatalf("Authorize(qa_reviewer, %q) error = %T, want *PermissionDeniedError", action, err)
		}
		if deniedErr.Role != "qa_reviewer" || deniedErr.Action != action {
			t.Fatalf("PermissionDeniedError = %#v, want role qa_reviewer action %q", deniedErr, action)
		}
	}
}

func TestAuthorizeMaintenanceVacuumMatchesMatrixExpectations(t *testing.T) {
	t.Parallel()

	workspace := newPolicyWorkspace(t)
	for _, role := range []string{"tl", "human"} {
		if err := Authorize(workspace, role, "maintenance.vacuum"); err != nil {
			t.Fatalf("Authorize(%q, maintenance.vacuum) error = %v, want allowed", role, err)
		}
	}

	for _, role := range []string{"developer", "qa_reviewer", "qa_lead", "tl_lite", "workflow_runner"} {
		err := Authorize(workspace, role, "maintenance.vacuum")
		var deniedErr *PermissionDeniedError
		if !errors.As(err, &deniedErr) {
			t.Fatalf("Authorize(%q, maintenance.vacuum) error = %T, want *PermissionDeniedError", role, err)
		}
	}
}
