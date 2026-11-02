package export

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"hash"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/mbell-nzl/fid/internal/config"
	"github.com/mbell-nzl/fid/internal/journal"
)

func journalArchiveEntries(workspaceRoot, scratch string) ([]archiveEntry, int, string, error) {
	journalDir := filepath.Join(config.LedgerDir(workspaceRoot), config.JournalDirName)
	files, err := os.ReadDir(journalDir)
	if err != nil {
		return nil, 0, "", fmt.Errorf("create export: read journal: %w", err)
	}
	chunks := &journalChunkWriter{directory: scratch}
	defer chunks.close()
	tracker := newJournalTracker()
	count, terminal := 0, ""
	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".jsonl") {
			continue
		}
		input, err := os.Open(filepath.Join(journalDir, file.Name()))
		if err != nil {
			return nil, 0, "", fmt.Errorf("create export: open journal: %w", err)
		}
		err = readEventLines(input, func(line []byte) error {
			hash, err := validateJournalLine(line, tracker)
			if err != nil {
				return err
			}
			if err := chunks.append(line); err != nil {
				return err
			}
			count++
			terminal = hash
			return nil
		})
		closeErr := input.Close()
		if err != nil {
			return nil, 0, "", fmt.Errorf("create export: journal %q: %w", file.Name(), err)
		}
		if closeErr != nil {
			return nil, 0, "", closeErr
		}
	}
	if count == 0 {
		return nil, 0, "", fmt.Errorf("create export: journal contains no Events")
	}
	if err := chunks.finish(); err != nil {
		return nil, 0, "", err
	}
	return chunks.entries, count, terminal, nil
}

type journalChunkWriter struct {
	directory string
	output    *os.File
	hash      hash.Hash
	size      int64
	total     int64
	entries   []archiveEntry
}

func (w *journalChunkWriter) append(line []byte) error {
	if int64(len(line))+1 > maxArchiveExpandedBytes-w.total {
		return errArchiveResourceLimit
	}
	if w.output != nil && w.size+int64(len(line))+1 > maxArchiveNonArtifactBytes {
		if err := w.finish(); err != nil {
			return err
		}
	}
	if w.output == nil {
		name := fmt.Sprintf("%08d.jsonl", len(w.entries)+1)
		output, err := os.Create(filepath.Join(w.directory, name))
		if err != nil {
			return err
		}
		w.output, w.hash, w.size = output, sha256.New(), 0
	}
	output := io.MultiWriter(w.output, w.hash)
	if _, err := output.Write(line); err != nil {
		return err
	}
	if _, err := output.Write([]byte{'\n'}); err != nil {
		return err
	}
	w.size += int64(len(line)) + 1
	w.total += int64(len(line)) + 1
	return nil
}

func (w *journalChunkWriter) finish() error {
	if w.output == nil {
		return nil
	}
	name := w.output.Name()
	if err := w.output.Close(); err != nil {
		return err
	}
	w.entries = append(w.entries, archiveEntry{entry: entry{Path: "journal/" + filepath.Base(name), Digest: "sha256:" + hex.EncodeToString(w.hash.Sum(nil)), SizeBytes: w.size, MediaType: "application/x-ndjson"}, filePath: name})
	w.output = nil
	return nil
}

func (w *journalChunkWriter) close() {
	if w.output != nil {
		_ = w.output.Close()
	}
}

func validateJournalLine(line []byte, tracker *journalTracker) (string, error) {
	var event journal.Event
	if err := event.UnmarshalJSON(line); err != nil {
		return "", err
	}
	if err := tracker.add(event); err != nil {
		return "", err
	}
	canonical, err := journal.MarshalCanonicalJSON(event)
	if err != nil {
		return "", err
	}
	if !bytes.Equal(line, canonical) {
		return "", fmt.Errorf("event is not canonical")
	}
	return journal.HashEvent(event)
}

type journalTracker struct {
	eventIDs   map[string]struct{}
	aggregates map[string]int
}

func newJournalTracker() *journalTracker {
	return &journalTracker{eventIDs: make(map[string]struct{}), aggregates: make(map[string]int)}
}

func (tracker *journalTracker) add(event journal.Event) error {
	if _, exists := tracker.eventIDs[event.EventID]; exists {
		return fmt.Errorf("duplicate Event ID %q", event.EventID)
	}
	tracker.eventIDs[event.EventID] = struct{}{}
	if event.EventSchema != journal.EventSchemaV2 {
		return fmt.Errorf("unsupported Event schema %d", event.EventSchema)
	}
	key := event.Aggregate.Type + "\x00" + event.Aggregate.ID
	if prior, found := tracker.aggregates[key]; found && event.Aggregate.Sequence != prior+1 {
		return fmt.Errorf("aggregate %q sequence %d follows %d", event.Aggregate.ID, event.Aggregate.Sequence, prior)
	}
	if _, found := tracker.aggregates[key]; !found && event.Aggregate.Sequence != 1 {
		return fmt.Errorf("aggregate %q starts at sequence %d", event.Aggregate.ID, event.Aggregate.Sequence)
	}
	tracker.aggregates[key] = event.Aggregate.Sequence
	return nil
}
