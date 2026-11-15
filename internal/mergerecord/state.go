package mergerecord

import (
	"database/sql"
	"errors"

	"github.com/mbell-nzl/fid/internal/github"
)

// ReviewState describes only durable knowledge for an unfinished review. CI
// attempts in progress are not receipts: this is an offline view, never fresh
// confirmation that the task can complete.
type ReviewState struct {
	MergeCommit string         `json:"merge_commit"`
	CI          github.Verdict `json:"ci"`
}

func (s ReviewState) Summary() string {
	switch s.CI {
	case github.Waiting:
		return "merged, waiting for CI"
	case github.NotGreen:
		return "merged, CI not green"
	default:
		return ""
	}
}

// ReadReviewState binds the merge to the latest hand-off and evaluates only
// receipts on its merge commit using the adapter's shared attempt-selection rule.
func ReadReviewState(conn Reader, task string) (*ReviewState, error) {
	m, err := readMerge(conn, `SELECT m.record_json FROM merge_observations m JOIN issues i ON i.id=m.task_id WHERE i.id=? AND i.kind='task' AND i.status='in_review' AND m.handoff_id=(SELECT id FROM handoffs WHERE task_id=i.id ORDER BY rowid DESC LIMIT 1) AND NOT EXISTS(SELECT 1 FROM completions WHERE task_id=i.id) ORDER BY m.rowid DESC LIMIT 1`, task)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	receipts, err := CIResults(conn, task, m.ID)
	if err != nil {
		return nil, err
	}
	runs := make([]github.WorkflowRun, 0, len(receipts))
	for _, r := range receipts {
		runs = append(runs, github.WorkflowRun{ID: r.RunID, Attempt: r.Attempt, WorkflowID: r.WorkflowID, Name: r.Name, Path: r.Path, Event: r.Event, HeadSHA: r.HeadSHA, Status: r.Status, Conclusion: r.Conclusion, URL: r.URL})
	}
	return &ReviewState{MergeCommit: m.MergeCommit, CI: github.EvaluateCI(m.MergeCommit, runs).Verdict}, nil
}
