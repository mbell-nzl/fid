package export

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/mbell-nzl/fid/internal/journal"
)

func TestArchiveWithValidHashesButInvalidCheckReferenceRejects(t *testing.T) {
	t.Parallel()
	fixture := newExportSourceFixture(t, addLeanArtifactExportFacts)
	if _, err := Validate(ValidateOptions{ArchivePath: fixture.archive, MigrationsDir: testMigrationsDir(t)}); err != nil {
		t.Fatalf("valid lean-history control: %v", err)
	}
	entries := archiveContents(t, fixture.archive)
	path := archivePathWithPrefix(t, entries, "journal/")
	var hash string
	entries[path], hash = mutateArchiveEvent(t, entries[path], "check.artifact-kept", func(e *journal.Event) { effectForTest(t, e, "check_files").Values["artifact_id"] = "ART-missing" })
	manifest := decodeManifestMap(t, entries["manifest.json"])
	manifest["terminal_event_hash"] = hash
	updateManifestEntry(t, manifest, path, entries[path])
	entries["manifest.json"] = encodeManifestMap(t, manifest)
	archive := filepath.Join(t.TempDir(), "missing-reference.zip")
	writeTestArchive(t, archive, entries)
	before := workspaceTreeDigest(t, fixture.workspace)
	_, err := Validate(ValidateOptions{ArchivePath: archive, MigrationsDir: testMigrationsDir(t)})
	assertArchiveCheck(t, err, "reference-invalid")
	if workspaceTreeDigest(t, fixture.workspace) != before {
		t.Fatal("reference rejection mutated source")
	}
}

func TestSemanticValidatorRejectsMissingHoldReference(t *testing.T) {
	t.Parallel()
	validator, err := newSemanticValidator(testMigrationsDir(t), 26)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = validator.close() }()
	raw, err := journal.MarshalCanonicalJSON(archiveHoldEvent(func(*journal.Event) {}))
	if err != nil {
		t.Fatal(err)
	}
	err = validator.accept(raw)
	if err == nil || !strings.Contains(err.Error(), "FOREIGN KEY constraint failed") || semanticCheckCode(err) != "reference-invalid" {
		t.Fatalf("missing hold reference=%v", err)
	}
}

func TestSemanticValidatorRejectsHoldEffectFamilyTampering(t *testing.T) {
	t.Parallel()
	for _, row := range []struct {
		name   string
		mutate func(*journal.Event)
	}{
		{"wrong operation", func(e *journal.Event) { e.Effects[0].Op = "update" }},
		{"reordered replacement hold", func(e *journal.Event) {
			e.Effects = append(e.Effects, journal.Effect{Op: "update", Table: "holds", PK: map[string]any{"id": "HLD-old"}, Values: map[string]any{"ended_at": "2026-11-10T09:00:00Z", "end_reason": "replaced"}, Preconditions: map[string]any{"ended_at": nil}})
		}},
		{"wrong aggregate", func(e *journal.Event) { e.Aggregate.Type = "issue" }},
	} {
		t.Run(row.name, func(t *testing.T) {
			t.Parallel()
			validator, err := newSemanticValidator(testMigrationsDir(t), 26)
			if err != nil {
				t.Fatal(err)
			}
			defer func() { _ = validator.close() }()
			raw, err := journal.MarshalCanonicalJSON(archiveHoldEvent(row.mutate))
			if err != nil {
				t.Fatal(err)
			}
			err = validator.accept(raw)
			if err == nil || semanticCheckCode(err) != "effect-invalid" {
				t.Fatalf("tampered hold Event=%v", err)
			}
		})
	}
}

func TestSemanticValidatorRejectsNestedDuplicateJSONKey(t *testing.T) {
	t.Parallel()
	validator, err := newSemanticValidator(testMigrationsDir(t), 17)
	if err != nil {
		t.Fatalf("new semantic validator: %v", err)
	}
	defer func() { _ = validator.close() }()

	raw := []byte(`{"actor":{"actor_id":"agent:one","actor_id":"agent:two","actor_type":"agent","producer":"fid","provenance_ref":"test:duplicate","session_ref":null},"aggregate":{"id":"PROJ-1","sequence":1,"type":"project"},"causation_id":null,"command":"init","correlation_id":null,"effects":[],"event_id":"EVT-duplicate","event_schema":2,"event_type":"project.initialized","idempotency":null,"occurred_at":"2026-09-08T00:00:00Z","repository":null}`)
	if err := validator.accept(raw); err == nil || !strings.Contains(err.Error(), `duplicate key "actor_id"`) {
		t.Fatalf("accept nested duplicate key error = %v, want duplicate actor key rejection", err)
	}
}
