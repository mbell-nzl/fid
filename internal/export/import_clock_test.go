package export

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/mbell-nzl/fid/internal/bootstrap"
	"github.com/mbell-nzl/fid/internal/config"
	"github.com/mbell-nzl/fid/internal/gittest"
	"github.com/mbell-nzl/fid/internal/issues"
	"github.com/mbell-nzl/fid/internal/journal"
)

// An import leaves the journal as the archive's numbered pieces, with no
// dated file. The journal still never appends before its latest date, so an
// append dated before the last imported Event is refused, as it would be in
// the ledger the archive came from.
func TestAppendAfterImportRefusesADateBeforeTheLastImportedEvent(t *testing.T) {
	t.Parallel()
	source := gittest.NewWorkspace(t)
	if _, err := bootstrap.InitializeWorkspace(bootstrap.InitOptions{WorkspaceRoot: source, MigrationsDir: testMigrationsDir(t), Prefix: "FID"}); err != nil {
		t.Fatalf("initialize source workspace: %v", err)
	}
	if _, err := issues.Create(issues.CreateOptions{WorkspaceRoot: source, MigrationsDir: testMigrationsDir(t), Kind: "task", Title: "Before the import", Actor: "agent:fixture", IdempotencyKey: "clock-issue"}); err != nil {
		t.Fatalf("create source issue: %v", err)
	}
	archive := filepath.Join(t.TempDir(), "recovery.zip")
	if _, err := Create(CreateOptions{WorkspaceRoot: source, ArchivePath: archive}); err != nil {
		t.Fatalf("create archive: %v", err)
	}
	target := gittest.NewWorkspace(t)
	if _, err := Import(ImportOptions{ArchivePath: archive, RepositoryRoot: target, MigrationsDir: testMigrationsDir(t), IdempotencyKey: "clock"}); err != nil {
		t.Fatalf("Import() error = %v", err)
	}
	journalDir := filepath.Join(config.LedgerDir(target), config.JournalDirName)
	before := journalDirectorySnapshot(t, journalDir)
	if names := strings.Join(sortedKeys(before), ","); names != "00000001.jsonl" {
		t.Fatalf("imported journal files = %s, want only the archive's numbered piece", names)
	}
	last := lastJournalEventTime(t, before["00000001.jsonl"])

	err := appendClockTestEvent(target, last.Add(-24*time.Hour))
	want := "system clock is behind latest journal date " + last.UTC().Format("2006-01-02") + "; refusing to append out of order"
	if err == nil || !strings.Contains(err.Error(), want) {
		t.Fatalf("append dated before the last imported Event = %v, want %q", err, want)
	}
	if after := journalDirectorySnapshot(t, journalDir); !equalSnapshots(before, after) {
		t.Fatalf("refused append changed the journal: %v, want %v", sortedKeys(after), sortedKeys(before))
	}

	if err := appendClockTestEvent(target, last); err != nil {
		t.Fatalf("append on the last imported Event's date = %v, want it appended", err)
	}
	if after := journalDirectorySnapshot(t, journalDir); len(after) != 2 || !bytes.Equal(after["00000001.jsonl"], before["00000001.jsonl"]) {
		t.Fatalf("journal files after an in-order append = %v, want the piece kept and one dated file", sortedKeys(after))
	}
}

func appendClockTestEvent(workspace string, at time.Time) error {
	return journal.WithLock(workspace, journal.LockOwner{Command: "clock test"}, func(lock *journal.WriteLock) error {
		_, err := journal.AppendEvent(lock, journal.Event{
			EventSchema: journal.EventSchemaV1, EventType: "issue.created", OccurredAt: at,
			Actor:   journal.Actor{Role: "developer"},
			Command: journal.Command{Argv: []string{"fid", "issue", "create"}, CWD: ".", FIDVersion: "0.1.0"},
			Effects: []journal.Effect{{Op: "insert", Table: "issues", PK: map[string]any{"id": "FID-clock"}, Values: map[string]any{"title": "Out of order"}}},
		})
		return err
	})
}

func journalDirectorySnapshot(t *testing.T, directory string) map[string][]byte {
	t.Helper()
	entries, err := os.ReadDir(directory)
	if err != nil {
		t.Fatalf("read journal directory: %v", err)
	}
	snapshot := map[string][]byte{}
	for _, entry := range entries {
		data, err := os.ReadFile(filepath.Join(directory, entry.Name()))
		if err != nil {
			t.Fatalf("read journal file %q: %v", entry.Name(), err)
		}
		snapshot[entry.Name()] = data
	}
	return snapshot
}

func lastJournalEventTime(t *testing.T, piece []byte) time.Time {
	t.Helper()
	lines := bytes.Split(bytes.TrimSpace(piece), []byte("\n"))
	var event struct {
		OccurredAt time.Time `json:"occurred_at"`
	}
	if err := json.Unmarshal(lines[len(lines)-1], &event); err != nil || event.OccurredAt.IsZero() {
		t.Fatalf("decode the last imported Event %q: %v", lines[len(lines)-1], err)
	}
	return event.OccurredAt
}

func sortedKeys(values map[string][]byte) []string {
	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func equalSnapshots(left, right map[string][]byte) bool {
	if len(left) != len(right) {
		return false
	}
	for name, data := range left {
		if other, ok := right[name]; !ok || !bytes.Equal(data, other) {
			return false
		}
	}
	return true
}
