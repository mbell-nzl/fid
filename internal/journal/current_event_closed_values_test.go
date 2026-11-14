package journal

import (
	"strings"
	"testing"
)

func TestCurrentValueDomains(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct{ name, kind, raw string }{
		{"byte-bound", "Text", `"` + strings.Repeat("界", 1400) + `"`},
		{"utc-time", "Timestamp", `"2026-09-09T00:00:00+02:00"`},
		{"invalid-utf8", "Text", "\"\xff\""},
	} {
		t.Run(tc.name, func(t *testing.T) {
			v, err := decodeCurrentJSON(tc.raw)
			if err == nil && validateCurrentValue(tc.kind, v) == nil {
				t.Fatal("invalid current scalar domain accepted")
			}
		})
	}
}

func TestCurrentNonJSONDomainsRejectMalformedIdentity(t *testing.T) {
	t.Parallel()
	for column, value := range map[string]any{"actor_id": "", "content_digest": "not-a-digest", "semantic_key": "plain-text", "declared_path": "../escape", "executor_id": "bad\nidentity"} {
		t.Run(column, func(t *testing.T) {
			if validateCurrentColumnDomain(column, value) == nil {
				t.Fatal("invalid non-JSON domain accepted")
			}
		})
	}
}
