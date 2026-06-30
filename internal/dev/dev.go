package dev

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/mbell-nzl/fid/internal/config"
	fiddb "github.com/mbell-nzl/fid/internal/db"
	"github.com/mbell-nzl/fid/internal/issues"
	"github.com/mbell-nzl/fid/internal/journal"
	"github.com/mbell-nzl/fid/internal/workflow"
)

const (
	actionDevStart    = "dev.start"
	actionDevGate     = "dev.gate"
	actionDevFinish   = "dev.finish"
	actionDevBlock    = "dev.block"
	defaultTTLSeconds = 1800
	fidVersion        = "0.1.0"
)

type issueState struct {
	ID         string
	Kind       string
	WorkflowID string
	Status     string
	UpdatedAt  string
}

type currentClaim struct {
	ClaimedBy  string
	ClaimedAt  time.Time
	TTLSeconds int
}

func withDevConn(
	workspaceRoot string,
	migrationsDir string,
	commandArgv []string,
	actorRole string,
	fn func(*journal.WriteLock, *sql.DB) error,
) error {
	return journal.WithLock(workspaceRoot, journal.LockOwner{
		Command:   strings.Join(commandArgv, " "),
		ActorRole: actorRole,
	}, func(lock *journal.WriteLock) error {
		if _, err := journal.EnsureProjectionFreshWithLock(lock, migrationsDir); err != nil {
			return err
		}

		conn, err := fiddb.Open(filepath.Join(workspaceRoot, config.WorkspaceDirName, config.DatabaseFileName))
		if err != nil {
			return err
		}
		defer closeDevConn(conn)

		return fn(lock, conn)
	})
}

func appendAndProjectEvent(workspaceRoot string, conn *sql.DB, lock *journal.WriteLock, event journal.Event) error {
	if err := journal.DryRunProjectEvent(conn, event); err != nil {
		return err
	}

	appended, err := journal.AppendEvent(lock, event)
	if err != nil {
		return err
	}

	eventHash, err := journal.HashEvent(appended)
	if err != nil {
		return fmt.Errorf("hash dev event %q: %w", appended.EventID, err)
	}

	journalFile := appended.OccurredAt.UTC().Format("2006-01-02") + ".jsonl"
	journalPath := filepath.Join(workspaceRoot, config.WorkspaceDirName, config.JournalDirName, journalFile)
	info, err := os.Stat(journalPath)
	if err != nil {
		return fmt.Errorf("stat dev journal %q: %w", journalPath, err)
	}

	return journal.ProjectEvent(conn, journal.ProjectionInput{
		Event:     appended,
		EventHash: eventHash,
		File:      journalFile,
		Line:      countJournalLines(journalPath),
		Offset:    info.Size(),
	})
}

func loadIssueState(conn *sql.DB, issueID string) (issueState, error) {
	var (
		state     issueState
		updatedAt sql.NullString
	)
	err := conn.QueryRow(`
		SELECT id, kind, workflow_id, status, updated_at
		FROM issues
		WHERE id = ?
	`, issueID).Scan(&state.ID, &state.Kind, &state.WorkflowID, &state.Status, &updatedAt)
	if err == sql.ErrNoRows {
		return issueState{}, &issues.NotFoundError{IssueID: issueID}
	}
	if err != nil {
		return issueState{}, fmt.Errorf("query issue %q: %w", issueID, err)
	}
	if !updatedAt.Valid || strings.TrimSpace(updatedAt.String) == "" {
		return issueState{}, fmt.Errorf("query issue %q: updated_at is required", issueID)
	}

	state.UpdatedAt = updatedAt.String
	return state, nil
}

func requireTaskOrSubtask(issueID, kind string) error {
	if kind == "task" || kind == "subtask" {
		return nil
	}
	return fmt.Errorf("dev command: issue %q is kind %q; developer composites require task or subtask", issueID, kind)
}

func loadCurrentClaim(conn *sql.DB, issueID string) (*currentClaim, error) {
	var (
		claimedBy string
		claimedAt string
		ttl       int
	)
	err := conn.QueryRow(`
		SELECT claimed_by, claimed_at, ttl_seconds
		FROM claims
		WHERE ref_type = 'issue' AND ref_id = ?
	`, issueID).Scan(&claimedBy, &claimedAt, &ttl)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("query claim for issue %q: %w", issueID, err)
	}

	parsedClaimedAt, err := time.Parse(time.RFC3339, claimedAt)
	if err != nil {
		return nil, fmt.Errorf("parse claim time for issue %q: %w", issueID, err)
	}

	return &currentClaim{
		ClaimedBy:  claimedBy,
		ClaimedAt:  parsedClaimedAt.UTC(),
		TTLSeconds: ttl,
	}, nil
}

func (claim *currentClaim) expiresAt() time.Time {
	return claim.ClaimedAt.Add(time.Duration(claim.TTLSeconds) * time.Second)
}

func (claim *currentClaim) isLiveAt(now time.Time) bool {
	return now.UTC().Before(claim.expiresAt())
}

func countJournalLines(path string) int {
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

func nextIssueUpdateTime(currentUpdatedAt string, now time.Time) time.Time {
	parsed, err := time.Parse(time.RFC3339, currentUpdatedAt)
	if err != nil {
		return now.UTC()
	}

	candidate := now.UTC()
	if candidate.Format(time.RFC3339) <= parsed.UTC().Format(time.RFC3339) {
		return parsed.UTC().Add(time.Second)
	}
	return candidate
}

func activityIDExists(conn *sql.DB, activityID string) (bool, error) {
	var existing string
	err := conn.QueryRow(`SELECT id FROM activity_log WHERE id = ? LIMIT 1`, activityID).Scan(&existing)
	if err == nil {
		return true, nil
	}
	if err == sql.ErrNoRows {
		return false, nil
	}
	return false, fmt.Errorf("query activity id %q: %w", activityID, err)
}

func runIDExists(conn *sql.DB, runID string) (bool, error) {
	var existing string
	err := conn.QueryRow(`SELECT id FROM agent_runs WHERE id = ? LIMIT 1`, runID).Scan(&existing)
	if err == nil {
		return true, nil
	}
	if err == sql.ErrNoRows {
		return false, nil
	}
	return false, fmt.Errorf("query run id %q: %w", runID, err)
}

func gateIDExists(conn *sql.DB, gateID string) (bool, error) {
	var existing string
	err := conn.QueryRow(`SELECT id FROM gate_runs WHERE id = ? LIMIT 1`, gateID).Scan(&existing)
	if err == nil {
		return true, nil
	}
	if err == sql.ErrNoRows {
		return false, nil
	}
	return false, fmt.Errorf("query gate id %q: %w", gateID, err)
}

func activeAttemptNumber(conn *sql.DB, issueID string) (int, error) {
	var attemptNumber int
	if err := conn.QueryRow(`
		SELECT COALESCE(MAX(attempt_number), 0) + 1
		FROM task_completions
		WHERE issue_id = ?
	`, issueID).Scan(&attemptNumber); err != nil {
		return 0, fmt.Errorf("query active attempt for %q: %w", issueID, err)
	}
	return attemptNumber, nil
}

func gateCounts(conn *sql.DB, issueID string, attemptNumber int) (int, int, error) {
	rows, err := conn.Query(`
		SELECT exit_code
		FROM gate_runs
		WHERE issue_id = ? AND attempt_number = ?
	`, issueID, attemptNumber)
	if err != nil {
		return 0, 0, fmt.Errorf("query gate counts for %q attempt %d: %w", issueID, attemptNumber, err)
	}
	defer closeDevRows(rows)

	gatesPassed := 0
	gatesFailed := 0
	for rows.Next() {
		var exitCode int
		if err := rows.Scan(&exitCode); err != nil {
			return 0, 0, fmt.Errorf("scan gate count for %q attempt %d: %w", issueID, attemptNumber, err)
		}
		if exitCode == 0 {
			gatesPassed++
		} else {
			gatesFailed++
		}
	}
	if err := rows.Err(); err != nil {
		return 0, 0, fmt.Errorf("iterate gate counts for %q attempt %d: %w", issueID, attemptNumber, err)
	}
	return gatesPassed, gatesFailed, nil
}

func openRunID(conn *sql.DB, issueID string) (*string, error) {
	var runID sql.NullString
	err := conn.QueryRow(`
		SELECT id
		FROM agent_runs
		WHERE issue_id = ? AND ended_at IS NULL
		ORDER BY started_at DESC, id DESC
		LIMIT 1
	`, issueID).Scan(&runID)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("query open run for %q: %w", issueID, err)
	}
	if !runID.Valid || strings.TrimSpace(runID.String) == "" {
		return nil, nil
	}
	return &runID.String, nil
}

func buildIssueTransitionEffects(current issueState, toState, actorRole, eventID, activityID, timestamp string) []journal.Effect {
	return []journal.Effect{
		{
			Op:    "update",
			Table: "issues",
			PK:    map[string]any{"id": current.ID},
			Values: map[string]any{
				"status":     toState,
				"updated_at": timestamp,
			},
			Prior: map[string]any{
				"status": current.Status,
			},
			Preconditions: map[string]any{
				"status":      current.Status,
				"workflow_id": current.WorkflowID,
			},
		},
		{
			Op:    "activity",
			Table: "activity_log",
			PK:    map[string]any{"id": activityID},
			Values: map[string]any{
				"journal_event_id": eventID,
				"ref_type":         "issue",
				"ref_id":           current.ID,
				"event":            "status_changed",
				"details": map[string]any{
					"from": current.Status,
					"to":   toState,
				},
				"actor":       actorRole,
				"occurred_at": timestamp,
			},
		},
	}
}

func issueStatusCategory(workflowID, state string) (string, error) {
	category, ok := workflow.CategoryForState(workflowID, state)
	if !ok {
		return "", fmt.Errorf("workflow %q does not define state %q", workflowID, state)
	}
	return category, nil
}

func closeDevConn(conn *sql.DB) {
	_ = conn.Close()
}

func closeDevRows(rows *sql.Rows) {
	_ = rows.Close()
}
