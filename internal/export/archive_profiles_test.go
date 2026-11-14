package export

import (
	"encoding/json"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/mbell-nzl/fid/internal/gittest"
)

func TestCreatedManifestHasOnlyPortableExportCapability(t *testing.T) {
	t.Parallel()
	fixture := newExportSourceFixture(t, nil)
	got := archiveManifest(t, fixture.archive)["capabilities"]
	if want := expectedLeanCapabilities(t); !reflect.DeepEqual(got, want) {
		t.Fatalf("created capability inventory = %#v, want %#v", got, want)
	}
}

// Independent wire contract: never derive this inventory from capabilitySnapshot.
func expectedLeanCapabilities(t *testing.T) []any {
	t.Helper()
	var expected []any
	if err := json.Unmarshal([]byte(`[{"available":true,"capability_id":"export.portable.v1","limits":{"create":true,"import":true,"validate":true,"import_mode":"empty-target"},"schema_max":1,"schema_min":1,"trust_source":"fid-journal"}]`), &expected); err != nil {
		t.Fatal(err)
	}
	return expected
}

func TestLeanArchiveRefusesAddedRetiredCapabilityWithoutMutation(t *testing.T) {
	t.Parallel()
	fixture := newExportSourceFixture(t, nil)
	if _, err := Validate(ValidateOptions{ArchivePath: fixture.archive, MigrationsDir: testMigrationsDir(t)}); err != nil {
		t.Fatalf("valid lean control: %v", err)
	}
	entries := archiveContents(t, fixture.archive)
	manifest := decodeManifestMap(t, entries["manifest.json"])
	manifest["capabilities"] = append(expectedLeanCapabilities(t), map[string]any{
		"available": true, "capability_id": "attention.v1", "limits": map[string]any{"state": true},
		"schema_max": 1, "schema_min": 1, "trust_source": "fid-journal",
	})
	entries["manifest.json"] = encodeManifestMap(t, manifest)
	archive := filepath.Join(t.TempDir(), "retired-capability.zip")
	writeTestArchive(t, archive, entries)
	target := gittest.NewWorkspace(t)
	sourceBefore, targetBefore := workspaceTreeDigest(t, fixture.workspace), workspaceTreeDigest(t, target)
	archiveBefore := workspaceTreeDigest(t, filepath.Dir(archive))
	for _, operation := range []string{"validate", "import"} {
		var err error
		if operation == "validate" {
			_, err = Validate(ValidateOptions{ArchivePath: archive, MigrationsDir: testMigrationsDir(t)})
		} else {
			_, err = Import(ImportOptions{ArchivePath: archive, RepositoryRoot: target, MigrationsDir: testMigrationsDir(t), IdempotencyKey: "retired-capability"})
		}
		assertArchiveCheck(t, err, "capability-snapshot-mismatch")
		if workspaceTreeDigest(t, fixture.workspace) != sourceBefore || workspaceTreeDigest(t, target) != targetBefore || workspaceTreeDigest(t, filepath.Dir(archive)) != archiveBefore {
			t.Fatalf("%s changed files", operation)
		}
	}
}

func TestRetiredArchiveMemberCaseVariantsRefusedWithoutMutation(t *testing.T) {
	t.Parallel()
	fixture := newExportSourceFixture(t, nil)
	for _, member := range []string{"runs/retired.json", "RUNS/retired.json", "Runs/retired.json", "sources/retired.json", "SOURCES/retired.json"} {
		t.Run(member, func(t *testing.T) {
			t.Parallel()
			entries := archiveContents(t, fixture.archive)
			path := archivePathWithPrefix(t, entries, "journal/")
			entries[member] = entries[path]
			delete(entries, path)
			manifest := decodeManifestMap(t, entries["manifest.json"])
			renameManifestEntry(t, manifest, path, member)
			entries["manifest.json"] = encodeManifestMap(t, manifest)
			archive := filepath.Join(t.TempDir(), "retired-member.zip")
			writeTestArchive(t, archive, entries)
			target := gittest.NewWorkspace(t)
			sourceBefore, targetBefore := workspaceTreeDigest(t, fixture.workspace), workspaceTreeDigest(t, target)
			archiveBefore := workspaceTreeDigest(t, filepath.Dir(archive))
			_, err := Validate(ValidateOptions{ArchivePath: archive, MigrationsDir: testMigrationsDir(t)})
			assertArchiveCheck(t, err, "entry-extra")
			_, err = Import(ImportOptions{ArchivePath: archive, RepositoryRoot: target, MigrationsDir: testMigrationsDir(t), IdempotencyKey: "retired-member"})
			assertArchiveCheck(t, err, "entry-extra")
			if workspaceTreeDigest(t, fixture.workspace) != sourceBefore || workspaceTreeDigest(t, target) != targetBefore || workspaceTreeDigest(t, filepath.Dir(archive)) != archiveBefore {
				t.Fatal("refusal changed files")
			}
		})
	}
}
