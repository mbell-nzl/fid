package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	fiddb "github.com/mbell-nzl/fid/internal/db"
)

func TestV01DClosureLifecycleSmokeAndRebuildParity(t *testing.T) {
	workspace := t.TempDir()
	runInitCommand(t, workspace, "init")

	journalPath := mustJournalPath(t, workspace)
	expectedJournalLines := len(splitJournalLines(mustReadCommandFile(t, journalPath)))

	_, labelStderr := runFIDCommand(t, workspace, "label", "add", "lifecycle", "--description", "Lifecycle smoke label", "--by", "developer")
	assertNoNestedLockFailure(t, labelStderr)
	expectedJournalLines++
	assertClosureJournalLines(t, journalPath, expectedJournalLines, "label add")

	_, primaryCreateStderr := runFIDCommand(
		t,
		workspace,
		"issue", "create",
		"--kind", "task",
		"--title", "Primary lifecycle issue",
		"--actor", "agent:developer",
	)
	assertNoNestedLockFailure(t, primaryCreateStderr)
	primaryIssueID := parseCreatedIssueID(t, primaryCreateStderr)
	expectedJournalLines++
	assertClosureJournalLines(t, journalPath, expectedJournalLines, "primary issue create")

	_, peerCreateStderr := runFIDCommand(
		t,
		workspace,
		"issue", "create",
		"--kind", "task",
		"--title", "Peer lifecycle issue",
		"--actor", "agent:developer",
	)
	assertNoNestedLockFailure(t, peerCreateStderr)
	peerIssueID := parseCreatedIssueID(t, peerCreateStderr)
	expectedJournalLines++
	assertClosureJournalLines(t, journalPath, expectedJournalLines, "peer issue create")

	_, labelAddStderr := runFIDCommand(t, workspace, "issue", "label", "add", primaryIssueID, "lifecycle", "--by", "developer")
	assertNoNestedLockFailure(t, labelAddStderr)
	expectedJournalLines++
	assertClosureJournalLines(t, journalPath, expectedJournalLines, "issue label add")
	assertClosureIssueState(t, workspace, primaryIssueID, "todo", "unstarted", []string{"lifecycle"})

	_, transitionStderr := runFIDCommand(
		t,
		workspace,
		"issue", "transition", primaryIssueID,
		"--to", "in_progress",
		"--actor", "agent:qa_reviewer",
	)
	assertNoNestedLockFailure(t, transitionStderr)
	if !strings.Contains(transitionStderr, "transitioned "+primaryIssueID+": todo → in_progress") {
		t.Fatalf("transition stderr = %q, want in_progress confirmation", transitionStderr)
	}
	expectedJournalLines++
	assertClosureJournalLines(t, journalPath, expectedJournalLines, "standalone issue transition")
	if claims := mustReadIssueClaimsSection(t, workspace, primaryIssueID); len(claims) != 0 {
		t.Fatalf("claims after standalone transition = %#v, want empty", claims)
	}
	assertClosureIssueState(t, workspace, primaryIssueID, "in_progress", "started", []string{"lifecycle"})

	rejectStdout, rejectStderr, rejectErr := runFIDCommandExpectError(
		t,
		workspace,
		"issue", "transition", primaryIssueID,
		"--to", "dev_complete",
		"--actor", "agent:developer",
	)
	if rejectErr == nil {
		t.Fatal("developer transition without completion report error = nil, want workflow rejection")
	}
	if rejectStdout != "" {
		t.Fatalf("developer transition without completion report stdout = %q, want empty", rejectStdout)
	}
	assertNoNestedLockFailure(t, rejectStderr)
	if !strings.Contains(rejectStderr, "has_completion_report") {
		t.Fatalf("developer transition without completion report stderr = %q, want workflow precondition", rejectStderr)
	}
	assertClosureJournalLines(t, journalPath, expectedJournalLines, "developer transition without completion report")

	_, linkStderr := runFIDCommand(t, workspace, "link", "add", peerIssueID, "discovered_from", primaryIssueID, "--from-type", "issue", "--to-type", "issue", "--by", "developer")
	assertNoNestedLockFailure(t, linkStderr)
	expectedJournalLines++
	assertClosureJournalLines(t, journalPath, expectedJournalLines, "link add")

	linkListStdout, linkListStderr := runFIDCommand(t, workspace, "--json", "link", "list", primaryIssueID, "--type", "issue")
	if linkListStderr != "" {
		t.Fatalf("link list stderr = %q, want empty", linkListStderr)
	}
	linkList := decodeCommandJSONStructure(t, linkListStdout).(map[string]any)["result"].(map[string]any)["links"].([]any)
	if len(linkList) != 1 {
		t.Fatalf("link list = %#v, want 1 row", linkList)
	}
	assertCommandLinkListRow(t, linkList[0], "incoming", "discovered_from", "issue", peerIssueID, "Peer lifecycle issue")

	_, completeStderr := runFIDCommand(
		t,
		workspace,
		"issue", "complete", primaryIssueID,
		"--by", "developer",
		"--notes-md", "project-docs/docs/phases/phase-13/tasks/task-2-done.md",
		"--tokens", "1600",
	)
	assertNoNestedLockFailure(t, completeStderr)
	if !strings.Contains(completeStderr, "completed "+primaryIssueID+" (attempt 1)") {
		t.Fatalf("complete stderr = %q, want confirmation", completeStderr)
	}
	expectedJournalLines++
	assertClosureJournalLines(t, journalPath, expectedJournalLines, "issue complete")

	assertClosureIssueState(t, workspace, primaryIssueID, "in_progress", "started", []string{"lifecycle"})

	beforeIssues := mustSnapshotClosureIssues(t, workspace, primaryIssueID, peerIssueID)
	beforeClaims := mustSnapshotClosureClaims(t, workspace)
	beforeLinks := mustSnapshotE2ELinks(t, workspace)
	beforeIssueLabels := mustSnapshotClosureIssueLabels(t, workspace)

	dbPath := filepath.Join(workspace, ".fid", "fid.sqlite")
	if err := os.Remove(dbPath); err != nil {
		t.Fatalf("Remove(%q) error = %v", dbPath, err)
	}

	rebuildStdout, rebuildStderr := runFIDCommand(t, workspace, "rebuild")
	if !strings.Contains(rebuildStdout, "rebuild complete") || rebuildStderr != "" {
		t.Fatalf("rebuild stdout/stderr = %q/%q, want stdout completion and empty stderr", rebuildStdout, rebuildStderr)
	}

	if got := mustSnapshotClosureIssues(t, workspace, primaryIssueID, peerIssueID); !reflect.DeepEqual(got, beforeIssues) {
		t.Fatalf("issues after rebuild = %#v, want %#v", got, beforeIssues)
	}
	if got := mustSnapshotClosureClaims(t, workspace); !reflect.DeepEqual(got, beforeClaims) {
		t.Fatalf("claims after rebuild = %#v, want %#v", got, beforeClaims)
	}
	if got := mustSnapshotE2ELinks(t, workspace); !reflect.DeepEqual(got, beforeLinks) {
		t.Fatalf("links after rebuild = %#v, want %#v", got, beforeLinks)
	}
	if got := mustSnapshotClosureIssueLabels(t, workspace); !reflect.DeepEqual(got, beforeIssueLabels) {
		t.Fatalf("issue_labels after rebuild = %#v, want %#v", got, beforeIssueLabels)
	}
}

type closureIssueView struct {
	Status         string   `json:"status"`
	StatusCategory string   `json:"status_category"`
	Labels         []string `json:"labels"`
}

type e2eClosureIssueRow struct {
	ID         string
	Title      string
	Kind       string
	ParentID   string
	WorkflowID string
	Status     string
	OwnerRole  string
}

type e2eClosureClaimRow struct {
	RefType    string
	RefID      string
	ClaimedBy  string
	ClaimedAt  string
	TTLSeconds int
	Note       string
}

type e2eClosureIssueLabelRow struct {
	IssueID string
	Label   string
}

func assertClosureJournalLines(t *testing.T, journalPath string, want int, step string) {
	t.Helper()

	if got := len(splitJournalLines(mustReadCommandFile(t, journalPath))); got != want {
		t.Fatalf("journal lines after %s = %d, want %d", step, got, want)
	}
}

func assertClosureIssueState(t *testing.T, workspace, issueID, wantStatus, wantStatusCategory string, wantLabels []string) {
	t.Helper()

	var got closureIssueView
	if err := json.Unmarshal([]byte(mustGetIssueJSON(t, workspace, issueID)), &got); err != nil {
		t.Fatalf("Unmarshal(%s) error = %v", issueID, err)
	}
	if got.Status != wantStatus {
		t.Fatalf("issue %s status = %q, want %q", issueID, got.Status, wantStatus)
	}
	if got.StatusCategory != wantStatusCategory {
		t.Fatalf("issue %s status_category = %q, want %q", issueID, got.StatusCategory, wantStatusCategory)
	}
	if !reflect.DeepEqual(got.Labels, wantLabels) {
		t.Fatalf("issue %s labels = %#v, want %#v", issueID, got.Labels, wantLabels)
	}
}

func mustSnapshotClosureIssues(t *testing.T, workspace string, issueIDs ...string) []e2eClosureIssueRow {
	t.Helper()

	conn, err := fiddb.Open(filepath.Join(workspace, ".fid", "fid.sqlite"))
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	defer closeCommandDB(t, conn)

	rows, err := conn.Query(`
		SELECT id, title, kind, COALESCE(parent_id, ''), workflow_id, status, COALESCE(owner_role, '')
		FROM issues
		WHERE id IN (?, ?)
		ORDER BY id
	`, issueIDs[0], issueIDs[1])
	if err != nil {
		t.Fatalf("query issues snapshot: %v", err)
	}
	defer closeE2ERows(t, rows)

	var snapshot []e2eClosureIssueRow
	for rows.Next() {
		var row e2eClosureIssueRow
		if err := rows.Scan(&row.ID, &row.Title, &row.Kind, &row.ParentID, &row.WorkflowID, &row.Status, &row.OwnerRole); err != nil {
			t.Fatalf("scan issues snapshot: %v", err)
		}
		snapshot = append(snapshot, row)
	}
	if err := rows.Err(); err != nil {
		t.Fatalf("iterate issues snapshot: %v", err)
	}
	return snapshot
}

func mustSnapshotClosureClaims(t *testing.T, workspace string) []e2eClosureClaimRow {
	t.Helper()

	conn, err := fiddb.Open(filepath.Join(workspace, ".fid", "fid.sqlite"))
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	defer closeCommandDB(t, conn)

	rows, err := conn.Query(`
		SELECT ref_type, ref_id, claimed_by, claimed_at, ttl_seconds, COALESCE(note, '')
		FROM claims
		ORDER BY ref_type, ref_id, claimed_by, claimed_at
	`)
	if err != nil {
		t.Fatalf("query claims snapshot: %v", err)
	}
	defer closeE2ERows(t, rows)

	var snapshot []e2eClosureClaimRow
	for rows.Next() {
		var row e2eClosureClaimRow
		if err := rows.Scan(&row.RefType, &row.RefID, &row.ClaimedBy, &row.ClaimedAt, &row.TTLSeconds, &row.Note); err != nil {
			t.Fatalf("scan claims snapshot: %v", err)
		}
		snapshot = append(snapshot, row)
	}
	if err := rows.Err(); err != nil {
		t.Fatalf("iterate claims snapshot: %v", err)
	}
	return snapshot
}

func mustSnapshotClosureIssueLabels(t *testing.T, workspace string) []e2eClosureIssueLabelRow {
	t.Helper()

	conn, err := fiddb.Open(filepath.Join(workspace, ".fid", "fid.sqlite"))
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	defer closeCommandDB(t, conn)

	rows, err := conn.Query(`
		SELECT issue_id, label
		FROM issue_labels
		ORDER BY issue_id, label
	`)
	if err != nil {
		t.Fatalf("query issue_labels snapshot: %v", err)
	}
	defer closeE2ERows(t, rows)

	var snapshot []e2eClosureIssueLabelRow
	for rows.Next() {
		var row e2eClosureIssueLabelRow
		if err := rows.Scan(&row.IssueID, &row.Label); err != nil {
			t.Fatalf("scan issue_labels snapshot: %v", err)
		}
		snapshot = append(snapshot, row)
	}
	if err := rows.Err(); err != nil {
		t.Fatalf("iterate issue_labels snapshot: %v", err)
	}
	return snapshot
}
