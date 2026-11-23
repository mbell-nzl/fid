package journal

import (
	"strings"
	"testing"
)

// These readers also run on a stored retry without replaying its Event.
func TestStoredRetryJSONKeys(t *testing.T) {
	t.Parallel()
	oid := strings.Repeat("a", 40)
	for _, tc := range []struct {
		name, raw, field, alias string
		validate                func(string) error
	}{
		{"create", `{"issue":{"kind":"task","lane":"standard"},"replayed":false}`, `"kind":"task"`, `"kind":"epic","Kind":"task"`, ValidateIssueCreateResult},
		{"update", `{"changes":[{"field":"lane","from":"standard","to":"assured"}],"replayed":false}`, `"field":"lane"`, `"field":"title","field":"lane"`, ValidateIssueUpdateResult},
		{"update struct alias", `{"changes":[{"field":"lane","from":"standard","to":"assured"}],"replayed":false}`, `"field":"lane"`, `"field":"title","Field":"lane"`, ValidateIssueUpdateResult},
		{"start", `{"fields":{"task_id":"FID-a","hold_id":"HLD-a","status":"in_progress","lane":"standard"},"occurred_at":"2026-11-10T00:00:00Z","replayed":false}`, `"fields":`, `"fieldſ":{},"fields":`, ValidateTaskStartResult},
		{"handoff", `{"fields":{"task_id":"FID-a","handoff_id":"HND-a","status":"in_review","pr":"https://github.com/mbell-nzl/fid/pull/101","head":"` + oid + `","base":"` + oid + `","required_checks":[],"result_ids":[]},"occurred_at":"2026-11-10T00:00:00Z","replayed":false}`, `"fields":`, `"fieldſ":{},"fields":`, ValidateHandoffResult},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if err := tc.validate(tc.raw); err != nil {
				t.Fatalf("valid stored result: %v", err)
			}
			for name, malformed := range map[string]string{
				"exact duplicate": strings.Replace(tc.raw, `"replayed":false`, `"replayed":true,"replayed":false`, 1),
				"typed field":     strings.Replace(tc.raw, tc.field, tc.alias, 1),
			} {
				t.Run(name, func(t *testing.T) {
					t.Parallel()
					if malformed == tc.raw {
						t.Fatal("probe did not change the retry JSON")
					}
					if err := tc.validate(malformed); err == nil || !strings.Contains(err.Error(), "duplicate key") {
						t.Fatalf("malformed stored result = %v, want duplicate-key refusal", err)
					}
				})
			}
			// Unknown fields remain allowed and map entries still match exactly.
			control := strings.Replace(tc.raw, `"replayed":false`, `"replayed":false,"extra":{"build":1,"Build":2,"s":3,"ſ":4}`, 1)
			if err := tc.validate(control); err != nil {
				t.Fatalf("distinct map entries in stored result: %v", err)
			}
		})
	}
}
