package journal

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/mbell-nzl/fid/internal/config"
)

func TestRemovedFamiliesRefusedWithoutMutation(t *testing.T) {
	t.Parallel()
	raw, err := os.ReadFile("testdata/removed-families.txt")
	if err != nil {
		t.Fatal(err)
	}
	for _, family := range strings.Fields(string(raw)) {
		for _, schema := range []int{EventSchemaV1, EventSchemaV2} {
			t.Run(fmt.Sprintf("%s/v%d", family, schema), func(t *testing.T) {
				t.Parallel()
				workspace := newReplayWorkspace(t)
				dir := filepath.Join(config.LedgerDir(workspace), "journal")
				path := filepath.Join(dir, "2026-11-15.jsonl")
				event := removedFamilyEvent(family, schema)
				// A valid prefix must not land when a later family is unsupported.
				prefix := removedFamilyEvent("project.initialized", EventSchemaV2)
				prefix.EventID, prefix.Aggregate.ID = "EVT-prefix", "PRJ-prefix"
				writeReplayEvents(t, path, prefix, event)
				beforeBytes, err := os.ReadFile(path)
				if err != nil {
					t.Fatal(err)
				}
				conn := openReplayDB(t, filepath.Join(config.LedgerDir(workspace), "fid.sqlite"))
				defer closeProjectDB(t, conn)
				before := currentProjectionSnapshot(t, conn)
				for _, mode := range []string{"project", "replay", "checkpointed-replay", "rebuild", "freshen"} {
					var err error
					switch mode {
					case "project":
						err = ProjectEvent(conn, ProjectionInput{Event: event, EventHash: strings.Repeat("a", 64), File: "2026-11-15.jsonl", Line: 1, Offset: 1})
					case "replay":
						_, err = ReplayJournal(conn, dir, nil)
					case "checkpointed-replay":
						_, err = ReplayJournal(conn, dir, &ReplayOptions{StartFile: filepath.Base(path), StartOffset: int64(len(beforeBytes))})
					case "rebuild":
						_, err = RebuildWorkspace(workspace, testProjectMigrationsDir(t))
					case "freshen":
						_, err = EnsureProjectionFresh(workspace, testProjectMigrationsDir(t))
					}
					if err == nil || !strings.Contains(err.Error(), "schema-unsupported") || !strings.Contains(err.Error(), family) {
						t.Fatalf("%s = %v, want schema-unsupported naming %s", mode, err, family)
					}
					if mode == "replay" {
						var corrupt *CorruptJournalError
						if !errors.As(err, &corrupt) || corrupt.Line != 2 || corrupt.File != filepath.ToSlash(path) {
							t.Fatalf("missing location: %v", err)
						}
					}
					if !reflect.DeepEqual(before, currentProjectionSnapshot(t, conn)) {
						t.Fatalf("%s changed projection", mode)
					}
					after, readErr := os.ReadFile(path)
					if readErr != nil || string(after) != string(beforeBytes) {
						t.Fatalf("%s changed journal: %v", mode, readErr)
					}
				}
			})
		}
	}
}

func removedFamilyEvent(family string, schema int) Event {
	event := Event{EventID: "EVT-removed", EventSchema: schema, EventType: family,
		OccurredAt: time.Date(2026, 11, 15, 0, 0, 0, 0, time.UTC),
		Actor:      Actor{Role: "test"}, Command: Command{Argv: []string{"fid", "test"}, CWD: "/test", FIDVersion: "test"},
		ActorV2:   EventActor{ActorID: "test", ActorType: "human", Producer: "fid", ProvenanceRef: "test:removed"},
		Aggregate: Aggregate{ID: "retired-subject", Type: "issue", Sequence: 1}, CommandName: "test.removed", Effects: []Effect{}}
	if family != "project.initialized" {
		event.Effects = []Effect{{Op: "activity", Table: "activity_log", PK: map[string]any{"id": "retired"}}}
	}
	return event
}

func TestRetiredEffectsCannotHideInLegacyIssueHistory(t *testing.T) {
	t.Parallel()
	for _, schema := range []int{EventSchemaV1, EventSchemaV2} {
		t.Run(fmt.Sprint(schema), func(t *testing.T) {
			conn := openReplayDB(t, filepath.Join(t.TempDir(), "projection.sqlite"))
			defer closeProjectDB(t, conn)
			event := removedFamilyEvent("issue.updated", schema)
			event.Effects = []Effect{{Op: "insert", Table: "concerns", PK: map[string]any{"id": "retired"}}}
			before := currentProjectionSnapshot(t, conn)
			err := ProjectEvent(conn, ProjectionInput{Event: event, EventHash: strings.Repeat("a", 64), File: "2026-11-15.jsonl", Line: 1, Offset: 1})
			if err == nil || !strings.Contains(err.Error(), "schema-unsupported") || !strings.Contains(err.Error(), "concerns") {
				t.Fatalf("project = %v", err)
			}
			if !reflect.DeepEqual(before, currentProjectionSnapshot(t, conn)) {
				t.Fatal("projection changed")
			}
		})
	}
}
