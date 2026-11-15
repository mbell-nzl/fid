package main

import (
	"encoding/json"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/mbell-nzl/fid/internal/config"
	fiddb "github.com/mbell-nzl/fid/internal/db"
	"github.com/mbell-nzl/fid/internal/github"
	"github.com/mbell-nzl/fid/internal/handoffrecord"
	"github.com/mbell-nzl/fid/internal/mergework"
)

func recordV8MergeCI(t *testing.T, workspace, task string) {
	t.Helper()
	conn, err := fiddb.Open(filepath.Join(config.LedgerDir(workspace), config.DatabaseFileName))
	if err != nil {
		t.Fatal(err)
	}
	var raw string
	if err := conn.QueryRow(`SELECT record_json FROM handoffs WHERE task_id=? ORDER BY rowid DESC LIMIT 1`, task).Scan(&raw); err != nil {
		t.Fatal(err)
	}
	if err := conn.Close(); err != nil {
		t.Fatal(err)
	}
	var handoff handoffrecord.Record
	if err := json.Unmarshal([]byte(raw), &handoff); err != nil {
		t.Fatal(err)
	}
	now := time.Now().UTC().Truncate(time.Second)
	options := mergework.Options{WorkspaceRoot: workspace, MigrationsDir: issueMigrationsDir(), TaskID: task, Actor: "agent:observer", Clock: func() time.Time { return now }}
	pr := github.PullRequest{Number: 63, URL: handoff.PR, State: "MERGED", Merged: true, MergedAt: &now, MergeCommit: strings.Repeat("d", 40), HeadSHA: handoff.Head, BaseRef: "main", BaseRepository: github.Repository{Owner: "mbell-nzl", Name: "fid"}}
	m, err := mergework.RecordMerge(options, handoff.ID, pr, handoff.Tree)
	if err != nil {
		t.Fatal(err)
	}
	facts := completedCI(pr)
	rerun := facts.Runs[0]
	rerun.Attempt = 2
	rerun.Conclusion = "success"
	success := rerun
	success.ID = 102
	success.Attempt = 1
	success.URL = "https://github.com/mbell-nzl/fid/actions/runs/102"
	facts.Runs = append(facts.Runs, rerun, success)
	if results, err := mergework.RecordCI(options, m.ID, pr.BaseRepository, facts); err != nil || len(results) != 3 {
		t.Fatalf("V8 CI=%+v/%v", results, err)
	}
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
	m, err := mergework.LatestMerge(conn, task)
	if err != nil {
		t.Fatal(err)
	}
	results, err := mergework.CIResults(conn, task, m.ID)
	if err != nil {
		t.Fatal(err)
	}
	if m.HandoffID == "" || m.Number != 63 || m.MergeCommit != strings.Repeat("d", 40) || len(results) != 3 {
		t.Fatalf("V8 merge/CI=%+v/%+v", m, results)
	}
	if results[0].Conclusion != "failure" || results[1].Conclusion != "success" || results[1].Attempt != 2 || results[0].RunID != results[1].RunID || results[2].RunID == results[0].RunID {
		t.Fatalf("V8 lost failure, success or rerun: %+v", results)
	}
}
