package export

import (
	"bytes"
	"encoding/json"
	"path/filepath"
	"testing"

	"github.com/mbell-nzl/fid/internal/journal"
)

func mutateArchiveEvent(t *testing.T, data []byte, eventType string, mutate func(*journal.Event)) ([]byte, string) {
	t.Helper()
	lines := bytes.Split(bytes.TrimSuffix(data, []byte{'\n'}), []byte{'\n'})
	changed := false
	var terminal journal.Event
	for i, line := range lines {
		var event journal.Event
		if err := event.UnmarshalJSON(line); err != nil {
			t.Fatal(err)
		}
		if event.EventType == eventType && !changed {
			mutate(&event)
			encoded, err := journal.MarshalCanonicalJSON(event)
			if err != nil {
				t.Fatal(err)
			}
			lines[i] = encoded
			changed = true
		}
		terminal = event
	}
	if !changed {
		t.Fatalf("no %s Event", eventType)
	}
	hash, err := journal.HashEvent(terminal)
	if err != nil {
		t.Fatal(err)
	}
	return append(bytes.Join(lines, []byte{'\n'}), '\n'), hash
}

func TestArchiveRejectsIntegrityValidSemanticCorruptions(t *testing.T) {
	t.Parallel()
	fixture := newExportSourceFixture(t, addLeanArtifactExportFacts)
	if _, err := Validate(ValidateOptions{ArchivePath: fixture.archive, MigrationsDir: testMigrationsDir(t)}); err != nil {
		t.Fatalf("valid control: %v", err)
	}
	original := archiveContents(t, fixture.archive)
	for _, row := range []struct {
		name, code, family string
		mutate             func(*journal.Event)
	}{
		{"unknown Event", "effect-invalid", "check.started", func(e *journal.Event) { e.EventType = "check.invented" }},
		{"incomplete nested check", "reference-invalid", "check.started", func(e *journal.Event) { effectForTest(t, e, "check_starts").Values["start_json"] = `{}` }},
		{"incomplete cached result", "reference-invalid", "check.result-recorded", func(e *journal.Event) {
			effectForTest(t, e, "idempotency_results").Values["result_json"] = `{"replayed":false}`
		}},
		{"unknown table", "effect-invalid", "check.started", func(e *journal.Event) { e.Effects[0].Table = "untrusted_table" }},
		{"unknown column", "effect-invalid", "check.started", func(e *journal.Event) { e.Effects[0].Values["untrusted_column"] = 1 }},
		{"wrong scalar type", "effect-invalid", "check.started", func(e *journal.Event) { effectForTest(t, e, "check_requests").Values["actor"] = 1 }},
		{"boolean used as number", "effect-invalid", "check.artifact-kept", func(e *journal.Event) { effectForTest(t, e, "referenced_artifacts").Values["size_bytes"] = true }},
		{"unresolved idempotency Event", "reference-invalid", "check.started", func(e *journal.Event) {
			effectForTest(t, e, "idempotency_results").Values["result_event_id"] = "EVT-absent"
		}},
		{"observation subject differs from its row", "reference-invalid", "check.started", func(e *journal.Event) {
			effect := effectForTest(t, e, "check_starts")
			var start map[string]any
			if err := json.Unmarshal([]byte(effect.Values["start_json"].(string)), &start); err != nil {
				t.Fatal(err)
			}
			start["task_id"] = "FID-forged"
			effect.Values["start_json"] = encodedLean(start)
		}},
		{"foreign check aggregate", "reference-invalid", "check.result-recorded", func(e *journal.Event) { e.Aggregate.ID = "CHK-foreign"; e.Aggregate.Sequence = 1 }},
		{"issue aggregate differs from immutable Effect", "reference-invalid", "issue.created", func(e *journal.Event) { e.Aggregate.ID = "FID-forged" }},
		{"immutable Artifact duplicate", "reference-invalid", "check.artifact-kept", func(e *journal.Event) { e.Effects = append([]journal.Effect{e.Effects[0]}, e.Effects...) }},
	} {
		t.Run(row.name, func(t *testing.T) {
			t.Parallel()
			entries := make(map[string][]byte, len(original))
			for path, data := range original {
				entries[path] = bytes.Clone(data)
			}
			path := archivePathWithPrefix(t, entries, "journal/")
			data, hash := mutateArchiveEvent(t, entries[path], row.family, row.mutate)
			entries[path] = data
			manifest := decodeManifestMap(t, entries["manifest.json"])
			manifest["terminal_event_hash"] = hash
			updateManifestEntry(t, manifest, path, data)
			entries["manifest.json"] = encodeManifestMap(t, manifest)
			archive := filepath.Join(t.TempDir(), "corrupt.zip")
			writeTestArchive(t, archive, entries)
			before := workspaceTreeDigest(t, fixture.workspace)
			_, err := Validate(ValidateOptions{ArchivePath: archive, MigrationsDir: testMigrationsDir(t)})
			assertArchiveCheck(t, err, row.code)
			if workspaceTreeDigest(t, fixture.workspace) != before {
				t.Fatal("semantic rejection mutated project")
			}
		})
	}
}

func effectForTest(t *testing.T, event *journal.Event, table string) *journal.Effect {
	t.Helper()
	for i := range event.Effects {
		if event.Effects[i].Table == table {
			return &event.Effects[i]
		}
	}
	t.Fatalf("missing %s Effect in %s", table, event.EventType)
	return nil
}
