package main

import (
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestLiveGoldenChildSkipsPristineBootstrap(t *testing.T) {
	t.Parallel()
	if containsTestArgument(os.Args, "--live-golden-child") {
		if pristineWorkspaceRoot != "" {
			t.Fatalf("live golden child initialized unused pristine workspace %q", pristineWorkspaceRoot)
		}
		return
	}
	cmd := exec.Command(os.Args[0], "-test.run=^TestLiveGoldenChildSkipsPristineBootstrap$", "--", "--live-golden-child")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("live golden child harness: %v; output: %s", err, strings.TrimSpace(string(output)))
	}
}
