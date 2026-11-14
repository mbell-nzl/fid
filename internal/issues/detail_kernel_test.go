package issues

// These assertions cover retiring kernel/historical behavior; retain until its owning cut.
import (
	"path/filepath"
	"testing"

	"github.com/mbell-nzl/fid/internal/config"
)

func TestKernelGetReturnsConcernDetailSections(t *testing.T) {
	t.Parallel()
	workspace := newIssueWorkspace(t)
	conn := openIssueDB(t, filepath.Join(config.LedgerDir(workspace), "fid.sqlite"))
	defer closeIssueTestDB(t, conn)
	seedIssueDetailFixture(t, workspace)
	if _, err := conn.Exec(`INSERT INTO concerns (id,subject_type,subject_id,actor_id,scope_json,severity,summary,reasoning,references_json,status,created_at,created_event_id)
		 VALUES ('CON-000401','issue','FID-a1b2c3','developer','{"paths":[]}','major','Need triage before closing','Needs review','[]','open','2026-05-27T10:18:00Z','EVT-concern-detail')`); err != nil {
		t.Fatal(err)
	}
	got, err := Get(GetOptions{WorkspaceRoot: workspace, MigrationsDir: testIssueMigrationsDir(t), ID: "FID-a1b2c3"})
	if err != nil {
		t.Fatal(err)
	}
	if !got.ConcernsRaised {
		t.Fatal("Get() concerns_raised = false, want true")
	}
	if got.OpenConcernsCount != 1 {
		t.Fatalf("Get() open_concerns_count = %d, want 1", got.OpenConcernsCount)
	}
	if len(got.Concerns) != 1 || got.Concerns[0].ID != "CON-000401" {
		t.Fatalf("Get() concerns = %#v, want one open concern", got.Concerns)
	}
}
