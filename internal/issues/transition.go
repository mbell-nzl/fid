package issues

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/mbell-nzl/fid/internal/config"
	fiddb "github.com/mbell-nzl/fid/internal/db"
	"github.com/mbell-nzl/fid/internal/journal"
	"github.com/mbell-nzl/fid/internal/policy"
	"github.com/mbell-nzl/fid/internal/workflow"
)

type TransitionOptions struct {
	WorkspaceRoot string
	MigrationsDir string
	ID            string
	ToState       string
	ActorRole     string
}

type TransitionResult struct {
	ID        string
	FromState string
	ToState   string
}

type InvalidTransitionError struct {
	IssueID        string
	CurrentState   string
	RequestedState string
}

func (err *InvalidTransitionError) Error() string {
	if err == nil {
		return "<nil>"
	}
	return fmt.Sprintf(
		"transition issue %q from %q to %q is not allowed",
		err.IssueID,
		err.CurrentState,
		err.RequestedState,
	)
}

type BlockedTransitionError struct {
	IssueID        string
	CurrentState   string
	RequestedState string
	Unmet          []workflow.UnmetCondition
}

func (err *BlockedTransitionError) Error() string {
	if err == nil {
		return "<nil>"
	}

	var builder strings.Builder
	fmt.Fprintf(
		&builder,
		"transition issue %q from %q to %q blocked by workflow rules:",
		err.IssueID,
		err.CurrentState,
		err.RequestedState,
	)
	for _, condition := range err.Unmet {
		builder.WriteString("\n  - ")
		builder.WriteString(condition.Rule)
		details := make([]string, 0, 2)
		if len(condition.EntityIDs) > 0 {
			details = append(details, "entities "+strings.Join(condition.EntityIDs, ", "))
		}
		if len(condition.MissingRoles) > 0 {
			details = append(details, "missing roles "+strings.Join(condition.MissingRoles, ", "))
		}
		if len(details) > 0 {
			builder.WriteString(": ")
			builder.WriteString(strings.Join(details, "; "))
		}
	}
	return builder.String()
}

type issueTransitionState struct {
	ID         string
	Kind       string
	WorkflowID string
	Status     string
	UpdatedAt  string
}

func Transition(options TransitionOptions) (TransitionResult, error) {
	if strings.TrimSpace(options.WorkspaceRoot) == "" {
		return TransitionResult{}, fmt.Errorf("transition issue: workspace root is required")
	}
	if strings.TrimSpace(options.MigrationsDir) == "" {
		return TransitionResult{}, fmt.Errorf("transition issue: migrations dir is required")
	}
	if strings.TrimSpace(options.ID) == "" {
		return TransitionResult{}, fmt.Errorf("transition issue: issue id is required")
	}
	if err := policy.RequireKnownRole(options.WorkspaceRoot, options.ActorRole); err != nil {
		return TransitionResult{}, err
	}

	toState := strings.TrimSpace(options.ToState)
	if toState == "" {
		return TransitionResult{}, fmt.Errorf("transition issue: target state is required")
	}

	if _, err := journal.EnsureProjectionFresh(options.WorkspaceRoot, options.MigrationsDir); err != nil {
		return TransitionResult{}, err
	}

	commandArgv := issueTransitionArgv(options.ID, toState, options.ActorRole)
	var result TransitionResult
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

		current, err := loadIssueTransitionState(conn, options.ID)
		if err != nil {
			return err
		}

		rules, allowed := workflow.RulesForEdge(current.WorkflowID, current.Status, toState)
		if !allowed {
			return &InvalidTransitionError{
				IssueID:        current.ID,
				CurrentState:   current.Status,
				RequestedState: toState,
			}
		}

		unmet, err := evaluateTransitionRules(conn, current, rules)
		if err != nil {
			return err
		}
		if len(unmet) > 0 {
			return &BlockedTransitionError{
				IssueID:        current.ID,
				CurrentState:   current.Status,
				RequestedState: toState,
				Unmet:          unmet,
			}
		}

		activityID, err := journal.GenerateID("ACT", func(candidate string) (bool, error) {
			return activityIDExists(conn, candidate)
		})
		if err != nil {
			return err
		}

		event, err := journal.PrepareEvent(lock, journal.Event{
			EventSchema: journal.CurrentEventSchema,
			EventType:   "issue.transitioned",
			OccurredAt:  nextIssueUpdateTime(current.UpdatedAt, time.Now().UTC()),
			Actor: journal.Actor{
				Role: options.ActorRole,
			},
			Command: journal.Command{
				Argv:       commandArgv,
				CWD:        ".",
				FIDVersion: fidVersion,
			},
		})
		if err != nil {
			return err
		}

		timestamp := event.OccurredAt.UTC().Format(time.RFC3339)
		event.Effects = []journal.Effect{
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
					"journal_event_id": event.EventID,
					"ref_type":         "issue",
					"ref_id":           current.ID,
					"event":            "status_changed",
					"details": map[string]any{
						"from": current.Status,
						"to":   toState,
					},
					"actor":       options.ActorRole,
					"occurred_at": timestamp,
				},
			},
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

		result = TransitionResult{
			ID:        current.ID,
			FromState: current.Status,
			ToState:   toState,
		}
		return nil
	})
	if err != nil {
		return TransitionResult{}, err
	}

	return result, nil
}

func loadIssueTransitionState(conn *sql.DB, issueID string) (issueTransitionState, error) {
	row := conn.QueryRow(`
		SELECT id, kind, workflow_id, status, updated_at
		FROM issues
		WHERE id = ?
	`, issueID)

	var (
		state     issueTransitionState
		updatedAt sql.NullString
	)
	if err := row.Scan(&state.ID, &state.Kind, &state.WorkflowID, &state.Status, &updatedAt); err != nil {
		if err == sql.ErrNoRows {
			return issueTransitionState{}, &NotFoundError{IssueID: issueID}
		}
		return issueTransitionState{}, fmt.Errorf("query issue %q: %w", issueID, err)
	}
	if !updatedAt.Valid || strings.TrimSpace(updatedAt.String) == "" {
		return issueTransitionState{}, fmt.Errorf("query issue %q: updated_at is required", issueID)
	}
	state.UpdatedAt = updatedAt.String
	return state, nil
}

func evaluateTransitionRules(conn *sql.DB, current issueTransitionState, rules []workflow.RuleReference) ([]workflow.UnmetCondition, error) {
	unmet := make([]workflow.UnmetCondition, 0)
	for _, rule := range rules {
		result, err := evaluateTransitionRule(conn, current, rule)
		if err != nil {
			return nil, err
		}
		unmet = append(unmet, result.Unmet...)
	}
	return unmet, nil
}

func evaluateTransitionRule(conn *sql.DB, current issueTransitionState, rule workflow.RuleReference) (workflow.RuleResult, error) {
	switch rule.Name {
	case workflow.RuleAllACsVerified:
		criteria, err := loadAcceptanceCriterionVerifications(conn, current.ID)
		if err != nil {
			return workflow.RuleResult{}, err
		}
		return workflow.AllACsVerified(rule.RequiredRoles, criteria), nil
	case workflow.RuleHasCompletionReport:
		status, err := loadCompletionReportStatus(conn, current.ID)
		if err != nil {
			return workflow.RuleResult{}, err
		}
		return workflow.HasCompletionReport(status), nil
	case workflow.RuleConcernsAddressed:
		status, err := loadConcernStatus(conn, current.ID)
		if err != nil {
			return workflow.RuleResult{}, err
		}
		return workflow.ConcernsAddressed(status), nil
	case workflow.RuleAllChildTasksDoneOrCanceled:
		children, err := loadChildTaskStatuses(conn, current.ID)
		if err != nil {
			return workflow.RuleResult{}, err
		}
		return workflow.AllChildTasksDoneOrCancelled(children), nil
	default:
		return workflow.RuleResult{}, fmt.Errorf("transition issue %q: unsupported workflow rule %q", current.ID, rule.Name)
	}
}

func loadAcceptanceCriterionVerifications(conn *sql.DB, issueID string) ([]workflow.AcceptanceCriterionVerification, error) {
	criteria, err := issueAcceptanceCriteria(conn, issueID, nil)
	if err != nil {
		return nil, err
	}

	result := make([]workflow.AcceptanceCriterionVerification, 0, len(criteria))
	for _, criterion := range criteria {
		passingRoles := make([]string, 0, len(criterion.LatestVerifications))
		for role, raw := range criterion.LatestVerifications {
			verification, ok := raw.(map[string]any)
			if !ok {
				continue
			}
			status, _ := verification["status"].(string)
			if strings.TrimSpace(status) == "pass" {
				passingRoles = append(passingRoles, role)
			}
		}
		sort.Strings(passingRoles)
		result = append(result, workflow.AcceptanceCriterionVerification{
			AcceptanceCriterionID: criterion.ID,
			PassingRoles:          passingRoles,
		})
	}
	return result, nil
}

func loadCompletionReportStatus(conn *sql.DB, issueID string) (workflow.CompletionReportStatus, error) {
	var count int
	if err := conn.QueryRow(`SELECT COUNT(*) FROM task_completions WHERE issue_id = ?`, issueID).Scan(&count); err != nil {
		return workflow.CompletionReportStatus{}, fmt.Errorf("query completion reports for %q: %w", issueID, err)
	}
	return workflow.CompletionReportStatus{
		IssueID: issueID,
		Exists:  count > 0,
	}, nil
}

func loadConcernStatus(conn *sql.DB, issueID string) (workflow.ConcernStatus, error) {
	rows, err := conn.Query(`
		SELECT id
		FROM concerns
		WHERE issue_id = ? AND status = ?
		ORDER BY id
	`, issueID, "open")
	if err != nil {
		return workflow.ConcernStatus{}, fmt.Errorf("query open concerns for %q: %w", issueID, err)
	}
	defer closeIssueRows(rows)

	openConcernIDs := []string{}
	for rows.Next() {
		var concernID string
		if err := rows.Scan(&concernID); err != nil {
			return workflow.ConcernStatus{}, fmt.Errorf("scan open concern for %q: %w", issueID, err)
		}
		openConcernIDs = append(openConcernIDs, concernID)
	}
	if err := rows.Err(); err != nil {
		return workflow.ConcernStatus{}, fmt.Errorf("iterate open concerns for %q: %w", issueID, err)
	}
	return workflow.ConcernStatus{OpenConcernIDs: openConcernIDs}, nil
}

func loadChildTaskStatuses(conn *sql.DB, issueID string) ([]workflow.ChildTaskStatus, error) {
	rows, err := conn.Query(`
		SELECT id, status
		FROM issues
		WHERE parent_id = ?
		ORDER BY id
	`, issueID)
	if err != nil {
		return nil, fmt.Errorf("query child tasks for %q: %w", issueID, err)
	}
	defer closeIssueRows(rows)

	children := []workflow.ChildTaskStatus{}
	for rows.Next() {
		var child workflow.ChildTaskStatus
		if err := rows.Scan(&child.IssueID, &child.State); err != nil {
			return nil, fmt.Errorf("scan child task for %q: %w", issueID, err)
		}
		children = append(children, child)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate child tasks for %q: %w", issueID, err)
	}
	return children, nil
}

func issueTransitionArgv(issueID, toState, actorRole string) []string {
	return []string{"fid", "issue", "transition", issueID, "--to", toState, "--by", actorRole}
}
