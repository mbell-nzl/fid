package db_test

import (
	"database/sql"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

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

	assertPlanContainsIndex(t, details, "concerns_by_subject_status")
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
	statements := []string{
		`INSERT INTO issues (id,kind,title,workflow_id,status,created_at,updated_at) VALUES
		 ('FID-epic01','epic','Enriched epic','epic-default','todo','2026-06-07T13:00:00Z','2026-06-07T13:00:00Z'),
		 ('FID-open','task','Open concern issue','task-default','todo','2026-06-07T13:01:00Z','2026-06-07T13:01:00Z'),
		 ('FID-addressed','task','Disposed concern issue','task-default','todo','2026-06-07T13:02:00Z','2026-06-07T13:02:00Z'),
		 ('FID-none','task','No concern issue','task-default','todo','2026-06-07T13:03:00Z','2026-06-07T13:03:00Z'),
		 ('FID-blocked','task','Blocked triage issue','task-default','todo','2026-06-07T13:04:00Z','2026-06-07T13:04:00Z')`,
		`INSERT INTO epic_metadata(issue_id,scheduling) VALUES('FID-epic01','sequential')`,
		`INSERT INTO concerns
		 (id,subject_type,subject_id,actor_id,scope_json,severity,summary,reasoning,references_json,status,created_at,created_event_id)
		 VALUES
		 ('CON-000501','issue','FID-open','developer','{"paths":[]}','major','First open concern','Needs review','[]','open','2026-06-07T13:05:30Z','EVT-501'),
		 ('CON-000502','issue','FID-open','developer','{"paths":[]}','major','Second open concern','Needs review','[]','open','2026-06-07T13:05:45Z','EVT-502'),
		 ('CON-000504','requirement','REQ-open','developer','{"paths":[]}','major','Non-Issue concern','Needs review','[]','open','2026-06-07T13:05:46Z','EVT-504')`,
		`INSERT INTO concerns
		 (id,subject_type,subject_id,actor_id,scope_json,severity,summary,reasoning,references_json,status,
		  disposition_action,disposition_reason,disposed_actor_id,disposed_at,disposed_event_id,created_at,created_event_id)
		 VALUES ('CON-000503','issue','FID-addressed','developer','{"paths":[]}','major','Disposed concern','Reviewed','[]',
		 'disposed','rejected','Resolved','qa_lead','2026-06-07T13:07:30Z','EVT-dispose-503','2026-06-07T13:06:30Z','EVT-503')`,
		`INSERT INTO triage_events (id,epic_id,issue_id,triager_role,source,source_ref_id,verdict,reasoning_md,created_at)
		 VALUES
		 ('TRG-000501','FID-epic01','FID-open','qa_lead','manual','REV-501','BLOCK','project-docs/triage/TRG-000501.md','2026-06-07T13:08:00Z'),
		 ('TRG-000502','FID-epic01','FID-open','qa_lead','manual','REV-502','RESUME','project-docs/triage/TRG-000502.md','2026-06-07T13:09:00Z'),
		 ('TRG-000503','FID-epic01','FID-blocked','qa_lead','manual','REV-503','RESUME','project-docs/triage/TRG-000503.md','2026-06-07T13:10:00Z'),
		 ('TRG-000504','FID-epic01','FID-blocked','qa_lead','manual','REV-504','BLOCK','project-docs/triage/TRG-000504.md','2026-06-07T13:11:00Z')`,
	}
	for _, statement := range statements {
		if _, err := conn.Exec(statement); err != nil {
			t.Fatalf("seed issues_enriched: %v; SQL=%s", err, statement)
		}
	}
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
