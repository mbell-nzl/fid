//go:build integration

package integration

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/mbell-nzl/fid/internal/config"
	"github.com/mbell-nzl/fid/internal/gittest"
)

// backupJourney is one machine for the built binary: HOME, user config and
// cache in temp directories, a local bare target, and an initialized project.
type backupJourney struct {
	base, workspace, target, configPath string
	env                                 []string
}

func newBackupJourney(t *testing.T) backupJourney {
	t.Helper()
	base := t.TempDir()
	journey := backupJourney{base: base, target: filepath.Join(base, "backups.git"), configPath: filepath.Join(base, "config", "fid", "config.toml")}
	if err := os.Mkdir(filepath.Join(base, "home"), 0o700); err != nil {
		t.Fatal(err)
	}
	journey.env = journey.environment(filepath.Join(base, "config"))
	if output, err := exec.Command("git", "init", "--bare", "-q", "--template=", journey.target).CombinedOutput(); err != nil {
		t.Fatalf("create backup target: %v: %s", err, output)
	}
	journey.workspace = gittest.NewWorkspace(t)
	backupJourneyFID(t, journey.env, journey.workspace, "init", "--prefix", "BKP", "--json").expectExit(t, 0)
	if err := config.SetBackupTarget(journey.configPath, journey.target); err != nil {
		t.Fatal(err)
	}
	return journey
}

// environment is the process environment with this machine's HOME, cache
// and configHome, less every inherited GIT_ variable but the guard.
func (journey backupJourney) environment(configHome string) []string {
	env := []string{}
	for _, entry := range os.Environ() {
		name, _, _ := strings.Cut(entry, "=")
		if name == "HOME" || name == "XDG_CONFIG_HOME" || name == "XDG_CACHE_HOME" ||
			strings.HasPrefix(name, "GIT_") && name != gittest.CeilingDirectoriesEnv {
			continue
		}
		env = append(env, entry)
	}
	return append(env, "HOME="+filepath.Join(journey.base, "home"), "XDG_CONFIG_HOME="+configHome, "XDG_CACHE_HOME="+filepath.Join(journey.base, "cache"))
}

// The loop and a daily cron job can run fid backup at the same moment, as
// separate processes. Both must succeed through fid's one clone, and both
// attempts must be recorded; a later failure exits non-zero and is recorded.
func TestBackupJourneyConcurrentProcessesAndRecordedFailure(t *testing.T) {
	t.Parallel()
	journey := newBackupJourney(t)
	env, workspace := journey.env, journey.workspace

	first, second := backupJourneyStart(t, env, workspace), backupJourneyStart(t, env, workspace)
	for _, run := range []*backupJourneyRun{first.wait(), second.wait()} {
		run.expectExit(t, 0)
	}
	outcomes := backupJourneyOutcomes(t, workspace)
	sort.Strings(outcomes)
	if strings.Join(outcomes, ",") != "pushed,unchanged" {
		t.Fatalf("concurrent outcomes = %v, want one push and one unchanged", outcomes)
	}

	if err := os.Remove(journey.configPath); err != nil {
		t.Fatal(err)
	}
	backupJourneyFID(t, env, workspace, "backup").expectExit(t, 2)
	if outcomes := backupJourneyOutcomes(t, workspace); len(outcomes) != 3 || outcomes[2] != "failed" {
		t.Fatalf("outcomes = %v, want the failure recorded", outcomes)
	}
}

// Processes that fail before fid's clone, such as with no target, record
// without its lock. When several meet a torn record at once, the record's
// own lock must keep every complete line and give each attempt its own.
func TestBackupJourneyConcurrentProcessesRepairATornRecordOnce(t *testing.T) {
	t.Parallel()
	const failing, seededPairs = 8, 5000
	journey := newBackupJourney(t)
	backupJourneyFID(t, journey.env, journey.workspace, "backup").expectExit(t, 0)
	backupJourneyFID(t, journey.env, journey.workspace, "backup").expectExit(t, 0)
	path := filepath.Join(config.LedgerDir(journey.workspace), "backups.jsonl")
	pair, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	// A long record widens the window in which unlocked repairs would overlap.
	seeded := bytes.Repeat(pair, seededPairs)
	if err := os.WriteFile(path, append(append([]byte{}, seeded...), pair[:48]...), 0o644); err != nil {
		t.Fatal(err)
	}
	withoutTarget := journey.environment(filepath.Join(journey.base, "config-without-target"))

	runs := []*backupJourneyRun{}
	for index := range failing + 1 {
		env := withoutTarget
		if index == failing {
			env = journey.env // this one records while it holds fid's clone
		}
		runs = append(runs, backupJourneyStart(t, env, journey.workspace))
	}
	exits := []int{}
	for _, run := range runs {
		exits = append(exits, run.wait().exit)
	}
	sort.Ints(exits)
	if fmt.Sprint(exits) != "[0 2 2 2 2 2 2 2 2]" {
		t.Fatalf("exit codes = %v, want one success and %d missing-target refusals", exits, failing)
	}
	after, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.HasPrefix(after, seeded) {
		t.Fatal("the complete records before the torn tail were not kept byte for byte")
	}
	outcomes := backupJourneyOutcomes(t, journey.workspace)
	if len(outcomes) != 2*seededPairs+failing+1 {
		t.Fatalf("records = %d, want the %d seeded and one per attempt", len(outcomes), 2*seededPairs)
	}
	added := outcomes[2*seededPairs:]
	sort.Strings(added)
	if strings.Join(added, ",") != "failed,failed,failed,failed,failed,failed,failed,failed,unchanged" {
		t.Fatalf("new outcomes = %v, want one per attempt", added)
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
