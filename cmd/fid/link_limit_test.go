package main

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"testing"

	"github.com/mbell-nzl/fid/internal/config"
	fiddb "github.com/mbell-nzl/fid/internal/db"
)

func TestLinkListJSONHonorsAtAndOverLimitWithoutTruncation(t *testing.T) {
	t.Parallel()
	workspace := t.TempDir()
	clonePristineWorkspaceInto(t, workspace)
	primaryID := createCommandIssue(t, workspace, "Link list limit primary")
	journalPath := mustJournalPath(t, workspace)
	conn, err := fiddb.Open(filepath.Join(config.LedgerDir(workspace), "fid.sqlite"))
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	defer closeCommandDB(t, conn)

	seedTx, err := conn.Begin()
	if err != nil {
		t.Fatalf("begin link limit fixture transaction: %v", err)
	}
	defer func() { _ = seedTx.Rollback() }()
	issueInsert, err := seedTx.Prepare(`
		INSERT INTO issues (id, kind, title, workflow_id, status)
		VALUES (?, 'task', ?, 'task-default', 'todo')
	`)
	if err != nil {
		t.Fatalf("prepare issues insert: %v", err)
	}
	defer func() { _ = issueInsert.Close() }()
	linkInsert, err := seedTx.Prepare(`
		INSERT INTO links (from_type, from_id, relation, to_type, to_id, created_at)
		VALUES ('issue', ?, 'relates_to', 'issue', ?, '2026-08-15T00:00:00Z')
	`)
	if err != nil {
		t.Fatalf("prepare links insert: %v", err)
	}
	defer func() { _ = linkInsert.Close() }()
	for index := 0; index < 10000; index++ {
		issueID := fmt.Sprintf("FID-link-limit-%05d", index)
		if _, err := issueInsert.Exec(issueID, "Issue "+issueID); err != nil {
			t.Fatalf("insert issue %d: %v", index, err)
		}
		if _, err := linkInsert.Exec(primaryID, issueID); err != nil {
			t.Fatalf("insert link %d: %v", index, err)
		}
	}
	if err := seedTx.Commit(); err != nil {
		t.Fatalf("commit link limit fixture transaction: %v", err)
	}

	atLimit := executeCommandForTest(workspace, "--json", "link", "list", primaryID, "--type", "issue")
	if atLimit.err != nil || atLimit.stderr != "" {
		t.Fatalf("at-limit link.list = %#v, want JSON success on stdout only", atLimit)
	}
	envelope := decodeCommandJSONStructure(t, atLimit.stdout).(map[string]any)
	result := envelope["result"].(map[string]any)
	if envelope["event"] != nil || result["count"] != json.Number("10000") || len(result["links"].([]any)) != 10000 {
		t.Fatalf("at-limit link.list envelope = %#v, want 10000 complete rows and null event", envelope)
	}

	if _, err := conn.Exec(`INSERT INTO issues (id, kind, title, workflow_id, status) VALUES (?, 'task', ?, 'task-default', 'todo')`, "FID-link-limit-over", "Over-limit issue"); err != nil {
		t.Fatalf("insert over-limit issue: %v", err)
	}
	if _, err := conn.Exec(`INSERT INTO links (from_type, from_id, relation, to_type, to_id, created_at) VALUES ('issue', ?, 'relates_to', 'issue', ?, '2026-08-15T00:00:00Z')`, primaryID, "FID-link-limit-over"); err != nil {
		t.Fatalf("insert over-limit link: %v", err)
	}
	before := len(splitJournalLines(mustReadCommandFile(t, journalPath)))
	overLimit := executeCommandForTest(workspace, "--json", "link", "list", primaryID, "--type", "issue")
	if overLimit.err == nil || overLimit.stdout != "" {
		t.Fatalf("over-limit link.list = %#v, want stderr-only rejection", overLimit)
	}
	assertLinkFailure(t, overLimit.stderr, "link.list", "result-limit-exceeded", nil, map[string]any{
		"collection": "links", "limit": json.Number("10000"), "observed": json.Number("10001"),
	})
	if got := len(splitJournalLines(mustReadCommandFile(t, journalPath))); got != before {
		t.Fatalf("journal lines after over-limit read = %d, want %d", got, before)
	}
	var rows int
	if err := conn.QueryRow(`SELECT COUNT(*) FROM links WHERE from_id = ?`, primaryID).Scan(&rows); err != nil {
		t.Fatalf("count links after over-limit read: %v", err)
	}
	if rows != 10001 {
		t.Fatalf("links after over-limit read = %d, want 10001 (no truncation/mutation)", rows)
	}
}
