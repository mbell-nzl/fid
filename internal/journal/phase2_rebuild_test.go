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

	rebuildStats, err := RebuildWorkspace(workspace, testProjectMigrationsDir(t), "developer")
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
				Table: "task_completions",
				PK:    map[string]any{"id": "TC-p2-001"},
				Values: map[string]any{
					"issue_id":               "FID-p2-task-a",
					"attempt_number":         1,
					"developer_role":         "developer",
					"model":                  "gpt-5",
					"status":                 "dev_complete",
					"notes_md_path":          "project-docs/completions/TC-p2-001.md",
					"actual_tokens":          1234,
					"commit_hashes":          `["abc123"]`,
					"files_touched":          `["internal/db/ref_triggers_test.go"]`,
					"expected_files_touched": `["migrations/0005_v01_ref_triggers.up.sql"]`,
					"gates_passed":           3,
					"gates_failed":           0,
					"gate_results":           `[{"cmd":"go test ./...","exit_code":0}]`,
					"human_intervention":     false,
					"completed_at":           "2026-06-07T15:45:00Z",
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
				Table: "ac_verifications",
				PK:    map[string]any{"id": "ACV-p2-001"},
				Values: map[string]any{
					"ac_id":          "AC-p2-001",
					"verifier_role":  "qa_lead",
					"verifier_actor": "agent",
					"status":         "pass",
					"evidence_type":  "test",
					"evidence_ref":   "internal/journal/phase2_rebuild_test.go",
					"notes":          "Deterministic rebuild proof",
					"completion_id":  "TC-p2-001",
					"verified_at":    "2026-06-07T15:46:00Z",
				},
			},
			{
				Op:    "insert",
				Table: "agent_runs",
				PK:    map[string]any{"id": "RUN-p2-001"},
				Values: map[string]any{
					"issue_id":           "FID-p2-task-a",
					"completion_id":      "TC-p2-001",
					"role":               "developer",
					"model":              "gpt-5",
					"tokens_in":          500,
					"tokens_out":         734,
					"price_in_per_mtok":  0.5,
					"price_out_per_mtok": 1.5,
					"cost_usd":           0.11,
					"tool_call_count":    4,
					"provider_run_id":    "run-p2-001",
					"provider_cost_usd":  0.11,
					"started_at":         "2026-06-07T15:30:00Z",
					"ended_at":           "2026-06-07T15:45:00Z",
				},
			},
			{
				Op:    "insert",
				Table: "gate_runs",
				PK:    map[string]any{"id": "GATE-p2-001"},
				Values: map[string]any{
					"issue_id":         "FID-p2-task-a",
					"attempt_number":   1,
					"agent_run_id":     "RUN-p2-001",
					"recorded_by_role": "developer",
					"command":          "go test ./...",
					"exit_code":        0,
					"output_excerpt":   "ok",
					"ran_at":           "2026-06-07T15:44:00Z",
				},
			},
			{
				Op:    "insert",
				Table: "concerns",
				PK:    map[string]any{"id": "CON-p2-001"},
				Values: map[string]any{
					"issue_id":                "FID-p2-task-a",
					"raised_in_completion_id": "TC-p2-001",
					"raised_by_role":          "developer",
					"summary":                 "Validate rebuild parity",
					"body_md_path":            "project-docs/concerns/CON-p2-001.md",
					"status":                  "open",
					"escalated_to_esc_id":     "ESC-p2-001",
					"triaged_in_event_id":     "TRG-p2-001",
					"created_at":              "2026-06-07T15:47:00Z",
					"updated_at":              "2026-06-07T15:47:00Z",
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
				Op:    "insert",
				Table: "claims",
				PK: map[string]any{
					"ref_type": "issue",
					"ref_id":   "FID-p2-task-a",
				},
				Values: map[string]any{
					"claimed_by":  "developer",
					"claimed_at":  "2026-06-07T15:49:00Z",
					"ttl_seconds": 1800,
					"note":        "active claim",
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
			name:  "roles",
			query: `SELECT name, COALESCE(description, 'NULL'), active FROM roles ORDER BY name`,
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
			name:  "task_completions",
			query: `SELECT id, issue_id, attempt_number, COALESCE(developer_role, 'NULL'), COALESCE(model, 'NULL'), COALESCE(status, 'NULL'), COALESCE(notes_md_path, 'NULL'), COALESCE(actual_tokens, -1), COALESCE(commit_hashes, 'NULL'), COALESCE(files_touched, 'NULL'), COALESCE(expected_files_touched, 'NULL'), COALESCE(gates_passed, -1), COALESCE(gates_failed, -1), COALESCE(gate_results, 'NULL'), human_intervention, COALESCE(human_intervention_note, 'NULL'), COALESCE(completed_at, 'NULL') FROM task_completions ORDER BY id`,
		},
		{
			name:  "acceptance_criteria",
			query: `SELECT id, issue_id, ordinal, text FROM acceptance_criteria ORDER BY id`,
		},
		{
			name:  "ac_verifications",
			query: `SELECT id, ac_id, verifier_role, COALESCE(verifier_actor, 'NULL'), status, COALESCE(evidence_type, 'NULL'), COALESCE(evidence_ref, 'NULL'), COALESCE(notes, 'NULL'), COALESCE(completion_id, 'NULL'), verified_at FROM ac_verifications ORDER BY id`,
		},
		{
			name:  "agent_runs",
			query: `SELECT id, COALESCE(issue_id, 'NULL'), COALESCE(completion_id, 'NULL'), role, model, COALESCE(tokens_in, -1), COALESCE(tokens_out, -1), COALESCE(price_in_per_mtok, -1), COALESCE(price_out_per_mtok, -1), COALESCE(cost_usd, -1), COALESCE(token_breakdown_json, 'NULL'), COALESCE(tool_call_count, -1), COALESCE(provider_run_id, 'NULL'), COALESCE(provider_cost_usd, -1), COALESCE(started_at, 'NULL'), COALESCE(ended_at, 'NULL') FROM agent_runs ORDER BY id`,
		},
		{
			name:  "gate_runs",
			query: `SELECT id, issue_id, attempt_number, COALESCE(agent_run_id, 'NULL'), recorded_by_role, command, exit_code, COALESCE(output_excerpt, 'NULL'), ran_at FROM gate_runs ORDER BY id`,
		},
		{
			name:  "concerns",
			query: `SELECT id, issue_id, raised_in_completion_id, raised_by_role, summary, COALESCE(body_md_path, 'NULL'), status, COALESCE(addressed_in_completion_id, 'NULL'), COALESCE(escalated_to_esc_id, 'NULL'), COALESCE(triaged_in_event_id, 'NULL'), COALESCE(created_at, 'NULL'), COALESCE(updated_at, 'NULL') FROM concerns ORDER BY id`,
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
			name:  "claims",
			query: `SELECT ref_type, ref_id, claimed_by, claimed_at, ttl_seconds, COALESCE(note, 'NULL') FROM claims ORDER BY ref_type, ref_id`,
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
