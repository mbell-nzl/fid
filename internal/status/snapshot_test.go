package status

import (
	"encoding/json"
	"path/filepath"
	"testing"
	"time"

	"github.com/mbell-nzl/fid/internal/bootstrap"
	"github.com/mbell-nzl/fid/internal/config"
	fiddb "github.com/mbell-nzl/fid/internal/db"
	"github.com/mbell-nzl/fid/internal/gittest"
	"github.com/mbell-nzl/fid/internal/issues"
	"github.com/mbell-nzl/fid/internal/next"
)

func TestStatusReadyTransitionBetweenReads(t *testing.T) {
	t.Parallel()
	for _, state := range []string{"cancelled", "in_progress"} {
		t.Run(state, func(t *testing.T) {
			t.Parallel()
			options, tasks := snapshotWorkspace(t)
			other := gittest.AddWorktree(t, options.WorkspaceRoot)
			result, err := build(options, func() {
				_, err := issues.Transition(issues.TransitionOptions{WorkspaceRoot: other,
					MigrationsDir: options.MigrationsDir, ID: tasks[0].ID, ToState: state, Reason: "Snapshot transition",
					Actor: "human:mark", IdempotencyKey: "transition-during-status"})
				if err != nil {
					t.Fatal(err)
				}
			})
			if err != nil {
				t.Fatal(err)
			}
			ready := snapshotReady(t, result)
			if len(ready.Tasks) != len(tasks) {
				t.Fatalf("READY tasks = %+v, want both original tasks", ready.Tasks)
			}
			for index, task := range ready.Tasks {
				if task.ID != tasks[index].ID || task.Title != tasks[index].Title || task.Status != "todo" || task.StatusCategory != "unstarted" {
					t.Errorf("READY task = %+v, want original task %+v", task, tasks[index])
				}
			}
			if ready.Next == nil || *ready.Next != tasks[0].ID {
				t.Errorf("READY next = %v, want original first task %s", ready.Next, tasks[0].ID)
			}
			// A later status sees the committed transition and selects the
			// remaining unstarted task, rather than retaining a stale snapshot.
			later, err := Build(options)
			if err != nil {
				t.Fatal(err)
			}
			updated := snapshotReady(t, later)
			if len(updated.Tasks) != 1 || updated.Tasks[0].ID != tasks[1].ID || updated.Next == nil || *updated.Next != tasks[1].ID {
				t.Fatalf("later READY = %+v, want only %s selected", updated, tasks[1].ID)
			}
		})
	}
}

func snapshotWorkspace(t *testing.T) (Options, []issues.CreateResult) {
	t.Helper()
	root := gittest.NewWorkspace(t)
	migrations, err := filepath.Abs(filepath.Join("..", "..", "migrations"))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := bootstrap.InitializeWorkspace(bootstrap.InitOptions{WorkspaceRoot: root, MigrationsDir: migrations}); err != nil {
		t.Fatal(err)
	}
	tasks := []issues.CreateResult{}
	for priority, title := range []string{"First task", "Second task"} {
		created, err := issues.Create(issues.CreateOptions{WorkspaceRoot: root, MigrationsDir: migrations, Kind: "task",
			Title: title, Priority: &priority, Actor: "human:mark", IdempotencyKey: title})
		if err != nil {
			t.Fatal(err)
		}
		tasks = append(tasks, created)
	}
	conn, err := fiddb.Open(filepath.Join(config.LedgerDir(root), config.DatabaseFileName))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = conn.Close() })
	// WAL lets the second connection commit while status holds its read
	// snapshot. The production database's journal mode is unchanged.
	var mode string
	if err := conn.QueryRow("PRAGMA journal_mode=WAL").Scan(&mode); err != nil || mode != "wal" {
		t.Fatalf("enable WAL: mode %q, error %v", mode, err)
	}
	return Options{WorkspaceRoot: root, MigrationsDir: migrations, Now: time.Now(),
		Environment: []string{"HOME=" + t.TempDir(), "XDG_CONFIG_HOME=" + t.TempDir()}}, tasks
}

type readySnapshot struct {
	Next  *string      `json:"next"`
	Tasks []next.Entry `json:"tasks"`
}

func snapshotReady(t *testing.T, result Result) readySnapshot {
	t.Helper()
	payload, err := json.Marshal(result.Fields()["ready"])
	if err != nil {
		t.Fatal(err)
	}
	var ready readySnapshot
	if err := json.Unmarshal(payload, &ready); err != nil {
		t.Fatal(err)
	}
	return ready
}
