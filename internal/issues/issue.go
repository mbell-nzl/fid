package issues

import (
	"database/sql"
	"errors"
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

const fidVersion = "0.1.0"

type CreateOptions struct {
	WorkspaceRoot string
	MigrationsDir string
	Kind          string
	Title         string
	ActorRole     string
	ParentID      *string
	Priority      *int
	Sequence      *int
	SpecMDPath    *string
	EstTokens     *int
	DueDate       *string
	Labels        []string
	InInbox       bool
	OwnerRole     *string
	OwnerActor    *string
	Scheduling    string
}

type CreateResult struct {
	ID    string
	Kind  string
	Title string
}

type GetOptions struct {
	WorkspaceRoot string
	MigrationsDir string
	ID            string
}

type UnsupportedKindError struct {
	Kind string
}

func (err *UnsupportedKindError) Error() string {
	if err == nil {
		return "<nil>"
	}
	return fmt.Sprintf("kind %q must be epic, task, or subtask", err.Kind)
}

type NotFoundError struct {
	IssueID string
}

func (err *NotFoundError) Error() string {
	if err == nil {
		return "<nil>"
	}
	return fmt.Sprintf("issue %q not found", err.IssueID)
}

type Issue struct {
	ID                 string                `json:"id"`
	Kind               string                `json:"kind"`
	Title              string                `json:"title"`
	Status             string                `json:"status"`
	StatusCategory     string                `json:"status_category"`
	WorkflowID         string                `json:"workflow_id"`
	ParentID           *string               `json:"parent_id"`
	Priority           *int                  `json:"priority"`
	Sequence           *int                  `json:"sequence"`
	SpecMDPath         *string               `json:"spec_md_path"`
	EstTokens          *int                  `json:"est_tokens"`
	OwnerRole          *string               `json:"owner_role"`
	OwnerActor         *string               `json:"owner_actor"`
	Labels             []string              `json:"labels"`
	InInbox            bool                  `json:"in_inbox"`
	ConcernsRaised     bool                  `json:"concerns_raised"`
	OpenConcernsCount  int                   `json:"open_concerns_count"`
	TriageStatus       *string               `json:"triage_status"`
	AcceptanceCriteria []AcceptanceCriterion `json:"acceptance_criteria"`
	CreatedAt          *string               `json:"created_at"`
	UpdatedAt          *string               `json:"updated_at"`
}

type AcceptanceCriterion struct {
	ID                  string         `json:"id"`
	Ordinal             int            `json:"ordinal"`
	Text                string         `json:"text"`
	LatestVerifications map[string]any `json:"latest_verifications"`
}

func Create(options CreateOptions) (CreateResult, error) {
	if strings.TrimSpace(options.WorkspaceRoot) == "" {
		return CreateResult{}, fmt.Errorf("create issue: workspace root is required")
	}
	if strings.TrimSpace(options.MigrationsDir) == "" {
		return CreateResult{}, fmt.Errorf("create issue: migrations dir is required")
	}
	if strings.TrimSpace(options.Title) == "" {
		return CreateResult{}, fmt.Errorf("create issue: title is required")
	}
	if err := policy.RequireKnownRole(options.WorkspaceRoot, options.ActorRole); err != nil {
		return CreateResult{}, err
	}

	kind := strings.TrimSpace(options.Kind)
	if !supportedIssueKind(kind) {
		return CreateResult{}, &UnsupportedKindError{Kind: kind}
	}
	parentID := trimmedOptionalString(options.ParentID)
	specMDPath := trimmedOptionalString(options.SpecMDPath)
	dueDate := trimmedOptionalString(options.DueDate)
	ownerRole := trimmedOptionalString(options.OwnerRole)
	ownerActor := trimmedOptionalString(options.OwnerActor)
	labels, err := normalizedLabels(options.Labels)
	if err != nil {
		return CreateResult{}, err
	}
	if err := validateCreateOptions(kind, parentID, options.Priority, options.Sequence, options.EstTokens, dueDate, ownerRole, strings.TrimSpace(options.Scheduling), options.WorkspaceRoot); err != nil {
		return CreateResult{}, err
	}

	if _, err := journal.EnsureProjectionFresh(options.WorkspaceRoot, options.MigrationsDir); err != nil {
		return CreateResult{}, err
	}

	projectConfig, err := config.ReadProjectConfig(filepath.Join(options.WorkspaceRoot, config.WorkspaceDirName, config.ConfigFileName))
	if err != nil {
		return CreateResult{}, err
	}

	var result CreateResult
	commandArgv := issueCreateArgv(kind, options.Title, options.ActorRole, parentID, options.Priority, options.Sequence, specMDPath, options.EstTokens, dueDate, labels, ownerRole, ownerActor, strings.TrimSpace(options.Scheduling), options.InInbox)
	err = journal.WithLock(options.WorkspaceRoot, journal.LockOwner{
		Command:   strings.Join(commandArgv, " "),
		ActorRole: options.ActorRole,
	}, func(lock *journal.WriteLock) error {
		dbPath := filepath.Join(options.WorkspaceRoot, config.WorkspaceDirName, config.DatabaseFileName)
		conn, err := fiddb.Open(dbPath)
		if err != nil {
			return err
		}
		defer closeIssueDB(conn)

		issueID, err := journal.GenerateID(projectConfig.Prefix, func(candidate string) (bool, error) {
			return issueIDExists(conn, candidate)
		})
		if err != nil {
			return err
		}

		if err := validateParentConstraint(conn, kind, parentID); err != nil {
			return err
		}
		if err := requireExistingLabels(conn, labels); err != nil {
			return err
		}

		workflowID := defaultWorkflowID(kind)
		defaultOwnerRole := defaultIssueOwnerRole(kind)
		if ownerRole != nil {
			defaultOwnerRole = *ownerRole
		}
		timestamp := time.Now().UTC().Format(time.RFC3339)
		issueValues := map[string]any{
			"kind":        kind,
			"title":       options.Title,
			"workflow_id": workflowID,
			"status":      "todo",
			"owner_role":  defaultOwnerRole,
			"created_at":  timestamp,
			"updated_at":  timestamp,
		}
		if parentID != nil {
			issueValues["parent_id"] = *parentID
		}
		if options.Priority != nil {
			issueValues["priority"] = *options.Priority
		}
		if options.Sequence != nil {
			issueValues["sequence"] = *options.Sequence
		}
		if specMDPath != nil {
			issueValues["spec_md_path"] = *specMDPath
		}
		if options.EstTokens != nil {
			issueValues["est_tokens"] = *options.EstTokens
		}
		if dueDate != nil {
			issueValues["due_date"] = *dueDate
		}
		if options.InInbox {
			issueValues["in_inbox"] = true
		}
		if ownerActor != nil {
			issueValues["owner_actor"] = *ownerActor
		}

		effects := []journal.Effect{
			{
				Op:    "insert",
				Table: "issues",
				PK: map[string]any{
					"id": issueID,
				},
				Values: issueValues,
			},
		}
		if kind == "epic" {
			effects = append(effects, journal.Effect{
				Op:    "insert",
				Table: "epic_metadata",
				PK: map[string]any{
					"issue_id": issueID,
				},
				Values: map[string]any{
					"scheduling": strings.TrimSpace(options.Scheduling),
				},
			})
		}
		for _, label := range labels {
			effects = append(effects, journal.Effect{
				Op:    "insert",
				Table: "issue_labels",
				PK: map[string]any{
					"issue_id": issueID,
					"label":    label,
				},
				Values: map[string]any{},
			})
		}

		event := journal.Event{
			EventSchema: journal.CurrentEventSchema,
			EventType:   "issue.created",
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

		result = CreateResult{
			ID:    issueID,
			Kind:  kind,
			Title: options.Title,
		}
		return nil
	})
	if err != nil {
		return CreateResult{}, err
	}

	return result, nil
}

func Get(options GetOptions) (Issue, error) {
	if strings.TrimSpace(options.WorkspaceRoot) == "" {
		return Issue{}, fmt.Errorf("get issue: workspace root is required")
	}
	if strings.TrimSpace(options.MigrationsDir) == "" {
		return Issue{}, fmt.Errorf("get issue: migrations dir is required")
	}
	if strings.TrimSpace(options.ID) == "" {
		return Issue{}, fmt.Errorf("get issue: issue id is required")
	}

	if _, err := journal.EnsureProjectionFresh(options.WorkspaceRoot, options.MigrationsDir); err != nil {
		return Issue{}, err
	}

	conn, err := fiddb.Open(filepath.Join(options.WorkspaceRoot, config.WorkspaceDirName, config.DatabaseFileName))
	if err != nil {
		return Issue{}, err
	}
	defer closeIssueDB(conn)

	row := conn.QueryRow(`
		SELECT id, kind, title, status, workflow_id, parent_id, priority, sequence,
		       spec_md_path, est_tokens, owner_role, owner_actor, in_inbox,
		       concerns_raised, open_concerns_count, triage_status, created_at, updated_at
		FROM issues_enriched
		WHERE id = ?
	`, options.ID)

	var (
		record                                         Issue
		parentID, specMDPath, ownerRole                sql.NullString
		ownerActor, triageStatus, createdAt, updatedAt sql.NullString
		priority, sequence, estTokens                  sql.NullInt64
	)
	if err := row.Scan(
		&record.ID,
		&record.Kind,
		&record.Title,
		&record.Status,
		&record.WorkflowID,
		&parentID,
		&priority,
		&sequence,
		&specMDPath,
		&estTokens,
		&ownerRole,
		&ownerActor,
		&record.InInbox,
		&record.ConcernsRaised,
		&record.OpenConcernsCount,
		&triageStatus,
		&createdAt,
		&updatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Issue{}, &NotFoundError{IssueID: options.ID}
		}
		return Issue{}, fmt.Errorf("query issue %q: %w", options.ID, err)
	}

	record.StatusCategory = statusCategory(record.Status)
	record.Labels, err = issueLabels(conn, record.ID)
	if err != nil {
		return Issue{}, err
	}
	record.AcceptanceCriteria = []AcceptanceCriterion{}
	record.ParentID = nullStringPointer(parentID)
	record.SpecMDPath = nullStringPointer(specMDPath)
	record.OwnerRole = nullStringPointer(ownerRole)
	record.OwnerActor = nullStringPointer(ownerActor)
	record.TriageStatus = nullStringPointer(triageStatus)
	record.CreatedAt = nullStringPointer(createdAt)
	record.UpdatedAt = nullStringPointer(updatedAt)
	record.Priority = nullIntPointer(priority)
	record.Sequence = nullIntPointer(sequence)
	record.EstTokens = nullIntPointer(estTokens)

	return record, nil
}

func issueIDExists(conn *sql.DB, issueID string) (bool, error) {
	var existing string
	err := conn.QueryRow(`SELECT id FROM issues WHERE id = ? LIMIT 1`, issueID).Scan(&existing)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	return false, fmt.Errorf("query issue id %q: %w", issueID, err)
}

func countJournalLines(path string) int {
	payload, err := os.ReadFile(path)
	if err != nil {
		return 1
	}
	trimmed := strings.TrimSpace(string(payload))
	if trimmed == "" {
		return 0
	}
	return len(strings.Split(trimmed, "\n"))
}

func supportedIssueKind(kind string) bool {
	switch kind {
	case "epic", "task", "subtask":
		return true
	default:
		return false
	}
}

func defaultWorkflowID(kind string) string {
	switch kind {
	case "epic":
		return "epic-default"
	case "subtask":
		return "subtask-default"
	default:
		return "task-default"
	}
}

func defaultIssueOwnerRole(kind string) string {
	if kind == "epic" {
		return "tl"
	}
	return "developer"
}

func validateCreateOptions(kind string, parentID *string, priority, sequence, estTokens *int, dueDate, ownerRole *string, scheduling, workspaceRoot string) error {
	switch kind {
	case "epic":
		if parentID != nil {
			return fmt.Errorf("issues: epic parent_id must be NULL")
		}
		if scheduling == "" {
			return fmt.Errorf("create issue: epic scheduling is required")
		}
	case "task":
		if scheduling != "" {
			return fmt.Errorf("create issue: scheduling is only valid for epics")
		}
	case "subtask":
		if scheduling != "" {
			return fmt.Errorf("create issue: scheduling is only valid for epics")
		}
	}

	if priority != nil && (*priority < 0 || *priority > 5) {
		return fmt.Errorf("create issue: priority must be between 0 and 5")
	}
	if sequence != nil && *sequence <= 0 {
		return fmt.Errorf("create issue: sequence must be positive")
	}
	if kind == "epic" && estTokens != nil {
		return fmt.Errorf("create issue: est_tokens is only valid for tasks and subtasks")
	}
	if dueDate != nil {
		if _, err := time.Parse("2006-01-02", *dueDate); err != nil {
			return fmt.Errorf("create issue: due_date must use YYYY-MM-DD")
		}
	}
	if ownerRole != nil {
		if err := policy.RequireKnownRole(workspaceRoot, *ownerRole); err != nil {
			return err
		}
	}
	if scheduling != "" && scheduling != "sequential" && scheduling != "prioritised" {
		return fmt.Errorf("create issue: scheduling must be sequential or prioritised")
	}
	return nil
}

func validateParentConstraint(conn *sql.DB, kind string, parentID *string) error {
	if parentID == nil {
		if kind == "subtask" {
			return fmt.Errorf("issues: subtask parent_id must reference a task")
		}
		return nil
	}

	parentKind, found, err := issueKindForID(conn, *parentID)
	if err != nil {
		return err
	}
	switch kind {
	case "epic":
		return fmt.Errorf("issues: epic parent_id must be NULL")
	case "task":
		if !found || parentKind != "epic" {
			return fmt.Errorf("issues: task parent_id must reference an epic or be NULL")
		}
	case "subtask":
		if !found || parentKind != "task" {
			return fmt.Errorf("issues: subtask parent_id must reference a task")
		}
	}
	return nil
}

func issueKindForID(conn *sql.DB, issueID string) (string, bool, error) {
	var kind string
	err := conn.QueryRow(`SELECT kind FROM issues WHERE id = ? LIMIT 1`, issueID).Scan(&kind)
	if err == nil {
		return kind, true, nil
	}
	if errors.Is(err, sql.ErrNoRows) {
		return "", false, nil
	}
	return "", false, fmt.Errorf("query parent issue %q: %w", issueID, err)
}

func requireExistingLabels(conn *sql.DB, labels []string) error {
	for _, label := range labels {
		var existing string
		err := conn.QueryRow(`SELECT name FROM labels WHERE name = ? LIMIT 1`, label).Scan(&existing)
		if err == nil {
			continue
		}
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("create issue: label %q does not exist", label)
		}
		return fmt.Errorf("query label %q: %w", label, err)
	}
	return nil
}

func issueLabels(conn *sql.DB, issueID string) ([]string, error) {
	rows, err := conn.Query(`SELECT label FROM issue_labels WHERE issue_id = ? ORDER BY label`, issueID)
	if err != nil {
		return nil, fmt.Errorf("query issue labels for %q: %w", issueID, err)
	}
	defer func() {
		_ = rows.Close()
	}()

	labels := []string{}
	for rows.Next() {
		var label string
		if err := rows.Scan(&label); err != nil {
			return nil, fmt.Errorf("scan issue label for %q: %w", issueID, err)
		}
		labels = append(labels, label)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate issue labels for %q: %w", issueID, err)
	}
	return labels, nil
}

func normalizedLabels(input []string) ([]string, error) {
	if len(input) == 0 {
		return nil, nil
	}

	seen := make(map[string]bool, len(input))
	labels := make([]string, 0, len(input))
	for _, raw := range input {
		label := strings.TrimSpace(raw)
		if label == "" {
			return nil, fmt.Errorf("create issue: label must not be blank")
		}
		if seen[label] {
			continue
		}
		seen[label] = true
		labels = append(labels, label)
	}
	return labels, nil
}

func trimmedOptionalString(value *string) *string {
	if value == nil {
		return nil
	}
	trimmed := strings.TrimSpace(*value)
	if trimmed == "" {
		return nil
	}
	return &trimmed
}

func issueCreateArgv(kind, title, actorRole string, parentID *string, priority, sequence *int, specMDPath *string, estTokens *int, dueDate *string, labels []string, ownerRole, ownerActor *string, scheduling string, inInbox bool) []string {
	argv := []string{"fid", "issue", "create", "--kind", kind, "--title", title}
	switch kind {
	case "task":
		if parentID != nil {
			argv = append(argv, "--epic", *parentID)
		}
	case "subtask":
		if parentID != nil {
			argv = append(argv, "--parent", *parentID)
		}
	}
	if priority != nil {
		argv = append(argv, "--priority", fmt.Sprintf("%d", *priority))
	}
	if sequence != nil {
		argv = append(argv, "--sequence", fmt.Sprintf("%d", *sequence))
	}
	if specMDPath != nil {
		argv = append(argv, "--spec-md", *specMDPath)
	}
	if estTokens != nil {
		argv = append(argv, "--est-tokens", fmt.Sprintf("%d", *estTokens))
	}
	if dueDate != nil {
		argv = append(argv, "--due-date", *dueDate)
	}
	for _, label := range labels {
		argv = append(argv, "--label", label)
	}
	if ownerRole != nil {
		argv = append(argv, "--owner-role", *ownerRole)
	}
	if ownerActor != nil {
		argv = append(argv, "--owner-actor", *ownerActor)
	}
	if scheduling != "" {
		argv = append(argv, "--scheduling", scheduling)
	}
	if inInbox {
		argv = append(argv, "--inbox")
	}
	argv = append(argv, "--by", actorRole)
	return argv
}

func statusCategory(status string) string {
	switch status {
	case "todo":
		return "unstarted"
	case "in_progress", "dev_complete", "qa_review", "qa_block", "failed_gate", "wrapped":
		return "started"
	case "qa_pass", "done", "closed":
		return "completed"
	case "cancelled":
		return "cancelled"
	default:
		return "unstarted"
	}
}

func nullStringPointer(value sql.NullString) *string {
	if !value.Valid {
		return nil
	}
	cloned := value.String
	return &cloned
}

func nullIntPointer(value sql.NullInt64) *int {
	if !value.Valid {
		return nil
	}
	cloned := int(value.Int64)
	return &cloned
}

func closeIssueDB(conn *sql.DB) {
	_ = conn.Close()
}
