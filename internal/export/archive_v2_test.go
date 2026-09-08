package export

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/mbell-nzl/fid/internal/bootstrap"
)

func TestValidateHasNoSourceOrDestinationEffects(t *testing.T) {
	workspace, destination := t.TempDir(), t.TempDir()
	if _, err := bootstrap.InitializeWorkspace(bootstrap.InitOptions{WorkspaceRoot: workspace, MigrationsDir: testMigrationsDir(t), Prefix: "FID"}); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(destination, "sentinel.txt"), []byte("operator files remain intact"), 0o600); err != nil {
		t.Fatal(err)
	}
	archive := filepath.Join(destination, "recovery.zip")
	if _, err := Create(CreateOptions{WorkspaceRoot: workspace, ArchivePath: archive}); err != nil {
		t.Fatal(err)
	}
	for _, invalid := range []bool{false, true} {
		if invalid {
			contents := archiveContents(t, archive)
			m := decodeManifestMap(t, contents["manifest.json"])
			m["project_id"] = "project-foreign"
			contents["manifest.json"] = encodeManifestMap(t, m)
			writeTestArchive(t, archive, contents)
		}
		sourceBefore, destinationBefore := workspaceTreeDigest(t, workspace), workspaceTreeDigest(t, destination)
		_, err := Validate(ValidateOptions{ArchivePath: archive, MigrationsDir: testMigrationsDir(t)})
		if (err != nil) != invalid {
			t.Fatalf("invalid=%v error=%v", invalid, err)
		}
		if workspaceTreeDigest(t, workspace) != sourceBefore || workspaceTreeDigest(t, destination) != destinationBefore {
			t.Fatal("validation changed source or destination")
		}
	}
}

func TestCreateRecoveryManifestV2ConfigAndOrdinalJournal(t *testing.T) {
	workspace := t.TempDir()
	if _, err := bootstrap.InitializeWorkspace(bootstrap.InitOptions{WorkspaceRoot: workspace, MigrationsDir: testMigrationsDir(t), Prefix: "FID"}); err != nil {
		t.Fatal(err)
	}
	archive := filepath.Join(t.TempDir(), "recovery.zip")
	if _, err := Create(CreateOptions{WorkspaceRoot: workspace, ArchivePath: archive}); err != nil {
		t.Fatal(err)
	}
	contents := archiveContents(t, archive)
	m := decodeManifestMap(t, contents["manifest.json"])
	if m["schema"] != "fid.export-manifest.v2" {
		t.Fatalf("schema = %v, want fid.export-manifest.v2", m["schema"])
	}
	project, ok := m["project_config"].(map[string]any)
	if !ok || project["project"].(map[string]any)["id"] != m["project_id"] {
		t.Fatalf("missing identity-bound project configuration: %v", project)
	}
	if len(contents["journal/00000001.jsonl"]) == 0 {
		t.Fatal("canonical ordinal journal chunk is missing")
	}
	if _, err := Validate(ValidateOptions{ArchivePath: archive, MigrationsDir: testMigrationsDir(t)}); err != nil {
		t.Fatalf("validate full v2 control: %v", err)
	}
}

func TestValidateRetainsPinnedV1Profile(t *testing.T) {
	workspace := t.TempDir()
	if _, err := bootstrap.InitializeWorkspace(bootstrap.InitOptions{WorkspaceRoot: workspace, MigrationsDir: testMigrationsDir(t), Prefix: "FID"}); err != nil {
		t.Fatal(err)
	}
	archive := filepath.Join(t.TempDir(), "historical.zip")
	if _, err := Create(CreateOptions{WorkspaceRoot: workspace, ArchivePath: archive}); err != nil {
		t.Fatal(err)
	}
	contents := archiveContents(t, archive)
	m := decodeManifestMap(t, contents["manifest.json"])
	m["schema"], m["projection_schema"] = "fid.export-manifest.v1", 14
	delete(m, "project_config")
	contents["journal/2026-08-30.jsonl"] = contents["journal/00000001.jsonl"]
	delete(contents, "journal/00000001.jsonl")
	renameManifestEntry(t, m, "journal/00000001.jsonl", "journal/2026-08-30.jsonl")
	contents["manifest.json"] = encodeManifestMap(t, m)
	writeTestArchive(t, archive, contents)
	before := workspaceTreeDigest(t, workspace)
	if _, err := Validate(ValidateOptions{ArchivePath: archive, MigrationsDir: testMigrationsDir(t)}); err != nil {
		t.Fatalf("retained v1 profile rejected: %v", err)
	}
	if workspaceTreeDigest(t, workspace) != before {
		t.Fatal("v1 validation mutated project")
	}
}

func TestValidateRecoveryManifestClosedConfigurationAndInventory(t *testing.T) {
	for _, tc := range []struct {
		name, code string
		change     func(map[string]any, map[string][]byte)
	}{
		{"configuration omitted", "configuration-invalid", func(m map[string]any, _ map[string][]byte) { delete(m, "project_config") }},
		{"configuration identity mismatch", "configuration-invalid", func(m map[string]any, _ map[string][]byte) {
			m["project_config"].(map[string]any)["project"].(map[string]any)["id"] = "project-foreign"
		}},
		{"journal project identity mismatch", "reference-invalid", func(m map[string]any, _ map[string][]byte) {
			m["project_id"] = "project-foreign"
			m["project_config"].(map[string]any)["project"].(map[string]any)["id"] = "project-foreign"
		}},
		{"configuration escaping storage", "configuration-invalid", func(m map[string]any, _ map[string][]byte) {
			m["project_config"].(map[string]any)["storage"].(map[string]any)["artifact_root"] = "../escape"
		}},
		{"configuration invalid limits", "configuration-invalid", func(m map[string]any, _ map[string][]byte) {
			m["project_config"].(map[string]any)["executor_limits"].(map[string]any)["max_output_bytes"] = 0
		}},
		{"unsupported projection", "schema-unsupported", func(m map[string]any, _ map[string][]byte) { m["projection_schema"] = 999 }},
		{"journal ordinal gap", "journal-invalid", func(m map[string]any, c map[string][]byte) {
			c["journal/00000002.jsonl"] = c["journal/00000001.jsonl"]
			delete(c, "journal/00000001.jsonl")
			renameManifestEntry(t, m, "journal/00000001.jsonl", "journal/00000002.jsonl")
		}},
		{"unknown entry family", "entry-extra", func(m map[string]any, c map[string][]byte) {
			c["project-prose.md"] = []byte("private prose")
			m["entries"] = append(m["entries"].([]any), map[string]any{"path": "project-prose.md", "digest": digest(c["project-prose.md"]), "size_bytes": len(c["project-prose.md"]), "media_type": "text/markdown"})
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			workspace := t.TempDir()
			if _, err := bootstrap.InitializeWorkspace(bootstrap.InitOptions{WorkspaceRoot: workspace, MigrationsDir: testMigrationsDir(t), Prefix: "FID"}); err != nil {
				t.Fatal(err)
			}
			archive := filepath.Join(t.TempDir(), "recovery.zip")
			if _, err := Create(CreateOptions{WorkspaceRoot: workspace, ArchivePath: archive}); err != nil {
				t.Fatal(err)
			}
			contents := archiveContents(t, archive)
			m := decodeManifestMap(t, contents["manifest.json"])
			tc.change(m, contents)
			contents["manifest.json"] = encodeManifestMap(t, m)
			writeTestArchive(t, archive, contents)
			before := workspaceTreeDigest(t, workspace)
			_, err := Validate(ValidateOptions{ArchivePath: archive, MigrationsDir: testMigrationsDir(t)})
			assertArchiveCheck(t, err, tc.code)
			if workspaceTreeDigest(t, workspace) != before {
				t.Fatal("invalid archive validation changed project")
			}
		})
	}
}
