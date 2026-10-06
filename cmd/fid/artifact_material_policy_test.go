package main

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestArtifactLossHonorsRequiredMaterialAndPreservesProofIdentity(t *testing.T) {
	t.Parallel()
	for _, required := range []bool{false, true} {
		name := "metadata-sufficient"
		if required {
			name = "required-material"
		}
		t.Run(name, func(t *testing.T) {
			workspace, issueID, envelopeID := gateReadyEvidenceEnvelopeForTest(t, func(definition map[string]any) {
				definition["argv"] = []string{"/bin/true"}
				definition["artifact_paths"] = []string{"evidence.txt"}
			})
			if err := os.WriteFile(filepath.Join(workspace, "evidence.txt"), []byte("retained proof material\n"), 0o600); err != nil {
				t.Fatal(err)
			}
			terminal, result := executeGateRunForTest(t, workspace, issueID, envelopeID, "material-policy-gate")
			if result.err != nil || terminal.Data.ReceiptID == "" {
				t.Fatalf("Gate with Artifact: result=%#v terminal=%#v", result, terminal)
			}
			artifacts := terminal.Data.Response["result"].(map[string]any)["receipt"].(map[string]any)["artifacts"].([]any)
			if len(artifacts) != 1 {
				t.Fatalf("Receipt Artifacts = %#v, want one", artifacts)
			}
			artifactID := artifacts[0].(map[string]any)["id"].(string)
			receiptID := terminal.Data.ReceiptID
			read := func(entity, id string) map[string]any {
				t.Helper()
				result := executeCommandForTest(workspace, "--project", workspace, "--json", entity, "get", id)
				if result.err != nil {
					t.Fatalf("read %s: %#v", entity, result)
				}
				return decodeCanonicalResult(t, result.stdout)[entity].(map[string]any)
			}
			artifact := read("receipt", receiptID)["artifacts"].([]any)[0].(map[string]any)
			receipt := read("receipt", receiptID)
			requiredIDs := []string{}
			if required {
				requiredIDs = []string{artifactID}
			}
			repository := runRepositoryForTest(t, workspace)
			linkPath := filepath.Join(workspace, "material-policy-link.json")
			writeActivationInput(t, linkPath, map[string]any{"obligation_id": testEvidenceObligationID, "evidence_type": "receipt", "receipt_id": receiptID, "review_id": nil, "human_outcome_id": nil, "required_artifact_ids": requiredIDs})
			linked := executeCommandForTest(workspace, "--project", workspace, "--json", "evidence", "link", "--actor", "actor-builder", "--input", linkPath, "--idempotency-key", "material-policy-link")
			if linked.err != nil {
				t.Fatalf("link Evidence: %#v", linked)
			}
			evidence := decodeCanonicalResult(t, linked.stdout)["evidence"].(map[string]any)
			fixture := linkedEvidenceFixture{workspace: workspace, obligationID: testEvidenceObligationID, evidenceID: evidence["id"].(string), artifactID: artifactID, digest: artifact["digest"].(string), repository: repository}
			assertEvaluation(t, fixture.evaluate(t), "valid", "fresh", "available", []string{})
			fixture.observeArtifactMissing(t)
			if required {
				assertEvaluation(t, fixture.evaluate(t), "valid", "fresh", "missing", []string{"artifact-unavailable"})
			} else {
				assertEvaluation(t, fixture.evaluate(t), "valid", "fresh", "available", []string{})
			}
			if got := read("evidence", fixture.evidenceID); !reflect.DeepEqual(got, evidence) {
				t.Fatalf("loss changed immutable Evidence: before=%#v after=%#v", evidence, got)
			}
			receipt["artifacts"].([]any)[0].(map[string]any)["availability"] = "missing"
			if got := read("receipt", receiptID); !reflect.DeepEqual(got, receipt) {
				t.Fatalf("loss changed immutable Receipt: before=%#v after=%#v", receipt, got)
			}
			afterArtifact := read("receipt", receiptID)["artifacts"].([]any)[0].(map[string]any)
			if afterArtifact["availability"] != "missing" {
				t.Fatalf("Artifact availability = %#v, want missing", afterArtifact)
			}
			artifact["availability"] = "missing"
			if !reflect.DeepEqual(afterArtifact, artifact) {
				t.Fatalf("loss changed Artifact metadata beyond availability: before=%#v after=%#v", artifact, afterArtifact)
			}
			if _, err := os.Stat(filepath.Join(workspace, ".fid", "artifacts", strings.TrimPrefix(fixture.digest, "sha256:"))); !os.IsNotExist(err) {
				t.Fatalf("missing Artifact was recreated: %v", err)
			}
		})
	}
}
