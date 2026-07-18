package maintenance

import (
	"bytes"
	"database/sql"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"

	"github.com/mbell-nzl/fid/internal/bootstrap"
	fiddb "github.com/mbell-nzl/fid/internal/db"
	"github.com/mbell-nzl/fid/internal/journal"
)

func TestVacuumPreservesProjectedStateJournalAndActivity(t *testing.T) {
	workspace := newVacuumWorkspace(t)
	journalPath := onlyVacuumJournalPath(t, workspace)
	journalBefore := readVacuumFile(t, journalPath)
	stateBefore := snapshotVacuumState(t, filepath.Join(workspace, ".fid", "fid.sqlite"))

	result, err := Vacuum(workspace, maintenanceMigrationsDir(t), "tl")
	if err != nil {
		t.Fatalf("Vacuum() error = %v", err)
	}
	if result.BeforeBytes <= 0 || result.AfterBytes <= 0 {
		t.Fatalf("Vacuum() result = %#v, want positive database sizes", result)
	}
	if journalAfter := readVacuumFile(t, journalPath); !bytes.Equal(journalAfter, journalBefore) {
		t.Fatal("Vacuum() changed journal bytes")
	}
	if stateAfter := snapshotVacuumState(t, filepath.Join(workspace, ".fid", "fid.sqlite")); !reflect.DeepEqual(stateAfter, stateBefore) {
		t.Fatalf("projected state after Vacuum() = %#v, want %#v", stateAfter, stateBefore)
	}
}

func TestVacuumWaitsForWriteLockThenCompletes(t *testing.T) {
	workspace := newVacuumWorkspace(t)
	completed := make(chan error, 1)

	err := journal.WithLock(workspace, journal.LockOwner{Command: "test holder", ActorRole: "tl"}, func(_ *journal.WriteLock) error {
		go func() {
			_, vacuumErr := Vacuum(workspace, maintenanceMigrationsDir(t), "tl")
			completed <- vacuumErr
		}()

		select {
		case vacuumErr := <-completed:
			return vacuumErr
		case <-time.After(100 * time.Millisecond):
			return nil
		}
	})
	if err != nil {
		t.Fatalf("WithLock() error = %v", err)
	}

	select {
	case vacuumErr := <-completed:
		if vacuumErr != nil {
			t.Fatalf("Vacuum() after lock release error = %v", vacuumErr)
		}
	case <-time.After(time.Second):
		t.Fatal("Vacuum() did not complete after write lock release")
	}
}

func TestVacuumRefreshesMissingProjectionBeforeCompacting(t *testing.T) {
	workspace := newVacuumWorkspace(t)
	dbPath := filepath.Join(workspace, ".fid", "fid.sqlite")
	journalPath := onlyVacuumJournalPath(t, workspace)
	journalBefore := readVacuumFile(t, journalPath)
	if err := os.Remove(dbPath); err != nil {
		t.Fatalf("remove projection: %v", err)
	}

	if _, err := Vacuum(workspace, maintenanceMigrationsDir(t), "human"); err != nil {
		t.Fatalf("Vacuum() error = %v", err)
	}

	state := snapshotVacuumState(t, dbPath)
	if len(state.Roles) != len(bootstrap.DefaultRoles()) {
		t.Fatalf("roles after missing-projection vacuum = %#v, want all default roles", state.Roles)
	}
	if journalAfter := readVacuumFile(t, journalPath); !bytes.Equal(journalAfter, journalBefore) {
		t.Fatal("Vacuum() changed journal bytes while refreshing projection")
	}
}

func newVacuumWorkspace(t *testing.T) string {
	t.Helper()

	workspace := t.TempDir()
	if _, err := bootstrap.InitializeWorkspace(bootstrap.InitOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: maintenanceMigrationsDir(t),
		Prefix:        "FID",
	}); err != nil {
		t.Fatalf("InitializeWorkspace() error = %v", err)
	}
	return workspace
}

func onlyVacuumJournalPath(t *testing.T, workspace string) string {
	t.Helper()

	matches, err := filepath.Glob(filepath.Join(workspace, ".fid", "journal", "*.jsonl"))
	if err != nil {
		t.Fatalf("glob journal files: %v", err)
	}
	if len(matches) != 1 {
		t.Fatalf("journal files = %#v, want exactly one", matches)
	}
	return matches[0]
}

func readVacuumFile(t *testing.T, path string) []byte {
	t.Helper()

	payload, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %q: %v", path, err)
	}
	return payload
}

type vacuumState struct {
	Roles    []string
	Activity int
}

func snapshotVacuumState(t *testing.T, dbPath string) vacuumState {
	t.Helper()

	conn, err := fiddb.Open(dbPath)
	if err != nil {
		t.Fatalf("open projection: %v", err)
	}
	defer closeVacuumDB(t, conn)

	rows, err := conn.Query(`SELECT name FROM roles ORDER BY name`)
	if err != nil {
		t.Fatalf("query roles: %v", err)
	}
	defer func() { _ = rows.Close() }()

	state := vacuumState{Roles: make([]string, 0)}
	for rows.Next() {
		var role string
		if err := rows.Scan(&role); err != nil {
			t.Fatalf("scan role: %v", err)
		}
		state.Roles = append(state.Roles, role)
	}
	if err := rows.Err(); err != nil {
		t.Fatalf("iterate roles: %v", err)
	}
	if err := conn.QueryRow(`SELECT COUNT(*) FROM activity_log`).Scan(&state.Activity); err != nil {
		t.Fatalf("count activity rows: %v", err)
	}
	return state
}

func closeVacuumDB(t *testing.T, conn *sql.DB) {
	t.Helper()
	if err := conn.Close(); err != nil {
		t.Fatalf("close projection: %v", err)
	}
}
