package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/mbell-nzl/fid/internal/checkrecord"
	"github.com/mbell-nzl/fid/internal/cli"
)

func requiredCheckFixture(t *testing.T, names string) *statusFixture {
	t.Helper()
	f := plannedStartFixture(t, "in_progress", stringPointer("standard"))
	path := filepath.Join(f.workspace, ".fid.toml")
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	text := strings.Replace(string(raw), "[lanes.standard]\nchecks = []", "[lanes.standard]\nchecks = "+names, 1)
	text += "\n[checks.go]\nrun=['true']\ntimeout='1m'\nwhen_changed=['**/*.go']\n"
	text += "\n[checks.always]\nrun=['true']\ntimeout='1m'\n"
	text += "\n[checks.empty]\nrun=['true']\ntimeout='1m'\nwhen_changed=[]\n"
	writeBackupTestFile(t, path, text)
	writeBackupTestFile(t, filepath.Join(f.workspace, "old.go"), "package old\n")
	writeBackupTestFile(t, filepath.Join(f.workspace, "notes.md"), "notes\n")
	runGitForTest(t, f.workspace, "add", ".")
	runGitForTest(t, f.workspace, "commit", "-m", "Required check policy and base files")
	identityGit(t, f.workspace, "push", "-q", "origin", "trunk")
	identityGit(t, f.workspace, "checkout", "-q", "-b", "task")
	return f
}

func checkRequired(t *testing.T, f *statusFixture, key string) commandExecutionResult {
	t.Helper()
	return f.run("--json", "check", "--hold", "HLD-aaaaaa", "--actor", "agent:claude-builder", "--idempotency-key", key)
}

func requiredNames(t *testing.T, got commandExecutionResult) []string {
	t.Helper()
	if got.err != nil {
		t.Fatalf("check=%#v", got)
	}
	var envelope struct {
		Result struct {
			Checks []checkrecord.Summary `json:"checks"`
		} `json:"result"`
	}
	if err := json.Unmarshal([]byte(got.stdout), &envelope); err != nil {
		t.Fatal(err)
	}
	names := []string{}
	for _, check := range envelope.Result.Checks {
		if !check.Counts {
			t.Fatalf("check did not count: %s", got.stdout)
		}
		names = append(names, check.Name)
	}
	return names
}

func TestRequiredChecksUseCommittedChangedPaths(t *testing.T) {
	t.Parallel()
	for _, row := range []struct {
		name string
		edit func(*testing.T, *statusFixture)
		want []string
	}{
		{"go", func(t *testing.T, f *statusFixture) {
			writeBackupTestFile(t, filepath.Join(f.workspace, "new.go"), "package new\n")
		}, []string{"always", "empty", "go"}},
		{"docs", func(t *testing.T, f *statusFixture) {
			writeBackupTestFile(t, filepath.Join(f.workspace, "notes.md"), "edited\n")
		}, []string{"always", "empty"}},
		{"deleted", func(t *testing.T, f *statusFixture) { identityGit(t, f.workspace, "rm", "old.go") }, []string{"always", "empty", "go"}},
		{"renamed-from", func(t *testing.T, f *statusFixture) { identityGit(t, f.workspace, "mv", "old.go", "old.txt") }, []string{"always", "empty", "go"}},
		{"renamed-to", func(t *testing.T, f *statusFixture) { identityGit(t, f.workspace, "mv", "notes.md", "notes.go") }, []string{"always", "empty", "go"}},
	} {
		t.Run(row.name, func(t *testing.T) {
			t.Parallel()
			f := requiredCheckFixture(t, "['go', 'empty', 'always']")
			row.edit(t, f)
			runGitForTest(t, f.workspace, "add", ".")
			runGitForTest(t, f.workspace, "commit", "-m", "Task change")
			if got := requiredNames(t, checkRequired(t, f, "required")); !slices.Equal(got, row.want) {
				t.Fatalf("checks=%v, want %v", got, row.want)
			}
		})
	}
}

func TestRequiredChecksRecomputeBaseAfterMergingOrigin(t *testing.T) {
	t.Parallel()
	f := requiredCheckFixture(t, "['go']")
	writeBackupTestFile(t, filepath.Join(f.workspace, "notes.md"), "task docs\n")
	runGitForTest(t, f.workspace, "add", ".")
	runGitForTest(t, f.workspace, "commit", "-m", "Task docs")
	identityGit(t, f.workspace, "checkout", "-q", "trunk")
	writeBackupTestFile(t, filepath.Join(f.workspace, "main.go"), "package main\n")
	runGitForTest(t, f.workspace, "add", ".")
	runGitForTest(t, f.workspace, "commit", "-m", "Main's change")
	identityGit(t, f.workspace, "push", "-q", "origin", "trunk")
	identityGit(t, f.workspace, "checkout", "-q", "task")
	for _, key := range []string{"before-merge", "after-merge"} {
		if key == "after-merge" {
			identityGit(t, f.workspace, "-c", "user.name=fid test", "-c", "user.email=fid@example.test", "merge", "--no-edit", "origin/trunk")
		}
		if got := requiredNames(t, checkRequired(t, f, key)); len(got) != 0 {
			t.Fatalf("%s checks=%v; main's files must not count", key, got)
		}
	}
}

func TestRequiredChecksRefuseUnavailableComparisonButNamedChecksDoNotNeedIt(t *testing.T) {
	t.Parallel()
	f := requiredCheckFixture(t, "['go']")
	identityGit(t, f.workspace, "remote", "remove", "origin")
	got := checkRequired(t, f, "no-origin")
	if cli.ExitCode(got.err) != 4 {
		t.Fatalf("no origin=%#v", got)
	}
	assertJSONFailureCode(t, got.stderr, "check.execute", "comparison-unavailable")
	// A configured origin with unrelated history still supplies trusted policy,
	// but cannot supply a merge-base. Explicit names bypass only comparison.
	identityGit(t, f.workspace, "remote", "add", "origin", f.workspace)
	identityGit(t, f.workspace, "checkout", "--orphan", "unrelated")
	runGitForTest(t, f.workspace, "add", ".")
	runGitForTest(t, f.workspace, "commit", "-m", "Unrelated root")
	identityGit(t, f.workspace, "update-ref", "refs/remotes/origin/trunk", "HEAD")
	identityGit(t, f.workspace, "symbolic-ref", "refs/remotes/origin/HEAD", "refs/remotes/origin/trunk")
	identityGit(t, f.workspace, "checkout", "-q", "task")
	got = checkRequired(t, f, "unrelated")
	assertJSONFailureCode(t, got.stderr, "check.execute", "comparison-unavailable")
	named := f.run("--json", "check", "go", "--hold", "HLD-aaaaaa", "--actor", "agent:claude-builder", "--idempotency-key", "named")
	if got := requiredNames(t, named); !slices.Equal(got, []string{"go"}) {
		t.Fatalf("named=%v", got)
	}
}

func TestRequiredChecksRefuseInsufficientShallowHistory(t *testing.T) {
	t.Parallel()
	f := requiredCheckFixture(t, "['go']")
	writeBackupTestFile(t, filepath.Join(f.workspace, "new.go"), "package new\n")
	runGitForTest(t, f.workspace, "add", ".")
	runGitForTest(t, f.workspace, "commit", "-m", "Task beyond shallow boundary")
	head := strings.TrimSpace(identityGit(t, f.workspace, "rev-parse", "HEAD"))
	writeBackupTestFile(t, filepath.Join(f.workspace, ".git", "shallow"), head+"\n")
	got := checkRequired(t, f, "shallow")
	if cli.ExitCode(got.err) != 4 {
		t.Fatalf("shallow comparison=%#v", got)
	}
	assertJSONFailureCode(t, got.stderr, "check.execute", "comparison-unavailable")
}

func TestRequiredChecksDoNotNeedTestAdequacyPaths(t *testing.T) {
	t.Parallel()
	f := requiredCheckFixture(t, "['always']")
	path := filepath.Join(f.workspace, ".fid.toml")
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	writeBackupTestFile(t, path, strings.Replace(string(raw), "[review]\ntest_adequacy = []", "", 1))
	policyCommit(t, f.workspace)
	identityGit(t, f.workspace, "push", "-q", "origin", "HEAD:trunk")
	if got := requiredNames(t, checkRequired(t, f, "no-test-paths")); !slices.Equal(got, []string{"always"}) {
		t.Fatalf("checks=%v", got)
	}
}
