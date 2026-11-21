package status

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/mbell-nzl/fid/internal/config"
	fiddb "github.com/mbell-nzl/fid/internal/db"
	"github.com/mbell-nzl/fid/internal/decisionview"
	"github.com/mbell-nzl/fid/internal/issues"
	"github.com/mbell-nzl/fid/internal/journal"
	"github.com/mbell-nzl/fid/internal/mergerecord"
	"github.com/mbell-nzl/fid/internal/next"
)

// readWork keeps status's grouping and next's queue on one projection snapshot.
// afterOverview lets tests commit a concurrent transition between the reads.
func readWork(options Options, afterOverview func()) (issues.Overview, next.Result, *Result, error) {
	if _, err := journal.EnsureProjectionFresh(options.WorkspaceRoot, options.MigrationsDir); err != nil {
		return issues.Overview{}, next.Result{}, nil, err
	}
	conn, err := fiddb.Open(filepath.Join(config.LedgerDir(options.WorkspaceRoot), config.DatabaseFileName))
	if err != nil {
		return issues.Overview{}, next.Result{}, nil, err
	}
	defer func() { _ = conn.Close() }()
	tx, err := conn.Begin()
	if err != nil {
		return issues.Overview{}, next.Result{}, nil, fmt.Errorf("begin status snapshot: %w", err)
	}
	defer func() { _ = tx.Rollback() }()
	overview, err := issues.ReadOverviewFrom(tx)
	if err != nil {
		return issues.Overview{}, next.Result{}, nil, err
	}
	if afterOverview != nil {
		afterOverview()
	}
	planned, err := next.ReadFrom(tx, next.Options{WorkspaceRoot: options.WorkspaceRoot, Clock: func() time.Time { return options.Now }})

	if err != nil {
		return issues.Overview{}, next.Result{}, nil, err
	}
	result := &Result{NeedsYou: []NeedsYouItem{}, Warnings: []Warning{}, merges: map[string]*mergerecord.ReviewState{}, approvals: map[string]decisionview.Approval{}}
	tasks := append([]issues.OverviewTask{}, overview.InProgress.Tasks...)
	for _, epic := range overview.InProgress.Epics {
		tasks = append(tasks, epic.Tasks...)
	}
	for _, task := range tasks {
		state, err := mergerecord.ReadReviewState(tx, task.ID)
		if err != nil {
			return issues.Overview{}, next.Result{}, nil, err
		}
		if state != nil {
			result.merges[task.ID] = state
		}
		if err := result.readApproval(options.WorkspaceRoot, tx, task.ID); err != nil {
			return issues.Overview{}, next.Result{}, nil, err
		}
		if task.Status == "in_review" {
			result.addApproval(task.ID)
		}
	}
	return overview, planned, result, nil
}
