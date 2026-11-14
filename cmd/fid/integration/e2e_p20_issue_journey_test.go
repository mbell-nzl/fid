//go:build integration

package integration

import (
	"bufio"
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/mbell-nzl/fid/internal/config"
	"github.com/mbell-nzl/fid/internal/gittest"
)

func TestP20StandaloneIssueJourneyJSONAndRebuildParity(t *testing.T) {
	t.Parallel()
	workspace := gittest.NewWorkspace(t)
	eventIDs := map[string]bool{}

	_, initEvent := p20JourneyMutation(t, workspace, "init", "--prefix", "P20")
	p20RememberEvent(t, eventIDs, initEvent)

	_, workflowEvent := p20JourneyMutation(t, workspace,
		"workflow", "define", "--name", "p20-task", "--applies-to", "subtask",
		"--states", "todo,in_progress,cancelled",
		"--categories", "unstarted,started,cancelled",
		"--transitions", "todo->in_progress,in_progress->cancelled", "--by", "tl",
	)
	p20RememberEvent(t, eventIDs, workflowEvent)
	_, assignmentEvent := p20JourneyMutation(t, workspace,
		"workflow", "assign", "--name", "p20-task", "--kind", "subtask", "--by", "tl",
	)
	p20RememberEvent(t, eventIDs, assignmentEvent)

	epicResult, epicEvent := p20JourneyMutation(t, workspace,
		"issue", "create", "--kind", "epic", "--title", "Phase 20 journey", "--scheduling", "sequential", "--actor", "agent:tl", "--idempotency-key", "p20-issue-create-epic",
	)
	p20AssertIdempotencyKey(t, epicResult, "p20-issue-create-epic")
	p20RememberEvent(t, eventIDs, epicEvent)
	epicID := p20IssueID(t, epicResult)

	taskResult, taskEvent := p20JourneyMutation(t, workspace,
		"issue", "create", "--kind", "task", "--title", "Exercise retained journey", "--epic", epicID,
		"--sequence", "1", "--priority", "1", "--actor", "agent:developer", "--idempotency-key", "p20-issue-create-task",
	)
	p20AssertIdempotencyKey(t, taskResult, "p20-issue-create-task")
	p20RememberEvent(t, eventIDs, taskEvent)
	taskID := p20IssueID(t, taskResult)

	followUpResult, followUpEvent := p20JourneyMutation(t, workspace,
		"issue", "create", "--kind", "task", "--title", "Follow-up retained journey", "--epic", epicID,
		"--sequence", "2", "--priority", "2", "--actor", "agent:developer", "--idempotency-key", "p20-issue-create-follow-up",
	)
	p20AssertIdempotencyKey(t, followUpResult, "p20-issue-create-follow-up")
	p20RememberEvent(t, eventIDs, followUpEvent)
	followUpID := p20IssueID(t, followUpResult)

	subtaskResult, subtaskEvent := p20JourneyMutation(t, workspace,
		"issue", "create", "--kind", "subtask", "--title", "Inspect canonical state", "--parent", taskID, "--actor", "agent:developer", "--idempotency-key", "p20-issue-create-subtask",
	)
	p20AssertIdempotencyKey(t, subtaskResult, "p20-issue-create-subtask")
	p20RememberEvent(t, eventIDs, subtaskEvent)
	subtaskID := p20IssueID(t, subtaskResult)

	for _, label := range []string{"zeta", "alpha"} {
		_, eventID := p20JourneyMutation(t, workspace, "label", "add", label, "--colour", "blue", "--by", "tl")
		p20RememberEvent(t, eventIDs, eventID)
	}
	for _, label := range []string{"zeta", "alpha"} {
		_, eventID := p20JourneyMutation(t, workspace, "issue", "label", "add", taskID, label, "--by", "developer")
		p20RememberEvent(t, eventIDs, eventID)
	}
	for _, criterion := range []struct {
		ordinal string
		text    string
	}{
		{ordinal: "2", text: "Second check is retained"},
		{ordinal: "1", text: "First check is retained"},
	} {
		_, eventID := p20JourneyMutation(t, workspace,
			"issue", "ac", "add", taskID, "--ordinal", criterion.ordinal, "--text", criterion.text, "--by", "developer",
		)
		p20RememberEvent(t, eventIDs, eventID)
	}
	for _, link := range []struct {
		relation string
		toID     string
	}{
		{relation: "relates_to", toID: subtaskID},
		{relation: "discovered_from", toID: followUpID},
	} {
		_, eventID := p20JourneyMutation(t, workspace,
			"link", "add", taskID, link.relation, link.toID, "--from-type", "issue", "--to-type", "issue", "--by", "developer",
		)
		p20RememberEvent(t, eventIDs, eventID)
	}

	issueList := p20JourneyRead(t, workspace, "issue", "list", "--epic", epicID, "--kind", "task", "--order", "sequence")
	p20AssertIDs(t, p20Objects(t, issueList["issues"]), []string{taskID, followUpID})
	p20AssertCount(t, issueList, 2)

	labelList := p20JourneyRead(t, workspace, "label", "list")
	p20AssertNames(t, p20Objects(t, labelList["labels"]), []string{"alpha", "zeta"})
	p20AssertCount(t, labelList, 2)

	acList := p20JourneyRead(t, workspace, "issue", "ac", "list", taskID)
	p20AssertOrdinals(t, p20Objects(t, acList["acceptance_criteria"]), []int{1, 2})

	linkList := p20JourneyRead(t, workspace, "link", "list", taskID, "--type", "issue")
	p20AssertDirections(t, p20Objects(t, linkList["links"]), []string{"outgoing", "outgoing", "incoming"})
	p20AssertRelations(t, p20Objects(t, linkList["links"]), []string{"discovered_from", "relates_to", "relates_to"})
	p20AssertCount(t, linkList, 3)

	ready := p20JourneyRead(t, workspace, "ready", "--epic", epicID)
	if got := p20Object(t, ready["task"])["id"]; got != taskID {
		t.Fatalf("ready task id = %q, want %q", got, taskID)
	}

	workflowList := p20JourneyRead(t, workspace, "workflow", "list", "--kind", "subtask")
	if !p20ContainsName(p20Objects(t, workflowList["workflows"]), "p20-task") {
		t.Fatalf("workflow list = %#v, want p20-task", workflowList)
	}

	beforeUpdate := p20IssueGet(t, workspace, taskID)
	currentRevision := p20IssueRevision(t, beforeUpdate)
	updateResult, updateEvent := p20JourneyMutation(t, workspace,
		"issue", "update", taskID, "--title", "Exercise retained journey, updated", "--expected-revision", strconv.Itoa(currentRevision), "--actor", "agent:developer",
	)
	p20RememberEvent(t, eventIDs, updateEvent)
	updatedRevision := p20Integer(t, p20Object(t, updateResult["issue"])["revision"])
	if updatedRevision != currentRevision+1 {
		t.Fatalf("updated revision = %d, want %d", updatedRevision, currentRevision+1)
	}
	p20AssertJournalEventAttribution(t, workspace, updateEvent, "developer")

	transitionResult, transitionEvent := p20JourneyMutation(t, workspace,
		"issue", "transition", taskID, "--to", "cancelled", "--reason", "Test cancellation", "--expected-revision", strconv.Itoa(updatedRevision), "--actor", "human:mark",
	)
	p20RememberEvent(t, eventIDs, transitionEvent)
	transitionedRevision := p20Integer(t, p20Object(t, transitionResult["issue"])["revision"])
	if transitionedRevision != updatedRevision+1 {
		t.Fatalf("transitioned revision = %d, want %d", transitionedRevision, updatedRevision+1)
	}

	p20AssertStaleRevisionFailure(t, workspace, taskID, updatedRevision, transitionedRevision)

	check := p20JourneyRead(t, workspace, "maintenance", "check")
	if status := check["status"]; status != "pass" && status != "warn" {
		t.Fatalf("maintenance check status = %q, want pass or warn without a divergence failure", status)
	}

	beforeRebuild := map[string]string{}
	for _, issueID := range []string{epicID, taskID, followUpID, subtaskID} {
		beforeRebuild[issueID] = p20IssueGet(t, workspace, issueID)
	}
	p20JourneyRead(t, workspace, "rebuild")
	for issueID, want := range beforeRebuild {
		if got := p20IssueGet(t, workspace, issueID); got != want {
			t.Fatalf("canonical issue.get bytes after rebuild for %s differ\n got: %q\nwant: %q", issueID, got, want)
		}
	}
}

type p20JourneyResponse struct {
	operation string
	result    json.RawMessage
	eventID   string
	raw       string
}

func p20JourneyMutation(t *testing.T, workspace string, args ...string) (map[string]any, string) {
	t.Helper()
	response := p20JourneyJSON(t, workspace, args...)
	if response.eventID == "" {
		t.Fatalf("fid --json %v event = nil, want mutation Event identity", args)
	}
	return p20JourneyResult(t, response), response.eventID
}

func p20JourneyRead(t *testing.T, workspace string, args ...string) map[string]any {
	t.Helper()
	response := p20JourneyJSON(t, workspace, args...)
	if response.eventID != "" {
		t.Fatalf("fid --json %v event = %q, want read event:null", args, response.eventID)
	}
	return p20JourneyResult(t, response)
}

func p20JourneyJSON(t *testing.T, workspace string, args ...string) p20JourneyResponse {
	t.Helper()
	stdout, stderr, err := closureRunFIDExpectError(
		t, integrationFIDBinary(t), workspace, append([]string{"--json"}, args...)...,
	)
	if err != nil || stderr != "" {
		t.Fatalf("fid --json %v = (%v, stdout %q, stderr %q), want stdout-only success", args, err, stdout, stderr)
	}
	if len(stdout) == 0 || stdout[len(stdout)-1] != '\n' {
		t.Fatalf("fid --json %v stdout = %q, want one JSON response with trailing newline", args, stdout)
	}
	var decoded struct {
		ContractVersion string          `json:"contract_version"`
		OK              bool            `json:"ok"`
		Operation       string          `json:"operation"`
		Result          json.RawMessage `json:"result"`
		Event           *struct {
			ID string `json:"id"`
		} `json:"event"`
	}
	if err := json.Unmarshal([]byte(stdout), &decoded); err != nil {
		t.Fatalf("decode fid --json %v: %v", args, err)
	}
	if decoded.ContractVersion != "fid-cli.v1" || !decoded.OK || decoded.Operation == "" || len(decoded.Result) == 0 {
		t.Fatalf("fid --json %v response = %s, want successful fid-cli.v1 envelope", args, stdout)
	}
	eventID := ""
	if decoded.Event != nil {
		eventID = decoded.Event.ID
	}
	return p20JourneyResponse{operation: decoded.Operation, result: decoded.Result, eventID: eventID, raw: stdout}
}

func p20JourneyResult(t *testing.T, response p20JourneyResponse) map[string]any {
	t.Helper()
	decoder := json.NewDecoder(bytes.NewReader(response.result))
	decoder.UseNumber()
	var result map[string]any
	if err := decoder.Decode(&result); err != nil {
		t.Fatalf("decode %s result: %v", response.operation, err)
	}
	return result
}

func p20IssueGet(t *testing.T, workspace, issueID string) string {
	t.Helper()
	response := p20JourneyJSON(t, workspace, "issue", "get", issueID)
	if response.operation != "issue.get" || response.eventID != "" {
		t.Fatalf("issue get response = %#v, want issue.get read envelope", response)
	}
	return response.raw
}

func p20IssueRevision(t *testing.T, issueJSON string) int {
	t.Helper()
	response := p20JourneyResponse{operation: "issue.get"}
	var decoded struct {
		Result json.RawMessage `json:"result"`
	}
	if err := json.Unmarshal([]byte(issueJSON), &decoded); err != nil {
		t.Fatalf("decode issue.get: %v", err)
	}
	response.result = decoded.Result
	return p20Integer(t, p20Object(t, p20JourneyResult(t, response)["issue"])["revision"])
}

func p20IssueID(t *testing.T, result map[string]any) string {
	t.Helper()
	id, ok := p20Object(t, result["issue"])["id"].(string)
	if !ok || id == "" {
		t.Fatalf("issue result = %#v, want non-empty issue id", result)
	}
	return id
}

func p20AssertIdempotencyKey(t *testing.T, result map[string]any, want string) {
	t.Helper()
	if got, _ := result["idempotency_key"].(string); got != want {
		t.Fatalf("idempotency key = %q, want %q", got, want)
	}
}

func p20RememberEvent(t *testing.T, eventIDs map[string]bool, eventID string) {
	t.Helper()
	if eventID == "" || eventIDs[eventID] {
		t.Fatalf("mutation event id = %q, want a new non-empty identity", eventID)
	}
	eventIDs[eventID] = true
}

func p20AssertStaleRevisionFailure(t *testing.T, workspace, issueID string, expected, actual int) {
	t.Helper()
	stdout, stderr, err := closureRunFIDExpectError(t, integrationFIDBinary(t), workspace,
		"--json", "issue", "update", issueID, "--title", "stale update must not apply", "--expected-revision", strconv.Itoa(expected), "--actor", "agent:developer",
	)
	if err == nil || stdout != "" {
		t.Fatalf("stale update = (%v, stdout %q, stderr %q), want stderr-only typed failure", err, stdout, stderr)
	}
	if len(stderr) == 0 || stderr[len(stderr)-1] != '\n' {
		t.Fatalf("stale update stderr = %q, want JSON failure with trailing newline", stderr)
	}
	var failure struct {
		ContractVersion string `json:"contract_version"`
		OK              bool   `json:"ok"`
		Operation       string `json:"operation"`
		Error           struct {
			Code    string `json:"code"`
			Details struct {
				Expected int `json:"expected_revision"`
				Actual   int `json:"actual_revision"`
			} `json:"details"`
		} `json:"error"`
	}
	if err := json.Unmarshal([]byte(stderr), &failure); err != nil {
		t.Fatalf("decode stale state-conflict: %v", err)
	}
	if failure.ContractVersion != "fid-cli.v1" || failure.OK || failure.Operation != "issue.update" || failure.Error.Code != "state-conflict" {
		t.Fatalf("stale failure = %#v, want issue.update state-conflict envelope", failure)
	}
	if failure.Error.Details.Expected != expected || failure.Error.Details.Actual != actual {
		t.Fatalf("state-conflict details = %#v, want expected %d actual %d", failure.Error.Details, expected, actual)
	}
}

func p20AssertJournalEventAttribution(t *testing.T, workspace, eventID, actorID string) {
	t.Helper()
	paths, err := filepath.Glob(filepath.Join(config.LedgerDir(workspace), "journal", "*.jsonl"))
	if err != nil || len(paths) != 1 {
		t.Fatalf("journal files = %v, %v; want one journal", paths, err)
	}
	file, err := os.Open(paths[0])
	if err != nil {
		t.Fatalf("open journal: %v", err)
	}
	defer func() { _ = file.Close() }()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var event struct {
			ID     string `json:"event_id"`
			Schema int    `json:"event_schema"`
			Actor  struct {
				ActorID       string  `json:"actor_id"`
				ActorType     string  `json:"actor_type"`
				Producer      string  `json:"producer"`
				ProvenanceRef string  `json:"provenance_ref"`
				SessionRef    *string `json:"session_ref"`
			} `json:"actor"`
		}
		if err := json.Unmarshal(scanner.Bytes(), &event); err != nil {
			t.Fatalf("decode journal event: %v", err)
		}
		if event.ID == eventID {
			if event.Schema != 2 || event.Actor.ActorID != actorID || event.Actor.ActorType != "agent" || event.Actor.Producer != "fid-standalone" || event.Actor.ProvenanceRef != "project-local-attribution" || event.Actor.SessionRef != nil {
				t.Fatalf("update Event-v2 attribution = %+v, want standalone agent %q with project-local provenance", event, actorID)
			}
			return
		}
	}
	if err := scanner.Err(); err != nil {
		t.Fatalf("read journal: %v", err)
	}
	t.Fatalf("journal has no event %q", eventID)
}

func p20Object(t *testing.T, value any) map[string]any {
	t.Helper()
	object, ok := value.(map[string]any)
	if !ok {
		t.Fatalf("value = %#v, want object", value)
	}
	return object
}

func p20Objects(t *testing.T, value any) []map[string]any {
	t.Helper()
	items, ok := value.([]any)
	if !ok {
		t.Fatalf("value = %#v, want array", value)
	}
	objects := make([]map[string]any, 0, len(items))
	for _, item := range items {
		objects = append(objects, p20Object(t, item))
	}
	return objects
}

func p20Integer(t *testing.T, value any) int {
	t.Helper()
	number, ok := value.(json.Number)
	if !ok {
		t.Fatalf("value = %#v, want JSON integer", value)
	}
	result, err := strconv.Atoi(number.String())
	if err != nil {
		t.Fatalf("parse JSON integer %q: %v", number, err)
	}
	return result
}

func p20AssertIDs(t *testing.T, objects []map[string]any, want []string) {
	t.Helper()
	if len(objects) != len(want) {
		t.Fatalf("collection length = %d, want %d", len(objects), len(want))
	}
	for index, wantID := range want {
		if got := objects[index]["id"]; got != wantID {
			t.Fatalf("collection id at %d = %q, want %q", index, got, wantID)
		}
	}
}

func p20AssertNames(t *testing.T, objects []map[string]any, want []string) {
	t.Helper()
	for index, wantName := range want {
		if index >= len(objects) || objects[index]["name"] != wantName {
			t.Fatalf("collection names = %#v, want ordered %v", objects, want)
		}
	}
}

func p20AssertOrdinals(t *testing.T, objects []map[string]any, want []int) {
	t.Helper()
	if len(objects) != len(want) {
		t.Fatalf("criteria length = %d, want %d", len(objects), len(want))
	}
	for index, wantOrdinal := range want {
		if got := p20Integer(t, objects[index]["ordinal"]); got != wantOrdinal {
			t.Fatalf("criterion ordinal at %d = %d, want %d", index, got, wantOrdinal)
		}
	}
}

func p20AssertRelations(t *testing.T, objects []map[string]any, want []string) {
	t.Helper()
	if len(objects) != len(want) {
		t.Fatalf("link length = %d, want %d", len(objects), len(want))
	}
	for index, wantRelation := range want {
		if got := objects[index]["relation"]; got != wantRelation {
			t.Fatalf("link relation at %d = %q, want %q", index, got, wantRelation)
		}
	}
}

func p20AssertDirections(t *testing.T, objects []map[string]any, want []string) {
	t.Helper()
	if len(objects) != len(want) {
		t.Fatalf("link length = %d, want %d", len(objects), len(want))
	}
	for index, wantDirection := range want {
		if got := objects[index]["direction"]; got != wantDirection {
			t.Fatalf("link direction at %d = %q, want %q", index, got, wantDirection)
		}
	}
}

func p20AssertCount(t *testing.T, result map[string]any, want int) {
	t.Helper()
	if got := p20Integer(t, result["count"]); got != want {
		t.Fatalf("collection count = %d, want %d", got, want)
	}
}

func p20ContainsName(objects []map[string]any, name string) bool {
	for _, object := range objects {
		if object["name"] == name {
			return true
		}
	}
	return false
}
