package runs

import (
	"database/sql"
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/mbell-nzl/fid/internal/config"
	fiddb "github.com/mbell-nzl/fid/internal/db"
	"github.com/mbell-nzl/fid/internal/issues"
	"github.com/mbell-nzl/fid/internal/journal"
	"github.com/mbell-nzl/fid/internal/policy"
)

const (
	actionRunRecord = "run.record"
	fidVersion      = "0.1.0"
)

type RecordOptions struct {
	WorkspaceRoot      string
	MigrationsDir      string
	IssueID            string
	CompletionID       *string
	ActorRole          string
	RunRole            string
	Model              string
	TokensIn           *int
	TokensOut          *int
	TokenBreakdownJSON *string
	ToolCallCount      *int
	ProviderRunID      *string
	ProviderCostUSD    *float64
}

type RecordResult struct {
	ID                 string
	IssueID            string
	CompletionID       *string
	RunRole            string
	Model              string
	TokensIn           *int
	TokensOut          *int
	PriceInPerMTok     *float64
	PriceOutPerMTok    *float64
	CostUSD            *float64
	TokenBreakdownJSON *string
	ToolCallCount      *int
	ProviderRunID      *string
	ProviderCostUSD    *float64
	StartedAt          string
	EndedAt            string
	Warning            *string
}

func Record(options RecordOptions) (RecordResult, error) {
	issueID, model, completionID, tokenBreakdownJSON, providerRunID, err := validateRecordOptions(options)
	if err != nil {
		return RecordResult{}, err
	}
	if err := policy.Authorize(options.WorkspaceRoot, options.ActorRole, actionRunRecord); err != nil {
		return RecordResult{}, err
	}
	runRole := strings.TrimSpace(options.RunRole)
	if err := policy.RequireKnownRole(options.WorkspaceRoot, runRole); err != nil {
		return RecordResult{}, err
	}

	commandArgv := runRecordArgv(issueID, completionID, runRole, model, options)
	pricingPath := config.WorkspacePathsForRoot(options.WorkspaceRoot).PricingPath

	var result RecordResult
	err = journal.WithLock(options.WorkspaceRoot, journal.LockOwner{
		Command:   strings.Join(commandArgv, " "),
		ActorRole: options.ActorRole,
	}, func(lock *journal.WriteLock) error {
		if _, err := journal.EnsureProjectionFreshWithLock(lock, options.MigrationsDir); err != nil {
			return err
		}

		conn, err := fiddb.Open(filepath.Join(options.WorkspaceRoot, config.WorkspaceDirName, config.DatabaseFileName))
		if err != nil {
			return err
		}
		defer closeRunsConn(conn)

		if err := requireRunIssue(conn, issueID); err != nil {
			return err
		}
		runID, err := journal.GenerateID("RUN", func(candidate string) (bool, error) {
			return runIDExists(conn, candidate)
		})
		if err != nil {
			return err
		}
		pricing, err := config.ReadPricingConfig(pricingPath)
		if err != nil {
			return err
		}

		event, built, err := buildRecordEvent(lock, conn, runID, pricing, commandArgv, issueID, completionID, runRole, model, tokenBreakdownJSON, providerRunID, options)
		if err != nil {
			return err
		}
		if err := journal.DryRunProjectEvent(conn, event); err != nil {
			return err
		}

		appended, err := journal.AppendEvent(lock, event)
		if err != nil {
			return err
		}
		if err := projectRunsEvent(options.WorkspaceRoot, conn, appended); err != nil {
			return err
		}

		result = built
		return nil
	})
	if err != nil {
		return RecordResult{}, err
	}

	return result, nil
}

func validateRecordOptions(options RecordOptions) (string, string, *string, *string, *string, error) {
	if strings.TrimSpace(options.WorkspaceRoot) == "" {
		return "", "", nil, nil, nil, fmt.Errorf("record run: workspace root is required")
	}
	if strings.TrimSpace(options.MigrationsDir) == "" {
		return "", "", nil, nil, nil, fmt.Errorf("record run: migrations dir is required")
	}
	issueID := strings.TrimSpace(options.IssueID)
	if issueID == "" {
		return "", "", nil, nil, nil, fmt.Errorf("record run: issue id is required")
	}
	if strings.TrimSpace(options.RunRole) == "" {
		return "", "", nil, nil, nil, fmt.Errorf("record run: role is required")
	}
	model := strings.TrimSpace(options.Model)
	if model == "" {
		return "", "", nil, nil, nil, fmt.Errorf("record run: model is required")
	}
	if err := validateNonNegativeInt("tokens_in", options.TokensIn); err != nil {
		return "", "", nil, nil, nil, err
	}
	if err := validateNonNegativeInt("tokens_out", options.TokensOut); err != nil {
		return "", "", nil, nil, nil, err
	}
	if err := validateNonNegativeInt("tool_call_count", options.ToolCallCount); err != nil {
		return "", "", nil, nil, nil, err
	}
	if options.ProviderCostUSD != nil && *options.ProviderCostUSD < 0 {
		return "", "", nil, nil, nil, fmt.Errorf("record run: provider_cost_usd must be >= 0")
	}
	return issueID, model, trimmedOptionalString(options.CompletionID), trimmedOptionalString(options.TokenBreakdownJSON), trimmedOptionalString(options.ProviderRunID), nil
}

func buildRecordEvent(
	lock *journal.WriteLock,
	conn *sql.DB,
	runID string,
	pricing config.PricingConfig,
	commandArgv []string,
	issueID string,
	completionID *string,
	runRole string,
	model string,
	tokenBreakdownJSON *string,
	providerRunID *string,
	options RecordOptions,
) (journal.Event, RecordResult, error) {
	event, err := journal.PrepareEvent(lock, journal.Event{
		EventSchema: journal.CurrentEventSchema,
		EventType:   "run.recorded",
		Actor:       journal.Actor{Role: options.ActorRole},
		Command: journal.Command{
			Argv:       commandArgv,
			CWD:        ".",
			FIDVersion: fidVersion,
		},
	})
	if err != nil {
		return journal.Event{}, RecordResult{}, err
	}

	startedAt := event.OccurredAt.UTC().Format(time.RFC3339)
	values, result := recordValues(issueID, completionID, runRole, model, tokenBreakdownJSON, providerRunID, startedAt, options)
	result.ID = runID
	if pricingModel, ok := pricing.Models[model]; ok {
		priceIn := pricingModel.InPerMTok
		priceOut := pricingModel.OutPerMTok
		costUSD := computeCost(optionalIntValue(options.TokensIn), optionalIntValue(options.TokensOut), priceIn, priceOut)
		values["price_in_per_mtok"] = priceIn
		values["price_out_per_mtok"] = priceOut
		values["cost_usd"] = costUSD
		result.PriceInPerMTok = floatPointer(priceIn)
		result.PriceOutPerMTok = floatPointer(priceOut)
		result.CostUSD = floatPointer(costUSD)
	} else {
		result.Warning = stringPointer(fmt.Sprintf("warning: model %q not found in %s; recorded NULL snapshot prices and cost", model, filepath.ToSlash(filepath.Join(config.WorkspaceDirName, config.PricingFileName))))
	}

	event.Effects = []journal.Effect{{
		Op:     "insert",
		Table:  "agent_runs",
		PK:     map[string]any{"id": runID},
		Values: values,
	}}
	if err := requireRunCompletion(conn, completionID); err != nil {
		return journal.Event{}, RecordResult{}, err
	}

	return event, result, nil
}

func recordValues(
	issueID string,
	completionID *string,
	runRole string,
	model string,
	tokenBreakdownJSON *string,
	providerRunID *string,
	timestamp string,
	options RecordOptions,
) (map[string]any, RecordResult) {
	values := map[string]any{
		"issue_id":   issueID,
		"role":       runRole,
		"model":      model,
		"started_at": timestamp,
		"ended_at":   timestamp,
	}
	result := RecordResult{
		ID:                 "",
		IssueID:            issueID,
		CompletionID:       completionID,
		RunRole:            runRole,
		Model:              model,
		TokensIn:           cloneOptionalInt(options.TokensIn),
		TokensOut:          cloneOptionalInt(options.TokensOut),
		TokenBreakdownJSON: tokenBreakdownJSON,
		ToolCallCount:      cloneOptionalInt(options.ToolCallCount),
		ProviderRunID:      providerRunID,
		ProviderCostUSD:    cloneOptionalFloat(options.ProviderCostUSD),
		StartedAt:          timestamp,
		EndedAt:            timestamp,
	}

	if completionID != nil {
		values["completion_id"] = *completionID
	}
	if options.TokensIn != nil {
		values["tokens_in"] = *options.TokensIn
	}
	if options.TokensOut != nil {
		values["tokens_out"] = *options.TokensOut
	}
	if tokenBreakdownJSON != nil {
		values["token_breakdown_json"] = *tokenBreakdownJSON
	}
	if options.ToolCallCount != nil {
		values["tool_call_count"] = *options.ToolCallCount
	}
	if providerRunID != nil {
		values["provider_run_id"] = *providerRunID
	}
	if options.ProviderCostUSD != nil {
		values["provider_cost_usd"] = *options.ProviderCostUSD
	}

	return values, result
}

func requireRunIssue(conn *sql.DB, issueID string) error {
	var existing string
	err := conn.QueryRow(`SELECT id FROM issues WHERE id = ? LIMIT 1`, issueID).Scan(&existing)
	if err == nil {
		return nil
	}
	if errors.Is(err, sql.ErrNoRows) {
		return &issues.NotFoundError{IssueID: issueID}
	}
	return fmt.Errorf("query issue id %q: %w", issueID, err)
}

func requireRunCompletion(conn *sql.DB, completionID *string) error {
	if completionID == nil {
		return nil
	}
	var existing string
	err := conn.QueryRow(`SELECT id FROM task_completions WHERE id = ? LIMIT 1`, *completionID).Scan(&existing)
	if err == nil {
		return nil
	}
	if errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("record run: completion %q not found", *completionID)
	}
	return fmt.Errorf("query completion id %q: %w", *completionID, err)
}

func runIDExists(conn *sql.DB, runID string) (bool, error) {
	var existing string
	err := conn.QueryRow(`SELECT id FROM agent_runs WHERE id = ? LIMIT 1`, runID).Scan(&existing)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	return false, fmt.Errorf("query run id %q: %w", runID, err)
}

func runRecordArgv(issueID string, completionID *string, runRole string, model string, options RecordOptions) []string {
	argv := []string{"fid", "run", "record", "--issue", issueID, "--role", runRole, "--model", model}
	if options.TokensIn != nil {
		argv = append(argv, "--tokens-in", fmt.Sprintf("%d", *options.TokensIn))
	}
	if options.TokensOut != nil {
		argv = append(argv, "--tokens-out", fmt.Sprintf("%d", *options.TokensOut))
	}
	if tokenBreakdownJSON := trimmedOptionalString(options.TokenBreakdownJSON); tokenBreakdownJSON != nil {
		argv = append(argv, "--token-breakdown-json", *tokenBreakdownJSON)
	}
	if options.ToolCallCount != nil {
		argv = append(argv, "--tool-call-count", fmt.Sprintf("%d", *options.ToolCallCount))
	}
	if providerRunID := trimmedOptionalString(options.ProviderRunID); providerRunID != nil {
		argv = append(argv, "--provider-run-id", *providerRunID)
	}
	if options.ProviderCostUSD != nil {
		argv = append(argv, "--provider-cost-usd", fmt.Sprintf("%g", *options.ProviderCostUSD))
	}
	if completionID != nil {
		argv = append(argv, "--completion", *completionID)
	}
	return append(argv, "--by", options.ActorRole)
}

func projectRunsEvent(workspaceRoot string, conn *sql.DB, event journal.Event) error {
	return projectEvent(workspaceRoot, conn, event, "run")
}
