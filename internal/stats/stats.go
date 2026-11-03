package stats

import (
	"database/sql"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/mbell-nzl/fid/internal/config"
	fiddb "github.com/mbell-nzl/fid/internal/db"
	"github.com/mbell-nzl/fid/internal/journal"
)

type GroupBy string

const (
	GroupByActor GroupBy = "actor"
	GroupByRole  GroupBy = GroupByActor
	GroupByModel GroupBy = "model"
	GroupByEpic  GroupBy = "epic"
)

type TokensOptions struct {
	WorkspaceRoot string
	MigrationsDir string
	By            GroupBy
	EpicID        string
}

type TokensRow struct {
	Role           *string `json:"role"`
	Model          *string `json:"model"`
	EpicID         *string `json:"epic_id"`
	TotalTokensIn  int     `json:"total_tokens_in"`
	TotalTokensOut int     `json:"total_tokens_out"`
	RunCount       int     `json:"run_count"`
}

type CostOptions struct {
	WorkspaceRoot string
	MigrationsDir string
	By            GroupBy
	EpicID        string
}

type CostRow struct {
	Role           *string `json:"role"`
	Model          *string `json:"model"`
	EpicID         *string `json:"epic_id"`
	TotalCostUSD   float64 `json:"total_cost_usd"`
	TotalTokensIn  int     `json:"total_tokens_in"`
	TotalTokensOut int     `json:"total_tokens_out"`
	RunCount       int     `json:"run_count"`
}

func Tokens(options TokensOptions) ([]TokensRow, error) {
	workspaceRoot, migrationsDir, by, epicID, err := validateStatsOptions(options.WorkspaceRoot, options.MigrationsDir, options.By, options.EpicID)
	if err != nil {
		return nil, err
	}

	if _, err := journal.EnsureProjectionFresh(workspaceRoot, migrationsDir); err != nil {
		return nil, err
	}

	conn, err := fiddb.Open(filepath.Join(config.LedgerDir(workspaceRoot), config.DatabaseFileName))
	if err != nil {
		return nil, err
	}
	defer closeStatsConn(conn)

	query, args := tokenStatsQuery(by, epicID)
	rows, err := conn.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("stats tokens: %w", err)
	}
	defer closeStatsRows(rows)

	items := []TokensRow{}
	for rows.Next() {
		item, err := scanTokensRow(rows, by)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate token stats rows: %w", err)
	}

	return items, nil
}

func Cost(options CostOptions) ([]CostRow, error) {
	workspaceRoot, migrationsDir, by, epicID, err := validateStatsOptions(options.WorkspaceRoot, options.MigrationsDir, options.By, options.EpicID)
	if err != nil {
		return nil, err
	}

	if _, err := journal.EnsureProjectionFresh(workspaceRoot, migrationsDir); err != nil {
		return nil, err
	}

	conn, err := fiddb.Open(filepath.Join(config.LedgerDir(workspaceRoot), config.DatabaseFileName))
	if err != nil {
		return nil, err
	}
	defer closeStatsConn(conn)

	query, args := costStatsQuery(by, epicID)
	rows, err := conn.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("stats cost: %w", err)
	}
	defer closeStatsRows(rows)

	items := []CostRow{}
	for rows.Next() {
		item, err := scanCostRow(rows, by)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate cost stats rows: %w", err)
	}

	return items, nil
}

func validateStatsOptions(workspaceRoot, migrationsDir string, by GroupBy, epicID string) (string, string, GroupBy, string, error) {
	trimmedWorkspaceRoot := strings.TrimSpace(workspaceRoot)
	if trimmedWorkspaceRoot == "" {
		return "", "", "", "", fmt.Errorf("stats: workspace root is required")
	}
	trimmedMigrationsDir := strings.TrimSpace(migrationsDir)
	if trimmedMigrationsDir == "" {
		return "", "", "", "", fmt.Errorf("stats: migrations dir is required")
	}

	normalizedBy := GroupBy(strings.TrimSpace(string(by)))
	if normalizedBy == "" {
		normalizedBy = GroupByActor
	}
	switch normalizedBy {
	case GroupByActor, GroupByModel, GroupByEpic:
	default:
		return "", "", "", "", fmt.Errorf("stats: unsupported group %q", normalizedBy)
	}

	return trimmedWorkspaceRoot, trimmedMigrationsDir, normalizedBy, strings.TrimSpace(epicID), nil
}

func tokenStatsQuery(by GroupBy, epicID string) (string, []any) {
	query := statsRunsCTE() + `
SELECT ` + statsGroupColumn(by) + ` AS group_value,
       COALESCE(SUM(COALESCE(tokens_in, 0)), 0) AS total_tokens_in,
       COALESCE(SUM(COALESCE(tokens_out, 0)), 0) AS total_tokens_out,
       COUNT(*) AS run_count
FROM stats_runs
`
	args := []any{}
	if epicID != "" {
		query += "WHERE owning_epic_id = ?\n"
		args = append(args, epicID)
	}
	query += `GROUP BY group_value
ORDER BY group_value IS NOT NULL, group_value ASC`
	return query, args
}

func costStatsQuery(by GroupBy, epicID string) (string, []any) {
	query := statsRunsCTE() + `
SELECT ` + statsGroupColumn(by) + ` AS group_value,
       COALESCE(SUM(cost_usd), 0) AS total_cost_usd,
       COALESCE(SUM(COALESCE(tokens_in, 0)), 0) AS total_tokens_in,
       COALESCE(SUM(COALESCE(tokens_out, 0)), 0) AS total_tokens_out,
       COUNT(*) AS run_count
FROM stats_runs
`
	args := []any{}
	if epicID != "" {
		query += "WHERE owning_epic_id = ?\n"
		args = append(args, epicID)
	}
	query += `GROUP BY group_value
ORDER BY group_value IS NOT NULL, group_value ASC`
	return query, args
}

func statsRunsCTE() string {
	return `WITH stats_runs AS (
	SELECT
		agent_runs.role,
		agent_runs.model,
		agent_runs.tokens_in,
		agent_runs.tokens_out,
		agent_runs.cost_usd,
		CASE
			WHEN issue.kind = 'epic' THEN issue.id
			WHEN parent.kind = 'epic' THEN parent.id
			WHEN grandparent.kind = 'epic' THEN grandparent.id
			ELSE NULL
		END AS owning_epic_id
	FROM agent_runs
	LEFT JOIN issues issue ON issue.id = agent_runs.issue_id
	LEFT JOIN issues parent ON parent.id = issue.parent_id
	LEFT JOIN issues grandparent ON grandparent.id = parent.parent_id
)
`
}

func statsGroupColumn(by GroupBy) string {
	switch by {
	case GroupByModel:
		return "model"
	case GroupByEpic:
		return "owning_epic_id"
	default:
		return "role"
	}
}

func scanTokensRow(rows *sql.Rows, by GroupBy) (TokensRow, error) {
	var (
		groupValue     sql.NullString
		totalTokensIn  int
		totalTokensOut int
		runCount       int
	)
	if err := rows.Scan(&groupValue, &totalTokensIn, &totalTokensOut, &runCount); err != nil {
		return TokensRow{}, fmt.Errorf("scan token stats row: %w", err)
	}

	item := TokensRow{
		TotalTokensIn:  totalTokensIn,
		TotalTokensOut: totalTokensOut,
		RunCount:       runCount,
	}
	switch by {
	case GroupByModel:
		item.Model = nullStringPointer(groupValue)
	case GroupByEpic:
		item.EpicID = nullStringPointer(groupValue)
	default:
		item.Role = nullStringPointer(groupValue)
	}
	return item, nil
}

func scanCostRow(rows *sql.Rows, by GroupBy) (CostRow, error) {
	var (
		groupValue     sql.NullString
		totalCostUSD   float64
		totalTokensIn  int
		totalTokensOut int
		runCount       int
	)
	if err := rows.Scan(&groupValue, &totalCostUSD, &totalTokensIn, &totalTokensOut, &runCount); err != nil {
		return CostRow{}, fmt.Errorf("scan cost stats row: %w", err)
	}

	item := CostRow{
		TotalCostUSD:   totalCostUSD,
		TotalTokensIn:  totalTokensIn,
		TotalTokensOut: totalTokensOut,
		RunCount:       runCount,
	}
	switch by {
	case GroupByModel:
		item.Model = nullStringPointer(groupValue)
	case GroupByEpic:
		item.EpicID = nullStringPointer(groupValue)
	default:
		item.Role = nullStringPointer(groupValue)
	}
	return item, nil
}

func nullStringPointer(value sql.NullString) *string {
	if !value.Valid {
		return nil
	}
	cloned := value.String
	return &cloned
}

func closeStatsConn(conn *sql.DB) {
	if conn != nil {
		_ = conn.Close()
	}
}

func closeStatsRows(rows *sql.Rows) {
	if rows != nil {
		_ = rows.Close()
	}
}
