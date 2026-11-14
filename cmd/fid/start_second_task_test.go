package main

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/mbell-nzl/fid/internal/gittest"
)

func TestStartSecondTaskKeepsStartedEpicRevisionAcrossReplay(t *testing.T) {
	t.Parallel()
	f := newJournalStatusFixture(t)
	now := time.Now().Add(time.Minute)
	f.clock = func() time.Time { return now }
	created, _ := lifecycleResult(t, f.run("--json", "issue", "create", "--kind", "epic", "--title", "Parallel work", "--scheduling", "prioritised", "--actor", "human:mark", "--idempotency-key", "epic"), "issue.create")
	epic := created["issue"].(map[string]any)["id"].(string)
	tasks := map[string]bool{}
	for _, title := range []string{"First task", "Second task"} {
		created, _ := lifecycleResult(t, f.run("--json", "issue", "create", "--kind", "task", "--epic", epic, "--title", title, "--lane", "standard", "--actor", "human:mark", "--idempotency-key", title), "issue.create")
		tasks[created["issue"].(map[string]any)["id"].(string)] = false
	}
	main := f.workspace
	for _, workspace := range []string{main, gittest.AddWorktree(t, main)} {
		f.workspace = workspace
		taken, _ := lifecycleResult(t, f.run("--json", "next", "--claim", "--actor", "agent:builder"), "next.claim")
		id := taken["task"].(map[string]any)["id"].(string)
		if started, exists := tasks[id]; !exists || started {
			t.Fatalf("next chose %s; task starts = %v", id, tasks)
		}
		hold := taken["hold"].(map[string]any)["id"].(string)
		started, event := lifecycleResult(t, f.run("--json", "start", "--hold", hold, "--actor", "agent:builder", "--idempotency-key", "start-"+id), "task.start")
		if started["task_id"] != id || started["status"] != "in_progress" || event == nil {
			t.Fatalf("start %s = %v/%v", id, started, event)
		}
		tasks[id] = true
		shown, _ := lifecycleResult(t, f.run("--json", "show", epic), "issue.show")
		issue := shown["issue"].(map[string]any)
		if issue["status"] != "in_progress" || issue["revision"] != float64(2) {
			t.Fatalf("epic after start %s = %v, want in_progress at revision 2", id, issue)
		}
	}
	query := "SELECT id,status,revision,lane FROM issues ORDER BY id"
	before := f.run("sql", query)
	if rebuilt := f.run("rebuild"); rebuilt.err != nil {
		t.Fatalf("rebuild = %#v", rebuilt)
	}
	if after := f.run("sql", query); before.err != nil || after.err != nil || before.stdout != after.stdout {
		t.Fatalf("replayed issues differ: before=%#v after=%#v", before, after)
	}
	for id := range tasks {
		shown, _ := lifecycleResult(t, f.run("--json", "show", id), "issue.show")
		issue := shown["issue"].(map[string]any)
		if issue["status"] != "in_progress" || issue["revision"] != float64(2) {
			t.Fatalf("replayed task %s = %v", id, issue)
		}
	}
	archive := filepath.Join(t.TempDir(), "two-starts.zip")
	if result := f.run("export", "create", archive); result.err != nil {
		t.Fatalf("export create = %#v", result)
	}
	if result := f.run("export", "validate", archive); result.err != nil {
		t.Fatalf("export validate = %#v", result)
	}
}
