package workflow

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/mbell-nzl/fid/internal/bootstrap"
	"github.com/mbell-nzl/fid/internal/config"
	"github.com/mbell-nzl/fid/internal/gittest"
)

var workflowPristineWorkspace string

// Workflow tests exercise independent mutations of an initialized workspace.
// Build the real bootstrap once, then copy its files (never share a database).
// Bootstrap/migration behavior retains its dedicated tests in those packages.
func TestMain(m *testing.M) {
	os.Exit(runWorkflowTestSuite(m))
}

func runWorkflowTestSuite(m *testing.M) int {
	if err := gittest.Guard(); err != nil {
		fmt.Fprintf(os.Stderr, "guard git discovery: %v\n", err)
		return 1
	}
	root, err := os.MkdirTemp("", "fid-workflow-pristine-")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	defer func() { _ = os.RemoveAll(root) }()
	if err := gittest.InitRepository(root); err != nil {
		fmt.Fprintf(os.Stderr, "make pristine workspace a git repository: %v\n", err)
		return 1
	}
	_, file, _, _ := runtime.Caller(0)
	migrationsDir := filepath.Join(filepath.Dir(file), "..", "..", "migrations")
	if _, err := bootstrap.InitializeWorkspace(bootstrap.InitOptions{WorkspaceRoot: root, MigrationsDir: migrationsDir}); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	workflowPristineWorkspace = root
	return m.Run()
}

// cloneWorkflowPristine copies the pristine ledger into destination, which
// must already be a guarded git repository. Some fixtures commit to their
// repository first, so the pristine repository itself is never copied.
func cloneWorkflowPristine(destination string) error {
	if err := gittest.CheckRepository(destination); err != nil {
		return err
	}
	source := config.LedgerDir(workflowPristineWorkspace)
	ledger := config.LedgerDir(destination)
	return filepath.WalkDir(source, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		relative, err := filepath.Rel(source, path)
		if err != nil {
			return err
		}
		info, err := entry.Info()
		if err != nil {
			return err
		}
		target := filepath.Join(ledger, relative)
		if entry.IsDir() {
			return os.MkdirAll(target, info.Mode().Perm())
		}
		if !info.Mode().IsRegular() {
			return fmt.Errorf("unexpected workflow fixture file type: %s", relative)
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		return os.WriteFile(target, data, info.Mode().Perm())
	})
}
