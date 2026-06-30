package issues

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
)

func requireCompletionIssue(conn *sql.DB, issueID string) error {
	kind, found, err := issueKindForID(conn, issueID)
	if err != nil {
		return err
	}
	if !found {
		return &NotFoundError{IssueID: issueID}
	}
	if kind == "task" || kind == "subtask" {
		return nil
	}
	return fmt.Errorf("complete issue: issue %q is kind %q; task completions require task or subtask", issueID, kind)
}

func completionIDExists(conn *sql.DB, completionID string) (bool, error) {
	var existing string
	err := conn.QueryRow(`SELECT id FROM task_completions WHERE id = ? LIMIT 1`, completionID).Scan(&existing)
	if err == nil {
		return true, nil
	}
	if err == sql.ErrNoRows {
		return false, nil
	}
	return false, fmt.Errorf("query completion id %q: %w", completionID, err)
}

func nextCompletionAttemptNumber(conn *sql.DB, issueID string) (int, error) {
	var attemptNumber int
	if err := conn.QueryRow(`
		SELECT COALESCE(MAX(attempt_number), 0) + 1
		FROM task_completions
		WHERE issue_id = ?
	`, issueID).Scan(&attemptNumber); err != nil {
		return 0, fmt.Errorf("query next completion attempt for %q: %w", issueID, err)
	}
	return attemptNumber, nil
}

func priorCompletionIDs(conn *sql.DB, issueID string) ([]string, error) {
	rows, err := conn.Query(`
		SELECT id
		FROM task_completions
		WHERE issue_id = ?
		ORDER BY attempt_number, id
	`, issueID)
	if err != nil {
		return nil, fmt.Errorf("query prior completions for %q: %w", issueID, err)
	}
	defer closeIssueRows(rows)

	result := []string{}
	for rows.Next() {
		var completionID string
		if err := rows.Scan(&completionID); err != nil {
			return nil, fmt.Errorf("scan prior completion for %q: %w", issueID, err)
		}
		result = append(result, completionID)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate prior completions for %q: %w", issueID, err)
	}
	return result, nil
}

func completionGateSnapshot(conn *sql.DB, issueID string, attemptNumber int) (string, int, int, error) {
	rows, err := conn.Query(`
		SELECT command, exit_code, ran_at, output_excerpt
		FROM gate_runs
		WHERE issue_id = ? AND attempt_number = ?
		ORDER BY ran_at, id
	`, issueID, attemptNumber)
	if err != nil {
		return "", 0, 0, fmt.Errorf("query gate runs for %q attempt %d: %w", issueID, attemptNumber, err)
	}
	defer closeIssueRows(rows)

	gateResults := []CompletionGateResult{}
	gatesPassed := 0
	gatesFailed := 0
	for rows.Next() {
		var gateResult CompletionGateResult
		var outputExcerpt sql.NullString
		if err := rows.Scan(&gateResult.Cmd, &gateResult.ExitCode, &gateResult.RanAt, &outputExcerpt); err != nil {
			return "", 0, 0, fmt.Errorf("scan gate run for %q attempt %d: %w", issueID, attemptNumber, err)
		}
		gateResult.OutputExcerpt = nullStringPointer(outputExcerpt)
		if gateResult.ExitCode == 0 {
			gatesPassed++
		} else {
			gatesFailed++
		}
		gateResults = append(gateResults, gateResult)
	}
	if err := rows.Err(); err != nil {
		return "", 0, 0, fmt.Errorf("iterate gate runs for %q attempt %d: %w", issueID, attemptNumber, err)
	}

	gateResultsJSON, err := json.Marshal(gateResults)
	if err != nil {
		return "", 0, 0, fmt.Errorf("marshal gate results for %q attempt %d: %w", issueID, attemptNumber, err)
	}
	return string(gateResultsJSON), gatesPassed, gatesFailed, nil
}

func marshalStringSlice(values []string) (string, error) {
	encoded, err := json.Marshal(values)
	if err != nil {
		return "", fmt.Errorf("marshal string slice: %w", err)
	}
	return string(encoded), nil
}

func compactTrimmedStrings(values []string) []string {
	result := make([]string, 0, len(values))
	for _, value := range values {
		trimmed := strings.TrimSpace(value)
		if trimmed == "" {
			continue
		}
		result = append(result, trimmed)
	}
	return result
}

func issueCompleteArgv(
	issueID, actorRole string,
	notesMDPath *string,
	actualTokens *int,
	commitHashes, filesTouched []string,
	humanIntervention bool,
	humanInterventionNote *string,
) []string {
	argv := []string{"fid", "issue", "complete", issueID, "--by", actorRole}
	if notesMDPath != nil {
		argv = append(argv, "--notes-md", *notesMDPath)
	}
	if actualTokens != nil {
		argv = append(argv, "--tokens", fmt.Sprintf("%d", *actualTokens))
	}
	if len(commitHashes) > 0 {
		argv = append(argv, "--commits", strings.Join(commitHashes, ", "))
	}
	for _, filePath := range filesTouched {
		argv = append(argv, "--files", filePath)
	}
	if humanIntervention {
		argv = append(argv, "--human-intervention")
	}
	if humanInterventionNote != nil {
		argv = append(argv, "--intervention-note", *humanInterventionNote)
	}
	return argv
}
