//go:build integration

package integration

import (
	"bytes"
	"path/filepath"
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
