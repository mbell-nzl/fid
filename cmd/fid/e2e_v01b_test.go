//go:build integration

package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestV01BSmokeCreateListGetUpdateTransitionReadyAndRebuildParity(t *testing.T) {
	workspace := t.TempDir()
	runInitCommand(t, workspace, "init")

	_, seqEpicStderr := runFIDCommand(
		t,
		workspace,
		"issue", "create",
		"--kind", "epic",
		"--title", "Sequential delivery epic",
		"--scheduling", "sequential",
		"--by", "tl",
	)
	seqEpicID := parseCreatedIssueID(t, seqEpicStderr)

	_, seqCurrentStderr := runFIDCommand(
		t,
		workspace,
		"issue", "create",
		"--kind", "task",
		"--title", "Draft implementation plan",
		"--epic", seqEpicID,
		"--sequence", "1",
		"--priority", "2",
		"--est-tokens", "1200",
		"--by", "developer",
	)
	seqCurrentID := parseCreatedIssueID(t, seqCurrentStderr)

	_, seqNextStderr := runFIDCommand(
		t,
		workspace,
		"issue", "create",
		"--kind", "task",
		"--title", "Ship follow-up work",
		"--epic", seqEpicID,
		"--sequence", "2",
		"--priority", "1",
		"--est-tokens", "900",
		"--by", "developer",
	)
	seqNextID := parseCreatedIssueID(t, seqNextStderr)

	_, priEpicStderr := runFIDCommand(
		t,
		workspace,
		"issue", "create",
		"--kind", "epic",
		"--title", "Prioritised incident epic",
		"--scheduling", "prioritised",
		"--by", "tl",
	)
	priEpicID := parseCreatedIssueID(t, priEpicStderr)

	_, priLowStderr := runFIDCommand(
		t,
		workspace,
		"issue", "create",
		"--kind", "task",
		"--title", "Clean up backlog fallout",
		"--epic", priEpicID,
		"--priority", "2",
		"--est-tokens", "700",
		"--by", "developer",
	)
	priLowID := parseCreatedIssueID(t, priLowStderr)

	_, priHighStderr := runFIDCommand(
		t,
		workspace,
		"issue", "create",
		"--kind", "task",
		"--title", "Patch customer outage",
		"--epic", priEpicID,
		"--priority", "0",
		"--est-tokens", "500",
		"--by", "developer",
	)
	priHighID := parseCreatedIssueID(t, priHighStderr)

	_, subtaskStderr := runFIDCommand(
		t,
		workspace,
		"issue", "create",
		"--kind", "subtask",
		"--title", "Capture parity checks",
		"--parent", seqCurrentID,
		"--by", "developer",
	)
	subtaskID := parseCreatedIssueID(t, subtaskStderr)

	seqListStdout, seqListStderr := runFIDCommand(
		t,
		workspace,
		"issue", "list",
		"--epic", seqEpicID,
		"--kind", "task",
		"--order", "sequence",
		"--format", "json",
	)
	if seqListStderr != "" {
		t.Fatalf("sequential list stderr = %q, want empty", seqListStderr)
	}
	if gotIDs := commandIssueListIDs(decodeCommandJSONStructure(t, seqListStdout).([]any)); !equalIssueIDs(gotIDs, []string{seqCurrentID, seqNextID}) {
		t.Fatalf("sequential list ids = %#v, want %#v", gotIDs, []string{seqCurrentID, seqNextID})
	}

	priListStdout, priListStderr := runFIDCommand(
		t,
		workspace,
		"issue", "list",
		"--epic", priEpicID,
		"--kind", "task",
		"--order", "priority",
		"--format", "json",
	)
	if priListStderr != "" {
		t.Fatalf("prioritised list stderr = %q, want empty", priListStderr)
	}
	if gotIDs := commandIssueListIDs(decodeCommandJSONStructure(t, priListStdout).([]any)); !equalIssueIDs(gotIDs, []string{priHighID, priLowID}) {
		t.Fatalf("prioritised list ids = %#v, want %#v", gotIDs, []string{priHighID, priLowID})
	}

	assertIssueReadBack(t, workspace, seqEpicID, issueReadBack{
		Kind:           "epic",
		Title:          "Sequential delivery epic",
		Status:         "todo",
		StatusCategory: "unstarted",
		WorkflowID:     "epic-default",
		OwnerRole:      "tl",
	})
	assertIssueReadBack(t, workspace, seqCurrentID, issueReadBack{
		Kind:           "task",
		ParentID:       seqEpicID,
		Title:          "Draft implementation plan",
		Status:         "todo",
		StatusCategory: "unstarted",
		WorkflowID:     "task-default",
		OwnerRole:      "developer",
		Sequence:       intPointer(1),
		Priority:       intPointer(2),
		EstTokens:      intPointer(1200),
	})
	assertIssueReadBack(t, workspace, priHighID, issueReadBack{
		Kind:           "task",
		ParentID:       priEpicID,
		Title:          "Patch customer outage",
		Status:         "todo",
		StatusCategory: "unstarted",
		WorkflowID:     "task-default",
		OwnerRole:      "developer",
		Priority:       intPointer(0),
		EstTokens:      intPointer(500),
	})
	assertIssueReadBack(t, workspace, subtaskID, issueReadBack{
		Kind:           "subtask",
		ParentID:       seqCurrentID,
		Title:          "Capture parity checks",
		Status:         "todo",
		StatusCategory: "unstarted",
		WorkflowID:     "subtask-default",
		OwnerRole:      "developer",
	})

	_, updateStderr := runFIDCommand(
		t,
		workspace,
		"issue", "update", seqCurrentID,
		"--title", "Draft implementation plan v2",
		"--priority", "1",
		"--est-tokens", "1500",
		"--by", "developer",
	)
	for _, want := range []string{
		"updated " + seqCurrentID + ":",
		`title "Draft implementation plan" -> "Draft implementation plan v2"`,
		"priority 2 -> 1",
		"est-tokens 1200 -> 1500",
	} {
		if !strings.Contains(updateStderr, want) {
			t.Fatalf("update stderr = %q, want substring %q", updateStderr, want)
		}
	}
	assertIssueReadBack(t, workspace, seqCurrentID, issueReadBack{
		Kind:           "task",
		ParentID:       seqEpicID,
		Title:          "Draft implementation plan v2",
		Status:         "todo",
		StatusCategory: "unstarted",
		WorkflowID:     "task-default",
		OwnerRole:      "developer",
		Sequence:       intPointer(1),
		Priority:       intPointer(1),
		EstTokens:      intPointer(1500),
	})

	_, transitionStderr := runFIDCommand(
		t,
		workspace,
		"issue", "transition", seqCurrentID,
		"--to", "cancelled",
		"--by", "tl",
	)
	if !strings.Contains(transitionStderr, "transitioned "+seqCurrentID+": todo → cancelled") {
		t.Fatalf("transition stderr = %q, want cancellation confirmation", transitionStderr)
	}
	assertIssueReadBack(t, workspace, seqCurrentID, issueReadBack{
		Kind:           "task",
		ParentID:       seqEpicID,
		Title:          "Draft implementation plan v2",
		Status:         "cancelled",
		StatusCategory: "cancelled",
		WorkflowID:     "task-default",
		OwnerRole:      "developer",
		Sequence:       intPointer(1),
		Priority:       intPointer(1),
		EstTokens:      intPointer(1500),
	})

	seqReadyStdout, seqReadyStderr := runFIDCommand(t, workspace, "ready", "--epic", seqEpicID, "--format", "json")
	if seqReadyStderr != "" {
		t.Fatalf("sequential ready stderr = %q, want empty", seqReadyStderr)
	}
	assertReadyTaskID(t, seqReadyStdout, seqNextID)

	priReadyStdout, priReadyStderr := runFIDCommand(t, workspace, "ready", "--epic", priEpicID, "--format", "json")
	if priReadyStderr != "" {
		t.Fatalf("prioritised ready stderr = %q, want empty", priReadyStderr)
	}
	assertReadyTaskID(t, priReadyStdout, priHighID)

	projectReadyStdout, projectReadyStderr := runFIDCommand(t, workspace, "ready", "--format", "json")
	if projectReadyStderr != "" {
		t.Fatalf("project ready stderr = %q, want empty", projectReadyStderr)
	}
	assertReadyTaskID(t, projectReadyStdout, priHighID)

	beforeIssueJSON := map[string]string{
		seqEpicID:    mustGetIssueJSON(t, workspace, seqEpicID),
		seqCurrentID: mustGetIssueJSON(t, workspace, seqCurrentID),
		seqNextID:    mustGetIssueJSON(t, workspace, seqNextID),
		priEpicID:    mustGetIssueJSON(t, workspace, priEpicID),
		priLowID:     mustGetIssueJSON(t, workspace, priLowID),
		priHighID:    mustGetIssueJSON(t, workspace, priHighID),
		subtaskID:    mustGetIssueJSON(t, workspace, subtaskID),
	}

	finalSeqListStdout, finalSeqListStderr := runFIDCommand(
		t,
		workspace,
		"issue", "list",
		"--epic", seqEpicID,
		"--kind", "task",
		"--order", "sequence",
		"--format", "json",
	)
	if finalSeqListStderr != "" {
		t.Fatalf("final sequential list stderr = %q, want empty", finalSeqListStderr)
	}

	finalPriListStdout, finalPriListStderr := runFIDCommand(
		t,
		workspace,
		"issue", "list",
		"--epic", priEpicID,
		"--kind", "task",
		"--order", "priority",
		"--format", "json",
	)
	if finalPriListStderr != "" {
		t.Fatalf("final prioritised list stderr = %q, want empty", finalPriListStderr)
	}

	beforeCommandOutputs := map[string]string{
		"seq_list":      finalSeqListStdout,
		"pri_list":      finalPriListStdout,
		"seq_ready":     seqReadyStdout,
		"pri_ready":     priReadyStdout,
		"project_ready": projectReadyStdout,
	}

	dbPath := filepath.Join(workspace, ".fid", "fid.sqlite")
	if err := os.Remove(dbPath); err != nil {
		t.Fatalf("Remove(%q) error = %v", dbPath, err)
	}

	runFIDCommand(t, workspace, "rebuild", "--by", "workflow_runner")

	for issueID, wantJSON := range beforeIssueJSON {
		if gotJSON := mustGetIssueJSON(t, workspace, issueID); gotJSON != wantJSON {
			t.Fatalf("issue JSON after rebuild for %s = %q, want %q", issueID, gotJSON, wantJSON)
		}
	}

	afterSeqListStdout, afterSeqListStderr := runFIDCommand(
		t,
		workspace,
		"issue", "list",
		"--epic", seqEpicID,
		"--kind", "task",
		"--order", "sequence",
		"--format", "json",
	)
	if afterSeqListStderr != "" {
		t.Fatalf("sequential list after rebuild stderr = %q, want empty", afterSeqListStderr)
	}
	if afterSeqListStdout != beforeCommandOutputs["seq_list"] {
		t.Fatalf("sequential list after rebuild = %q, want %q", afterSeqListStdout, beforeCommandOutputs["seq_list"])
	}

	afterPriListStdout, afterPriListStderr := runFIDCommand(
		t,
		workspace,
		"issue", "list",
		"--epic", priEpicID,
		"--kind", "task",
		"--order", "priority",
		"--format", "json",
	)
	if afterPriListStderr != "" {
		t.Fatalf("prioritised list after rebuild stderr = %q, want empty", afterPriListStderr)
	}
	if afterPriListStdout != beforeCommandOutputs["pri_list"] {
		t.Fatalf("prioritised list after rebuild = %q, want %q", afterPriListStdout, beforeCommandOutputs["pri_list"])
	}

	for key, args := range map[string][]string{
		"seq_ready":     {"ready", "--epic", seqEpicID, "--format", "json"},
		"pri_ready":     {"ready", "--epic", priEpicID, "--format", "json"},
		"project_ready": {"ready", "--format", "json"},
	} {
		stdout, stderr := runFIDCommand(t, workspace, args...)
		if stderr != "" {
			t.Fatalf("%s stderr after rebuild = %q, want empty", key, stderr)
		}
		if stdout != beforeCommandOutputs[key] {
			t.Fatalf("%s stdout after rebuild = %q, want %q", key, stdout, beforeCommandOutputs[key])
		}
	}
}

func TestV01BReadyReturnsEmptyAfterAllTasksFinishedEndToEnd(t *testing.T) {
	workspace := t.TempDir()
	runInitCommand(t, workspace, "init")

	_, createStderr := runFIDCommand(
		t,
		workspace,
		"issue", "create",
		"--kind", "task",
		"--title", "Only actionable task",
		"--priority", "1",
		"--by", "developer",
	)
	issueID := parseCreatedIssueID(t, createStderr)

	beforeReadyStdout, beforeReadyStderr := runFIDCommand(t, workspace, "ready", "--format", "json")
	if beforeReadyStderr != "" {
		t.Fatalf("ready before completion stderr = %q, want empty", beforeReadyStderr)
	}
	assertReadyTaskID(t, beforeReadyStdout, issueID)

	_, transitionStderr := runFIDCommand(
		t,
		workspace,
		"issue", "transition", issueID,
		"--to", "cancelled",
		"--by", "tl",
	)
	if !strings.Contains(transitionStderr, "transitioned "+issueID+": todo → cancelled") {
		t.Fatalf("transition stderr = %q, want cancellation confirmation", transitionStderr)
	}

	emptyReadyStdout, emptyReadyStderr := runFIDCommand(t, workspace, "ready", "--format", "json")
	if strings.TrimSpace(emptyReadyStdout) != `{"task":null}` {
		t.Fatalf("empty ready stdout = %q, want exact null payload", emptyReadyStdout)
	}
	if !strings.Contains(emptyReadyStderr, "note: no ready tasks in project") {
		t.Fatalf("empty ready stderr = %q, want project note", emptyReadyStderr)
	}
}

type issueReadBack struct {
	Kind           string
	ParentID       string
	Title          string
	Status         string
	StatusCategory string
	WorkflowID     string
	OwnerRole      string
	Sequence       *int
	Priority       *int
	EstTokens      *int
}

func assertIssueReadBack(t *testing.T, workspace, issueID string, want issueReadBack) {
	t.Helper()

	var got struct {
		ID             string  `json:"id"`
		Kind           string  `json:"kind"`
		ParentID       *string `json:"parent_id"`
		Title          string  `json:"title"`
		Status         string  `json:"status"`
		StatusCategory string  `json:"status_category"`
		WorkflowID     string  `json:"workflow_id"`
		OwnerRole      *string `json:"owner_role"`
		Sequence       *int    `json:"sequence"`
		Priority       *int    `json:"priority"`
		EstTokens      *int    `json:"est_tokens"`
	}

	if err := json.Unmarshal([]byte(mustGetIssueJSON(t, workspace, issueID)), &got); err != nil {
		t.Fatalf("Unmarshal(%s) error = %v", issueID, err)
	}

	if got.ID != issueID {
		t.Fatalf("issue %s id = %q, want %q", issueID, got.ID, issueID)
	}
	if got.Kind != want.Kind {
		t.Fatalf("issue %s kind = %q, want %q", issueID, got.Kind, want.Kind)
	}
	if want.ParentID == "" {
		if got.ParentID != nil {
			t.Fatalf("issue %s parent_id = %#v, want nil", issueID, got.ParentID)
		}
	} else if got.ParentID == nil || *got.ParentID != want.ParentID {
		t.Fatalf("issue %s parent_id = %#v, want %q", issueID, got.ParentID, want.ParentID)
	}
	if got.Title != want.Title {
		t.Fatalf("issue %s title = %q, want %q", issueID, got.Title, want.Title)
	}
	if got.Status != want.Status {
		t.Fatalf("issue %s status = %q, want %q", issueID, got.Status, want.Status)
	}
	if got.StatusCategory != want.StatusCategory {
		t.Fatalf("issue %s status_category = %q, want %q", issueID, got.StatusCategory, want.StatusCategory)
	}
	if got.WorkflowID != want.WorkflowID {
		t.Fatalf("issue %s workflow_id = %q, want %q", issueID, got.WorkflowID, want.WorkflowID)
	}
	if got.OwnerRole == nil || *got.OwnerRole != want.OwnerRole {
		t.Fatalf("issue %s owner_role = %#v, want %q", issueID, got.OwnerRole, want.OwnerRole)
	}
	assertOptionalInt(t, issueID, "sequence", got.Sequence, want.Sequence)
	assertOptionalInt(t, issueID, "priority", got.Priority, want.Priority)
	assertOptionalInt(t, issueID, "est_tokens", got.EstTokens, want.EstTokens)
}

func assertOptionalInt(t *testing.T, issueID, field string, got, want *int) {
	t.Helper()

	if want == nil {
		if got != nil {
			t.Fatalf("issue %s %s = %#v, want nil", issueID, field, got)
		}
		return
	}
	if got == nil || *got != *want {
		t.Fatalf("issue %s %s = %#v, want %d", issueID, field, got, *want)
	}
}

func assertReadyTaskID(t *testing.T, stdout, wantTaskID string) {
	t.Helper()

	payload := decodeCommandJSONStructure(t, stdout).(map[string]any)
	task, ok := payload["task"].(map[string]any)
	if !ok {
		t.Fatalf("ready payload = %#v, want task object", payload)
	}
	if gotTaskID, _ := task["id"].(string); gotTaskID != wantTaskID {
		t.Fatalf("ready task id = %q, want %q", gotTaskID, wantTaskID)
	}
}
