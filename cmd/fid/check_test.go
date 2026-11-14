package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func leanCheckFixture(t *testing.T) *statusFixture {
	t.Helper()
	f := plannedStartFixture(t, "in_progress", stringPointer("standard"))
	appendPolicy(t, f.workspace, "\n[checks.lean]\nrun = ['true']\ntimeout = '15m'\nartifacts = ['coverage.out']\n")
	policyCommit(t, f.workspace)
	identityGit(t, f.workspace, "push", "-q", "origin", "trunk")
	if err := os.WriteFile(filepath.Join(f.workspace, "coverage.out"), []byte("kept"), 0o600); err != nil {
		t.Fatal(err)
	}
	runGitForTest(t, f.workspace, "add", ".")
	runGitForTest(t, f.workspace, "-c", "user.name=fid test", "-c", "user.email=fid@example.test", "-c", "commit.gpgsign=false", "commit", "-m", "Check fixture")
	return f
}

func TestLeanCheckRecordsResultAndReplaysWithoutSpawning(t *testing.T) {
	t.Parallel()
	f := leanCheckFixture(t)
	args := []string{"--json", "check", "lean", "--hold", "HLD-aaaaaa", "--actor", "agent:claude-builder", "--idempotency-key", "check-1"}
	first := f.run(args...)
	if first.err != nil {
		t.Fatalf("check=%#v", first)
	}
	var envelope struct {
		Result struct {
			TaskID string `json:"task_id"`
			Checks []struct {
				ID        string   `json:"id"`
				Counts    bool     `json:"counts"`
				Artifacts []string `json:"artifact_ids"`
			} `json:"checks"`
		} `json:"result"`
	}
	if err := json.Unmarshal([]byte(first.stdout), &envelope); err != nil {
		t.Fatal(err)
	}
	if envelope.Result.TaskID != "FID-9f13aa" || len(envelope.Result.Checks) != 1 || !envelope.Result.Checks[0].Counts || len(envelope.Result.Checks[0].Artifacts) != 1 {
		t.Fatalf("response=%s", first.stdout)
	}
	before := startEventCount(t, f.workspace)
	if err := os.Remove(filepath.Join(f.workspace, "coverage.out")); err != nil {
		t.Fatal(err)
	}
	second := f.run(args...)
	if second.err != nil || second.stdout != first.stdout || startEventCount(t, f.workspace) != before {
		t.Fatalf("retry=%#v; first=%#v", second, first)
	}
	shown := f.run("--json", "show", "FID-9f13aa")
	if shown.err != nil || !strings.Contains(shown.stdout, `"checks":`) || !strings.Contains(shown.stdout, envelope.Result.Checks[0].ID) {
		t.Fatalf("show=%#v", shown)
	}
}

func TestLeanCheckDefaultsToAnEmptyRequiredSetAndRecordsTheRetry(t *testing.T) {
	t.Parallel()
	f := plannedStartFixture(t, "in_progress", stringPointer("standard"))
	got := f.run("--json", "check", "--hold", "HLD-aaaaaa", "--actor", "agent:claude-builder", "--idempotency-key", "check-1")
	if got.err != nil || !strings.Contains(got.stdout, `"checks":[]`) || !strings.Contains(got.stdout, `"event":{"id":`) {
		t.Fatalf("empty check=%#v", got)
	}
	before := startEventCount(t, f.workspace)
	identityGit(t, f.workspace, "remote", "remove", "origin")
	repeat := f.run("--json", "check", "--hold", "HLD-aaaaaa", "--actor", "agent:claude-builder", "--idempotency-key", "check-1")
	if repeat.err != nil || repeat.stdout != got.stdout || startEventCount(t, f.workspace) != before {
		t.Fatalf("empty retry=%#v, original=%#v", repeat, got)
	}
}
