package main

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/mbell-nzl/fid/internal/config"
	"github.com/mbell-nzl/fid/internal/journal"
)

func TestPersistenceSnapshotsAcrossJournalDates(t *testing.T) {
	t.Parallel()
	workspace := twoDateJournalWorkspace(t)
	for _, helper := range []struct {
		name     string
		snapshot func(*testing.T, string) workspacePersistenceSnapshot
	}{
		{"workspace", snapshotWorkspacePersistence},
		{"binding", p24BindingPersistenceSnapshot},
		{"cancellation", snapshotEpicCancellationPersistence},
	} {
		t.Run(helper.name, func(t *testing.T) {
			t.Parallel()
			before := helper.snapshot(t, workspace)
			if got := before.eventCount(); got != 2 {
				t.Fatalf("snapshot events = %d, want both dates' 2 events", got)
			}
			if !before.equal(helper.snapshot(t, workspace)) {
				t.Fatal("unchanged two-date ledger differs from its snapshot")
			}
			for _, mutation := range []string{"first file", "last file", "add empty file", "remove file", "rename file", "move file boundary", "database"} {
				t.Run(mutation, func(t *testing.T) {
					t.Parallel()
					root := t.TempDir()
					cloneWorkspaceTreeInto(t, workspace, root)
					mutatePersistenceFixture(t, root, mutation)
					if before.equal(helper.snapshot(t, root)) {
						t.Fatalf("snapshot missed mutation: %s", mutation)
					}
				})
			}
		})
	}
}

func TestJournalEventCountAcrossDates(t *testing.T) {
	t.Parallel()
	if got := journalEventCount(t, twoDateJournalWorkspace(t)); got != 2 {
		t.Fatalf("journal events = %d, want both dates' 2 events", got)
	}
}

// Fixed times reproduce a suite whose pristine fixture predates UTC midnight.
func twoDateJournalWorkspace(t *testing.T) string {
	t.Helper()
	workspace := clonePristineWorkspace(t)
	backdateGenesis(t, workspace, time.Date(2026, 11, 13, 23, 59, 59, 0, time.UTC))
	appendCommandIssueEvent(t, workspace, journal.Event{
		EventSchema: journal.EventSchemaV1,
		EventType:   "label.seeded",
		OccurredAt:  time.Date(2026, 11, 14, 0, 0, 1, 0, time.UTC),
		Actor:       journal.Actor{Role: "tl"},
		Command:     journal.Command{Argv: []string{"fid", "maintenance", "seed"}, CWD: ".", FIDVersion: "0.1.0"},
		Effects: []journal.Effect{{Op: "insert", Table: "labels", PK: map[string]any{"name": "after-midnight"},
			Values: map[string]any{"description": "After UTC midnight"}}},
	})
	if result := executeCommandForTest(workspace, "rebuild"); result.err != nil {
		t.Fatalf("rebuild two-date ledger: %v: %s", result.err, result.stderr)
	}
	return workspace
}

func mutatePersistenceFixture(t *testing.T, workspace, mutation string) {
	t.Helper()
	paths := config.WorkspacePathsForRoot(workspace)
	first := filepath.Join(paths.JournalDir, "2026-11-13.jsonl")
	last := filepath.Join(paths.JournalDir, "2026-11-14.jsonl")
	write := func(path string, data []byte) {
		t.Helper()
		if err := os.WriteFile(path, data, 0o600); err != nil {
			t.Fatal(err)
		}
	}
	switch mutation {
	case "first file":
		write(first, append(mustReadCommandFile(t, first), '\n'))
	case "last file":
		write(last, append(mustReadCommandFile(t, last), '\n'))
	case "add empty file":
		write(filepath.Join(paths.JournalDir, "2026-11-15.jsonl"), nil)
	case "remove file":
		if err := os.Remove(first); err != nil {
			t.Fatal(err)
		}
	case "rename file":
		if err := os.Rename(first, filepath.Join(paths.JournalDir, "2026-11-12.jsonl")); err != nil {
			t.Fatal(err)
		}
	case "move file boundary":
		// Keep the concatenated bytes identical while changing both files.
		write(first, append(mustReadCommandFile(t, first), mustReadCommandFile(t, last)...))
		write(last, nil)
	case "database":
		write(paths.DatabasePath, append(mustReadCommandFile(t, paths.DatabasePath), 0))
	default:
		t.Fatalf("unknown persistence mutation %q", mutation)
	}
}
