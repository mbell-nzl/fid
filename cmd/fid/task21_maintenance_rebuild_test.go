package main

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/mbell-nzl/fid/internal/config"
)

func TestFinalSchemaNamedReadsAndExportMatchAfterRebuild(t *testing.T) {
	t.Parallel()
	workspace, issueID := task21CompletedWorkspace(t)
	archivePath := filepath.Join(t.TempDir(), "task21-export.zip")
	before := task21SnapshotReads(t, workspace, issueID, archivePath)

	if err := os.Remove(filepath.Join(config.LedgerDir(workspace), "fid.sqlite")); err != nil {
		t.Fatalf("delete projection: %v", err)
	}
	rebuilt := executeCommandForTest(workspace, "--project", workspace, "--json", "rebuild")
	if rebuilt.err != nil || rebuilt.stderr != "" {
		t.Fatalf("rebuild final-schema projection = %#v", rebuilt)
	}
	after := task21SnapshotReads(t, workspace, issueID, archivePath)

	if !bytes.Equal(before.exportBytes, after.exportBytes) {
		t.Fatal("export bytes differ after final-schema rebuild")
	}
	if before.issue != after.issue || before.show != after.show || before.checks != after.checks || before.next != after.next || before.maintenance != after.maintenance || before.validation != after.validation {
		t.Fatalf("named reads differ after rebuild: before=%#v after=%#v", before, after)
	}
}

type task21ReadSnapshot struct {
	issue       string
	show        string
	checks      string
	next        string
	maintenance string
	validation  string
	exportBytes []byte
}

func task21SnapshotReads(t *testing.T, workspace, issueID, archivePath string) task21ReadSnapshot {
	t.Helper()
	read := func(args ...string) string {
		if args[0] == "sql" {
			return task21Read(t, workspace, args...)
		}
		return task21Read(t, workspace, append([]string{"--project", workspace, "--json"}, args...)...)
	}
	issue := read("issue", "get", issueID)
	shown := read("show", issueID)
	checks := ""
	for _, table := range []string{"holds", "questions", "person_decisions", "policy_snapshots", "check_requests", "check_starts", "check_processes", "check_results", "check_files", "referenced_artifacts", "handoffs", "idempotency_results"} {
		checks += read("sql", "SELECT * FROM "+table+" ORDER BY rowid")
	}

	next := read("next")
	maintenance := read("maintenance", "check")
	_ = read("export", "create", archivePath)
	validation := read("export", "validate", archivePath)
	exportBytes, err := os.ReadFile(archivePath)
	if err != nil {
		t.Fatal(err)
	}
	return task21ReadSnapshot{issue, shown, checks, next, maintenance, validation, exportBytes}
}

func task21Read(t *testing.T, workspace string, args ...string) string {
	t.Helper()
	result := executeCommandForTest(workspace, args...)
	if result.err != nil || result.stderr != "" {
		t.Fatalf("read %v = %#v", args, result)
	}
	return result.stdout
}

func task21CompletedWorkspace(t *testing.T) (string, string) {
	t.Helper()
	workspace := clonePristineWorkspace(t)
	appendPolicy(t, workspace, "\n[checks.parity]\nrun=['true']\ntimeout='1m'\nartifacts=['parity.txt']\n[lanes.standard]\nchecks=['parity']\n[review]\ntest_adequacy=[]\n")
	if err := os.WriteFile(filepath.Join(workspace, "parity.txt"), []byte("retained parity Artifact\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	runGitForTest(t, workspace, "add", ".")
	runGitForTest(t, workspace, "commit", "-qm", "Parity fixture")
	runGitForTest(t, workspace, "remote", "add", "origin", workspace)
	runGitForTest(t, workspace, "update-ref", "refs/remotes/origin/trunk", "HEAD")
	runGitForTest(t, workspace, "symbolic-ref", "refs/remotes/origin/HEAD", "refs/remotes/origin/trunk")
	created := decodeCanonicalResult(t, task21Read(t, workspace, "--json", "issue", "create", "--kind", "task", "--title", "Retained parity", "--lane", "standard", "--actor", "human:mark", "--idempotency-key", "parity-create"))
	task := created["issue"].(map[string]any)["id"].(string)
	taken := decodeCanonicalResult(t, task21Read(t, workspace, "--json", "takeover", task, "--actor", "human:mark", "--reason", "Build parity fixture", "--idempotency-key", "parity-takeover"))
	hold := taken["hold_id"].(string)
	task21Mutate(t, workspace, "start", "--hold", hold, "--actor", "human:mark", "--idempotency-key", "parity-start")
	task21Mutate(t, workspace, "check", "--hold", hold, "--actor", "human:mark", "--idempotency-key", "parity-check")
	task21Mutate(t, workspace, "handoff", "--hold", hold, "--actor", "human:mark", "--pr", "https://github.com/mbell-nzl/fid/pull/63", "--idempotency-key", "parity-handoff")
	assertV8HasNoKernelRecords(t, workspace)
	return workspace, task
}

func task21Mutate(t *testing.T, workspace string, args ...string) {
	t.Helper()
	result := executeCommandForTest(workspace, append([]string{"--project", workspace, "--json"}, args...)...)
	if result.err != nil || result.stderr != "" {
		t.Fatalf("mutation %v = %#v", args, result)
	}
}
