package issues

import (
	"path/filepath"
	"testing"
)

// The removed judgment writers have no projection target in schema 18.
func TestV1JudgmentProjectionIsAbsentAfterReplacement(t *testing.T) {
	t.Parallel()
	workspace := newIssueWorkspace(t)
	conn := openIssueDB(t, filepath.Join(workspace, ".fid", "fid.sqlite"))
	defer closeIssueTestDB(t, conn)
	for _, name := range []string{"task_completions", "ac_verifications", "judgment_concerns"} {
		var count int
		if err := conn.QueryRow(`SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name=?`, name).Scan(&count); err != nil || count != 0 {
			t.Fatalf("removed table %q count = %d, error %v", name, count, err)
		}
	}
	var count int
	if err := conn.QueryRow(`SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name='concerns'`).Scan(&count); err != nil || count != 1 {
		t.Fatalf("generic concerns table count = %d, error %v", count, err)
	}
}
