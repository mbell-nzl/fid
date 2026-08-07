package workflow

import (
	"reflect"
	"testing"

	"github.com/mbell-nzl/fid/internal/journal"
)

func TestBuildDefinitionEffectsIsOrderedAndDeterministic(t *testing.T) {
	t.Parallel()
	conn := migratedWorkflowRegistryDB(t)
	defer closeRegistryDB(t, conn)

	request := DefinitionRequest{
		Name:      "release-flow",
		AppliesTo: "task",
		States:    []string{"todo", "in_progress", "done", "cancelled"},
		Categories: []string{
			"unstarted", "started", "completed", "cancelled",
		},
		Transitions: []Transition{
			{From: "todo", To: "in_progress"},
			{From: "in_progress", To: "done"},
			{From: "todo", To: "cancelled"},
		},
		Requirements: []Requirement{
			{From: "in_progress", To: "done", Rule: RuleAllACsVerified, ParamsJSON: `{"roles":["developer"]}`},
			{From: "in_progress", To: "done", Rule: RuleHasCompletionReport},
		},
	}

	got, err := BuildDefinitionEffects(conn, request)
	if err != nil {
		t.Fatalf("BuildDefinitionEffects() error = %v", err)
	}
	again, err := BuildDefinitionEffects(conn, request)
	if err != nil {
		t.Fatalf("BuildDefinitionEffects() second error = %v", err)
	}
	if !reflect.DeepEqual(again, got) {
		t.Fatalf("second effect build = %#v, want %#v", again, got)
	}

	want := []journal.Effect{
		{Op: "insert", Table: "workflows", PK: map[string]any{"id": "release-flow"}, Values: map[string]any{"name": "release-flow", "applies_to": "task", "active": true}},
		{Op: "insert", Table: "workflow_states", PK: map[string]any{"workflow_id": "release-flow", "state": "todo"}, Values: map[string]any{"category": "unstarted", "ordinal": 1}},
		{Op: "insert", Table: "workflow_states", PK: map[string]any{"workflow_id": "release-flow", "state": "in_progress"}, Values: map[string]any{"category": "started", "ordinal": 2}},
		{Op: "insert", Table: "workflow_states", PK: map[string]any{"workflow_id": "release-flow", "state": "done"}, Values: map[string]any{"category": "completed", "ordinal": 3}},
		{Op: "insert", Table: "workflow_states", PK: map[string]any{"workflow_id": "release-flow", "state": "cancelled"}, Values: map[string]any{"category": "cancelled", "ordinal": 4}},
		{Op: "insert", Table: "workflow_transitions", PK: map[string]any{"workflow_id": "release-flow", "from_state": "todo", "to_state": "in_progress"}},
		{Op: "insert", Table: "workflow_transitions", PK: map[string]any{"workflow_id": "release-flow", "from_state": "in_progress", "to_state": "done"}},
		{Op: "insert", Table: "workflow_transitions", PK: map[string]any{"workflow_id": "release-flow", "from_state": "todo", "to_state": "cancelled"}},
		{Op: "insert", Table: "workflow_transitions", PK: map[string]any{"workflow_id": "release-flow", "from_state": "in_progress", "to_state": "cancelled"}},
		{Op: "insert", Table: "workflow_transitions", PK: map[string]any{"workflow_id": "release-flow", "from_state": "done", "to_state": "cancelled"}},
		{Op: "insert", Table: "workflow_transition_requires", PK: map[string]any{"workflow_id": "release-flow", "from_state": "in_progress", "to_state": "done", "rule": RuleAllACsVerified}, Values: map[string]any{"params_json": `{"roles":["developer"]}`}},
		{Op: "insert", Table: "workflow_transition_requires", PK: map[string]any{"workflow_id": "release-flow", "from_state": "in_progress", "to_state": "done", "rule": RuleHasCompletionReport}, Values: map[string]any{"params_json": nil}},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("effects = %#v, want %#v", got, want)
	}
}

func TestBuildDefinitionEffectsAcceptsRequirementOnMaterializedEdge(t *testing.T) {
	t.Parallel()
	conn := migratedWorkflowRegistryDB(t)
	defer closeRegistryDB(t, conn)

	request := DefinitionRequest{
		Name:        "materialized-edge-flow",
		AppliesTo:   "task",
		States:      []string{"todo", "done", "cancelled"},
		Categories:  []string{"unstarted", "completed", "cancelled"},
		Transitions: []Transition{{From: "todo", To: "done"}},
		Requirements: []Requirement{{
			From: "done", To: "cancelled", Rule: RuleHasCompletionReport,
		}},
	}

	effects, err := BuildDefinitionEffects(conn, request)
	if err != nil {
		t.Fatalf("BuildDefinitionEffects() error = %v", err)
	}
	if !hasRequirementEffect(effects, "done", "cancelled", RuleHasCompletionReport, nil) {
		t.Fatalf("effects = %#v, want requirement on materialized done -> cancelled edge", effects)
	}
}

func hasRequirementEffect(effects []journal.Effect, from, to, rule string, params any) bool {
	for _, effect := range effects {
		if effect.Table != "workflow_transition_requires" || effect.PK["from_state"] != from || effect.PK["to_state"] != to || effect.PK["rule"] != rule {
			continue
		}
		return reflect.DeepEqual(effect.Values["params_json"], params)
	}
	return false
}
