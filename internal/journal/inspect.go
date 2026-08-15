package journal

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type FoundEvent struct {
	Event       Event
	JournalFile string
	Line        int
	Hash        string
}

type EventNotFoundError struct {
	EventID string
}

func (err *EventNotFoundError) Error() string {
	if err == nil {
		return "<nil>"
	}
	return fmt.Sprintf("journal event %q not found", err.EventID)
}

type journalInspection struct {
	checkpointSeen bool
	checkpoint     journalCheckpoint
	tail           journalTail
}

type journalCheckpoint struct {
	EventID   string
	EventHash string
}

type journalTail struct {
	HasEvents bool
	File      string
	Offset    int64
}

func FindEvent(journalDir, eventID string) (FoundEvent, error) {
	if strings.TrimSpace(journalDir) == "" {
		return FoundEvent{}, fmt.Errorf("find event: journal dir is required")
	}
	if strings.TrimSpace(eventID) == "" {
		return FoundEvent{}, fmt.Errorf("find event: event id is required")
	}

	files, err := orderedJournalFiles(journalDir)
	if err != nil {
		return FoundEvent{}, err
	}

	for _, fileName := range files {
		found, ok, err := findEventInFile(journalDir, fileName, eventID)
		if err != nil {
			return FoundEvent{}, err
		}
		if ok {
			return found, nil
		}
	}

	return FoundEvent{}, &EventNotFoundError{EventID: eventID}
}

func inspectJournal(journalDir, checkpointFile string, checkpointOffset int64) (journalInspection, error) {
	files, err := orderedJournalFiles(journalDir)
	if err != nil {
		return journalInspection{}, err
	}

	var inspection journalInspection
	for _, fileName := range files {
		if err := inspectJournalFile(journalDir, fileName, checkpointFile, checkpointOffset, &inspection); err != nil {
			return journalInspection{}, err
		}
	}
	return inspection, nil
}

func findEventInFile(journalDir, fileName, eventID string) (FoundEvent, bool, error) {
	path := filepath.Join(journalDir, fileName)
	file, err := os.Open(path)
	if err != nil {
		return FoundEvent{}, false, fmt.Errorf("open journal file %q: %w", path, err)
	}
	defer closeReplayFile(file)

	reader := bufio.NewReader(file)
	for lineNumber := 1; ; lineNumber++ {
		lineBytes, readErr := reader.ReadBytes('\n')
		if readErr != nil && !errors.Is(readErr, io.EOF) {
			return FoundEvent{}, false, fmt.Errorf("read journal file %q: %w", path, readErr)
		}
		if len(lineBytes) == 0 && errors.Is(readErr, io.EOF) {
			return FoundEvent{}, false, nil
		}

		found, ok, err := inspectJournalLine(fileName, lineNumber, lineBytes, eventID)
		if err != nil {
			return FoundEvent{}, false, err
		}
		if ok {
			return found, true, nil
		}

		if errors.Is(readErr, io.EOF) {
			return FoundEvent{}, false, nil
		}
	}
}

func inspectJournalFile(journalDir, fileName, checkpointFile string, checkpointOffset int64, inspection *journalInspection) error {
	path := filepath.Join(journalDir, fileName)
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("open journal file %q: %w", path, err)
	}
	defer closeReplayFile(file)

	reader := bufio.NewReader(file)
	var offset int64
	for lineNumber := 1; ; lineNumber++ {
		lineBytes, readErr := reader.ReadBytes('\n')
		if readErr != nil && !errors.Is(readErr, io.EOF) {
			return fmt.Errorf("read journal file %q: %w", path, readErr)
		}
		if len(lineBytes) == 0 && errors.Is(readErr, io.EOF) {
			return nil
		}

		offset += int64(len(lineBytes))
		checkpoint, hasEvent, err := inspectCheckpointLine(fileName, lineNumber, offset, lineBytes)
		if err != nil {
			return err
		}
		if hasEvent {
			inspection.tail = journalTail{HasEvents: true, File: fileName, Offset: offset}
		}
		if checkpoint.FileMatches(fileName, checkpointFile, offset, checkpointOffset) {
			inspection.checkpointSeen = true
			inspection.checkpoint = journalCheckpoint(checkpoint)
		}

		if errors.Is(readErr, io.EOF) {
			return nil
		}
	}
}

type checkpointLine struct {
	EventID   string
	EventHash string
}

func (line checkpointLine) FileMatches(fileName, checkpointFile string, offset, checkpointOffset int64) bool {
	return fileName == checkpointFile && offset == checkpointOffset
}

func inspectCheckpointLine(fileName string, lineNumber int, offset int64, lineBytes []byte) (checkpointLine, bool, error) {
	trimmed := bytesTrimSpace(lineBytes)
	if len(trimmed) == 0 {
		return checkpointLine{}, false, nil
	}

	var stored Event
	if err := json.Unmarshal(trimmed, &stored); err != nil {
		return checkpointLine{}, false, &CorruptJournalError{
			File: filepath.ToSlash(filepath.Join(".fid", "journal", fileName)),
			Line: lineNumber,
			Err:  err,
		}
	}

	eventHash, err := HashEvent(stored)
	if err != nil {
		return checkpointLine{}, false, fmt.Errorf("hash event %q from %s:%d: %w", stored.EventID, fileName, lineNumber, err)
	}
	_ = offset
	return checkpointLine{EventID: stored.EventID, EventHash: eventHash}, true, nil
}

func inspectJournalLine(fileName string, lineNumber int, lineBytes []byte, eventID string) (FoundEvent, bool, error) {
	trimmed := bytesTrimSpace(lineBytes)
	if len(trimmed) == 0 {
		return FoundEvent{}, false, nil
	}

	var stored Event
	if err := json.Unmarshal(trimmed, &stored); err != nil {
		return FoundEvent{}, false, &CorruptJournalError{
			File: filepath.ToSlash(filepath.Join(".fid", "journal", fileName)),
			Line: lineNumber,
			Err:  err,
		}
	}
	if stored.EventID != eventID {
		return FoundEvent{}, false, nil
	}

	eventHash, err := HashEvent(stored)
	if err != nil {
		return FoundEvent{}, false, fmt.Errorf("hash event %q from %s:%d: %w", stored.EventID, fileName, lineNumber, err)
	}

	return FoundEvent{
		Event:       stored,
		JournalFile: filepath.ToSlash(filepath.Join(".fid", "journal", fileName)),
		Line:        lineNumber,
		Hash:        eventHash,
	}, true, nil
}
