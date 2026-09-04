//go:build integration

package integration

import (
	"database/sql"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	fiddb "github.com/mbell-nzl/fid/internal/db"
)

func TestV01ParityGoldenOutputs(t *testing.T) {
	binary := integrationFIDBinary(t)
	workspace := t.TempDir()
	copyParityFixture(t, filepath.Join("..", "testdata", "v01_parity", "fixture"), workspace)
	for _, args := range [][]string{{"--json", "issue", "get", "FID-a10001"}, {"--json", "issue", "get", "FID-b20002"}, {"--json", "issue", "get", "FID-c30003"}, {"--json", "issue", "list", "--kind", "epic"}, {"--json", "issue", "list", "--kind", "task"}, {"--json", "issue", "list", "--kind", "subtask"}} {
		got, _ := closureRunFID(t, binary, workspace, args...)
		if !strings.Contains(got, `"contract_version":"fid-cli.v1"`) || !strings.Contains(got, `"event":null`) {
			t.Fatalf("canonical read = %q", got)
		}
	}
}

func TestV01ParityFixtureRebuildPreservesEveryProjectionTable(t *testing.T) {
	binary := integrationFIDBinary(t)
	workspace := t.TempDir()
	copyParityFixture(t, filepath.Join("..", "testdata", "v01_parity", "fixture"), workspace)

	// A read materializes the live projection from the regenerated fixture journal.
	closureRunFID(t, binary, workspace, "--json", "issue", "get", "FID-a10001")
	before := parityProjectionDump(t, filepath.Join(workspace, ".fid", "fid.sqlite"))

	closureRunFID(t, binary, workspace, "rebuild")
	after := parityProjectionDump(t, filepath.Join(workspace, ".fid", "fid.sqlite"))
	if before != after {
		t.Fatalf("fixture projection after rebuild differs:\n%s", parityDiff([]byte(before), []byte(after)))
	}
}

func parityProjectionDump(t *testing.T, dbPath string) string {
	t.Helper()

	conn, err := fiddb.Open(dbPath)
	if err != nil {
		t.Fatalf("open fixture projection: %v", err)
	}
	defer func() { _ = conn.Close() }()

	tables, err := conn.Query(`
		SELECT name FROM sqlite_schema
		WHERE type = 'table' AND name NOT LIKE 'sqlite_%' AND name != 'schema_migrations'
		ORDER BY name
	`)
	if err != nil {
		t.Fatalf("list projection tables: %v", err)
	}
	defer func() { _ = tables.Close() }()

	var dump strings.Builder
	for tables.Next() {
		var table string
		if err := tables.Scan(&table); err != nil {
			t.Fatalf("scan projection table: %v", err)
		}
		fmt.Fprintf(&dump, "[%s]\n", table)
		appendParityTableRows(t, &dump, conn, table)
	}
	if err := tables.Err(); err != nil {
		t.Fatalf("iterate projection tables: %v", err)
	}

	return dump.String()
}

func appendParityTableRows(t *testing.T, dump *strings.Builder, conn *sql.DB, table string) {
	t.Helper()
	query := `SELECT * FROM "` + strings.ReplaceAll(table, `"`, `""`) + `" ORDER BY 1`
	rows, err := conn.Query(query)
	if err != nil {
		t.Fatalf("read projection table %q: %v", table, err)
	}
	defer func() { _ = rows.Close() }()

	columns, err := rows.Columns()
	if err != nil {
		t.Fatalf("read projection table %q columns: %v", table, err)
	}
	for rows.Next() {
		values := make([]any, len(columns))
		destinations := make([]any, len(columns))
		for index := range values {
			destinations[index] = &values[index]
		}
		if err := rows.Scan(destinations...); err != nil {
			t.Fatalf("scan projection table %q: %v", table, err)
		}
		for index, value := range values {
			if index > 0 {
				dump.WriteByte('|')
			}
			fmt.Fprint(dump, value)
		}
		dump.WriteByte('\n')
	}
	if err := rows.Err(); err != nil {
		t.Fatalf("iterate projection table %q: %v", table, err)
	}
}

func copyParityFixture(t *testing.T, source, destination string) {
	t.Helper()

	entries, err := os.ReadDir(source)
	if err != nil {
		t.Fatalf("read fixture directory %s: %v", source, err)
	}
	for _, entry := range entries {
		sourcePath := filepath.Join(source, entry.Name())
		destinationPath := filepath.Join(destination, entry.Name())
		if entry.IsDir() {
			if err := os.MkdirAll(destinationPath, 0o755); err != nil {
				t.Fatalf("create fixture directory %s: %v", destinationPath, err)
			}
			copyParityFixture(t, sourcePath, destinationPath)
			continue
		}
		copyParityFixtureFile(t, sourcePath, destinationPath)
	}
}

func copyParityFixtureFile(t *testing.T, source, destination string) {
	t.Helper()

	input, err := os.Open(source)
	if err != nil {
		t.Fatalf("open fixture file %s: %v", source, err)
	}
	defer func() { _ = input.Close() }()
	if err := os.MkdirAll(filepath.Dir(destination), 0o755); err != nil {
		t.Fatalf("create fixture parent %s: %v", destination, err)
	}
	output, err := os.Create(destination)
	if err != nil {
		t.Fatalf("create fixture file %s: %v", destination, err)
	}
	if _, err := io.Copy(output, input); err != nil {
		_ = output.Close()
		t.Fatalf("copy fixture file %s: %v", source, err)
	}
	if err := output.Close(); err != nil {
		t.Fatalf("close fixture file %s: %v", destination, err)
	}
}

func parityDiff(want, got []byte) string {
	limit := len(want)
	if len(got) < limit {
		limit = len(got)
	}
	for index := 0; index < limit; index++ {
		if want[index] != got[index] {
			return fmt.Sprintf("first difference at byte %d: want %q, got %q", index, want[index], got[index])
		}
	}
	return fmt.Sprintf("length differs: want %d bytes, got %d bytes", len(want), len(got))
}
