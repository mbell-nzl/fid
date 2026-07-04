package issues

import (
	"database/sql"
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

type CompletionBuildOptions struct {
	IssueID               string
	ActorRole             string
	Status                string
	NotesMDPath           *string
	ActualTokens          *int
	CommitHashes          []string
	FilesTouched          []string
	HumanIntervention     bool
	HumanInterventionNote *string
}

func Complete(options CompleteOptions) (CompleteResult, error) {
	buildOptions, commandArgv, err := prepareCompleteInputs(options)
	if err != nil {
		return CompleteResult{}, err
	}
	if err := policy.Authorize(options.WorkspaceRoot, options.ActorRole, actionIssueComplete); err != nil {
		return CompleteResult{}, err
	}
	if _, err := journal.EnsureProjectionFresh(options.WorkspaceRoot, options.MigrationsDir); err != nil {
		return CompleteResult{}, err
	}
	return writeCompletionEvent(
		options.WorkspaceRoot,
		options.ActorRole,
		commandArgv,
		"issue.completed",
		buildOptions,
	)
}

func prepareCompleteInputs(options CompleteOptions) (CompletionBuildOptions, []string, error) {
	if strings.TrimSpace(options.WorkspaceRoot) == "" {
		return CompletionBuildOptions{}, nil, fmt.Errorf("complete issue: workspace root is required")
	}
	if strings.TrimSpace(options.MigrationsDir) == "" {
		return CompletionBuildOptions{}, nil, fmt.Errorf("complete issue: migrations dir is required")
	}
	issueID := strings.TrimSpace(options.IssueID)
	if issueID == "" {
		return CompletionBuildOptions{}, nil, fmt.Errorf("complete issue: issue id is required")
	}
	buildOptions := CompletionBuildOptions{
		IssueID:               issueID,
		ActorRole:             options.ActorRole,
		Status:                "dev_complete",
		NotesMDPath:           trimmedOptionalString(options.NotesMDPath),
		ActualTokens:          options.ActualTokens,
		CommitHashes:          compactTrimmedStrings(options.CommitHashes),
		FilesTouched:          compactTrimmedStrings(options.FilesTouched),
		HumanIntervention:     options.HumanIntervention,
		HumanInterventionNote: trimmedOptionalString(options.HumanInterventionNote),
	}
	return buildOptions, issueCompleteArgv(
		issueID,
		options.ActorRole,
		buildOptions.NotesMDPath,
		buildOptions.ActualTokens,
		buildOptions.CommitHashes,
		buildOptions.FilesTouched,
		buildOptions.HumanIntervention,
		buildOptions.HumanInterventionNote,
	), nil
}

func writeCompletionEvent(
	workspaceRoot string,
	actorRole string,
	commandArgv []string,
	eventType string,
	buildOptions CompletionBuildOptions,
) (CompleteResult, error) {
	var result CompleteResult
	err := journal.WithLock(workspaceRoot, journal.LockOwner{
		Command:   strings.Join(commandArgv, " "),
		ActorRole: actorRole,
	}, func(lock *journal.WriteLock) error {
		conn, err := openCompletionDB(workspaceRoot)
		if err != nil {
			return err
		}
		defer closeIssueDB(conn)

		completionResult, err := appendCompletionEvent(
			workspaceRoot,
			conn,
			lock,
			actorRole,
			commandArgv,
			eventType,
			buildOptions,
		)
		if err != nil {
			return err
		}
		result = completionResult
		return nil
	})
	if err != nil {
		return CompleteResult{}, err
	}

	return result, nil
}

func openCompletionDB(workspaceRoot string) (*sql.DB, error) {
	dbPath := filepath.Join(workspaceRoot, config.WorkspaceDirName, config.DatabaseFileName)
	return fiddb.Open(dbPath)
}

func appendCompletionEvent(
	workspaceRoot string,
	conn *sql.DB,
	lock *journal.WriteLock,
	actorRole string,
	commandArgv []string,
	eventType string,
	buildOptions CompletionBuildOptions,
) (CompleteResult, error) {
	event, err := journal.PrepareEvent(lock, journal.Event{
		EventSchema: journal.CurrentEventSchema,
		EventType:   eventType,
		Actor:       journal.Actor{Role: actorRole},
		Command: journal.Command{
			Argv:       commandArgv,
			CWD:        ".",
			FIDVersion: fidVersion,
		},
	})
	if err != nil {
		return CompleteResult{}, err
	}

	timestamp := event.OccurredAt.UTC().Format(time.RFC3339)
	effects, completionResult, err := BuildCompletionEffects(conn, buildOptions, event.EventID, timestamp)
	if err != nil {
		return CompleteResult{}, err
	}
	event.Effects = effects
	if err := journal.DryRunProjectEvent(conn, event); err != nil {
		return CompleteResult{}, err
	}
	appended, err := journal.AppendEvent(lock, event)
	if err != nil {
		return CompleteResult{}, err
	}
	return completionResult, projectCompletionEvent(workspaceRoot, conn, appended)
}

func projectCompletionEvent(workspaceRoot string, conn *sql.DB, appended journal.Event) error {
	eventHash, err := journal.HashEvent(appended)
	if err != nil {
		return fmt.Errorf("hash issue event %q: %w", appended.EventID, err)
	}
	journalFile := appended.OccurredAt.UTC().Format("2006-01-02") + ".jsonl"
	journalPath := filepath.Join(workspaceRoot, config.WorkspaceDirName, config.JournalDirName, journalFile)
	info, err := os.Stat(journalPath)
	if err != nil {
		return fmt.Errorf("stat issue journal %q: %w", journalPath, err)
	}
	return journal.ProjectEvent(conn, journal.ProjectionInput{
		Event:     appended,
		EventHash: eventHash,
		File:      journalFile,
		Line:      countJournalLines(journalPath),
		Offset:    info.Size(),
	})
}

func BuildCompletionEffects(
	conn *sql.DB,
	options CompletionBuildOptions,
	eventID string,
	completedAt string,
) ([]journal.Effect, CompleteResult, error) {
	status := strings.TrimSpace(options.Status)
	if status == "" {
		return nil, CompleteResult{}, fmt.Errorf("build completion effects: status is required")
	}
	if err := requireCompletionIssue(conn, options.IssueID); err != nil {
		return nil, CompleteResult{}, err
	}

	data, err := loadCompletionBuildData(conn, options.IssueID, completedAt)
	if err != nil {
		return nil, CompleteResult{}, err
	}
	values, err := buildCompletionValues(options, data)
	if err != nil {
		return nil, CompleteResult{}, err
	}
	return completionEffects(data, options, eventID, values),
		completeResult(options, data),
		nil
}

type completionBuildData struct {
	ActivityID         string
	CompletionID       string
	AttemptNumber      int
	PriorCompletionIDs []string
	GateResultsJSON    string
	GatesPassed        int
	GatesFailed        int
	CompletedAt        string
}

func loadCompletionBuildData(conn *sql.DB, issueID, completedAt string) (completionBuildData, error) {
	activityID, err := nextIssueActivityID(conn)
	if err != nil {
		return completionBuildData{}, err
	}
	completionID, err := journal.GenerateID("TC", func(candidate string) (bool, error) {
		return completionIDExists(conn, candidate)
	})
	if err != nil {
		return completionBuildData{}, err
	}
	attemptNumber, err := nextCompletionAttemptNumber(conn, issueID)
	if err != nil {
		return completionBuildData{}, err
	}
	priorIDs, err := priorCompletionIDs(conn, issueID)
	if err != nil {
		return completionBuildData{}, err
	}
	gateResultsJSON, gatesPassed, gatesFailed, err := completionGateSnapshot(conn, issueID, attemptNumber)
	if err != nil {
		return completionBuildData{}, err
	}
	return completionBuildData{
		ActivityID:         activityID,
		CompletionID:       completionID,
		AttemptNumber:      attemptNumber,
		PriorCompletionIDs: priorIDs,
		GateResultsJSON:    gateResultsJSON,
		GatesPassed:        gatesPassed,
		GatesFailed:        gatesFailed,
		CompletedAt:        completedAt,
	}, nil
}

func buildCompletionValues(options CompletionBuildOptions, data completionBuildData) (map[string]any, error) {
	values := map[string]any{
		"issue_id":           options.IssueID,
		"attempt_number":     data.AttemptNumber,
		"developer_role":     options.ActorRole,
		"status":             options.Status,
		"gates_passed":       data.GatesPassed,
		"gates_failed":       data.GatesFailed,
		"gate_results":       data.GateResultsJSON,
		"human_intervention": options.HumanIntervention,
		"completed_at":       data.CompletedAt,
	}
	if err := addOptionalCompletionValues(values, options); err != nil {
		return nil, err
	}
	return values, nil
}

func addOptionalCompletionValues(values map[string]any, options CompletionBuildOptions) error {
	if options.NotesMDPath != nil {
		values["notes_md_path"] = *options.NotesMDPath
	}
	if options.ActualTokens != nil {
		values["actual_tokens"] = *options.ActualTokens
	}
	if len(options.CommitHashes) > 0 {
		commitHashesJSON, err := marshalStringSlice(compactTrimmedStrings(options.CommitHashes))
		if err != nil {
			return err
		}
		values["commit_hashes"] = commitHashesJSON
	}
	if len(options.FilesTouched) > 0 {
		filesTouchedJSON, err := marshalStringSlice(compactTrimmedStrings(options.FilesTouched))
		if err != nil {
			return err
		}
		values["files_touched"] = filesTouchedJSON
	}
	if options.HumanIntervention && options.HumanInterventionNote != nil {
		values["human_intervention_note"] = *options.HumanInterventionNote
	}
	return nil
}

func completionEffects(
	data completionBuildData,
	options CompletionBuildOptions,
	eventID string,
	values map[string]any,
) []journal.Effect {
	effects := []journal.Effect{{
		Op:     "insert",
		Table:  "task_completions",
		PK:     map[string]any{"id": data.CompletionID},
		Values: values,
	}, issueActivityEffect(
		data.ActivityID,
		eventID,
		options.ActorRole,
		data.CompletedAt,
		options.IssueID,
		"task_completed",
		map[string]any{
			"completion_id":  data.CompletionID,
			"status":         options.Status,
			"attempt_number": data.AttemptNumber,
		},
	)}
	for _, priorID := range data.PriorCompletionIDs {
		effects = append(effects, journal.Effect{
			Op:    "update",
			Table: "task_completions",
			PK:    map[string]any{"id": priorID},
			Values: map[string]any{
				"status": "superseded",
			},
		})
	}
	return effects
}

func completeResult(options CompletionBuildOptions, data completionBuildData) CompleteResult {
	return CompleteResult{
		ID:                    data.CompletionID,
		IssueID:               options.IssueID,
		AttemptNumber:         data.AttemptNumber,
		Status:                options.Status,
		NotesMDPath:           trimmedOptionalString(options.NotesMDPath),
		ActualTokens:          options.ActualTokens,
		CommitHashes:          compactTrimmedStrings(options.CommitHashes),
		FilesTouched:          compactTrimmedStrings(options.FilesTouched),
		HumanIntervention:     options.HumanIntervention,
		HumanInterventionNote: trimmedOptionalString(options.HumanInterventionNote),
		CompletedAt:           data.CompletedAt,
	}
}
