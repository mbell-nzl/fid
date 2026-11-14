package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/mbell-nzl/fid/internal/cli"
	"github.com/mbell-nzl/fid/internal/journal"
)

func TestRemovedFamiliesRenderSchemaRefusal(t *testing.T) {
	t.Parallel()
	raw, err := os.ReadFile("../../internal/journal/testdata/removed-families.txt")
	if err != nil {
		t.Fatal(err)
	}
	for _, family := range strings.Fields(string(raw)) {
		t.Run(family, func(t *testing.T) {
			t.Parallel()
			fixture := newJournalStatusFixture(t)
			event := journal.Event{EventSchema: journal.EventSchemaV1, EventID: "EVT-removed", EventType: family, OccurredAt: time.Now().UTC(), Effects: []journal.Effect{}}
			encoded, err := journal.MarshalCanonicalJSON(event)
			if err != nil {
				t.Fatal(err)
			}
			path := mustJournalPath(t, fixture.workspace)
			appendRawJournalLine(t, path, string(encoded))
			before := snapshotWorkspacePersistence(t, fixture.workspace)
			for _, command := range []string{"status", "rebuild"} {
				result := fixture.run("--json", command)
				if cli.ExitCode(result.err) != 6 || result.stdout != "" {
					t.Fatalf("%s = %#v", command, result)
				}
				var response struct {
					Error struct {
						Code    string
						Details map[string]any
					}
				}
				if err := json.Unmarshal([]byte(result.stderr), &response); err != nil {
					t.Fatal(err)
				}
				details := response.Error.Details
				if response.Error.Code != "journal-replay-stopped" || details["reason"] != "schema-unsupported" || details["event_type"] != family || details["event_id"] != "EVT-removed" || details["journal_file"] != filepath.ToSlash(path) {
					t.Fatalf("%s failure = %s", command, result.stderr)
				}
				text := fixture.run(command)
				if cli.ExitCode(text.err) != 6 || text.stdout != "" || !strings.Contains(text.stderr, "schema-unsupported") {
					t.Fatalf("%s human failure = %#v", command, text)
				}
				for _, word := range []string{"claim", "runs"} {
					if strings.Contains(strings.ToLower(text.stderr), word) {
						t.Fatalf("kernel word in human refusal: %s", text.stderr)
					}
				}
				if !before.equal(snapshotWorkspacePersistence(t, fixture.workspace)) {
					t.Fatal("refusal changed persistence")
				}
			}
		})
	}
}
