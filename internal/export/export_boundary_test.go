package export

import (
	"archive/zip"
	"os"
	"path/filepath"
	"sort"
	"testing"

	"github.com/mbell-nzl/fid/internal/bootstrap"
)

func TestArchiveRejectsSpecialEntriesAndExpandedMember(t *testing.T) {
	for _, tc := range []struct {
		name  string
		mode  os.FileMode
		large bool
		code  string
	}{
		{name: "regular control", mode: 0o644},
		{name: "symlink", mode: os.ModeSymlink | 0o777, code: "path-unsafe"},
		{name: "device", mode: os.ModeDevice | 0o600, code: "path-unsafe"},
		{name: "directory", mode: os.ModeDir | 0o755, code: "path-unsafe"},
		{name: "expanded non-Artifact limit plus one", mode: 0o644, large: true, code: "resource-limit"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			workspace := t.TempDir()
			if _, err := bootstrap.InitializeWorkspace(bootstrap.InitOptions{WorkspaceRoot: workspace, MigrationsDir: testMigrationsDir(t), Prefix: "FID"}); err != nil {
				t.Fatal(err)
			}
			archivePath := filepath.Join(t.TempDir(), "boundary.zip")
			if _, err := Create(CreateOptions{WorkspaceRoot: workspace, ArchivePath: archivePath}); err != nil {
				t.Fatal(err)
			}
			if _, err := Validate(ValidateOptions{ArchivePath: archivePath, MigrationsDir: testMigrationsDir(t)}); err != nil {
				t.Fatalf("full valid control: %v", err)
			}
			contents := archiveContents(t, archivePath)
			if tc.large {
				payload := make([]byte, (64<<20)+1)
				contents["material.bin"] = payload
				manifest := decodeManifestMap(t, contents["manifest.json"])
				manifest["entries"] = append(manifest["entries"].([]any), map[string]any{"path": "material.bin", "digest": digest(payload), "size_bytes": len(payload), "media_type": "application/octet-stream"})
				contents["manifest.json"] = encodeManifestMap(t, manifest)
			}
			writeBoundaryContents(t, archivePath, contents, tc.mode)
			before := workspaceTreeDigest(t, workspace)
			_, err := Validate(ValidateOptions{ArchivePath: archivePath, MigrationsDir: testMigrationsDir(t)})
			if after := workspaceTreeDigest(t, workspace); after != before {
				t.Fatal("validation changed source project")
			}
			if tc.code == "" {
				if err != nil {
					t.Fatalf("regular control: %v", err)
				}
				return
			}
			assertArchiveCheck(t, err, tc.code)
		})
	}
}

func writeBoundaryContents(t *testing.T, path string, contents map[string][]byte, manifestMode os.FileMode) {
	t.Helper()
	file, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	writer := zip.NewWriter(file)
	names := make([]string, 0, len(contents))
	for name := range contents {
		names = append(names, name)
	}
	sort.Strings(names)
	for _, name := range names {
		header := &zip.FileHeader{Name: name, Method: zip.Deflate, Modified: zipEpoch}
		mode := os.FileMode(0o644)
		if name == "manifest.json" {
			mode = manifestMode
		}
		header.SetMode(mode)
		output, err := writer.CreateHeader(header)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := output.Write(contents[name]); err != nil {
			t.Fatal(err)
		}
	}
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}
	if err := file.Close(); err != nil {
		t.Fatal(err)
	}
}
