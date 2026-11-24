package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mbell-nzl/fid/internal/checkread"
	"github.com/mbell-nzl/fid/internal/cli"
	"github.com/mbell-nzl/fid/internal/config"
	"github.com/mbell-nzl/fid/internal/journal"
	"github.com/mbell-nzl/fid/internal/reviewwork"
	fidmigrations "github.com/mbell-nzl/fid/migrations"
)

func keptTaskCheckPath(t *testing.T, f *statusFixture, task string) string {
	t.Helper()
	conn, err := reviewwork.OpenProjection(f.workspace)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = conn.Close() }()
	records, err := checkread.Read(conn, task)
	if err != nil {
		t.Fatal(err)
	}
	return filepath.Join(config.LedgerDir(f.workspace), "artifacts", strings.TrimPrefix(records[0].Result.Artifacts[0].Digest, "sha256:"))
}

func TestHandoffRequiresKeptArtifactsNowButReplayDoesNot(t *testing.T) {
	t.Parallel()
	for _, change := range []string{"missing", "changed"} {
		t.Run(change, func(t *testing.T) {
			t.Parallel()
			f, task, hold := doneBuildFixture(t, "standard", false)
			path := keptTaskCheckPath(t, f, task)
			original, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			if change == "missing" {
				if err := os.Remove(path); err != nil {
					t.Fatal(err)
				}
			} else {
				writeBackupTestFile(t, path, "changed kept bytes")
			}
			before := startEventCount(t, f.workspace)
			for _, jsonMode := range []bool{false, true} {
				args := doneHandoffArgs(hold)
				if jsonMode {
					args = append([]string{"--json"}, args...)
				}
				got := f.run(args...)
				if cli.ExitCode(got.err) != 4 || got.stdout != "" || !strings.Contains(got.stderr, "kept-artifact-unavailable") {
					t.Fatalf("%s refusal=%#v", change, got)
				}
				if jsonMode {
					assertJSONFailureCode(t, got.stderr, "task.handoff", "kept-artifact-unavailable")
				} else if !strings.Contains(got.stderr, "kept Artifacts are missing or changed") {
					t.Fatalf("human refusal=%s", got.stderr)
				}
			}
			if startEventCount(t, f.workspace) != before {
				t.Fatal("refusal wrote an event")
			}
			writeBackupTestFile(t, path, string(original))
			args := append([]string{"--json"}, doneHandoffArgs(hold)...)
			first := f.run(args...)
			lifecycleResult(t, first, "task.handoff")
			if err := os.Remove(path); err != nil {
				t.Fatal(err)
			}
			if _, err := journal.RebuildWorkspace(f.workspace, fidmigrations.FS); err != nil {
				t.Fatal(err)
			}
			if got := f.run(args...); got.err != nil || got.stdout != first.stdout {
				t.Fatalf("replay re-read kept bytes: %#v", got)
			}
		})
	}
}
