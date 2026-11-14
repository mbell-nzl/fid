//go:build integration

package integration

import (
	"encoding/json"
	"os"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"testing"
)

func TestP24All51KeyedLeavesRebuildInOneJournal(t *testing.T) {
	t.Parallel()
	binary := integrationFIDBinary(t)
	workspace := t.TempDir()
	actions := append(append([]string(nil), p22FinalActions...),
		"concern.record", "concern.dispose", "deviation.record", "human-outcome.record",
		"decision.record-reference",
		"review.record", "finding.address")
	a := p22NewFIDRecoveryAdapterAt(t, binary, workspace, "human", 3600, []string{"artifact.json"}, actions)
	p24SingleFlagsAndLinks(t, a)
	p24SingleRecovery(t, a)
	claim := p23AcquireClaim(t, a, a.actor, "mutation", "p24-single-judgment-claim")
	run := p23StartRun(t, a, a.actor, claim, "p24-single-judgment-run")
	a.claimID, a.runID = p21PathString(t, claim, "id"), p21PathString(t, run, "id")
	a.claimVersion, a.runVersion = p23Number(t, claim["version"]), p23Number(t, run["version"])
	p24SingleGateAndArtifact(t, a)
	p24SingleTransitions(t, a)
	_ = p23GenerateLeafPairsForAdapter(t, a)
	p24SingleRenew(t, a)
	p24SingleCompletion(t, a)
	p24SingleActivationRecovery(t, a)
	p24SingleWorkflow(t, a)
	p24SingleLifecycle(t, a)
	p24AssertSingleJournalInventory(t, workspace)
	p24AssertRebuildMapParity(t, workspace, binary)
}

func p24SingleFlagsAndLinks(t *testing.T, a *p22FIDRecoveryAdapter) {
	t.Helper()
	call := func(args ...string) map[string]any {
		t.Helper()
		return p21Run(t, a.binary, a.workspace, append([]string{"--project", a.workspace, "--json"}, args...)...)
	}
	call("label", "add", "p24-label", "--colour", "red", "--actor", "agent:developer", "--idempotency-key", "p24-single-label")
	call("issue", "transition", a.issueID, "--to", "building", "--expected-revision", "1", "--actor", "agent:"+a.actor, "--idempotency-key", "p24-single-transition")
	call("issue", "ac", "add", a.issueID, "--text", "first criterion", "--ordinal", "1", "--actor", "agent:developer", "--idempotency-key", "p24-single-ac")
	call("issue", "update", a.issueID, "--title", "Phase 24 updated", "--expected-revision", "3", "--actor", "agent:developer", "--idempotency-key", "p24-single-update")
	call("issue", "label", "add", a.issueID, "p24-label", "--actor", "agent:developer", "--idempotency-key", "p24-single-label-add")
	call("issue", "label", "remove", a.issueID, "p24-label", "--actor", "agent:developer", "--idempotency-key", "p24-single-label-remove")
	other := call("issue", "create", "--kind", "task", "--title", "Phase 24 link peer", "--actor", "agent:developer", "--idempotency-key", "p24-single-peer")
	peer := p21PathString(t, other, "result", "issue", "id")
	call("link", "add", a.issueID, "blocks", peer, "--from-type", "issue", "--to-type", "issue", "--actor", "agent:developer", "--idempotency-key", "p24-single-link-add")
	call("link", "remove", a.issueID, "blocks", peer, "--from-type", "issue", "--to-type", "issue", "--actor", "agent:developer", "--idempotency-key", "p24-single-link-remove")
}

func p24SingleGateAndArtifact(t *testing.T, a *p22FIDRecoveryAdapter) {
	t.Helper()
	p22Write(t, a.workspace, "artifact.json", "phase 24 artifact\n")
	p22AcceptedCall(t, a, "p24-single-budget", "run.observe-budget", map[string]any{"value": 0})
	gate := p23RunInput(t, a, "gate", "run", a.actor, "p24-single-gate", map[string]any{"run_id": a.runID, "registration_id": a.registrationID, "registration_version": a.registrationVersion, "obligation_id": "OBL-p22-corpus", "expected_repository": map[string]any{"head_oid": a.repository["head_oid"], "dirty": a.repository["dirty"]}})
	if gate.ExitCode != 0 || gate.Stderr != "" {
		t.Fatalf("single-journal gate.run = %d stderr=%q", gate.ExitCode, gate.Stderr)
	}
	response := p24GateTerminal(t, gate.Stdout)["response"].(map[string]any)
	receipt := p21Path(t, response, "result", "receipt").(map[string]any)
	artifacts := receipt["artifacts"].([]any)
	if len(artifacts) != 1 {
		t.Fatalf("single-journal gate artifacts = %#v", artifacts)
	}
	artifact := artifacts[0].(map[string]any)
	p24SingleInput(t, a, "artifact", "verify", "p24-single-artifact", map[string]any{"artifact_id": artifact["id"]})
	p24SingleInput(t, a, "evidence", "link", "p24-single-evidence", map[string]any{"obligation_id": "OBL-p22-corpus", "evidence_type": "receipt", "receipt_id": receipt["id"], "review_id": nil, "human_outcome_id": nil, "required_artifact_ids": []string{}})
}

func p24SingleRenew(t *testing.T, a *p22FIDRecoveryAdapter) {
	t.Helper()
	claim := p23AcquireClaim(t, a, a.actor, "mutation", "p24-single-renew-claim")
	id := p21PathString(t, claim, "id")
	p24SingleInput(t, a, "claim", "renew", "p24-single-renew", map[string]any{"claim_id": id, "expected_version": 1, "ttl_seconds": 3600})
	p24SingleInput(t, a, "claim", "release", "p24-single-renew-release", map[string]any{"claim_id": id, "expected_version": 2, "reason": "completed"})
}

func p24SingleRecovery(t *testing.T, a *p22FIDRecoveryAdapter) {
	t.Helper()
	for _, fixture := range p22RecoveryFixtures(t) {
		if fixture.Kind != "accepted" || fixture.Reset || fixture.Name == "gate-committed-denial" {
			continue
		}
		if fixture.Operation == "run.close" {
			continue
		}
		response := a.Execute(t, fixture)
		if response.ExitCode != fixture.ExitCode {
			t.Fatalf("single-journal %s = %d stderr=%q, want %d", fixture.Name, response.ExitCode, response.Stderr, fixture.ExitCode)
		}
	}
}

func p24SingleTransitions(t *testing.T, a *p22FIDRecoveryAdapter) {
	t.Helper()
	issue := p21Run(t, a.binary, a.workspace, "--project", a.workspace, "--json", "issue", "get", a.issueID)
	revision := int(p21Path(t, issue, "result", "issue", "revision").(float64))
	for _, state := range []string{"dev_complete", "qa_review", "qa_pass"} {
		p21Run(t, a.binary, a.workspace, "--project", a.workspace, "--json", "issue", "transition", a.issueID, "--to", state, "--actor", "agent:"+a.actor, "--expected-revision", strconv.Itoa(revision), "--idempotency-key", "p24-single-issue-"+state)
		revision++
	}
}

func p24SingleCompletion(t *testing.T, a *p22FIDRecoveryAdapter) {
	t.Helper()
	p24SingleInput(t, a, "envelope", "complete", "p24-single-complete", map[string]any{"envelope_id": p23EnvelopeID(t, a), "expected_revision": 1})
}

func p24SingleActivationRecovery(t *testing.T, a *p22FIDRecoveryAdapter) {
	t.Helper()
	prepare := func(key string) string {
		t.Helper()
		prior := p21ActivationStatus(t, a.workspace, a.activation["id"].(string))
		resources := make([]any, 0, len(prior.Resources))
		for _, resource := range prior.Resources {
			resources = append(resources, map[string]any{"component": resource.Component, "old_identity": resource.NewIdentity, "new_identity": resource.NewIdentity + "-p24-next"})
		}
		response := p23RunInput(t, a, "activation", "prepare", "bootstrap:corpus", key, map[string]any{
			"backend_identity": "fid", "direction": "activate", "kit_identity": key, "phase": 21,
			"expected_prior_activation_id": prior.ID, "expected_prior_version": prior.Version, "resources": resources,
		})
		if response.ExitCode != 0 || response.Stderr != "" {
			t.Fatalf("single-journal activation.prepare %s = %d stderr=%q", key, response.ExitCode, response.Stderr)
		}
		return p21PathString(t, p23Envelope(t, response.Stdout), "result", "activation", "id")
	}
	abortedID := prepare("p24-single-abort-prep")
	p24SingleActivationInput(t, a, "abort", "p24-single-abort", map[string]any{"activation_id": abortedID, "expected_status": "prepared", "expected_version": 1, "reason": "operator-abort"})
	recoveryID := prepare("p24-single-recover-prep")
	status := p21Run(t, a.binary, a.workspace, "--project", a.workspace, "--json", "activation", "status", recoveryID)
	adapter := *a
	adapter.actor = "bootstrap:corpus"
	adapter.activation = p21Path(t, status, "result", "activation").(map[string]any)
	adapter.requests = nil
	data, err := os.ReadFile("testdata/p22_conformance/activation-cases.json")
	if err != nil {
		t.Fatal(err)
	}
	var cases map[string][]p22ConformanceCase
	if err := json.Unmarshal(data, &cases); err != nil {
		t.Fatal(err)
	}
	for _, fixture := range cases["prepared"] {
		if fixture.Kind != "accepted" || fixture.Operation == "activation.status" {
			continue
		}
		response := adapter.Execute(t, fixture)
		if response.ExitCode != 0 || response.Stderr != "" {
			t.Fatalf("single-journal %s = %d stderr=%q", fixture.Name, response.ExitCode, response.Stderr)
		}
	}
}

func p24SingleActivationInput(t *testing.T, a *p22FIDRecoveryAdapter, leaf, key string, input map[string]any) {
	t.Helper()
	response := p23RunInput(t, a, "activation", leaf, "bootstrap:corpus", key, input)
	if response.ExitCode != 0 || response.Stderr != "" {
		t.Fatalf("single-journal activation.%s = %d stderr=%q", leaf, response.ExitCode, response.Stderr)
	}
}

func p24SingleWorkflow(t *testing.T, a *p22FIDRecoveryAdapter) {
	t.Helper()
	for _, args := range [][]string{
		{"workflow", "define", "--name", "task-p24-single", "--applies-to", "subtask", "--states", "todo,building,done,cancelled", "--categories", "unstarted,started,completed,cancelled", "--transitions", "todo->building,building->done", "--actor", "agent:developer", "--idempotency-key", "p24-single-workflow-define"},
		{"workflow", "assign", "--name", "task-p24-single", "--kind", "subtask", "--actor", "agent:developer", "--idempotency-key", "p24-single-workflow-assign"},
	} {
		response := p24WorkflowCall(t, a.workspace, a.binary, append([]string{"--project", a.workspace, "--json"}, args...))
		if response.ExitCode != 0 || response.Stderr != "" {
			t.Fatalf("single-journal %v = %d stderr=%q", args[:2], response.ExitCode, response.Stderr)
		}
	}
}

func p24SingleInput(t *testing.T, a *p22FIDRecoveryAdapter, family, leaf, key string, input map[string]any) map[string]any {
	t.Helper()
	response := p23RunInput(t, a, family, leaf, a.actor, key, input)
	if response.ExitCode != 0 || response.Stderr != "" {
		t.Fatalf("single-journal %s.%s = %d stderr=%q", family, leaf, response.ExitCode, response.Stderr)
	}
	var envelope map[string]any
	if err := json.Unmarshal([]byte(response.Stdout), &envelope); err != nil {
		t.Fatalf("decode %s.%s: %v", family, leaf, err)
	}
	return envelope
}

func p24AssertSingleJournalInventory(t *testing.T, workspace string) {
	t.Helper()
	data, err := os.ReadFile("testdata/p24_keyed_leaves.tsv")
	if err != nil {
		t.Fatal(err)
	}
	want := make([]string, 0, 49)
	for _, line := range strings.Split(strings.TrimSpace(string(data)), "\n") {
		name, _, ok := strings.Cut(line, "\t")
		if !ok {
			t.Fatalf("invalid keyed-leaf inventory row %q", line)
		}
		operation := strings.ReplaceAll(name, " ", ".")
		if name == "start" {
			operation = "task.start"
		}
		if name == "park" {
			operation = "hold.park"
		}
		if name == "answer" {
			operation = "task.answer"
		}
		if name == "takeover" {
			operation = "task.takeover"
		}
		want = append(want, operation)
	}
	conn := p22FaultDB(t, workspace)
	defer func() { p22FaultOK(t, conn.Close()) }()
	rows, err := conn.Query(`SELECT DISTINCT operation_scope FROM idempotency_results ORDER BY operation_scope`)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { p22FaultOK(t, rows.Close()) }()
	got := make([]string, 0, 49)
	for rows.Next() {
		var scope string
		if err := rows.Scan(&scope); err != nil {
			t.Fatal(err)
		}
		operation, _, ok := strings.Cut(scope, ":")
		if !ok {
			t.Fatalf("unqualified operation scope %q", scope)
		}
		got = append(got, operation)
	}
	if err := rows.Err(); err != nil {
		t.Fatal(err)
	}
	sort.Strings(got)
	sort.Strings(want)
	if len(want) != 53 || !reflect.DeepEqual(got, want) {
		t.Fatalf("single-journal keyed operation inventory = %v; want exact 53 = %v", got, want)
	}
}
