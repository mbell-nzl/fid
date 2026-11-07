//go:build integration

package integration

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/mbell-nzl/fid/internal/config"
	"github.com/mbell-nzl/fid/internal/gittest"
)

// The loop and a daily cron job can run fid backup at the same moment, as
// separate processes. Both must succeed through fid's one clone, and both
// attempts must be recorded; a later failure exits non-zero and is recorded.
func TestBackupJourneyConcurrentProcessesAndRecordedFailure(t *testing.T) {
	t.Parallel()
	base := t.TempDir()
	home := filepath.Join(base, "home")
	if err := os.Mkdir(home, 0o700); err != nil {
		t.Fatal(err)
	}
	env := []string{}
	for _, entry := range os.Environ() {
		name, _, _ := strings.Cut(entry, "=")
		if name == "HOME" || name == "XDG_CONFIG_HOME" || name == "XDG_CACHE_HOME" ||
			strings.HasPrefix(name, "GIT_") && name != gittest.CeilingDirectoriesEnv {
			continue
		}
		env = append(env, entry)
	}
	env = append(env, "HOME="+home, "XDG_CONFIG_HOME="+filepath.Join(base, "config"), "XDG_CACHE_HOME="+filepath.Join(base, "cache"))
	target := filepath.Join(base, "backups.git")
	if output, err := exec.Command("git", "init", "--bare", "-q", "--template=", target).CombinedOutput(); err != nil {
		t.Fatalf("create backup target: %v: %s", err, output)
	}
	workspace := gittest.NewWorkspace(t)
	backupJourneyFID(t, env, workspace, "init", "--prefix", "BKP", "--json").expectExit(t, 0)
	configPath := filepath.Join(base, "config", "fid", "config.toml")
	if err := config.SetBackupTarget(configPath, target); err != nil {
		t.Fatal(err)
	}

	first, second := backupJourneyStart(t, env, workspace), backupJourneyStart(t, env, workspace)
	for _, run := range []*backupJourneyRun{first.wait(), second.wait()} {
		run.expectExit(t, 0)
	}
	outcomes := backupJourneyOutcomes(t, workspace)
	sort.Strings(outcomes)
	if strings.Join(outcomes, ",") != "pushed,unchanged" {
		t.Fatalf("concurrent outcomes = %v, want one push and one unchanged", outcomes)
	}

	if err := os.Remove(configPath); err != nil {
		t.Fatal(err)
	}
	backupJourneyFID(t, env, workspace, "backup").expectExit(t, 2)
	if outcomes := backupJourneyOutcomes(t, workspace); len(outcomes) != 3 || outcomes[2] != "failed" {
		t.Fatalf("outcomes = %v, want the failure recorded", outcomes)
	}
}

type backupJourneyRun struct {
	command        *exec.Cmd
	stdout, stderr bytes.Buffer
	exit           int
}

func backupJourneyStart(t *testing.T, env []string, workspace string) *backupJourneyRun {
	t.Helper()
	run := &backupJourneyRun{command: exec.Command(integrationFIDBinary(t), "--json", "backup")}
	run.command.Dir, run.command.Env = workspace, env
	run.command.Stdout, run.command.Stderr = &run.stdout, &run.stderr
	if err := run.command.Start(); err != nil {
		t.Fatal(err)
	}
	return run
}

func (run *backupJourneyRun) wait() *backupJourneyRun {
	var exit *exec.ExitError
	if err := run.command.Wait(); errors.As(err, &exit) {
		run.exit = exit.ExitCode()
	} else if err != nil {
		run.exit = -1
	}
	return run
}

func (run *backupJourneyRun) expectExit(t *testing.T, want int) {
	t.Helper()
	if run.exit != want {
		t.Fatalf("fid %v exit = %d, want %d; stdout=%q stderr=%q", run.command.Args[1:], run.exit, want, run.stdout.String(), run.stderr.String())
	}
}

func backupJourneyFID(t *testing.T, env []string, workspace string, args ...string) *backupJourneyRun {
	t.Helper()
	run := &backupJourneyRun{command: exec.Command(integrationFIDBinary(t), args...)}
	run.command.Dir, run.command.Env = workspace, env
	run.command.Stdout, run.command.Stderr = &run.stdout, &run.stderr
	if err := run.command.Start(); err != nil {
		t.Fatal(err)
	}
	return run.wait()
}

func backupJourneyOutcomes(t *testing.T, workspace string) []string {
	t.Helper()
	file, err := os.Open(filepath.Join(config.LedgerDir(workspace), "backups.jsonl"))
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = file.Close() }()
	outcomes := []string{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var record struct {
			Outcome string `json:"outcome"`
		}
		if err := json.Unmarshal(scanner.Bytes(), &record); err != nil {
			t.Fatalf("record line %q: %v", scanner.Text(), err)
		}
		outcomes = append(outcomes, record.Outcome)
	}
	if err := scanner.Err(); err != nil {
		t.Fatal(err)
	}
	return outcomes
}
