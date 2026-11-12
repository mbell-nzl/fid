package status

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/mbell-nzl/fid/internal/config"
	fiddb "github.com/mbell-nzl/fid/internal/db"
	"github.com/mbell-nzl/fid/internal/issues"
	"github.com/mbell-nzl/fid/internal/journal"
	"github.com/mbell-nzl/fid/internal/next"
)

// readWork keeps status's grouping and next's queue on one projection snapshot.
// afterOverview lets tests commit a concurrent transition between the reads.
func readWork(options Options, afterOverview func()) (issues.Overview, next.Result, error) {
	if _, err := journal.EnsureProjectionFresh(options.WorkspaceRoot, options.MigrationsDir); err != nil {
		return issues.Overview{}, next.Result{}, err
	}
	conn, err := fiddb.Open(filepath.Join(config.LedgerDir(options.WorkspaceRoot), config.DatabaseFileName))
	if err != nil {
		return issues.Overview{}, next.Result{}, err
	}
	defer func() { _ = conn.Close() }()
	tx, err := conn.Begin()
	if err != nil {
		return issues.Overview{}, next.Result{}, fmt.Errorf("begin status snapshot: %w", err)
	}
	defer func() { _ = tx.Rollback() }()
	overview, err := issues.ReadOverviewFrom(tx)
	if err != nil {
		return issues.Overview{}, next.Result{}, err
	}
	if afterOverview != nil {
		afterOverview()
	}
	planned, err := next.ReadFrom(tx, next.Options{WorkspaceRoot: options.WorkspaceRoot, Clock: func() time.Time { return options.Now }})
	return overview, planned, err
}
