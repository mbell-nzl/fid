//go:build integration

package integration

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mbell-nzl/fid/internal/config"
	"github.com/mbell-nzl/fid/internal/gittest"
)

func handoffJourneyRepo(t *testing.T, testPaths bool) (string, string) {
	t.Helper()
	root, binary := gittest.NewWorkspace(t), integrationFIDBinary(t)
	p21Run(t, binary, root, "init", "--prefix", "FID", "--json")
	path := filepath.Join(root, ".fid.toml")
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	text := string(raw) + "\n[lanes.standard]\nchecks=['test']\n[checks.test]\nrun=['true']\ntimeout='1m'\n"
	if testPaths {
		text += "[review]\ntest_adequacy=[]\n"
	}
	if err := os.WriteFile(path, []byte(text), 0o600); err != nil {
		t.Fatal(err)
	}
	closureRun(t, root, "git", "add", ".")
	closureRun(t, root, "git", "-c", "user.name=fid test", "-c", "user.email=fid@example.test", "commit", "-qm", "Trusted policy")
	origin := filepath.Join(t.TempDir(), "origin.git")
	closureRun(t, root, "git", "init", "--bare", origin)
	closureRun(t, root, "git", "remote", "add", "origin", origin)
	closureRun(t, root, "git", "push", "-q", "origin", "HEAD:trunk")
	closureRun(t, root, "git", "symbolic-ref", "refs/remotes/origin/HEAD", "refs/remotes/origin/trunk")
	return root, binary
}

func TestLeanTaskJourneyFromNextThroughHandoff(t *testing.T) {
	t.Parallel()
	root, binary := handoffJourneyRepo(t, true)
	call := func(args ...string) map[string]any {
		t.Helper()
		result := p21Run(t, binary, root, append([]string{"--json"}, args...)...)
		assertLeanHandoffJournal(t, root)
		return result
	}
	created := call("issue", "create", "--kind", "task", "--title", "Complete lean task", "--lane", "standard", "--actor", "agent:builder", "--idempotency-key", "create")
	task := p21PathString(t, created, "result", "issue", "id")
	taken := call("next", "--claim", "--actor", "agent:builder")
	hold := p21PathString(t, taken, "result", "hold", "id")
	call("start", "--hold", hold, "--actor", "agent:builder", "--idempotency-key", "start")
	call("heartbeat", "--hold", hold, "--actor", "agent:builder")
	call("park", "--hold", hold, "--actor", "agent:builder", "--interrupted", "--reason", "Usage limit", "--idempotency-key", "interrupted")
	taken = call("next", "--claim", "--actor", "agent:builder")
	if p21PathString(t, taken, "result", "kind") != "resume" || p21PathString(t, taken, "result", "hold", "replaces") != hold {
		t.Fatalf("resume=%v", taken)
	}
	hold = p21PathString(t, taken, "result", "hold", "id")
	call("park", "--hold", hold, "--actor", "agent:builder", "--waiting", "Include retries?", "--idempotency-key", "waiting")
	waiting := call("next")
	if p21PathString(t, waiting, "result", "kind") != "none" {
		t.Fatalf("waiting=%v", waiting)
	}
	call("answer", task, "Include them.", "--actor", "human:mark")
	taken = call("next", "--claim", "--actor", "agent:builder")
	if p21PathString(t, taken, "result", "answer", "text") != "Include them." {
		t.Fatalf("answer=%v", taken)
	}
	hold = p21PathString(t, taken, "result", "hold", "id")
	checked := call("check", "--hold", hold, "--actor", "agent:builder", "--idempotency-key", "check")
	closureRun(t, root, "git", "-c", "user.name=fid test", "-c", "user.email=fid@example.test", "commit", "--allow-empty", "-qm", "Ready for review")
	args := []string{"handoff", "--hold", hold, "--actor", "agent:builder", "--pr", "https://github.com/mbell-nzl/fid/pull/63", "--idempotency-key", "handoff"}
	handed := call(args...)
	checks := checked["result"].(map[string]any)["checks"].([]any)
	results := handed["result"].(map[string]any)["result_ids"].([]any)
	if len(results) != 1 || results[0] != checks[0].(map[string]any)["id"] {
		t.Fatalf("hand-off evidence=%v", handed)
	}
	if p21PathString(t, handed, "result", "status") != "in_review" {
		t.Fatalf("handoff=%v", handed)
	}
	waiting = call("next", "--claim", "--actor", "agent:builder")
	if p21PathString(t, waiting, "result", "kind") != "none" || p21PathString(t, waiting, "result", "waiting", "kind") != "review" {
		t.Fatalf("review=%v", waiting)
	}
	shown := call("show", task)
	if p21PathString(t, shown, "result", "handoff", "id") != p21PathString(t, handed, "result", "handoff_id") {
		t.Fatalf("show=%v", shown)
	}
	call("rebuild")
	retry := call(args...)
	before, _ := json.Marshal(handed)
	after, _ := json.Marshal(retry)
	if string(before) != string(after) {
		t.Fatalf("retry after replay=%s, original=%s", after, before)
	}
}

func assertLeanHandoffJournal(t *testing.T, root string) {
	t.Helper()
	files, err := filepath.Glob(filepath.Join(config.LedgerDir(root), "journal", "*.jsonl"))
	if err != nil {
		t.Fatal(err)
	}
	allowed := map[string]bool{"project.initialized": true, "issue.created": true, "policy.snapshotted": true, "hold.taken": true, "task.started": true, "hold.heartbeat": true, "hold.parked": true, "hold.waiting": true, "task.answered": true, "check.started": true, "check.process-started": true, "check.result-recorded": true, "task.handed-off": true}
	if len(files) == 0 {
		t.Fatal("no journal files")
	}
	for _, path := range files {
		raw, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		for _, line := range strings.Split(strings.TrimSpace(string(raw)), "\n") {
			var event struct {
				Type string `json:"event_type"`
			}
			if err := json.Unmarshal([]byte(line), &event); err != nil {
				t.Fatal(err)
			}
			if !allowed[event.Type] {
				t.Fatalf("unexpected event family in lean journey: %s", event.Type)
			}
		}
	}
}

func handoffCodePair(code string) func(*testing.T) {
	return func(t *testing.T) {
		t.Helper()
		root, binary := handoffJourneyRepo(t, code != "test-paths-not-configured")
		created := p21Run(t, binary, root, "--json", "issue", "create", "--kind", "task", "--title", "Hand-off refusal", "--lane", "standard", "--actor", "agent:builder", "--idempotency-key", "create")
		task := p21PathString(t, created, "result", "issue", "id")
		taken := p21Run(t, binary, root, "--json", "next", "--claim", "--actor", "agent:builder")
		hold := p21PathString(t, taken, "result", "hold", "id")
		p21Run(t, binary, root, "--json", "start", "--hold", hold, "--actor", "agent:builder", "--idempotency-key", "start")
		if code == "tree-dirty" {
			if err := os.WriteFile(filepath.Join(root, "dirty"), []byte("dirty"), 0o600); err != nil {
				t.Fatal(err)
			}
		}
		args := []string{"handoff", "--hold", hold, "--actor", "agent:builder", "--pr", "https://github.com/mbell-nzl/fid/pull/63", "--idempotency-key", "handoff"}
		stdout, stderr, err := closureRunCommand(root, binary, append([]string{"--json"}, args...)...)
		if err == nil || stdout != "" || processExitCode(t, err) != 4 {
			t.Fatalf("JSON=%v/%s/%s", err, stdout, stderr)
		}
		envelope := p23Envelope(t, stderr)
		if envelope["operation"] != "task.handoff" || envelope["error"].(map[string]any)["code"] != code {
			t.Fatalf("failure=%v", envelope)
		}
		stdout, stderr, err = closureRunCommand(root, binary, args...)
		wantText := p24ExpectedFailureText(t, envelope, "task.handoff", false)
		if code == "checks-unsatisfied" {
			wantText += "next: fid check --hold " + hold + " --actor agent:builder --idempotency-key <key>\n"
		}
		if err == nil || stdout != "" || processExitCode(t, err) != 4 || stderr != wantText {
			t.Fatalf("text=%v/%s/%s", err, stdout, stderr)
		}
		shown := p21Run(t, binary, root, "--json", "show", task)
		if shown["result"].(map[string]any)["handoff"] != nil {
			t.Fatal("refusal recorded a hand-off")
		}
	}
}
