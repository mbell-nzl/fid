package requirements

import (
	"database/sql"
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"testing"

	"github.com/mbell-nzl/fid/internal/bootstrap"
	fiddb "github.com/mbell-nzl/fid/internal/db"
	"github.com/mbell-nzl/fid/internal/policy"
)

func TestSyncCreatesThenSkipsUnchangedThenUpdatesChangedRows(t *testing.T) {
	t.Parallel()

	workspace := newRequirementsWorkspace(t)
	sourceDoc := writeRequirementsDoc(t, workspace, "project-docs/specs/sample.md", `# Sample

## 3.11 Governance

| ID | Requirement | Priority | Notes |
|---|---|---|---|
| PRD-GOV-006 | The system shall support requirement pointer records. | Must | v0.1 |
| TRD-ARCH-010 | Every mutation flows through the journal write path. | Must | v0.1 |
`)

	beforeEvents := countRequirementJournalEvents(t, workspace)

	first, err := Sync(SyncOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testRequirementsMigrationsDir(t),
		SourceDoc:     sourceDoc,
		ActorRole:     "tl",
	})
	if err != nil {
		t.Fatalf("Sync(first) error = %v", err)
	}
	if first.Created != 2 || first.Updated != 0 || first.Unchanged != 0 {
		t.Fatalf("first counts = %#v, want created=2 updated=0 unchanged=0", first)
	}
	if got := countRequirementJournalEvents(t, workspace); got != beforeEvents+1 {
		t.Fatalf("journal events after first sync = %d, want %d", got, beforeEvents+1)
	}

	conn := openRequirementsDB(t, workspace)
	defer closeRequirementsDB(t, conn)

	firstRows := loadRequirementRowsFromDB(t, conn)
	if got := requirementIDs(firstRows); !equalRequirementIDs(got, []string{"PRD-GOV-006", "TRD-ARCH-010"}) {
		t.Fatalf("first synced ids = %#v, want %#v", got, []string{"PRD-GOV-006", "TRD-ARCH-010"})
	}
	firstHash := firstRows["PRD-GOV-006"].LastSyncedHash

	second, err := Sync(SyncOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testRequirementsMigrationsDir(t),
		SourceDoc:     sourceDoc,
		ActorRole:     "tl",
	})
	if err != nil {
		t.Fatalf("Sync(second) error = %v", err)
	}
	if second.Created != 0 || second.Updated != 0 || second.Unchanged != 2 {
		t.Fatalf("second counts = %#v, want created=0 updated=0 unchanged=2", second)
	}
	if got := countRequirementJournalEvents(t, workspace); got != beforeEvents+1 {
		t.Fatalf("journal events after unchanged sync = %d, want %d", got, beforeEvents+1)
	}

	writeRequirementsDoc(t, workspace, sourceDoc, `# Sample

## 3.11 Governance

| ID | Requirement | Priority | Notes |
|---|---|---|---|
| PRD-GOV-006 | The system shall support requirement pointers with a drift hash. | Must | v0.1 |
| TRD-ARCH-010 | Every mutation flows through the journal write path. | Must | v0.1 |
`)

	third, err := Sync(SyncOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testRequirementsMigrationsDir(t),
		SourceDoc:     sourceDoc,
		ActorRole:     "tl",
	})
	if err != nil {
		t.Fatalf("Sync(third) error = %v", err)
	}
	if third.Created != 0 || third.Updated != 1 || third.Unchanged != 1 {
		t.Fatalf("third counts = %#v, want created=0 updated=1 unchanged=1", third)
	}
	if got := countRequirementJournalEvents(t, workspace); got != beforeEvents+2 {
		t.Fatalf("journal events after update sync = %d, want %d", got, beforeEvents+2)
	}

	updatedRows := loadRequirementRowsFromDB(t, conn)
	if updatedRows["PRD-GOV-006"].Summary != "The system shall support requirement pointers with a drift hash." {
		t.Fatalf("updated summary = %q, want changed summary", updatedRows["PRD-GOV-006"].Summary)
	}
	if updatedRows["PRD-GOV-006"].LastSyncedHash == firstHash {
		t.Fatal("last_synced_hash unchanged after summary update, want changed hash")
	}
}

func TestSyncSupportsTypeAndIDFilters(t *testing.T) {
	t.Parallel()

	t.Run("type filter", func(t *testing.T) {
		t.Parallel()

		workspace := newRequirementsWorkspace(t)
		sourceDoc := writeRequirementsDoc(t, workspace, "project-docs/specs/mixed.md", `# Sample

## 3.11 Governance

| ID | Requirement | Priority | Notes |
|---|---|---|---|
| PRD-GOV-006 | The system shall support requirement pointer records. | Must | v0.1 |
| PRD-LINK-002 | The system shall support links from issues to requirements. | Must | v0.1 |
| TRD-ARCH-010 | Every mutation flows through the journal write path. | Must | v0.1 |
`)

		result, err := Sync(SyncOptions{
			WorkspaceRoot: workspace,
			MigrationsDir: testRequirementsMigrationsDir(t),
			SourceDoc:     sourceDoc,
			ActorRole:     "human",
			TypeFilter:    "PRD",
		})
		if err != nil {
			t.Fatalf("Sync(type filter) error = %v", err)
		}
		if result.Created != 2 || result.Updated != 0 || result.Unchanged != 0 {
			t.Fatalf("type-filter counts = %#v, want created=2 updated=0 unchanged=0", result)
		}

		conn := openRequirementsDB(t, workspace)
		defer closeRequirementsDB(t, conn)
		rows := loadRequirementRowsFromDB(t, conn)
		if got := requirementIDs(rows); !equalRequirementIDs(got, []string{"PRD-GOV-006", "PRD-LINK-002"}) {
			t.Fatalf("type-filter ids = %#v, want %#v", got, []string{"PRD-GOV-006", "PRD-LINK-002"})
		}
	})

	t.Run("single id", func(t *testing.T) {
		t.Parallel()

		workspace := newRequirementsWorkspace(t)
		sourceDoc := writeRequirementsDoc(t, workspace, "project-docs/specs/mixed.md", `# Sample

## 4.3 Persistence Rules

| ID | Requirement | Priority | Notes |
|---|---|---|---|
| PRD-GOV-006 | The system shall support requirement pointer records. | Must | v0.1 |
| TRD-ARCH-010 | Every mutation flows through the journal write path. | Must | v0.1 |
`)

		result, err := Sync(SyncOptions{
			WorkspaceRoot: workspace,
			MigrationsDir: testRequirementsMigrationsDir(t),
			SourceDoc:     sourceDoc,
			ActorRole:     "tl",
			RequirementID: "TRD-ARCH-010",
		})
		if err != nil {
			t.Fatalf("Sync(single id) error = %v", err)
		}
		if result.Created != 1 || result.Updated != 0 || result.Unchanged != 0 {
			t.Fatalf("single-id counts = %#v, want created=1 updated=0 unchanged=0", result)
		}

		conn := openRequirementsDB(t, workspace)
		defer closeRequirementsDB(t, conn)
		rows := loadRequirementRowsFromDB(t, conn)
		if got := requirementIDs(rows); !equalRequirementIDs(got, []string{"TRD-ARCH-010"}) {
			t.Fatalf("single-id ids = %#v, want %#v", got, []string{"TRD-ARCH-010"})
		}
	})
}

func TestSyncRejectsMalformedInputAndDeniedRolesBeforeAppend(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		sourceBody  string
		typeFilter  string
		requirement string
		role        string
		assertError func(*testing.T, error)
	}{
		{
			name: "no parseable ids",
			sourceBody: `# Sample

## Notes

This file has no canonical requirement rows.
`,
			role: "tl",
			assertError: func(t *testing.T, err error) {
				t.Helper()
				var target *NoParseableIDsError
				if !errors.As(err, &target) {
					t.Fatalf("error = %T, want *NoParseableIDsError", err)
				}
			},
		},
		{
			name: "duplicate ids",
			sourceBody: `# Sample

## 3.11 Governance

| ID | Requirement | Priority | Notes |
|---|---|---|---|
| PRD-GOV-006 | First pointer requirement. | Must | v0.1 |
| PRD-GOV-006 | Duplicate pointer requirement. | Must | v0.1 |
`,
			role: "tl",
			assertError: func(t *testing.T, err error) {
				t.Helper()
				var target *DuplicateIDError
				if !errors.As(err, &target) {
					t.Fatalf("error = %T, want *DuplicateIDError", err)
				}
			},
		},
		{
			name: "unknown type",
			sourceBody: `# Sample

## 3.11 Governance

| ID | Requirement | Priority | Notes |
|---|---|---|---|
| PRD-GOV-006 | Pointer requirement. | Must | v0.1 |
`,
			typeFilter: "OPS",
			role:       "tl",
			assertError: func(t *testing.T, err error) {
				t.Helper()
				var target *UnknownTypeError
				if !errors.As(err, &target) {
					t.Fatalf("error = %T, want *UnknownTypeError", err)
				}
			},
		},
		{
			name: "missing id from source",
			sourceBody: `# Sample

## 3.11 Governance

| ID | Requirement | Priority | Notes |
|---|---|---|---|
| PRD-GOV-006 | Pointer requirement. | Must | v0.1 |
`,
			requirement: "TRD-ARCH-010",
			role:        "tl",
			assertError: func(t *testing.T, err error) {
				t.Helper()
				var target *RequirementNotFoundInSourceError
				if !errors.As(err, &target) {
					t.Fatalf("error = %T, want *RequirementNotFoundInSourceError", err)
				}
			},
		},
		{
			name: "denied role",
			sourceBody: `# Sample

## 3.11 Governance

| ID | Requirement | Priority | Notes |
|---|---|---|---|
| PRD-GOV-006 | Pointer requirement. | Must | v0.1 |
`,
			role: "developer",
			assertError: func(t *testing.T, err error) {
				t.Helper()
				var target *policy.PermissionDeniedError
				if !errors.As(err, &target) {
					t.Fatalf("error = %T, want *PermissionDeniedError", err)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			workspace := newRequirementsWorkspace(t)
			sourceDoc := writeRequirementsDoc(t, workspace, "project-docs/specs/bad.md", tt.sourceBody)
			beforeEvents := countRequirementJournalEvents(t, workspace)

			_, err := Sync(SyncOptions{
				WorkspaceRoot: workspace,
				MigrationsDir: testRequirementsMigrationsDir(t),
				SourceDoc:     sourceDoc,
				ActorRole:     tt.role,
				TypeFilter:    tt.typeFilter,
				RequirementID: tt.requirement,
			})
			if err == nil {
				t.Fatal("Sync() error = nil, want rejection")
			}
			tt.assertError(t, err)
			if got := countRequirementJournalEvents(t, workspace); got != beforeEvents {
				t.Fatalf("journal events after rejection = %d, want %d", got, beforeEvents)
			}

			conn := openRequirementsDB(t, workspace)
			defer closeRequirementsDB(t, conn)
			if got := countRequirementRows(t, conn); got != 0 {
				t.Fatalf("requirements row count after rejection = %d, want 0", got)
			}
		})
	}
}

func newRequirementsWorkspace(t *testing.T) string {
	t.Helper()

	workspace := t.TempDir()
	if _, err := bootstrap.InitializeWorkspace(bootstrap.InitOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testRequirementsMigrationsDir(t),
	}); err != nil {
		t.Fatalf("InitializeWorkspace() error = %v", err)
	}
	return workspace
}

func testRequirementsMigrationsDir(t *testing.T) string {
	t.Helper()

	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller(0) failed")
	}
	return filepath.Clean(filepath.Join(filepath.Dir(currentFile), "..", "..", "migrations"))
}

func writeRequirementsDoc(t *testing.T, workspace, relativePath, body string) string {
	t.Helper()

	var path string
	if filepath.IsAbs(relativePath) {
		path = relativePath
	} else {
		path = filepath.Join(workspace, relativePath)
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("MkdirAll(%q) error = %v", filepath.Dir(path), err)
	}
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatalf("WriteFile(%q) error = %v", path, err)
	}
	return filepath.ToSlash(relativePath)
}

func countRequirementJournalEvents(t *testing.T, workspace string) int {
	t.Helper()

	entries, err := os.ReadDir(filepath.Join(workspace, ".fid", "journal"))
	if err != nil {
		t.Fatalf("ReadDir(journal) error = %v", err)
	}
	total := 0
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".jsonl") {
			continue
		}
		payload, err := os.ReadFile(filepath.Join(workspace, ".fid", "journal", entry.Name()))
		if err != nil {
			t.Fatalf("ReadFile(%q) error = %v", entry.Name(), err)
		}
		trimmed := strings.TrimSpace(string(payload))
		if trimmed == "" {
			continue
		}
		total += len(strings.Split(trimmed, "\n"))
	}
	return total
}

func openRequirementsDB(t *testing.T, workspace string) *sql.DB {
	t.Helper()

	conn, err := fiddb.Open(filepath.Join(workspace, ".fid", "fid.sqlite"))
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	return conn
}

func closeRequirementsDB(t *testing.T, conn *sql.DB) {
	t.Helper()
	if err := conn.Close(); err != nil {
		t.Fatalf("Close() error = %v", err)
	}
}

func loadRequirementRowsFromDB(t *testing.T, conn *sql.DB) map[string]requirementRow {
	t.Helper()

	rows, err := conn.Query(`
		SELECT id, type, summary, source_doc, source_section, COALESCE(last_synced_hash, ''), COALESCE(status, '')
		FROM requirements
	`)
	if err != nil {
		t.Fatalf("Query(requirements) error = %v", err)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			t.Fatalf("Rows.Close() error = %v", err)
		}
	}()

	items := make(map[string]requirementRow)
	for rows.Next() {
		var (
			row           requirementRow
			sourceSection sql.NullString
		)
		if err := rows.Scan(&row.ID, &row.Type, &row.Summary, &row.SourceDoc, &sourceSection, &row.LastSyncedHash, &row.Status); err != nil {
			t.Fatalf("Scan(requirements) error = %v", err)
		}
		row.SourceSection = nullRequirementString(sourceSection)
		items[row.ID] = row
	}
	if err := rows.Err(); err != nil {
		t.Fatalf("Rows(requirements) error = %v", err)
	}
	return items
}

func countRequirementRows(t *testing.T, conn *sql.DB) int {
	t.Helper()

	var count int
	if err := conn.QueryRow(`SELECT COUNT(*) FROM requirements`).Scan(&count); err != nil {
		t.Fatalf("QueryRow(count requirements) error = %v", err)
	}
	return count
}

func requirementIDs(rows map[string]requirementRow) []string {
	ids := make([]string, 0, len(rows))
	for id := range rows {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	return ids
}

func equalRequirementIDs(left, right []string) bool {
	if len(left) != len(right) {
		return false
	}
	for index := range left {
		if left[index] != right[index] {
			return false
		}
	}
	return true
}
