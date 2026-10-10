//go:build integration

package integration

import (
	"encoding/json"
	"path/filepath"
	"reflect"
	"testing"
)

func TestP24TakeControlRebuildExportImportParity(t *testing.T) {
	t.Parallel()
	start := p24NewTakeControlStart(t, integrationFIDBinary(t), t.TempDir())
	// Setup seeds the two additional detail-read subjects. These IDs and the
	// setup's JSON never supply a value to the Human journey below.
	for _, row := range []struct {
		family, id string
		input      map[string]any
	}{
		{"concern", "CON-parity", map[string]any{"id": "CON-parity", "subject": map[string]any{"type": "issue", "id": start.issueID}, "scope": map[string]any{"paths": []string{}}, "severity": "minor", "summary": "Parity concern", "reasoning": "Retain recorded reasoning", "references": []any{}}},
		{"deviation", "DEV-parity", map[string]any{"id": "DEV-parity", "subject": map[string]any{"type": "issue", "id": start.issueID}, "expected": "Original plan", "observed": "Human approval", "authority": map[string]any{"actor_id": "human", "reference": nil}, "disposition": "pending-decision", "reasoning": "Retain approved context", "references": []any{}, "supersedes_deviation_id": nil}},
	} {
		p21Run(t, start.binary, start.workspace, "--json", row.family, "record", "--actor", "human", "--input", p21Input(t, t.TempDir(), row.id, row.input), "--idempotency-key", row.id)
	}
	p24CompleteTakeControl(t, start)
	// Parity inspection begins after the Human journey has completed. JSON
	// and SQL here are never journey provenance sources.
	live := p23PortableRows(t, start.workspace)
	maps := p24ProjectedFieldMaps(t, start.workspace)
	reads := p24ParityDetailReads(t, start.binary, start.workspace)
	archive := filepath.Join(t.TempDir(), "phase24.zip")
	p21Run(t, start.binary, start.workspace, "--json", "export", "create", archive)
	entries := p22ArchiveEntries(t, archive)
	var manifest map[string]any
	if err := json.Unmarshal(entries["manifest.json"], &manifest); err != nil {
		t.Fatal(err)
	}
	if manifest["projection_schema"] != float64(19) {
		t.Fatalf("manifest schema = %v", manifest["projection_schema"])
	}
	p21Run(t, start.binary, start.workspace, "--json", "export", "validate", archive)
	p21Run(t, start.binary, start.workspace, "--json", "rebuild")
	p24AssertJourneyParity(t, "rebuilt", start.binary, start.workspace, live, maps, reads)
	target := filepath.Join(t.TempDir(), "imported")
	closureRun(t, filepath.Dir(target), "git", "clone", "-q", start.workspace, target)
	imported := p21Run(t, start.binary, target, "--json", "export", "import", archive, "--repository", target, "--idempotency-key", "phase24-import")
	if p21Path(t, imported, "result", "import", "projection_schema") != float64(19) {
		t.Fatal("import schema differs")
	}
	p24AssertJourneyParity(t, "imported", start.binary, target, live, maps, reads)
}

func p24ParityDetailReads(t *testing.T, binary, workspace string) map[string]any {
	t.Helper()
	result := map[string]any{}
	for family, id := range map[string]string{"concern": "CON-parity", "deviation": "DEV-parity", "human-outcome": "OUT-take-control"} {
		response := p21Run(t, binary, workspace, "--json", family, "get", id)
		result[family] = response["result"]
	}
	return result
}

func p24AssertJourneyParity(t *testing.T, label, binary, workspace string, live map[string][]string, maps map[string]string, reads map[string]any) {
	t.Helper()
	p23PortableAssertRows(t, label, live, p23PortableRows(t, workspace))
	p23AssertRebuiltProjectionIntegrity(t, workspace)
	if !reflect.DeepEqual(maps, p24ProjectedFieldMaps(t, workspace)) {
		t.Fatalf("%s field-digest maps differ", label)
	}
	if !reflect.DeepEqual(reads, p24ParityDetailReads(t, binary, workspace)) {
		t.Fatalf("%s new detail reads differ", label)
	}
}
