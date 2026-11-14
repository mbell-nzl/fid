package issuestate

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestRegistryNextStates(t *testing.T) {
	t.Parallel()
	conn := migratedWorkflowRegistryDB(t)
	defer closeRegistryDB(t, conn)
	registry, err := Load(conn)
	if err != nil {
		t.Fatal(err)
	}
	for _, row := range []struct {
		state string
		want  []NextState
	}{
		{"in_review", []NextState{{"cancelled", "cancelled"}, {"done", "completed"}, {"in_progress", "started"}}},
		{"cancelled", []NextState{}},
		{"missing", []NextState{}},
	} {
		t.Run(row.state, func(t *testing.T) {
			t.Parallel()
			got := registry.NextStates("task-default", row.state)
			if !reflect.DeepEqual(got, row.want) {
				t.Fatalf("NextStates = %#v, want %#v", got, row.want)
			}
			definition, _ := registry.Workflow("task-default")
			for state := range definition.States {
				_, edge := registry.RulesForEdge("task-default", row.state, state)
				found := false
				for _, next := range got {
					found = found || next.State == state
				}
				if edge != found {
					t.Fatalf("edge predicate differs for %s", state)
				}
			}
			if len(got) == 0 {
				data, err := json.Marshal(got)
				if err != nil || string(data) != "[]" {
					t.Fatalf("empty JSON = %s, error %v", data, err)
				}
			}
		})
	}
	if got := registry.NextStates("unknown", "in_review"); got == nil || len(got) != 0 {
		t.Fatalf("unknown workflow = %#v", got)
	}
}
