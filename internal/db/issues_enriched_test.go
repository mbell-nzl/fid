package db_test

import (
	"database/sql"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	fiddb "github.com/mbell-nzl/fid/internal/db"
	"github.com/mbell-nzl/fid/internal/journal"
)

func TestIssuesEnrichedReportsConcernCountsAndLatestTriageStatus(t *testing.T) {
	t.Parallel()

	conn := openIssuesEnrichedDB(t)
	defer closeIssuesEnrichedDB(t, conn)

	projectIssuesEnrichedFixture(t, conn)

	tests := []struct {
		issueID          string
		wantRaised       bool
		wantOpenCount    int
		wantTriageStatus *string
	}{
		{
			issueID:          "FID-open",
			wantRaised:       true,
			wantOpenCount:    2,
			wantTriageStatus: stringPointer("resumed"),
		},
		{
			issueID:          "FID-addressed",
			wantRaised:       false,
			wantOpenCount:    0,
			wantTriageStatus: nil,
		},
		{
			issueID:          "FID-none",
			wantRaised:       false,
			wantOpenCount:    0,
			wantTriageStatus: nil,
		},
		{
			issueID:          "FID-blocked",
			wantRaised:       false,
			wantOpenCount:    0,
			wantTriageStatus: stringPointer("blocked"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.issueID, func(t *testing.T) {
			var (
				gotRaised       bool
				gotOpenCount    int
				gotTriageStatus sql.NullString
			)
			if err := conn.QueryRow(`
				SELECT concerns_raised, open_concerns_count, triage_status
				FROM issues_enriched
				WHERE id = ?
			`, tt.issueID).Scan(&gotRaised, &gotOpenCount, &gotTriageStatus); err != nil {
				t.Fatalf("query issues_enriched(%q): %v", tt.issueID, err)
			}

			if gotRaised != tt.wantRaised {
				t.Fatalf("concerns_raised = %t, want %t", gotRaised, tt.wantRaised)
			}
			if gotOpenCount != tt.wantOpenCount {
				t.Fatalf("open_concerns_count = %d, want %d", gotOpenCount, tt.wantOpenCount)
			}

			if tt.wantTriageStatus == nil {
				if gotTriageStatus.Valid {
					t.Fatalf("triage_status = %q, want NULL", gotTriageStatus.String)
				}
				return
			}
			if !gotTriageStatus.Valid || gotTriageStatus.String != *tt.wantTriageStatus {
				t.Fatalf("triage_status = %#v, want %q", gotTriageStatus, *tt.wantTriageStatus)
			}
		})
	}
}

func TestIssuesEnrichedQueryPlanUsesConcernAndTriageIndexes(t *testing.T) {
	t.Parallel()

	conn := openIssuesEnrichedDB(t)
	defer closeIssuesEnrichedDB(t, conn)

	rows, err := conn.Query(`
		EXPLAIN QUERY PLAN
		SELECT concerns_raised, open_concerns_count, triage_status
		FROM issues_enriched
		WHERE id = ?
	`, "FID-open")
	if err != nil {
		t.Fatalf("EXPLAIN QUERY PLAN error = %v", err)
	}
	defer closeIssuesEnrichedRows(t, rows)

	var details []string
	for rows.Next() {
		var (
			id     int
			parent int
			aux    int
			detail string
		)
		if err := rows.Scan(&id, &parent, &aux, &detail); err != nil {
			t.Fatalf("scan EXPLAIN QUERY PLAN row: %v", err)
		}
		details = append(details, detail)
	}
	if err := rows.Err(); err != nil {
		t.Fatalf("iterate EXPLAIN QUERY PLAN rows: %v", err)
	}

	assertPlanContainsIndex(t, details, "idx_concerns_issue_status")
	assertPlanContainsIndex(t, details, "idx_triage_issue_created")
	assertPlanLacksScan(t, details, "concerns")
	assertPlanLacksScan(t, details, "triage_events")
}

func openIssuesEnrichedDB(t *testing.T) *sql.DB {
	t.Helper()

	dbPath := filepath.Join(t.TempDir(), "fid.sqlite")
	if err := fiddb.MigrateUp(dbPath, testIssuesEnrichedMigrationsDir(t)); err != nil {
		t.Fatalf("MigrateUp() error = %v", err)
	}

	conn, err := fiddb.Open(dbPath)
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	return conn
}

func projectIssuesEnrichedFixture(t *testing.T, conn *sql.DB) {
	t.Helper()

	event := journal.Event{
		EventID:     "EVT-20260607-000501",
		EventSchema: journal.EventSchemaV1,
		EventType:   "phase2.enriched.seeded",
		OccurredAt:  time.Date(2026, time.June, 7, 13, 0, 0, 0, time.UTC),
		Actor: journal.Actor{
			Role: "developer",
		},
		Command: journal.Command{
			Argv:       []string{"fid", "maintenance", "replay"},
			CWD:        ".",
			FIDVersion: "0.1.0",
		},
		Effects: []journal.Effect{
			{
				Op:    "insert",
				Table: "issues",
				PK: map[string]any{
					"id": "FID-epic01",
				},
				Values: map[string]any{
					"kind":        "epic",
					"title":       "Phase 2 enriched epic",
					"workflow_id": "epic-default",
					"status":      "todo",
					"created_at":  "2026-06-07T13:00:00Z",
					"updated_at":  "2026-06-07T13:00:00Z",
				},
			},
			{
				Op:    "insert",
				Table: "issues",
				PK: map[string]any{
					"id": "FID-open",
				},
				Values: map[string]any{
					"kind":        "task",
					"title":       "Open concern issue",
					"workflow_id": "task-default",
					"status":      "todo",
					"created_at":  "2026-06-07T13:01:00Z",
					"updated_at":  "2026-06-07T13:01:00Z",
				},
			},
			{
				Op:    "insert",
				Table: "issues",
				PK: map[string]any{
					"id": "FID-addressed",
				},
				Values: map[string]any{
					"kind":        "task",
					"title":       "Addressed concern issue",
					"workflow_id": "task-default",
					"status":      "todo",
					"created_at":  "2026-06-07T13:02:00Z",
					"updated_at":  "2026-06-07T13:02:00Z",
				},
			},
			{
				Op:    "insert",
				Table: "issues",
				PK: map[string]any{
					"id": "FID-none",
				},
				Values: map[string]any{
					"kind":        "task",
					"title":       "No concern issue",
					"workflow_id": "task-default",
					"status":      "todo",
					"created_at":  "2026-06-07T13:03:00Z",
					"updated_at":  "2026-06-07T13:03:00Z",
				},
			},
			{
				Op:    "insert",
				Table: "issues",
				PK: map[string]any{
					"id": "FID-blocked",
				},
				Values: map[string]any{
					"kind":        "task",
					"title":       "Blocked triage issue",
					"workflow_id": "task-default",
					"status":      "todo",
					"created_at":  "2026-06-07T13:04:00Z",
					"updated_at":  "2026-06-07T13:04:00Z",
				},
			},
			{
				Op:    "insert",
				Table: "task_completions",
				PK: map[string]any{
					"id": "TC-000501",
				},
				Values: map[string]any{
					"issue_id":           "FID-open",
					"attempt_number":     1,
					"developer_role":     "developer",
					"model":              "gpt-5",
					"status":             "dev_complete",
					"gate_results":       `["go test ./..."]`,
					"human_intervention": false,
					"completed_at":       "2026-06-07T13:05:00Z",
				},
			},
			{
				Op:    "insert",
				Table: "task_completions",
				PK: map[string]any{
					"id": "TC-000502",
				},
				Values: map[string]any{
					"issue_id":           "FID-addressed",
					"attempt_number":     1,
					"developer_role":     "developer",
					"model":              "gpt-5",
					"status":             "dev_complete",
					"gate_results":       `["go test ./..."]`,
					"human_intervention": false,
					"completed_at":       "2026-06-07T13:06:00Z",
				},
			},
			{
				Op:    "insert",
				Table: "task_completions",
				PK: map[string]any{
					"id": "TC-000503",
				},
				Values: map[string]any{
					"issue_id":           "FID-addressed",
					"attempt_number":     2,
					"developer_role":     "developer",
					"model":              "gpt-5",
					"status":             "qa_pass",
					"gate_results":       `["go test ./..."]`,
					"human_intervention": false,
					"completed_at":       "2026-06-07T13:07:00Z",
				},
			},
			{
				Op:    "insert",
				Table: "concerns",
				PK: map[string]any{
					"id": "CON-000501",
				},
				Values: map[string]any{
					"issue_id":                "FID-open",
					"raised_in_completion_id": "TC-000501",
					"raised_by_role":          "developer",
					"summary":                 "First open concern",
					"body_md_path":            "project-docs/concerns/CON-000501.md",
					"status":                  "open",
					"created_at":              "2026-06-07T13:05:30Z",
					"updated_at":              "2026-06-07T13:05:30Z",
				},
			},
			{
				Op:    "insert",
				Table: "concerns",
				PK: map[string]any{
					"id": "CON-000502",
				},
				Values: map[string]any{
					"issue_id":                "FID-open",
					"raised_in_completion_id": "TC-000501",
					"raised_by_role":          "developer",
					"summary":                 "Second open concern",
					"body_md_path":            "project-docs/concerns/CON-000502.md",
					"status":                  "open",
					"created_at":              "2026-06-07T13:05:45Z",
					"updated_at":              "2026-06-07T13:05:45Z",
				},
			},
			{
				Op:    "insert",
				Table: "concerns",
				PK: map[string]any{
					"id": "CON-000503",
				},
				Values: map[string]any{
					"issue_id":                   "FID-addressed",
					"raised_in_completion_id":    "TC-000502",
					"raised_by_role":             "developer",
					"summary":                    "Addressed concern",
					"body_md_path":               "project-docs/concerns/CON-000503.md",
					"status":                     "addressed",
					"addressed_in_completion_id": "TC-000503",
					"created_at":                 "2026-06-07T13:06:30Z",
					"updated_at":                 "2026-06-07T13:07:30Z",
				},
			},
			{
				Op:    "insert",
				Table: "triage_events",
				PK: map[string]any{
					"id": "TRG-000501",
				},
				Values: map[string]any{
					"epic_id":       "FID-epic01",
					"issue_id":      "FID-open",
					"triager_role":  "qa_lead",
					"source":        "task_completion",
					"source_ref_id": "TC-000501",
					"verdict":       "BLOCK",
					"reasoning_md":  "project-docs/triage/TRG-000501.md",
					"created_at":    "2026-06-07T13:08:00Z",
				},
			},
			{
				Op:    "insert",
				Table: "triage_events",
				PK: map[string]any{
					"id": "TRG-000502",
				},
				Values: map[string]any{
					"epic_id":       "FID-epic01",
					"issue_id":      "FID-open",
					"triager_role":  "qa_lead",
					"source":        "task_completion",
					"source_ref_id": "TC-000501",
					"verdict":       "RESUME",
					"reasoning_md":  "project-docs/triage/TRG-000502.md",
					"created_at":    "2026-06-07T13:09:00Z",
				},
			},
			{
				Op:    "insert",
				Table: "triage_events",
				PK: map[string]any{
					"id": "TRG-000503",
				},
				Values: map[string]any{
					"epic_id":       "FID-epic01",
					"issue_id":      "FID-blocked",
					"triager_role":  "qa_lead",
					"source":        "task_completion",
					"source_ref_id": "TC-000501",
					"verdict":       "RESUME",
					"reasoning_md":  "project-docs/triage/TRG-000503.md",
					"created_at":    "2026-06-07T13:10:00Z",
				},
			},
			{
				Op:    "insert",
				Table: "triage_events",
				PK: map[string]any{
					"id": "TRG-000504",
				},
				Values: map[string]any{
					"epic_id":       "FID-epic01",
					"issue_id":      "FID-blocked",
					"triager_role":  "qa_lead",
					"source":        "task_completion",
					"source_ref_id": "TC-000501",
					"verdict":       "BLOCK",
					"reasoning_md":  "project-docs/triage/TRG-000504.md",
					"created_at":    "2026-06-07T13:11:00Z",
				},
			},
		},
	}

	projectDBEvent(t, conn, event)
}

func projectDBEvent(t *testing.T, conn *sql.DB, event journal.Event) {
	t.Helper()

	eventHash, err := journal.HashEvent(event)
	if err != nil {
		t.Fatalf("HashEvent() error = %v", err)
	}

	if err := journal.ProjectEvent(conn, journal.ProjectionInput{
		Event:     event,
		EventHash: eventHash,
		File:      "2026-06-07.jsonl",
		Line:      1,
		Offset:    4096,
	}); err != nil {
		t.Fatalf("ProjectEvent() error = %v", err)
	}
}

func assertPlanContainsIndex(t *testing.T, details []string, indexName string) {
	t.Helper()

	for _, detail := range details {
		if strings.Contains(strings.ToLower(detail), strings.ToLower(indexName)) {
			return
		}
	}
	t.Fatalf("query plan details %q did not reference index %q", details, indexName)
}

func assertPlanLacksScan(t *testing.T, details []string, table string) {
	t.Helper()

	wantMissing := "scan " + strings.ToLower(table)
	for _, detail := range details {
		if strings.Contains(strings.ToLower(detail), wantMissing) {
			t.Fatalf("query plan detail %q unexpectedly scanned table %q", detail, table)
		}
	}
}

func stringPointer(value string) *string {
	return &value
}

func testIssuesEnrichedMigrationsDir(t *testing.T) string {
	t.Helper()

	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller(0) failed")
	}

	return filepath.Clean(filepath.Join(filepath.Dir(currentFile), "..", "..", "migrations"))
}

func closeIssuesEnrichedDB(t *testing.T, conn *sql.DB) {
	t.Helper()

	if err := conn.Close(); err != nil {
		t.Fatalf("Close() error = %v", err)
	}
}

func closeIssuesEnrichedRows(t *testing.T, rows *sql.Rows) {
	t.Helper()

	if err := rows.Close(); err != nil {
		t.Fatalf("Rows.Close() error = %v", err)
	}
}
