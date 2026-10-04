//go:build integration

package integration

import (
	"reflect"
	"testing"
)

func TestP24RunStartConflictNamesChangedClaimVersion(t *testing.T) {
	t.Parallel()
	a := p22NewFIDRecoveryAdapter(t)
	before := p21JournalEventCount(t, a.workspace)
	input := map[string]any{
		"issue_id": a.issueID, "envelope_id": p23EnvelopeID(t, a),
		"envelope_revision": 1, "claim_id": a.claimID, "claim_version": 2,
		"executor":      map[string]any{"id": "executor-corpus-prior", "source": "activated-harness"},
		"obligation_id": nil, "repository": a.repository,
	}
	response := p23RunInput(t, a, "run", "start", a.actor, "p22-corpus-run-start", input)
	assertP22ProcessFailure(t, "run.start", "idempotency-conflict", 4, response)
	fields := p21Path(t, p23Envelope(t, response.Stderr), "error", "details", "differing_fields")
	if !reflect.DeepEqual(fields, []any{"claim_version"}) {
		t.Fatalf("differing_fields = %#v, want [claim_version]", fields)
	}
	if got := p21JournalEventCount(t, a.workspace); got != before {
		t.Fatalf("conflict appended Event: %d -> %d", before, got)
	}
}
