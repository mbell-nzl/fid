//go:build linux

package main

import (
	"strings"
	"testing"
	"time"
)

func TestWaiveRealTerminalBoundaries(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name                string
		input, prompt       bool
		actor, answer, want string
	}{
		{"piped input", false, true, "human:mark", "", "standard input and standard error must be terminals"},
		{"piped prompt", true, false, "human:mark", "", "standard input and standard error must be terminals"},
		{"agent", true, true, "agent:builder", "", "only a person"},
		{"service", true, true, "service:ci", "", "only a person"},
		{"cancel", true, true, "human:mark", "no\n", "Waiver cancelled"},
		{"confirmed", true, true, "human:mark", "TASK", "Type "},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			f, task := canonicalReviewFixtureWithClock(t, "assured", false, time.Now)
			before := reviewProjection(t, f)
			answer := tc.answer
			if answer == "TASK" {
				answer = task + "\n"
			}
			got := personPublicForTest(t, f.workspace, task, answer, tc.input, tc.prompt, true, "waive", "--lane", "standard", "--actor", tc.actor, "--idempotency-key", "terminal")
			if !strings.Contains(got.stderr, tc.want) {
				t.Fatal(got)
			}
			if tc.name == "confirmed" {
				waiveSuccess(t, got)
			} else if got.stdout != "" || reviewProjection(t, f) != before {
				t.Fatal("refusal changed ledger", got)
			}
		})
	}
}
