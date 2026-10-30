package db

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestMigration0019AddsOnlyFieldDigestColumn(t *testing.T) {
	t.Parallel()
	migrator, conn := newTestMigrator(t, filepath.Join(t.TempDir(), "fid.sqlite"), testMigrationsDir(t))
	defer closeMigrator(t, migrator, conn)
	if err := migrator.Migrate(18); err != nil {
		t.Fatal(err)
	}
	before := migration0018SchemaObjects(t, conn)
	if err := migrator.Migrate(19); err != nil {
		t.Fatal(err)
	}
	after := migration0018SchemaObjects(t, conn)
	key := "table:idempotency_results:idempotency_results"
	want := before[key]
	column := "request_field_digests_json TEXT NOT NULL CHECK (json_valid(request_field_digests_json))"
	expectedUp := strings.Replace(want, "created_at DATETIME NOT NULL,", "created_at DATETIME NOT NULL, "+column+",", 1)
	if after[key] != expectedUp {
		t.Errorf("migration 19 idempotency DDL differs:\ngot %s\nwant %s", after[key], expectedUp)
	}
	delete(before, key)
	delete(after, key)
	if len(before) != len(after) {
		t.Errorf("other schema object count changed: %d to %d", len(before), len(after))
	}
	for name, sql := range before {
		if after[name] != sql {
			t.Errorf("other schema object %s changed", name)
		}
	}
	if err := migrator.Migrate(18); err != nil {
		t.Fatal(err)
	}
	restored := migration0018SchemaObjects(t, conn)
	before[key] = want
	if len(restored) != len(before) {
		t.Errorf("down schema object count = %d, want %d", len(restored), len(before))
	}
	for name, sql := range before {
		if restored[name] != sql {
			t.Errorf("down schema object %s differs:\ngot %s\nwant %s", name, restored[name], sql)
		}
	}
}
