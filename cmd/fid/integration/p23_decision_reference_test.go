//go:build integration

package integration

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/mbell-nzl/fid/internal/evidence"
)

func TestP23DecisionReferenceProtectedLifecycleRealProcess(t *testing.T) {
	t.Parallel()
	actions := []string{"claim.acquire", "claim.release", "decision.record-reference", "envelope.activate", "envelope.complete", "envelope.create", "evidence.link", "gate.register", "gate.run", "human-outcome.record", "issue.transition", "run.close", "run.observe", "run.observe-budget", "run.start", "stop.resolve"}
	decisions := []any{
		map[string]any{"ordinal": 0, "decision_class": "release", "required_capabilities": []string{"decision.record-reference"}, "affected_actions": []string{"envelope.complete"}},
		map[string]any{"ordinal": 1, "decision_class": "security", "required_capabilities": []string{"decision.record-reference"}, "affected_actions": []string{"issue.transition-terminal-success"}},
	}
	a := p22NewFIDRecoveryAdapterAt(t, integrationFIDBinary(t), t.TempDir(), "agent", 3600, []string{}, actions, p22ExecutionFixture{ProtectedDecisions: decisions})
	a.task13ReplayCorpus = true
	assertP23DecisionProjectionShape(t, a)
	assertP23DecisionStopStates(t, a, map[string]string{"release": "open", "security": "open"})
	initialResume := p21Run(t, a.binary, a.workspace, "--json", "resume", a.issueID, "--repository", a.workspace)
	resumeActions := p21Path(t, initialResume, "result", "resume", "next_actions").([]any)
	foundOutcomeAction := false
	for _, raw := range resumeActions {
		action := raw.(map[string]any)
		if action["operation"] == "human-outcome.record" {
			foundOutcomeAction = true
			if len(action["reasons"].([]any)) != 0 {
				t.Fatalf("activated Human Outcome action should be executable: %#v", action)
			}
		}
	}
	if !foundOutcomeAction {
		t.Fatal("resume omits Human Outcome recording action")
	}
	for revision, state := range []string{"in_progress", "dev_complete", "qa_review"} {
		response := p23IssueTransition(t, a, state, revision+1, "task9-"+state)
		assertP23DecisionCommandSuccess(t, "transition Issue to "+state, response)
	}
	beforeSecurityDenial := p21JournalEventCount(t, a.workspace)
	securityDenial := p23IssueTransition(t, a, "qa_pass", 4, "task9-security-denial")
	assertP23ProtectedDecisionDenial(t, a, beforeSecurityDenial, "issue.transition", securityDenial, "security", []string{"issue.transition-terminal-success"})
	p23AssertNoncommittingDenialRetry(t, a, securityDenial, func() p22ConformanceResponse {
		return p23IssueTransition(t, a, "qa_pass", 4, "task9-security-denial")
	})

	for _, outcome := range []string{"rejected", "deferred"} {
		response := p23RecordDecision(t, a, "release-"+outcome, "release", outcome, []string{"decision.record-reference"}, p23DecisionSource(a, true))
		assertP23DecisionSuccess(t, response, 0)
	}
	wrongAuthority := p23RecordDecision(t, a, "release-wrong-authority", "release", "approved", []string{"stop.resolve"}, p23DecisionSource(a, true))
	assertP23DecisionSuccess(t, wrongAuthority, 0)
	assertP23DecisionStopStates(t, a, map[string]string{"release": "open", "security": "open"})
	beforeRebuild := a.p22Run(t, "decision", "list", "--issue", a.issueID)
	if beforeRebuild.ExitCode != 0 || beforeRebuild.Stderr != "" {
		t.Fatalf("pre-rebuild decision list = %#v", beforeRebuild)
	}
	rebuilt := a.p22Run(t, "rebuild")
	if rebuilt.ExitCode != 0 || rebuilt.Stderr != "" {
		t.Fatalf("rebuild decision journal = %#v", rebuilt)
	}
	afterRebuild := a.p22Run(t, "decision", "list", "--issue", a.issueID)
	if afterRebuild.ExitCode != 0 || afterRebuild.Stderr != "" || afterRebuild.Stdout != beforeRebuild.Stdout {
		t.Fatalf("decision live/replay parity: before=%q after=%#v", beforeRebuild.Stdout, afterRebuild)
	}

	// Task 10 converges these fixtures onto the real Human Outcome writer,
	// superseding the seeded task 9 projection rows: recording provides the
	// same durable human_outcomes facts protected-Decision satisfaction reads.
	releaseOutcome := p23RecordHumanOutcome(t, a, "decision-release-approved", "actor-human-outcome", "release", "approved")
	assertP23HumanOutcomeSuccess(t, releaseOutcome)
	securityBeforeOutcome := p23RecordDecision(t, a, "security-reference-only", "security", "approved", []string{"decision.record-reference"}, p23DecisionSource(a, true))
	assertP23DecisionSuccess(t, securityBeforeOutcome, 0)
	assertP23DecisionStopStates(t, a, map[string]string{"release": "open", "security": "open"})
	// Recording the real Human Outcome now resolves the security Stop
	// atomically with this exact Event, because the approved Decision
	// Reference from securityBeforeOutcome is already durable: the later of
	// the two qualifying facts is the resolver (contracts.md section 9),
	// which the old seeded-row fixture could never observe since seeding
	// bypassed this resolution path entirely.
	securityOutcome := p23RecordHumanOutcome(t, a, "decision-security-approved", "actor-human-outcome", "security", "approved")
	assertP23DecisionSuccess(t, securityOutcome, 1)
	security := p23RecordDecision(t, a, "security-approved", "security", "approved", []string{"decision.record-reference"}, p23DecisionSource(a, true))
	assertP23DecisionSuccess(t, security, 0)
	assertP23DecisionStopStates(t, a, map[string]string{"release": "open", "security": "resolved"})

	p22AcceptedCall(t, a, "task9-budget", "run.observe-budget", map[string]any{"value": 0})
	gate, err := evidence.RunGate(context.Background(), p22BudgetGateOptions(t, a, "task9-proof"))
	if err != nil || gate.Receipt == nil {
		t.Fatalf("run qualifying gate: result=%#v error=%v", gate, err)
	}
	linked := p23RunInput(t, a, "evidence", "link", a.actor, "task9-proof-link", map[string]any{"obligation_id": "OBL-p22-corpus", "evidence_type": "receipt", "receipt_id": gate.Receipt.ID, "review_id": nil, "human_outcome_id": nil, "target_repository": a.repository, "required_artifact_ids": []string{}})
	assertP23DecisionCommandSuccess(t, "link qualifying Evidence", linked)
	assertP23DecisionCommandSuccess(t, "transition terminal with unrelated release Decision open", p23IssueTransition(t, a, "qa_pass", 4, "task9-qa-pass"))
	p23CloseAndRelease(t, a, a.actor, a.runID, a.claimID, "task9-close")
	beforeReleaseDenial := p21JournalEventCount(t, a.workspace)
	releaseDenialInput := map[string]any{"envelope_id": p23EnvelopeID(t, a), "expected_revision": 1, "repository": a.repository}
	releaseDenial := p23RunInput(t, a, "envelope", "complete", a.actor, "task9-release-denial", releaseDenialInput)
	assertP23ProtectedDecisionDenial(t, a, beforeReleaseDenial, "envelope.complete", releaseDenial, "release", []string{"envelope.complete"})
	p23AssertNoncommittingDenialRetry(t, a, releaseDenial, func() p22ConformanceResponse {
		return p23RunInput(t, a, "envelope", "complete", a.actor, "task9-release-denial", releaseDenialInput)
	})

	release := p23RecordDecision(t, a, "release-approved", "release", "approved", []string{"decision.record-reference"}, p23DecisionSource(a, true))
	assertP23DecisionSuccess(t, release, 1)
	assertP23DecisionStopStates(t, a, map[string]string{"release": "resolved", "security": "resolved"})
	assertP23DecisionCommandSuccess(t, "complete Envelope after both Decisions", p23RunInput(t, a, "envelope", "complete", a.actor, "task9-complete", map[string]any{"envelope_id": p23EnvelopeID(t, a), "expected_revision": 1, "repository": a.repository}))

	listed := a.p22Run(t, "decision", "list", "--issue", a.issueID)
	if listed.ExitCode != 0 || listed.Stderr != "" {
		t.Fatalf("decision list = %#v", listed)
	}
	var listEnvelope map[string]any
	if err := json.Unmarshal([]byte(listed.Stdout), &listEnvelope); err != nil {
		t.Fatal(err)
	}
	if items := p21Path(t, listEnvelope, "result", "items").([]any); len(items) != 6 {
		t.Fatalf("decision list count = %d, want 6", len(items))
	}

	for _, failure := range []struct {
		name   string
		source map[string]any
	}{
		{"digest-mismatch", p23DecisionSource(a, false)},
		{"source-missing", p23MissingDecisionSource(a)},
	} {
		before := p21JournalEventCount(t, a.workspace)
		response := p23RecordDecision(t, a, failure.name, "release", "approved", []string{"decision.record-reference"}, failure.source)
		assertP22ProcessFailure(t, "decision.record-reference", "source-digest-mismatch", 4, response)
		if response.ExitCode != 4 || response.Stdout != "" || !p23HasJSON(response.Stderr, `"code":"source-digest-mismatch"`) || p21JournalEventCount(t, a.workspace) != before {
			t.Fatalf("%s response = %#v; events %d -> %d", failure.name, response, before, p21JournalEventCount(t, a.workspace))
		}
	}
}

func p23IssueTransition(t *testing.T, a *p22FIDRecoveryAdapter, state string, revision int, key string) p22ConformanceResponse {
	t.Helper()
	stdout, stderr, err := closureRunCommand(a.workspace, a.binary, "--project", a.workspace, "--json", "issue", "transition", a.issueID, "--to", state, "--actor", "agent:"+a.actor, "--expected-revision", fmt.Sprint(revision), "--idempotency-key", key)
	response := p22ConformanceResponse{Stdout: stdout, Stderr: stderr}
	if err != nil {
		response.ExitCode = processExitCode(t, err)
	}
	return response
}

func assertP23DecisionCommandSuccess(t *testing.T, operation string, response p22ConformanceResponse) {
	t.Helper()
	if response.ExitCode != 0 || response.Stderr != "" {
		t.Fatalf("%s = %#v", operation, response)
	}
}

func assertP23ProtectedDecisionDenial(t *testing.T, a *p22FIDRecoveryAdapter, before int, operation string, response p22ConformanceResponse, class string, actions []string) {
	t.Helper()
	assertP22ProcessFailure(t, operation, "protected-decision-required", 5, response)
	if response.ExitCode != 5 || response.Stdout != "" {
		t.Fatalf("%s protected Decision denial = %#v", operation, response)
	}
	var envelope map[string]any
	if err := json.Unmarshal([]byte(response.Stderr), &envelope); err != nil {
		t.Fatalf("decode %s denial: %v; stderr=%q", operation, err, response.Stderr)
	}
	if got := p21PathString(t, envelope, "error", "code"); got != "protected-decision-required" {
		t.Fatalf("%s denial code = %q", operation, got)
	}
	details := p21Path(t, envelope, "error", "details").(map[string]any)
	want := map[string]any{"envelope_id": p23EnvelopeID(t, a), "decision_class": class, "required_capabilities": []any{"decision.record-reference"}, "affected_actions": stringSliceAsAny(actions), "stop_ids": details["stop_ids"], "attention_ids": details["attention_ids"], "event": nil, "replayed": false}
	if fmt.Sprint(details) != fmt.Sprint(want) {
		t.Fatalf("%s denial details = %v, want %v", operation, details, want)
	}
	p23AssertStopBlockedDetails(t, a, response, "protected-decision-required", "decision_class", class)
	if after := p21JournalEventCount(t, a.workspace); after != before {
		t.Fatalf("%s denial appended Event(s): %d -> %d", operation, before, after)
	}
}

func stringSliceAsAny(values []string) []any {
	items := make([]any, len(values))
	for index, value := range values {
		items[index] = value
	}
	return items
}

func p23RecordDecision(t *testing.T, a *p22FIDRecoveryAdapter, suffix, class, outcome string, capabilities []string, source map[string]any) p22ConformanceResponse {
	t.Helper()
	input := map[string]any{
		"id": "DREF-" + suffix, "class": class, "source": source,
		"authority": map[string]any{"actor_id": a.actor, "required_capabilities": capabilities},
		"outcome":   outcome, "scope": map[string]any{"issue_id": a.issueID},
	}
	return p23RunInput(t, a, "decision", "record-reference", a.actor, "p23-decision-"+suffix, input)
}

func p23DecisionSource(a *p22FIDRecoveryAdapter, matching bool) map[string]any {
	digest := sha256.Sum256([]byte("base\n"))
	encoded := hex.EncodeToString(digest[:])
	if !matching {
		encoded = p22Zeros(64)
	}
	return map[string]any{"source_type": "project-file", "display_label": "Task 9 decision", "path": "base.txt", "section": nil, "external_id": nil, "digest": "sha256:" + encoded, "repository": a.repository}
}

func p23MissingDecisionSource(a *p22FIDRecoveryAdapter) map[string]any {
	source := p23DecisionSource(a, true)
	source["path"] = "missing-decision.md"
	return source
}

func assertP23DecisionSuccess(t *testing.T, response p22ConformanceResponse, resolvedStops int) {
	t.Helper()
	if response.ExitCode != 0 || response.Stderr != "" {
		t.Fatalf("decision mutation = %#v", response)
	}
	var envelope map[string]any
	if err := json.Unmarshal([]byte(response.Stdout), &envelope); err != nil {
		t.Fatal(err)
	}
	stops := p21Path(t, envelope, "result", "stops").([]any)
	if len(stops) != resolvedStops {
		t.Fatalf("resolved Stops = %d, want %d; response=%s", len(stops), resolvedStops, response.Stdout)
	}
	attention := p21Path(t, envelope, "result", "attention").([]any)
	if len(attention) != resolvedStops {
		t.Fatalf("resolved Attention = %d, want %d; response=%s", len(attention), resolvedStops, response.Stdout)
	}
	if resolvedStops > 0 {
		eventID := p21PathString(t, envelope, "event", "id")
		stop := stops[0].(map[string]any)
		item := attention[0].(map[string]any)
		if p21PathString(t, stop, "resolution", "event_id") != eventID || item["resolution_event_id"] != eventID {
			t.Fatalf("selective resolution did not share Decision Event %s: stop=%#v attention=%#v", eventID, stop, item)
		}
	}
}

func assertP23DecisionStopStates(t *testing.T, a *p22FIDRecoveryAdapter, want map[string]string) {
	t.Helper()
	conn := p22FaultDB(t, a.workspace)
	defer func() { p22FaultOK(t, conn.Close()) }()
	rows, err := conn.Query(`SELECT json_extract(s.details_json,'$.decision_class'),CASE WHEN r.stop_id IS NULL THEN 'open' ELSE 'resolved' END,a.status FROM stops s LEFT JOIN stop_resolutions r ON r.stop_id=s.id JOIN attention_stop_links l ON l.stop_id=s.id JOIN attention_items a ON a.id=l.attention_id WHERE s.code='protected-decision-required' ORDER BY 1`)
	p22FaultOK(t, err)
	defer func() { p22FaultOK(t, rows.Close()) }()
	got := map[string]string{}
	for rows.Next() {
		var class, stopStatus, attentionStatus string
		p22FaultOK(t, rows.Scan(&class, &stopStatus, &attentionStatus))
		if stopStatus != attentionStatus {
			t.Fatalf("%s Stop=%s Attention=%s", class, stopStatus, attentionStatus)
		}
		got[class] = stopStatus
	}
	p22FaultOK(t, rows.Err())
	if fmt.Sprint(got) != fmt.Sprint(want) {
		t.Fatalf("protected Decision states = %v, want %v", got, want)
	}
}

func assertP23DecisionProjectionShape(t *testing.T, a *p22FIDRecoveryAdapter) {
	t.Helper()
	conn := p22FaultDB(t, a.workspace)
	defer func() { p22FaultOK(t, conn.Close()) }()
	rows, err := conn.Query(`PRAGMA table_info(decision_refs)`)
	p22FaultOK(t, err)
	defer func() { p22FaultOK(t, rows.Close()) }()
	for rows.Next() {
		var cid int
		var name, kind string
		var notNull, primaryKey int
		var defaultValue any
		p22FaultOK(t, rows.Scan(&cid, &name, &kind, &notNull, &defaultValue, &primaryKey))
		if name == "body" || name == "prose" || name == "content" {
			t.Fatalf("decision_refs stores forbidden body column %q", name)
		}
	}
	p22FaultOK(t, rows.Err())
}
