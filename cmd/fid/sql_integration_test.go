package main

import (
	"path/filepath"
	"testing"

	fiddb "github.com/mbell-nzl/fid/internal/db"
)

func TestSQLCommandRejectsMutationAttemptsWithoutChangingProjectedCounts(t *testing.T) {
	t.Parallel()
	workspace := t.TempDir()
	clonePristineWorkspaceInto(t, workspace)

	_, createStderr := runFIDCommand(
		t,
		workspace,
		"issue", "create",
		"--kind", "task",
		"--title", "SQL integration issue",
		"--by", "developer",
	)
	issueID := parseCreatedIssueID(t, createStderr)
	if issueID == "" {
		t.Fatal("issue id = empty, want created issue")
	}

	before := countCommandTableRows(t, workspace, "issues")
	commands := []string{
		"INSERT INTO issues(id, kind, title, status, workflow_id) VALUES('FID-sql', 'task', 'Injected', 'todo', 'default_task')",
		"UPDATE issues SET title = 'changed'",
		"DELETE FROM issues",
		"DROP TABLE issues",
	}

	for _, query := range commands {
		t.Run(query, func(t *testing.T) {
			stdout, stderr, err := runFIDCommandExpectError(t, workspace, "sql", query)
			if err == nil {
				t.Fatalf("Execute() error = nil, want rejection for %q", query)
			}
			if stdout != "" {
				t.Fatalf("stdout = %q, want empty", stdout)
			}
			if stderr == "" {
				t.Fatalf("stderr = empty, want rejection for %q", query)
			}

			after := countCommandTableRows(t, workspace, "issues")
			if after != before {
				t.Fatalf("issues row count after %q = %d, want %d", query, after, before)
			}
		})
	}
}

func countCommandTableRows(t *testing.T, workspace, table string) int {
	t.Helper()

	conn, err := fiddb.Open(filepath.Join(workspace, ".fid", "fid.sqlite"))
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	defer closeCommandDB(t, conn)

	var count int
	query := "SELECT COUNT(*) FROM " + table
	if err := conn.QueryRow(query).Scan(&count); err != nil {
		t.Fatalf("QueryRow(%q) error = %v", query, err)
	}
	return count
}
