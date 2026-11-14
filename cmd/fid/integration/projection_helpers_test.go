//go:build integration

package integration

// Shared process/JSON/filesystem helpers have no kernel fixture dependency.
import (
	"archive/zip"
	"crypto/sha256"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/mbell-nzl/fid/internal/config"
)

func p22FaultDB(t *testing.T, root string) *sql.DB {
	t.Helper()
	conn, err := sql.Open("sqlite", filepath.Join(config.LedgerDir(root), config.DatabaseFileName))
	p22FaultOK(t, err)
	return conn
}

func p22FaultProjection(t *testing.T, root string) map[string][]string {
	t.Helper()
	conn := p22FaultDB(t, root)
	defer conn.Close()
	return p22FaultProjectionFromDB(t, conn)
}

func p22FaultProjectionFromDB(t *testing.T, conn *sql.DB) map[string][]string {
	t.Helper()
	names, err := conn.Query(`SELECT name FROM sqlite_master WHERE type = 'table' AND name NOT LIKE 'sqlite_%' ORDER BY name`)
	p22FaultOK(t, err)
	var tables []string
	for names.Next() {
		var name string
		p22FaultOK(t, names.Scan(&name))
		tables = append(tables, name)
	}
	p22FaultOK(t, names.Err())
	p22FaultOK(t, names.Close())
	result := map[string][]string{}
	for _, table := range tables {
		rows, err := conn.Query(`SELECT * FROM "` + table + `"`)
		p22FaultOK(t, err)
		columns, err := rows.Columns()
		p22FaultOK(t, err)
		result[table] = []string{}
		for rows.Next() {
			values, pointers := make([]any, len(columns)), make([]any, len(columns))
			for i := range values {
				pointers[i] = &values[i]
			}
			p22FaultOK(t, rows.Scan(pointers...))
			encoded, err := json.Marshal(values)
			p22FaultOK(t, err)
			result[table] = append(result[table], string(encoded))
		}
		p22FaultOK(t, rows.Err())
		p22FaultOK(t, rows.Close())
		sort.Strings(result[table])
	}
	return result
}

func p22FaultOK(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatal(err)
	}
}

func p22ArchiveEntries(t *testing.T, path string) map[string][]byte {
	t.Helper()
	reader, err := zip.OpenReader(path)
	if err != nil {
		t.Fatalf("open archive %s: %v", path, err)
	}
	defer func() { _ = reader.Close() }()
	files := append([]*zip.File(nil), reader.File...)
	sort.Slice(files, func(left, right int) bool { return files[left].Name < files[right].Name })
	entries := make(map[string][]byte, len(files))
	for _, file := range files {
		input, err := file.Open()
		if err != nil {
			t.Fatalf("open archive entry %s: %v", file.Name, err)
		}
		data, readErr := io.ReadAll(input)
		closeErr := input.Close()
		if readErr != nil || closeErr != nil {
			t.Fatalf("read archive entry %s: read=%v close=%v", file.Name, readErr, closeErr)
		}
		entries[file.Name] = data
	}
	return entries
}

func p23PortableAllRows(t *testing.T, conn *sql.DB) map[string][]string {
	t.Helper()
	names, err := conn.Query(`SELECT name FROM sqlite_master WHERE type='table' AND name NOT LIKE 'sqlite_%' ORDER BY name`)
	p22FaultOK(t, err)
	var tables []string
	for names.Next() {
		var table string
		p22FaultOK(t, names.Scan(&table))
		tables = append(tables, table)
	}
	p22FaultOK(t, names.Err())
	p22FaultOK(t, names.Close())
	result := make(map[string][]string, len(tables))
	for _, table := range tables {
		result[table] = p23PortableTableRows(t, conn, table)
	}
	return result
}

func p23PortableTableRows(t *testing.T, conn *sql.DB, table string) []string {
	t.Helper()
	rows, err := conn.Query(`SELECT * FROM "` + table + `"`)
	p22FaultOK(t, err)
	defer func() { p22FaultOK(t, rows.Close()) }()
	columns, err := rows.Columns()
	p22FaultOK(t, err)
	result := []string{}
	for rows.Next() {
		values := make([]any, len(columns))
		pointers := make([]any, len(columns))
		for i := range values {
			pointers[i] = &values[i]
		}
		p22FaultOK(t, rows.Scan(pointers...))
		encoded, err := json.Marshal(values)
		p22FaultOK(t, err)
		result = append(result, string(encoded))
	}
	p22FaultOK(t, rows.Err())
	sort.Strings(result)
	return result
}

func p23PortableTreeDigest(t *testing.T, root string) string {
	t.Helper()
	hash := sha256.New()
	err := filepath.WalkDir(root, func(path string, item fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if item.IsDir() {
			return nil
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		_, _ = hash.Write([]byte(strings.TrimPrefix(path, root)))
		_, _ = hash.Write(data)
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	return fmt.Sprintf("%x", hash.Sum(nil))
}

func p23PortableAssertRows(t *testing.T, phase string, before, after map[string][]string) {
	t.Helper()
	if len(before) != len(after) {
		t.Errorf("%s table inventory changed: live=%d current=%d", phase, len(before), len(after))
	}
	for table := range before {
		if !reflect.DeepEqual(before[table], after[table]) {
			t.Errorf("%s %s rows changed: live=%v current=%v", phase, table, before[table], after[table])
		}
	}
}
