package journal

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func AppendEvent(lock *WriteLock, event Event) (Event, error) {
	return appendEventWithClock(lock, event, time.Now)
}

func PrepareEvent(lock *WriteLock, event Event) (Event, error) {
	return prepareEventWithClock(lock, event, time.Now)
}

func appendEventWithClock(lock *WriteLock, event Event, now func() time.Time) (Event, error) {
	prepared, err := prepareEventWithClock(lock, event, now)
	if err != nil {
		return Event{}, err
	}

	payload, err := MarshalCanonicalJSON(prepared)
	if err != nil {
		return Event{}, err
	}
	payload = append(payload, '\n')

	journalPath := filepath.Join(lock.workspaceRoot, ".fid", "journal", prepared.OccurredAt.UTC().Format("2006-01-02")+".jsonl")
	file, err := os.OpenFile(journalPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return Event{}, fmt.Errorf("open journal file %q: %w", journalPath, err)
	}

	if _, err := file.Write(payload); err != nil {
		_ = file.Close()
		return Event{}, fmt.Errorf("append journal file %q: %w", journalPath, err)
	}
	if err := file.Sync(); err != nil {
		_ = file.Close()
		return Event{}, fmt.Errorf("sync journal file %q: %w", journalPath, err)
	}
	if err := file.Close(); err != nil {
		return Event{}, fmt.Errorf("close journal file %q: %w", journalPath, err)
	}

	return prepared, nil
}

func prepareEventWithClock(lock *WriteLock, event Event, now func() time.Time) (Event, error) {
	if lock == nil {
		return Event{}, fmt.Errorf("append event: lock is nil")
	}
	if lock.released {
		return Event{}, fmt.Errorf("append event: lock %q is already released", lock.path)
	}
	if now == nil {
		return Event{}, fmt.Errorf("append event: clock is nil")
	}

	occurredAt := now().UTC()
	if !event.OccurredAt.IsZero() {
		occurredAt = event.OccurredAt.UTC()
	}
	journalDir := filepath.Join(lock.workspaceRoot, ".fid", "journal")
	latestDate, hasLatest, err := greatestJournalDate(journalDir)
	if err != nil {
		return Event{}, err
	}

	chosenDate := occurredAt.Format("2006-01-02")
	if hasLatest && chosenDate < latestDate {
		return Event{}, fmt.Errorf("system clock is behind latest journal date %s; refusing to append out of order", latestDate)
	}

	journalPath := filepath.Join(journalDir, chosenDate+".jsonl")
	sequence, err := nextSequence(journalPath)
	if err != nil {
		return Event{}, err
	}

	if event.OccurredAt.IsZero() {
		event.OccurredAt = occurredAt
	}
	if event.EventID == "" {
		event.EventID, err = FormatEventID(event.OccurredAt, sequence)
		if err != nil {
			return Event{}, err
		}
	}
	if event.IsV2() {
		if err := validateV2Event(event); err != nil {
			return Event{}, err
		}
		if err := validateNextAggregateSequence(journalDir, event); err != nil {
			return Event{}, err
		}
	}

	return event, nil
}

func greatestJournalDate(journalDir string) (string, bool, error) {
	entries, err := os.ReadDir(journalDir)
	if err != nil {
		return "", false, fmt.Errorf("read journal directory %q: %w", journalDir, err)
	}

	var latest string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		if !strings.HasSuffix(name, ".jsonl") {
			continue
		}

		dateText := strings.TrimSuffix(name, ".jsonl")
		if _, err := time.Parse("2006-01-02", dateText); err != nil {
			continue
		}
		if latest == "" || dateText > latest {
			latest = dateText
		}
	}

	return latest, latest != "", nil
}

func nextSequence(journalPath string) (int, error) {
	file, err := os.Open(journalPath)
	if err != nil {
		if os.IsNotExist(err) {
			return 1, nil
		}
		return 0, fmt.Errorf("open journal file %q for sequence scan: %w", journalPath, err)
	}
	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)

	highest := 0
	for scanner.Scan() {
		line := bytes.TrimSpace(scanner.Bytes())
		if len(line) == 0 {
			continue
		}

		var event Event
		if err := json.Unmarshal(line, &event); err != nil {
			_ = file.Close()
			return 0, fmt.Errorf("scan journal file %q: decode event: %w", journalPath, err)
		}

		sequence, err := parseEventSequence(event.EventID)
		if err != nil {
			_ = file.Close()
			return 0, fmt.Errorf("scan journal file %q: %w", journalPath, err)
		}
		if sequence > highest {
			highest = sequence
		}
	}
	if err := scanner.Err(); err != nil {
		_ = file.Close()
		return 0, fmt.Errorf("scan journal file %q: %w", journalPath, err)
	}
	if err := file.Close(); err != nil {
		return 0, fmt.Errorf("close journal file %q after sequence scan: %w", journalPath, err)
	}

	return highest + 1, nil
}

func parseEventSequence(eventID string) (int, error) {
	parts := strings.Split(eventID, "-")
	if len(parts) != 3 {
		return 0, fmt.Errorf("parse event id %q: want EVT-<date>-<sequence>", eventID)
	}
	if parts[0] != "EVT" {
		return 0, fmt.Errorf("parse event id %q: want EVT prefix", eventID)
	}
	if _, err := time.Parse("20060102", parts[1]); err != nil {
		return 0, fmt.Errorf("parse event id %q date: %w", eventID, err)
	}

	sequence, err := strconv.Atoi(parts[2])
	if err != nil {
		return 0, fmt.Errorf("parse event id %q sequence: %w", eventID, err)
	}
	if sequence <= 0 {
		return 0, fmt.Errorf("parse event id %q sequence: must be positive", eventID)
	}

	return sequence, nil
}
