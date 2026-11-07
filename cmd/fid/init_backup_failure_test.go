package main

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mbell-nzl/fid/internal/cli"
	"github.com/mbell-nzl/fid/internal/config"
	"github.com/mbell-nzl/fid/internal/gittest"
)

func TestInitReportsIdentityWhenBackupTargetSaveFails(t *testing.T) { //nolint:paralleltest // t.Setenv isolates HOME and user config.
	configHome := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", configHome)
	t.Setenv("HOME", t.TempDir())
	directory := filepath.Join(configHome, "fid")
	if err := os.Mkdir(directory, 0o500); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := os.Chmod(directory, 0o700); err != nil {
			t.Error(err)
		}
	})
	if file, err := os.CreateTemp(directory, "permission-probe-"); !errors.Is(err, os.ErrPermission) {
		if err == nil {
			_ = file.Close()
			_ = os.Remove(file.Name())
		}
		t.Fatalf("config directory must be unwritable, got %v", err)
	}
	workspace := gittest.NewWorkspace(t)
	root := newRootCommandAt(workspace)
	old, _, err := root.Find([]string{"init"})
	if err != nil {
		t.Fatal(err)
	}
	root.RemoveCommand(old)
	root.AddCommand(newInitCommandWithTerminalCheck(func(uintptr) bool { return true }))
	stdout, stderr := &initTestOutput{}, &bytes.Buffer{}
	root.SetIn(&initTestInput{Reader: strings.NewReader("../backups.git\n")})
	root.SetOut(stdout)
	root.SetErr(stderr)
	root.SetArgs([]string{"init", "--prefix", "FID"})
	if err := root.Execute(); cli.ExitCode(err) != 7 {
		t.Fatalf("init = %v, want save failure exit 7; stderr=%q", err, stderr)
	}
	id, err := config.ReadProjectIdentity(workspace)
	if err != nil {
		t.Fatal(err)
	}
	identityPath := filepath.Join(gittest.ResolvedPath(t, workspace), config.ProjectIdentityFileName)
	for _, want := range []string{"initialized fid project in ", "project ID: " + id, "commit " + identityPath + " to keep this identity in clones"} {
		if !strings.Contains(stdout.String(), want) {
			t.Errorf("stdout = %q, want %q after target-save failure", stdout, want)
		}
	}
	if !strings.Contains(stderr.String(), "ledger initialized, but backup target setup failed:") || strings.Contains(stderr.String(), "project initialization failed") {
		t.Errorf("stderr = %q, want an unambiguous post-ledger failure", stderr)
	}
	project, err := config.ReadProjectConfig(config.WorkspacePathsForRoot(workspace).ConfigPath)
	if err != nil || project.ProjectID != id {
		t.Errorf("initialized ledger = %+v, %v, want identity %q", project, err, id)
	}
	if _, err := os.Lstat(filepath.Join(directory, "config.toml")); !os.IsNotExist(err) {
		t.Errorf("failed save published config: %v", err)
	}
}
