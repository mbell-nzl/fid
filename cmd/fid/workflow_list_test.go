package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"sort"
	"testing"
)

func TestWorkflowListJSONMatchesContractAndUsesNullAssignmentScope(t *testing.T) {
	t.Parallel()
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
	t.Parallel()
	workspace := t.TempDir()
	clonePristineWorkspaceInto(t, workspace)
	result := executeCommandForTest(workspace, "--json", "workflow", "list", "--kind", "qa_review")
	if result.err == nil || result.stdout != "" {
		t.Fatalf("list rejection = %#v", result)
	}
	assertWorkflowInvalidEnvelope(t, result.stderr, "workflow.list")
}

func TestWorkflowListEnforcesCollectionLimitWithoutTruncation(t *testing.T) {
	t.Parallel()
	workspace := t.TempDir()
	clonePristineWorkspaceInto(t, workspace)
	conn := openWorkflowCommandDB(t, workspace)
	defer closeCommandDB(t, conn)
	resetWorkflowListFixture(t, conn)

	for index := 0; index < 10000; index++ {
		insertWorkflowListFixture(t, conn, index)
	}
	atLimit := executeCommandForTest(workspace, "--json", "workflow", "list")
	if atLimit.err != nil || atLimit.stderr != "" {
		t.Fatalf("at-limit workflow.list = %#v, want stdout-only success", atLimit)
	}
	var success struct {
		Result struct {
			Workflows []struct {
				ID   string `json:"id"`
				Name string `json:"name"`
			} `json:"workflows"`
			Count int `json:"count"`
		} `json:"result"`
		Event any `json:"event"`
	}
	if err := json.Unmarshal([]byte(atLimit.stdout), &success); err != nil || success.Event != nil || success.Result.Count != 10000 || len(success.Result.Workflows) != 10000 || success.Result.Workflows[0].Name != "limit workflow 00000" || success.Result.Workflows[len(success.Result.Workflows)-1].Name != "limit workflow 09999" {
		t.Fatalf("at-limit workflow.list = %+v, err = %v; want 10000 complete ordered rows and null event", success, err)
	}

	insertWorkflowListFixture(t, conn, 10000)
	journalPath := mustJournalPath(t, workspace)
	before := mustReadCommandFile(t, journalPath)
	overLimit := executeCommandForTest(workspace, "--json", "workflow", "list")
	if overLimit.err == nil || overLimit.stdout != "" {
		t.Fatalf("over-limit workflow.list = %#v, want stderr-only rejection", overLimit)
	}
	var failure struct {
		OK        bool   `json:"ok"`
		Operation string `json:"operation"`
		Error     struct {
			Code     string `json:"code"`
			Category string `json:"category"`
			Subject  any    `json:"subject"`
			Details  struct {
				Collection string `json:"collection"`
				Limit      int    `json:"limit"`
				Observed   int    `json:"observed"`
			} `json:"details"`
		} `json:"error"`
	}
	if err := json.Unmarshal([]byte(overLimit.stderr), &failure); err != nil || failure.OK || failure.Operation != "workflow.list" || failure.Error.Code != "result-limit-exceeded" || failure.Error.Category != "integrity" || failure.Error.Subject != nil || failure.Error.Details.Collection != "workflows" || failure.Error.Details.Limit != 10000 || failure.Error.Details.Observed != 10001 {
		t.Fatalf("over-limit failure = %+v, err = %v; want closed result-limit-exceeded envelope", failure, err)
	}
	if after := mustReadCommandFile(t, journalPath); string(after) != string(before) {
		t.Fatal("over-limit workflow list appended a journal event")
	}
	var count int
	if err := conn.QueryRow(`SELECT COUNT(*) FROM workflows`).Scan(&count); err != nil || count != 10001 {
		t.Fatalf("workflows after over-limit list = %d, err = %v; want 10001 without truncation", count, err)
	}
}

func resetWorkflowListFixture(t *testing.T, conn *sql.DB) {
	t.Helper()
	for _, table := range []string{"workflow_transition_requires", "workflow_transitions", "workflow_states", "workflow_assignments", "workflows"} {
		if _, err := conn.Exec("DELETE FROM " + table); err != nil {
			t.Fatalf("clear %s: %v", table, err)
		}
	}
}

func insertWorkflowListFixture(t *testing.T, conn *sql.DB, index int) {
	t.Helper()
	id := fmt.Sprintf("limit-%05d", index)
	name := fmt.Sprintf("limit workflow %05d", index)
	if _, err := conn.Exec(`INSERT INTO workflows (id, name, applies_to, active) VALUES (?, ?, 'task', 1)`, id, name); err != nil {
		t.Fatalf("insert workflow %q: %v", id, err)
	}
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
