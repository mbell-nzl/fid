package workflow

import (
	"database/sql"
	"encoding/json"
	"testing"

	_ "modernc.org/sqlite"
)

func TestProtectedDecisionBlocksOnlyDeclaredAffectedActions(t *testing.T) {
	t.Parallel()
	conn, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = conn.Close() }()
	for _, statement := range []string{
		`CREATE TABLE decision_refs (id TEXT PRIMARY KEY, class TEXT, authority_capabilities_json TEXT, outcome TEXT, scope_issue_id TEXT, created_event_id TEXT)`,
		`CREATE TABLE human_outcomes (id TEXT PRIMARY KEY, class TEXT, result TEXT, scope_issue_id TEXT, created_event_id TEXT)`,
	} {
		if _, err := conn.Exec(statement); err != nil {
			t.Fatal(err)
		}
	}

	envelope := Envelope{
		ID:              "ENV-protected",
		PrimaryIssueID:  "ISS-protected",
		CurrentRevision: 1,
		ProtectedDecisions: []ProtectedDecision{
			{ProtectedDecisionInput: ProtectedDecisionInput{Ordinal: 0, DecisionClass: "release", RequiredCapabilities: []string{"decision.approve"}, AffectedActions: []string{"envelope.complete"}}},
			{ProtectedDecisionInput: ProtectedDecisionInput{Ordinal: 1, DecisionClass: "quality", RequiredCapabilities: []string{"decision.approve"}, AffectedActions: []string{"issue.transition-terminal-success"}}},
		},
	}

	reasons, err := evaluateProtectedDecisions(conn, envelope, "envelope.complete")
	if err != nil {
		t.Fatal(err)
	}
	if len(reasons) != 1 || reasons[0].Code != "protected-decision-required" || reasons[0].SubjectID != "release" || reasons[0].SubjectType != "decision" {
		t.Fatalf("envelope.complete reasons = %#v", reasons)
	}
}

func TestProtectedDecisionSatisfactionRequiresBothApprovedFactsAndExactAuthority(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name, referenceOutcome, outcomeResult string
		referenceCapabilities                 string
		wantSatisfied                         bool
	}{
		{name: "neither"},
		{name: "reference-only", referenceOutcome: "approved", referenceCapabilities: `["decision.record-reference"]`},
		{name: "outcome-only", outcomeResult: "approved"},
		{name: "both-approved", referenceOutcome: "approved", referenceCapabilities: `["decision.record-reference"]`, outcomeResult: "approved", wantSatisfied: true},
		{name: "rejected-reference", referenceOutcome: "rejected", referenceCapabilities: `["decision.record-reference"]`, outcomeResult: "approved"},
		{name: "deferred-reference", referenceOutcome: "deferred", referenceCapabilities: `["decision.record-reference"]`, outcomeResult: "approved"},
		{name: "wrong-authority", referenceOutcome: "approved", referenceCapabilities: `["stop.resolve"]`, outcomeResult: "approved"},
		{name: "rejected-outcome", referenceOutcome: "approved", referenceCapabilities: `["decision.record-reference"]`, outcomeResult: "rejected"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			conn := protectedDecisionTestDB(t)
			defer func() { _ = conn.Close() }()
			if tc.referenceOutcome != "" {
				_, err := conn.Exec(`INSERT INTO decision_refs VALUES ('DREF-test','release',?,?,?,'EVT-reference')`, tc.referenceCapabilities, tc.referenceOutcome, "ISS-protected")
				if err != nil {
					t.Fatal(err)
				}
			}
			if tc.outcomeResult != "" {
				_, err := conn.Exec(`INSERT INTO human_outcomes VALUES ('OUT-test','release',?,?,'EVT-outcome')`, tc.outcomeResult, "ISS-protected")
				if err != nil {
					t.Fatal(err)
				}
			}
			decision := ProtectedDecision{ProtectedDecisionInput: ProtectedDecisionInput{DecisionClass: "release", RequiredCapabilities: []string{"decision.record-reference"}}}
			if err := resolveProtectedDecision(conn, "ISS-protected", &decision); err != nil {
				t.Fatal(err)
			}
			if got := decision.Status == "satisfied"; got != tc.wantSatisfied {
				t.Fatalf("status = %q, reference=%v outcome=%v", decision.Status, decision.DecisionRefID, decision.HumanOutcomeID)
			}
		})
	}
}

func TestDecisionConditionPreservesRetainedBytesAndSeparatesAttentionKeys(t *testing.T) {
	t.Parallel()
	subject := Subject{Type: "envelope", ID: "ENV-protected"}
	retained := Condition{Type: "work-ended", Subject: subject}
	raw, err := json.Marshal(retained)
	if err != nil {
		t.Fatal(err)
	}
	const wantRetained = `{"type":"work-ended","subject":{"type":"envelope","id":"ENV-protected"},"budget_name":null,"capability_id":null}`
	if string(raw) != wantRetained {
		t.Fatalf("retained Condition bytes = %s, want %s", raw, wantRetained)
	}
	classOne, classTwo := "release", "security"
	first := Condition{Type: "decision-satisfied", Subject: subject, DecisionClass: &classOne}
	second := Condition{Type: "decision-satisfied", Subject: subject, DecisionClass: &classTwo}
	firstKey, err := attentionSemanticKey(subject, "protected-decision-required", first)
	if err != nil {
		t.Fatal(err)
	}
	secondKey, err := attentionSemanticKey(subject, "protected-decision-required", second)
	if err != nil {
		t.Fatal(err)
	}
	if firstKey == secondKey {
		t.Fatalf("distinct Decision classes share Attention key %s", firstKey)
	}
	again, err := attentionSemanticKey(subject, "protected-decision-required", first)
	if err != nil || again != firstKey {
		t.Fatalf("Attention key changed: %q %v", again, err)
	}
}

func TestTask9RetainedAttentionKeyMatchesPrechangeReference(t *testing.T) {
	t.Parallel()
	const reference = "74c83b7033729785eb99403ecec9dcde30ac4280"
	name := "attempts"
	subject := Subject{Type: "run", ID: "RUN-pinned"}
	condition := Condition{Type: "budget-observed", Subject: subject, BudgetName: &name}
	key, err := attentionSemanticKey(subject, "budget-exhausted", condition)
	if err != nil {
		t.Fatal(err)
	}
	const prechangeKey = "sha256:d4b51bb4a6404d4a2579826ad5c08d802ed83c41814cd5fa74d3ee8c7796c3de"
	if reference == "" || key != prechangeKey {
		t.Fatalf("Attention key at %s = %s, want current key %s", reference, prechangeKey, key)
	}
}

func protectedDecisionTestDB(t *testing.T) *sql.DB {
	t.Helper()
	conn, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	statements := []string{
		`CREATE TABLE decision_refs (id TEXT PRIMARY KEY, class TEXT, authority_capabilities_json TEXT, outcome TEXT, scope_issue_id TEXT, created_event_id TEXT)`,
		`CREATE TABLE human_outcomes (id TEXT PRIMARY KEY, class TEXT, result TEXT, scope_issue_id TEXT, created_event_id TEXT)`,
	}
	for _, statement := range statements {
		if _, err := conn.Exec(statement); err != nil {
			t.Fatal(err)
		}
	}
	return conn
}
