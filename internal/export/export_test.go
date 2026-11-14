package export

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/mbell-nzl/fid/internal/bootstrap"
	"github.com/mbell-nzl/fid/internal/config"
	"github.com/mbell-nzl/fid/internal/db"
	"github.com/mbell-nzl/fid/internal/gittest"
	"github.com/mbell-nzl/fid/internal/journal"
)

func TestCreateIsByteIdenticalWithCompleteDeterministicInventory(t *testing.T) {
	t.Parallel()
	workspace := gittest.NewWorkspace(t)
	if _, err := bootstrap.InitializeWorkspace(bootstrap.InitOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testMigrationsDir(t),
		Prefix:        "FID",
	}); err != nil {
		t.Fatalf("initialize workspace: %v", err)
	}
	artifactEntry := addLeanArtifactExportFacts(t, workspace)

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
	if sources, ok := manifest["sources"].([]any); !ok || len(sources) != 0 {
		t.Fatalf("manifest sources = %#v, want no kernel sources", manifest["sources"])
	}
}

func TestSchema26ArchiveImportsAndSchema18ManifestFailsClosed(t *testing.T) {
	t.Parallel()
	source := gittest.NewWorkspace(t)
	if _, err := bootstrap.InitializeWorkspace(bootstrap.InitOptions{WorkspaceRoot: source, MigrationsDir: testMigrationsDir(t), Prefix: "FID"}); err != nil {
		t.Fatal(err)
	}
	addLeanArtifactExportFacts(t, source)
	archivePath := filepath.Join(t.TempDir(), "schema26.zip")
	if _, err := Create(CreateOptions{WorkspaceRoot: source, ArchivePath: archivePath}); err != nil {
		t.Fatal(err)
	}
	if got := archiveManifest(t, archivePath)["projection_schema"]; got != float64(26) {
		t.Fatalf("manifest projection_schema = %v, want 26", got)
	}
	if _, err := Validate(ValidateOptions{ArchivePath: archivePath, MigrationsDir: testMigrationsDir(t)}); err != nil {
		t.Fatalf("validate schema-26 archive: %#v", err)
	}
	if _, err := Import(ImportOptions{ArchivePath: archivePath, RepositoryRoot: gittest.NewWorkspace(t), MigrationsDir: testMigrationsDir(t), IdempotencyKey: "schema26"}); err != nil {
		t.Fatal(err)
	}
	entries := archiveContents(t, archivePath)
	manifest := decodeManifestMap(t, entries["manifest.json"])
	manifest["projection_schema"] = json.Number("18")
	entries["manifest.json"] = encodeManifestMap(t, manifest)
	invalidArchive := filepath.Join(t.TempDir(), "schema18.zip")
	writeTestArchive(t, invalidArchive, entries)
	target := gittest.NewWorkspace(t)
	before := workspaceTreeDigest(t, target)
	_, err := Import(ImportOptions{ArchivePath: invalidArchive, RepositoryRoot: target, MigrationsDir: testMigrationsDir(t), IdempotencyKey: "schema18"})
	assertArchiveCheck(t, err, "schema-unsupported")
	invalid, ok := err.(*InvalidArchiveError)
	if !ok {
		t.Fatalf("schema-18 import error = %T, want archive-invalid", err)
	}
	foundManifest := false
	for _, check := range invalid.Checks {
		if check.Code == "schema-unsupported" && check.Entry != nil && *check.Entry == "manifest.json" {
			foundManifest = true
		}
	}
	if !foundManifest {
		t.Fatalf("schema-18 archive checks = %#v, want schema-unsupported on manifest.json", invalid.Checks)
	}
	if after := workspaceTreeDigest(t, target); after != before {
		t.Fatalf("invalid import changed target tree: %s to %s", before, after)
	}
}

func TestCheckResultFieldDigestGuardAndMalformedInsertFailClosedThroughCLI(t *testing.T) {
	t.Parallel()
	workspace := gittest.NewWorkspace(t)
	if _, err := bootstrap.InitializeWorkspace(bootstrap.InitOptions{WorkspaceRoot: workspace, MigrationsDir: testMigrationsDir(t), Prefix: "FID"}); err != nil {
		t.Fatal(err)
	}
	addLeanArtifactExportFacts(t, workspace)
	archive := filepath.Join(t.TempDir(), "control.zip")
	if _, err := Create(CreateOptions{WorkspaceRoot: workspace, ArchivePath: archive}); err != nil {
		t.Fatal(err)
	}
	archiveMembers := archiveContents(t, archive)
	journalFiles, err := filepath.Glob(filepath.Join(config.LedgerDir(workspace), "journal", "*.jsonl"))
	if err != nil || len(journalFiles) == 0 {
		t.Fatalf("journal files = %v, error %v", journalFiles, err)
	}
	original := make(map[string][]byte, len(journalFiles))
	for _, path := range journalFiles {
		original[path], err = os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
	}
	binary := filepath.Join(t.TempDir(), "fid")
	build := exec.Command("go", "build", "-o", binary, "./cmd/fid")
	build.Dir = filepath.Join("..", "..")
	if output, err := build.CombinedOutput(); err != nil {
		t.Fatalf("build fid: %v\n%s", err, output)
	}
	cases := []struct {
		name, eventType string
		mutate          func(*journal.Effect)
		wantExit        int
	}{
		{"valid-check-result", "", nil, 0},
		{"missing-insert-map", "check.started", func(e *journal.Effect) { delete(e.Values, "request_field_digests_json") }, 6},
		{"array-insert-map", "check.started", func(e *journal.Effect) { e.Values["request_field_digests_json"] = `[]` }, 6},
		{"non-digest-insert-map", "check.started", func(e *journal.Effect) { e.Values["request_field_digests_json"] = `{"--actor":"plain"}` }, 6},
		{"mismatched-result-guard", "check.result-recorded", func(e *journal.Effect) {
			e.Preconditions["request_field_digests_json"] = `{"different-field":"sha256:` + strings.Repeat("a", 64) + `"}`
		}, 6},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			for path, contents := range original {
				if err := os.WriteFile(path, contents, 0o600); err != nil {
					t.Fatal(err)
				}
			}
			var journalFile string
			var line int
			if tc.mutate != nil {
				journalFile, line = mutateIdempotencyEffect(t, original, tc.eventType, tc.mutate)
			}
			before := workspaceTreeDigest(t, workspace)
			cmd := exec.Command(binary, "--project", workspace, "--json", "rebuild")
			var stdout, stderr bytes.Buffer
			cmd.Stdout, cmd.Stderr = &stdout, &stderr
			err := cmd.Run()
			if tc.wantExit == 0 {
				if err != nil || stderr.Len() != 0 {
					t.Fatalf("valid check result rebuild = %v, stderr %q", err, stderr.String())
				}
				assertIdempotencyArchiveAdmission(t, binary, archiveMembers, journalFiles, true)
				return
			}
			var exit *exec.ExitError
			if !errors.As(err, &exit) || exit.ExitCode() != tc.wantExit || stdout.Len() != 0 {
				t.Fatalf("rebuild exit = %v, stdout = %q, stderr = %q", err, stdout.String(), stderr.String())
			}
			var failure struct {
				Operation string `json:"operation"`
				Error     struct {
					Code     string `json:"code"`
					Category string `json:"category"`
					Details  struct {
						JournalFile string `json:"journal_file"`
						Line        int    `json:"line"`
					} `json:"details"`
				} `json:"error"`
			}
			if err := json.Unmarshal(stderr.Bytes(), &failure); err != nil {
				t.Fatalf("decode rebuild failure: %v: %s", err, stderr.String())
			}
			if failure.Operation != "project.rebuild" || failure.Error.Code != "journal-replay-stopped" || failure.Error.Category != "integrity" || failure.Error.Details.JournalFile != journalFile || failure.Error.Details.Line != line {
				t.Fatalf("rebuild failure = %+v, want journal-replay-stopped at %s:%d", failure, journalFile, line)
			}
			if after := workspaceTreeDigest(t, workspace); after != before {
				t.Fatalf("rejected rebuild changed target tree: %s to %s", before, after)
			}
			assertIdempotencyArchiveAdmission(t, binary, archiveMembers, journalFiles, false)
		})
	}
}

func mutateIdempotencyEffect(t *testing.T, original map[string][]byte, eventType string, mutate func(*journal.Effect)) (string, int) {
	t.Helper()
	for path, contents := range original {
		lines := bytes.Split(bytes.TrimSuffix(contents, []byte{'\n'}), []byte{'\n'})
		for index, raw := range lines {
			var event journal.Event
			if err := event.UnmarshalJSON(raw); err != nil {
				t.Fatal(err)
			}
			if event.EventType != eventType {
				continue
			}
			for effectIndex := range event.Effects {
				if event.Effects[effectIndex].Table != "idempotency_results" {
					continue
				}
				mutate(&event.Effects[effectIndex])
				encoded, err := journal.MarshalCanonicalJSON(event)
				if err != nil {
					t.Fatal(err)
				}
				lines[index] = encoded
				if err := os.WriteFile(path, append(bytes.Join(lines, []byte{'\n'}), '\n'), 0o600); err != nil {
					t.Fatal(err)
				}
				// fid names the file by absolute path, under the root git
				// reports, so with symlinks resolved.
				return filepath.ToSlash(gittest.ResolvedPath(t, path)), index + 1
			}
		}
	}
	t.Fatalf("missing %s idempotency Effect", eventType)
	return "", 0
}

func TestImportCrashBeforePublicationLeavesNoPartialProject(t *testing.T) {
	t.Parallel()
	source := gittest.NewWorkspace(t)
	if _, err := bootstrap.InitializeWorkspace(bootstrap.InitOptions{WorkspaceRoot: source, MigrationsDir: testMigrationsDir(t), Prefix: "FID"}); err != nil {
		t.Fatalf("initialize source workspace: %v", err)
	}
	addLeanArtifactExportFacts(t, source)
	archive := filepath.Join(t.TempDir(), "recovery.zip")
	if _, err := Create(CreateOptions{WorkspaceRoot: source, ArchivePath: archive}); err != nil {
		t.Fatalf("create source archive: %v", err)
	}
	target := gittest.NewWorkspace(t)
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
	if _, err := os.Lstat(config.LedgerDir(target)); !os.IsNotExist(err) {
		t.Fatalf("target contains a partial project after pre-publication fault: %v", err)
	}
}

func TestValidateRejectsJournalAggregateSequenceCollisionWithoutMutation(t *testing.T) {
	t.Parallel()
	workspace := gittest.NewWorkspace(t)
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
	workspace := gittest.NewWorkspace(t)
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
	fixture := newExportSourceFixture(t, addLeanArtifactExportFacts)
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

// treeListing names every entry under root, directories included, with each
// regular file's contents, so a comparison also catches an empty directory
// left behind.
func treeListing(t *testing.T, root string) string {
	t.Helper()
	var listing strings.Builder
	if err := filepath.WalkDir(root, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		relative, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		listing.WriteString(entry.Type().String() + " " + relative)
		if entry.Type().IsRegular() {
			data, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			listing.WriteString(" " + digest(data))
		}
		listing.WriteByte('\n')
		return nil
	}); err != nil {
		t.Fatalf("walk %q: %v", root, err)
	}
	return listing.String()
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
