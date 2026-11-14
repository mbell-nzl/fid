package github

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"
)

//nolint:paralleltest // Disposable child owns PATH and other process environment.
func TestGHRunnerChild(t *testing.T) {
	if os.Getenv("FID_GH_RUNNER_CHILD") != "1" {
		return
	}
	if testing.CoverMode() != "" && os.Getenv("GOCOVERDIR") == "" {
		if err := os.Setenv("GOCOVERDIR", t.TempDir()); err != nil {
			t.Fatal(err)
		}
	}
	ctx := context.Background()
	if os.Getenv("FID_GH_TEST_MODE") == "timeout" {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 150*time.Millisecond)
		defer cancel()
	}
	result := RunGH(ctx, []string{"pr", "view", "74", "--repo", "mbell-nzl/fid"}, os.Environ())
	output := struct {
		Result                      Result
		Oversized, Timeout, Missing bool
	}{
		Result: result, Oversized: errors.Is(result.Err, ErrOutputLimit),
		Timeout: errors.Is(result.Err, context.DeadlineExceeded), Missing: errors.Is(result.Err, exec.ErrNotFound),
	}
	// Errors are inspected above, before JSON replaces the concrete error type.
	output.Result.Err = nil
	if err := json.NewEncoder(os.Stdout).Encode(output); err != nil {
		t.Fatal(err)
	}
	os.Exit(0)
}

func TestProductionRunner(t *testing.T) {
	t.Parallel()
	for _, mode := range []string{"environment", "timeout", "stdout", "stderr", "combined", "missing", "exit", "inherited-pipe"} {
		t.Run(mode, func(t *testing.T) {
			t.Parallel()
			bin := t.TempDir()
			script := `#!/bin/sh
case "$FID_GH_TEST_MODE" in
 environment) printf '%s\n' "$@"; /usr/bin/env ;;
 timeout) exec /bin/sleep 60 ;;
 stdout) exec /usr/bin/head -c 9000000 /dev/zero ;;
 stderr) exec /usr/bin/head -c 9000000 /dev/zero >&2 ;;
 combined) /usr/bin/head -c 5000000 /dev/zero; /usr/bin/head -c 5000000 /dev/zero >&2 ;;
 exit) echo 'API unavailable' >&2; exit 9 ;;
 inherited-pipe) /bin/sleep 2 & exit 0 ;;
 *) exit 97 ;;
esac
`
			if mode != "missing" {
				if err := os.WriteFile(filepath.Join(bin, "gh"), []byte(script), 0o700); err != nil {
					t.Fatal(err)
				}
			}
			command := exec.Command(os.Args[0], "-test.run=^TestGHRunnerChild$")
			command.Env = append(os.Environ(), "FID_GH_RUNNER_CHILD=1", "FID_GH_TEST_MODE="+mode, "PATH="+bin,
				"GH_PROMPT_DISABLED=0", "GH_PAGER=bad", "PAGER=bad", "NO_COLOR=0", "GH_HOST=wrong.example", "GH_FORCE_TTY=1", "GH_DEBUG=api", "DEBUG=1", "CLICOLOR_FORCE=1")
			began := time.Now()
			data, err := command.CombinedOutput()
			if err != nil {
				t.Fatalf("runner child: %v: %s", err, data)
			}
			var got struct {
				Result                      Result
				Oversized, Timeout, Missing bool
			}
			if err := json.Unmarshal(data, &got); err != nil {
				t.Fatalf("decode child: %v: %.300s", err, data)
			}
			switch mode {
			case "environment":
				lines := strings.Split(string(got.Result.Stdout), "\n")
				if !slices.Equal(lines[:5], []string{"pr", "view", "74", "--repo", "mbell-nzl/fid"}) {
					t.Fatalf("argv = %v", lines[:5])
				}
				for _, value := range []string{"GH_PROMPT_DISABLED=1", "GH_PAGER=cat", "PAGER=cat", "NO_COLOR=1", "GH_HOST=github.com", "GH_NO_UPDATE_NOTIFIER=1", "GH_NO_EXTENSION_UPDATE_NOTIFIER=1", "GH_TELEMETRY=0", "CLICOLOR=0"} {
					if !slices.Contains(lines, value) {
						t.Errorf("missing environment %s", value)
					}
				}
				for _, key := range []string{"GH_FORCE_TTY=", "GH_DEBUG=", "DEBUG=", "CLICOLOR_FORCE="} {
					for _, line := range lines {
						if strings.HasPrefix(line, key) {
							t.Errorf("unsafe environment %s", line)
						}
					}
				}
			case "timeout":
				if !got.Timeout || time.Since(began) > 3*time.Second {
					t.Fatalf("timeout = %#v after %s", got, time.Since(began))
				}
			case "stdout", "stderr", "combined":
				if !got.Oversized || len(got.Result.Stdout)+len(got.Result.Stderr) > DefaultOutputLimit {
					t.Fatalf("output bound: oversized=%v bytes=%d", got.Oversized, len(got.Result.Stdout)+len(got.Result.Stderr))
				}
			case "missing":
				if !got.Missing {
					t.Fatalf("missing = %#v", got)
				}
			case "exit":
				if got.Result.ExitCode != 9 || string(got.Result.Stderr) != "API unavailable\n" {
					t.Fatalf("exit = %#v", got)
				}
			case "inherited-pipe":
				if time.Since(began) > time.Second {
					t.Fatalf("inherited pipe blocked exit for %s", time.Since(began))
				}
			}
		})
	}
}
