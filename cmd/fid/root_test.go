package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestPromisedGovernanceCommandsAreAbsent(t *testing.T) {
	t.Parallel()

	for _, commandName := range []string{"decision", "escalation", "triage"} {
		t.Run(commandName, func(t *testing.T) {
			stdout := &bytes.Buffer{}
			stderr := &bytes.Buffer{}
			cmd := newRootCommandAt(t.TempDir())
			cmd.SetOut(stdout)
			cmd.SetErr(stderr)
			cmd.SetArgs([]string{commandName})

			err := cmd.Execute()
			if err == nil {
				t.Fatalf("Execute(%q) error = nil, want unknown command", commandName)
			}
			if !strings.Contains(err.Error(), `unknown command "`+commandName+`" for "fid"`) {
				t.Fatalf("Execute(%q) error = %q, want Cobra unknown-command error", commandName, err)
			}
			if stdout.Len() != 0 {
				t.Fatalf("stdout = %q, want empty", stdout.String())
			}
			if stderr.Len() != 0 {
				t.Fatalf("stderr = %q, want empty", stderr.String())
			}
		})
	}
}
