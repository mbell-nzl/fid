package db

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"testing/fstest"
)

func TestMigrationUsesProvidedFS(t *testing.T) {
	t.Parallel()
	files := fstest.MapFS{
		"0001_fixture.up.sql":   {Data: []byte("CREATE TABLE fixture (n INTEGER); INSERT INTO fixture VALUES (1);")},
		"0001_fixture.down.sql": {Data: []byte("DROP TABLE fixture;")},
		"0007_fixture.up.sql":   {Data: []byte("INSERT INTO fixture VALUES (7);")},
		"0007_fixture.down.sql": {Data: []byte("DELETE FROM fixture WHERE n = 7;")},
	}
	directory := t.TempDir()
	for name, file := range files {
		if err := os.WriteFile(filepath.Join(directory, name), file.Data, 0o600); err != nil {
			t.Fatal(err)
		}
	}
	for _, test := range []struct {
		name   string
		source fs.FS
	}{{"memory", files}, {"directory", os.DirFS(directory)}} {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			path := filepath.Join(t.TempDir(), "projection.sqlite")
			for range 2 { // A repeated up remains a no-op.
				if err := MigrateUp(path, test.source); err != nil {
					t.Fatal(err)
				}
			}
			conn, err := Open(path)
			if err != nil {
				t.Fatal(err)
			}
			defer closeDB(t, conn)
			var version, total int
			if err := conn.QueryRow("SELECT version FROM schema_migrations").Scan(&version); err != nil {
				t.Fatal(err)
			}
			if err := conn.QueryRow("SELECT sum(n) FROM fixture").Scan(&total); err != nil {
				t.Fatal(err)
			}
			if version != 7 || total != 8 {
				t.Fatalf("version=%d total=%d, want staged version 7 and total 8", version, total)
			}
			assertTableDoesNotExist(t, conn, "issues")
			if err := MigrateDown(path, test.source); err != nil {
				t.Fatal(err)
			}
			assertTableDoesNotExist(t, conn, "fixture")
		})
	}
}

func TestMigrationRefusesNilFS(t *testing.T) {
	t.Parallel()
	for _, test := range []struct {
		name string
		run  func(string, fs.FS) error
	}{{"up", MigrateUp}, {"down", MigrateDown}} {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			path := filepath.Join(t.TempDir(), "projection.sqlite")
			if err := test.run(path, nil); err == nil || !strings.Contains(err.Error(), "migrations filesystem is required") {
				t.Fatalf("nil FS error = %v, want required filesystem refusal", err)
			}
			if _, err := os.Stat(path); !os.IsNotExist(err) {
				t.Fatalf("database stat after refusal = %v, want no database", err)
			}
		})
	}
}
