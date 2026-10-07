package main

import (
	"bytes"
	"strings"
	"testing"

	"github.com/mbell-nzl/fid/internal/cli"
)

func TestPromisedGovernanceCommandsAreAbsent(t *testing.T) {
	t.Parallel()

	for _, commandName := range []string{"escalation", "triage"} {
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
			allowed := []string{}
			for _, child := range cmd.Commands() {
				if child.IsAvailableCommand() || child.Name() == "help" {
					allowed = append(allowed, child.Name())
				}
			}
			row := grammarRow{operation: "fid", code: "operation-unavailable", argument: "command", value: commandName,
				allowed: allowed, message: `unknown command "` + commandName + `" for "fid"`}
			if stderr.String() != p24GrammarHumanText(t, row) || cli.ExitCode(err) != 2 {
				t.Fatalf("stderr = %q exit=%d, want typed unknown-command rejection", stderr.String(), cli.ExitCode(err))
			}
		})
	}
}
