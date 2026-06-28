package stats

import (
	"path/filepath"
	"reflect"
	"runtime"
	"testing"

	"github.com/mbell-nzl/fid/internal/bootstrap"
	"github.com/mbell-nzl/fid/internal/issues"
	"github.com/mbell-nzl/fid/internal/runs"
)

func TestTokensGroupsByRole(t *testing.T) {
	t.Parallel()

	workspace := newStatsWorkspace(t)
	epicID := createStatsIssue(t, workspace, "epic", "Telemetry epic", nil)
	taskID := createStatsIssue(t, workspace, "task", "Child task", stringPointer(epicID))
	otherTaskID := createStatsIssue(t, workspace, "task", "Standalone task", nil)

	recordStatsRun(t, workspace, taskID, "developer", "claude-sonnet-4-6", 120, 30)
	recordStatsRun(t, workspace, taskID, "developer", "gpt-5", 80, 20)
	recordStatsRun(t, workspace, otherTaskID, "qa_reviewer", "claude-sonnet-4-6", 50, 10)

	got, err := Tokens(TokensOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testStatsMigrationsDir(t),
		By:            GroupByRole,
	})
	if err != nil {
		t.Fatalf("Tokens() error = %v", err)
	}

	want := []TokensRow{
		{
			Role:           stringPointer("developer"),
			TotalTokensIn:  200,
			TotalTokensOut: 50,
			RunCount:       2,
		},
		{
			Role:           stringPointer("qa_reviewer"),
			TotalTokensIn:  50,
			TotalTokensOut: 10,
			RunCount:       1,
		},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Tokens() = %#v, want %#v", got, want)
	}
}

func TestTokensGroupsByModelAndEpic(t *testing.T) {
	t.Parallel()

	workspace := newStatsWorkspace(t)
	epicID := createStatsIssue(t, workspace, "epic", "Telemetry epic", nil)
	taskID := createStatsIssue(t, workspace, "task", "Epic task", stringPointer(epicID))
	subtaskID := createStatsIssue(t, workspace, "subtask", "Epic subtask", stringPointer(taskID))
	standaloneTaskID := createStatsIssue(t, workspace, "task", "Standalone task", nil)

	recordStatsRun(t, workspace, taskID, "developer", "claude-sonnet-4-6", 100, 10)
	recordStatsRun(t, workspace, subtaskID, "qa_reviewer", "claude-sonnet-4-6", 50, 20)
	recordStatsRun(t, workspace, epicID, "tl", "gemini-2.5-pro", 70, 7)
	recordStatsRun(t, workspace, standaloneTaskID, "developer", "gpt-5", 30, 5)

	modelRows, err := Tokens(TokensOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testStatsMigrationsDir(t),
		By:            GroupByModel,
	})
	if err != nil {
		t.Fatalf("Tokens(by model) error = %v", err)
	}

	wantModelRows := []TokensRow{
		{
			Model:          stringPointer("claude-sonnet-4-6"),
			TotalTokensIn:  150,
			TotalTokensOut: 30,
			RunCount:       2,
		},
		{
			Model:          stringPointer("gemini-2.5-pro"),
			TotalTokensIn:  70,
			TotalTokensOut: 7,
			RunCount:       1,
		},
		{
			Model:          stringPointer("gpt-5"),
			TotalTokensIn:  30,
			TotalTokensOut: 5,
			RunCount:       1,
		},
	}
	if !reflect.DeepEqual(modelRows, wantModelRows) {
		t.Fatalf("Tokens(by model) = %#v, want %#v", modelRows, wantModelRows)
	}

	epicRows, err := Tokens(TokensOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testStatsMigrationsDir(t),
		By:            GroupByEpic,
	})
	if err != nil {
		t.Fatalf("Tokens(by epic) error = %v", err)
	}

	wantEpicRows := []TokensRow{
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
	if !reflect.DeepEqual(epicRows, wantEpicRows) {
		t.Fatalf("Tokens(by epic) = %#v, want %#v", epicRows, wantEpicRows)
	}

	filteredRows, err := Tokens(TokensOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testStatsMigrationsDir(t),
		By:            GroupByRole,
		EpicID:        epicID,
	})
	if err != nil {
		t.Fatalf("Tokens(filtered by epic) error = %v", err)
	}

	wantFilteredRows := []TokensRow{
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
	}
	if !reflect.DeepEqual(filteredRows, wantFilteredRows) {
		t.Fatalf("Tokens(filtered by epic) = %#v, want %#v", filteredRows, wantFilteredRows)
	}
}

func TestCostGroupsByEpicAndExcludesNullCostsFromSum(t *testing.T) {
	t.Parallel()

	workspace := newStatsWorkspace(t)
	epicID := createStatsIssue(t, workspace, "epic", "Telemetry epic", nil)
	taskID := createStatsIssue(t, workspace, "task", "Epic task", stringPointer(epicID))
	standaloneTaskID := createStatsIssue(t, workspace, "task", "Standalone task", nil)

	knownEpicRun := recordStatsRun(t, workspace, epicID, "tl", "claude-sonnet-4-6", 120, 30)
	recordStatsRun(t, workspace, taskID, "developer", "mystery-model", 80, 20)
	knownStandaloneRun := recordStatsRun(t, workspace, standaloneTaskID, "developer", "gpt-5", 40, 10)

	epicRows, err := Cost(CostOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testStatsMigrationsDir(t),
		By:            GroupByEpic,
	})
	if err != nil {
		t.Fatalf("Cost(by epic) error = %v", err)
	}

	wantEpicRows := []CostRow{
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
	if !reflect.DeepEqual(epicRows, wantEpicRows) {
		t.Fatalf("Cost(by epic) = %#v, want %#v", epicRows, wantEpicRows)
	}

	filteredRows, err := Cost(CostOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testStatsMigrationsDir(t),
		By:            GroupByRole,
		EpicID:        epicID,
	})
	if err != nil {
		t.Fatalf("Cost(filtered by epic) error = %v", err)
	}

	wantFilteredRows := []CostRow{
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
	}
	if !reflect.DeepEqual(filteredRows, wantFilteredRows) {
		t.Fatalf("Cost(filtered by epic) = %#v, want %#v", filteredRows, wantFilteredRows)
	}
}

func newStatsWorkspace(t *testing.T) string {
	t.Helper()

	workspace := t.TempDir()
	if _, err := bootstrap.InitializeWorkspace(bootstrap.InitOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testStatsMigrationsDir(t),
	}); err != nil {
		t.Fatalf("InitializeWorkspace() error = %v", err)
	}

	return workspace
}

func createStatsIssue(t *testing.T, workspace, kind, title string, parentID *string) string {
	t.Helper()

	options := issues.CreateOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testStatsMigrationsDir(t),
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

func recordStatsRun(t *testing.T, workspace, issueID, runRole, model string, tokensIn, tokensOut int) runs.RecordResult {
	t.Helper()

	result, err := runs.Record(runs.RecordOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testStatsMigrationsDir(t),
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

func testStatsMigrationsDir(t *testing.T) string {
	t.Helper()

	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller(0) failed")
	}

	return filepath.Clean(filepath.Join(filepath.Dir(currentFile), "..", "..", "migrations"))
}

func intPointer(value int) *int {
	return &value
}

func stringPointer(value string) *string {
	return &value
}

func derefFloat(value *float64) float64 {
	if value == nil {
		return 0
	}
	return *value
}
