package main

import (
	"errors"
	"path/filepath"
	"testing"

	"github.com/mbell-nzl/fid/internal/policy"
	"github.com/mbell-nzl/fid/internal/workflow"
)

func TestEnvelopeAuthorityRequiresEveryBindingAndRejectsRoleAttribution(t *testing.T) {
	t.Parallel()
	workspace, issueID, registration := finalizedEnvelopeWorkspace(t)
	inputPath := filepath.Join(workspace, "authority-envelope.json")
	input := envelopeCreateInput(t, workspace, issueID, registration)
	grant := input["authority_grants"].([]any)[0].(map[string]any)
	grant["actions"] = []string{"claim.acquire"}
	grant["actor_selector"] = map[string]any{
		"actor_id": "actor-builder", "required_capabilities": []string{"claim.acquire"},
	}
	writeActivationInput(t, inputPath, input)

	created := executeCommandForTest(workspace,
		"--project", workspace, "--json", "envelope", "create", "--actor", "actor-builder",
		"--input", inputPath, "--idempotency-key", "authority-envelope-create",
	)
	if created.err != nil || created.stderr != "" {
		t.Fatalf("envelope create = %#v, want stdout-only success", created)
	}
	envelopeID := assertEnvelopeShape(t, created.stdout, "draft")

	activatePath := filepath.Join(workspace, "authority-envelope-activate.json")
	writeActivationInput(t, activatePath, map[string]any{"envelope_id": envelopeID, "expected_revision": 1})
	activated := executeCommandForTest(workspace,
		"--project", workspace, "--json", "envelope", "activate", "--actor", "actor-builder",
		"--input", activatePath, "--idempotency-key", "authority-envelope-activate",
	)
	if activated.err != nil || activated.stderr != "" {
		t.Fatalf("envelope activate = %#v, want stdout-only success", activated)
	}

	request := workflow.EnvelopeAuthorityRequest{
		ActorID:    "actor-builder",
		Action:     "claim.acquire",
		Capability: "claim.acquire",
		IssueID:    issueID,
		EnvelopeID: envelopeID,
		Revision:   1,
	}
	if err := workflow.AuthorizeEnvelopeGrant(workspace, issueMigrationsDir(), request); err != nil {
		t.Fatalf("AuthorizeEnvelopeGrant(all match) error = %v", err)
	}

	for _, testCase := range []struct {
		name  string
		apply func(*workflow.EnvelopeAuthorityRequest)
	}{
		{"actor", func(request *workflow.EnvelopeAuthorityRequest) { request.ActorID = "actor-other" }},
		{"action", func(request *workflow.EnvelopeAuthorityRequest) { request.Action = "run.start" }},
		{"capability", func(request *workflow.EnvelopeAuthorityRequest) { request.Capability = "run.start" }},
		{"issue", func(request *workflow.EnvelopeAuthorityRequest) { request.IssueID = "ISS-other" }},
		{"envelope", func(request *workflow.EnvelopeAuthorityRequest) { request.EnvelopeID = "ENV-other" }},
		{"revision", func(request *workflow.EnvelopeAuthorityRequest) { request.Revision = 2 }},
	} {
		t.Run("denies mismatched "+testCase.name, func(t *testing.T) {
			candidate := request
			testCase.apply(&candidate)
			assertEnvelopeAuthorityDenied(t, workspace, candidate)
		})
	}

	t.Run("denies without an active activation", func(t *testing.T) {
		candidate := request
		assertEnvelopeAuthorityDenied(t, clonePristineWorkspace(t), candidate)
	})

	for _, operation := range []string{
		"artifact.verify", "claim.acquire", "claim.release", "claim.renew", "envelope.activate", "envelope.complete",
		"envelope.create", "evidence.link", "gate.register", "gate.run", "issue.transition", "run.close", "run.start",
	} {
		t.Run("seeded developer attribution cannot authorize "+operation, func(t *testing.T) {
			candidate := request
			candidate.ActorID = "developer"
			candidate.Action = operation
			candidate.Capability = operation
			assertEnvelopeAuthorityDenied(t, workspace, candidate)
		})
	}
}

func assertEnvelopeAuthorityDenied(t *testing.T, workspace string, request workflow.EnvelopeAuthorityRequest) {
	t.Helper()
	err := workflow.AuthorizeEnvelopeGrant(workspace, issueMigrationsDir(), request)
	var denied *policy.EnvelopeAuthorityDeniedError
	if !errors.As(err, &denied) {
		t.Fatalf("AuthorizeEnvelopeGrant(%+v) error = %T %v, want *EnvelopeAuthorityDeniedError", request, err, err)
	}
}
