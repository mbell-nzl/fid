package journal

import (
	"bytes"
	"encoding/json"
	"errors"
	"os"
	"reflect"
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
		"header":     bytes.Replace(raw, []byte(`"event_schema":2`), []byte(`"event_schema":1,"event_schema":2`), 1),
		"nested key": bytes.Replace(raw, []byte(`"id":"FID-a"`), []byte(`"id":"FID-b","id":"FID-a"`), 1),
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			var decoded Event
			if err := json.Unmarshal(malformed, &decoded); err == nil {
				t.Error("malformed duplicate keys decoded")
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
					if err := w.replayStored(t, mode); !errors.As(err, &corrupt) {
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
