package journal

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type aggregateSequenceTracker struct {
	sequences map[string]int
}

func newAggregateSequenceTracker() *aggregateSequenceTracker {
	return &aggregateSequenceTracker{sequences: make(map[string]int)}
}

func (tracker *aggregateSequenceTracker) add(event Event) error {
	if !event.IsV2() {
		return nil
	}
	key := event.Aggregate.Type + "/" + event.Aggregate.ID
	want := tracker.sequences[key] + 1
	if event.Aggregate.Sequence != want {
		return fmt.Errorf("aggregate %s sequence %d, want %d", key, event.Aggregate.Sequence, want)
	}
	tracker.sequences[key] = event.Aggregate.Sequence
	return nil
}

func validateNextAggregateSequence(journalDir string, event Event) error {
	if !event.IsV2() {
		return nil
	}
	sequence, err := NextAggregateSequence(journalDir, event.Aggregate)
	if err != nil {
		return err
	}
	if event.Aggregate.Sequence != sequence {
		return fmt.Errorf("aggregate %s/%s sequence %d, want %d", event.Aggregate.Type, event.Aggregate.ID, event.Aggregate.Sequence, sequence)
	}
	return nil
}

// NextAggregateSequence returns the next Event-v2 sequence for one aggregate.
// Callers hold the journal write lock while deriving and appending the Event.
func NextAggregateSequence(journalDir string, aggregate Aggregate) (int, error) {
	files, err := orderedJournalFiles(journalDir)
	if err != nil {
		return 0, err
	}
	tracker := newAggregateSequenceTracker()
	for _, fileName := range files {
		if err := addAggregateSequencesFromFile(tracker, filepath.Join(journalDir, fileName)); err != nil {
			return 0, err
		}
	}
	return tracker.sequences[aggregate.Type+"/"+aggregate.ID] + 1, nil
}

func addAggregateSequencesFromFile(tracker *aggregateSequenceTracker, path string) error {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("open journal file %q: %w", path, err)
	}
	defer closeReplayFile(file)

	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	for scanner.Scan() {
		line := bytesTrimSpace(scanner.Bytes())
		if len(line) == 0 {
			continue
		}
		var stored Event
		if err := json.Unmarshal(line, &stored); err != nil {
			return fmt.Errorf("scan journal file %q: decode Event: %w", path, err)
		}
		if err := tracker.add(stored); err != nil {
			return err
		}
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("scan journal file %q: %w", path, err)
	}
	return nil
}
