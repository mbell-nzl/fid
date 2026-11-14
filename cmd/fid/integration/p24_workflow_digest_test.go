//go:build integration

package integration

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"

	"github.com/mbell-nzl/fid/internal/gittest"
)

var p24WorkflowFields = map[string]string{"workflow.define": "--name", "workflow.assign": "--name"}

func TestP24WorkflowInputAsWrittenAndTwoFieldConflict(t *testing.T) {
	t.Parallel()
	binary := integrationFIDBinary(t)
	workspace := gittest.NewWorkspace(t)
	p21Run(t, binary, workspace, "--project", workspace, "--json", "init", "--prefix", "P24")
	define := []string{"--project", workspace, "--json", "workflow", "define", "--name", "task-p24", "--applies-to", "subtask", "--states", "todo,building,done,cancelled", "--categories", "unstarted,started,completed,cancelled", "--transitions", "todo->building,building->done", "--actor", "agent:developer", "--idempotency-key", "p24-workflow-define"}
	first := p24WorkflowCall(t, workspace, binary, define)
	if first.ExitCode != 0 || first.Stderr != "" {
		t.Fatalf("workflow.define first = %d %q", first.ExitCode, first.Stderr)
	}
	p24AssertWorkflowFields(t, workspace, "p24-workflow-define", []string{"--name", "--applies-to", "--states", "--categories", "--transitions", "--actor"})
	defineEvents := p21JournalEventCount(t, workspace)
	definePermuted := []string{"--project", workspace, "--json", "workflow", "define", "--actor", "agent:developer", "--transitions", "todo->building,building->done", "--categories", "unstarted,started,completed,cancelled", "--states", "todo,building,done,cancelled", "--applies-to", "subtask", "--idempotency-key", "p24-workflow-define", "--name", "task-p24"}
	defineReplay := p24WorkflowCall(t, workspace, binary, definePermuted)
	if defineReplay.ExitCode != 0 || defineReplay.Stderr != "" {
		t.Fatalf("workflow.define replay = %d %q", defineReplay.ExitCode, defineReplay.Stderr)
	}
	if want := strings.Replace(first.Stdout, `"replayed":false`, `"replayed":true`, 1); want == first.Stdout || defineReplay.Stdout != want {
		t.Fatal("workflow.define replay bytes changed beyond replayed")
	}
	defineWant := p23Envelope(t, first.Stdout)
	defineWant["result"].(map[string]any)["replayed"] = true
	if got := p23Envelope(t, defineReplay.Stdout); !reflect.DeepEqual(got, defineWant) {
		t.Fatalf("workflow.define replay = %#v, want %#v", got, defineWant)
	}
	defineChanged := append([]string(nil), define...)
	for i, arg := range defineChanged {
		if arg == "--name" {
			defineChanged[i+1] = "task-other"
		}
	}
	defineConflict := p24WorkflowCall(t, workspace, binary, defineChanged)
	assertP22ProcessFailure(t, "workflow.define", "idempotency-conflict", 4, defineConflict)
	if got := p21Path(t, p23Envelope(t, defineConflict.Stderr), "error", "details", "differing_fields"); !reflect.DeepEqual(got, []any{p24WorkflowFields["workflow.define"]}) {
		t.Fatalf("workflow.define differing_fields = %#v", got)
	}
	if after := p21JournalEventCount(t, workspace); after != defineEvents {
		t.Fatalf("workflow.define replay/conflict appended Events: %d -> %d", defineEvents, after)
	}
	assign := []string{"--project", workspace, "--json", "workflow", "assign", "--name", "task-p24", "--kind", "subtask", "--actor", "agent:developer", "--idempotency-key", "p24-workflow-assign"}
	assigned := p24WorkflowCall(t, workspace, binary, assign)
	if assigned.ExitCode != 0 || assigned.Stderr != "" {
		t.Fatalf("workflow.assign first = %d %q", assigned.ExitCode, assigned.Stderr)
	}
	p24AssertWorkflowFields(t, workspace, "p24-workflow-assign", []string{"--name", "--kind", "--actor"})
	before := p21JournalEventCount(t, workspace)
	permuted := []string{"--project", workspace, "--json", "workflow", "assign", "--actor", "agent:developer", "--kind", "subtask", "--idempotency-key", "p24-workflow-assign", "--name", "task-p24"}
	replay := p24WorkflowCall(t, workspace, binary, permuted)
	if replay.ExitCode != 0 || replay.Stderr != "" {
		t.Fatalf("workflow.assign replay = %d %q", replay.ExitCode, replay.Stderr)
	}
	if byteWant := strings.Replace(assigned.Stdout, `"replayed":false`, `"replayed":true`, 1); byteWant == assigned.Stdout || replay.Stdout != byteWant {
		t.Fatal("workflow.assign replay bytes changed beyond replayed")
	}
	want := p23Envelope(t, assigned.Stdout)
	want["result"].(map[string]any)["replayed"] = true
	if got := p23Envelope(t, replay.Stdout); !reflect.DeepEqual(got, want) {
		t.Fatalf("workflow.assign replay = %#v, want %#v", got, want)
	}
	assignChanged := append([]string(nil), assign...)
	for i, arg := range assignChanged {
		if arg == "--name" {
			assignChanged[i+1] = "task-other"
		}
	}
	oneField := p24WorkflowCall(t, workspace, binary, assignChanged)
	assertP22ProcessFailure(t, "workflow.assign", "idempotency-conflict", 4, oneField)
	if got := p21Path(t, p23Envelope(t, oneField.Stderr), "error", "details", "differing_fields"); !reflect.DeepEqual(got, []any{p24WorkflowFields["workflow.assign"]}) {
		t.Fatalf("workflow.assign one-field differing_fields = %#v", got)
	}
	actorChanged := append([]string(nil), assign...)
	for i, arg := range actorChanged {
		if arg == "--actor" {
			actorChanged[i+1] = "agent:other"
		}
	}
	actorConflict := p24WorkflowCall(t, workspace, binary, actorChanged)
	assertP22ProcessFailure(t, "workflow.assign", "idempotency-conflict", 4, actorConflict)
	if got := p21Path(t, p23Envelope(t, actorConflict.Stderr), "error", "details", "differing_fields"); !reflect.DeepEqual(got, []any{"--actor"}) {
		t.Fatalf("workflow.assign Actor differing_fields = %#v", got)
	}
	changed := []string{"--project", workspace, "--json", "workflow", "assign", "--name", "task-other", "--kind", "subtask", "--actor", "agent:reviewer", "--idempotency-key", "p24-workflow-assign"}
	conflict := p24WorkflowCall(t, workspace, binary, changed)
	assertP22ProcessFailure(t, "workflow.assign", "idempotency-conflict", 4, conflict)
	if got := p21Path(t, p23Envelope(t, conflict.Stderr), "error", "details", "differing_fields"); !reflect.DeepEqual(got, []any{"--actor", "--name"}) {
		t.Fatalf("workflow.assign differing_fields = %#v", got)
	}
	if after := p21JournalEventCount(t, workspace); after != before {
		t.Fatalf("workflow.assign replay/conflict appended Events: %d -> %d", before, after)
	}
	p24AssertRebuildMapParity(t, workspace, binary)
}

func p24WorkflowCall(t *testing.T, workspace, binary string, args []string) p22ConformanceResponse {
	t.Helper()
	stdout, stderr, err := closureRunCommand(workspace, binary, args...)
	response := p22ConformanceResponse{Stdout: stdout, Stderr: stderr}
	if err != nil {
		response.ExitCode = processExitCode(t, err)
	}
	return response
}

func p24AssertWorkflowFields(t *testing.T, workspace, key string, want []string) {
	t.Helper()
	conn := p22FaultDB(t, workspace)
	defer func() { p22FaultOK(t, conn.Close()) }()
	var encoded string
	if err := conn.QueryRow(`SELECT request_field_digests_json FROM idempotency_results WHERE idempotency_key = ?`, key).Scan(&encoded); err != nil {
		t.Fatal(err)
	}
	var fields map[string]string
	if err := json.Unmarshal([]byte(encoded), &fields); err != nil {
		t.Fatal(err)
	}
	names := make([]string, 0, len(fields))
	for name := range fields {
		names = append(names, name)
	}
	if len(names) != len(want) {
		t.Fatalf("%s field map = %v, want %v", key, fields, want)
	}
	for _, name := range want {
		if fields[name] == "" {
			t.Fatalf("%s missing %s: %v", key, name, fields)
		}
	}
}
