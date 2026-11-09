package main

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mbell-nzl/fid/internal/config"
	"github.com/mbell-nzl/fid/internal/gittest"
)

// retiredPricingFileName is the model price list fid once wrote into every
// new ledger for V1's cost telemetry. Nothing read it.
const retiredPricingFileName = "pricing.toml"

// However a ledger is made, it has no pricing.toml.
func TestNewLedgersHaveNoPricingFile(t *testing.T) {
	t.Parallel()
	for _, test := range []struct {
		name string
		// ledger makes a ledger and returns its project root.
		ledger func(t *testing.T) string
	}{
		{"init", func(t *testing.T) string {
			workspace := gittest.NewWorkspace(t)
			runInitCommand(t, workspace, "init", "--json")
			return workspace
		}},
		{"export import", func(t *testing.T) string {
			archive := filepath.Join(t.TempDir(), "export.zip")
			if created := executeCommandForTest(clonePristineWorkspace(t), "--json", "export", "create", archive); created.err != nil {
				t.Fatalf("export create = %#v", created)
			}
			target := gittest.NewWorkspace(t)
			imported := executeCommandForTest(target, "--json", "export", "import", archive, "--repository", target, "--idempotency-key", "no-pricing")
			if imported.err != nil || !strings.Contains(imported.stdout, `"ok":true`) {
				t.Fatalf("export import = %#v", imported)
			}
			return target
		}},
		{"restore", func(t *testing.T) string {
			fixture := newBackupFixture(t)
			original, id := fixture.project()
			fixture.setTarget(fixture.target)
			decodeBackupSuccess(t, fixture.backup(original, "--json"))
			workspace := fixture.restoreWorkspace(id)
			decodeRestoreSuccess(t, fixture.restore(workspace, "--json"))
			return workspace
		}},
	} {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			ledger := config.LedgerDir(test.ledger(t))
			if _, err := os.Stat(filepath.Join(ledger, config.ConfigFileName)); err != nil {
				t.Fatalf("the new ledger has no config: %v", err)
			}
			if _, err := os.Lstat(filepath.Join(ledger, retiredPricingFileName)); !errors.Is(err, fs.ErrNotExist) {
				t.Fatalf("the new ledger has a %s: %v", retiredPricingFileName, err)
			}
		})
	}
}

// A ledger made while fid wrote pricing.toml keeps the file, which fid
// neither reads nor removes, and works as before.
func TestALedgerWithAPricingFileStillWorks(t *testing.T) {
	t.Parallel()
	workspace := clonePristineWorkspace(t)
	path := filepath.Join(config.LedgerDir(workspace), retiredPricingFileName)
	contents := []byte("[models.claude-opus-4-7]\nin_per_mtok = 15\nout_per_mtok = 75\n")
	if err := os.WriteFile(path, contents, 0o644); err != nil {
		t.Fatal(err)
	}
	runFIDCommand(t, workspace, "issue", "create", "--kind", "task", "--title", "Made beside a pricing.toml", "--actor", "agent:developer")
	for _, args := range [][]string{{"issue", "list"}, {"status"}, {"export", "create", filepath.Join(t.TempDir(), "export.zip")}} {
		if result := executeCommandForTest(workspace, append([]string{"--json"}, args...)...); result.err != nil || result.stderr != "" {
			t.Fatalf("fid %v = %#v", args, result)
		}
	}
	if checked := executeCommandForTest(workspace, "--json", "maintenance", "check"); checked.err != nil || !strings.Contains(checked.stdout, `"status":"pass"`) {
		t.Fatalf("maintenance check = %#v, want a pass", checked)
	}
	if kept, err := os.ReadFile(path); err != nil || string(kept) != string(contents) {
		t.Fatalf("%s after the commands = %q, %v; want it kept unchanged", retiredPricingFileName, kept, err)
	}
}
