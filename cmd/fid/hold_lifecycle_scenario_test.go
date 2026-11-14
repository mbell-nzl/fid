package main

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/mbell-nzl/fid/internal/cli"
	"github.com/mbell-nzl/fid/internal/gittest"
)

func TestInterruptedParkKeepsOnlyStartedWorkForItsWorktree(t *testing.T) {
	t.Parallel()
	for _, state := range []string{"todo", "in_progress"} {
		t.Run(state, func(t *testing.T) {
			f := newLifecycleFixture(t)
			f.seed([]readyCommandRow{statusTask("FID-aaaaaa", "Interrupted work", state, nil, nil, nil)}, nil)
			main := f.workspace
			other := gittest.AddWorktree(t, main)
			f.hold("HLD-aaaaaa", "FID-aaaaaa", "", time.Minute)
			args := []string{"--json", "park", "--hold", "HLD-aaaaaa", "--actor", "agent:claude-builder", "--interrupted", "--reason", "Usage limit", "--idempotency-key", "park"}
			first := f.run(args...)
			lifecycleResult(t, first, "hold.park")
			f.workspace = other
			next, _ := lifecycleResult(t, f.run("--json", "next", "--claim", "--actor", "agent:other"), "next.claim")
			if state == "todo" {
				if next["kind"] != "new" || next["task"].(map[string]any)["id"] != "FID-aaaaaa" {
					t.Fatalf("released task elsewhere = %v", next)
				}
			} else {
				queue := next["queue"].([]any)
				entry := queue[0].(map[string]any)
				if next["kind"] != "none" || entry["reason"] != "interrupted-elsewhere" || entry["working"] != nil || entry["kept"].(map[string]any)["here"] != false {
					t.Fatalf("kept task elsewhere = %v", next)
				}
				f.workspace = main
				resumed, _ := lifecycleResult(t, f.run("--json", "next", "--claim", "--actor", "agent:new"), "next.claim")
				if resumed["kind"] != "resume" || resumed["hold"].(map[string]any)["replaces"] != "HLD-aaaaaa" {
					t.Fatalf("resumed = %v", resumed)
				}
				if again := f.run(args...); again.stdout != first.stdout || again.err != nil {
					t.Fatalf("park retry after resume = %#v, want %s", again, first.stdout)
				}
			}
		})
	}
}

func TestHeartbeatSetsTheLastSignOfLifeAndRefusesAtExpiry(t *testing.T) {
	t.Parallel()
	f := newLifecycleFixture(t)
	f.seed([]readyCommandRow{statusTask("FID-aaaaaa", "Interrupted work", "in_progress", nil, nil, nil)}, nil)
	f.hold("HLD-aaaaaa", "FID-aaaaaa", "", 30*time.Minute)
	lifecycleResult(t, f.run("--json", "heartbeat", "--hold", "HLD-aaaaaa", "--actor", "agent:claude-builder"), "hold.heartbeat")
	f.clock = func() time.Time { return statusNow.Add(time.Hour) }
	refused := f.run("--json", "heartbeat", "--hold", "HLD-aaaaaa", "--actor", "agent:claude-builder")
	if cli.ExitCode(refused.err) != 4 || !strings.Contains(refused.stderr, `"reason":"expired"`) {
		t.Fatalf("heartbeat at expiry = %#v", refused)
	}
	shown, _ := lifecycleResult(t, f.run("--json", "show", "FID-aaaaaa"), "issue.show")
	working := shown["working"].(map[string]any)
	kept := shown["kept"].(map[string]any)
	if working["since"] != "2026-11-08T08:30:00Z" || working["last_heartbeat_at"] != "2026-11-08T09:00:00Z" || kept["since"] != "2026-11-08T09:00:00Z" || kept["reason"] != nil {
		t.Fatalf("heartbeat lapse = %v", shown)
	}
	text := f.run("show", "FID-aaaaaa")
	if text.err != nil || !strings.Contains(strings.Join(strings.Fields(text.stdout), " "), "work stopped at 22:00 without a park") {
		t.Fatalf("show last heartbeat = %#v", text)
	}
	lifecycleGolden(t, f, "heartbeat-lapsed-show.txt", text.stdout)
	lifecycleGolden(t, f, "heartbeat-lapsed-show.json", f.run("--json", "show", "FID-aaaaaa").stdout)
	lifecycleGolden(t, f, "heartbeat-lapsed-next.txt", f.run("next").stdout)
}

func TestTakeoverOfAbsentLapsedAndKeptHolds(t *testing.T) {
	t.Parallel()
	for _, mode := range []string{"absent", "lapsed", "parked"} {
		t.Run(mode, func(t *testing.T) {
			f := newLifecycleFixture(t)
			f.seed([]readyCommandRow{statusTask("FID-aaaaaa", "Interrupted work", "in_progress", nil, nil, nil)}, nil)
			main := f.workspace
			other := gittest.AddWorktree(t, main)
			if mode != "absent" {
				ago := time.Hour
				if mode == "parked" {
					ago = time.Minute
				}
				f.hold("HLD-aaaaaa", "FID-aaaaaa", "", ago)
				if mode == "parked" {
					lifecycleResult(t, f.run("--json", "park", "--hold", "HLD-aaaaaa", "--actor", "agent:claude-builder", "--interrupted", "--reason", "Usage limit", "--idempotency-key", "park"), "hold.park")
				}
			}
			f.workspace = other
			taken, _ := lifecycleResult(t, f.run("--json", "takeover", "FID-aaaaaa", "--actor", "human:mark", "--reason", "Old worktree deleted"), "task.takeover")
			if mode == "absent" && taken["previous_hold_id"] != nil || mode != "absent" && taken["previous_hold_id"] != "HLD-aaaaaa" {
				t.Fatalf("%s takeover = %v", mode, taken)
			}
			shown, _ := lifecycleResult(t, f.run("--json", "show", "FID-aaaaaa"), "issue.show")
			if shown["issue"].(map[string]any)["status"] != "in_progress" || shown["working"].(map[string]any)["here"] != true || shown["kept"] != nil {
				t.Fatalf("takeover changed state = %v", shown)
			}
			f.workspace = main
			next, _ := lifecycleResult(t, f.run("--json", "next"), "next.get")
			if next["kind"] != "none" || next["queue"].([]any)[0].(map[string]any)["reason"] != "held-elsewhere" {
				t.Fatalf("old worktree after takeover = %v", next)
			}
		})
	}
}

func TestLifecycleRebuildAndExport(t *testing.T) {
	t.Parallel()
	workspace := gittest.NewWorkspace(t)
	runFIDCommand(t, workspace, "init", "--prefix", "FID")
	_, created := runFIDCommand(t, workspace, "issue", "create", "--kind", "task", "--title", "Lifecycle replay", "--lane", "standard", "--actor", "human:mark")
	task := parseCreatedIssueID(t, created)
	takenBefore, _ := lifecycleResult(t, executeCommandForTest(workspace, "--json", "next", "--claim", "--actor", "agent:builder"), "next.claim")
	runFIDCommand(t, workspace, "start", "--hold", takenBefore["hold"].(map[string]any)["id"].(string), "--actor", "agent:builder", "--idempotency-key", "lifecycle-start")
	// Takeover also starts holds for already-started tasks with no earlier hold.
	taken, _ := lifecycleResult(t, executeCommandForTest(workspace, "--json", "takeover", task, "--actor", "human:mark", "--reason", "Resume task", "--idempotency-key", "take"), "task.takeover")
	hold := fmt.Sprint(taken["hold_id"])
	lifecycleResult(t, executeCommandForTest(workspace, "--json", "heartbeat", "--hold", hold, "--actor", "human:mark"), "hold.heartbeat")
	park := []string{"--json", "park", "--hold", hold, "--actor", "human:mark", "--interrupted", "--reason", "Usage limit", "--idempotency-key", "park"}
	first := executeCommandForTest(workspace, park...)
	lifecycleResult(t, first, "hold.park")
	query := "SELECT * FROM holds ORDER BY id"
	before := executeCommandForTest(workspace, "sql", query)
	decisions := executeCommandForTest(workspace, "sql", "SELECT * FROM person_decisions ORDER BY id")
	runFIDCommand(t, workspace, "rebuild")
	if after := executeCommandForTest(workspace, "sql", query); before.err != nil || after.err != nil || before.stdout != after.stdout {
		t.Fatalf("hold rebuild before=%#v after=%#v", before, after)
	}
	if after := executeCommandForTest(workspace, "sql", "SELECT * FROM person_decisions ORDER BY id"); decisions.err != nil || after.err != nil || decisions.stdout != after.stdout {
		t.Fatalf("decision rebuild before=%#v after=%#v", decisions, after)
	}
	if again := executeCommandForTest(workspace, park...); again.err != nil || again.stdout != first.stdout {
		t.Fatalf("park retry after rebuild = %#v, want %s", again, first.stdout)
	}
	archive := filepath.Join(t.TempDir(), "lifecycle.zip")
	runFIDCommand(t, workspace, "export", "create", archive)
	if answer := executeCommandForTest(workspace, "--json", "export", "validate", archive); answer.err != nil {
		t.Fatalf("export validate = %#v", answer)
	}
}
