package main

import (
	"strings"
	"testing"
)

func TestHandoffAfterWaitingTakeoverKeepsReviewReservation(t *testing.T) {
	t.Parallel()
	f := requiredCheckFixture(t, "[]")
	parked := parkWaiting(t, f, "HLD-aaaaaa", waitingQuestion, "wait")
	taken, _ := lifecycleResult(t, f.run("--json", "takeover", "FID-9f13aa", "--actor", "human:mark", "--reason", "Continue here"), "task.takeover")
	lifecycleResult(t, f.run("--json", "handoff", "--hold", taken["hold_id"].(string), "--actor", "human:mark", "--pr", "https://github.com/mbell-nzl/fid/pull/66", "--idempotency-key", "handoff"), "task.handoff")
	before := startEventCount(t, f.workspace)
	for _, row := range []struct {
		operation string
		args      []string
	}{
		{"next.get", []string{"next"}},
		{"next.claim", []string{"next", "--claim", "--actor", "human:mark"}},
	} {
		t.Run(row.operation, func(t *testing.T) {
			got, event := lifecycleResult(t, f.run(append([]string{"--json"}, row.args...)...), row.operation)
			waiting, ok := got["waiting"].(map[string]any)
			if !ok || len(waiting) != 2 || waiting["kind"] != "review" || waiting["task_id"] != "FID-9f13aa" {
				t.Errorf("review waiting = %v", got["waiting"])
			}
			if got["kind"] != "none" || event != nil {
				t.Errorf("review next = %v; event = %v", got, event)
			}
			for _, field := range []string{"hold", "task", "epic", "answer"} {
				if got[field] != nil {
					t.Errorf("review %s = %v, want null", field, got[field])
				}
			}
			text := f.run(row.args...)
			if text.err != nil || !strings.Contains(text.stdout, "Nothing here: FID-9f13aa is in review") {
				t.Errorf("review text = %#v", text)
			}
		})
	}
	status := f.run("status")
	if status.err != nil || !strings.Contains(status.stdout, "in review") {
		t.Fatalf("review status = %#v", status)
	}
	shown, _ := lifecycleResult(t, f.run("--json", "show", "FID-9f13aa"), "issue.show")
	questions := shown["questions"].([]any)
	if len(questions) != 1 {
		t.Fatalf("question history = %v", questions)
	}
	question := questions[0].(map[string]any)
	if question["id"] != parked["question_id"] || question["question"] != waitingQuestion || question["resolved_at"] != nil {
		t.Fatalf("hand-off changed unanswered question = %v", question)
	}
	if startEventCount(t, f.workspace) != before {
		t.Fatal("review next wrote an event")
	}
}
