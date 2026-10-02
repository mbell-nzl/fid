//go:build integration

package integration

import (
	"reflect"
	"testing"
)

// public-interface.md section 3 fixes the complete Concern and disposition
// shapes. Values below come from the two caller requests, not fid responses.
func p23ConcernModel(t *testing.T, pairs map[string]p23HumanLeafPair) map[string]map[string]any {
	t.Helper()
	record := p23LeafRequest(t, pairs["concern.record"])
	dispose := p23LeafRequest(t, pairs["concern.dispose"])
	base := map[string]any{
		"id": record["id"], "subject": record["subject"], "actor_id": p23LeafFlag(t, pairs["concern.record"], "--actor"),
		"scope": record["scope"], "severity": record["severity"], "summary": record["summary"],
		"reasoning": record["reasoning"], "references": record["references"], "status": "open",
		"disposition": nil, "created_at": "<generated:timestamp>", "event_id": "<generated:EventRef>",
	}
	disposed := make(map[string]any, len(base))
	for key, value := range base {
		disposed[key] = value
	}
	disposed["status"] = "disposed"
	disposed["disposition"] = map[string]any{
		"action": dispose["action"], "issue_id": dispose["issue_id"],
		"supersedes_concern_id": dispose["supersedes_concern_id"],
		"supersedes_finding_id": dispose["supersedes_finding_id"],
		"reason":                dispose["reason"], "actor_id": p23LeafFlag(t, pairs["concern.dispose"], "--actor"),
		"disposed_at": "<generated:timestamp>", "event_id": "<generated:EventRef>",
	}
	return map[string]map[string]any{
		"concern.record":  {"concern": base, "idempotency_key": p23LeafFlag(t, pairs["concern.record"], "--idempotency-key"), "replayed": true},
		"concern.dispose": {"concern": disposed, "idempotency_key": p23LeafFlag(t, pairs["concern.dispose"], "--idempotency-key"), "replayed": true},
		"concern.list":    {"items": []any{disposed}},
	}
}

func p23LeafFlag(t *testing.T, pair p23HumanLeafPair, flag string) string {
	t.Helper()
	for i, arg := range pair.JSONArgs {
		if arg == flag && i+1 < len(pair.JSONArgs) {
			return pair.JSONArgs[i+1]
		}
	}
	t.Fatalf("%s has no %s", pair.Operation, flag)
	return ""
}

func p23CheckConcernIdentity(t *testing.T, accepted map[string]p23LeafOutcome) {
	t.Helper()
	for _, operation := range []string{"concern.record", "concern.dispose", "concern.list"} {
		if _, ok := accepted[operation]; !ok {
			return // A -run subtest selection may omit this family.
		}
	}
	recorded := accepted["concern.record"]
	disposed := accepted["concern.dispose"].Semantics["concern"].(map[string]any)
	listed := accepted["concern.list"].Semantics["items"].([]any)
	base := recorded.Semantics["concern"].(map[string]any)
	if len(listed) != 1 || !reflect.DeepEqual(disposed, listed[0]) ||
		base["id"] != disposed["id"] || base["created_at"] != disposed["created_at"] ||
		base["event_id"] != disposed["event_id"] {
		t.Fatalf("Concern record/dispose/list identity drift: %#v %#v %#v", base, disposed, listed)
	}
}
