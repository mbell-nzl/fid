package export

import (
	"strings"
	"testing"
)

// Pin the archive's JSON-text boundary independently of Event decoding, which
// also validates embedded text before the archive's schema checks run.
func TestSemanticJSONTextValidation(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name, raw, reason string
	}{
		{"distinct map keys", `{"build":1,"Build":2}`, ""},
		{"depth boundary", strings.Repeat("[", 10_000) + "0" + strings.Repeat("]", 10_000), ""},
		{"deep nesting", `{"deep":` + strings.Repeat("[", 10_001) + "0" + strings.Repeat("]", 10_001) + "}", "exceeded max depth"},
		{"duplicate key", `{"probe":1,"probe":2}`, "duplicate key"},
		{"overflow", `{"probe":1e999}`, "cannot unmarshal number"},
		{"trailing value", `{} {}`, "multiple JSON values"},
		{"malformed syntax", `{"probe":}`, "invalid character"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			err := validateSemanticJSONText(tc.raw)
			if tc.reason == "" {
				if err != nil {
					t.Fatalf("valid JSON refused: %v", err)
				}
			} else if err == nil || !strings.Contains(err.Error(), tc.reason) {
				t.Fatalf("JSON validation = %v, want %s refusal", err, tc.reason)
			}
		})
	}
}
