package db

import (
	"path/filepath"
	"testing"
)

func TestEpicValidationMigrationUpDownPreservesExistingDecisions(t *testing.T) {
	t.Parallel()
	m, q := newTestMigrator(t, filepath.Join(t.TempDir(), "fid.sqlite"), testMigrationsFS(t))
	defer closeMigrator(t, m, q)
	if err := m.Migrate(44); err != nil {
		t.Fatal(err)
	}
	for _, sql := range []string{
		`INSERT INTO issues(id,kind,title,workflow_id,status,created_at,updated_at) VALUES('FID-e','epic','Epic','epic-default','todo','2026-11-28T00:00:00Z','2026-11-28T00:00:00Z')`,
		`INSERT INTO person_decisions(id,kind,actor,subject_type,subject_id,reason,binding_id,decided_at) VALUES('PDN-old','close-epic','human:mark','epic','FID-e','Historical close','FID-e','2026-11-28T00:00:00Z')`,
	} {
		if _, err := q.Exec(sql); err != nil {
			t.Fatal(err)
		}
	}
	for _, version := range []uint{45, 44, 45} {
		if err := m.Migrate(version); err != nil {
			t.Fatal(err)
		}
		var reason string
		if err := q.QueryRow(`SELECT reason FROM person_decisions WHERE id='PDN-old'`).Scan(&reason); err != nil || reason != "Historical close" {
			t.Fatalf("old decision=%q,%v", reason, err)
		}
	}
	for _, table := range []string{"epic_validation_changes", "epic_validation_results", "epic_plan_acceptances"} {
		var n int
		if err := q.QueryRow("SELECT count(*) FROM " + table).Scan(&n); err != nil || n != 0 {
			t.Fatalf("%s=%d,%v", table, n, err)
		}
	}
	if _, err := q.Exec(`INSERT INTO person_decisions(id,kind,actor,subject_type,subject_id,reason,binding_id,decided_at,binding_json) VALUES('PDN-new','approve-epic','human:mark','epic','FID-e','Approve','digest','2026-11-28T00:00:00Z','{}')`); err != nil {
		t.Fatal(err)
	}
	if _, err := q.Exec(`UPDATE person_decisions SET binding_json='{"approval_id":"PDN-new","results":[],"merges":[]}' WHERE id='PDN-old'`); err != nil {
		t.Fatal(err)
	}
	if err := m.Migrate(44); err != nil {
		t.Fatal(err)
	}
	var binding any
	if err := q.QueryRow(`SELECT binding_json FROM person_decisions WHERE id='PDN-old'`).Scan(&binding); err != nil || binding != nil {
		t.Fatalf("down close binding=%v,%v", binding, err)
	}
	var n int
	if err := q.QueryRow(`SELECT count(*) FROM person_decisions`).Scan(&n); err != nil || n != 1 {
		t.Fatalf("down decisions=%d,%v", n, err)
	}
}
