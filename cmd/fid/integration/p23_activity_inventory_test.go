//go:build integration

package integration

import (
	"reflect"
	"sort"
	"strings"
	"testing"
)

func TestP23TenNameActivityInventoryFromProducingLeaves(t *testing.T) {
	t.Parallel()
	actions := []string{
		"claim.acquire", "claim.release", "concern.dispose", "concern.record", "decision.record-reference",
		"deviation.record", "envelope.activate", "envelope.create", "finding.address", "gate.register",
		"human-outcome.record", "review.record", "run.close", "run.observe", "run.start",
	}
	a := p22NewFIDRecoveryAdapterAt(t, integrationFIDBinary(t), t.TempDir(), "agent", 3600, []string{}, actions)
	a.task13ReplayCorpus = true
	var producers []p22ConformanceResponse
	appendSuccess := func(name string, response p22ConformanceResponse) p22ConformanceResponse {
		t.Helper()
		p23AssertReviewSuccess(t, name, response)
		producers = append(producers, response)
		return response
	}
	requirement := appendSuccess("requirement.record-reference", p23RunInput(t, a, "requirement", "record-reference", a.actor, "activity-requirement", map[string]any{"id": "REQ-activity", "source": p23DecisionSource(a, true), "status": "active"}))
	p23AssertActivitySubject(t, a.workspace, requirement, "requirement-referenced", "requirement", "REQ-activity", 1)
	linked := appendSuccess("requirement.link", p23RunInput(t, a, "requirement", "link", a.actor, "activity-link", map[string]any{"requirement_id": "REQ-activity", "issue_id": a.issueID, "relation": "verifies"}))
	p23AssertActivitySubject(t, a.workspace, linked, "requirement-linked", "requirement", "REQ-activity", 1)
	decision := appendSuccess("decision.record-reference", p23RecordDecision(t, a, "activity", "release", "approved", []string{"decision.record-reference"}, p23DecisionSource(a, true)))
	p23AssertActivitySubject(t, a.workspace, decision, "decision-referenced", "decision", "DREF-activity", 1)
	concern := appendSuccess("concern.record", p23RunInput(t, a, "concern", "record", a.actor, "activity-concern", map[string]any{"id": "CON-activity", "subject": map[string]any{"type": "issue", "id": a.issueID}, "scope": map[string]any{"paths": []string{}}, "severity": "minor", "summary": "activity inventory", "reasoning": "check producer", "references": []any{}}))
	p23AssertActivity(t, a.workspace, concern, "concern-recorded", a.issueID, 1)
	disposed := appendSuccess("concern.dispose", p23RunInput(t, a, "concern", "dispose", a.actor, "activity-dispose", map[string]any{"concern_id": "CON-activity", "action": "rejected", "issue_id": nil, "supersedes_concern_id": nil, "supersedes_finding_id": nil, "reason": "inventory complete"}))
	p23AssertActivity(t, a.workspace, disposed, "concern-disposed", a.issueID, 1)
	deviation := appendSuccess("deviation.record", p23RunInput(t, a, "deviation", "record", a.actor, "activity-deviation", map[string]any{"id": "DEV-activity", "subject": map[string]any{"type": "issue", "id": a.issueID}, "expected": "contract", "observed": "departure", "authority": map[string]any{"actor_id": "human-reviewer", "reference": nil}, "disposition": "pending-decision", "reasoning": "inventory", "references": []any{}, "supersedes_deviation_id": nil}))
	p23AssertActivity(t, a.workspace, deviation, "deviation-recorded", a.issueID, 1)
	human := appendSuccess("human-outcome.record", p23RecordHumanOutcome(t, a, "activity", "actor-human-outcome", "release", "approved"))
	p23AssertActivity(t, a.workspace, human, "human-outcome-recorded", a.issueID, 1)

	p23CloseAndRelease(t, a, a.actor, a.runID, a.claimID, "activity-builder")
	claim := p23AcquireClaim(t, a, "actor-denied", "review", "activity-first")
	run := p23StartRun(t, a, "actor-denied", claim, "activity-first")
	first := appendSuccess("review.record", p23ReviewRecord(t, a, "actor-denied", "activity-first", p23ReviewInput(a, p21PathString(t, run, "id"), p23Number(t, run["version"]), "changes-required", []any{p23Finding(a.issueID, "FND-activity", "major", true)}, nil)))
	p23AssertActivity(t, a.workspace, first, "review-recorded", a.issueID, 1)
	p23CloseAndRelease(t, a, "actor-denied", p21PathString(t, run, "id"), p21PathString(t, claim, "id"), "activity-first")
	mutationClaim := p23AcquireClaim(t, a, a.actor, "mutation", "activity-address")
	addressed := appendSuccess("finding.address", p23FindingAddress(t, a, a.actor, "activity-address", map[string]any{"finding_id": "FND-activity", "claim_id": mutationClaim["id"], "expected_claim_version": mutationClaim["version"], "reason": "addressed for inventory"}))
	p23AssertActivity(t, a.workspace, addressed, "finding-addressed", a.issueID, 1)
	p23ReleaseClaim(t, a, a.actor, mutationClaim, "activity-address")
	secondClaim := p23AcquireClaim(t, a, "actor-reviewer-two", "review", "activity-second")
	secondRun := p23StartRun(t, a, "actor-reviewer-two", secondClaim, "activity-second")
	second := appendSuccess("review.record resolve", p23ReviewRecord(t, a, "actor-reviewer-two", "activity-second", p23ReviewInput(a, p21PathString(t, secondRun, "id"), p23Number(t, secondRun["version"]), "pass", []any{}, []any{"FND-activity"})))
	p23AssertActivity(t, a.workspace, second, "review-recorded", a.issueID, 1)
	p23AssertActivity(t, a.workspace, second, "finding-resolved", a.issueID, 1)

	ids := make([]string, 0, len(producers))
	for _, response := range producers {
		ids = append(ids, p21PathString(t, p23Envelope(t, response.Stdout), "event", "id"))
	}
	conn := p22FaultDB(t, a.workspace)
	defer func() { p22FaultOK(t, conn.Close()) }()
	placeholders := strings.TrimSuffix(strings.Repeat("?,", len(ids)), ",")
	args := make([]any, len(ids))
	for index, id := range ids {
		args[index] = id
	}
	rows, err := conn.Query(`SELECT DISTINCT event FROM activity_log WHERE journal_event_id IN (`+placeholders+`) ORDER BY event`, args...)
	p22FaultOK(t, err)
	defer func() { p22FaultOK(t, rows.Close()) }()
	var got []string
	for rows.Next() {
		var name string
		p22FaultOK(t, rows.Scan(&name))
		got = append(got, name)
	}
	p22FaultOK(t, rows.Err())
	want := []string{"requirement-referenced", "requirement-linked", "decision-referenced", "review-recorded", "finding-addressed", "finding-resolved", "human-outcome-recorded", "concern-recorded", "concern-disposed", "deviation-recorded"}
	sort.Strings(want)
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("producing-leaf activity names = %v, want exactly %v", got, want)
	}
}
