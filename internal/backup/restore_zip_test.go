package backup

import (
	"archive/zip"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mbell-nzl/fid/internal/gittest"
)

func TestMain(m *testing.M) {
	os.Exit(gittest.Main(m))
}

// Inspect the rebuilt ZIP while it still exists. Restore removes its
// private directory afterwards, which the command-level test checks.
func TestRebuiltArchiveHasPrivatePermissions(t *testing.T) {
	t.Parallel()
	root := gittest.NewWorkspace(t)
	repository := clone{gitDir: filepath.Join(root, ".git"), workingDir: root,
		environment: gitEnvironment(append(os.Environ(), "HOME="+t.TempDir(), "XDG_CONFIG_HOME="+t.TempDir()))}
	blob, err := repository.output(strings.NewReader("member contents\n"), "hash-object", "-w", "--stdin")
	if err != nil {
		t.Fatal(err)
	}
	tree, err := repository.output(strings.NewReader("100644 blob "+blob+"\tmanifest.json\n"), "mktree")
	if err != nil {
		t.Fatal(err)
	}
	tree, err = repository.output(strings.NewReader("040000 tree "+tree+"\tproject\n"), "mktree")
	if err != nil {
		t.Fatal(err)
	}
	scratch := t.TempDir()
	// A mismatched expected digest leaves the completed private archive
	// available to inspect, without requiring this fixture to be an export.
	_, err = zipFolder(repository, tree, "project", "mismatched", scratch)
	var failed *Error
	if !errors.As(err, &failed) || failed.Kind != KindArchiveInvalid {
		t.Fatalf("rebuild = %v, want an archive digest mismatch", err)
	}
	archive := filepath.Join(scratch, ArchiveName)
	info, err := os.Stat(archive)
	if err != nil || info.Mode().Perm() != 0o600 {
		t.Fatalf("rebuilt archive info = %v, %v, want mode 0600", info, err)
	}
	reader, err := zip.OpenReader(archive)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = reader.Close() }()
	if len(reader.File) != 1 || reader.File[0].Mode().Perm() != 0o644 {
		t.Fatalf("rebuilt archive members = %v, want a 0644 export member", reader.File)
	}
}
