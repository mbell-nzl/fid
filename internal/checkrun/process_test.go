package checkrun

import (
	"context"
	"fmt"
	"os"
	"strings"
	"syscall"
	"testing"
	"time"
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
			for j := 0; j < 11; j++ {
				fmt.Print(strings.Repeat("x", 1<<20))
			}
			os.Exit(0)
		case "env":
			fmt.Printf("secret=%s; home=%s", os.Getenv("FID_CHECK_SECRET"), os.Getenv("HOME"))
			os.Exit(0)
		}
		os.Exit(0)
	}
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
			got, err := ExecuteProcess(context.Background(), ProcessOptions{Root: t.TempDir(), Definition: definition})
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
}
