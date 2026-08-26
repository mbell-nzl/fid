package main

import (
	"database/sql"
	"encoding/json"
	"path/filepath"
	"sort"
	"strings"
	"testing"
	"time"

	fiddb "github.com/mbell-nzl/fid/internal/db"
)

func TestLinkAddCreatesDeclaredAndInverseRowsInOneJournalEvent(t *testing.T) {
	t.Parallel()
	workspace := t.TempDir()
	clonePristineWorkspaceInto(t, workspace)

	_, firstStderr := runFIDCommand(t, workspace, "issue", "create", "--kind", "task", "--title", "First issue", "--actor", "agent:developer")
	firstIssueID := parseCreatedIssueID(t, firstStderr)

	_, secondStderr := runFIDCommand(t, workspace, "issue", "create", "--kind", "task", "--title", "Second issue", "--actor", "agent:developer")
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
	if len(event.Effects) != 6 {
		t.Fatalf("effect count = %d, want link effects plus idempotency result", len(event.Effects))
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
		case effect.Op == "insert" && effect.Table == "idempotency_results":
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
	if activityEvent != "link-added" {
		t.Fatalf("activity_log event = %q, want link-added", activityEvent)
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
	if len(event.Effects) != 6 {
		t.Fatalf("effect count = %d, want link effects plus idempotency result", len(event.Effects))
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
		case effect.Op == "insert" && effect.Table == "idempotency_results":
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
	if activityEvent != "link-removed" {
		t.Fatalf("activity_log event = %q, want %q", activityEvent, "link-removed")
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
		wantSuccess   bool
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
			name: "project-local actor",
			setup: func(t *testing.T, workspace string) []string {
				firstIssueID := createCommandIssue(t, workspace, "First issue")
				secondIssueID := createCommandIssue(t, workspace, "Second issue")
				runFIDCommand(t, workspace, "link", "add", firstIssueID, "blocks", secondIssueID, "--from-type", "issue", "--to-type", "issue", "--by", "developer")
				return []string{firstIssueID, secondIssueID}
			},
			command:     []string{"link", "remove", "", "blocks", "", "--from-type", "issue", "--to-type", "issue", "--by", "intern"},
			wantSuccess: true,
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
			if tt.wantSuccess {
				if err != nil || stdout == "" || stderr != "" {
					t.Fatalf("project-local actor result = (%q, %q, %v), want stdout success", stdout, stderr, err)
				}
				if got := len(splitJournalLines(mustReadCommandFile(t, journalPath))); got != linesBefore+1 {
					t.Fatalf("journal lines after project-local actor = %d, want %d", got, linesBefore+1)
				}
				return
			}
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

	rows := legacyPolicyRows(t)
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
	if err != nil || rejectedStderr != "" || rejectedStdout == "" {
		t.Fatalf("project-local actor result = (%q, %q, %v), want stdout success", rejectedStdout, rejectedStderr, err)
	}
	if got := len(splitJournalLines(mustReadCommandFile(t, journalPath))); got != linesBefore+1 {
		t.Fatalf("journal lines after project-local actor = %d, want %d", got, linesBefore+1)
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

func TestLinkAddUsesSuppliedTypeForAmbiguousSameIDEndpoints(t *testing.T) {
	workspace := t.TempDir()
	clonePristineWorkspaceInto(t, workspace)

	issueID := createCommandIssue(t, workspace, "Issue and requirement share an opaque ID")
	conn, err := fiddb.Open(filepath.Join(workspace, ".fid", "fid.sqlite"))
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	defer closeCommandDB(t, conn)
	if _, err := conn.Exec(`
		INSERT INTO requirements (id, type, summary, source_doc, source_section, status)
		VALUES (?, 'PRD', 'Requirement with the same opaque ID', 'project-docs/specs/link_requirements.md', '3.8', 'active')
	`, issueID); err != nil {
		t.Fatalf("insert same-ID requirement: %v", err)
	}
	journalPath := mustJournalPath(t, workspace)
	before := len(splitJournalLines(mustReadCommandFile(t, journalPath)))

	result := executeCommandForTest(workspace,
		"--json", "link", "add", issueID, "implements", issueID,
		"--from-type", "issue", "--to-type", "requirement", "--by", "developer",
	)
	if result.err != nil || result.stderr != "" {
		t.Fatalf("link.add result = %#v, want JSON success on stdout only", result)
	}

	envelope := decodeCommandJSONStructure(t, result.stdout).(map[string]any)
	if envelope["operation"] != "link.add" || envelope["event"] == nil {
		t.Fatalf("link.add envelope = %#v, want mutation event", envelope)
	}
	link := envelope["result"].(map[string]any)["link"].(map[string]any)
	if link["from_type"] != "issue" || link["to_type"] != "requirement" || link["from_id"] != issueID || link["to_id"] != issueID {
		t.Fatalf("declared link = %#v, want supplied issue/requirement types for %q", link, issueID)
	}
	if got := len(splitJournalLines(mustReadCommandFile(t, journalPath))); got != before+1 {
		t.Fatalf("journal lines after typed ambiguous add = %d, want %d", got, before+1)
	}
}

func TestLinkMutationsJSONMatchContract(t *testing.T) {
	workspace := t.TempDir()
	clonePristineWorkspaceInto(t, workspace)
	firstID := createCommandIssue(t, workspace, "First JSON link issue")
	secondID := createCommandIssue(t, workspace, "Second JSON link issue")
	journalPath := mustJournalPath(t, workspace)

	add := executeCommandForTest(workspace,
		"--json", "link", "add", firstID, "blocks", secondID,
		"--from-type", "issue", "--to-type", "issue", "--by", "developer",
	)
	assertLinkMutationEnvelope(t, add, journalPath, "link.add", "link.added", firstID, "blocks", secondID, 2)

	remove := executeCommandForTest(workspace,
		"--json", "link", "remove", firstID, "blocks", secondID,
		"--from-type", "issue", "--to-type", "issue", "--by", "developer",
	)
	assertLinkMutationEnvelope(t, remove, journalPath, "link.remove", "link.removed", firstID, "blocks", secondID, 3)
}

func TestLinkJSONRequestAndTypedFailuresAreAtomic(t *testing.T) {
	workspace := t.TempDir()
	clonePristineWorkspaceInto(t, workspace)
	firstID := createCommandIssue(t, workspace, "First link failure issue")
	secondID := createCommandIssue(t, workspace, "Second link failure issue")
	journalPath := mustJournalPath(t, workspace)

	tests := []struct {
		name      string
		operation string
		args      []string
		assert    func(t *testing.T, output string)
	}{
		{
			name: "add omitted from type", operation: "link.add",
			args:   []string{"--json", "link", "add", firstID, "blocks", secondID, "--to-type", "issue", "--by", "developer"},
			assert: func(t *testing.T, output string) { assertLinkInvalidRequestEnvelope(t, output, "link.add") },
		},
		{
			name: "add omitted to type", operation: "link.add",
			args:   []string{"--json", "link", "add", firstID, "blocks", secondID, "--from-type", "issue", "--by", "developer"},
			assert: func(t *testing.T, output string) { assertLinkInvalidRequestEnvelope(t, output, "link.add") },
		},
		{
			name: "remove omitted from type", operation: "link.remove",
			args:   []string{"--json", "link", "remove", firstID, "blocks", secondID, "--to-type", "issue", "--by", "developer"},
			assert: func(t *testing.T, output string) { assertLinkInvalidRequestEnvelope(t, output, "link.remove") },
		},
		{
			name: "remove omitted to type", operation: "link.remove",
			args:   []string{"--json", "link", "remove", firstID, "blocks", secondID, "--from-type", "issue", "--by", "developer"},
			assert: func(t *testing.T, output string) { assertLinkInvalidRequestEnvelope(t, output, "link.remove") },
		},
		{
			name: "list omitted type", operation: "link.list",
			args:   []string{"--json", "link", "list", firstID},
			assert: func(t *testing.T, output string) { assertLinkInvalidRequestEnvelope(t, output, "link.list") },
		},
		{
			name: "self link", operation: "link.add",
			args: []string{"--json", "link", "add", firstID, "relates_to", firstID, "--from-type", "issue", "--to-type", "issue", "--by", "developer"},
			assert: func(t *testing.T, output string) {
				assertLinkFailure(t, output, "link.add", "invalid-argument", nil, map[string]any{
					"argument": "link", "value": firstID, "allowed": []any{"distinct endpoint ids"},
				})
			},
		},
		{
			name: "missing link", operation: "link.remove",
			args: []string{"--json", "link", "remove", firstID, "blocks", secondID, "--from-type", "issue", "--to-type", "issue", "--by", "developer"},
			assert: func(t *testing.T, output string) {
				assertLinkFailure(t, output, "link.remove", "link-not-found", map[string]any{"type": "link", "id": firstID}, map[string]any{
					"from_type": "issue", "from_id": firstID, "relation": "blocks", "to_type": "issue", "to_id": secondID,
				})
			},
		},
		{
			name: "missing issue endpoint", operation: "link.add",
			args: []string{"--json", "link", "add", "FID-missing", "blocks", secondID, "--from-type", "issue", "--to-type", "issue", "--by", "developer"},
			assert: func(t *testing.T, output string) {
				assertLinkFailure(t, output, "link.add", "issue-not-found", map[string]any{"type": "issue", "id": "FID-missing"}, map[string]any{"id": "FID-missing"})
			},
		},
		{
			name: "missing requirement endpoint", operation: "link.add",
			args: []string{"--json", "link", "add", firstID, "implements", "PRD-missing", "--from-type", "issue", "--to-type", "requirement", "--by", "developer"},
			assert: func(t *testing.T, output string) {
				assertLinkFailure(t, output, "link.add", "requirement-not-found", map[string]any{"type": "requirement", "id": "PRD-missing"}, map[string]any{"id": "PRD-missing"})
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			before := len(splitJournalLines(mustReadCommandFile(t, journalPath)))
			result := executeCommandForTest(workspace, test.args...)
			if result.err == nil || result.stdout != "" {
				t.Fatalf("%s result = %#v, want stderr-only rejection", test.operation, result)
			}
			test.assert(t, result.stderr)
			if got := len(splitJournalLines(mustReadCommandFile(t, journalPath))); got != before {
				t.Fatalf("journal lines after %s rejection = %d, want %d", test.name, got, before)
			}
		})
	}
}

func TestLinkJSONDependencyCycleHasExactPathAndNoAppend(t *testing.T) {
	workspace := t.TempDir()
	clonePristineWorkspaceInto(t, workspace)
	firstID := createCommandIssue(t, workspace, "First cycle issue")
	secondID := createCommandIssue(t, workspace, "Second cycle issue")
	runFIDCommand(t, workspace, "link", "add", firstID, "blocks", secondID, "--from-type", "issue", "--to-type", "issue", "--by", "developer")
	journalPath := mustJournalPath(t, workspace)
	before := len(splitJournalLines(mustReadCommandFile(t, journalPath)))

	result := executeCommandForTest(workspace,
		"--json", "link", "add", secondID, "blocks", firstID,
		"--from-type", "issue", "--to-type", "issue", "--by", "developer",
	)
	if result.err == nil || result.stdout != "" {
		t.Fatalf("dependency-cycle result = %#v, want stderr-only rejection", result)
	}
	assertLinkFailure(t, result.stderr, "link.add", "dependency-cycle", nil, map[string]any{
		"path": []any{secondID, firstID, secondID},
	})
	if got := len(splitJournalLines(mustReadCommandFile(t, journalPath))); got != before {
		t.Fatalf("journal lines after dependency-cycle rejection = %d, want %d", got, before)
	}
}

func assertLinkMutationEnvelope(t *testing.T, result commandExecutionResult, journalPath, operation, eventType, fromID, relation, toID string, revision int) {
	t.Helper()
	if result.err != nil || result.stderr != "" {
		t.Fatalf("%s result = %#v, want JSON success on stdout only", operation, result)
	}
	var envelope struct {
		ContractVersion string         `json:"contract_version"`
		OK              bool           `json:"ok"`
		Operation       string         `json:"operation"`
		Result          map[string]any `json:"result"`
		Event           map[string]any `json:"event"`
	}
	if err := json.Unmarshal([]byte(result.stdout), &envelope); err != nil {
		t.Fatalf("Unmarshal(%s stdout) error = %v; stdout = %q", operation, err, result.stdout)
	}
	assertEnvelopeKeySet(t, []byte(result.stdout), []string{"contract_version", "event", "ok", "operation", "result"})
	assertExactKeys(t, envelope.Result, []string{"idempotency_key", "issue_revisions", "link", "replayed"})
	assertExactKeys(t, envelope.Event, []string{"id", "occurred_at"})
	link := envelope.Result["link"].(map[string]any)
	assertExactKeys(t, link, []string{"from_id", "from_type", "relation", "to_id", "to_type"})
	if envelope.ContractVersion != "fid-cli.v1" || !envelope.OK || envelope.Operation != operation ||
		link["from_type"] != "issue" || link["from_id"] != fromID || link["relation"] != relation || link["to_type"] != "issue" || link["to_id"] != toID {
		t.Fatalf("%s envelope = %#v, want exact typed link result", operation, envelope)
	}

	issueRevisions := envelope.Result["issue_revisions"].([]any)
	if len(issueRevisions) != 2 {
		t.Fatalf("issue_revisions = %#v, want two sorted revisions", issueRevisions)
	}
	wantIDs := []string{fromID, toID}
	sort.Strings(wantIDs)
	for index, wantID := range wantIDs {
		item := issueRevisions[index].(map[string]any)
		assertExactKeys(t, item, []string{"id", "revision"})
		if item["id"] != wantID || item["revision"] != float64(revision) {
			t.Fatalf("issue_revisions[%d] = %#v, want id %q revision %d", index, item, wantID, revision)
		}
	}

	last := mustReadLastJournalEvent(t, journalPath)
	if last.EventType != eventType || last.ActorV2.ActorID != "developer" || last.ActorV2.ActorType != "agent" || envelope.Event["id"] != last.EventID || envelope.Event["occurred_at"] != last.OccurredAt.UTC().Format(time.RFC3339Nano) {
		t.Fatalf("%s envelope event = %#v, journal tail = %#v", operation, envelope.Event, last)
	}
	if len(last.Effects) != 6 {
		t.Fatalf("%s effect count = %d, want link effects plus idempotency result", operation, len(last.Effects))
	}
	for _, effect := range last.Effects {
		if effect.Table == "evidence" || effect.Table == "outcomes" || effect.Table == "claims" {
			t.Fatalf("%s created protected/trusted effect %#v", operation, effect)
		}
	}
}

func assertLinkInvalidRequestEnvelope(t *testing.T, output, operation string) {
	t.Helper()
	assertLinkFailure(t, output, operation, "invalid-argument", nil, map[string]any{
		"argument": "request", "value": nil, "allowed": []any{"accepted " + operation + " input"},
	})
}

func assertLinkFailure(t *testing.T, output, operation, code string, wantSubject, wantDetails map[string]any) {
	t.Helper()
	envelope := decodeCommandJSONStructure(t, output).(map[string]any)
	assertExactKeys(t, envelope, []string{"contract_version", "error", "ok", "operation"})
	errorObject := envelope["error"].(map[string]any)
	assertExactKeys(t, errorObject, []string{"category", "code", "details", "message", "retryable", "subject"})
	if envelope["contract_version"] != "fid-cli.v1" || envelope["ok"] != false || envelope["operation"] != operation || errorObject["code"] != code || errorObject["retryable"] != false {
		t.Fatalf("failure envelope = %#v, want %s %s", envelope, operation, code)
	}
	if wantSubject == nil {
		if errorObject["subject"] != nil {
			t.Fatalf("failure subject = %#v, want null", errorObject["subject"])
		}
	} else if got := errorObject["subject"].(map[string]any); !sameJSONMap(got, wantSubject) {
		t.Fatalf("failure subject = %#v, want %#v", got, wantSubject)
	}
	if got := errorObject["details"].(map[string]any); !sameJSONMap(got, wantDetails) {
		t.Fatalf("failure details = %#v, want %#v", got, wantDetails)
	}
}

func sameJSONMap(got, want map[string]any) bool {
	if len(got) != len(want) {
		return false
	}
	for key, wantValue := range want {
		gotValue, ok := got[key]
		if !ok || !sameJSONValue(gotValue, wantValue) {
			return false
		}
	}
	return true
}

func sameJSONValue(got, want any) bool {
	gotJSON, gotErr := json.Marshal(got)
	wantJSON, wantErr := json.Marshal(want)
	return gotErr == nil && wantErr == nil && string(gotJSON) == string(wantJSON)
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

	_, stderr := runFIDCommand(t, workspace, "issue", "create", "--kind", "task", "--title", title, "--actor", "agent:developer")
	return parseCreatedIssueID(t, stderr)
}

func syncCommandRequirement(t *testing.T, workspace, requirementID string) {
	t.Helper()

	syncCommandRequirementWithSummary(t, workspace, requirementID, "The system shall support links from issues to requirements.")
}

func syncCommandRequirementWithSummary(t *testing.T, workspace, requirementID, summary string) {
	t.Helper()

	writeRequirementFixture(t, workspace, "project-docs/specs/link_requirements.md", "# Requirements\n\n## 3.8 Cross-Entity Links\n\n| ID | Requirement | Priority | Notes |\n|---|---|---|---|\n| PRD-LINK-002 | "+summary+" | Must | v0.1 |\n")

	stdout, stderr := runFIDCommand(t, workspace, "requirement", "sync", "--source-doc", "project-docs/specs/link_requirements.md", "--id", requirementID, "--by", "tl")
	if stderr != "" || !strings.Contains(stdout, "synced requirements from project-docs/specs/link_requirements.md") {
		t.Fatalf("requirement sync = (%q, %q), want stdout confirmation", stdout, stderr)
	}
}
