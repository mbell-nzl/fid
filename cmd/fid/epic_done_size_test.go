package main

import (
	"strings"
	"testing"
	"time"
)

func TestEpicDoneEventSizeBudget(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name, reason, key string
		wantSuccess       bool
	}{
		{"escaped reason fits without retry", strings.Repeat("<", 100000), "", true},
		{"escaped reason exceeds with retry", strings.Repeat("<", 100000), "close", false},
		{"escaped reason exceeds without retry", strings.Repeat("<", 200000), "", false},
		{"plain reason exceeds with retry", strings.Repeat("x", 600000), "close", false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			f, epic := newEpicCancellationFixture(t)
			f.clock = time.Now
			args := []string{"--json", "done", epic, "--actor", "human:mark", "--reason", tc.reason}
			if tc.key != "" {
				args = append(args, "--idempotency-key", tc.key)
			}
			before := snapshotWorkspacePersistence(t, f.workspace)
			got := f.run(args...)
			if tc.wantSuccess {
				result, _ := lifecycleResult(t, got, "task.done")
				if result["status"] != "done" || result["decision"].(map[string]any)["reason"] != tc.reason {
					t.Fatal("close lost its reason or resulting state")
				}
			} else {
				if got.err == nil || got.stdout != "" || !strings.Contains(got.stderr, "epic close Event exceeds the journal's 1 MiB line budget") {
					t.Fatalf("oversized close: error=%v, stdout bytes=%d, stderr=%s", got.err, len(got.stdout), got.stderr)
				}
				if !before.equal(snapshotWorkspacePersistence(t, f.workspace)) {
					t.Fatal("oversized close changed the journal or projection")
				}
			}
			// A later write must still be able to scan the journal and allocate IDs.
			lifecycleResult(t, f.run("--json", "issue", "create", "--kind", "task", "--title", "After close", "--lane", "routine", "--actor", "human:mark", "--idempotency-key", "after-close"), "issue.create")
			if !tc.wantSuccess {
				// The rejected request must not consume its key or close the epic.
				args[6] = "Finished"
				lifecycleResult(t, f.run(args...), "task.done")
			}
		})
	}
}
