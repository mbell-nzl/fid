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
	t.Parallel()
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
		Command: "fid issue create --by developer",
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

func TestAppendEventWrapperUsesRealClockPath(t *testing.T) {
	t.Parallel()
	workspace := t.TempDir()
	mustMkdirAll(t, filepath.Join(workspace, ".fid", "journal"))

	var appended Event
	err := withLockOptions(workspace, LockOwner{
		Command: "fid issue create --by developer",
	}, defaultLockOptions(), func(lock *WriteLock) error {
		var appendErr error
		appended, appendErr = AppendEvent(lock, testAppendEvent("Append wrapper"))
		return appendErr
	})
	if err != nil {
		t.Fatalf("AppendEvent() error = %v", err)
	}
	if appended.EventID == "" {
		t.Fatal("AppendEvent() event id = empty, want generated id")
	}
	if appended.OccurredAt.IsZero() {
		t.Fatal("AppendEvent() occurred_at = zero, want current timestamp")
	}
}

func TestPrepareEventAssignsNextIDWithoutAppendingJournalBytes(t *testing.T) {
	t.Parallel()
	workspace := t.TempDir()
	journalDir := filepath.Join(workspace, ".fid", "journal")
	mustMkdirAll(t, journalDir)

	journalPath := filepath.Join(journalDir, "2026-06-05.jsonl")
	writeWriterEvents(t, journalPath,
		writerRoleEvent("EVT-20260605-000001"),
		writerRoleEvent("EVT-20260605-000002"),
	)
	before, err := os.ReadFile(journalPath)
	if err != nil {
		t.Fatalf("ReadFile(%q) before preparation error = %v", journalPath, err)
	}

	input := testAppendEvent("Prepare without append")
	input.OccurredAt = time.Date(2026, time.June, 5, 14, 23, 11, 0, time.UTC)

	var prepared Event
	err = WithLock(workspace, LockOwner{Command: "fid issue create --by developer"}, func(lock *WriteLock) error {
		var prepareErr error
		prepared, prepareErr = PrepareEvent(lock, input)
		return prepareErr
	})
	if err != nil {
		t.Fatalf("PrepareEvent() error = %v", err)
	}

	if prepared.EventID != "EVT-20260605-000003" {
		t.Fatalf("prepared event id = %q, want %q", prepared.EventID, "EVT-20260605-000003")
	}
	if !prepared.OccurredAt.Equal(input.OccurredAt) {
		t.Fatalf("prepared occurred_at = %s, want %s", prepared.OccurredAt, input.OccurredAt)
	}
	if !reflect.DeepEqual(prepared.Command, input.Command) {
		t.Fatalf("prepared command = %#v, want %#v", prepared.Command, input.Command)
	}

	after, err := os.ReadFile(journalPath)
	if err != nil {
		t.Fatalf("ReadFile(%q) after preparation error = %v", journalPath, err)
	}
	if !bytes.Equal(after, before) {
		t.Fatalf("PrepareEvent() changed journal bytes = %q, want %q", after, before)
	}
}

func TestAppendEventIncrementsSequenceWithinUTCJournalFile(t *testing.T) {
	t.Parallel()
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
	t.Parallel()
	workspace := t.TempDir()
	journalDir := filepath.Join(workspace, ".fid", "journal")
	mustMkdirAll(t, journalDir)

	laterPath := filepath.Join(journalDir, "2026-06-06.jsonl")
	if err := os.WriteFile(laterPath, []byte{}, 0o644); err != nil {
		t.Fatalf("WriteFile(%q) error = %v", laterPath, err)
	}

	fixedNow := time.Date(2026, time.June, 5, 9, 0, 0, 0, time.UTC)
	err := withLockOptions(workspace, LockOwner{
		Command: "fid issue create --by developer",
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
	t.Parallel()
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
				Command: "fid issue create --by developer",
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

func TestWriterHelpersParseEventIDsAndDates(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		eventID string
		want    int
		wantErr string
	}{
		{name: "valid event id", eventID: "EVT-20260605-000123", want: 123},
		{name: "wrong parts", eventID: "EVT-20260605", wantErr: `parse event id "EVT-20260605": want EVT-<date>-<sequence>`},
		{name: "wrong prefix", eventID: "BAD-20260605-000123", wantErr: `parse event id "BAD-20260605-000123": want EVT prefix`},
		{name: "bad date", eventID: "EVT-20261305-000123", wantErr: `parse event id "EVT-20261305-000123" date:`},
		{name: "bad sequence", eventID: "EVT-20260605-nope", wantErr: `parse event id "EVT-20260605-nope" sequence:`},
		{name: "non-positive sequence", eventID: "EVT-20260605-000000", wantErr: `parse event id "EVT-20260605-000000" sequence: must be positive`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseEventSequence(tt.eventID)
			if tt.wantErr == "" {
				if err != nil {
					t.Fatalf("parseEventSequence() error = %v", err)
				}
				if got != tt.want {
					t.Fatalf("parseEventSequence() = %d, want %d", got, tt.want)
				}
				return
			}

			if err == nil {
				t.Fatalf("parseEventSequence() error = nil, want %q", tt.wantErr)
			}
			if !strings.Contains(err.Error(), tt.wantErr) {
				t.Fatalf("parseEventSequence() error = %q, want substring %q", err.Error(), tt.wantErr)
			}
		})
	}
}

func TestWriterHelpersDetermineLatestDateAndNextSequence(t *testing.T) {
	t.Parallel()
	journalDir := filepath.Join(t.TempDir(), ".fid", "journal")
	mustMkdirAll(t, journalDir)

	for _, name := range []string{"2026-06-04.jsonl", "2026-06-05.jsonl", "notes.txt", "invalid.jsonl"} {
		path := filepath.Join(journalDir, name)
		if err := os.WriteFile(path, []byte{}, 0o644); err != nil {
			t.Fatalf("WriteFile(%q) error = %v", path, err)
		}
	}

	latest, found, err := greatestJournalDate(journalDir)
	if err != nil {
		t.Fatalf("greatestJournalDate() error = %v", err)
	}
	if !found {
		t.Fatal("greatestJournalDate() found = false, want true")
	}
	if latest != "2026-06-05" {
		t.Fatalf("greatestJournalDate() latest = %q, want %q", latest, "2026-06-05")
	}

	if next, err := nextSequence(filepath.Join(journalDir, "2026-06-06.jsonl")); err != nil {
		t.Fatalf("nextSequence(missing) error = %v", err)
	} else if next != 1 {
		t.Fatalf("nextSequence(missing) = %d, want 1", next)
	}

	path := filepath.Join(journalDir, "2026-06-05.jsonl")
	writeWriterEvents(t, path,
		writerRoleEvent("EVT-20260605-000001"),
		writerRoleEvent("EVT-20260605-000002"),
	)

	next, err := nextSequence(path)
	if err != nil {
		t.Fatalf("nextSequence(existing) error = %v", err)
	}
	if next != 3 {
		t.Fatalf("nextSequence(existing) = %d, want 3", next)
	}
}

func appendTestEvent(workspace string, fixedNow time.Time, title string) (Event, error) {
	var appended Event
	err := withLockOptions(workspace, LockOwner{
		Command: "fid issue create --by developer",
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

func writeWriterEvents(t *testing.T, path string, events ...Event) {
	t.Helper()

	var payload []byte
	for _, event := range events {
		line, err := MarshalCanonicalJSON(event)
		if err != nil {
			t.Fatalf("MarshalCanonicalJSON() error = %v", err)
		}
		payload = append(payload, line...)
		payload = append(payload, '\n')
	}

	if err := os.WriteFile(path, payload, 0o644); err != nil {
		t.Fatalf("WriteFile(%q) error = %v", path, err)
	}
}

func writerRoleEvent(eventID string) Event {
	return Event{
		EventID:     eventID,
		EventSchema: EventSchemaV1,
		EventType:   "project.initialized",
		OccurredAt:  time.Date(2026, time.June, 5, 14, 23, 11, 0, time.UTC),
		Actor: Actor{
			Role: "human",
		},
		Command: Command{
			Argv:       []string{"fid", "init"},
			CWD:        ".",
			FIDVersion: "0.1.0",
		},
		Effects: []Effect{
			{
				Op:    "insert",
				Table: "labels",
				PK: map[string]any{
					"name": eventID,
				},
				Values: map[string]any{
					"description": eventID,
					"colour":      "#00ff00",
				},
			},
		},
	}
}

func mustMkdirAll(t *testing.T, path string) {
	t.Helper()

	if err := os.MkdirAll(path, 0o755); err != nil {
		t.Fatalf("MkdirAll(%q) error = %v", path, err)
	}
}
