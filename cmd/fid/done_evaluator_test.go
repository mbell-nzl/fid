package main

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/mbell-nzl/fid/internal/doneeval"
	"github.com/mbell-nzl/fid/internal/github"
	"github.com/mbell-nzl/fid/internal/journal"
	"github.com/mbell-nzl/fid/internal/mergework"
	"github.com/mbell-nzl/fid/internal/reviewrecord"
	"github.com/mbell-nzl/fid/internal/reviewwork"
)

// All evaluator evidence is journal-backed: no SQL seeds, even for work holds.
func doneBuildFixture(t *testing.T, lane string, trigger bool, prepare ...func(*statusFixture)) (*statusFixture, string, string) {
	t.Helper()
	f := newJournalStatusFixture(t)
	f.clock = time.Now
	appendPolicy(t, f.workspace, "\n[review]\ntest_adequacy=['**/*_test.go']\n[checks.evidence]\nrun=['true']\ntimeout='1m'\nartifacts=['notes.md']\n[lanes.routine]\npaths=['**/*.md','**/*_test.go']\nchecks=['evidence']\n[lanes.standard]\nchecks=['evidence']\n[lanes.assured]\nchecks=['evidence']\n")
	for _, setup := range prepare {
		setup(f)
	}
	policyOrigin(t, f.workspace)
	runGitForTest(t, f.workspace, "checkout", "-qb", "task")
	writeBackupTestFile(t, filepath.Join(f.workspace, "notes.md"), "A task change\n")
	if trigger {
		writeBackupTestFile(t, filepath.Join(f.workspace, "change_test.go"), "package change\n")
	}
	runGitForTest(t, f.workspace, "add", ".")
	runGitForTest(t, f.workspace, "commit", "-qm", "Task change")
	created, _ := lifecycleResult(t, f.run("--json", "issue", "create", "--kind", "task", "--title", "Done evaluator subject", "--lane", lane, "--actor", "human:mark", "--idempotency-key", "task"), "issue.create")
	task := created["issue"].(map[string]any)["id"].(string)
	took, _ := lifecycleResult(t, f.run("--json", "takeover", task, "--actor", "human:builder", "--model", "codex/gpt-6.1-sol", "--reason", "Build the task", "--idempotency-key", "takeover"), "task.takeover")
	hold := took["hold_id"].(string)
	lifecycleResult(t, f.run("--json", "start", "--hold", hold, "--actor", "human:builder", "--idempotency-key", "start"), "task.start")
	lifecycleResult(t, f.run("--json", "check", "--hold", hold, "--actor", "human:builder", "--idempotency-key", "checks"), "check.execute")
	return f, task, hold
}

func doneHandoffArgs(hold string) []string {
	return []string{"handoff", "--hold", hold, "--actor", "human:builder", "--pr", "https://github.com/MBELL-NZL/FID/pull/63", "--idempotency-key", "handoff"}
}

func doneRemoteFixture(t *testing.T, f *statusFixture, task string) (doneeval.Target, doneeval.GitHubFacts) {
	t.Helper()
	conn, err := reviewwork.OpenProjection(f.workspace)
	if err != nil {
		t.Fatal(err)
	}
	h, err := reviewrecord.LatestHandoff(conn, task)
	_ = conn.Close()
	if err != nil {
		t.Fatal(err)
	}
	now := time.Now().UTC()
	pr := github.PullRequest{Number: 63, URL: "https://github.com/mbell-nzl/fid/pull/63", State: "MERGED", Merged: true, MergedAt: &now, MergeCommit: strings.Repeat("a", 40), HeadSHA: h.Head, BaseRef: "trunk", BaseRepository: github.Repository{Owner: "mbell-nzl", Name: "fid"}, Body: "## Rollback\nRevert the change."}
	options := mergework.Options{WorkspaceRoot: f.workspace, MigrationsDir: issueMigrationsDir(), TaskID: task, Actor: "human:observer", Clock: time.Now}
	m, err := mergework.RecordMerge(options, h.ID, pr, h.Tree)
	if err != nil {
		t.Fatal(err)
	}
	ci := completedCI(pr)
	ci.Runs[0].Conclusion = "success"
	if _, err := mergework.RecordCI(options, m.ID, pr.BaseRepository, ci); err != nil {
		t.Fatal(err)
	}
	return doneeval.Target{Repository: pr.BaseRepository, DefaultBranch: "trunk"}, doneeval.GitHubFacts{PR: &pr, HeadTree: h.Tree, CI: &ci}
}

func gatherDone(t *testing.T, f *statusFixture, task string, target doneeval.Target, remote doneeval.GitHubFacts) doneeval.Snapshot {
	t.Helper()
	conn, err := reviewwork.OpenProjection(f.workspace)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = conn.Close() }()
	// The command can hold this transaction under its write lock, then revalidate
	// task revision and latest hand-off before it appends completion in 1.6i.
	tx, err := conn.Begin()
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = tx.Rollback() }()
	s, err := doneeval.Gather(f.workspace, tx, task, target, remote)
	if err != nil {
		t.Fatal(err)
	}
	return s
}

func TestDoneGatherLaneMatrixAndReplay(t *testing.T) {
	t.Parallel()
	for _, name := range []string{"routine", "routine-tests", "standard", "assured"} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			lane := strings.TrimSuffix(name, "-tests")
			f, task, hold := doneBuildFixture(t, lane, name == "routine-tests")
			lifecycleResult(t, f.run(append([]string{"--json"}, doneHandoffArgs(hold)...)...), "task.handoff")
			if name != "routine" {
				lifecycleResult(t, recordReviewTask(t, f, reviewInput(), "primary", "human:courier", task), "review.record")
			}
			if lane == "assured" {
				in := reviewInput()
				in["stream"] = "separate"
				in["reviewer"].(map[string]any)["session_id"] = "separate-session"
				lifecycleResult(t, recordReviewTask(t, f, in, "separate", "human:courier", task), "review.record")
				approveSuccess(t, approveForTest(f, task, task+"\n", true))
			}
			target, remote := doneRemoteFixture(t, f, task)
			before := startEventCount(t, f.workspace)
			snapshot := gatherDone(t, f, task, target, remote)
			got := doneeval.Evaluate(snapshot)
			if got.Status != doneeval.Complete || got.NoOtherModelReview {
				t.Fatalf("outcome=%+v", got)
			}
			if snapshot.Merge == nil || len(snapshot.CIResults) != 1 || len(snapshot.Checks) != 1 || snapshot.Handoff.PolicySnapshot == 0 {
				t.Fatalf("snapshot=%+v", snapshot)
			}
			if snapshot.Handoff.PR != "https://github.com/MBELL-NZL/FID/pull/63" || snapshot.Merge.PR != remote.PR.URL {
				t.Fatalf("PR casing fixture changed: hand-off=%s merge=%s", snapshot.Handoff.PR, snapshot.Merge.PR)
			}
			if _, err := journal.RebuildWorkspace(f.workspace, issueMigrationsDir()); err != nil {
				t.Fatal(err)
			}
			replay := gatherDone(t, f, task, target, remote)
			if !reflect.DeepEqual(snapshot, replay) {
				t.Fatalf("snapshot changed on replay\nbefore=%+v\nafter=%+v", snapshot, replay)
			}
			path := keptTaskCheckPath(t, f, task)
			if err := os.Remove(path); err != nil {
				t.Fatal(err)
			}
			missing := doneeval.Evaluate(gatherDone(t, f, task, target, remote))
			if missing.Status != doneeval.Refused || missing.Gates[doneeval.ArtifactsGate].Reason != doneeval.KeptArtifactUnavailable {
				t.Fatalf("missing Artifact=%+v", missing)
			}
			if startEventCount(t, f.workspace) != before {
				t.Fatal("gather or evaluation wrote records")
			}
		})
	}
}
