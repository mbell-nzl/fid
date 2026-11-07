package journal

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
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

func prepareIdempotentEvent(lock *WriteLock, event Event) (Event, error) {
	return prepareEventWithClockOptions(lock, event, time.Now, true)
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

	journalPath := filepath.Join(lock.ledgerDir, "journal", prepared.OccurredAt.UTC().Format("2006-01-02")+".jsonl")
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
	return prepareEventWithClockOptions(lock, event, now, false)
}

func prepareEventWithClockOptions(lock *WriteLock, event Event, now func() time.Time, allowIdempotencyEffectOnly bool) (Event, error) {
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
	journalDir := filepath.Join(lock.ledgerDir, "journal")
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
		if err := validateV2EventFields(event, !allowIdempotencyEffectOnly); err != nil {
			return Event{}, err
		}
		if err := validateNextAggregateSequence(journalDir, event); err != nil {
			return Event{}, err
		}
	}

	return event, nil
}

// greatestJournalDate returns the latest UTC date the journal holds: the
// greatest dated file's name or, where later, the date of the last Event in
// the numbered pieces an import installs, which have no date in their names.
// Pieces hold their Events in order and sort in replay order, so only the
// tail of the last piece that holds an Event is read.
func greatestJournalDate(journalDir string) (string, bool, error) {
	entries, err := os.ReadDir(journalDir)
	if err != nil {
		return "", false, fmt.Errorf("read journal directory %q: %w", journalDir, err)
	}

	var latest string
	pieces := []string{}
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
			if journalFileName(name) {
				pieces = append(pieces, name)
			}
			continue
		}
		if latest == "" || dateText > latest {
			latest = dateText
		}
	}

	sort.Sort(sort.Reverse(sort.StringSlice(pieces)))
	for _, piece := range pieces {
		date, found, err := lastEventDate(filepath.Join(journalDir, piece))
		if err != nil {
			return "", false, err
		}
		if found {
			if date > latest {
				latest = date
			}
			break
		}
	}

	return latest, latest != "", nil
}

// lastEventDate returns the UTC date of the last Event in a journal file,
// reading backwards from its end in blocks, so a large piece costs only its
// tail. It reports false for a file that holds no Event.
func lastEventDate(path string) (string, bool, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", false, fmt.Errorf("open journal file %q for its latest date: %w", path, err)
	}
	defer func() { _ = file.Close() }()
	info, err := file.Stat()
	if err != nil {
		return "", false, fmt.Errorf("stat journal file %q for its latest date: %w", path, err)
	}
	const block int64 = 64 << 10
	var tail []byte
	for end := info.Size(); end > 0; {
		start := end - min(end, block)
		chunk := make([]byte, end-start)
		if _, err := file.ReadAt(chunk, start); err != nil {
			return "", false, fmt.Errorf("read journal file %q for its latest date: %w", path, err)
		}
		tail, end = append(chunk, tail...), start
		trimmed := bytes.TrimSpace(tail)
		newline := bytes.LastIndexByte(trimmed, '\n')
		if len(trimmed) == 0 || newline < 0 && end > 0 {
			continue
		}
		var event struct {
			OccurredAt time.Time `json:"occurred_at"`
		}
		if err := json.Unmarshal(bytes.TrimSpace(trimmed[newline+1:]), &event); err != nil {
			return "", false, fmt.Errorf("read the last Event of journal file %q: %w", path, err)
		}
		return event.OccurredAt.UTC().Format("2006-01-02"), true, nil
	}
	return "", false, nil
}

func nextSequence(journalPath string) (int, error) {
	date := strings.TrimSuffix(filepath.Base(journalPath), ".jsonl")
	files, err := orderedJournalFiles(filepath.Dir(journalPath))
	if err != nil {
		return 0, err
	}
	highest := 0
	for _, name := range files {
		file, err := os.Open(filepath.Join(filepath.Dir(journalPath), name))
		if err != nil {
			return 0, fmt.Errorf("open journal file %q for sequence scan: %w", name, err)
		}
		scanner := bufio.NewScanner(file)
		scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)
		for scanner.Scan() {
			line := bytes.TrimSpace(scanner.Bytes())
			if len(line) == 0 {
				continue
			}
			var event Event
			if err := json.Unmarshal(line, &event); err != nil {
				_ = file.Close()
				return 0, fmt.Errorf("scan journal file %q: decode Event: %w", name, err)
			}
			parts := strings.Split(event.EventID, "-")
			if len(parts) != 3 || parts[1] != strings.ReplaceAll(date, "-", "") {
				continue
			}
			sequence, err := parseEventSequence(event.EventID)
			if err != nil {
				_ = file.Close()
				return 0, fmt.Errorf("scan journal file %q: %w", name, err)
			}
			if sequence > highest {
				highest = sequence
			}
		}
		if err := scanner.Err(); err != nil {
			_ = file.Close()
			return 0, fmt.Errorf("scan journal file %q: %w", name, err)
		}
		if err := file.Close(); err != nil {
			return 0, fmt.Errorf("close journal file %q after sequence scan: %w", name, err)
		}
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
