package db_test

import (
	"database/sql"
	"path/filepath"
	"testing"

	"github.com/golang-migrate/migrate/v4"
	migratesqlite "github.com/golang-migrate/migrate/v4/database/sqlite"
	fiddb "github.com/mbell-nzl/fid/internal/db"
)

// openSchema17DB preserves the pre-0018 subject of historical trigger tests.
func openSchema17DB(t *testing.T, migrationsDir string) *sql.DB {
	t.Helper()
	path := filepath.Join(t.TempDir(), "fid.sqlite")
	conn, err := fiddb.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	driver, err := migratesqlite.WithInstance(conn, &migratesqlite.Config{})
	if err != nil {
		t.Fatal(err)
	}
	absoluteDir, err := filepath.Abs(migrationsDir)
	if err != nil {
		t.Fatal(err)
	}
	migrator, err := migrate.NewWithDatabaseInstance("file://"+absoluteDir, "sqlite", driver)
	if err != nil {
		t.Fatal(err)
	}
	if err := migrator.Migrate(17); err != nil {
		t.Fatalf("migrate historical schema 17: %v", err)
	}
	sourceErr, databaseErr := migrator.Close()
	if sourceErr != nil || databaseErr != nil {
		t.Fatalf("close historical migrator: %v / %v", sourceErr, databaseErr)
	}
	if err := conn.Close(); err != nil {
		t.Fatal(err)
	}
	conn, err = fiddb.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	return conn
}
