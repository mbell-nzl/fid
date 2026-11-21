//go:build integration

package integration

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mbell-nzl/fid/internal/gittest"
)

func TestLaneFloorPublicJourney(t *testing.T) {
	t.Parallel()
	root, binary := gittest.NewWorkspace(t), integrationFIDBinary(t)
	call := func(args ...string) map[string]any {
		t.Helper()
		return p21Run(t, binary, root, append([]string{"--json"}, args...)...)
	}
	call("init", "--prefix", "FID")
	path := filepath.Join(root, ".fid.toml")
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	raw = append(raw, []byte("\n[lanes.routine]\npaths=['**/*.md']\nchecks=[]\n[lanes.standard]\nchecks=[]\n[lanes.assured]\npaths=['secure/**']\nchecks=['pass']\n[checks.pass]\nrun=['true']\ntimeout='1m'\n[review]\ntest_adequacy=[]\n")...)
	if err := os.WriteFile(path, raw, 0o600); err != nil {
		t.Fatal(err)
	}
	closureRun(t, root, "git", "add", ".")
	closureRun(t, root, "git", "-c", "user.name=fid test", "-c", "user.email=fid@example.test", "commit", "-qm", "Policy base")
	closureRun(t, root, "git", "remote", "add", "origin", root)
	closureRun(t, root, "git", "update-ref", "refs/remotes/origin/trunk", "HEAD")
	closureRun(t, root, "git", "symbolic-ref", "refs/remotes/origin/HEAD", "refs/remotes/origin/trunk")
	task := p21PathString(t, call("issue", "create", "--kind", "task", "--title", "Lane floor journey", "--lane", "routine", "--actor", "human:mark", "--idempotency-key", "create"), "result", "issue", "id")
	hold := p21PathString(t, call("takeover", task, "--actor", "human:mark", "--reason", "Build the journey"), "result", "hold_id")
	call("start", "--hold", hold, "--actor", "human:mark", "--idempotency-key", "start")
	if err := os.MkdirAll(filepath.Join(root, "secure"), 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "secure/task.go"), []byte("package secure\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	closureRun(t, root, "git", "add", ".")
	closureRun(t, root, "git", "-c", "user.name=fid test", "-c", "user.email=fid@example.test", "commit", "-qm", "Assured task change")
	checked := call("check", "--hold", hold, "--actor", "human:mark", "--idempotency-key", "floor-check")
	checks := checked["result"].(map[string]any)["checks"].([]any)
	if len(checks) != 1 || checks[0].(map[string]any)["name"] != "pass" {
		t.Fatalf("floor checks=%v", checks)
	}
	if lane := p21PathString(t, call("show", task), "result", "issue", "lane"); lane != "routine" {
		t.Fatalf("check changed lane: %s", lane)
	}
	call("handoff", "--hold", hold, "--actor", "human:mark", "--pr", "https://github.com/owner/repo/pull/1", "--idempotency-key", "handoff")
	for _, rebuild := range []bool{false, true} {
		if rebuild {
			call("rebuild")
		}
		shown := call("show", task)
		if lane := p21PathString(t, shown, "result", "issue", "lane"); lane != "assured" {
			t.Fatalf("hand-off did not raise: %s", lane)
		}
		if lane := p21PathString(t, shown, "result", "floor", "lane"); lane != "assured" {
			t.Fatalf("floor=%s", lane)
		}
		text, stderr, err := closureRunCommand(root, binary, "show", task)
		if err != nil || stderr != "" || !strings.Contains(strings.Join(strings.Fields(text), " "), "floor: secure/task.go · raised from Routine") {
			t.Fatalf("show=%v/%s/%s", err, text, stderr)
		}
	}
}
