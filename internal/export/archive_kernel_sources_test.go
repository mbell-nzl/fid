package export

import (
	"path/filepath"
	"testing"
)

// Source inventories belong to the retiring kernel profile; keep their
// assertions until RESET 1.5d1 closes that profile.
func TestKernelArchivePreservesSourceMetadataWithoutProse(t *testing.T) {
	t.Parallel()
	fixture := newExportSourceFixture(t, addKernelExportFacts)
	archivePath := filepath.Join(t.TempDir(), "sources.zip")
	if _, err := Create(CreateOptions{WorkspaceRoot: fixture.workspace, ArchivePath: archivePath}); err != nil {
		t.Fatal(err)
	}
	if _, err := Validate(ValidateOptions{ArchivePath: archivePath, MigrationsDir: testMigrationsDir(t)}); err != nil {
		t.Fatal(err)
	}
	entries := archiveContents(t, archivePath)
	if _, found := entries["project-docs/gates/fixture.json"]; found {
		t.Fatal("archive includes Source file prose instead of source metadata only")
	}
	manifest := archiveManifest(t, archivePath)
	sources, ok := manifest["sources"].([]any)
	if !ok || len(sources) != 1 {
		t.Fatalf("manifest Sources = %#v, want one source metadata object", manifest["sources"])
	}
	source, ok := sources[0].(map[string]any)
	if !ok || source["path"] != "project-docs/gates/fixture.json" || source["source_type"] != "project-file" {
		t.Fatalf("manifest Source = %#v, want canonical project-file metadata", sources[0])
	}
}
