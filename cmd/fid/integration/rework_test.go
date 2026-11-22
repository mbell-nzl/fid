//go:build integration

package integration

import (
	"encoding/json"
	"github.com/mbell-nzl/fid/internal/gittest"
	"os"
	"path/filepath"
	"testing"
)

func TestPublicReviewReworkJourney(t *testing.T) {
	t.Parallel()
	root, binary := handoffJourneyRepo(t, true)
	builder := root
	reviewer, other := gittest.AddWorktree(t, root), gittest.AddWorktree(t, root)
	call := func(args ...string) map[string]any {
		t.Helper()
		return p21Run(t, binary, root, append([]string{"--json"}, args...)...)
	}
	commit := func(contents string) {
		t.Helper()
		if err := os.WriteFile(filepath.Join(root, "change.txt"), []byte(contents), 0o600); err != nil {
			t.Fatal(err)
		}
		closureRun(t, root, "git", "add", ".")
		closureRun(t, root, "git", "-c", "user.name=fid test", "-c", "user.email=fid@example.test", "commit", "-qm", contents)
	}
	created := call("issue", "create", "--kind", "task", "--title", "Review and fix", "--lane", "standard", "--actor", "agent:builder", "--idempotency-key", "create")
	task := p21PathString(t, created, "result", "issue", "id")
	taken := call("next", "--claim", "--actor", "agent:builder", "--model", "codex/gpt-6.1-sol")
	original := p21PathString(t, taken, "result", "hold", "id")
	call("start", "--hold", original, "--actor", "agent:builder", "--idempotency-key", "start")
	commit("Initial implementation\n")
	call("check", "--hold", original, "--actor", "agent:builder", "--idempotency-key", "check")
	call("handoff", "--hold", original, "--actor", "agent:builder", "--pr", "https://github.com/mbell-nzl/fid/pull/63", "--idempotency-key", "handoff")
	root = reviewer
	selected := call("next", "--claim", "--actor", "agent:reviewer", "--model", "claude/opus-5.5")
	if p21PathString(t, selected, "result", "kind") != "review" {
		t.Fatalf("review work=%v", selected)
	}
	root = other
	if got := call("next", "--claim", "--actor", "agent:second", "--model", "claude/opus-5.5"); p21PathString(t, got, "result", "kind") != "none" {
		t.Fatalf("review handed to two worktrees: %v", got)
	}
	root = reviewer

	input := map[string]any{
		"reviewer": map[string]any{"actor": "agent:reviewer", "model": "claude/opus-5.5", "session_id": "first"},
		"stream":   "primary", "profiles": []string{"code"}, "verdict": "CHANGES", "purpose": "Fix the task", "verification": "Public journey exercised",
		"findings": []any{map[string]any{"profile": "code", "severity": "major", "summary": "Missing case", "location": nil, "evidence": "Case reproduced", "fix": "Handle the case"}},
	}
	record := func(key string) map[string]any {
		t.Helper()
		raw, err := json.Marshal(input)
		if err != nil {
			t.Fatal(err)
		}
		path := filepath.Join(t.TempDir(), "round.json")
		if err := os.WriteFile(path, raw, 0o600); err != nil {
			t.Fatal(err)
		}
		return call("review", "record", task, "--input", path, "--actor", "agent:courier", "--idempotency-key", key)
	}
	first := record("changes")
	finding := first["result"].(map[string]any)["round"].(map[string]any)["findings"].([]any)[0].(map[string]any)["id"].(string)
	shown := call("show", task)
	if p21PathString(t, shown, "result", "issue", "status") != "in_progress" || p21PathString(t, shown, "result", "kept", "mode") != "rework" {
		t.Fatalf("CHANGES did not return work: %v", shown)
	}
	root = builder
	taken = call("next", "--claim", "--actor", "agent:builder", "--model", "codex/gpt-6.1-sol")
	if p21PathString(t, taken, "result", "kind") != "rework" || p21PathString(t, taken, "result", "hold", "replaces") != original {
		t.Fatalf("rework=%v", taken)
	}
	hold := p21PathString(t, taken, "result", "hold", "id")
	commit("Handle the missing case\n")
	call("check", "--hold", hold, "--actor", "agent:builder", "--idempotency-key", "fix-check")
	call("handoff", "--hold", hold, "--actor", "agent:builder", "--pr", "https://github.com/mbell-nzl/fid/pull/63", "--idempotency-key", "fix-handoff")
	root = reviewer
	if got := call("next", "--claim", "--actor", "agent:reviewer", "--model", "claude/opus-5.5"); p21PathString(t, got, "result", "kind") != "review" {
		t.Fatalf("review after rework=%v", got)
	}
	input["verdict"], input["findings"] = "APPROVE", []any{}
	input["reviewer"].(map[string]any)["session_id"] = "fixed"
	input["dispositions"] = []any{map[string]any{"finding_id": finding, "disposition": "resolved", "reason": "Regression verifies the fix"}}
	record("approve")
	call("rebuild")
	listed := call("review", "list", task)["result"].(map[string]any)["rounds"].([]any)
	if len(listed) != 2 || listed[0].(map[string]any)["current"] != false || listed[1].(map[string]any)["current"] != true || listed[0].(map[string]any)["findings"].([]any)[0].(map[string]any)["open"] != false {
		t.Fatalf("replayed rework=%v", listed)
	}
}
