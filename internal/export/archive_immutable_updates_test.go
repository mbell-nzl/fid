package export

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/mbell-nzl/fid/internal/journal"
)

// Profile 20 predates person_decisions and the later lifecycle guard. Its
// private projection has no command workspace; this is the same historical
// hold shape exercised by TestSemanticValidatorChecksHoldEvents.
func TestImmutableUpdatesAcrossOldProfilePaths(t *testing.T) {
	t.Parallel()
	for _, mode := range []string{"dry-run", "project", "replay", "import"} {
		t.Run(mode, func(t *testing.T) {
			t.Parallel()
			v, err := newSemanticValidator(testMigrationsDir(t), 20)
			if err != nil {
				t.Fatal(err)
			}
			defer func() { _ = v.close() }()
			for _, id := range []string{"FID-a", "FID-b"} {
				if _, err := v.conn.Exec(`INSERT INTO issues(id,kind,title,workflow_id,status) VALUES(?,'task','Task','task-default','todo')`, id); err != nil {
					t.Fatal(err)
				}
			}
			if _, err := v.conn.Exec(`INSERT INTO holds(id,task_id,worktree_id,worktree_root,actor,since,expires_at) VALUES('HLD-0','FID-b','sha256:two','/w/two','agent:two','2026-11-10T08:00:00Z','2026-11-10T09:00:00Z')`); err != nil {
				t.Fatal(err)
			}
			e := archiveHoldEvent(func(e *journal.Event) {
				e.EventID = "EVT-20261110-000001"
				e.Effects = append([]journal.Effect{{Op: "update", Table: "holds", PK: map[string]any{"id": "HLD-0"}, Values: map[string]any{"task_id": "FID-b"}}}, e.Effects...)
			})
			raw, err := journal.MarshalCanonicalJSON(e)
			if err != nil {
				t.Fatal(err)
			}
			before := handoffProjectionSnapshot(t, v)
			switch mode {
			case "dry-run":
				err = journal.DryRunProjectEvent(v.conn, e)
			case "project":
				err = applyLeanEmptyEvent(v, "replay", e)
			case "import":
				err = applyLeanEmptyEvent(v, "import", e)
			case "replay":
				dir := t.TempDir()
				path := filepath.Join(dir, "2026-11-10.jsonl")
				stored := append(bytes.Clone(raw), '\n')
				if err := os.WriteFile(path, stored, 0o644); err != nil {
					t.Fatal(err)
				}
				_, err = journal.ReplayJournal(v.conn, dir, nil)
				got, readErr := os.ReadFile(path)
				if readErr != nil || !bytes.Equal(got, stored) {
					t.Fatalf("refusal changed journal: %v", readErr)
				}
			}
			var refused *journal.RowIdentityError
			if !errors.As(err, &refused) || refused.Table != "holds" {
				t.Fatalf("%s malformed hold update = %v, want shared typed refusal", mode, err)
			}
			if before != handoffProjectionSnapshot(t, v) {
				t.Fatal("refusal changed projection")
			}
			e.Effects = e.Effects[1:]
			if err := applyLeanEmptyEvent(v, "import", e); err != nil {
				t.Fatalf("valid write after refusal: %v", err)
			}
		})
	}
}

func TestTransitionHoldSubjectAcrossOldProfilePaths(t *testing.T) {
	t.Parallel()
	for _, mode := range []string{"dry-run", "project", "replay", "import"} {
		t.Run(mode, func(t *testing.T) {
			t.Parallel()
			v, err := newSemanticValidator(testMigrationsDir(t), 20)
			if err != nil {
				t.Fatal(err)
			}
			defer func() { _ = v.close() }()
			for _, id := range []string{"FID-a", "FID-b"} {
				if _, err := v.conn.Exec(`INSERT INTO issues(id,kind,title,workflow_id,status) VALUES(?,'task','Task','task-default','todo')`, id); err != nil {
					t.Fatal(err)
				}
			}
			if _, err := v.conn.Exec(`INSERT INTO issues(id,kind,parent_id,title,workflow_id,status) VALUES('FID-s','subtask','FID-a','Subtask','subtask-default','todo')`); err != nil {
				t.Fatal(err)
			}
			if _, err := v.conn.Exec(`INSERT INTO holds(id,task_id,worktree_id,worktree_root,actor,since,expires_at) VALUES('HLD-0','FID-b','sha256:two','/w/two','agent:two','2026-11-10T08:00:00Z','2026-11-10T09:00:00Z')`); err != nil {
				t.Fatal(err)
			}
			e := archiveHoldEvent(func(e *journal.Event) {
				e.EventID = "EVT-20261110-000001"
				e.EventType = "issue.transitioned"
				e.CommandName = "issue.transition"
				e.Aggregate = journal.Aggregate{ID: "FID-s", Type: "issue", Sequence: 1}
				e.Effects = []journal.Effect{
					{Op: "update", Table: "issues", PK: map[string]any{"id": "FID-s"}, Values: map[string]any{"status": "in_progress", "updated_at": "2026-11-10T09:00:00Z"}, Preconditions: map[string]any{"status": "todo"}},
					{Op: "update", Table: "issues", PK: map[string]any{"id": "FID-s"}, Values: map[string]any{"revision": 2}, Prior: map[string]any{"revision": 1}, Preconditions: map[string]any{"revision": 1}},
					{Op: "update", Table: "holds", PK: map[string]any{"id": "HLD-0"}, Values: map[string]any{"ended_at": "2026-11-10T09:00:00Z", "end_reason": "task-cancelled"}},
					{Op: "activity", Table: "activity_log", PK: map[string]any{"id": "ACT-move"}, Values: map[string]any{"journal_event_id": e.EventID, "ref_type": "issue", "ref_id": "FID-s", "event": "status-changed"}},
				}
			})
			before := handoffProjectionSnapshot(t, v)
			raw, err := journal.MarshalCanonicalJSON(e)
			if err != nil {
				t.Fatal(err)
			}
			switch mode {
			case "dry-run":
				err = journal.DryRunProjectEvent(v.conn, e)
			case "project":
				err = applyLeanEmptyEvent(v, "replay", e)
			case "import":
				err = applyLeanEmptyEvent(v, "import", e)
			case "replay":
				dir := t.TempDir()
				path := filepath.Join(dir, "2026-11-10.jsonl")
				stored := append(bytes.Clone(raw), '\n')
				if err := os.WriteFile(path, stored, 0o644); err != nil {
					t.Fatal(err)
				}
				_, err = journal.ReplayJournal(v.conn, dir, nil)
				got, readErr := os.ReadFile(path)
				if readErr != nil || !bytes.Equal(got, stored) {
					t.Fatalf("refusal changed journal: %v", readErr)
				}
			}
			var refused *journal.RowIdentityError
			if !errors.As(err, &refused) || refused.Table != "holds" {
				t.Fatalf("%s changed another task's hold: %v", mode, err)
			}
			if before != handoffProjectionSnapshot(t, v) {
				t.Fatal("refusal changed projection")
			}
			e.Effects = append(e.Effects[:2:2], e.Effects[3])
			if err := applyLeanEmptyEvent(v, "import", e); err != nil {
				t.Fatalf("valid subtask move after refusal: %v", err)
			}
		})
	}
}
