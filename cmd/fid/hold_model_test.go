package main

import (
	"strings"
	"testing"
	"time"
)

func TestBuildHoldModelAdmission(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name, actor, model string
		want               bool
	}{
		{"agent missing", "agent:builder", "", false},
		{"agent model", "agent:builder", "codex/gpt-6.1-sol", true},
		{"human by hand", "human:mark", "", true},
		{"human model", "human:mark", "claude/opus-5.5", true},
		{"malformed", "agent:builder", "codex", false},
		{"extra slash", "agent:builder", "codex/a/b", false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			f := newStatusFixture(t, 100*time.Hour)
			f.seed([]readyCommandRow{statusTask("FID-aaaaaa", "Build", "todo", nil, nil, nil)}, nil)
			args := []string{"--json", "next", "--claim", "--actor", tc.actor}
			if tc.model != "" {
				args = append(args, "--model", tc.model)
			}
			got := f.run(args...)
			if !tc.want {
				if got.err == nil {
					t.Fatalf("accepted invalid model: %#v", got)
				}
				return
			}
			fields, _ := decodeNext(t, got, "next.claim")
			if fields["hold"].(map[string]any)["model"] != tc.model {
				t.Fatalf("model not recorded: %v", fields)
			}
		})
	}
}

func TestBuildModelFlagsAreNeverIgnored(t *testing.T) {
	t.Parallel()
	for _, args := range [][]string{{"next", "--model", "codex/gpt-6.1-sol"}, {"takeover", "FID-aaaaaa", "--model", "malformed", "--actor", "human:mark", "--reason", "Continue"}} {
		t.Run(args[0], func(t *testing.T) {
			t.Parallel()
			f := newStatusFixture(t, 100*time.Hour)
			before := startEventCount(t, f.workspace)
			got := f.run(append([]string{"--json"}, args...)...)
			if got.err == nil || startEventCount(t, f.workspace) != before {
				t.Fatalf("invalid model flags=%#v", got)
			}
		})
	}
}

func TestTakeoverModelParticipatesInRetryAndReads(t *testing.T) {
	t.Parallel()
	f := newStatusFixture(t, 100*time.Hour)
	f.seed([]readyCommandRow{statusTask("FID-aaaaaa", "Build", "todo", nil, nil, nil)}, nil)
	args := []string{"--json", "takeover", "FID-aaaaaa", "--actor", "human:mark", "--reason", "Build with a model", "--idempotency-key", "model-takeover", "--model", "codex/gpt-6.1-sol"}
	first := f.run(args...)
	lifecycleResult(t, first, "task.takeover")
	if again := f.run(args...); again.err != nil || again.stdout != first.stdout {
		t.Fatalf("retry=%#v", again)
	}
	shown, _ := lifecycleResult(t, f.run("--json", "show", "FID-aaaaaa"), "issue.show")
	if shown["working"].(map[string]any)["model"] != "codex/gpt-6.1-sol" {
		t.Fatal(shown)
	}
	before := startEventCount(t, f.workspace)
	args[len(args)-1] = "claude/opus-5.5"
	conflict := f.run(args...)
	assertJSONFailureCode(t, conflict.stderr, "task.takeover", "idempotency-conflict")
	if !strings.Contains(conflict.stderr, `"differing_fields":["--model"]`) || startEventCount(t, f.workspace) != before {
		t.Fatalf("conflict=%#v", conflict)
	}
}
