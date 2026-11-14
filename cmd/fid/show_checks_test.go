package main

import (
	"context"
	"encoding/json"
	"errors"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mbell-nzl/fid/internal/checkrecord"
	"github.com/mbell-nzl/fid/internal/checkrun"
	"github.com/mbell-nzl/fid/internal/gittest"
)

type shownCheck struct {
	Name      string               `json:"name"`
	Start     *checkrecord.Start   `json:"start"`
	Process   *checkrecord.Process `json:"process"`
	Result    *checkrecord.Result  `json:"result"`
	Required  *bool                `json:"required"`
	Satisfied *bool                `json:"satisfied"`
}

type shownChecks struct {
	Checks   []shownCheck     `json:"checks"`
	Base     *string          `json:"checks_base"`
	Warnings []map[string]any `json:"warnings"`
}

func readShownChecks(t *testing.T, f *statusFixture) shownChecks {
	t.Helper()
	got := f.run("--json", "show", "FID-9f13aa")
	var envelope struct {
		Result shownChecks `json:"result"`
	}
	if got.err != nil || got.stderr != "" {
		t.Fatalf("show=%#v", got)
	}
	if err := json.Unmarshal([]byte(got.stdout), &envelope); err != nil {
		t.Fatal(err)
	}
	if envelope.Result.Checks == nil || envelope.Result.Warnings == nil {
		t.Fatalf("show omitted its arrays: %s", got.stdout)
	}
	return envelope.Result
}

func checkBool(t *testing.T, field string, got *bool, want bool) {
	t.Helper()
	if got == nil || *got != want {
		t.Fatalf("%s=%v, want %v", field, got, want)
	}
}

func trustedShowPolicy(t *testing.T, f *statusFixture, change func(string) string) string {
	t.Helper()
	root := gittest.AddWorktree(t, f.workspace)
	identityGit(t, root, "checkout", "-q", "--detach", "origin/trunk")
	path := filepath.Join(root, ".fid.toml")
	raw := string(mustReadCommandFile(t, path))
	writeBackupTestFile(t, path, change(raw))
	policyCommit(t, root)
	identityGit(t, root, "push", "-q", "origin", "HEAD:trunk")
	return root
}

func TestShowChecksRequiredBeforeFirstResult(t *testing.T) {
	t.Parallel()
	f := requiredCheckFixture(t, "['go', 'always']")
	writeBackupTestFile(t, filepath.Join(f.workspace, "new.go"), "package new\n")
	runGitForTest(t, f.workspace, "add", ".")
	runGitForTest(t, f.workspace, "commit", "-m", "Task Go file")
	before := startEventCount(t, f.workspace)
	got := readShownChecks(t, f)
	if len(got.Checks) != 2 || got.Checks[0].Name != "always" || got.Checks[1].Name != "go" || len(got.Warnings) != 0 || got.Base == nil || *got.Base != strings.TrimSpace(identityGit(t, f.workspace, "rev-parse", "origin/trunk")) {
		t.Fatalf("checks=%+v", got)
	}
	for _, check := range got.Checks {
		checkBool(t, "required", check.Required, true)
		checkBool(t, "satisfied", check.Satisfied, false)
		if check.Start != nil || check.Process != nil || check.Result != nil {
			t.Fatalf("invented record: %+v", check)
		}
	}
	text := f.run("show", "FID-9f13aa")
	if text.err != nil || !strings.Contains(text.stdout, "○ always required · not checked yet") || startEventCount(t, f.workspace) != before {
		t.Fatalf("read-only show=%#v", text)
	}
}

func TestShowChecksSelectTheCountingResultOnHeadsTree(t *testing.T) {
	t.Parallel()
	f := resultCheckFixture(t)
	first := requiredNames(t, checkRequired(t, f, "first"))
	if len(first) != 1 {
		t.Fatal(first)
	}
	original := readShownChecks(t, f).Checks[0]
	checkBool(t, "first satisfied", original.Satisfied, true)
	runGitForTest(t, f.workspace, "commit", "--allow-empty", "-m", "Same tree, new commit")
	checkBool(t, "same tree", readShownChecks(t, f).Checks[0].Satisfied, true)
	writeBackupTestFile(t, filepath.Join(f.workspace, "task.md"), "new tree\n")
	runGitForTest(t, f.workspace, "add", ".")
	runGitForTest(t, f.workspace, "commit", "-m", "Different tree")
	checkBool(t, "new tree", readShownChecks(t, f).Checks[0].Satisfied, false)
	requiredNames(t, checkRequired(t, f, "second-tree"))
	identityGit(t, f.workspace, "rm", "task.md")
	runGitForTest(t, f.workspace, "commit", "-m", "Return to original tree")
	returned := readShownChecks(t, f).Checks[0]
	checkBool(t, "returned tree", returned.Satisfied, true)
	if returned.Start == nil || returned.Start.ID != original.Start.ID {
		t.Fatalf("wrong counting record: %+v", returned)
	}
	options := leanOptions(f, "unfinished-show")
	options.Names = []string{"test"}
	options.Before = func(string) error { return errors.New("interrupted") }
	if _, err := checkrun.Run(context.Background(), options); err == nil {
		t.Fatal("expected interruption")
	}
	got := readShownChecks(t, f)
	if len(got.Checks) != 2 || got.Checks[1].Result != nil {
		t.Fatalf("unfinished start lost: %+v", got)
	}
	for _, check := range got.Checks {
		checkBool(t, "unfinished preserves satisfaction", check.Satisfied, true)
	}
	trustedShowPolicy(t, f, func(raw string) string { return raw + "\ndefined_by=['Makefile']\nwhen_changed=[]\n" })
	for _, check := range readShownChecks(t, f).Checks {
		checkBool(t, "scheduling metadata preserves satisfaction", check.Satisfied, true)
	}
	trustedShowPolicy(t, f, func(raw string) string { return strings.ReplaceAll(raw, "timeout='1m'", "timeout='2m'") })
	for _, check := range readShownChecks(t, f).Checks {
		checkBool(t, "definition change makes the result stale", check.Satisfied, false)
	}
}

func TestShowChecksRecomputeBaseAfterMergingOrigin(t *testing.T) {
	t.Parallel()
	f := requiredCheckFixture(t, "['go']")
	writeBackupTestFile(t, filepath.Join(f.workspace, "task.md"), "task docs\n")
	runGitForTest(t, f.workspace, "add", ".")
	runGitForTest(t, f.workspace, "commit", "-m", "Task docs")
	main := gittest.AddWorktree(t, f.workspace)
	identityGit(t, main, "checkout", "-q", "--detach", "origin/trunk")
	writeBackupTestFile(t, filepath.Join(main, "main.go"), "package main\n")
	runGitForTest(t, main, "add", ".")
	runGitForTest(t, main, "commit", "-m", "Main's Go change")
	identityGit(t, main, "push", "-q", "origin", "HEAD:trunk")
	before := readShownChecks(t, f)
	identityGit(t, f.workspace, "-c", "user.name=fid test", "-c", "user.email=fid@example.test", "merge", "--no-edit", "origin/trunk")
	after := readShownChecks(t, f)
	if before.Base == nil || after.Base == nil || *before.Base == *after.Base || *after.Base != strings.TrimSpace(identityGit(t, f.workspace, "rev-parse", "origin/trunk")) {
		t.Fatalf("base did not advance: before=%+v after=%+v", before, after)
	}
	if len(before.Checks) != 0 || len(after.Checks) != 0 || len(after.Warnings) != 0 {
		t.Fatalf("main's files required a check: before=%+v after=%+v", before, after)
	}
}

func TestShowChecksKeepUnknownSeparateFromFalse(t *testing.T) {
	t.Parallel()
	for _, mode := range []string{"policy-unavailable", "checks-not-configured", "comparison-unavailable", "undefined"} {
		t.Run(mode, func(t *testing.T) {
			t.Parallel()
			f := requiredCheckFixture(t, "['always', 'go']")
			got := f.run("check", "always", "--hold", "HLD-aaaaaa", "--actor", "agent:claude-builder", "--idempotency-key", "pass")
			if got.err != nil {
				t.Fatalf("check=%#v", got)
			}
			switch mode {
			case "policy-unavailable":
				trustedShowPolicy(t, f, func(raw string) string { return raw + "\nunknown=true\n" })
			case "checks-not-configured":
				trustedShowPolicy(t, f, func(raw string) string { return strings.Replace(raw, "checks = ['always', 'go']", "", 1) })
			case "undefined":
				trustedShowPolicy(t, f, func(raw string) string {
					raw = strings.Replace(raw, "checks = ['always', 'go']", "checks=[]", 1)
					return strings.Replace(raw, "[checks.always]\nrun=['true']\ntimeout='1m'", "", 1)
				})
			case "comparison-unavailable":
				root := gittest.AddWorktree(t, f.workspace)
				identityGit(t, root, "checkout", "--orphan", "unrelated-show")
				runGitForTest(t, root, "add", ".")
				runGitForTest(t, root, "commit", "-m", "Unrelated policy root")
				identityGit(t, root, "update-ref", "refs/remotes/origin/trunk", "HEAD")
			}
			shown := readShownChecks(t, f)
			if len(shown.Checks) == 0 {
				t.Fatalf("record lost: %+v", shown)
			}
			check := shown.Checks[0]
			if check.Name != "always" || check.Start == nil {
				t.Fatalf("record=%+v", check)
			}
			switch mode {
			case "policy-unavailable":
				if check.Required != nil || check.Satisfied != nil {
					t.Fatalf("unknown policy=%+v", check)
				}
			case "checks-not-configured":
				if check.Required != nil {
					t.Fatalf("unset checks=%+v", check)
				}
				checkBool(t, "known definition", check.Satisfied, true)
			case "undefined":
				if check.Satisfied != nil {
					t.Fatalf("undefined check=%+v", check)
				}
				checkBool(t, "no longer required", check.Required, false)
			case "comparison-unavailable":
				checkBool(t, "unconditional required", check.Required, true)
				checkBool(t, "known definition", check.Satisfied, true)
				if len(shown.Checks) != 2 || shown.Checks[1].Name != "go" || shown.Checks[1].Required != nil {
					t.Fatalf("conditional unknown=%+v", shown)
				}
				if len(shown.Warnings) != 1 || !strings.Contains(shown.Warnings[0]["problem"].(string), "common ancestor") {
					t.Fatalf("comparison warning must explain the missing base: %+v", shown.Warnings)
				}
			}
			if mode != "undefined" && (len(shown.Warnings) != 1 || shown.Warnings[0]["kind"] != mode) {
				t.Fatalf("warnings=%+v", shown.Warnings)
			}
			if shown.Base != nil {
				t.Fatalf("base should be absent: %+v", shown)
			}
		})
	}
}
