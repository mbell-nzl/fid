package main

import (
	"encoding/json"
	"strings"
	"testing"
	"time"
)

func TestWorkflowDefineJSONMatchesContract(t *testing.T) {
	workspace := t.TempDir()
	clonePristineWorkspaceInto(t, workspace)
	result := executeCommandForTest(workspace, "--json", "workflow", "define", "--name", "contract-flow", "--applies-to", "task", "--states", "todo,done,cancelled", "--categories", "unstarted,completed,cancelled", "--transitions", "todo->done", "--by", "tl")
	if result.err != nil || result.stderr != "" {
		t.Fatalf("workflow define result = %#v, want stdout-only JSON success", result)
	}
	var envelope struct {
		ContractVersion string `json:"contract_version"`
		OK              bool   `json:"ok"`
		Operation       string `json:"operation"`
		Result          struct {
			Workflow struct {
				ID              string `json:"id"`
				Name            string `json:"name"`
				AppliesTo       string `json:"applies_to"`
				StateCount      int    `json:"state_count"`
				TransitionCount int    `json:"transition_count"`
			} `json:"workflow"`
		} `json:"result"`
		Event *struct {
			ID         string `json:"id"`
			OccurredAt string `json:"occurred_at"`
		} `json:"event"`
	}
	if err := json.Unmarshal([]byte(result.stdout), &envelope); err != nil {
		t.Fatalf("decode workflow define response: %v; stdout = %q", err, result.stdout)
	}
	assertEnvelopeKeySet(t, []byte(result.stdout), []string{"contract_version", "event", "ok", "operation", "result"})
	assertNestedEnvelopeKeySet(t, []byte(result.stdout), "result", []string{"workflow"})
	assertNestedEnvelopeKeySet(t, []byte(result.stdout), "result.workflow", []string{"applies_to", "id", "name", "state_count", "transition_count"})
	if envelope.ContractVersion != "fid-cli.v1" || !envelope.OK || envelope.Operation != "workflow.define" || envelope.Event == nil || envelope.Event.ID == "" || envelope.Event.OccurredAt == "" || envelope.Result.Workflow.ID != "contract-flow" || envelope.Result.Workflow.AppliesTo != "task" || envelope.Result.Workflow.StateCount != 3 || envelope.Result.Workflow.TransitionCount != 3 {
		t.Fatalf("envelope = %+v, want closed workflow.define success", envelope)
	}
	event := mustReadLastJournalEvent(t, mustJournalPath(t, workspace))
	if event.EventID != envelope.Event.ID || event.OccurredAt.UTC().Format(time.RFC3339) != envelope.Event.OccurredAt || event.Actor.Role != "tl" || len(event.Effects) != 7 {
		t.Fatalf("event = %#v, want supplied attribution and no requirement effects", event)
	}
}

func TestWorkflowDefineRejectsRemovedAndInvalidRequestsBeforeAppend(t *testing.T) {
	for _, test := range []struct {
		name string
		args []string
	}{
		{"removed require", []string{"--require", "todo->done:all_acs_verified"}},
		{"legacy applicability", []string{"--applies-to", "qa_review"}},
		{"blank token", []string{"--states", "todo,,cancelled"}},
		{"noncanonical arrow", []string{"--transitions", "todo→done"}},
		{"too many states", []string{"--states", strings.Repeat("s,", 64) + "cancelled", "--categories", strings.Repeat("unstarted,", 64) + "cancelled"}},
	} {
		t.Run(test.name, func(t *testing.T) {
			workspace := t.TempDir()
			clonePristineWorkspaceInto(t, workspace)
			journal := mustJournalPath(t, workspace)
			before := mustReadCommandFile(t, journal)
			args := []string{"--json", "workflow", "define", "--name", "reject-flow", "--applies-to", "task", "--states", "todo,done,cancelled", "--categories", "unstarted,completed,cancelled", "--transitions", "todo->done", "--by", "tl"}
			args = append(args, test.args...)
			result := executeCommandForTest(workspace, args...)
			if result.err == nil || result.stdout != "" {
				t.Fatalf("result = %#v, want stderr-only rejection", result)
			}
			assertWorkflowInvalidEnvelope(t, result.stderr, "workflow.define")
			if after := mustReadCommandFile(t, journal); string(after) != string(before) {
				t.Fatal("rejected define appended an event")
			}
		})
	}
}

func TestWorkflowDefineDuplicateIsTyped(t *testing.T) {
	workspace := t.TempDir()
	clonePristineWorkspaceInto(t, workspace)
	args := []string{"--json", "workflow", "define", "--name", "once-flow", "--applies-to", "task", "--states", "todo,done,cancelled", "--categories", "unstarted,completed,cancelled", "--transitions", "todo->done", "--by", "tl"}
	if result := executeCommandForTest(workspace, args...); result.err != nil {
		t.Fatalf("initial define = %#v", result)
	}
	result := executeCommandForTest(workspace, args...)
	if result.err == nil || result.stdout != "" {
		t.Fatalf("duplicate result = %#v", result)
	}
	var failure struct {
		Error struct {
			Code    string `json:"code"`
			Subject struct {
				Type string `json:"type"`
				ID   string `json:"id"`
			} `json:"subject"`
		} `json:"error"`
	}
	if err := json.Unmarshal([]byte(result.stderr), &failure); err != nil || failure.Error.Code != "duplicate-entity" || failure.Error.Subject.Type != "workflow" || failure.Error.Subject.ID != "once-flow" {
		t.Fatalf("duplicate failure = %+v, err = %v", failure, err)
	}
}

func TestWorkflowDefineHumanSuccessUsesStdoutOnly(t *testing.T) {
	workspace := t.TempDir()
	clonePristineWorkspaceInto(t, workspace)
	result := executeCommandForTest(workspace, "workflow", "define", "--name", "human-flow", "--applies-to", "task", "--states", "todo,done,cancelled", "--categories", "unstarted,completed,cancelled", "--transitions", "todo->done", "--by", "tl")
	if result.err != nil || result.stdout != "defined workflow human-flow\n" || result.stderr != "" {
		t.Fatalf("Human workflow define = %#v, want stdout-only confirmation", result)
	}
}

func assertWorkflowInvalidEnvelope(t *testing.T, output, operation string) {
	t.Helper()
	var failure struct {
		OK        bool   `json:"ok"`
		Operation string `json:"operation"`
		Error     struct {
			Code     string `json:"code"`
			Category string `json:"category"`
			Subject  any    `json:"subject"`
		} `json:"error"`
	}
	if err := json.Unmarshal([]byte(output), &failure); err != nil || failure.OK || failure.Operation != operation || failure.Error.Code != "invalid-argument" || failure.Error.Category != "invalid-input" || failure.Error.Subject != nil {
		t.Fatalf("failure = %+v, err = %v", failure, err)
	}
}
