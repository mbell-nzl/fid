package main

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestReviewRequiredProfilesSurviveInvalidUTF8Handoff(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name, path     string
		needTestReview bool
	}{
		{"invalid UTF-8", "notes\xff.md", true},
		{"readable replacement", "notes\ufffd.md", false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			f := requiredCheckFixture(t, "[]")
			// Neither readable .md name matches the configured test paths.
			reviewPolicy(t, f, "standard", "['**/*_test.go']")
			writeBackupTestFile(t, filepath.Join(f.workspace, tc.path), "review this change\n")
			runGitForTest(t, f.workspace, "add", ".")
			runGitForTest(t, f.workspace, "commit", "-m", "Review path bytes")
			reviewHandoff(t, f, "path-handoff")

			// Review reloads the hand-off from its persisted JSON receipt.
			in := reviewInput()
			in["profiles"] = []string{"code"}
			before := snapshotWorkspacePersistence(t, f.workspace)
			got := recordReview(t, f, in, "path-review", "human:courier")
			if tc.needTestReview {
				if got.err == nil || !strings.Contains(got.stderr, "missing required profile test-adequacy") {
					t.Fatalf("code-only review: err=%v stdout=%s stderr=%s", got.err, got.stdout, got.stderr)
				}
				if !before.equal(snapshotWorkspacePersistence(t, f.workspace)) {
					t.Fatal("refused review changed journal or projection")
				}
				in["profiles"] = []string{"code", "test-adequacy"}
				got = recordReview(t, f, in, "path-review", "human:courier")
			}
			lifecycleResult(t, got, "review.record")
			if after := startEventCount(t, f.workspace); after != before.eventCount()+1 {
				t.Fatalf("accepted review appended %d Events, want one", after-before.eventCount())
			}
		})
	}
}
