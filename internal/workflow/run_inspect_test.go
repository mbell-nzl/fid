package workflow

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRunInspectIsStableAndCheckpointBecomesComparisonBase(t *testing.T) {
	options, request := checkpointFixture(t)
	first, err := InspectRun(options.WorkspaceRoot, options.MigrationsDir, request.RunID, options.WorkspaceRoot)
	if err != nil {
		t.Fatalf("InspectRun(first) error = %v", err)
	}
	second, err := InspectRun(options.WorkspaceRoot, options.MigrationsDir, request.RunID, options.WorkspaceRoot)
	if err != nil {
		t.Fatalf("InspectRun(second) error = %v", err)
	}
	if first.Reconciliation.Digest != second.Reconciliation.Digest || second.Reconciliation.CheckpointID != nil {
		t.Fatalf("unchanged inspection = %#v then %#v, want stable digest without checkpoint", first, second)
	}
	checkpoint, err := RunCheckpoint(options, request)
	if err != nil {
		t.Fatalf("RunCheckpoint() error = %v", err)
	}
	inspected, err := InspectRun(options.WorkspaceRoot, options.MigrationsDir, request.RunID, options.WorkspaceRoot)
	if err != nil {
		t.Fatalf("InspectRun(after checkpoint) error = %v", err)
	}
	if inspected.Reconciliation.CheckpointID == nil || *inspected.Reconciliation.CheckpointID != checkpoint.Checkpoint.ID {
		t.Fatalf("InspectRun() checkpoint = %#v, want %q comparison base", inspected.Reconciliation.CheckpointID, checkpoint.Checkpoint.ID)
	}
}

func TestRunInspectLeavesRepositoryByteInert(t *testing.T) {
	options, request := checkpointFixture(t)
	before := inspectRepositorySnapshot(t, options.WorkspaceRoot)
	if _, err := InspectRun(options.WorkspaceRoot, options.MigrationsDir, request.RunID, options.WorkspaceRoot); err != nil {
		t.Fatalf("InspectRun() error = %v", err)
	}
	after := inspectRepositorySnapshot(t, options.WorkspaceRoot)
	if before != after {
		t.Fatalf("InspectRun() changed refs, index, stash, or worktree: before=%q after=%q", before, after)
	}
}

func inspectRepositorySnapshot(t *testing.T, root string) string {
	t.Helper()
	base, err := os.ReadFile(filepath.Join(root, "base.txt"))
	if err != nil {
		t.Fatalf("read worktree: %v", err)
	}
	index, err := os.ReadFile(filepath.Join(root, ".git", "index"))
	if err != nil {
		t.Fatalf("read index: %v", err)
	}
	head, err := os.ReadFile(filepath.Join(root, ".git", "HEAD"))
	if err != nil {
		t.Fatalf("read HEAD: %v", err)
	}
	stash, err := os.ReadFile(filepath.Join(root, ".git", "refs", "stash"))
	if err != nil && !os.IsNotExist(err) {
		t.Fatalf("read stash: %v", err)
	}
	return string(base) + "\x00" + string(index) + "\x00" + string(head) + "\x00" + string(stash)
}
