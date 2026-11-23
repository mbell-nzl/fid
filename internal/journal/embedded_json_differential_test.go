package journal_test

import (
	"archive/zip"
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/mbell-nzl/fid/internal/bootstrap"
	"github.com/mbell-nzl/fid/internal/config"
	"github.com/mbell-nzl/fid/internal/export"
	"github.com/mbell-nzl/fid/internal/gittest"
	"github.com/mbell-nzl/fid/internal/issues"
	"github.com/mbell-nzl/fid/internal/journal"
)

func TestEmbeddedJSONDifferential(t *testing.T) {
	t.Parallel()
	t.Run("readers", journal.CheckEmbeddedJSONReadersForTest)
	source := gittest.NewWorkspace(t)
	migrations := journal.JSONMigrationFixtureForTest(t)
	if _, err := bootstrap.InitializeWorkspace(bootstrap.InitOptions{WorkspaceRoot: source, MigrationsDir: migrations, Prefix: "FID"}); err != nil {
		t.Fatal(err)
	}
	archive := filepath.Join(t.TempDir(), "valid.zip")
	if _, err := export.Create(export.CreateOptions{WorkspaceRoot: source, ArchivePath: archive}); err != nil {
		t.Fatal(err)
	}
	original := readJSONTestArchive(t, archive)
	probes := journal.EmbeddedJSONRepresentativesForTest(journal.EmbeddedJSONMatrixForTest(t))
	for i := range probes {
		probes[i].MigrationsDir = migrations
	}
	t.Run("paths", func(t *testing.T) { journal.CheckEmbeddedJSONMatrixPathsForTest(t, probes) })
	// Bound import concurrency alongside the seven journal paths. Each worker
	// keeps refusal, unchanged-target and recovery checks for its representatives.
	const importWorkers = 4
	for worker := range importWorkers {
		t.Run(fmt.Sprintf("imports-%d", worker), func(t *testing.T) {
			t.Parallel()
			for i := worker; i < len(probes); i += importWorkers {
				probe := probes[i]
				t.Run(probe.Name, func(t *testing.T) {
					event := probe.Event
					event.EventID = "EVT-20261110-000001"
					malformed := appendJSONTestArchive(t, original, event)
					target := gittest.NewWorkspace(t)
					before := jsonTestTree(t, target)
					_, err := export.Import(export.ImportOptions{ArchivePath: malformed, RepositoryRoot: target, MigrationsDir: migrations, IdempotencyKey: "malformed"})
					var refusal *export.InvalidArchiveError
					if !errors.As(err, &refusal) {
						t.Fatalf("import = %v, want archive refusal", err)
					}
					if len(refusal.Checks) != 1 || !bytes.Contains([]byte(refusal.Checks[0].Detail), []byte(probe.Reason)) {
						t.Fatalf("import checks = %#v, want %s refusal", refusal.Checks, probe.Reason)
					}

					if !bytes.Equal(before, jsonTestTree(t, target)) {
						t.Fatal("refusal changed import target")
					}
					if _, err := export.Import(export.ImportOptions{ArchivePath: archive, RepositoryRoot: target, MigrationsDir: migrations, IdempotencyKey: "valid"}); err != nil {
						t.Fatal(err)
					}
				})
			}
		})
	}
	// Every representative retains its refusal, unchanged-target assertion and
	// successful recovery import. The accepted archive is identical for every
	// case, so exercise its write/export/validate/restore journey once overall.
	t.Run("accepted", func(t *testing.T) {
		t.Parallel()
		target := gittest.NewWorkspace(t)
		if _, err := export.Import(export.ImportOptions{ArchivePath: archive, RepositoryRoot: target, MigrationsDir: migrations, IdempotencyKey: "valid"}); err != nil {
			t.Fatal(err)
		}
		if _, err := issues.Create(issues.CreateOptions{WorkspaceRoot: target, MigrationsDir: migrations, Kind: "task", Title: "After restore", Actor: "human:mark", IdempotencyKey: "later"}); err != nil {
			t.Fatalf("valid follow-up: %v", err)
		}
		accepted := filepath.Join(t.TempDir(), "accepted.zip")
		if _, err := export.Create(export.CreateOptions{WorkspaceRoot: target, ArchivePath: accepted}); err != nil {
			t.Fatal(err)
		}
		if _, err := export.Validate(export.ValidateOptions{ArchivePath: accepted, MigrationsDir: migrations}); err != nil {
			t.Fatal(err)
		}
		restored := gittest.NewWorkspace(t)
		if _, err := export.Import(export.ImportOptions{ArchivePath: accepted, RepositoryRoot: restored, MigrationsDir: migrations, IdempotencyKey: "round-trip"}); err != nil {
			t.Fatal(err)
		}
	})
}

func readJSONTestArchive(t *testing.T, path string) map[string][]byte {
	t.Helper()
	r, err := zip.OpenReader(path)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = r.Close() }()
	entries := map[string][]byte{}
	for _, f := range r.File {
		input, err := f.Open()
		if err != nil {
			t.Fatal(err)
		}
		entries[f.Name], err = io.ReadAll(input)
		if err != nil {
			t.Fatal(err)
		}
		if err := input.Close(); err != nil {
			t.Fatal(err)
		}
	}
	return entries
}

func appendJSONTestArchive(t *testing.T, original map[string][]byte, event journal.Event) string {
	t.Helper()
	var manifest map[string]any
	if err := json.Unmarshal(original["manifest.json"], &manifest); err != nil {
		t.Fatal(err)
	}
	raw, err := journal.MarshalCanonicalJSON(event)
	if err != nil {
		t.Fatal(err)
	}
	hash, err := journal.HashEvent(event)
	if err != nil {
		t.Fatal(err)
	}
	manifest["event_count"] = manifest["event_count"].(float64) + 1
	manifest["terminal_event_hash"] = hash
	entries := map[string][]byte{}
	for name, data := range original {
		entries[name] = bytes.Clone(data)
	}
	for _, item := range manifest["entries"].([]any) {
		entry := item.(map[string]any)
		path := entry["path"].(string)
		if !bytes.HasPrefix([]byte(path), []byte("journal/")) {
			continue
		}
		entries[path] = append(entries[path], append(raw, '\n')...)
		entry["digest"] = fmt.Sprintf("sha256:%x", sha256.Sum256(entries[path]))
		entry["size_bytes"] = len(entries[path])
	}
	entries["manifest.json"], err = journal.MarshalCanonicalJSON(manifest)
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(t.TempDir(), "malformed.zip")
	output, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	writer := zip.NewWriter(output)
	for name, data := range entries {
		w, err := writer.Create(name)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := w.Write(data); err != nil {
			t.Fatal(err)
		}
	}
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}
	if err := output.Close(); err != nil {
		t.Fatal(err)
	}
	return path
}

func jsonTestTree(t *testing.T, root string) []byte {
	t.Helper()
	data := []byte{}
	err := filepath.WalkDir(root, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			return nil
		}
		raw, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		data = append(data, []byte(path)...)
		data = append(data, raw...)
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	// Also pin a ledger path in this helper: a refused public import must
	// leave the journal and projection absent as well as preserve git files.
	if _, err := os.Stat(config.LedgerDir(root)); err == nil {
		t.Fatal("unexpected ledger before/after refused import")
	}
	return data
}
