package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/mbell-nzl/fid/internal/cli"
	"github.com/mbell-nzl/fid/internal/gittest"
)

func lifecycleResult(t *testing.T, result commandExecutionResult, operation string) (map[string]any, any) {
	t.Helper()
	return decodeNext(t, result, operation)
}

func TestHoldHeartbeatAndInterruptedPark(t *testing.T) {
	t.Parallel()
	f := newLifecycleFixture(t)
	nextWork(f)
	f.hold("HLD-aaaaaa", "FID-71a0c2", "", 30*time.Minute)
	result, event := lifecycleResult(t, f.run("--json", "heartbeat", "--hold", "HLD-aaaaaa", "--actor", "agent:claude-builder"), "hold.heartbeat")
	if result["last_heartbeat_at"] != "2026-11-08T09:00:00Z" || result["expires_at"] != "2026-11-08T10:00:00Z" || result["heartbeat_interval_seconds"] != float64(300) || event == nil {
		t.Fatalf("heartbeat = %v, %v", result, event)
	}
	args := []string{"--json", "park", "--hold", "HLD-aaaaaa", "--actor", "agent:claude-builder", "--interrupted", "--reason", "Usage limit", "--idempotency-key", "park-1"}
	first := f.run(args...)
	parked, _ := lifecycleResult(t, first, "hold.park")
	if parked["kept_for_worktree"] != true || parked["mode"] != "interrupted" {
		t.Fatalf("park = %v", parked)
	}
	if again := f.run(args...); again.stdout != first.stdout || again.err != nil {
		t.Fatalf("same-key park = %#v, want %s", again, first.stdout)
	}
	shown, _ := lifecycleResult(t, f.run("--json", "show", "FID-71a0c2"), "issue.show")
	if shown["working"] != nil || shown["kept"].(map[string]any)["reason"] != "Usage limit" {
		t.Fatalf("show parked = %v", shown)
	}
	resumed, _ := lifecycleResult(t, f.run("--json", "next", "--claim", "--actor", "agent:new"), "next.claim")
	if resumed["kind"] != "resume" || resumed["hold"].(map[string]any)["replaces"] != "HLD-aaaaaa" {
		t.Fatalf("resume = %v", resumed)
	}
}

func TestTakeoverEndsTheOldSessionAndReplays(t *testing.T) {
	t.Parallel()
	f := newLifecycleFixture(t)
	nextWork(f)
	f.hold("HLD-aaaaaa", "FID-71a0c2", "/w/old", 10*time.Minute)
	args := []string{"--json", "takeover", "FID-71a0c2", "--actor", "human:mark", "--reason", "Agent stuck", "--idempotency-key", "take-1"}
	first := f.run(args...)
	taken, _ := lifecycleResult(t, first, "task.takeover")
	if taken["previous_hold_id"] != "HLD-aaaaaa" || taken["decision_id"] == nil {
		t.Fatalf("takeover = %v", taken)
	}
	if again := f.run(args...); again.stdout != first.stdout || again.err != nil {
		t.Fatalf("same-key takeover = %#v, want %s", again, first.stdout)
	}
	again := f.run("--json", "takeover", "FID-71a0c2", "--actor", "human:mark", "--reason", "Again")
	if cli.ExitCode(again.err) != 4 {
		t.Fatalf("repeat takeover = %#v", again)
	}
	assertJSONFailureCode(t, again.stderr, "task.takeover", "already-held")
	for _, command := range [][]string{{"heartbeat"}, {"park", "--interrupted", "--reason", "Late", "--idempotency-key", "late"}} {
		attempt := f.run(append(append([]string{"--json"}, command...), "--hold", "HLD-aaaaaa", "--actor", "agent:claude-builder")...)
		if cli.ExitCode(attempt.err) != 4 || !strings.Contains(attempt.stderr, `"reason":"ended"`) {
			t.Fatalf("stale %v = %#v", command, attempt)
		}
	}
}

func TestCancelRequiresAPersonAndReleasesKeptWork(t *testing.T) {
	t.Parallel()
	f := newLifecycleFixture(t)
	nextWork(f)
	f.hold("HLD-aaaaaa", "FID-71a0c2", "", 10*time.Minute)
	refused := f.run("--json", "issue", "transition", "FID-71a0c2", "--to", "cancelled", "--actor", "agent:claude-builder", "--reason", "No longer needed", "--idempotency-key", "cancel-agent")
	if cli.ExitCode(refused.err) != 5 {
		t.Fatalf("agent cancel = %#v", refused)
	}
	assertJSONFailureCode(t, refused.stderr, "issue.transition", "authority-denied")
	lifecycleResult(t, f.run("--json", "park", "--hold", "HLD-aaaaaa", "--actor", "agent:claude-builder", "--interrupted", "--reason", "Usage limit", "--idempotency-key", "park"), "hold.park")
	cancelled, _ := lifecycleResult(t, f.run("--json", "issue", "transition", "FID-71a0c2", "--to", "cancelled", "--actor", "human:mark", "--reason", "No longer needed"), "issue.transition")
	if cancelled["decision_id"] == nil {
		t.Fatalf("cancel = %v", cancelled)
	}
	shown, _ := lifecycleResult(t, f.run("--json", "show", "FID-71a0c2"), "issue.show")
	if shown["kept"] != nil || shown["working"] != nil {
		t.Fatalf("cancel kept work = %v", shown)
	}
	decisions := f.run("sql", "SELECT kind, actor, reason FROM person_decisions")
	var rows any
	if decisions.err != nil || json.Unmarshal([]byte(decisions.stdout), &rows) != nil || !strings.Contains(decisions.stdout, `"kind":"cancel"`) {
		t.Fatalf("decisions = %#v", decisions)
	}
}

// A short root keeps wrapping independent of the test runner's temporary
// directory suffix, while still executing the actual CLI renderers.
func newLifecycleFixture(t *testing.T) *statusFixture {
	t.Helper()
	f := newStatusFixture(t, 100*time.Hour)
	base, err := os.MkdirTemp("", "fhl-")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := os.RemoveAll(base); err != nil {
			t.Error(err)
		}
	})
	root := filepath.Join(base, "worktree")
	if err := os.Rename(f.workspace, root); err != nil {
		t.Fatal(err)
	}
	f.workspace = root
	if err := gittest.CheckRepository(root); err != nil {
		t.Fatal(err)
	}
	return f
}
