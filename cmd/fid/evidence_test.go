package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/mbell-nzl/fid/internal/config"
	fiddb "github.com/mbell-nzl/fid/internal/db"
)

func TestEvidenceLinkAndEvaluateKeepsTrustedDimensionsIndependent(t *testing.T) {
	workspace, issueID, envelopeID := runEvidenceEnvelopeForTest(t, func(definition map[string]any) {
		definition["argv"] = []string{"/bin/true"}
		definition["artifact_paths"] = []string{"evidence.txt"}
	})
	if err := os.WriteFile(filepath.Join(workspace, "evidence.txt"), []byte("trusted evidence\n"), 0o600); err != nil {
		t.Fatalf("write declared Artifact: %v", err)
	}
	terminal, gateRun := executeGateRunForTest(t, workspace, issueID, envelopeID, "evidence-source")
	if gateRun.err != nil || terminal.Data.ReceiptID == "" {
		t.Fatalf("gate run = %#v, terminal=%#v, want successful Receipt", gateRun, terminal)
	}

	conn, err := fiddb.Open(filepath.Join(workspace, config.WorkspaceDirName, config.DatabaseFileName))
	if err != nil {
		t.Fatalf("open projection: %v", err)
	}
	defer func() { _ = conn.Close() }()
	var artifactID string
	if err := conn.QueryRow(`SELECT artifact_id FROM receipt_artifacts WHERE receipt_id = ?`, terminal.Data.ReceiptID).Scan(&artifactID); err != nil {
		t.Fatalf("read Receipt Artifact: %v", err)
	}
	repository := runRepositoryForTest(t, workspace, issueID, envelopeID)
	obligationID := evidenceObligationID(t, workspace, envelopeID)
	linkPath := filepath.Join(workspace, "evidence-link.json")
	writeActivationInput(t, linkPath, map[string]any{
		"obligation_id": obligationID, "receipt_id": terminal.Data.ReceiptID,
		"target_repository": repository, "required_artifact_ids": []string{artifactID},
	})

	linked := executeCommandForTest(workspace,
		"--project", workspace, "--json", "evidence", "link", "--actor", "actor-builder",
		"--input", linkPath, "--idempotency-key", "evidence-link-qualified",
	)
	if linked.err != nil || linked.stderr != "" {
		t.Fatalf("evidence link = %#v, want trusted Receipt link", linked)
	}
	linkedResult := decodeCanonicalResult(t, linked.stdout)
	evidence := linkedResult["evidence"].(map[string]any)
	if evidence["receipt_id"] != terminal.Data.ReceiptID || evidence["producer"] != "fid" || evidence["trust_source"] != "fid-executed" {
		t.Fatalf("linked Evidence = %#v, want exact fid-executed Receipt binding", evidence)
	}
	assertEvaluation(t, linkedResult["evaluation"].(map[string]any), "valid", "fresh", "available", []string{})

	evaluatePath := filepath.Join(workspace, "evidence-evaluate.json")
	writeActivationInput(t, evaluatePath, map[string]any{"obligation_id": obligationID, "target_repository": repository})
	first := executeCommandForTest(workspace, "--project", workspace, "--json", "evidence", "evaluate", "--input", evaluatePath)
	if first.err != nil || first.stderr != "" {
		t.Fatalf("fresh evidence evaluate = %#v", first)
	}
	assertEvaluation(t, decodeCanonicalResult(t, first.stdout)["evaluation"].(map[string]any), "valid", "fresh", "available", []string{})

	if _, err := conn.Exec(`UPDATE evidence_records SET validity = 'invalid' WHERE id = ?`, evidence["id"]); err != nil {
		t.Fatalf("make Evidence invalid: %v", err)
	}
	if _, err := conn.Exec(`UPDATE referenced_artifacts SET availability = 'missing' WHERE id = ?`, artifactID); err != nil {
		t.Fatalf("make required Artifact missing: %v", err)
	}
	staleRepository := cloneJSONMap(t, repository)
	staleRepository["tree_oid"] = strings.Repeat("c", 40)
	writeActivationInput(t, evaluatePath, map[string]any{"obligation_id": obligationID, "target_repository": staleRepository})
	second := executeCommandForTest(workspace, "--project", workspace, "--json", "evidence", "evaluate", "--input", evaluatePath)
	if second.err != nil || second.stderr != "" {
		t.Fatalf("multi-failure evidence evaluate = %#v", second)
	}
	assertEvaluation(t, decodeCanonicalResult(t, second.stdout)["evaluation"].(map[string]any), "invalid", "stale", "missing", []string{"evidence-invalid", "evidence-stale", "artifact-unavailable"})
}

func TestEvidenceLinkRejectsEveryCallerAuthoredProof(t *testing.T) {
	workspace, issueID, envelopeID := runEvidenceEnvelopeForTest(t, func(definition map[string]any) {
		definition["argv"] = []string{"/bin/true"}
	})
	terminal, gateRun := executeGateRunForTest(t, workspace, issueID, envelopeID, "caller-proof-source")
	if gateRun.err != nil || terminal.Data.ReceiptID == "" {
		t.Fatalf("gate run = %#v, terminal=%#v, want successful Receipt", gateRun, terminal)
	}
	repository := runRepositoryForTest(t, workspace, issueID, envelopeID)
	obligationID := evidenceObligationID(t, workspace, envelopeID)
	for _, proof := range []string{"pass", "command", "import", "attestation", "role"} {
		t.Run(proof, func(t *testing.T) {
			path := filepath.Join(workspace, proof+"-evidence.json")
			input := map[string]any{"obligation_id": obligationID, "receipt_id": terminal.Data.ReceiptID, "target_repository": repository, "required_artifact_ids": []string{}}
			input[proof] = true
			writeActivationInput(t, path, input)
			result := executeCommandForTest(workspace, "--project", workspace, "--json", "evidence", "link", "--actor", "actor-builder", "--input", path, "--idempotency-key", "caller-"+proof)
			if result.err == nil || result.stdout != "" {
				t.Fatalf("caller-authored %s proof = %#v, want typed rejection and no Evidence", proof, result)
			}
		})
	}
}

func runEvidenceEnvelopeForTest(t *testing.T, modifyDefinition func(map[string]any)) (string, string, string) {
	t.Helper()
	workspace, issueID, registration := finalizedEnvelopeWorkspaceWithDefinition(t, modifyDefinition)
	input := envelopeCreateInput(t, workspace, issueID, registration)
	grant := input["authority_grants"].([]any)[0].(map[string]any)
	grant["actions"] = []string{"claim.acquire", "claim.release", "envelope.activate", "envelope.complete", "evidence.link", "gate.run", "issue.transition", "run.finish", "run.start"}
	grant["actor_selector"] = map[string]any{"actor_id": "actor-builder", "required_capabilities": []string{"claim.acquire", "claim.release", "envelope.activate", "envelope.complete", "evidence.link", "gate.run", "issue.transition", "run.finish", "run.start"}}
	path := filepath.Join(workspace, "evidence-envelope.json")
	writeActivationInput(t, path, input)
	created := executeCommandForTest(workspace, "--project", workspace, "--json", "envelope", "create", "--actor", "actor-builder", "--input", path, "--idempotency-key", "evidence-envelope-create")
	if created.err != nil {
		t.Fatalf("create evidence envelope: %v; stderr=%q", created.err, created.stderr)
	}
	envelopeID := assertEnvelopeShape(t, created.stdout, "draft")
	activatePath := filepath.Join(workspace, "evidence-envelope-activate.json")
	writeActivationInput(t, activatePath, map[string]any{"envelope_id": envelopeID, "expected_revision": 1})
	activated := executeCommandForTest(workspace, "--project", workspace, "--json", "envelope", "activate", "--actor", "actor-builder", "--input", activatePath, "--idempotency-key", "evidence-envelope-activate")
	if activated.err != nil {
		t.Fatalf("activate evidence envelope: %v; stderr=%q", activated.err, activated.stderr)
	}
	return workspace, issueID, envelopeID
}

func evidenceObligationID(t *testing.T, workspace, envelopeID string) string {
	t.Helper()
	output := executeCommandForTest(workspace, "--project", workspace, "--json", "envelope", "get", envelopeID)
	if output.err != nil {
		t.Fatalf("get evidence envelope: %v; stderr=%q", output.err, output.stderr)
	}
	return decodeCanonicalResult(t, output.stdout)["envelope"].(map[string]any)["obligations"].([]any)[0].(map[string]any)["id"].(string)
}

func assertEvaluation(t *testing.T, evaluation map[string]any, validity, freshness, availability string, reasonCodes []string) {
	t.Helper()
	if evaluation["validity"] != validity || evaluation["freshness"] != freshness || evaluation["artifact_availability"] != availability {
		t.Fatalf("evaluation dimensions = %#v, want validity=%s freshness=%s artifact_availability=%s", evaluation, validity, freshness, availability)
	}
	got := make([]string, 0, len(evaluation["reasons"].([]any)))
	for _, reason := range evaluation["reasons"].([]any) {
		got = append(got, reason.(map[string]any)["code"].(string))
	}
	if !reflect.DeepEqual(got, reasonCodes) {
		t.Fatalf("evaluation reasons = %#v, want %#v", got, reasonCodes)
	}
}

func cloneJSONMap(t *testing.T, value map[string]any) map[string]any {
	t.Helper()
	payload, err := json.Marshal(value)
	if err != nil {
		t.Fatalf("marshal map: %v", err)
	}
	var cloned map[string]any
	if err := json.Unmarshal(payload, &cloned); err != nil {
		t.Fatalf("unmarshal map: %v", err)
	}
	return cloned
}
