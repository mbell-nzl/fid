//go:build integration

package integration

import "testing"

// contracts.md section 8 requires each disposition target to resolve before
// the Concern can leave its open state.
func TestP23ConcernDispositionReferenceRejections(t *testing.T) {
	t.Parallel()
	a := p22NewFIDRecoveryAdapterAt(t, integrationFIDBinary(t), t.TempDir(), "agent", 3600, []string{},
		[]string{"claim.acquire", "concern.dispose", "concern.record", "envelope.activate", "envelope.create", "gate.register", "run.start"})
	created := p23RunInput(t, a, "concern", "record", a.actor, "p23-reference-concern", map[string]any{
		"id": "CON-reference-probe", "subject": map[string]any{"type": "issue", "id": a.issueID},
		"scope": map[string]any{"paths": []string{}}, "severity": "minor", "summary": "reference probe",
		"reasoning": "check disposition references", "references": []any{},
	})
	p23AssertReviewSuccess(t, "record reference-probe Concern", created)
	for _, test := range []struct {
		name, action, field, id, code, subject string
	}{
		{"issue", "carried-forward", "issue_id", "ISS-missing", "issue-not-found", "issue"},
		{"concern", "superseded", "supersedes_concern_id", "CON-missing", "concern-not-found", "concern"},
		{"finding", "superseded", "supersedes_finding_id", "FND-missing", "finding-not-found", "finding"},
	} {
		t.Run(test.name, func(t *testing.T) {
			input := map[string]any{"concern_id": "CON-reference-probe", "action": test.action,
				"issue_id": nil, "supersedes_concern_id": nil, "supersedes_finding_id": nil, "reason": "unresolved target"}
			input[test.field] = test.id
			args := []string{"--project", a.workspace, "--json", "concern", "dispose", "--input",
				p21Input(t, a.workspace, "p23-missing-disposition-"+test.name, input), "--actor", a.actor,
				"--idempotency-key", "p23-missing-disposition-" + test.name}
			fid := p23FIDLeafAdapter{pair: p23HumanLeafPair{Operation: "concern.dispose", CWD: a.workspace}, binary: a.binary}
			p23CompareLeafOutcome(t, test.name, p23ExpectedNotFound("concern.dispose", p23NotFoundCase{
				ID: test.id, Code: test.code, SubjectType: test.subject,
			}), fid.rejection(t, args))
		})
	}
}
