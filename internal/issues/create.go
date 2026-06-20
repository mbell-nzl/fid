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
