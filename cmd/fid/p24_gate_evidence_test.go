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
	fiddb "github.com/mbell-nzl/fid/internal/db"
	"github.com/mbell-nzl/fid/internal/journal"
)

func TestP24GatePreSpawnObservationAfterEdit(t *testing.T) {
	t.Parallel()
	root, issueID, envelopeID := runEnvelopeForTestWithDefinition(t, func(definition map[string]any) {
		definition["argv"] = []string{"/bin/touch", "spawned.txt"}
	})
	path, key := gateRunInputForTest(t, root, issueID, envelopeID, "p24-observed-gate")
	if err := os.WriteFile(filepath.Join(root, "tracked.txt"), []byte("original"), 0o600); err != nil {
		t.Fatal(err)
	}
	runGitForTest(t, root, "add", "tracked.txt")
	runGitForTest(t, root, "commit", "-m", "advance after Run")
	if err := os.WriteFile(filepath.Join(root, "tracked.txt"), []byte("edited"), 0o600); err != nil {
		t.Fatal(err)
	}
	observed := p24GateDirtyRepository(t, root, "tracked.txt")
	git := func(args ...string) string {
		t.Helper()
		data, err := exec.Command("git", append([]string{"-C", root}, args...)...).Output()
		if err != nil {
			t.Fatal(err)
		}
		return strings.TrimSpace(string(data))
	}
	absolute, err := filepath.Abs(root)
	if err != nil {
		t.Fatal(err)
	}
	if observed["head_oid"] != git("rev-parse", "HEAD") || observed["base_oid"] != git("rev-parse", "HEAD") || observed["tree_oid"] != git("rev-parse", "HEAD^{tree}") || observed["dirty"] != true || git("status", "--porcelain", "--", "tracked.txt") == "" || observed["worktree_id"] != fmt.Sprintf("sha256:%x", sha256.Sum256([]byte(absolute))) {
		t.Fatalf("independent observation disagrees: %#v", observed)
	}
	input := decodeCommandJSONStructure(t, string(mustReadCommandFile(t, path))).(map[string]any)
	delete(input, "repository")
	input["expected_repository"] = map[string]any{"head_oid": observed["head_oid"], "dirty": observed["dirty"]}
	writeActivationInput(t, path, input)
	response := executeCommandForTest(root, "--project", root, "--json", "gate", "run", "--actor", "actor-builder", "--input", path, "--idempotency-key", key)
	if response.err != nil {
		t.Fatalf("gate after edit: %#v", response)
	}
	lines := strings.Split(strings.TrimSpace(response.stdout), "\n")
	var terminal gateTerminalFrame
	if err := json.Unmarshal([]byte(lines[len(lines)-1]), &terminal); err != nil {
		t.Fatal(err)
	}
	receipt := terminal.Data.Response["result"].(map[string]any)["receipt"].(map[string]any)
	if !reflect.DeepEqual(receipt["terminal_repository"], observed) {
		t.Fatalf("Receipt repository = %#v, want pre-spawn %#v", receipt["terminal_repository"], observed)
	}
	p24AssertGatePersistedRepositories(t, root, terminal, observed)
	count := 0
	eventCounts := map[string]int{}
	attemptEffects, receiptEffects := 0, 0
	for _, line := range strings.Split(strings.TrimSpace(string(mustReadCommandFile(t, mustJournalPath(t, root)))), "\n") {
		var event journal.Event
		if err := json.Unmarshal([]byte(line), &event); err != nil {
			t.Fatal(err)
		}
		if event.EventType != "gate.execution-started" && event.EventType != "gate.receipt-recorded" {
			continue
		}
		count++
		eventCounts[event.EventType]++
		data, err := json.Marshal(event.Repository)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(decodeCommandJSONStructure(t, string(data)), observed) {
			t.Fatalf("%s repository = %s", event.EventType, data)
		}
		for _, effect := range event.Effects {
			if effect.Table == "execution_attempts" {
				attemptEffects++
				if !reflect.DeepEqual(decodeCommandJSONStructure(t, effect.Values["start_repository_json"].(string)), observed) {
					t.Fatal("Attempt differs from pre-spawn observation")
				}
			}
			if effect.Table == "execution_receipts" {
				receiptEffects++
				if !reflect.DeepEqual(decodeCommandJSONStructure(t, effect.Values["terminal_repository_json"].(string)), observed) {
					t.Fatal("Receipt Effect differs from pre-spawn observation")
				}
			}
		}
	}
	if count != 2 {
		t.Fatalf("gate Events = %d", count)
	}
	if eventCounts["gate.execution-started"] != 1 || eventCounts["gate.receipt-recorded"] != 1 || attemptEffects != 1 || receiptEffects != 1 {
		t.Fatalf("gate Event/Effect counts = %v / %d Attempts / %d Receipts", eventCounts, attemptEffects, receiptEffects)
	}
	if _, err := os.Stat(filepath.Join(root, "spawned.txt")); err != nil {
		t.Fatal(err)
	}
}

// This fixture has one unstaged tracked edit. Build its expected repository
// from Git, config and file images, without fid's observer or repository output.
func p24GateDirtyRepository(t *testing.T, workspace, path string) map[string]any {
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
	projectConfig, err := config.ReadProjectConfig(filepath.Join(root, ".fid", "config.toml"))
	if err != nil {
		t.Fatal(err)
	}
	head := git("rev-parse", "HEAD")
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
	image := func(mode string, data []byte) map[string]any {
		return map[string]any{"mode": mode, "digest": fmt.Sprintf("sha256:%x", sha256.Sum256(data))}
	}
	// Stdlib JSON sorts these map keys to independently construct the accepted
	// fid.repository-dirty.v1 image; no production observer/canonicalizer is used.
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

func p24AssertGatePersistedRepositories(t *testing.T, root string, terminal gateTerminalFrame, repository map[string]any) {
	t.Helper()
	conn, err := fiddb.Open(filepath.Join(root, ".fid", config.DatabaseFileName))
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = conn.Close() }()
	for _, record := range []struct {
		name, query, id string
	}{
		{"Attempt", `SELECT COUNT(*), MIN(id), MIN(start_repository_json) FROM execution_attempts`, terminal.Data.AttemptID},
		{"Receipt", `SELECT COUNT(*), MIN(id), MIN(terminal_repository_json) FROM execution_receipts`, terminal.Data.ReceiptID},
	} {
		var count int
		var id, data string
		if err := conn.QueryRow(record.query).Scan(&count, &id, &data); err != nil {
			t.Fatal(err)
		}
		if count != 1 || id == "" || id != record.id {
			t.Fatalf("persisted %s count/id = %d/%q, want 1/%q", record.name, count, id, record.id)
		}
		if got := decodeCommandJSONStructure(t, data); !reflect.DeepEqual(got, repository) {
			t.Fatalf("persisted %s repository = %#v, want pre-spawn %#v", record.name, got, repository)
		}
	}
}

func TestP24GateExpectationBeforeAttempt(t *testing.T) {
	t.Parallel()
	for _, field := range []string{"head_oid", "dirty"} {
		t.Run(field, func(t *testing.T) {
			t.Parallel()
			root, issueID, envelopeID := runEnvelopeForTestWithDefinition(t, func(definition map[string]any) {
				definition["argv"] = []string{"/bin/touch", "spawned.txt"}
			})
			path, key := gateRunInputForTest(t, root, issueID, envelopeID, "p24-mismatch")
			observed := p24IndependentBindingRepository(t, root)
			input := decodeCommandJSONStructure(t, string(mustReadCommandFile(t, path))).(map[string]any)
			delete(input, "repository")
			expected := map[string]any{"head_oid": observed["head_oid"], "dirty": observed["dirty"]}
			if field == "head_oid" {
				expected[field] = strings.Repeat("f", 40)
			} else {
				expected[field] = !observed[field].(bool)
			}
			input["expected_repository"] = expected
			writeActivationInput(t, path, input)
			before := snapshotWorkspacePersistence(t, root)
			response := executeCommandForTest(root, "--project", root, "--json", "gate", "run", "--actor", "actor-builder", "--input", path, "--idempotency-key", key)
			if cli.ExitCode(response.err) != 4 || response.stdout != "" {
				t.Fatalf("mismatch: %#v", response)
			}
			failure := decodeCommandJSONStructure(t, response.stderr).(map[string]any)["error"].(map[string]any)
			if failure["code"] != "repository-mismatch" || !reflect.DeepEqual(failure["details"], map[string]any{"expected": expected, "observed": observed}) {
				t.Fatalf("mismatch details: %#v", failure)
			}
			assertRunStateUnchanged(t, root, before)
			if _, err := os.Stat(filepath.Join(root, "spawned.txt")); !os.IsNotExist(err) {
				t.Fatalf("process spawned: %v", err)
			}
		})
	}
}

func TestP24ReceiptEvidenceUsesSubjectAndCurrentObservation(t *testing.T) {
	t.Parallel()
	root, issueID, envelopeID := gateReadyEvidenceEnvelopeForTest(t, func(definition map[string]any) {
		definition["argv"] = []string{"/bin/true"}
	})
	path, key := gateRunInputForTest(t, root, issueID, envelopeID, "p24-evidence-gate")
	input := decodeCommandJSONStructure(t, string(mustReadCommandFile(t, path))).(map[string]any)
	delete(input, "repository")
	observed := p24IndependentBindingRepository(t, root)
	input["expected_repository"] = map[string]any{"head_oid": observed["head_oid"], "dirty": observed["dirty"]}
	writeActivationInput(t, path, input)
	response := executeCommandForTest(root, "--project", root, "--json", "gate", "run", "--actor", "actor-builder", "--input", path, "--idempotency-key", key)
	if response.err != nil {
		t.Fatalf("gate: %#v", response)
	}
	lines := strings.Split(strings.TrimSpace(response.stdout), "\n")
	var terminal gateTerminalFrame
	if err := json.Unmarshal([]byte(lines[len(lines)-1]), &terminal); err != nil {
		t.Fatal(err)
	}
	// Linking must copy the subject even when the current tree has moved.
	if err := os.WriteFile(filepath.Join(root, "later.txt"), []byte("later"), 0o600); err != nil {
		t.Fatal(err)
	}
	writeActivationInput(t, path, map[string]any{"obligation_id": input["obligation_id"], "evidence_type": "receipt", "receipt_id": terminal.Data.ReceiptID, "review_id": nil, "human_outcome_id": nil, "required_artifact_ids": []string{}})
	savedGit := filepath.Join(t.TempDir(), "saved-git")
	if err := os.Rename(filepath.Join(root, ".git"), savedGit); err != nil {
		t.Fatal(err)
	}
	linked := executeCommandForTest(root, "--project", root, "--json", "evidence", "link", "--actor", "actor-builder", "--input", path, "--idempotency-key", "p24-derived-link")
	if linked.err != nil {
		t.Fatalf("link: %#v", linked)
	}
	if got := decodeCanonicalResult(t, linked.stdout)["evidence"].(map[string]any)["target_repository"]; !reflect.DeepEqual(got, observed) {
		t.Fatalf("target = %#v, want subject %#v", got, observed)
	}
	p24AssertBindingEventRepository(t, root, observed)
	if err := os.Rename(savedGit, filepath.Join(root, ".git")); err != nil {
		t.Fatal(err)
	}
	for _, state := range []string{"edited", "reverted"} {
		if state == "reverted" {
			if err := os.Remove(filepath.Join(root, "later.txt")); err != nil {
				t.Fatal(err)
			}
		}
		writeActivationInput(t, path, map[string]any{"obligation_id": input["obligation_id"]})
		before := snapshotWorkspacePersistence(t, root)
		evaluated := executeCommandForTest(root, "--project", root, "--json", "evidence", "evaluate", "--input", path)
		if evaluated.err != nil {
			t.Fatalf("evaluate %s: %#v", state, evaluated)
		}
		evaluation := decodeCanonicalResult(t, evaluated.stdout)["evaluation"].(map[string]any)
		want := "fresh"
		if state == "edited" {
			want = "stale"
		}
		if evaluation["freshness"] != want {
			t.Fatalf("%s freshness = %#v", state, evaluation)
		}
		assertRunStateUnchanged(t, root, before)
	}
	// evidence.get remains a stored-record read; only evaluate observes Git.
	if err := os.Rename(filepath.Join(root, ".git"), filepath.Join(t.TempDir(), "saved-git")); err != nil {
		t.Fatal(err)
	}
	evidenceID := decodeCanonicalResult(t, linked.stdout)["evidence"].(map[string]any)["id"].(string)
	before := snapshotWorkspacePersistence(t, root)
	got := executeCommandForTest(root, "--project", root, "--json", "evidence", "get", evidenceID)
	if got.err != nil || got.stderr != "" {
		t.Fatalf("stored Evidence read without Git: %#v", got)
	}
	if evaluation := decodeCanonicalResult(t, got.stdout)["evaluation"].(map[string]any); evaluation["freshness"] != "fresh" {
		t.Fatalf("stored Evidence evaluation: %#v", evaluation)
	}
	assertRunStateUnchanged(t, root, before)
}

func TestP24GateEvidenceDescriptorProcessCorpus(t *testing.T) {
	t.Parallel()
	leaves := map[string]cli.InputDescriptor{}
	for _, leaf := range []string{"gate run", "evidence link", "evidence evaluate", "artifact verify"} {
		leaves[leaf] = inputDescriptorRegistry[leaf]
	}
	p24DescriptorProcessCorpus(t, leaves)
}

func TestP24GateOriginalCopyAndImport(t *testing.T) {
	t.Parallel()
	root, issueID, envelopeID := runEnvelopeForTestWithDefinition(t, func(definition map[string]any) {
		definition["argv"] = []string{"/bin/true"}
	})
	path, key := gateRunInputForTest(t, root, issueID, envelopeID, "p24-binding-gate")
	input := decodeCommandJSONStructure(t, string(mustReadCommandFile(t, path))).(map[string]any)
	delete(input, "repository")
	copied := t.TempDir()
	if err := copyWorkspaceTree(root, copied); err != nil {
		t.Fatal(err)
	}
	archive := filepath.Join(t.TempDir(), "history.zip")
	if response := executeCommandForTest(root, "export", "create", archive, "--json"); response.err != nil {
		t.Fatalf("export: %#v", response)
	}
	target := filepath.Join(t.TempDir(), "imported")
	if data, err := exec.Command("git", "clone", "-q", root, target).CombinedOutput(); err != nil {
		t.Fatalf("clone: %v %s", err, data)
	}
	if response := executeCommandForTest(target, "export", "import", archive, "--repository", target, "--idempotency-key", "import", "--json"); response.err != nil {
		t.Fatalf("import: %#v", response)
	}
	for _, workspace := range []string{copied, target, root} {
		observed := p24IndependentBindingRepository(t, workspace)
		input["expected_repository"] = map[string]any{"head_oid": observed["head_oid"], "dirty": observed["dirty"]}
		path = filepath.Join(workspace, "binding-gate.json")
		writeActivationInput(t, path, input)
		before := snapshotWorkspacePersistence(t, workspace)
		response := executeCommandForTest(workspace, "--project", workspace, "--json", "gate", "run", "--actor", "actor-builder", "--input", path, "--idempotency-key", key)
		if workspace == copied {
			if cli.ExitCode(response.err) != 2 || response.stdout != "" {
				t.Fatalf("plain copy: %#v", response)
			}
			failure := decodeCommandJSONStructure(t, response.stderr).(map[string]any)["error"].(map[string]any)
			if failure["code"] != "invalid-argument" || failure["message"] != "run gate: repository binding is invalid" || !reflect.DeepEqual(failure["details"], map[string]any{"argument": "input", "value": nil, "allowed": []any{"accepted gate registration input"}}) {
				t.Fatalf("plain-copy refusal: %#v", failure)
			}
			assertRunStateUnchanged(t, workspace, before)
		} else if response.err != nil {
			t.Fatalf("admitted original/import: %#v", response)
		}
	}
}

func TestP24RemovedGateEvidenceMembers(t *testing.T) {
	t.Parallel()
	for leaf, member := range map[string]string{"gate run": "repository", "evidence link": "target_repository", "evidence evaluate": "target_repository"} {
		t.Run(leaf, func(t *testing.T) {
			t.Parallel()
			root := t.TempDir()
			path := filepath.Join(t.TempDir(), "input.json")
			writeActivationInput(t, path, map[string]any{member: nil})
			args := append([]string{"--project", root, "--json"}, strings.Split(leaf, " ")...)
			args = append(args, "--input", path)
			if leaf != "evidence evaluate" {
				args = append(args, "--actor", "ACT-1", "--idempotency-key", "removed")
			}
			response := executeCommandForTest(root, args...)
			if cli.ExitCode(response.err) != 2 || response.stdout != "" {
				t.Fatalf("removed member: %#v", response)
			}
			failure := decodeCommandJSONStructure(t, response.stderr).(map[string]any)["error"].(map[string]any)
			if failure["code"] != "invalid-argument" || failure["details"].(map[string]any)["argument"] != "input."+member {
				t.Fatalf("removed-member path: %#v", failure)
			}
			entries, err := os.ReadDir(root)
			if err != nil || len(entries) != 0 {
				t.Fatalf("structural rejection changed directory: %v %v", entries, err)
			}
		})
	}
}
