package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mbell-nzl/fid/internal/doneeval"
	"github.com/mbell-nzl/fid/internal/journal"
	"github.com/mbell-nzl/fid/internal/reviewwork"
	fidmigrations "github.com/mbell-nzl/fid/migrations"
)

func TestDoneUsesCurrentWaiverAndReplaysCompletion(t *testing.T) {
	t.Parallel()
	f, task := canonicalReviewFixtureWithLane(t, "assured")
	lifecycleResult(t, recordReviewTask(t, f, reviewInput(), "primary", "human:courier", task), "review.record")
	facts := doneFacts(t, f, task)
	facts.PR.Body = "No rollback section"
	refused := runDone(t, f, doneRunner(t, facts, nil), "--json", "done", task, "--actor", "agent:finisher")
	assertJSONFailureCode(t, refused.stderr, "task.done", "done-separate-review-required")
	waiveSuccess(t, waiveForTest(f, task, task+"\n", true, "--lane", "standard"))
	lifecycleResult(t, runDone(t, f, doneRunner(t, facts, nil), "--json", "done", task, "--actor", "agent:finisher"), "task.done")
	if _, err := journal.RebuildWorkspace(f.workspace, fidmigrations.FS); err != nil {
		t.Fatal(err)
	}
	fields, _ := lifecycleResult(t, f.run("--json", "show", task), "issue.show")
	if fields["issue"].(map[string]any)["status"] != "done" {
		t.Fatal("completion not replayed", fields)
	}
}

func TestDoneWaiverSkipsChecksAndArtifactsNotRequiredByNewLane(t *testing.T) {
	t.Parallel()
	f, task, hold := doneBuildFixture(t, "assured", false, func(f *statusFixture) {
		path := filepath.Join(f.workspace, ".fid.toml")
		text := string(mustReadCommandFile(t, path))
		text = strings.Replace(text, "checks=['evidence']", "checks=[]", 1)
		writeBackupTestFile(t, path, text)
	})
	lifecycleResult(t, f.run(append([]string{"--json"}, doneHandoffArgs(hold)...)...), "task.handoff")
	waiveSuccess(t, waiveForTest(f, task, task+"\n", true))
	if err := os.Remove(keptTaskCheckPath(t, f, task)); err != nil {
		t.Fatal(err)
	}
	facts := doneFacts(t, f, task)
	facts.PR.Body = ""
	lifecycleResult(t, runDone(t, f, doneRunner(t, facts, nil), "--json", "done", task, "--actor", "agent:finisher"), "task.done")
	if _, err := journal.RebuildWorkspace(f.workspace, fidmigrations.FS); err != nil {
		t.Fatal(err)
	}
}

func TestDoneGatherMissingCheckConfigurationStaysEvidence(t *testing.T) {
	t.Parallel()
	for _, waived := range []bool{false, true} {
		t.Run(map[bool]string{false: "handoff lane", true: "waived lane"}[waived], func(t *testing.T) {
			t.Parallel()
			f, task := canonicalReviewFixtureWithLane(t, "assured")
			lane := "assured"
			if waived {
				waiveSuccess(t, waiveForTest(f, task, task+"\n", true))
				lane = "routine"
			}
			startSQL(t, f, `UPDATE policy_snapshots SET effective_json=json_set(effective_json,'$.lanes.`+lane+`.checks_set',json('false'))`)
			conn, err := reviewwork.OpenProjection(f.workspace)
			if err != nil {
				t.Fatal(err)
			}
			defer func() { _ = conn.Close() }()
			snapshot, err := doneeval.Gather(f.workspace, conn, task, doneeval.Target{}, doneeval.GitHubFacts{})
			if err != nil {
				t.Fatalf("missing check configuration escaped the evaluator: %v", err)
			}
			if got := doneeval.Evaluate(snapshot).Gates[doneeval.PolicyGate]; got.Reason != doneeval.ChecksNotConfigured {
				t.Fatalf("missing configuration gate=%+v", got)
			}
		})
	}
}
