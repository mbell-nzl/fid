package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/mbell-nzl/fid/internal/checkrecord"
	"github.com/mbell-nzl/fid/internal/checkrun"
	"github.com/mbell-nzl/fid/internal/cli"
)

//nolint:paralleltest // A disposable child process for the named check tests.
func TestLeanCheckChild(t *testing.T) {
	for i, arg := range os.Args {
		if arg != "--check-child" {
			continue
		}
		switch os.Args[i+1] {
		case "pass":
			fmt.Println("check output\x1b[31m")
		case "fail":
			fmt.Fprintln(os.Stderr, "check failed")
			os.Exit(3)
		case "signal":
			_ = syscall.Kill(os.Getpid(), syscall.SIGTERM)
			time.Sleep(time.Hour)
		case "pause":
			time.Sleep(150 * time.Millisecond)
		case "wait":
			time.Sleep(time.Hour)
		case "change":
			if err := os.WriteFile("coverage.out", []byte("changed"), 0o600); err != nil {
				os.Exit(2)
			}
		}
		os.Exit(0)
	}
}

func TestLeanCheckStartNamesThePolicyItActuallyResolved(t *testing.T) {
	t.Parallel()
	f := leanJourneyFixture(t, leanHelperDefinition(t, "pass", "1m", nil)+"\n[checks.zpass]\nrun=['true']\ntimeout='1m'\n")
	options := leanOptions(f, "policy-binding")
	options.Names = []string{"lean", "zpass"}
	options.After = func(record checkrecord.Record) error {
		if record.Start.Name != "lean" {
			return nil
		}
		appendPolicy(t, f.workspace, "defined_by=['Makefile']\n")
		policyCommit(t, f.workspace)
		identityGit(t, f.workspace, "push", "-q", "origin", "trunk")
		got := f.run("--json", "next", "--claim", "--actor", "agent:claude-builder")
		return got.err
	}
	if _, err := checkrun.Run(context.Background(), options); err != nil {
		t.Fatal(err)
	}
	shown := f.run("--json", "show", "FID-9f13aa")
	var response struct {
		Result struct {
			Checks []checkrecord.Record `json:"checks"`
		} `json:"result"`
	}
	if err := json.Unmarshal([]byte(shown.stdout), &response); err != nil {
		t.Fatal(err)
	}
	if len(response.Result.Checks) != 2 || response.Result.Checks[0].Start.PolicySnapshot != response.Result.Checks[1].Start.PolicySnapshot {
		t.Fatalf("check starts mixed policy snapshots: %s", shown.stdout)
	}
}

func leanHelperDefinition(t *testing.T, mode, timeout string, artifacts []string) string {
	t.Helper()
	binary, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	paths := []string{}
	for _, path := range artifacts {
		paths = append(paths, fmt.Sprintf("%q", path))
	}
	return fmt.Sprintf("\n[checks.lean]\nrun = [%q, '-test.run=^TestLeanCheckChild$', '--', '--check-child', %q]\ntimeout = %q\nartifacts = [%s]\n", binary, mode, timeout, strings.Join(paths, ","))
}

func leanJourneyFixture(t *testing.T, definition string) *statusFixture {
	t.Helper()
	f := plannedStartFixture(t, "in_progress", stringPointer("standard"))
	appendPolicy(t, f.workspace, definition)
	policyCommit(t, f.workspace)
	identityGit(t, f.workspace, "push", "-q", "origin", "trunk")
	if err := os.WriteFile(filepath.Join(f.workspace, "coverage.out"), []byte("kept"), 0o600); err != nil {
		t.Fatal(err)
	}
	runGitForTest(t, f.workspace, "add", ".")
	runGitForTest(t, f.workspace, "commit", "-m", "Check fixture")
	return f
}

func latestLeanRecord(t *testing.T, f *statusFixture) checkrecord.Record {
	t.Helper()
	shown := f.run("--json", "show", "FID-9f13aa")
	if shown.err != nil {
		t.Fatalf("show=%#v", shown)
	}
	var envelope struct {
		Result struct {
			Checks []checkrecord.Record `json:"checks"`
		} `json:"result"`
	}
	if err := json.Unmarshal([]byte(shown.stdout), &envelope); err != nil {
		t.Fatal(err)
	}
	if len(envelope.Result.Checks) == 0 {
		t.Fatal("no check record")
	}
	return envelope.Result.Checks[len(envelope.Result.Checks)-1]
}

func TestLeanCheckOutcomesQualificationAndArtifacts(t *testing.T) {
	t.Parallel()
	for _, row := range []struct {
		name, mode, timeout, outcome, reason string
		exit                                 int
		setup                                func(*statusFixture)
	}{
		{"pass", "pass", "15m", "passed", "", 0, nil},
		{"fail", "fail", "15m", "failed", "", 8, nil},
		{"timeout", "wait", "100ms", "timed-out", "", 8, nil},
		{"signal", "signal", "15m", "signaled", "", 8, nil},
		{"dirty", "pass", "15m", "passed", "tree-dirty", 4, func(f *statusFixture) { writeBackupTestFile(f.t, filepath.Join(f.workspace, "uncommitted"), "dirty") }},
		{"changed", "change", "15m", "passed", "tree-changed", 4, nil},
		{"missing", "pass", "15m", "passed", "artifact-not-kept", 4, func(f *statusFixture) {
			if err := os.Remove(filepath.Join(f.workspace, "coverage.out")); err != nil {
				f.t.Fatal(err)
			}
		}},
	} {
		t.Run(row.name, func(t *testing.T) {
			t.Parallel()
			f := leanJourneyFixture(t, leanHelperDefinition(t, row.mode, row.timeout, []string{"coverage.out"}))
			if row.setup != nil {
				row.setup(f)
			}
			got := f.run("--json", "check", "lean", "--hold", "HLD-aaaaaa", "--actor", "agent:claude-builder", "--idempotency-key", "check")
			exit := 0
			if got.err != nil {
				exit = cli.ExitCode(got.err)
			}
			if exit != row.exit {
				t.Fatalf("check=%#v, want exit %d", got, row.exit)
			}
			record := latestLeanRecord(t, f)
			if record.Result == nil || record.Result.Outcome != row.outcome || record.Result.Counts != (row.exit == 0) || row.reason != "" && !strings.Contains(strings.Join(record.Result.Reasons, ","), row.reason) {
				t.Fatalf("record=%#v", record)
			}
			if record.Start.DefinitionDigest == "" || record.Start.EnvironmentDigest == "" || record.Process == nil {
				t.Fatalf("incomplete start/process=%#v", record)
			}
			if row.mode == "pass" && !strings.Contains(record.Result.OutputExcerpt, "check output\x1b[31m") {
				t.Fatalf("output=%q", record.Result.OutputExcerpt)
			}
			if row.exit != 0 && strings.Contains(got.stderr, `"result"`) {
				t.Fatalf("failure carries a result: %s", got.stderr)
			}
			repeat := f.run("--json", "check", "lean", "--hold", "HLD-aaaaaa", "--actor", "agent:claude-builder", "--idempotency-key", "check")
			if repeat.stdout != got.stdout || repeat.stderr != got.stderr {
				t.Fatalf("retry=%#v, original=%#v", repeat, got)
			}
		})
	}
}

func TestLeanCheckContinuesAfterFailureAndDeduplicatesNames(t *testing.T) {
	t.Parallel()
	definition := leanHelperDefinition(t, "fail", "1m", nil) + "\n[checks.zpass]\nrun=['true']\ntimeout='1m'\n"
	f := leanJourneyFixture(t, definition)
	got := f.run("check", "zpass", "lean", "lean", "--hold", "HLD-aaaaaa", "--actor", "agent:claude-builder", "--idempotency-key", "several")
	if cli.ExitCode(got.err) != 8 || !strings.Contains(got.stdout, "✗ lean failed\n● Checking zpass…\n✓ zpass passed · counts for hand-off") || strings.Count(got.stdout, "Checking lean") != 1 {
		t.Fatalf("several=%#v", got)
	}
	repeat := f.run("check", "zpass", "lean", "lean", "--hold", "HLD-aaaaaa", "--actor", "agent:claude-builder", "--idempotency-key", "several")
	if repeat.stdout != got.stdout || repeat.stderr != got.stderr {
		t.Fatalf("human retry=%#v, original=%#v", repeat, got)
	}
}

func TestLeanCheckRefusesUnsafeAndOversizedArtifacts(t *testing.T) {
	t.Parallel()
	for _, row := range []struct{ path, reason string }{{"../outside", "outside-worktree"}, {"link", "outside-worktree"}, {"large", "oversized"}, {"missing", "missing"}} {
		t.Run(row.reason+row.path, func(t *testing.T) {
			t.Parallel()
			f := leanJourneyFixture(t, leanHelperDefinition(t, "pass", "1m", []string{row.path}))
			writeBackupTestFile(t, filepath.Join(f.workspace, ".gitignore"), "link\nlarge\n")
			runGitForTest(t, f.workspace, "add", ".gitignore")
			runGitForTest(t, f.workspace, "commit", "-m", "Ignore generated files")
			if row.path == "link" {
				outside := filepath.Join(t.TempDir(), "outside")
				writeBackupTestFile(t, outside, "outside bytes")
				if err := os.Symlink(outside, filepath.Join(f.workspace, "link")); err != nil {
					t.Fatal(err)
				}
			}
			if row.path == "large" {
				file, err := os.Create(filepath.Join(f.workspace, "large"))
				if err != nil {
					t.Fatal(err)
				}
				if err := file.Truncate((64 << 20) + 1); err != nil {
					t.Fatal(err)
				}
				if err := file.Close(); err != nil {
					t.Fatal(err)
				}
			}
			got := f.run("--json", "check", "lean", "--hold", "HLD-aaaaaa", "--actor", "agent:claude-builder", "--idempotency-key", "unsafe")
			assertJSONFailureCode(t, got.stderr, "check.execute", "check-not-qualified")
			record := latestLeanRecord(t, f)
			if record.Result == nil || record.Result.Counts || len(record.Result.RefusedArtifacts) != 1 || record.Result.RefusedArtifacts[0].Reason != row.reason || strings.Join(record.Result.Reasons, ",") != "artifact-not-kept" {
				t.Fatalf("unsafe file=%#v", record)
			}
		})
	}
}

func TestLeanCheckDefinitionDigestExcludesSchedulingMetadata(t *testing.T) {
	t.Parallel()
	f := leanJourneyFixture(t, leanHelperDefinition(t, "pass", "1m", nil))
	execute := func(key string) checkrecord.Record {
		t.Helper()
		got := f.run("--json", "check", "lean", "--hold", "HLD-aaaaaa", "--actor", "agent:claude-builder", "--idempotency-key", key)
		if got.err != nil {
			t.Fatalf("check=%#v", got)
		}
		return latestLeanRecord(t, f)
	}
	first := execute("first")
	appendPolicy(t, f.workspace, "defined_by=['Makefile']\nwhen_changed=['**/*.go']\n")
	policyCommit(t, f.workspace)
	identityGit(t, f.workspace, "push", "-q", "origin", "trunk")
	second := execute("second")
	if second.Start.DefinitionDigest != first.Start.DefinitionDigest || second.Start.PolicySnapshot == first.Start.PolicySnapshot {
		t.Fatalf("unrelated policy changed definition: %#v/%#v", first.Start, second.Start)
	}
	path := filepath.Join(f.workspace, ".fid.toml")
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(strings.Replace(string(raw), "artifacts = []", "artifacts = ['coverage.out']", 1)), 0o600); err != nil {
		t.Fatal(err)
	}
	policyCommit(t, f.workspace)
	identityGit(t, f.workspace, "push", "-q", "origin", "trunk")
	third := execute("third")
	if third.Start.DefinitionDigest == second.Start.DefinitionDigest || len(third.Result.Artifacts) != 1 {
		t.Fatalf("Artifact list did not bind the definition: %#v", third)
	}
}
