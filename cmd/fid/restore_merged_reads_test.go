package main

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/mbell-nzl/fid/internal/github"
	"github.com/mbell-nzl/fid/internal/gittest"
)

func recordV8OfflineMerge(t *testing.T, root, task string, number int64, failed bool) {
	t.Helper()
	f := &statusFixture{t: t, workspace: root, clock: time.Now}
	facts := doneFacts(t, f, task)
	facts.PR.Number = number
	facts.PR.MergeCommit = fmt.Sprintf("%040x", number)
	facts.CI = ptrV8CI(completedCI(*facts.PR))
	if !failed {
		facts.CI.Runs = []github.WorkflowRun{}
	}
	stub := newReconcileStub(t, root)
	stub.add(facts)
	assertJSONFailureCode(t, stub.run("--json", "done", task, "--actor", "agent:observer").stderr, "task.done", "done-review-required")
}

func addV8FailedMerge(t *testing.T, root string) string {
	t.Helper()
	other := gittest.AddWorktree(t, root)
	f := &statusFixture{t: t, workspace: other, clock: time.Now}
	created, _ := lifecycleResult(t, f.run("--json", "issue", "create", "--kind", "task", "--title", "Merged with failed CI at backup", "--lane", "standard", "--actor", "human:mark", "--idempotency-key", "v8-failed-ci"), "issue.create")
	task := created["issue"].(map[string]any)["id"].(string)
	took, _ := lifecycleResult(t, f.run("--json", "takeover", task, "--actor", "human:mark", "--reason", "Keep failed merge evidence"), "task.takeover")
	hold := took["hold_id"].(string)
	lifecycleResult(t, f.run("--json", "start", "--hold", hold, "--actor", "human:mark", "--idempotency-key", "v8-failed-start"), "task.start")
	lifecycleResult(t, f.run("--json", "handoff", "--hold", hold, "--actor", "human:mark", "--pr", "https://github.com/mbell-nzl/fid/pull/66", "--idempotency-key", "v8-failed-handoff"), "task.handoff")
	recordV8OfflineMerge(t, other, task, 66, true)
	return task
}

// Compare the new offline answers, alongside V8's existing full read comparisons.
// Other status/show fields depend on the caller's worktree and backup health;
// those intentionally differ when restoring into a new clone.
func captureV8MergeReads(t *testing.T, workspace string, ids v8IDs, state *v8State) {
	t.Helper()
	f := &statusFixture{t: t, workspace: workspace, clock: func() time.Time { return ids.now }, home: t.TempDir(), configHome: filepath.Join(t.TempDir(), "config")}
	for _, task := range ids.issues {
		fields, _ := lifecycleResult(t, f.run("--json", "show", task), "issue.show")
		raw, err := json.Marshal(fields["merge"])
		if err != nil {
			t.Fatal(err)
		}
		state.reads["show "+task+" merge"] = string(raw)
	}
	fields, _ := lifecycleResult(t, f.run("--json", "status"), "project.status")
	raw, err := json.Marshal(fields["in_progress"])
	if err != nil {
		t.Fatal(err)
	}
	// Keep only the new merge fields; group/order and worktree-relative data
	// retain their existing V8 coverage rather than changing the comparison.
	var section struct {
		Tasks []map[string]any `json:"tasks"`
		Epics []struct {
			Tasks []map[string]any `json:"tasks"`
		} `json:"epics"`
	}
	if err := json.Unmarshal(raw, &section); err != nil {
		t.Fatal(err)
	}
	for _, epic := range section.Epics {
		section.Tasks = append(section.Tasks, epic.Tasks...)
	}
	merges := map[string]any{}
	for _, task := range section.Tasks {
		if merge, ok := task["merge"]; ok {
			merges[task["id"].(string)] = merge
		}
	}
	raw, err = json.Marshal(merges)
	if err != nil {
		t.Fatal(err)
	}
	state.reads["status merges"] = string(raw)
	if !strings.Contains(string(raw), `"ci":"waiting"`) || !strings.Contains(string(raw), `"ci":"not-green"`) {
		t.Fatalf("V8 missing offline states: %s", raw)
	}
}
