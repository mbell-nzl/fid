package main

import (
	"encoding/json"
	"sort"
	"testing"
)

func TestWorkflowListJSONMatchesContractAndUsesNullAssignmentScope(t *testing.T) {
	workspace := t.TempDir()
	clonePristineWorkspaceInto(t, workspace)
	defineWorkflowForTest(t, workspace, "list-flow", "any")
	if result := executeCommandForTest(workspace, "--json", "workflow", "assign", "--name", "list-flow", "--project-wide", "--by", "tl"); result.err != nil {
		t.Fatalf("assign project scope = %#v", result)
	}
	result := executeCommandForTest(workspace, "--json", "workflow", "list")
	if result.err != nil || result.stderr != "" {
		t.Fatalf("workflow list = %#v", result)
	}
	var envelope struct {
		ContractVersion string `json:"contract_version"`
		OK              bool   `json:"ok"`
		Operation       string `json:"operation"`
		Result          struct {
			Workflows []map[string]any `json:"workflows"`
			Count     int              `json:"count"`
		} `json:"result"`
		Event any `json:"event"`
	}
	if err := json.Unmarshal([]byte(result.stdout), &envelope); err != nil || !envelope.OK || envelope.ContractVersion != "fid-cli.v1" || envelope.Operation != "workflow.list" || envelope.Event != nil || envelope.Result.Count != len(envelope.Result.Workflows) {
		t.Fatalf("list envelope = %+v, err = %v", envelope, err)
	}
	assertEnvelopeKeySet(t, []byte(result.stdout), []string{"contract_version", "event", "ok", "operation", "result"})
	assertNestedEnvelopeKeySet(t, []byte(result.stdout), "result", []string{"count", "workflows"})
	for _, item := range envelope.Result.Workflows {
		keys := make([]string, 0, len(item))
		for key := range item {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		if got, want := joinWorkflowKeys(keys), "active,applies_to,assigned_to,id,name,states,transitions"; got != want {
			t.Fatalf("workflow keys = %s, want %s", got, want)
		}
	}
	for _, item := range envelope.Result.Workflows {
		if item["id"] == "list-flow" {
			assignments := item["assigned_to"].([]any)
			if len(assignments) != 1 || assignments[0].(map[string]any)["scope_type"] != "project" || assignments[0].(map[string]any)["scope_value"] != nil {
				t.Fatalf("project assignment = %#v", assignments)
			}
			return
		}
	}
	t.Fatal("list-flow missing from result")
}

func TestWorkflowListRejectsNonIssueKind(t *testing.T) {
	workspace := t.TempDir()
	clonePristineWorkspaceInto(t, workspace)
	result := executeCommandForTest(workspace, "--json", "workflow", "list", "--kind", "qa_review")
	if result.err == nil || result.stdout != "" {
		t.Fatalf("list rejection = %#v", result)
	}
	assertWorkflowInvalidEnvelope(t, result.stderr, "workflow.list")
}

func joinWorkflowKeys(keys []string) string {
	result := ""
	for i, key := range keys {
		if i > 0 {
			result += ","
		}
		result += key
	}
	return result
}
