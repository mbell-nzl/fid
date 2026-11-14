package issuestate

import "testing"

func TestIsAllowedEdge(t *testing.T) {
	t.Parallel()
	if !IsAllowedEdge("task-default", "todo", "in_progress") {
		t.Fatalf("expected todo -> in_progress to be allowed for task-default")
	}
	if IsAllowedEdge("task-default", "todo", "done") {
		t.Fatalf("expected todo -> done to be disallowed for task-default")
	}
	if !IsAllowedEdge("task-default", "in_review", "cancelled") {
		t.Fatalf("expected qa_review -> cancelled to be allowed for task-default")
	}
	if IsAllowedEdge("task-default", "cancelled", "cancelled") {
		t.Fatalf("expected cancelled -> cancelled to be disallowed for task-default")
	}
}

func TestCategoryForStateMatchesDefaultWorkflows(t *testing.T) {
	t.Parallel()

	tests := []struct {
		workflowID string
		state      string
		want       string
	}{
		{workflowID: "task-default", state: "todo", want: "unstarted"},
		{workflowID: "task-default", state: "done", want: "completed"},
		{workflowID: "task-default", state: "cancelled", want: "cancelled"},
		{workflowID: "epic-default", state: "in_progress", want: "started"},
		{workflowID: "epic-default", state: "done", want: "completed"},
		{workflowID: "subtask-default", state: "in_progress", want: "started"},
		{workflowID: "subtask-default", state: "done", want: "completed"},
	}

	for _, tt := range tests {
		t.Run(tt.workflowID+"-"+tt.state, func(t *testing.T) {
			got, ok := CategoryForState(tt.workflowID, tt.state)
			if !ok {
				t.Fatalf("CategoryForState(%q, %q) reported missing state", tt.workflowID, tt.state)
			}
			if got != tt.want {
				t.Fatalf("CategoryForState(%q, %q) = %q, want %q", tt.workflowID, tt.state, got, tt.want)
			}
		})
	}

	if _, ok := CategoryForState("task-default", "missing"); ok {
		t.Fatalf("expected missing task state lookup to fail")
	}
	if _, ok := CategoryForState("missing-workflow", "todo"); ok {
		t.Fatalf("expected missing workflow lookup to fail")
	}
}

func TestRulesForEdgeMatchesTRDRequirements(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		workflowID string
		fromState  string
		toState    string
		want       []RuleReference
	}{
		{
			name:       "task dev complete requires developer ac verification and completion report",
			workflowID: "task-default",
			fromState:  "in_progress",
			toState:    "in_review",
			want:       nil,
		},
		{
			name:       "task qa pass requires qa lead verification",
			workflowID: "task-default",
			fromState:  "in_review",
			toState:    "done",
			want:       nil,
		},
		{
			name:       "epic wrapped requires child task rule",
			workflowID: "epic-default",
			fromState:  "in_progress",
			toState:    "done",
			want:       nil,
		},
		{
			name:       "cancel edge has no per edge rules",
			workflowID: "subtask-default",
			fromState:  "in_progress",
			toState:    "cancelled",
			want:       nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := RulesForEdge(tt.workflowID, tt.fromState, tt.toState)
			if !ok {
				t.Fatalf("RulesForEdge(%q, %q, %q) reported missing edge", tt.workflowID, tt.fromState, tt.toState)
			}
			if len(got) != len(tt.want) {
				t.Fatalf("RulesForEdge(%q, %q, %q) len = %d, want %d", tt.workflowID, tt.fromState, tt.toState, len(got), len(tt.want))
			}
			for i := range got {
				if got[i].Name != tt.want[i].Name {
					t.Fatalf("RulesForEdge(%q, %q, %q)[%d].Name = %q, want %q", tt.workflowID, tt.fromState, tt.toState, i, got[i].Name, tt.want[i].Name)
				}
				if got[i].Unimplemented != tt.want[i].Unimplemented {
					t.Fatalf("rule %d unimplemented = %v, want %v", i, got[i].Unimplemented, tt.want[i].Unimplemented)
				}
			}
		})
	}

	if _, ok := RulesForEdge("task-default", "todo", "done"); ok {
		t.Fatalf("expected todo -> done to be rejected")
	}
	if _, ok := RulesForEdge("task-default", "missing", "cancelled"); ok {
		t.Fatalf("expected missing from-state to be rejected")
	}
}

func TestWorkflowDefinitionsMatchTRDExactly(t *testing.T) {
	t.Parallel()

	tests := []struct {
		workflowID  string
		states      map[string]string
		transitions map[string]map[string][]RuleReference
	}{
		{workflowID: "task-default", states: map[string]string{"todo": "unstarted", "in_progress": "started", "in_review": "started", "done": "completed", "cancelled": "cancelled"}, transitions: map[string]map[string][]RuleReference{"todo": {"in_progress": nil}, "in_progress": {"in_review": nil}, "in_review": {"in_progress": nil, "done": nil}}},
		{workflowID: "epic-default", states: map[string]string{"todo": "unstarted", "in_progress": "started", "done": "completed", "cancelled": "cancelled"}, transitions: map[string]map[string][]RuleReference{"todo": {"in_progress": nil}, "in_progress": {"done": nil}}},
		{
			workflowID: "subtask-default",
			states: map[string]string{
				"todo":        "unstarted",
				"in_progress": "started",
				"done":        "completed",
				"cancelled":   "cancelled",
			},
			transitions: map[string]map[string][]RuleReference{
				"todo": {
					"in_progress": nil,
				},
				"in_progress": {
					"done": nil,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.workflowID, func(t *testing.T) {
			definition, ok := Workflow(tt.workflowID)
			if !ok {
				t.Fatalf("Workflow(%q) reported missing workflow", tt.workflowID)
			}
			if len(definition.States) != len(tt.states) {
				t.Fatalf("len(Workflow(%q).States) = %d, want %d", tt.workflowID, len(definition.States), len(tt.states))
			}
			for state, wantCategory := range tt.states {
				gotCategory, ok := definition.States[state]
				if !ok {
					t.Fatalf("Workflow(%q).States missing %q", tt.workflowID, state)
				}
				if gotCategory != wantCategory {
					t.Fatalf("Workflow(%q).States[%q] = %q, want %q", tt.workflowID, state, gotCategory, wantCategory)
				}
			}

			for fromState := range tt.states {
				for toState := range tt.states {
					wantRules, wantAllowed := expectedRules(tt.transitions, fromState, toState)
					if toState == "cancelled" && fromState != "cancelled" {
						wantRules = nil
						wantAllowed = true
					}

					if gotAllowed := IsAllowedEdge(tt.workflowID, fromState, toState); gotAllowed != wantAllowed {
						t.Fatalf("IsAllowedEdge(%q, %q, %q) = %t, want %t", tt.workflowID, fromState, toState, gotAllowed, wantAllowed)
					}

					gotRules, gotAllowed := RulesForEdge(tt.workflowID, fromState, toState)
					if gotAllowed != wantAllowed {
						t.Fatalf("RulesForEdge(%q, %q, %q) allowed = %t, want %t", tt.workflowID, fromState, toState, gotAllowed, wantAllowed)
					}
					if !wantAllowed {
						continue
					}
					assertRuleReferencesEqual(t, tt.workflowID, fromState, toState, gotRules, wantRules)
				}
			}
		})
	}
}

func expectedRules(transitions map[string]map[string][]RuleReference, fromState, toState string) ([]RuleReference, bool) {
	destinations, ok := transitions[fromState]
	if !ok {
		return nil, false
	}
	rules, ok := destinations[toState]
	if !ok {
		return nil, false
	}
	return rules, true
}

func assertRuleReferencesEqual(t *testing.T, workflowID, fromState, toState string, got, want []RuleReference) {
	t.Helper()

	if len(got) != len(want) {
		t.Fatalf("RulesForEdge(%q, %q, %q) len = %d, want %d", workflowID, fromState, toState, len(got), len(want))
	}
	for i := range got {
		if got[i].Name != want[i].Name {
			t.Fatalf("RulesForEdge(%q, %q, %q)[%d].Name = %q, want %q", workflowID, fromState, toState, i, got[i].Name, want[i].Name)
		}
		if got[i].Unimplemented != want[i].Unimplemented {
			t.Fatalf("rule %d unimplemented = %v, want %v", i, got[i].Unimplemented, want[i].Unimplemented)
		}
	}
}
