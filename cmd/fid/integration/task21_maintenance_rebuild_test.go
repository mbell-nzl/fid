//go:build integration

package integration

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/mbell-nzl/fid/internal/config"
	fiddb "github.com/mbell-nzl/fid/internal/db"
	"github.com/mbell-nzl/fid/internal/gittest"
)

func TestFinalBinaryInitAndRebuildRemainRoleFree(t *testing.T) {
	t.Parallel()
	workspace := gittest.NewWorkspace(t)
	binary := integrationFIDBinary(t)

	closureRun(t, workspace, binary, "--json", "init", "--prefix", "P21")
	journalPath := filepath.Join(config.LedgerDir(workspace), "journal")
	bootstrap := onlyJournalEvent(t, journalPath)
	assertRoleFreeBootstrap(t, bootstrap)
	assertRolesTableAbsent(t, filepath.Join(config.LedgerDir(workspace), "fid.sqlite"))

	capabilitiesBefore, _ := closureRun(t, workspace, binary, "--json", "capabilities")
	if err := os.Remove(filepath.Join(config.LedgerDir(workspace), "fid.sqlite")); err != nil {
		t.Fatalf("delete projection: %v", err)
	}
	closureRun(t, workspace, binary, "--json", "rebuild")
	capabilitiesAfter, _ := closureRun(t, workspace, binary, "--json", "capabilities")

	if capabilitiesAfter != capabilitiesBefore {
		t.Fatalf("capabilities after rebuild = %s, want %s", capabilitiesAfter, capabilitiesBefore)
	}
	assertRoleFreeBootstrap(t, onlyJournalEvent(t, journalPath))
	assertRolesTableAbsent(t, filepath.Join(config.LedgerDir(workspace), "fid.sqlite"))
}

func onlyJournalEvent(t *testing.T, journalDir string) map[string]any {
	t.Helper()
	entries, err := os.ReadDir(journalDir)
	if err != nil {
		t.Fatalf("read journal: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("journal files = %d, want 1", len(entries))
	}
	payload, err := os.ReadFile(filepath.Join(journalDir, entries[0].Name()))
	if err != nil {
		t.Fatalf("read bootstrap journal: %v", err)
	}
	var event map[string]any
	if err := json.NewDecoder(bytes.NewReader(payload)).Decode(&event); err != nil {
		t.Fatalf("decode bootstrap Event: %v", err)
	}
	return event
}

func assertRoleFreeBootstrap(t *testing.T, event map[string]any) {
	t.Helper()
	if event["event_schema"] != float64(2) || event["event_type"] != "project.initialized" {
		t.Fatalf("bootstrap Event = %#v, want Event-v2 project.initialized", event)
	}
	effects, ok := event["effects"].([]any)
	if !ok || len(effects) != 0 {
		t.Fatalf("bootstrap Effects = %#v, want empty", event["effects"])
	}
	actor, ok := event["actor"].(map[string]any)
	if !ok || actor["actor_type"] != "service" || actor["actor_id"] != "fid-bootstrap" {
		t.Fatalf("bootstrap Actor = %#v, want fid bootstrap service Actor", event["actor"])
	}
}

func assertRolesTableAbsent(t *testing.T, dbPath string) {
	t.Helper()
	database, err := fiddb.Open(dbPath)
	if err != nil {
		t.Fatalf("open projection: %v", err)
	}
	defer func() { _ = database.Close() }()

	var tables int
	if err := database.QueryRow(`SELECT COUNT(*) FROM sqlite_schema WHERE type = 'table' AND name = 'roles'`).Scan(&tables); err != nil {
		t.Fatalf("inspect roles table: %v", err)
	}
	if tables != 0 {
		t.Fatalf("roles table count = %d, want 0", tables)
	}
}
