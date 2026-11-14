//go:build integration

package integration

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/mbell-nzl/fid/internal/config"
	"github.com/mbell-nzl/fid/internal/gittest"
)

func p24InternalPair(t *testing.T, _ string) p24CodePair {
	t.Helper()
	root := t.TempDir()
	// A directory is not an export archive file. This is a real OS read error,
	// not a fabricated Failure or domain-mapper hook.
	return p24CodePair{root, "export.validate", []string{"export", "validate", root}}
}

func p24JournalPair(t *testing.T, _ string) p24CodePair {
	t.Helper()
	root := gittest.NewWorkspace(t)
	p21Run(t, integrationFIDBinary(t), root, "init", "--prefix", "P24", "--json")
	if err := os.WriteFile(filepath.Join(config.LedgerDir(root), "journal", "9999-12-01.jsonl"), []byte("invalid JSON\n"), 0600); err != nil {
		t.Fatal(err)
	}
	return p24CodePair{root, "project.rebuild", []string{"rebuild"}}
}

func p24LockPair(t *testing.T, _ string) p24CodePair {
	t.Helper()
	root := gittest.NewWorkspace(t)
	p21Run(t, integrationFIDBinary(t), root, "init", "--prefix", "P24", "--json")
	if err := os.Mkdir(filepath.Join(config.LedgerDir(root), "write.lock"), 0700); err != nil {
		t.Fatal(err)
	}
	return p24CodePair{root, "issue.create", []string{"issue", "create", "--kind", "task", "--title", "locked", "--actor", "agent:a", "--idempotency-key", "locked"}}
}

func p24ProjectionPair(t *testing.T, _ string) p24CodePair {
	t.Helper()
	root := gittest.NewWorkspace(t)
	binary := integrationFIDBinary(t)
	p21Run(t, binary, root, "init", "--prefix", "FID", "--json")
	conn := p22FaultDB(t, root)
	conn.SetMaxOpenConns(1)
	_, err := conn.Exec("PRAGMA foreign_keys=OFF")
	p22FaultOK(t, err)
	_, err = conn.Exec("PRAGMA ignore_check_constraints=ON")
	p22FaultOK(t, err)
	_, err = conn.Exec("INSERT INTO issues (id,kind,title,workflow_id,status,created_at,updated_at) VALUES ('P24-corrupt','task','Corrupt workflow','missing','todo','2026-01-01T00:00:00Z','2026-01-01T00:00:00Z')")
	p22FaultOK(t, err)
	_, err = conn.Exec("PRAGMA ignore_check_constraints=OFF")
	p22FaultOK(t, err)
	_, err = conn.Exec("PRAGMA foreign_keys=ON")
	p22FaultOK(t, err)
	p22FaultOK(t, conn.Close())
	return p24CodePair{root, "maintenance.check", []string{"maintenance", "check"}}
}

func p24ImportPair(t *testing.T, _ string) p24CodePair {
	t.Helper()
	root := gittest.NewWorkspace(t)
	binary := integrationFIDBinary(t)
	p21Run(t, binary, root, "init", "--prefix", "FID", "--json")
	archive := filepath.Join(t.TempDir(), "history.zip")
	p21Run(t, binary, root, "--json", "export", "create", archive)
	return p24CodePair{root, "export.import", []string{"export", "import", archive, "--repository", root, "--idempotency-key", "code-pair"}}
}
