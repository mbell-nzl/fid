package journal

// Snapshot every projected table for generic atomicity and rebuild assertions.
import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"testing"
	"time"
)

func currentProjectionSnapshot(t *testing.T, conn *sql.DB) map[string][]string {
	t.Helper()
	names, err := conn.Query(`SELECT name FROM sqlite_master WHERE type='table' AND name NOT LIKE 'sqlite_%' ORDER BY name`)
	if err != nil {
		t.Fatal(err)
	}
	var tables []string
	for names.Next() {
		var name string
		if err := names.Scan(&name); err != nil {
			t.Fatal(err)
		}
		tables = append(tables, name)
	}
	if err := names.Err(); err != nil {
		t.Fatal(err)
	}
	if err := names.Close(); err != nil {
		t.Fatal(err)
	}
	result := map[string][]string{}
	for _, table := range tables {
		rows, err := conn.Query(`SELECT * FROM "` + table + `" ORDER BY rowid`)
		if err != nil {
			t.Fatal(err)
		}
		columns, err := rows.Columns()
		if err != nil {
			t.Fatal(err)
		}
		result[table] = []string{}
		for rows.Next() {
			values := make([]any, len(columns))
			pointers := make([]any, len(columns))
			for i := range values {
				pointers[i] = &values[i]
			}
			if err := rows.Scan(pointers...); err != nil {
				t.Fatal(err)
			}
			raw, err := json.Marshal(values)
			if err != nil {
				t.Fatal(err)
			}
			result[table] = append(result[table], string(raw))
		}
		if err := rows.Err(); err != nil {
			t.Fatal(err)
		}
		if err := rows.Close(); err != nil {
			t.Fatal(err)
		}
	}
	return result
}

func dumpProjectionQuery(t *testing.T, conn *sql.DB, name, query string) string {
	t.Helper()

	rows, err := conn.Query(query)
	if err != nil {
		t.Fatalf("query %s: %v", name, err)
	}
	defer closeProjectDB(t, rows)

	columns, err := rows.Columns()
	if err != nil {
		t.Fatalf("columns %s: %v", name, err)
	}

	var result []string
	for rows.Next() {
		values := make([]any, len(columns))
		destinations := make([]any, len(columns))
		for index := range destinations {
			destinations[index] = &values[index]
		}
		if err := rows.Scan(destinations...); err != nil {
			t.Fatalf("scan %s: %v", name, err)
		}

		row := make([]string, len(values))
		for index, value := range values {
			row[index] = formatProjectionSnapshotValue(value)
		}
		result = append(result, strings.Join(row, "|"))
	}
	if err := rows.Err(); err != nil {
		t.Fatalf("iterate %s: %v", name, err)
	}

	return name + "=" + strings.Join(result, ";")
}

func formatProjectionSnapshotValue(value any) string {
	switch typed := value.(type) {
	case nil:
		return "NULL"
	case []byte:
		return string(typed)
	case time.Time:
		return typed.UTC().Format(time.RFC3339Nano)
	default:
		return fmt.Sprint(typed)
	}
}
