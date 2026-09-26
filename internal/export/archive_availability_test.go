package export

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mbell-nzl/fid/internal/bootstrap"
	fiddb "github.com/mbell-nzl/fid/internal/db"
	"github.com/mbell-nzl/fid/internal/evidence"
)

func TestCreateAndValidatePreserveUnavailableReceiptArtifactWithoutBody(t *testing.T) {
	t.Parallel()
	workspace := t.TempDir()
	if _, err := bootstrap.InitializeWorkspace(bootstrap.InitOptions{WorkspaceRoot: workspace, MigrationsDir: testMigrationsDir(t), Prefix: "FID"}); err != nil {
		t.Fatalf("initialize workspace: %v", err)
	}
	artifactEntry := addJournalBackedExportFixtureFacts(t, workspace)
	artifactDigest := "sha256:" + strings.TrimPrefix(filepath.Base(artifactEntry), "sha256:")
	artifactID := exportFixtureArtifactID(t, workspace, artifactDigest)
	if err := os.Remove(filepath.Join(workspace, ".fid", artifactEntry)); err != nil {
		t.Fatalf("remove retained Artifact body: %v", err)
	}
	verified, err := evidence.VerifyArtifact(evidence.ArtifactVerifyOptions{
		WorkspaceRoot:  workspace,
		MigrationsDir:  testMigrationsDir(t),
		ActorID:        "actor-builder",
		IdempotencyKey: "fixture-artifact-unavailable",
		Request:        evidence.ArtifactVerifyRequest{ArtifactID: artifactID},
	})
	if err != nil {
		t.Fatalf("observe missing Artifact availability: %v", err)
	}
	if !verified.Changed || verified.Artifact.Availability != "missing" {
		t.Fatalf("artifact availability observation = %#v, want changed missing", verified)
	}

	archivePath := filepath.Join(t.TempDir(), "unavailable-artifact.zip")
	if _, err := Create(CreateOptions{WorkspaceRoot: workspace, ArchivePath: archivePath}); err != nil {
		t.Fatalf("create unavailable Artifact archive: %v", err)
	}
	if _, err := Validate(ValidateOptions{ArchivePath: archivePath, MigrationsDir: testMigrationsDir(t)}); err != nil {
		t.Fatalf("validate unavailable Artifact archive: %#v", err)
	}

	entries := archiveContents(t, archivePath)
	if _, found := entries[artifactEntry]; found {
		t.Fatalf("archive retains unavailable Artifact body %q", artifactEntry)
	}
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
	journal := strings.Join(journalArchiveContents(t, entries), "\n")
	if !strings.Contains(journal, `"event_type":"artifact.availability-observed"`) || !strings.Contains(journal, artifactID) || !strings.Contains(journal, `"availability":"missing"`) {
		t.Fatalf("archive journal omits unavailable Artifact observation for %q", artifactID)
	}
}

func exportFixtureArtifactID(t *testing.T, workspace, digest string) string {
	t.Helper()
	database, err := fiddb.Open(filepath.Join(workspace, ".fid", "fid.sqlite"))
	if err != nil {
		t.Fatalf("open fixture projection: %v", err)
	}
	defer func() { _ = database.Close() }()
	var id string
	if err := database.QueryRow(`SELECT id FROM referenced_artifacts WHERE digest = ?`, digest).Scan(&id); err != nil {
		t.Fatalf("read retained fixture Artifact: %v", err)
	}
	return id
}

func journalArchiveContents(t *testing.T, entries map[string][]byte) []string {
	t.Helper()
	contents := make([]string, 0)
	for path, data := range entries {
		if strings.HasPrefix(path, "journal/") {
			contents = append(contents, string(data))
		}
	}
	return contents
}
