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
		"unicode header":       bytes.Replace(raw, []byte(`"event_schema":2`), []byte(`"event_schema":1,"event_ſchema":2`), 1),
		"unicode effects":      bytes.Replace(raw, []byte(`"effects":[`), []byte(`"effects":[],"effectſ":[`), 1),
		"unicode Effect":       bytes.Replace(raw, []byte(`"values":{`), []byte(`"values":{},"valueſ":{`), 1),
		"unicode aggregate":    bytes.Replace(raw, []byte(`"sequence":2`), []byte(`"sequence":1,"ſequence":2`), 1),
		"escaped unicode":      bytes.Replace(raw, []byte(`"event_schema":2`), []byte(`"event_schema":1,"event_\u017fchema":2`), 1),
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
					after := stateProjection(t, w.root)
					if !reflect.DeepEqual(before, after) {
						t.Fatal("refusal changed projection")
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

func TestStructuredEventDuplicateJSONKeys(t *testing.T) {
	t.Parallel()
	for name, details := range map[string]any{
		"case variant map": map[string]any{"A": 2, "a": 1},
		"nested in array":  []any{map[string]any{"A": 2, "a": 1}},
		"long s":           map[string]any{"s": 2, "ſ": 1},
		"Kelvin sign":      map[string]any{"K": 2, "K": 1},
		"Greek sigma":      map[string]any{"σ": 2, "ς": 1},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			e := admissionLaneRaise()
			e.Effects[2].Values["details"] = details
			for _, mode := range admissionModes(e) {
				t.Run(mode, func(t *testing.T) {
					w := newAdmissionWorkspace(t, admissionPrefix(), mode)
					before := stateProjection(t, w.root)
					err := w.admit(t, mode, e)
					if err == nil || !strings.Contains(err.Error(), "duplicate key") {
						t.Fatalf("structured duplicate keys = %v, want duplicate-key refusal", err)
					}
					live := mode == "commit" || mode == "dry-run" || mode == "project"
					var corrupt *CorruptJournalError
					if !live && !errors.As(err, &corrupt) {
						t.Fatalf("replay refusal = %v, want CorruptJournalError", err)
					}
					if !reflect.DeepEqual(before, stateProjection(t, w.root)) {
						t.Fatal("refusal changed projection")
					}
					want := w.prefix
					if !live {
						want = w.journalWith(t, e)
					}
					got, err := os.ReadFile(w.file)
					if err != nil || !bytes.Equal(want, got) {
						t.Fatalf("refusal changed journal: %v", err)
					}
					if err := os.WriteFile(w.file, w.prefix, 0o644); err != nil {
						t.Fatal(err)
					}
					if _, err := EnsureProjectionFresh(w.root, testProjectMigrationsDir(t)); err != nil {
						t.Fatal(err)
					}
					if err := w.admit(t, "commit", admissionLaneRaise()); err != nil {
						t.Fatalf("valid write after refusal: %v", err)
					}
				})
			}
		})
	}
}

func TestStructuredEventDistinctUnicodeKeys(t *testing.T) {
	t.Parallel()
	e := admissionLaneRaise()
	// Unicode simple folding keeps dotted/dotless I distinct from ASCII i.
	// Lowercasing would incorrectly treat dotted I as the same key.
	e.Effects[2].Values["details"] = map[string]any{"i": 1, "İ": 2, "ı": 3, "é": 4}
	for _, mode := range admissionModes(e) {
		t.Run(mode, func(t *testing.T) {
			t.Parallel()
			w := newAdmissionWorkspace(t, admissionPrefix(), mode)
			if err := w.admit(t, mode, e); err != nil {
				t.Fatalf("distinct Unicode keys refused: %v", err)
			}
		})
	}
}
