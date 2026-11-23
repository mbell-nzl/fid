package main

import (
	"fmt"
	"testing"

	"github.com/mbell-nzl/fid/internal/journal"
)

func TestWaivedDefaultConflictSelection(t *testing.T) {
	t.Parallel()
	for _, lane := range []string{"routine", "standard"} {
		for _, conflict := range []bool{false, true} {
			t.Run(fmt.Sprintf("%s/conflict=%t", lane, conflict), func(t *testing.T) {
				t.Parallel()
				f, task := canonicalReviewFixtureWithLane(t, "assured")
				waiveSuccess(t, waiveForTest(f, task, task+"\n", true, "--lane", lane))
				if lane == "standard" {
					approveConflictReviews(t, f, task, lane)
				}
				moveConflictDefault(t, f, conflict)
				status, _ := lifecycleResult(t, f.run("--json", "status"), "project.status")
				wantAction, wantWork := "merge", "none"
				if conflict {
					wantAction, wantWork = "default-conflict", "rework"
				}
				found := false
				for _, raw := range status["needs_you"].([]any) {
					item := raw.(map[string]any)
					subject, ok := item["subject"].(map[string]any)
					if ok && subject["id"] == task {
						if item["kind"] != wantAction {
							t.Fatalf("waived %s action=%v, want %s", lane, item, wantAction)
						}
						found = true
					}
				}
				if !found {
					t.Fatalf("waived %s has no %s action: %v", lane, wantAction, status)
				}
				got, _ := lifecycleResult(t, f.run("--json", "next", "--claim", "--actor", "agent:builder", "--model", "codex/gpt-6.1-sol"), "next.claim")
				if got["kind"] != wantWork {
					t.Fatalf("waived %s work=%v, want %s", lane, got, wantWork)
				}
				if _, err := journal.RebuildWorkspace(f.workspace, issueMigrationsDir()); err != nil {
					t.Fatalf("waiver and default-conflict replay: %v", err)
				}
			})
		}
	}
}
