package main

import (
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/mbell-nzl/fid/internal/gittest"
)

func checkHumanJourneyForKernelWords(t *testing.T) {
	t.Helper()
	base := t.TempDir()
	f := &statusFixture{t: t, workspace: gittest.NewWorkspace(t), home: filepath.Join(base, "home"),
		configHome: filepath.Join(base, "config"), clock: time.Now}
	call := func(exit int, args ...string) commandExecutionResult {
		t.Helper()
		got := f.run(args...)
		assertHumanCommand(t, got, exit)
		if exit == 0 && got.stdout == "" {
			t.Fatalf("%v produced no human output", args)
		}
		return got
	}
	call(0, "init", "--prefix", "FID")
	appendPolicy(t, f.workspace, "\n[lanes.standard]\nchecks=['lean']\n[review]\ntest_adequacy=[]\n"+leanHelperDefinition(t, "pass", "1m", nil))
	policyOrigin(t, f.workspace)
	created := call(0, "issue", "create", "--kind", "task", "--title", "Finish the text journey", "--lane", "standard", "--actor", "human:mark", "--idempotency-key", "create")
	task := parseCreatedIssueID(t, created.stdout)
	call(0, "next", "--claim", "--actor", "agent:builder")
	hold := humanJourneyHold(t, f, task)
	call(0, "start", "--hold", hold, "--actor", "agent:builder", "--idempotency-key", "start")
	call(0, "heartbeat", "--hold", hold, "--actor", "agent:builder")
	call(0, "check", "--hold", hold, "--actor", "agent:builder", "--idempotency-key", "check")
	call(0, "park", "--hold", hold, "--actor", "agent:builder", "--interrupted", "--reason", "Session stopped", "--idempotency-key", "interrupted")
	call(0, "status")
	call(0, "show", task)
	call(0, "next", "--claim", "--actor", "agent:builder")
	hold = humanJourneyHold(t, f, task)
	call(0, "park", "--hold", hold, "--actor", "agent:builder", "--waiting", "Include retries?", "--idempotency-key", "waiting")
	call(0, "status")
	call(0, "show", task)
	call(0, "next")
	call(0, "answer", task, "Include them.", "--actor", "human:mark")
	call(0, "next", "--claim", "--actor", "agent:builder")
	hold = humanJourneyHold(t, f, task)
	call(0, "takeover", task, "--actor", "human:mark", "--reason", "Finish this session")
	call(4, "heartbeat", "--hold", hold, "--actor", "agent:builder")
	hold = humanJourneyHold(t, f, task)
	call(0, "check", "--hold", hold, "--actor", "human:mark", "--idempotency-key", "person-check")
	call(0, "handoff", "--hold", hold, "--actor", "human:mark", "--pr", "https://github.com/mbell-nzl/fid/pull/63", "--idempotency-key", "handoff")
	call(0, "next", "--claim", "--actor", "agent:builder")
	call(0, "status")
	call(0, "show", task)
}

// The text deliberately hides hold IDs. A JSON read supplies only the
// session handle; every action and every output checked above is human mode.
func humanJourneyHold(t *testing.T, f *statusFixture, task string) string {
	t.Helper()
	fields, _ := lifecycleResult(t, f.run("--json", "show", task), "issue.show")
	working, ok := fields["working"].(map[string]any)
	if !ok {
		t.Fatalf("show has no working session: %v", fields)
	}
	hold, ok := working["hold_id"].(string)
	if !ok || !strings.HasPrefix(hold, "HLD-") {
		t.Fatalf("show has no hold ID: %v", working)
	}
	return hold
}

func TestV2HumanRefusalsNameNoKernelWords(t *testing.T) {
	t.Parallel()
	for _, row := range []struct {
		name string
		args []string
		exit int
	}{
		{"taking without actor", []string{"next", "--claim"}, 2},
		{"actor without taking", []string{"next", "--actor", "agent:builder"}, 2},
		{"service taking", []string{"next", "--claim", "--actor", "service:ci"}, 5},
		{"missing hold", []string{"start", "--hold", "HLD-abcdef", "--actor", "agent:builder", "--idempotency-key", "start"}, 4},
		{"missing issue", []string{"show", "FID-abcdef"}, 3},
		{"agent takeover", []string{"takeover", "FID-abcdef", "--actor", "agent:builder", "--reason", "Continue"}, 5},
		{"agent answer", []string{"answer", "FID-abcdef", "Continue", "--actor", "agent:builder"}, 5},
	} {
		t.Run(row.name, func(t *testing.T) {
			t.Parallel()
			f := newJournalStatusFixture(t)
			f.clock = time.Now
			assertHumanCommand(t, f.run(row.args...), row.exit)
		})
	}
}
