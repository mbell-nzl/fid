package db

import (
	"database/sql"
	"fmt"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestMigration0018DeletesOnlyLegacyNonIssueReferences(t *testing.T) {
	t.Parallel()
	migrator, conn := newTestMigrator(t, filepath.Join(t.TempDir(), "fid.sqlite"), testMigrationsDir(t))
	defer closeMigrator(t, migrator, conn)
	if err := migrator.Migrate(17); err != nil {
		t.Fatal(err)
	}
	for _, query := range []string{
		`INSERT INTO issues(id,kind,title,workflow_id,status) VALUES('ISS-1','task','one','task-default','todo')`,
		`INSERT INTO issues(id,kind,title,workflow_id,status) VALUES('ISS-2','task','two','task-default','todo')`,
		`INSERT INTO requirements(id,type,summary,source_doc) VALUES('REQ-1','functional','one','source')`,
		`INSERT INTO decisions(id,author_role,summary,body_md_path,status) VALUES('DEC-1','developer','one','decision.md','accepted')`,
		`INSERT INTO links(from_type,from_id,relation,to_type,to_id) VALUES('issue','ISS-1','blocks','issue','ISS-2')`,
		`INSERT INTO links(from_type,from_id,relation,to_type,to_id) VALUES('issue','ISS-1','implements','requirement','REQ-1')`,
		`INSERT INTO links(from_type,from_id,relation,to_type,to_id) VALUES('requirement','REQ-1','implemented_by','issue','ISS-1')`,
		`INSERT INTO activity_log(id,journal_event_id,ref_type,ref_id,event) VALUES('ACT-I','event-i','issue','ISS-1','issue-created')`,
		`INSERT INTO activity_log(id,journal_event_id,ref_type,ref_id,event) VALUES('ACT-R','event-r','requirement','REQ-1','requirement-synced')`,
		`INSERT INTO activity_log(id,journal_event_id,ref_type,ref_id,event) VALUES('ACT-D','event-d','decision','DEC-1','decision-created')`,
	} {
		if _, err := conn.Exec(query); err != nil {
			t.Fatalf("seed 17: %s: %v", query, err)
		}
	}
	if err := migrator.Migrate(18); err != nil {
		t.Fatal(err)
	}
	for _, probe := range []struct {
		query string
		want  []string
	}{
		{`SELECT from_type || ':' || from_id || ':' || relation || ':' || to_type || ':' || to_id FROM links ORDER BY 1`, []string{"issue:ISS-1:blocks:issue:ISS-2"}},
		{`SELECT id FROM activity_log ORDER BY id`, []string{"ACT-I"}},
	} {
		rows, err := conn.Query(probe.query)
		if err != nil {
			t.Fatal(err)
		}
		var got []string
		for rows.Next() {
			var value string
			if err := rows.Scan(&value); err != nil {
				t.Fatal(err)
			}
			got = append(got, value)
		}
		if err := rows.Close(); err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(got, probe.want) {
			t.Errorf("%s = %v, want %v", probe.query, got, probe.want)
		}
	}
	assertForeignKeyCheckClean(t, conn)
	assertIntegrityCheckClean(t, conn)
}

func TestMigration0018CurrentReferenceGuards(t *testing.T) {
	t.Parallel()
	migrator, conn := newTestMigrator(t, filepath.Join(t.TempDir(), "fid.sqlite"), testMigrationsDir(t))
	defer closeMigrator(t, migrator, conn)
	if err := migrator.Migrate(18); err != nil {
		t.Fatal(err)
	}
	for _, query := range []string{
		`INSERT INTO issues(id,kind,title,workflow_id,status) VALUES('ISS-1','task','one','task-default','todo')`,
		`INSERT INTO issues(id,kind,title,workflow_id,status) VALUES('ISS-2','task','two','task-default','todo')`,
		`INSERT INTO requirement_refs(id,requirement_id,source_json,status,created_at,created_event_id) VALUES('RREF-1','REQ-1','{}','active','now','event-r')`,
		`INSERT INTO decision_refs(id,class,source_json,authority_actor_id,authority_capabilities_json,outcome,scope_issue_id,created_at,created_event_id) VALUES('DREF-1','release','{}','actor','[]','approved','ISS-1','now','event-d')`,
	} {
		if _, err := conn.Exec(query); err != nil {
			t.Fatalf("seed 18: %s: %v", query, err)
		}
	}
	for _, probe := range []struct {
		name, query, errorText string
	}{
		{"issue link", `INSERT INTO links(from_type,from_id,relation,to_type,to_id) VALUES('issue','ISS-1','blocks','issue','ISS-2')`, ""},
		{"nonissue link", `INSERT INTO links(from_type,from_id,relation,to_type,to_id) VALUES('issue','ISS-1','relates_to','requirement','ISS-2')`, "endpoint types must be issue"},
		{"missing source", `INSERT INTO links(from_type,from_id,relation,to_type,to_id) VALUES('issue','MISSING','blocks','issue','ISS-2')`, "from endpoint must resolve"},
		{"missing target", `INSERT INTO links(from_type,from_id,relation,to_type,to_id) VALUES('issue','ISS-1','blocks','issue','MISSING')`, "to endpoint must resolve"},
		{"invalid relation", `INSERT INTO links(from_type,from_id,relation,to_type,to_id) VALUES('issue','ISS-1','implements','issue','ISS-2')`, "relation/endpoint combination"},
		{"issue activity", `INSERT INTO activity_log(id,journal_event_id,ref_type,ref_id,event) VALUES('ACT-I','i','issue','ISS-1','issue-created')`, ""},
		{"requirement activity", `INSERT INTO activity_log(id,journal_event_id,ref_type,ref_id,event) VALUES('ACT-R','r','requirement','REQ-1','requirement-referenced')`, ""},
		{"decision activity", `INSERT INTO activity_log(id,journal_event_id,ref_type,ref_id,event) VALUES('ACT-D','d','decision','DREF-1','decision-referenced')`, ""},
		{"missing requirement", `INSERT INTO activity_log(id,journal_event_id,ref_type,ref_id,event) VALUES('ACT-BAD-R','r','requirement','MISSING','requirement-referenced')`, "must resolve"},
		{"missing decision", `INSERT INTO activity_log(id,journal_event_id,ref_type,ref_id,event) VALUES('ACT-BAD-D','d','decision','MISSING','decision-referenced')`, "must resolve"},
		{"invalid type", `INSERT INTO activity_log(id,journal_event_id,ref_type,ref_id,event) VALUES('ACT-BAD-T','t','envelope','ISS-1','test')`, "must resolve"},
	} {
		_, err := conn.Exec(probe.query)
		if probe.errorText == "" && err != nil || probe.errorText != "" && (err == nil || !strings.Contains(err.Error(), probe.errorText)) {
			t.Errorf("%s: error = %v, want %q", probe.name, err, probe.errorText)
		}
	}
	assertForeignKeyCheckClean(t, conn)
	assertIntegrityCheckClean(t, conn)
}

func TestMigration0018IssuesEnrichedColumnsAndSubjectCount(t *testing.T) {
	t.Parallel()
	migrator, conn := newTestMigrator(t, filepath.Join(t.TempDir(), "fid.sqlite"), testMigrationsDir(t))
	defer closeMigrator(t, migrator, conn)
	if err := migrator.Migrate(17); err != nil {
		t.Fatal(err)
	}
	before := migration0018ColumnNames(t, conn, "issues_enriched")
	if err := migrator.Migrate(18); err != nil {
		t.Fatal(err)
	}
	after := migration0018ColumnNames(t, conn, "issues_enriched")
	if !reflect.DeepEqual(before, after) {
		t.Fatalf("issues_enriched columns changed across 0018: before=%v after=%v", before, after)
	}
	if _, err := conn.Exec(`INSERT INTO issues(id,kind,title,workflow_id,status) VALUES('ISS-1','task','one','task-default','todo')`); err != nil {
		t.Fatal(err)
	}
	for _, subject := range []string{"issue", "requirement"} {
		query := fmt.Sprintf(`INSERT INTO concerns(id,subject_type,subject_id,actor_id,scope_json,severity,summary,reasoning,references_json,status,created_at,created_event_id) VALUES('CON-%s','%s','ISS-1','actor','{"paths":[]}','minor','summary','reason','[]','open','now','event-%s')`, subject, subject, subject)
		if _, err := conn.Exec(query); err != nil {
			t.Fatal(err)
		}
	}
	var count int
	if err := conn.QueryRow(`SELECT open_concerns_count FROM issues_enriched WHERE id='ISS-1'`).Scan(&count); err != nil {
		t.Fatal(err)
	}
	if count != 1 {
		t.Fatalf("Issue open Concerns count = %d, want 1", count)
	}
}

func migration0018ColumnNames(t *testing.T, conn *sql.DB, table string) []string {
	t.Helper()
	rows, err := conn.Query(`SELECT name FROM pragma_table_info(?) ORDER BY cid`, table)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = rows.Close() }()
	var names []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			t.Fatal(err)
		}
		names = append(names, name)
	}
	if err := rows.Err(); err != nil {
		t.Fatal(err)
	}
	return names
}

func TestLatestSchemaHasNoV1JudgmentTables(t *testing.T) {
	t.Parallel()
	migrator, conn := newTestMigrator(t, filepath.Join(t.TempDir(), "fid.sqlite"), testMigrationsDir(t))
	defer closeMigrator(t, migrator, conn)
	if err := migrator.Migrate(ActiveMigrationVersion); err != nil {
		t.Fatalf("migrate latest: %v", err)
	}
	for _, name := range []string{"task_completions", "ac_verifications", "requirements", "decisions", "judgment_concerns", "concerns", "requirement_refs", "decision_refs"} {
		if _, exists := schemaObjectSQL(t, conn, "table", name); exists {
			t.Errorf("removed V1 table %q survives latest schema", name)
		}
	}
	var count int
	for _, object := range []struct{ kind, name string }{
		{"index", "idx_concerns_issue_status"}, {"index", "judgment_concerns_by_subject"},
		{"trigger", "decisions_epic_kind_guard_insert"}, {"trigger", "task_completions_issue_kind_guard_insert"},
	} {
		if _, exists := schemaObjectSQL(t, conn, object.kind, object.name); exists {
			t.Errorf("removed %s %q survives latest schema", object.kind, object.name)
		}
	}
	for _, object := range []struct{ table, column string }{
		{"agent_runs", "completion_id"}, {"concerns", "completion_id"}, {"requirement_refs", "last_synced_hash"},
	} {
		if migration0018HasColumn(t, conn, object.table, object.column) {
			t.Errorf("removed column %s.%s survives latest schema", object.table, object.column)
		}
	}
	for _, rule := range []string{"all_acs_verified", "has_completion_report", "concerns_addressed"} {
		if err := conn.QueryRow(`SELECT COUNT(*) FROM workflow_transition_requires WHERE rule=?`, rule).Scan(&count); err != nil {
			t.Fatal(err)
		}
		if count != 0 {
			t.Errorf("removed transition rule %q has %d rows", rule, count)
		}
	}
	for _, name := range []string{"issues_enriched", "links_endpoint_scope_guard_insert", "links_from_ref_integrity_guard_insert", "links_to_ref_integrity_guard_insert", "links_relation_compatibility_guard_insert", "activity_log_ref_integrity_guard_insert"} {
		kind := "trigger"
		switch name {
		case "issues_enriched":
			kind = "view"
		case "concerns_by_subject_status":
			kind = "index"
		}
		if _, exists := schemaObjectSQL(t, conn, kind, name); !exists {
			t.Errorf("required %s %q is missing", kind, name)
		}
	}
}

func migration0018HasColumn(t *testing.T, conn *sql.DB, table, column string) bool {
	t.Helper()
	var count int
	if err := conn.QueryRow(`SELECT COUNT(*) FROM pragma_table_info(?) WHERE name=?`, table, column).Scan(&count); err != nil {
		t.Fatal(err)
	}
	return count != 0
}

func TestMigration0018DownRestoresEmptyV1Tables(t *testing.T) {
	t.Parallel()
	migrator, conn := newTestMigrator(t, filepath.Join(t.TempDir(), "fid.sqlite"), testMigrationsDir(t))
	defer closeMigrator(t, migrator, conn)
	if err := migrator.Migrate(18); err != nil {
		t.Fatalf("migrate 18: %v", err)
	}
	if err := migrator.Migrate(17); err != nil {
		t.Fatalf("migrate down to 17: %v", err)
	}
	for _, name := range []string{"task_completions", "ac_verifications", "requirements", "decisions", "judgment_concerns"} {
		if _, exists := schemaObjectSQL(t, conn, "table", name); !exists {
			t.Errorf("restored V1 table %q is missing", name)
		}
	}
}

func TestMigration0018DownMatchesFreshSchema17(t *testing.T) {
	t.Parallel()
	dir := testMigrationsDir(t)
	fresh, freshDB := newTestMigrator(t, filepath.Join(t.TempDir(), "fresh.sqlite"), dir)
	defer closeMigrator(t, fresh, freshDB)
	if err := fresh.Migrate(17); err != nil {
		t.Fatal(err)
	}
	round, roundDB := newTestMigrator(t, filepath.Join(t.TempDir(), "round.sqlite"), dir)
	defer closeMigrator(t, round, roundDB)
	if err := round.Migrate(18); err != nil {
		t.Fatal(err)
	}
	if err := round.Migrate(17); err != nil {
		t.Fatal(err)
	}
	want := migration0018SchemaObjects(t, freshDB)
	got := migration0018SchemaObjects(t, roundDB)
	if len(got) != len(want) {
		t.Errorf("schema object count = %d, want %d", len(got), len(want))
	}
	for key, expected := range want {
		actual, ok := got[key]
		if !ok {
			t.Errorf("missing schema object %q", key)
			continue
		}
		// SQLite quotes names touched by ALTER TABLE RENAME.
		normalize := func(s string) string {
			return strings.NewReplacer(`"judgment_concerns"`, "judgment_concerns", `"agent_runs"`, "agent_runs").Replace(s)
		}
		if normalize(actual) != normalize(expected) {
			t.Errorf("schema object %q differs:\ngot  %s\nwant %s", key, actual, expected)
		}
	}
}

func migration0018SchemaObjects(t *testing.T, conn *sql.DB) map[string]string {
	t.Helper()
	rows, err := conn.Query(`SELECT type, name, tbl_name, sql FROM sqlite_master WHERE name NOT LIKE 'sqlite_%' AND name != 'schema_migrations' ORDER BY type, name`)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = rows.Close() }()
	objects := map[string]string{}
	for rows.Next() {
		var kind, name, table, definition string
		if err := rows.Scan(&kind, &name, &table, &definition); err != nil {
			t.Fatal(err)
		}
		objects[kind+":"+name+":"+table] = definition
	}
	if err := rows.Err(); err != nil {
		t.Fatal(err)
	}
	return objects
}

func TestMigration0018DownPreservesAgentRunAndFullRollback(t *testing.T) {
	t.Parallel()
	path := filepath.Join(t.TempDir(), "fid.sqlite")
	migrator, conn := newTestMigrator(t, path, testMigrationsDir(t))
	if err := migrator.Migrate(18); err != nil {
		t.Fatal(err)
	}
	if _, err := conn.Exec(`INSERT INTO agent_runs(id,role,model) VALUES('RUN-18','developer','test')`); err != nil {
		t.Fatal(err)
	}
	if err := migrator.Migrate(17); err != nil {
		t.Fatal(err)
	}
	var role, model string
	var completion sql.NullString
	if err := conn.QueryRow(`SELECT role,model,completion_id FROM agent_runs WHERE id='RUN-18'`).Scan(&role, &model, &completion); err != nil {
		t.Fatal(err)
	}
	if role != "developer" || model != "test" || completion.Valid {
		t.Errorf("restored run = %q %q %v", role, model, completion)
	}
	closeMigrator(t, migrator, conn)
	if err := MigrateDown(path, testMigrationsDir(t)); err != nil {
		t.Fatalf("full rollback: %v", err)
	}
}

func TestMigration0018RoundTripIntegrity(t *testing.T) {
	t.Parallel()
	migrator, conn := newTestMigrator(t, filepath.Join(t.TempDir(), "fid.sqlite"), testMigrationsDir(t))
	defer closeMigrator(t, migrator, conn)
	for _, version := range []uint{17, 18, 17, 18} {
		if err := migrator.Migrate(version); err != nil {
			t.Fatalf("migrate %d: %v", version, err)
		}
	}
	assertForeignKeyCheckClean(t, conn)
	assertIntegrityCheckClean(t, conn)
}
