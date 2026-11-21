package main

import (
	"path/filepath"
	"testing"

	"github.com/mbell-nzl/fid/internal/config"
	fiddb "github.com/mbell-nzl/fid/internal/db"
)

// setIssueStatus puts a task or an epic in a state through the projection.
// Only fid start, the hand-off, a review round asking for changes, fid done
// and a person's cancellation move one, and admission refuses any other
// Event that tries, creation outside todo included (RESET 1.9b). Fixtures
// that need a state those commands don't reach, or that seed read rows, set
// it here, as the start and hand-off refusal fixtures set theirs.
func setIssueStatus(t *testing.T, root, id, status string) {
	t.Helper()
	conn, err := fiddb.Open(filepath.Join(config.LedgerDir(root), config.DatabaseFileName))
	if err != nil {
		t.Fatal(err)
	}
	defer closeCommandDB(t, conn)
	if _, err := conn.Exec(`UPDATE issues SET status = ? WHERE id = ?`, status, id); err != nil {
		t.Fatal(err)
	}
}

// seedReviewState puts task in review without a hand-off, a state no
// command reaches, so approve's own hand-off check stays covered.
func seedReviewState(t *testing.T, root, task string) {
	t.Helper()
	setIssueStatus(t, root, task, "in_review")
}
