package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/mbell-nzl/fid/internal/cli"
	"github.com/mbell-nzl/fid/internal/gittest"
)

func handoffArgs(key string) []string {
	return []string{"handoff", "--hold", "HLD-aaaaaa", "--actor", "agent:claude-builder", "--pr", "https://github.com/mbell-nzl/fid/pull/63", "--idempotency-key", key}
}

func handoffGolden(t *testing.T, f *statusFixture, name, output string) {
	t.Helper()
	output = regexp.MustCompile(` on [0-9a-f]{8}\b`).ReplaceAllString(output, " on <head>")
	output = regexp.MustCompile(`sha256:[0-9a-f]{64}`).ReplaceAllString(output, "sha256:<digest>")
	for _, field := range []string{"machine", "boot_id", "start_time", "os", "arch", "executable", "fid_version"} {
		output = regexp.MustCompile(`"`+field+`":"[^"]*"`).ReplaceAllString(output, `"`+field+`":"<`+field+`>"`)
	}
	output = regexp.MustCompile(`"pid":[0-9]+`).ReplaceAllString(output, `"pid":10001`)
	output = regexp.MustCompile(`\b[0-9a-f]{40}\b`).ReplaceAllString(output, "<commit-or-tree>")
	output = regexp.MustCompile(`\b(?:HOF|CHK)-[0-9a-f]{6}\b`).ReplaceAllStringFunc(output, func(id string) string { return id[:3] + "-<new>" })
	lifecycleGolden(t, f, filepath.Join("..", "handoff", name), output)
}

func TestHandoffGoldenAndRetry(t *testing.T) {
	t.Parallel()
	for _, names := range []string{"[]", "['always', 'empty']"} {
		t.Run(names, func(t *testing.T) {
			t.Parallel()
			f := requiredCheckFixture(t, names)
			requiredNames(t, checkRequired(t, f, "checks"))
			// A commit with the same tree preserves the recorded checks.
			runGitForTest(t, f.workspace, "commit", "--allow-empty", "-m", "Ready for review")
			args := handoffArgs("handoff")
			first := f.run(append([]string{"--json"}, args...)...)
			fields, _ := lifecycleResult(t, first, "task.handoff")
			if fields["status"] != "in_review" || fields["handoff_id"] == "" {
				t.Fatalf("handoff=%v", fields)
			}
			name := "passed"
			if names == "[]" {
				name = "empty"
			}
			handoffGolden(t, f, name+".json", first.stdout)
			before := startEventCount(t, f.workspace)
			// Replay precedes hold, policy, git and clock checks.
			f.clock = func() time.Time { return statusNow.Add(2 * time.Hour) }
			again := f.run(append([]string{"--json"}, args...)...)
			if again.err != nil || again.stdout != first.stdout || startEventCount(t, f.workspace) != before {
				t.Fatalf("retry=%#v; original=%s", again, first.stdout)
			}
			f.clock = func() time.Time { return statusNow }
			human := f.run(args...)
			if human.err != nil {
				t.Fatalf("text=%#v", human)
			}
			handoffGolden(t, f, name+".txt", human.stdout)
			for _, command := range []string{"next", "show", "status"} {
				for _, jsonMode := range []bool{false, true} {
					invoke, ext := []string{command}, ".txt"
					if command == "show" {
						invoke = append(invoke, "FID-9f13aa")
					}
					if jsonMode {
						invoke, ext = append([]string{"--json"}, invoke...), ".json"
					}
					got := f.run(invoke...)
					if got.err != nil {
						t.Fatalf("%s=%#v", command, got)
					}
					handoffGolden(t, f, name+"-"+command+ext, got.stdout)
				}
			}
			other := gittest.AddWorktree(t, f.workspace)
			f.seed([]readyCommandRow{statusTask("FID-abcdef", "Other work", "todo", nil, nil, nil)}, nil)
			got, _ := decodeNext(t, executeCommandForTest(other, "--json", "next"), "next.get")
			if got["kind"] != "new" {
				t.Fatalf("other worktree=%v", got)
			}
			lifecycleResult(t, f.run("--json", "issue", "transition", "FID-9f13aa", "--to", "cancelled", "--actor", "human:mark", "--reason", "No longer needed"), "issue.transition")
			got, _ = decodeNext(t, f.run("--json", "next"), "next.get")
			if got["kind"] != "new" {
				t.Fatalf("cancellation did not free worktree: %v", got)
			}
		})
	}
}

func TestHandoffRefusalGolden(t *testing.T) {
	t.Parallel()
	for _, row := range []struct{ name, code string }{
		{"checks-not-configured", "checks-not-configured"}, {"test-paths-not-configured", "test-paths-not-configured"},
		{"tree-dirty", "tree-dirty"}, {"checks-unsatisfied", "checks-unsatisfied"}, {"comparison-unavailable", "comparison-unavailable"},
		{"policy-unavailable", "policy-unavailable"}, {"not-started", "transition-not-permitted"},
		{"missing-hold", "hold-not-current"}, {"ended-hold", "hold-not-current"}, {"expired-hold", "hold-not-current"},
		{"wrong-actor", "hold-not-current"}, {"wrong-worktree", "hold-not-current"},
		{"invalid-pr", "invalid-argument"}, {"missing-pr", "invalid-argument"}, {"missing-key", "invalid-argument"},
	} {
		t.Run(row.name, func(t *testing.T) {
			t.Parallel()
			f := requiredCheckFixture(t, "['always']")
			args := handoffArgs("refused")
			switch row.name {
			case "checks-not-configured", "test-paths-not-configured", "policy-unavailable":
				path := filepath.Join(f.workspace, ".fid.toml")
				raw, err := os.ReadFile(path)
				if err != nil {
					t.Fatal(err)
				}
				text := string(raw)
				if row.name == "checks-not-configured" {
					text = strings.Replace(text, "checks = ['always']", "", 1)
				}
				if row.name == "test-paths-not-configured" {
					text = strings.Replace(text, "[review]\ntest_adequacy = []", "", 1)
				}
				if row.name == "policy-unavailable" {
					text += "\nunknown=1\n"
				}
				writeBackupTestFile(t, path, text)
				policyCommit(t, f.workspace)
				identityGit(t, f.workspace, "push", "-q", "origin", "HEAD:trunk")
			case "tree-dirty":
				writeBackupTestFile(t, filepath.Join(f.workspace, "dirty"), "uncommitted")
			case "comparison-unavailable":
				runGitForTest(t, f.workspace, "commit", "--allow-empty", "-m", "Beyond shallow boundary")
				head := strings.TrimSpace(identityGit(t, f.workspace, "rev-parse", "HEAD"))
				writeBackupTestFile(t, filepath.Join(f.workspace, ".git", "shallow"), head+"\n")
			case "not-started":
				startSQL(t, f, `UPDATE issues SET status='todo' WHERE id='FID-9f13aa'`)
			case "missing-hold":
				args[2] = "HLD-bbbbbb"
			case "ended-hold":
				startSQL(t, f, `UPDATE holds SET ended_at='2026-11-08T09:00:00Z',end_reason='takeover' WHERE id='HLD-aaaaaa'`)
			case "expired-hold":
				f.clock = func() time.Time { return statusNow.Add(time.Hour) }
			case "wrong-actor":
				args[4] = "agent:other"
			case "wrong-worktree":
				startSQL(t, f, `UPDATE holds SET worktree_id='sha256:elsewhere',worktree_root='/w/elsewhere' WHERE id='HLD-aaaaaa'`)
			case "invalid-pr":
				args[6] = "https://example.com/owner/repo/pull/63"
			case "missing-pr":
				args = append(args[:5], args[7:]...)
			case "missing-key":
				args = args[:7]
			}
			before := startEventCount(t, f.workspace)
			for _, jsonMode := range []bool{false, true} {
				invoke, ext := args, ".txt"
				if jsonMode {
					invoke, ext = append([]string{"--json"}, args...), ".json"
				}
				got := f.run(invoke...)
				want := 4
				if row.code == "invalid-argument" {
					want = 2
				}
				if cli.ExitCode(got.err) != want || got.stdout != "" {
					t.Fatalf("refusal=%#v", got)
				}
				if jsonMode {
					assertJSONFailureCode(t, got.stderr, "task.handoff", row.code)
				}
				handoffGolden(t, f, row.name+ext, got.stderr)
			}
			if startEventCount(t, f.workspace) != before {
				t.Fatal("refusal wrote an event")
			}
		})
	}
}

func TestHandoffReviewElsewhereAndMissingWorktree(t *testing.T) {
	t.Parallel()
	f := requiredCheckFixture(t, "[]")
	lifecycleResult(t, f.run(append([]string{"--json"}, handoffArgs("handoff")...)...), "task.handoff")
	other := gittest.AddWorktree(t, f.workspace)
	for _, args := range [][]string{{"next"}, {"next", "--explain"}} {
		got := executeCommandForTest(other, args...)
		if got.err != nil || !strings.Contains(got.stdout, "in review") {
			t.Fatalf("review elsewhere=%#v", got)
		}
	}
	startSQL(t, f, `UPDATE holds SET worktree_root='/w/no-longer-present' WHERE id='HLD-aaaaaa'`)
	got := f.run("--json", "status")
	if got.err != nil || !strings.Contains(got.stdout, `"kind":"worktree-missing"`) {
		t.Fatalf("missing review worktree=%#v", got)
	}
}

func TestHandoffRefusesUnmergedChangesAsTreeDirty(t *testing.T) {
	t.Parallel()
	f := requiredCheckFixture(t, "[]")
	oid := strings.TrimSpace(identityGit(t, f.workspace, "rev-parse", "HEAD:old.go"))
	runGitForTest(t, f.workspace, "update-index", "--force-remove", "old.go")
	command := exec.Command("git", "-C", f.workspace, "update-index", "--index-info")
	command.Stdin = strings.NewReader("100644 " + oid + " 1\told.go\n100644 " + oid + " 2\told.go\n100644 " + oid + " 3\told.go\n")
	if output, err := command.CombinedOutput(); err != nil {
		t.Fatalf("make unmerged index: %v: %s", err, output)
	}
	before := startEventCount(t, f.workspace)
	got := f.run(append([]string{"--json"}, handoffArgs("unmerged")...)...)
	if cli.ExitCode(got.err) != 4 {
		t.Fatalf("unmerged hand-off=%#v", got)
	}
	assertJSONFailureCode(t, got.stderr, "task.handoff", "tree-dirty")
	if startEventCount(t, f.workspace) != before {
		t.Fatal("refusal wrote an event")
	}
}
