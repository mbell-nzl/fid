package links

import (
	"database/sql"
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/mbell-nzl/fid/internal/bootstrap"
	"github.com/mbell-nzl/fid/internal/config"
	fiddb "github.com/mbell-nzl/fid/internal/db"
	"github.com/mbell-nzl/fid/internal/issues"
)

func TestRelationInverseMapMatchesIssueOnlySet(t *testing.T) {
	t.Parallel()

	if got, want := SupportedRelations(), []string{
		"blocks",
		"blocked_by",
		"relates_to",
		"duplicates",
		"duplicated_by",
		"supersedes",
		"superseded_by",
		"caused_by",
		"causes",
		"discovered_from",
		"discovered",
	}; !equalStrings(got, want) {
		t.Fatalf("SupportedRelations() = %#v, want %#v", got, want)
	}

	want := map[string]string{
		"blocks":          "blocked_by",
		"blocked_by":      "blocks",
		"relates_to":      "relates_to",
		"duplicates":      "duplicated_by",
		"duplicated_by":   "duplicates",
		"supersedes":      "superseded_by",
		"superseded_by":   "supersedes",
		"caused_by":       "causes",
		"causes":          "caused_by",
		"discovered_from": "discovered",
		"discovered":      "discovered_from",
	}

	if got := RelationInverseMap(); !mapsEqual(got, want) {
		t.Fatalf("RelationInverseMap() = %#v, want %#v", got, want)
	}
}

func TestAddSupportsInverseDeclaredRelationAndRejectsRequirementEndpoints(t *testing.T) {
	t.Parallel()

	t.Run("declared inverse relation", func(t *testing.T) {
		workspace := newLinksWorkspace(t)
		firstIssueID := createLinksIssue(t, workspace, "First issue")
		secondIssueID := createLinksIssue(t, workspace, "Second issue")

		beforeEvents := countLinksJournalEvents(t, workspace)

		result, err := Add(AddOptions{
			WorkspaceRoot: workspace,
			MigrationsDir: testLinksMigrationsDir(t),
			FromID:        secondIssueID,
			Relation:      "blocked_by",
			ToID:          firstIssueID,
			Actor:         testActorForRole("developer"), IdempotencyKey: testIdempotencyKey(),
		})
		if err != nil {
			t.Fatalf("Add(blocked_by) error = %v", err)
		}
		if got := countLinksJournalEvents(t, workspace); got != beforeEvents+1 {
			t.Fatalf("journal events after Add(blocked_by) = %d, want %d", got, beforeEvents+1)
		}
		if result.Declared.Relation != "blocked_by" || result.Inverse.Relation != "blocks" {
			t.Fatalf("Add(blocked_by) = %#v, want declared blocked_by + inverse blocks", result)
		}

		conn := openLinksDB(t, workspace)
		defer closeLinksDB(t, conn)

		rows := loadLinkRows(t, conn)
		if !containsLink(rows, Link{FromType: "issue", FromID: secondIssueID, Relation: "blocked_by", ToType: "issue", ToID: firstIssueID}) {
			t.Fatalf("links rows = %#v, want declared blocked_by row", rows)
		}
		if !containsLink(rows, Link{FromType: "issue", FromID: firstIssueID, Relation: "blocks", ToType: "issue", ToID: secondIssueID}) {
			t.Fatalf("links rows = %#v, want inverse blocks row", rows)
		}
	})

	t.Run("issue implements requirement", func(t *testing.T) {
		workspace := newLinksWorkspace(t)
		issueID := createLinksIssue(t, workspace, "Implements issue")
		before := countLinksJournalEvents(t, workspace)
		_, err := Add(AddOptions{
			WorkspaceRoot: workspace,
			MigrationsDir: testLinksMigrationsDir(t),
			FromID:        issueID,
			Relation:      "implements",
			ToType:        "requirement",
			ToID:          "PRD-LINK-002",
			Actor:         testActorForRole("developer"), IdempotencyKey: testIdempotencyKey(),
		})
		var invalid *InvalidEndpointTypeError
		if !errors.As(err, &invalid) || countLinksJournalEvents(t, workspace) != before {
			t.Fatalf("removed requirement endpoint = %v; Event count %d -> %d", err, before, countLinksJournalEvents(t, workspace))
		}
	})
}

func TestAddRejectsInvalidRequestsBeforeAppend(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		setup func(t *testing.T, workspace string) (AddOptions, func(*testing.T, error))
	}{
		{
			name: "duplicate link",
			setup: func(t *testing.T, workspace string) (AddOptions, func(*testing.T, error)) {
				firstIssueID := createLinksIssue(t, workspace, "First issue")
				secondIssueID := createLinksIssue(t, workspace, "Second issue")
				if _, err := Add(AddOptions{
					WorkspaceRoot: workspace,
					MigrationsDir: testLinksMigrationsDir(t),
					FromID:        firstIssueID,
					Relation:      "blocks",
					ToID:          secondIssueID,
					Actor:         testActorForRole("developer"), IdempotencyKey: testIdempotencyKey(),
				}); err != nil {
					t.Fatalf("seed Add(blocks) error = %v", err)
				}
				return AddOptions{
						WorkspaceRoot: workspace,
						MigrationsDir: testLinksMigrationsDir(t),
						FromID:        firstIssueID,
						Relation:      "blocks",
						ToID:          secondIssueID,
						Actor:         testActorForRole("developer"), IdempotencyKey: testIdempotencyKey(),
					}, func(t *testing.T, err error) {
						t.Helper()
						var target *DuplicateLinkError
						if !errors.As(err, &target) {
							t.Fatalf("error = %T, want *DuplicateLinkError", err)
						}
					}
			},
		},
		{
			name: "unknown relation",
			setup: func(t *testing.T, workspace string) (AddOptions, func(*testing.T, error)) {
				firstIssueID := createLinksIssue(t, workspace, "First issue")
				secondIssueID := createLinksIssue(t, workspace, "Second issue")
				return AddOptions{
						WorkspaceRoot: workspace,
						MigrationsDir: testLinksMigrationsDir(t),
						FromID:        firstIssueID,
						Relation:      "parent_of",
						ToID:          secondIssueID,
						Actor:         testActorForRole("developer"), IdempotencyKey: testIdempotencyKey(),
					}, func(t *testing.T, err error) {
						t.Helper()
						var target *UnknownRelationError
						if !errors.As(err, &target) {
							t.Fatalf("error = %T, want *UnknownRelationError", err)
						}
					}
			},
		},
		{
			name: "self link",
			setup: func(t *testing.T, workspace string) (AddOptions, func(*testing.T, error)) {
				issueID := createLinksIssue(t, workspace, "Self issue")
				return AddOptions{
						WorkspaceRoot: workspace,
						MigrationsDir: testLinksMigrationsDir(t),
						FromID:        issueID,
						Relation:      "relates_to",
						ToID:          issueID,
						Actor:         testActorForRole("developer"), IdempotencyKey: testIdempotencyKey(),
					}, func(t *testing.T, err error) {
						t.Helper()
						var target *SelfLinkError
						if !errors.As(err, &target) {
							t.Fatalf("error = %T, want *SelfLinkError", err)
						}
					}
			},
		},
		{
			name: "incompatible relation",
			setup: func(t *testing.T, workspace string) (AddOptions, func(*testing.T, error)) {
				firstIssueID := createLinksIssue(t, workspace, "First issue")
				secondIssueID := createLinksIssue(t, workspace, "Second issue")
				return AddOptions{
						WorkspaceRoot: workspace,
						MigrationsDir: testLinksMigrationsDir(t),
						FromID:        firstIssueID,
						Relation:      "implements",
						ToID:          secondIssueID,
						Actor:         testActorForRole("developer"), IdempotencyKey: testIdempotencyKey(),
					}, func(t *testing.T, err error) {
						t.Helper()
						var target *UnknownRelationError
						if !errors.As(err, &target) {
							t.Fatalf("error = %T, want *UnknownRelationError", err)
						}
					}
			},
		},
		{
			name: "missing requirement endpoint",
			setup: func(t *testing.T, workspace string) (AddOptions, func(*testing.T, error)) {
				issueID := createLinksIssue(t, workspace, "Implements issue")
				return AddOptions{
						WorkspaceRoot: workspace,
						MigrationsDir: testLinksMigrationsDir(t),
						FromType:      "issue",
						FromID:        issueID,
						Relation:      "implements",
						ToType:        "requirement",
						ToID:          "PRD-LINK-002",
						Actor:         testActorForRole("developer"), IdempotencyKey: testIdempotencyKey(),
					}, func(t *testing.T, err error) {
						t.Helper()
						var target *InvalidEndpointTypeError
						if !errors.As(err, &target) {
							t.Fatalf("error = %T, want *InvalidEndpointTypeError", err)
						}
					}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			workspace := newLinksWorkspace(t)
			options, assertError := tt.setup(t, workspace)
			beforeEvents := countLinksJournalEvents(t, workspace)

			_, err := Add(options)
			if err == nil {
				t.Fatal("Add() error = nil, want rejection")
			}
			assertError(t, err)

			if got := countLinksJournalEvents(t, workspace); got != beforeEvents {
				t.Fatalf("journal events after rejection = %d, want %d", got, beforeEvents)
			}
		})
	}
}

func newLinksWorkspace(t *testing.T) string {
	t.Helper()

	workspace := t.TempDir()
	if _, err := bootstrap.InitializeWorkspace(bootstrap.InitOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testLinksMigrationsDir(t),
	}); err != nil {
		t.Fatalf("InitializeWorkspace() error = %v", err)
	}
	return workspace
}

func testLinksMigrationsDir(t *testing.T) string {
	t.Helper()

	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller(0) failed")
	}
	return filepath.Clean(filepath.Join(filepath.Dir(currentFile), "..", "..", "migrations"))
}

func createLinksIssue(t *testing.T, workspace, title string) string {
	t.Helper()

	result, err := issues.Create(issues.CreateOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testLinksMigrationsDir(t),
		Kind:          "task",
		Title:         title,
		Actor:         testActorForRole("developer"), IdempotencyKey: testIdempotencyKey(),
	})
	if err != nil {
		t.Fatalf("Create(%q) error = %v", title, err)
	}
	return result.ID
}

func countLinksJournalEvents(t *testing.T, workspace string) int {
	t.Helper()

	entries, err := os.ReadDir(filepath.Join(config.LedgerDir(workspace), "journal"))
	if err != nil {
		t.Fatalf("ReadDir(journal) error = %v", err)
	}
	total := 0
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".jsonl") {
			continue
		}
		payload, err := os.ReadFile(filepath.Join(config.LedgerDir(workspace), "journal", entry.Name()))
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

func openLinksDB(t *testing.T, workspace string) *sql.DB {
	t.Helper()

	conn, err := fiddb.Open(filepath.Join(config.LedgerDir(workspace), "fid.sqlite"))
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	return conn
}

func closeLinksDB(t *testing.T, conn *sql.DB) {
	t.Helper()
	if err := conn.Close(); err != nil {
		t.Fatalf("Close() error = %v", err)
	}
}

func loadLinkRows(t *testing.T, conn *sql.DB) []Link {
	t.Helper()

	rows, err := conn.Query(`
		SELECT from_type, from_id, relation, to_type, to_id
		FROM links
		ORDER BY from_type, from_id, relation, to_type, to_id
	`)
	if err != nil {
		t.Fatalf("Query(links) error = %v", err)
	}
	defer closeLinksRows(t, rows)

	items := []Link{}
	for rows.Next() {
		var item Link
		if err := rows.Scan(&item.FromType, &item.FromID, &item.Relation, &item.ToType, &item.ToID); err != nil {
			t.Fatalf("Scan(links) error = %v", err)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		t.Fatalf("Rows(links) error = %v", err)
	}
	return items
}

func closeLinksRows(t *testing.T, rows *sql.Rows) {
	t.Helper()
	if err := rows.Close(); err != nil {
		t.Fatalf("Rows.Close() error = %v", err)
	}
}

func containsLink(rows []Link, want Link) bool {
	for _, row := range rows {
		if row == want {
			return true
		}
	}
	return false
}

func mapsEqual(left, right map[string]string) bool {
	if len(left) != len(right) {
		return false
	}
	for key, value := range left {
		if right[key] != value {
			return false
		}
	}
	return true
}

func equalStrings(left, right []string) bool {
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
