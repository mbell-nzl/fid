package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestRootFlagsRejectBeforeLegacyAndUnmigratedLeafHandlers(t *testing.T) {
	for _, fixture := range []struct {
		name string
		args []string
		flag string
	}{
		{"legacy json", []string{"--json", "journal", "inspect", "EVT-missing"}, "--json"},
		{"legacy project", []string{"--project", t.TempDir(), "maintenance", "vacuum"}, "--project"},
		{"unmigrated target", []string{"--project", t.TempDir(), "init"}, "--project"},
	} {
		fixture := fixture
		t.Run(fixture.name, func(t *testing.T) {
			stdout := &bytes.Buffer{}
			stderr := &bytes.Buffer{}
			cmd := newRootCommand()
			cmd.SetOut(stdout)
			cmd.SetErr(stderr)
			cmd.SetArgs(fixture.args)

			err := cmd.Execute()
			if err == nil {
				t.Fatal("Execute() error = nil, want root-flag rejection")
			}
			if stdout.Len() != 0 {
				t.Fatalf("stdout = %q, want no handler output", stdout.String())
			}
			if !strings.Contains(err.Error(), fixture.flag+" is not available") {
				t.Fatalf("Execute() error = %q, want explicit %s rejection", err, fixture.flag)
			}
			if got := stderr.String(); !strings.Contains(got, "error: "+fixture.flag+" is not available") {
				t.Fatalf("stderr = %q, want concise explicit %s rejection", got, fixture.flag)
			}
		})
	}
}
