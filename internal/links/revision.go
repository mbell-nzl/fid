package links

import (
	"database/sql"
	"fmt"
	"sort"

	"github.com/mbell-nzl/fid/internal/issues"
)

func linkIssueRevisions(conn *sql.DB, declared Link) ([]issues.IssueRevision, error) {
	issueIDs := []string{}
	if declared.FromType == "issue" {
		issueIDs = append(issueIDs, declared.FromID)
	}
	if declared.ToType == "issue" && declared.ToID != declared.FromID {
		issueIDs = append(issueIDs, declared.ToID)
	}
	if len(issueIDs) == 0 {
		return nil, fmt.Errorf("link %q %q %q has no Issue endpoint", declared.FromID, declared.Relation, declared.ToID)
	}

	revisions := make([]issues.IssueRevision, 0, len(issueIDs))
	for _, issueID := range issueIDs {
		var revision int
		if err := conn.QueryRow(`SELECT revision FROM issues WHERE id = ?`, issueID).Scan(&revision); err != nil {
			return nil, fmt.Errorf("read Issue revision for %q: %w", issueID, err)
		}
		revisions = append(revisions, issues.IssueRevision{IssueID: issueID, Observed: revision})
	}
	sort.Slice(revisions, func(left, right int) bool {
		return revisions[left].IssueID < revisions[right].IssueID
	})
	return revisions, nil
}

func resultIssueRevisions(observed []issues.IssueRevision) []IssueRevision {
	result := make([]IssueRevision, len(observed))
	for index, revision := range observed {
		result[index] = IssueRevision{ID: revision.IssueID, Revision: revision.Observed + 1}
	}
	return result
}
