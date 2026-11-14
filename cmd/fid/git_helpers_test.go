package main

import (
	"os"
	"os/exec"
	"testing"
)

func runGitForTest(t *testing.T, workspace string, args ...string) {
	t.Helper()
	command := exec.Command("git", args...)
	command.Dir = workspace
	command.Env = append(os.Environ(), "GIT_AUTHOR_NAME=fid test", "GIT_AUTHOR_EMAIL=fid@example.test",
		"GIT_COMMITTER_NAME=fid test", "GIT_COMMITTER_EMAIL=fid@example.test")
	if output, err := command.CombinedOutput(); err != nil {
		t.Fatalf("git %v: %v: %s", args, err, output)
	}
}
