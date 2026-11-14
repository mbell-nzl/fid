package main

// RESET 1.5b2 exception: these own-command tests keep Activation authority
// unchanged until 1.5c2 switches the adapter and converts their fixtures.
import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mbell-nzl/fid/internal/config"
	fiddb "github.com/mbell-nzl/fid/internal/db"
	"github.com/mbell-nzl/fid/internal/journal"
)

func TestArtifactVerifyObservesLossAndRestorationWithoutIdentityMutation(t *testing.T) {
	t.Parallel()
	// Task 10's readiness-Stop re-evaluation (DEC-224) makes artifact.verify
	// observe the current repository whenever a flip to available reaches an
	// Obligation with linked Evidence, so this journey needs a real Git
	// worktree rather than the cheaper unobserved fixture's synthetic OIDs.
	content := []byte("Artifact verification content\n")
	workspace, issueID, envelopeID := observedGateReadyEvidenceEnvelopeWithSourceForTest(t, "artifact.txt", content, func(definition map[string]any) {
		definition["argv"] = []string{"/bin/true"}
		definition["artifact_paths"] = []string{"artifact.txt"}
	})
	terminal, result := executeGateRunForTest(t, workspace, issueID, envelopeID, "artifact-verify-source")
	if result.err != nil || terminal.Data.ReceiptID == "" {
		t.Fatalf("gate run = %#v, terminal=%#v, want successful Receipt", result, terminal)
	}

	conn, err := fiddb.Open(filepath.Join(config.LedgerDir(workspace), config.DatabaseFileName))
	if err != nil {
		t.Fatalf("open database: %v", err)
	}
	defer func() { _ = conn.Close() }()
	artifactID, digest, receiptID := receiptArtifactIdentity(t, conn)
	repository := runRepositoryForTest(t, workspace)
	obligationID := testEvidenceObligationID
	linkPath := filepath.Join(workspace, "artifact-lifecycle-evidence.json")
	writeActivationInput(t, linkPath, map[string]any{
		"obligation_id": obligationID, "evidence_type": "receipt", "receipt_id": receiptID, "review_id": nil, "human_outcome_id": nil,
		"required_artifact_ids": []string{artifactID},
	})
	linked := executeCommandForTest(workspace,
		"--project", workspace, "--json", "evidence", "link", "--actor", "actor-builder",
		"--input", linkPath, "--idempotency-key", "artifact-lifecycle-evidence-link",
	)
	if linked.err != nil || linked.stderr != "" {
		t.Fatalf("link lifecycle Evidence = %#v, want qualifying Evidence", linked)
	}
	linkedResult := decodeCanonicalResult(t, linked.stdout)
	evidenceID := linkedResult["evidence"].(map[string]any)["id"].(string)
	assertEvaluation(t, linkedResult["evaluation"].(map[string]any), "valid", "fresh", "available", []string{})
	before := artifactIdentitySnapshot(t, conn, artifactID, receiptID)
	if before.EvidenceID != evidenceID || before.EvidenceReceiptID != receiptID {
		t.Fatalf("initial Artifact/Receipt/Evidence snapshot = %#v, want Evidence %q bound to Receipt %q", before, evidenceID, receiptID)
	}
	path := filepath.Join(config.LedgerDir(workspace), "artifacts", digest[len("sha256:"):])
	if err := os.Remove(path); err != nil {
		t.Fatalf("remove stored Artifact: %v", err)
	}

	inputPath := filepath.Join(workspace, "artifact-verify.json")
	writeActivationInput(t, inputPath, map[string]any{"artifact_id": artifactID})
	missing := executeCommandForTest(workspace,
		"--project", workspace, "--json", "artifact", "verify", "--actor", "actor-builder",
		"--input", inputPath, "--idempotency-key", "artifact-verify-missing",
	)
	if missing.err != nil || missing.stderr != "" {
		t.Fatalf("missing Artifact verify = %#v, want availability observation", missing)
	}
	assertArtifactVerifyResult(t, missing.stdout, "missing", true)
	assertEvaluation(t, evaluateEvidenceForArtifactTest(t, workspace, obligationID, repository), "valid", "fresh", "missing", []string{"artifact-unavailable"})
	wantMissing := before
	wantMissing.Availability = "missing"
	if afterMissing := artifactIdentitySnapshot(t, conn, artifactID, receiptID); afterMissing != wantMissing {
		t.Fatalf("missing observation snapshot = %#v, want %#v", afterMissing, wantMissing)
	}

	if err := os.WriteFile(path, content, 0o600); err != nil {
		t.Fatalf("restore stored Artifact: %v", err)
	}
	restored := executeCommandForTest(workspace,
		"--project", workspace, "--json", "artifact", "verify", "--actor", "actor-builder",
		"--input", inputPath, "--idempotency-key", "artifact-verify-restored",
	)
	if restored.err != nil || restored.stderr != "" {
		t.Fatalf("restored Artifact verify = %#v, want availability observation", restored)
	}
	assertArtifactVerifyResult(t, restored.stdout, "available", true)
	assertEvaluation(t, evaluateEvidenceForArtifactTest(t, workspace, obligationID, repository), "valid", "fresh", "available", []string{})
	if afterRestored := artifactIdentitySnapshot(t, conn, artifactID, receiptID); afterRestored != before {
		t.Fatalf("restored observation snapshot = %#v, want %#v", afterRestored, before)
	}
	unchanged := executeCommandForTest(workspace,
		"--project", workspace, "--json", "artifact", "verify", "--actor", "actor-builder",
		"--input", inputPath, "--idempotency-key", "artifact-verify-unchanged",
	)
	if unchanged.err != nil || unchanged.stderr != "" {
		t.Fatalf("unchanged Artifact verify = %#v, want idempotent observation Event", unchanged)
	}
	assertArtifactVerifyResult(t, unchanged.stdout, "available", false)
	var unchangedEnvelope struct {
		Event struct {
			ID string `json:"id"`
		} `json:"event"`
	}
	if err := json.Unmarshal([]byte(unchanged.stdout), &unchangedEnvelope); err != nil {
		t.Fatalf("decode unchanged Artifact Event: %v", err)
	}
	event, err := journal.FindEvent(filepath.Join(config.LedgerDir(workspace), "journal"), unchangedEnvelope.Event.ID)
	if err != nil {
		t.Fatalf("read unchanged Artifact Event: %v", err)
	}
	if event.Event.EventType != "artifact.availability-observed" || len(event.Event.Effects) != 1 || event.Event.Effects[0].Table != "idempotency_results" {
		t.Fatalf("unchanged Artifact Event = %#v, want only idempotency observation effect", event.Event)
	}
}

func TestArtifactVerifyAuthorityFailureNamesArtifactVerify(t *testing.T) {
	t.Parallel()
	workspace := clonePristineWorkspace(t)
	inputPath := filepath.Join(workspace, "artifact-verify.json")
	writeActivationInput(t, inputPath, map[string]any{"artifact_id": "ART-missing"})
	result := executeCommandForTest(workspace,
		"--project", workspace, "--json", "artifact", "verify", "--actor", "actor-builder",
		"--input", inputPath, "--idempotency-key", "artifact-verify-authority",
	)
	if result.err == nil || result.stdout != "" || !strings.Contains(result.stderr, `"operation":"artifact.verify"`) || !strings.Contains(result.stderr, `"required_capabilities":["artifact.verify"]`) {
		t.Fatalf("Artifact verify authority failure = %#v, want artifact.verify capability metadata", result)
	}
}
