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

type journalCheckpoint struct {
	File      string
	Offset    int64
	EventID   string
	EventHash string
	HasEvents bool
}

type journalInspection struct {
	tail           journalCheckpoint
	checkpointSeen bool
	checkpoint     journalCheckpoint
}

func inspectJournal(journalDir, targetFile string, targetOffset int64) (journalInspection, error) {
	files, err := orderedJournalFiles(journalDir)
	if err != nil {
		return journalInspection{}, err
	}

	var inspection journalInspection
	for _, fileName := range files {
		fileInspection, err := inspectJournalFile(journalDir, fileName, targetFile, targetOffset)
		if err != nil {
			return journalInspection{}, err
		}
		if fileInspection.tail.HasEvents {
			inspection.tail = fileInspection.tail
		}
		if fileInspection.checkpointSeen {
			inspection.checkpointSeen = true
			inspection.checkpoint = fileInspection.checkpoint
		}
	}

	return inspection, nil
}

func inspectJournalFile(journalDir, fileName, targetFile string, targetOffset int64) (journalInspection, error) {
	path := filepath.Join(journalDir, fileName)
	file, err := os.Open(path)
	if err != nil {
		return journalInspection{}, fmt.Errorf("open journal file %q: %w", path, err)
	}
	defer closeReplayFile(file)

	reader := bufio.NewReader(file)
	var (
		inspection journalInspection
		offset     int64
		lineNumber int
	)
	for {
		lineBytes, readErr := reader.ReadBytes('\n')
		if readErr != nil && !errors.Is(readErr, io.EOF) {
			return journalInspection{}, fmt.Errorf("read journal file %q: %w", path, readErr)
		}
		if len(lineBytes) == 0 && errors.Is(readErr, io.EOF) {
			break
		}

		lineNumber++
		offset += int64(len(lineBytes))
		checkpoint, isEmpty, err := inspectJournalLine(fileName, lineNumber, offset, lineBytes)
		if err != nil {
			return journalInspection{}, err
		}
		if isEmpty {
			if errors.Is(readErr, io.EOF) {
				break
			}
			continue
		}

		inspection.tail = checkpoint
		if fileName == targetFile && offset == targetOffset {
			inspection.checkpointSeen = true
			inspection.checkpoint = checkpoint
		}

		if errors.Is(readErr, io.EOF) {
			break
		}
	}

	return inspection, nil
}

func inspectJournalLine(fileName string, lineNumber int, offset int64, lineBytes []byte) (journalCheckpoint, bool, error) {
	trimmed := bytesTrimSpace(lineBytes)
	if len(trimmed) == 0 {
		return journalCheckpoint{}, true, nil
	}

	var stored Event
	if err := json.Unmarshal(trimmed, &stored); err != nil {
		return journalCheckpoint{}, false, &CorruptJournalError{
			File: filepath.ToSlash(filepath.Join(".fid", "journal", fileName)),
			Line: lineNumber,
			Err:  err,
		}
	}

	hash, err := HashEvent(stored)
	if err != nil {
		return journalCheckpoint{}, false, fmt.Errorf("hash event %q from %s: %w", stored.EventID, fileName, err)
	}

	return journalCheckpoint{
		File:      fileName,
		Offset:    offset,
		EventID:   stored.EventID,
		EventHash: hash,
		HasEvents: true,
	}, false, nil
}
