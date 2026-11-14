package journal

import (
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mbell-nzl/fid/internal/config"
	"github.com/mbell-nzl/fid/internal/gittest"
)

const (
	task9PrechangeReference     = "74c83b7033729785eb99403ecec9dcde30ac4280"
	task9PrechangeJournalDigest = "8fdf19380c9cf8e1627bbde730e5554e3d0f3b3006ffcff7e21d6162197f11df"
)

func TestTask9PrechangeBudgetObservedJournalFailsClosedWithoutMutation(t *testing.T) {
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

	workspace := gittest.NewWorkspace(t)
	journalDir := filepath.Join(config.LedgerDir(workspace), "journal")
	if err := os.MkdirAll(journalDir, 0o755); err != nil {
		t.Fatal(err)
	}
	journalPath := filepath.Join(journalDir, "2026-09-20.jsonl")
	if err := os.WriteFile(journalPath, journalBytes, 0o600); err != nil {
		t.Fatal(err)
	}
	_, err = RebuildWorkspace(workspace, testProjectMigrationsDir(t))
	var corrupt *CorruptJournalError
	if !errors.As(err, &corrupt) || corrupt.Line != 3 || !strings.Contains(err.Error(), "schema-unsupported") {
		t.Fatalf("rebuild pre-change journal = %v, want line-3 unsupported-family refusal", err)
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
