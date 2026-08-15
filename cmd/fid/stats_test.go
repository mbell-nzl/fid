package main

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"

	"github.com/mbell-nzl/fid/internal/issues"
	"github.com/mbell-nzl/fid/internal/runs"
	internalstats "github.com/mbell-nzl/fid/internal/stats"
)

func TestStatsTokensDefaultsToJSONByRole(t *testing.T) {
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

	want := expectedCommandStructure(t, []internalstats.TokensRow{
		{
			Role:           stringPointer("developer"),
			TotalTokensIn:  120,
			TotalTokensOut: 30,
			RunCount:       1,
		},
		{
			Role:           stringPointer("qa_reviewer"),
			TotalTokensIn:  40,
			TotalTokensOut: 10,
			RunCount:       1,
		},
	})
	if got := decodeCommandJSONStructure(t, stdout); !reflect.DeepEqual(got, want) {
		t.Fatalf("stats tokens json = %#v, want %#v", got, want)
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

	jsonStdout, jsonStderr := runFIDCommand(t, workspace, "stats", "tokens", "--by", "epic")
	if jsonStderr != "" {
		t.Fatalf("json stderr = %q, want empty", jsonStderr)
	}

	wantRows := []internalstats.TokensRow{
		{
			TotalTokensIn:  30,
			TotalTokensOut: 5,
			RunCount:       1,
		},
		{
			EpicID:         stringPointer(epicID),
			TotalTokensIn:  220,
			TotalTokensOut: 37,
			RunCount:       3,
		},
	}
	want := expectedCommandStructure(t, wantRows)
	if got := decodeCommandJSONStructure(t, jsonStdout); !reflect.DeepEqual(got, want) {
		t.Fatalf("stats tokens by epic json = %#v, want %#v", got, want)
	}

	filteredStdout, filteredStderr := runFIDCommand(t, workspace, "stats", "tokens", "--by", "role", "--epic", epicID)
	if filteredStderr != "" {
		t.Fatalf("filtered stderr = %q, want empty", filteredStderr)
	}
	wantFiltered := expectedCommandStructure(t, []internalstats.TokensRow{
		{
			Role:           stringPointer("developer"),
			TotalTokensIn:  100,
			TotalTokensOut: 10,
			RunCount:       1,
		},
		{
			Role:           stringPointer("qa_reviewer"),
			TotalTokensIn:  50,
			TotalTokensOut: 20,
			RunCount:       1,
		},
		{
			Role:           stringPointer("tl"),
			TotalTokensIn:  70,
			TotalTokensOut: 7,
			RunCount:       1,
		},
	})
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

	jsonStdout, jsonStderr := runFIDCommand(t, workspace, "stats", "cost", "--by", "epic")
	if jsonStderr != "" {
		t.Fatalf("json stderr = %q, want empty", jsonStderr)
	}

	wantRows := []internalstats.CostRow{
		{
			TotalCostUSD:   derefFloat(knownStandaloneRun.CostUSD),
			TotalTokensIn:  40,
			TotalTokensOut: 10,
			RunCount:       1,
		},
		{
			EpicID:         stringPointer(epicID),
			TotalCostUSD:   derefFloat(knownEpicRun.CostUSD),
			TotalTokensIn:  200,
			TotalTokensOut: 50,
			RunCount:       2,
		},
	}
	want := expectedCommandStructure(t, wantRows)
	if got := decodeCommandJSONStructure(t, jsonStdout); !reflect.DeepEqual(got, want) {
		t.Fatalf("stats cost by epic json = %#v, want %#v", got, want)
	}

	filteredStdout, filteredStderr := runFIDCommand(t, workspace, "stats", "cost", "--by", "role", "--epic", epicID)
	if filteredStderr != "" {
		t.Fatalf("filtered stderr = %q, want empty", filteredStderr)
	}
	wantFiltered := expectedCommandStructure(t, []internalstats.CostRow{
		{
			Role:           stringPointer("developer"),
			TotalCostUSD:   0,
			TotalTokensIn:  80,
			TotalTokensOut: 20,
			RunCount:       1,
		},
		{
			Role:           stringPointer("tl"),
			TotalCostUSD:   derefFloat(knownEpicRun.CostUSD),
			TotalTokensIn:  120,
			TotalTokensOut: 30,
			RunCount:       1,
		},
	})
	if got := decodeCommandJSONStructure(t, filteredStdout); !reflect.DeepEqual(got, wantFiltered) {
		t.Fatalf("stats cost filtered json = %#v, want %#v", got, wantFiltered)
	}
}

func createStatsCommandIssue(t *testing.T, workspace, kind, title string, parentID *string) string {
	t.Helper()

	options := issues.CreateOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: issueMigrationsDir(),
		Kind:          kind,
		Title:         title,
		ActorRole:     "developer",
		ParentID:      parentID,
	}
	if kind == "epic" {
		options.ActorRole = "tl"
		options.Scheduling = "prioritised"
	}

	created, err := issues.Create(options)
	if err != nil {
		t.Fatalf("Create(%s) error = %v", kind, err)
	}
	return created.ID
}

func recordStatsCommandRun(t *testing.T, workspace, issueID, runRole, model string, tokensIn, tokensOut int) runs.RecordResult {
	t.Helper()

	result, err := runs.Record(runs.RecordOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: issueMigrationsDir(),
		IssueID:       issueID,
		ActorRole:     runRole,
		RunRole:       runRole,
		Model:         model,
		TokensIn:      intPointer(tokensIn),
		TokensOut:     intPointer(tokensOut),
	})
	if err != nil {
		t.Fatalf("Record() error = %v", err)
	}
	return result
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
