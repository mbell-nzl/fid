package journal

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/mbell-nzl/fid/internal/config"
)

func TestRetiredIdentifiersRefusedBeforeMutation(t *testing.T) {
	t.Parallel()
	for _, schema := range []int{EventSchemaV1, EventSchemaV2} {
		for _, name := range []string{"workflows", "WORKFLOWS", "WoRkFlOwS", "CONCERNS", "RUN.STARTED", "Run.Started"} {
			t.Run(fmt.Sprintf("%s/v%d", name, schema), func(t *testing.T) {
				t.Parallel()
				event := removedFamilyEvent("issue.updated", schema)
				table := name
				if strings.Contains(name, ".") {
					event.EventType, table = name, "labels"
				}
				event.Effects = []Effect{{Op: "insert", Table: table, PK: map[string]any{"id": "removed"}, Values: map[string]any{"name": "removed"}}}
				if table == "labels" {
					event.Effects[0].PK = map[string]any{"name": "removed"}
				}
				assertIdentifierRefusalUnchanged(t, event, "schema-unsupported")
			})
		}
	}
}

func TestEffectTableSyntaxRefusedBeforeMutation(t *testing.T) {
	t.Parallel()
	// Effect identifiers are bare ASCII names. Quoting, qualification and SQL
	// whitespace are rejected by the same grammar that protects SQL execution.
	for _, table := range []string{`"WORKFLOWS"`, "`WoRkFlOwS`", "[WORKFLOWS]", "'WORKFLOWS'", "main.WORKFLOWS", " WORKFLOWS", "WORKFLOWS "} {
		for _, schema := range []int{EventSchemaV1, EventSchemaV2} {
			t.Run(fmt.Sprintf("%s/v%d", table, schema), func(t *testing.T) {
				t.Parallel()
				event := removedFamilyEvent("issue.updated", schema)
				event.Effects = []Effect{{Op: "insert", Table: table, PK: map[string]any{"id": "removed"}, Values: map[string]any{"name": "removed"}}}
				assertIdentifierRefusalUnchanged(t, event, "invalid SQL identifier")
			})
		}
	}
}

func assertIdentifierRefusalUnchanged(t *testing.T, event Event, reason string) {
	t.Helper()
	// This prefix writes a real label, so partial replay is observable even
	// without inspecting the projection checkpoint.
	prefix := replayRoleEvent("EVT-prefix", "supported-prefix")
	assertReplayPrefix(t, prefix)
	for _, mode := range []string{"project", "dry-run", "replay", "checkpointed-replay", "rebuild", "freshen"} {
		t.Run(mode, func(t *testing.T) {
			workspace := newReplayWorkspace(t)
			dir := filepath.Join(config.LedgerDir(workspace), "journal")
			path := filepath.Join(dir, "2026-11-15.jsonl")
			writeReplayEvents(t, path, prefix, event)
			journalBefore, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			conn := openReplayDB(t, filepath.Join(config.LedgerDir(workspace), "fid.sqlite"))
			defer closeProjectDB(t, conn)
			before := currentProjectionSnapshot(t, conn)
			switch mode {
			case "project":
				err = ProjectEvent(conn, ProjectionInput{Event: event, EventHash: strings.Repeat("a", 64), File: filepath.Base(path), Line: 2, Offset: int64(len(journalBefore))})
			case "dry-run":
				err = DryRunProjectEvent(conn, event)
			case "replay":
				_, err = ReplayJournal(conn, dir, nil)
			case "checkpointed-replay":
				_, err = ReplayJournal(conn, dir, &ReplayOptions{StartFile: filepath.Base(path), StartOffset: int64(len(journalBefore))})
			case "rebuild":
				_, err = RebuildWorkspace(workspace, testProjectMigrationsDir(t))
			case "freshen":
				_, err = EnsureProjectionFresh(workspace, testProjectMigrationsDir(t))
			}
			if err == nil || !strings.Contains(err.Error(), reason) {
				t.Errorf("%s = %v, want %s", mode, err, reason)
			}
			if reason == "schema-unsupported" {
				var unsupported *SchemaUnsupportedError
				if !errors.As(err, &unsupported) || unsupported.EventID != event.EventID || unsupported.EventType != event.EventType {
					t.Errorf("missing schema refusal identity: %v", err)
				}
			}
			if !reflect.DeepEqual(before, currentProjectionSnapshot(t, conn)) {
				t.Errorf("%s changed projection", mode)
			}
			after, readErr := os.ReadFile(path)
			if readErr != nil || !reflect.DeepEqual(journalBefore, after) {
				t.Errorf("%s changed journal: %v", mode, readErr)
			}
		})
	}
}
