package issues

import (
	"database/sql"
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/mbell-nzl/fid/internal/config"
	fiddb "github.com/mbell-nzl/fid/internal/db"
	"github.com/mbell-nzl/fid/internal/journal"
	"github.com/mbell-nzl/fid/internal/workflow"
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

	kind := strings.TrimSpace(options.Kind)
	if !supportedIssueKind(kind) {
		return CreateResult{}, &UnsupportedKindError{Kind: kind}
	}
	actor, err := parseStandaloneActor(options.Actor)
	if err != nil {
		return CreateResult{}, err
	}
	parentID := trimmedOptionalString(options.ParentID)
	specMDPath := trimmedOptionalString(options.SpecMDPath)
	labels, err := normalizedLabels(options.Labels)
	if err != nil {
		return CreateResult{}, err
	}
	if err := validateCreateOptions(kind, parentID, options.Priority, options.Sequence, strings.TrimSpace(options.Scheduling)); err != nil {
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
	commandArgv := issueCreateArgv(kind, options.Title, options.Actor, options.IdempotencyKey, parentID, options.Priority, options.Sequence, specMDPath, labels, strings.TrimSpace(options.Scheduling))
	err = journal.WithLock(options.WorkspaceRoot, journal.LockOwner{
		Command:   strings.Join(commandArgv, " "),
		ActorRole: actor.id,
	}, func(lock *journal.WriteLock) error {
		dbPath := filepath.Join(options.WorkspaceRoot, config.WorkspaceDirName, config.DatabaseFileName)
		conn, err := fiddb.Open(dbPath)
		if err != nil {
			return err
		}
		defer closeIssueDB(conn)

		issueID, err := journal.GenerateID(projectConfig.IssuePrefix, func(candidate string) (bool, error) {
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

		activityID, err := nextIssueActivityID(conn)
		if err != nil {
			return err
		}

		resolvedWorkflow, err := workflow.ResolveIssueWorkflow(conn, kind)
		if err != nil {
			return err
		}
		defaultOwnerRole := defaultIssueOwnerRole(kind)
		timestamp := time.Now().UTC().Format(time.RFC3339)
		issueValues := map[string]any{
			"kind":        kind,
			"title":       options.Title,
			"workflow_id": resolvedWorkflow.ID,
			"status":      resolvedWorkflow.InitialState,
			"revision":    1,
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

		event, err := journal.PrepareEvent(lock, journal.Event{
			EventSchema: journal.EventSchemaV2,
			EventType:   "issue.created",
			ActorV2:     actor.eventActor(),
			Aggregate:   journal.Aggregate{ID: issueID, Sequence: 1, Type: "issue"},
			CommandName: "issue.create",
			Effects:     effects,
		})
		if err != nil {
			return err
		}

		timestamp = event.OccurredAt.UTC().Format(time.RFC3339)
		issueValues["created_at"] = timestamp
		issueValues["updated_at"] = timestamp
		effects = append(effects, issueActivityEffect(
			activityID,
			event.EventID,
			actor.id,
			timestamp,
			issueID,
			"issue_created",
			map[string]any{
				"kind":      kind,
				"title":     options.Title,
				"parent_id": parentIDValue(parentID),
			},
		))
		event.Effects = effects

		if err := journal.DryRunProjectEvent(conn, event); err != nil {
			return err
		}

		first := CreateResult{
			ID:             issueID,
			Kind:           kind,
			Title:          options.Title,
			Status:         resolvedWorkflow.InitialState,
			Revision:       1,
			EventID:        event.EventID,
			OccurredAt:     event.OccurredAt.UTC().Format(time.RFC3339Nano),
			IdempotencyKey: options.IdempotencyKey,
		}
		resultJSON, err := createResultJSON(first)
		if err != nil {
			return err
		}
		requestDigest, err := issueRequestDigest("issue.create", projectConfig.ProjectID, actor, map[string]any{"kind": kind, "labels": labels, "parent_id": parentIDValue(parentID), "priority": options.Priority, "scheduling": options.Scheduling, "sequence": options.Sequence, "spec_md_path": specMDPath, "title": options.Title})
		if err != nil {
			return err
		}
		outcome, err := journal.CommitIdempotent(lock, conn, journal.IdempotencyRequest{OperationScope: "issue.create:" + projectConfig.ProjectID, IdempotencyKey: options.IdempotencyKey, RequestDigest: requestDigest}, event, func(journal.Event) (journal.IdempotencyMutation, error) {
			return journal.IdempotencyMutation{EntityID: issueID, ResultJSON: resultJSON}, nil
		})
		if err != nil {
			return err
		}
		result, err = createResultFromJSON(outcome.ResultJSON, outcome.EventID, outcome.Replayed)
		return err
	})
	if err != nil {
		return CreateResult{}, err
	}

	return result, nil
}

func parentIDValue(parentID *string) any {
	if parentID == nil {
		return nil
	}
	return *parentID
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

func defaultIssueOwnerRole(kind string) string {
	if kind == "epic" {
		return "tl"
	}
	return "developer"
}

func validateCreateOptions(kind string, parentID *string, priority, sequence *int, scheduling string) error {
	switch kind {
	case "epic":
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
	if scheduling != "" && scheduling != "sequential" && scheduling != "prioritised" {
		return fmt.Errorf("create issue: scheduling must be sequential or prioritised")
	}
	return nil
}

func validateParentConstraint(conn *sql.DB, kind string, parentID *string) error {
	if parentID == nil {
		if kind == "subtask" {
			return hierarchyInvalidError(kind, nil, nil)
		}
		return nil
	}

	parentKind, found, err := issueKindForID(conn, *parentID)
	if err != nil {
		return err
	}
	if !found {
		return &NotFoundError{IssueID: *parentID}
	}
	switch kind {
	case "epic":
		return hierarchyInvalidError(kind, parentID, &parentKind)
	case "task":
		if parentKind != "epic" {
			return hierarchyInvalidError(kind, parentID, &parentKind)
		}
	case "subtask":
		if parentKind != "task" {
			return hierarchyInvalidError(kind, parentID, &parentKind)
		}
	}
	return nil
}

func hierarchyInvalidError(kind string, parentID, parentKind *string) *HierarchyInvalidError {
	rule := ""
	switch kind {
	case "epic":
		rule = "epics must not have a parent"
	case "task":
		rule = "tasks must have an epic parent or no parent"
	case "subtask":
		rule = "subtasks must have a task parent"
	}
	return &HierarchyInvalidError{
		Kind:       kind,
		ParentID:   parentID,
		ParentKind: parentKind,
		Rule:       rule,
	}
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

func issueCreateArgv(kind, title, actor, idempotencyKey string, parentID *string, priority, sequence *int, specMDPath *string, labels []string, scheduling string) []string {
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
	for _, label := range labels {
		argv = append(argv, "--label", label)
	}
	if scheduling != "" {
		argv = append(argv, "--scheduling", scheduling)
	}
	argv = append(argv, "--actor", actor, "--idempotency-key", idempotencyKey)
	return argv
}
