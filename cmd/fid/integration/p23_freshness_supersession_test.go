//go:build integration

package integration

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"
)

func TestP23FreshnessAndOutcomeSupersessionRealProcessJourney(t *testing.T) {
	actions := []string{"claim.acquire", "claim.release", "envelope.activate", "envelope.create", "evidence.link", "gate.register", "review.record", "run.close", "run.observe", "run.start"}
	a := p22NewFIDRecoveryAdapterWithClaimModeAt(t, integrationFIDBinary(t), t.TempDir(), "agent", 3600, []string{}, actions, "mutation", p22ExecutionFixture{ReviewObligation: true, HumanObligations: true})
	p23CloseAndRelease(t, a, a.actor, a.runID, a.claimID, "freshness-builder")

	claim := p23AcquireClaim(t, a, "actor-reviewer-two", "review", "freshness-review")
	run := p23StartRun(t, a, "actor-reviewer-two", claim, "freshness-review")
	recorded := p23ReviewRecord(t, a, "actor-reviewer-two", "freshness-review", p23ReviewInput(a, p21PathString(t, run, "id"), p23Number(t, run["version"]), "pass", []any{}, nil))
	p23AssertReviewSuccess(t, "freshness Review", recorded)
	reviewID := p23ReviewFromResponse(t, recorded)["id"].(string)
	linked := p23EvidenceLink(t, a, "freshness-link", reviewID)
	p23AssertReviewSuccess(t, "freshness Review Evidence", linked)
	beforeReview := p23ReviewGet(t, a, reviewID)

	changed := cloneP23Repository(a.repository)
	changed["tree_oid"] = strings.Repeat("f", 40)
	initial := p23EvidenceEvaluateTarget(t, a, "freshness-initial", changed)
	reviewEvaluation := p23EvaluationByID(t, initial, "OBL-p22-corpus")
	p23AssertEvaluation(t, reviewEvaluation, "unsatisfied", "valid", "stale", "not-applicable", "evidence-stale")
	if afterReview := p23ReviewGet(t, a, reviewID); !reflect.DeepEqual(afterReview, beforeReview) {
		t.Fatalf("tree change rewrote Review: before=%#v after=%#v", beforeReview, afterReview)
	}

	seedP23HumanOutcome(t, a, "OUT-release-approved", "release", "approved", "2026-09-20T23:00:00Z", "EVT-20991231-900001", "EVD-release-approved")
	seedP23HumanOutcome(t, a, "OUT-security-rejected", "security", "rejected", "2026-09-20T23:00:00Z", "EVT-20991231-900003", "")
	initialOutcomes := p23OutcomeRows(t, a)
	initial = p23EvidenceEvaluateTarget(t, a, "freshness-human-initial", changed)
	p23AssertEvaluation(t, p23EvaluationByID(t, initial, "OBL-human-release"), "satisfied", "valid", "fresh", "not-applicable", "")
	p23AssertEvaluation(t, p23EvaluationByID(t, initial, "OBL-human-security"), "unsatisfied", "invalid", "fresh", "not-applicable", "human-outcome-absent")

	seedP23HumanOutcome(t, a, "OUT-release-rejected", "release", "rejected", "2026-09-20T01:00:00Z", "EVT-20991231-900002", "")
	seedP23HumanOutcome(t, a, "OUT-security-approved", "security", "approved", "2026-09-20T01:00:00Z", "EVT-20991231-900004", "EVD-security-approved")
	allOutcomes := p23OutcomeRows(t, a)
	after := p23EvidenceEvaluateTarget(t, a, "freshness-human-superseded", changed)
	p23AssertEvaluation(t, p23EvaluationByID(t, after, "OBL-human-release"), "unsatisfied", "valid", "fresh", "not-applicable", "outcome-superseded")
	p23AssertEvaluation(t, p23EvaluationByID(t, after, "OBL-human-security"), "satisfied", "valid", "fresh", "not-applicable", "")
	if reread := p23OutcomeRows(t, a); !reflect.DeepEqual(reread, allOutcomes) || !p23RowsRetained(reread, initialOutcomes) {
		t.Fatalf("Outcome evaluation rewrote readable records: initial=%#v all=%#v reread=%#v", initialOutcomes, allOutcomes, reread)
	}
}

func p23RowsRetained(all, retained [][]string) bool {
	for _, want := range retained {
		found := false
		for _, got := range all {
			if reflect.DeepEqual(got, want) {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

func cloneP23Repository(repository map[string]any) map[string]any {
	cloned := make(map[string]any, len(repository))
	for key, value := range repository {
		cloned[key] = value
	}
	return cloned
}

func p23EvidenceEvaluateTarget(t *testing.T, a *p22FIDRecoveryAdapter, key string, repository map[string]any) p22ConformanceResponse {
	t.Helper()
	response := a.p22Run(t, "evidence", "evaluate", "--input", p21Input(t, a.workspace, key, map[string]any{"obligation_id": nil, "target_repository": repository}))
	p23AssertReviewSuccess(t, "evidence evaluate "+key, response)
	return response
}

func p23EvaluationByID(t *testing.T, response p22ConformanceResponse, obligationID string) map[string]any {
	t.Helper()
	for _, item := range p21Path(t, p23Envelope(t, response.Stdout), "result", "evaluations").([]any) {
		evaluation := item.(map[string]any)
		if evaluation["obligation_id"] == obligationID {
			return evaluation
		}
	}
	t.Fatalf("evaluation %q missing from %s", obligationID, response.Stdout)
	return nil
}

func p23AssertEvaluation(t *testing.T, evaluation map[string]any, status, validity, freshness, availability, reason string) {
	t.Helper()
	if evaluation["status"] != status || evaluation["validity"] != validity || evaluation["freshness"] != freshness || evaluation["artifact_availability"] != availability {
		t.Fatalf("evaluation = %#v, want %s/%s/%s/%s", evaluation, status, validity, freshness, availability)
	}
	reasons := evaluation["reasons"].([]any)
	if reason == "" && len(reasons) != 0 {
		t.Fatalf("evaluation reasons = %#v, want none", reasons)
	}
	if reason != "" && (len(reasons) != 1 || reasons[0].(map[string]any)["code"] != reason) {
		t.Fatalf("evaluation reasons = %#v, want %s", reasons, reason)
	}
}

func seedP23HumanOutcome(t *testing.T, a *p22FIDRecoveryAdapter, id, class, result, createdAt, eventID, evidenceID string) {
	t.Helper()
	conn := p22FaultDB(t, a.workspace)
	defer func() { p22FaultOK(t, conn.Close()) }()
	var activationID string
	p22FaultOK(t, conn.QueryRow(`SELECT id FROM kit_activations WHERE status='finalized'`).Scan(&activationID))
	_, err := conn.Exec(`INSERT INTO human_outcomes (id,actor_id,activation_id,activation_version,profile_ref,class,subject_type,subject_id,result,scope_issue_id,reasoning,references_json,created_at,created_event_id) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?)`, id, "human:task7-fixture", activationID, 1, "profile:task7-fixture", class, "issue", a.issueID, result, a.issueID, "seeded Task 7 projection fixture", `[]`, createdAt, eventID)
	p22FaultOK(t, err)
	if evidenceID == "" {
		return
	}
	repositoryJSON, err := json.Marshal(a.repository)
	p22FaultOK(t, err)
	_, err = conn.Exec(`INSERT INTO evidence_records (id,obligation_id,evidence_type,receipt_id,review_id,human_outcome_id,producer,trust_source,target_repository_json,validity,created_at) VALUES (?,?, 'human-outcome',NULL,NULL,?,'human','activated-human',?,'valid',?)`, evidenceID, "OBL-human-"+class, id, string(repositoryJSON), createdAt)
	p22FaultOK(t, err)
}

func p23OutcomeRows(t *testing.T, a *p22FIDRecoveryAdapter) [][]string {
	t.Helper()
	conn := p22FaultDB(t, a.workspace)
	defer func() { p22FaultOK(t, conn.Close()) }()
	rows, err := conn.Query(`SELECT id,actor_id,activation_id,CAST(activation_version AS TEXT),profile_ref,class,subject_type,subject_id,result,scope_issue_id,reasoning,references_json,created_at,created_event_id FROM human_outcomes ORDER BY created_event_id`)
	p22FaultOK(t, err)
	defer func() { p22FaultOK(t, rows.Close()) }()
	result := [][]string{}
	for rows.Next() {
		values := make([]string, 14)
		p22FaultOK(t, rows.Scan(p23StringPointers(values)...))
		result = append(result, values)
	}
	p22FaultOK(t, rows.Err())
	return result
}

func p23StringPointers(values []string) []any {
	result := make([]any, len(values))
	for index := range values {
		result[index] = &values[index]
	}
	return result
}
