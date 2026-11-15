package main

import (
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/mbell-nzl/fid/internal/completionrecord"
	"github.com/mbell-nzl/fid/internal/config"
	fiddb "github.com/mbell-nzl/fid/internal/db"
	"github.com/mbell-nzl/fid/internal/github"
	"github.com/mbell-nzl/fid/internal/mergework"
)

func recordV8MergeCI(t *testing.T, workspace, task string) {
	t.Helper()
	f := &statusFixture{t: t, workspace: workspace, clock: time.Now}
	facts := doneFacts(t, f, task)
	facts.PR.MergeCommit = strings.Repeat("d", 40)
	facts.CI = ptrV8CI(completedCI(*facts.PR))
	first := doneChild(t, workspace, facts, "--json", "done", task, "--actor", "agent:observer")
	assertJSONFailureCode(t, first.stderr, "task.done", "done-review-required")
	rerun := facts.CI.Runs[0]
	rerun.Attempt = 2
	rerun.Conclusion = "success"
	success := rerun
	success.ID = 102
	success.Attempt = 1
	success.URL = "https://github.com/mbell-nzl/fid/actions/runs/102"
	facts.CI.Runs = append(facts.CI.Runs, rerun, success)
	second := doneChild(t, workspace, facts, "--json", "done", task, "--actor", "agent:observer")
	assertJSONFailureCode(t, second.stderr, "task.done", "done-review-required")

}

func assertV8MergeCI(t *testing.T, workspace string) {
	t.Helper()
	conn, err := fiddb.Open(filepath.Join(config.LedgerDir(workspace), config.DatabaseFileName))
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = conn.Close() }()
	var task string
	if err := conn.QueryRow(`SELECT task_id FROM merge_observations`).Scan(&task); err != nil {
		t.Fatal(err)
	}
	var firstMerge string
	if err := conn.QueryRow(`SELECT id FROM merge_observations WHERE task_id=? ORDER BY rowid LIMIT 1`, task).Scan(&firstMerge); err != nil {
		t.Fatal(err)
	}
	m, err := mergework.LatestMerge(conn, task)
	if err != nil {
		t.Fatal(err)
	}
	results, err := mergework.CIResults(conn, task, firstMerge)
	if err != nil {
		t.Fatal(err)
	}
	if m.HandoffID == "" || m.Number != 63 || m.MergeCommit != strings.Repeat("d", 40) || len(results) != 3 {
		t.Fatalf("V8 merge/CI=%+v/%+v", m, results)
	}
	completion, err := completionrecord.Read(conn, task)
	if err != nil || completion.MergeID != m.ID || completion.HandoffID != m.HandoffID || len(completion.CIResultIDs) != 1 || len(completion.RoundIDs) != 2 || completion.ApprovalID == "" {
		t.Fatalf("V8 completion=%+v/%v", completion, err)
	}
	var status string
	if err := conn.QueryRow(`SELECT status FROM issues WHERE id=?`, task).Scan(&status); err != nil {
		t.Fatal(err)
	}
	if status != "done" {
		t.Fatalf("V8 completed task=%s", status)
	}
	if results[0].Conclusion != "failure" || results[1].Conclusion != "success" || results[1].Attempt != 2 || results[0].RunID != results[1].RunID || results[2].RunID == results[0].RunID {
		t.Fatalf("V8 lost failure, success or rerun: %+v", results)
	}
}

func ptrV8CI(ci github.CI) *github.CI { return &ci }

func recordV8Completion(t *testing.T, root, task string) {
	t.Helper()
	f := &statusFixture{t: t, workspace: root, clock: time.Now}
	facts := doneFacts(t, f, task)
	facts.PR.MergeCommit = strings.Repeat("d", 40)
	ci := completedCI(*facts.PR)
	ci.Runs[0].Conclusion = "success"
	facts.CI = &ci
	lifecycleResult(t, doneChild(t, root, facts, "--json", "done", task, "--actor", "agent:observer", "--idempotency-key", "v8-done"), "task.done")
}
