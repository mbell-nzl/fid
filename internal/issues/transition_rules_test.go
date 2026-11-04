package issues

import (
	"path/filepath"
	"testing"

	"github.com/mbell-nzl/fid/internal/config"
)

func TestRemovedConcernTransitionRuleHasNoSeededRequirement(t *testing.T) {
	t.Parallel()
	workspace := newIssueWorkspace(t)
	conn := openIssueDB(t, filepath.Join(config.LedgerDir(workspace), "fid.sqlite"))
	defer closeIssueTestDB(t, conn)
	for _, rule := range []string{"all_acs_verified", "has_completion_report", "concerns_addressed"} {
		var count int
		if err := conn.QueryRow(`SELECT COUNT(*) FROM workflow_transition_requires WHERE rule=?`, rule).Scan(&count); err != nil || count != 0 {
			t.Fatalf("removed rule %q count = %d, error %v", rule, count, err)
		}
	}
}
