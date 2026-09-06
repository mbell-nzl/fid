package workflow

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/mbell-nzl/fid/internal/journal"
)

func GetRun(workspaceRoot, migrationsDir, id string) (Run, error) {
	if _, err := journal.EnsureProjectionFresh(workspaceRoot, migrationsDir); err != nil {
		return Run{}, err
	}
	conn, err := openWorkflowDB(workspaceRoot)
	if err != nil {
		return Run{}, err
	}
	defer closeWorkflowDB(conn)
	return loadRun(conn, strings.TrimSpace(id))
}

func ListRuns(options RunListOptions) ([]Run, error) {
	if strings.TrimSpace(options.WorkspaceRoot) == "" || strings.TrimSpace(options.MigrationsDir) == "" || options.Limit < 1 || options.Limit > 10000 || (options.Status != "" && options.Status != "abandoned" && options.Status != "active" && options.Status != "failed" && options.Status != "interrupted" && options.Status != "succeeded") {
		return nil, invalidRun("run list input is invalid")
	}
	if _, err := journal.EnsureProjectionFresh(options.WorkspaceRoot, options.MigrationsDir); err != nil {
		return nil, err
	}
	conn, err := openWorkflowDB(options.WorkspaceRoot)
	if err != nil {
		return nil, err
	}
	defer closeWorkflowDB(conn)
	query, args := `SELECT id FROM runs WHERE 1 = 1`, []any{}
	if options.IssueID != "" {
		query, args = query+` AND issue_id = ?`, append(args, options.IssueID)
	}
	if options.EnvelopeID != "" {
		query, args = query+` AND envelope_id = ?`, append(args, options.EnvelopeID)
	}
	if options.Status != "" {
		query, args = query+` AND status = ?`, append(args, options.Status)
	}
	rows, err := conn.Query(query+` ORDER BY issue_id, id LIMIT ?`, append(args, options.Limit)...)
	if err != nil {
		return nil, fmt.Errorf("list runs: %w", err)
	}
	defer func() { _ = rows.Close() }()
	runs := []Run{}
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		run, err := loadRun(conn, id)
		if err != nil {
			return nil, err
		}
		runs = append(runs, run)
	}
	return runs, rows.Err()
}

func loadRun(conn *sql.DB, id string) (Run, error) {
	var run Run
	var start string
	var terminal, ended sql.NullString
	err := conn.QueryRow(`SELECT id, issue_id, envelope_id, envelope_revision, claim_id, claim_version, actor_id, start_repository_json, terminal_repository_json, status, started_at, ended_at FROM runs WHERE id = ?`, id).Scan(&run.ID, &run.IssueID, &run.EnvelopeID, &run.EnvelopeRevision, &run.ClaimID, &run.ClaimVersion, &run.ActorID, &start, &terminal, &run.Status, &run.StartedAt, &ended)
	if err == sql.ErrNoRows {
		return Run{}, &RunNotFoundError{ID: id}
	}
	if err != nil {
		return Run{}, fmt.Errorf("load run: %w", err)
	}
	if err := json.Unmarshal([]byte(start), &run.StartRepository); err != nil {
		return Run{}, fmt.Errorf("decode run start repository: %w", err)
	}
	if terminal.Valid {
		var repository Repository
		if err := json.Unmarshal([]byte(terminal.String), &repository); err != nil {
			return Run{}, fmt.Errorf("decode run terminal repository: %w", err)
		}
		run.TerminalRepository = &repository
	}
	if ended.Valid {
		run.EndedAt = &ended.String
	}
	var obligation sql.NullString
	if err := conn.QueryRow(`SELECT version, executor_id, obligation_id FROM run_runtime WHERE run_id = ?`, run.ID).Scan(&run.Version, &run.Executor.ID, &obligation); err != nil {
		return Run{}, fmt.Errorf("load run runtime: %w", err)
	}
	run.Executor.Source = "activated-harness"
	if obligation.Valid {
		run.ObligationID = &obligation.String
	}
	rows, err := conn.Query(`SELECT budget_json FROM run_budget_bindings WHERE run_id = ? ORDER BY name`, run.ID)
	if err != nil {
		return Run{}, fmt.Errorf("load run budgets: %w", err)
	}
	defer func() { _ = rows.Close() }()
	run.Budgets = []Budget{}
	for rows.Next() {
		var raw string
		var budget Budget
		if err := rows.Scan(&raw); err != nil || json.Unmarshal([]byte(raw), &budget) != nil {
			return Run{}, fmt.Errorf("decode run budget")
		}
		run.Budgets = append(run.Budgets, budget)
	}
	if err := rows.Err(); err != nil {
		return Run{}, err
	}
	if err := conn.QueryRow(`SELECT prior_run_id FROM run_recoveries WHERE successor_run_id = ?`, run.ID).Scan(&run.PredecessorRunID); err != nil && err != sql.ErrNoRows {
		return Run{}, fmt.Errorf("load run predecessor: %w", err)
	}
	if err := conn.QueryRow(`SELECT successor_run_id FROM run_recoveries WHERE prior_run_id = ?`, run.ID).Scan(&run.SuccessorRunID); err != nil && err != sql.ErrNoRows {
		return Run{}, fmt.Errorf("load run successor: %w", err)
	}
	return run, nil
}

func sameRepository(left, right Repository) bool {
	return sameStringPointer(left.BaseOID, right.BaseOID) && left.CaptureSource == right.CaptureSource && left.Dirty == right.Dirty && sameStringPointer(left.DirtyDigest, right.DirtyDigest) && left.HeadOID == right.HeadOID && left.ProjectID == right.ProjectID && left.TreeOID == right.TreeOID && left.WorktreeID == right.WorktreeID
}

func sameRepositoryLocation(left, right Repository) bool {
	return validEnvelopeRepository(right) && left.ProjectID == right.ProjectID && left.WorktreeID == right.WorktreeID
}

func sameStringPointer(left, right *string) bool {
	if left == nil || right == nil {
		return left == right
	}
	return *left == *right
}
