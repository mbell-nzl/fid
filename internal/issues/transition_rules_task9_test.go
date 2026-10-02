package issues

import (
	"database/sql"
	"testing"

	"github.com/mbell-nzl/fid/internal/workflow"
	_ "modernc.org/sqlite"
)

func TestLinkedToDecisionRequiresApprovedReferenceInIssueScope(t *testing.T) {
	t.Parallel()
	conn, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = conn.Close() }()
	if _, err := conn.Exec(`CREATE TABLE decision_refs (id TEXT PRIMARY KEY,class TEXT,outcome TEXT,scope_issue_id TEXT); INSERT INTO decision_refs VALUES ('DEC-approved','release','approved','ISS-target'),('DEC-rejected','release','rejected','ISS-other')`); err != nil {
		t.Fatal(err)
	}
	rule := workflow.RuleReference{Name: workflow.RuleLinkedToDecision}
	accepted, err := evaluateTransitionRule(conn, issueTransitionState{ID: "ISS-target"}, rule)
	if err != nil || !accepted.Satisfied {
		t.Fatalf("approved scoped Decision result = %#v, error = %v", accepted, err)
	}
	rejected, err := evaluateTransitionRule(conn, issueTransitionState{ID: "ISS-other"}, rule)
	if err != nil || rejected.Satisfied || len(rejected.Unmet) != 1 || rejected.Unmet[0].Rule != workflow.RuleLinkedToDecision {
		t.Fatalf("unqualified Decision result = %#v, error = %v", rejected, err)
	}
}
