package journal

import (
	"database/sql"
	"os"
	"path/filepath"
	"testing"

	"github.com/mbell-nzl/fid/internal/config"
	fiddb "github.com/mbell-nzl/fid/internal/db"
)

// Keep the migrated connection setup out of the per-refusal path. Production
// admission/replay APIs remain the same; rebuild and freshness still open their
// own databases using the verified current-schema fixture.
type jsonAdmissionWorkspace struct {
	admissionWorkspace
	migrations string
}

func newJSONAdmissionWorkspace(t *testing.T, mode, migrations string) jsonAdmissionWorkspace {
	t.Helper()
	root := newReplayWorkspace(t)
	dir := filepath.Join(config.LedgerDir(root), "journal")
	file := filepath.Join(dir, "2026-11-10.jsonl")
	writeReplayEvents(t, file, admissionPrefix()...)
	path := filepath.Join(config.LedgerDir(root), "fid.sqlite")
	if err := fiddb.MigrateUp(path, migrations); err != nil {
		t.Fatal(err)
	}
	w := jsonAdmissionWorkspace{admissionWorkspace: admissionWorkspace{root: root, journal: dir, file: file}, migrations: migrations}
	conn := w.open(t)
	defer closeProjectDB(t, conn)
	if mode != "replay" {
		if _, err := ReplayJournal(conn, dir, nil); err != nil {
			t.Fatal(err)
		}
	}
	var err error
	w.prefix, err = os.ReadFile(file)
	if err != nil {
		t.Fatal(err)
	}
	return w
}

func (w jsonAdmissionWorkspace) open(t *testing.T) *sql.DB {
	t.Helper()
	conn, err := fiddb.Open(filepath.Join(config.LedgerDir(w.root), "fid.sqlite"))
	if err != nil {
		t.Fatal(err)
	}
	return conn
}

func (w jsonAdmissionWorkspace) snapshot(t *testing.T) map[string][]string {
	t.Helper()
	conn := w.open(t)
	defer closeProjectDB(t, conn)
	return currentProjectionSnapshot(t, conn)
}

func (w jsonAdmissionWorkspace) admit(t *testing.T, mode string, event Event) error {
	t.Helper()
	if mode != "commit" && mode != "dry-run" && mode != "project" {
		if err := os.WriteFile(w.file, w.journalWith(t, event), 0o600); err != nil {
			t.Fatal(err)
		}
		return w.replay(t, mode)
	}
	conn := w.open(t)
	defer closeProjectDB(t, conn)
	switch mode {
	case "commit":
		return WithLedgerLock(config.LedgerDir(w.root), LockOwner{Command: "test"}, func(lock *WriteLock) error { _, err := CommitEvent(lock, conn, event); return err })
	case "dry-run":
		return DryRunProjectEvent(conn, event)
	default:
		hash, err := HashEvent(event)
		if err != nil {
			return err
		}
		return ProjectEvent(conn, ProjectionInput{Event: event, EventHash: hash, File: filepath.Base(w.file), Line: 4, Offset: int64(len(w.prefix)) + 1})
	}
}

func (w jsonAdmissionWorkspace) replay(t *testing.T, mode string) error {
	t.Helper()
	switch mode {
	case "freshen":
		_, err := EnsureProjectionFresh(w.root, w.migrations)
		return err
	case "rebuild":
		_, err := RebuildWorkspace(w.root, w.migrations)
		return err
	}
	conn := w.open(t)
	defer closeProjectDB(t, conn)
	var options *ReplayOptions
	if mode == "catch-up" {
		options = &ReplayOptions{StartFile: filepath.Base(w.file), StartOffset: int64(len(w.prefix))}
	}
	_, err := ReplayJournal(conn, w.journal, options)
	return err
}
