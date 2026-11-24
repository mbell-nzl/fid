package export

import (
	"database/sql"
	"strings"
	"testing"
	"testing/fstest"
)

func TestSemanticProjectionUsesProvidedFSAndProfile(t *testing.T) {
	t.Parallel()
	files := fstest.MapFS{
		"0019_fixture.up.sql": {Data: []byte("CREATE TABLE fixture (n INTEGER); INSERT INTO fixture VALUES (19);")},
		"0020_fixture.up.sql": {Data: []byte("INSERT INTO fixture VALUES (20);")},
		"0021_fixture.up.sql": {Data: []byte("CREATE TABLE later (n INTEGER);")},
	}
	validator, err := newSemanticValidator(files, 20)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = validator.close() }()
	var version, total, later int
	if err := validator.conn.QueryRow("SELECT version FROM schema_migrations").Scan(&version); err != nil {
		t.Fatal(err)
	}
	if err := validator.conn.QueryRow("SELECT sum(n) FROM fixture").Scan(&total); err != nil {
		t.Fatal(err)
	}
	if err := validator.conn.QueryRow("SELECT count(*) FROM sqlite_schema WHERE name = 'later'").Scan(&later); err != nil {
		t.Fatal(err)
	}
	if version != 20 || total != 39 || later != 0 {
		t.Fatalf("version=%d total=%d later=%d, want 20, 39, 0", version, total, later)
	}
	if len(validator.tables) != 1 || validator.tables["fixture"].columns["n"] != "INTEGER" {
		t.Fatalf("trusted tables = %v, want only the supplied fixture", validator.tables)
	}
}

func TestSemanticProjectionRefusesNilFS(t *testing.T) {
	t.Parallel()
	validator, err := newSemanticValidator(nil, 20)
	if err == nil || !strings.Contains(err.Error(), "migrations filesystem is required") {
		t.Fatalf("nil FS validator = %v, %v; want required filesystem refusal", validator, err)
	}
	conn, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = conn.Close() }()
	if err := migrateSemanticProjection(conn, nil, 20); err == nil || !strings.Contains(err.Error(), "migrations filesystem is required") {
		t.Fatalf("nil FS migration = %v; want required filesystem refusal", err)
	}
}
