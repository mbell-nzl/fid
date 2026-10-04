//go:build integration

package integration

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"testing"

	"github.com/mbell-nzl/fid/internal/journal"
)

// Exercise the closed-value promise using genuine producer payloads. Each
// changed value is applied in a disposable transaction against its real prefix.
// Existing companion tests check publication/replay rollback across all tables.
func p22RejectNestedCurrentValues(t *testing.T, conn *sql.DB, event journal.Event) {
	t.Helper()
	for index, effect := range event.Effects {
		for column, raw := range effect.Values {
			if !strings.HasSuffix(column, "_json") || raw == nil {
				continue
			}
			var root any
			p22FaultOK(t, json.Unmarshal([]byte(raw.(string)), &root))
			check := func(label string) {
				// Task 2 field maps are open by the caller's written optional
				// members. Their keys cannot be inferred from the Event alone.
				if column == "request_field_digests_json" && (strings.HasSuffix(label, "/unknown") || strings.Contains(label, "/missing-")) {
					return
				}
				body, err := json.Marshal(event)
				p22FaultOK(t, err)
				var changed journal.Event
				p22FaultOK(t, json.Unmarshal(body, &changed))
				body, err = json.Marshal(root)
				p22FaultOK(t, err)
				changed.Effects[index].Values[column] = string(body)
				t.Run(fmt.Sprintf("closed/%d/%s/%s", index, column, label), func(t *testing.T) {
					defer func() {
						if p := recover(); p != nil {
							t.Errorf("malformed captured value panicked: %v", p)
						}
					}()
					if err := journal.DryRunProjectEvent(conn, changed); err == nil {
						t.Fatal("malformed nested value accepted")
					}
				})
			}
			p22MutateClosedJSON(&root, "$", check, column == "result_json")
		}
	}
}

func p22MutateClosedJSON(value *any, path string, check func(string), semantic bool) {
	original := *value
	switch v := original.(type) {
	case map[string]any:
		v["__unexpected_contract_field__"] = true
		check(path + "/unknown")
		delete(v, "__unexpected_contract_field__")
		keys := make([]string, 0, len(v))
		for key := range v {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		for _, key := range keys {
			child := v[key]
			delete(v, key)
			check(path + "/missing-" + key)
			v[key] = child
			p22MutateClosedJSON(&child, path+"/"+key, func(label string) { v[key] = child; check(label) }, semantic)
			v[key] = child
		}
	case []any:
		for i := range v {
			p22MutateClosedJSON(&v[i], fmt.Sprintf("%s/%d", path, i), check, semantic)
		}
	}
	if _, boolean := original.(bool); boolean {
		*value = "false"
	} else {
		*value = true
	}
	check(path + "/wrong-type")
	*value = original
	if !semantic {
		return
	}
	switch v := original.(type) {
	case string:
		*value = v + "-altered"
	case float64:
		*value = v + 1
	case bool:
		*value = !v
	case []any:
		if len(v) == 0 {
			return
		}
		*value = append([]any{}, v[:len(v)-1]...)
	default:
		return
	}
	check(path + "/changed-recorded-value")
	*value = original
}
