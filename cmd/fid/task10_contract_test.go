package main

import (
	"encoding/json"
	"testing"
)

func TestTask10OperationsUseSelectedProjectOutsideWorkingDirectory(t *testing.T) {
	workspace := t.TempDir()
	clonePristineWorkspaceInto(t, workspace)
	outside := t.TempDir()

	assertTask10ProjectOperation(t, executeCommandForTest(outside, "--project", workspace, "--json", "label", "add", "core", "--by", "developer"), "label.add")
	assertTask10ProjectOperation(t, executeCommandForTest(outside, "--project", workspace, "--json", "label", "list"), "label.list")

	created := executeCommandForTest(outside, "--project", workspace, "--json", "issue", "create", "--kind", "task", "--title", "Selected project", "--by", "developer")
	if created.err != nil || created.stderr != "" {
		t.Fatalf("issue create = %#v, want selected-project success", created)
	}
	var createEnvelope struct {
		Result struct {
			Issue struct {
				ID string `json:"id"`
			} `json:"issue"`
		} `json:"result"`
	}
	if err := json.Unmarshal([]byte(created.stdout), &createEnvelope); err != nil || createEnvelope.Result.Issue.ID == "" {
		t.Fatalf("issue create envelope = %#v, error = %v", createEnvelope, err)
	}
	issueID := createEnvelope.Result.Issue.ID

	assertTask10ProjectOperation(t, executeCommandForTest(outside, "--project", workspace, "--json", "issue", "ac", "add", issueID, "--text", "Selected project AC", "--ordinal", "1", "--by", "developer"), "issue.ac.add")
	assertTask10ProjectOperation(t, executeCommandForTest(outside, "--project", workspace, "--json", "issue", "ac", "list", issueID), "issue.ac.list")
	assertTask10ProjectOperation(t, executeCommandForTest(outside, "--project", workspace, "--json", "issue", "label", "add", issueID, "core", "--by", "developer"), "issue.label.add")
	assertTask10ProjectOperation(t, executeCommandForTest(outside, "--project", workspace, "--json", "issue", "label", "remove", issueID, "core", "--by", "developer"), "issue.label.remove")
}

func TestIssueACVerifyRemainsGuardedLegacyLeaf(t *testing.T) {
	workspace := t.TempDir()
	clonePristineWorkspaceInto(t, workspace)
	issueID, criterionID := seedIssueAcceptanceCriterion(t, workspace, "Legacy verify guard")
	journalPath := mustJournalPath(t, workspace)
	before := len(splitJournalLines(mustReadCommandFile(t, journalPath)))

	result := executeCommandForTest(workspace, "--json", "issue", "ac", "verify", issueID, criterionID, "--by", "developer", "--status", "pass")
	if result.err == nil || result.stdout != "" || result.stderr == "" {
		t.Fatalf("result = %#v, want guarded legacy rejection", result)
	}
	if got := len(splitJournalLines(mustReadCommandFile(t, journalPath))); got != before {
		t.Fatalf("journal lines after guarded rejection = %d, want %d", got, before)
	}
}

func assertTask10ProjectOperation(t *testing.T, result commandExecutionResult, operation string) {
	t.Helper()
	if result.err != nil || result.stderr != "" {
		t.Fatalf("%s result = %#v, want selected-project JSON success", operation, result)
	}
	var envelope struct {
		ContractVersion string `json:"contract_version"`
		OK              bool   `json:"ok"`
		Operation       string `json:"operation"`
	}
	if err := json.Unmarshal([]byte(result.stdout), &envelope); err != nil {
		t.Fatalf("Unmarshal(%s) error = %v", operation, err)
	}
	if envelope.ContractVersion != "fid-cli.v1" || !envelope.OK || envelope.Operation != operation {
		t.Fatalf("%s envelope = %#v, want selected-project success", operation, envelope)
	}
}
