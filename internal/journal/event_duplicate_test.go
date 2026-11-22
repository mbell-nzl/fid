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
	legacy := e
	legacy.EventSchema = EventSchemaV1
	legacy.Actor = Actor{Role: "developer"}
	legacy.Command = Command{Argv: []string{"fid"}, CWD: ".", FIDVersion: "test"}
	legacyRaw, err := MarshalCanonicalJSON(legacy)
	if err != nil {
		t.Fatal(err)
	}
	legacyMalformed := bytes.Replace(legacyRaw, []byte(`"cwd":"."`), []byte(`"cwd":"/","CWD":"."`), 1)
	if bytes.Equal(legacyRaw, legacyMalformed) {
		t.Fatal("probe did not change the Event-v1 command")
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
		"exact map duplicate":  bytes.Replace(raw, []byte(`"details":null`), []byte(`"details":{"a":1,"a":2}`), 1),
		"nested map duplicate": bytes.Replace(raw, []byte(`"details":null`), []byte(`"details":[{"a":1,"a":2}]`), 1),
		"Event-v1 command":     legacyMalformed,
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
	// A Go map cannot hold an exact duplicate. Embedded JSON text can, and
	// must be refused before append as well as through every replay path.
	// Distinct fold-equivalent map keys are positive controls in
	// TestStructuredEventCaseSensitiveMapKeys.
	for name, tc := range map[string]struct{ column, raw, family string }{
		"ASCII map":                     {"request_field_digests_json", `{"a":2,"a":1}`, ""},
		"nested in array":               {"request_field_digests_json", `{"nested":[{"a":2,"a":1}]}`, ""},
		"long s":                        {"request_field_digests_json", `{"ſ":2,"ſ":1}`, ""},
		"Kelvin sign":                   {"request_field_digests_json", `{"K":2,"K":1}`, ""},
		"Greek sigma":                   {"request_field_digests_json", `{"σ":2,"σ":1}`, ""},
		"retry result":                  {"result_json", `{"replayed":true,"replayed":false}`, ""},
		"nested retry result":           {"result_json", `{"replayed":false,"changes":[{"field":"lane","from":"standard","to":"standard","to":"assured"}]}`, ""},
		"criterion retry result":        {"result_json", `{"replayed":true,"replayed":false}`, "issue.ac-added"},
		"criterion nested retry result": {"result_json", `{"replayed":false,"criterion":{"text":"overwritten","text":"original"}}`, "issue.ac-added"},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			e := admissionLaneRaise()
			if tc.family == "issue.ac-added" {
				e.EventType, e.CommandName = tc.family, "issue.ac.add"
				e.Effects[0] = Effect{Op: "insert", Table: "acceptance_criteria", PK: map[string]any{"id": "AC-one"}, Values: map[string]any{"issue_id": e.Aggregate.ID, "ordinal": 1, "text": "Works"}}
			}
			request := IdempotencyRequest{OperationScope: e.CommandName + ":project-1", IdempotencyKey: "duplicate", RequestDigest: "sha256:" + strings.Repeat("b", 64)}
			e.Idempotency = &EventIdempotency{Key: request.IdempotencyKey, Scope: request.OperationScope, RequestDigest: request.RequestDigest}
			result := idempotencyResultEffect(request, e, e.Aggregate.ID, `{"replayed":false}`)
			result.Values[tc.column] = tc.raw
			e.Effects = append(e.Effects, result)
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
