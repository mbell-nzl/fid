package journal

import (
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const (
	task9PrechangeReference     = "74c83b7033729785eb99403ecec9dcde30ac4280"
	task9PrechangeJournalDigest = "8fdf19380c9cf8e1627bbde730e5554e3d0f3b3006ffcff7e21d6162197f11df"
)

func TestTask9PrechangeBudgetObservedJournalReplaysWithoutMutation(t *testing.T) {
	t.Parallel()
	encoded, err := os.ReadFile(filepath.Join("testdata", "task9-prechange-budget-observed.jsonl.gz.b64"))
	if err != nil {
		t.Fatal(err)
	}
	compressed, err := base64.StdEncoding.DecodeString(strings.Join(strings.Fields(string(encoded)), ""))
	if err != nil {
		t.Fatal(err)
	}
	reader, err := gzip.NewReader(bytes.NewReader(compressed))
	if err != nil {
		t.Fatal(err)
	}
	journalBytes, err := io.ReadAll(reader)
	if err != nil {
		t.Fatal(err)
	}
	if err := reader.Close(); err != nil {
		t.Fatal(err)
	}
	if task9PrechangeReference == "" || !bytes.Contains(journalBytes, []byte(`\"type\":\"budget-observed\"`)) {
		t.Fatal("fixture lacks named reference or budget-observed Stop")
	}
	if got := task9JournalDigest(journalBytes); got != task9PrechangeJournalDigest {
		t.Fatalf("source journal digest = %s, want %s", got, task9PrechangeJournalDigest)
	}

	workspace := t.TempDir()
	journalDir := filepath.Join(workspace, ".fid", "journal")
	if err := os.MkdirAll(journalDir, 0o755); err != nil {
		t.Fatal(err)
	}
	journalPath := filepath.Join(journalDir, "2026-09-20.jsonl")
	if err := os.WriteFile(journalPath, journalBytes, 0o600); err != nil {
		t.Fatal(err)
	}
	stats, err := RebuildWorkspace(workspace, testProjectMigrationsDir(t))
	if err != nil {
		t.Fatalf("replay pre-change journal: %v", err)
	}
	if stats.EventsReplayed != 18 {
		t.Fatalf("replayed Events = %d, want 18", stats.EventsReplayed)
	}
	after, err := os.ReadFile(journalPath)
	if err != nil {
		t.Fatal(err)
	}
	if got := task9JournalDigest(after); got != task9PrechangeJournalDigest || !bytes.Equal(after, journalBytes) {
		t.Fatalf("post-replay journal digest = %s, want unchanged %s", got, task9PrechangeJournalDigest)
	}
}

func task9JournalDigest(payload []byte) string {
	sum := sha256.Sum256(payload)
	return hex.EncodeToString(sum[:])
}
