package main

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/mbell-nzl/fid/internal/cli"
	"github.com/mbell-nzl/fid/internal/gittest"
)

const waitingQuestion = "Should retries include authentication failures?"

func parkWaiting(t *testing.T, f *statusFixture, hold, question, key string) map[string]any {
	t.Helper()
	result, _ := lifecycleResult(t, f.run("--json", "park", "--hold", hold, "--actor", "agent:claude-builder", "--waiting", question, "--idempotency-key", key), "hold.park")
	return result
}

func TestWaitingAnswerResume(t *testing.T) {
	t.Parallel()
	for _, state := range []string{"todo", "in_progress"} {
		t.Run(state, func(t *testing.T) {
			f := newLifecycleFixture(t)
			f.seed([]readyCommandRow{statusTask("FID-aaaaaa", "Waiting work", state, nil, nil, nil), statusTask("FID-bbbbbb", "Other work", "todo", nil, nil, nil)}, nil)
			other := gittest.AddWorktree(t, f.workspace)
			main := f.workspace
			f.hold("HLD-aaaaaa", "FID-aaaaaa", "", time.Minute)
			parked := parkWaiting(t, f, "HLD-aaaaaa", waitingQuestion, "wait-1")
			question := parked["question_id"]
			if question == nil || parked["kept_for_worktree"] != true {
				t.Fatalf("waiting park = %v", parked)
			}
			blocked, event := lifecycleResult(t, f.run("--json", "next", "--claim", "--actor", "agent:builder"), "next.claim")
			if blocked["kind"] != "none" || blocked["task"] != nil || event != nil || blocked["waiting"].(map[string]any)["question_id"] != question {
				t.Fatalf("waiting next = %v, %v", blocked, event)
			}
			f.workspace = other
			available, _ := lifecycleResult(t, f.run("--json", "next"), "next.get")
			if available["task"].(map[string]any)["id"] != "FID-bbbbbb" || available["queue"].([]any)[0].(map[string]any)["reason"] != "waiting-for-you" {
				t.Fatalf("other worktree = %v", available)
			}
			args := []string{"--json", "answer", "FID-aaaaaa", "Exclude authentication failures.", "--actor", "human:mark", "--idempotency-key", "answer-1"}
			first := f.run(args...)
			answered, _ := lifecycleResult(t, first, "task.answer")
			if answered["question_id"] != question || answered["answer"] != "Exclude authentication failures." {
				t.Fatalf("answer = %v", answered)
			}
			if repeat := f.run(args...); repeat.err != nil || repeat.stdout != first.stdout {
				t.Fatalf("answer retry = %#v", repeat)
			}
			f.workspace = main
			resumed, _ := lifecycleResult(t, f.run("--json", "next", "--claim", "--actor", "agent:claude-builder"), "next.claim")
			hold := resumed["hold"].(map[string]any)
			answer := resumed["answer"].(map[string]any)
			if resumed["kind"] != "resume" || hold["replaces"] != "HLD-aaaaaa" || answer["decision_id"] != answered["decision_id"] || answer["question_id"] != question {
				t.Fatalf("resume = %v", resumed)
			}
			repeated, _ := lifecycleResult(t, f.run("--json", "next", "--claim", "--actor", "agent:claude-builder"), "next.claim")
			if fmt.Sprint(repeated["answer"]) != fmt.Sprint(resumed["answer"]) {
				t.Fatalf("lost answer = %v", repeated)
			}
			parkWaiting(t, f, hold["id"].(string), "Include timeouts?", "wait-2")
			second, _ := lifecycleResult(t, f.run("--json", "next"), "next.get")
			if second["answer"] != nil || second["waiting"].(map[string]any)["question_id"] == question {
				t.Fatalf("new question = %v", second)
			}
			shown, _ := lifecycleResult(t, f.run("--json", "show", "FID-aaaaaa"), "issue.show")
			if len(shown["questions"].([]any)) != 2 {
				t.Fatalf("question history = %v", shown)
			}
		})
	}
}

func TestWaitingTakeoverAndCancel(t *testing.T) {
	t.Parallel()
	for _, action := range []string{"answer", "cancel", "cancel-without-takeover"} {
		t.Run(action, func(t *testing.T) {
			f := newLifecycleFixture(t)
			nextWork(f)
			f.hold("HLD-aaaaaa", "FID-71a0c2", "", time.Minute)
			parked := parkWaiting(t, f, "HLD-aaaaaa", waitingQuestion, "wait")
			if action != "cancel-without-takeover" {
				f.workspace = gittest.AddWorktree(t, f.workspace)
				lifecycleResult(t, f.run("--json", "takeover", "FID-71a0c2", "--actor", "human:mark", "--reason", "Continue here"), "task.takeover")
				shown, _ := lifecycleResult(t, f.run("--json", "show", "FID-71a0c2"), "issue.show")
				if shown["questions"].([]any)[0].(map[string]any)["resolved_at"] != nil || shown["working"] == nil {
					t.Fatalf("takeover lost question = %v", shown)
				}
			}
			var shown map[string]any
			if action == "answer" {
				answered, _ := lifecycleResult(t, f.run("--json", "answer", "FID-71a0c2", "Yes", "--actor", "human:mark"), "task.answer")
				if answered["question_id"] != parked["question_id"] {
					t.Fatalf("wrong question = %v", answered)
				}
			} else {
				lifecycleResult(t, f.run("--json", "issue", "transition", "FID-71a0c2", "--to", "cancelled", "--reason", "No longer needed", "--actor", "human:mark"), "issue.transition")
				shown, _ = lifecycleResult(t, f.run("--json", "show", "FID-71a0c2"), "issue.show")
				if shown["questions"].([]any)[0].(map[string]any)["resolution"] != "cancelled" || shown["kept"] != nil {
					t.Fatalf("cancel left waiting = %v", shown)
				}
				refused := f.run("--json", "answer", "FID-71a0c2", "Late", "--actor", "human:mark")
				assertJSONFailureCode(t, refused.stderr, "task.answer", "not-waiting")
			}
		})
	}
}

func TestWaitingRefusals(t *testing.T) {
	t.Parallel()
	for _, test := range []struct {
		name, actor, code string
		waiting           bool
		exit              int
	}{
		{"not-waiting", "human:mark", "not-waiting", false, 4},
		{"agent-answer", "agent:builder", "authority-denied", true, 5},
		{"service-answer", "service:ci", "authority-denied", true, 5},
	} {
		t.Run(test.name, func(t *testing.T) {
			f := newLifecycleFixture(t)
			nextWork(f)
			if test.waiting {
				f.hold("HLD-aaaaaa", "FID-71a0c2", "", time.Minute)
				parkWaiting(t, f, "HLD-aaaaaa", waitingQuestion, "wait")
			}
			for _, jsonMode := range []bool{false, true} {
				args := []string{"answer", "FID-71a0c2", "Yes", "--actor", test.actor}
				ext := ".txt"
				if jsonMode {
					args = append([]string{"--json"}, args...)
					ext = ".json"
				}
				result := f.run(args...)
				if cli.ExitCode(result.err) != test.exit || result.stdout != "" {
					t.Fatalf("refusal = %#v", result)
				}
				if jsonMode {
					assertJSONFailureCode(t, result.stderr, "task.answer", test.code)
				}
				waitingGolden(t, f, test.name+ext, result.stderr)
				if strings.Contains(result.stderr, "--question") {
					t.Fatal("question flag leaked")
				}
			}
		})
	}
}

func waitingGolden(t *testing.T, f *statusFixture, name, output string) {
	t.Helper()
	output = regexp.MustCompile(`QST-[0-9a-f]{6}`).ReplaceAllString(output, "QST-<new>")
	lifecycleGolden(t, f, filepath.Join("..", "waiting", name), output)
}

func TestWaitingGolden(t *testing.T) {
	t.Parallel()
	f := newLifecycleFixture(t)
	nextWork(f)
	f.hold("HLD-aaaaaa", "FID-71a0c2", "", time.Minute)
	args := []string{"park", "--hold", "HLD-aaaaaa", "--actor", "agent:claude-builder", "--waiting", waitingQuestion, "--idempotency-key", "wait"}
	parked := f.run(append([]string{"--json"}, args...)...)
	lifecycleResult(t, parked, "hold.park")
	waitingGolden(t, f, "park.json", parked.stdout)
	waitingGolden(t, f, "park.txt", f.run(args...).stdout)
	for _, command := range []string{"next", "status", "show"} {
		args := []string{command}
		if command == "show" {
			args = append(args, "FID-71a0c2")
		}
		for _, jsonMode := range []bool{false, true} {
			ext := ".txt"
			if jsonMode {
				args = append([]string{"--json"}, args...)
				ext = ".json"
			}
			result := f.run(args...)
			if result.err != nil {
				t.Fatalf("%s = %#v", command, result)
			}
			if command == "status" && !strings.Contains(result.stdout, "waiting for you") {
				t.Fatalf("status lost waiting question: %s", result.stdout)
			}
			waitingGolden(t, f, command+ext, result.stdout)
		}
	}
	answerArgs := []string{"answer", "FID-71a0c2", "Exclude authentication failures.", "--actor", "human:mark", "--idempotency-key", "answer"}
	answer := f.run(append([]string{"--json"}, answerArgs...)...)
	lifecycleResult(t, answer, "task.answer")
	waitingGolden(t, f, "answer.json", answer.stdout)
	waitingGolden(t, f, "answer.txt", f.run(answerArgs...).stdout)
	for _, jsonMode := range []bool{false, true} {
		args := []string{"next"}
		ext := ".txt"
		if jsonMode {
			args = append([]string{"--json"}, args...)
			ext = ".json"
		}
		waitingGolden(t, f, "resume"+ext, f.run(args...).stdout)
	}
}

func TestAnswerRetryKeepsOriginalWorktree(t *testing.T) {
	t.Parallel()
	f := newLifecycleFixture(t)
	nextWork(f)
	f.hold("HLD-aaaaaa", "FID-71a0c2", "", time.Minute)
	parkWaiting(t, f, "HLD-aaaaaa", waitingQuestion, "wait")
	args := []string{"answer", "FID-71a0c2", "Yes", "--actor", "human:mark", "--idempotency-key", "answer"}
	first := f.run(args...)
	if first.err != nil {
		t.Fatal(first.err)
	}
	lifecycleResult(t, f.run("--json", "issue", "transition", "FID-71a0c2", "--to", "cancelled", "--actor", "human:mark", "--reason", "No longer needed"), "issue.transition")
	if replay := f.run(args...); replay.err != nil || replay.stdout != first.stdout {
		t.Fatalf("answer retry = %#v; original = %#v", replay, first)
	}
}
