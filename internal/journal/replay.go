package journal

import (
	"bufio"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	fiddb "github.com/mbell-nzl/fid/internal/db"
)

type ReplayStats struct {
	EventsReplayed int
	FilesReplayed  int
	UpgradedEvents int
}

type ReplayOptions struct {
	StartFile   string
	StartOffset int64
}

type CorruptJournalError struct {
	File string
	Line int
	Err  error
}

func (err *CorruptJournalError) Error() string {
	if err == nil {
		return "<nil>"
	}

	return fmt.Sprintf("malformed journal entry at %s:%d: %v", err.File, err.Line, err.Err)
}

func (err *CorruptJournalError) Unwrap() error {
	if err == nil {
		return nil
	}
	return err.Err
}

func ReplayJournal(conn *sql.DB, journalDir string, options *ReplayOptions) (ReplayStats, error) {
	if conn == nil {
		return ReplayStats{}, fmt.Errorf("replay journal: connection is nil")
	}

	return replayJournal(journalDir, options, func(input ProjectionInput) error {
		return ProjectEvent(conn, input)
	})
}

func replayJournal(journalDir string, options *ReplayOptions, project func(ProjectionInput) error) (ReplayStats, error) {
	files, err := orderedJournalFiles(journalDir)
	if err != nil {
		return ReplayStats{}, err
	}

	var stats ReplayStats
	tracker := newAggregateSequenceTracker()
	for _, fileName := range files {
		fileReplayed, err := replayJournalFile(project, journalDir, fileName, options, &stats, tracker)
		if err != nil {
			return ReplayStats{}, err
		}
		if fileReplayed {
			stats.FilesReplayed++
		}
	}

	return stats, nil
}

func replayJournalFile(project func(ProjectionInput) error, journalDir, fileName string, options *ReplayOptions, stats *ReplayStats, tracker *aggregateSequenceTracker) (bool, error) {
	path := filepath.Join(journalDir, fileName)
	file, err := os.Open(path)
	if err != nil {
		return false, fmt.Errorf("open journal file %q: %w", path, err)
	}
	defer closeReplayFile(file)

	reader := bufio.NewReader(file)
	var (
		lineNumber   int
		offset       int64
		fileReplayed bool
		startFile    string
		startOffset  int64
	)
	if options != nil {
		startFile = options.StartFile
		startOffset = options.StartOffset
	}

	for {
		lineBytes, readErr := reader.ReadBytes('\n')
		if readErr != nil && !errors.Is(readErr, io.EOF) {
			return false, fmt.Errorf("read journal file %q: %w", path, readErr)
		}
		if len(lineBytes) == 0 && errors.Is(readErr, io.EOF) {
			break
		}

		lineNumber++
		offset += int64(len(lineBytes))

		trimmed := bytesTrimSpace(lineBytes)
		if len(trimmed) == 0 {
			if errors.Is(readErr, io.EOF) {
				break
			}
			continue
		}

		var stored Event
		if err := json.Unmarshal(trimmed, &stored); err != nil {
			return false, &CorruptJournalError{
				File: filepath.ToSlash(filepath.Join(".fid", "journal", fileName)),
				Line: lineNumber,
				Err:  err,
			}
		}
		if err := tracker.add(stored); err != nil {
			return false, &CorruptJournalError{
				File: filepath.ToSlash(filepath.Join(".fid", "journal", fileName)),
				Line: lineNumber,
				Err:  err,
			}
		}
		if startFile != "" && (fileName < startFile || (fileName == startFile && offset <= startOffset)) {
			if errors.Is(readErr, io.EOF) {
				break
			}
			continue
		}

		eventHash, err := HashEvent(stored)
		if err != nil {
			return false, fmt.Errorf("hash event %q from %s:%d: %w", stored.EventID, fileName, lineNumber, err)
		}

		upgraded, err := UpgradeToCurrent(stored)
		if err != nil {
			return false, fmt.Errorf("upgrade event %q from %s:%d: %w", stored.EventID, fileName, lineNumber, err)
		}
		if upgraded.EventSchema != stored.EventSchema {
			stats.UpgradedEvents++
		}

		if err := project(ProjectionInput{
			Event:      upgraded,
			EventHash:  eventHash,
			File:       fileName,
			Line:       lineNumber,
			Offset:     offset,
			ReplayMode: fiddb.ReplayStatusClean,
		}); err != nil {
			var conflictErr *PreconditionConflictError
			if errors.As(err, &conflictErr) {
				return false, err
			}
			return false, &CorruptJournalError{
				File: filepath.ToSlash(filepath.Join(".fid", "journal", fileName)),
				Line: lineNumber,
				Err:  fmt.Errorf("project event %q: %w", upgraded.EventID, err),
			}
		}

		stats.EventsReplayed++
		fileReplayed = true

		if errors.Is(readErr, io.EOF) {
			break
		}
	}

	return fileReplayed, nil
}

func orderedJournalFiles(journalDir string) ([]string, error) {
	entries, err := os.ReadDir(journalDir)
	if err != nil {
		return nil, fmt.Errorf("read journal directory %q: %w", journalDir, err)
	}

	files := make([]string, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if !strings.HasSuffix(name, ".jsonl") {
			continue
		}
		if !journalFileName(name) {
			continue
		}
		files = append(files, name)
	}
	sort.Strings(files)
	return files, nil
}

func journalFileName(name string) bool {
	if !strings.HasSuffix(name, ".jsonl") {
		return false
	}
	stem := strings.TrimSuffix(name, ".jsonl")
	if _, err := time.Parse("2006-01-02", stem); err == nil {
		return true
	}
	if len(stem) != 8 {
		return false
	}
	_, err := strconv.Atoi(stem)
	return err == nil
}

func bytesTrimSpace(value []byte) []byte {
	return []byte(strings.TrimSpace(string(value)))
}

func closeReplayDB(conn *sql.DB) {
	if conn != nil {
		_ = conn.Close()
	}
}

func closeReplayFile(file *os.File) {
	if file != nil {
		_ = file.Close()
	}
}
