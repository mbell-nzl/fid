package main

import (
	"database/sql"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/mbell-nzl/fid/internal/config"
	fiddb "github.com/mbell-nzl/fid/internal/db"
	"github.com/mbell-nzl/fid/internal/gittest"
)

func TestV01DLinkFlowAndRebuildParity(t *testing.T) {
	t.Parallel()
	workspace := gittest.NewWorkspace(t)
	runInitCommand(t, workspace, "init")

	primaryIssueID := createCommandIssue(t, workspace, "Primary issue")
	peerIssueID := createCommandIssue(t, workspace, "Add index on concerns(issue_id, status)")

	journalPath := mustJournalPath(t, workspace)
	linesBeforePhase12Writes := len(splitJournalLines(mustReadCommandFile(t, journalPath)))

	_, linkIssueStderr := runFIDCommand(t, workspace, "link", "add", peerIssueID, "discovered_from", primaryIssueID, "--from-type", "issue", "--to-type", "issue", "--by", "developer")
	assertNoNestedLockFailure(t, linkIssueStderr)
	if got := len(splitJournalLines(mustReadCommandFile(t, journalPath))); got != linesBeforePhase12Writes+1 {
		t.Fatalf("journal lines after issue link add = %d, want %d", got, linesBeforePhase12Writes+1)
	}

	_, linkRequirementStderr, linkRequirementErr := runFIDCommandExpectError(t, workspace, "link", "add", primaryIssueID, "implements", "PRD-LINK-002", "--from-type", "issue", "--to-type", "requirement", "--by", "developer")
	if linkRequirementErr == nil || !strings.Contains(linkRequirementErr.Error(), "to-type must be issue") ||
		linkRequirementStderr != "error: invalid link endpoint type\nreason: invalid-argument; argument=to-type, allowed=issue\n" {
		t.Fatalf("removed requirement endpoint = (%q, %v)", linkRequirementStderr, linkRequirementErr)
	}
	if got := len(splitJournalLines(mustReadCommandFile(t, journalPath))); got != linesBeforePhase12Writes+1 {
		t.Fatalf("journal lines after rejected requirement link = %d, want %d", got, linesBeforePhase12Writes+1)
	}

	primaryListStdout, primaryListStderr := runFIDCommand(t, workspace, "--json", "link", "list", primaryIssueID, "--type", "issue")
	if primaryListStderr != "" {
		t.Fatalf("primary list stderr = %q, want empty", primaryListStderr)
	}
	primaryList := decodeCommandJSONStructure(t, primaryListStdout).(map[string]any)["result"].(map[string]any)["links"].([]any)
	if len(primaryList) != 1 {
		t.Fatalf("primary link list = %#v, want 1 row", primaryList)
	}
	assertCommandLinkListRow(t, primaryList[0], "incoming", "discovered_from", "issue", peerIssueID, "Add index on concerns(issue_id, status)")

	_, _, requirementListErr := runFIDCommandExpectError(t, workspace, "--json", "link", "list", "PRD-LINK-002", "--type", "requirement")
	if requirementListErr == nil {
		t.Fatal("link list accepted removed requirement endpoint")
	}

	_, removeStderr := runFIDCommand(t, workspace, "link", "remove", peerIssueID, "discovered_from", primaryIssueID, "--from-type", "issue", "--to-type", "issue", "--by", "developer")
	assertNoNestedLockFailure(t, removeStderr)
	if got := len(splitJournalLines(mustReadCommandFile(t, journalPath))); got != linesBeforePhase12Writes+2 {
		t.Fatalf("journal lines after link remove = %d, want %d", got, linesBeforePhase12Writes+2)
	}

	beforeLinks := mustSnapshotE2ELinks(t, workspace)
	assertE2ERequirementsAbsent(t, workspace)

	dbPath := filepath.Join(config.LedgerDir(workspace), "fid.sqlite")
	if err := os.Remove(dbPath); err != nil {
		t.Fatalf("Remove(%q) error = %v", dbPath, err)
	}

	rebuildStdout, rebuildStderr := runFIDCommand(t, workspace, "rebuild")
	if !strings.Contains(rebuildStdout, "rebuild complete") || rebuildStderr != "" {
		t.Fatalf("rebuild stdout/stderr = %q/%q, want stdout completion and empty stderr", rebuildStdout, rebuildStderr)
	}

	if got := mustSnapshotE2ELinks(t, workspace); !reflect.DeepEqual(got, beforeLinks) {
		t.Fatalf("links after rebuild = %#v, want %#v", got, beforeLinks)
	}
	assertE2ERequirementsAbsent(t, workspace)
}

type e2eLinkRow struct {
	FromType string
	FromID   string
	Relation string
	ToType   string
	ToID     string
}

func assertNoNestedLockFailure(t *testing.T, stderr string) {
	t.Helper()

	lower := strings.ToLower(stderr)
	for _, blocked := range []string{"write.lock", "nested lock"} {
		if strings.Contains(lower, blocked) {
			t.Fatalf("stderr = %q, want no nested-lock failure", stderr)
		}
	}
}

func assertCommandLinkListRow(t *testing.T, value any, direction, relation, entityType, entityID, entityTitle string) {
	t.Helper()

	row, ok := value.(map[string]any)
	if !ok {
		t.Fatalf("row = %#v, want object", value)
	}
	want := map[string]string{
		"direction":    direction,
		"relation":     relation,
		"entity_type":  entityType,
		"entity_id":    entityID,
		"entity_title": entityTitle,
	}
	for key, wantValue := range want {
		if gotValue := strings.TrimSpace(row[key].(string)); gotValue != wantValue {
			t.Fatalf("row[%s] = %q, want %q", key, gotValue, wantValue)
		}
	}
}

func mustSnapshotE2ELinks(t *testing.T, workspace string) []e2eLinkRow {
	t.Helper()

	conn, err := fiddb.Open(filepath.Join(config.LedgerDir(workspace), "fid.sqlite"))
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	defer closeCommandDB(t, conn)

	rows, err := conn.Query(`
		SELECT from_type, from_id, relation, to_type, to_id
		FROM links
		ORDER BY from_type, from_id, relation, to_type, to_id
	`)
	if err != nil {
		t.Fatalf("query links snapshot: %v", err)
	}
	defer closeE2ERows(t, rows)

	var snapshot []e2eLinkRow
	for rows.Next() {
		var row e2eLinkRow
		if err := rows.Scan(&row.FromType, &row.FromID, &row.Relation, &row.ToType, &row.ToID); err != nil {
			t.Fatalf("scan links snapshot: %v", err)
		}
		snapshot = append(snapshot, row)
	}
	if err := rows.Err(); err != nil {
		t.Fatalf("iterate links snapshot: %v", err)
	}
	return snapshot
}

func assertE2ERequirementsAbsent(t *testing.T, workspace string) {
	t.Helper()
	conn, err := fiddb.Open(filepath.Join(config.LedgerDir(workspace), "fid.sqlite"))
	if err != nil {
		t.Fatal(err)
	}
	defer closeCommandDB(t, conn)
	var count int
	if err := conn.QueryRow(`SELECT count(*) FROM sqlite_schema WHERE type='table' AND name='requirement_refs'`).Scan(&count); err != nil || count != 0 {
		t.Fatalf("requirement table count = %d, %v; want absent", count, err)
	}
}

func closeE2ERows(t *testing.T, rows *sql.Rows) {
	t.Helper()
	if err := rows.Close(); err != nil {
		t.Fatalf("Rows.Close() error = %v", err)
	}
}
