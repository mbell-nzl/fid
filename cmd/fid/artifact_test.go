package main

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mbell-nzl/fid/internal/config"
	fiddb "github.com/mbell-nzl/fid/internal/db"
	"github.com/mbell-nzl/fid/internal/journal"
)

func TestGateRunPublishesSanitizedArtifactByDigest(t *testing.T) {
	workspace, issueID, envelopeID := runEnvelopeForTestWithDefinition(t, func(definition map[string]any) {
		definition["argv"] = []string{"/bin/true"}
		definition["artifact_paths"] = []string{"artifact.txt"}
		definition["redact_literals"] = []string{"SECRET-artifact"}
	})
	source := []byte("safe SECRET-artifact content\n")
	if err := os.WriteFile(filepath.Join(workspace, "artifact.txt"), source, 0o600); err != nil {
		t.Fatalf("write declared Artifact: %v", err)
	}

	terminal, result := executeGateRunForTest(t, workspace, issueID, envelopeID, "artifact-safe")
	if result.err != nil {
		t.Fatalf("gate run: %v; stderr=%q", result.err, result.stderr)
	}
	if terminal.Data.ReceiptID == "" {
		t.Fatal("gate run Receipt ID is empty")
	}

	want := []byte("safe [REDACTED] content\n")
	digest := sha256.Sum256(want)
	path := filepath.Join(workspace, ".fid", "artifacts", hex.EncodeToString(digest[:]))
	stored, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read content-addressed Artifact %q: %v", path, err)
	}
	if string(stored) != string(want) {
		t.Fatalf("stored Artifact = %q, want sanitized %q", stored, want)
	}
	conn, err := fiddb.Open(filepath.Join(workspace, ".fid", config.DatabaseFileName))
	if err != nil {
		t.Fatalf("open database: %v", err)
	}
	defer func() { _ = conn.Close() }()
	var recordedDigest, mediaType, retentionClass, availability, locationPolicy, declaredPath string
	var sizeBytes int64
	if err := conn.QueryRow(`SELECT a.digest, a.media_type, a.size_bytes, a.retention_class, a.availability, a.location_policy, ra.declared_path FROM referenced_artifacts a JOIN receipt_artifacts ra ON ra.artifact_id = a.id WHERE ra.receipt_id = ?`, terminal.Data.ReceiptID).Scan(&recordedDigest, &mediaType, &sizeBytes, &retentionClass, &availability, &locationPolicy, &declaredPath); err != nil {
		t.Fatalf("read recorded Artifact metadata: %v", err)
	}
	if recordedDigest != "sha256:"+hex.EncodeToString(digest[:]) || mediaType != "text/plain; charset=utf-8" || sizeBytes != int64(len(want)) || retentionClass != "evidence" || availability != "available" || locationPolicy != "content-addressed" || declaredPath != "artifact.txt" {
		t.Fatalf("recorded Artifact metadata = digest=%q media=%q size=%d retention=%q availability=%q location=%q path=%q", recordedDigest, mediaType, sizeBytes, retentionClass, availability, locationPolicy, declaredPath)
	}
	entries, err := os.ReadDir(filepath.Join(workspace, ".fid", "artifacts"))
	if err != nil {
		t.Fatalf("read Artifact store: %v", err)
	}
	for _, entry := range entries {
		if len(entry.Name()) >= len(".artifact-") && entry.Name()[:len(".artifact-")] == ".artifact-" {
			t.Fatalf("partial Artifact staging file remains: %q", entry.Name())
		}
	}
}

func TestArtifactVerifyObservesLossAndRestorationWithoutIdentityMutation(t *testing.T) {
	workspace, issueID, envelopeID := runEnvelopeForTestWithDefinition(t, func(definition map[string]any) {
		definition["argv"] = []string{"/bin/true"}
		definition["artifact_paths"] = []string{"artifact.txt"}
	})
	content := []byte("Artifact verification content\n")
	if err := os.WriteFile(filepath.Join(workspace, "artifact.txt"), content, 0o600); err != nil {
		t.Fatalf("write declared Artifact: %v", err)
	}
	_, result := executeGateRunForTest(t, workspace, issueID, envelopeID, "artifact-verify-source")
	if result.err != nil {
		t.Fatalf("gate run: %v; stderr=%q", result.err, result.stderr)
	}

	conn, err := fiddb.Open(filepath.Join(workspace, ".fid", config.DatabaseFileName))
	if err != nil {
		t.Fatalf("open database: %v", err)
	}
	defer func() { _ = conn.Close() }()
	artifactID, digest, receiptID := receiptArtifactIdentity(t, conn)
	before := artifactIdentitySnapshot(t, conn, artifactID, receiptID)
	path := filepath.Join(workspace, ".fid", "artifacts", digest[len("sha256:"):])
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
	if afterMissing := artifactIdentitySnapshot(t, conn, artifactID, receiptID); afterMissing != before {
		t.Fatalf("missing observation changed Artifact/Receipt identity = %#v, want %#v", afterMissing, before)
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
	if afterRestored := artifactIdentitySnapshot(t, conn, artifactID, receiptID); afterRestored != before {
		t.Fatalf("restored observation changed Artifact/Receipt identity = %#v, want %#v", afterRestored, before)
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
	event, err := journal.FindEvent(filepath.Join(workspace, ".fid", "journal"), unchangedEnvelope.Event.ID)
	if err != nil {
		t.Fatalf("read unchanged Artifact Event: %v", err)
	}
	if event.Event.EventType != "artifact.availability-observed" || len(event.Event.Effects) != 1 || event.Event.Effects[0].Table != "idempotency_results" {
		t.Fatalf("unchanged Artifact Event = %#v, want only idempotency observation effect", event.Event)
	}
}

func TestArtifactVerifyAuthorityFailureNamesArtifactVerify(t *testing.T) {
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

type artifactSnapshot struct {
	ArtifactID string
	Digest     string
	ReceiptID  string
}

func receiptArtifactIdentity(t *testing.T, conn *sql.DB) (string, string, string) {
	t.Helper()
	var artifactID, digest, receiptID string
	if err := conn.QueryRow(`SELECT a.id, a.digest, ra.receipt_id FROM referenced_artifacts a JOIN receipt_artifacts ra ON ra.artifact_id = a.id`).Scan(&artifactID, &digest, &receiptID); err != nil {
		t.Fatalf("read Artifact identity: %v", err)
	}
	return artifactID, digest, receiptID
}

func artifactIdentitySnapshot(t *testing.T, conn *sql.DB, artifactID, receiptID string) artifactSnapshot {
	t.Helper()
	var snapshot artifactSnapshot
	if err := conn.QueryRow(`SELECT id, digest FROM referenced_artifacts WHERE id = ?`, artifactID).Scan(&snapshot.ArtifactID, &snapshot.Digest); err != nil {
		t.Fatalf("read Artifact snapshot: %v", err)
	}
	if err := conn.QueryRow(`SELECT id FROM execution_receipts WHERE id = ?`, receiptID).Scan(&snapshot.ReceiptID); err != nil {
		t.Fatalf("read Receipt snapshot: %v", err)
	}
	return snapshot
}

func assertArtifactVerifyResult(t *testing.T, output, availability string, changed bool) {
	t.Helper()
	var envelope struct {
		Result struct {
			Artifact struct {
				Availability string `json:"availability"`
			} `json:"artifact"`
			Changed bool `json:"changed"`
		} `json:"result"`
	}
	if err := json.Unmarshal([]byte(output), &envelope); err != nil {
		t.Fatalf("decode Artifact verify output: %v; output=%q", err, output)
	}
	if envelope.Result.Artifact.Availability != availability || envelope.Result.Changed != changed {
		t.Fatalf("Artifact verify = %#v, want availability=%q changed=%t", envelope.Result, availability, changed)
	}
}
