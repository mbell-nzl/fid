package main

import (
	"encoding/json"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/mbell-nzl/fid/internal/config"
	fiddb "github.com/mbell-nzl/fid/internal/db"
)

// READY keeps status's order even when next serves unblocking work below
// the first row, or the worktree already holds a task.
func TestStatusReadyAgreesWithNext(t *testing.T) {
	t.Parallel()
	for _, test := range []struct {
		name, task string
		ago        time.Duration
		wantNext   any
	}{
		{"unblocking below first row", "", 0, "FID-77aa01"},
		{"held here", "FID-9f13aa", time.Minute, nil},
		{"interrupted here", "FID-71a0c2", 2 * time.Hour, nil},
		{"lapsed unstarted here", "FID-9f13aa", 2 * time.Hour, "FID-77aa01"},
	} {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			fixture := newStatusFixture(t, time.Hour)
			nextWork(fixture)
			fixture.seed(nil, [][3]string{{"FID-77aa01", "blocks", "FID-c2d4e5"}})
			if test.task != "" {
				fixture.hold("HLD-0a0a0a", test.task, "", test.ago)
			}
			status := commandResultFields(t, fixture, "status")
			ready, ok := status["ready"].(map[string]any)
			if !ok || ready["next"] != test.wantNext {
				t.Fatalf("READY = %#v, want next %v", status["ready"], test.wantNext)
			}
			if _, old := status["to_do"]; old {
				t.Fatal("status still has to_do")
			}
			next := commandResultFields(t, fixture, "next")
			if test.wantNext != nil && next["task"].(map[string]any)["id"] != ready["next"] {
				t.Fatalf("status next %v disagrees with next's task %v", ready["next"], next["task"])
			}
			if test.wantNext == nil && (next["kind"] != "resume" || next["task"].(map[string]any)["id"] != test.task) {
				t.Fatalf("next doesn't return this worktree's task: %v", next)
			}
			assertReadyMatchesQueue(t, ready, next)
			text := fixture.run("status")
			if text.err != nil || !strings.Contains(text.stdout, "  READY\n") || test.wantNext == nil && strings.Contains(text.stdout, "  next\n") {
				t.Fatalf("status = %v:\n%s%s", text.err, text.stdout, text.stderr)
			}
			if test.task != "" && test.ago < time.Hour && !strings.Contains(text.stdout, "working here") {
				t.Fatalf("status doesn't show working here:\n%s", text.stdout)
			}
			if test.wantNext != nil {
				taken := commandResultFields(t, fixture, "next", "--claim", "--actor", "agent:codex")
				if taken["task"].(map[string]any)["id"] != ready["next"] {
					t.Fatalf("READY next %v disagrees with the task handed out: %v", ready["next"], taken["task"])
				}
			}
		})
	}
}

func assertReadyMatchesQueue(t *testing.T, ready, next map[string]any) {
	t.Helper()
	queue := map[any]any{}
	for _, item := range next["queue"].([]any) {
		queue[item.(map[string]any)["id"]] = item
	}
	tasks := append([]any{}, ready["tasks"].([]any)...)
	for _, epic := range ready["epics"].([]any) {
		tasks = append(tasks, epic.(map[string]any)["tasks"].([]any)...)
	}
	for _, item := range tasks {
		if task := item.(map[string]any); !reflect.DeepEqual(task, queue[task["id"]]) {
			t.Fatalf("READY's task differs from next's queue entry:\n%v\n%v", task, queue[task["id"]])
		}
	}
}

func TestStatusReadyWhenNothingCanBeServed(t *testing.T) {
	t.Parallel()
	fixture := newStatusFixture(t, time.Hour)
	nextWork(fixture)
	for index, task := range []string{"FID-9f13aa", "FID-0d4e8b", "FID-77aa01"} {
		fixture.hold([]string{"HLD-0a0a0a", "HLD-0b0b0b", "HLD-0c0c0c"}[index], task, "/tmp/other-"+task, time.Minute)
	}
	ready := commandResultFields(t, fixture, "status")["ready"].(map[string]any)
	next := commandResultFields(t, fixture, "next")
	if ready["next"] != nil || next["kind"] != "none" {
		t.Fatalf("READY next = %v, next kind = %v", ready["next"], next["kind"])
	}
	assertReadyMatchesQueue(t, ready, next)
}

func commandResultFields(t *testing.T, fixture *statusFixture, args ...string) map[string]any {
	t.Helper()
	result := fixture.run(append([]string{"--json"}, args...)...)
	var reply struct{ Result map[string]any }
	if result.err != nil || json.Unmarshal([]byte(result.stdout), &reply) != nil {
		t.Fatalf("%v = %v:\n%s%s", args, result.err, result.stdout, result.stderr)
	}
	return reply.Result
}

func (fixture *statusFixture) transitionSeededTask(id, state string) {
	fixture.t.Helper()
	conn, err := fiddb.Open(filepath.Join(config.LedgerDir(fixture.workspace), config.DatabaseFileName))
	if err != nil {
		fixture.t.Fatal(err)
	}
	defer closeCommandDB(fixture.t, conn)
	if _, err := conn.Exec(`UPDATE issues SET status = ? WHERE id = ?`, state, id); err != nil {
		fixture.t.Fatal(err)
	}
}

// Each verdict is covered in text and JSON; existing status goldens also
// keep the cap, grouping, backup line and empty-project behavior pinned.
func TestStatusReadyGolden(t *testing.T) {
	t.Parallel()
	for _, name := range []string{"verdicts", "working-here", "interrupted-here"} {
		for _, format := range []string{"txt", "json"} {
			t.Run(name+"."+format, func(t *testing.T) {
				t.Parallel()
				fixture := newStatusFixture(t, time.Hour)
				fixture.setTarget()
				nextWork(fixture)
				fixture.seed([]readyCommandRow{
					statusTask("FID-999999", "Another blocker", "in_progress", nil, nil, nil),
					statusEpic("FID-eeeeee", "Waiting epic", "todo", intPointer(3)),
					statusTask("FID-aaaaaa", "Waiting on the epic", "todo", stringPointer("FID-eeeeee"), nil, nil),
					statusEpic("FID-ffffff", "Closed epic", "done", intPointer(4)),
					statusTask("FID-bbbbbb", "Left in a closed epic", "todo", stringPointer("FID-ffffff"), nil, nil),
				}, [][3]string{{"FID-eeeeee", "blocked_by", "FID-71a0c2"}, {"FID-eeeeee", "blocked_by", "FID-c2d4e5"}, {"FID-c2d4e5", "blocked_by", "FID-999999"}})
				fixture.hold("HLD-0a0a0a", "FID-0d4e8b", "/home/mark/projects/fid-wt3", time.Minute)
				fixture.hold("HLD-0b0b0b", "FID-71a0c2", "/home/mark/projects/fid-wt4", 2*time.Hour)
				if name == "working-here" {
					fixture.hold("HLD-0c0c0c", "FID-9f13aa", "", time.Minute)
				}
				if name == "interrupted-here" {
					fixture.hold("HLD-0c0c0c", "FID-77aa01", "", 2*time.Hour)
					// A started task's lapsed hold keeps it for this worktree.
					fixture.transitionSeededTask("FID-77aa01", "in_progress")
				}
				args := []string{"status"}
				if format == "json" {
					args = append([]string{"--json"}, args...)
				}
				fixture.golden(filepath.Join("status", "ready-"+name+"."+format), args...)
			})
		}
	}
}
