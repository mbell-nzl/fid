//go:build integration

package integration

import (
	"errors"
	"os/exec"
	"strings"
	"testing"
)

func TestProcessExitContract(t *testing.T) {
	t.Parallel()
	workspace := t.TempDir()
	binary := integrationFIDBinary(t)
	closureRun(t, workspace, binary, "init", "--prefix", "FID")

	for _, test := range []struct {
		name     string
		args     []string
		code     string
		exitCode int
	}{
		{"AC verification", []string{"--json", "issue", "ac", "verify", "FID-404", "AC-404", "--actor", "agent:task-9"}, "capability-unavailable", 5},
		{"concern raise", []string{"--json", "concern", "raise", "--on", "FID-404", "--actor", "agent:task-9"}, "capability-unavailable", 5},
		{"concern address", []string{"--json", "concern", "address", "CON-404", "--actor", "agent:task-9"}, "capability-unavailable", 5},
		{"maintenance vacuum", []string{"--json", "maintenance", "vacuum", "--actor", "agent:task-9"}, "capability-unavailable", 5},
		{"invalid input", []string{"maintenance", "vacuum", "--by", "tl"}, "", 2},
		{"unknown flag", []string{"maintenance", "vacuum", "--unexpected"}, "", 2},
		{"usage", []string{"issue", "get"}, "", 2},
		{"not found", []string{"--json", "issue", "get", "FID-404"}, "issue-not-found", 3},
		{"conflict", []string{"--json", "init", "--prefix", "FID"}, "project-already-initialized", 4},
		{"unknown command", []string{"unknown-command"}, "", 2},
	} {
		test := test
		t.Run(test.name, func(t *testing.T) {
			stdout, stderr, err := closureRunCommand(workspace, binary, test.args...)
			if got := processExitCode(t, err); got != test.exitCode {
				t.Fatalf("%v exit code = %d, want %d; stdout=%q stderr=%q", test.args, got, test.exitCode, stdout, stderr)
			}
			if test.code != "" && !strings.Contains(stderr, `"code":"`+test.code+`"`) {
				t.Fatalf("%v stderr = %q, want code %q", test.args, stderr, test.code)
			}
		})
	}
}

func processExitCode(t *testing.T, err error) int {
	t.Helper()
	var exitError *exec.ExitError
	if !errors.As(err, &exitError) {
		t.Fatalf("process error = %v, want exit error", err)
	}
	return exitError.ExitCode()
}
