package export

import (
	"bytes"
	"path/filepath"
	"testing"

	"github.com/mbell-nzl/fid/internal/bootstrap"
	"github.com/mbell-nzl/fid/internal/gittest"
	"github.com/mbell-nzl/fid/internal/journal"
)

func TestKernelArchiveRejectsRecoveryAndRuntimeCorruptions(t *testing.T) {
	t.Parallel()
	workspace := gittest.NewWorkspace(t)
	if _, err := bootstrap.InitializeWorkspace(bootstrap.InitOptions{WorkspaceRoot: workspace, MigrationsDir: testMigrationsDir(t), Prefix: "FID"}); err != nil {
		t.Fatal(err)
	}
	addRecoveredKernelExportFacts(t, workspace)
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
