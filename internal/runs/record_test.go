package runs

import (
	"database/sql"
	"math"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/mbell-nzl/fid/internal/bootstrap"
	fiddb "github.com/mbell-nzl/fid/internal/db"
	"github.com/mbell-nzl/fid/internal/issues"
)

func TestRecordKnownModelPersistsAgentRunAndComputedCost(t *testing.T) {
	t.Parallel()

	workspace := newRunWorkspace(t)

	created, err := issues.Create(issues.CreateOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testRunMigrationsDir(t),
		Kind:          "task",
		Title:         "Record telemetry run",
		ActorRole:     "developer",
	})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	conn := openRunDB(t, filepath.Join(workspace, ".fid", "fid.sqlite"))
	defer closeRunDB(t, conn)

	var activityBefore int
	if err := conn.QueryRow(`SELECT COUNT(*) FROM activity_log`).Scan(&activityBefore); err != nil {
		t.Fatalf("count activity_log before: %v", err)
	}

	journalPath := journalPathForRunWorkspace(t, workspace)
	journalLinesBefore := countRunJournalLines(journalPath)

	result, err := Record(RecordOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testRunMigrationsDir(t),
		IssueID:       created.ID,
		ActorRole:     "developer",
		RunRole:       "developer",
		Model:         "claude-sonnet-4-6",
		TokensIn:      intPointer(1200),
		TokensOut:     intPointer(3400),
	})
	if err != nil {
		t.Fatalf("Record() error = %v", err)
	}

	if !strings.HasPrefix(result.ID, "RUN-") {
		t.Fatalf("Record().ID = %q, want RUN- prefix", result.ID)
	}
	if result.Warning != nil {
		t.Fatalf("Record().Warning = %#v, want nil", result.Warning)
	}
	if result.CostUSD == nil {
		t.Fatal("Record().CostUSD = nil, want computed cost")
	}

	var row struct {
		IssueID         sql.NullString
		CompletionID    sql.NullString
		Role            string
		Model           string
		TokensIn        sql.NullInt64
		TokensOut       sql.NullInt64
		PriceInPerMTok  sql.NullFloat64
		PriceOutPerMTok sql.NullFloat64
		CostUSD         sql.NullFloat64
		TokenBreakdown  sql.NullString
		ToolCallCount   sql.NullInt64
		ProviderRunID   sql.NullString
		ProviderCostUSD sql.NullFloat64
		StartedAt       sql.NullString
		EndedAt         sql.NullString
	}
	if err := conn.QueryRow(`
		SELECT issue_id, completion_id, role, model, tokens_in, tokens_out,
		       price_in_per_mtok, price_out_per_mtok, cost_usd,
		       token_breakdown_json, tool_call_count, provider_run_id, provider_cost_usd,
		       started_at, ended_at
		FROM agent_runs
		WHERE id = ?
	`, result.ID).Scan(
		&row.IssueID,
		&row.CompletionID,
		&row.Role,
		&row.Model,
		&row.TokensIn,
		&row.TokensOut,
		&row.PriceInPerMTok,
		&row.PriceOutPerMTok,
		&row.CostUSD,
		&row.TokenBreakdown,
		&row.ToolCallCount,
		&row.ProviderRunID,
		&row.ProviderCostUSD,
		&row.StartedAt,
		&row.EndedAt,
	); err != nil {
		t.Fatalf("query agent_runs row: %v", err)
	}

	if !row.IssueID.Valid || row.IssueID.String != created.ID {
		t.Fatalf("agent_runs.issue_id = %#v, want %q", row.IssueID, created.ID)
	}
	if row.CompletionID.Valid {
		t.Fatalf("agent_runs.completion_id = %#v, want NULL", row.CompletionID)
	}
	if row.Role != "developer" {
		t.Fatalf("agent_runs.role = %q, want developer", row.Role)
	}
	if row.Model != "claude-sonnet-4-6" {
		t.Fatalf("agent_runs.model = %q, want claude-sonnet-4-6", row.Model)
	}
	if !row.TokensIn.Valid || row.TokensIn.Int64 != 1200 {
		t.Fatalf("agent_runs.tokens_in = %#v, want 1200", row.TokensIn)
	}
	if !row.TokensOut.Valid || row.TokensOut.Int64 != 3400 {
		t.Fatalf("agent_runs.tokens_out = %#v, want 3400", row.TokensOut)
	}
	if !row.PriceInPerMTok.Valid || !row.PriceOutPerMTok.Valid {
		t.Fatalf("snapshot prices = (%#v, %#v), want both populated", row.PriceInPerMTok, row.PriceOutPerMTok)
	}
	if !row.CostUSD.Valid {
		t.Fatal("agent_runs.cost_usd = NULL, want computed value")
	}
	wantCost := (1200*row.PriceInPerMTok.Float64 + 3400*row.PriceOutPerMTok.Float64) / 1_000_000
	if math.Abs(row.CostUSD.Float64-wantCost) > 1e-12 {
		t.Fatalf("agent_runs.cost_usd = %v, want %v", row.CostUSD.Float64, wantCost)
	}
	if math.Abs(*result.CostUSD-row.CostUSD.Float64) > 1e-12 {
		t.Fatalf("Record().CostUSD = %v, want %v", *result.CostUSD, row.CostUSD.Float64)
	}
	if row.TokenBreakdown.Valid {
		t.Fatalf("agent_runs.token_breakdown_json = %#v, want NULL", row.TokenBreakdown)
	}
	if row.ToolCallCount.Valid {
		t.Fatalf("agent_runs.tool_call_count = %#v, want NULL", row.ToolCallCount)
	}
	if row.ProviderRunID.Valid {
		t.Fatalf("agent_runs.provider_run_id = %#v, want NULL", row.ProviderRunID)
	}
	if row.ProviderCostUSD.Valid {
		t.Fatalf("agent_runs.provider_cost_usd = %#v, want NULL", row.ProviderCostUSD)
	}
	if !row.StartedAt.Valid || row.StartedAt.String == "" {
		t.Fatalf("agent_runs.started_at = %#v, want timestamp", row.StartedAt)
	}
	if !row.EndedAt.Valid || row.EndedAt.String == "" {
		t.Fatalf("agent_runs.ended_at = %#v, want timestamp", row.EndedAt)
	}

	if got := countRunJournalLines(journalPath); got != journalLinesBefore+1 {
		t.Fatalf("countRunJournalLines(%q) = %d, want %d", journalPath, got, journalLinesBefore+1)
	}

	var activityAfter int
	if err := conn.QueryRow(`SELECT COUNT(*) FROM activity_log`).Scan(&activityAfter); err != nil {
		t.Fatalf("count activity_log after: %v", err)
	}
	if activityAfter != activityBefore {
		t.Fatalf("activity_log rows after Record() = %d, want unchanged %d", activityAfter, activityBefore)
	}
}

func TestRecordUnknownModelStoresNullCostAndWarning(t *testing.T) {
	t.Parallel()

	workspace := newRunWorkspace(t)

	created, err := issues.Create(issues.CreateOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testRunMigrationsDir(t),
		Kind:          "task",
		Title:         "Record unknown-model telemetry run",
		ActorRole:     "developer",
	})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	conn := openRunDB(t, filepath.Join(workspace, ".fid", "fid.sqlite"))
	defer closeRunDB(t, conn)

	result, err := Record(RecordOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testRunMigrationsDir(t),
		IssueID:       created.ID,
		ActorRole:     "developer",
		RunRole:       "developer",
		Model:         "mystery-model",
		TokensIn:      intPointer(500),
		TokensOut:     intPointer(700),
	})
	if err != nil {
		t.Fatalf("Record() error = %v", err)
	}

	if result.Warning == nil {
		t.Fatal("Record().Warning = nil, want unknown-model warning")
	}
	if !strings.Contains(*result.Warning, `mystery-model`) {
		t.Fatalf("Record().Warning = %q, want model name", *result.Warning)
	}
	if result.CostUSD != nil {
		t.Fatalf("Record().CostUSD = %#v, want nil", result.CostUSD)
	}

	var priceIn sql.NullFloat64
	var priceOut sql.NullFloat64
	var cost sql.NullFloat64
	if err := conn.QueryRow(`
		SELECT price_in_per_mtok, price_out_per_mtok, cost_usd
		FROM agent_runs
		WHERE id = ?
	`, result.ID).Scan(&priceIn, &priceOut, &cost); err != nil {
		t.Fatalf("query unknown-model row: %v", err)
	}
	if priceIn.Valid {
		t.Fatalf("agent_runs.price_in_per_mtok = %#v, want NULL", priceIn)
	}
	if priceOut.Valid {
		t.Fatalf("agent_runs.price_out_per_mtok = %#v, want NULL", priceOut)
	}
	if cost.Valid {
		t.Fatalf("agent_runs.cost_usd = %#v, want NULL", cost)
	}
}

func TestRecordProviderFieldsPersistWithoutChangingNormalizedSummary(t *testing.T) {
	t.Parallel()

	workspace := newRunWorkspace(t)

	created, err := issues.Create(issues.CreateOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testRunMigrationsDir(t),
		Kind:          "task",
		Title:         "Record provider fields",
		ActorRole:     "developer",
	})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	conn := openRunDB(t, filepath.Join(workspace, ".fid", "fid.sqlite"))
	defer closeRunDB(t, conn)

	withProvider, err := Record(RecordOptions{
		WorkspaceRoot:      workspace,
		MigrationsDir:      testRunMigrationsDir(t),
		IssueID:            created.ID,
		ActorRole:          "developer",
		RunRole:            "developer",
		Model:              "claude-sonnet-4-6",
		TokensIn:           intPointer(1200),
		TokensOut:          intPointer(3400),
		TokenBreakdownJSON: stringPointer(`{"input":900,"cache_read":300}`),
		ToolCallCount:      intPointer(4),
		ProviderRunID:      stringPointer("run_123"),
		ProviderCostUSD:    floatPointer(0.055),
	})
	if err != nil {
		t.Fatalf("Record(with provider fields) error = %v", err)
	}

	withoutProvider, err := Record(RecordOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testRunMigrationsDir(t),
		IssueID:       created.ID,
		ActorRole:     "developer",
		RunRole:       "developer",
		Model:         "claude-sonnet-4-6",
		TokensIn:      intPointer(1200),
		TokensOut:     intPointer(3400),
	})
	if err != nil {
		t.Fatalf("Record(without provider fields) error = %v", err)
	}

	type runRow struct {
		TokensIn        sql.NullInt64
		TokensOut       sql.NullInt64
		PriceInPerMTok  sql.NullFloat64
		PriceOutPerMTok sql.NullFloat64
		CostUSD         sql.NullFloat64
		TokenBreakdown  sql.NullString
		ToolCallCount   sql.NullInt64
		ProviderRunID   sql.NullString
		ProviderCostUSD sql.NullFloat64
	}

	loadRow := func(runID string) runRow {
		t.Helper()

		var row runRow
		if err := conn.QueryRow(`
			SELECT tokens_in, tokens_out, price_in_per_mtok, price_out_per_mtok, cost_usd,
			       token_breakdown_json, tool_call_count, provider_run_id, provider_cost_usd
			FROM agent_runs
			WHERE id = ?
		`, runID).Scan(
			&row.TokensIn,
			&row.TokensOut,
			&row.PriceInPerMTok,
			&row.PriceOutPerMTok,
			&row.CostUSD,
			&row.TokenBreakdown,
			&row.ToolCallCount,
			&row.ProviderRunID,
			&row.ProviderCostUSD,
		); err != nil {
			t.Fatalf("query agent_runs row %q: %v", runID, err)
		}
		return row
	}

	withRow := loadRow(withProvider.ID)
	withoutRow := loadRow(withoutProvider.ID)

	if !withRow.TokenBreakdown.Valid || withRow.TokenBreakdown.String != `{"input":900,"cache_read":300}` {
		t.Fatalf("token_breakdown_json = %#v, want provider JSON", withRow.TokenBreakdown)
	}
	if !withRow.ToolCallCount.Valid || withRow.ToolCallCount.Int64 != 4 {
		t.Fatalf("tool_call_count = %#v, want 4", withRow.ToolCallCount)
	}
	if !withRow.ProviderRunID.Valid || withRow.ProviderRunID.String != "run_123" {
		t.Fatalf("provider_run_id = %#v, want run_123", withRow.ProviderRunID)
	}
	if !withRow.ProviderCostUSD.Valid || math.Abs(withRow.ProviderCostUSD.Float64-0.055) > 1e-12 {
		t.Fatalf("provider_cost_usd = %#v, want 0.055", withRow.ProviderCostUSD)
	}
	if withoutRow.TokenBreakdown.Valid || withoutRow.ToolCallCount.Valid || withoutRow.ProviderRunID.Valid || withoutRow.ProviderCostUSD.Valid {
		t.Fatalf("provider fields without options = %#v, want NULLs", withoutRow)
	}

	if withRow.TokensIn != withoutRow.TokensIn || withRow.TokensOut != withoutRow.TokensOut {
		t.Fatalf("token summary differs: with=%#v without=%#v", withRow, withoutRow)
	}
	if withRow.PriceInPerMTok != withoutRow.PriceInPerMTok || withRow.PriceOutPerMTok != withoutRow.PriceOutPerMTok {
		t.Fatalf("snapshot prices differ: with=%#v without=%#v", withRow, withoutRow)
	}
	if math.Abs(withRow.CostUSD.Float64-withoutRow.CostUSD.Float64) > 1e-12 {
		t.Fatalf("cost_usd differs: with=%v without=%v", withRow.CostUSD.Float64, withoutRow.CostUSD.Float64)
	}
}

func testRunMigrationsDir(t *testing.T) string {
	t.Helper()

	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller(0) failed")
	}

	return filepath.Clean(filepath.Join(filepath.Dir(currentFile), "..", "..", "migrations"))
}

func newRunWorkspace(t *testing.T) string {
	t.Helper()

	workspace := t.TempDir()
	if _, err := bootstrap.InitializeWorkspace(bootstrap.InitOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testRunMigrationsDir(t),
	}); err != nil {
		t.Fatalf("InitializeWorkspace() error = %v", err)
	}

	return workspace
}

func openRunDB(t *testing.T, dbPath string) *sql.DB {
	t.Helper()

	conn, err := fiddb.Open(dbPath)
	if err != nil {
		t.Fatalf("Open(%q) error = %v", dbPath, err)
	}

	return conn
}

func closeRunDB(t *testing.T, conn *sql.DB) {
	t.Helper()

	if err := conn.Close(); err != nil {
		t.Fatalf("Close() error = %v", err)
	}
}

func journalPathForRunWorkspace(t *testing.T, workspace string) string {
	t.Helper()

	return filepath.Join(workspace, ".fid", "journal", firstRunJournalFile(t, workspace))
}

func firstRunJournalFile(t *testing.T, workspace string) string {
	t.Helper()

	entries, err := os.ReadDir(filepath.Join(workspace, ".fid", "journal"))
	if err != nil {
		t.Fatalf("ReadDir(journal) error = %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("journal file count = %d, want 1", len(entries))
	}
	if entries[0].IsDir() {
		t.Fatalf("journal entry %q is a directory, want file", entries[0].Name())
	}
	return entries[0].Name()
}

func countRunJournalLines(path string) int {
	payload, err := os.ReadFile(path)
	if err != nil {
		return 0
	}
	trimmed := strings.TrimSpace(string(payload))
	if trimmed == "" {
		return 0
	}
	return len(strings.Split(trimmed, "\n"))
}

func intPointer(value int) *int {
	return &value
}
