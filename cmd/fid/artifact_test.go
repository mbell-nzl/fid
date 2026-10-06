package main

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/mbell-nzl/fid/internal/config"
	fiddb "github.com/mbell-nzl/fid/internal/db"
	"github.com/mbell-nzl/fid/internal/journal"
)

const artifactMutationWait = 15 * time.Second

func TestGateRunPublishesSanitizedArtifactByDigest(t *testing.T) {
	t.Parallel()
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

func TestGateRunUnsafeArtifactsEmitControlsAndRemainRequiredUnavailable(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name       string
		path       string
		reason     string
		definition func(map[string]any)
		setup      func(*testing.T, string)
	}{
		{
			name:   "changed during copy",
			path:   "changed.txt",
			reason: "changed-during-copy",
			definition: func(definition map[string]any) {
				definition["argv"] = []string{os.Args[0], "-test.run=^TestArtifactMutationHelperProcess$", "--"}
				definition["artifact_bytes"] = 32 * 1024 * 1024
				definition["timeout_ms"] = 30_000
				definition["environment"] = []any{map[string]any{"name": "FID_ARTIFACT_HELPER", "source": "literal", "value": "writer"}}
			},
		},
		{
			name:   "directory",
			path:   "directory",
			reason: "directory",
			setup: func(t *testing.T, workspace string) {
				t.Helper()
				if err := os.Mkdir(filepath.Join(workspace, "directory"), 0o700); err != nil {
					t.Fatalf("create declared directory: %v", err)
				}
			},
		},
		{
			name:   "missing",
			path:   "missing.txt",
			reason: "missing",
		},
		{
			name:   "non regular",
			path:   "pipe",
			reason: "non-regular",
			setup: func(t *testing.T, workspace string) {
				t.Helper()
				if err := syscall.Mkfifo(filepath.Join(workspace, "pipe"), 0o600); err != nil {
					t.Fatalf("create declared FIFO: %v", err)
				}
			},
		},
		{
			name:   "outside cwd",
			path:   "escape/outside.txt",
			reason: "outside-cwd",
			definition: func(definition map[string]any) {
				definition["argv"] = []string{os.Args[0], "-test.run=^TestArtifactMutationHelperProcess$", "--"}
				definition["timeout_ms"] = 10_000
				definition["environment"] = []any{map[string]any{"name": "FID_ARTIFACT_HELPER", "source": "literal", "value": "outside"}}
			},
		},
		{
			name:   "oversized",
			path:   "oversized.txt",
			reason: "oversized",
			definition: func(definition map[string]any) {
				definition["artifact_bytes"] = 1
			},
			setup: func(t *testing.T, workspace string) {
				t.Helper()
				if err := os.WriteFile(filepath.Join(workspace, "oversized.txt"), []byte("too large"), 0o600); err != nil {
					t.Fatalf("write oversized candidate: %v", err)
				}
			},
		},
		{
			name:   "symlink",
			path:   "symlink",
			reason: "symlink",
			setup: func(t *testing.T, workspace string) {
				t.Helper()
				target := filepath.Join(workspace, "target.txt")
				if err := os.WriteFile(target, []byte("target\n"), 0o600); err != nil {
					t.Fatalf("write symlink target: %v", err)
				}
				if err := os.Symlink(target, filepath.Join(workspace, "symlink")); err != nil {
					t.Fatalf("create declared symlink: %v", err)
				}
			},
		},
	}

	for _, testCase := range cases {
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			workspace, issueID, envelopeID := gateReadyEvidenceEnvelopeForTest(t, func(definition map[string]any) {
				definition["argv"] = []string{"/bin/true"}
				definition["artifact_paths"] = []string{testCase.path}
				if testCase.definition != nil {
					testCase.definition(definition)
				}
			})
			if testCase.setup != nil {
				testCase.setup(t, workspace)
			}

			terminal, result := executeGateRunForTest(t, workspace, issueID, envelopeID, "artifact-rejected-"+strings.ReplaceAll(testCase.name, " ", "-"))
			if result.err != nil || terminal.Data.ReceiptID == "" {
				t.Fatalf("gate run = %#v, terminal=%#v, want successful Receipt with rejected Artifact", result, terminal)
			}
			assertArtifactRejectedControl(t, result.stdout, testCase.path, testCase.reason)
			assertReceiptHasNoArtifacts(t, workspace, terminal.Data.ReceiptID)

			repository := runRepositoryForTest(t, workspace)
			obligationID := testEvidenceObligationID
			assertEvaluation(t, evaluateEvidenceForArtifactTest(t, workspace, obligationID, repository), "invalid", "stale", "missing", []string{"obligation-evidence-absent"})
			if testCase.name == "changed during copy" {
				waitForArtifactMutation(t, workspace)
			}
		})
	}
}

func TestArtifactMutationHelperProcess(t *testing.T) {
	t.Parallel()
	switch os.Getenv("FID_ARTIFACT_HELPER") {
	case "writer":
		candidate, err := os.OpenFile("changed.txt", os.O_CREATE|os.O_WRONLY, 0o600)
		if err != nil {
			t.Fatalf("create changed candidate: %v", err)
		}
		if err := candidate.Truncate(8 * 1024 * 1024); err != nil {
			_ = candidate.Close()
			t.Fatalf("create changed candidate: %v", err)
		}
		if err := candidate.Close(); err != nil {
			t.Fatalf("close changed candidate: %v", err)
		}
		mutator := exec.Command(os.Args[0], "-test.run=^TestArtifactMutationHelperProcess$", "--")
		mutator.Env = append(os.Environ(), "FID_ARTIFACT_HELPER=mutator")
		if err := mutator.Start(); err != nil {
			t.Fatalf("start Artifact mutator: %v", err)
		}
		waitForArtifactMutationMarker(t, "changed.mutator.ready", artifactMutationWait)
	case "mutator":
		defer func() { _ = os.WriteFile("changed.mutator.done", []byte("done"), 0o600) }()
		if err := os.WriteFile("changed.mutator.ready", []byte("ready"), 0o600); err != nil {
			t.Fatalf("signal Artifact mutator readiness: %v", err)
		}
		deadline := time.Now().Add(artifactMutationWait)
		for {
			staging, err := filepath.Glob(filepath.Join(".fid", "artifacts", ".artifact-*"))
			if err != nil {
				t.Fatalf("inspect Artifact staging files: %v", err)
			}
			if len(staging) > 0 {
				break
			}
			if time.Now().After(deadline) {
				t.Fatal("Artifact staging file did not appear")
			}
			time.Sleep(time.Millisecond)
		}
		if err := os.Truncate("changed.txt", 12*1024*1024); err != nil {
			t.Fatalf("grow changed candidate: %v", err)
		}
	case "outside":
		outside := filepath.Join("..", "outside-artifact")
		if err := os.MkdirAll(outside, 0o700); err != nil {
			t.Fatalf("create outside candidate directory: %v", err)
		}
		if err := os.WriteFile(filepath.Join(outside, "outside.txt"), []byte("outside\n"), 0o600); err != nil {
			t.Fatalf("write outside candidate: %v", err)
		}
		if err := os.Symlink(outside, "escape"); err != nil {
			t.Fatalf("create escaping parent symlink: %v", err)
		}
	default:
		return
	}
}

func waitForArtifactMutation(t *testing.T, workspace string) {
	t.Helper()
	waitForArtifactMutationMarker(t, filepath.Join(workspace, "changed.mutator.done"), artifactMutationWait)
}

func waitForArtifactMutationMarker(t *testing.T, path string, timeout time.Duration) {
	t.Helper()
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if _, err := os.Stat(path); err == nil {
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
	t.Fatalf("Artifact mutation marker %q did not appear", path)
}

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

	conn, err := fiddb.Open(filepath.Join(workspace, ".fid", config.DatabaseFileName))
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
	event, err := journal.FindEvent(filepath.Join(workspace, ".fid", "journal"), unchangedEnvelope.Event.ID)
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

type artifactSnapshot struct {
	ArtifactID        string
	Digest            string
	Availability      string
	ReceiptID         string
	EvidenceID        string
	EvidenceReceiptID string
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
	if err := conn.QueryRow(`SELECT id, digest, availability FROM referenced_artifacts WHERE id = ?`, artifactID).Scan(&snapshot.ArtifactID, &snapshot.Digest, &snapshot.Availability); err != nil {
		t.Fatalf("read Artifact snapshot: %v", err)
	}
	if err := conn.QueryRow(`SELECT id FROM execution_receipts WHERE id = ?`, receiptID).Scan(&snapshot.ReceiptID); err != nil {
		t.Fatalf("read Receipt snapshot: %v", err)
	}
	if err := conn.QueryRow(`SELECT id, receipt_id FROM evidence_records WHERE receipt_id = ?`, receiptID).Scan(&snapshot.EvidenceID, &snapshot.EvidenceReceiptID); err != nil {
		t.Fatalf("read Evidence snapshot: %v", err)
	}
	return snapshot
}

func evaluateEvidenceForArtifactTest(t *testing.T, workspace, obligationID string, repository map[string]any) map[string]any {
	t.Helper()
	path := filepath.Join(workspace, "artifact-lifecycle-evaluate.json")
	writeActivationInput(t, path, map[string]any{"obligation_id": obligationID})
	result := executeCommandForTest(workspace, "--project", workspace, "--json", "evidence", "evaluate", "--input", path)
	if result.err != nil || result.stderr != "" {
		t.Fatalf("evaluate lifecycle Evidence = %#v, want evaluation", result)
	}
	return decodeCanonicalResult(t, result.stdout)["evaluation"].(map[string]any)
}

func assertArtifactRejectedControl(t *testing.T, output, path, reason string) {
	t.Helper()
	for _, line := range strings.Split(strings.TrimSpace(output), "\n") {
		var frame struct {
			Stream string         `json:"stream"`
			Data   map[string]any `json:"data"`
		}
		if err := json.Unmarshal([]byte(line), &frame); err != nil {
			t.Fatalf("decode gate frame: %v; frame=%q", err, line)
		}
		if frame.Stream == "control" && frame.Data["kind"] == "artifact-rejected" && frame.Data["path"] == path && frame.Data["reason"] == reason {
			return
		}
	}
	t.Fatalf("artifact-rejected control for path=%q reason=%q is absent from %q", path, reason, output)
}

func assertReceiptHasNoArtifacts(t *testing.T, workspace, receiptID string) {
	t.Helper()
	conn, err := fiddb.Open(filepath.Join(workspace, config.WorkspaceDirName, config.DatabaseFileName))
	if err != nil {
		t.Fatalf("open projection: %v", err)
	}
	defer func() { _ = conn.Close() }()
	var count int
	if err := conn.QueryRow(`SELECT COUNT(*) FROM receipt_artifacts WHERE receipt_id = ?`, receiptID).Scan(&count); err != nil {
		t.Fatalf("count rejected Receipt Artifacts: %v", err)
	}
	if count != 0 {
		t.Fatalf("rejected Receipt Artifact count = %d, want 0", count)
	}
}

// observedGateReadyEvidenceEnvelopeWithSourceForTest builds a gate-ready
// Evidence Envelope on a real observed Git worktree, exactly like
// observedGateReadyEvidenceEnvelopeForTest, but declares sourcePath excluded
// from the observed worktree from the Envelope's very first repository
// capture onward, then writes it. Every repository-bound step downstream of
// Envelope activation - run.start, gate.run, and now, since task 10's
// artifact.availability-observed readiness re-evaluation (DEC-224),
// artifact.verify itself - compares a fresh Git observation byte-for-byte
// (including HeadOID) against a state captured earlier in the chain. Writing
// or committing sourcePath after that first capture would change every later
// observation and break that exact-match binding; declaring it ignored
// before the first capture keeps every later observation identical no matter
// when the file is written, removed, or restored.
func observedGateReadyEvidenceEnvelopeWithSourceForTest(t *testing.T, sourcePath string, content []byte, modifyDefinition func(map[string]any)) (string, string, string) {
	t.Helper()
	definitions := gateDefinitionsForTest(modifyDefinition)
	workspace, issueID, registration := finalizedObservedEnvelopeWorkspaceWithDefinitions(t, definitions)
	ignorePath := filepath.Join(workspace, ".gitignore")
	existing, err := os.ReadFile(ignorePath)
	if err != nil {
		t.Fatalf("read observed repository ignore: %v", err)
	}
	if err := os.WriteFile(ignorePath, append(existing, []byte(sourcePath+"\n")...), 0o600); err != nil {
		t.Fatalf("extend observed repository ignore: %v", err)
	}
	for _, arguments := range [][]string{{"add", ".gitignore"}, {"commit", "-qm", "ignore declared Artifact source"}} {
		if output, err := exec.Command("git", append([]string{"-C", workspace}, arguments...)...).CombinedOutput(); err != nil {
			t.Fatalf("commit observed repository ignore %v: %v (%s)", arguments, err, output)
		}
	}
	if err := os.WriteFile(filepath.Join(workspace, sourcePath), content, 0o600); err != nil {
		t.Fatalf("write declared Artifact source: %v", err)
	}

	input := envelopeCreateInput(t, workspace, issueID, registration)
	grant := input["authority_grants"].([]any)[0].(map[string]any)
	grant["actions"] = []string{"claim.acquire", "claim.release", "envelope.activate", "envelope.complete", "evidence.link", "gate.run", "issue.transition", "run.close", "run.start"}
	grant["actor_selector"] = map[string]any{"actor_id": "actor-builder", "required_capabilities": []string{"claim.acquire", "claim.release", "envelope.activate", "envelope.complete", "evidence.link", "gate.run", "issue.transition", "run.close", "run.start"}}
	input["authority_grants"] = append(input["authority_grants"].([]any), map[string]any{"actor_selector": map[string]any{"actor_id": "actor-harness", "required_capabilities": []string{"run.observe"}}, "actions": []string{"run.observe"}, "scope": map[string]any{"issue_id": issueID}})
	createPath := filepath.Join(workspace, "evidence-envelope.json")
	writeActivationInput(t, createPath, input)
	created := executeCommandForTest(workspace, "--project", workspace, "--json", "envelope", "create", "--actor", "actor-builder", "--input", createPath, "--idempotency-key", "evidence-envelope-create")
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

	claim, _ := acquireGenericClaimForTestWithTTL(t, workspace, issueID, envelopeID, "gate-ready-claim", 86400)
	repository := testEnvelopeRepository(t, workspace)
	runPath := filepath.Join(workspace, "gate-ready-run.json")
	writeActivationInput(t, runPath, runStartInput(issueID, envelopeID, claim["id"], repository))
	started := executeCommandForTest(workspace, "--project", workspace, "--json", "run", "start", "--actor", "actor-builder", "--input", runPath, "--idempotency-key", "gate-ready-run")
	if started.err != nil {
		t.Fatalf("start observed gate-ready Run: %v; stderr=%q", started.err, started.stderr)
	}
	runID := runFromResponse(t, started.stdout)["id"].(string)
	rememberEnvelopeBindingForTest(t, workspace, envelopeBindingForTest{registrationID: registration.ID, registrationVersion: registration.Version, obligationID: testEvidenceObligationID, runID: runID, claimID: claim["id"].(string)})
	return workspace, issueID, envelopeID
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
