//go:build integration

package integration

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/mbell-nzl/fid/internal/config"
	fiddb "github.com/mbell-nzl/fid/internal/db"
	"github.com/mbell-nzl/fid/internal/gittest"
	"github.com/mbell-nzl/fid/internal/journal"
)

// Migrate and replay the identical genuine prefix once. Each negative still
// gets its own physical DB and independent observer; only immutable setup bytes
// are shared. Positive controls and actual rebuild checks retain their paths.
func p22ForceEventChecker(t *testing.T, binary, source string, prefix []byte) func(*testing.T, journal.Event, string, bool) {
	t.Helper()
	setup := t.TempDir()
	dbPath := filepath.Join(setup, "prefix.sqlite")
	p22FaultOK(t, fiddb.MigrateUp(dbPath, filepath.Join(integrationRepositoryRoot(), "migrations")))
	conn, err := fiddb.Open(dbPath)
	p22FaultOK(t, err)
	dir := t.TempDir()
	p22FaultOK(t, os.WriteFile(filepath.Join(dir, "2026-09-09.jsonl"), prefix, 0600))
	_, err = journal.ReplayJournal(conn, dir, nil)
	p22FaultOK(t, err)
	_, err = conn.Exec(`PRAGMA wal_checkpoint(TRUNCATE)`)
	p22FaultOK(t, err)
	p22FaultOK(t, conn.Close())
	pristine, err := os.ReadFile(dbPath)
	p22FaultOK(t, err)
	return func(t *testing.T, event journal.Event, mode string, accepted bool) {
		t.Helper()
		if accepted || mode == "rebuild" {
			p22CheckCurrentEventCase(t, binary, source, prefix, event, mode, accepted)
			return
		}
		workspace := gittest.NewWorkspace(t)
		p22FaultOK(t, os.MkdirAll(config.LedgerDir(workspace), 0700))
		path := filepath.Join(config.LedgerDir(workspace), "fid.sqlite")
		p22FaultOK(t, os.WriteFile(path, pristine, 0600))
		conn, err := fiddb.Open(path)
		p22FaultOK(t, err)
		defer conn.Close()
		observer := p22FaultDB(t, workspace)
		defer observer.Close()
		before := p22FaultProjectionFromDB(t, observer)
		if mode == "live" {
			hash, hashErr := journal.HashEvent(event)
			p22FaultOK(t, hashErr)
			err = journal.ProjectEvent(conn, journal.ProjectionInput{Event: event, EventHash: hash})
		} else {
			raw, marshalErr := json.Marshal(event)
			p22FaultOK(t, marshalErr)
			body := append(append(append([]byte{}, prefix...), raw...), '\n')
			dir := t.TempDir()
			p22FaultOK(t, os.WriteFile(filepath.Join(dir, "2026-09-09.jsonl"), body, 0600))
			_, err = journal.ReplayJournal(conn, dir, &journal.ReplayOptions{StartFile: "2026-09-09.jsonl", StartOffset: int64(len(prefix))})
		}
		if err == nil {
			t.Fatal("malformed current force Event accepted")
		}
		if !reflect.DeepEqual(before, p22FaultProjectionFromDB(t, observer)) {
			t.Fatal("rejected force Event changed projection")
		}
	}
}
