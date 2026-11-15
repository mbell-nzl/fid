//go:build integration

package integration

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestReconcileBuiltBinaryExitContract(t *testing.T) {
	t.Parallel()
	j := newDoneJourney(t, "standard")
	// Start pending: a prior green receipt would correctly remain green offline.
	j.ci = []any{}
	for _, phase := range []string{"refused", "waiting", "completed", "empty"} {
		switch phase {
		case "waiting":
			j.review(t, "primary", "APPROVE", "review")
			j.ci = []any{}
		case "completed":
			j.ci = []any{map[string]any{"id": 101, "run_attempt": 1, "workflow_id": 7, "name": "CI", "path": ".github/workflows/ci.yml", "event": "push", "head_sha": strings.Repeat("a", 40), "status": "completed", "conclusion": "success", "html_url": "https://github.com/mbell-nzl/fid/actions/runs/101"}}
		}
		out, stderr, err := j.run(t, "--json", "done", "--actor", "agent:finisher")
		if err != nil || stderr != "" {
			t.Fatalf("%s exit=%v stderr=%s", phase, err, stderr)
		}
		fields := p23Envelope(t, out)
		outcomes, ok := p21Path(t, fields, "result", "outcomes").([]any)
		if !ok {
			t.Fatalf("%s: %s", phase, out)
		}
		if phase == "empty" {
			if len(outcomes) != 0 {
				t.Fatalf("completed task touched: %s", out)
			}
			continue
		}
		if len(outcomes) != 1 || outcomes[0].(map[string]any)["task_id"] != j.task || outcomes[0].(map[string]any)["status"] != phase {
			t.Fatalf("%s: %s", phase, out)
		}
		if phase == "waiting" {
			assertReconcileOfflineReads(t, j)
		}
	}
	out, stderr, err := j.run(t, "--json", "done", "--actor", "agent:finisher", "--idempotency-key", "batch")
	if processExitCode(t, err) != 2 || out != "" || !strings.Contains(stderr, `"code":"invalid-argument"`) {
		t.Fatalf("batch key=%v %s %s", err, out, stderr)
	}
}

func assertReconcileOfflineReads(t *testing.T, j *doneJourney) {
	t.Helper()
	git, err := exec.LookPath("git")
	if err != nil {
		t.Fatal(err)
	}
	bin := t.TempDir()
	if err := os.Symlink(git, filepath.Join(bin, "git")); err != nil {
		t.Fatal(err)
	}
	before := p21JournalEventCount(t, j.root)
	for _, args := range [][]string{{"status"}, {"show", j.task}} {
		cmd := exec.CommandContext(t.Context(), j.binary, args...)
		cmd.Dir = j.root
		cmd.Env = append(os.Environ(), "PATH="+bin)
		out, err := cmd.CombinedOutput()
		if err != nil || !strings.Contains(string(out), "merged, waiting for CI") {
			t.Fatalf("offline %v=%v %s", args, err, out)
		}
	}
	if p21JournalEventCount(t, j.root) != before {
		t.Fatal("offline read wrote Events")
	}
}
