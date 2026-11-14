//go:build integration

package integration

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/mbell-nzl/fid/internal/evidence"
)

type p23HumanStep struct {
	Name         string   `json:"name"`
	CWD          string   `json:"cwd"`
	Args         []string `json:"args"`
	ExpectedExit int      `json:"expected_exit"`
}

type p23HumanCheckpoint struct {
	Name          string               `json:"name"`
	Project       string               `json:"project"`
	IDs           map[string]string    `json:"ids"`
	Inspections   []p22HumanInspection `json:"inspections"`
	ExpectedHuman map[string]string    `json:"expected_human"`
	Steps         []p23HumanStep       `json:"steps,omitempty"`
}

type p23HumanManifest struct {
	Version      int                  `json:"version"`
	GeneratedAt  string               `json:"generated_at"`
	SourceCommit string               `json:"source_commit"`
	SourceDirty  bool                 `json:"source_dirty"`
	Binary       string               `json:"binary"`
	BinarySHA256 string               `json:"binary_sha256"`
	Checkpoints  []p23HumanCheckpoint `json:"checkpoints"`
	LeafPairs    []p23HumanLeafPair   `json:"leaf_pairs"`
}

type p23HumanLeafPair struct {
	Operation     string         `json:"operation"`
	CWD           string         `json:"cwd"`
	HumanArgs     []string       `json:"human_args"`
	JSONArgs      []string       `json:"json_args"`
	ExpectedHuman string         `json:"expected_human"`
	JSONFacts     map[string]any `json:"json_facts"`
}

func TestP23HumanFixtures(t *testing.T) {
	t.Parallel()
	if root := os.Getenv("FID_P23_HUMAN_FIXTURE_CHILD"); root != "" {
		p23GenerateHumanFixtures(t, root)
		return
	}
	root := os.Getenv("FID_P23_HUMAN_FIXTURE_ROOT")
	if root == "" {
		root = t.TempDir()
	}
	if err := os.MkdirAll(root, 0o700); err != nil {
		t.Fatal(err)
	}
	fresh, err := os.MkdirTemp(root, "checkpoints-")
	if err != nil {
		t.Fatal(err)
	}
	fresh, err = filepath.Abs(fresh)
	if err != nil {
		t.Fatal(err)
	}
	executable, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	child := exec.Command(executable, "-test.run=^TestP23HumanFixtures$", "-test.v")
	child.Env = append(os.Environ(), "FID_P23_HUMAN_FIXTURE_CHILD="+fresh)
	child.Dir = integrationRepositoryRoot()
	if output, err := child.CombinedOutput(); err != nil {
		t.Fatalf("fixture generator child: %v\n%s", err, output)
	}
	manifestPath := filepath.Join(fresh, "manifest.json")
	data, err := os.ReadFile(manifestPath)
	if err != nil {
		t.Fatal(err)
	}
	var manifest p23HumanManifest
	if err := json.Unmarshal(data, &manifest); err != nil {
		t.Fatal(err)
	}
	binary, err := os.ReadFile(manifest.Binary)
	if err != nil || fmt.Sprintf("sha256:%x", sha256.Sum256(binary)) != manifest.BinarySHA256 {
		t.Fatalf("retained binary changed: %v", err)
	}
	if len(manifest.Checkpoints) != 3 {
		t.Fatalf("checkpoint count = %d, want 3", len(manifest.Checkpoints))
	}
	wantLeaves := []string{"review.record", "review.get", "review.list", "finding.address", "finding.list", "concern.record", "concern.dispose", "concern.list", "deviation.record", "deviation.list", "decision.record-reference", "decision.list", "human-outcome.record", "human-outcome.list", "requirement.record-reference", "requirement.link", "requirement.get", "requirement.list"}
	seenLeaves := map[string]bool{}
	for _, pair := range manifest.LeafPairs {
		if seenLeaves[pair.Operation] {
			t.Fatalf("duplicate Human/JSON leaf pair %s", pair.Operation)
		}
		seenLeaves[pair.Operation] = true
	}
	for _, operation := range wantLeaves {
		if !seenLeaves[operation] {
			t.Errorf("missing Human/JSON leaf pair %s", operation)
		}
	}
	if len(manifest.LeafPairs) != len(wantLeaves) {
		t.Errorf("Human/JSON leaf pair count = %d, want %d", len(manifest.LeafPairs), len(wantLeaves))
	}
	for _, pair := range manifest.LeafPairs {
		if pair.CWD == "" || len(pair.HumanArgs) == 0 || len(pair.JSONArgs) == 0 || pair.ExpectedHuman == "" || len(pair.JSONFacts) == 0 {
			t.Fatalf("incomplete Human/JSON leaf pair %s", pair.Operation)
		}
		human, stderr, err := closureRunCommand(pair.CWD, manifest.Binary, pair.HumanArgs...)
		if err != nil || stderr != "" || human != pair.ExpectedHuman {
			t.Fatalf("post-exit %s Human: error=%v stderr=%q got=%q want=%q", pair.Operation, err, stderr, human, pair.ExpectedHuman)
		}
		jsonText, stderr, err := closureRunCommand(pair.CWD, manifest.Binary, pair.JSONArgs...)
		if err != nil || stderr != "" {
			t.Fatalf("post-exit %s JSON: error=%v stderr=%q", pair.Operation, err, stderr)
		}
		var envelope map[string]any
		if err := json.Unmarshal([]byte(jsonText), &envelope); err != nil {
			t.Fatal(err)
		}
		facts, ok := envelope["result"].(map[string]any)
		if envelope["operation"] != pair.Operation || !ok || !reflect.DeepEqual(facts, pair.JSONFacts) {
			t.Fatalf("post-exit %s JSON facts drifted: got=%#v want=%#v", pair.Operation, facts, pair.JSONFacts)
		}
		if want := p23ExpectedLeafText(t, pair.Operation, facts); human != want {
			t.Fatalf("post-exit %s Human/JSON fact mismatch: got=%q want=%q", pair.Operation, human, want)
		}
		p24CheckProjectContentSentinel(t, pair.Operation, human)
		p23CheckJudgmentQualifiers(t, pair.Operation, human)
	}
	for _, checkpoint := range manifest.Checkpoints {
		if checkpoint.Project == "" || len(checkpoint.IDs) == 0 || len(checkpoint.Inspections) < 2 {
			t.Fatalf("incomplete checkpoint: %#v", checkpoint)
		}
		for _, inspection := range checkpoint.Inspections {
			stdout, stderr, err := closureRunCommand(inspection.CWD, manifest.Binary, inspection.Args...)
			if err != nil || stderr != "" || strings.TrimSpace(stdout) == "" {
				t.Fatalf("post-exit inspection %s %v: %v stdout=%q stderr=%q", checkpoint.Name, inspection.Args, err, stdout, stderr)
			}
			if !containsAction(inspection.Args, "--json") {
				p24CheckProjectContentSentinel(t, checkpoint.Name, stdout)
			}
			if !containsAction(inspection.Args, "--json") && stdout != checkpoint.ExpectedHuman[strings.Join(inspection.Args, "\x00")] {
				t.Fatalf("post-exit Human text changed for %s %v", checkpoint.Name, inspection.Args)
			}
		}
	}
	t.Logf("Retained Human fixture manifest: %s", manifestPath)
}

func TestP23HumanFixtureStepsRemainRunnable(t *testing.T) {
	t.Parallel()
	binary := integrationFIDBinary(t)
	root := t.TempDir()
	checkpoints := []p23HumanCheckpoint{
		p23ProtectedDecisionFixture(t, binary, root),
		p23ReviewExhaustionFixture(t, binary, root),
	}
	count := 0
	for _, checkpoint := range checkpoints {
		for _, step := range checkpoint.Steps {
			stdout, stderr, err := closureRunCommand(step.CWD, binary, step.Args...)
			exit := 0
			if err != nil {
				exit = processExitCode(t, err)
			}
			if exit != step.ExpectedExit || stdout == "" && stderr == "" {
				t.Fatalf("%s / %s: exit=%d want=%d stdout=%q stderr=%q", checkpoint.Name, step.Name, exit, step.ExpectedExit, stdout, stderr)
			}
			count++
		}
	}
	if count != 6 {
		t.Fatalf("Human fixture steps exercised = %d, want 6", count)
	}
}

func p23GenerateHumanFixtures(t *testing.T, root string) {
	t.Helper()
	existing, err := os.ReadDir(root)
	if err != nil || len(existing) != 0 {
		t.Fatalf("generator requires fresh empty directory: %v", err)
	}
	contents, err := os.ReadFile(integrationFIDBinary(t))
	if err != nil {
		t.Fatal(err)
	}
	binary := filepath.Join(root, "fid")
	if err := os.WriteFile(binary, contents, 0o700); err != nil {
		t.Fatal(err)
	}
	commit, _ := closureRun(t, integrationRepositoryRoot(), "git", "rev-parse", "HEAD")
	status, _ := closureRun(t, integrationRepositoryRoot(), "git", "status", "--porcelain")
	manifest := p23HumanManifest{Version: 1, GeneratedAt: time.Now().UTC().Format(time.RFC3339Nano), SourceCommit: strings.TrimSpace(commit), SourceDirty: status != "", Binary: binary, BinarySHA256: fmt.Sprintf("sha256:%x", sha256.Sum256(contents))}
	manifest.Checkpoints = append(manifest.Checkpoints, p23MixedFindingFixture(t, binary, root))
	manifest.Checkpoints = append(manifest.Checkpoints, p23ProtectedDecisionFixture(t, binary, root))
	manifest.Checkpoints = append(manifest.Checkpoints, p23ReviewExhaustionFixture(t, binary, root))
	manifest.LeafPairs = p23GenerateLeafPairs(t, binary, root)
	for index := range manifest.Checkpoints {
		checkpoint := &manifest.Checkpoints[index]
		checkpoint.ExpectedHuman = map[string]string{}
		for _, inspection := range checkpoint.Inspections {
			if containsAction(inspection.Args, "--json") {
				continue
			}
			stdout, stderr, err := closureRunCommand(inspection.CWD, binary, inspection.Args...)
			if err != nil || stderr != "" {
				t.Fatalf("seed Human text %s: %v %q", checkpoint.Name, err, stderr)
			}
			p24CheckProjectContentSentinel(t, checkpoint.Name, stdout)
			checkpoint.ExpectedHuman[strings.Join(inspection.Args, "\x00")] = stdout
		}
	}
	data, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "manifest.json"), append(data, '\n'), 0o600); err != nil {
		t.Fatal(err)
	}
}

func p23MixedFindingFixture(t *testing.T, binary, root string) p23HumanCheckpoint {
	t.Helper()
	workspace := p22HumanDirectory(t, root, "mixed-findings")
	actions := []string{"claim.acquire", "claim.release", "envelope.activate", "envelope.create", "finding.address", "gate.register", "review.record", "run.close", "run.observe", "run.start"}
	a := p22NewFIDRecoveryAdapterWithClaimModeAt(t, binary, workspace, "agent", 3600, []string{}, actions, "mutation")
	p23CloseAndRelease(t, a, a.actor, a.runID, a.claimID, "fixture-before-review")
	claim := p23AcquireClaim(t, a, "actor-denied", "review", "fixture-first-review")
	run := p23StartRun(t, a, "actor-denied", claim, "fixture-first-review")
	first := p23ReviewRecord(t, a, "actor-denied", "fixture-first-review", p23ReviewInput(a, p21PathString(t, run, "id"), p23Number(t, run["version"]), "changes-required", []any{
		p23Finding(a.issueID, "FND-fixture-addressed", "major", true), p23Finding(a.issueID, "FND-fixture-blocking", "major", true), p23Finding(a.issueID, "FND-fixture-resolved", "major", true),
	}, nil))
	p23AssertReviewSuccess(t, "fixture first Review", first)
	firstID := p23ReviewFromResponse(t, first)["id"].(string)
	p23CloseAndRelease(t, a, "actor-denied", p21PathString(t, run, "id"), p21PathString(t, claim, "id"), "fixture-first-review")
	mutation := p23AcquireClaim(t, a, a.actor, "mutation", "fixture-address")
	address := p23FindingAddress(t, a, a.actor, "fixture-address", map[string]any{"finding_id": "FND-fixture-addressed", "claim_id": mutation["id"], "expected_claim_version": mutation["version"], "reason": "Builder says correction was made"})
	p23AssertReviewSuccess(t, "fixture addressed Finding", address)
	p23ReleaseClaim(t, a, a.actor, mutation, "fixture-address")
	secondClaim := p23AcquireClaim(t, a, "actor-reviewer-two", "review", "fixture-second-review")
	secondRun := p23StartRun(t, a, "actor-reviewer-two", secondClaim, "fixture-second-review")
	second := p23ReviewRecord(t, a, "actor-reviewer-two", "fixture-second-review", p23ReviewInput(a, p21PathString(t, secondRun, "id"), p23Number(t, secondRun["version"]), "changes-required", []any{p23Finding(a.issueID, "FND-fixture-new", "major", true)}, []any{"FND-fixture-resolved"}))
	p23AssertReviewSuccess(t, "fixture second Review", second)
	return p23HumanCheckpoint{Name: "VAL-070 mixed Findings", Project: workspace, IDs: map[string]string{"issue_id": a.issueID, "first_review_id": firstID, "second_review_id": p23ReviewFromResponse(t, second)["id"].(string)}, Inspections: p22HumanInspections(workspace,
		[]string{"review", "get", firstID}, []string{"review", "list", "--issue", a.issueID}, []string{"finding", "list", "--issue", a.issueID}, []string{"resume", a.issueID, "--repository", workspace})}
}

func p23ProtectedDecisionFixture(t *testing.T, binary, root string) p23HumanCheckpoint {
	t.Helper()
	workspace := p22HumanDirectory(t, root, "protected-decision")
	actions := []string{"claim.acquire", "claim.release", "decision.record-reference", "envelope.activate", "envelope.create", "evidence.link", "gate.register", "gate.run", "human-outcome.record", "issue.transition", "run.close", "run.observe", "run.observe-budget", "run.start"}
	decisions := []any{map[string]any{"ordinal": 0, "decision_class": "security", "required_capabilities": []string{"decision.record-reference"}, "affected_actions": []string{"issue.transition-terminal-success"}}}
	a := p22NewFIDRecoveryAdapterAt(t, binary, workspace, "agent", 3600, []string{}, actions, p22ExecutionFixture{ProtectedDecisions: decisions})
	source := p24SeedJudgmentSource(t, a)
	for revision, state := range []string{"building", "dev_complete", "qa_review"} {
		assertP23DecisionCommandSuccess(t, "fixture transition "+state, p23IssueTransition(t, a, state, revision+1, "fixture-"+state))
	}
	p22AcceptedCall(t, a, "fixture-budget", "run.observe-budget", map[string]any{"value": 0})
	gate, err := evidence.RunGate(context.Background(), p22BudgetGateOptions(t, a, "fixture-proof"))
	if err != nil || gate.Receipt == nil {
		t.Fatalf("fixture gate: %#v %v", gate, err)
	}
	linked := p23RunInput(t, a, "evidence", "link", a.actor, "fixture-proof-link", map[string]any{"obligation_id": "OBL-p22-corpus", "evidence_type": "receipt", "receipt_id": gate.Receipt.ID, "review_id": nil, "human_outcome_id": nil, "required_artifact_ids": []string{}})
	assertP23DecisionCommandSuccess(t, "fixture evidence", linked)
	reference := p23RecordDecision(t, a, "fixture-security", "security", "approved", []string{"decision.record-reference"}, source)
	assertP23DecisionSuccess(t, reference, 0)
	blockedArgs := []string{"--project", workspace, "--json", "issue", "transition", a.issueID, "--to", "qa_pass", "--actor", "agent:" + a.actor, "--expected-revision", "4", "--idempotency-key", "fixture-human-transition"}
	stdout, stderr, err := closureRunCommand(workspace, binary, blockedArgs...)
	if err == nil || stdout != "" || !strings.Contains(stderr, "protected-decision-required") {
		t.Fatalf("fixture protected denial: %v %q %q", err, stdout, stderr)
	}
	input := p21Input(t, workspace, "fixture-human-outcome", map[string]any{"id": "OUT-fixture-security", "class": "security", "subject": map[string]any{"type": "issue", "id": a.issueID}, "result": "approved", "scope": map[string]any{"issue_id": a.issueID}, "reasoning": "Human approves this guarded action", "references": []any{}})
	recordArgs := []string{"--project", workspace, "--json", "human-outcome", "record", "--input", input, "--actor", "actor-human-outcome", "--idempotency-key", "fixture-human-outcome"}
	return p23HumanCheckpoint{Name: "VAL-071 protected Decision", Project: workspace, IDs: map[string]string{"issue_id": a.issueID, "envelope_id": p23EnvelopeID(t, a), "human_actor_id": "actor-human-outcome"}, Inspections: p22HumanInspections(workspace,
		[]string{"resume", a.issueID, "--repository", workspace}, []string{"decision", "list", "--issue", a.issueID}, []string{"stop", "list", "--subject", "envelope:" + p23EnvelopeID(t, a)}),
		Steps: []p23HumanStep{{Name: "blocked action before Outcome", CWD: workspace, Args: blockedArgs, ExpectedExit: 5}, {Name: "record authorized Human Outcome", CWD: workspace, Args: recordArgs, ExpectedExit: 0}, {Name: "retry guarded action", CWD: workspace, Args: blockedArgs, ExpectedExit: 0}}}
}

func p23ReviewExhaustionFixture(t *testing.T, binary, root string) p23HumanCheckpoint {
	t.Helper()
	workspace := p22HumanDirectory(t, root, "review-exhaustion")
	actions := []string{"claim.acquire", "claim.release", "envelope.activate", "envelope.create", "gate.register", "human-outcome.record", "review.record", "run.close", "run.observe", "run.start"}
	a := p22NewFIDRecoveryAdapterWithClaimModeAt(t, binary, workspace, "agent", 3600, []string{}, actions, "mutation", p22ExecutionFixture{ReviewObligation: true, ReviewCyclesMaximum: 1})
	p23CloseAndRelease(t, a, a.actor, a.runID, a.claimID, "fixture-builder")
	claim := p23AcquireClaim(t, a, "actor-denied", "review", "fixture-cycle")
	run := p23StartRun(t, a, "actor-denied", claim, "fixture-cycle")
	recorded := p23ReviewRecord(t, a, "actor-denied", "fixture-cycle", p23ReviewInput(a, p21PathString(t, run, "id"), p23Number(t, run["version"]), "changes-required", []any{p23Finding(a.issueID, "FND-fixture-cycle", "major", true)}, nil))
	p23AssertReviewSuccess(t, "fixture exhaustion Review", recorded)
	p23CloseAndRelease(t, a, "actor-denied", p21PathString(t, run, "id"), p21PathString(t, claim, "id"), "fixture-cycle")
	input := p21Input(t, workspace, "fixture-rework-claim", map[string]any{"issue_id": a.issueID, "envelope_id": p23EnvelopeID(t, a), "envelope_revision": 1, "scope": map[string]any{"paths": []string{"base.txt"}}, "mode": "mutation", "ttl_seconds": 60})
	deniedArgs := []string{"--project", workspace, "--json", "claim", "acquire", "--input", input, "--actor", a.actor, "--idempotency-key", "fixture-exhausted-claim"}
	stdout, stderr, err := closureRunCommand(workspace, binary, deniedArgs...)
	if err == nil || stdout != "" || !strings.Contains(stderr, "review-rework-exhausted") {
		t.Fatalf("fixture exhaustion denial: %v %q %q", err, stdout, stderr)
	}
	outcomeInput := p21Input(t, workspace, "fixture-rework-outcome", map[string]any{"id": "OUT-fixture-rework", "class": "review-rework-extension", "subject": map[string]any{"type": "issue", "id": a.issueID}, "result": "approved", "scope": map[string]any{"issue_id": a.issueID}, "reasoning": "Human authorizes rework beyond the original cycle budget", "references": []any{}})
	resolveArgs := []string{"--project", workspace, "--json", "human-outcome", "record", "--input", outcomeInput, "--actor", "actor-human-outcome", "--idempotency-key", "fixture-rework-outcome"}
	return p23HumanCheckpoint{Name: "VAL-072 review rework exhaustion", Project: workspace, IDs: map[string]string{"issue_id": a.issueID, "envelope_id": p23EnvelopeID(t, a), "human_actor_id": "actor-human-outcome", "review_id": p23ReviewFromResponse(t, recorded)["id"].(string)}, Inspections: p22HumanInspections(workspace,
		[]string{"resume", a.issueID, "--repository", workspace}, []string{"stop", "list", "--subject", "envelope:" + p23EnvelopeID(t, a)}, []string{"attention", "list", "--subject", "envelope:" + p23EnvelopeID(t, a)}),
		Steps: []p23HumanStep{{Name: "exhausted rework denial", CWD: workspace, Args: deniedArgs, ExpectedExit: 4}, {Name: "record approved review-rework-extension Outcome", CWD: workspace, Args: resolveArgs, ExpectedExit: 0}, {Name: "retry rework Claim", CWD: workspace, Args: deniedArgs, ExpectedExit: 0}}}
}
