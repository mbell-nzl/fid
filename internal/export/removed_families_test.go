package export

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/mbell-nzl/fid/internal/config"
	"github.com/mbell-nzl/fid/internal/gittest"
	"github.com/mbell-nzl/fid/internal/journal"
)

func TestRemovedFamiliesArchiveRefusalWithoutMutation(t *testing.T) {
	t.Parallel()
	fixture := newExportSourceFixture(t, nil)
	original := archiveContents(t, fixture.archive)
	raw, err := os.ReadFile("../journal/testdata/removed-families.txt")
	if err != nil {
		t.Fatal(err)
	}
	for _, family := range strings.Fields(string(raw)) {
		for _, schema := range []int{journal.EventSchemaV1, journal.EventSchemaV2} {
			for _, retiredProfile := range []bool{false, true} {
				t.Run(fmt.Sprintf("%s/v%d/retired-profile=%t", family, schema, retiredProfile), func(t *testing.T) {
					t.Parallel()
					entries := make(map[string][]byte, len(original))
					for path, data := range original {
						entries[path] = bytes.Clone(data)
					}
					path := archivePathWithPrefix(t, entries, "journal/")
					event := journal.Event{EventID: "EVT-retired", EventSchema: schema, EventType: family,
						OccurredAt: time.Date(2026, 11, 15, 0, 0, 0, 0, time.UTC),
						ActorV2:    journal.EventActor{ActorID: "test", ActorType: "human", Producer: "fid", ProvenanceRef: "test:removed"},
						Aggregate:  journal.Aggregate{ID: "retired", Type: "issue", Sequence: 1}, CommandName: "test.removed", Effects: []journal.Effect{{Op: "activity", Table: "activity_log", PK: map[string]any{"id": "retired"}}}}
					encoded, err := journal.MarshalCanonicalJSON(event)
					if err != nil {
						t.Fatal(err)
					}
					entries[path] = append(entries[path], append(encoded, '\n')...)
					manifest := decodeManifestMap(t, entries["manifest.json"])
					manifest["terminal_event_hash"], err = journal.HashEvent(event)
					if err != nil {
						t.Fatal(err)
					}
					manifest["event_count"] = json.Number(fmt.Sprint(len(bytes.Split(bytes.TrimSpace(entries[path]), []byte{'\n'}))))
					updateManifestEntry(t, manifest, path, entries[path])
					if retiredProfile {
						manifest["capabilities"] = retiredV1CapabilityFixture()
						manifest["sources"] = []any{map[string]any{"source_type": "external-reference", "external_id": "retired"}}
					}
					entries["manifest.json"] = encodeManifestMap(t, manifest)
					archive := filepath.Join(t.TempDir(), "removed.zip")
					writeTestArchive(t, archive, entries)
					target := gittest.NewWorkspace(t)
					sourceBefore, targetBefore := workspaceTreeDigest(t, fixture.workspace), workspaceTreeDigest(t, target)
					archiveBefore := workspaceTreeDigest(t, filepath.Dir(archive))
					for _, operation := range []string{"validate", "import"} {
						var err error
						if operation == "validate" {
							_, err = Validate(ValidateOptions{ArchivePath: archive, MigrationsDir: testMigrationsDir(t)})
						} else {
							_, err = Import(ImportOptions{ArchivePath: archive, RepositoryRoot: target, MigrationsDir: testMigrationsDir(t), IdempotencyKey: "removed"})
						}
						var invalid *InvalidArchiveError
						if !errors.As(err, &invalid) || len(invalid.Checks) != 1 || invalid.Checks[0].Code != "schema-unsupported" || invalid.Checks[0].Entry == nil || *invalid.Checks[0].Entry != path || !strings.Contains(invalid.Checks[0].Detail, family) {
							t.Fatalf("%s = %#v, want schema-unsupported naming %s at %s", operation, err, family, path)
						}
						if workspaceTreeDigest(t, fixture.workspace) != sourceBefore || workspaceTreeDigest(t, target) != targetBefore || workspaceTreeDigest(t, filepath.Dir(archive)) != archiveBefore {
							t.Fatalf("%s refusal mutated files", operation)
						}
					}
				})
			}
		}
	}
}

func TestCreateRefusesRemovedJournalBeforePublication(t *testing.T) {
	t.Parallel()
	fixture := newExportSourceFixture(t, nil)
	files, err := filepath.Glob(filepath.Join(config.LedgerDir(fixture.workspace), "journal", "*.jsonl"))
	if err != nil || len(files) != 1 {
		t.Fatalf("journal=%v/%v", files, err)
	}
	event := journal.Event{EventSchema: journal.EventSchemaV2, EventID: "EVT-removed", EventType: "run.started", OccurredAt: time.Now().UTC(), ActorV2: journal.EventActor{ActorID: "test", ActorType: "human", Producer: "fid", ProvenanceRef: "test:removed"}, Aggregate: journal.Aggregate{Type: "run", ID: "RUN-old", Sequence: 1}, CommandName: "run.start", Effects: []journal.Effect{{Op: "insert", Table: "runs", PK: map[string]any{"id": "RUN-old"}}}}
	encoded, err := journal.MarshalCanonicalJSON(event)
	if err != nil {
		t.Fatal(err)
	}
	body, err := os.ReadFile(files[0])
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(files[0], append(body, append(encoded, '\n')...), 0o600); err != nil {
		t.Fatal(err)
	}
	sourceBefore, archiveBefore := workspaceTreeDigest(t, fixture.workspace), workspaceTreeDigest(t, filepath.Dir(fixture.archive))
	_, err = Create(CreateOptions{WorkspaceRoot: fixture.workspace, ArchivePath: fixture.archive})
	var unsupported *journal.SchemaUnsupportedError
	var corrupt *journal.CorruptJournalError
	if !errors.As(err, &unsupported) || unsupported.EventType != "run.started" || !errors.As(err, &corrupt) || corrupt.Line != 2 || corrupt.File != filepath.ToSlash(files[0]) {
		t.Fatalf("create=%v, want located schema-unsupported", err)
	}
	if workspaceTreeDigest(t, fixture.workspace) != sourceBefore || workspaceTreeDigest(t, filepath.Dir(fixture.archive)) != archiveBefore {
		t.Fatal("refused create changed source or published archive")
	}
}
