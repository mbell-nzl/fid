package checkread

import (
	"encoding/json"
	"path/filepath"
	"slices"
	"testing"
	"time"

	"github.com/mbell-nzl/fid/internal/checkrecord"
	"github.com/mbell-nzl/fid/internal/config"
	fiddb "github.com/mbell-nzl/fid/internal/db"
)

func countingRecord(id, task, outcome string) checkrecord.Record {
	definition := checkrecord.Definition{Run: []string{"true"}, Timeout: time.Minute, Environment: checkrecord.EnvironmentNames(), Artifacts: []string{}}
	return checkrecord.Record{Start: checkrecord.Start{ID: id, TaskID: task, Name: "test", Definition: definition, DefinitionDigest: checkrecord.Digest(definition), Tree: checkrecord.Tree{Tree: "tree"}}, Result: &checkrecord.Result{Outcome: outcome, Counts: outcome == "passed", Tree: checkrecord.Tree{Tree: "tree"}, FinishedAt: "2026-11-14T00:00:00Z"}}
}

func TestEvaluateLatestFinishedResultOnTaskTreeAndDefinition(t *testing.T) {
	t.Parallel()
	passed := countingRecord("CHK-100001", "task", "passed")
	failed := countingRecord("CHK-100002", "task", "failed")
	// Caller times are not journal order.
	failed.Result.FinishedAt = "2026-11-13T00:00:00Z"
	unfinished := countingRecord("CHK-100003", "task", "passed")
	unfinished.Result = nil
	other := countingRecord("CHK-100004", "other", "passed")
	stale := countingRecord("CHK-100005", "task", "passed")
	stale.Start.DefinitionDigest = "old definition"
	changed := countingRecord("CHK-100006", "task", "passed")
	changed.Start.Tree.Tree, changed.Result.Tree.Tree = "other tree", "other tree"
	unqualified := countingRecord("CHK-100007", "task", "passed")
	unqualified.Result.Counts = false
	for _, row := range []struct {
		name    string
		records []checkrecord.Record
		chosen  string
		missing bool
	}{
		{"passed then failed", []checkrecord.Record{passed, failed}, failed.Start.ID, true},
		{"failed then passed", []checkrecord.Record{failed, passed}, passed.Start.ID, false},
		{"unfinished newer start", []checkrecord.Record{passed, unfinished}, passed.Start.ID, false},
		{"another task", []checkrecord.Record{other}, "", true},
		{"stale definition", []checkrecord.Record{stale}, "", true},
		{"different tree", []checkrecord.Record{changed}, "", true},
		{"latest unqualified pass", []checkrecord.Record{passed, unqualified}, unqualified.Start.ID, true},
		{"irrelevant later records", []checkrecord.Record{passed, other, stale, changed}, passed.Start.ID, false},
	} {
		t.Run(row.name, func(t *testing.T) {
			t.Parallel()
			checks := map[string]config.Check{"test": {Run: []string{"true"}, Timeout: time.Minute}}
			got := Evaluate("task", "tree", []string{"test"}, checks, row.records)
			if !slices.Equal(got.Required, []string{"test"}) || got.LatestResultIDs["test"] != row.chosen || (len(got.Missing) > 0) != row.missing {
				t.Fatalf("evaluation=%+v", got)
			}
			refusal := got.Refusal("task")
			if (refusal != nil) != row.missing {
				t.Fatalf("refusal=%v", refusal)
			}
			if refusal != nil && (refusal.Code != "checks-unsatisfied" || refusal.TaskID != "task") {
				t.Fatalf("refusal=%+v", refusal)
			}
		})
	}
}

func TestReadUsesFinishJournalOrderRatherThanStartOrder(t *testing.T) {
	t.Parallel()
	conn, err := fiddb.Open(filepath.Join(t.TempDir(), "checks.sqlite"))
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = conn.Close() }()
	for _, query := range []string{
		`CREATE TABLE check_starts (id TEXT PRIMARY KEY,task_id TEXT,start_json TEXT)`,
		`CREATE TABLE check_processes (check_id TEXT PRIMARY KEY,process_json TEXT)`,
		`CREATE TABLE check_results (check_id TEXT PRIMARY KEY,result_json TEXT)`,
	} {
		if _, err := conn.Exec(query); err != nil {
			t.Fatal(err)
		}
	}
	first := countingRecord("CHK-100001", "task", "failed")
	second := countingRecord("CHK-100002", "task", "passed")
	for _, record := range []checkrecord.Record{first, second} {
		raw, _ := json.Marshal(record.Start)
		if _, err := conn.Exec(`INSERT INTO check_starts VALUES (?,?,?)`, record.Start.ID, "task", string(raw)); err != nil {
			t.Fatal(err)
		}
	}
	for _, record := range []checkrecord.Record{second, first} {
		raw, _ := json.Marshal(record.Result)
		if _, err := conn.Exec(`INSERT INTO check_results VALUES (?,?)`, record.Start.ID, string(raw)); err != nil {
			t.Fatal(err)
		}
	}
	records, err := Read(conn, "task")
	if err != nil {
		t.Fatal(err)
	}
	if got := Latest(records); len(got) != 1 || got[0].Start.ID != first.Start.ID {
		t.Fatalf("latest finished result=%+v, want %s", got, first.Start.ID)
	}
}
