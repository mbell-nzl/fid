package main

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/mbell-nzl/fid/internal/gittest"
)

func TestCommandGoTelemetryIsOff(t *testing.T) {
	t.Parallel()
	output, err := exec.Command("go", "env", "GOTELEMETRY").CombinedOutput()
	if err != nil || strings.TrimSpace(string(output)) != "off" {
		t.Fatalf("go env GOTELEMETRY = %q, %v; want off in the suite's isolated config", output, err)
	}
}

// Re-execute the real TestMain, rather than a substitute for suite setup.
// Each child owns a private TMPDIR; everything it leaves there is a leak.
func TestCommandSuiteRemovesTemporaryDirectories(t *testing.T) {
	t.Parallel()
	binary, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	for _, row := range []struct {
		name string
		args []string
		env  []string
	}{
		{"go-env", []string{"-test.run=^TestCommandTestsKeepTheSharedGoBuildCache$"}, nil},
		{"as-fid", []string{"--help"}, []string{runAsFIDVariable + "=1"}},
	} {
		t.Run(row.name, func(t *testing.T) {
			t.Parallel()
			scratch := t.TempDir()
			command := exec.CommandContext(t.Context(), binary, row.args...)
			command.Env = append(os.Environ(), "TMPDIR="+scratch, "GOCOVERDIR="+t.TempDir())
			command.Env = append(command.Env, row.env...)
			if output, err := command.CombinedOutput(); err != nil {
				t.Fatalf("child suite: %v\n%s", err, output)
			}
			assertNoSuiteTempEntries(t, scratch)
		})
	}
}

func TestLeanCheckChildrenSkipSuiteSetup(t *testing.T) {
	t.Parallel()
	binary, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	for _, mode := range []string{"pass", "fail", "signal", "wait", "pause", "change"} {
		t.Run(mode, func(t *testing.T) {
			t.Parallel()
			scratch, workspace := t.TempDir(), gittest.NewWorkspace(t)
			ctx, cancel := context.WithTimeout(t.Context(), 2*time.Second)
			defer cancel()
			command := exec.CommandContext(ctx, binary, "-test.run=^TestLeanCheckChild$", "--", "--check-child", mode)
			command.Dir = workspace
			// Exercise the coverage fallback too: checks deliberately strip
			// GOCOVERDIR, so any fallback must belong to the parent fixture.
			command.Env = append(os.Environ(), "TMPDIR="+scratch, "GOCOVERDIR=")
			output, err := command.CombinedOutput()
			switch mode {
			case "fail":
				var exit *exec.ExitError
				if !errors.As(err, &exit) || exit.ExitCode() != 3 {
					t.Fatalf("failing child: %v\n%s", err, output)
				}
			case "signal", "wait":
				var exit *exec.ExitError
				if !errors.As(err, &exit) || !exit.Sys().(syscall.WaitStatus).Signaled() {
					t.Fatalf("stopped child: %v\n%s", err, output)
				}
			default:
				if err != nil {
					t.Fatalf("child: %v\n%s", err, output)
				}
			}
			if strings.Contains(string(output), "GOCOVERDIR") {
				t.Fatalf("child coverage warning: %s", output)
			}
			assertNoSuiteTempEntries(t, scratch)
		})
	}
}

func assertNoSuiteTempEntries(t *testing.T, scratch string) {
	t.Helper()
	// Go's telemetry sidecar is detached and may recreate config after the
	// test process exits. Observe briefly after cleanup, not just at exit.
	deadline := time.Now().Add(250 * time.Millisecond)
	for {
		entries, err := os.ReadDir(scratch)
		if err != nil {
			t.Fatal(err)
		}
		if len(entries) != 0 {
			var names []string
			for _, entry := range entries {
				names = append(names, filepath.Join(scratch, entry.Name()))
			}
			t.Fatalf("temporary directory has %d leaked entries after child exit: %v", len(entries), names)
		}
		if time.Now().After(deadline) {
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
}
