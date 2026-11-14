//go:build integration

package integration

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/mbell-nzl/fid/internal/cli"
)

const p21CorpusDir = "testdata/p21_conformance"

func p21CorpusEntries(root string) ([]string, error) {
	var entries []string
	err := filepath.WalkDir(root, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			return nil
		}
		relative, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		entries = append(entries, relative)
		return nil
	})
	sort.Strings(entries)
	return entries, err
}

func p21CorpusDiff(want, got []byte) string {
	return fmt.Sprintf("want:\n%s\ngot:\n%s", want, got)
}

func TestRetainedConformanceCorpusIsCompleteAndByteExact(t *testing.T) {
	t.Parallel()
	operations := []map[string]string{}
	for _, operation := range []string{"activity.list", "artifact.verify", "export.create", "export.validate", "issue.ac.add", "issue.ac.list", "issue.create", "issue.get", "issue.label.add", "issue.label.remove", "issue.list", "issue.transition", "issue.update", "label.add", "label.list", "link.add", "link.list", "link.remove", "maintenance.check", "maintenance.vacuum", "project.init", "project.rebuild"} {
		command := "fid " + strings.ReplaceAll(operation, ".", " ")
		if positional := map[string]string{"export.create": "fid export create <archive-path>", "export.validate": "fid export validate <archive-path>", "issue.ac.add": "fid issue ac add <issue-id>", "issue.ac.list": "fid issue ac list <issue-id>", "issue.get": "fid issue get <id>", "issue.label.add": "fid issue label add <issue-id> <label>", "issue.label.remove": "fid issue label remove <issue-id> <label>", "issue.transition": "fid issue transition <id>", "issue.update": "fid issue update <id>", "label.add": "fid label add <name>", "link.add": "fid link add <from-id> <LinkRelation> <to-id>", "link.list": "fid link list <id>", "link.remove": "fid link remove <from-id> <LinkRelation> <to-id>"}[operation]; positional != "" {
			command = positional
		}
		if operation == "project.init" {
			command = "fid init"
		}
		if operation == "project.rebuild" {
			command = "fid rebuild"
		}
		kind := "accepted"
		if operation == "maintenance.vacuum" {
			kind = "fail-closed"
		}
		operations = append(operations, map[string]string{"operation": operation, "command": command, "kind": kind})
	}
	failures := []string{"idempotency-conflict"}
	for _, code := range cli.RegisteredCodes() {
		failures = append(failures, string(code))
	}
	sort.Strings(failures)
	expected := map[string][]byte{"manifest.json": p21JSON(t, map[string]any{"operations": operations, "failures": failures})}
	for _, row := range operations {
		expected["operations/"+row["operation"]+".json"] = p21JSON(t, row)
	}
	for _, code := range failures {
		expected["failures/"+code+".json"] = p21JSON(t, map[string]string{"failure_code": code})
	}

	entries, err := p21CorpusEntries(p21CorpusDir)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != len(expected) {
		t.Fatalf("corpus entries=%d, want %d", len(entries), len(expected))
	}
	for _, entry := range entries {
		want, found := expected[entry]
		if !found {
			t.Fatalf("unexpected corpus entry %s", entry)
		}
		if got := readP21File(t, filepath.Join(p21CorpusDir, entry)); !bytes.Equal(got, want) {
			t.Fatalf("corpus %s differs: %s", entry, p21CorpusDiff(want, got))
		}
	}
}
