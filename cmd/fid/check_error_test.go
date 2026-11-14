package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/mbell-nzl/fid/internal/config"
)

func TestLeanCheckFidErrorStopsAndLeavesIncompleteStart(t *testing.T) {
	t.Parallel()
	f := leanJourneyFixture(t, leanHelperDefinition(t, "pass", "1m", []string{"coverage.out"})+"\n[checks.zpass]\nrun=['true']\ntimeout='1m'\n")
	path := filepath.Join(config.LedgerDir(f.workspace), "artifacts")
	if err := os.RemoveAll(path); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte("not a directory"), 0o600); err != nil {
		t.Fatal(err)
	}
	args := []string{"--json", "check", "lean", "zpass", "--hold", "HLD-aaaaaa", "--actor", "agent:claude-builder", "--idempotency-key", "fid-error"}
	got := f.run(args...)
	assertJSONFailureCode(t, got.stderr, "check.execute", "internal-error")
	record := latestLeanRecord(t, f)
	if record.Result != nil || record.Start.Name != "lean" {
		t.Fatalf("fid error invented a result or started later work: %#v", record)
	}
	got = f.run(args...)
	assertJSONFailureCode(t, got.stderr, "check.execute", "check-incomplete")
}
