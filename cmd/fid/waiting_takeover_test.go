package main

import (
	"testing"
	"time"

	"github.com/mbell-nzl/fid/internal/cli"
	"github.com/mbell-nzl/fid/internal/gittest"
)

func TestWaitingTakeoverKeepsUnstartedWorkWithoutALiveSession(t *testing.T) {
	t.Parallel()
	for _, mode := range []string{"expired", "interrupted"} {
		t.Run(mode, func(t *testing.T) {
			f := newLifecycleFixture(t)
			f.seed([]readyCommandRow{statusTask("FID-aaaaaa", "Waiting work", "todo", nil, nil, nil), statusTask("FID-bbbbbb", "Other work", "todo", nil, nil, nil)}, nil)
			f.hold("HLD-aaaaaa", "FID-aaaaaa", "", time.Minute)
			parkWaiting(t, f, "HLD-aaaaaa", waitingQuestion, "wait")
			f.workspace = gittest.AddWorktree(t, f.workspace)
			taken, _ := lifecycleResult(t, f.run("--json", "takeover", "FID-aaaaaa", "--actor", "human:mark", "--reason", "Continue here"), "task.takeover")
			if mode == "expired" {
				f.clock = func() time.Time { return statusNow.Add(time.Hour) }
			} else {
				parked, _ := lifecycleResult(t, f.run("--json", "park", "--hold", taken["hold_id"].(string), "--actor", "human:mark", "--interrupted", "--reason", "Need to stop", "--idempotency-key", "stop"), "hold.park")
				if parked["kept_for_worktree"] != true {
					t.Fatalf("unstarted waiting work was released: %v", parked)
				}
			}
			shown, _ := lifecycleResult(t, f.run("--json", "show", "FID-aaaaaa"), "issue.show")
			kept, ok := shown["kept"].(map[string]any)
			if !ok || kept["mode"] != "waiting" {
				t.Fatalf("waiting takeover lost reservation: %v", shown)
			}
			busy := f.run("--json", "takeover", "FID-bbbbbb", "--actor", "human:mark", "--reason", "Different task")
			if cli.ExitCode(busy.err) != 4 {
				t.Fatalf("waiting worktree took another task: %#v", busy)
			}
			assertJSONFailureCode(t, busy.stderr, "task.takeover", "worktree-busy")
			blocked, _ := lifecycleResult(t, f.run("--json", "next", "--claim", "--actor", "human:mark"), "next.claim")
			if blocked["kind"] != "none" || blocked["waiting"] == nil {
				t.Fatalf("waiting takeover resumed before answer: %v", blocked)
			}
		})
	}
}
