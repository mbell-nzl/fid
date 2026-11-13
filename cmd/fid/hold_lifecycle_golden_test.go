package main

import (
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/mbell-nzl/fid/internal/cli"
)

var lifecycleIDs = regexp.MustCompile(`\b(?:EVT-\d{8}-\d+|PDN-[0-9a-f]{6}|HLD-[0-9a-f]{6})\b`)

func lifecycleGolden(t *testing.T, f *statusFixture, name, output string) {
	t.Helper()
	if strings.HasPrefix(name, "cancel.") {
		output = regexp.MustCompile(`"occurred_at":"[^"]+"`).ReplaceAllString(output, `"occurred_at":"<TIME>"`)
	}
	output = lifecycleIDs.ReplaceAllStringFunc(f.placeholders(output), func(id string) string {
		if id == "HLD-aaaaaa" || id == "HLD-bbbbbb" {
			return id
		}
		return strings.Split(id, "-")[0] + "-<new>"
	})
	assertGolden(t, filepath.Join("hold-lifecycle", name), output)
}

func TestHoldLifecycleGolden(t *testing.T) {
	t.Parallel()
	for _, test := range []struct {
		name      string
		setup     func(*statusFixture)
		args      []string
		operation string
	}{
		{"heartbeat", func(f *statusFixture) { f.hold("HLD-aaaaaa", "FID-71a0c2", "", 30*time.Minute) }, []string{"heartbeat", "--hold", "HLD-aaaaaa", "--actor", "agent:claude-builder"}, "hold.heartbeat"},
		{"park-started", func(f *statusFixture) { f.hold("HLD-aaaaaa", "FID-71a0c2", "", 30*time.Minute) }, []string{"park", "--hold", "HLD-aaaaaa", "--actor", "agent:claude-builder", "--interrupted", "--reason", "Usage limit", "--idempotency-key", "park-1"}, "hold.park"},
		{"park-unstarted", func(f *statusFixture) { f.hold("HLD-aaaaaa", "FID-9f13aa", "", 30*time.Minute) }, []string{"park", "--hold", "HLD-aaaaaa", "--actor", "agent:claude-builder", "--interrupted", "--reason", "Usage limit", "--idempotency-key", "park-1"}, "hold.park"},
		{"takeover", func(f *statusFixture) { f.hold("HLD-aaaaaa", "FID-71a0c2", "/w/old", 30*time.Minute) }, []string{"takeover", "FID-71a0c2", "--actor", "human:mark", "--reason", "Agent stuck on a flaky test", "--idempotency-key", "take-1"}, "task.takeover"},
		{"cancel", func(f *statusFixture) { f.hold("HLD-aaaaaa", "FID-71a0c2", "", 30*time.Minute) }, []string{"issue", "transition", "FID-71a0c2", "--to", "cancelled", "--actor", "human:mark", "--reason", "No longer needed", "--idempotency-key", "cancel-1"}, "issue.transition"},
	} {
		t.Run(test.name, func(t *testing.T) {
			f := newLifecycleFixture(t)
			nextWork(f)
			test.setup(f)
			answer := f.run(append([]string{"--json"}, test.args...)...)
			lifecycleResult(t, answer, test.operation)
			lifecycleGolden(t, f, test.name+".json", answer.stdout)
			human := f.run(test.args...)
			if human.err != nil || human.stderr != "" {
				t.Fatalf("text = %#v", human)
			}
			lifecycleGolden(t, f, test.name+".txt", human.stdout)
			if test.name == "park-started" {
				for _, command := range []struct {
					name string
					args []string
				}{{"show", []string{"show", "FID-71a0c2"}}, {"next", []string{"next"}}, {"status", []string{"status", "--all"}}} {
					f.golden(filepath.Join("hold-lifecycle", "parked-"+command.name+".txt"), command.args...)
					f.golden(filepath.Join("hold-lifecycle", "parked-"+command.name+".json"), append([]string{"--json"}, command.args...)...)
				}
			}
		})
	}
}

func TestHoldLifecycleRefusalGolden(t *testing.T) {
	t.Parallel()
	for _, test := range []struct {
		name, code string
		exit       int
		setup      func(*statusFixture)
		args       []string
	}{
		{"missing", "hold-not-current", 4, nil, []string{"heartbeat", "--hold", "HLD-aaaaaa", "--actor", "agent:claude-builder"}},
		{"expired", "hold-not-current", 4, func(f *statusFixture) { f.hold("HLD-aaaaaa", "FID-71a0c2", "", time.Hour) }, []string{"heartbeat", "--hold", "HLD-aaaaaa", "--actor", "agent:claude-builder"}},
		{"actor-mismatch", "hold-not-current", 4, func(f *statusFixture) { f.hold("HLD-aaaaaa", "FID-71a0c2", "", time.Minute) }, []string{"heartbeat", "--hold", "HLD-aaaaaa", "--actor", "agent:other"}},
		{"worktree-mismatch", "hold-not-current", 4, func(f *statusFixture) { f.hold("HLD-aaaaaa", "FID-71a0c2", "/w/old", time.Minute) }, []string{"heartbeat", "--hold", "HLD-aaaaaa", "--actor", "agent:claude-builder"}},
		{"ended", "hold-not-current", 4, func(f *statusFixture) {
			f.hold("HLD-aaaaaa", "FID-71a0c2", "", time.Minute)
			lifecycleResult(f.t, f.run("--json", "takeover", "FID-71a0c2", "--actor", "human:mark", "--reason", "Agent stuck"), "task.takeover")
		}, []string{"heartbeat", "--hold", "HLD-aaaaaa", "--actor", "agent:claude-builder"}},
		{"busy", "worktree-busy", 4, func(f *statusFixture) { f.hold("HLD-bbbbbb", "FID-c2d4e5", "", time.Minute) }, []string{"takeover", "FID-71a0c2", "--actor", "human:mark", "--reason", "Agent stuck"}},
		{"already-held", "already-held", 4, func(f *statusFixture) {
			lifecycleResult(f.t, f.run("--json", "takeover", "FID-71a0c2", "--actor", "human:mark", "--reason", "Agent stuck"), "task.takeover")
		}, []string{"takeover", "FID-71a0c2", "--actor", "human:mark", "--reason", "Agent stuck"}},
		{"closed", "transition-not-permitted", 4, nil, []string{"takeover", "FID-77aa03", "--actor", "human:mark", "--reason", "Agent stuck"}},
		{"epic", "unsupported-kind", 2, nil, []string{"takeover", "FID-e3a1c0", "--actor", "human:mark", "--reason", "Agent stuck"}},
		{"agent-takeover", "authority-denied", 5, nil, []string{"takeover", "FID-71a0c2", "--actor", "agent:builder", "--reason", "Agent stuck"}},
		{"service-takeover", "authority-denied", 5, nil, []string{"takeover", "FID-71a0c2", "--actor", "service:ci", "--reason", "Agent stuck"}},
		{"agent-cancel", "authority-denied", 5, nil, []string{"issue", "transition", "FID-71a0c2", "--to", "cancelled", "--actor", "agent:builder", "--reason", "No longer needed", "--idempotency-key", "cancel"}},
		{"bare-park", "invalid-argument", 2, nil, []string{"park"}},
		{"park-reason-required", "invalid-argument", 2, nil, []string{"park", "--hold", "HLD-aaaaaa", "--actor", "agent:builder", "--interrupted", "--idempotency-key", "park"}},
		{"park-key-required", "invalid-argument", 2, nil, []string{"park", "--hold", "HLD-aaaaaa", "--actor", "agent:builder", "--interrupted", "--reason", "Usage limit"}},
		{"cancel-reason-required", "invalid-argument", 2, nil, []string{"issue", "transition", "FID-71a0c2", "--to", "cancelled", "--actor", "human:mark"}},
	} {
		t.Run(test.name, func(t *testing.T) {
			f := newLifecycleFixture(t)
			nextWork(f)
			if test.setup != nil {
				test.setup(f)
			}
			for _, jsonMode := range []bool{false, true} {
				args := append([]string{}, test.args...)
				ext := ".txt"
				if jsonMode {
					args = append([]string{"--json"}, args...)
					ext = ".json"
				}
				answer := f.run(args...)
				if cli.ExitCode(answer.err) != test.exit || answer.stdout != "" {
					t.Fatalf("refusal = %#v, want exit %d", answer, test.exit)
				}
				if jsonMode {
					assertJSONFailureCode(t, answer.stderr, leafForLifecycle(test.args), test.code)
				}
				lifecycleGolden(t, f, test.name+ext, answer.stderr)
			}
		})
	}
}

func leafForLifecycle(args []string) string {
	switch args[0] {
	case "takeover":
		return "task.takeover"
	case "park":
		return "hold.park"
	case "issue":
		return "issue.transition"
	default:
		return "hold.heartbeat"
	}
}
