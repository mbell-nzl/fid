package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/mbell-nzl/fid/internal/cli"
	"github.com/mbell-nzl/fid/internal/config"
)

// The ledger is no longer in the worktree, so a .fid/ there is an ordinary
// path (D4 / F1, decided 6 November 2026), and the dirty help no longer
// singles it out.
func TestDirtyHelpCountsEveryPath(t *testing.T) {
	t.Parallel()
	const want = "input.expected_repository.dirty: boolean, required — True when git status --porcelain lists any path."
	for _, leaf := range []string{"run start", "run close", "run abandon", "review record"} {
		t.Run(leaf, func(t *testing.T) {
			t.Parallel()
			help := executeCommandForTest(t.TempDir(), append(strings.Fields(leaf), "--help")...)
			if help.err != nil || !strings.Contains(help.stdout, want) || strings.Contains(help.stdout, ".fid") {
				t.Fatalf("%s help (error %v):\n%s\nwant %q and no mention of .fid", leaf, help.err, help.stdout, want)
			}
		})
	}
}

// A leftover .fid/ makes the worktree dirty, as git status reports it.
func TestRunStartCountsALeftoverFidDirectoryAsDirty(t *testing.T) {
	t.Parallel()
	root, issueID, envelopeID := runEnvelopeForTest(t)
	claim, _ := acquireGenericClaimForTest(t, root, issueID, envelopeID, "leftover-fid-claim")
	const leftover = ".fid/config.toml"
	content := []byte("leftover V1 config\n")
	name := filepath.Join(root, filepath.FromSlash(leftover))
	if err := os.MkdirAll(filepath.Dir(name), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(name, content, 0o644); err != nil {
		t.Fatal(err)
	}
	observed := leftoverFidRepository(t, root, leftover, content)

	input := runStartInput(issueID, envelopeID, claim["id"], observed)
	expected := input["expected_repository"].(map[string]any)
	expected["dirty"] = false
	before := snapshotWorkspacePersistence(t, root)
	path := filepath.Join(root, "leftover-fid-start.json")
	writeActivationInput(t, path, input)
	refused := executeCommandForTest(root, "--project", root, "--json", "run", "start", "--actor", "actor-builder", "--input", path, "--idempotency-key", "leftover-fid-clean")
	if cli.ExitCode(refused.err) != 4 {
		t.Fatalf("start expecting a clean worktree: %#v", refused)
	}
	failure := decodeCommandJSONStructure(t, refused.stderr).(map[string]any)["error"].(map[string]any)
	if failure["code"] != "repository-mismatch" || !reflect.DeepEqual(failure["details"], map[string]any{"expected": expected, "observed": observed}) {
		t.Fatalf("mismatch = %#v, want expected %#v and observed %#v", failure, expected, observed)
	}
	assertRunStateUnchanged(t, root, before)

	expected["dirty"] = true
	writeActivationInput(t, path, input)
	started := executeCommandForTest(root, "--project", root, "--json", "run", "start", "--actor", "actor-builder", "--input", path, "--idempotency-key", "leftover-fid-dirty")
	if started.err != nil {
		t.Fatalf("start expecting a dirty worktree: %#v", started)
	}
	if got := runFromResponse(t, started.stdout)["start_repository"]; !reflect.DeepEqual(got, observed) {
		t.Fatalf("start repository = %#v, want %#v", got, observed)
	}
}

// leftoverFidRepository builds fid's expected observation of a worktree whose
// only change is one untracked file, from git and the file's known bytes.
func leftoverFidRepository(t *testing.T, workspace, path string, content []byte) map[string]any {
	t.Helper()
	git := func(args ...string) string {
		t.Helper()
		data, err := exec.Command("git", append([]string{"-C", workspace}, args...)...).Output()
		if err != nil {
			t.Fatalf("git %v: %v", args, err)
		}
		return strings.TrimSuffix(string(data), "\n")
	}
	if status := git("status", "--porcelain", "--untracked-files=all"); status != "?? "+path {
		t.Fatalf("git status = %q, want one untracked %s", status, path)
	}
	projectConfig, err := config.ReadProjectConfig(filepath.Join(config.LedgerDir(workspace), config.ConfigFileName))
	if err != nil {
		t.Fatal(err)
	}
	entries := []any{map[string]any{"path": path, "head": nil, "index": nil, "worktree": map[string]any{"mode": "100644", "digest": fmt.Sprintf("sha256:%x", sha256.Sum256(content))}}}
	payload, err := json.Marshal(map[string]any{"schema": "fid.repository-dirty.v1", "entries": entries})
	if err != nil {
		t.Fatal(err)
	}
	head := git("rev-parse", "HEAD")
	return map[string]any{"base_oid": head, "capture_source": "git", "dirty": true, "dirty_digest": fmt.Sprintf("sha256:%x", sha256.Sum256(payload)),
		"head_oid": head, "project_id": projectConfig.ProjectID, "tree_oid": git("rev-parse", "HEAD^{tree}"),
		"worktree_id": fmt.Sprintf("sha256:%x", sha256.Sum256([]byte(git("rev-parse", "--show-toplevel"))))}
}
