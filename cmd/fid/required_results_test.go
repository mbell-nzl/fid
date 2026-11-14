package main

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/mbell-nzl/fid/internal/checkread"
	"github.com/mbell-nzl/fid/internal/checkrun"
	"github.com/mbell-nzl/fid/internal/cli"
	"github.com/mbell-nzl/fid/internal/config"
	fiddb "github.com/mbell-nzl/fid/internal/db"
	"github.com/mbell-nzl/fid/internal/gittest"
	"github.com/mbell-nzl/fid/internal/policy"
)

func resultCheckFixture(t *testing.T) *statusFixture {
	t.Helper()
	f := requiredCheckFixture(t, "['test']")
	appendPolicy(t, f.workspace, "\n[checks.test]\nrun=['sh','-c','test ! -f .fail']\ntimeout='1m'\n")
	writeBackupTestFile(t, filepath.Join(f.workspace, ".gitignore"), ".fail\n")
	runGitForTest(t, f.workspace, "add", ".")
	runGitForTest(t, f.workspace, "commit", "-m", "Check with controllable outcome")
	identityGit(t, f.workspace, "push", "-q", "origin", "HEAD:trunk")
	return f
}

func evaluateRequired(t *testing.T, f *statusFixture, task string) checkread.Evaluation {
	t.Helper()
	loaded, err := policy.Load(f.workspace)
	if err != nil {
		t.Fatal(err)
	}
	comparison, err := policy.Compare(f.workspace, task)
	if err != nil {
		t.Fatal(err)
	}
	names, err := policy.Required(loaded.Effective, "standard", task, comparison.Paths)
	if err != nil {
		t.Fatal(err)
	}
	conn, err := fiddb.Open(filepath.Join(config.LedgerDir(f.workspace), config.DatabaseFileName))
	if err != nil {
		t.Fatal(err)
	}
	defer closeCommandDB(t, conn)
	records, err := checkread.Read(conn, task)
	if err != nil {
		t.Fatal(err)
	}
	tree := identityGit(t, f.workspace, "rev-parse", "HEAD^{tree}")
	return checkread.Evaluate(task, strings.TrimSpace(tree), names, loaded.Effective.Checks, records)
}

func TestRequiredResultUsesLatestFinishedOutcome(t *testing.T) {
	t.Parallel()
	for _, outcomes := range [][]bool{{true, false}, {false, true}} {
		t.Run(map[bool]string{true: "pass-then-fail", false: "fail-then-pass"}[outcomes[0]], func(t *testing.T) {
			t.Parallel()
			f := resultCheckFixture(t)
			for index, pass := range outcomes {
				if pass {
					if err := os.Remove(filepath.Join(f.workspace, ".fail")); err != nil && !os.IsNotExist(err) {
						t.Fatal(err)
					}
				} else {
					writeBackupTestFile(t, filepath.Join(f.workspace, ".fail"), "fail")
				}
				key := []string{"first", "second"}[index]
				got := checkRequired(t, f, key)
				if pass && got.err != nil || !pass && cli.ExitCode(got.err) != 8 {
					t.Fatalf("check=%#v", got)
				}
				verdict := evaluateRequired(t, f, "FID-9f13aa")
				if verdict.Satisfied("test") != pass || (verdict.Refusal("FID-9f13aa") == nil) != pass {
					t.Fatalf("verdict=%+v, pass=%v", verdict, pass)
				}
			}
			// An interrupted newer start must preserve the last finished verdict.
			options := leanOptions(f, "unfinished-required")
			options.Names = nil
			options.Before = func(string) error { return errors.New("interrupted before spawn") }
			before := evaluateRequired(t, f, "FID-9f13aa")
			if _, err := checkrun.Run(context.Background(), options); err == nil {
				t.Fatal("expected interruption")
			}
			after := evaluateRequired(t, f, "FID-9f13aa")
			if after.LatestResultIDs["test"] != before.LatestResultIDs["test"] || after.Satisfied("test") != outcomes[1] {
				t.Fatalf("unfinished start changed verdict: before=%+v after=%+v", before, after)
			}
		})
	}
}

func TestRequiredResultBelongsToItsTaskOnTheSameTree(t *testing.T) {
	t.Parallel()
	f := resultCheckFixture(t)
	other := *f
	other.workspace = gittest.AddWorktree(t, f.workspace)
	other.seed([]readyCommandRow{statusTask("FID-cccccc", "Other task", "in_progress", nil, nil, nil)}, nil)
	startSQL(t, &other, `UPDATE issues SET lane='standard' WHERE id='FID-cccccc'`)
	other.hold("HLD-bbbbbb", "FID-cccccc", "", time.Minute)
	got := other.run("--json", "check", "--hold", "HLD-bbbbbb", "--actor", "agent:claude-builder", "--idempotency-key", "other-check")
	if names := requiredNames(t, got); !slices.Equal(names, []string{"test"}) {
		t.Fatalf("other checks=%v", names)
	}
	if result := evaluateRequired(t, f, "FID-9f13aa"); result.Satisfied("test") || !slices.Equal(result.Missing, []string{"test"}) || len(result.LatestResultIDs) != 0 {
		t.Fatalf("another task's pass counted: %+v", result)
	}
}
