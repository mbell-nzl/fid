package main

// These assertions cover retiring kernel/historical behavior; retain until its owning cut.
import (
	"database/sql"
	"path/filepath"
	"testing"

	"github.com/mbell-nzl/fid/internal/config"
	fiddb "github.com/mbell-nzl/fid/internal/db"
)

func seedIssueListCommandFixtureConcerns(t *testing.T, conn *sql.DB) {
	t.Helper()
	for _, statement := range []struct {
		query string
		args  []any
	}{
		{
			query: `INSERT INTO concerns
				(id, subject_type, subject_id, actor_id, scope_json, severity, summary, reasoning, references_json,
				 status, created_at, created_event_id)
				VALUES (?, 'issue', ?, 'developer', '{"paths":[]}', 'major', ?, 'Needs review', '[]',
				        'open', ?, ?)`,
			args: []any{"CON-open01", "FID-task01", "Needs triage", "2026-06-09T09:05:00Z", "EVT-open01"},
		},
		{
			query: `INSERT INTO concerns
				(id, subject_type, subject_id, actor_id, scope_json, severity, summary, reasoning, references_json,
				 status, disposition_action, disposition_reason, disposed_actor_id, disposed_at,
				 disposed_event_id, created_at, created_event_id)
				VALUES (?, 'issue', ?, 'developer', '{"paths":[]}', 'major', ?, 'Needs review', '[]',
				        'disposed', 'rejected', 'Addressed', 'qa_lead', ?, ?, ?, ?)`,
			args: []any{"CON-addressed01", "FID-task03", "Already addressed", "2026-06-09T09:06:00Z",
				"EVT-dispose01", "2026-06-09T09:06:00Z", "EVT-addressed01"},
		},
	} {
		if _, err := conn.Exec(statement.query, statement.args...); err != nil {
			t.Fatal(err)
		}
	}
}

func TestKernelIssueListCombinedConcernFilters(t *testing.T) {
	t.Parallel()
	workspace := clonePristineWorkspace(t)
	conn, err := fiddb.Open(filepath.Join(config.LedgerDir(workspace), "fid.sqlite"))
	if err != nil {
		t.Fatal(err)
	}
	defer closeCommandDB(t, conn)
	seedIssueListCommandFixture(t, workspace)
	seedIssueListCommandFixtureConcerns(t, conn)
	stdout, stderr := runFIDCommand(t, workspace, "--json", "issue", "list", "--epic", "FID-epic01", "--label", "core", "--priority", "1", "--has-concerns")
	if stderr != "" {
		t.Fatalf("stderr=%q", stderr)
	}
	rows := decodeCanonicalResult(t, stdout)["issues"].([]any)
	if ids := commandIssueListIDs(rows); !equalIssueIDs(ids, []string{"FID-task01"}) {
		t.Fatalf("combined concern filter=%v", ids)
	}
}
