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
	"time"

	"github.com/mbell-nzl/fid/internal/cli"
	"github.com/mbell-nzl/fid/internal/config"
	"github.com/mbell-nzl/fid/internal/workflow"
)

func TestP24RunStartAfterCommitAndEdit(t *testing.T) {
	t.Parallel()
	workspace, issueID, envelopeID := runEnvelopeForTest(t)
	claim, _ := acquireGenericClaimForTest(t, workspace, issueID, envelopeID, "p24-start-claim")
	if err := os.WriteFile(filepath.Join(workspace, "tracked.txt"), []byte("committed"), 0o600); err != nil {
		t.Fatal(err)
	}
	runGitForTest(t, workspace, "add", "tracked.txt")
	runGitForTest(t, workspace, "commit", "-m", "advance after Envelope")
	if err := os.WriteFile(filepath.Join(workspace, "tracked.txt"), []byte("edited"), 0o600); err != nil {
		t.Fatal(err)
	}
	repository := p24RunDirtyRepository(t, workspace, "tracked.txt")
	input := runStartInput(issueID, envelopeID, claim["id"], repository)
	delete(input, "repository")
	input["expected_repository"] = map[string]any{"head_oid": repository["head_oid"], "dirty": repository["dirty"]}
	path := filepath.Join(workspace, "p24-start.json")
	writeActivationInput(t, path, input)
	result := executeCommandForTest(workspace, "--project", workspace, "--json", "run", "start", "--actor", "actor-builder", "--input", path, "--idempotency-key", "p24-start")
	if result.err != nil {
		t.Fatalf("start after commit/edit: %#v", result)
	}
	if got := runFromResponse(t, result.stdout)["start_repository"]; !reflect.DeepEqual(got, repository) {
		t.Fatalf("repository = %#v, want observation %#v", got, repository)
	}
	git := func(args ...string) string {
		t.Helper()
		data, err := exec.Command("git", append([]string{"-C", workspace}, args...)...).Output()
		if err != nil {
			t.Fatal(err)
		}
		return strings.TrimSpace(string(data))
	}
	root, err := filepath.Abs(workspace)
	if err != nil {
		t.Fatal(err)
	}
	if repository["head_oid"] != git("rev-parse", "HEAD") || repository["base_oid"] != git("rev-parse", "HEAD") || repository["tree_oid"] != git("rev-parse", "HEAD^{tree}") || repository["dirty"] != true || git("status", "--porcelain", "--", "tracked.txt") == "" || repository["worktree_id"] != fmt.Sprintf("sha256:%x", sha256.Sum256([]byte(root))) {
		t.Fatalf("independent Git observation disagrees: %#v", repository)
	}
	p24AssertBindingEventRepository(t, workspace, repository)
}

func runGitForTest(t *testing.T, workspace string, args ...string) {
	t.Helper()
	command := exec.Command("git", args...)
	command.Dir = workspace
	if output, err := command.CombinedOutput(); err != nil {
		t.Fatalf("git %v: %v: %s", args, err, output)
	}
}

func TestP24RunCloseAfterCommitAndEdit(t *testing.T) {
	t.Parallel()
	workspace, issueID, envelopeID := runEnvelopeForTest(t)
	claim, _ := acquireGenericClaimForTest(t, workspace, issueID, envelopeID, "p24-close-claim")
	path := filepath.Join(workspace, "p24-close.json")
	writeActivationInput(t, path, runStartInput(issueID, envelopeID, claim["id"], runRepositoryForTest(t, workspace)))
	started := executeCommandForTest(workspace, "--project", workspace, "--json", "run", "start", "--actor", "actor-builder", "--input", path, "--idempotency-key", "p24-close-start")
	if started.err != nil {
		t.Fatalf("start: %#v", started)
	}
	run := runFromResponse(t, started.stdout)
	writeActivationInput(t, path, map[string]any{"run_id": run["id"], "expected_version": 1, "liveness": map[string]any{"executor_id": "executor-builder", "state": "terminated", "source_ref": "process:p24", "observed_at": time.Now().UTC().Add(-time.Second).Format(time.RFC3339Nano)}})
	observed := executeCommandForTest(workspace, "--project", workspace, "--json", "run", "observe", "--actor", "actor-harness", "--input", path, "--idempotency-key", "p24-close-observe")
	if observed.err != nil {
		t.Fatalf("observe: %#v", observed)
	}
	observation := decodeCommandJSONStructure(t, observed.stdout).(map[string]any)["result"].(map[string]any)["observation"].(map[string]any)
	if err := os.WriteFile(filepath.Join(workspace, "tracked.txt"), []byte("committed"), 0o600); err != nil {
		t.Fatal(err)
	}
	runGitForTest(t, workspace, "add", "tracked.txt")
	runGitForTest(t, workspace, "commit", "-m", "advance after Run")
	if err := os.WriteFile(filepath.Join(workspace, "tracked.txt"), []byte("edited"), 0o600); err != nil {
		t.Fatal(err)
	}
	repository := p24RunDirtyRepository(t, workspace, "tracked.txt")
	writeActivationInput(t, path, map[string]any{"run_id": run["id"], "expected_version": 2, "observation_id": observation["id"], "expected_repository": map[string]any{"head_oid": repository["head_oid"], "dirty": repository["dirty"]}})
	closed := executeCommandForTest(workspace, "--project", workspace, "--json", "run", "close", "--actor", "actor-builder", "--input", path, "--idempotency-key", "p24-close")
	if closed.err != nil {
		t.Fatalf("close: %#v", closed)
	}
	if got := runFromResponse(t, closed.stdout)["terminal_repository"]; !reflect.DeepEqual(got, repository) {
		t.Fatalf("terminal = %#v, want %#v", got, repository)
	}
	p24AssertBindingEventRepository(t, workspace, repository)
}

// The fixture has exactly one unstaged tracked edit. Require that complete
// status before constructing the expected repository without fid observation.
func p24RunDirtyRepository(t *testing.T, workspace, path string) map[string]any {
	t.Helper()
	gitBytes := func(args ...string) []byte {
		t.Helper()
		data, err := exec.Command("git", append([]string{"-C", workspace}, args...)...).Output()
		if err != nil {
			t.Fatal(err)
		}
		return data
	}
	git := func(args ...string) string {
		t.Helper()
		return strings.TrimSuffix(string(gitBytes(args...)), "\n")
	}
	if status := git("status", "--porcelain", "--untracked-files=all", "--", ".", ":!.fid"); status != " M "+path {
		t.Fatalf("dirty fixture status = %q, want one tracked edit %s", status, path)
	}
	root, err := filepath.Abs(git("rev-parse", "--show-toplevel"))
	if err != nil {
		t.Fatal(err)
	}
	projectConfig, err := config.ReadProjectConfig(filepath.Join(config.LedgerDir(workspace), "config.toml"))
	if err != nil {
		t.Fatal(err)
	}
	head := git("rev-parse", "HEAD")
	image := func(mode string, data []byte) map[string]any {
		return map[string]any{"mode": mode, "digest": fmt.Sprintf("sha256:%x", sha256.Sum256(data))}
	}
	headEntry := strings.Fields(git("ls-tree", "HEAD", "--", path))
	indexEntry := strings.Fields(git("ls-files", "-s", "--", path))
	if len(headEntry) != 4 || headEntry[1] != "blob" || len(indexEntry) != 4 || indexEntry[2] != "0" {
		t.Fatalf("unsupported tracked fixture: HEAD %v / index %v", headEntry, indexEntry)
	}
	info, err := os.Lstat(filepath.Join(root, path))
	if err != nil || !info.Mode().IsRegular() {
		t.Fatalf("worktree fixture must be regular: %v", err)
	}
	content, err := os.ReadFile(filepath.Join(root, path))
	if err != nil {
		t.Fatal(err)
	}
	mode := "100644"
	if info.Mode().Perm()&0o111 != 0 {
		mode = "100755"
	}
	// encoding/json sorts map keys, independently producing the accepted
	// fid.repository-dirty.v1 canonical image (no production JSON helper).
	entries := []any{map[string]any{"path": path,
		"head":     image(headEntry[0], gitBytes("cat-file", "blob", headEntry[2])),
		"index":    image(indexEntry[0], gitBytes("cat-file", "blob", indexEntry[1])),
		"worktree": image(mode, content)}}
	payload, err := json.Marshal(map[string]any{"schema": "fid.repository-dirty.v1", "entries": entries})
	if err != nil {
		t.Fatal(err)
	}
	return map[string]any{"base_oid": head, "capture_source": "git", "dirty": true, "dirty_digest": fmt.Sprintf("sha256:%x", sha256.Sum256(payload)),
		"head_oid": head, "project_id": projectConfig.ProjectID, "tree_oid": git("rev-parse", "HEAD^{tree}"),
		"worktree_id": fmt.Sprintf("sha256:%x", sha256.Sum256([]byte(root)))}
}

func TestP24RunStartOriginalCopyAndImport(t *testing.T) {
	t.Parallel()
	workspace, issueID, envelopeID := runEnvelopeForTest(t)
	claim, _ := acquireGenericClaimForTest(t, workspace, issueID, envelopeID, "p24-copy-claim")
	copied := t.TempDir()
	if err := copyWorkspaceTree(workspace, copied); err != nil {
		t.Fatal(err)
	}
	archive := filepath.Join(t.TempDir(), "history.zip")
	exported := executeCommandForTest(workspace, "export", "create", archive, "--json")
	if exported.err != nil {
		t.Fatalf("export: %#v", exported)
	}
	target := filepath.Join(t.TempDir(), "imported")
	if data, err := exec.Command("git", "clone", "-q", workspace, target).CombinedOutput(); err != nil {
		t.Fatalf("clone: %v %s", err, data)
	}
	imported := executeCommandForTest(target, "export", "import", archive, "--repository", target, "--idempotency-key", "import", "--json")
	if imported.err != nil {
		t.Fatalf("import: %#v", imported)
	}
	for _, root := range []string{copied, target, workspace} {
		repository := p24IndependentBindingRepository(t, root)
		input := runStartInput(issueID, envelopeID, claim["id"], repository)
		before := snapshotWorkspacePersistence(t, root)
		path := filepath.Join(root, "p24-binding-start.json")
		writeActivationInput(t, path, input)
		response := executeCommandForTest(root, "--project", root, "--json", "run", "start", "--actor", "actor-builder", "--input", path, "--idempotency-key", "p24-binding-start")
		if root == copied {
			assertRunStartInvalidArgument(t, response, "run claim binding is invalid")
			if cli.ExitCode(response.err) != 2 {
				t.Fatalf("copy exit: %v", response.err)
			}
			var request workflow.RunStartRequest
			data, err := json.Marshal(input)
			if err != nil {
				t.Fatal(err)
			}
			if err := json.Unmarshal(data, &request); err != nil {
				t.Fatal(err)
			}
			_, err = workflow.StartRun(workflow.RunOptions{WorkspaceRoot: root, MigrationsDir: issueMigrationsDir(), ActorID: "actor-builder", IdempotencyKey: "p24-domain-cause"}, request)
			if err == nil || err.Error() != "run claim binding is invalid" {
				t.Fatalf("copy domain cause = %v", err)
			}
			assertRunStateUnchanged(t, root, before)
		} else {
			if response.err != nil {
				t.Fatalf("%s: %#v", root, response)
			}
			if !reflect.DeepEqual(runFromResponse(t, response.stdout)["start_repository"], repository) {
				t.Fatal("Run did not record independent observation")
			}
		}
	}
}

func TestP24RunStartExpectationMatrix(t *testing.T) {
	t.Parallel()
	for _, kind := range []string{"match", "HEAD", "dirty"} {
		t.Run(kind, func(t *testing.T) {
			t.Parallel()
			root, issueID, envelopeID := runEnvelopeForTest(t)
			claim, _ := acquireGenericClaimForTest(t, root, issueID, envelopeID, "expectation-claim")
			repository := p24IndependentBindingRepository(t, root)
			input := runStartInput(issueID, envelopeID, claim["id"], repository)
			expected := input["expected_repository"].(map[string]any)
			if kind == "HEAD" {
				expected["head_oid"] = strings.Repeat("f", 40)
			}
			if kind == "dirty" {
				expected["dirty"] = !repository["dirty"].(bool)
			}
			before := snapshotWorkspacePersistence(t, root)
			path := filepath.Join(root, "expectation.json")
			writeActivationInput(t, path, input)
			response := executeCommandForTest(root, "--project", root, "--json", "run", "start", "--actor", "actor-builder", "--input", path, "--idempotency-key", "expectation")
			if kind == "match" {
				if response.err != nil {
					t.Fatalf("match: %#v", response)
				}
				return
			}
			if cli.ExitCode(response.err) != 4 {
				t.Fatalf("mismatch: %#v", response)
			}
			failure := decodeCommandJSONStructure(t, response.stderr).(map[string]any)["error"].(map[string]any)
			if failure["code"] != "repository-mismatch" || !reflect.DeepEqual(failure["details"], map[string]any{"expected": expected, "observed": repository}) {
				t.Fatalf("mismatch details: %#v", failure)
			}
			assertRunStateUnchanged(t, root, before)
		})
	}
}
