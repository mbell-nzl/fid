package export

import (
	"bytes"
	"encoding/json"
	"path/filepath"
	"testing"

	"github.com/mbell-nzl/fid/internal/bootstrap"
	"github.com/mbell-nzl/fid/internal/journal"
)

func TestArchiveRejectsIntegrityValidSemanticCorruptions(t *testing.T) {
	t.Parallel()
	workspace := t.TempDir()
	if _, err := bootstrap.InitializeWorkspace(bootstrap.InitOptions{WorkspaceRoot: workspace, MigrationsDir: testMigrationsDir(t), Prefix: "FID"}); err != nil {
		t.Fatal(err)
	}
	addRecoveredExportFixtureFacts(t, workspace)
	archive := filepath.Join(t.TempDir(), "semantic.zip")
	if _, err := Create(CreateOptions{WorkspaceRoot: workspace, ArchivePath: archive}); err != nil {
		t.Fatal(err)
	}
	if _, err := Validate(ValidateOptions{ArchivePath: archive, MigrationsDir: testMigrationsDir(t)}); err != nil {
		t.Fatalf("valid control: %#v", err)
	}
	original := archiveContents(t, archive)
	for _, tc := range []struct {
		name, code, eventType string
		mutate                func(*journal.Event)
	}{
		{"unknown Event", "effect-invalid", "run.recovered", func(e *journal.Event) { e.EventType = "run.invented" }},
		{"incomplete nested recovery", "reference-invalid", "run.recovered", func(e *journal.Event) {
			for i := range e.Effects {
				if e.Effects[i].Table == "run_recoveries" {
					e.Effects[i].Values["reconciliation_json"] = `{}`
				}
			}
		}},
		{"incomplete cached recovery", "reference-invalid", "run.recovered", func(e *journal.Event) {
			for i := range e.Effects {
				if e.Effects[i].Table == "idempotency_results" {
					e.Effects[i].Values["result_json"] = `{"replayed":false}`
				}
			}
		}},
		{"unknown table", "effect-invalid", "run.recovered", func(e *journal.Event) { e.Effects[0].Table = "untrusted_table" }},
		{"unknown column", "effect-invalid", "run.recovered", func(e *journal.Event) { e.Effects[0].Values["untrusted_column"] = 1 }},
		{"wrong scalar type", "effect-invalid", "run.recovered", func(e *journal.Event) { e.Effects[0].Values["status"] = 1 }},
		{"boolean used as version", "effect-invalid", "run.recovered", func(e *journal.Event) {
			for i := range e.Effects {
				if e.Effects[i].Table == "run_runtime" {
					e.Effects[i].Values["version"] = true
					break
				}
			}
		}},
		{"executor identity rewritten", "effect-invalid", "run.recovered", func(e *journal.Event) {
			for i := range e.Effects {
				if e.Effects[i].Table == "run_runtime" && e.Effects[i].Op == "update" {
					e.Effects[i].Values["executor_id"] = "executor-forged"
					break
				}
			}
		}},
		{"unresolved idempotency Event", "reference-invalid", "run.recovered", func(e *journal.Event) {
			for i := range e.Effects {
				if e.Effects[i].Table == "idempotency_results" {
					e.Effects[i].Values["result_event_id"] = "EVT-absent"
				}
			}
		}},
		{"observation subject differs from its row", "reference-invalid", "run.observed", func(e *journal.Event) {
			for i := range e.Effects {
				if e.Effects[i].Table != "run_observations" {
					continue
				}
				var observation map[string]any
				if err := json.Unmarshal([]byte(e.Effects[i].Values["observation_json"].(string)), &observation); err != nil {
					t.Fatal(err)
				}
				observation["run_id"] = "RUN-forged"
				encoded, err := journal.MarshalCanonicalJSON(observation)
				if err != nil {
					t.Fatal(err)
				}
				e.Effects[i].Values["observation_json"] = string(encoded)
				break
			}
		}},
		{"foreign recovery aggregate", "reference-invalid", "run.recovered", func(e *journal.Event) { e.Aggregate.ID = "RUN-foreign"; e.Aggregate.Sequence = 1 }},
		{"issue aggregate differs from immutable Effect", "reference-invalid", "issue.created", func(e *journal.Event) { e.Aggregate.ID = "ISS-forged" }},
		{"immutable Artifact duplicate", "reference-invalid", "gate.receipt-recorded", func(e *journal.Event) { e.Effects = append([]journal.Effect{e.Effects[0]}, e.Effects...) }},
		{"retired caller outcome", "effect-invalid", "run.recovered", func(e *journal.Event) {
			e.EventType = "run.finished"
			var effects []journal.Effect
			for _, effect := range e.Effects {
				if (effect.Table == "runs" && effect.Op == "update") || effect.Table == "idempotency_results" {
					effects = append(effects, effect)
				}
			}
			e.Effects = effects
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			contents := make(map[string][]byte, len(original))
			for path, data := range original {
				contents[path] = bytes.Clone(data)
			}
			path := archivePathWithPrefix(t, contents, "journal/")
			data, hash := mutateArchiveEvent(t, contents[path], tc.eventType, tc.mutate)
			contents[path] = data
			m := decodeManifestMap(t, contents["manifest.json"])
			m["terminal_event_hash"] = hash
			updateManifestEntry(t, m, path, data)
			contents["manifest.json"] = encodeManifestMap(t, m)
			writeTestArchive(t, archive, contents)
			before := workspaceTreeDigest(t, workspace)
			_, err := Validate(ValidateOptions{ArchivePath: archive, MigrationsDir: testMigrationsDir(t)})
			assertArchiveCheck(t, err, tc.code)
			if workspaceTreeDigest(t, workspace) != before {
				t.Fatal("semantic rejection mutated project")
			}
		})
	}
}

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
