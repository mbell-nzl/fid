package export

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/mbell-nzl/fid/internal/bootstrap"
	"github.com/mbell-nzl/fid/internal/config"
	fiddb "github.com/mbell-nzl/fid/internal/db"
	"github.com/mbell-nzl/fid/internal/gittest"
)

// A source fixture is built through real journal writers. Copies are used only
// as input to archive tests; imports and replay under test remain real operations.
// Never use a copied path-bound source as a newly authorized live executor.
type exportSourceFixture struct {
	workspace string
	archive   string
}

func newExportSourceFixture(t *testing.T, populate func(*testing.T, string) string) exportSourceFixture {
	t.Helper()
	workspace := gittest.NewWorkspace(t)
	if _, err := bootstrap.InitializeWorkspace(bootstrap.InitOptions{WorkspaceRoot: workspace, MigrationsDir: testMigrationsDir(t), Prefix: "FID"}); err != nil {
		t.Fatal(err)
	}
	if populate != nil {
		populate(t, workspace)
	}
	archive := filepath.Join(t.TempDir(), "source.zip")
	if _, err := Create(CreateOptions{WorkspaceRoot: workspace, ArchivePath: archive}); err != nil {
		t.Fatal(err)
	}
	return exportSourceFixture{workspace, archive}
}

func (fixture exportSourceFixture) clone(t *testing.T) (string, string) {
	t.Helper()
	workspace := t.TempDir()
	err := filepath.WalkDir(fixture.workspace, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		relative, err := filepath.Rel(fixture.workspace, path)
		if err != nil {
			return err
		}
		info, err := entry.Info()
		if err != nil {
			return err
		}
		target := filepath.Join(workspace, relative)
		if entry.IsDir() {
			return os.MkdirAll(target, info.Mode().Perm())
		}
		if !info.Mode().IsRegular() {
			return fmt.Errorf("unexpected source fixture entry: %s", relative)
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		return os.WriteFile(target, data, info.Mode().Perm())
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := gittest.CheckRepository(workspace); err != nil {
		t.Fatalf("clone export source fixture: %v", err)
	}
	archive := filepath.Join(t.TempDir(), "copy.zip")
	data, err := os.ReadFile(fixture.archive)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(archive, data, 0o600); err != nil {
		t.Fatal(err)
	}
	return workspace, archive
}

func TestExportSourceFixtureCopiesKeepDatabaseAndArchiveMutationsIsolated(t *testing.T) {
	t.Parallel()
	fixture := newExportSourceFixture(t, nil)
	sourceBefore := workspaceTreeDigest(t, fixture.workspace)
	archiveBefore, err := os.ReadFile(fixture.archive)
	if err != nil {
		t.Fatal(err)
	}
	first, firstArchive := fixture.clone(t)
	second, secondArchive := fixture.clone(t)
	if first == second || first == fixture.workspace || firstArchive == secondArchive || firstArchive == fixture.archive {
		t.Fatal("fixture clones share a mutable workspace or archive path")
	}
	for _, root := range []string{first, second} {
		if workspaceTreeDigest(t, root) != sourceBefore {
			t.Fatal("clone did not retain complete source files")
		}
	}
	dbPath := func(root string) string { return config.WorkspacePathsForRoot(root).DatabasePath }
	a, err := os.Stat(dbPath(first))
	if err != nil {
		t.Fatal(err)
	}
	b, err := os.Stat(dbPath(second))
	if err != nil {
		t.Fatal(err)
	}
	if os.SameFile(a, b) {
		t.Fatal("clones share the same physical database")
	}
	conn, err := fiddb.Open(dbPath(first))
	if err != nil {
		t.Fatal(err)
	}
	_, writeErr := conn.Exec("PRAGMA user_version = 77")
	closeErr := conn.Close()
	if writeErr != nil || closeErr != nil {
		t.Fatalf("mutate clone database: %v / %v", writeErr, closeErr)
	}
	if err := os.WriteFile(firstArchive, []byte("corrupted archive"), 0o600); err != nil {
		t.Fatal(err)
	}
	if workspaceTreeDigest(t, first) == sourceBefore {
		t.Fatal("database mutation did not affect its clone")
	}
	third, thirdArchive := fixture.clone(t)
	for _, root := range []string{fixture.workspace, second, third} {
		if workspaceTreeDigest(t, root) != sourceBefore {
			t.Fatal("clone mutation changed source or sibling")
		}
	}
	for _, path := range []string{fixture.archive, secondArchive, thirdArchive} {
		data, err := os.ReadFile(path)
		if err != nil || string(data) != string(archiveBefore) {
			t.Fatalf("archive copy changed: %s: %v", path, err)
		}
	}
	if _, err := Validate(ValidateOptions{ArchivePath: thirdArchive, MigrationsDir: testMigrationsDir(t)}); err != nil {
		t.Fatalf("untouched cloned archive invalid: %v", err)
	}
}
