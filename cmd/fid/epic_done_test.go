package main

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/mbell-nzl/fid/internal/cli"
)

func TestEpicDoneDecisionRetryAndReplay(t *testing.T) {
	t.Parallel()
	for _, key := range []string{"", "close-epic"} {
		t.Run("key="+key, func(t *testing.T) {
			t.Parallel()
			f, epic := newEpicCancellationFixture(t)
			f.clock = time.Now
			args := []string{"--json", "done", epic, "--actor", "human:mark", "--reason", "All work finished"}
			if key != "" {
				args = append(args, "--idempotency-key", key)
			}
			before := startEventCount(t, f.workspace)
			got, event := lifecycleResult(t, f.run(args...), "task.done")
			decision := got["decision"].(map[string]any)
			if got["status"] != "done" || got["revision"] != float64(2) || decision["kind"] != "close-epic" || decision["subject_id"] != epic || decision["reason"] != "All work finished" || decision["actor"] != "human:mark" {
				t.Fatalf("close=%v", got)
			}
			if startEventCount(t, f.workspace) != before+1 {
				t.Fatal("close must append one Event")
			}
			for _, rebuild := range []bool{false, true} {
				if rebuild {
					if r := f.run("rebuild"); r.err != nil {
						t.Fatalf("rebuild=%#v", r)
					}
				}
				snapshot := snapshotWorkspacePersistence(t, f.workspace)
				retry := f.run(args...)
				if key == "" {
					assertJSONFailureCode(t, retry.stderr, "task.done", "done-epic-not-open")
				} else {
					retried, retriedEvent := lifecycleResult(t, retry, "task.done")
					if !reflect.DeepEqual(got, retried) || !reflect.DeepEqual(event, retriedEvent) {
						t.Fatalf("retry=%v,%v want %v,%v", retried, retriedEvent, got, event)
					}
				}
				if !snapshot.equal(snapshotWorkspacePersistence(t, f.workspace)) {
					t.Fatal("retry changed ledger")
				}
			}
		})
	}
}

func TestEpicDoneRefusals(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name, actor, reason, state, code string
		exit                             int
	}{
		{"agent", "agent:builder", "Finished", "", "authority-denied", 5},
		{"service", "service:ci", "Finished", "", "authority-denied", 5},
		{"no reason", "human:mark", "", "", "invalid-argument", 2},
		{"blank reason", "human:mark", " \n ", "", "invalid-argument", 2},
		{"cancelled epic", "human:mark", "Finished", "cancelled", "done-epic-not-open", 4},
		{"todo child", "human:mark", "Finished", "todo", "done-epic-open-tasks", 4},
		{"started child", "human:mark", "Finished", "in_progress", "done-epic-open-tasks", 4},
		{"review child", "human:mark", "Finished", "in_review", "done-epic-open-tasks", 4},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			f, epic := newEpicCancellationFixture(t)
			f.clock = time.Now
			if tc.state == "cancelled" {
				lifecycleResult(t, f.run("--json", "issue", "transition", epic, "--to", "cancelled", "--actor", "human:mark", "--reason", "Abandoned"), "issue.transition")
			} else if tc.state != "" {
				task := epicDoneChild(t, f, epic, "child")
				if tc.state != "todo" {
					took, _ := lifecycleResult(t, f.run("--json", "takeover", task, "--actor", "human:mark", "--reason", "Build"), "task.takeover")
					hold := took["hold_id"].(string)
					lifecycleResult(t, f.run("--json", "start", "--hold", hold, "--actor", "human:mark", "--idempotency-key", "start"), "task.start")
					if tc.state == "in_review" {
						lifecycleResult(t, f.run("--json", "handoff", "--hold", hold, "--actor", "human:mark", "--pr", "https://github.com/mbell-nzl/fid/pull/88", "--idempotency-key", "handoff"), "task.handoff")
					}
				}
			}
			before := snapshotWorkspacePersistence(t, f.workspace)
			got := f.run("--json", "done", epic, "--actor", tc.actor, "--reason", tc.reason)
			assertJSONFailureCode(t, got.stderr, "task.done", tc.code)
			if cli.ExitCode(got.err) != tc.exit || got.stdout != "" {
				t.Fatalf("done=%#v", got)
			}
			if !before.equal(snapshotWorkspacePersistence(t, f.workspace)) {
				t.Fatal("refusal changed ledger")
			}
		})
	}
}

func TestEpicDoneAllCancelledNeedsYouAndText(t *testing.T) {
	t.Parallel()
	f, epic := newEpicCancellationFixture(t)
	f.clock = time.Now
	assertEpicNeedsYou(t, f, epic, false)
	child := epicDoneChild(t, f, epic, "cancelled-child")
	assertEpicNeedsYou(t, f, epic, false)
	lifecycleResult(t, f.run("--json", "issue", "transition", child, "--to", "cancelled", "--actor", "human:mark", "--reason", "Not needed"), "issue.transition")
	assertEpicNeedsYou(t, f, epic, true)
	before := startEventCount(t, f.workspace)
	text := f.run("status")
	if text.err != nil || !strings.Contains(text.stdout, "fid done "+epic) {
		t.Fatalf("status=%#v", text)
	}
	if startEventCount(t, f.workspace) != before {
		t.Fatal("status wrote Event")
	}
	got := f.run("done", epic, "--actor", "human:mark", "--reason", "Finished")
	if got.err != nil || !strings.Contains(got.stdout, "Closed "+epic) || strings.Contains(got.stdout, "claim") || strings.Contains(got.stdout, "runs") {
		t.Fatalf("close=%#v", got)
	}
	assertEpicNeedsYou(t, f, epic, false)
}

func TestEpicTransitionRequiresDone(t *testing.T) {
	t.Parallel()
	for _, actor := range []string{"human:mark", "agent:builder", "service:ci"} {
		t.Run(actor, func(t *testing.T) {
			t.Parallel()
			f, epic := newEpicCancellationFixture(t)
			before := snapshotWorkspacePersistence(t, f.workspace)
			got := f.run("--json", "issue", "transition", epic, "--to", "done", "--actor", actor, "--idempotency-key", "bypass")
			assertJSONFailureCode(t, got.stderr, "issue.transition", "transition-not-permitted")
			if !strings.Contains(got.stderr, "use-done") || !strings.Contains(got.stderr, "fid done "+epic) {
				t.Fatalf("transition=%#v", got)
			}
			if !before.equal(snapshotWorkspacePersistence(t, f.workspace)) {
				t.Fatal("transition changed ledger")
			}
		})
	}
}

func epicDoneChild(t *testing.T, f *statusFixture, epic, key string) string {
	t.Helper()
	r, _ := lifecycleResult(t, f.run("--json", "issue", "create", "--kind", "task", "--epic", epic, "--title", key, "--lane", "routine", "--actor", "human:mark", "--idempotency-key", key), "issue.create")
	return r["issue"].(map[string]any)["id"].(string)
}
func assertEpicNeedsYou(t *testing.T, f *statusFixture, epic string, want bool) {
	t.Helper()
	r, _ := lifecycleResult(t, f.run("--json", "status"), "project.status")
	raw, _ := json.Marshal(r["needs_you"])
	found := false
	for _, v := range r["needs_you"].([]any) {
		item := v.(map[string]any)
		if item["kind"] == "epic-ready-to-close" && item["subject"].(map[string]any)["id"] == epic {
			found = true
		}
	}
	if found != want {
		t.Fatalf("needs_you=%s want epic=%v", raw, want)
	}
}

func TestDoneReasonOnlyAppliesToEpics(t *testing.T) {
	t.Parallel()
	f, task, _ := doneBuildFixture(t, "routine", false)
	for _, id := range []string{"", task} {
		t.Run("id="+id, func(t *testing.T) {
			t.Parallel()
			args := []string{"--json", "done"}
			if id != "" {
				args = append(args, id)
			}
			args = append(args, "--actor", "human:mark", "--reason", "Finished")
			got := f.run(args...)
			assertJSONFailureCode(t, got.stderr, "task.done", "invalid-argument")
			if !strings.Contains(got.stderr, `"argument":"--reason"`) {
				t.Fatalf("reason=%#v", got)
			}
		})
	}
}

func TestEpicDoneKeyBindsSubjectReasonAndActor(t *testing.T) {
	t.Parallel()
	f, epic := newEpicCancellationFixture(t)
	f.clock = time.Now
	args := []string{"--json", "done", epic, "--actor", "human:mark", "--reason", "Finished", "--idempotency-key", "close"}
	lifecycleResult(t, f.run(args...), "task.done")
	for _, field := range []int{2, 4, 6} {
		changed := append([]string{}, args...)
		changed[field] = "human:other"
		before := snapshotWorkspacePersistence(t, f.workspace)
		got := f.run(changed...)
		assertJSONFailureCode(t, got.stderr, "task.done", "idempotency-conflict")
		if !before.equal(snapshotWorkspacePersistence(t, f.workspace)) {
			t.Fatal("conflicting retry changed ledger")
		}
	}
}
