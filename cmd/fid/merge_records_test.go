package main

import (
	"database/sql"
	"encoding/json"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/mbell-nzl/fid/internal/config"
	fiddb "github.com/mbell-nzl/fid/internal/db"
	"github.com/mbell-nzl/fid/internal/github"
	"github.com/mbell-nzl/fid/internal/mergework"
)

func mergeFixture(t *testing.T) (*statusFixture, mergework.Options, string, github.PullRequest, string) {
	t.Helper()
	f := requiredCheckFixture(t, "[]")
	fields, _ := lifecycleResult(t, f.run(append([]string{"--json"}, handoffArgs("merge-fixture")...)...), "task.handoff")
	conn, err := fiddb.Open(filepath.Join(config.LedgerDir(f.workspace), config.DatabaseFileName))
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = conn.Close() }()
	var raw string
	if err := conn.QueryRow(`SELECT record_json FROM handoffs WHERE id=?`, fields["handoff_id"]).Scan(&raw); err != nil {
		t.Fatal(err)
	}
	var handoff struct {
		Tree string `json:"tree"`
	}
	if err := json.Unmarshal([]byte(raw), &handoff); err != nil {
		t.Fatal(err)
	}
	now := statusNow.Add(time.Minute)
	pr := github.PullRequest{Number: 63, URL: fields["pr"].(string), State: "MERGED", Merged: true, MergedAt: &now, MergeCommit: strings.Repeat("a", 40), HeadSHA: fields["head"].(string), BaseRef: "main", BaseRepository: github.Repository{Owner: "mbell-nzl", Name: "fid"}}
	options := mergework.Options{WorkspaceRoot: f.workspace, MigrationsDir: issueMigrationsDir(), TaskID: "FID-9f13aa", Actor: "agent:observer", Clock: func() time.Time { return now }}
	return f, options, fields["handoff_id"].(string), pr, handoff.Tree
}

func mergeRows(t *testing.T, workspace string) string {
	t.Helper()
	var state strings.Builder
	for _, table := range []string{"merge_observations", "ci_results", "check_requests", "check_starts", "check_processes", "check_results"} {
		got := executeCommandForTest(workspace, "sql", "SELECT * FROM "+table+" ORDER BY rowid")
		if got.err != nil {
			t.Fatal(got.err)
		}
		state.WriteString(got.stdout)
	}
	return state.String()
}

func completedCI(pr github.PullRequest) github.CI {
	return github.CI{Commit: pr.MergeCommit, Runs: []github.WorkflowRun{{ID: 101, Attempt: 1, WorkflowID: 10, Name: "CI", Path: ".github/workflows/ci.yml", Event: "push", HeadSHA: pr.MergeCommit, Status: "completed", Conclusion: "failure", URL: "https://github.com/mbell-nzl/fid/actions/runs/101"}}}
}

func TestMergeObservationRefusesInvalidBindingsWithoutMutation(t *testing.T) {
	t.Parallel()
	for _, name := range []string{"missing task", "not task", "missing handoff", "other task handoff", "wrong PR", "wrong number", "wrong repository", "wrong tree", "short merge", "unmerged", "service"} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			f, options, handoff, pr, tree := mergeFixture(t)
			switch name {
			case "missing task":
				options.TaskID = "FID-absent"
			case "not task":
				_, created := runFIDCommand(t, f.workspace, "issue", "create", "--kind", "epic", "--title", "Other epic", "--scheduling", "prioritised", "--actor", "human:mark")
				options.TaskID = parseCreatedIssueID(t, created)
			case "missing handoff":
				handoff = "HOF-absent"
			case "other task handoff":
				_, created := runFIDCommand(t, f.workspace, "issue", "create", "--kind", "task", "--title", "Other task", "--actor", "human:mark")
				options.TaskID = parseCreatedIssueID(t, created)
			case "wrong PR":
				pr.URL = "https://github.com/mbell-nzl/fid/pull/64"
				pr.Number = 64
			case "wrong number":
				pr.Number++
			case "wrong repository":
				pr.BaseRepository.Name = "other"
			case "wrong tree":
				tree = strings.Repeat("b", 40)
			case "short merge":
				pr.MergeCommit = "abc"
			case "unmerged":
				pr.Merged = false
			case "service":
				options.Actor = "service:observer"
			}
			before, events := mergeRows(t, f.workspace), startEventCount(t, f.workspace)
			if _, err := mergework.RecordMerge(options, handoff, pr, tree); err == nil {
				t.Fatal("accepted invalid merge")
			}
			if after := mergeRows(t, f.workspace); after != before || startEventCount(t, f.workspace) != events {
				t.Fatal("refusal changed ledger")
			}
		})
	}
}

func TestCIObservationsAreImmutableAndDedupeAttempts(t *testing.T) {
	t.Parallel()
	f, options, handoff, pr, tree := mergeFixture(t)
	first, err := mergework.RecordMerge(options, handoff, pr, tree)
	if err != nil {
		t.Fatal(err)
	}
	events := startEventCount(t, f.workspace)
	options.Clock = func() time.Time { return statusNow.Add(2 * time.Minute) }
	again, err := mergework.RecordMerge(options, handoff, pr, tree)
	if err != nil || !reflect.DeepEqual(first, again) || startEventCount(t, f.workspace) != events {
		t.Fatalf("merge retry=%+v/%v", again, err)
	}
	ci := completedCI(pr)
	failed, err := mergework.RecordCI(options, first.ID, pr.BaseRepository, ci)
	if err != nil || len(failed) != 1 {
		t.Fatalf("failed=%+v/%v", failed, err)
	}
	ci.Runs[0].Attempt = 2
	ci.Runs[0].Conclusion = "success"
	passed, err := mergework.RecordCI(options, first.ID, pr.BaseRepository, ci)
	if err != nil || len(passed) != 1 || passed[0].ID == failed[0].ID {
		t.Fatalf("rerun=%+v/%v", passed, err)
	}
	events = startEventCount(t, f.workspace)
	options.Clock = func() time.Time { return statusNow.Add(3 * time.Minute) }
	retried, err := mergework.RecordCI(options, first.ID, pr.BaseRepository, ci)
	if err != nil || !reflect.DeepEqual(passed, retried) || startEventCount(t, f.workspace) != events {
		t.Fatalf("CI retry=%+v/%v", retried, err)
	}
	ci.Runs[0].Attempt = 3
	ci.Runs[0].Status = "in_progress"
	ci.Runs[0].Conclusion = ""
	pending, err := mergework.RecordCI(options, first.ID, pr.BaseRepository, ci)
	if err != nil || len(pending) != 0 || startEventCount(t, f.workspace) != events {
		t.Fatalf("pending=%+v/%v", pending, err)
	}
	conn, err := fiddb.Open(filepath.Join(config.LedgerDir(f.workspace), config.DatabaseFileName))
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = conn.Close() }()
	latest, err := mergework.LatestMerge(conn, options.TaskID)
	if err != nil || !reflect.DeepEqual(latest, first) {
		t.Fatalf("latest=%+v/%v", latest, err)
	}
	results, err := mergework.CIResults(conn, options.TaskID, first.ID)
	if err != nil || !reflect.DeepEqual(results, append(failed, passed...)) {
		t.Fatalf("results=%+v/%v", results, err)
	}
	if _, err := mergework.LatestMerge(conn, "absent"); err != sql.ErrNoRows {
		t.Fatalf("absent=%v", err)
	}
}

func TestCIObservationRefusesInvalidBindingsAtomically(t *testing.T) {
	t.Parallel()
	for _, name := range []string{"wrong task", "missing merge", "wrong repository", "wrong commit", "wrong head", "bad conclusion", "bad status", "bad attempt", "bad workflow", "bad URL", "bad name", "bad path", "changed retry"} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			f, options, handoff, pr, tree := mergeFixture(t)
			m, err := mergework.RecordMerge(options, handoff, pr, tree)
			if err != nil {
				t.Fatal(err)
			}
			ci, repo := completedCI(pr), pr.BaseRepository
			id := m.ID
			switch name {
			case "wrong task":
				options.TaskID = "FID-abcdef"
			case "missing merge":
				id = "MRG-absent"
			case "wrong repository":
				repo.Name = "other"
			case "wrong commit":
				ci.Commit = strings.Repeat("c", 40)
			case "wrong head":
				ci.Runs[0].HeadSHA = strings.Repeat("c", 40)
			case "bad conclusion":
				ci.Runs[0].Conclusion = "green"
			case "bad status":
				ci.Runs[0].Status = "green"
			case "bad attempt":
				ci.Runs[0].Attempt = 0
			case "bad workflow":
				ci.Runs[0].WorkflowID = 0
			case "bad URL":
				ci.Runs[0].URL = "https://example.com"
			case "bad name":
				ci.Runs[0].Name = ""
			case "bad path":
				ci.Runs[0].Path = ""
			case "changed retry":
				if _, err := mergework.RecordCI(options, id, repo, ci); err != nil {
					t.Fatal(err)
				}
				ci.Runs[0].Conclusion = "success"
			}
			// Put a valid new result before the invalid one to detect partial writes.
			valid := completedCI(pr).Runs[0]
			valid.ID = 102
			valid.URL = "https://github.com/mbell-nzl/fid/actions/runs/102"
			ci.Runs = append([]github.WorkflowRun{valid}, ci.Runs...)
			before, events := mergeRows(t, f.workspace), startEventCount(t, f.workspace)
			if _, err := mergework.RecordCI(options, id, repo, ci); err == nil {
				t.Fatal("accepted invalid CI")
			}
			if mergeRows(t, f.workspace) != before || startEventCount(t, f.workspace) != events {
				t.Fatal("refusal changed ledger")
			}
		})
	}
}

func TestMergeHistoryUsesJournalOrderAndScopesCIToObservation(t *testing.T) {
	t.Parallel()
	f, options, handoff, pr, tree := mergeFixture(t)
	first, err := mergework.RecordMerge(options, handoff, pr, tree)
	if err != nil {
		t.Fatal(err)
	}
	ci := completedCI(pr)
	if _, err := mergework.RecordCI(options, first.ID, pr.BaseRepository, ci); err != nil {
		t.Fatal(err)
	}
	// The later observation can describe an older merge time. A rebased
	// head remains acceptable when it preserves the hand-off's tree.
	older := pr.MergedAt.Add(-time.Minute)
	pr.MergedAt = &older
	pr.MergeCommit = strings.Repeat("e", 40)
	pr.HeadSHA = strings.Repeat("f", 40)
	second, err := mergework.RecordMerge(options, handoff, pr, tree)
	if err != nil {
		t.Fatal(err)
	}
	ci.Commit = pr.MergeCommit
	ci.Runs[0].HeadSHA = pr.MergeCommit
	if _, err := mergework.RecordCI(options, second.ID, pr.BaseRepository, ci); err != nil {
		t.Fatal(err)
	}
	conn, err := fiddb.Open(filepath.Join(config.LedgerDir(f.workspace), config.DatabaseFileName))
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = conn.Close() }()
	latest, err := mergework.LatestMerge(conn, options.TaskID)
	if err != nil || !reflect.DeepEqual(latest, second) {
		t.Fatalf("latest=%+v/%v", latest, err)
	}
	for _, m := range []struct{ id, head string }{{first.ID, first.MergeCommit}, {second.ID, second.MergeCommit}} {
		results, err := mergework.CIResults(conn, options.TaskID, m.id)
		if err != nil || len(results) != 1 || results[0].HeadSHA != m.head {
			t.Fatalf("scoped CI=%+v/%v", results, err)
		}
	}
}
