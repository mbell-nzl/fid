package main

import (
	"reflect"
	"testing"
	"time"

	"github.com/mbell-nzl/fid/internal/gittest"
)

func TestWaitingAnswerSurvivesUnstartedInterruption(t *testing.T) {
	t.Parallel()
	for _, test := range []struct {
		name, interruption string
		elsewhere          bool
	}{
		{"expiry-here", "expiry", false},
		{"expiry-elsewhere", "expiry", true},
		{"park-here", "park", false},
		{"park-elsewhere", "park", true},
	} {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			f := newLifecycleFixture(t)
			f.seed([]readyCommandRow{statusTask("FID-aaaaaa", "Waiting work", "todo", nil, nil, nil)}, nil)
			other := gittest.AddWorktree(t, f.workspace)
			f.hold("HLD-aaaaaa", "FID-aaaaaa", "", time.Minute)
			parked := parkWaiting(t, f, "HLD-aaaaaa", waitingQuestion, "wait")
			answered, _ := lifecycleResult(t, f.run("--json", "answer", "FID-aaaaaa", "Exclude authentication failures.", "--actor", "human:mark"), "task.answer")
			want := map[string]any{"question_id": parked["question_id"], "question": waitingQuestion,
				"text": answered["answer"], "decision_id": answered["decision_id"]}
			resumed, _ := lifecycleResult(t, f.run("--json", "next", "--claim", "--actor", "agent:claude-builder"), "next.claim")
			if resumed["kind"] != "resume" || !reflect.DeepEqual(resumed["answer"], want) {
				t.Fatalf("initial resume = %v, want answer %v", resumed, want)
			}
			if test.interruption == "expiry" {
				f.clock = func() time.Time { return statusNow.Add(time.Hour) }
			} else {
				hold := resumed["hold"].(map[string]any)["id"].(string)
				lifecycleResult(t, f.run("--json", "park", "--hold", hold, "--actor", "agent:claude-builder", "--interrupted", "--reason", "Usage limit", "--idempotency-key", "interrupt"), "hold.park")
			}
			if test.elsewhere {
				f.workspace = other
			}
			for _, claim := range []bool{false, true} {
				args, operation := []string{"--json", "next"}, "next.get"
				if claim {
					args, operation = append(args, "--claim", "--actor", "agent:claude-builder"), "next.claim"
				}
				next, _ := lifecycleResult(t, f.run(args...), operation)
				if next["kind"] != "new" || next["task"].(map[string]any)["id"] != "FID-aaaaaa" || next["waiting"] != nil {
					t.Fatalf("unstarted task scheduling changed: %v", next)
				}
				if claim && next["hold"].(map[string]any)["replaces"] != nil {
					t.Fatalf("unstarted task unexpectedly got a linked hold: %v", next)
				}
				if !reflect.DeepEqual(next["answer"], want) {
					t.Fatalf("answer disappeared without a new question: kind=%v answer=%v, want %v", next["kind"], next["answer"], want)
				}
			}
		})
	}
}
