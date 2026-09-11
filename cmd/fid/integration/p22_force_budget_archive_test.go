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

	"github.com/mbell-nzl/fid/internal/journal"
)

func TestP22ForceRecoveryBudgetRejectsIntegrityValidMissingChargeArchive(t *testing.T) {
	a := p22ForceBudgetFixture(t, 1)
	p22AcceptedCall(t, a, "archive-force", "run.recover", p22ForceBudgetInput(t, a, "archive-force"))
	archive := filepath.Join(t.TempDir(), "valid.zip")
	p21Run(t, a.binary, a.workspace, "--json", "export", "create", archive)
	p21Run(t, a.binary, a.workspace, "--json", "export", "validate", archive)
	entries := p22ArchiveEntries(t, archive)
	changed := 0
	for name, data := range entries {
		if !strings.HasPrefix(name, "journal/") {
			continue
		}
		lines := bytes.Split(bytes.TrimSuffix(data, []byte{'\n'}), []byte{'\n'})
		for i, line := range lines {
			var event journal.Event
			p22FaultOK(t, json.Unmarshal(line, &event))
			if event.EventType != "run.recovered" {
				continue
			}
			p22RemoveForceEffects(&event, "budget_samples", "stops", "attention_items", "attention_stop_links")
			p22ForceCachedBudget(&event, 0, 1, "available")
			encoded, err := json.Marshal(event)
			p22FaultOK(t, err)
			lines[i] = encoded
			changed++
		}
		entries[name] = append(bytes.Join(lines, []byte{'\n'}), '\n')
	}
	if changed != 1 {
		t.Fatalf("changed %d force Events, want 1", changed)
	}
	var manifest map[string]any
	p22FaultOK(t, json.Unmarshal(entries["manifest.json"], &manifest))
	for _, raw := range manifest["entries"].([]any) {
		entry := raw.(map[string]any)
		data := entries[entry["path"].(string)]
		entry["size_bytes"] = len(data)
		entry["digest"] = fmt.Sprintf("sha256:%x", sha256.Sum256(data))
	}
	encoded, err := json.Marshal(manifest)
	p22FaultOK(t, err)
	entries["manifest.json"] = encoded
	malformed := filepath.Join(t.TempDir(), "invalid.zip")
	p22WritePortableArchive(t, malformed, entries)
	for _, operation := range []string{"validate", "import"} {
		t.Run(operation, func(t *testing.T) {
			target := filepath.Join(t.TempDir(), "destination")
			closureRun(t, filepath.Dir(target), "git", "clone", "-q", a.workspace, target)
			args := []string{"--json", "export", operation, malformed}
			if operation == "import" {
				args = append(args, "--repository", target, "--idempotency-key", "invalid-force")
			}
			stdout, stderr, err := closureRunCommand(target, a.binary, args...)
			if err == nil {
				t.Fatal("accepted missing force accounting")
			}
			response := p22ConformanceResponse{Stdout: stdout, Stderr: stderr, ExitCode: processExitCode(t, err)}
			assertP22ConformanceResponse(t, p22ConformanceCase{Operation: "export." + operation, Kind: "invalid", ExitCode: 6, FailureCode: "archive-invalid"}, response)
			if _, err := os.Stat(filepath.Join(target, ".fid")); !os.IsNotExist(err) {
				t.Fatalf("rejected archive installed target state: %v", err)
			}
		})
	}
}
