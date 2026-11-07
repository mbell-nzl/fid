package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mbell-nzl/fid/internal/config"
	"github.com/mbell-nzl/fid/internal/gittest"
)

func TestInitBackupTargetPrompt(t *testing.T) { //nolint:paralleltest // t.Setenv isolates user config for the CLI.
	for _, test := range []struct {
		name, input, existing, wantTarget string
		stdinTTY, stdoutTTY, json, prompt bool
	}{
		{name: "terminal URL", stdinTTY: true, stdoutTTY: true, input: "git@example.test:private/ledgers.git\n", prompt: true, wantTarget: "git@example.test:private/ledgers.git"},
		{name: "terminal local path", stdinTTY: true, stdoutTTY: true, input: "/tmp/backups with spaces/#ledger\n", prompt: true, wantTarget: "/tmp/backups with spaces/#ledger"},
		{name: "empty answer skips", stdinTTY: true, stdoutTTY: true, input: "\n", prompt: true},
		{name: "EOF skips", stdinTTY: true, stdoutTTY: true, prompt: true},
		{name: "piped input", stdoutTTY: true, input: "must not read\n"},
		{name: "redirected output", stdinTTY: true, input: "must not read\n"},
		{name: "no terminal", input: "must not read\n"},
		{name: "noninteractive ignores malformed config", existing: "[backup]\ntarget = broken\n", input: "must not read\n"},
		{name: "JSON suppresses prompt", stdinTTY: true, stdoutTTY: true, json: true, input: "must not read\n"},
		{name: "JSON ignores malformed config", stdinTTY: true, stdoutTTY: true, json: true, existing: "[backup]\ntarget = broken\n", input: "must not read\n"},
		{name: "configured target", stdinTTY: true, stdoutTTY: true, existing: "# keep this file\n[backup]\ntarget = 'existing:ledgers'\n", wantTarget: "existing:ledgers", input: "must not read\n"},
		{name: "empty configured target", stdinTTY: true, stdoutTTY: true, existing: "[backup]\ntarget = \"\"\n", input: "../backup.git\n", prompt: true, wantTarget: "../backup.git"},
	} {
		t.Run(test.name, func(t *testing.T) {
			configHome := t.TempDir()
			t.Setenv("XDG_CONFIG_HOME", configHome)
			t.Setenv("HOME", t.TempDir())
			path := filepath.Join(configHome, "fid", "config.toml")
			if test.existing != "" {
				if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(path, []byte(test.existing), 0o600); err != nil {
					t.Fatal(err)
				}
			}
			workspace := gittest.NewWorkspace(t)
			input := &initTestInput{Reader: strings.NewReader(test.input)}
			output := &initTestOutput{}
			stderr := &bytes.Buffer{}
			root := newRootCommandAt(workspace)
			old, _, err := root.Find([]string{"init"})
			if err != nil {
				t.Fatal(err)
			}
			root.RemoveCommand(old)
			root.AddCommand(newInitCommandWithTerminalCheck(func(fd uintptr) bool {
				return fd == input.Fd() && test.stdinTTY || fd == output.Fd() && test.stdoutTTY
			}))
			root.SetIn(input)
			root.SetOut(output)
			root.SetErr(stderr)
			args := []string{"init", "--prefix", "FID"}
			if test.json {
				args = append(args, "--json")
			}
			root.SetArgs(args)
			if err := root.Execute(); err != nil {
				t.Fatalf("init: %v; stderr=%q", err, stderr)
			}
			if got := strings.Contains(output.String(), "Backup repository"); got != test.prompt {
				t.Fatalf("prompt present = %t, want %t; output=%q", got, test.prompt, output)
			}
			if !test.prompt && input.Len() != len(test.input) {
				t.Fatal("init consumed stdin without prompting")
			}
			payload, err := os.ReadFile(path)
			if !test.prompt && test.existing != "" {
				if err != nil || string(payload) != test.existing {
					t.Fatalf("init changed existing config: %q, %v", payload, err)
				}
			} else if test.wantTarget == "" {
				if !os.IsNotExist(err) {
					t.Fatalf("skipped config = %q, %v; want absent", payload, err)
				}
			} else {
				if err != nil || !strings.Contains(string(payload), test.wantTarget) {
					t.Fatalf("user config = %q, %v; want target %q", payload, err, test.wantTarget)
				}
				if !test.prompt && string(payload) != test.existing {
					t.Fatal("init rewrote configured target")
				}
			}
			if stderr.Len() != 0 {
				t.Fatalf("stderr = %q", stderr)
			}
		})
	}
}

func TestInitRefusesMalformedInteractiveUserConfig(t *testing.T) { //nolint:paralleltest // t.Setenv isolates user config.
	configHome := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", configHome)
	t.Setenv("HOME", t.TempDir())
	path := filepath.Join(configHome, "fid", "config.toml")
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		t.Fatal(err)
	}
	contents := "[backup]\ntarget = broken\n"
	if err := os.WriteFile(path, []byte(contents), 0o600); err != nil {
		t.Fatal(err)
	}
	workspace := gittest.NewWorkspace(t)
	cmd := newInitCommandWithTerminalCheck(func(uintptr) bool { return true })
	root := newRootCommandAt(workspace)
	old, _, err := root.Find([]string{"init"})
	if err != nil {
		t.Fatal(err)
	}
	root.RemoveCommand(old)
	root.AddCommand(cmd)
	input, output, stderr := &initTestInput{Reader: strings.NewReader("must not read\n")}, &initTestOutput{}, &bytes.Buffer{}
	root.SetIn(input)
	root.SetOut(output)
	root.SetErr(stderr)
	root.SetArgs([]string{"init", "--prefix", "FID"})
	if err := root.Execute(); err == nil || !strings.Contains(stderr.String(), path) {
		t.Fatalf("init = %v, stderr=%q; want config refusal", err, stderr)
	}
	if output.Len() != 0 || input.Len() != len("must not read\n") {
		t.Fatalf("refused init prompted or consumed input: %q", output)
	}
	entries, err := os.ReadDir(workspace)
	if err != nil || len(entries) != 1 || entries[0].Name() != ".git" {
		t.Fatalf("refused init changed worktree: %v, %v", entries, err)
	}
	if _, err := os.Lstat(config.LedgerDir(workspace)); !os.IsNotExist(err) {
		t.Fatalf("refused init created ledger: %v", err)
	}
	if got, err := os.ReadFile(path); err != nil || string(got) != contents {
		t.Fatalf("refused init changed config: %q, %v", got, err)
	}
}

type initTestInput struct{ *strings.Reader }

func (*initTestInput) Fd() uintptr { return 10 }

type initTestOutput struct{ bytes.Buffer }

func (*initTestOutput) Fd() uintptr { return 11 }
