package main

import (
	"path/filepath"
	"reflect"
	"testing"

	"github.com/mbell-nzl/fid/internal/gittest"
)

func TestWaivedReviewSelectionProfiles(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		state    string
		profiles []any
	}{
		{"current", []any{"test-adequacy"}},
		{"unchanged", []any{"test-adequacy"}},
		{"stale", []any{"code", "test-adequacy"}},
	} {
		t.Run(tc.state, func(t *testing.T) {
			t.Parallel()
			f, task := waivedRoutineTestAdequacyFixture(t)
			if tc.state != "current" {
				// Keep primary review due even when the unchanged round stays current.
				in := changesInput("separate")
				in["profiles"] = []string{"test-adequacy"}
				in["findings"].([]any)[0].(map[string]any)["profile"] = "test-adequacy"
				lifecycleResult(t, recordReviewTask(t, f, in, "changes", "human:courier", task), "review.record")
				prepareSecondReviewHandoff(t, f, task)
				if tc.state == "stale" {
					writeBackupTestFile(t, filepath.Join(f.workspace, "notes.md"), "edited waived change\n")
					runGitForTest(t, f.workspace, "add", ".")
					runGitForTest(t, f.workspace, "commit", "-m", "Edit waived change")
				}
				reviewHandoffTask(t, f, "next-handoff", task)
			}
			f.workspace = gittest.AddWorktree(t, f.workspace)
			got := selectionNext(t, f, "")
			review, reviewOK := got["review"].(map[string]any)
			selected, taskOK := got["task"].(map[string]any)
			if got["kind"] != "review" || !reviewOK || !taskOK || selected["id"] != task || review["stream"] != "primary" {
				t.Fatalf("want primary review for %s: %v", task, got)
			}
			if !reflect.DeepEqual(review["required_profiles"], tc.profiles) {
				t.Fatalf("%s waiver: next required_profiles = %v, want %v", tc.state, review["required_profiles"], tc.profiles)
			}
		})
	}
}
