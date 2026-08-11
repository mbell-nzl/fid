package issues

import (
	"fmt"
	"sort"
	"strings"

	"github.com/mbell-nzl/fid/internal/journal"
)

// IssueRevision is the observed aggregate revision while a write lock is held.
// The resulting revision is always derived by BuildIssueRevisionEffects.
type IssueRevision struct {
	IssueID  string
	Observed int
}

// BuildIssueRevisionEffects returns one absolute revision update per observed
// Issue aggregate, sorted by Issue ID for deterministic Event payloads.
func BuildIssueRevisionEffects(observed []IssueRevision) ([]journal.Effect, error) {
	if len(observed) == 0 {
		return nil, fmt.Errorf("issue revision effects require at least one Issue")
	}

	revisions := append([]IssueRevision(nil), observed...)
	sort.Slice(revisions, func(left, right int) bool {
		return revisions[left].IssueID < revisions[right].IssueID
	})

	effects := make([]journal.Effect, 0, len(revisions))
	for index, revision := range revisions {
		if strings.TrimSpace(revision.IssueID) == "" {
			return nil, fmt.Errorf("issue revision effect %d: Issue ID is required", index)
		}
		if revision.Observed < 1 {
			return nil, fmt.Errorf("issue revision effect %q: observed revision must be positive", revision.IssueID)
		}
		if index > 0 && revision.IssueID == revisions[index-1].IssueID {
			return nil, fmt.Errorf("issue revision effects duplicate Issue %q", revision.IssueID)
		}
		effects = append(effects, journal.Effect{
			Op:    "update",
			Table: "issues",
			PK:    map[string]any{"id": revision.IssueID},
			Values: map[string]any{
				"revision": revision.Observed + 1,
			},
			Prior: map[string]any{
				"revision": revision.Observed,
			},
			Preconditions: map[string]any{
				"revision": revision.Observed,
			},
		})
	}
	return effects, nil
}
