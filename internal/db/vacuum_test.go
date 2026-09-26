package db

import (
	"path/filepath"
	"testing"
)

func TestVacuumCompactsSQLiteDatabaseWithoutChangingRows(t *testing.T) {
	t.Parallel()
	dbPath := filepath.Join(t.TempDir(), "fid.sqlite")
	if err := MigrateUp(dbPath, testMigrationsDir(t)); err != nil {
		t.Fatalf("MigrateUp() error = %v", err)
	}

	conn, err := Open(dbPath)
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	defer closeDB(t, conn)

	if _, err := conn.Exec(`INSERT INTO labels (name, description) VALUES ('vacuum', 'vacuum adapter test')`); err != nil {
		t.Fatalf("seed label: %v", err)
	}
	if err := Vacuum(conn); err != nil {
		t.Fatalf("Vacuum() error = %v", err)
	}

	var description string
	if err := conn.QueryRow(`SELECT description FROM labels WHERE name = 'vacuum'`).Scan(&description); err != nil {
		t.Fatalf("query label after vacuum: %v", err)
	}
	if description != "vacuum adapter test" {
		t.Fatalf("label description = %q, want vacuum adapter test", description)
	}
}

func TestVacuumRejectsNilConnection(t *testing.T) {
	t.Parallel()
	if err := Vacuum(nil); err == nil {
		t.Fatal("Vacuum(nil) error = nil, want connection error")
	}
}
