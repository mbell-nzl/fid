//go:build integration

package integration

import (
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
	"reflect"
	"strings"
	"testing"

	"github.com/mbell-nzl/fid/internal/gittest"
)

func TestProcessExitContract(t *testing.T) {
	t.Parallel()
	workspace := gittest.NewWorkspace(t)
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
			if test.removed {
				assertRemovedCommandFailure(t, test.args, stdout, stderr, err)
				if p21JournalEventCount(t, workspace) != before {
					t.Fatalf("removed leaf %v changed journal", test.args)
				}
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

// Retained absent-command calls keep their flags: parsing rejects an unknown
// flag before the group's unknown-word check. Neither rejection opens a project.
func assertRemovedCommandFailure(t *testing.T, args []string, stdout, stderr string, cause error) {
	t.Helper()
	if processExitCode(t, cause) != 2 || stdout != "" {
		t.Fatalf("removed command %v: stdout=%q stderr=%q cause=%v", args, stdout, stderr, cause)
	}
	jsonMode := args[0] == "--json"
	if jsonMode {
		args = args[1:]
	}
	group, word, allowed := args[0], args[1], []string{}
	switch group {
	case "policy":
		group, word = "fid", "policy"
		allowed = []string{"activation", "activity", "artifact", "attention", "capabilities", "claim", "completion", "concern", "decision", "deviation", "envelope", "evidence", "export", "finding", "gate", "help", "human-outcome", "init", "issue", "journal", "label", "link", "maintenance", "ready", "rebuild", "receipt", "requirement", "resume", "review", "run", "sql", "stats", "stop", "workflow"}
	case "issue":
		group, word, allowed = "issue.ac", args[2], []string{"add", "list"}
	case "concern":
		allowed = []string{"dispose", "get", "list", "record"}
	case "requirement":
		allowed = []string{"get", "link", "list", "record-reference"}
	}
	code, argument, value := "operation-unavailable", "command", word
	path := "fid " + strings.ReplaceAll(group, ".", " ")
	if group == "fid" {
		path = "fid"
	}
	message := fmt.Sprintf("unknown command %q for %q", word, path)
	for _, arg := range args {
		if strings.HasPrefix(arg, "--") {
			code, argument, value, allowed = "invalid-argument", arg, arg, []string{}
			message = "unknown flag: " + arg
			break
		}
	}
	if !jsonMode {
		detail := "argument=" + argument + ", allowed=" + strings.Join(allowed, "|")
		if code == "operation-unavailable" {
			detail = "allowed=" + strings.Join(allowed, ",") + ", argument=command, value=" + value
		}
		reason := "reason: " + code + "; " + detail
		// These retained invocations contain ASCII only; their independent
		// expectation counts the complete label and retains 200 scalars.
		if scalars := []rune(reason); len(scalars) > 200 {
			reason = string(scalars[:200]) + "…"
		}
		want := "error: " + message + "\n" + reason + "\n"
		if stderr != want {
			t.Fatalf("removed command %v stderr=%q, want %q", args, stderr, want)
		}
		return
	}
	want := map[string]any{
		"contract_version": "fid-cli.v1", "operation": group, "ok": false,
		"error": map[string]any{
			"code": code, "category": "invalid-input", "retryable": false, "subject": nil,
			"message": message, "details": map[string]any{"argument": argument, "value": value, "allowed": allowed},
		},
	}
	wantBytes, marshalErr := json.Marshal(want)
	if marshalErr != nil {
		t.Fatal(marshalErr)
	}
	var got, normalizedWant map[string]any
	if err := json.Unmarshal([]byte(stderr), &got); err != nil {
		t.Fatalf("removed command failure JSON: %v; stderr=%q", err, stderr)
	}
	if err := json.Unmarshal(wantBytes, &normalizedWant); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, normalizedWant) {
		t.Fatalf("removed command %v envelope=%s, want %s", args, stderr, wantBytes)
	}
}
