package main

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/mbell-nzl/fid/internal/config"
	fiddb "github.com/mbell-nzl/fid/internal/db"
	"github.com/mbell-nzl/fid/internal/journal"
)

func TestMaintenanceCheckRefreshesProjectionWithoutChangingJournal(t *testing.T) {
	t.Parallel()
	for _, testCase := range []struct {
		name  string
		setup func(*testing.T, string)
	}{
		{
			name: "missing projection rebuilds",
			setup: func(t *testing.T, workspace string) {
				t.Helper()
				if err := os.Remove(filepath.Join(config.LedgerDir(workspace), "fid.sqlite")); err != nil {
					t.Fatalf("remove projection: %v", err)
				}
			},
		},
		{
			name: "stale projection catches up",
			setup: func(t *testing.T, workspace string) {
				t.Helper()
				appendCommandIssueEvent(t, workspace, journal.Event{
					EventSchema: journal.EventSchemaV1,
					EventType:   "label.seeded",
					Actor:       journal.Actor{Role: "tl"},
					Command: journal.Command{
						Argv:       []string{"fid", "maintenance", "seed"},
						CWD:        ".",
						FIDVersion: "0.1.0",
					},
					Effects: []journal.Effect{{
						Op:     "insert",
						Table:  "labels",
						PK:     map[string]any{"name": "maintenance"},
						Values: map[string]any{"description": "maintenance fixture"},
					}},
				})
			},
		},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			workspace := clonePristineWorkspace(t)
			testCase.setup(t, workspace)

			journalPath := mustJournalPath(t, workspace)
			before := mustReadCommandFile(t, journalPath)
			stdout, stderr, err := runFIDCommandExpectError(t, workspace, "--json", "maintenance", "check")
			if err != nil {
				t.Fatalf("maintenance check error = %v", err)
			}
			if stderr != "" {
				t.Fatalf("maintenance check stderr = %q, want empty", stderr)
			}
			envelope := decodeCommandJSONStructure(t, stdout).(map[string]any)
			if envelope["event"] != nil || envelope["operation"] != "maintenance.check" || envelope["result"].(map[string]any)["status"] != "pass" {
				t.Fatalf("maintenance envelope = %#v, want successful maintenance.check read envelope", envelope)
			}
			if after := mustReadCommandFile(t, journalPath); !bytes.Equal(after, before) {
				t.Fatal("maintenance check changed journal bytes")
			}

			conn, err := fiddb.Open(filepath.Join(config.LedgerDir(workspace), "fid.sqlite"))
			if err != nil {
				t.Fatalf("open refreshed projection: %v", err)
			}
			defer closeCommandDB(t, conn)
			if testCase.name == "stale projection catches up" {
				var count int
				if err := conn.QueryRow(`SELECT COUNT(*) FROM labels WHERE name = 'maintenance'`).Scan(&count); err != nil {
					t.Fatalf("query caught-up label: %v", err)
				}
				if count != 1 {
					t.Fatalf("caught-up label count = %d, want 1", count)
				}
			}
		})
	}
}
