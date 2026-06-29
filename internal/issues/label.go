package issues

import (
	"database/sql"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/mbell-nzl/fid/internal/config"
	fiddb "github.com/mbell-nzl/fid/internal/db"
	"github.com/mbell-nzl/fid/internal/journal"
	"github.com/mbell-nzl/fid/internal/policy"
)

type AddLabelOptions struct {
	WorkspaceRoot string
	MigrationsDir string
	IssueID       string
	Label         string
	ActorRole     string
}

type AddLabelResult struct {
	IssueID string
	Label   string
}

type RemoveLabelOptions struct {
	WorkspaceRoot string
	MigrationsDir string
	IssueID       string
	Label         string
	ActorRole     string
}

type RemoveLabelResult struct {
	IssueID string
	Label   string
}

func AddLabel(options AddLabelOptions) (AddLabelResult, error) {
	request, err := prepareIssueLabelWrite(
		options.WorkspaceRoot,
		options.MigrationsDir,
		options.IssueID,
		options.Label,
		options.ActorRole,
		actionIssueLabelAdd,
		issueLabelAddArgv,
		"issue label add",
		issueLabelOperationSpec{
			eventType: "issue.label.added",
			effectOp:  "insert",
			direction: "added",
			eventName: "issue_label_added",
			validate:  assertIssueLabelAddable,
		},
	)
	if err != nil {
		return AddLabelResult{}, err
	}

	result, err := issueLabelWrite(request)
	if err != nil {
		return AddLabelResult{}, err
	}
	return AddLabelResult{IssueID: result.issueID, Label: result.label}, nil
}

func RemoveLabel(options RemoveLabelOptions) (RemoveLabelResult, error) {
	request, err := prepareIssueLabelWrite(
		options.WorkspaceRoot,
		options.MigrationsDir,
		options.IssueID,
		options.Label,
		options.ActorRole,
		actionIssueLabelRemove,
		issueLabelRemoveArgv,
		"issue label remove",
		issueLabelOperationSpec{
			eventType: "issue.label.removed",
			effectOp:  "delete",
			direction: "removed",
			eventName: "issue_label_removed",
			validate:  assertIssueLabelRemovable,
		},
	)
	if err != nil {
		return RemoveLabelResult{}, err
	}

	result, err := issueLabelWrite(request)
	if err != nil {
		return RemoveLabelResult{}, err
	}
	return RemoveLabelResult{IssueID: result.issueID, Label: result.label}, nil
}

type issueLabelOperationSpec struct {
	eventType string
	effectOp  string
	direction string
	eventName string
	validate  func(*sql.DB, string, string) error
}

type issueLabelWriteRequest struct {
	workspaceRoot string
	migrationsDir string
	issueID       string
	label         string
	actorRole     string
	commandArgv   []string
	spec          issueLabelOperationSpec
}

type issueLabelWriteResult struct {
	issueID string
	label   string
}

func prepareIssueLabelWrite(
	workspaceRoot, migrationsDir, issueID, label, actorRole, action string,
	argvBuilder func(string, string, string) []string,
	opName string,
	spec issueLabelOperationSpec,
) (issueLabelWriteRequest, error) {
	if strings.TrimSpace(workspaceRoot) == "" {
		return issueLabelWriteRequest{}, fmt.Errorf("%s: workspace root is required", opName)
	}
	if strings.TrimSpace(migrationsDir) == "" {
		return issueLabelWriteRequest{}, fmt.Errorf("%s: migrations dir is required", opName)
	}

	trimmedIssueID := strings.TrimSpace(issueID)
	if trimmedIssueID == "" {
		return issueLabelWriteRequest{}, fmt.Errorf("%s: issue id is required", opName)
	}
	trimmedLabel := strings.TrimSpace(label)
	if trimmedLabel == "" {
		return issueLabelWriteRequest{}, fmt.Errorf("%s: label is required", opName)
	}
	if err := policy.Authorize(workspaceRoot, actorRole, action); err != nil {
		return issueLabelWriteRequest{}, err
	}

	return issueLabelWriteRequest{
		workspaceRoot: workspaceRoot,
		migrationsDir: migrationsDir,
		issueID:       trimmedIssueID,
		label:         trimmedLabel,
		actorRole:     actorRole,
		commandArgv:   argvBuilder(trimmedIssueID, trimmedLabel, actorRole),
		spec:          spec,
	}, nil
}

func issueLabelWrite(request issueLabelWriteRequest) (issueLabelWriteResult, error) {
	var result issueLabelWriteResult
	err := journal.WithLock(request.workspaceRoot, journal.LockOwner{
		Command:   strings.Join(request.commandArgv, " "),
		ActorRole: request.actorRole,
	}, func(lock *journal.WriteLock) error {
		return issueLabelWriteLocked(lock, request, &result)
	})
	if err != nil {
		return issueLabelWriteResult{}, err
	}
	return result, nil
}

func issueLabelWriteLocked(lock *journal.WriteLock, request issueLabelWriteRequest, result *issueLabelWriteResult) error {
	if _, err := journal.EnsureProjectionFreshWithLock(lock, request.migrationsDir); err != nil {
		return err
	}

	conn, err := openIssueLabelDB(request.workspaceRoot)
	if err != nil {
		return err
	}
	defer closeIssueDB(conn)

	if err := request.spec.validate(conn, request.issueID, request.label); err != nil {
		return err
	}

	event, err := prepareIssueLabelEvent(lock, conn, request)
	if err != nil {
		return err
	}
	if err := appendAndProjectIssueLabelEvent(conn, request.workspaceRoot, lock, event); err != nil {
		return err
	}

	*result = issueLabelWriteResult{issueID: request.issueID, label: request.label}
	return nil
}

func openIssueLabelDB(workspaceRoot string) (*sql.DB, error) {
	return fiddb.Open(filepath.Join(workspaceRoot, config.WorkspaceDirName, config.DatabaseFileName))
}

func prepareIssueLabelEvent(lock *journal.WriteLock, conn *sql.DB, request issueLabelWriteRequest) (journal.Event, error) {
	activityID, err := journal.GenerateID("ACT", func(candidate string) (bool, error) {
		return activityIDExists(conn, candidate)
	})
	if err != nil {
		return journal.Event{}, err
	}

	event, err := journal.PrepareEvent(lock, journal.Event{
		EventSchema: journal.CurrentEventSchema,
		EventType:   request.spec.eventType,
		Actor: journal.Actor{
			Role: request.actorRole,
		},
		Command: journal.Command{
			Argv:       request.commandArgv,
			CWD:        ".",
			FIDVersion: fidVersion,
		},
	})
	if err != nil {
		return journal.Event{}, err
	}

	timestamp := event.OccurredAt.UTC().Format(time.RFC3339)
	event.Effects = issueLabelEffects(request, activityID, event.EventID, timestamp)
	return event, nil
}

func issueLabelEffects(request issueLabelWriteRequest, activityID, eventID, timestamp string) []journal.Effect {
	effects := []journal.Effect{
		{
			Op:    request.spec.effectOp,
			Table: "issue_labels",
			PK: map[string]any{
				"issue_id": request.issueID,
				"label":    request.label,
			},
		},
	}
	if request.spec.effectOp == "insert" {
		effects[0].Values = map[string]any{}
	}
	effects = append(effects, issueLabelActivityEffect(
		activityID,
		eventID,
		request.actorRole,
		timestamp,
		request.issueID,
		request.label,
		request.spec.direction,
		request.spec.eventName,
	))
	return effects
}

func appendAndProjectIssueLabelEvent(conn *sql.DB, workspaceRoot string, lock *journal.WriteLock, event journal.Event) error {
	if err := journal.DryRunProjectEvent(conn, event); err != nil {
		return err
	}

	appended, err := journal.AppendEvent(lock, event)
	if err != nil {
		return err
	}
	return projectIssueEvent(conn, workspaceRoot, appended)
}

func assertIssueLabelAddable(conn *sql.DB, issueID, label string) error {
	if _, found, err := issueKindForID(conn, issueID); err != nil {
		return err
	} else if !found {
		return &NotFoundError{IssueID: issueID}
	}

	if err := requireExistingLabels(conn, []string{label}); err != nil {
		return fmt.Errorf("%s", strings.Replace(err.Error(), "create issue:", "issue label add:", 1))
	}

	labels, err := issueLabels(conn, issueID)
	if err != nil {
		return err
	}
	for _, existing := range labels {
		if existing == label {
			return fmt.Errorf("issue label add: issue %q already has label %q", issueID, label)
		}
	}
	return nil
}

func assertIssueLabelRemovable(conn *sql.DB, issueID, label string) error {
	if _, found, err := issueKindForID(conn, issueID); err != nil {
		return err
	} else if !found {
		return &NotFoundError{IssueID: issueID}
	}

	labels, err := issueLabels(conn, issueID)
	if err != nil {
		return err
	}
	for _, existing := range labels {
		if existing == label {
			return nil
		}
	}
	return fmt.Errorf("issue label remove: issue %q does not have label %q", issueID, label)
}

func issueLabelAddArgv(issueID, label, actorRole string) []string {
	return []string{"fid", "issue", "label", "add", issueID, label, "--by", actorRole}
}

func issueLabelRemoveArgv(issueID, label, actorRole string) []string {
	return []string{"fid", "issue", "label", "remove", issueID, label, "--by", actorRole}
}

func issueLabelActivityEffect(activityID, eventID, actorRole, occurredAt, issueID, label, direction, eventName string) journal.Effect {
	return journal.Effect{
		Op:    "activity",
		Table: "activity_log",
		PK:    map[string]any{"id": activityID},
		Values: map[string]any{
			"journal_event_id": eventID,
			"ref_type":         "issue",
			"ref_id":           issueID,
			"event":            eventName,
			"details": map[string]any{
				"label":     label,
				"direction": direction,
			},
			"actor":       actorRole,
			"occurred_at": occurredAt,
		},
	}
}
