package main

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRootCommandUsesInjectedWorkingDirectoryWithoutChangingProcess(t *testing.T) {
	workspace := clonePristineWorkspace(t)
	originalDirectory, err := os.Getwd()
	if err != nil {
		t.Fatalf("Getwd() error = %v", err)
	}

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	cmd := newRootCommandAt(workspace)
	cmd.SetOut(stdout)
	cmd.SetErr(stderr)
	cmd.SetArgs([]string{"issue", "list"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if stderr.Len() != 0 {
		t.Fatalf("stderr = %q, want empty", stderr.String())
	}
	if stdout.String() != "" {
		t.Fatalf("stdout = %q, want empty human issue list", stdout.String())
	}
	afterDirectory, err := os.Getwd()
	if err != nil {
		t.Fatalf("Getwd() after command error = %v", err)
	}
	if afterDirectory != originalDirectory {
		t.Fatalf("process working directory = %q, want unchanged %q", afterDirectory, originalDirectory)
	}
}

func TestRootCommandCapturesCurrentWorkingDirectory(t *testing.T) {
	want, err := os.Getwd()
	if err != nil {
		t.Fatalf("Getwd() error = %v", err)
	}
	cmd, err := newRootCommand()
	if err != nil {
		t.Fatalf("newRootCommand() error = %v", err)
	}
	got, err := commandWorkingDirectory(cmd)
	if err != nil {
		t.Fatalf("commandWorkingDirectory() error = %v", err)
	}
	if got != want {
		t.Fatalf("command working directory = %q, want %q", got, want)
	}
}

func TestRootCommandResolvesRelativeProjectFromInjectedWorkingDirectory(t *testing.T) {
	base := t.TempDir()
	workspace := filepath.Join(base, "relative-project")
	if err := os.Mkdir(workspace, 0o755); err != nil {
		t.Fatalf("Mkdir() error = %v", err)
	}
	cmd := newRootCommandAt(base)
	cmd.SetOut(&bytes.Buffer{})
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetArgs([]string{"--project", "relative-project", "init", "--prefix", "FID"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if _, err := os.Stat(filepath.Join(workspace, ".fid")); err != nil {
		t.Fatalf("relative project workspace not initialized: %v", err)
	}
}

func TestRootFlagsAllowRebuildForAnExplicitProject(t *testing.T) {
	workspace := t.TempDir()
	runInitCommand(t, workspace, "init")

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	cmd := newRootCommandAt(t.TempDir())
	cmd.SetOut(stdout)
	cmd.SetErr(stderr)
	cmd.SetArgs([]string{"--project", workspace, "--json", "rebuild", "--by", "workflow_runner"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if stderr.Len() != 0 {
		t.Fatalf("stderr = %q, want empty", stderr.String())
	}
	var envelope struct {
		OK        bool   `json:"ok"`
		Operation string `json:"operation"`
		Event     *any   `json:"event"`
	}
	if err := json.Unmarshal(stdout.Bytes(), &envelope); err != nil {
		t.Fatalf("Unmarshal(stdout) error = %v", err)
	}
	if !envelope.OK || envelope.Operation != "project.rebuild" || envelope.Event != nil {
		t.Fatalf("rebuild envelope = %+v, want project.rebuild read response", envelope)
	}
}

func TestRootFlagsRejectBeforeLegacyAndUnmigratedLeafHandlers(t *testing.T) {
	for _, fixture := range []struct {
		name string
		args []string
		flag string
	}{
		{"legacy json", []string{"--json", "journal", "inspect", "EVT-missing"}, "--json"},
		{"legacy project", []string{"--project", t.TempDir(), "maintenance", "vacuum"}, "--project"},
	} {
		fixture := fixture
		t.Run(fixture.name, func(t *testing.T) {
			stdout := &bytes.Buffer{}
			stderr := &bytes.Buffer{}
			cmd := newRootCommandAt(t.TempDir())
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
