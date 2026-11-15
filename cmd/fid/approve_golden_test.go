package main

import (
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

func TestApprovePromptGolden(t *testing.T) {
	t.Parallel()
	f := requiredCheckFixture(t, "['always', 'empty']")
	writeBackupTestFile(t, filepath.Join(f.workspace, "notes.md"), "review this change\n")
	runGitForTest(t, f.workspace, "add", ".")
	runGitForTest(t, f.workspace, "commit", "-m", "Change for approval")
	requiredNames(t, checkRequired(t, f, "checks"))
	reviewHandoff(t, f, "approval-handoff")
	in := reviewInput()
	lifecycleResult(t, recordReview(t, f, in, "review", "human:courier"), "review.record")
	normalize := func(s string) string {
		s = f.placeholders(s)
		// Mark the input position while preserving the prompt's literal space.
		s = strings.ReplaceAll(s, " to approve: ", " to approve: <input>")
		s = regexp.MustCompile(`[0-9]{4}-[0-9]{2}-[0-9]{2}T[0-9:.]+Z`).ReplaceAllString(s, "<time>")
		s = regexp.MustCompile(`EVT-[0-9]{8}-[0-9]{6}`).ReplaceAllString(s, "EVT-<new>")
		s = regexp.MustCompile(`\b(?:RND|HOF|CHK|PDN)-[0-9a-f]{6}\b`).ReplaceAllStringFunc(s, func(id string) string { return id[:3] + "-<new>" })
		return regexp.MustCompile(`\b(?:[0-9a-f]{64}|[0-9a-f]{40})\b`).ReplaceAllString(s, "<fingerprint-or-commit>")
	}
	first := approveForTest(f, "FID-9f13aa", "FID-9f13aa\n", true, "--idempotency-key", "golden-approval")
	approveSuccess(t, first)
	assertGolden(t, filepath.Join("approve", "prompt.txt"), normalize(first.stderr))
	assertGolden(t, filepath.Join("approve", "approved.json"), normalize(first.stdout))
	again := approveForTest(f, "FID-9f13aa", "no\n", true)
	assertGolden(t, filepath.Join("approve", "earlier-current.txt"), normalize(again.stderr))
}
