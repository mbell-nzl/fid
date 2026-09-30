//go:build integration

package integration

import (
	"encoding/json"
	"os"
	"reflect"
	"testing"
)

// This model follows public-interface.md sections 3 and 4: a reference record
// has a generated row ID, while the caller's Requirement ID keys its links.
type p23RequirementExpectation struct {
	results map[string]map[string]any
}

func p23RequirementModel(t *testing.T, pairs map[string]p23HumanLeafPair) p23RequirementExpectation {
	t.Helper()
	record := p23LeafRequest(t, pairs["requirement.record-reference"])
	link := p23LeafRequest(t, pairs["requirement.link"])
	reference := map[string]any{
		"id": "<generated:RREF>", "source": record["source"], "status": record["status"],
		"superseded_by_id": nil, "created_at": "<generated:timestamp>", "event_id": "<generated:EventRef>",
	}
	linked := map[string]any{
		"requirement_id": link["requirement_id"], "issue_id": link["issue_id"],
		"relation": link["relation"], "created_at": "<generated:timestamp>", "event_id": "<generated:EventRef>",
	}
	return p23RequirementExpectation{results: map[string]map[string]any{
		"requirement.record-reference": {"requirement": reference, "idempotency_key": "leaf-requirement", "replayed": true},
		"requirement.link":             {"link": linked, "idempotency_key": "leaf-link", "replayed": true},
		"requirement.get":              {"requirement": reference, "links": []any{linked}},
		"requirement.list":             {"items": []any{reference}},
	}}
}

func p23LeafRequest(t *testing.T, pair p23HumanLeafPair) map[string]any {
	t.Helper()
	for i, arg := range pair.JSONArgs {
		if arg != "--input" || i+1 >= len(pair.JSONArgs) {
			continue
		}
		data, err := os.ReadFile(pair.JSONArgs[i+1])
		if err != nil {
			t.Fatal(err)
		}
		var request map[string]any
		if err := json.Unmarshal(data, &request); err != nil {
			t.Fatal(err)
		}
		return request
	}
	t.Fatalf("%s has no input", pair.Operation)
	return nil
}

func p23CheckRequirementIdentity(t *testing.T, accepted map[string]p23LeafOutcome) {
	t.Helper()
	for _, operation := range []string{"requirement.record-reference", "requirement.link", "requirement.get", "requirement.list"} {
		if _, ok := accepted[operation]; !ok {
			return // A -run subtest selection may omit this family.
		}
	}
	record := accepted["requirement.record-reference"].Semantics["requirement"]
	get := accepted["requirement.get"].Semantics
	items := accepted["requirement.list"].Semantics["items"].([]any)
	link := accepted["requirement.link"].Semantics["link"]
	links := get["links"].([]any)
	if len(items) != 1 || len(links) != 1 || !reflect.DeepEqual(record, get["requirement"]) ||
		!reflect.DeepEqual(record, items[0]) || !reflect.DeepEqual(link, links[0]) {
		t.Fatalf("Requirement mutation/read identity drift: record=%#v get=%#v list=%#v link=%#v", record, get, items, link)
	}
}
