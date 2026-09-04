package main

import (
	"database/sql"
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/mbell-nzl/fid/internal/config"
	fiddb "github.com/mbell-nzl/fid/internal/db"
	"github.com/mbell-nzl/fid/internal/journal"
)

func TestEvidenceLinkAndEvaluateKeepsTrustedDimensionsIndependent(t *testing.T) {
	t.Parallel()
	t.Run("validity changes only validity", func(t *testing.T) {
		fixture := newLinkedEvidenceFixture(t, "validity")
		assertEvaluation(t, fixture.evaluate(t, fixture.repository), "valid", "fresh", "available", []string{})
		fixture.updateValidity(t, "invalid")
		assertEvaluation(t, fixture.evaluate(t, fixture.repository), "invalid", "fresh", "available", []string{"evidence-invalid"})
	})

	t.Run("tree state changes freshness while branch label does not", func(t *testing.T) {
		fixture := newLinkedEvidenceFixture(t, "freshness")
		assertEvaluation(t, fixture.evaluate(t, fixture.repository), "valid", "fresh", "available", []string{})

		branchMoved := cloneJSONMap(t, fixture.repository)
		branchMoved["head_oid"] = strings.Repeat("c", 40)
		assertEvaluation(t, fixture.evaluate(t, branchMoved), "valid", "fresh", "available", []string{})

		changedTree := cloneJSONMap(t, fixture.repository)
		changedTree["tree_oid"] = strings.Repeat("d", 40)
		assertEvaluation(t, fixture.evaluate(t, changedTree), "valid", "stale", "available", []string{"evidence-stale"})
	})

	t.Run("Artifact availability changes only Artifact availability", func(t *testing.T) {
		fixture := newLinkedEvidenceFixture(t, "artifact")
		assertEvaluation(t, fixture.evaluate(t, fixture.repository), "valid", "fresh", "available", []string{})
		fixture.updateArtifactAvailability(t, "missing")
		assertEvaluation(t, fixture.evaluate(t, fixture.repository), "valid", "fresh", "missing", []string{"artifact-unavailable"})
	})
}

func TestEvidenceEvaluationMultiFailureGoldenSurvivesRepeatAndRebuild(t *testing.T) {
	t.Parallel()
	fixture := newLinkedEvidenceFixture(t, "golden")
	fixture.observeArtifactMissing(t)
	staleRepository := cloneJSONMap(t, fixture.repository)
	staleRepository["tree_oid"] = strings.Repeat("e", 40)

	assertEvaluation(t, fixture.evaluate(t, staleRepository), "valid", "stale", "missing", []string{"evidence-stale", "artifact-unavailable"})
	first := fixture.evaluateGolden(t, staleRepository)
	if repeated := fixture.evaluateGolden(t, staleRepository); repeated != first {
		t.Fatalf("repeated multi-failure Evaluation = %s, want byte-identical %s", repeated, first)
	}

	rebuilt := executeCommandForTest(fixture.workspace, "--project", fixture.workspace, "--json", "rebuild")
	if rebuilt.err != nil || rebuilt.stderr != "" {
		t.Fatalf("rebuild Evidence fixture = %#v, want journal rebuild", rebuilt)
	}
	if afterRebuild := fixture.evaluateGolden(t, staleRepository); afterRebuild != first {
		t.Fatalf("post-rebuild multi-failure Evaluation = %s, want byte-identical %s", afterRebuild, first)
	}
}

type linkedEvidenceFixture struct {
	workspace    string
	obligationID string
	evidenceID   string
	artifactID   string
	digest       string
	repository   map[string]any
}

func newLinkedEvidenceFixture(t *testing.T, _ string) linkedEvidenceFixture {
	t.Helper()
	workspace, values := cloneReusableWorkspaceFixture(t, "linked-evidence", func() (string, []string) {
		fixture := buildLinkedEvidenceFixture(t, "cached")
		repositoryJSON, err := json.Marshal(fixture.repository)
		if err != nil {
			t.Fatalf("encode linked Evidence repository: %v", err)
		}
		return fixture.workspace, []string{fixture.obligationID, fixture.evidenceID, fixture.artifactID, fixture.digest, string(repositoryJSON)}
	})
	repository := make(map[string]any)
	if err := json.Unmarshal([]byte(values[4]), &repository); err != nil {
		t.Fatalf("decode linked Evidence repository: %v", err)
	}
	return linkedEvidenceFixture{workspace: workspace, obligationID: values[0], evidenceID: values[1], artifactID: values[2], digest: values[3], repository: repository}
}

func buildLinkedEvidenceFixture(t *testing.T, key string) linkedEvidenceFixture {
	t.Helper()
	workspace, issueID, envelopeID := runEvidenceEnvelopeForTest(t, func(definition map[string]any) {
		definition["argv"] = []string{"/bin/true"}
		definition["artifact_paths"] = []string{"evidence.txt"}
	})
	if err := os.WriteFile(filepath.Join(workspace, "evidence.txt"), []byte("trusted evidence\n"), 0o600); err != nil {
		t.Fatalf("write declared Artifact: %v", err)
	}
	terminal, gateRun := executeGateRunForTest(t, workspace, issueID, envelopeID, "evidence-"+key)
	if gateRun.err != nil || terminal.Data.ReceiptID == "" {
		t.Fatalf("gate run = %#v, terminal=%#v, want successful Receipt", gateRun, terminal)
	}

	conn, err := fiddb.Open(filepath.Join(workspace, config.WorkspaceDirName, config.DatabaseFileName))
	if err != nil {
		t.Fatalf("open projection: %v", err)
	}
	defer func() { _ = conn.Close() }()
	var artifactID, digest string
	if err := conn.QueryRow(`SELECT a.id, a.digest FROM referenced_artifacts a JOIN receipt_artifacts l ON l.artifact_id = a.id WHERE l.receipt_id = ?`, terminal.Data.ReceiptID).Scan(&artifactID, &digest); err != nil {
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
		"--input", linkPath, "--idempotency-key", "evidence-link-"+key,
	)
	if linked.err != nil || linked.stderr != "" {
		t.Fatalf("evidence link = %#v, want trusted Receipt link", linked)
	}
	evidence := decodeCanonicalResult(t, linked.stdout)["evidence"].(map[string]any)
	if evidence["receipt_id"] != terminal.Data.ReceiptID || evidence["producer"] != "fid" || evidence["trust_source"] != "fid-executed" {
		t.Fatalf("linked Evidence = %#v, want exact fid-executed Receipt binding", evidence)
	}
	return linkedEvidenceFixture{workspace: workspace, obligationID: obligationID, evidenceID: evidence["id"].(string), artifactID: artifactID, digest: digest, repository: repository}
}

func (fixture linkedEvidenceFixture) evaluate(t *testing.T, repository map[string]any) map[string]any {
	t.Helper()
	path := filepath.Join(fixture.workspace, "evidence-evaluate.json")
	writeActivationInput(t, path, map[string]any{"obligation_id": fixture.obligationID, "target_repository": repository})
	result := executeCommandForTest(fixture.workspace, "--project", fixture.workspace, "--json", "evidence", "evaluate", "--input", path)
	if result.err != nil || result.stderr != "" {
		t.Fatalf("evidence evaluate = %#v", result)
	}
	return decodeCanonicalResult(t, result.stdout)["evaluation"].(map[string]any)
}

func (fixture linkedEvidenceFixture) evaluateGolden(t *testing.T, repository map[string]any) string {
	t.Helper()
	payload, err := journal.MarshalCanonicalJSON(fixture.evaluate(t, repository))
	if err != nil {
		t.Fatalf("canonicalize Evaluation: %v", err)
	}
	return string(payload)
}

func (fixture linkedEvidenceFixture) updateValidity(t *testing.T, validity string) {
	t.Helper()
	fixture.updateProjection(t, `UPDATE evidence_records SET validity = ? WHERE id = ?`, validity, fixture.evidenceID)
}

func (fixture linkedEvidenceFixture) updateArtifactAvailability(t *testing.T, availability string) {
	t.Helper()
	fixture.updateProjection(t, `UPDATE referenced_artifacts SET availability = ? WHERE id = ?`, availability, fixture.artifactID)
}

func (fixture linkedEvidenceFixture) updateProjection(t *testing.T, query string, arguments ...any) {
	t.Helper()
	conn, err := fiddb.Open(filepath.Join(fixture.workspace, config.WorkspaceDirName, config.DatabaseFileName))
	if err != nil {
		t.Fatalf("open projection: %v", err)
	}
	defer func() { _ = conn.Close() }()
	if _, err := conn.Exec(query, arguments...); err != nil {
		t.Fatalf("update Evidence fixture: %v", err)
	}
}

func (fixture linkedEvidenceFixture) observeArtifactMissing(t *testing.T) {
	t.Helper()
	path := filepath.Join(fixture.workspace, config.WorkspaceDirName, "artifacts", strings.TrimPrefix(fixture.digest, "sha256:"))
	if err := os.Remove(path); err != nil {
		t.Fatalf("remove stored Artifact: %v", err)
	}
	inputPath := filepath.Join(fixture.workspace, "artifact-verify.json")
	writeActivationInput(t, inputPath, map[string]any{"artifact_id": fixture.artifactID})
	result := executeCommandForTest(fixture.workspace,
		"--project", fixture.workspace, "--json", "artifact", "verify", "--actor", "actor-builder",
		"--input", inputPath, "--idempotency-key", "artifact-verify-missing",
	)
	if result.err != nil || result.stderr != "" {
		t.Fatalf("missing Artifact verify = %#v, want availability observation", result)
	}
	assertArtifactVerifyResult(t, result.stdout, "missing", true)
}

func TestEvidenceLinkRejectsEveryCallerAuthoredProof(t *testing.T) {
	t.Parallel()
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
			beforeEvents, beforeEvidence := evidenceStateCounts(t, workspace)
			result := executeCommandForTest(workspace, "--project", workspace, "--json", "evidence", "link", "--actor", "actor-builder", "--input", path, "--idempotency-key", "caller-"+proof)
			assertEvidenceLinkRejectedWithoutEvidenceEvent(t, workspace, result, beforeEvents, beforeEvidence)
		})
	}
}

func TestEvidenceLinkRejectsPersistedReceiptBindingMismatches(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name       string
		mutate     func(t *testing.T, conn *sql.DB, receiptID string)
		mutateLink func(t *testing.T, repository map[string]any)
	}{
		{
			name: "Obligation",
			mutate: func(t *testing.T, conn *sql.DB, receiptID string) {
				t.Helper()
				updateReceiptAttemptField(t, conn, receiptID, "obligation_id", "OBL-mismatched")
			},
		},
		{
			name: "Registration",
			mutate: func(t *testing.T, conn *sql.DB, receiptID string) {
				t.Helper()
				updateReceiptAttemptField(t, conn, receiptID, "registration_id", "REG-mismatched")
			},
		},
		{
			name: "Run",
			mutate: func(t *testing.T, conn *sql.DB, receiptID string) {
				t.Helper()
				if _, err := conn.Exec(`UPDATE runs SET envelope_revision = 2 WHERE id = (SELECT a.run_id FROM execution_receipts r JOIN execution_attempts a ON a.id = r.attempt_id WHERE r.id = ?)`, receiptID); err != nil {
					t.Fatalf("mismatch Receipt Run binding: %v", err)
				}
			},
		},
		{
			name: "Issue",
			mutate: func(t *testing.T, conn *sql.DB, receiptID string) {
				t.Helper()
				if _, err := conn.Exec(`UPDATE runs SET issue_id = 'ISS-mismatched' WHERE id = (SELECT a.run_id FROM execution_receipts r JOIN execution_attempts a ON a.id = r.attempt_id WHERE r.id = ?)`, receiptID); err != nil {
					t.Fatalf("mismatch Receipt Issue binding: %v", err)
				}
			},
		},
		{
			name: "repository target",
			mutateLink: func(t *testing.T, repository map[string]any) {
				t.Helper()
				repository["tree_oid"] = strings.Repeat("f", 40)
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			workspace, issueID, envelopeID := runEvidenceEnvelopeForTest(t, func(definition map[string]any) {
				definition["argv"] = []string{"/bin/true"}
			})
			terminal, gateRun := executeGateRunForTest(t, workspace, issueID, envelopeID, "binding-mismatch-"+strings.ReplaceAll(testCase.name, " ", "-"))
			if gateRun.err != nil || terminal.Data.ReceiptID == "" {
				t.Fatalf("gate run = %#v, terminal=%#v, want successful Receipt", gateRun, terminal)
			}

			conn, err := fiddb.Open(filepath.Join(workspace, config.WorkspaceDirName, config.DatabaseFileName))
			if err != nil {
				t.Fatalf("open projection: %v", err)
			}
			if testCase.mutate != nil {
				if _, err := conn.Exec(`PRAGMA foreign_keys = OFF`); err != nil {
					_ = conn.Close()
					t.Fatalf("disable foreign keys for persisted mismatch fixture: %v", err)
				}
				testCase.mutate(t, conn, terminal.Data.ReceiptID)
			}
			if err := conn.Close(); err != nil {
				t.Fatalf("close projection: %v", err)
			}

			repository := runRepositoryForTest(t, workspace, issueID, envelopeID)
			if testCase.mutateLink != nil {
				testCase.mutateLink(t, repository)
			}
			linkPath := filepath.Join(workspace, "binding-mismatch-evidence.json")
			writeActivationInput(t, linkPath, map[string]any{
				"obligation_id": evidenceObligationID(t, workspace, envelopeID), "receipt_id": terminal.Data.ReceiptID,
				"target_repository": repository, "required_artifact_ids": []string{},
			})
			beforeEvents, beforeEvidence := evidenceStateCounts(t, workspace)
			result := executeCommandForTest(workspace,
				"--project", workspace, "--json", "evidence", "link", "--actor", "actor-builder",
				"--input", linkPath, "--idempotency-key", "binding-mismatch-"+strings.ReplaceAll(testCase.name, " ", "-"),
			)
			assertEvidenceLinkRejectedWithoutEvidenceEvent(t, workspace, result, beforeEvents, beforeEvidence)
		})
	}
}

func runEvidenceEnvelopeForTest(t *testing.T, modifyDefinition func(map[string]any)) (string, string, string) {
	t.Helper()
	definitions := gateDefinitionsForTest(modifyDefinition)
	key := reusableWorkspaceFixtureKey(t, "evidence-envelope", definitions)
	workspace, values := cloneReusableWorkspaceFixture(t, key, func() (string, []string) {
		workspace, issueID, envelopeID := buildEvidenceEnvelopeForTest(t, definitions)
		return workspace, []string{issueID, envelopeID}
	})
	return workspace, values[0], values[1]
}

func buildEvidenceEnvelopeForTest(t *testing.T, definitions []any) (string, string, string) {
	t.Helper()
	workspace, issueID, registration := finalizedEnvelopeWorkspaceWithDefinitions(t, definitions)
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

func updateReceiptAttemptField(t *testing.T, conn *sql.DB, receiptID, field, value string) {
	t.Helper()
	if _, err := conn.Exec(`UPDATE execution_attempts SET `+field+` = ? WHERE id = (SELECT attempt_id FROM execution_receipts WHERE id = ?)`, value, receiptID); err != nil {
		t.Fatalf("mismatch Receipt %s binding: %v", field, err)
	}
}

func evidenceStateCounts(t *testing.T, workspace string) (int, int) {
	t.Helper()
	conn, err := fiddb.Open(filepath.Join(workspace, config.WorkspaceDirName, config.DatabaseFileName))
	if err != nil {
		t.Fatalf("open projection for Evidence count: %v", err)
	}
	defer func() { _ = conn.Close() }()
	var evidenceRecords int
	if err := conn.QueryRow(`SELECT COUNT(*) FROM evidence_records`).Scan(&evidenceRecords); err != nil {
		t.Fatalf("count Evidence records: %v", err)
	}
	return journalEventCount(t, workspace), evidenceRecords
}

func assertEvidenceLinkRejectedWithoutEvidenceEvent(t *testing.T, workspace string, result commandExecutionResult, wantEvents, wantEvidence int) {
	t.Helper()
	if result.err == nil || result.stdout != "" || !strings.Contains(result.stderr, `"code":"invalid-argument"`) {
		t.Fatalf("evidence link = %#v, want typed rejection", result)
	}
	if gotEvents, gotEvidence := evidenceStateCounts(t, workspace); gotEvents != wantEvents || gotEvidence != wantEvidence {
		t.Fatalf("Evidence rejection state = events:%d evidence:%d, want events:%d evidence:%d", gotEvents, gotEvidence, wantEvents, wantEvidence)
	}
}
