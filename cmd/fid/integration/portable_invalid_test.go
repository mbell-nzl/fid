//go:build integration

package integration

import (
	"archive/zip"
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/mbell-nzl/fid/internal/gittest"
)

func TestPortableRejectsUnknownEffectsWithoutMutation(t *testing.T) {
	t.Parallel()
	workspace := gittest.NewWorkspace(t)
	binary := integrationFIDBinary(t)
	p21Run(t, binary, workspace, "init", "--prefix", "FID", "--json")
	p21Run(t, binary, workspace, "--json", "label", "add", "portable", "--actor", "agent:builder", "--colour", "blue", "--idempotency-key", "portable-label")
	p21Run(t, binary, workspace, "--json", "label", "add", "portable-tail", "--actor", "agent:builder", "--idempotency-key", "portable-tail")
	assertLeanIntegrationJournal(t, workspace)
	path := filepath.Join(t.TempDir(), "valid.zip")
	p21Run(t, binary, workspace, "--json", "export", "create", path)
	control := p22ArchiveEntries(t, path)
	p21Run(t, binary, workspace, "--json", "export", "validate", path)
	for _, kind := range []string{"table", "column", "operation", "event-type", "duplicate-event"} {
		t.Run(kind, func(t *testing.T) {
			entries := map[string][]byte{}
			for name, data := range control {
				entries[name] = bytes.Clone(data)
			}
			p22CorruptPortableEffect(t, entries, kind)
			archive := filepath.Join(t.TempDir(), "invalid.zip")
			p22WritePortableArchive(t, archive, entries)
			beforeJournal := p21JournalBytes(t, workspace)
			beforeArchive := readP21File(t, archive)
			stdout, stderr, err := closureRunCommand(workspace, binary, "--json", "export", "validate", archive)
			if err == nil || processExitCode(t, err) != 6 || stdout != "" {
				t.Fatalf("malformed portable archive=%v/%q/%q", err, stdout, stderr)
			}
			envelope := p23Envelope(t, stderr)
			if envelope["contract_version"] != "fid-cli.v1" || envelope["ok"] != false || envelope["operation"] != "export.validate" || p21PathString(t, envelope, "error", "code") != "archive-invalid" {
				t.Fatalf("archive refusal=%v", envelope)
			}

			if !bytes.Equal(beforeJournal, p21JournalBytes(t, workspace)) || !bytes.Equal(beforeArchive, readP21File(t, archive)) {
				t.Fatal("refusal changed source journal or archive")
			}
		})
	}
}

func p22CorruptPortableEffect(t *testing.T, entries map[string][]byte, kind string) {
	t.Helper()
	names := []string{}
	for name := range entries {
		if strings.HasPrefix(name, "journal/") {
			names = append(names, name)
		}
	}
	sort.Strings(names)
	firstID := ""
	changed := false
	for _, name := range names {
		lines := bytes.Split(bytes.TrimSuffix(entries[name], []byte{'\n'}), []byte{'\n'})
		for index, line := range lines {
			var event map[string]any
			if err := json.Unmarshal(line, &event); err != nil {
				t.Fatal(err)
			}
			if firstID == "" {
				firstID = event["event_id"].(string)
			}
			effects := event["effects"].([]any)
			if len(effects) == 0 || event["event_type"] == "project.initialized" {
				continue
			}
			effect := effects[0].(map[string]any)
			switch kind {
			case "table":
				effect["table"] = "unknown_portable_table"
			case "column":
				effect["values"].(map[string]any)["unknown_portable_column"] = "invalid"
			case "operation":
				effect["op"] = "execute"
			case "event-type":
				event["event_type"] = "run.finished"
			case "duplicate-event":
				event["event_id"] = firstID
			}
			lines[index], _ = json.Marshal(event)
			changed = true
			break
		}
		if changed {
			entries[name] = append(bytes.Join(lines, []byte{'\n'}), '\n')
			break
		}
	}
	if !changed {
		t.Fatal("portable negative has no changed Effect")
	}
	var manifest map[string]any
	if err := json.Unmarshal(entries["manifest.json"], &manifest); err != nil {
		t.Fatal(err)
	}
	for _, value := range manifest["entries"].([]any) {
		entry := value.(map[string]any)
		data := entries[entry["path"].(string)]
		entry["size_bytes"] = len(data)
		entry["digest"] = fmt.Sprintf("sha256:%x", sha256.Sum256(data))
	}
	entries["manifest.json"], _ = json.Marshal(manifest)
}

func p22WritePortableArchive(t *testing.T, path string, entries map[string][]byte) {
	t.Helper()
	file, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	writer := zip.NewWriter(file)
	names := []string{}
	for name := range entries {
		names = append(names, name)
	}
	sort.Strings(names)
	for _, name := range names {
		output, err := writer.Create(name)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := output.Write(entries[name]); err != nil {
			t.Fatal(err)
		}
	}
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}
	if err := file.Close(); err != nil {
		t.Fatal(err)
	}
}
