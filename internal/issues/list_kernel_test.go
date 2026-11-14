package issues

// These assertions cover retiring kernel/historical behavior; retain until its owning cut.
import (
	"database/sql"
	"path/filepath"
	"testing"

	"github.com/mbell-nzl/fid/internal/config"
)

func seedIssueListFixtureConcerns(t *testing.T, conn *sql.DB) {
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

func TestKernelListFiltersOpenConcerns(t *testing.T) {
	t.Parallel()
	workspace := newIssueWorkspace(t)
	conn := openIssueDB(t, filepath.Join(config.LedgerDir(workspace), "fid.sqlite"))
	defer closeIssueTestDB(t, conn)
	seedIssueListFixture(t, workspace)
	seedIssueListFixtureConcerns(t, conn)
	got, err := List(ListOptions{WorkspaceRoot: workspace, MigrationsDir: testIssueMigrationsDir(t), HasOpenConcerns: boolPointer(true)})
	if err != nil {
		t.Fatal(err)
	}
	if ids := issueListIDs(got); !equalStrings(ids, []string{"FID-task01"}) {
		t.Fatalf("open-concern filter=%v", ids)
	}
}
