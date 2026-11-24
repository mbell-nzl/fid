package main

import (
	"fmt"
	"math"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/mbell-nzl/fid/internal/doneeval"
	"github.com/mbell-nzl/fid/internal/gittest"
	"github.com/mbell-nzl/fid/internal/journal"
	"github.com/mbell-nzl/fid/internal/reviewwork"
	fidmigrations "github.com/mbell-nzl/fid/migrations"
)

func allowanceFixture(t *testing.T) (*statusFixture, string) {
	t.Helper()
	f, task := canonicalReviewFixtureWithClock(t, "standard", false, time.Now, func(f *statusFixture) {
		trustedShowPolicy(t, f, func(raw string) string {
			return strings.Replace(raw, "[lanes.standard]", "[lanes.standard]\nrework_cap=0", 1)
		})
	})
	lifecycleResult(t, recordReviewTask(t, f, changesInput("primary"), "changes", "human:courier", task), "review.record")
	return f, task
}

func TestReworkAllowanceServesBuilderWithGuidance(t *testing.T) {
	t.Parallel()
	f, task := allowanceFixture(t)
	stopped := commandResultFields(t, f, "next")
	if stopped["kind"] != "none" {
		t.Fatal(stopped)
	}
	item := selectionItem(t, f, "rework-cap")
	if item == nil || !strings.HasPrefix(item["next_action"].(string), "fid rework allow "+task+" --rounds <n>") {
		t.Fatalf("cap not actionable: %v", item)
	}
	for _, tc := range []struct{ actor, rounds, reason, code string }{{"agent:builder", "1", "Fix another round", "authority-denied"}, {"service:ci", "1", "Fix another round", "authority-denied"}, {"human:mark", "0", "Fix another round", "invalid-argument"}, {"human:mark", "1", " ", "invalid-argument"}} {
		before := startEventCount(t, f.workspace)
		got := f.run("--json", "rework", "allow", task, "--rounds", tc.rounds, "--actor", tc.actor, "--reason", tc.reason)
		assertJSONFailureCode(t, got.stderr, "task.rework-allow", ""+tc.code)
		if startEventCount(t, f.workspace) != before {
			t.Fatal("refusal changed journal")
		}
	}
	fields, _ := lifecycleResult(t, f.run("--json", "rework", "allow", task, "--rounds", "1", "--actor", "human:mark", "--reason", "Fix the uncovered case"), "task.rework-allow")
	if fields["rounds"] != float64(1) {
		t.Fatal(fields)
	}
	if selectionItem(t, f, "rework-cap") != nil {
		t.Fatal("cap item remains after allowance")
	}
	for _, rebuild := range []bool{false, true} {
		if rebuild {
			if _, err := journal.RebuildWorkspace(f.workspace, fidmigrations.FS); err != nil {
				t.Fatal(err)
			}
		}
		got := commandResultFields(t, f, "next")
		if got["kind"] != "rework" || got["task"].(map[string]any)["id"] != task || !strings.Contains(fmt.Sprint(got["rework_allowance"]), "Fix the uncovered case") {
			t.Fatalf("allowance did not return builder's work: %v", got)
		}
	}
	main := f.workspace
	f.workspace = gittest.AddWorktree(t, main)
	other := commandResultFields(t, f, "next")
	if other["kind"] != "none" {
		t.Fatal("another worktree took rework", other)
	}
	f.workspace = main
	lifecycleResult(t, f.run("--json", "rework", "allow", task, "--rounds", "2", "--actor", "human:mark", "--reason", "Cover both edge cases"), "task.rework-allow")
	conn, err := reviewwork.OpenProjection(main)
	if err != nil {
		t.Fatal(err)
	}
	defer closeCommandDB(t, conn)
	limit, err := doneeval.ReadReworkLimit(conn, task)
	if err != nil || limit.Count != 1 || limit.Cap != 3 {
		t.Fatalf("count was reset or allowances not summed: %v/%v", limit, err)
	}
	got := commandResultFields(t, f, "next", "--claim", "--actor", "agent:builder", "--model", "codex/gpt-6.1-sol")
	if got["kind"] != "rework" || !strings.Contains(fmt.Sprint(got["rework_allowance"]), "Cover both edge cases") {
		t.Fatal(got)
	}
	again := commandResultFields(t, f, "next")
	if again["kind"] != "resume" || !strings.Contains(fmt.Sprint(again["rework_allowance"]), "Cover both edge cases") {
		t.Fatal(again)
	}
	// A changed fingerprint and another hand-off retain all allowances and
	// earlier changes-requested rounds: the authority belongs to the task.
	writeBackupTestFile(t, filepath.Join(f.workspace, "notes.md"), "Another fix\n")
	runGitForTest(t, f.workspace, "add", ".")
	runGitForTest(t, f.workspace, "commit", "-m", "Another fix")
	reviewHandoffTask(t, f, "allowance-handoff", task)
	in := changesInput("primary")
	in["findings"] = []any{}
	lifecycleResult(t, recordReviewTask(t, f, in, "second-changes", "human:courier", task), "review.record")
	limit, err = doneeval.ReadReworkLimit(conn, task)
	if err != nil || limit.Count != 2 || limit.Cap != 3 {
		t.Fatalf("allowance lost across a hand-off: %v/%v", limit, err)
	}
}

func TestReworkAllowanceCannotOverflowEffectiveCap(t *testing.T) {
	t.Parallel()
	f, task := allowanceFixture(t)
	for _, reason := range []string{"First allowance", "Latest allowance"} {
		lifecycleResult(t, f.run("--json", "rework", "allow", task, "--rounds", "9223372036854775807", "--actor", "human:mark", "--reason", reason), "task.rework-allow")
	}
	conn, err := reviewwork.OpenProjection(f.workspace)
	if err != nil {
		t.Fatal(err)
	}
	defer closeCommandDB(t, conn)
	limit, err := doneeval.ReadReworkLimit(conn, task)
	if err != nil || limit.Cap != math.MaxInt64 || limit.Count != 1 || limit.Allowance.Reason != "Latest allowance" {
		t.Fatalf("overflowing allowances: %v/%v", limit, err)
	}
}
