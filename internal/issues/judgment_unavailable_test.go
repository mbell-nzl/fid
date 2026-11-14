package issues

import (
	"path/filepath"
	"testing"

	"github.com/mbell-nzl/fid/internal/config"
)

// Retired judgment writers have no projection target in the current schema.
func TestV1JudgmentProjectionIsAbsentAfterReplacement(t *testing.T) {
	t.Parallel()
	workspace := newIssueWorkspace(t)
	conn := openIssueDB(t, filepath.Join(config.LedgerDir(workspace), "fid.sqlite"))
	defer closeIssueTestDB(t, conn)
	for _, name := range []string{"task_completions", "ac_verifications", "judgment_concerns", "concerns"} {
		var count int
		if err := conn.QueryRow(`SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name=?`, name).Scan(&count); err != nil || count != 0 {
			t.Fatalf("removed table %q count = %d, error %v", name, count, err)
		}
	}
}
