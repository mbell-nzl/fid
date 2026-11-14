package export

import (
	"encoding/json"
	"testing"

	"github.com/mbell-nzl/fid/internal/journal"
)

func TestStartReplayAndArchiveRejectMalformedResultFields(t *testing.T) {
	t.Parallel()
	for _, field := range []string{"lane", "task_id", "hold_id", "status", "occurred_at", "fields"} {
		for _, row := range []struct {
			name  string
			value any
		}{
			{"missing", nil}, {"null", nil}, {"number", 42}, {"empty", ""},
			{"array", []any{}}, {"object", map[string]any{}}, {"boolean", true}, {"unknown", "urgent"},
		} {
			for _, mode := range []string{"replay", "archive"} {
				t.Run(field+"/"+row.name+"/"+mode, func(t *testing.T) {
					t.Parallel()
					validator := startArchiveValidator(t)
					event := archiveStartEvent(t)
					last := &event.Effects[len(event.Effects)-1]
					var result map[string]any
					if err := json.Unmarshal([]byte(last.Values["result_json"].(string)), &result); err != nil {
						t.Fatal(err)
					}
					fields := result
					if field != "fields" && field != "occurred_at" {
						fields = result["fields"].(map[string]any)
					}
					fields[field] = row.value
					if row.name == "missing" {
						delete(fields, field)
					}
					raw, err := journal.MarshalCanonicalJSON(result)
					if err != nil {
						t.Fatal(err)
					}
					last.Values["result_json"] = string(raw)
					if err := acceptRetryEvent(t, validator, event, mode); err == nil {
						t.Fatal("malformed start retry result accepted")
					}
					var unchanged, bindings int
					if err := validator.conn.QueryRow(`SELECT (SELECT count(*) FROM issues WHERE status='todo' AND revision=1), (SELECT count(*) FROM idempotency_results)`).Scan(&unchanged, &bindings); err != nil {
						t.Fatal(err)
					}
					if unchanged != 2 || bindings != 0 {
						t.Fatalf("rejected start left partial state = %d unchanged/%d bindings", unchanged, bindings)
					}
				})
			}
		}
	}
}
