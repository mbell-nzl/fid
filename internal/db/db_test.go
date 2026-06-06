package db

import (
	"context"
	"database/sql"
	"fmt"
	"path/filepath"
	"reflect"
	"runtime"
	"sync"
	"testing"
)

func TestOpenEnablesRequiredPragmas(t *testing.T) {
	t.Parallel()

	dbPath := t.TempDir() + "/fid.sqlite"

	conn, err := Open(dbPath)
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	defer closeDB(t, conn)

	assertPragmaValue(t, conn, "foreign_keys", 1)
	assertPragmaValue(t, conn, "recursive_triggers", 1)
}

func TestOpenEnablesRequiredPragmasOnEveryPhysicalConnection(t *testing.T) {
	t.Parallel()

	const connectionCount = 3

	dbPath := t.TempDir() + "/fid.sqlite"

	conn, err := Open(dbPath)
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	defer closeDB(t, conn)

	conn.SetMaxOpenConns(connectionCount)

	ctx := context.Background()
	heldConns := make(chan *sql.Conn, connectionCount)
	errs := make(chan error, connectionCount)

	var wg sync.WaitGroup
	for i := 0; i < connectionCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			heldConn, err := conn.Conn(ctx)
			if err != nil {
				errs <- fmt.Errorf("Conn() error: %w", err)
				return
			}

			heldConns <- heldConn
		}()
	}

	wg.Wait()
	close(errs)
	close(heldConns)

	for err := range errs {
		if err != nil {
			t.Fatal(err)
		}
	}

	var physicalConnections []*sql.Conn
	for heldConn := range heldConns {
		physicalConnections = append(physicalConnections, heldConn)
	}
	defer closeSQLConns(t, physicalConnections)

	seenPhysicalConnections := make(map[uintptr]struct{}, len(physicalConnections))
	for _, heldConn := range physicalConnections {
		physicalConnectionID := physicalConnectionPointer(t, heldConn)
		seenPhysicalConnections[physicalConnectionID] = struct{}{}
		assertConnPragmaValue(t, heldConn, "foreign_keys", 1)
		assertConnPragmaValue(t, heldConn, "recursive_triggers", 1)
	}

	if len(seenPhysicalConnections) < 2 {
		t.Fatalf("distinct physical connections = %d, want >= 2", len(seenPhysicalConnections))
	}
}

func TestMigrateUpCreatesMinimalSchema(t *testing.T) {
	t.Parallel()

	dbPath := t.TempDir() + "/fid.sqlite"

	if err := MigrateUp(dbPath, testMigrationsDir(t)); err != nil {
		t.Fatalf("MigrateUp() error = %v", err)
	}

	conn, err := Open(dbPath)
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	defer closeDB(t, conn)

	assertTableExists(t, conn, "projection_meta")
	assertTableExists(t, conn, "roles")
	assertTableExists(t, conn, "issues")
	assertTableDoesNotExist(t, conn, "epic_metadata")
}

func TestMigrateDownDropsMinimalSchema(t *testing.T) {
	t.Parallel()

	dbPath := t.TempDir() + "/fid.sqlite"
	migrationsDir := testMigrationsDir(t)

	if err := MigrateUp(dbPath, migrationsDir); err != nil {
		t.Fatalf("MigrateUp() error = %v", err)
	}

	if err := MigrateDown(dbPath, migrationsDir); err != nil {
		t.Fatalf("MigrateDown() error = %v", err)
	}

	conn, err := Open(dbPath)
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	defer closeDB(t, conn)

	assertTableDoesNotExist(t, conn, "projection_meta")
	assertTableDoesNotExist(t, conn, "roles")
	assertTableDoesNotExist(t, conn, "issues")
}

func TestMigrateRoundTripRestoresMinimalSchema(t *testing.T) {
	t.Parallel()

	dbPath := t.TempDir() + "/fid.sqlite"
	migrationsDir := testMigrationsDir(t)

	if err := MigrateUp(dbPath, migrationsDir); err != nil {
		t.Fatalf("initial MigrateUp() error = %v", err)
	}
	if err := MigrateDown(dbPath, migrationsDir); err != nil {
		t.Fatalf("MigrateDown() error = %v", err)
	}
	if err := MigrateUp(dbPath, migrationsDir); err != nil {
		t.Fatalf("second MigrateUp() error = %v", err)
	}

	conn, err := Open(dbPath)
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	defer closeDB(t, conn)

	assertTableExists(t, conn, "projection_meta")
	assertTableExists(t, conn, "roles")
	assertTableExists(t, conn, "issues")
}

func TestMigrateUpIsIdempotentOnReopen(t *testing.T) {
	t.Parallel()

	dbPath := t.TempDir() + "/fid.sqlite"
	migrationsDir := testMigrationsDir(t)

	if err := MigrateUp(dbPath, migrationsDir); err != nil {
		t.Fatalf("first MigrateUp() error = %v", err)
	}

	if err := MigrateUp(dbPath, migrationsDir); err != nil {
		t.Fatalf("second MigrateUp() error = %v", err)
	}
}

func TestIssuesCheckConstraintRejectsInvalidWorkflowForKind(t *testing.T) {
	t.Parallel()

	dbPath := t.TempDir() + "/fid.sqlite"

	if err := MigrateUp(dbPath, testMigrationsDir(t)); err != nil {
		t.Fatalf("MigrateUp() error = %v", err)
	}

	conn, err := Open(dbPath)
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	defer closeDB(t, conn)

	_, err = conn.Exec(`
		INSERT INTO issues (id, kind, title, workflow_id, status)
		VALUES (?, ?, ?, ?, ?)
	`, "FID-abc123", "epic", "Invalid epic", "task-default", "todo")
	if err == nil {
		t.Fatal("Exec() error = nil, want CHECK constraint failure")
	}
}

func assertPragmaValue(t *testing.T, conn *sql.DB, pragma string, want int) {
	t.Helper()

	var got int
	if err := conn.QueryRow("PRAGMA " + pragma + ";").Scan(&got); err != nil {
		t.Fatalf("query pragma %q: %v", pragma, err)
	}

	if got != want {
		t.Fatalf("pragma %q = %d, want %d", pragma, got, want)
	}
}

func assertConnPragmaValue(t *testing.T, conn *sql.Conn, pragma string, want int) {
	t.Helper()

	var got int
	if err := conn.QueryRowContext(context.Background(), "PRAGMA "+pragma+";").Scan(&got); err != nil {
		t.Fatalf("query pragma %q: %v", pragma, err)
	}

	if got != want {
		t.Fatalf("pragma %q = %d, want %d", pragma, got, want)
	}
}

func physicalConnectionPointer(t *testing.T, conn *sql.Conn) uintptr {
	t.Helper()

	var pointer uintptr
	if err := conn.Raw(func(driverConn any) error {
		value := reflect.ValueOf(driverConn)
		if value.Kind() != reflect.Pointer {
			return fmt.Errorf("driver connection kind = %s, want pointer", value.Kind())
		}

		pointer = value.Pointer()
		return nil
	}); err != nil {
		t.Fatalf("Raw() error = %v", err)
	}

	return pointer
}

func assertTableExists(t *testing.T, conn *sql.DB, tableName string) {
	t.Helper()

	if !tableExists(t, conn, tableName) {
		t.Fatalf("table %q does not exist", tableName)
	}
}

func assertTableDoesNotExist(t *testing.T, conn *sql.DB, tableName string) {
	t.Helper()

	if tableExists(t, conn, tableName) {
		t.Fatalf("table %q exists", tableName)
	}
}

func tableExists(t *testing.T, conn *sql.DB, tableName string) bool {
	t.Helper()

	var count int
	if err := conn.QueryRow(`
		SELECT COUNT(*)
		FROM sqlite_master
		WHERE type = 'table' AND name = ?
	`, tableName).Scan(&count); err != nil {
		t.Fatalf("query sqlite_master for %q: %v", tableName, err)
	}

	return count == 1
}

func testMigrationsDir(t *testing.T) string {
	t.Helper()

	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller(0) failed")
	}

	return filepath.Clean(filepath.Join(filepath.Dir(currentFile), "..", "..", "migrations"))
}

func closeDB(t *testing.T, conn *sql.DB) {
	t.Helper()

	if err := conn.Close(); err != nil {
		t.Fatalf("Close() error = %v", err)
	}
}

func closeSQLConns(t *testing.T, conns []*sql.Conn) {
	t.Helper()

	for _, conn := range conns {
		if err := conn.Close(); err != nil {
			t.Fatalf("Conn.Close() error = %v", err)
		}
	}
}
