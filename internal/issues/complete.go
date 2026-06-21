package issues

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/mbell-nzl/fid/internal/config"
	fiddb "github.com/mbell-nzl/fid/internal/db"
	"github.com/mbell-nzl/fid/internal/journal"
	"github.com/mbell-nzl/fid/internal/policy"
)

type CompleteOptions struct {
	WorkspaceRoot         string
	MigrationsDir         string
	IssueID               string
	ActorRole             string
	NotesMDPath           *string
	ActualTokens          *int
	CommitHashes          []string
	FilesTouched          []string
	HumanIntervention     bool
	HumanInterventionNote *string
}

type CompleteResult struct {
	ID                    string
	IssueID               string
	AttemptNumber         int
	Status                string
	NotesMDPath           *string
	ActualTokens          *int
	CommitHashes          []string
	FilesTouched          []string
	HumanIntervention     bool
	HumanInterventionNote *string
	CompletedAt           string
}

func Complete(options CompleteOptions) (CompleteResult, error) {
	if strings.TrimSpace(options.WorkspaceRoot) == "" {
		return CompleteResult{}, fmt.Errorf("complete issue: workspace root is required")
	}
	if strings.TrimSpace(options.MigrationsDir) == "" {
		return CompleteResult{}, fmt.Errorf("complete issue: migrations dir is required")
	}
	issueID := strings.TrimSpace(options.IssueID)
	if issueID == "" {
		return CompleteResult{}, fmt.Errorf("complete issue: issue id is required")
	}
	if err := policy.Authorize(options.WorkspaceRoot, options.ActorRole, actionIssueComplete); err != nil {
		return CompleteResult{}, err
	}

	notesMDPath := trimmedOptionalString(options.NotesMDPath)
	actualTokens := options.ActualTokens
	commitHashes := compactTrimmedStrings(options.CommitHashes)
	filesTouched := compactTrimmedStrings(options.FilesTouched)
	humanInterventionNote := trimmedOptionalString(options.HumanInterventionNote)

	if _, err := journal.EnsureProjectionFresh(options.WorkspaceRoot, options.MigrationsDir); err != nil {
		return CompleteResult{}, err
	}

	commandArgv := issueCompleteArgv(
		issueID,
		options.ActorRole,
		notesMDPath,
		actualTokens,
		commitHashes,
		filesTouched,
		options.HumanIntervention,
		humanInterventionNote,
	)

	var result CompleteResult
	err := journal.WithLock(options.WorkspaceRoot, journal.LockOwner{
		Command:   strings.Join(commandArgv, " "),
		ActorRole: options.ActorRole,
	}, func(lock *journal.WriteLock) error {
		dbPath := filepath.Join(options.WorkspaceRoot, config.WorkspaceDirName, config.DatabaseFileName)
		conn, err := fiddb.Open(dbPath)
		if err != nil {
			return err
		}
		defer closeIssueDB(conn)

		if err := requireCompletionIssue(conn, issueID); err != nil {
			return err
		}

		completionID, err := journal.GenerateID("TC", func(candidate string) (bool, error) {
			return completionIDExists(conn, candidate)
		})
		if err != nil {
			return err
		}

		attemptNumber, err := nextCompletionAttemptNumber(conn, issueID)
		if err != nil {
			return err
		}
		priorCompletionIDs, err := priorCompletionIDs(conn, issueID)
		if err != nil {
			return err
		}
		gateResultsJSON, gatesPassed, gatesFailed, err := completionGateSnapshot(conn, issueID, attemptNumber)
		if err != nil {
			return err
		}

		commitHashesJSON, err := marshalStringSlice(commitHashes)
		if err != nil {
			return err
		}
		filesTouchedJSON, err := marshalStringSlice(filesTouched)
		if err != nil {
			return err
		}

		completedAt := time.Now().UTC().Format(time.RFC3339)
		values := map[string]any{
			"issue_id":           issueID,
			"attempt_number":     attemptNumber,
			"developer_role":     options.ActorRole,
			"status":             "dev_complete",
			"gates_passed":       gatesPassed,
			"gates_failed":       gatesFailed,
			"gate_results":       gateResultsJSON,
			"human_intervention": options.HumanIntervention,
			"completed_at":       completedAt,
		}
		if notesMDPath != nil {
			values["notes_md_path"] = *notesMDPath
		}
		if actualTokens != nil {
			values["actual_tokens"] = *actualTokens
		}
		if len(commitHashes) > 0 {
			values["commit_hashes"] = commitHashesJSON
		}
		if len(filesTouched) > 0 {
			values["files_touched"] = filesTouchedJSON
		}
		if options.HumanIntervention && humanInterventionNote != nil {
			values["human_intervention_note"] = *humanInterventionNote
		}

		effects := make([]journal.Effect, 0, 1+len(priorCompletionIDs))
		effects = append(effects, journal.Effect{
			Op:     "insert",
			Table:  "task_completions",
			PK:     map[string]any{"id": completionID},
			Values: values,
		})
		for _, priorID := range priorCompletionIDs {
			effects = append(effects, journal.Effect{
				Op:    "update",
				Table: "task_completions",
				PK:    map[string]any{"id": priorID},
				Values: map[string]any{
					"status": "superseded",
				},
			})
		}

		event := journal.Event{
			EventSchema: journal.CurrentEventSchema,
			EventType:   "issue.completed",
			Actor: journal.Actor{
				Role: options.ActorRole,
			},
			Command: journal.Command{
				Argv:       commandArgv,
				CWD:        ".",
				FIDVersion: fidVersion,
			},
			Effects: effects,
		}

		if err := journal.DryRunProjectEvent(conn, event); err != nil {
			return err
		}

		appended, err := journal.AppendEvent(lock, event)
		if err != nil {
			return err
		}

		eventHash, err := journal.HashEvent(appended)
		if err != nil {
			return fmt.Errorf("hash issue event %q: %w", appended.EventID, err)
		}

		journalFile := appended.OccurredAt.UTC().Format("2006-01-02") + ".jsonl"
		journalPath := filepath.Join(options.WorkspaceRoot, config.WorkspaceDirName, config.JournalDirName, journalFile)
		info, err := os.Stat(journalPath)
		if err != nil {
			return fmt.Errorf("stat issue journal %q: %w", journalPath, err)
		}

		if err := journal.ProjectEvent(conn, journal.ProjectionInput{
			Event:     appended,
			EventHash: eventHash,
			File:      journalFile,
			Line:      countJournalLines(journalPath),
			Offset:    info.Size(),
		}); err != nil {
			return err
		}

		result = CompleteResult{
			ID:                    completionID,
			IssueID:               issueID,
			AttemptNumber:         attemptNumber,
			Status:                "dev_complete",
			NotesMDPath:           notesMDPath,
			ActualTokens:          actualTokens,
			CommitHashes:          commitHashes,
			FilesTouched:          filesTouched,
			HumanIntervention:     options.HumanIntervention,
			HumanInterventionNote: humanInterventionNote,
			CompletedAt:           completedAt,
		}
		return nil
	})
	if err != nil {
		return CompleteResult{}, err
	}

	return result, nil
}

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
