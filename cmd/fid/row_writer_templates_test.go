package main

import (
	"encoding/json"
	"os"
	"reflect"
	"testing"

	"github.com/mbell-nzl/fid/internal/gittest"
	"github.com/mbell-nzl/fid/internal/journal"
)

const rowWriterTemplatePath = "../../internal/journal/testdata/row-writer-templates.json"

type rowWriterTemplates struct {
	First, Second string
	Templates     map[string]journal.Event
}

// TestRowWriterTemplatesMatchCommands always recaptures in a guarded temporary
// repository. Comparing table/op sequences detects writer drift without pinning
// randomly allocated IDs or timestamps. Regenerate explicitly with:
// FID_REGENERATE_ROW_TEMPLATES=1 go test ./cmd/fid -run '^TestRowWriterTemplatesMatchCommands$' -count=1
func TestRowWriterTemplatesMatchCommands(t *testing.T) {
	t.Parallel()
	fresh := captureRowWriterTemplates(t)
	if os.Getenv("FID_REGENERATE_ROW_TEMPLATES") == "1" {
		raw, err := json.MarshalIndent(fresh, "", "  ")
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(rowWriterTemplatePath, append(raw, '\n'), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	raw, err := os.ReadFile(rowWriterTemplatePath)
	if err != nil {
		t.Fatal(err)
	}
	var stored rowWriterTemplates
	if err := json.Unmarshal(raw, &stored); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(writerSequences(stored), writerSequences(fresh)) {
		t.Fatalf("writer table/op sequences drifted:\nstored: %v\nfresh: %v\nregenerate with the command above", writerSequences(stored), writerSequences(fresh))
	}
}

func captureRowWriterTemplates(t *testing.T) rowWriterTemplates {
	t.Helper()
	workspace := gittest.NewWorkspace(t)
	captured := rowWriterTemplates{Templates: map[string]journal.Event{}}
	run := func(args ...string) journal.Event {
		t.Helper()
		result := executeCommandForTest(workspace, args...)
		if result.err != nil {
			t.Fatalf("%v: %v: %s", args, result.err, result.stderr)
		}
		e := lastCommandJournalEvent(t, workspace)
		if _, found := captured.Templates[e.EventType]; !found {
			captured.Templates[e.EventType] = e
		}
		return e
	}
	run("init", "--json", "--prefix", "FID")
	run("label", "add", "perf", "--actor", "human:mark", "--idempotency-key", "perf")
	captured.First = run("issue", "create", "--kind", "task", "--title", "First", "--actor", "human:mark", "--idempotency-key", "first").Aggregate.ID
	captured.Second = run("issue", "create", "--kind", "task", "--title", "Second", "--actor", "human:mark", "--idempotency-key", "second").Aggregate.ID
	run("issue", "ac", "add", captured.First, "--text", "Works", "--ordinal", "1", "--actor", "human:mark", "--idempotency-key", "ac")
	run("issue", "label", "add", captured.First, "perf", "--actor", "human:mark", "--idempotency-key", "label")
	run("issue", "transition", captured.First, "--to", "cancelled", "--actor", "human:mark", "--reason", "Out of scope", "--idempotency-key", "transition")
	run("link", "add", captured.First, "blocks", captured.Second, "--from-type", "issue", "--to-type", "issue", "--actor", "human:mark", "--idempotency-key", "link")
	return captured
}

func writerSequences(templates rowWriterTemplates) map[string][]string {
	result := make(map[string][]string, len(templates.Templates))
	for family, e := range templates.Templates {
		result[family] = []string{}
		for _, row := range e.Effects {
			result[family] = append(result[family], row.Table+"/"+row.Op)
		}
	}
	return result
}
