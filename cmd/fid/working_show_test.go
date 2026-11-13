package main

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/mbell-nzl/fid/internal/config"
	fiddb "github.com/mbell-nzl/fid/internal/db"
	"github.com/mbell-nzl/fid/internal/holds"
)

func TestShowWorkingGolden(t *testing.T) {
	t.Parallel()
	for _, test := range []struct {
		name, task, root string
		ago              time.Duration
	}{
		{"working-here", "FID-9f13aa", "", time.Minute},
		{"working-elsewhere", "FID-0d4e8b", "/home/mark/projects/fid-wt3", 20 * time.Minute},
		{"interrupted", "FID-1f2e3d", "/home/mark/projects/fid-wt2", 2 * time.Hour},
		{"interrupted-yesterday", "FID-1f2e3d", "/home/mark/projects/fid-wt2", 26 * time.Hour},
		{"lapsed-unstarted", "FID-9f13aa", "", 2 * time.Hour},
		{"sequential", "FID-8a1b2c", "", 0},
		{"prioritised", "FID-e3a1c0", "", 0},
	} {
		for _, format := range []string{"txt", "json"} {
			t.Run(test.name+"."+format, func(t *testing.T) {
				t.Parallel()
				fixture := newStatusFixture(t, time.Hour)
				nextWork(fixture)
				if test.ago > 0 {
					// nextWork already holds the first sequential task.
					if test.task == "FID-1f2e3d" {
						fixture.retimeHold(test.task, test.ago)
					} else {
						fixture.hold("HLD-0a0a0a", test.task, test.root, test.ago)
					}
				}
				args := []string{"show", test.task}
				if format == "json" {
					args = append([]string{"--json"}, args...)
				}
				fields := commandResultFields(t, fixture, "show", test.task)
				if _, ok := fields["working"]; !ok {
					t.Fatal("show has no working field")
				}
				if _, ok := fields["scheduling"]; !ok {
					t.Fatal("show has no scheduling field")
				}
				switch test.name {
				case "lapsed-unstarted":
					if fields["working"] != nil {
						t.Fatalf("lapsed unstarted task is still held: %v", fields["working"])
					}
				case "sequential", "prioritised":
					if fields["scheduling"] != test.name {
						t.Fatalf("scheduling = %v, want %s", fields["scheduling"], test.name)
					}
				case "interrupted":
					if fields["working"].(map[string]any)["lapsed"] != true {
						t.Fatal("interrupted task's working isn't lapsed")
					}
				}
				fixture.golden(filepath.Join("show", test.name+"."+format), args...)
			})
		}
	}
}

func (fixture *statusFixture) retimeHold(task string, ago time.Duration) {
	fixture.t.Helper()
	conn, err := fiddb.Open(filepath.Join(config.LedgerDir(fixture.workspace), config.DatabaseFileName))
	if err != nil {
		fixture.t.Fatal(err)
	}
	defer closeCommandDB(fixture.t, conn)
	since := statusNow.Add(-ago)
	if _, err := conn.Exec(`UPDATE holds SET since = ?, expires_at = ? WHERE task_id = ?`, holds.Timestamp(since), holds.Timestamp(since.Add(holds.Lifetime)), task); err != nil {
		fixture.t.Fatal(err)
	}
}
