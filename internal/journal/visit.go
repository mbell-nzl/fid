package journal

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// Visit reads Events in replay order until visit returns true. It does not
// project or admit Events; callers first refresh their projection as usual.
// In particular, callers can stop at a canonical receipt without relying on
// caller-supplied timestamps or comparing rowids from different tables.
func Visit(journalDir string, visit func(Event) (bool, error)) error {
	files, err := orderedJournalFiles(journalDir)
	if err != nil {
		return err
	}
	for _, name := range files {
		stop, err := visitFile(filepath.Join(journalDir, name), visit)
		if err != nil || stop {
			return err
		}
	}
	return nil
}

func visitFile(path string, visit func(Event) (bool, error)) (bool, error) {
	file, err := os.Open(path)
	if err != nil {
		return false, fmt.Errorf("open journal file %q: %w", path, err)
	}
	defer closeReplayFile(file)
	reader := bufio.NewReader(file)
	for line := 1; ; line++ {
		data, readErr := reader.ReadBytes('\n')
		if readErr != nil && !errors.Is(readErr, io.EOF) {
			return false, fmt.Errorf("read journal file %q: %w", path, readErr)
		}
		if trimmed := bytesTrimSpace(data); len(trimmed) != 0 {
			var event Event
			if err := json.Unmarshal(trimmed, &event); err != nil {
				return false, &CorruptJournalError{File: filepath.ToSlash(path), Line: line, Err: err}
			}
			stop, err := visit(event)
			if err != nil || stop {
				return stop, err
			}
		}
		if errors.Is(readErr, io.EOF) {
			return false, nil
		}
	}
}
