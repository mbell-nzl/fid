package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mbell-nzl/fid/internal/config"
)

func TestRestoreRefusesOversizedManifestsBeforeImport(t *testing.T) {
	t.Parallel()
	for _, name := range []string{"manifest.json", "backup.json"} {
		t.Run(name, func(t *testing.T) {
			fixture := newBackupFixture(t)
			id := "project-" + strings.Repeat("a", 32)
			wrapper := `{"schema":"fid.backup.v2","project_id":"` + id + `","archive_digest":"sha256:` + strings.Repeat("0", 64) + `"}`
			exported := `{"project_id":"` + id + `"}`
			files := map[string]string{id + "/backup.json": wrapper, id + "/manifest.json": exported}
			// The existing non-Artifact ceiling applies before either manifest
			// is buffered, even when the eventual ZIP digest would not match.
			files[id+"/"+name] += strings.Repeat(" ", (64<<20)+1-len(files[id+"/"+name]))
			fixture.setTarget(fixture.target)
			fixture.commitToTarget(files)
			workspace := fixture.restoreWorkspace(id)
			failure := decodeRestoreFailure(t, fixture.restore(workspace, "--json"), 6)
			checks, _ := failure.Details["checks"].([]any)
			if failure.Code != "archive-invalid" || len(checks) != 1 || checks[0].(map[string]any)["code"] != "resource-limit" {
				t.Fatalf("failure = %+v, want a pre-import resource-limit refusal", failure)
			}
			if _, err := os.Lstat(config.LedgerDir(workspace)); !os.IsNotExist(err) {
				t.Fatalf("refused restore installed a ledger: %v", err)
			}
			scratch, err := filepath.Glob(filepath.Join(filepath.Dir(config.LedgerDir(workspace)), "fid-restore-*"))
			if err != nil || len(scratch) != 0 {
				t.Fatalf("restore scratch = %v, %v, want cleanup", scratch, err)
			}
		})
	}
}
