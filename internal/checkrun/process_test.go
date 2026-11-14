package checkrun

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/mbell-nzl/fid/internal/checkrecord"
)

//nolint:paralleltest // This entry point is also a disposable child process.
func TestCheckProcessHelper(t *testing.T) {
	for i, arg := range os.Args {
		if arg != "--check-helper" {
			continue
		}
		switch os.Args[i+1] {
		case "pass":
			fmt.Println("hello\x1b[31m")
			os.Exit(0)
		case "fail":
			os.Exit(3)
		case "sleep":
			time.Sleep(time.Hour)
		case "signal":
			_ = syscall.Kill(os.Getpid(), syscall.SIGTERM)
			time.Sleep(time.Hour)
		case "large":
			for j := 0; j < 10; j++ {
				fmt.Print(strings.Repeat("x", 1<<20))
			}
			fmt.Print(strings.Repeat("y", 1<<20))
			os.Exit(0)
		case "env":
			fmt.Printf("secret=%s; home=%s", os.Getenv("FID_CHECK_SECRET"), os.Getenv("HOME"))
			os.Exit(0)
		case "touch":
			if err := os.WriteFile("executed", []byte("executed"), 0o600); err != nil {
				os.Exit(2)
			}
			os.Exit(0)
		}
		os.Exit(0)
	}
}

func TestProcessCannotExecuteBeforeItsIdentityIsRecorded(t *testing.T) {
	t.Parallel()
	binary, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	root := t.TempDir()
	unrecorded := errors.New("process identity was not recorded")
	_, err = ExecuteProcess(context.Background(), ProcessOptions{Root: root, Definition: Definition{Run: []string{binary, "-test.run=^TestCheckProcessHelper$", "--", "--check-helper", "touch"}, Timeout: time.Minute, Environment: EnvironmentNames(), Artifacts: []string{}}, Started: func(_ checkrecord.Process) error { time.Sleep(100 * time.Millisecond); return unrecorded }})
	if !errors.Is(err, unrecorded) {
		t.Fatalf("recording failure=%v", err)
	}
	if _, err := os.Stat(filepath.Join(root, "executed")); !os.IsNotExist(err) {
		t.Fatalf("child executed without a durable identity: %v", err)
	}
}

func TestLaunchPipeClosingWithoutPermissionDoesNotExecuteTheCheck(t *testing.T) {
	t.Parallel()
	binary, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	root := t.TempDir()
	command, pipes, err := launchCommand(binary, []string{binary, "-test.run=^TestCheckProcessHelper$", "--", "--check-helper", "touch"})
	if err != nil {
		t.Fatal(err)
	}
	defer pipes.close()
	command.Dir = root
	if err := command.Start(); err != nil {
		t.Fatal(err)
	}
	pipes.closeChildEnds()
	if err := pipes.permitWriter.Close(); err != nil {
		t.Fatal(err)
	}
	if err := command.Wait(); err == nil {
		t.Fatal("unpermitted child succeeded")
	}
	if _, err := os.Stat(filepath.Join(root, "executed")); !os.IsNotExist(err) {
		t.Fatalf("unpermitted check executed: %v", err)
	}
}

// coverEnvironment gives the helper child a coverage directory when tests run
// with coverage. The child is this test binary; without GOCOVERDIR its exit
// adds a "GOCOVERDIR not set" warning to the output these tests compare.
func coverEnvironment(t *testing.T, env []string) []string {
	t.Helper()
	if testing.CoverMode() == "" {
		return env
	}
	return append(env, "GOCOVERDIR="+t.TempDir())
}

func TestProcessOutcomesAndBoundedOutput(t *testing.T) {
	t.Parallel()
	for _, row := range []struct {
		mode, outcome string
		timeout       time.Duration
		exit          int
	}{
		{"pass", "passed", time.Minute, 0}, {"fail", "failed", time.Minute, 3},
		{"sleep", "timed-out", 100 * time.Millisecond, -1}, {"signal", "signaled", time.Minute, -1},
		{"large", "passed", time.Minute, 0},
	} {
		t.Run(row.mode, func(t *testing.T) {
			t.Parallel()
			executable, err := os.Executable()
			if err != nil {
				t.Fatal(err)
			}
			definition := Definition{Run: []string{executable, "-test.run=^TestCheckProcessHelper$", "--", "--check-helper", row.mode}, Timeout: row.timeout, Environment: EnvironmentNames(), Artifacts: []string{}}
			got, err := ExecuteProcess(context.Background(), ProcessOptions{Root: t.TempDir(), Definition: definition, Environment: coverEnvironment(t, environment())})
			if err != nil {
				t.Fatal(err)
			}
			if got.Outcome != row.outcome {
				t.Fatalf("outcome=%s, want %s", got.Outcome, row.outcome)
			}
			if row.exit >= 0 && (got.ExitCode == nil || *got.ExitCode != row.exit) {
				t.Fatalf("exit=%v", got.ExitCode)
			}
			if got.Toolchain.Executable != executable || got.Toolchain.Digest == "" {
				t.Fatalf("toolchain=%#v", got.Toolchain)
			}
			if row.mode == "large" && (!got.OutputTruncated || got.OutputBytes != 11<<20 || len(got.OutputExcerpt) != 16<<10) {
				t.Fatalf("output=%d/%d/%v", got.OutputBytes, len(got.OutputExcerpt), got.OutputTruncated)
			}
			if row.mode == "pass" && (got.OutputBytes != int64(len("hello\x1b[31m\n")) || got.OutputExcerpt != "hello\x1b[31m\n" || got.OutputTruncated) {
				t.Fatalf("raw output=%#v", got.Result)
			}
			if row.mode == "pass" || row.mode == "large" {
				raw := "hello\x1b[31m\n"
				if row.mode == "large" {
					raw = strings.Repeat("x", 10<<20)
				}
				want := fmt.Sprintf("sha256:%x", sha256.Sum256([]byte(raw)))
				if got.OutputDigest != want {
					t.Fatalf("raw output digest=%s, want %s", got.OutputDigest, want)
				}
			}
		})
	}
}

//nolint:paralleltest // Sets the parent environment to prove it is filtered.
func TestProcessEnvironmentAndSpawnFailure(t *testing.T) {
	t.Setenv("FID_CHECK_SECRET", "must-not-reach-child")
	executable, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	definition := Definition{Run: []string{executable, "-test.run=^TestCheckProcessHelper$", "--", "--check-helper", "env"}, Timeout: time.Minute, Environment: EnvironmentNames(), Artifacts: []string{}}
	got, err := ExecuteProcess(context.Background(), ProcessOptions{Root: t.TempDir(), Definition: definition})
	if err != nil || got.Outcome != "passed" || strings.Contains(got.OutputExcerpt, "must-not-reach-child") || !strings.Contains(got.OutputExcerpt, "home="+os.Getenv("HOME")) {
		t.Fatalf("process=%#v, %v", got, err)
	}
	definition.Run = []string{"fid-check-executable-does-not-exist"}
	got, err = ExecuteProcess(context.Background(), ProcessOptions{Root: t.TempDir(), Definition: definition})
	if err != nil || got.Outcome != "spawn-failed" {
		t.Fatalf("spawn=%#v, %v", got, err)
	}
	bad := filepath.Join(t.TempDir(), "invalid-executable")
	if err := os.WriteFile(bad, []byte("not an executable image\n"), 0o700); err != nil {
		t.Fatal(err)
	}
	definition.Run = []string{bad}
	got, err = ExecuteProcess(context.Background(), ProcessOptions{Root: t.TempDir(), Definition: definition})
	if err != nil || got.Outcome != "spawn-failed" || got.ExitCode != nil {
		t.Fatalf("exec failure=%#v/%v", got, err)
	}
}

func TestProcessResolvesExecutableWithTheChildPath(t *testing.T) {
	t.Parallel()
	binary, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	directory := t.TempDir()
	if err := os.Symlink(binary, filepath.Join(directory, "check-helper")); err != nil {
		t.Fatal(err)
	}
	got, err := ExecuteProcess(context.Background(), ProcessOptions{Root: t.TempDir(), Definition: Definition{Run: []string{"check-helper", "-test.run=^TestCheckProcessHelper$", "--", "--check-helper", "env"}, Timeout: time.Minute, Environment: EnvironmentNames(), Artifacts: []string{}}, Environment: coverEnvironment(t, []string{"PATH=" + directory, "HOME=allowed"})})
	if err != nil || got.Outcome != "passed" || got.OutputExcerpt != "secret=; home=allowed" || got.Toolchain.Executable != binary {
		t.Fatalf("PATH=%#v/%v", got, err)
	}
}

type checkErrorWriter struct{ err error }

func (writer checkErrorWriter) Write([]byte) (int, error) { return 0, writer.err }

func TestProcessStopsOnAnOutputError(t *testing.T) {
	t.Parallel()
	binary, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	sinkError := errors.New("output sink failed")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	_, err = ExecuteProcess(ctx, ProcessOptions{Root: t.TempDir(), Definition: Definition{Run: []string{binary, "-test.run=^TestCheckProcessHelper$", "--", "--check-helper", "large"}, Timeout: time.Minute, Environment: EnvironmentNames(), Artifacts: []string{}}, Output: checkErrorWriter{sinkError}})
	if !errors.Is(err, sinkError) {
		t.Fatalf("output failure=%v", err)
	}
}
