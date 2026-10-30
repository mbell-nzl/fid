//go:build integration

package integration

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestP22CompleteOneIssueRecoveryJourney(t *testing.T) {
	t.Parallel()
	p22CompleteOneIssueRecoveryJourney(t)
}

func p22CompleteOneIssueRecoveryJourney(t *testing.T) *p22FIDRecoveryAdapter {
	t.Helper()
	a := p22NewFIDRecoveryAdapterAt(t, integrationFIDBinary(t), t.TempDir(), "human", 3600, []string{}, p22FinalActions)
	prior := a.runID
	denial := assertP22LiveDenial(t, a.Execute(t, p22ConformanceCase{Name: "journey-denial", Operation: "gate.run", Kind: "denied"}), false)
	replayed := assertP22LiveDenial(t, a.Execute(t, p22ConformanceCase{Name: "journey-denial", Operation: "gate.run", Kind: "denied"}), true)
	if !reflect.DeepEqual(p22ResultPath(denial, "error.details.stop_ids"), p22ResultPath(replayed, "error.details.stop_ids")) {
		t.Fatal("denial replay changed Stop identities")
	}
	a.stopID = p22ResultPath(denial, "error.details.stop_ids.0").(string)
	a.attentionID = p22ResultPath(denial, "error.details.attention_ids.0").(string)
	a.attentionVersion = 1
	current := p21Run(t, a.binary, a.workspace, "--json", "run", "get", a.runID)
	a.runVersion = p22Number(p21Path(t, current, "result", "run", "version"))
	baseline := p22AcceptedCall(t, a, "journey-baseline", "run.observe-budget", map[string]any{"value": 0})
	a.terminalEvent = p21PathString(t, baseline, "event", "id")
	disposition := p22AcceptedCall(t, a, "journey-stop-disposition", "stop.resolve", nil)
	before := p21JournalEventCount(t, a.workspace)
	replayedDisposition := p22AcceptedCall(t, a, "journey-stop-disposition", "stop.resolve", nil)
	if p21PathString(t, disposition, "event", "id") != p21PathString(t, replayedDisposition, "event", "id") || p21JournalEventCount(t, a.workspace) != before {
		t.Fatal("embedded Stop disposition replay changed durable facts")
	}
	p22AcceptedCall(t, a, "journey-attention-disposition", "attention.resolve", nil)
	firstProof := p22JourneyProof(t, a, "journey-initial")
	p22Write(t, a.workspace, "base.txt", "adopted interrupted work\n")
	p22AcceptedCall(t, a, "journey-terminated", "run.observe", nil)
	p22AcceptedCall(t, a, "journey-inspect", "run.inspect", nil)
	p22AcceptedCall(t, a, "journey-checkpoint", "run.checkpoint", nil)
	p22AcceptedCall(t, a, "journey-inspect-checkpoint", "run.inspect", nil)
	p22AcceptedCall(t, a, "journey-expire", "claim.expire", nil)
	recovered := p22AcceptedCall(t, a, "journey-recover", "run.recover", map[string]any{"adopted_paths": []string{"base.txt"}, "ttl_seconds": 3600})
	if p21PathString(t, recovered, "result", "run", "predecessor_run_id") != prior || a.runID == prior {
		t.Fatal("recovery did not retain the predecessor link")
	}
	a.repository = p21Path(t, recovered, "result", "run", "start_repository").(map[string]any)
	resume := p21Run(t, a.binary, a.workspace, "--json", "resume", a.issueID, "--repository", a.workspace)
	if p22ResultPath(resume, "result.resume.obligations.0.freshness") != "stale" {
		t.Fatalf("adoption alone refreshed proof: %#v", resume)
	}
	secondProof := p22JourneyProof(t, a, "journey-successor")
	if firstProof == secondProof {
		t.Fatal("successor reused prior proof identity")
	}
	resume = p21Run(t, a.binary, a.workspace, "--json", "resume", a.issueID, "--repository", a.workspace)
	if p22ResultPath(resume, "result.resume.obligations.0.freshness") != "fresh" {
		t.Fatalf("successor proof not fresh: %#v", resume)
	}
	// This continues the same Issue and linked Runs through actual export,
	// validation, import, both rebuilds, named resume reads and Event-byte parity.
	entries, views := p22PortableThreeWay(t, a, []string{a.activation["id"].(string)})
	if !reflect.DeepEqual(p22PortableNamedSummary(views, entries), p22PortableFakeSummary(t, entries, a.runID)) {
		t.Fatal("final journey fake/fid summaries disagree")
	}
	if p22ResultPath(views["resume"].(map[string]any), "issue.id") != a.issueID {
		t.Fatal("roundtrip changed the journey Issue")
	}
	return a
}

func p22AcceptedCall(t *testing.T, a *p22FIDRecoveryAdapter, name, operation string, input map[string]any) map[string]any {
	t.Helper()
	response := a.Execute(t, p22ConformanceCase{Name: name, Operation: operation, Kind: "accepted", Input: input})
	if response.ExitCode != 0 || response.Stderr != "" {
		t.Fatalf("%s: %#v", operation, response)
	}
	var envelope map[string]any
	if err := json.Unmarshal([]byte(response.Stdout), &envelope); err != nil {
		t.Fatal(err)
	}
	return envelope
}

func p22JourneyProof(t *testing.T, a *p22FIDRecoveryAdapter, key string) string {
	t.Helper()
	gate := p22GateRun(t, a.binary, a.workspace, a.actor, p21Input(t, a.workspace, key+"-gate", map[string]any{"run_id": a.runID, "registration_id": a.registrationID, "registration_version": a.registrationVersion, "obligation_id": "OBL-p22-corpus", "expected_repository": map[string]any{"head_oid": a.repository["head_oid"], "dirty": a.repository["dirty"]}}), key+"-gate")
	proof := p21Run(t, a.binary, a.workspace, "--json", "evidence", "link", "--actor", a.actor, "--input", p21Input(t, a.workspace, key+"-proof", map[string]any{"obligation_id": "OBL-p22-corpus", "evidence_type": "receipt", "receipt_id": p21PathString(t, gate, "result", "receipt", "id"), "review_id": nil, "human_outcome_id": nil, "required_artifact_ids": []string{}}), "--idempotency-key", key+"-proof")
	if a.firstResults != nil {
		a.firstResults["evidence.link"] = proof["result"].(map[string]any)
	}
	run := p21Run(t, a.binary, a.workspace, "--json", "run", "get", a.runID)
	a.runVersion = p22Number(p21Path(t, run, "result", "run", "version"))
	return p21PathString(t, proof, "result", "evidence", "id")
}
