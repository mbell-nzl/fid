//go:build integration

package integration

import (
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/mbell-nzl/fid/internal/config"
	"github.com/mbell-nzl/fid/internal/gittest"
)

// Each binary must keep its own migration history after its build checkout
// disappears, including release builds whose source locations are trimmed.
func TestBinaryMigrationsSurviveRemovedSource(t *testing.T) {
	t.Parallel()
	for _, trimpath := range []bool{false, true} {
		name := "ordinary"
		if trimpath {
			name = "trimpath"
		}
		t.Run(name, func(t *testing.T) {
			binary := buildWithoutSource(t, trimpath)
			journey := backupJourney{base: t.TempDir()}
			env := journey.environment(filepath.Join(journey.base, "config"))
			workspace := gittest.NewWorkspace(t)
			call := func(command string, args ...string) string {
				t.Helper()
				cmd := exec.Command(command, args...)
				cmd.Dir, cmd.Env = workspace, env
				output, err := cmd.CombinedOutput()
				if err != nil {
					t.Fatalf("%s %v: %v\n%s", filepath.Base(command), args, err, output)
				}
				return string(output)
			}
			call("git", "-c", "user.name=fid test", "-c", "user.email=fid@example.test", "-c", "commit.gpgsign=false",
				"commit", "-q", "--no-verify", "--allow-empty", "-m", "base")
			initialized := p23Envelope(t, call(binary, "init", "--prefix", "EMB", "--json"))
			projectID := p21PathString(t, initialized, "result", "project", "id")
			created := p23Envelope(t, call(binary, "issue", "create", "--kind", "task", "--title", "Self-contained migrations",
				"--actor", "human:mark", "--idempotency-key", "create", "--json"))
			issueID := p21PathString(t, created, "result", "issue", "id")
			call(binary, "rebuild", "--json")
			got := p23Envelope(t, call(binary, "issue", "get", issueID, "--json"))
			if title := p21PathString(t, got, "result", "issue", "title"); title != "Self-contained migrations" {
				t.Fatalf("rebuilt title = %q", title)
			}
			archive := filepath.Join(journey.base, "export.zip")
			call(binary, "export", "create", archive, "--json")
			validated := p23Envelope(t, call(binary, "export", "validate", archive, "--json"))
			if status := p21PathString(t, validated, "result", "validation", "status"); status != "valid" {
				t.Fatalf("archive validation = %q, want valid", status)
			}
			target := filepath.Join(journey.base, "backups.git")
			call("git", "init", "--bare", "-q", "--template=", target)
			if err := config.SetBackupTarget(filepath.Join(journey.base, "config", "fid", "config.toml"), target); err != nil {
				t.Fatal(err)
			}
			call(binary, "backup", "--json")
			if outcomes := backupJourneyOutcomes(t, workspace); len(outcomes) != 1 || outcomes[0] != "pushed" {
				t.Fatalf("backup outcomes = %v, want [pushed]", outcomes)
			}
			call("git", "--git-dir="+target, "cat-file", "-e", "refs/heads/main:"+projectID+"/manifest.json")
		})
	}
}

func buildWithoutSource(t *testing.T, trimpath bool) string {
	t.Helper()
	source := filepath.Join(t.TempDir(), "source")
	if err := os.Mkdir(source, 0o700); err != nil {
		t.Fatal(err)
	}
	root := integrationRepositoryRoot()
	for _, dir := range []string{"cmd", "internal", "migrations"} {
		if err := os.CopyFS(filepath.Join(source, dir), os.DirFS(filepath.Join(root, dir))); err != nil {
			t.Fatal(err)
		}
	}
	for _, name := range []string{"go.mod", "go.sum"} {
		contents, err := fs.ReadFile(os.DirFS(root), name)
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(source, name), contents, 0o600); err != nil {
			t.Fatal(err)
		}
	}
	binary := filepath.Join(t.TempDir(), filepath.Base(integrationFIDBinary(t)))
	args := []string{"build", "-buildvcs=false", "-o", binary}
	if trimpath {
		args = append(args, "-trimpath")
	}
	cmd := exec.Command("go", append(args, "./cmd/fid")...)
	cmd.Dir = source
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("build isolated fid: %v\n%s", err, output)
	}
	if err := os.RemoveAll(source); err != nil {
		t.Fatal(err)
	}
	return binary
}
