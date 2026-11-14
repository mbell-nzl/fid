package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mbell-nzl/fid/internal/cli"
)

func TestHarmlessStartLeavesFreshKeyUnbound(t *testing.T) {
	t.Parallel()
	f := plannedStartFixture(t, "todo", stringPointer("standard"))
	firstArgs := []string{"--json", "start", "--hold", "HLD-aaaaaa", "--actor", "agent:claude-builder", "--idempotency-key", "recorded"}
	first := f.run(firstArgs...)
	lifecycleResult(t, first, "task.start")
	args := append([]string(nil), firstArgs...)
	args[len(args)-1] = "fresh"
	before := startEventCount(t, f.workspace)
	again := f.run(args...)
	_, event := lifecycleResult(t, again, "task.start")
	if event != nil || startEventCount(t, f.workspace) != before {
		t.Fatal("harmless start recorded an Event")
	}
	if repeat := f.run(args...); repeat.err != nil || repeat.stdout != again.stdout {
		t.Fatalf("state repeat = %#v", repeat)
	}
	query := f.run("sql", "SELECT count(*) AS n FROM idempotency_results WHERE idempotency_key='fresh'")
	if query.err != nil || !strings.Contains(query.stdout, `"n":0`) {
		t.Fatalf("no-op retry binding = %#v", query)
	}
	lifecycleResult(t, f.run("--json", "issue", "transition", "FID-9f13aa", "--to", "cancelled", "--actor", "human:mark", "--reason", "No longer needed"), "issue.transition")
	refused := f.run(args...)
	assertJSONFailureCode(t, refused.stderr, "task.start", "hold-not-current")
	if repeat := f.run(firstArgs...); repeat.err != nil || repeat.stdout != first.stdout {
		t.Fatalf("recorded start lost original replay: %#v", repeat)
	}
}

func TestStartReadsPolicySnapshotsChangesAndPreservesDirtyGit(t *testing.T) {
	t.Parallel()
	f := plannedStartFixture(t, "todo", stringPointer("standard"))
	if err := os.WriteFile(filepath.Join(f.workspace, "unfinished.txt"), []byte("work in progress\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	beforeGit := identityGit(t, f.workspace, "status", "--porcelain=v1")
	beforeHead := identityGit(t, f.workspace, "rev-parse", "HEAD")
	args := []string{"--json", "start", "--hold", "HLD-aaaaaa", "--actor", "agent:claude-builder", "--idempotency-key", "start"}
	lifecycleResult(t, f.run(args...), "task.start")
	if beforeGit == "" || identityGit(t, f.workspace, "status", "--porcelain=v1") != beforeGit || identityGit(t, f.workspace, "rev-parse", "HEAD") != beforeHead {
		t.Fatal("start changed dirty git state")
	}
	first := policySnapshots(t, f.workspace)
	if len(first) != 1 {
		t.Fatalf("start snapshots = %v", first)
	}
	// Changing only the worktree has no policy authority.
	addStartPolicyLine(t, f.workspace, "token_budget = 42")
	args[len(args)-1] = "no-op-worktree"
	_, event := lifecycleResult(t, f.run(args...), "task.start")
	if event != nil || len(policySnapshots(t, f.workspace)) != 1 {
		t.Fatal("worktree policy changed start policy")
	}
	commit := policyCommit(t, f.workspace)
	identityGit(t, f.workspace, "update-ref", "refs/remotes/origin/trunk", commit)
	args[len(args)-1] = "no-op-policy-change"
	_, event = lifecycleResult(t, f.run(args...), "task.start")
	snapshots := policySnapshots(t, f.workspace)
	if event != nil || len(snapshots) != 2 || snapshots[0].Digest == snapshots[1].Digest {
		t.Fatalf("changed policy snapshots = %v, event=%v", snapshots, event)
	}
}

func TestStartPolicyUnavailableGolden(t *testing.T) {
	t.Parallel()
	f := plannedStartFixture(t, "todo", stringPointer("standard"))
	addStartPolicyLine(t, f.workspace, "unknown_key = true")
	commit := policyCommit(t, f.workspace)
	identityGit(t, f.workspace, "update-ref", "refs/remotes/origin/trunk", commit)
	before := startEventCount(t, f.workspace)
	args := []string{"start", "--hold", "HLD-aaaaaa", "--actor", "agent:claude-builder", "--idempotency-key", "start"}
	for _, jsonMode := range []bool{false, true} {
		invoke := args
		ext := ".txt"
		if jsonMode {
			invoke = append([]string{"--json"}, args...)
			ext = ".json"
		}
		result := f.run(invoke...)
		if cli.ExitCode(result.err) != 4 || result.stdout != "" {
			t.Fatalf("policy refusal = %#v", result)
		}
		if jsonMode {
			assertJSONFailureCode(t, result.stderr, "task.start", "policy-unavailable")
		}
		output := strings.ReplaceAll(result.stderr, commit, "<COMMIT>")
		lifecycleGolden(t, f, filepath.Join("..", "start", "policy-unavailable"+ext), output)
	}
	if startEventCount(t, f.workspace) != before || len(policySnapshots(t, f.workspace)) != 0 {
		t.Fatal("policy failure wrote a start or snapshot")
	}
}

func addStartPolicyLine(t *testing.T, root, line string) {
	t.Helper()
	path := filepath.Join(root, ".fid.toml")
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	text := strings.Replace(string(raw), "[lanes.standard]\n", "[lanes.standard]\n"+line+"\n", 1)
	if err := os.WriteFile(path, []byte(text), 0o600); err != nil {
		t.Fatal(err)
	}
}
