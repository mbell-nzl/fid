package cli

import "testing"

func TestJournalSchemaRefusalDetailsAreClosed(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name  string
		edit  func(map[string]any)
		valid bool
	}{
		{"schema refusal", func(map[string]any) {}, true},
		{"unknown reason", func(d map[string]any) { d["reason"] = "unknown" }, false},
		{"missing family", func(d map[string]any) { delete(d, "event_type") }, false},
		{"empty family", func(d map[string]any) { d["event_type"] = "" }, false},
		{"missing reason", func(d map[string]any) { delete(d, "reason") }, false},
		{"extra detail", func(d map[string]any) { d["extra"] = "no" }, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			d := map[string]any{"journal_file": "/journal/2026-11-15.jsonl", "line": 2, "event_id": "EVT-removed", "reason": "schema-unsupported", "event_type": "claim.acquired"}
			tc.edit(d)
			_, err := NewFailure(CodeJournalReplayStopped, "journal replay stopped", nil, d)
			if (err == nil) != tc.valid {
				t.Fatalf("valid=%t error=%v", tc.valid, err)
			}
		})
	}
}
