package maintenance

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/mbell-nzl/fid/internal/config"
	fiddb "github.com/mbell-nzl/fid/internal/db"
	"github.com/mbell-nzl/fid/internal/journal"
	"github.com/mbell-nzl/fid/internal/policy"
)

// VacuumResult reports the derived SQLite projection size before and after compaction.
type VacuumResult struct {
	BeforeBytes int64
	AfterBytes  int64
}

// Vacuum compacts the derived SQLite projection without changing journal history.
func Vacuum(workspaceRoot, migrationsDir, byRole string) (VacuumResult, error) {
	if err := policy.Authorize(workspaceRoot, byRole, "maintenance.vacuum"); err != nil {
		return VacuumResult{}, err
	}

	var result VacuumResult
	err := journal.WithLock(workspaceRoot, journal.LockOwner{
		Command:   "fid maintenance vacuum --by " + byRole,
		ActorRole: byRole,
	}, func(lock *journal.WriteLock) error {
		if _, err := journal.EnsureProjectionFreshWithLock(lock, migrationsDir); err != nil {
			return err
		}
		vacuumResult, err := vacuumProjection(filepath.Join(workspaceRoot, config.WorkspaceDirName, "fid.sqlite"))
		if err != nil {
			return err
		}
		result = vacuumResult
		return nil
	})
	if err != nil {
		return VacuumResult{}, err
	}
	return result, nil
}

func vacuumProjection(dbPath string) (VacuumResult, error) {
	before, err := databaseFileSize(dbPath)
	if err != nil {
		return VacuumResult{}, err
	}

	conn, err := fiddb.Open(dbPath)
	if err != nil {
		return VacuumResult{}, err
	}
	defer func() { _ = conn.Close() }()

	if err := fiddb.Vacuum(conn); err != nil {
		return VacuumResult{}, err
	}

	after, err := databaseFileSize(dbPath)
	if err != nil {
		return VacuumResult{}, err
	}
	return VacuumResult{BeforeBytes: before, AfterBytes: after}, nil
}

func databaseFileSize(path string) (int64, error) {
	info, err := os.Stat(path)
	if err != nil {
		return 0, fmt.Errorf("stat sqlite database %q: %w", path, err)
	}
	return info.Size(), nil
}
