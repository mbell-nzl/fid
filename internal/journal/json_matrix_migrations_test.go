package journal

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	fiddb "github.com/mbell-nzl/fid/internal/db"
)

// JSONMigrationFixtureForTest captures the current schema and seeded rows from
// the real migration runner. Matrix cases need the final profile, not thousands
// of repetitions of its historical table rewrites. No production SQL is edited.
// Compare schema, indexes, triggers, views and every seeded row before use.
func JSONMigrationFixtureForTest(t *testing.T) string {
	t.Helper()
	canonical := openReplayDB(t, filepath.Join(t.TempDir(), "canonical.sqlite"))
	defer closeProjectDB(t, canonical)
	var version int
	if err := canonical.QueryRow(`SELECT version FROM schema_migrations`).Scan(&version); err != nil {
		t.Fatal(err)
	}
	statements := jsonSchemaStatements(t, canonical)
	var script strings.Builder
	script.WriteString("PRAGMA defer_foreign_keys=ON;\n")
	for _, s := range statements {
		if s.kind == "table" {
			script.WriteString(s.sql + ";\n")
		}
	}
	for _, s := range statements {
		if s.kind == "table" {
			script.WriteString(jsonSeedStatements(t, canonical, s.name))
		}
	}
	for _, kind := range []string{"view", "index", "trigger"} {
		for _, s := range statements {
			if s.kind == kind {
				script.WriteString(s.sql + ";\n")
			}
		}
	}
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, fmt.Sprintf("%04d_json_matrix.up.sql", version)), []byte(script.String()), 0o600); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(t.TempDir(), "snapshot.sqlite")
	if err := fiddb.MigrateUp(path, dir); err != nil {
		t.Fatal(err)
	}
	snapshot, err := fiddb.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer closeProjectDB(t, snapshot)
	if !reflect.DeepEqual(statements, jsonSchemaStatements(t, snapshot)) || !reflect.DeepEqual(currentProjectionSnapshot(t, canonical), currentProjectionSnapshot(t, snapshot)) {
		t.Fatal("current-schema fixture differs from the real migrations")
	}
	return dir
}

type jsonSchemaStatement struct{ kind, name, sql string }

func jsonSchemaStatements(t *testing.T, conn *sql.DB) []jsonSchemaStatement {
	t.Helper()
	rows, err := conn.Query(`SELECT type,name,sql FROM sqlite_schema WHERE sql IS NOT NULL AND name NOT LIKE 'sqlite_%' AND tbl_name!='schema_migrations' ORDER BY type,name`)
	if err != nil {
		t.Fatal(err)
	}
	defer closeProjectDB(t, rows)
	var result []jsonSchemaStatement
	for rows.Next() {
		var s jsonSchemaStatement
		if err := rows.Scan(&s.kind, &s.name, &s.sql); err != nil {
			t.Fatal(err)
		}
		result = append(result, s)
	}
	if err := rows.Err(); err != nil {
		t.Fatal(err)
	}
	return result
}

func jsonSeedStatements(t *testing.T, conn *sql.DB, table string) string {
	t.Helper()
	quoted := `"` + strings.ReplaceAll(table, `"`, `""`) + `"`
	rows, err := conn.Query(`SELECT * FROM ` + quoted + ` LIMIT 0`)
	if err != nil {
		t.Fatal(err)
	}
	columns, err := rows.Columns()
	if err != nil {
		t.Fatal(err)
	}
	if err := rows.Close(); err != nil {
		t.Fatal(err)
	}
	expressions := make([]string, len(columns))
	for i, column := range columns {
		expressions[i] = `quote("` + strings.ReplaceAll(column, `"`, `""`) + `")`
	}
	rows, err = conn.Query(`SELECT ` + strings.Join(expressions, `||','||`) + ` FROM ` + quoted + ` ORDER BY rowid`)
	if err != nil {
		t.Fatal(err)
	}
	defer closeProjectDB(t, rows)
	var result strings.Builder
	for rows.Next() {
		var values string
		if err := rows.Scan(&values); err != nil {
			t.Fatal(err)
		}
		result.WriteString("INSERT INTO " + quoted + " VALUES(" + values + ");\n")
	}
	if err := rows.Err(); err != nil {
		t.Fatal(err)
	}
	return result.String()
}

func TestJSONMigrationFixtureMatchesCanonical(t *testing.T) {
	t.Parallel()
	JSONMigrationFixtureForTest(t)
}
