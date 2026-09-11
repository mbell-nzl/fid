package workflow

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/mbell-nzl/fid/internal/bootstrap"
)

var workflowPristineWorkspace string

// Workflow tests exercise independent mutations of an initialized workspace.
// Build the real bootstrap once, then copy its files (never share a database).
// Bootstrap/migration behavior retains its dedicated tests in those packages.
func TestMain(m *testing.M) {
	os.Exit(runWorkflowTestSuite(m))
}

func runWorkflowTestSuite(m *testing.M) int {
	root, err := os.MkdirTemp("", "fid-workflow-pristine-")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	defer func() { _ = os.RemoveAll(root) }()
	_, file, _, _ := runtime.Caller(0)
	migrationsDir := filepath.Join(filepath.Dir(file), "..", "..", "migrations")
	if _, err := bootstrap.InitializeWorkspace(bootstrap.InitOptions{WorkspaceRoot: root, MigrationsDir: migrationsDir}); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	workflowPristineWorkspace = root
	return m.Run()
}

func cloneWorkflowPristine(destination string) error {
	return filepath.WalkDir(workflowPristineWorkspace, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		relative, err := filepath.Rel(workflowPristineWorkspace, path)
		if err != nil {
			return err
		}
		info, err := entry.Info()
		if err != nil {
			return err
		}
		target := filepath.Join(destination, relative)
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
