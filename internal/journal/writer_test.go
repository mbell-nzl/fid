package journal

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestAppendEventWritesOneJSONLLineWithFirstSequence(t *testing.T) {
	workspace := t.TempDir()
	mustMkdirAll(t, filepath.Join(workspace, ".fid", "journal"))

	fixedNow := time.Date(2026, time.June, 5, 14, 23, 11, 0, time.UTC)
	event := Event{
		EventSchema: EventSchemaV1,
		EventType:   "issue.created",
		Actor: Actor{
			Role: "developer",
		},
		Command: Command{
			Argv:       []string{"fid", "issue", "create", "--by", "developer"},
			CWD:        ".",
			FIDVersion: "0.1.0",
		},
		Effects: []Effect{
			{
				Op:    "insert",
				Table: "issues",
				PK: map[string]any{
					"id": "FID-a1b2c3",
				},
				Values: map[string]any{
					"title": "Implement append writer",
				},
			},
		},
	}

	var appended Event
	clockReads := 0
	err := withLockOptions(workspace, LockOwner{
		Command:   "fid issue create --by developer",
		ActorRole: "developer",
	}, defaultLockOptions(), func(lock *WriteLock) error {
		var appendErr error
		appended, appendErr = appendEventWithClock(lock, event, func() time.Time {
			clockReads++
			return fixedNow
		})
		return appendErr
	})
	if err != nil {
		t.Fatalf("withLockOptions() error = %v", err)
	}

	if clockReads != 1 {
		t.Fatalf("clock reads = %d, want 1", clockReads)
	}
	if appended.EventID != "EVT-20260605-000001" {
		t.Fatalf("appended event id = %q, want %q", appended.EventID, "EVT-20260605-000001")
	}
	if !appended.OccurredAt.Equal(fixedNow) {
		t.Fatalf("appended occurred_at = %s, want %s", appended.OccurredAt, fixedNow)
	}

	journalPath := filepath.Join(workspace, ".fid", "journal", "2026-06-05.jsonl")
	raw, err := os.ReadFile(journalPath)
	if err != nil {
		t.Fatalf("ReadFile(%q) error = %v", journalPath, err)
	}

	lines := bytes.Split(bytes.TrimSpace(raw), []byte("\n"))
	if len(lines) != 1 {
		t.Fatalf("journal lines = %d, want 1", len(lines))
	}

	var stored Event
	if err := json.Unmarshal(lines[0], &stored); err != nil {
		t.Fatalf("Unmarshal() error = %v", err)
	}

	if stored.EventID != appended.EventID {
		t.Fatalf("stored event id = %q, want %q", stored.EventID, appended.EventID)
	}
	if stored.EventSchema != EventSchemaV1 {
		t.Fatalf("stored event schema = %d, want %d", stored.EventSchema, EventSchemaV1)
	}
}

func TestAppendEventIncrementsSequenceWithinUTCJournalFile(t *testing.T) {
	workspace := t.TempDir()
	mustMkdirAll(t, filepath.Join(workspace, ".fid", "journal"))

	fixedNow := time.Date(2026, time.June, 5, 14, 23, 11, 0, time.UTC)
	first, err := appendTestEvent(workspace, fixedNow, "Implement append writer")
	if err != nil {
		t.Fatalf("append first event error = %v", err)
	}
	second, err := appendTestEvent(workspace, fixedNow, "Implement lock writer")
	if err != nil {
		t.Fatalf("append second event error = %v", err)
	}

	if first.EventID != "EVT-20260605-000001" {
		t.Fatalf("first event id = %q, want %q", first.EventID, "EVT-20260605-000001")
	}
	if second.EventID != "EVT-20260605-000002" {
		t.Fatalf("second event id = %q, want %q", second.EventID, "EVT-20260605-000002")
	}
}

func TestAppendEventRejectsClockSkewWithoutAppending(t *testing.T) {
	workspace := t.TempDir()
	journalDir := filepath.Join(workspace, ".fid", "journal")
	mustMkdirAll(t, journalDir)

	laterPath := filepath.Join(journalDir, "2026-06-06.jsonl")
	if err := os.WriteFile(laterPath, []byte{}, 0o644); err != nil {
		t.Fatalf("WriteFile(%q) error = %v", laterPath, err)
	}

	fixedNow := time.Date(2026, time.June, 5, 9, 0, 0, 0, time.UTC)
	err := withLockOptions(workspace, LockOwner{
		Command:   "fid issue create --by developer",
		ActorRole: "developer",
	}, defaultLockOptions(), func(lock *WriteLock) error {
		_, appendErr := appendEventWithClock(lock, testAppendEvent("Clock skew rejected"), func() time.Time {
			return fixedNow
		})
		return appendErr
	})
	if err == nil {
		t.Fatal("withLockOptions() error = nil, want clock skew error")
	}

	const want = "system clock is behind latest journal date 2026-06-06; refusing to append out of order"
	if err.Error() != want {
		t.Fatalf("clock skew error = %q, want %q", err, want)
	}

	earlierPath := filepath.Join(journalDir, "2026-06-05.jsonl")
	if _, statErr := os.Stat(earlierPath); !os.IsNotExist(statErr) {
		t.Fatalf("earlier journal path stat error = %v, want not exist", statErr)
	}

	laterRaw, err := os.ReadFile(laterPath)
	if err != nil {
		t.Fatalf("ReadFile(%q) error = %v", laterPath, err)
	}
	if len(laterRaw) != 0 {
		t.Fatalf("later journal bytes = %q, want unchanged empty file", laterRaw)
	}
}

func TestConcurrentWritesAppendDistinctEventsWithMonotonicSequences(t *testing.T) {
	workspace := t.TempDir()
	mustMkdirAll(t, filepath.Join(workspace, ".fid", "journal"))

	fixedNow := time.Date(2026, time.June, 5, 14, 23, 11, 0, time.UTC)
	start := make(chan struct{})
	results := make(chan Event, 2)
	errs := make(chan error, 2)

	var wg sync.WaitGroup
	for _, title := range []string{"Implement append writer", "Implement lock writer"} {
		wg.Add(1)
		go func(title string) {
			defer wg.Done()
			<-start

			err := withLockOptions(workspace, LockOwner{
				Command:   "fid issue create --by developer",
				ActorRole: "developer",
			}, defaultLockOptions(), func(lock *WriteLock) error {
				appended, appendErr := appendEventWithClock(lock, testAppendEvent(title), func() time.Time {
					return fixedNow
				})
				if appendErr == nil {
					results <- appended
				}
				return appendErr
			})
			errs <- err
		}(title)
	}

	close(start)
	wg.Wait()
	close(results)
	close(errs)

	for err := range errs {
		if err != nil {
			t.Fatalf("concurrent withLockOptions() error = %v", err)
		}
	}

	var eventIDs []string
	for result := range results {
		eventIDs = append(eventIDs, result.EventID)
	}
	if len(eventIDs) != 2 {
		t.Fatalf("concurrent results = %d, want 2", len(eventIDs))
	}
	if eventIDs[0] == eventIDs[1] {
		t.Fatalf("event ids = %v, want distinct ids", eventIDs)
	}

	stored, err := readStoredEvents(filepath.Join(workspace, ".fid", "journal", "2026-06-05.jsonl"))
	if err != nil {
		t.Fatalf("readStoredEvents() error = %v", err)
	}
	if len(stored) != 2 {
		t.Fatalf("stored events = %d, want 2", len(stored))
	}

	got := []string{stored[0].EventID, stored[1].EventID}
	want := []string{"EVT-20260605-000001", "EVT-20260605-000002"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("stored event ids = %v, want %v", got, want)
	}
}

func appendTestEvent(workspace string, fixedNow time.Time, title string) (Event, error) {
	var appended Event
	err := withLockOptions(workspace, LockOwner{
		Command:   "fid issue create --by developer",
		ActorRole: "developer",
	}, defaultLockOptions(), func(lock *WriteLock) error {
		var appendErr error
		appended, appendErr = appendEventWithClock(lock, testAppendEvent(title), func() time.Time {
			return fixedNow
		})
		return appendErr
	})

	return appended, err
}

func testAppendEvent(title string) Event {
	return Event{
		EventSchema: EventSchemaV1,
		EventType:   "issue.created",
		Actor: Actor{
			Role: "developer",
		},
		Command: Command{
			Argv:       []string{"fid", "issue", "create", "--by", "developer"},
			CWD:        ".",
			FIDVersion: "0.1.0",
		},
		Effects: []Effect{
			{
				Op:    "insert",
				Table: "issues",
				PK: map[string]any{
					"id": "FID-a1b2c3",
				},
				Values: map[string]any{
					"title": title,
				},
			},
		},
	}
}

func readStoredEvents(path string) ([]Event, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(strings.TrimSpace(string(raw)), "\n")
	events := make([]Event, 0, len(lines))
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}

		var event Event
		if err := json.Unmarshal([]byte(line), &event); err != nil {
			return nil, err
		}
		events = append(events, event)
	}

	return events, nil
}

func mustMkdirAll(t *testing.T, path string) {
	t.Helper()

	if err := os.MkdirAll(path, 0o755); err != nil {
		t.Fatalf("MkdirAll(%q) error = %v", path, err)
	}
}
