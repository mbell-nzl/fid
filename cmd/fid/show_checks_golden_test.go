package main

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/mbell-nzl/fid/internal/checkrun"
	"github.com/mbell-nzl/fid/internal/gittest"
)

func showChecksGoldenFixture(t *testing.T, mode string) *statusFixture {
	t.Helper()
	f := requiredCheckFixture(t, "['always', 'go']")
	switch mode {
	case "unchecked", "passed":
		writeBackupTestFile(t, filepath.Join(f.workspace, "task.go"), "package task\n")
		runGitForTest(t, f.workspace, "add", ".")
		runGitForTest(t, f.workspace, "commit", "-m", "Task Go change")
		if mode == "passed" {
			requiredNames(t, checkRequired(t, f, "pass"))
		}
	case "docs-only":
		writeBackupTestFile(t, filepath.Join(f.workspace, "task.md"), "docs only\n")
		runGitForTest(t, f.workspace, "add", ".")
		runGitForTest(t, f.workspace, "commit", "-m", "Task docs change")
	case "base-unavailable":
		root := gittest.AddWorktree(t, f.workspace)
		identityGit(t, root, "checkout", "--orphan", "unrelated-show")
		runGitForTest(t, root, "add", ".")
		runGitForTest(t, root, "commit", "-m", "Unrelated policy root")
		identityGit(t, root, "update-ref", "refs/remotes/origin/trunk", "HEAD")
	case "policy-unavailable", "checks-not-configured", "definition-removed":
		got := f.run("check", "always", "--hold", "HLD-aaaaaa", "--actor", "agent:claude-builder", "--idempotency-key", "pass")
		if got.err != nil {
			t.Fatalf("check=%#v", got)
		}
		trustedShowPolicy(t, f, func(raw string) string {
			switch mode {
			case "policy-unavailable":
				return raw + "\nunknown=true\n"
			case "checks-not-configured":
				return strings.Replace(raw, "checks = ['always', 'go']", "", 1)
			default:
				raw = strings.Replace(raw, "checks = ['always', 'go']", "checks=[]", 1)
				return strings.Replace(raw, "[checks.always]\nrun=['true']\ntimeout='1m'", "", 1)
			}
		})
	default:
		f = resultCheckFixture(t)
		for index, fail := range map[string][]bool{"pass-then-fail": {false, true}, "fail-then-pass": {true, false}, "unfinished": {false}, "definition-changed": {false}}[mode] {
			if fail {
				writeBackupTestFile(t, filepath.Join(f.workspace, ".fail"), "fail")
			} else {
				if err := os.Remove(filepath.Join(f.workspace, ".fail")); err != nil && !os.IsNotExist(err) {
					t.Fatal(err)
				}
			}
			got := checkRequired(t, f, []string{"first", "second"}[index])
			if fail {
				assertJSONFailureCode(t, got.stderr, "check.execute", "check-unsuccessful")
			} else {
				requiredNames(t, got)
			}
		}
		switch mode {
		case "unfinished":
			options := leanOptions(f, "unfinished")
			options.Names = []string{"test"}
			options.Before = func(string) error { return errors.New("interrupted") }
			if _, err := checkrun.Run(context.Background(), options); err == nil {
				t.Fatal("expected interruption")
			}
		case "definition-changed":
			trustedShowPolicy(t, f, func(raw string) string { return strings.ReplaceAll(raw, "timeout='1m'", "timeout='2m'") })
		}
	}
	return f
}

func TestShowCheckGolden(t *testing.T) {
	t.Parallel()
	for _, mode := range []string{"unchecked", "docs-only", "passed", "base-unavailable", "policy-unavailable", "checks-not-configured", "definition-removed", "definition-changed", "pass-then-fail", "fail-then-pass", "unfinished"} {
		t.Run(mode, func(t *testing.T) {
			t.Parallel()
			f := showChecksGoldenFixture(t, mode)
			before := startEventCount(t, f.workspace)
			for _, jsonMode := range []bool{false, true} {
				args, ext := []string{"show", "FID-9f13aa"}, ".txt"
				if jsonMode {
					args, ext = append([]string{"--json"}, args...), ".json"
				}
				got := f.run(args...)
				if got.err != nil || got.stderr != "" || startEventCount(t, f.workspace) != before {
					t.Fatalf("read-only show=%#v", got)
				}
				output := regexp.MustCompile(`sha256:[0-9a-f]{64}`).ReplaceAllString(got.stdout, "sha256:<digest>")
				for _, field := range []string{"machine", "boot_id", "start_time", "executable", "os", "arch"} {
					output = regexp.MustCompile(`"`+field+`":"[^"]*"`).ReplaceAllString(output, `"`+field+`":"<`+field+`>"`)
				}
				checkGolden(t, f, "show-"+mode+ext, output)
			}
		})
	}
}
