package export

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestArchiveWriterRejectsOversizedInventoryBeforePublication(t *testing.T) {
	for _, tc := range []struct {
		name, path string
		size       int64
	}{
		{"non-Artifact", "journal/00000001.jsonl", maxArchiveNonArtifactBytes + 1},
		{"Artifact", "artifacts/0000000000000000000000000000000000000000000000000000000000000000", maxArchiveArtifactBytes + 1},
	} {
		t.Run(tc.name, func(t *testing.T) {
			err := writeArchive(filepath.Join(t.TempDir(), "oversized.zip"), []archiveEntry{{entry: entry{Path: tc.path, SizeBytes: tc.size}}})
			if !errors.Is(err, errArchiveResourceLimit) {
				t.Fatalf("oversized writer inventory error=%v", err)
			}
		})
	}
}

func TestJournalWriterChunksAtLineBoundaryPreservingBytes(t *testing.T) {
	w := &journalChunkWriter{directory: t.TempDir()}
	defer w.close()
	line := bytes.Repeat([]byte("x"), 1<<20)
	for i := 0; i < 63; i++ {
		if err := w.append(line); err != nil {
			t.Fatal(err)
		}
	}
	if err := w.append(line[:len(line)-64]); err != nil {
		t.Fatal(err)
	}
	if err := w.append([]byte("last Event")); err != nil {
		t.Fatal(err)
	}
	if err := w.finish(); err != nil {
		t.Fatal(err)
	}
	if len(w.entries) != 2 || w.entries[0].SizeBytes != maxArchiveNonArtifactBytes || w.entries[0].Path != "journal/00000001.jsonl" || w.entries[1].Path != "journal/00000002.jsonl" {
		t.Fatalf("chunk inventory=%+v", w.entries)
	}
	count := 0
	for _, item := range w.entries {
		input, err := os.Open(item.filePath)
		if err != nil {
			t.Fatal(err)
		}
		err = readEventLines(input, func(actual []byte) error {
			want := line
			if count == 63 {
				want = line[:len(line)-64]
			}
			if count == 64 {
				want = []byte("last Event")
			}
			if !bytes.Equal(actual, want) {
				t.Fatalf("Event line %d changed", count)
			}
			count++
			return nil
		})
		_ = input.Close()
		if err != nil {
			t.Fatal(err)
		}
		if got, err := fileDigest(item.filePath); err != nil || got != item.Digest {
			t.Fatalf("chunk digest=%s err=%v", got, err)
		}
	}
	if count != 65 {
		t.Fatalf("Event lines=%d", count)
	}
}
