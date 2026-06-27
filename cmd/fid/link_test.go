package main

import (
	"database/sql"
	"path/filepath"
	"strings"
	"testing"

	fiddb "github.com/mbell-nzl/fid/internal/db"
)

func TestLinkAddCreatesDeclaredAndInverseRowsInOneJournalEvent(t *testing.T) {
	workspace := t.TempDir()
	runInitCommand(t, workspace, "init")

	_, firstStderr := runFIDCommand(t, workspace, "issue", "create", "--kind", "task", "--title", "First issue", "--by", "developer")
	firstIssueID := parseCreatedIssueID(t, firstStderr)

	_, secondStderr := runFIDCommand(t, workspace, "issue", "create", "--kind", "task", "--title", "Second issue", "--by", "developer")
	secondIssueID := parseCreatedIssueID(t, secondStderr)

	journalPath := mustJournalPath(t, workspace)
	linesBefore := len(splitJournalLines(mustReadCommandFile(t, journalPath)))

	stdout, stderr := runFIDCommand(t, workspace, "link", "add", firstIssueID, "blocks", secondIssueID, "--by", "developer")
	if stdout != "" {
		t.Fatalf("stdout = %q, want empty", stdout)
	}
	for _, want := range []string{
		"linked " + firstIssueID + " blocks " + secondIssueID,
		"inverse: " + secondIssueID + " blocked_by " + firstIssueID,
	} {
		if !strings.Contains(stderr, want) {
			t.Fatalf("stderr = %q, want substring %q", stderr, want)
		}
	}

	linesAfter := len(splitJournalLines(mustReadCommandFile(t, journalPath)))
	if linesAfter != linesBefore+1 {
		t.Fatalf("journal lines after link add = %d, want %d", linesAfter, linesBefore+1)
	}

	conn, err := fiddb.Open(filepath.Join(workspace, ".fid", "fid.sqlite"))
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	defer closeCommandDB(t, conn)

	rows := loadCommandLinks(t, conn)
	if len(rows) != 2 {
		t.Fatalf("link row count = %d, want 2", len(rows))
	}

	wantRows := []commandLinkRow{
		{FromType: "issue", FromID: firstIssueID, Relation: "blocks", ToType: "issue", ToID: secondIssueID},
		{FromType: "issue", FromID: secondIssueID, Relation: "blocked_by", ToType: "issue", ToID: firstIssueID},
	}
	for _, want := range wantRows {
		if !containsCommandLinkRow(rows, want) {
			t.Fatalf("links rows = %#v, want row %#v", rows, want)
		}
	}

	event := mustReadLastJournalEvent(t, journalPath)
	if len(event.Effects) != 2 {
		t.Fatalf("effect count = %d, want 2", len(event.Effects))
	}
	for _, effect := range event.Effects {
		if effect.Op != "insert" {
			t.Fatalf("effect op = %q, want insert", effect.Op)
		}
		if effect.Table != "links" {
			t.Fatalf("effect table = %q, want links", effect.Table)
		}
	}
}

func TestLinkAddImplementsRequirementAndRejectsInvalidInputBeforeAppend(t *testing.T) {
	t.Run("implements requirement", func(t *testing.T) {
		workspace := t.TempDir()
		runInitCommand(t, workspace, "init")

		issueID := createCommandIssue(t, workspace, "Implements issue")
		syncCommandRequirement(t, workspace, "PRD-LINK-002")

		stdout, stderr := runFIDCommand(t, workspace, "link", "add", issueID, "implements", "PRD-LINK-002", "--by", "developer")
		if stdout != "" {
			t.Fatalf("stdout = %q, want empty", stdout)
		}
		for _, want := range []string{
			"linked " + issueID + " implements PRD-LINK-002",
			"inverse: PRD-LINK-002 implemented_by " + issueID,
		} {
			if !strings.Contains(stderr, want) {
				t.Fatalf("stderr = %q, want substring %q", stderr, want)
			}
		}

		conn, err := fiddb.Open(filepath.Join(workspace, ".fid", "fid.sqlite"))
		if err != nil {
			t.Fatalf("Open() error = %v", err)
		}
		defer closeCommandDB(t, conn)

		rows := loadCommandLinks(t, conn)
		if !containsCommandLinkRow(rows, commandLinkRow{FromType: "issue", FromID: issueID, Relation: "implements", ToType: "requirement", ToID: "PRD-LINK-002"}) {
			t.Fatalf("links rows = %#v, want implements row", rows)
		}
		if !containsCommandLinkRow(rows, commandLinkRow{FromType: "requirement", FromID: "PRD-LINK-002", Relation: "implemented_by", ToType: "issue", ToID: issueID}) {
			t.Fatalf("links rows = %#v, want implemented_by row", rows)
		}
	})

	t.Run("rejections keep journal unchanged", func(t *testing.T) {
		tests := []struct {
			name          string
			setup         func(t *testing.T, workspace string) []string
			command       []string
			wantSubstring string
		}{
			{
				name: "duplicate add",
				setup: func(t *testing.T, workspace string) []string {
					firstIssueID := createCommandIssue(t, workspace, "First issue")
					secondIssueID := createCommandIssue(t, workspace, "Second issue")
					runFIDCommand(t, workspace, "link", "add", firstIssueID, "blocks", secondIssueID, "--by", "developer")
					return []string{firstIssueID, secondIssueID}
				},
				command:       []string{"link", "add", "", "blocks", "", "--by", "developer"},
				wantSubstring: "error: link add: link already exists:",
			},
			{
				name: "unknown relation",
				setup: func(t *testing.T, workspace string) []string {
					return []string{createCommandIssue(t, workspace, "First issue"), createCommandIssue(t, workspace, "Second issue")}
				},
				command:       []string{"link", "add", "", "parent_of", "", "--by", "developer"},
				wantSubstring: `error: link add: unknown relation "parent_of"`,
			},
			{
				name: "self link",
				setup: func(t *testing.T, workspace string) []string {
					issueID := createCommandIssue(t, workspace, "Self issue")
					return []string{issueID, issueID}
				},
				command:       []string{"link", "add", "", "relates_to", "", "--by", "developer"},
				wantSubstring: "error: link add: self link is not allowed:",
			},
			{
				name: "implements issue to issue",
				setup: func(t *testing.T, workspace string) []string {
					return []string{createCommandIssue(t, workspace, "First issue"), createCommandIssue(t, workspace, "Second issue")}
				},
				command:       []string{"link", "add", "", "implements", "", "--by", "developer"},
				wantSubstring: `error: link add: relation "implements" is not valid between issue`,
			},
			{
				name: "blocks issue to requirement",
				setup: func(t *testing.T, workspace string) []string {
					issueID := createCommandIssue(t, workspace, "Issue")
					syncCommandRequirement(t, workspace, "PRD-LINK-002")
					return []string{issueID, "PRD-LINK-002"}
				},
				command:       []string{"link", "add", "", "blocks", "", "--by", "developer"},
				wantSubstring: `error: link add: relation "blocks" is not valid between issue`,
			},
			{
				name: "missing requirement endpoint",
				setup: func(t *testing.T, workspace string) []string {
					return []string{createCommandIssue(t, workspace, "Issue"), "PRD-LINK-002"}
				},
				command:       []string{"link", "add", "", "implements", "", "--by", "developer"},
				wantSubstring: `error: link add: requirement "PRD-LINK-002" does not resolve`,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				workspace := t.TempDir()
				runInitCommand(t, workspace, "init")
				ids := tt.setup(t, workspace)

				command := append([]string(nil), tt.command...)
				command[2] = ids[0]
				command[4] = ids[1]

				journalPath := mustJournalPath(t, workspace)
				linesBefore := len(splitJournalLines(mustReadCommandFile(t, journalPath)))

				stdout, stderr, err := runFIDCommandExpectError(t, workspace, command...)
				if err == nil {
					t.Fatal("Execute() error = nil, want rejection")
				}
				if stdout != "" {
					t.Fatalf("stdout = %q, want empty", stdout)
				}
				if !strings.Contains(stderr, tt.wantSubstring) {
					t.Fatalf("stderr = %q, want substring %q", stderr, tt.wantSubstring)
				}

				linesAfter := len(splitJournalLines(mustReadCommandFile(t, journalPath)))
				if linesAfter != linesBefore {
					t.Fatalf("journal lines after rejection = %d, want %d", linesAfter, linesBefore)
				}
			})
		}
	})
}

func TestLinkAddPolicyShowParityAndUnknownRoleRejection(t *testing.T) {
	workspace := t.TempDir()
	runInitCommand(t, workspace, "init")

	stdout, stderr := runFIDCommand(t, workspace, "policy", "show", "--format", "json")
	if stderr != "" {
		t.Fatalf("policy show stderr = %q, want empty", stderr)
	}

	rows := decodeCommandJSONStructure(t, stdout).([]any)
	for _, action := range []string{"link.add", "link.remove"} {
		for _, role := range []string{"developer", "qa_reviewer", "qa_lead", "tl_lite", "tl", "workflow_runner", "human"} {
			if got := commandPolicyCellAllowed(t, rows, action, role); !got {
				t.Fatalf("policy show %s/%s = false, want true", action, role)
			}
		}
	}

	firstIssueID := createCommandIssue(t, workspace, "First issue")
	secondIssueID := createCommandIssue(t, workspace, "Second issue")

	journalPath := mustJournalPath(t, workspace)
	linesBefore := len(splitJournalLines(mustReadCommandFile(t, journalPath)))

	rejectedStdout, rejectedStderr, err := runFIDCommandExpectError(t, workspace, "link", "add", firstIssueID, "blocks", secondIssueID, "--by", "intern")
	if err == nil {
		t.Fatal("Execute() error = nil, want unknown-role rejection")
	}
	if rejectedStdout != "" {
		t.Fatalf("stdout = %q, want empty", rejectedStdout)
	}
	if !strings.Contains(rejectedStderr, `error: role "intern" not in project profile`) {
		t.Fatalf("stderr = %q, want unknown-role error", rejectedStderr)
	}
	if got := len(splitJournalLines(mustReadCommandFile(t, journalPath))); got != linesBefore {
		t.Fatalf("journal lines after unknown-role rejection = %d, want %d", got, linesBefore)
	}
}

type commandLinkRow struct {
	FromType string
	FromID   string
	Relation string
	ToType   string
	ToID     string
}

func loadCommandLinks(t *testing.T, conn *sql.DB) []commandLinkRow {
	t.Helper()

	rows, err := conn.Query(`
		SELECT from_type, from_id, relation, to_type, to_id
		FROM links
		ORDER BY from_type, from_id, relation, to_type, to_id
	`)
	if err != nil {
		t.Fatalf("Query(links) error = %v", err)
	}
	defer closeCommandLinkRows(t, rows)

	items := []commandLinkRow{}
	for rows.Next() {
		var item commandLinkRow
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

func closeCommandLinkRows(t *testing.T, rows *sql.Rows) {
	t.Helper()
	if err := rows.Close(); err != nil {
		t.Fatalf("Rows.Close() error = %v", err)
	}
}

func containsCommandLinkRow(rows []commandLinkRow, want commandLinkRow) bool {
	for _, row := range rows {
		if row == want {
			return true
		}
	}
	return false
}

func createCommandIssue(t *testing.T, workspace, title string) string {
	t.Helper()

	_, stderr := runFIDCommand(t, workspace, "issue", "create", "--kind", "task", "--title", title, "--by", "developer")
	return parseCreatedIssueID(t, stderr)
}

func syncCommandRequirement(t *testing.T, workspace, requirementID string) {
	t.Helper()

	writeRequirementFixture(t, workspace, "project-docs/specs/link_requirements.md", `# Requirements

## 3.8 Cross-Entity Links

| ID | Requirement | Priority | Notes |
|---|---|---|---|
| PRD-LINK-002 | The system shall support links from issues to requirements. | Must | v0.1 |
`)

	_, stderr := runFIDCommand(t, workspace, "requirement", "sync", "--source-doc", "project-docs/specs/link_requirements.md", "--id", requirementID, "--by", "tl")
	if !strings.Contains(stderr, "synced requirements from project-docs/specs/link_requirements.md") {
		t.Fatalf("stderr = %q, want requirement sync confirmation", stderr)
	}
}
