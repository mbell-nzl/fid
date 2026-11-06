package main

import (
	"database/sql"
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/mbell-nzl/fid/internal/cli"
	"github.com/mbell-nzl/fid/internal/config"
	fiddb "github.com/mbell-nzl/fid/internal/db"
	"github.com/mbell-nzl/fid/internal/journal"
)

func TestEvidenceLinkAndEvaluateKeepsTrustedDimensionsIndependent(t *testing.T) {
	t.Parallel()
	t.Run("validity changes only validity", func(t *testing.T) {
		fixture := newLinkedEvidenceFixture(t, "validity")
		assertEvaluation(t, fixture.evaluate(t), "valid", "fresh", "available", []string{})
		fixture.updateValidity(t, "invalid")
		assertEvaluation(t, fixture.evaluate(t), "invalid", "fresh", "available", []string{"evidence-invalid"})
	})

	t.Run("tree state changes freshness while branch label does not", func(t *testing.T) {
		fixture := newLinkedEvidenceFixture(t, "freshness")
		assertEvaluation(t, fixture.evaluate(t), "valid", "fresh", "available", []string{})

		before := testEnvelopeRepository(t, fixture.workspace)
		runGitForTest(t, fixture.workspace, "commit", "--allow-empty", "-m", "branch label only")
		after := testEnvelopeRepository(t, fixture.workspace)
		if before["head_oid"] == after["head_oid"] || before["tree_oid"] != after["tree_oid"] {
			t.Fatal("branch control must change HEAD without changing tree")
		}
		assertEvaluation(t, fixture.evaluate(t), "valid", "fresh", "available", []string{})

		if err := os.WriteFile(filepath.Join(fixture.workspace, "tracked-proof.txt"), []byte("changed"), 0o600); err != nil {
			t.Fatal(err)
		}
		assertEvaluation(t, fixture.evaluate(t), "valid", "stale", "available", []string{"evidence-stale"})
	})

	t.Run("Artifact availability changes only Artifact availability", func(t *testing.T) {
		fixture := newLinkedEvidenceFixture(t, "artifact")
		assertEvaluation(t, fixture.evaluate(t), "valid", "fresh", "available", []string{})
		fixture.observeArtifactMissing(t)
		assertEvaluation(t, fixture.evaluate(t), "valid", "fresh", "missing", []string{"artifact-unavailable"})
		got := executeCommandForTest(fixture.workspace, "--project", fixture.workspace, "--json", "evidence", "get", fixture.evidenceID)
		if got.err != nil || got.stderr != "" {
			t.Fatalf("Artifact loss changed Evidence identity: %#v", got)
		}
		evidence := decodeCanonicalResult(t, got.stdout)["evidence"].(map[string]any)
		if evidence["id"] != fixture.evidenceID {
			t.Fatalf("Artifact loss changed Evidence identity: %#v", evidence)
		}
	})
}

func TestEvidenceEvaluationMultiFailureGoldenSurvivesRepeatAndRebuild(t *testing.T) {
	t.Parallel()
	fixture := newLinkedEvidenceFixture(t, "golden")
	fixture.observeArtifactMissing(t)
	if err := os.WriteFile(filepath.Join(fixture.workspace, "tracked-proof.txt"), []byte("changed"), 0o600); err != nil {
		t.Fatal(err)
	}

	assertEvaluation(t, fixture.evaluate(t), "valid", "stale", "missing", []string{"evidence-stale", "artifact-unavailable"})
	first := fixture.evaluateGolden(t)
	if repeated := fixture.evaluateGolden(t); repeated != first {
		t.Fatalf("repeated multi-failure Evaluation = %s, want byte-identical %s", repeated, first)
	}

	rebuilt := executeCommandForTest(fixture.workspace, "--project", fixture.workspace, "--json", "rebuild")
	if rebuilt.err != nil || rebuilt.stderr != "" {
		t.Fatalf("rebuild Evidence fixture = %#v, want journal rebuild", rebuilt)
	}
	if afterRebuild := fixture.evaluateGolden(t); afterRebuild != first {
		t.Fatalf("post-rebuild multi-failure Evaluation = %s, want byte-identical %s", afterRebuild, first)
	}
}

func TestEvidenceEvaluateNullObligationReturnsActiveRevisionInIDOrder(t *testing.T) {
	t.Parallel()
	workspace, issueID, registration := finalizedEnvelopeWorkspace(t)
	input := typedEnvelopeCreateInput(t, workspace, issueID, registration)
	path := filepath.Join(workspace, "evidence-evaluate-all-envelope.json")
	writeActivationInput(t, path, input)
	created := executeCommandForTest(workspace, "--project", workspace, "--json", "envelope", "create", "--actor", "actor-builder", "--input", path, "--idempotency-key", "evidence-evaluate-all-envelope")
	if created.err != nil || created.stderr != "" {
		t.Fatalf("create multi-Obligation Envelope = %#v", created)
	}
	envelopeID := assertEnvelopeShape(t, created.stdout, "draft")
	writeActivationInput(t, path, map[string]any{"envelope_id": envelopeID, "expected_revision": 1})
	activated := executeCommandForTest(workspace, "--project", workspace, "--json", "envelope", "activate", "--actor", "actor-builder", "--input", path, "--idempotency-key", "evidence-evaluate-all-activate")
	if activated.err != nil || activated.stderr != "" {
		t.Fatalf("activate multi-Obligation Envelope = %#v", activated)
	}

	writeActivationInput(t, path, map[string]any{"obligation_id": nil})
	result := executeCommandForTest(workspace, "--project", workspace, "--json", "evidence", "evaluate", "--input", path)
	if result.err != nil || result.stderr != "" {
		t.Fatalf("evaluate active Envelope = %#v", result)
	}
	evaluations := decodeCanonicalResult(t, result.stdout)["evaluations"].([]any)
	gotIDs := make([]string, 0, len(evaluations))
	for _, item := range evaluations {
		gotIDs = append(gotIDs, item.(map[string]any)["obligation_id"].(string))
	}
	wantIDs := []string{"OBL-1", "OBL-human", "OBL-review"}
	if !reflect.DeepEqual(gotIDs, wantIDs) {
		t.Fatalf("evaluation IDs = %#v, want ID order %#v", gotIDs, wantIDs)
	}
	if evaluations[0].(map[string]any)["artifact_availability"] != "missing" || evaluations[1].(map[string]any)["artifact_availability"] != "not-applicable" || evaluations[2].(map[string]any)["artifact_availability"] != "not-applicable" {
		t.Fatalf("evaluation dimensions = %#v, want independent per-policy Artifact states", evaluations)
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
	return buildLinkedEvidenceFixture(t, "local")
}

func buildLinkedEvidenceFixture(t *testing.T, key string) linkedEvidenceFixture {
	t.Helper()
	workspace, issueID, envelopeID := gateReadyEvidenceEnvelopeForTest(t, func(definition map[string]any) {
		definition["argv"] = []string{"/bin/true"}
		definition["artifact_paths"] = []string{"evidence.txt"}
	})
	if err := os.WriteFile(filepath.Join(workspace, "evidence.txt"), []byte("trusted evidence\n"), 0o600); err != nil {
		t.Fatalf("write declared Artifact: %v", err)
	}
	if err := os.WriteFile(filepath.Join(workspace, ".gitignore"), []byte("*.json\nevidence.txt\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(workspace, "tracked-proof.txt"), []byte("original"), 0o600); err != nil {
		t.Fatal(err)
	}
	runGitForTest(t, workspace, "add", ".gitignore", "tracked-proof.txt")
	runGitForTest(t, workspace, "commit", "-m", "proof fixture")
	terminal, gateRun := executeGateRunForTest(t, workspace, issueID, envelopeID, "evidence-"+key)
	if gateRun.err != nil || terminal.Data.ReceiptID == "" {
		t.Fatalf("gate run = %#v, terminal=%#v, want successful Receipt", gateRun, terminal)
	}

	conn, err := fiddb.Open(filepath.Join(config.LedgerDir(workspace), config.DatabaseFileName))
	if err != nil {
		t.Fatalf("open projection: %v", err)
	}
	defer func() { _ = conn.Close() }()
	var artifactID, digest string
	if err := conn.QueryRow(`SELECT a.id, a.digest FROM referenced_artifacts a JOIN receipt_artifacts l ON l.artifact_id = a.id WHERE l.receipt_id = ?`, terminal.Data.ReceiptID).Scan(&artifactID, &digest); err != nil {
		t.Fatalf("read Receipt Artifact: %v", err)
	}
	repository := runRepositoryForTest(t, workspace)
	obligationID := testEvidenceObligationID
	linkPath := filepath.Join(workspace, "evidence-link.json")
	writeActivationInput(t, linkPath, map[string]any{
		"obligation_id": obligationID, "evidence_type": "receipt", "receipt_id": terminal.Data.ReceiptID, "review_id": nil, "human_outcome_id": nil,
		"required_artifact_ids": []string{artifactID},
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

func (fixture linkedEvidenceFixture) evaluate(t *testing.T) map[string]any {
	t.Helper()
	path := filepath.Join(fixture.workspace, "evidence-evaluate.json")
	writeActivationInput(t, path, map[string]any{"obligation_id": fixture.obligationID})
	result := executeCommandForTest(fixture.workspace, "--project", fixture.workspace, "--json", "evidence", "evaluate", "--input", path)
	if result.err != nil || result.stderr != "" {
		t.Fatalf("evidence evaluate = %#v", result)
	}
	return decodeCanonicalResult(t, result.stdout)["evaluation"].(map[string]any)
}

func (fixture linkedEvidenceFixture) evaluateGolden(t *testing.T) string {
	t.Helper()
	payload, err := journal.MarshalCanonicalJSON(fixture.evaluate(t))
	if err != nil {
		t.Fatalf("canonicalize Evaluation: %v", err)
	}
	return string(payload)
}

func (fixture linkedEvidenceFixture) updateValidity(t *testing.T, validity string) {
	t.Helper()
	fixture.updateProjection(t, `UPDATE evidence_records SET validity = ? WHERE id = ?`, validity, fixture.evidenceID)
}

func (fixture linkedEvidenceFixture) updateProjection(t *testing.T, query string, arguments ...any) {
	t.Helper()
	conn, err := fiddb.Open(filepath.Join(config.LedgerDir(fixture.workspace), config.DatabaseFileName))
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
	path := filepath.Join(config.LedgerDir(fixture.workspace), "artifacts", strings.TrimPrefix(fixture.digest, "sha256:"))
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
	workspace, issueID, envelopeID := gateReadyEvidenceEnvelopeForTest(t, func(definition map[string]any) {
		definition["argv"] = []string{"/bin/true"}
	})
	terminal, gateRun := executeGateRunForTest(t, workspace, issueID, envelopeID, "caller-proof-source")
	if gateRun.err != nil || terminal.Data.ReceiptID == "" {
		t.Fatalf("gate run = %#v, terminal=%#v, want successful Receipt", gateRun, terminal)
	}
	obligationID := testEvidenceObligationID
	for _, proof := range []string{"pass", "command", "import", "attestation", "role"} {
		t.Run(proof, func(t *testing.T) {
			path := filepath.Join(workspace, proof+"-evidence.json")
			input := map[string]any{"obligation_id": obligationID, "evidence_type": "receipt", "receipt_id": terminal.Data.ReceiptID, "review_id": nil, "human_outcome_id": nil, "required_artifact_ids": []string{}}
			input[proof] = true
			writeActivationInput(t, path, input)
			before := snapshotWorkspacePersistence(t, workspace)
			result := executeCommandForTest(workspace, "--project", workspace, "--json", "evidence", "link", "--actor", "actor-builder", "--input", path, "--idempotency-key", "caller-"+proof)
			assertEvidenceLinkRejectedWithoutEvidenceEvent(t, workspace, result, before)

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
			workspace, issueID, envelopeID := gateReadyEvidenceEnvelopeForTest(t, func(definition map[string]any) {
				definition["argv"] = []string{"/bin/true"}
			})
			terminal, gateRun := executeGateRunForTest(t, workspace, issueID, envelopeID, "binding-mismatch-"+strings.ReplaceAll(testCase.name, " ", "-"))
			if gateRun.err != nil || terminal.Data.ReceiptID == "" {
				t.Fatalf("gate run = %#v, terminal=%#v, want successful Receipt", gateRun, terminal)
			}

			conn, err := fiddb.Open(filepath.Join(config.LedgerDir(workspace), config.DatabaseFileName))
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

			repository := runRepositoryForTest(t, workspace)
			if testCase.mutateLink != nil {
				testCase.mutateLink(t, repository)
			}
			linkPath := filepath.Join(workspace, "binding-mismatch-evidence.json")
			writeActivationInput(t, linkPath, map[string]any{
				"obligation_id": testEvidenceObligationID, "evidence_type": "receipt", "receipt_id": terminal.Data.ReceiptID, "review_id": nil, "human_outcome_id": nil,
				"required_artifact_ids": []string{},
			})
			if testCase.mutateLink != nil {
				request := decodeCommandJSONStructure(t, string(mustReadCommandFile(t, linkPath))).(map[string]any)
				request["target_repository"] = repository
				writeActivationInput(t, linkPath, request)
			}
			before := snapshotWorkspacePersistence(t, workspace)
			result := executeCommandForTest(workspace,
				"--project", workspace, "--json", "evidence", "link", "--actor", "actor-builder",
				"--input", linkPath, "--idempotency-key", "binding-mismatch-"+strings.ReplaceAll(testCase.name, " ", "-"),
			)
			assertEvidenceLinkRejectedWithoutEvidenceEvent(t, workspace, result, before)
			if testCase.mutateLink != nil {
				failure := decodeCommandJSONStructure(t, result.stderr).(map[string]any)
				want := map[string]any{"contract_version": "fid-cli.v1", "ok": false, "operation": "evidence.link",
					"error": map[string]any{"code": "invalid-argument", "category": "invalid-input", "retryable": false, "subject": nil,
						"message": "invalid input.target_repository", "details": map[string]any{
							"argument": "input.target_repository", "value": "target_repository",
							"allowed": []any{"evidence_type", "human_outcome_id", "obligation_id", "receipt_id", "required_artifact_ids", "review_id"}}}}
				if cli.ExitCode(result.err) != 2 || !reflect.DeepEqual(failure, want) {
					t.Fatalf("removed repository target = %#v, want %#v with exit 2", failure, want)
				}
			}
		})
	}
}

func observedEvidenceEnvelopeForTest(t *testing.T, modifyDefinition func(map[string]any)) (string, string, string) {
	return evidenceEnvelopeForTest(t, modifyDefinition, true)
}

func evidenceEnvelopeForTest(t *testing.T, modifyDefinition func(map[string]any), observedRepository bool) (string, string, string) {
	t.Helper()
	workspace, issueID, envelopeID, registration := buildEvidenceEnvelopeForTest(t, gateDefinitionsForTest(modifyDefinition), observedRepository)
	rememberEnvelopeBindingForTest(t, workspace, envelopeBindingForTest{
		registrationID: registration.ID, registrationVersion: registration.Version, obligationID: testEvidenceObligationID,
	})
	return workspace, issueID, envelopeID
}

func gateReadyEvidenceEnvelopeForTest(t *testing.T, modifyDefinition func(map[string]any)) (string, string, string) {
	return readyEvidenceEnvelopeForTest(t, modifyDefinition, false)
}

func observedGateReadyEvidenceEnvelopeForTest(t *testing.T, modifyDefinition func(map[string]any)) (string, string, string) {
	return readyEvidenceEnvelopeForTest(t, modifyDefinition, true)
}

func readyEvidenceEnvelopeForTest(t *testing.T, modifyDefinition func(map[string]any), observedRepository bool) (string, string, string) {
	t.Helper()
	workspace, issueID, envelopeID, registration := buildEvidenceEnvelopeForTest(t, gateDefinitionsForTest(modifyDefinition), observedRepository)
	claim, _ := acquireGenericClaimForTestWithTTL(t, workspace, issueID, envelopeID, "gate-ready-claim", 86400)
	repository := testEnvelopeRepository(t, workspace)
	path := filepath.Join(workspace, "gate-ready-run.json")
	writeActivationInput(t, path, runStartInput(issueID, envelopeID, claim["id"], repository))
	started := executeCommandForTest(workspace, "--project", workspace, "--json", "run", "start", "--actor", "actor-builder", "--input", path, "--idempotency-key", "gate-ready-run")
	if started.err != nil {
		t.Fatalf("start gate-ready Run: %v; stderr=%q", started.err, started.stderr)
	}
	runID := runFromResponse(t, started.stdout)["id"].(string)
	rememberEnvelopeBindingForTest(t, workspace, envelopeBindingForTest{
		registrationID: registration.ID, registrationVersion: registration.Version,
		obligationID: testEvidenceObligationID, runID: runID, claimID: claim["id"].(string),
	})
	return workspace, issueID, envelopeID
}

func buildEvidenceEnvelopeForTest(t *testing.T, definitions []any, observedRepository bool) (string, string, string, gateRegistrationView) {
	t.Helper()
	var workspace, issueID string
	var registration gateRegistrationView
	if observedRepository {
		workspace, issueID, registration = finalizedObservedEnvelopeWorkspaceWithDefinitions(t, definitions)
	} else {
		workspace, issueID, registration = finalizedEnvelopeWorkspaceWithDefinitions(t, definitions)
	}
	input := envelopeCreateInput(t, workspace, issueID, registration)
	grant := input["authority_grants"].([]any)[0].(map[string]any)
	grant["actions"] = []string{"claim.acquire", "claim.release", "envelope.activate", "envelope.complete", "evidence.link", "gate.run", "issue.transition", "run.close", "run.start"}
	grant["actor_selector"] = map[string]any{"actor_id": "actor-builder", "required_capabilities": []string{"claim.acquire", "claim.release", "envelope.activate", "envelope.complete", "evidence.link", "gate.run", "issue.transition", "run.close", "run.start"}}
	input["authority_grants"] = append(input["authority_grants"].([]any), map[string]any{"actor_selector": map[string]any{"actor_id": "actor-harness", "required_capabilities": []string{"run.observe"}}, "actions": []string{"run.observe"}, "scope": map[string]any{"issue_id": issueID}})
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
	return workspace, issueID, envelopeID, registration
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

func assertEvidenceLinkRejectedWithoutEvidenceEvent(t *testing.T, workspace string, result commandExecutionResult, want workspacePersistenceSnapshot) {
	t.Helper()
	if result.err == nil || result.stdout != "" || !strings.Contains(result.stderr, `"code":"invalid-argument"`) {
		t.Fatalf("evidence link = %#v, want typed rejection", result)
	}
	got := snapshotWorkspacePersistence(t, workspace)
	if !got.equal(want) {
		t.Fatalf("Evidence rejection changed persisted projection or journal; events=%d, want %d", got.eventCount(), want.eventCount())
	}
}
