package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/mbell-nzl/fid/internal/epicplan"
	"github.com/mbell-nzl/fid/internal/epicwork"
	"github.com/mbell-nzl/fid/internal/reviewwork"
	fidmigrations "github.com/mbell-nzl/fid/migrations"
)

func TestEpicValidationVerifyUsesCleanDefaultAndRecordsFailures(t *testing.T) {
	t.Parallel()
	for _, exit := range []string{"0", "1"} {
		t.Run(exit, func(t *testing.T) {
			t.Parallel()
			f, epic := newEpicCancellationFixture(t)
			f.clock = time.Now
			appendPolicy(t, f.workspace, "\n[checks.epicproof]\nrun = ['sh', '-c', 'test \"$(cat fixture.txt)\" = canonical && printf receipt > receipt.txt; exit "+exit+"']\ntimeout = '5s'\nartifacts = ['receipt.txt']\n")
			for path, contents := range map[string]string{"fixture.txt": "canonical\n", ".gitignore": "receipt.txt\n"} {
				if err := os.WriteFile(filepath.Join(f.workspace, path), []byte(contents), 0o600); err != nil {
					t.Fatal(err)
				}
			}
			identityGit(t, f.workspace, "add", ".fid.toml", "fixture.txt", ".gitignore")
			identityGit(t, f.workspace, "commit", "-qm", "Validation policy")
			commit := strings.TrimSpace(identityGit(t, f.workspace, "rev-parse", "HEAD"))
			identityGit(t, f.workspace, "update-ref", "refs/remotes/origin/trunk", commit)
			lifecycleResult(t, f.run("--json", "epic", "validation", "add", epic, "--claim", "Canonical fixture verified", "--verifier", "check:epicproof", "--actor", "agent:builder", "--idempotency-key", "claim"), "epic.validation.add")
			if err := os.WriteFile(filepath.Join(f.workspace, "fixture.txt"), []byte("caller edits\n"), 0o600); err != nil {
				t.Fatal(err)
			}
			before := identityGit(t, f.workspace, "status", "--porcelain=v1")
			worktrees := identityGit(t, f.workspace, "worktree", "list", "--porcelain")
			r, _ := lifecycleResult(t, f.run("--json", "epic", "validation", "verify", epic, "--actor", "agent:builder", "--idempotency-key", "verify"), "epic.validation.verify")
			checks := r["checks"].([]any)
			if len(checks) != 1 {
				t.Fatalf("checks=%v", checks)
			}
			check := checks[0].(map[string]any)
			want := "passed"
			if exit == "1" {
				want = "failed"
			}
			if check["outcome"] != want || check["commit"] != commit {
				t.Fatalf("result=%v", check)
			}
			text := f.run("epic", "validation", "verify", epic, "--actor", "agent:builder", "--idempotency-key", "verify")
			if text.err != nil || !strings.Contains(text.stdout, want) {
				t.Fatalf("text must report %s: %#v", want, text)
			}
			if len(check["check"].(map[string]any)["artifacts"].([]any)) != 1 {
				t.Fatalf("missing artifact: %v", check)
			}
			if after := identityGit(t, f.workspace, "status", "--porcelain=v1"); after != before {
				t.Fatalf("caller tree changed: %s -> %s", before, after)
			}
			if after := identityGit(t, f.workspace, "worktree", "list", "--porcelain"); after != worktrees {
				t.Fatalf("worktree leaked: %s", after)
			}
			data, err := os.ReadFile(filepath.Join(f.workspace, "fixture.txt"))
			if err != nil || string(data) != "caller edits\n" {
				t.Fatalf("caller edit=%s,%v", data, err)
			}
			if r := f.run("rebuild"); r.err != nil {
				t.Fatalf("rebuild=%#v", r)
			}
			q, err := reviewwork.OpenProjection(f.workspace)
			if err != nil {
				t.Fatal(err)
			}
			defer func() { _ = q.Close() }()
			views, err := epicplan.Views(q, epic)
			state := "current"
			if exit == "1" {
				state = "failed"
			}
			if err != nil || len(views) != 1 || views[0].State != state {
				t.Fatalf("views=%+v,%v", views, err)
			}
		})
	}
}

func TestEpicValidationPersonPassGuards(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct{ name, actor, verifier, recipe string }{{"agent", "agent:builder", "review:code", ""}, {"service", "service:ci", "review:code", ""}, {"check", "human:mark", "check:proof", ""}, {"missing recipe", "human:mark", "human", "missing.md"}, {"literal recipe", "human:mark", "human", ":(top).fid.toml"}} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			f, epic := newEpicCancellationFixture(t)
			f.clock = time.Now
			appendPolicy(t, f.workspace, "\n[checks.proof]\nrun=['true']\ntimeout='5s'\n")
			commit := policyCommit(t, f.workspace)
			identityGit(t, f.workspace, "update-ref", "refs/remotes/origin/trunk", commit)
			_, err := epicwork.Change(epicwork.Options{WorkspaceRoot: f.workspace, MigrationsFS: fidmigrations.FS, EpicID: epic, Actor: "agent:builder", IdempotencyKey: "claim", Validation: epicplan.Validation{Claim: "Claim", Verifier: tc.verifier, Recipe: tc.recipe}}, "add")
			if err != nil {
				t.Fatal(err)
			}
			before := snapshotWorkspacePersistence(t, f.workspace)
			r := f.run("--json", "epic", "validation", "pass", epic, "V1", "--actor", tc.actor, "--reason", "Verified")
			if r.err == nil {
				t.Fatalf("invalid pass succeeded: %#v", r)
			}
			if !before.equal(snapshotWorkspacePersistence(t, f.workspace)) {
				t.Fatal("refused pass changed ledger")
			}
		})
	}
}

func TestEpicValidationSoftCapAndDefaultPolicy(t *testing.T) {
	t.Parallel()
	f, epic := newEpicCancellationFixture(t)
	f.clock = time.Now
	appendPolicy(t, f.workspace, "\n[checks.unmerged]\nrun=['true']\ntimeout='1s'\n")
	bad := f.run("--json", "epic", "validation", "add", epic, "--claim", "Unmerged policy", "--verifier", "check:unmerged", "--actor", "agent:builder", "--idempotency-key", "unmerged")
	if bad.err == nil {
		t.Fatal("unmerged check accepted")
	}
	for n := 1; n <= 11; n++ {
		r, _ := lifecycleResult(t, f.run("--json", "epic", "validation", "add", epic, "--claim", "Claim", "--verifier", "review:code", "--actor", "agent:builder", "--idempotency-key", strings.Repeat("x", n)), "epic.validation.add")
		if n == 11 && len(r["warnings"].([]any)) != 1 {
			t.Fatalf("missing soft cap warning: %v", r)
		}
	}
}
