package main

import (
	"database/sql"
	"encoding/json"
	"path/filepath"
	"strings"
	"testing"

	fiddb "github.com/mbell-nzl/fid/internal/db"
)

func TestLinkAddCreatesDeclaredAndInverseRowsInOneJournalEvent(t *testing.T) {
	t.Parallel()
	workspace := t.TempDir()
	clonePristineWorkspaceInto(t, workspace)

	_, firstStderr := runFIDCommand(t, workspace, "issue", "create", "--kind", "task", "--title", "First issue", "--by", "developer")
	firstIssueID := parseCreatedIssueID(t, firstStderr)

	_, secondStderr := runFIDCommand(t, workspace, "issue", "create", "--kind", "task", "--title", "Second issue", "--by", "developer")
	secondIssueID := parseCreatedIssueID(t, secondStderr)

	journalPath := mustJournalPath(t, workspace)
	linesBefore := len(splitJournalLines(mustReadCommandFile(t, journalPath)))

	stdout, stderr := runFIDCommand(t, workspace, "link", "add", firstIssueID, "blocks", secondIssueID, "--from-type", "issue", "--to-type", "issue", "--by", "developer")
	if stderr != "" {
		t.Fatalf("stderr = %q, want empty", stderr)
	}
	for _, want := range []string{
		"linked " + firstIssueID + " blocks " + secondIssueID,
		"inverse: " + secondIssueID + " blocked_by " + firstIssueID,
	} {
		if !strings.Contains(stdout, want) {
			t.Fatalf("stdout = %q, want substring %q", stdout, want)
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
	if len(event.Effects) != 5 {
		t.Fatalf("effect count = %d, want 5", len(event.Effects))
	}
	insertEffects := 0
	activityEffects := 0
	for _, effect := range event.Effects {
		switch {
		case effect.Op == "insert" && effect.Table == "links":
			insertEffects++
		case effect.Op == "activity" && effect.Table == "activity_log":
			activityEffects++
		case effect.Op == "update" && effect.Table == "issues":
		default:
			t.Fatalf("unexpected effect = %#v", effect)
		}
	}
	if insertEffects != 2 || activityEffects != 1 {
		t.Fatalf("effect counts = %d inserts / %d activity, want 2 / 1", insertEffects, activityEffects)
	}

	var (
		activityEvent string
		actor         string
		details       string
	)
	if err := conn.QueryRow(`
		SELECT event, COALESCE(actor, ''), COALESCE(details, '')
		FROM activity_log
		WHERE journal_event_id = ?
	`, event.EventID).Scan(&activityEvent, &actor, &details); err != nil {
		t.Fatalf("query activity_log row: %v", err)
	}
	if activityEvent != "link_added" {
		t.Fatalf("activity_log event = %q, want link_added", activityEvent)
	}
	if actor != "developer" {
		t.Fatalf("activity_log actor = %q, want developer", actor)
	}
	for _, want := range []string{
		`"relation":"blocks"`,
		`"to_id":"` + secondIssueID + `"`,
		`"inverse_relation":"blocked_by"`,
	} {
		if !strings.Contains(details, want) {
			t.Fatalf("activity_log details = %q, want substring %q", details, want)
		}
	}
}

func TestLinkAddImplementsRequirementAndRejectsInvalidInputBeforeAppend(t *testing.T) {
	t.Parallel()
	t.Run("implements requirement", func(t *testing.T) {
		workspace := t.TempDir()
		clonePristineWorkspaceInto(t, workspace)

		issueID := createCommandIssue(t, workspace, "Implements issue")
		syncCommandRequirement(t, workspace, "PRD-LINK-002")

		stdout, stderr := runFIDCommand(t, workspace, "link", "add", issueID, "implements", "PRD-LINK-002", "--from-type", "issue", "--to-type", "requirement", "--by", "developer")
		if stderr != "" {
			t.Fatalf("stderr = %q, want empty", stderr)
		}
		for _, want := range []string{
			"linked " + issueID + " implements PRD-LINK-002",
			"inverse: PRD-LINK-002 implemented_by " + issueID,
		} {
			if !strings.Contains(stdout, want) {
				t.Fatalf("stdout = %q, want substring %q", stdout, want)
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
					runFIDCommand(t, workspace, "link", "add", firstIssueID, "blocks", secondIssueID, "--from-type", "issue", "--to-type", "issue", "--by", "developer")
					return []string{firstIssueID, secondIssueID}
				},
				command:       []string{"link", "add", "", "blocks", "", "--from-type", "issue", "--to-type", "issue", "--by", "developer"},
				wantSubstring: "error: link add: link already exists:",
			},
			{
				name: "unknown relation",
				setup: func(t *testing.T, workspace string) []string {
					return []string{createCommandIssue(t, workspace, "First issue"), createCommandIssue(t, workspace, "Second issue")}
				},
				command:       []string{"link", "add", "", "parent_of", "", "--from-type", "issue", "--to-type", "issue", "--by", "developer"},
				wantSubstring: `error: link add: unknown relation "parent_of"`,
			},
			{
				name: "self link",
				setup: func(t *testing.T, workspace string) []string {
					issueID := createCommandIssue(t, workspace, "Self issue")
					return []string{issueID, issueID}
				},
				command:       []string{"link", "add", "", "relates_to", "", "--from-type", "issue", "--to-type", "issue", "--by", "developer"},
				wantSubstring: "error: link add: self link is not allowed:",
			},
			{
				name: "implements issue to issue",
				setup: func(t *testing.T, workspace string) []string {
					return []string{createCommandIssue(t, workspace, "First issue"), createCommandIssue(t, workspace, "Second issue")}
				},
				command:       []string{"link", "add", "", "implements", "", "--from-type", "issue", "--to-type", "issue", "--by", "developer"},
				wantSubstring: `error: link add: relation "implements" is not valid between issue`,
			},
			{
				name: "blocks issue to requirement",
				setup: func(t *testing.T, workspace string) []string {
					issueID := createCommandIssue(t, workspace, "Issue")
					syncCommandRequirement(t, workspace, "PRD-LINK-002")
					return []string{issueID, "PRD-LINK-002"}
				},
				command:       []string{"link", "add", "", "blocks", "", "--from-type", "issue", "--to-type", "requirement", "--by", "developer"},
				wantSubstring: `error: link add: relation "blocks" is not valid between issue`,
			},
			{
				name: "missing requirement endpoint",
				setup: func(t *testing.T, workspace string) []string {
					return []string{createCommandIssue(t, workspace, "Issue"), "PRD-LINK-002"}
				},
				command:       []string{"link", "add", "", "implements", "", "--from-type", "issue", "--to-type", "requirement", "--by", "developer"},
				wantSubstring: `error: link add: requirement "PRD-LINK-002" does not resolve`,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				workspace := t.TempDir()
				clonePristineWorkspaceInto(t, workspace)
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

func TestLinkRemoveDeletesDeclaredAndInverseRowsInOneJournalEvent(t *testing.T) {
	t.Parallel()
	workspace := t.TempDir()
	clonePristineWorkspaceInto(t, workspace)

	firstIssueID := createCommandIssue(t, workspace, "First issue")
	secondIssueID := createCommandIssue(t, workspace, "Second issue")
	runFIDCommand(t, workspace, "link", "add", firstIssueID, "blocks", secondIssueID, "--from-type", "issue", "--to-type", "issue", "--by", "developer")

	journalPath := mustJournalPath(t, workspace)
	linesBefore := len(splitJournalLines(mustReadCommandFile(t, journalPath)))

	conn, err := fiddb.Open(filepath.Join(workspace, ".fid", "fid.sqlite"))
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	defer closeCommandDB(t, conn)

	activityBefore := countActivityLogRows(t, conn)

	stdout, stderr := runFIDCommand(t, workspace, "link", "remove", firstIssueID, "blocks", secondIssueID, "--from-type", "issue", "--to-type", "issue", "--by", "developer")
	if stderr != "" {
		t.Fatalf("stderr = %q, want empty", stderr)
	}
	for _, want := range []string{
		"removed link " + firstIssueID + " blocks " + secondIssueID,
		"inverse removed: " + secondIssueID + " blocked_by " + firstIssueID,
	} {
		if !strings.Contains(stdout, want) {
			t.Fatalf("stdout = %q, want substring %q", stdout, want)
		}
	}

	linesAfter := len(splitJournalLines(mustReadCommandFile(t, journalPath)))
	if linesAfter != linesBefore+1 {
		t.Fatalf("journal lines after link remove = %d, want %d", linesAfter, linesBefore+1)
	}

	if rows := loadCommandLinks(t, conn); len(rows) != 0 {
		t.Fatalf("links rows after remove = %#v, want none", rows)
	}
	if got := countActivityLogRows(t, conn); got != activityBefore+1 {
		t.Fatalf("activity_log rows after remove = %d, want %d", got, activityBefore+1)
	}

	event := mustReadLastJournalEvent(t, journalPath)
	if len(event.Effects) != 5 {
		t.Fatalf("effect count = %d, want 5", len(event.Effects))
	}

	deleteEffects := 0
	activityEffects := 0
	for _, effect := range event.Effects {
		switch {
		case effect.Op == "delete" && effect.Table == "links":
			deleteEffects++
		case effect.Op == "activity" && effect.Table == "activity_log":
			activityEffects++
		case effect.Op == "update" && effect.Table == "issues":
		default:
			t.Fatalf("unexpected effect = %#v", effect)
		}
	}
	if deleteEffects != 2 || activityEffects != 1 {
		t.Fatalf("effect counts = %d deletes / %d activity, want 2 / 1", deleteEffects, activityEffects)
	}

	var (
		activityEvent string
		actor         string
		details       string
	)
	if err := conn.QueryRow(`
		SELECT event, COALESCE(actor, ''), COALESCE(details, '')
		FROM activity_log
		WHERE journal_event_id = ?
	`, event.EventID).Scan(&activityEvent, &actor, &details); err != nil {
		t.Fatalf("query activity_log row: %v", err)
	}
	if activityEvent != "link_removed" {
		t.Fatalf("activity_log event = %q, want %q", activityEvent, "link_removed")
	}
	if actor != "developer" {
		t.Fatalf("activity_log actor = %q, want %q", actor, "developer")
	}
	for _, want := range []string{
		`"relation":"blocks"`,
		`"to_id":"` + secondIssueID + `"`,
		`"inverse_relation":"blocked_by"`,
	} {
		if !strings.Contains(details, want) {
			t.Fatalf("activity_log details = %q, want substring %q", details, want)
		}
	}
}

func TestLinkRemoveRejectsInvalidInputBeforeAppend(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name          string
		setup         func(t *testing.T, workspace string) []string
		command       []string
		wantSubstring string
	}{
		{
			name: "missing link",
			setup: func(t *testing.T, workspace string) []string {
				return []string{createCommandIssue(t, workspace, "First issue"), createCommandIssue(t, workspace, "Second issue")}
			},
			command:       []string{"link", "remove", "", "blocks", "", "--from-type", "issue", "--to-type", "issue", "--by", "developer"},
			wantSubstring: "error: link remove: link does not exist:",
		},
		{
			name: "unknown role",
			setup: func(t *testing.T, workspace string) []string {
				firstIssueID := createCommandIssue(t, workspace, "First issue")
				secondIssueID := createCommandIssue(t, workspace, "Second issue")
				runFIDCommand(t, workspace, "link", "add", firstIssueID, "blocks", secondIssueID, "--from-type", "issue", "--to-type", "issue", "--by", "developer")
				return []string{firstIssueID, secondIssueID}
			},
			command:       []string{"link", "remove", "", "blocks", "", "--from-type", "issue", "--to-type", "issue", "--by", "intern"},
			wantSubstring: `error: role "intern" not in project profile`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			workspace := t.TempDir()
			clonePristineWorkspaceInto(t, workspace)
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
}

func TestLinkAddPolicyShowParityAndUnknownRoleRejection(t *testing.T) {
	t.Parallel()
	workspace := t.TempDir()
	clonePristineWorkspaceInto(t, workspace)

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

	rejectedStdout, rejectedStderr, err := runFIDCommandExpectError(t, workspace, "link", "add", firstIssueID, "blocks", secondIssueID, "--from-type", "issue", "--to-type", "issue", "--by", "intern")
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

func TestLinkListReturnsTypedCanonicalJSONInContractOrder(t *testing.T) {
	t.Parallel()
	workspace := t.TempDir()
	clonePristineWorkspaceInto(t, workspace)

	primaryIssueID := createCommandIssue(t, workspace, "Primary issue")
	outgoingIssueID := createCommandIssue(t, workspace, "Schema migration runner")
	incomingIssueID := createCommandIssue(t, workspace, "Add index on concerns(issue_id, status)")
	syncCommandRequirementWithSummary(t, workspace, "PRD-LINK-002", "Unified polymorphic issue table")

	runFIDCommand(t, workspace, "link", "add", primaryIssueID, "blocks", outgoingIssueID, "--from-type", "issue", "--to-type", "issue", "--by", "developer")
	runFIDCommand(t, workspace, "link", "add", primaryIssueID, "implements", "PRD-LINK-002", "--from-type", "issue", "--to-type", "requirement", "--by", "developer")
	runFIDCommand(t, workspace, "link", "add", incomingIssueID, "discovered_from", primaryIssueID, "--from-type", "issue", "--to-type", "issue", "--by", "developer")

	jsonStdout, jsonStderr := runFIDCommand(t, workspace, "--json", "link", "list", primaryIssueID, "--type", "issue")
	if jsonStderr != "" {
		t.Fatalf("json stderr = %q, want empty", jsonStderr)
	}

	envelope := decodeCommandJSONStructure(t, jsonStdout).(map[string]any)
	if envelope["operation"] != "link.list" || envelope["event"] != nil {
		t.Fatalf("link list envelope = %#v, want link.list read with null event", envelope)
	}
	result := envelope["result"].(map[string]any)
	if result["entity_type"] != "issue" || result["entity_id"] != primaryIssueID || result["count"] != json.Number("3") {
		t.Fatalf("link list result header = %#v", result)
	}
	rows := result["links"].([]any)
	if len(rows) != 3 {
		t.Fatalf("link list rows = %#v, want 3 rows", rows)
	}

	wantRows := []map[string]any{
		{
			"direction":    "outgoing",
			"relation":     "blocks",
			"entity_type":  "issue",
			"entity_id":    outgoingIssueID,
			"entity_title": "Schema migration runner",
		},
		{
			"direction":    "outgoing",
			"relation":     "implements",
			"entity_type":  "requirement",
			"entity_id":    "PRD-LINK-002",
			"entity_title": "Unified polymorphic issue table",
		},
		{
			"direction":    "incoming",
			"relation":     "discovered_from",
			"entity_type":  "issue",
			"entity_id":    incomingIssueID,
			"entity_title": "Add index on concerns(issue_id, status)",
		},
	}
	for index, wantRow := range wantRows {
		gotRow, ok := rows[index].(map[string]any)
		if !ok {
			t.Fatalf("row %d = %#v, want object", index, rows[index])
		}
		for key, wantValue := range wantRow {
			if gotValue := gotRow[key]; gotValue != wantValue {
				t.Fatalf("row %d field %s = %v, want %v", index, key, gotRow[key], wantValue)
			}
		}
	}
}

func TestLinkListEmptyResult(t *testing.T) {
	t.Parallel()
	workspace := t.TempDir()
	clonePristineWorkspaceInto(t, workspace)

	issueID := createCommandIssue(t, workspace, "Unlinked issue")

	stdout, stderr := runFIDCommand(t, workspace, "--json", "link", "list", issueID, "--type", "issue")
	if stderr != "" {
		t.Fatalf("stderr = %q, want empty", stderr)
	}
	result := decodeCommandJSONStructure(t, stdout).(map[string]any)["result"].(map[string]any)
	if result["count"] != json.Number("0") || len(result["links"].([]any)) != 0 {
		t.Fatalf("empty link result = %#v, want no links", result)
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

	syncCommandRequirementWithSummary(t, workspace, requirementID, "The system shall support links from issues to requirements.")
}

func syncCommandRequirementWithSummary(t *testing.T, workspace, requirementID, summary string) {
	t.Helper()

	writeRequirementFixture(t, workspace, "project-docs/specs/link_requirements.md", "# Requirements\n\n## 3.8 Cross-Entity Links\n\n| ID | Requirement | Priority | Notes |\n|---|---|---|---|\n| PRD-LINK-002 | "+summary+" | Must | v0.1 |\n")

	_, stderr := runFIDCommand(t, workspace, "requirement", "sync", "--source-doc", "project-docs/specs/link_requirements.md", "--id", requirementID, "--by", "tl")
	if !strings.Contains(stderr, "synced requirements from project-docs/specs/link_requirements.md") {
		t.Fatalf("stderr = %q, want requirement sync confirmation", stderr)
	}
}

func splitMarkdownTableRow(line string) []string {
	trimmed := strings.TrimSpace(line)
	trimmed = strings.TrimPrefix(trimmed, "|")
	trimmed = strings.TrimSuffix(trimmed, "|")
	parts := strings.Split(trimmed, "|")
	for index := range parts {
		parts[index] = strings.TrimSpace(parts[index])
	}
	return parts
}
