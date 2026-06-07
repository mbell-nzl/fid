package journal

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
	"time"
)

func TestMarshalCanonicalJSONProducesStableV1EventBytes(t *testing.T) {
	t.Parallel()

	event := testV1Event()

	got, err := MarshalCanonicalJSON(event)
	if err != nil {
		t.Fatalf("MarshalCanonicalJSON() error = %v", err)
	}

	want := `{"actor":{"actor":null,"role":"developer"},"command":{"argv":["fid","issue","transition","FID-a1b2c3","--to","in_progress","--by","developer"],"cwd":".","fid_version":"0.1.0"},"effects":[{"op":"update","pk":{"id":"FID-a1b2c3"},"preconditions":{"status":"todo","workflow_id":"task-default"},"prior":{"status":"todo"},"table":"issues","values":{"status":"in_progress","updated_at":"2026-05-29T03:12:45Z"}},{"op":"activity","pk":{"id":"ACT-d4e5f6"},"table":"activity_log","values":{"actor":"developer","details":{"from":"todo","to":"in_progress"},"event":"status_changed","journal_event_id":"EVT-20260529-000001","occurred_at":"2026-05-29T03:12:45Z","ref_id":"FID-a1b2c3","ref_type":"issue"}}],"event_id":"EVT-20260529-000001","event_schema":1,"event_type":"issue.transitioned","occurred_at":"2026-05-29T03:12:45Z"}`

	if string(got) != want {
		t.Fatalf("MarshalCanonicalJSON() = %s, want %s", got, want)
	}
}

func TestHashEventReturnsStableDigest(t *testing.T) {
	t.Parallel()

	event := testV1Event()

	got, err := HashEvent(event)
	if err != nil {
		t.Fatalf("HashEvent() error = %v", err)
	}

	const want = "sha256:e24e6bd01efe5a21417b6eae9532d7014ac4e5508a68cce03976c32b7c7db371"
	if got != want {
		t.Fatalf("HashEvent() = %q, want %q", got, want)
	}
}

func TestV1GoldenFixtureRoundTrips(t *testing.T) {
	t.Parallel()

	golden, err := os.ReadFile(testGoldenEventPath(t))
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}
	golden = bytes.TrimSpace(golden)

	var event Event
	if err := json.Unmarshal(golden, &event); err != nil {
		t.Fatalf("Unmarshal() error = %v", err)
	}

	if !reflect.DeepEqual(event, testV1Event()) {
		t.Fatalf("fixture event = %#v, want %#v", event, testV1Event())
	}

	roundTripped, err := MarshalCanonicalJSON(event)
	if err != nil {
		t.Fatalf("MarshalCanonicalJSON() error = %v", err)
	}

	if !bytes.Equal(roundTripped, golden) {
		t.Fatalf("round-tripped bytes = %s, want %s", roundTripped, golden)
	}
}

func testV1Event() Event {
	return Event{
		EventID:     "EVT-20260529-000001",
		EventSchema: EventSchemaV1,
		EventType:   "issue.transitioned",
		OccurredAt:  time.Date(2026, time.May, 29, 3, 12, 45, 0, time.UTC),
		Actor: Actor{
			Role:  "developer",
			Actor: nil,
		},
		Command: Command{
			Argv:       []string{"fid", "issue", "transition", "FID-a1b2c3", "--to", "in_progress", "--by", "developer"},
			CWD:        ".",
			FIDVersion: "0.1.0",
		},
		Effects: []Effect{
			{
				Op:    "update",
				Table: "issues",
				PK: map[string]any{
					"id": "FID-a1b2c3",
				},
				Preconditions: map[string]any{
					"status":      "todo",
					"workflow_id": "task-default",
				},
				Prior: map[string]any{
					"status": "todo",
				},
				Values: map[string]any{
					"status":     "in_progress",
					"updated_at": "2026-05-29T03:12:45Z",
				},
			},
			{
				Op:    "activity",
				Table: "activity_log",
				PK: map[string]any{
					"id": "ACT-d4e5f6",
				},
				Values: map[string]any{
					"journal_event_id": "EVT-20260529-000001",
					"ref_type":         "issue",
					"ref_id":           "FID-a1b2c3",
					"event":            "status_changed",
					"details": map[string]any{
						"from": "todo",
						"to":   "in_progress",
					},
					"actor":       "developer",
					"occurred_at": "2026-05-29T03:12:45Z",
				},
			},
		},
	}
}

func testGoldenEventPath(t *testing.T) string {
	t.Helper()

	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller(0) failed")
	}

	return filepath.Clean(filepath.Join(filepath.Dir(currentFile), "..", "..", "testdata", "journal", "event-v1.golden.json"))
}
