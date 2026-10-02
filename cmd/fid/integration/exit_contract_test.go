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
		removed  bool
	}{
		{"AC verification", []string{"--json", "issue", "ac", "verify", "FID-404", "AC-404", "--actor", "agent:task-9"}, "", 2, true},
		{"concern raise", []string{"--json", "concern", "raise", "--on", "FID-404", "--actor", "agent:task-9"}, "", 2, true},
		{"concern address", []string{"--json", "concern", "address", "CON-404", "--actor", "agent:task-9"}, "", 2, true},
		{"maintenance vacuum", []string{"--json", "maintenance", "vacuum", "--actor", "agent:task-9"}, "capability-unavailable", 5, false},
		{"invalid input", []string{"maintenance", "vacuum", "--by", "tl"}, "", 2, false},
		{"unknown flag", []string{"maintenance", "vacuum", "--unexpected"}, "", 2, false},
		{"usage", []string{"issue", "get"}, "", 2, false},
		{"not found", []string{"--json", "issue", "get", "FID-404"}, "issue-not-found", 3, false},
		{"conflict", []string{"--json", "init", "--prefix", "FID"}, "project-already-initialized", 4, false},
		{"unknown command", []string{"unknown-command"}, "", 2, false},
	} {
		test := test
		t.Run(test.name, func(t *testing.T) {
			before := p21JournalEventCount(t, workspace)
			stdout, stderr, err := closureRunCommand(workspace, binary, test.args...)
			if got := processExitCode(t, err); got != test.exitCode {
				t.Fatalf("%v exit code = %d, want %d; stdout=%q stderr=%q", test.args, got, test.exitCode, stdout, stderr)
			}
			if test.code != "" && !strings.Contains(stderr, `"code":"`+test.code+`"`) {
				t.Fatalf("%v stderr = %q, want code %q", test.args, stderr, test.code)
			}
			if test.removed && (stdout != "" || stderr != "" || p21JournalEventCount(t, workspace) != before) {
				t.Fatalf("removed leaf %v changed output or journal: stdout=%q stderr=%q", test.args, stdout, stderr)
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
