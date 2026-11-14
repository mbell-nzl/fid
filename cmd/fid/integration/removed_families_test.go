//go:build integration

package integration

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
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

func TestRemovedFamiliesRefuseHumanAndJSONWithoutMutation(t *testing.T) {
	t.Parallel()
	binary := integrationFIDBinary(t)
	raw, err := os.ReadFile("../../../internal/journal/testdata/removed-families.txt")
	if err != nil {
		t.Fatal(err)
	}
	for _, family := range strings.Fields(string(raw)) {
		t.Run(family, func(t *testing.T) {
			t.Parallel()
			source, target := gittest.NewWorkspace(t), gittest.NewWorkspace(t)
			p21Run(t, binary, source, "--json", "init", "--prefix", "FID")
			archiveDir, err := os.MkdirTemp("", "fid-schema-refusal-")
			if err != nil {
				t.Fatal(err)
			}
			t.Cleanup(func() { _ = os.RemoveAll(archiveDir) })
			archive := filepath.Join(archiveDir, "removed.zip")
			p21Run(t, binary, source, "--json", "export", "create", archive)
			entries := p22ArchiveEntries(t, archive)
			event := journal.Event{EventSchema: journal.EventSchemaV2, EventID: "EVT-removed", EventType: family, OccurredAt: time.Now().UTC(),
				ActorV2:   journal.EventActor{ActorID: "test", ActorType: "human", Producer: "fid", ProvenanceRef: "test:removed"},
				Aggregate: journal.Aggregate{Type: "issue", ID: "retired", Sequence: 1}, CommandName: "test.removed",
				Effects: []journal.Effect{{Op: "activity", Table: "activity_log", PK: map[string]any{"id": "removed"}}}}
			encoded, err := journal.MarshalCanonicalJSON(event)
			if err != nil {
				t.Fatal(err)
			}
			files, err := filepath.Glob(filepath.Join(config.LedgerDir(source), "journal", "*.jsonl"))
			if err != nil || len(files) != 1 {
				t.Fatalf("journal paths=%v/%v", files, err)
			}
			prior := readP21File(t, files[0])
			if err := os.WriteFile(files[0], append(prior, append(encoded, '\n')...), 0o600); err != nil {
				t.Fatal(err)
			}
			appendRemovedArchiveEvent(t, entries, event, encoded)
			p22WritePortableArchive(t, archive, entries)
			sourceBefore, targetBefore, archiveBefore := p23PortableTreeDigest(t, source), p23PortableTreeDigest(t, target), readP21File(t, archive)
			for _, args := range [][]string{{"status"}, {"rebuild"}, {"export", "create", archive}, {"export", "validate", archive}, {"export", "import", archive, "--repository", target, "--idempotency-key", "removed"}} {
				archiveOperation := args[0] == "export" && args[1] != "create"
				root, code := source, "journal-replay-stopped"
				if archiveOperation {
					root, code = target, "archive-invalid"
				}
				for _, jsonMode := range []bool{false, true} {
					command := append([]string{}, args...)
					if jsonMode {
						command = append([]string{"--json"}, command...)
					}
					stdout, stderr, err := closureRunCommand(root, binary, command...)
					if err == nil || processExitCode(t, err) != 6 || stdout != "" {
						t.Fatalf("%v = %v/%q/%q", command, err, stdout, stderr)
					}
					if jsonMode {
						response := p23Envelope(t, stderr)
						if p21PathString(t, response, "error", "code") != code {
							t.Fatalf("failure=%s", stderr)
						}
						if archiveOperation {
							checks := p21Path(t, response, "error", "details", "checks").([]any)
							check := checks[0].(map[string]any)
							if len(checks) != 1 || check["code"] != "schema-unsupported" || check["entry"] != "journal/00000001.jsonl" || !strings.Contains(check["detail"].(string), family) {
								t.Fatalf("archive checks=%v", checks)
							}
						} else if p21PathString(t, response, "error", "details", "reason") != "schema-unsupported" || p21PathString(t, response, "error", "details", "event_type") != family {
							t.Fatalf("journal refusal=%s", stderr)
						}
					} else {
						if !strings.Contains(stderr, "reason: "+code) {
							t.Fatalf("human refusal=%s", stderr)
						}
						for _, word := range []string{"claim", "runs"} {
							if strings.Contains(strings.ToLower(stderr), word) {
								t.Fatalf("kernel word in human refusal=%s", stderr)
							}
						}
					}
					if p23PortableTreeDigest(t, source) != sourceBefore || p23PortableTreeDigest(t, target) != targetBefore || !bytes.Equal(archiveBefore, readP21File(t, archive)) {
						t.Fatalf("%v refusal mutated files", command)
					}
				}
			}
		})
	}
}

func appendRemovedArchiveEvent(t *testing.T, entries map[string][]byte, event journal.Event, encoded []byte) {
	t.Helper()
	const path = "journal/00000001.jsonl"
	entries[path] = append(entries[path], append(encoded, '\n')...)
	var manifest map[string]any
	if err := json.Unmarshal(entries["manifest.json"], &manifest); err != nil {
		t.Fatal(err)
	}
	manifest["event_count"] = manifest["event_count"].(float64) + 1
	hash, err := journal.HashEvent(event)
	if err != nil {
		t.Fatal(err)
	}
	manifest["terminal_event_hash"] = hash
	for _, value := range manifest["entries"].([]any) {
		entry := value.(map[string]any)
		if entry["path"] == path {
			entry["digest"] = fmt.Sprintf("sha256:%x", sha256.Sum256(entries[path]))
			entry["size_bytes"] = len(entries[path])
		}
	}
	entries["manifest.json"], err = journal.MarshalCanonicalJSON(manifest)
	if err != nil {
		t.Fatal(err)
	}
}
