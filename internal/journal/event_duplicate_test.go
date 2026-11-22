package journal

import (
	"bytes"
	"encoding/json"
	"errors"
	"os"
	"reflect"
	"strings"
	"testing"
)

func TestMalformedEventDuplicateJSONKeys(t *testing.T) {
	t.Parallel()
	e := admissionLaneRaise()
	raw, err := MarshalCanonicalJSON(e)
	if err != nil {
		t.Fatal(err)
	}
	for name, malformed := range map[string][]byte{
		"header":               bytes.Replace(raw, []byte(`"event_schema":2`), []byte(`"event_schema":1,"event_schema":2`), 1),
		"nested key":           bytes.Replace(raw, []byte(`"id":"FID-a"`), []byte(`"id":"FID-b","id":"FID-a"`), 1),
		"case variant family":  bytes.Replace(raw, []byte(`"event_type":"issue.updated"`), []byte(`"event_type":"label.added","EVENT_TYPE":"issue.updated"`), 1),
		"case variant command": bytes.Replace(raw, []byte(`"command":"issue.update"`), []byte(`"command":"x.y","Command":"issue.update"`), 1),
		"case variant nested":  bytes.Replace(raw, []byte(`"id":"FID-a"`), []byte(`"id":"FID-b","ID":"FID-a"`), 1),
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			if bytes.Equal(raw, malformed) {
				t.Fatal("probe did not change the Event JSON")
			}
			var decoded Event
			if err := json.Unmarshal(malformed, &decoded); err == nil || !strings.Contains(err.Error(), "duplicate key") {
				t.Errorf("decode = %v, want duplicate-key refusal", err)
			}
			for _, mode := range []string{"replay", "catch-up", "freshen", "rebuild"} {
				t.Run(mode, func(t *testing.T) {
					w := newAdmissionWorkspace(t, admissionPrefix(), mode)
					before := stateProjection(t, w.root)
					stored := append(append(bytes.Clone(w.prefix), malformed...), '\n')
					if err := os.WriteFile(w.file, stored, 0o644); err != nil {
						t.Fatal(err)
					}
					var corrupt *CorruptJournalError
					if err := w.replayStored(t, mode); !errors.As(err, &corrupt) || !strings.Contains(err.Error(), "duplicate key") {
						t.Fatalf("duplicate keys = %v, want CorruptJournalError", err)
					}
					if mode != "replay" {
						after := stateProjection(t, w.root)
						if !reflect.DeepEqual(before, after) {
							t.Fatal("refusal changed projection")
						}
					}
					got, err := os.ReadFile(w.file)
					if err != nil || !bytes.Equal(stored, got) {
						t.Fatalf("refusal changed journal: %v", err)
					}
					if err := os.WriteFile(w.file, w.prefix, 0o644); err != nil {
						t.Fatal(err)
					}
					if _, err := EnsureProjectionFresh(w.root, testProjectMigrationsDir(t)); err != nil {
						t.Fatal(err)
					}
					if err := w.admit(t, "commit", admissionLaneRaise()); err != nil {
						t.Fatal(err)
					}
				})
			}
		})
	}
}
