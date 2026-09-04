package workflow

import (
	"errors"
	"fmt"

	"github.com/mbell-nzl/fid/internal/policy"
)

// EnvelopeAuthorityRequest identifies one post-creation protected operation.
type EnvelopeAuthorityRequest struct {
	ActorID    string
	Action     string
	Capability string
	IssueID    string
	EnvelopeID string
	Revision   int
}

// AuthorizeEnvelopeGrant resolves the active Activation Actor before evaluating an Envelope Grant.
func AuthorizeEnvelopeGrant(workspaceRoot, migrationsDir string, request EnvelopeAuthorityRequest) error {
	actor, err := ResolveTrustedActor(workspaceRoot, migrationsDir, request.ActorID)
	if err != nil {
		var unavailable *TrustedActorUnavailableError
		if errors.As(err, &unavailable) {
			return deniedEnvelopeAuthority(request)
		}
		return fmt.Errorf("authorize envelope grant: resolve active actor: %w", err)
	}
	envelope, err := GetEnvelope(workspaceRoot, migrationsDir, request.EnvelopeID)
	if err != nil {
		var notFound *EnvelopeNotFoundError
		if errors.As(err, &notFound) {
			return deniedEnvelopeAuthority(request)
		}
		return fmt.Errorf("authorize envelope grant: load envelope: %w", err)
	}
	grants := make([]policy.EnvelopeGrant, 0, len(envelope.AuthorityGrants))
	for _, grant := range envelope.AuthorityGrants {
		grants = append(grants, policy.EnvelopeGrant{
			ActorID:              grant.ActorSelector.ActorID,
			RequiredCapabilities: grant.ActorSelector.RequiredCapabilities,
			Actions:              grant.Actions,
			IssueID:              grant.Scope.IssueID,
		})
	}
	return policy.EvaluateEnvelopeGrant(policy.EnvelopeGrantEvaluation{
		Active:            true,
		ActorID:           actor.ActorID,
		ActorCapabilities: actor.Capabilities,
		Action:            request.Action,
		Capability:        request.Capability,
		IssueID:           request.IssueID,
		EnvelopeID:        request.EnvelopeID,
		EnvelopeStatus:    envelope.Status,
		EnvelopeRevision:  envelope.CurrentRevision,
		Revision:          request.Revision,
		Grants:            grants,
	})
}

func deniedEnvelopeAuthority(request EnvelopeAuthorityRequest) error {
	return &policy.EnvelopeAuthorityDeniedError{ActorID: request.ActorID, Action: request.Action, Capability: request.Capability}
}
