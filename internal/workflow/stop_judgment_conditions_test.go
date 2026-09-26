package workflow

import (
	"database/sql"
	"fmt"
	"sort"
	"testing"

	"github.com/mbell-nzl/fid/internal/journal"
)

// judgmentConditionFixture opens the seeded recovery workspace (ENV-recovery,
// ISS-recovery, RUN-prior, ACT-recovery) so a test can seed judgment facts and
// Stops directly. Each case uses its own class name,
// so cases never interact through supersession or pair binding.
func judgmentConditionFixture(t *testing.T) *sql.DB {
	t.Helper()
	options, _ := recoveryFixture(t)
	conn, err := openWorkflowDB(options.WorkspaceRoot)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { closeWorkflowDB(conn) })
	for _, statement := range []string{
		`INSERT INTO issues (id, kind, title, workflow_id, status) VALUES ('ISS-other', 'task', 'other scope', 'task-default', 'todo')`,
		`INSERT INTO kit_activations (id, kit_identity, backend_identity, phase, status, point_of_no_return, direction, created_at, updated_at, version) VALUES ('ACT-other', 'kit', 'backend', 22, 'superseded', 0, 'activate', '2026-09-26T00:00:00Z', '2026-09-26T00:00:00Z', 1)`,
	} {
		if _, err := conn.Exec(statement); err != nil {
			t.Fatalf("seed judgment fixture: %v", err)
		}
	}
	return conn
}

func seedHumanOutcome(t *testing.T, conn *sql.DB, id, class, result, scopeIssue, activation, eventID string) {
	t.Helper()
	if _, err := conn.Exec(`INSERT INTO human_outcomes (id, actor_id, activation_id, activation_version, profile_ref, class, subject_type, subject_id, result, scope_issue_id, reasoning, references_json, created_at, created_event_id) VALUES (?, 'actor-human', ?, 1, 'profile', ?, 'issue', ?, ?, ?, 'reasoning', '[]', '2026-09-26T00:00:00Z', ?)`, id, activation, class, scopeIssue, result, scopeIssue, eventID); err != nil {
		t.Fatalf("seed Human Outcome %s: %v", id, err)
	}
}

func seedDecisionReference(t *testing.T, conn *sql.DB, id, class, outcome, scopeIssue, capabilities, eventID string) {
	t.Helper()
	if _, err := conn.Exec(`INSERT INTO decision_refs (id, class, source_json, authority_actor_id, authority_capabilities_json, outcome, scope_issue_id, created_at, created_event_id) VALUES (?, ?, '{}', 'actor-authority', ?, ?, ?, '2026-09-26T00:00:00Z', ?)`, id, class, capabilities, outcome, scopeIssue, eventID); err != nil {
		t.Fatalf("seed Decision Reference %s: %v", id, err)
	}
}

func declareProtectedDecision(t *testing.T, conn *sql.DB, class string, ordinal int) {
	t.Helper()
	if _, err := conn.Exec(`INSERT INTO envelope_protected_decisions (envelope_id, revision, ordinal, decision_class, required_capabilities_json, affected_actions_json) VALUES ('ENV-recovery', 1, ?, ?, '["decision.record-reference"]', '["envelope.complete"]')`, ordinal, class); err != nil {
		t.Fatalf("declare protected Decision %s: %v", class, err)
	}
}

func outcomeFact(id, class, result, scopeIssue, activation, eventID string) JudgmentFact {
	return JudgmentFact{EventType: "human-outcome.recorded", ID: id, Class: class, Result: result, ScopeIssueID: scopeIssue, ActivationID: activation, EventID: eventID}
}

func referenceFact(id, class, outcome, scopeIssue string, capabilities []string, eventID string) JudgmentFact {
	return JudgmentFact{EventType: "decision.referenced", ID: id, Class: class, Result: outcome, ScopeIssueID: scopeIssue, RequiredCapabilities: capabilities, EventID: eventID}
}

func humanOutcomeCondition(subject Subject, class string) Condition {
	return Condition{Type: "human-outcome-recorded", Subject: subject, OutcomeClass: &class}
}

func decisionCondition(subject Subject, class string) Condition {
	return Condition{Type: "decision-satisfied", Subject: subject, DecisionClass: &class}
}

// TestHumanOutcomeConditionQualifyingFactRules is DEC-250 R2: an Outcome
// qualifies only when it is a human-outcome.recorded fact, of the Condition's
// class, approved, not superseded, and inside the subject's scope.
func TestHumanOutcomeConditionQualifyingFactRules(t *testing.T) {
	conn := judgmentConditionFixture(t)
	envelope := Subject{Type: "envelope", ID: "ENV-recovery"}
	run := Subject{Type: "run", ID: "RUN-prior"}
	activation := Subject{Type: "activation", ID: "ACT-recovery"}
	cases := []struct {
		name    string
		subject Subject
		mutate  func(fact *JudgmentFact)
		seed    func(class string)
		want    bool
	}{
		{name: "envelope subject with selected scope Issue qualifies", subject: envelope, want: true},
		{name: "run subject with the Run's Issue qualifies", subject: run, want: true},
		{name: "activation subject with the provenance Activation qualifies", subject: activation, want: true},
		{name: "wrong class rejects", subject: envelope, mutate: func(f *JudgmentFact) { f.Class = "another-class" }},
		{name: "rejected result rejects", subject: envelope, mutate: func(f *JudgmentFact) { f.Result = "rejected" }},
		{name: "changes-required result rejects", subject: envelope, mutate: func(f *JudgmentFact) { f.Result = "changes-required" }},
		{name: "acknowledged result rejects", subject: envelope, mutate: func(f *JudgmentFact) { f.Result = "acknowledged" }},
		{name: "decision reference Event rejects", subject: envelope, mutate: func(f *JudgmentFact) { f.EventType = "decision.referenced" }},
		{name: "envelope subject with unselected scope Issue rejects", subject: envelope, mutate: func(f *JudgmentFact) { f.ScopeIssueID = "ISS-other" }},
		{name: "run subject with another Issue rejects", subject: run, mutate: func(f *JudgmentFact) { f.ScopeIssueID = "ISS-other" }},
		{name: "activation subject with another Activation rejects", subject: activation, mutate: func(f *JudgmentFact) { f.ActivationID = "ACT-other" }},
		{name: "superseded by a later different-result Outcome rejects", subject: envelope, want: false, seed: func(class string) {
			seedHumanOutcome(t, conn, "OUT-later-"+class, class, "rejected", "ISS-recovery", "ACT-recovery", "EVT-"+class+"-9")
		}},
		{name: "a later same-result Outcome does not supersede", subject: envelope, want: true, seed: func(class string) {
			seedHumanOutcome(t, conn, "OUT-later-"+class, class, "approved", "ISS-recovery", "ACT-recovery", "EVT-"+class+"-9")
		}},
	}
	for index, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			class := fmt.Sprintf("class-%d", index)
			fact := outcomeFact("OUT-"+class, class, "approved", "ISS-recovery", "ACT-recovery", "EVT-"+class+"-1")
			seedHumanOutcome(t, conn, fact.ID, class, "approved", "ISS-recovery", "ACT-recovery", "EVT-"+class+"-1")
			if tc.mutate != nil {
				tc.mutate(&fact)
			}
			if tc.seed != nil {
				tc.seed(class)
			}
			condition := humanOutcomeCondition(tc.subject, class)
			err := qualifyingHumanOutcomeCondition(conn, condition, fact)
			if (err == nil) != tc.want {
				t.Fatalf("qualifyingHumanOutcomeCondition() = %v, want qualifies=%v", err, tc.want)
			}
		})
	}
}

// TestDecisionSatisfiedConditionQualifyingFactRules is DEC-250 R3: only the
// pair-completing Event of a protected Decision that is satisfied under section
// 9 qualifies, an Envelope or a Run's Envelope is the scope, and an activation
// subject never qualifies.
func TestDecisionSatisfiedConditionQualifyingFactRules(t *testing.T) {
	conn := judgmentConditionFixture(t)
	envelope := Subject{Type: "envelope", ID: "ENV-recovery"}
	run := Subject{Type: "run", ID: "RUN-prior"}
	activation := Subject{Type: "activation", ID: "ACT-recovery"}
	authority := `["decision.record-reference"]`
	type pair struct {
		refOutcome, refCapabilities, outcomeResult string
		outcomeFirst                               bool
	}
	cases := []struct {
		name    string
		subject Subject
		pair    pair
		declare bool
		fact    string // "reference", "outcome"
		want    bool
	}{
		{name: "reference completing the pair qualifies", subject: envelope, pair: pair{"approved", authority, "approved", true}, declare: true, fact: "reference", want: true},
		{name: "outcome completing the pair qualifies", subject: envelope, pair: pair{"approved", authority, "approved", false}, declare: true, fact: "outcome", want: true},
		{name: "run subject resolves to the Run's Envelope", subject: run, pair: pair{"approved", authority, "approved", false}, declare: true, fact: "outcome", want: true},
		{name: "earlier half of the pair (reference) rejects", subject: envelope, pair: pair{"approved", authority, "approved", false}, declare: true, fact: "reference"},
		{name: "earlier half of the pair (outcome) rejects", subject: envelope, pair: pair{"approved", authority, "approved", true}, declare: true, fact: "outcome"},
		{name: "rejected reference never satisfies", subject: envelope, pair: pair{"rejected", authority, "approved", false}, declare: true, fact: "outcome"},
		{name: "deferred reference never satisfies", subject: envelope, pair: pair{"deferred", authority, "approved", false}, declare: true, fact: "reference"},
		{name: "reference without the declared authority never satisfies", subject: envelope, pair: pair{"approved", `["stop.resolve"]`, "approved", false}, declare: true, fact: "outcome"},
		{name: "rejected outcome never satisfies", subject: envelope, pair: pair{"approved", authority, "rejected", false}, declare: true, fact: "outcome"},
		{name: "class the Envelope does not declare rejects", subject: envelope, pair: pair{"approved", authority, "approved", false}, declare: false, fact: "outcome"},
		{name: "activation subject has no Envelope and rejects", subject: activation, pair: pair{"approved", authority, "approved", false}, declare: true, fact: "outcome"},
	}
	for index, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			class := fmt.Sprintf("decision-%d", index)
			if tc.declare {
				declareProtectedDecision(t, conn, class, index)
			}
			refEvent, outcomeEvent := "EVT-"+class+"-1", "EVT-"+class+"-2"
			if tc.pair.outcomeFirst {
				refEvent, outcomeEvent = "EVT-"+class+"-2", "EVT-"+class+"-1"
			}
			seedDecisionReference(t, conn, "DREF-"+class, class, tc.pair.refOutcome, "ISS-recovery", tc.pair.refCapabilities, refEvent)
			seedHumanOutcome(t, conn, "OUT-"+class, class, tc.pair.outcomeResult, "ISS-recovery", "ACT-recovery", outcomeEvent)
			fact := outcomeFact("OUT-"+class, class, tc.pair.outcomeResult, "ISS-recovery", "ACT-recovery", outcomeEvent)
			if tc.fact == "reference" {
				fact = referenceFact("DREF-"+class, class, tc.pair.refOutcome, "ISS-recovery", []string{"decision.record-reference"}, refEvent)
			}
			err := qualifyingDecisionCondition(conn, decisionCondition(tc.subject, class), fact)
			if (err == nil) != tc.want {
				t.Fatalf("qualifyingDecisionCondition() = %v, want qualifies=%v", err, tc.want)
			}
		})
	}
}

// TestManualAndAtomicResolutionShareOnePredicate is DEC-250 R1: the atomic
// resolution judges the Event that is about to commit (a pending fact) and the
// manual `stop resolve` judges the committed Event, through the same function.
// For every fixture below both judgments must agree.
func TestManualAndAtomicResolutionShareOnePredicate(t *testing.T) {
	conn := judgmentConditionFixture(t)
	envelope := Subject{Type: "envelope", ID: "ENV-recovery"}
	run := Subject{Type: "run", ID: "RUN-prior"}
	activation := Subject{Type: "activation", ID: "ACT-recovery"}
	for index, tc := range []struct {
		name    string
		subject Subject
		result  string
		class   string // when non-empty, the fact's class differs from the Condition's
		scope   string
		want    bool
	}{
		{name: "envelope approved", subject: envelope, result: "approved", scope: "ISS-recovery", want: true},
		{name: "run approved", subject: run, result: "approved", scope: "ISS-recovery", want: true},
		{name: "activation approved", subject: activation, result: "approved", scope: "ISS-recovery", want: true},
		{name: "envelope rejected", subject: envelope, result: "rejected", scope: "ISS-recovery"},
		{name: "envelope wrong class", subject: envelope, result: "approved", class: "elsewhere", scope: "ISS-recovery"},
		{name: "envelope out of scope", subject: envelope, result: "approved", scope: "ISS-other"},
		{name: "run out of scope", subject: run, result: "approved", scope: "ISS-other"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			class := fmt.Sprintf("equivalence-%d", index)
			factClass := class
			if tc.class != "" {
				factClass = tc.class
			}
			condition := humanOutcomeCondition(tc.subject, class)
			pending := PendingHumanOutcomeFact("OUT-"+class, factClass, tc.result, tc.scope, "ACT-recovery")
			atomic := qualifyingHumanOutcomeCondition(conn, condition, pending)
			seedHumanOutcome(t, conn, "OUT-"+class, factClass, tc.result, tc.scope, "ACT-recovery", "EVT-"+class)
			manual := qualifyingHumanOutcomeCondition(conn, condition, outcomeFact("OUT-"+class, factClass, tc.result, tc.scope, "ACT-recovery", "EVT-"+class))
			if (atomic == nil) != tc.want || (manual == nil) != tc.want {
				t.Fatalf("pending (atomic) = %v, committed (manual) = %v, want both qualifies=%v", atomic, manual, tc.want)
			}
		})
	}

	// The same holds for the protected-Decision pair: pending before the
	// completing row is inserted, committed after.
	class := "equivalence-pair"
	declareProtectedDecision(t, conn, class, 90)
	condition := decisionCondition(envelope, class)
	capabilities := []string{"decision.record-reference"}
	seedHumanOutcome(t, conn, "OUT-"+class, class, "approved", "ISS-recovery", "ACT-recovery", "EVT-3100")
	if err := qualifyingDecisionCondition(conn, condition, outcomeFact("OUT-"+class, class, "approved", "ISS-recovery", "ACT-recovery", "EVT-3100")); err == nil {
		t.Fatal("outcome alone qualified before its pair was complete")
	}
	pendingReference := PendingDecisionReferenceFact("DREF-"+class, class, "approved", "ISS-recovery", capabilities)
	atomic := qualifyingDecisionCondition(conn, condition, pendingReference)
	seedDecisionReference(t, conn, "DREF-"+class, class, "approved", "ISS-recovery", `["decision.record-reference"]`, "EVT-3101")
	manual := qualifyingDecisionCondition(conn, condition, referenceFact("DREF-"+class, class, "approved", "ISS-recovery", capabilities, "EVT-3101"))
	if atomic != nil || manual != nil {
		t.Fatalf("pair-completing reference: pending (atomic) = %v, committed (manual) = %v, want both to qualify", atomic, manual)
	}
	if err := qualifyingDecisionCondition(conn, condition, outcomeFact("OUT-"+class, class, "approved", "ISS-recovery", "ACT-recovery", "EVT-3100")); err == nil {
		t.Fatal("earlier-half outcome qualified after the reference completed the pair")
	}
}

// TestQualifyingJudgmentStopsKeysOnConditionNotRaisingCode is DEC-250 R1: a
// producer's readiness Stop and an operator-stop carrying the same Condition
// are resolved by the same fact, whatever their subject type, while a Stop
// carrying any other Condition, another class, or an already-resolved Stop is
// never touched.
func TestQualifyingJudgmentStopsKeysOnConditionNotRaisingCode(t *testing.T) {
	conn := judgmentConditionFixture(t)
	envelope := Subject{Type: "envelope", ID: "ENV-recovery"}
	run := Subject{Type: "run", ID: "RUN-prior"}
	activation := Subject{Type: "activation", ID: "ACT-recovery"}
	seed := func(id, code string, subject Subject, condition Condition) {
		t.Helper()
		raw, err := journal.MarshalCanonicalJSON(condition)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := conn.Exec(`INSERT INTO stops (id, subject_type, subject_id, code, blocked_actions_json, details_json, resolution_condition_json, resolver_capability, actor_id, raised_at, created_event_id) VALUES (?, ?, ?, ?, '["run.recover"]', '{}', ?, 'stop.resolve', 'actor', '2026-09-26T00:00:00Z', 'EVT-raise')`, id, subject.Type, subject.ID, code, string(raw)); err != nil {
			t.Fatalf("seed Stop %s: %v", id, err)
		}
	}
	seed("STP-01-producer", "human-outcome-required", envelope, humanOutcomeCondition(envelope, "security"))
	seed("STP-02-operator-envelope", "operator-stop", envelope, humanOutcomeCondition(envelope, "security"))
	seed("STP-03-operator-run", "operator-stop", run, humanOutcomeCondition(run, "security"))
	seed("STP-04-operator-activation", "operator-stop", activation, humanOutcomeCondition(activation, "security"))
	seed("STP-05-other-class", "operator-stop", envelope, humanOutcomeCondition(envelope, "release"))
	seed("STP-06-other-condition", "operator-stop", envelope, Condition{Type: "review-recorded", Subject: envelope})
	seed("STP-07-resolved", "operator-stop", envelope, humanOutcomeCondition(envelope, "security"))
	seed("STP-08-decision-condition", "operator-stop", envelope, decisionCondition(envelope, "security"))
	if _, err := conn.Exec(`INSERT INTO stop_resolutions (stop_id, actor_id, fact_event_id, resolved_at, created_event_id) VALUES ('STP-07-resolved', 'actor', 'EVT-old', '2026-09-26T00:00:00Z', 'EVT-old-resolve')`); err != nil {
		t.Fatal(err)
	}

	stops, err := QualifyingJudgmentStops(conn, PendingHumanOutcomeFact("OUT-new", "security", "approved", "ISS-recovery", "ACT-recovery"))
	if err != nil {
		t.Fatalf("QualifyingJudgmentStops() error = %v", err)
	}
	got := []string{}
	for _, stop := range stops {
		got = append(got, stop.ID)
	}
	sort.Strings(got)
	want := []string{"STP-01-producer", "STP-02-operator-envelope", "STP-03-operator-run", "STP-04-operator-activation"}
	if fmt.Sprint(got) != fmt.Sprint(want) {
		t.Fatalf("Outcome resolved Stops %v, want exactly the open human-outcome-recorded Stops in scope %v", got, want)
	}

	stops, err = QualifyingJudgmentStops(conn, PendingHumanOutcomeFact("OUT-wrong", "security", "rejected", "ISS-recovery", "ACT-recovery"))
	if err != nil || len(stops) != 0 {
		t.Fatalf("rejected Outcome resolved Stops %v (%v), want none", stops, err)
	}
}

func TestJudgmentFactFromEventReadsTheCommittedEffectRow(t *testing.T) {
	outcome := journal.Event{EventID: "EVT-outcome", EventType: "human-outcome.recorded", Effects: []journal.Effect{{Op: "insert", Table: "human_outcomes", PK: map[string]any{"id": "OUT-1"}, Values: map[string]any{"class": "security", "result": "approved", "scope_issue_id": "ISS-1", "activation_id": "ACT-1"}}}}
	fact, err := judgmentFactFromEvent(outcome)
	if err != nil || fmt.Sprint(fact) != fmt.Sprint(outcomeFact("OUT-1", "security", "approved", "ISS-1", "ACT-1", "EVT-outcome")) {
		t.Fatalf("Human Outcome fact = %#v, %v", fact, err)
	}
	reference := journal.Event{EventID: "EVT-reference", EventType: "decision.referenced", Effects: []journal.Effect{{Op: "insert", Table: "decision_refs", PK: map[string]any{"id": "DREF-1"}, Values: map[string]any{"class": "release", "outcome": "approved", "scope_issue_id": "ISS-1", "authority_capabilities_json": `["decision.record-reference"]`}}}}
	fact, err = judgmentFactFromEvent(reference)
	if err != nil || fmt.Sprint(fact) != fmt.Sprint(referenceFact("DREF-1", "release", "approved", "ISS-1", []string{"decision.record-reference"}, "EVT-reference")) {
		t.Fatalf("Decision Reference fact = %#v, %v", fact, err)
	}
	if _, err := judgmentFactFromEvent(journal.Event{EventType: "human-outcome.recorded"}); err == nil {
		t.Fatal("Event without its judgment row produced a fact")
	}
}

func TestEventMatchesJudgmentConditionTypes(t *testing.T) {
	envelope := Subject{Type: "envelope", ID: "ENV-recovery"}
	for _, tc := range []struct {
		condition Condition
		eventType string
		want      bool
	}{
		{humanOutcomeCondition(envelope, "security"), "human-outcome.recorded", true},
		{humanOutcomeCondition(envelope, "security"), "decision.referenced", false},
		{humanOutcomeCondition(envelope, "security"), "evidence.linked", false},
		{decisionCondition(envelope, "security"), "human-outcome.recorded", true},
		{decisionCondition(envelope, "security"), "decision.referenced", true},
		{decisionCondition(envelope, "security"), "review.recorded", false},
	} {
		if got := eventMatchesCondition(journal.Event{EventType: tc.eventType}, tc.condition); got != tc.want {
			t.Fatalf("eventMatchesCondition(%s, %s) = %v, want %v", tc.condition.Type, tc.eventType, got, tc.want)
		}
	}
}

// TestStopRaiseRejectsDecisionSatisfiedOnActivationSubject is the DEC-250 R3
// raise rule: an activation subject has no Envelope scope, so such a Stop
// could never resolve.
func TestStopRaiseRejectsDecisionSatisfiedOnActivationSubject(t *testing.T) {
	for _, tc := range []struct {
		name    string
		subject Subject
		want    bool
	}{
		{"activation subject is rejected", Subject{Type: "activation", ID: "ACT-stop-test"}, false},
		{"envelope subject is accepted", Subject{Type: "envelope", ID: "ENV-1"}, true},
		{"run subject is accepted", Subject{Type: "run", ID: "RUN-1"}, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			request := stopTestRequest()
			request.Subject = tc.subject
			request.ResolutionCondition = decisionCondition(tc.subject, "release")
			err := validateStopRaise(request)
			if (err == nil) != tc.want {
				t.Fatalf("validateStopRaise() = %v, want valid=%v", err, tc.want)
			}
		})
	}
	// human-outcome-recorded keeps an activation scope (R2).
	request := stopTestRequest()
	request.ResolutionCondition = humanOutcomeCondition(request.Subject, "security")
	if err := validateStopRaise(request); err != nil {
		t.Fatalf("human-outcome-recorded on an activation subject rejected: %v", err)
	}
}
