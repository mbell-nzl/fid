package main

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"testing"

	"github.com/mbell-nzl/fid/internal/config"
)

// The folder contains every export member verbatim, including kept Artifacts.
func TestBackupStoresTheExportMembersUnzipped(t *testing.T) {
	t.Parallel()
	fixture := newBackupFixture(t)
	workspace, _ := v8Ledger(t)
	id, err := config.ReadProjectIdentity(workspace)
	if err != nil {
		t.Fatal(err)
	}
	fixture.setTarget(fixture.target)
	archive, _ := exportForTest(t, workspace)
	members := exportMembersForBackup(t, archive)
	decodeBackupSuccess(t, fixture.backup(workspace, "--json"))
	want := []string{id + "/backup.json"}
	artifacts := 0
	for name, content := range members {
		want = append(want, id+"/"+name)
		stored, err := fixture.gitBytes("--git-dir="+fixture.target, "cat-file", "blob", "main:"+id+"/"+name)
		if err != nil || !bytes.Equal(stored, content) {
			t.Errorf("stored %s differs from the export: %v", name, err)
		}
		if strings.HasPrefix(name, "journal/") && len(content) > 64<<20 {
			t.Errorf("journal member %s has %d bytes, exceeds 64 MiB", name, len(content))
		}
		if strings.HasPrefix(name, "artifacts/") {
			artifacts++
		}
	}
	if artifacts == 0 {
		t.Fatal("fixture has no kept Artifact")
	}
	slices.Sort(want)
	if got := fixture.targetGit("ls-tree", "-r", "--name-only", "main"); got != strings.Join(want, "\n") {
		t.Fatalf("backup files = %s, want %v", got, want)
	}
}

func exportMembersForBackup(t *testing.T, archive []byte) map[string][]byte {
	t.Helper()
	reader, err := zip.NewReader(bytes.NewReader(archive), int64(len(archive)))
	if err != nil {
		t.Fatal(err)
	}
	members := map[string][]byte{}
	for _, file := range reader.File {
		input, err := file.Open()
		if err != nil {
			t.Fatal(err)
		}
		content, err := io.ReadAll(input)
		_ = input.Close()
		if err != nil {
			t.Fatal(err)
		}
		members[file.Name] = content
	}
	return members
}

// A hundred changed backups must not leave a loose full copy per backup,
// even when the user's git config disables automatic maintenance.
func TestBackupPacksItsCacheAfterOneHundredBackups(t *testing.T) {
	t.Parallel()
	fixture := newBackupFixture(t)
	workspace, _ := fixture.project()
	fixture.setTarget(fixture.target)
	writeBackupTestFile(t, filepath.Join(fixture.home, ".gitconfig"), "[gc]\n\tauto = 0\n\tautoDetach = true\n[maintenance]\n\tauto = false\n\tautoDetach = true\n")
	for index := range 100 {
		runFIDCommand(t, workspace, "issue", "create", "--kind", "task", "--title", fmt.Sprintf("Backup %d", index), "--actor", "agent:developer")
		if record, _ := decodeBackupSuccess(t, fixture.backup(workspace, "--json")); record.Outcome != "pushed" {
			t.Fatalf("backup %d = %+v, want pushed", index+1, record)
		}
	}
	if count := fixture.targetGit("rev-list", "--count", "main"); count != "100" {
		t.Fatalf("backup commits = %s, want 100", count)
	}
	counts := map[string]int{}
	for _, line := range strings.Split(fixture.git("--git-dir="+fixture.clone(), "count-objects", "-v"), "\n") {
		name, value, _ := strings.Cut(line, ": ")
		counts[name], _ = strconv.Atoi(value)
	}
	if counts["packs"] < 1 || counts["in-pack"] < 300 || counts["count"] >= 256 {
		t.Fatalf("cache object counts = %v, want history packed and fewer than 256 loose objects", counts)
	}
	if got := fixture.git("--git-dir="+fixture.clone(), "fsck", "--no-dangling"); got != "" {
		t.Fatalf("cache fsck: %s", got)
	}
	if leftovers, err := filepath.Glob(filepath.Join(filepath.Dir(fixture.clone()), "export")); err != nil || len(leftovers) != 0 {
		t.Fatalf("backup left scratch files: %v, %v", leftovers, err)
	}
	if _, err := os.Stat(filepath.Join(fixture.clone(), "objects", "pack")); err != nil {
		t.Fatal(err)
	}
}

func TestBackupRecordsPackingFailureBeforePushing(t *testing.T) {
	t.Parallel()
	fixture := newBackupFixture(t)
	workspace, _ := fixture.project()
	fixture.setTarget(fixture.target)
	first, _ := decodeBackupSuccess(t, fixture.backup(workspace, "--json"))
	// Seed enough distinct loose objects to require packing on the next
	// backup. No refs are needed: they model objects from failed pushes.
	input := filepath.Join(fixture.base, "loose-input")
	for index := range 300 {
		writeBackupTestFile(t, input, fmt.Sprintf("unreferenced object %d\n", index))
		fixture.git("--git-dir="+fixture.clone(), "hash-object", "-w", "--no-filters", input)
	}
	fixture.git("--git-dir="+fixture.clone(), "config", "pack.indexVersion", "77")
	runFIDCommand(t, workspace, "issue", "create", "--kind", "task", "--title", "Not pushed", "--actor", "agent:developer")

	failure := decodeBackupFailure(t, fixture.backup(workspace, "--json"), 7)
	if failure.Code != "internal-error" || !strings.Contains(failure.Message, "pack the backup clone") {
		t.Fatalf("backup failure = %+v, want a packing internal-error", failure)
	}
	if fixture.targetMain() != *first.Commit {
		t.Fatal("failed packing changed the target")
	}
	records := backupRecords(t, workspace)
	if len(records) != 2 || records[1].Outcome != "failed" || records[1].Failure == nil || records[1].Failure.Kind != "internal" || records[1].Commit != nil {
		t.Fatalf("records = %+v, want the failed packing recorded without a commit", records)
	}
	fixture.git("--git-dir="+fixture.clone(), "config", "--unset", "pack.indexVersion")
	if retry, _ := decodeBackupSuccess(t, fixture.backup(workspace, "--json")); retry.Outcome != "pushed" {
		t.Fatalf("retry after fixing packing = %+v, want pushed", retry)
	}
}
