package main

import (
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/mbell-nzl/fid/internal/bootstrap"
	"github.com/mbell-nzl/fid/internal/config"
	"github.com/mbell-nzl/fid/internal/journal"
)

var pristineWorkspaceRoot string

var issueMutationTestSequence uint64

func TestMain(m *testing.M) {
	exitCode := runCommandTestSuite(m)
	os.Exit(exitCode)
}

func runCommandTestSuite(m interface{ Run() int }) int {
	templateRoot, err := os.MkdirTemp("", "fid-pristine-workspace-")
	if err != nil {
		fmt.Fprintf(os.Stderr, "create pristine workspace root: %v\n", err)
		return 1
	}
	pristineWorkspaceRoot = templateRoot
	_, err = bootstrap.InitializeWorkspace(bootstrap.InitOptions{
		WorkspaceRoot: templateRoot,
		MigrationsDir: filepath.Join(commandRepositoryRoot(), "migrations"),
		Prefix:        "FID",
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "initialize pristine workspace: %v\n", err)
		_ = os.RemoveAll(templateRoot)
		return 1
	}

	exitCode := m.Run()
	if err := os.RemoveAll(templateRoot); err != nil {
		fmt.Fprintf(os.Stderr, "remove pristine workspace root: %v\n", err)
		if exitCode == 0 {
			exitCode = 1
		}
	}
	return exitCode
}

func commandRepositoryRoot() string {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return ""
	}
	return filepath.Clean(filepath.Join(filepath.Dir(filename), "..", ".."))
}

func clonePristineWorkspace(t *testing.T) string {
	t.Helper()
	destination := t.TempDir()
	clonePristineWorkspaceInto(t, destination)
	return destination
}

func clonePristineWorkspaceInto(t *testing.T, destination string) {
	t.Helper()
	if pristineWorkspaceRoot == "" {
		t.Fatal("pristine workspace was not initialized")
	}
	if err := copyWorkspaceTree(pristineWorkspaceRoot, destination); err != nil {
		t.Fatalf("clone pristine workspace: %v", err)
	}
}

func copyWorkspaceTree(sourceRoot, destinationRoot string) error {
	return filepath.WalkDir(sourceRoot, func(sourcePath string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		relative, err := filepath.Rel(sourceRoot, sourcePath)
		if err != nil {
			return err
		}
		if relative == "." {
			return nil
		}
		destinationPath := filepath.Join(destinationRoot, relative)
		info, err := entry.Info()
		if err != nil {
			return err
		}
		switch {
		case entry.IsDir():
			return os.Mkdir(destinationPath, info.Mode().Perm())
		case info.Mode().IsRegular():
			return copyWorkspaceFile(sourcePath, destinationPath, info.Mode().Perm())
		case info.Mode()&os.ModeSymlink != 0:
			target, err := os.Readlink(sourcePath)
			if err != nil {
				return err
			}
			return os.Symlink(target, destinationPath)
		default:
			return fmt.Errorf("unsupported workspace entry %q with mode %v", sourcePath, info.Mode())
		}
	})
}

func copyWorkspaceFile(source, destination string, mode fs.FileMode) error {
	input, err := os.Open(source)
	if err != nil {
		return err
	}
	defer func() { _ = input.Close() }()
	output, err := os.OpenFile(destination, os.O_CREATE|os.O_EXCL|os.O_WRONLY, mode)
	if err != nil {
		return err
	}
	if _, err := io.Copy(output, input); err != nil {
		_ = output.Close()
		return err
	}
	return output.Close()
}

type workspaceSnapshotEntry struct {
	Mode fs.FileMode
	Data []byte
}

type commandExecutionResult struct {
	stdout string
	stderr string
	err    error
}

func executeCommandForTest(workingDirectory string, args ...string) commandExecutionResult {
	args = withIssueMutationTestKey(args)
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	cmd := newRootCommandAt(workingDirectory)
	cmd.SetOut(stdout)
	cmd.SetErr(stderr)
	cmd.SetArgs(args)
	err := cmd.Execute()
	return commandExecutionResult{stdout: stdout.String(), stderr: stderr.String(), err: err}
}

func withIssueMutationTestKey(args []string) []string {
	args = migrateTask9RetainedMutationFixture(args)
	if (!isTypedIssueMutation(args) && !isTask9RetainedMutation(args)) || containsTestArgument(args, "--idempotency-key") {
		return args
	}
	return append(args, "--idempotency-key", fmt.Sprintf("mutation-test-%d", atomic.AddUint64(&issueMutationTestSequence, 1)))
}

func isTypedIssueMutation(args []string) bool {
	for index := 0; index+1 < len(args); index++ {
		if args[index] == "issue" && (args[index+1] == "create" || args[index+1] == "update" || args[index+1] == "transition") {
			return containsTestArgument(args, "--actor")
		}
	}
	return false
}

func migrateTask9RetainedMutationFixture(args []string) []string {
	if !isTask9RetainedMutation(args) || containsTestArgument(args, "--actor") {
		return args
	}
	for index := 0; index+1 < len(args); index++ {
		if args[index] != "--by" {
			continue
		}
		migrated := append([]string{}, args...)
		migrated[index] = "--actor"
		migrated[index+1] = "agent:" + args[index+1]
		return migrated
	}
	return args
}

func isTask9RetainedMutation(args []string) bool {
	for index := range args {
		if index+2 < len(args) && args[index] == "issue" && args[index+1] == "ac" && args[index+2] == "add" {
			return true
		}
		if index+2 < len(args) && args[index] == "issue" && args[index+1] == "label" && (args[index+2] == "add" || args[index+2] == "remove") {
			return true
		}
		if index+1 < len(args) && args[index] == "label" && args[index+1] == "add" {
			return true
		}
		if index+1 < len(args) && args[index] == "link" && (args[index+1] == "add" || args[index+1] == "remove") {
			return true
		}
		if index+1 < len(args) && args[index] == "requirement" && args[index+1] == "sync" {
			return true
		}
		if index+1 < len(args) && args[index] == "workflow" && (args[index+1] == "define" || args[index+1] == "assign") {
			return true
		}
	}
	return false
}

func containsTestArgument(args []string, target string) bool {
	for _, argument := range args {
		if argument == target {
			return true
		}
	}
	return false
}

func TestTask9RetainedMutationFixtureArgumentsUseActorAndIdempotency(t *testing.T) {
	args := withIssueMutationTestKey([]string{"label", "add", "core", "--by", "developer"})
	if containsTestArgument(args, "--by") {
		t.Fatalf("retained fixture args = %v, want no V1 --by", args)
	}
	if !containsTestArgument(args, "--actor") || !containsTestArgument(args, "--idempotency-key") {
		t.Fatalf("retained fixture args = %v, want standalone actor and idempotency key", args)
	}
	for index := range args {
		if args[index] == "--actor" && (index+1 == len(args) || args[index+1] != "agent:developer") {
			t.Fatalf("retained fixture actor = %v, want agent:developer", args)
		}
	}
}

func TestTask9FixtureAdapterLeavesJudgmentAndPolicyCommandsUntouched(t *testing.T) {
	for _, args := range [][]string{
		{"concern", "raise", "--by", "developer"},
		{"policy", "show"},
	} {
		got := withIssueMutationTestKey(args)
		if !reflect.DeepEqual(got, args) {
			t.Fatalf("fixture adapter(%v) = %v, want unchanged", args, got)
		}
	}
}

func TestPristineWorkspaceClonePreservesRebuildParityAndIsolation(t *testing.T) {
	t.Parallel()
	cloneA := clonePristineWorkspace(t)
	cloneB := clonePristineWorkspace(t)
	want := snapshotWorkspace(t, pristineWorkspaceRoot)
	if got := snapshotWorkspace(t, cloneA); !reflect.DeepEqual(got, want) {
		t.Fatal("first clone differs from pristine template")
	}
	if got := snapshotWorkspace(t, cloneB); !reflect.DeepEqual(got, want) {
		t.Fatal("second clone differs from pristine template")
	}

	runFIDCommand(t, cloneA, "issue", "create", "--kind", "task", "--title", "Fixture isolation", "--actor", "agent:developer")
	if got := snapshotWorkspace(t, cloneA); reflect.DeepEqual(got, want) {
		t.Fatal("mutated clone still equals pristine template")
	}
	if got := snapshotWorkspace(t, pristineWorkspaceRoot); !reflect.DeepEqual(got, want) {
		t.Fatal("mutation changed pristine template")
	}
	if got := snapshotWorkspace(t, cloneB); !reflect.DeepEqual(got, want) {
		t.Fatal("mutation crossed into sibling clone")
	}

	databasePath := config.WorkspacePathsForRoot(cloneA).DatabasePath
	beforeRebuild := readWorkspaceFixtureFile(t, databasePath)
	if _, err := journal.RebuildWorkspace(cloneA, filepath.Join("..", "..", "migrations")); err != nil {
		t.Fatalf("rebuild cloned workspace: %v", err)
	}
	afterRebuild := readWorkspaceFixtureFile(t, databasePath)
	if !bytes.Equal(afterRebuild, beforeRebuild) {
		t.Fatal("cloned workspace projection bytes differ after rebuild")
	}
}

func TestConcurrentClonedWorkspaceCommandsRemainIsolated(t *testing.T) {
	t.Parallel()
	templateBefore := snapshotWorkspace(t, pristineWorkspaceRoot)
	type fixture struct {
		workspace string
		title     string
	}
	fixtures := []fixture{
		{clonePristineWorkspace(t), "Concurrent fixture alpha"},
		{clonePristineWorkspace(t), "Concurrent fixture beta"},
	}
	start := make(chan struct{})
	results := make(chan commandExecutionResult, len(fixtures))
	for _, item := range fixtures {
		item := item
		go func() {
			<-start
			results <- executeCommandForTest(item.workspace,
				"issue", "create", "--kind", "task", "--title", item.title, "--actor", "agent:developer")
		}()
	}
	close(start)
	for range fixtures {
		result := <-results
		if result.err != nil {
			t.Fatalf("concurrent command error = %v, stderr = %q", result.err, result.stderr)
		}
	}

	for index, item := range fixtures {
		journal := string(readWorkspaceFixtureFile(t, mustJournalPath(t, item.workspace)))
		if !strings.Contains(journal, item.title) {
			t.Fatalf("workspace %d journal does not contain own title %q", index, item.title)
		}
		other := fixtures[(index+1)%len(fixtures)].title
		if strings.Contains(journal, other) {
			t.Fatalf("workspace %d journal contains sibling title %q", index, other)
		}
	}
	if got := snapshotWorkspace(t, pristineWorkspaceRoot); !reflect.DeepEqual(got, templateBefore) {
		t.Fatal("concurrent commands changed pristine template")
	}
}

func snapshotWorkspace(t *testing.T, root string) map[string]workspaceSnapshotEntry {
	t.Helper()
	snapshot := make(map[string]workspaceSnapshotEntry)
	err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		relative, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		if relative == "." {
			return nil
		}
		info, err := entry.Info()
		if err != nil {
			return err
		}
		item := workspaceSnapshotEntry{Mode: info.Mode()}
		if info.Mode().IsRegular() {
			item.Data, err = os.ReadFile(path)
			if err != nil {
				return err
			}
		}
		snapshot[relative] = item
		return nil
	})
	if err != nil {
		t.Fatalf("snapshot workspace %q: %v", root, err)
	}
	return snapshot
}

func readWorkspaceFixtureFile(t *testing.T, path string) []byte {
	t.Helper()
	contents, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read fixture file %q: %v", path, err)
	}
	return contents
}
