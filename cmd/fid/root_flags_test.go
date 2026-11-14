package main

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mbell-nzl/fid/internal/cli"
	"github.com/mbell-nzl/fid/internal/config"
	"github.com/mbell-nzl/fid/internal/gittest"
)

const p24UnavailableProjectJSON = `{"contract_version":"fid-cli.v1","error":{"category":"invalid-input","code":"invalid-argument","details":{"allowed":["--idempotency-key","--repository"],"argument":"--project","value":"--project"},"message":"--project is not available for export import","retryable":false,"subject":null},"ok":false,"operation":"export.import"}` + "\n"

func TestP24UnavailableRootFlagMapping(t *testing.T) {
	t.Parallel()
	for _, missing := range []bool{false, true} {
		for _, jsonMode := range []bool{false, true} {
			root := t.TempDir()
			project := root
			if missing {
				project = filepath.Join(root, "nonexistent")
			}
			args := []string{"--project", project, "export", "import", "missing.zip", "--repository", root, "--idempotency-key", "refusal"}
			want := "error: --project is not available for export import\nreason: invalid-argument; argument=--project, allowed=--idempotency-key|--repository\n"
			if jsonMode {
				args = append([]string{"--json"}, args...)
				want = p24UnavailableProjectJSON
			}
			result := executeCommandForTest(root, args...)
			if result.err == nil || cli.ExitCode(result.err) != 2 || result.stdout != "" || result.stderr != want {
				t.Fatalf("missing=%t json=%t: %#v, want stderr=%q exit=2", missing, jsonMode, result, want)
			}
			entries, err := os.ReadDir(root)
			if err != nil || len(entries) != 0 {
				t.Fatalf("refusal changed directory: %v %v", entries, err)
			}
		}
	}
}

func TestP24DiagnosticRootFlagBytesPreserved(t *testing.T) {
	t.Parallel()
	for _, leaf := range [][]string{{"sql", "SELECT 1"}, {"journal", "inspect", "EVT-missing"}} {
		for _, flag := range []string{"--json", "--project"} {
			args := []string{flag}
			if flag == "--project" {
				args = append(args, filepath.Join(t.TempDir(), "missing"))
			}
			result := executeCommandForTest(t.TempDir(), append(args, leaf...)...)
			path := strings.Join(leaf[:len(leaf)-1], " ")
			want := "error: " + flag + " is not available for " + path + " until its fid-cli.v1 transport is migrated\n"
			if result.err == nil || cli.ExitCode(result.err) != 7 || result.stdout != "" || result.stderr != want {
				t.Fatalf("diagnostic refusal: %#v want %q", result, want)
			}
		}
	}
}

func TestRootCommandUsesInjectedWorkingDirectoryWithoutChangingProcess(t *testing.T) {
	t.Parallel()
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
	t.Parallel()
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
	t.Parallel()
	workspace := gittest.NewNamedWorkspace(t, "relative-project")
	cmd := newRootCommandAt(filepath.Dir(workspace))
	cmd.SetOut(&bytes.Buffer{})
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetArgs([]string{"--project", "relative-project", "init", "--prefix", "FID"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if _, err := os.Stat(config.LedgerDir(workspace)); err != nil {
		t.Fatalf("relative project workspace not initialized: %v", err)
	}
}

func TestRootFlagsAllowRebuildForAnExplicitProject(t *testing.T) {
	t.Parallel()
	workspace := gittest.NewWorkspace(t)
	runInitCommand(t, workspace, "init")

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	cmd := newRootCommandAt(t.TempDir())
	cmd.SetOut(stdout)
	cmd.SetErr(stderr)
	cmd.SetArgs([]string{"--project", workspace, "--json", "rebuild"})
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
	t.Parallel()
	for _, fixture := range []struct {
		name string
		args []string
		flag string
	}{
		{"legacy json", []string{"--json", "journal", "inspect", "EVT-missing"}, "--json"},
		{"unavailable maintenance", []string{"--project", t.TempDir(), "maintenance", "vacuum"}, "capability"},
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
			if !strings.Contains(err.Error(), fixture.flag+" is not available") && !strings.Contains(err.Error(), "is unavailable") {
				t.Fatalf("Execute() error = %q, want explicit %s rejection", err, fixture.flag)
			}
			if got := stderr.String(); !strings.Contains(got, "error: "+fixture.flag+" is not available") && !strings.Contains(got, "error: capability is unavailable") {
				t.Fatalf("stderr = %q, want concise explicit %s rejection", got, fixture.flag)
			}
		})
	}
}

func TestRootFlagsAcceptConvertedReadsAndRejectOnlyLegacyDiagnostics(t *testing.T) {
	t.Parallel()
	for _, fixture := range []struct {
		name    string
		command []string
	}{
		{"activity list", []string{"activity", "list"}},
	} {
		for _, rootFlag := range []string{"--json", "--project"} {
			t.Run(fixture.name+" "+rootFlag, func(t *testing.T) {
				workspace := clonePristineWorkspace(t)
				args := append([]string{rootFlag}, fixture.command...)
				if rootFlag == "--project" {
					args = append([]string{rootFlag, workspace, "--json"}, fixture.command...)
				}
				result := executeCommandForTest(workspace, args...)
				if result.err != nil || result.stderr != "" || result.stdout == "" {
					t.Fatalf("%s %s = %#v, want accepted JSON read", fixture.name, rootFlag, result)
				}
			})
		}
	}

	for _, fixture := range []struct {
		name    string
		command []string
	}{
		{"sql", []string{"sql", "SELECT 1"}},
		{"journal inspect", []string{"journal", "inspect", "EVT-missing"}},
	} {
		for _, rootFlag := range []string{"--json", "--project"} {
			t.Run(fixture.name+" "+rootFlag, func(t *testing.T) {
				stdout := &bytes.Buffer{}
				stderr := &bytes.Buffer{}
				cmd := newRootCommandAt(t.TempDir())
				cmd.SetOut(stdout)
				cmd.SetErr(stderr)

				args := append([]string{rootFlag}, fixture.command...)
				if rootFlag == "--project" {
					args = append([]string{rootFlag, t.TempDir()}, fixture.command...)
				}
				cmd.SetArgs(args)

				err := cmd.Execute()
				if err == nil {
					t.Fatal("Execute() error = nil, want root-flag rejection")
				}
				if stdout.Len() != 0 {
					t.Fatalf("stdout = %q, want no handler output", stdout.String())
				}
				if !strings.Contains(stderr.String(), "error: "+rootFlag+" is not available") {
					t.Fatalf("stderr = %q, want explicit %s rejection", stderr.String(), rootFlag)
				}
			})
		}
	}

	for _, rootFlag := range []string{"--json", "--project"} {
		t.Run("policy absent "+rootFlag, func(t *testing.T) {
			args := []string{rootFlag, "policy", "show"}
			if rootFlag == "--project" {
				args = []string{rootFlag, clonePristineWorkspace(t), "--json", "policy", "show"}
			}
			result := executeCommandForTest(clonePristineWorkspace(t), args...)
			if result.err == nil || result.stdout != "" || !strings.Contains(result.err.Error(), `unknown command "policy"`) {
				t.Fatalf("policy show %s = %#v, want absent command", rootFlag, result)
			}
		})
	}
}
