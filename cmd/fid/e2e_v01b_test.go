package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	fiddb "github.com/mbell-nzl/fid/internal/db"
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
		"--json", "issue", "list",
		"--epic", seqEpicID,
		"--kind", "task",
		"--order", "sequence",
	)
	if seqListStderr != "" {
		t.Fatalf("sequential list stderr = %q, want empty", seqListStderr)
	}
	if gotIDs := commandIssueListIDs(decodeCanonicalResult(t, seqListStdout)["issues"].([]any)); !equalIssueIDs(gotIDs, []string{seqCurrentID, seqNextID}) {
		t.Fatalf("sequential list ids = %#v, want %#v", gotIDs, []string{seqCurrentID, seqNextID})
	}

	priListStdout, priListStderr := runFIDCommand(
		t,
		workspace,
		"--json", "issue", "list",
		"--epic", priEpicID,
		"--kind", "task",
		"--order", "priority",
	)
	if priListStderr != "" {
		t.Fatalf("prioritised list stderr = %q, want empty", priListStderr)
	}
	if gotIDs := commandIssueListIDs(decodeCanonicalResult(t, priListStdout)["issues"].([]any)); !equalIssueIDs(gotIDs, []string{priHighID, priLowID}) {
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
		"--by", "developer",
	)
	for _, want := range []string{
		"updated " + seqCurrentID + ":",
		`title "Draft implementation plan" -> "Draft implementation plan v2"`,
		"priority 2 -> 1",
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
	})

	seqReadyStdout, seqReadyStderr := runFIDCommand(t, workspace, "--json", "ready", "--epic", seqEpicID)
	if seqReadyStderr != "" {
		t.Fatalf("sequential ready stderr = %q, want empty", seqReadyStderr)
	}
	assertReadyTaskID(t, seqReadyStdout, seqNextID)

	priReadyStdout, priReadyStderr := runFIDCommand(t, workspace, "--json", "ready", "--epic", priEpicID)
	if priReadyStderr != "" {
		t.Fatalf("prioritised ready stderr = %q, want empty", priReadyStderr)
	}
	assertReadyTaskID(t, priReadyStdout, priHighID)

	projectReadyStdout, projectReadyStderr := runFIDCommand(t, workspace, "--json", "ready")
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
		"--json", "issue", "list",
		"--epic", seqEpicID,
		"--kind", "task",
		"--order", "sequence",
	)
	if finalSeqListStderr != "" {
		t.Fatalf("final sequential list stderr = %q, want empty", finalSeqListStderr)
	}

	finalPriListStdout, finalPriListStderr := runFIDCommand(
		t,
		workspace,
		"--json", "issue", "list",
		"--epic", priEpicID,
		"--kind", "task",
		"--order", "priority",
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

	runFIDCommand(t, workspace, "rebuild")

	for issueID, wantJSON := range beforeIssueJSON {
		if gotJSON := mustGetIssueJSON(t, workspace, issueID); gotJSON != wantJSON {
			t.Fatalf("issue JSON after rebuild for %s = %q, want %q", issueID, gotJSON, wantJSON)
		}
	}

	afterSeqListStdout, afterSeqListStderr := runFIDCommand(
		t,
		workspace,
		"--json", "issue", "list",
		"--epic", seqEpicID,
		"--kind", "task",
		"--order", "sequence",
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
		"--json", "issue", "list",
		"--epic", priEpicID,
		"--kind", "task",
		"--order", "priority",
	)
	if afterPriListStderr != "" {
		t.Fatalf("prioritised list after rebuild stderr = %q, want empty", afterPriListStderr)
	}
	if afterPriListStdout != beforeCommandOutputs["pri_list"] {
		t.Fatalf("prioritised list after rebuild = %q, want %q", afterPriListStdout, beforeCommandOutputs["pri_list"])
	}

	for key, args := range map[string][]string{
		"seq_ready":     {"--json", "ready", "--epic", seqEpicID},
		"pri_ready":     {"--json", "ready", "--epic", priEpicID},
		"project_ready": {"--json", "ready"},
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

	beforeReadyStdout, beforeReadyStderr := runFIDCommand(t, workspace, "--json", "ready")
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

	emptyReadyStdout, emptyReadyStderr := runFIDCommand(t, workspace, "--json", "ready")
	if decodeCanonicalResult(t, emptyReadyStdout)["task"] != nil {
		t.Fatalf("empty ready stdout = %q, want exact null payload", emptyReadyStdout)
	}
	if emptyReadyStderr != "" {
		t.Fatalf("empty ready stderr = %q, want empty", emptyReadyStderr)
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
	if got.OwnerRole != nil {
		t.Fatalf("issue %s canonical owner_role = %#v, want field absent", issueID, got.OwnerRole)
	}
	if gotRole := mustIssueOwnerRole(t, workspace, issueID); gotRole != want.OwnerRole {
		t.Fatalf("issue %s stored owner_role = %q, want %q", issueID, gotRole, want.OwnerRole)
	}
	assertOptionalInt(t, issueID, "sequence", got.Sequence, want.Sequence)
	assertOptionalInt(t, issueID, "priority", got.Priority, want.Priority)
}

func mustIssueOwnerRole(t *testing.T, workspace, issueID string) string {
	t.Helper()
	conn, err := fiddb.Open(filepath.Join(workspace, ".fid", "fid.sqlite"))
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	defer closeCommandDB(t, conn)
	var role string
	if err := conn.QueryRow("SELECT owner_role FROM issues WHERE id = ?", issueID).Scan(&role); err != nil {
		t.Fatalf("query owner_role for %s: %v", issueID, err)
	}
	return role
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

	payload := decodeCanonicalResult(t, stdout)
	task, ok := payload["task"].(map[string]any)
	if !ok {
		t.Fatalf("ready payload = %#v, want task object", payload)
	}
	if gotTaskID, _ := task["id"].(string); gotTaskID != wantTaskID {
		t.Fatalf("ready task id = %q, want %q", gotTaskID, wantTaskID)
	}
}
