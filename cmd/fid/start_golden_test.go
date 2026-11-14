package main

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/mbell-nzl/fid/internal/cli"
	"github.com/mbell-nzl/fid/internal/config"
	fiddb "github.com/mbell-nzl/fid/internal/db"
)

func plannedStartFixture(t *testing.T, state string, lane *string) *statusFixture {
	t.Helper()
	f := newLifecycleFixture(t)
	f.seed([]readyCommandRow{statusEpic("FID-e3a1c0", "Planned epic", "todo", nil), statusTask("FID-9f13aa", "Harden login rate limiting", state, stringPointer("FID-e3a1c0"), nil, nil)}, nil)
	conn, err := fiddb.Open(filepath.Join(config.LedgerDir(f.workspace), config.DatabaseFileName))
	if err != nil {
		t.Fatal(err)
	}
	defer closeCommandDB(t, conn)
	if _, err := conn.Exec(`UPDATE issues SET lane=? WHERE id='FID-9f13aa'`, lane); err != nil {
		t.Fatal(err)
	}
	f.hold("HLD-aaaaaa", "FID-9f13aa", "", time.Minute)
	return f
}

func TestStartGolden(t *testing.T) {
	t.Parallel()
	for _, repeat := range []bool{false, true} {
		t.Run(map[bool]string{false: "start", true: "already-started"}[repeat], func(t *testing.T) {
			f := plannedStartFixture(t, "todo", stringPointer("standard"))
			args := []string{"start", "--hold", "HLD-aaaaaa", "--actor", "agent:claude-builder", "--idempotency-key", "start-1"}
			if repeat {
				lifecycleResult(t, f.run(append([]string{"--json"}, args...)...), "task.start")
				args[len(args)-1] = "start-2"
			}
			for _, jsonMode := range []bool{false, true} {
				ext := ".txt"
				invoke := args
				if jsonMode {
					ext = ".json"
					invoke = append([]string{"--json"}, args...)
				}
				// Independent fixtures pin both the first-start text and JSON.
				current := f
				if !repeat && jsonMode {
					current = plannedStartFixture(t, "todo", stringPointer("standard"))
				}
				result := current.run(invoke...)
				if result.err != nil {
					t.Fatalf("start = %#v", result)
				}
				lifecycleGolden(t, current, filepath.Join("..", "start", map[bool]string{false: "start", true: "already-started"}[repeat]+ext), result.stdout)
			}
		})
	}
}

func TestStartRefusalGolden(t *testing.T) {
	t.Parallel()
	for _, row := range []struct {
		name, code string
		exit       int
		setup      func(*statusFixture)
		args       []string
	}{
		{"lane-not-set", "lane-not-set", 4, func(f *statusFixture) { startSQL(f.t, f, `UPDATE issues SET lane=NULL WHERE id='FID-9f13aa'`) }, []string{"start", "--hold", "HLD-aaaaaa", "--actor", "agent:claude-builder", "--idempotency-key", "start"}},
		{"in-review", "transition-not-permitted", 4, func(f *statusFixture) { startSQL(f.t, f, `UPDATE issues SET status='in_review' WHERE id='FID-9f13aa'`) }, []string{"start", "--hold", "HLD-aaaaaa", "--actor", "agent:claude-builder", "--idempotency-key", "start"}},
		{"epic-lane-create", "unsupported-kind", 2, nil, []string{"issue", "create", "--kind", "epic", "--title", "Epics have no lane", "--scheduling", "prioritised", "--lane", "standard", "--actor", "human:mark", "--idempotency-key", "lane"}},
		{"epic-lane-update", "unsupported-kind", 2, nil, []string{"issue", "update", "FID-e3a1c0", "--lane", "standard", "--actor", "human:mark", "--idempotency-key", "lane"}},
		{"agent-lowering", "authority-denied", 5, nil, []string{"issue", "update", "FID-9f13aa", "--lane", "routine", "--actor", "agent:claude-builder", "--idempotency-key", "lower"}},
		{"started-lowering", "transition-not-permitted", 4, func(f *statusFixture) {
			startSQL(f.t, f, `UPDATE issues SET status='in_progress' WHERE id='FID-9f13aa'`)
		}, []string{"issue", "update", "FID-9f13aa", "--lane", "routine", "--actor", "human:mark", "--idempotency-key", "lower"}},
		{"use-start", "transition-not-permitted", 4, nil, []string{"issue", "transition", "FID-9f13aa", "--to", "in_progress", "--actor", "agent:claude-builder", "--idempotency-key", "manual"}},
		{"use-hand-off", "transition-not-permitted", 4, func(f *statusFixture) {
			startSQL(f.t, f, `UPDATE issues SET status='in_progress' WHERE id='FID-9f13aa'`)
		}, []string{"issue", "transition", "FID-9f13aa", "--to", "in_review", "--actor", "human:mark", "--idempotency-key", "manual"}},
		{"no-return-to-todo", "transition-not-permitted", 4, func(f *statusFixture) {
			startSQL(f.t, f, `UPDATE issues SET status='in_progress' WHERE id='FID-9f13aa'`)
		}, []string{"issue", "transition", "FID-9f13aa", "--to", "todo", "--actor", "human:mark", "--idempotency-key", "manual"}},
		{"epic-use-start", "transition-not-permitted", 4, nil, []string{"issue", "transition", "FID-e3a1c0", "--to", "in_progress", "--actor", "human:mark", "--idempotency-key", "manual"}},
		{"ended-hold", "hold-not-current", 4, func(f *statusFixture) {
			startSQL(f.t, f, `UPDATE holds SET ended_at='2026-11-08T08:59:30Z',end_reason='takeover' WHERE id='HLD-aaaaaa'`)
		}, []string{"start", "--hold", "HLD-aaaaaa", "--actor", "agent:claude-builder", "--idempotency-key", "start"}},
		{"wrong-worktree", "hold-not-current", 4, func(f *statusFixture) {
			startSQL(f.t, f, `UPDATE holds SET worktree_id='sha256:elsewhere',worktree_root='/w/elsewhere' WHERE id='HLD-aaaaaa'`)
		}, []string{"start", "--hold", "HLD-aaaaaa", "--actor", "agent:claude-builder", "--idempotency-key", "start"}},
		{"missing-key", "invalid-argument", 2, nil, []string{"start", "--hold", "HLD-aaaaaa", "--actor", "agent:claude-builder"}},
		{"missing-hold", "hold-not-current", 4, nil, []string{"start", "--hold", "HLD-bbbbbb", "--actor", "agent:claude-builder", "--idempotency-key", "start"}},
		{"wrong-actor", "hold-not-current", 4, nil, []string{"start", "--hold", "HLD-aaaaaa", "--actor", "agent:other", "--idempotency-key", "start"}},
		{"expired-hold", "hold-not-current", 4, func(f *statusFixture) { f.clock = func() time.Time { return statusNow.Add(time.Hour) } }, []string{"start", "--hold", "HLD-aaaaaa", "--actor", "agent:claude-builder", "--idempotency-key", "start"}},
	} {
		t.Run(row.name, func(t *testing.T) {
			f := plannedStartFixture(t, "todo", stringPointer("standard"))
			if row.setup != nil {
				row.setup(f)
			}
			operation := "task.start"
			if row.args[0] == "issue" {
				operation = "issue." + row.args[1]
			}
			before := startEventCount(t, f.workspace)
			for _, jsonMode := range []bool{false, true} {
				args := row.args
				ext := ".txt"
				if jsonMode {
					args = append([]string{"--json"}, args...)
					ext = ".json"
				}
				result := f.run(args...)
				if cli.ExitCode(result.err) != row.exit || result.stdout != "" {
					t.Fatalf("refusal = %#v", result)
				}
				if jsonMode {
					assertJSONFailureCode(t, result.stderr, operation, row.code)
				}
				lifecycleGolden(t, f, filepath.Join("..", "start", row.name+ext), result.stderr)
			}
			if startEventCount(t, f.workspace) != before {
				t.Fatal("refusal wrote an Event")
			}
		})
	}
}

func startSQL(t *testing.T, f *statusFixture, query string) {
	t.Helper()
	conn, err := fiddb.Open(filepath.Join(config.LedgerDir(f.workspace), config.DatabaseFileName))
	if err != nil {
		t.Fatal(err)
	}
	defer closeCommandDB(t, conn)
	if _, err := conn.Exec(query); err != nil {
		t.Fatal(err)
	}
}
