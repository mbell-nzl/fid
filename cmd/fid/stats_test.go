package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"path/filepath"
	"reflect"
	"strings"
	"sync/atomic"
	"testing"

	fiddb "github.com/mbell-nzl/fid/internal/db"
	"github.com/mbell-nzl/fid/internal/issues"
)

func TestStatsTokensDefaultsToJSONByActorAttribution(t *testing.T) {
	t.Parallel()
	workspace := t.TempDir()
	clonePristineWorkspaceInto(t, workspace)

	taskID := createStatsCommandIssue(t, workspace, "task", "CLI stats task", nil)
	recordStatsCommandRun(t, workspace, taskID, "developer", "claude-sonnet-4-6", 120, 30)
	recordStatsCommandRun(t, workspace, taskID, "qa_reviewer", "gpt-5", 40, 10)

	stdout, stderr := runFIDCommand(t, workspace, "stats", "tokens")
	if stderr != "" {
		t.Fatalf("stderr = %q, want empty", stderr)
	}

	want := expectedCommandStructure(t, map[string]any{"contract_version": "fid-cli.v1", "event": nil, "ok": true, "operation": "stats.tokens", "result": map[string]any{"items": []any{
		map[string]any{"actor_attribution": "developer", "epic_id": nil, "model": nil, "total_tokens_in": 120, "total_tokens_out": 30, "run_count": 1},
		map[string]any{"actor_attribution": "qa_reviewer", "epic_id": nil, "model": nil, "total_tokens_in": 40, "total_tokens_out": 10, "run_count": 1},
	}}})
	if got := decodeCommandJSONStructure(t, stdout); !reflect.DeepEqual(got, want) {
		t.Fatalf("stats tokens json = %#v, want %#v", got, want)
	}
}

func TestStatsTokensAcceptsOnlyClosedGroupByValues(t *testing.T) {
	workspace := clonePristineWorkspace(t)
	epicID := createStatsCommandIssue(t, workspace, "epic", "Closed stats grouping epic", nil)
	taskID := createStatsCommandIssue(t, workspace, "task", "Closed stats grouping task", stringPointer(epicID))
	recordStatsCommandRun(t, workspace, taskID, "developer", "gpt-5", 100, 25)
	journalPath := mustJournalPath(t, workspace)
	before := mustReadCommandFile(t, journalPath)

	accepted := []struct {
		group string
		item  map[string]any
	}{
		{
			group: "actor",
			item:  map[string]any{"actor_attribution": "developer", "epic_id": nil, "model": nil, "total_tokens_in": 100, "total_tokens_out": 25, "run_count": 1},
		},
		{
			group: "epic",
			item:  map[string]any{"actor_attribution": nil, "epic_id": epicID, "model": nil, "total_tokens_in": 100, "total_tokens_out": 25, "run_count": 1},
		},
		{
			group: "model",
			item:  map[string]any{"actor_attribution": nil, "epic_id": nil, "model": "gpt-5", "total_tokens_in": 100, "total_tokens_out": 25, "run_count": 1},
		},
	}
	for _, testCase := range accepted {
		t.Run("accepts "+testCase.group, func(t *testing.T) {
			result := executeCommandForTest(workspace, "--project", workspace, "--json", "stats", "tokens", "--group-by", testCase.group)
			if result.err != nil || result.stderr != "" {
				t.Fatalf("stats tokens --group-by %s = %#v, want successful JSON response", testCase.group, result)
			}
			want := expectedCommandStructure(t, map[string]any{
				"contract_version": "fid-cli.v1",
				"event":            nil,
				"ok":               true,
				"operation":        "stats.tokens",
				"result":           map[string]any{"items": []any{testCase.item}},
			})
			if got := decodeCommandJSONStructure(t, result.stdout); !reflect.DeepEqual(got, want) {
				t.Fatalf("stats tokens --group-by %s = %#v, want %#v", testCase.group, got, want)
			}
			if after := mustReadCommandFile(t, journalPath); !bytes.Equal(after, before) {
				t.Fatalf("journal changed after accepted %q grouping", testCase.group)
			}
		})
	}

	for _, testCase := range []struct {
		name string
		args []string
	}{
		{name: "retired role", args: []string{"--project", workspace, "--json", "stats", "tokens", "--group-by", "role"}},
		{name: "unknown group", args: []string{"--project", workspace, "--json", "stats", "tokens", "--group-by", "team"}},
		{name: "retired by flag", args: []string{"--project", workspace, "--json", "stats", "tokens", "--by", "developer"}},
	} {
		t.Run("rejects "+testCase.name, func(t *testing.T) {
			result := executeCommandForTest(workspace, testCase.args...)
			if result.err == nil || result.stdout != "" {
				t.Fatalf("%q result = %#v, want stderr-only rejection", testCase.args, result)
			}
			if after := mustReadCommandFile(t, journalPath); !bytes.Equal(after, before) {
				t.Fatalf("journal changed after rejected %q", testCase.name)
			}
		})
	}
}

func TestStatsTokensUsesDirectLegacyJSONAndEpicFilter(t *testing.T) {
	t.Parallel()
	workspace := t.TempDir()
	clonePristineWorkspaceInto(t, workspace)

	epicID := createStatsCommandIssue(t, workspace, "epic", "CLI telemetry epic", nil)
	taskID := createStatsCommandIssue(t, workspace, "task", "CLI epic task", stringPointer(epicID))
	subtaskID := createStatsCommandIssue(t, workspace, "subtask", "CLI epic subtask", stringPointer(taskID))
	standaloneTaskID := createStatsCommandIssue(t, workspace, "task", "CLI standalone task", nil)

	recordStatsCommandRun(t, workspace, taskID, "developer", "claude-sonnet-4-6", 100, 10)
	recordStatsCommandRun(t, workspace, subtaskID, "qa_reviewer", "claude-sonnet-4-6", 50, 20)
	recordStatsCommandRun(t, workspace, epicID, "tl", "gemini-2.5-pro", 70, 7)
	recordStatsCommandRun(t, workspace, standaloneTaskID, "developer", "gpt-5", 30, 5)

	jsonStdout, jsonStderr := runFIDCommand(t, workspace, "stats", "tokens", "--group-by", "epic")
	if jsonStderr != "" {
		t.Fatalf("json stderr = %q, want empty", jsonStderr)
	}

	want := expectedCommandStructure(t, map[string]any{"contract_version": "fid-cli.v1", "event": nil, "ok": true, "operation": "stats.tokens", "result": map[string]any{"items": []any{
		map[string]any{"actor_attribution": nil, "epic_id": nil, "model": nil, "total_tokens_in": 30, "total_tokens_out": 5, "run_count": 1},
		map[string]any{"actor_attribution": nil, "epic_id": epicID, "model": nil, "total_tokens_in": 220, "total_tokens_out": 37, "run_count": 3},
	}}})
	if got := decodeCommandJSONStructure(t, jsonStdout); !reflect.DeepEqual(got, want) {
		t.Fatalf("stats tokens by epic json = %#v, want %#v", got, want)
	}

	filteredStdout, filteredStderr := runFIDCommand(t, workspace, "stats", "tokens", "--group-by", "actor", "--epic", epicID)
	if filteredStderr != "" {
		t.Fatalf("filtered stderr = %q, want empty", filteredStderr)
	}
	wantFiltered := expectedCommandStructure(t, map[string]any{"contract_version": "fid-cli.v1", "event": nil, "ok": true, "operation": "stats.tokens", "result": map[string]any{"items": []any{
		map[string]any{"actor_attribution": "developer", "epic_id": nil, "model": nil, "total_tokens_in": 100, "total_tokens_out": 10, "run_count": 1},
		map[string]any{"actor_attribution": "qa_reviewer", "epic_id": nil, "model": nil, "total_tokens_in": 50, "total_tokens_out": 20, "run_count": 1},
		map[string]any{"actor_attribution": "tl", "epic_id": nil, "model": nil, "total_tokens_in": 70, "total_tokens_out": 7, "run_count": 1},
	}}})
	if got := decodeCommandJSONStructure(t, filteredStdout); !reflect.DeepEqual(got, wantFiltered) {
		t.Fatalf("stats tokens filtered json = %#v, want %#v", got, wantFiltered)
	}
}

func TestStatsCostUsesDirectLegacyJSONAndExcludesNullCostsFromSum(t *testing.T) {
	t.Parallel()
	workspace := t.TempDir()
	clonePristineWorkspaceInto(t, workspace)

	epicID := createStatsCommandIssue(t, workspace, "epic", "CLI cost epic", nil)
	taskID := createStatsCommandIssue(t, workspace, "task", "CLI cost task", stringPointer(epicID))
	standaloneTaskID := createStatsCommandIssue(t, workspace, "task", "CLI standalone task", nil)

	knownEpicRun := recordStatsCommandRun(t, workspace, epicID, "tl", "claude-sonnet-4-6", 120, 30)
	recordStatsCommandRun(t, workspace, taskID, "developer", "mystery-model", 80, 20)
	knownStandaloneRun := recordStatsCommandRun(t, workspace, standaloneTaskID, "developer", "gpt-5", 40, 10)

	jsonStdout, jsonStderr := runFIDCommand(t, workspace, "stats", "cost", "--group-by", "epic")
	if jsonStderr != "" {
		t.Fatalf("json stderr = %q, want empty", jsonStderr)
	}

	want := expectedCommandStructure(t, map[string]any{"contract_version": "fid-cli.v1", "event": nil, "ok": true, "operation": "stats.cost", "result": map[string]any{"items": []any{
		map[string]any{"actor_attribution": nil, "epic_id": nil, "model": nil, "total_cost_usd": derefFloat(knownStandaloneRun.CostUSD), "total_tokens_in": 40, "total_tokens_out": 10, "run_count": 1},
		map[string]any{"actor_attribution": nil, "epic_id": epicID, "model": nil, "total_cost_usd": derefFloat(knownEpicRun.CostUSD), "total_tokens_in": 200, "total_tokens_out": 50, "run_count": 2},
	}}})
	if got := decodeCommandJSONStructure(t, jsonStdout); !reflect.DeepEqual(got, want) {
		t.Fatalf("stats cost by epic json = %#v, want %#v", got, want)
	}

	filteredStdout, filteredStderr := runFIDCommand(t, workspace, "stats", "cost", "--group-by", "actor", "--epic", epicID)
	if filteredStderr != "" {
		t.Fatalf("filtered stderr = %q, want empty", filteredStderr)
	}
	wantFiltered := expectedCommandStructure(t, map[string]any{"contract_version": "fid-cli.v1", "event": nil, "ok": true, "operation": "stats.cost", "result": map[string]any{"items": []any{
		map[string]any{"actor_attribution": "developer", "epic_id": nil, "model": nil, "total_cost_usd": 0, "total_tokens_in": 80, "total_tokens_out": 20, "run_count": 1},
		map[string]any{"actor_attribution": "tl", "epic_id": nil, "model": nil, "total_cost_usd": derefFloat(knownEpicRun.CostUSD), "total_tokens_in": 120, "total_tokens_out": 30, "run_count": 1},
	}}})
	if got := decodeCommandJSONStructure(t, filteredStdout); !reflect.DeepEqual(got, wantFiltered) {
		t.Fatalf("stats cost filtered json = %#v, want %#v", got, wantFiltered)
	}
}

func createStatsCommandIssue(t *testing.T, workspace, kind, title string, parentID *string) string {
	t.Helper()

	options := issues.CreateOptions{
		WorkspaceRoot:  workspace,
		MigrationsDir:  issueMigrationsDir(),
		Kind:           kind,
		Title:          title,
		Actor:          "agent:developer",
		IdempotencyKey: fmt.Sprintf("stats-command-create-%d", atomic.AddUint64(&statsCommandFixtureSequence, 1)),
		ParentID:       parentID,
	}
	if kind == "epic" {
		options.Actor = "agent:tl"
		options.Scheduling = "prioritised"
	}

	created, err := issues.Create(options)
	if err != nil {
		t.Fatalf("Create(%s) error = %v", kind, err)
	}
	return created.ID
}

var statsCommandFixtureSequence uint64

type statsCommandRunFixture struct {
	CostUSD *float64
}

func recordStatsCommandRun(t *testing.T, workspace, issueID, runRole, model string, tokensIn, tokensOut int) statsCommandRunFixture {
	t.Helper()

	conn, err := fiddb.Open(filepath.Join(workspace, ".fid", "fid.sqlite"))
	if err != nil {
		t.Fatalf("open frozen telemetry fixture: %v", err)
	}
	defer closeCommandDB(t, conn)

	var count int
	if err := conn.QueryRow(`SELECT COUNT(*) FROM agent_runs`).Scan(&count); err != nil {
		t.Fatalf("count frozen telemetry rows: %v", err)
	}
	var cost any
	var costResult *float64
	if model != "mystery-model" {
		value := float64(tokensIn+tokensOut) / 1_000_000
		cost = value
		costResult = &value
	}
	if _, err := conn.Exec(`
		INSERT INTO agent_runs (id, issue_id, role, model, tokens_in, tokens_out, cost_usd)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`, fmt.Sprintf("STATS-%03d", count+1), issueID, runRole, model, tokensIn, tokensOut, cost); err != nil {
		t.Fatalf("insert frozen telemetry fixture: %v", err)
	}
	return statsCommandRunFixture{CostUSD: costResult}
}

func expectedCommandStructure(t *testing.T, value any) any {
	t.Helper()

	payload, err := json.Marshal(value)
	if err != nil {
		t.Fatalf("Marshal(expected structure) error = %v", err)
	}

	decoder := json.NewDecoder(strings.NewReader(string(payload)))
	decoder.UseNumber()

	var result any
	if err := decoder.Decode(&result); err != nil {
		t.Fatalf("Decode(expected structure) error = %v", err)
	}
	return result
}

func derefFloat(value *float64) float64 {
	if value == nil {
		return 0
	}
	return *value
}
