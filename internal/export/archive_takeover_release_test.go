package export

import (
	"database/sql"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/mbell-nzl/fid/internal/holds"
	"github.com/mbell-nzl/fid/internal/journal"
	"github.com/mbell-nzl/fid/internal/project"
)

func TestTakeoverReplayAndImportValidateEveryReleasedReservation(t *testing.T) {
	t.Parallel()
	for _, path := range []string{"replay", "import"} {
		for _, test := range []struct {
			name, status, root string
			expires            string
			kept, valid        bool
		}{
			{"live unstarted", "todo", "/w/new", "2026-11-10T10:00:00Z", false, false},
			{"live started", "in_progress", "/w/new", "2026-11-10T10:00:00Z", false, false},
			{"lapsed started", "in_progress", "/w/new", "2026-11-10T09:05:00Z", false, false},
			{"parked reservation", "in_progress", "/w/new", "2026-11-10T10:00:00Z", true, false},
			{"expired unstarted elsewhere", "todo", "/w/elsewhere", "2026-11-10T09:05:00Z", false, false},
			{"expired unstarted destination", "todo", "/w/new", "2026-11-10T09:05:00Z", false, true},
		} {
			t.Run(path+"/"+test.name, func(t *testing.T) {
				t.Parallel()
				validator, err := newSemanticValidator(testMigrationsDir(t), 22)
				if err != nil {
					t.Fatal(err)
				}
				defer func() { _ = validator.close() }()
				seedTakeoverRelease(t, validator.conn, test.status, test.root, test.expires, test.kept)
				before := takeoverProjectionSnapshot(t, validator.conn)
				event := takeoverArchiveEvent()
				other := holds.Hold{ID: "HLD-b"}
				if test.kept {
					ended := "2026-11-10T09:01:00Z"
					other.EndedAt = &ended
				}
				release := other.Release(holds.EndReplaced, event.OccurredAt)
				event.Effects = append([]journal.Effect{release}, event.Effects...)
				if path == "import" {
					var raw []byte
					raw, err = journal.MarshalCanonicalJSON(event)
					if err != nil {
						t.Fatal(err)
					}
					err = validator.accept(raw)
				} else {
					err = journal.ProjectEvent(validator.conn, journal.ProjectionInput{Event: event,
						EventHash: "sha256:" + strings.Repeat("0", 64), File: "2026-11-10.jsonl", Line: 1})
				}
				if !test.valid {
					if err == nil {
						t.Error("takeover released another task's busy or unrelated reservation")
					}
					if after := takeoverProjectionSnapshot(t, validator.conn); after != before {
						t.Errorf("refused takeover changed holds or decisions:\nbefore: %s\nafter: %s", before, after)
					}
					return
				}
				if err != nil {
					t.Fatalf("legitimate expired-unstarted cleanup: %v", err)
				}
				var reason, newTask, binding string
				if err := validator.conn.QueryRow(`SELECT end_reason FROM holds WHERE id='HLD-b'`).Scan(&reason); err != nil || reason != "replaced" {
					t.Fatalf("cleanup reason = %q, %v", reason, err)
				}
				if err := validator.conn.QueryRow(`SELECT task_id FROM holds WHERE id='HLD-2' AND ended_at IS NULL`).Scan(&newTask); err != nil || newTask != "FID-a" {
					t.Fatalf("new hold task = %q, %v", newTask, err)
				}
				if err := validator.conn.QueryRow(`SELECT binding_id FROM person_decisions WHERE id='PDN-1'`).Scan(&binding); err != nil || binding != "FID-a" {
					t.Fatalf("takeover decision binding = %q, %v", binding, err)
				}
			})
		}
	}
}

func seedTakeoverRelease(t *testing.T, conn *sql.DB, status, root, expires string, kept bool) {
	t.Helper()
	for _, task := range []struct{ id, status string }{{"FID-a", "todo"}, {"FID-b", status}} {
		if _, err := conn.Exec(`INSERT INTO issues (id,kind,title,workflow_id,status,owner_role,created_at,updated_at)
			VALUES (?,'task','Work','task-default',?,'developer','2026-11-10T00:00:00Z','2026-11-10T00:00:00Z')`, task.id, task.status); err != nil {
			t.Fatal(err)
		}
	}
	for _, hold := range []struct{ id, task, root, expires string }{
		{"HLD-1", "FID-a", "/w/old", "2026-11-10T10:00:00Z"}, {"HLD-b", "FID-b", root, expires},
	} {
		if _, err := conn.Exec(`INSERT INTO holds (id,task_id,worktree_id,worktree_root,actor,since,expires_at)
			VALUES (?,?,?,?,'agent:builder','2026-11-10T09:00:00Z',?)`, hold.id, hold.task, project.WorktreeID(hold.root), hold.root, hold.expires); err != nil {
			t.Fatal(err)
		}
	}
	if kept {
		if _, err := conn.Exec(`UPDATE holds SET ended_at='2026-11-10T09:01:00Z',end_reason='interrupted',kept=1,park_reason='Usage limit' WHERE id='HLD-b'`); err != nil {
			t.Fatal(err)
		}
	}
	if _, err := conn.Exec(`INSERT INTO person_decisions VALUES ('PDN-old','takeover','human:mark','task','FID-b','Earlier decision','FID-b',?)`, holds.Timestamp(time.Date(2026, 11, 10, 8, 0, 0, 0, time.UTC))); err != nil {
		t.Fatal(err)
	}
}

func takeoverProjectionSnapshot(t *testing.T, conn *sql.DB) string {
	t.Helper()
	var snapshot []any
	for _, query := range []string{"SELECT * FROM holds ORDER BY id", "SELECT * FROM person_decisions ORDER BY id"} {
		rows, err := conn.Query(query)
		if err != nil {
			t.Fatal(err)
		}
		columns, err := rows.Columns()
		if err != nil {
			t.Fatal(err)
		}
		for rows.Next() {
			values := make([]any, len(columns))
			pointers := make([]any, len(columns))
			for i := range values {
				pointers[i] = &values[i]
			}
			if err := rows.Scan(pointers...); err != nil {
				t.Fatal(err)
			}
			snapshot = append(snapshot, values)
		}
		if err := rows.Err(); err != nil {
			t.Fatal(err)
		}
		if err := rows.Close(); err != nil {
			t.Fatal(err)
		}
	}
	raw, err := json.Marshal(snapshot)
	if err != nil {
		t.Fatal(err)
	}
	return string(raw)
}
