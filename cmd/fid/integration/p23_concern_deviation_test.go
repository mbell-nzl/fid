//go:build integration

package integration

import (
	"context"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/mbell-nzl/fid/internal/config"
	fiddb "github.com/mbell-nzl/fid/internal/db"
	"github.com/mbell-nzl/fid/internal/evidence"
)

func TestP23ConcernDeviationRealProcessJourney(t *testing.T) {
	t.Parallel()
	actions := []string{
		"claim.acquire", "claim.release", "concern.dispose", "concern.record", "deviation.record",
		"envelope.activate", "envelope.complete", "envelope.create", "evidence.link", "gate.register",
		"gate.run", "issue.transition", "run.close", "run.observe", "run.observe-budget", "run.start",
	}
	a := p22NewFIDRecoveryAdapterAt(t, integrationFIDBinary(t), t.TempDir(), "agent", 3600, []string{}, actions)
	a.task13ReplayCorpus = true

	subjects := []string{"issue", "envelope", "run", "requirement"}
	concerns := map[string]map[string]any{}
	for _, subjectType := range subjects {
		subjectID := "unresolved-" + subjectType
		switch subjectType {
		case "issue":
			subjectID = a.issueID
		case "envelope":
			subjectID = p23EnvelopeID(t, a)
		case "run":
			subjectID = a.runID
		}
		response := p23ConcernMutation(t, a, "record", "record-"+subjectType, map[string]any{
			"id": "CON-" + subjectType, "subject": map[string]any{"type": subjectType, "id": subjectID},
			"scope": map[string]any{"paths": []string{}}, "severity": "blocker", "summary": "Task 8 " + subjectType,
			"reasoning": "records a generic observation", "references": []any{},
		})
		p23AssertReviewSuccess(t, "record "+subjectType+" Concern", response)
		if subjectType != "requirement" {
			p23AssertActivity(t, a.workspace, response, "concern-recorded", a.issueID, 1)
		} else {
			p23AssertActivity(t, a.workspace, response, "concern-recorded", a.issueID, 0)
		}
		concern := p21Path(t, p23Envelope(t, response.Stdout), "result", "concern").(map[string]any)
		if concern["status"] != "open" || concern["disposition"] != nil || concern["actor_id"] != a.actor {
			t.Fatalf("recorded %s Concern = %#v", subjectType, concern)
		}
		if _, found := concern["completion_id"]; found {
			t.Fatal("generic Concern retained a completion link")
		}
		concerns[subjectType] = concern
		p23AssertReplayWithoutEvent(t, a, "record Concern", func() p22ConformanceResponse {
			return p23ConcernMutation(t, a, "record", "record-"+subjectType, map[string]any{
				"id": "CON-" + subjectType, "subject": map[string]any{"type": subjectType, "id": subjectID},
				"scope": map[string]any{"paths": []string{}}, "severity": "blocker", "summary": "Task 8 " + subjectType,
				"reasoning": "records a generic observation", "references": []any{},
			})
		})
	}
	p23AssertConcernHasNoCompletionColumn(t, a.workspace)

	self := assertConcernDisposal(t, a, "CON-issue", "carry", map[string]any{"concern_id": "CON-issue", "action": "carried-forward", "issue_id": a.issueID, "supersedes_concern_id": nil, "supersedes_finding_id": nil, "reason": "continue as this Issue"}, "carried-forward")
	if self["actor_id"] != a.actor || self["disposition"].(map[string]any)["actor_id"] != a.actor {
		t.Fatalf("self-disposed Concern lost Actor provenance: %#v", self)
	}
	differentResponse := p23RunInput(t, a, "concern", "dispose", "actor-reviewer-two", "reject-different-actor", map[string]any{"concern_id": "CON-requirement", "action": "rejected", "issue_id": nil, "supersedes_concern_id": nil, "supersedes_finding_id": nil, "reason": "separate Actor cleared the observation"})
	p23AssertReviewSuccess(t, "dispose Concern with a different Actor", differentResponse)
	p23AssertActivity(t, a.workspace, differentResponse, "concern-disposed", a.issueID, 0)
	different := p21Path(t, p23Envelope(t, differentResponse.Stdout), "result", "concern").(map[string]any)
	differentDisposition := different["disposition"].(map[string]any)
	if different["id"] != "CON-requirement" || different["status"] != "disposed" || differentDisposition["action"] != "rejected" || different["actor_id"] != a.actor || differentDisposition["actor_id"] != "actor-reviewer-two" {
		t.Fatalf("different-Actor disposition merged Actor provenance: %#v", different)
	}
	assertConcernDisposal(t, a, "CON-envelope", "supersede", map[string]any{"concern_id": "CON-envelope", "action": "superseded", "issue_id": nil, "supersedes_concern_id": "CON-run", "supersedes_finding_id": nil, "reason": "newer Concern owns the observation"}, "superseded")

	negativeConcernDisposals(t, a)
	negativeDeviationPredecessors(t, a)
	deviations := recordDeviationChain(t, a)
	issueDeviation := p23RunInput(t, a, "deviation", "record", a.actor, "deviation-issue-activity", map[string]any{"id": "DEV-issue-activity", "subject": map[string]any{"type": "issue", "id": a.issueID}, "expected": "accepted contract", "observed": "issue departure", "authority": map[string]any{"actor_id": "human-reviewer", "reference": nil}, "disposition": "pending-decision", "reasoning": "issue observation", "references": []any{}, "supersedes_deviation_id": nil})
	p23AssertReviewSuccess(t, "record Issue Deviation", issueDeviation)
	p23AssertActivity(t, a.workspace, issueDeviation, "deviation-recorded", a.issueID, 1)
	assertConcernDeviationEvidenceRejected(t, a, "CON-envelope", deviations[0]["id"].(string))

	beforeRebuild := p23DeviationList(t, a)
	p21Run(t, a.binary, a.workspace, "--project", a.workspace, "--json", "rebuild")
	if after := p23DeviationList(t, a); !reflect.DeepEqual(after, beforeRebuild) {
		t.Fatalf("Deviation chain changed across rebuild: before=%#v after=%#v", beforeRebuild, after)
	}

	p22AcceptedCall(t, a, "task8-budget", "run.observe-budget", map[string]any{"value": 0})
	gate, err := evidence.RunGate(context.Background(), p22BudgetGateOptions(t, a, "task8-proof"))
	if err != nil || gate.Receipt == nil {
		t.Fatalf("run qualifying gate: result=%#v error=%v", gate, err)
	}
	linked := p23RunInput(t, a, "evidence", "link", a.actor, "task8-proof-link", map[string]any{"obligation_id": "OBL-p22-corpus", "evidence_type": "receipt", "receipt_id": gate.Receipt.ID, "review_id": nil, "human_outcome_id": nil, "target_repository": a.repository, "required_artifact_ids": []string{}})
	p23AssertReviewSuccess(t, "link qualifying proof", linked)
	p21Run(t, a.binary, a.workspace, "--project", a.workspace, "--json", "issue", "transition", a.issueID, "--to", "in_progress", "--actor", "agent:"+a.actor, "--expected-revision", "1", "--idempotency-key", "task8-in-progress")
	p21Run(t, a.binary, a.workspace, "--project", a.workspace, "--json", "issue", "transition", a.issueID, "--to", "dev_complete", "--actor", "agent:"+a.actor, "--expected-revision", "2", "--idempotency-key", "task8-dev-complete")
	p21Run(t, a.binary, a.workspace, "--project", a.workspace, "--json", "issue", "transition", a.issueID, "--to", "qa_review", "--actor", "agent:"+a.actor, "--expected-revision", "3", "--idempotency-key", "task8-qa-review")
	p21Run(t, a.binary, a.workspace, "--project", a.workspace, "--json", "issue", "transition", a.issueID, "--to", "qa_pass", "--actor", "agent:"+a.actor, "--expected-revision", "4", "--idempotency-key", "task8-qa-pass")
	p23CloseAndRelease(t, a, a.actor, a.runID, a.claimID, "task8-close")
	completed := p23RunInput(t, a, "envelope", "complete", a.actor, "task8-envelope-complete", map[string]any{"envelope_id": p23EnvelopeID(t, a), "expected_revision": 1, "repository": a.repository})
	p23AssertReviewSuccess(t, "complete Envelope with open blocker Concern", completed)
}

func negativeConcernDisposals(t *testing.T, a *p22FIDRecoveryAdapter) {
	t.Helper()
	cases := []struct {
		name, code string
		input      map[string]any
	}{
		{name: "second-disposal", code: "concern-already-disposed", input: map[string]any{"concern_id": "CON-issue", "action": "rejected", "issue_id": nil, "supersedes_concern_id": nil, "supersedes_finding_id": nil, "reason": "second"}},
		{name: "missing-reason", code: "invalid-argument", input: map[string]any{"concern_id": "CON-run", "action": "rejected", "issue_id": nil, "supersedes_concern_id": nil, "supersedes_finding_id": nil, "reason": ""}},
		{name: "unknown-issue", code: "issue-not-found", input: map[string]any{"concern_id": "CON-run", "action": "carried-forward", "issue_id": "ISS-missing", "supersedes_concern_id": nil, "supersedes_finding_id": nil, "reason": "missing"}},
		{name: "unknown-concern", code: "concern-not-found", input: map[string]any{"concern_id": "CON-missing", "action": "rejected", "issue_id": nil, "supersedes_concern_id": nil, "supersedes_finding_id": nil, "reason": "missing"}},
		{name: "unknown-finding", code: "finding-not-found", input: map[string]any{"concern_id": "CON-run", "action": "superseded", "issue_id": nil, "supersedes_concern_id": nil, "supersedes_finding_id": "FND-missing", "reason": "missing"}},
		{name: "self-supersession", code: "invalid-argument", input: map[string]any{"concern_id": "CON-run", "action": "superseded", "issue_id": nil, "supersedes_concern_id": "CON-run", "supersedes_finding_id": nil, "reason": "self"}},
		{name: "outside-action", code: "invalid-argument", input: map[string]any{"concern_id": "CON-run", "action": "addressed", "issue_id": nil, "supersedes_concern_id": nil, "supersedes_finding_id": nil, "reason": "old action"}},
	}
	for _, test := range cases {
		before := p21JournalEventCount(t, a.workspace)
		response := p23ConcernMutation(t, a, "dispose", "reject-"+test.name, test.input)
		if response.ExitCode == 0 || response.Stdout != "" || !strings.Contains(response.Stderr, `"code":"`+test.code+`"`) {
			t.Fatalf("%s response = %#v", test.name, response)
		}
		if after := p21JournalEventCount(t, a.workspace); after != before {
			t.Fatalf("%s appended Event(s): %d -> %d", test.name, before, after)
		}
	}
}

func recordDeviationChain(t *testing.T, a *p22FIDRecoveryAdapter) []map[string]any {
	t.Helper()
	inputs := []map[string]any{
		{"id": "DEV-prior", "subject": map[string]any{"type": "requirement", "id": "REQ-unresolved"}, "expected": "accepted contract", "observed": "first departure", "authority": map[string]any{"actor_id": "human-reviewer", "reference": nil}, "disposition": "pending-decision", "reasoning": "awaiting judgment", "references": []any{}, "supersedes_deviation_id": nil},
		{"id": "DEV-next", "subject": map[string]any{"type": "requirement", "id": "REQ-unresolved"}, "expected": "accepted contract", "observed": "approved departure", "authority": map[string]any{"actor_id": "human-reviewer", "reference": "DEC-8"}, "disposition": "accepted", "reasoning": "judgment recorded", "references": []any{}, "supersedes_deviation_id": "DEV-prior"},
	}
	for index, input := range inputs {
		response := p23RunInput(t, a, "deviation", "record", a.actor, "deviation-"+input["id"].(string), input)
		p23AssertReviewSuccess(t, "record Deviation", response)
		p23AssertActivity(t, a.workspace, response, "deviation-recorded", a.issueID, 0)
		value := p21Path(t, p23Envelope(t, response.Stdout), "result", "deviation").(map[string]any)
		if index == 1 && value["supersedes_deviation_id"] != "DEV-prior" {
			t.Fatalf("superseding Deviation = %#v", value)
		}
		p23AssertReplayWithoutEvent(t, a, "record Deviation", func() p22ConformanceResponse {
			return p23RunInput(t, a, "deviation", "record", a.actor, "deviation-"+input["id"].(string), input)
		})
	}
	values := p23DeviationList(t, a)
	if len(values) != 2 || values[0]["id"] != "DEV-prior" || values[1]["id"] != "DEV-next" || values[1]["supersedes_deviation_id"] != "DEV-prior" {
		t.Fatalf("readable Deviation chain = %#v", values)
	}
	return values
}

func p23AssertConcernHasNoCompletionColumn(t *testing.T, workspace string) {
	t.Helper()
	conn, err := fiddb.Open(filepath.Join(workspace, config.WorkspaceDirName, config.DatabaseFileName))
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = conn.Close() }()
	var found int
	err = conn.QueryRow(`SELECT COUNT(*) FROM pragma_table_info('judgment_concerns') WHERE name='completion_id'`).Scan(&found)
	if err != nil {
		t.Fatal(err)
	}
	if found != 0 {
		t.Fatal("judgment_concerns retains a completion_id column")
	}
}

func negativeDeviationPredecessors(t *testing.T, a *p22FIDRecoveryAdapter) {
	t.Helper()
	cases := []struct {
		name, code, id, predecessor string
		exit                        int
	}{
		{name: "missing", code: "deviation-not-found", id: "DEV-missing-child", predecessor: "DEV-missing", exit: 3},
		{name: "self", code: "invalid-argument", id: "DEV-self", predecessor: "DEV-self", exit: 2},
	}
	for _, test := range cases {
		input := map[string]any{"id": test.id, "subject": map[string]any{"type": "requirement", "id": "REQ-unresolved"}, "expected": "accepted contract", "observed": "departure", "authority": map[string]any{"actor_id": "human-reviewer", "reference": nil}, "disposition": "pending-decision", "reasoning": "invalid predecessor", "references": []any{}, "supersedes_deviation_id": test.predecessor}
		before := p21JournalEventCount(t, a.workspace)
		response := p23RunInput(t, a, "deviation", "record", a.actor, "deviation-reject-"+test.name, input)
		assertP22ProcessFailure(t, "deviation.record", test.code, test.exit, response)
		if after := p21JournalEventCount(t, a.workspace); after != before {
			t.Fatalf("%s Deviation predecessor rejection appended Event(s): %d -> %d", test.name, before, after)
		}
	}
}

func assertConcernDeviationEvidenceRejected(t *testing.T, a *p22FIDRecoveryAdapter, concernID, deviationID string) {
	t.Helper()
	for _, test := range []struct{ evidenceType, subjectID string }{{"concern", concernID}, {"deviation", deviationID}} {
		before := p21JournalEventCount(t, a.workspace)
		response := p23RunInput(t, a, "evidence", "link", a.actor, "evidence-reject-"+test.evidenceType, map[string]any{"obligation_id": "OBL-p22-corpus", "evidence_type": test.evidenceType, "receipt_id": test.subjectID, "review_id": nil, "human_outcome_id": nil, "target_repository": a.repository, "required_artifact_ids": []string{}})
		assertP22ProcessFailure(t, "evidence.link", "invalid-argument", 2, response)
		if p21JournalEventCount(t, a.workspace) != before {
			t.Fatalf("%s Evidence rejection appended an Event", test.evidenceType)
		}
	}
}

func assertConcernDisposal(t *testing.T, a *p22FIDRecoveryAdapter, concernID, key string, input map[string]any, action string) map[string]any {
	t.Helper()
	response := p23ConcernMutation(t, a, "dispose", key, input)
	p23AssertReviewSuccess(t, "dispose Concern", response)
	p23AssertActivity(t, a.workspace, response, "concern-disposed", a.issueID, 1)
	value := p21Path(t, p23Envelope(t, response.Stdout), "result", "concern").(map[string]any)
	if value["id"] != concernID || value["status"] != "disposed" || value["disposition"].(map[string]any)["action"] != action {
		t.Fatalf("disposed Concern = %#v", value)
	}
	p23AssertReplayWithoutEvent(t, a, "dispose Concern", func() p22ConformanceResponse {
		return p23ConcernMutation(t, a, "dispose", key, input)
	})
	return value
}

func p23AssertReplayWithoutEvent(t *testing.T, a *p22FIDRecoveryAdapter, operation string, run func() p22ConformanceResponse) {
	t.Helper()
	before := p21JournalEventCount(t, a.workspace)
	response := run()
	p23AssertReviewSuccess(t, operation+" replay", response)
	if replayed := p21Path(t, p23Envelope(t, response.Stdout), "result", "replayed"); replayed != true {
		t.Fatalf("%s replayed = %#v, want true", operation, replayed)
	}
	if after := p21JournalEventCount(t, a.workspace); after != before {
		t.Fatalf("%s replay appended Event(s): %d -> %d", operation, before, after)
	}
}

func p23ConcernMutation(t *testing.T, a *p22FIDRecoveryAdapter, command, key string, input map[string]any) p22ConformanceResponse {
	t.Helper()
	return p23RunInput(t, a, "concern", command, a.actor, key, input)
}

func p23DeviationList(t *testing.T, a *p22FIDRecoveryAdapter) []map[string]any {
	t.Helper()
	response := a.p22Run(t, "deviation", "list", "--subject", "requirement:REQ-unresolved")
	p23AssertReviewSuccess(t, "list Deviations", response)
	items := p21Path(t, p23Envelope(t, response.Stdout), "result", "items").([]any)
	result := make([]map[string]any, len(items))
	for index := range items {
		result[index] = items[index].(map[string]any)
	}
	return result
}
