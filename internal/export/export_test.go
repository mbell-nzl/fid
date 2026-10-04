package export

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"testing"

	"github.com/mbell-nzl/fid/internal/bootstrap"
	"github.com/mbell-nzl/fid/internal/db"
	"github.com/mbell-nzl/fid/internal/journal"
)

func TestCreateIsByteIdenticalWithCompleteDeterministicInventory(t *testing.T) {
	t.Parallel()
	workspace := t.TempDir()
	if _, err := bootstrap.InitializeWorkspace(bootstrap.InitOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testMigrationsDir(t),
		Prefix:        "FID",
	}); err != nil {
		t.Fatalf("initialize workspace: %v", err)
	}
	artifactEntry := addJournalBackedExportFixtureFacts(t, workspace)

	first := filepath.Join(t.TempDir(), "first.zip")
	second := filepath.Join(t.TempDir(), "second.zip")
	for _, archivePath := range []string{first, second} {
		if _, err := Create(CreateOptions{WorkspaceRoot: workspace, ArchivePath: archivePath}); err != nil {
			t.Fatalf("create %q: %v", archivePath, err)
		}
	}

	firstBytes, err := os.ReadFile(first)
	if err != nil {
		t.Fatalf("read first archive: %v", err)
	}
	secondBytes, err := os.ReadFile(second)
	if err != nil {
		t.Fatalf("read second archive: %v", err)
	}
	if !bytes.Equal(firstBytes, secondBytes) {
		t.Fatal("repeated creates from identical facts produced different archive bytes")
	}
	validation, err := Validate(ValidateOptions{ArchivePath: first, MigrationsDir: testMigrationsDir(t)})
	if err != nil {
		t.Fatalf("validate deterministic archive: %#v", err)
	}
	if validation.Status != "valid" {
		t.Fatalf("validation status = %q, want valid", validation.Status)
	}

	entries := archiveEntries(t, first)
	wantEntries := append(journalEntries(t, workspace), artifactEntry, "manifest.json")
	sort.Strings(wantEntries)
	if !reflect.DeepEqual(entries, wantEntries) {
		t.Fatalf("archive entries = %v, want %v", entries, wantEntries)
	}
	manifest := archiveManifest(t, first)
	if schema, ok := manifest["projection_schema"].(float64); !ok || int(schema) != db.ActiveMigrationVersion {
		t.Fatalf("manifest projection_schema = %#v, want active schema %d", manifest["projection_schema"], db.ActiveMigrationVersion)
	}
	if sources, ok := manifest["sources"].([]any); !ok || len(sources) != 1 {
		t.Fatalf("manifest sources = %#v, want one canonical source", manifest["sources"])
	}
}

func TestSchema19ArchiveImportsAndSchema18ManifestFailsClosed(t *testing.T) {
	t.Parallel()
	source := t.TempDir()
	if _, err := bootstrap.InitializeWorkspace(bootstrap.InitOptions{WorkspaceRoot: source, MigrationsDir: testMigrationsDir(t), Prefix: "FID"}); err != nil {
		t.Fatal(err)
	}
	addJournalBackedExportFixtureFacts(t, source)
	archivePath := filepath.Join(t.TempDir(), "schema19.zip")
	if _, err := Create(CreateOptions{WorkspaceRoot: source, ArchivePath: archivePath}); err != nil {
		t.Fatal(err)
	}
	if got := archiveManifest(t, archivePath)["projection_schema"]; got != float64(19) {
		t.Fatalf("manifest projection_schema = %v, want 19", got)
	}
	if _, err := Validate(ValidateOptions{ArchivePath: archivePath, MigrationsDir: testMigrationsDir(t)}); err != nil {
		t.Fatalf("validate schema-19 archive: %#v", err)
	}
	if _, err := Import(ImportOptions{ArchivePath: archivePath, RepositoryRoot: t.TempDir(), MigrationsDir: testMigrationsDir(t), IdempotencyKey: "schema19"}); err != nil {
		t.Fatal(err)
	}
	entries := archiveContents(t, archivePath)
	manifest := decodeManifestMap(t, entries["manifest.json"])
	manifest["projection_schema"] = json.Number("18")
	entries["manifest.json"] = encodeManifestMap(t, manifest)
	invalidArchive := filepath.Join(t.TempDir(), "schema18.zip")
	writeTestArchive(t, invalidArchive, entries)
	target := t.TempDir()
	before := workspaceTreeDigest(t, target)
	_, err := Import(ImportOptions{ArchivePath: invalidArchive, RepositoryRoot: target, MigrationsDir: testMigrationsDir(t), IdempotencyKey: "schema18"})
	assertArchiveCheck(t, err, "schema-unsupported")
	if after := workspaceTreeDigest(t, target); after != before {
		t.Fatalf("invalid import changed target tree: %s to %s", before, after)
	}
}

func TestImportCrashBeforePublicationLeavesNoPartialProject(t *testing.T) {
	t.Parallel()
	source := t.TempDir()
	if _, err := bootstrap.InitializeWorkspace(bootstrap.InitOptions{WorkspaceRoot: source, MigrationsDir: testMigrationsDir(t), Prefix: "FID"}); err != nil {
		t.Fatalf("initialize source workspace: %v", err)
	}
	addJournalBackedExportFixtureFacts(t, source)
	archive := filepath.Join(t.TempDir(), "recovery.zip")
	if _, err := Create(CreateOptions{WorkspaceRoot: source, ArchivePath: archive}); err != nil {
		t.Fatalf("create source archive: %v", err)
	}
	target := t.TempDir()
	_, err := Import(ImportOptions{
		ArchivePath: archive, RepositoryRoot: target, MigrationsDir: testMigrationsDir(t), IdempotencyKey: "import-crash",
		Fault: func(stage ImportStage) error {
			if stage == ImportBeforePublication {
				return ErrImportFault
			}
			return nil
		},
	})
	if !errors.Is(err, ErrImportFault) {
		t.Fatalf("Import() error = %v, want injected publication fault", err)
	}
	if _, err := os.Lstat(filepath.Join(target, ".fid")); !os.IsNotExist(err) {
		t.Fatalf("target contains a partial project after pre-publication fault: %v", err)
	}
}

func TestValidateRejectsJournalAggregateSequenceCollisionWithoutMutation(t *testing.T) {
	t.Parallel()
	workspace := t.TempDir()
	if _, err := bootstrap.InitializeWorkspace(bootstrap.InitOptions{WorkspaceRoot: workspace, MigrationsDir: testMigrationsDir(t), Prefix: "FID"}); err != nil {
		t.Fatalf("initialize workspace: %v", err)
	}
	archivePath := filepath.Join(t.TempDir(), "archive.zip")
	if _, err := Create(CreateOptions{WorkspaceRoot: workspace, ArchivePath: archivePath}); err != nil {
		t.Fatalf("create archive: %v", err)
	}
	entries := archiveContents(t, archivePath)
	journalPath := journalEntries(t, workspace)[0]
	entries[journalPath] = append(append([]byte{}, entries[journalPath]...), entries[journalPath]...)
	manifest := decodeManifestMap(t, entries["manifest.json"])
	manifest["event_count"] = json.Number("2")
	updateManifestEntry(t, manifest, journalPath, entries[journalPath])
	entries["manifest.json"] = encodeManifestMap(t, manifest)
	writeTestArchive(t, archivePath, entries)

	before := workspaceTreeDigest(t, workspace)
	_, err := Validate(ValidateOptions{ArchivePath: archivePath, MigrationsDir: testMigrationsDir(t)})
	after := workspaceTreeDigest(t, workspace)
	if before != after {
		t.Fatal("validate mutated project state")
	}
	assertArchiveCheck(t, err, "journal-invalid")
}

// This retained negative covers malformed Source metadata; the real Recovery
// Event-reference regression is separate.
func TestArchiveRejectsUnresolvedSourceMetadata(t *testing.T) {
	t.Parallel()
	workspace := t.TempDir()
	if _, err := bootstrap.InitializeWorkspace(bootstrap.InitOptions{WorkspaceRoot: workspace, MigrationsDir: testMigrationsDir(t), Prefix: "FID"}); err != nil {
		t.Fatal(err)
	}
	archivePath := filepath.Join(t.TempDir(), "reference.zip")
	if _, err := Create(CreateOptions{WorkspaceRoot: workspace, ArchivePath: archivePath}); err != nil {
		t.Fatal(err)
	}
	entries := archiveContents(t, archivePath)
	manifest := decodeManifestMap(t, entries["manifest.json"])
	manifest["sources"] = []any{map[string]any{"recovery_id": "REC-missing"}}
	entries["manifest.json"] = encodeManifestMap(t, manifest)
	writeTestArchive(t, archivePath, entries)
	// A Source inventory must match the journal-backed Source facts.
	_, err := Validate(ValidateOptions{ArchivePath: archivePath, MigrationsDir: testMigrationsDir(t)})
	assertArchiveCheck(t, err, "reference-invalid")
}

func TestValidateRejectsEveryArchiveInvalidCheckWithoutMutation(t *testing.T) {
	t.Parallel()
	fixture := newExportSourceFixture(t, addJournalBackedExportFixtureFacts)
	for _, corruption := range []string{
		"artifact-reference-invalid",
		"capability-snapshot-mismatch",
		"digest-mismatch",
		"duplicate-entry",
		"entry-extra",
		"entry-missing",
		"journal-invalid",
		"manifest-invalid",
		"path-unsafe",
		"size-mismatch",
		"terminal-hash-mismatch",
	} {
		t.Run(corruption, func(t *testing.T) {
			workspace, archivePath := fixture.clone(t)
			corruptArchiveForCheck(t, corruption, workspace, archivePath)

			before := workspaceTreeDigest(t, workspace)
			_, err := Validate(ValidateOptions{ArchivePath: archivePath, MigrationsDir: testMigrationsDir(t)})
			after := workspaceTreeDigest(t, workspace)
			if before != after {
				t.Fatal("validate mutated project state")
			}
			assertArchiveCheck(t, err, corruption)
		})
	}
}

func corruptArchiveForCheck(t *testing.T, code, workspace, archivePath string) {
	t.Helper()
	entries := archiveContents(t, archivePath)
	manifest := decodeManifestMap(t, entries["manifest.json"])
	artifactPath := archivePathWithPrefix(t, entries, "artifacts/")
	journalPath := journalEntries(t, workspace)[0]

	switch code {
	case "artifact-reference-invalid":
		wrongPath := "artifacts/0000000000000000000000000000000000000000000000000000000000000000"
		entries[wrongPath] = entries[artifactPath]
		delete(entries, artifactPath)
		renameManifestEntry(t, manifest, artifactPath, wrongPath)
	case "capability-snapshot-mismatch":
		manifest["capabilities"] = []any{}
	case "digest-mismatch":
		entries[artifactPath] = append(entries[artifactPath], 'x')
	case "duplicate-entry":
		writeTestArchiveEntries(t, archivePath, append(archiveEntriesForContents(entries), archiveEntry{entry: entry{Path: "manifest.json"}, data: entries["manifest.json"]}))
		return
	case "entry-extra":
		entries["unexpected.txt"] = []byte("unexpected\n")
	case "entry-missing":
		delete(entries, artifactPath)
	case "journal-invalid":
		entries[journalPath] = append(append([]byte{}, entries[journalPath]...), entries[journalPath]...)
		manifest["event_count"] = json.Number("2")
		updateManifestEntry(t, manifest, journalPath, entries[journalPath])
	case "manifest-invalid":
		manifest["schema"] = "invalid"
	case "path-unsafe":
		entries["../unsafe.txt"] = []byte("unsafe\n")
	case "size-mismatch":
		setManifestEntrySize(t, manifest, artifactPath, len(entries[artifactPath])+1)
	case "terminal-hash-mismatch":
		manifest["terminal_event_hash"] = "sha256:0000000000000000000000000000000000000000000000000000000000000000"
	default:
		t.Fatalf("unexpected archive corruption code %q", code)
	}
	entries["manifest.json"] = encodeManifestMap(t, manifest)
	writeTestArchive(t, archivePath, entries)
}

func archiveManifest(t *testing.T, archivePath string) map[string]any {
	t.Helper()
	reader, err := zip.OpenReader(archivePath)
	if err != nil {
		t.Fatalf("open archive: %v", err)
	}
	defer func() { _ = reader.Close() }()
	for _, entry := range reader.File {
		if entry.Name != "manifest.json" {
			continue
		}
		input, err := entry.Open()
		if err != nil {
			t.Fatalf("open manifest: %v", err)
		}
		defer func() { _ = input.Close() }()
		var manifest map[string]any
		if err := json.NewDecoder(input).Decode(&manifest); err != nil {
			t.Fatalf("decode manifest: %v", err)
		}
		return manifest
	}
	t.Fatal("manifest entry missing")
	return nil
}

func archiveContents(t *testing.T, archivePath string) map[string][]byte {
	t.Helper()
	reader, err := zip.OpenReader(archivePath)
	if err != nil {
		t.Fatalf("open archive: %v", err)
	}
	defer func() { _ = reader.Close() }()
	entries := make(map[string][]byte, len(reader.File))
	for _, entry := range reader.File {
		input, err := entry.Open()
		if err != nil {
			t.Fatalf("open entry %q: %v", entry.Name, err)
		}
		data, err := io.ReadAll(input)
		if err != nil {
			_ = input.Close()
			t.Fatalf("read entry %q: %v", entry.Name, err)
		}
		if err := input.Close(); err != nil {
			t.Fatalf("close entry %q: %v", entry.Name, err)
		}
		entries[entry.Name] = data
	}
	return entries
}

func decodeManifestMap(t *testing.T, data []byte) map[string]any {
	t.Helper()
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.UseNumber()
	var manifest map[string]any
	if err := decoder.Decode(&manifest); err != nil {
		t.Fatalf("decode manifest map: %v", err)
	}
	return manifest
}

func encodeManifestMap(t *testing.T, manifest map[string]any) []byte {
	t.Helper()
	data, err := journal.MarshalCanonicalJSON(manifest)
	if err != nil {
		t.Fatalf("encode manifest map: %v", err)
	}
	return data
}

func updateManifestEntry(t *testing.T, manifest map[string]any, path string, data []byte) {
	t.Helper()
	entry := manifestEntry(t, manifest, path)
	entry["digest"] = digest(data)
	entry["size_bytes"] = json.Number(fmt.Sprint(len(data)))
}

func renameManifestEntry(t *testing.T, manifest map[string]any, oldPath, newPath string) {
	t.Helper()
	manifestEntry(t, manifest, oldPath)["path"] = newPath
}

func setManifestEntrySize(t *testing.T, manifest map[string]any, path string, size int) {
	t.Helper()
	manifestEntry(t, manifest, path)["size_bytes"] = json.Number(fmt.Sprint(size))
}

func manifestEntry(t *testing.T, manifest map[string]any, path string) map[string]any {
	t.Helper()
	entries, ok := manifest["entries"].([]any)
	if !ok {
		t.Fatal("manifest entries are absent")
	}
	for _, value := range entries {
		entry, ok := value.(map[string]any)
		if ok && entry["path"] == path {
			return entry
		}
	}
	t.Fatalf("manifest entry %q is absent", path)
	return nil
}

func writeTestArchive(t *testing.T, path string, contents map[string][]byte) {
	t.Helper()
	writeTestArchiveEntries(t, path, archiveEntriesForContents(contents))
}

func writeTestArchiveEntries(t *testing.T, path string, entries []archiveEntry) {
	t.Helper()
	if err := writeArchive(path, entries); err != nil {
		t.Fatalf("write archive: %v", err)
	}
}

func archiveEntriesForContents(contents map[string][]byte) []archiveEntry {
	entries := make([]archiveEntry, 0, len(contents))
	for name, data := range contents {
		entries = append(entries, archiveEntry{entry: entry{Path: name}, data: data})
	}
	sort.Slice(entries, func(i, j int) bool { return entries[i].Path < entries[j].Path })
	return entries
}

func archivePathWithPrefix(t *testing.T, entries map[string][]byte, prefix string) string {
	t.Helper()
	for path := range entries {
		if len(path) >= len(prefix) && path[:len(prefix)] == prefix {
			return path
		}
	}
	t.Fatalf("archive entry with prefix %q is absent", prefix)
	return ""
}

func workspaceTreeDigest(t *testing.T, root string) string {
	t.Helper()
	paths := make([]string, 0)
	if err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.Mode().IsRegular() {
			return nil
		}
		relative, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		paths = append(paths, relative)
		return nil
	}); err != nil {
		t.Fatalf("walk workspace: %v", err)
	}
	sort.Strings(paths)
	var snapshot bytes.Buffer
	for _, relative := range paths {
		data, err := os.ReadFile(filepath.Join(root, relative))
		if err != nil {
			t.Fatalf("read workspace file %q: %v", relative, err)
		}
		snapshot.WriteString(relative)
		snapshot.WriteByte(0)
		snapshot.Write(data)
		snapshot.WriteByte(0)
	}
	return digest(snapshot.Bytes())
}

func assertArchiveCheck(t *testing.T, err error, want string) {
	t.Helper()
	invalid, ok := err.(*InvalidArchiveError)
	if !ok {
		t.Fatalf("validate error = %v, want invalid archive with %q", err, want)
	}
	for _, check := range invalid.Checks {
		if check.Code == want {
			return
		}
	}
	t.Fatalf("archive checks = %#v, want %q", invalid.Checks, want)
}

func journalEntries(t *testing.T, workspace string) []string {
	t.Helper()
	return []string{"journal/00000001.jsonl"}
}

func archiveEntries(t *testing.T, archivePath string) []string {
	t.Helper()
	reader, err := zip.OpenReader(archivePath)
	if err != nil {
		t.Fatalf("open archive: %v", err)
	}
	defer func() { _ = reader.Close() }()
	entries := make([]string, 0, len(reader.File))
	for _, entry := range reader.File {
		entries = append(entries, entry.Name)
	}
	sort.Strings(entries)
	return entries
}

func testMigrationsDir(t *testing.T) string {
	t.Helper()
	return filepath.Join("..", "..", "migrations")
}
