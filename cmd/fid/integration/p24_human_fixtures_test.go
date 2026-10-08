//go:build integration

package integration

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

type p24HumanFixture struct {
	Name         string   `json:"name"`
	Operation    string   `json:"operation"`
	Code         string   `json:"code,omitempty"`
	Workspace    string   `json:"workspace"`
	CWD          string   `json:"cwd"`
	HumanArgs    []string `json:"human_args"`
	JSONArgs     []string `json:"json_args"`
	HumanCommand string   `json:"human_command"`
	JSONCommand  string   `json:"json_command"`
	ExpectedExit int      `json:"expected_exit"`
	BinarySHA256 string   `json:"binary_sha256"`
	HumanStdout  string   `json:"human_stdout"`
	HumanStderr  string   `json:"human_stderr"`
	JSONStdout   string   `json:"json_stdout"`
	JSONStderr   string   `json:"json_stderr"`
}
type p24HumanManifest struct {
	Version      int               `json:"version"`
	Binary       string            `json:"binary"`
	BinarySHA256 string            `json:"binary_sha256"`
	Fixtures     []p24HumanFixture `json:"fixtures"`
}

func TestP24HumanFixtures(t *testing.T) {
	t.Parallel()
	root := os.Getenv("FID_P24_HUMAN_FIXTURE_ROOT")
	if root == "" {
		root = t.TempDir()
	}
	if err := os.MkdirAll(root, 0700); err != nil {
		t.Fatal(err)
	}
	fresh, err := os.MkdirTemp(root, "fixtures-")
	if err != nil {
		t.Fatal(err)
	}
	fresh, err = filepath.Abs(fresh)
	if err != nil {
		t.Fatal(err)
	}
	p24GenerateHumanFixtures(t, fresh)
	data, err := os.ReadFile(filepath.Join(fresh, "manifest.json"))
	if err != nil {
		t.Fatal(err)
	}
	var manifest p24HumanManifest
	if err := json.Unmarshal(data, &manifest); err != nil {
		t.Fatal(err)
	}
	if len(manifest.Fixtures) < 21 {
		t.Fatalf("manifest has %d fixtures, want five denials and four views in four states", len(manifest.Fixtures))
	}
	p24AssertHumanManifest(t, manifest)
	t.Logf("Retained Human fixture manifest: %s", filepath.Join(fresh, "manifest.json"))
}

var p24FixtureActions = []string{"attention.resolve", "claim.acquire", "claim.expire", "claim.release", "decision.record-reference", "envelope.activate", "envelope.create", "gate.register", "gate.run", "human-outcome.record", "issue.transition", "review.record", "run.close", "run.observe", "run.observe-budget", "run.recover", "run.start", "stop.raise", "stop.resolve"}

func p24GenerateHumanFixtures(t *testing.T, root string) {
	t.Helper()
	contents, err := os.ReadFile(integrationFIDBinary(t))
	if err != nil {
		t.Fatal(err)
	}
	manifest := p24HumanManifest{Version: 1, Binary: filepath.Join(root, "fid"), BinarySHA256: fmt.Sprintf("sha256:%x", sha256.Sum256(contents))}
	if err := os.WriteFile(manifest.Binary, contents, 0700); err != nil {
		t.Fatal(err)
	}
	p24HumanDenialFixtures(t, root, &manifest)
	for _, dirty := range []bool{false, true} {
		p24HumanViewFixtures(t, root, &manifest, dirty)
	}
	data, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "manifest.json"), append(data, '\n'), 0600); err != nil {
		t.Fatal(err)
	}
}

func p24HumanDenialFixtures(t *testing.T, root string, manifest *p24HumanManifest) {
	t.Helper()
	workspace := p22HumanDirectory(t, root, "denials")
	a := p22NewFIDRecoveryAdapterAt(t, manifest.Binary, workspace, "human", 86400, []string{}, p24FixtureActions)
	p22AcceptedCall(t, a, "human-conflict", "run.observe-budget", map[string]any{"value": 0})
	fixture := p22ConformanceCase{Name: "human-conflict", Operation: "run.observe-budget"}
	args, input := a.p22BuildArgs(t, fixture)
	input["value"] = 1
	args = append(p22WithoutJSON(args), "--input", p21Input(t, workspace, "conflict", input), "--idempotency-key", "p22-corpus-human-conflict")
	p24AddHumanFixture(t, manifest, "idempotency-conflict", "run.observe-budget", "idempotency-conflict", workspace, args, 4)

	a.p22Observe(t, "terminated")
	args, input = a.p22BuildArgs(t, p22ConformanceCase{Operation: "run.close"})
	input["expected_repository"] = map[string]any{"head_oid": p22Zeros(40), "dirty": false}
	args = append(p22WithoutJSON(args), "--input", p21Input(t, workspace, "mismatch", input), "--idempotency-key", "human-mismatch")
	p24AddHumanFixture(t, manifest, "repository-mismatch", "run.close", "repository-mismatch", workspace, args, 4)

	input = p24ObservedGateInputWithoutObservation(a)
	input["expected_repository"].(map[string]any)["dirty"] = "yes"
	args = []string{"--project", workspace, "gate", "run", "--input", p21Input(t, workspace, "nested", input), "--actor", a.actor, "--idempotency-key", "human-nested"}
	p24AddHumanFixture(t, manifest, "invalid-argument", "gate.run", "invalid-argument", workspace, args, 2)

	args = []string{"--project", workspace, "requirement", "record-reference", "--input", p21Input(t, workspace, "digest", map[string]any{"id": "REQ-human", "source": p23DecisionSource(a, false), "status": "active"}), "--actor", a.actor, "--idempotency-key", "human-digest"}
	p24AddHumanFixture(t, manifest, "source-digest-mismatch", "requirement.record-reference", "source-digest-mismatch", workspace, args, 4)

	checkpoint := p23ProtectedDecisionFixture(t, manifest.Binary, root)
	p24AddHumanFixture(t, manifest, "protected-decision-required", "issue.transition", "protected-decision-required", checkpoint.Project, p22WithoutJSON(checkpoint.Steps[0].Args), 5)
}

func p24HumanViewFixtures(t *testing.T, root string, manifest *p24HumanManifest, dirty bool) {
	t.Helper()
	state := "clean"
	if dirty {
		state = "dirty"
	}
	workspace := p22HumanDirectory(t, root, state+"-same")
	a := p22NewFIDRecoveryAdapterAt(t, manifest.Binary, workspace, "human", 86400, []string{}, p24FixtureActions)
	p23CloseAndRelease(t, a, a.actor, a.runID, a.claimID, "human-builder")
	claim := p23AcquireClaim(t, a, "actor-denied", "review", "human-review")
	run := p23StartRun(t, a, "actor-denied", claim, "human-review")
	runID := p21PathString(t, run, "id")
	if dirty {
		p22Write(t, workspace, "base.txt", "uncommitted Human work\n")
	}
	a.repository = p23CurrentRepository(t, a)
	input := p23ReviewInput(a, runID, p23Number(t, run["version"]), "pass", nil, nil)
	review := p23ReviewRecord(t, a, "actor-denied", "human-review", input)
	p23AssertReviewSuccess(t, "Human fixture Review", review)
	reviewID := p23ReviewFromResponse(t, review)["id"].(string)
	p22AcceptedCall(t, a, "human-stop", "stop.raise", nil)
	copied := p22HumanDirectory(t, root, state+"-different")
	copyParityFixture(t, workspace, copied)
	for _, same := range []bool{true, false} {
		repository, relation := workspace, "same"
		if !same {
			repository, relation = copied, "different"
		}
		viewArgs := [][]string{
			{"stop", "list"}, {"attention", "list"},
			{"resume", a.issueID, "--repository", repository},
			{"review", "get", reviewID}, {"review", "list", "--issue", a.issueID},
			{"run", "inspect", runID, "--repository", repository},
			{"review", "record", "--input", filepath.Join(workspace, "human-review.json"), "--actor", "actor-denied", "--idempotency-key", "human-review"},
		}
		// Replay the exact accepted input, retaining its own key and original facts.
		viewArgs[6][3] = p21Input(t, workspace, "human-review-replay", input)
		viewArgs[6][7] = "human-review"
		for _, leaf := range viewArgs {
			operation := strings.Join(leaf[:2], ".")
			if leaf[0] == "resume" {
				operation = "resume"
			}
			project := workspace
			if leaf[0] == "review" && leaf[1] != "record" || leaf[0] == "stop" || leaf[0] == "attention" {
				project = repository
			}
			args := append([]string{"--project", project}, leaf...)
			row := p24AddHumanFixture(t, manifest, state+"-"+relation+"-"+operation, operation, "", workspace, args, 0)
			if leaf[0] == "resume" || leaf[0] == "run" || leaf[0] == "review" {
				p24AssertShortHumanRepository(t, row.HumanStdout, p24FixtureHEAD(t, workspace), dirty)
			}
			if leaf[0] == "resume" || leaf[0] == "run" {
				if !strings.Contains(row.HumanStdout, relation+" worktree") {
					t.Fatalf("%s comparison missing: %s", row.Name, row.HumanStdout)
				}
			}
		}
	}
}

func p24AddHumanFixture(t *testing.T, manifest *p24HumanManifest, name, operation, code, workspace string, args []string, expected int) p24HumanFixture {
	t.Helper()
	row := p24HumanFixture{Name: name, Operation: operation, Code: code, Workspace: workspace, CWD: workspace, HumanArgs: args, JSONArgs: append([]string{"--json"}, args...), ExpectedExit: expected, BinarySHA256: manifest.BinarySHA256}
	row.HumanCommand = p24FixtureCommand(manifest.Binary, row.HumanArgs)
	row.JSONCommand = p24FixtureCommand(manifest.Binary, row.JSONArgs)
	row.JSONStdout, row.JSONStderr = p24FixtureRun(t, manifest.Binary, row, true)
	row.HumanStdout, row.HumanStderr = p24FixtureRun(t, manifest.Binary, row, false)
	if code != "" {
		failure := p23Envelope(t, row.JSONStderr)
		if p21PathString(t, failure, "error", "code") != code {
			t.Fatalf("%s wrong denial: %s", name, row.JSONStderr)
		}
		if row.HumanStderr != p24ExpectedFailureText(t, failure, operation, containsAction(args, "--input")) {
			t.Fatalf("%s Human/JSON mismatch: %s", name, row.HumanStderr)
		}
		if code == "invalid-argument" && p21PathString(t, failure, "error", "details", "argument") != "input.expected_repository.dirty" {
			t.Fatal("nested path missing")
		}
		if code == "protected-decision-required" && !strings.Contains(row.HumanStderr, "blocked by:") {
			t.Fatal("blocked by missing")
		}
	} else if row.HumanStdout == "" || strings.Contains(row.HumanStdout, "{") {
		t.Fatalf("%s empty or raw JSON view", name)
	}
	manifest.Fixtures = append(manifest.Fixtures, row)
	return row
}

func p24FixtureRun(t *testing.T, binary string, row p24HumanFixture, jsonMode bool) (string, string) {
	t.Helper()
	args := row.HumanArgs
	if jsonMode {
		args = row.JSONArgs
	}
	out, stderr, err := closureRunCommand(row.CWD, binary, args...)
	exit := 0
	if err != nil {
		exit = processExitCode(t, err)
	}
	if exit != row.ExpectedExit {
		t.Fatalf("%s json=%t exit=%d want=%d: %s %s", row.Name, jsonMode, exit, row.ExpectedExit, out, stderr)
	}
	if exit != 0 && out != "" || exit == 0 && stderr != "" {
		t.Fatalf("%s invalid streams: %q %q", row.Name, out, stderr)
	}
	return out, stderr
}

func p24FixtureCommand(binary string, args []string) string {
	parts := append([]string{binary}, args...)
	for i, part := range parts {
		parts[i] = "'" + strings.ReplaceAll(part, "'", "'\"'\"'") + "'"
	}
	return strings.Join(parts, " ")
}

func p24AssertShortHumanRepository(t *testing.T, text, head string, dirty bool) {
	t.Helper()
	state := "clean"
	if dirty {
		state = "dirty"
	}
	if !strings.Contains(text, "HEAD "+head[:12]+"; "+state) {
		t.Fatalf("missing short %s repository: %s", state, text)
	}
	for _, forbidden := range []string{"worktree_id", "tree_oid", "dirty_digest", "capture_source", "<nil>", head} {
		if strings.Contains(text, forbidden) {
			t.Fatalf("repository leaked %q: %s", forbidden, text)
		}
	}
}

func p24FixtureHEAD(t *testing.T, workspace string) string {
	t.Helper()
	out, _ := closureRun(t, workspace, "git", "rev-parse", "HEAD")
	return strings.TrimSpace(out)
}

func p24AssertHumanManifest(t *testing.T, manifest p24HumanManifest) {
	t.Helper()
	data, err := os.ReadFile(manifest.Binary)
	if err != nil || fmt.Sprintf("sha256:%x", sha256.Sum256(data)) != manifest.BinarySHA256 {
		t.Fatalf("binary digest mismatch: %v", err)
	}
	seen := map[string]bool{}
	for _, row := range manifest.Fixtures {
		if seen[row.Name] || row.CWD == "" || row.Workspace == "" || row.BinarySHA256 != manifest.BinarySHA256 || row.HumanCommand != p24FixtureCommand(manifest.Binary, row.HumanArgs) || row.JSONCommand != p24FixtureCommand(manifest.Binary, row.JSONArgs) {
			t.Fatalf("incomplete/duplicate manifest row %s", row.Name)
		}
		seen[row.Name] = true
		out, stderr := p24FixtureRun(t, manifest.Binary, row, false)
		if out != row.HumanStdout || stderr != row.HumanStderr {
			t.Fatalf("retained fixture %s Human drifted", row.Name)
		}
		out, stderr = p24FixtureRun(t, manifest.Binary, row, true)
		if row.ExpectedExit == 0 {
			if row.Operation == "resume" {
				continue
			} // The lease projection is live; byte parity is checked at generation.
			if !reflect.DeepEqual(p23Envelope(t, out), p23Envelope(t, row.JSONStdout)) {
				t.Fatalf("retained fixture %s JSON drifted", row.Name)
			}
		} else if out != row.JSONStdout || stderr != row.JSONStderr {
			t.Fatalf("retained fixture %s JSON drifted", row.Name)
		}
	}
	for _, name := range []string{"idempotency-conflict", "repository-mismatch", "invalid-argument", "protected-decision-required", "source-digest-mismatch"} {
		if !seen[name] {
			t.Errorf("missing denial %s", name)
		}
	}
	for _, state := range []string{"clean", "dirty"} {
		for _, relation := range []string{"same", "different"} {
			for _, view := range []string{"stop.list", "attention.list", "resume", "review.get", "review.list", "review.record", "run.inspect"} {
				name := state + "-" + relation + "-" + view
				if !seen[name] {
					t.Errorf("missing view %s", name)
				}
			}
		}
	}
}
