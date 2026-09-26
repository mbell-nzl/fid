//go:build integration

package integration

import (
	"context"
	"encoding/json"
	"fmt"
	"path/filepath"
	"sort"
	"testing"

	"github.com/mbell-nzl/fid/internal/evidence"
)

func mustP23JSON(t *testing.T, value any) string {
	t.Helper()
	encoded, err := json.Marshal(value)
	p22FaultOK(t, err)
	return string(encoded)
}

// TestP23GateRunAuthorizesEveryEnvelopeGrantThroughBuiltBinary is the review-7
// I-5 conformance evidence (DEC-252). One active Envelope carries three Grants:
// actor-denied (first sorted) and actor-reviewer-two (last sorted) are granted
// gate.run; actor-recovery (between them) holds the gate.run capability in its
// activated profile but has no gate.run Grant. Before the repair only the Actor
// whose Grant sorts first could run a gate.
func TestP23GateRunAuthorizesEveryEnvelopeGrantThroughBuiltBinary(t *testing.T) {
	t.Parallel()
	actions := []string{"claim.acquire", "claim.release", "envelope.activate", "envelope.create", "gate.register", "gate.run", "review.record", "run.close", "run.observe", "run.start"}
	a := p22NewFIDRecoveryAdapterWithClaimModeAt(t, integrationFIDBinary(t), t.TempDir(), "agent", 3600, []string{}, actions, "mutation")
	issue := p21Run(t, a.binary, a.workspace, "--project", a.workspace, "--json", "issue", "create", "--kind", "task", "--title", "Multi-grant gate.run", "--actor", "agent:developer", "--idempotency-key", "p23-gate-grants-issue")
	issueID := p21PathString(t, issue, "result", "issue", "id")
	registered := p23RunInput(t, a, "gate", "register", a.actor, "p23-gate-grants-register", map[string]any{"definition_id": "gate-p22-corpus", "definition_version": 1, "worktree_id": a.repository["worktree_id"], "scope": map[string]any{"issue_id": issueID}, "repository": a.repository})
	p23AssertReviewSuccess(t, "register gate for the multi-grant Issue", registered)
	registration := p21Path(t, p23Envelope(t, registered.Stdout), "result", "registration").(map[string]any)

	all := []string{"claim.acquire", "claim.release", "envelope.activate", "envelope.create", "gate.register", "gate.run", "review.record", "run.close", "run.observe", "run.start"}
	withoutGate := []string{"claim.acquire", "claim.release", "envelope.activate", "envelope.create", "gate.register", "review.record", "run.close", "run.observe", "run.start"}
	grant := func(actor string, granted []string) map[string]any {
		return map[string]any{"actions": granted, "actor_selector": map[string]any{"actor_id": actor, "required_capabilities": granted}, "scope": map[string]any{"issue_id": issueID}}
	}
	obligation := map[string]any{"id": "OBL-p23-gate-grants", "type": "gate-evidence", "applicability": "required", "guarded_actions": []string{"envelope.complete", "issue.transition-terminal-success"}, "subject": map[string]any{"type": "issue", "id": issueID}, "registration": map[string]any{"id": registration["id"], "version": registration["version"]}, "accepted_evidence_types": []string{"receipt"}, "accepted_trust_sources": []string{"fid-executed"}, "freshness_target": "repository-tree", "artifact_policy": "required-available", "independence": "fid-executed", "required_lane": nil, "source": a.measurement}
	created := p23RunInput(t, a, "envelope", "create", a.actor, "p23-gate-grants-envelope", map[string]any{"outcome": "Authorize gate.run per Grant", "non_goals": []string{}, "risk_result": map[string]any{"source": a.measurement, "result_digest": "sha256:" + p22Ones(64)}, "repository_scope": a.repository, "sources": []any{a.measurement}, "primary_issue_id": issueID, "related_issue_ids": []string{}, "authority_grants": []any{grant("actor-denied", all), grant("actor-recovery", withoutGate), grant("actor-reviewer-two", all)}, "obligations": []any{obligation}, "budgets": []any{map[string]any{"name": "attempts", "scope": "envelope", "unit": "attempts", "maximum": 5, "warning_threshold": nil, "exhaustion_action": "budget-exhausted", "consumed": nil, "remaining": nil}}})
	p23AssertReviewSuccess(t, "create the three-Grant Envelope", created)
	envelopeID := p21PathString(t, p23Envelope(t, created.Stdout), "result", "envelope", "id")
	activated := p23RunInput(t, a, "envelope", "activate", a.actor, "p23-gate-grants-activate", map[string]any{"envelope_id": envelopeID, "expected_revision": 1})
	p23AssertReviewSuccess(t, "activate the three-Grant Envelope", activated)

	gateInput := func(runID string) map[string]any {
		return map[string]any{"run_id": runID, "registration_id": registration["id"], "registration_version": registration["version"], "obligation_id": "OBL-p23-gate-grants", "repository": a.repository}
	}
	session := func(actor string) (runID, claimID string) {
		claim := p23RunInput(t, a, "claim", "acquire", actor, "p23-gate-grants-claim-"+actor, map[string]any{"issue_id": issueID, "envelope_id": envelopeID, "envelope_revision": 1, "worktree_id": a.repository["worktree_id"], "scope": map[string]any{"paths": []string{"base.txt"}}, "mode": "mutation", "ttl_seconds": 3600})
		p23AssertReviewSuccess(t, actor+" acquires a mutation Claim", claim)
		claimValue := p21Path(t, p23Envelope(t, claim.Stdout), "result", "claim").(map[string]any)
		run := p23RunInput(t, a, "run", "start", actor, "p23-gate-grants-run-"+actor, map[string]any{"issue_id": issueID, "envelope_id": envelopeID, "envelope_revision": 1, "claim_id": claimValue["id"], "claim_version": claimValue["version"], "executor": map[string]any{"id": "executor-gate-grants-" + actor, "source": "activated-harness"}, "obligation_id": nil, "repository": a.repository})
		p23AssertReviewSuccess(t, actor+" starts a Run", run)
		return p21PathString(t, p23Envelope(t, run.Stdout), "result", "run", "id"), claimValue["id"].(string)
	}

	for _, actor := range []string{"actor-denied", "actor-reviewer-two"} {
		runID, claimID := session(actor)
		ran := p23RunInput(t, a, "gate", "run", actor, "p23-gate-grants-gate-"+actor, gateInput(runID))
		p23AssertReviewSuccess(t, actor+" gate.run with its own gate.run Grant", ran)
		if !p23HasJSON(ran.Stdout, `"trust_source":"fid-executed"`) || !p23HasJSON(ran.Stdout, `"actor_id":"`+actor+`"`) || !p23HasJSON(ran.Stdout, `"run_id":"`+runID+`"`) || !p23HasJSON(ran.Stdout, `"outcome":"succeeded"`) {
			t.Fatalf("%s gate.run = %q, want an fid-executed succeeded Receipt for its own Run", actor, ran.Stdout)
		}
		p23CloseAndRelease(t, a, actor, runID, claimID, "p23-gate-grants-close-"+actor)
	}

	runID, claimID := session("actor-recovery")
	events, projection := p21JournalEventCount(t, a.workspace), p23GateProjection(t, a)
	denied := p23RunInput(t, a, "gate", "run", "actor-recovery", "p23-gate-grants-gate-actor-recovery", gateInput(runID))
	// The unchanged adapter maps the grant denial to invalid-argument; the
	// in-process call below pins the underlying cause to the Grant match.
	assertP22ProcessFailure(t, "gate.run", "invalid-argument", 2, denied)
	if after := p21JournalEventCount(t, a.workspace); after != events {
		t.Fatalf("ungranted gate.run Event count = %d, want unchanged %d", after, events)
	}
	if after := p23GateProjection(t, a); after != projection {
		t.Fatalf("ungranted gate.run projection = %s, want unchanged %s", after, projection)
	}
	options := evidence.GateRunOptions{WorkspaceRoot: a.workspace, MigrationsDir: filepath.Join(integrationRepositoryRoot(), "migrations"), ActorID: "actor-recovery", IdempotencyKey: "p23-gate-grants-gate-actor-recovery-inprocess"}
	p22FaultOK(t, json.Unmarshal([]byte(mustP23JSON(t, gateInput(runID))), &options.Request))
	if _, err := evidence.RunGate(context.Background(), options); err == nil || err.Error() != "run gate: authority grant does not permit gate.run" {
		t.Fatalf("ungranted gate.run cause = %v, want the Grant denial", err)
	}
	if after := p21JournalEventCount(t, a.workspace); after != events || p23GateProjection(t, a) != projection {
		t.Fatalf("in-process ungranted gate.run changed Event count %d -> %d or the projection", events, after)
	}
	p23CloseAndRelease(t, a, "actor-recovery", runID, claimID, "p23-gate-grants-close-actor-recovery")
}

// p23GateProjection summarizes the gate-execution projection: row counts of the
// tables gate.run writes, plus every Run's identity and status.
func p23GateProjection(t *testing.T, a *p22FIDRecoveryAdapter) string {
	t.Helper()
	conn := p22FaultDB(t, a.workspace)
	defer func() { p22FaultOK(t, conn.Close()) }()
	var summary []string
	for _, table := range []string{"execution_attempts", "execution_receipts", "evidence_records", "budget_samples", "idempotency_results", "run_budget_bindings"} {
		var count int
		p22FaultOK(t, conn.QueryRow(`SELECT count(*) FROM `+table).Scan(&count))
		summary = append(summary, fmt.Sprintf("%s=%d", table, count))
	}
	rows, err := conn.Query(`SELECT id, status FROM runs ORDER BY id`)
	p22FaultOK(t, err)
	defer func() { p22FaultOK(t, rows.Close()) }()
	for rows.Next() {
		var id, status string
		p22FaultOK(t, rows.Scan(&id, &status))
		summary = append(summary, id+":"+status)
	}
	p22FaultOK(t, rows.Err())
	sort.Strings(summary)
	return fmt.Sprint(summary)
}
