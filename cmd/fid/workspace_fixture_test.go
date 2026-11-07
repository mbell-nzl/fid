package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/mbell-nzl/fid/internal/bootstrap"
	"github.com/mbell-nzl/fid/internal/config"
	"github.com/mbell-nzl/fid/internal/gittest"
	"github.com/mbell-nzl/fid/internal/journal"
)

var pristineWorkspaceRoot string

var issueMutationTestSequence uint64

var reusableWorkspaceFixtures = struct {
	sync.Mutex
	items map[string]*reusableWorkspaceFixture
}{items: make(map[string]*reusableWorkspaceFixture)}

var envelopeBindingsForTest sync.Map

type reusableWorkspaceFixture struct {
	sync.Mutex
	root   string
	values []string
}

type envelopeBindingForTest struct {
	registrationID      string
	registrationVersion int
	obligationID        string
	claimID             string
	runID               string
}

func rememberEnvelopeBindingForTest(t *testing.T, workspace string, binding envelopeBindingForTest) {
	t.Helper()
	envelopeBindingsForTest.Store(workspace, binding)
	t.Cleanup(func() { envelopeBindingsForTest.Delete(workspace) })
}

func envelopeBindingForWorkspace(t *testing.T, workspace string) envelopeBindingForTest {
	t.Helper()
	value, ok := envelopeBindingsForTest.Load(workspace)
	if !ok {
		t.Fatalf("workspace %q has no cached Envelope binding", workspace)
	}
	return value.(envelopeBindingForTest)
}

func TestMain(m *testing.M) {
	exitCode := runCommandTestSuite(m)
	os.Exit(exitCode)
}

func runCommandTestSuite(m interface{ Run() int }) int {
	if err := gittest.Guard(); err != nil {
		fmt.Fprintf(os.Stderr, "guard git discovery: %v\n", err)
		return 1
	}
	if containsTestArgument(os.Args, "--live-golden-child") || commandTestHelperChild() {
		// Children inherit the isolated config directory from their parent.
		return m.Run()
	}
	configHome, err := os.MkdirTemp("", "fid-test-config-")
	if err != nil {
		fmt.Fprintf(os.Stderr, "create isolated user config: %v\n", err)
		return 1
	}
	defer func() { _ = os.RemoveAll(configHome) }()
	if err := os.Setenv("XDG_CONFIG_HOME", configHome); err != nil {
		fmt.Fprintf(os.Stderr, "isolate user config: %v\n", err)
		return 1
	}

	templateRoot, err := os.MkdirTemp("", "fid-pristine-workspace-")
	if err != nil {
		fmt.Fprintf(os.Stderr, "create pristine workspace root: %v\n", err)
		return 1
	}
	pristineWorkspaceRoot = templateRoot
	if err := gittest.InitRepository(templateRoot); err != nil {
		fmt.Fprintf(os.Stderr, "make pristine workspace a git repository: %v\n", err)
		_ = os.RemoveAll(templateRoot)
		return 1
	}
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
	if err := cleanupReusableWorkspaceFixtures(); err != nil {
		fmt.Fprintf(os.Stderr, "remove reusable workspace fixtures: %v\n", err)
		if exitCode == 0 {
			exitCode = 1
		}
	}
	if err := os.RemoveAll(templateRoot); err != nil {
		fmt.Fprintf(os.Stderr, "remove pristine workspace root: %v\n", err)
		if exitCode == 0 {
			exitCode = 1
		}
	}
	return exitCode
}

// commandTestHelperChild reports whether this process is a helper started by
// a test: a gate or Artifact helper, which runs under a gate's environment
// allowlist with no PATH for git. Neither uses the pristine workspace.
func commandTestHelperChild() bool {
	return os.Getenv("FID_GATE_HELPER") != "" || os.Getenv("FID_ARTIFACT_HELPER") != ""
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
	if err := gittest.CheckRepository(destination); err != nil {
		t.Fatalf("clone pristine workspace: %v", err)
	}
}

func cloneReusableWorkspaceFixture(t *testing.T, key string, build func() (string, []string)) (string, []string) {
	t.Helper()
	fixture := reusableWorkspaceFixtureFor(key)
	fixture.Lock()
	defer fixture.Unlock()
	if fixture.root == "" {
		workspace, values := build()
		if err := gittest.CheckRepository(workspace); err != nil {
			t.Fatalf("build reusable workspace fixture: %v", err)
		}
		templateRoot, err := os.MkdirTemp("", "fid-reusable-workspace-")
		if err != nil {
			t.Fatalf("create reusable workspace fixture: %v", err)
		}
		if err := copyWorkspaceTree(workspace, templateRoot); err != nil {
			_ = os.RemoveAll(templateRoot)
			t.Fatalf("snapshot reusable workspace fixture: %v", err)
		}
		fixture.root = templateRoot
		fixture.values = append([]string(nil), values...)
		return workspace, append([]string(nil), values...)
	}

	workspace := t.TempDir()
	cloneWorkspaceTreeInto(t, fixture.root, workspace)
	if err := gittest.CheckRepository(workspace); err != nil {
		t.Fatalf("clone reusable workspace fixture: %v", err)
	}
	return workspace, append([]string(nil), fixture.values...)
}

func reusableWorkspaceFixtureKey(t *testing.T, prefix string, value any) string {
	t.Helper()
	payload, err := json.Marshal(value)
	if err != nil {
		t.Fatalf("encode reusable workspace fixture key: %v", err)
	}
	return prefix + ":" + string(payload)
}

func reusableWorkspaceFixtureFor(key string) *reusableWorkspaceFixture {
	reusableWorkspaceFixtures.Lock()
	defer reusableWorkspaceFixtures.Unlock()
	fixture := reusableWorkspaceFixtures.items[key]
	if fixture == nil {
		fixture = &reusableWorkspaceFixture{}
		reusableWorkspaceFixtures.items[key] = fixture
	}
	return fixture
}

func cleanupReusableWorkspaceFixtures() error {
	reusableWorkspaceFixtures.Lock()
	defer reusableWorkspaceFixtures.Unlock()
	var firstErr error
	for key, fixture := range reusableWorkspaceFixtures.items {
		if fixture.root != "" {
			if err := os.RemoveAll(fixture.root); err != nil && firstErr == nil {
				firstErr = fmt.Errorf("remove %q: %w", key, err)
			}
		}
		delete(reusableWorkspaceFixtures.items, key)
	}
	return firstErr
}

func cloneWorkspaceTreeInto(t *testing.T, source, destination string) {
	t.Helper()
	if err := copyWorkspaceTree(source, destination); err != nil {
		t.Fatalf("clone reusable workspace fixture: %v", err)
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

type workspacePersistenceSnapshot struct {
	database []byte
	journal  []byte
}

func snapshotWorkspacePersistence(t *testing.T, workspace string) workspacePersistenceSnapshot {
	t.Helper()
	return workspacePersistenceSnapshot{
		database: mustReadCommandFile(t, config.WorkspacePathsForRoot(workspace).DatabasePath),
		journal:  mustReadCommandFile(t, mustJournalPath(t, workspace)),
	}
}

func (snapshot workspacePersistenceSnapshot) equal(other workspacePersistenceSnapshot) bool {
	return bytes.Equal(snapshot.database, other.database) && bytes.Equal(snapshot.journal, other.journal)
}

func (snapshot workspacePersistenceSnapshot) eventCount() int {
	return len(splitJournalLines(snapshot.journal))
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
	t.Parallel()
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
	t.Parallel()
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

func TestReusableWorkspaceFixtureBuildsOnceAndKeepsClonesIsolated(t *testing.T) {
	t.Parallel()
	var builds atomic.Int32
	key := filepath.Join(t.TempDir(), "reusable-fixture")
	build := func() (string, []string) {
		builds.Add(1)
		workspace := clonePristineWorkspace(t)
		if err := os.WriteFile(filepath.Join(workspace, "checkpoint.txt"), []byte("cached\n"), 0o600); err != nil {
			t.Fatalf("write reusable fixture checkpoint: %v", err)
		}
		return workspace, []string{"fixture-id"}
	}

	first, firstValues := cloneReusableWorkspaceFixture(t, key, build)
	second, secondValues := cloneReusableWorkspaceFixture(t, key, build)
	if builds.Load() != 1 {
		t.Fatalf("reusable fixture builds = %d, want 1", builds.Load())
	}
	if first == second || firstValues[0] != "fixture-id" || secondValues[0] != "fixture-id" {
		t.Fatalf("reusable fixture clones = (%q, %v) and (%q, %v), want distinct roots and preserved metadata", first, firstValues, second, secondValues)
	}

	firstValues[0] = "mutated"
	if err := os.WriteFile(filepath.Join(first, "checkpoint.txt"), []byte("changed\n"), 0o600); err != nil {
		t.Fatalf("mutate first reusable fixture clone: %v", err)
	}
	third, thirdValues := cloneReusableWorkspaceFixture(t, key, build)
	if builds.Load() != 1 || thirdValues[0] != "fixture-id" {
		t.Fatalf("reused fixture = builds %d metadata %v, want one build and immutable metadata", builds.Load(), thirdValues)
	}
	if got := string(readWorkspaceFixtureFile(t, filepath.Join(second, "checkpoint.txt"))); got != "cached\n" {
		t.Fatalf("second reusable fixture checkpoint = %q, want cached state", got)
	}
	if got := string(readWorkspaceFixtureFile(t, filepath.Join(third, "checkpoint.txt"))); got != "cached\n" {
		t.Fatalf("third reusable fixture checkpoint = %q, want cached state", got)
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
