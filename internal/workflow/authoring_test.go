package workflow

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/mbell-nzl/fid/internal/journal"
)

func TestValidateDefinitionRequestIdentifierMatrix(t *testing.T) {
	conn := migratedWorkflowRegistryDB(t)
	defer closeRegistryDB(t, conn)

	tests := []struct {
		name    string
		request DefinitionRequest
		wantErr string
	}{
		{"invalid slug", withDefinition(func(request *DefinitionRequest) { request.Name = "Bad Flow" }), `"Bad Flow"`},
		{"duplicate slug", withDefinition(func(request *DefinitionRequest) { request.Name = "task-default" }), `"task-default"`},
		{"invalid state", withDefinition(func(request *DefinitionRequest) { request.States[0] = "Todo" }), `"Todo"`},
		{"duplicate state", withDefinition(func(request *DefinitionRequest) { request.States[1] = "todo" }), `"todo"`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertDefinitionError(t, conn, tt.request, tt.wantErr)
		})
	}
	if _, err := BuildDefinitionEffects(conn, validDefinition()); err != nil {
		t.Fatalf("BuildDefinitionEffects(validDefinition()) error = %v", err)
	}
}

func TestBuildDefinitionEffectsRequiresReadOnlyDatabase(t *testing.T) {
	if _, err := BuildDefinitionEffects(nil, validDefinition()); err == nil || !strings.Contains(err.Error(), "database connection is required") {
		t.Fatalf("BuildDefinitionEffects(nil) error = %v, want database connection error", err)
	}
}

func TestValidateDefinitionRequestStatesAndTransitions(t *testing.T) {
	conn := migratedWorkflowRegistryDB(t)
	defer closeRegistryDB(t, conn)
	tests := []struct {
		name    string
		request DefinitionRequest
		wantErr string
	}{
		{"empty states", withDefinition(func(request *DefinitionRequest) { request.States = nil; request.Categories = nil }), "must not be empty"},
		{"count mismatch", withDefinition(func(request *DefinitionRequest) { request.Categories = request.Categories[:2] }), "3 states, 2 categories"},
		{"unknown category", withDefinition(func(request *DefinitionRequest) { request.Categories[1] = "waiting" }), `"waiting"`},
		{"unknown transition state", withDefinition(func(request *DefinitionRequest) { request.Transitions = []Transition{{From: "todo", To: "missing"}} }), `"todo" -> "missing"`},
		{"duplicate transition", withDefinition(func(request *DefinitionRequest) {
			request.Transitions = []Transition{{From: "todo", To: "done"}, {From: "todo", To: "done"}}
		}), `"todo" -> "done"`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertDefinitionError(t, conn, tt.request, tt.wantErr)
		})
	}
}

func TestDefinitionRuleVocabularyAndParams(t *testing.T) {
	conn := migratedWorkflowRegistryDB(t)
	defer closeRegistryDB(t, conn)
	tests := []struct {
		name    string
		require Requirement
		wantErr string
	}{
		{"reserved v0.2b", Requirement{From: "todo", To: "done", Rule: RuleRequiredGatesCurrent}, "v0.2b"},
		{"reserved v0.2c", Requirement{From: "todo", To: "done", Rule: RuleHasCompletion}, "v0.2c"},
		{"reserved v0.2d", Requirement{From: "todo", To: "done", Rule: RuleLinkedToDecision}, "v0.2d"},
		{"unknown rule", Requirement{From: "todo", To: "done", Rule: "invented_rule"}, `"invented_rule"`},
		{"invalid json", Requirement{From: "todo", To: "done", Rule: RuleHasCompletionReport, ParamsJSON: "{"}, "invalid JSON"},
		{"all acs without roles", Requirement{From: "todo", To: "done", Rule: RuleAllACsVerified, ParamsJSON: `{"roles":[]}`}, "non-empty roles"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := withDefinition(func(request *DefinitionRequest) { request.Requirements = []Requirement{tt.require} })
			assertDefinitionError(t, conn, request, tt.wantErr)
		})
	}

	request := withDefinition(func(request *DefinitionRequest) {
		request.Requirements = []Requirement{
			{From: "todo", To: "done", Rule: RuleAllACsVerified, ParamsJSON: `{"roles":["developer"]}`},
			{From: "todo", To: "done", Rule: RuleHasCompletionReport, ParamsJSON: `{"note":"complete"}`},
			{From: "todo", To: "done", Rule: RuleConcernsAddressed, ParamsJSON: `{"note":"clear"}`},
			{From: "todo", To: "done", Rule: RuleAllChildTasksDoneOrCanceled, ParamsJSON: `{"children":true}`},
		}
	})
	effects, err := BuildDefinitionEffects(conn, request)
	if err != nil {
		t.Fatalf("BuildDefinitionEffects() error = %v", err)
	}
	for _, requirement := range request.Requirements {
		if !hasRequirementEffect(effects, requirement.From, requirement.To, requirement.Rule, requirement.ParamsJSON) {
			t.Fatalf("effects = %#v, missing %#v", effects, requirement)
		}
	}
}

func TestValidateDefinitionRequestRequirementEdges(t *testing.T) {
	conn := migratedWorkflowRegistryDB(t)
	defer closeRegistryDB(t, conn)
	unknown := withDefinition(func(request *DefinitionRequest) {
		request.Requirements = []Requirement{{From: "done", To: "todo", Rule: RuleHasCompletionReport}}
	})
	assertDefinitionError(t, conn, unknown, `"done" -> "todo"`)

	declared := withDefinition(func(request *DefinitionRequest) {
		request.Requirements = []Requirement{{From: "todo", To: "done", Rule: RuleHasCompletionReport}}
	})
	if _, err := BuildDefinitionEffects(conn, declared); err != nil {
		t.Fatalf("BuildDefinitionEffects(declared requirement) error = %v", err)
	}
}

func TestBuildDefinitionEffectsCancelledStatePolicy(t *testing.T) {
	conn := migratedWorkflowRegistryDB(t)
	defer closeRegistryDB(t, conn)
	request := DefinitionRequest{
		Name:        "cancel-flow",
		AppliesTo:   "task",
		States:      []string{"todo", "doing", "done", "cancelled"},
		Categories:  []string{"unstarted", "started", "completed", "cancelled"},
		Transitions: []Transition{{From: "todo", To: "doing"}, {From: "todo", To: "cancelled"}},
	}
	effects, err := BuildDefinitionEffects(conn, request)
	if err != nil {
		t.Fatalf("BuildDefinitionEffects() error = %v", err)
	}
	if got, want := transitionEffects(effects), []Transition{{From: "todo", To: "doing"}, {From: "todo", To: "cancelled"}, {From: "doing", To: "cancelled"}, {From: "done", To: "cancelled"}}; !reflect.DeepEqual(got, want) {
		t.Fatalf("transition effects = %#v, want %#v", got, want)
	}
}

func TestBuildDefinitionEffectsZeroCancelledApplicability(t *testing.T) {
	conn := migratedWorkflowRegistryDB(t)
	defer closeRegistryDB(t, conn)
	for _, appliesTo := range []string{"escalation", "qa_review"} {
		t.Run("accept "+appliesTo, func(t *testing.T) {
			request := zeroCancelledDefinition(appliesTo)
			effects, err := BuildDefinitionEffects(conn, request)
			if err != nil {
				t.Fatalf("BuildDefinitionEffects() error = %v", err)
			}
			if got, want := transitionEffects(effects), request.Transitions; !reflect.DeepEqual(got, want) {
				t.Fatalf("transition effects = %#v, want %#v", got, want)
			}
		})
	}
	for _, appliesTo := range []string{"epic", "task", "subtask", "any"} {
		t.Run("reject "+appliesTo, func(t *testing.T) {
			assertDefinitionError(t, conn, zeroCancelledDefinition(appliesTo), appliesTo)
			assertDefinitionError(t, conn, zeroCancelledDefinition(appliesTo), "PRD-WF-013")
		})
	}
}

func TestBuildDefinitionEffectsRejectsMultipleCancelledStates(t *testing.T) {
	conn := migratedWorkflowRegistryDB(t)
	defer closeRegistryDB(t, conn)
	request := withDefinition(func(request *DefinitionRequest) {
		request.States = append(request.States, "abandoned")
		request.Categories = append(request.Categories, "cancelled")
	})
	effects, err := BuildDefinitionEffects(conn, request)
	if err == nil || !strings.Contains(err.Error(), `"cancelled"`) || !strings.Contains(err.Error(), `"abandoned"`) {
		t.Fatalf("BuildDefinitionEffects() error = %v, want both cancelled states", err)
	}
	if effects != nil {
		t.Fatalf("effects = %#v, want no effects", effects)
	}
}

func TestInitialState(t *testing.T) {
	if got := InitialState([]State{{Name: "started", Category: "started"}, {Name: "todo", Category: "unstarted"}, {Name: "queued", Category: "unstarted"}}); got != "todo" {
		t.Fatalf("InitialState(first unstarted) = %q, want todo", got)
	}
	if got := InitialState([]State{{Name: "started", Category: "started"}, {Name: "done", Category: "completed"}}); got != "started" {
		t.Fatalf("InitialState(no unstarted) = %q, want started", got)
	}
}

func TestAuthoringCoreDoesNotWrite(t *testing.T) {
	conn := migratedWorkflowRegistryDB(t)
	defer closeRegistryDB(t, conn)
	before := workflowTableSnapshot(t, conn)
	journalPath := filepath.Join(t.TempDir(), "events.jsonl")
	if err := os.WriteFile(journalPath, []byte("journal before\n"), 0o600); err != nil {
		t.Fatalf("write journal fixture: %v", err)
	}
	beforeJournal, err := os.ReadFile(journalPath)
	if err != nil {
		t.Fatalf("read journal before: %v", err)
	}
	if _, err := conn.Exec(`PRAGMA query_only = ON`); err != nil {
		t.Fatalf("set query_only: %v", err)
	}
	if _, err := BuildDefinitionEffects(conn, validDefinition()); err != nil {
		t.Fatalf("BuildDefinitionEffects() error = %v", err)
	}
	if after := workflowTableSnapshot(t, conn); !reflect.DeepEqual(after, before) {
		t.Fatalf("workflow tables after core use = %#v, want %#v", after, before)
	}
	afterJournal, err := os.ReadFile(journalPath)
	if err != nil {
		t.Fatalf("read journal after: %v", err)
	}
	if !reflect.DeepEqual(afterJournal, beforeJournal) {
		t.Fatalf("journal bytes after core use = %q, want %q", afterJournal, beforeJournal)
	}
}

func validDefinition() DefinitionRequest {
	return DefinitionRequest{
		Name:        "valid-flow",
		AppliesTo:   "task",
		States:      []string{"todo", "done", "cancelled"},
		Categories:  []string{"unstarted", "completed", "cancelled"},
		Transitions: []Transition{{From: "todo", To: "done"}},
	}
}

func withDefinition(change func(*DefinitionRequest)) DefinitionRequest {
	request := validDefinition()
	change(&request)
	return request
}

func zeroCancelledDefinition(appliesTo string) DefinitionRequest {
	name := "zero-" + strings.ReplaceAll(appliesTo, "_", "-")
	return DefinitionRequest{Name: name, AppliesTo: appliesTo, States: []string{"open", "resolved"}, Categories: []string{"unstarted", "completed"}, Transitions: []Transition{{From: "open", To: "resolved"}}}
}

func assertDefinitionError(t *testing.T, conn *sql.DB, request DefinitionRequest, want string) {
	t.Helper()
	_, err := BuildDefinitionEffects(conn, request)
	if err == nil || !strings.Contains(err.Error(), want) {
		t.Fatalf("BuildDefinitionEffects(%#v) error = %v, want containing %q", request, err, want)
	}
}

func transitionEffects(effects []journal.Effect) []Transition {
	transitions := make([]Transition, 0)
	for _, effect := range effects {
		if effect.Table == "workflow_transitions" {
			transitions = append(transitions, Transition{From: effect.PK["from_state"].(string), To: effect.PK["to_state"].(string)})
		}
	}
	return transitions
}

func workflowTableSnapshot(t *testing.T, conn *sql.DB) map[string][]string {
	t.Helper()
	snapshot := make(map[string][]string)
	for _, table := range []string{"workflows", "workflow_states", "workflow_transitions", "workflow_transition_requires", "workflow_assignments"} {
		rows, err := conn.Query(fmt.Sprintf("SELECT * FROM %s ORDER BY rowid", table))
		if err != nil {
			t.Fatalf("query %s: %v", table, err)
		}
		columns, err := rows.Columns()
		if err != nil {
			if closeErr := rows.Close(); closeErr != nil {
				t.Fatalf("close %s after columns: %v", table, closeErr)
			}
			t.Fatalf("columns %s: %v", table, err)
		}
		for rows.Next() {
			values := make([]any, len(columns))
			pointers := make([]any, len(columns))
			for index := range values {
				pointers[index] = &values[index]
			}
			if err := rows.Scan(pointers...); err != nil {
				if closeErr := rows.Close(); closeErr != nil {
					t.Fatalf("close %s after scan: %v", table, closeErr)
				}
				t.Fatalf("scan %s: %v", table, err)
			}
			snapshot[table] = append(snapshot[table], fmt.Sprint(values...))
		}
		if err := rows.Err(); err != nil {
			if closeErr := rows.Close(); closeErr != nil {
				t.Fatalf("close %s after iterate: %v", table, closeErr)
			}
			t.Fatalf("iterate %s: %v", table, err)
		}
		if err := rows.Close(); err != nil {
			t.Fatalf("close %s: %v", table, err)
		}
	}
	return snapshot
}
