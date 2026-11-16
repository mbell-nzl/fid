package main

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

func addV8Epics(t *testing.T, root string) []string {
	t.Helper()
	f := &statusFixture{t: t, workspace: root, clock: time.Now}
	ids := []string{}
	for _, state := range []string{"closed", "needs-you", "open"} {
		created, _ := lifecycleResult(t, f.run("--json", "issue", "create", "--kind", "epic", "--scheduling", "prioritised", "--title", "Restore epic "+state, "--actor", "human:mark", "--idempotency-key", "v8-epic-"+state), "issue.create")
		epic := created["issue"].(map[string]any)["id"].(string)
		child := epicDoneChild(t, f, epic, "v8-epic-child-"+state)
		if state != "open" {
			lifecycleResult(t, f.run("--json", "issue", "transition", child, "--to", "cancelled", "--actor", "human:mark", "--reason", "Finished without delivery"), "issue.transition")
		}
		if state == "closed" {
			lifecycleResult(t, f.run("--json", "done", epic, "--actor", "human:mark", "--reason", "Epic finished", "--idempotency-key", "v8-epic-close"), "task.done")
		}
		assertEpicNeedsYou(t, f, epic, state == "needs-you")
		ids = append(ids, epic, child)
	}
	return ids
}
func captureV8EpicReads(t *testing.T, root string, ids v8IDs, state *v8State) {
	t.Helper()
	f := &statusFixture{t: t, workspace: root, clock: func() time.Time { return ids.now }, home: t.TempDir(), configHome: t.TempDir()}
	result, _ := lifecycleResult(t, f.run("--json", "status"), "project.status")
	epics := []any{}
	for _, v := range result["needs_you"].([]any) {
		item := v.(map[string]any)
		if item["kind"] == "epic-ready-to-close" {
			epics = append(epics, item)
		}
	}
	if len(epics) != 1 {
		t.Fatal("V8 must have exactly one epic needing a person")
	}
	raw, err := json.Marshal(epics)
	if err != nil {
		t.Fatal(err)
	}
	state.reads["status epic needs you"] = string(raw)
	decisions := f.run("sql", `SELECT * FROM person_decisions WHERE kind='close-epic' ORDER BY id`)
	if decisions.err != nil || decisions.stdout == "[]\n" {
		t.Fatalf("V8 closed epic decision: %s, %v", decisions.stdout, decisions.err)
	}
	state.reads["epic close decisions"] = fmt.Sprint(decisions.stdout)
}
