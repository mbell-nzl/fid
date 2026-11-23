package journal

import (
	"encoding/json"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/mbell-nzl/fid/internal/config"
	"github.com/mbell-nzl/fid/internal/handoffrecord"
	"github.com/mbell-nzl/fid/internal/policy"
)

func TestStructuredEventCaseSensitiveMapKeys(t *testing.T) {
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
					if err := w.admit(t, mode, e); err != nil {
						t.Fatalf("distinct map keys refused: %v", err)
					}
					if mode == "dry-run" {
						return
					}
					conn := openReplayDB(t, filepath.Join(config.LedgerDir(w.root), "fid.sqlite"))
					defer closeProjectDB(t, conn)
					var raw string
					if err := conn.QueryRow(`SELECT details FROM activity_log WHERE id='ACT-update'`).Scan(&raw); err != nil {
						t.Fatal(err)
					}
					want, err := MarshalCanonicalJSON(details)
					if err != nil || raw != string(want) {
						t.Fatalf("projected map = %s, want %s: %v", raw, want, err)
					}
				})
			}
		})
	}
}

func TestJSONMapKeysMatchExactly(t *testing.T) {
	t.Parallel()
	for _, raw := range []string{
		`{"build":1,"Build":2}`,
		`{"nested":[{"s":1,"ſ":2,"S":3}]}`,
		`{"Θ":1,"θ":2,"ϑ":3,"ϴ":4}`,
	} {
		t.Run(raw, func(t *testing.T) {
			t.Parallel()
			if _, err := decodeCurrentJSON(raw); err != nil {
				t.Fatalf("distinct map keys refused: %v", err)
			}
		})
	}
}

func TestTypedJSONKeysFollowDestination(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name, raw string
		target    func() any
		duplicate bool
	}{
		{"policy checks", `{"checks":{"build":{"run":[]},"Build":{"run":[]}}}`, func() any { return new(policy.Effective) }, false},
		{"policy lanes", `{"lanes":{"standard":{},"Standard":{}}}`, func() any { return new(policy.Effective) }, false},
		{"hand-off digests", `{"definition_digests":{"build":"one","Build":"two"}}`, func() any { return new(handoffrecord.Record) }, false},
		{"policy struct alias", `{"checks":{"build":{"run":[],"RUN":[]}}}`, func() any { return new(policy.Effective) }, true},
		{"aliased parent field", `{"CHECKS":{"build":{"run":[],"RUN":[]}}}`, func() any { return new(policy.Effective) }, true},
		{"hand-off embedded field", `{"result_ids":[],"reſult_ids":[]}`, func() any { return new(handoffrecord.Record) }, true},
		{"hand-off embedded object", `{"floor":{"paths":[],"pathſ":[]}}`, func() any { return new(handoffrecord.Record) }, true},
		{"exact map duplicate", `{"definition_digests":{"build":"one","build":"two"}}`, func() any { return new(handoffrecord.Record) }, true},
		{"Event-v1 command", `{"command":{"cwd":".","CWD":"/"}}`, func() any { return new(eventV1Wire) }, true},
		{"actor", `{"actor":{"session_ref":null,"ſession_ref":null}}`, func() any { return new(eventV2Wire) }, true},
		{"idempotency", `{"idempotency":{"scope":"one","ſcope":"two"}}`, func() any { return new(eventV2Wire) }, true},
		{"repository", `{"repository":{"capture_source":"one","capture_ſource":"two"}}`, func() any { return new(eventV2Wire) }, true},
		{"Effect struct", `{"effects":[{"values":{},"valueſ":{}}]}`, func() any { return new(eventV2Wire) }, true},
		{"Effect maps", `{"effects":[{"pk":{"key":1,"Key":2},"values":{"key":1,"Key":2},"prior":{"key":1,"Key":2},"preconditions":{"key":1,"Key":2}}]}`, func() any { return new(eventV2Wire) }, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			err := decodeCheckJSON(tc.raw, tc.target())
			if tc.duplicate {
				if err == nil || !strings.Contains(err.Error(), "duplicate key") {
					t.Fatalf("decode = %v, want duplicate-key refusal", err)
				}
			} else if err != nil {
				t.Fatalf("distinct map keys refused: %v", err)
			}
		})
	}
}

func TestPolicyPayloadStructAliases(t *testing.T) {
	t.Parallel()
	e := snapshotEvent(t, 1)
	raw := e.Effects[0].Values["effective_json"].(string)
	malformed := strings.Replace(raw, `"rules_version":3`, `"rules_version":1,"ruleſ_version":3`, 1)
	if malformed == raw {
		t.Fatal("probe did not change policy JSON")
	}
	if err := validatePolicyPayload(malformed, e.Effects[0].Values["digest"]); err == nil || !strings.Contains(err.Error(), "duplicate key") {
		t.Fatalf("policy payload = %v, want duplicate-key refusal", err)
	}
}

func TestFoldJSONKeyMatchesStructFields(t *testing.T) {
	t.Parallel()
	for _, group := range []string{"Ssſ", "KkK", "Σσς", "Θθϑϴ", "Ттᲄᲅ", "iİı", "éÉ"} {
		for _, field := range group {
			for _, key := range group {
				t.Run(string(field)+"/"+string(key), func(t *testing.T) {
					t.Parallel()
					typ := reflect.StructOf([]reflect.StructField{{Name: "Value", Type: reflect.TypeFor[int](), Tag: reflect.StructTag(`json:"` + string(field) + `"`)}})
					target := reflect.New(typ)
					if err := json.Unmarshal([]byte(`{"`+string(key)+`":1}`), target.Interface()); err != nil {
						t.Fatal(err)
					}
					matched := target.Elem().Field(0).Int() == 1
					if got := foldJSONKey(string(field)) == foldJSONKey(string(key)); got != matched {
						t.Fatalf("fold equivalence = %t, encoding/json struct match = %t", got, matched)
					}
				})
			}
		}
	}
}
