package db

import (
	"database/sql"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"
)

func TestMigration0008EnforcesWorkflowAndStateForeignKeys(t *testing.T) {
	t.Parallel()

	dbPath := filepath.Join(t.TempDir(), "fid.sqlite")
	if err := MigrateUp(dbPath, testMigrationsDir(t)); err != nil {
		t.Fatalf("MigrateUp() error = %v", err)
	}

	conn, err := Open(dbPath)
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	defer closeDB(t, conn)

	if _, err := conn.Exec(`
		INSERT INTO workflows (id, name, applies_to, active)
		VALUES ('custom-test', 'Custom test', 'task', 1)
	`); err != nil {
		t.Fatalf("insert custom workflow: %v", err)
	}
	if _, err := conn.Exec(`
		INSERT INTO workflow_states (workflow_id, state, category, ordinal)
		VALUES ('custom-test', 'review', 'started', 1)
	`); err != nil {
		t.Fatalf("insert custom workflow state: %v", err)
	}

	for _, issue := range []struct {
		name       string
		workflowID string
		status     string
	}{
		{name: "unknown-workflow", workflowID: "missing-workflow", status: "review"},
		{name: "unknown-state", workflowID: "custom-test", status: "missing-state"},
	} {
		t.Run(issue.name, func(t *testing.T) {
			err := execIssueInsert(conn, map[string]any{
				"id":          "FID-" + issue.name,
				"kind":        "task",
				"title":       issue.name,
				"workflow_id": issue.workflowID,
				"status":      issue.status,
			})
			if err == nil {
				t.Fatal("insert error = nil, want foreign key rejection")
			}
		})
	}

	if err := execIssueInsert(conn, map[string]any{
		"id":          "FID-valid-custom-pair",
		"kind":        "task",
		"title":       "Valid custom workflow pair",
		"workflow_id": "custom-test",
		"status":      "review",
	}); err != nil {
		t.Fatalf("insert valid workflow/state pair: %v", err)
	}

	for _, query := range []string{
		`UPDATE issues SET workflow_id = 'missing-workflow' WHERE id = 'FID-valid-custom-pair'`,
		`UPDATE issues SET status = 'missing-state' WHERE id = 'FID-valid-custom-pair'`,
	} {
		if _, err := conn.Exec(query); err == nil {
			t.Fatalf("Exec(%q) error = nil, want foreign key rejection", query)
		}
	}

	assertIssuesWorkflowForeignKeys(t, conn)
	assertForeignKeyCheckClean(t, conn)
}

func TestMigration0008PreservesIssuesAndDependentObjectsThroughDownUp(t *testing.T) {
	t.Parallel()

	dbPath := filepath.Join(t.TempDir(), "fid.sqlite")
	migrator, conn := newTestMigrator(t, dbPath, testMigrationsDir(t))
	defer closeMigrator(t, migrator, conn)

	if err := migrator.Migrate(7); err != nil {
		t.Fatalf("Migrate(7) error = %v", err)
	}
	seedMigration0008Fixture(t, conn)

	beforeRows := issueSnapshot(t, conn, "issues")
	beforeViewRows := issueSnapshot(t, conn, "issues_enriched")
	beforeFingerprint := schemaFingerprintExceptIssuesTable(t, conn)

	if err := migrator.Migrate(8); err != nil {
		t.Fatalf("Migrate(8) error = %v", err)
	}

	if got := issueSnapshot(t, conn, "issues"); !reflect.DeepEqual(got, beforeRows) {
		t.Fatalf("issues snapshot after migration = %#v, want %#v", got, beforeRows)
	}
	if got := issueSnapshot(t, conn, "issues_enriched"); !reflect.DeepEqual(got, beforeViewRows) {
		t.Fatalf("issues_enriched snapshot after migration = %#v, want %#v", got, beforeViewRows)
	}
	if got := schemaFingerprintExceptIssuesTable(t, conn); !reflect.DeepEqual(got, beforeFingerprint) {
		t.Fatalf("sqlite_schema fingerprint after migration = %#v, want %#v", got, beforeFingerprint)
	}

	assertIssuesWorkflowForeignKeys(t, conn)
	assertIssuesWorkflowChecksReplaced(t, conn)
	assertIssuesAutoIndexes(t, conn)
	assertForeignKeyCheckClean(t, conn)
	assertRetainedIssueReferenceTriggersFunctional(t, conn)

	if err := migrator.Migrate(7); err != nil {
		t.Fatalf("Migrate(7) down error = %v", err)
	}
	if err := migrator.Migrate(8); err != nil {
		t.Fatalf("Migrate(8) second up error = %v", err)
	}

	if got := schemaFingerprintExceptIssuesTable(t, conn); !reflect.DeepEqual(got, beforeFingerprint) {
		t.Fatalf("sqlite_schema fingerprint after down/up = %#v, want %#v", got, beforeFingerprint)
	}
	if got := issueSnapshot(t, conn, "issues"); !reflect.DeepEqual(got, beforeRows) {
		t.Fatalf("issues snapshot after down/up = %#v, want %#v", got, beforeRows)
	}
	assertIssuesWorkflowForeignKeys(t, conn)
	assertIssuesAutoIndexes(t, conn)
	assertForeignKeyCheckClean(t, conn)
}

func seedMigration0008Fixture(t *testing.T, conn *sql.DB) {
	t.Helper()

	issues := []struct {
		id, kind, parentID, title, workflowID, status, specPath, dueDate, ownerRole, ownerActor, createdAt, updatedAt string
		priority, sequence, estTokens                                                                                 any
		inInbox                                                                                                       bool
	}{
		{
			id: "FID-epic-0001", kind: "epic", title: "Migration fixture epic", workflowID: "epic-default", status: "todo",
			priority: 5, specPath: "project-docs/epics/FID-epic-0001.md", dueDate: "2026-07-31", inInbox: true,
			ownerRole: "tl", ownerActor: "human:mark", createdAt: "2026-07-19T01:00:00Z", updatedAt: "2026-07-19T01:01:00Z",
		},
		{
			id: "FID-task-0001", kind: "task", parentID: "FID-epic-0001", title: "Migration fixture task", workflowID: "task-default", status: "in_progress",
			priority: 0, sequence: 1, specPath: "project-docs/tasks/FID-task-0001.md", estTokens: 1234, dueDate: "2026-07-30",
			ownerRole: "developer", ownerActor: "agent", createdAt: "2026-07-19T01:02:00Z", updatedAt: "2026-07-19T01:03:00Z",
		},
		{
			id: "FID-subtask-0001", kind: "subtask", parentID: "FID-task-0001", title: "Migration fixture subtask", workflowID: "subtask-default", status: "done",
			priority: 2, sequence: 1, estTokens: 21, inInbox: true, ownerRole: "developer", createdAt: "2026-07-19T01:04:00Z", updatedAt: "2026-07-19T01:05:00Z",
		},
	}

	for _, issue := range issues {
		if _, err := conn.Exec(`
			INSERT INTO issues (
				id, kind, parent_id, title, workflow_id, status, priority, sequence,
				spec_md_path, est_tokens, due_date, in_inbox, owner_role, owner_actor,
				created_at, updated_at
			) VALUES (?, ?, NULLIF(?, ''), ?, ?, ?, ?, ?, NULLIF(?, ''), ?, NULLIF(?, ''), ?, NULLIF(?, ''), NULLIF(?, ''), ?, ?)
		`, issue.id, issue.kind, issue.parentID, issue.title, issue.workflowID, issue.status, issue.priority, issue.sequence,
			issue.specPath, issue.estTokens, issue.dueDate, issue.inInbox, issue.ownerRole, issue.ownerActor, issue.createdAt, issue.updatedAt); err != nil {
			t.Fatalf("insert fixture issue %q: %v", issue.id, err)
		}
		if issue.id == "FID-epic-0001" {
			if _, err := conn.Exec(`INSERT INTO epic_metadata (issue_id, scheduling) VALUES (?, ?)`, issue.id, "sequential"); err != nil {
				t.Fatalf("insert fixture epic metadata: %v", err)
			}
		}
	}
}

func issueSnapshot(t *testing.T, conn *sql.DB, table string) []string {
	t.Helper()

	columns := `quote(id), quote(kind), quote(parent_id), quote(title), quote(workflow_id), quote(status), quote(priority), quote(sequence), quote(spec_md_path), quote(est_tokens), quote(due_date), quote(in_inbox), quote(owner_role), quote(owner_actor), quote(created_at), quote(updated_at)`
	if table == "issues_enriched" {
		columns += `, quote(concerns_raised), quote(open_concerns_count), quote(triage_status)`
	}
	rows, err := conn.Query(`SELECT ` + columns + ` FROM ` + table + ` ORDER BY id`)
	if err != nil {
		t.Fatalf("snapshot %s: %v", table, err)
	}
	defer closeRows(t, rows)

	var snapshot []string
	for rows.Next() {
		values := make([]string, 16)
		if table == "issues_enriched" {
			values = make([]string, 19)
		}
		pointers := make([]any, len(values))
		for index := range values {
			pointers[index] = &values[index]
		}
		if err := rows.Scan(pointers...); err != nil {
			t.Fatalf("scan %s snapshot: %v", table, err)
		}
		snapshot = append(snapshot, strings.Join(values, "\x1f"))
	}
	if err := rows.Err(); err != nil {
		t.Fatalf("iterate %s snapshot: %v", table, err)
	}
	return snapshot
}

type schemaFingerprintRow struct {
	typeName string
	name     string
	table    string
	sql      string
}

func schemaFingerprintExceptIssuesTable(t *testing.T, conn *sql.DB) []schemaFingerprintRow {
	t.Helper()

	rows, err := conn.Query(`
		SELECT type, name, tbl_name, COALESCE(sql, '')
		FROM sqlite_schema
		WHERE NOT (type = 'table' AND name = 'issues')
		ORDER BY type, name
	`)
	if err != nil {
		t.Fatalf("query sqlite_schema fingerprint: %v", err)
	}
	defer closeRows(t, rows)

	var fingerprint []schemaFingerprintRow
	for rows.Next() {
		var row schemaFingerprintRow
		if err := rows.Scan(&row.typeName, &row.name, &row.table, &row.sql); err != nil {
			t.Fatalf("scan sqlite_schema fingerprint: %v", err)
		}
		fingerprint = append(fingerprint, row)
	}
	if err := rows.Err(); err != nil {
		t.Fatalf("iterate sqlite_schema fingerprint: %v", err)
	}
	return fingerprint
}

func assertIssuesWorkflowForeignKeys(t *testing.T, conn *sql.DB) {
	t.Helper()

	type foreignKeyColumn struct{ from, to string }
	columnsByID := make(map[int][]foreignKeyColumn)
	tableByID := make(map[int]string)
	rows, err := conn.Query(`PRAGMA foreign_key_list(issues)`)
	if err != nil {
		t.Fatalf("PRAGMA foreign_key_list(issues): %v", err)
	}
	defer closeRows(t, rows)

	for rows.Next() {
		var id, sequence int
		var table, from, to, onUpdate, onDelete, match string
		if err := rows.Scan(&id, &sequence, &table, &from, &to, &onUpdate, &onDelete, &match); err != nil {
			t.Fatalf("scan issues foreign key: %v", err)
		}
		tableByID[id] = table
		columnsByID[id] = append(columnsByID[id], foreignKeyColumn{from: from, to: to})
	}
	if err := rows.Err(); err != nil {
		t.Fatalf("iterate issues foreign keys: %v", err)
	}

	want := map[string][]foreignKeyColumn{
		"workflows":       {{from: "workflow_id", to: "id"}},
		"workflow_states": {{from: "workflow_id", to: "workflow_id"}, {from: "status", to: "state"}},
	}
	for table, wantColumns := range want {
		found := false
		for id, gotColumns := range columnsByID {
			if tableByID[id] == table && reflect.DeepEqual(gotColumns, wantColumns) {
				found = true
				break
			}
		}
		if !found {
			t.Fatalf("issues foreign keys = %#v to %#v, want %s %#v", columnsByID, tableByID, table, wantColumns)
		}
	}
}

func assertIssuesWorkflowChecksReplaced(t *testing.T, conn *sql.DB) {
	t.Helper()

	var schemaSQL string
	if err := conn.QueryRow(`SELECT sql FROM sqlite_schema WHERE type = 'table' AND name = 'issues'`).Scan(&schemaSQL); err != nil {
		t.Fatalf("load issues schema: %v", err)
	}
	normalized := strings.ToUpper(strings.Join(strings.Fields(schemaSQL), " "))
	if strings.Contains(normalized, "WORKFLOW_ID =") || strings.Contains(normalized, "STATUS IN (") {
		t.Fatalf("issues schema retains workflow enum CHECK: %s", normalized)
	}
	for _, fragment := range []string{
		"CHECK (kind IN ('epic', 'task', 'subtask'))",
		"CHECK (priority IS NULL OR (priority >= 0 AND priority <= 5))",
		"CHECK (kind != 'epic' OR est_tokens IS NULL)",
		"CHECK (sequence IS NULL OR sequence > 0)",
		"UNIQUE (parent_id, sequence)",
	} {
		if !strings.Contains(strings.Join(strings.Fields(schemaSQL), " "), strings.Join(strings.Fields(fragment), " ")) {
			t.Fatalf("issues schema missing preserved constraint %q: %s", fragment, schemaSQL)
		}
	}
}

func assertIssuesAutoIndexes(t *testing.T, conn *sql.DB) {
	t.Helper()

	rows, err := conn.Query(`PRAGMA index_list(issues)`)
	if err != nil {
		t.Fatalf("PRAGMA index_list(issues): %v", err)
	}
	defer closeRows(t, rows)

	var indexes []string
	for rows.Next() {
		var sequence, unique, partial int
		var name, origin string
		if err := rows.Scan(&sequence, &name, &unique, &origin, &partial); err != nil {
			t.Fatalf("scan issues index list: %v", err)
		}
		columns := issueIndexColumns(t, conn, name)
		if unique == 1 && reflect.DeepEqual(columns, []string{"id"}) && origin == "pk" {
			indexes = append(indexes, "pk")
		}
		if unique == 1 && reflect.DeepEqual(columns, []string{"parent_id", "sequence"}) && origin == "u" {
			indexes = append(indexes, "unique-parent-sequence")
		}
	}
	if err := rows.Err(); err != nil {
		t.Fatalf("iterate issues index list: %v", err)
	}
	sort.Strings(indexes)
	if !reflect.DeepEqual(indexes, []string{"pk", "unique-parent-sequence"}) {
		t.Fatalf("issues structural autoindexes = %#v, want PK and UNIQUE(parent_id, sequence)", indexes)
	}
}

func issueIndexColumns(t *testing.T, conn *sql.DB, index string) []string {
	t.Helper()

	rows, err := conn.Query(`PRAGMA index_info('` + index + `')`)
	if err != nil {
		t.Fatalf("PRAGMA index_info(%q): %v", index, err)
	}
	defer closeRows(t, rows)

	var columns []string
	for rows.Next() {
		var sequence, columnID int
		var name string
		if err := rows.Scan(&sequence, &columnID, &name); err != nil {
			t.Fatalf("scan index info for %q: %v", index, err)
		}
		columns = append(columns, name)
	}
	if err := rows.Err(); err != nil {
		t.Fatalf("iterate index info for %q: %v", index, err)
	}
	return columns
}

func assertForeignKeyCheckClean(t *testing.T, conn *sql.DB) {
	t.Helper()

	rows, err := conn.Query(`PRAGMA foreign_key_check`)
	if err != nil {
		t.Fatalf("PRAGMA foreign_key_check: %v", err)
	}
	defer closeRows(t, rows)
	if rows.Next() {
		t.Fatal("PRAGMA foreign_key_check returned a violation")
	}
	if err := rows.Err(); err != nil {
		t.Fatalf("iterate foreign_key_check: %v", err)
	}
}

func assertRetainedIssueReferenceTriggersFunctional(t *testing.T, conn *sql.DB) {
	t.Helper()

	tests := []struct {
		name, query, want string
		args              []any
	}{
		{"acceptance criteria kind guard", `INSERT INTO acceptance_criteria (id, issue_id, ordinal, text) VALUES (?, ?, ?, ?)`, "acceptance_criteria: issue_id must reference a task or subtask", []any{"AC-fk-rebuild", "FID-epic-0001", 1, "Invalid"}},
		{"task completion kind guard", `INSERT INTO task_completions (id, issue_id, attempt_number, completed_at) VALUES (?, ?, ?, ?)`, "task_completions: issue_id must reference a task or subtask", []any{"TC-fk-rebuild", "FID-epic-0001", 1, "2026-07-19T02:00:00Z"}},
		{"epic metadata kind guard", `INSERT INTO epic_metadata (issue_id, scheduling) VALUES (?, ?)`, "epic_metadata: issue_id must reference an epic", []any{"FID-task-0001", "sequential"}},
		{"triage epic kind guard", `INSERT INTO triage_events (id, epic_id, issue_id, triager_role, source, verdict, created_at) VALUES (?, ?, ?, ?, ?, ?, ?)`, "triage_events: epic_id must reference an epic", []any{"TRG-fk-rebuild", "FID-task-0001", "FID-task-0001", "qa_lead", "task_completion", "BLOCK", "2026-07-19T02:01:00Z"}},
		{"decision epic kind guard", `INSERT INTO decisions (id, epic_id, author_role, summary, body_md_path, status) VALUES (?, ?, ?, ?, ?, ?)`, "decisions: epic_id must reference an epic when set", []any{"DEC-fk-rebuild", "FID-task-0001", "tl", "Invalid", "project-docs/decisions/DEC-fk-rebuild.md", "open"}},
		{"link from ref guard", `INSERT INTO links (from_type, from_id, relation, to_type, to_id) VALUES (?, ?, ?, ?, ?)`, "links: from endpoint must resolve to an existing row", []any{"issue", "FID-missing", "relates_to", "issue", "FID-task-0001"}},
		{"link to ref guard", `INSERT INTO links (from_type, from_id, relation, to_type, to_id) VALUES (?, ?, ?, ?, ?)`, "links: to endpoint must resolve to an existing row", []any{"issue", "FID-task-0001", "relates_to", "issue", "FID-missing"}},
		{"claim ref guard", `INSERT INTO claims (ref_type, ref_id, claimed_by, claimed_at) VALUES (?, ?, ?, ?)`, "claims: ref_type/ref_id must resolve to an existing row", []any{"issue", "FID-missing", "developer", "2026-07-19T02:02:00Z"}},
		{"activity ref guard", `INSERT INTO activity_log (id, journal_event_id, ref_type, ref_id, event, occurred_at) VALUES (?, ?, ?, ?, ?, ?)`, "activity_log: ref_type/ref_id must resolve to an existing row", []any{"ACT-fk-rebuild", "EVT-20260719-000001", "issue", "FID-missing", "status_changed", "2026-07-19T02:03:00Z"}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := conn.Exec(test.query, test.args...)
			if err == nil || !strings.Contains(err.Error(), test.want) {
				t.Fatalf("Exec() error = %v, want %q", err, test.want)
			}
		})
	}

	if _, err := conn.Exec(`INSERT INTO claims (ref_type, ref_id, claimed_by, claimed_at) VALUES (?, ?, ?, ?)`, "issue", "FID-task-0001", "developer", "2026-07-19T02:04:00Z"); err != nil {
		t.Fatalf("insert valid claim: %v", err)
	}
	if _, err := conn.Exec(`UPDATE claims SET ref_id = ? WHERE ref_type = ? AND ref_id = ?`, "FID-missing", "issue", "FID-task-0001"); err == nil || !strings.Contains(err.Error(), "claims: ref_type/ref_id must resolve to an existing row") {
		t.Fatalf("update claim dangling ref error = %v, want ref-integrity rejection", err)
	}
}
