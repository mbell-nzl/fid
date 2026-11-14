//go:build integration

package integration

import (
	"os"
	"reflect"
	"strings"
	"testing"
)

func TestP24ExecutedKeyedCorpusMatches51LeafInventory(t *testing.T) {
	t.Parallel()
	rows := map[string]string{}
	add := func(operation, field string) {
		t.Helper()
		name := strings.ReplaceAll(operation, ".", " ")
		if operation == "hold.park" {
			name = "park"
		}
		if operation == "task.takeover" {
			name = "takeover"
		}
		if prior, found := rows[name]; found {
			t.Fatalf("duplicate executable row %s: %s and %s", name, prior, field)
		}
		rows[name] = field
	}
	for _, group := range []map[string]string{p24LifecycleFields, p24FoundationFields, p24RecoveryFields, p24FlagFields, p24ActivationRecoveryFields, p24ActivationAbortFields, p24GateEvidenceFields, p24ArtifactFields, p24WorkflowFields} {
		for operation, field := range group {
			add(operation, field)
		}
	}
	for operation, names := range p24JudgmentFields {
		add(operation, names[0])
	}
	data, err := os.ReadFile("testdata/p24_keyed_leaves.tsv")
	if err != nil {
		t.Fatal(err)
	}
	planned := map[string]string{}
	for _, line := range strings.Split(strings.TrimSpace(string(data)), "\n") {
		parts := strings.Split(line, "\t")
		if len(parts) != 2 {
			t.Fatalf("invalid inventory row %q", line)
		}
		planned[parts[0]] = parts[1]
	}
	if len(rows) != 51 || !reflect.DeepEqual(rows, planned) {
		t.Fatalf("executable rows = %v; planned inventory = %v", rows, planned)
	}
}
