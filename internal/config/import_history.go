package config

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// Imported ordinal chunks are immutable. Later local Events use dated files.
// Verify both their bytes and the identity list used to fence old observations.
func verifyImportedHistory(workspace string, binding ImportBinding) error {
	if len(binding.JournalEntries) == 0 {
		return fmt.Errorf("read import binding: imported journal inventory is missing")
	}
	position := 0
	for index, entry := range binding.JournalEntries {
		if entry.Name != fmt.Sprintf("%08d.jsonl", index+1) || entry.SizeBytes <= 0 || entry.SizeBytes > 64<<20 {
			return fmt.Errorf("read import binding: invalid imported journal inventory")
		}
		if err := verifyImportedChunk(filepath.Join(workspace, JournalDirName, entry.Name), entry, binding.ImportedEventIDs, &position); err != nil {
			return err
		}
	}
	if position != len(binding.ImportedEventIDs) {
		return fmt.Errorf("read import binding: imported Event inventory differs from journal")
	}
	return nil
}

func verifyImportedChunk(path string, entry ImportJournalEntry, ids []string, position *int) error {
	info, err := os.Lstat(path)
	if err != nil || !info.Mode().IsRegular() || info.Size() != entry.SizeBytes {
		return fmt.Errorf("read import binding: imported journal file differs from inventory")
	}
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("read import binding: open journal: %w", err)
	}
	defer func() { _ = file.Close() }()
	hash := sha256.New()
	scanner := bufio.NewScanner(io.TeeReader(io.LimitReader(file, entry.SizeBytes+1), hash))
	scanner.Buffer(make([]byte, 64<<10), (1<<20)+1)
	for scanner.Scan() {
		var event struct {
			ID string `json:"event_id"`
		}
		if err := json.Unmarshal(scanner.Bytes(), &event); err != nil || *position >= len(ids) || event.ID != ids[*position] {
			return fmt.Errorf("read import binding: imported Event identity differs from journal")
		}
		*position++
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("read import binding: read journal: %w", err)
	}
	if "sha256:"+hex.EncodeToString(hash.Sum(nil)) != entry.Digest {
		return fmt.Errorf("read import binding: imported journal digest differs from inventory")
	}
	return nil
}
