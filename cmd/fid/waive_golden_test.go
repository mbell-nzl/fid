package main

import (
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

func TestWaivePromptAndReplyGoldens(t *testing.T) {
	t.Parallel()
	f := laneFloorFixture(t, "standard")
	writeBackupTestFile(t, filepath.Join(f.workspace, "secure/task.go"), "package secure\n")
	runGitForTest(t, f.workspace, "add", ".")
	runGitForTest(t, f.workspace, "commit", "-m", "Configured Assured path")
	requiredNames(t, checkRequired(t, f, "checks"))
	reviewHandoff(t, f, "handoff")
	normalize := func(s string) string {
		s = f.placeholders(s)
		s = strings.ReplaceAll(s, " to waive: ", " to waive: <input>")
		s = regexp.MustCompile(`[0-9]{4}-[0-9]{2}-[0-9]{2}T[0-9:.]+Z`).ReplaceAllString(s, "<time>")
		s = regexp.MustCompile(`EVT-[0-9]{8}-[0-9]{6}`).ReplaceAllString(s, "EVT-<new>")
		s = regexp.MustCompile(`\b(?:PDN|HOF)-[0-9a-f]{6}\b`).ReplaceAllStringFunc(s, func(id string) string { return id[:3] + "-<new>" })
		return regexp.MustCompile(`\b(?:[0-9a-f]{64}|[0-9a-f]{40})\b`).ReplaceAllString(s, "<fingerprint-or-commit>")
	}
	got := waiveForTest(f, "FID-9f13aa", "FID-9f13aa\n", true, "--lane", "standard", "--idempotency-key", "waiver")
	waiveSuccess(t, got)
	assertGolden(t, filepath.Join("waive", "prompt.txt"), normalize(got.stderr))
	assertGolden(t, filepath.Join("waive", "waived.json"), normalize(got.stdout))
	refusal := waiveForTest(f, "FID-9f13aa", "", true, "--lane", "assured")
	assertGolden(t, filepath.Join("waive", "not-lower.json"), normalize(refusal.stderr))
}
