package journal

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	fiddb "github.com/mbell-nzl/fid/internal/db"
)

func TestRebuildWorkspacePreservesFullPhaseTwoProjectionWithRefTriggersActive(t *testing.T) {
	t.Parallel()

	workspace := newReplayWorkspace(t)
	writeReplayEvents(
		t,
		filepath.Join(workspace, ".fid", "journal", "2026-06-07.jsonl"),
		replayRoleEvent("EVT-20260607-000901", "developer"),
		phaseTwoFullSchemaEvent(),
	)

	dbPath := filepath.Join(workspace, ".fid", "fid.sqlite")
	conn := openReplayDB(t, dbPath)

	stats, err := ReplayJournal(conn, filepath.Join(workspace, ".fid", "journal"), nil)
	if err != nil {
		t.Fatalf("ReplayJournal() error = %v", err)
	}
	if stats.EventsReplayed != 2 {
		t.Fatalf("initial replayed events = %d, want %d", stats.EventsReplayed, 2)
	}

	before := snapshotPhaseTwoProjectionState(t, conn)
	closeProjectDB(t, conn)

	if err := os.Remove(dbPath); err != nil {
		t.Fatalf("Remove(%q) error = %v", dbPath, err)
	}

	rebuildStats, err := RebuildWorkspace(workspace, testProjectMigrationsDir(t))
	if err != nil {
		t.Fatalf("RebuildWorkspace() error = %v", err)
	}
	if rebuildStats.EventsReplayed != 2 {
		t.Fatalf("rebuilt events = %d, want %d", rebuildStats.EventsReplayed, 2)
	}

	rebuiltConn, err := fiddb.Open(dbPath)
	if err != nil {
		t.Fatalf("Open() rebuilt error = %v", err)
	}
	defer closeProjectDB(t, rebuiltConn)

	after := snapshotPhaseTwoProjectionState(t, rebuiltConn)
	if before != after {
		t.Fatalf("rebuilt phase two state = %q, want %q", after, before)
	}
}

func phaseTwoFullSchemaEvent() Event {
	return Event{
		EventID:     "EVT-20260607-000902",
		EventSchema: EventSchemaV1,
		EventType:   "phase2.full-schema.seeded",
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
				Table: "requirements",
				PK:    map[string]any{"id": "REQ-p2-001"},
				Values: map[string]any{
					"type":             "PRD",
					"summary":          "Full-schema rebuild stays deterministic",
					"source_doc":       "01_PRD.md",
					"source_section":   "OPS",
					"last_synced_hash": "sha256:req-phase2",
					"status":           "active",
				},
			},
			{
				Op:    "insert",
				Table: "decisions",
				PK:    map[string]any{"id": "DEC-p2-001"},
				Values: map[string]any{
					"epic_id":      "FID-p2-epic",
					"author_role":  "tl",
					"summary":      "Keep trigger deferrals explicit",
					"body_md_path": "project-docs/decisions/DEC-p2-001.md",
					"area":         "schema",
					"status":       "open",
					"created_at":   "2026-06-07T15:43:00Z",
				},
			},
			{
				Op:    "insert",
				Table: "escalations",
				PK:    map[string]any{"id": "ESC-p2-001"},
				Values: map[string]any{
					"status":         "open",
					"trigger_class":  "policy",
					"raised_by_role": "qa_lead",
					"raised_at":      "2026-06-07T15:44:00Z",
				},
			},
			{
				Op:    "insert",
				Table: "triage_events",
				PK:    map[string]any{"id": "TRG-p2-001"},
				Values: map[string]any{
					"epic_id":       "FID-p2-epic",
					"issue_id":      "FID-p2-task-a",
					"triager_role":  "qa_lead",
					"source":        "task_completion",
					"source_ref_id": "TC-p2-001",
					"verdict":       "BLOCK",
					"esc_id":        "ESC-p2-001",
					"reasoning_md":  "project-docs/triage/TRG-p2-001.md",
					"created_at":    "2026-06-07T15:50:00Z",
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
					"relation":  "implements",
					"to_type":   "requirement",
					"to_id":     "REQ-p2-001",
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

func snapshotPhaseTwoProjectionState(t *testing.T, conn *sql.DB) string {
	t.Helper()

	queries := []struct {
		name  string
		query string
	}{
		{
			name: "sqlite_master",
			query: `
				SELECT type, name, COALESCE(sql, 'NULL')
				FROM sqlite_master
				WHERE type IN ('table', 'view', 'trigger', 'index')
				  AND name NOT LIKE 'sqlite_%'
				ORDER BY type, name
			`,
		},
		{
			name:  "labels",
			query: `SELECT name, COALESCE(description, 'NULL'), COALESCE(colour, 'NULL') FROM labels ORDER BY name`,
		},
		{
			name: "issues",
			query: `
				SELECT id, kind, COALESCE(parent_id, 'NULL'), title, workflow_id, status
				FROM issues
				ORDER BY id
			`,
		},
		{
			name: "epic_metadata",
			query: `
				SELECT issue_id, scheduling, COALESCE(brief_md_path, 'NULL'), COALESCE(date(target_start), 'NULL'), COALESCE(date(target_end), 'NULL')
				FROM epic_metadata
				ORDER BY issue_id
			`,
		},
		{
			name:  "requirements",
			query: `SELECT id, type, summary, source_doc, COALESCE(source_section, 'NULL'), COALESCE(last_synced_hash, 'NULL'), COALESCE(status, 'NULL') FROM requirements ORDER BY id`,
		},
		{
			name:  "decisions",
			query: `SELECT id, COALESCE(epic_id, 'NULL'), author_role, summary, body_md_path, COALESCE(area, 'NULL'), status, COALESCE(created_at, 'NULL') FROM decisions ORDER BY id`,
		},
		{
			name:  "escalations",
			query: `SELECT id, status, COALESCE(trigger_class, 'NULL'), COALESCE(raised_by_role, 'NULL'), COALESCE(raised_at, 'NULL'), COALESCE(resolved_at, 'NULL'), COALESCE(resolution_dec, 'NULL') FROM escalations ORDER BY id`,
		},
		{
			name:  "triage_events",
			query: `SELECT id, epic_id, issue_id, triager_role, source, COALESCE(source_ref_id, 'NULL'), verdict, COALESCE(esc_id, 'NULL'), COALESCE(reasoning_md, 'NULL'), COALESCE(created_at, 'NULL') FROM triage_events ORDER BY id`,
		},
		{
			name:  "acceptance_criteria",
			query: `SELECT id, issue_id, ordinal, text FROM acceptance_criteria ORDER BY id`,
		},
		{
			name:  "labels",
			query: `SELECT name, COALESCE(colour, 'NULL'), COALESCE(description, 'NULL') FROM labels ORDER BY name`,
		},
		{
			name:  "issue_labels",
			query: `SELECT issue_id, label FROM issue_labels ORDER BY issue_id, label`,
		},
		{
			name:  "links",
			query: `SELECT from_type, from_id, relation, to_type, to_id, COALESCE(created_at, 'NULL') FROM links ORDER BY from_type, from_id, relation, to_type, to_id`,
		},
		{
			name:  "activity_log",
			query: `SELECT id, journal_event_id, ref_type, ref_id, event, COALESCE(details, 'NULL'), COALESCE(actor, 'NULL'), COALESCE(occurred_at, 'NULL') FROM activity_log ORDER BY id`,
		},
		{
			name:  "issues_enriched",
			query: `SELECT id, concerns_raised, open_concerns_count, COALESCE(triage_status, 'NULL') FROM issues_enriched ORDER BY id`,
		},
		{
			name:  "projection_meta",
			query: `SELECT key, COALESCE(journal_file, 'NULL'), journal_offset, COALESCE(last_event_id, 'NULL'), COALESCE(last_event_hash, 'NULL'), event_schema, projection_schema_version, replay_status, updated_at FROM projection_meta ORDER BY key`,
		},
	}

	parts := make([]string, 0, len(queries))
	for _, spec := range queries {
		parts = append(parts, dumpPhaseTwoQuery(t, conn, spec.name, spec.query))
	}

	return strings.Join(parts, "\n")
}

func dumpPhaseTwoQuery(t *testing.T, conn *sql.DB, name, query string) string {
	t.Helper()

	rows, err := conn.Query(query)
	if err != nil {
		t.Fatalf("query %s: %v", name, err)
	}
	defer closeProjectDB(t, rows)

	columns, err := rows.Columns()
	if err != nil {
		t.Fatalf("columns %s: %v", name, err)
	}

	var result []string
	for rows.Next() {
		values := make([]any, len(columns))
		destinations := make([]any, len(columns))
		for index := range destinations {
			destinations[index] = &values[index]
		}
		if err := rows.Scan(destinations...); err != nil {
			t.Fatalf("scan %s: %v", name, err)
		}

		row := make([]string, len(values))
		for index, value := range values {
			row[index] = formatPhaseTwoSnapshotValue(value)
		}
		result = append(result, strings.Join(row, "|"))
	}
	if err := rows.Err(); err != nil {
		t.Fatalf("iterate %s: %v", name, err)
	}

	return name + "=" + strings.Join(result, ";")
}

func formatPhaseTwoSnapshotValue(value any) string {
	switch typed := value.(type) {
	case nil:
		return "NULL"
	case []byte:
		return string(typed)
	case time.Time:
		return typed.UTC().Format(time.RFC3339Nano)
	default:
		return fmt.Sprint(typed)
	}
}
