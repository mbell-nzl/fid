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
)

func TestP22PortableFakeAndFIDRejectUnknownEffects(t *testing.T) {
	t.Parallel()
	a := p22NewFIDRecoveryAdapter(t)
	path := filepath.Join(t.TempDir(), "valid.zip")
	p21Run(t, a.binary, a.workspace, "--json", "export", "create", path)
	control := p22ArchiveEntries(t, path)
	if _, err := p22PortableFakeImport(control); err != nil {
		t.Fatal(err)
	}
	p21Run(t, a.binary, a.workspace, "--json", "export", "validate", path)
	for _, kind := range []string{"table", "column", "operation", "event-type", "duplicate-event"} {
		t.Run(kind, func(t *testing.T) {
			entries := map[string][]byte{}
			for name, data := range control {
				entries[name] = bytes.Clone(data)
			}
			p22CorruptPortableEffect(t, entries, kind)
			if _, err := p22PortableFakeImport(entries); err == nil {
				t.Fatal("fake accepted malformed logical Effect archive")
			}
			archive := filepath.Join(t.TempDir(), "invalid.zip")
			p22WritePortableArchive(t, archive, entries)
			stdout, stderr, err := closureRunCommand(a.workspace, a.binary, "--json", "export", "validate", archive)
			response := p22ConformanceResponse{Stdout: stdout, Stderr: stderr}
			if err != nil {
				response.ExitCode = processExitCode(t, err)
			}
			assertP22ConformanceResponse(t, p22ConformanceCase{Operation: "export.validate", Kind: "invalid", ExitCode: 6, FailureCode: "archive-invalid"}, response)
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
			if len(effects) == 0 {
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
