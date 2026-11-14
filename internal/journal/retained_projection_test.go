package journal

import (
	"database/sql"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"

	"github.com/mbell-nzl/fid/internal/config"
)

func TestRebuildWorkspacePreservesFullRetainedProjectionWithRefTriggersActive(t *testing.T) {
	t.Parallel()
	root := newReplayWorkspace(t)
	path := filepath.Join(config.LedgerDir(root), "fid.sqlite")
	writeReplayEvents(t, filepath.Join(config.LedgerDir(root), "journal", "2026-06-07.jsonl"), retainedProjectionEvent())
	conn := openReplayDB(t, path)
	stats, err := ReplayJournal(conn, filepath.Join(config.LedgerDir(root), "journal"), nil)
	if err != nil || stats.EventsReplayed != 1 {
		t.Fatalf("initial replay = %#v/%v", stats, err)
	}
	before := currentProjectionSnapshot(t, conn)
	assertRetainedProjectionRows(t, conn)
	closeProjectDB(t, conn)
	if err := os.Remove(path); err != nil {
		t.Fatal(err)
	}
	stats, err = RebuildWorkspace(root, testProjectMigrationsDir(t))
	if err != nil || stats.EventsReplayed != 1 {
		t.Fatalf("rebuild = %#v/%v", stats, err)
	}
	conn = openReplayDB(t, path)
	defer closeProjectDB(t, conn)
	assertRetainedProjectionRows(t, conn)
	if after := currentProjectionSnapshot(t, conn); !reflect.DeepEqual(before, after) {
		t.Fatalf("full retained projection differs: before=%v after=%v", before, after)
	}
}

func TestProjectEventPopulatesRetainedRepresentativeTablesViaJournalEffects(t *testing.T) {
	t.Parallel()
	conn := openReplayDB(t, filepath.Join(t.TempDir(), "projection.sqlite"))
	defer closeProjectDB(t, conn)
	event := retainedProjectionEvent()
	hash, err := HashEvent(event)
	if err != nil {
		t.Fatal(err)
	}
	if err := ProjectEvent(conn, ProjectionInput{Event: event, EventHash: hash, File: "2026-06-07.jsonl", Line: 1, Offset: 4096}); err != nil {
		t.Fatal(err)
	}
	assertRetainedProjectionRows(t, conn)
}

func retainedProjectionEvent() Event {
	return Event{
		EventID:     "EVT-20260607-000902",
		EventSchema: EventSchemaV1,
		EventType:   "retained.projection.seeded",
		OccurredAt:  time.Date(2026, time.June, 7, 16, 0, 0, 0, time.UTC),
		Actor: Actor{
			Role: "developer",
		},
		Command: Command{
			Argv:       []string{"fid", "maintenance", "replay"},
			CWD:        ".",
			FIDVersion: "0.1.0",
		},
		Effects: []Effect{
			{
				Op:    "insert",
				Table: "issues",
				PK:    map[string]any{"id": "FID-p2-epic"},
				Values: map[string]any{
					"kind":        "epic",
					"title":       "Phase 2 epic",
					"workflow_id": "epic-default",
					"status":      "todo",
					"created_at":  "2026-06-07T15:40:00Z",
					"updated_at":  "2026-06-07T15:40:00Z",
				},
			},
			{
				Op:    "insert",
				Table: "epic_metadata",
				PK:    map[string]any{"issue_id": "FID-p2-epic"},
				Values: map[string]any{
					"scheduling":    "sequential",
					"brief_md_path": "project-docs/epics/phase-2.md",
					"target_start":  "2026-06-07",
					"target_end":    "2026-06-21",
				},
			},
			{
				Op:    "insert",
				Table: "issues",
				PK:    map[string]any{"id": "FID-p2-task-a"},
				Values: map[string]any{
					"kind":        "task",
					"parent_id":   "FID-p2-epic",
					"sequence":    1,
					"title":       "Phase 2 task A",
					"workflow_id": "task-default",
					"status":      "todo",
					"created_at":  "2026-06-07T15:41:00Z",
					"updated_at":  "2026-06-07T15:41:00Z",
				},
			},
			{
				Op:    "insert",
				Table: "issues",
				PK:    map[string]any{"id": "FID-p2-task-b"},
				Values: map[string]any{
					"kind":        "task",
					"parent_id":   "FID-p2-epic",
					"sequence":    2,
					"title":       "Phase 2 task B",
					"workflow_id": "task-default",
					"status":      "todo",
					"created_at":  "2026-06-07T15:42:00Z",
					"updated_at":  "2026-06-07T15:42:00Z",
				},
			},
			{
				Op:    "insert",
				Table: "acceptance_criteria",
				PK:    map[string]any{"id": "AC-p2-001"},
				Values: map[string]any{
					"issue_id": "FID-p2-task-a",
					"ordinal":  1,
					"text":     "Ref triggers protect widened schema rows",
				},
			},
			{
				Op:    "insert",
				Table: "labels",
				PK:    map[string]any{"name": "phase-2"},
				Values: map[string]any{
					"colour":      "#0a9396",
					"description": "Phase 2 scope",
				},
			},
			{
				Op:    "insert",
				Table: "issue_labels",
				PK: map[string]any{
					"issue_id": "FID-p2-task-a",
					"label":    "phase-2",
				},
				Values: map[string]any{},
			},
			{
				Op:    "insert",
				Table: "links",
				PK: map[string]any{
					"from_type": "issue",
					"from_id":   "FID-p2-task-a",
					"relation":  "blocks",
					"to_type":   "issue",
					"to_id":     "FID-p2-task-b",
				},
				Values: map[string]any{
					"created_at": "2026-06-07T15:48:00Z",
				},
			},
			{
				Op:    "activity",
				Table: "activity_log",
				PK:    map[string]any{"id": "ACT-p2-001"},
				Values: map[string]any{
					"journal_event_id": "EVT-20260607-000902",
					"ref_type":         "issue",
					"ref_id":           "FID-p2-task-a",
					"event":            "status_changed",
					"details":          `{"from":"todo","to":"blocked"}`,
					"actor":            "developer",
					"occurred_at":      "2026-06-07T15:51:00Z",
				},
			},
		},
	}
}

func assertRetainedProjectionRows(t *testing.T, conn *sql.DB) {
	t.Helper()
	want := map[string]int{"issues": 3, "epic_metadata": 1, "labels": 1, "acceptance_criteria": 1, "issue_labels": 1, "links": 1, "activity_log": 1}
	for table, count := range want {
		var got int
		if err := conn.QueryRow("SELECT count(*) FROM " + table).Scan(&got); err != nil || got != count {
			t.Fatalf("%s count=%d/%v, want %d", table, got, err, count)
		}
	}
	var violations int
	if err := conn.QueryRow("SELECT count(*) FROM pragma_foreign_key_check").Scan(&violations); err != nil || violations != 0 {
		t.Fatalf("foreign key violations=%d/%v", violations, err)
	}
	for _, row := range []struct{ query, want string }{
		{`SELECT name||'|'||colour||'|'||description FROM labels`, "phase-2|#0a9396|Phase 2 scope"},
		{`SELECT issue_id||'|'||scheduling||'|'||brief_md_path||'|'||date(target_start)||'|'||date(target_end) FROM epic_metadata`, "FID-p2-epic|sequential|project-docs/epics/phase-2.md|2026-06-07|2026-06-21"},
		{`SELECT id||'|'||issue_id||'|'||ordinal||'|'||text FROM acceptance_criteria`, "AC-p2-001|FID-p2-task-a|1|Ref triggers protect widened schema rows"},
		{`SELECT from_type||'|'||from_id||'|'||relation||'|'||to_type||'|'||to_id FROM links`, "issue|FID-p2-task-a|blocks|issue|FID-p2-task-b"},
		{`SELECT id||'|'||journal_event_id||'|'||ref_type||'|'||ref_id||'|'||event||'|'||details||'|'||actor FROM activity_log`, `ACT-p2-001|EVT-20260607-000902|issue|FID-p2-task-a|status_changed|{"from":"todo","to":"blocked"}|developer`},
	} {
		var got string
		if err := conn.QueryRow(row.query).Scan(&got); err != nil || got != row.want {
			t.Fatalf("retained row=%q/%v, want %q", got, err, row.want)
		}
	}
}
