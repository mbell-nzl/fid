package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
	"unicode"

	"github.com/mbell-nzl/fid/internal/config"
	"github.com/mbell-nzl/fid/internal/gittest"
)

// statusControls is text fid didn't write that a terminal would act on: it
// erases the screen and starts a line of its own.
const statusControls = "\x1b[2J\nINJECTED"

// statusEscaped is statusControls as status text shows it.
const statusEscaped = `\u001b[2J\nINJECTED`

// statusBackupJSON is the backup object of fid --json status.
func statusBackupJSON(t *testing.T, stdout string) map[string]any {
	t.Helper()
	var envelope struct {
		Result struct {
			Backup map[string]any `json:"backup"`
		} `json:"result"`
	}
	if err := json.Unmarshal([]byte(stdout), &envelope); err != nil || envelope.Result.Backup == nil {
		t.Fatalf("status --json = %v:\n%s", err, stdout)
	}
	return envelope.Result.Backup
}

// Paths and backup record values carry their controls into JSON, and are
// escaped in the text, which then holds no control but its line breaks.
func TestStatusEscapesWhatFidDidNotWrite(t *testing.T) {
	t.Parallel()
	for _, test := range []struct {
		name string
		// setup makes the fixture and returns the raw value status shows.
		setup func(t *testing.T) (*statusFixture, string)
		text  func(raw string) []string
		json  func(backup map[string]any) string
	}{
		{"user config path", func(t *testing.T) (*statusFixture, string) {
			fixture := newStatusFixture(t, 2*time.Hour)
			fixture.configHome = filepath.Join(t.TempDir(), "config"+statusControls)
			return fixture, filepath.Join(fixture.configHome, "fid", "config.toml")
		}, func(raw string) []string {
			return []string{"  ⚠ No backup target is set; not backed up yet\n", "      → add [backup] target to\n        " + escapedPath(raw) + ",\n        then fid backup\n"}
		}, func(backup map[string]any) string { return backup["user_config_path"].(string) }},
		{"backup record path", func(t *testing.T) (*statusFixture, string) {
			workspace := gittest.NewNamedWorkspace(t, "work"+statusControls)
			runInitCommand(t, workspace, "init", "--json")
			base := t.TempDir()
			fixture := &statusFixture{t: t, workspace: workspace, configHome: filepath.Join(base, "config"), home: filepath.Join(base, "home")}
			fixture.setTarget()
			record := filepath.Join(config.LedgerDir(gittest.ResolvedPath(t, workspace)), "backups.jsonl")
			if err := os.Mkdir(record, 0o755); err != nil {
				t.Fatal(err)
			}
			return fixture, record
		}, func(raw string) []string {
			return []string{"  ⚠ Can't read the backup record\n    " + escapedPath(raw) + "\n", "      → fix or remove\n        " + escapedPath(raw) + ",\n        then fid backup\n"}
		}, func(backup map[string]any) string { return backup["record_path"].(string) }},
		{"failure kind and message", func(t *testing.T) (*statusFixture, string) {
			fixture := newStatusFixture(t, 100*time.Hour)
			fixture.setTarget()
			kind := "push" + statusControls
			fixture.writeRecord(statusRecordLine{"pushed", 2 * time.Hour, ""}, statusRecordLine{"failed", time.Hour, kind})
			path := filepath.Join(config.LedgerDir(fixture.workspace), "backups.jsonl")
			payload, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			message, err := json.Marshal("fatal: MESSAGE" + statusControls)
			if err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(path, []byte(strings.Replace(string(payload), `"the backup failed"`, string(message), 1)), 0o644); err != nil {
				t.Fatal(err)
			}
			return fixture, kind
		}, func(raw string) []string {
			return []string{"  ⚠ Last backup failed (push" + statusEscaped + "); last backed up 2\n    hours ago\n"}
		}, func(backup map[string]any) string {
			failure := backup["latest"].(map[string]any)["failure"].(map[string]any)
			if message := failure["message"].(string); message != "fatal: MESSAGE"+statusControls {
				return "message " + message
			}
			return failure["kind"].(string)
		}},
	} {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			fixture, raw := test.setup(t)
			text := fixture.run("status")
			if text.err != nil || text.stderr != "" {
				t.Fatalf("status = %v, stderr %q", text.err, text.stderr)
			}
			for _, want := range test.text(raw) {
				if !strings.Contains(text.stdout, want) {
					t.Fatalf("status text doesn't show %q:\n%s", want, text.stdout)
				}
			}
			for _, value := range text.stdout {
				if value != '\n' && unicode.IsControl(value) {
					t.Fatalf("status text holds control %U:\n%q", value, text.stdout)
				}
			}
			if strings.Contains(text.stdout, "\nINJECTED") || strings.Contains(text.stdout, "MESSAGE") {
				t.Fatalf("status text shows an injected line or the failure's message:\n%q", text.stdout)
			}
			result := fixture.run("--json", "status")
			if result.err != nil {
				t.Fatalf("status --json = %v, stderr %q", result.err, result.stderr)
			}
			if got := test.json(statusBackupJSON(t, result.stdout)); got != raw {
				t.Fatalf("status JSON = %q, want %q as it is", got, raw)
			}
		})
	}
}

// escapedPath is path as status text shows it, its controls escaped.
func escapedPath(path string) string {
	return strings.ReplaceAll(path, statusControls, statusEscaped)
}
