package journal

import (
	"fmt"
	"sort"
	"strings"
)

func validateCurrentUpdate(event Event, effect Effect) error {
	if effect.Op != "update" {
		return nil
	}
	values, guards := "", ""
	switch effect.Table {
	case "referenced_artifacts":
		values, guards = "availability", "availability"
	case "check_requests":
		values, guards = "response_json", "response_json"
	case "questions":
		values, guards = "resolved_at resolution decision_id", "resolved_at"
	case "idempotency_results":
		values, guards = "result_entity_id result_event_id aggregate_sequence result_json", "request_digest request_field_digests_json"
	default:
		return fmt.Errorf("current Event update table not admitted")
	}
	if !currentKeys(effect.Values, values) || !currentKeys(effect.Preconditions, guards) {
		return fmt.Errorf("current Event %s: wrong changed fields or guards for %s", event.EventType, effect.Table)
	}
	return nil
}

func currentKeys(values map[string]any, names string) bool {
	want := strings.Fields(names)
	sort.Strings(want)
	return strings.Join(sortedMapKeys(values), " ") == strings.Join(want, " ")
}
