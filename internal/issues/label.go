package issues

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/mbell-nzl/fid/internal/config"
	fiddb "github.com/mbell-nzl/fid/internal/db"
	"github.com/mbell-nzl/fid/internal/journal"
	"github.com/mbell-nzl/fid/internal/validation"
)

const IssueLabelsLimit = 100

type AddLabelOptions struct {
	WorkspaceRoot  string
	MigrationsDir  string
	IssueID        string
	Label          string
	Actor          string
	IdempotencyKey string
	ActorRole      string
}

type AddLabelResult struct {
	IssueID        string
	Label          string
	IssueRevision  int
	EventID        string
	OccurredAt     string
	IdempotencyKey string
	Replayed       bool
}

type RemoveLabelOptions struct {
	WorkspaceRoot  string
	MigrationsDir  string
	IssueID        string
	Label          string
	Actor          string
	IdempotencyKey string
	ActorRole      string
}

type RemoveLabelResult struct {
	IssueID        string
	Label          string
	IssueRevision  int
	EventID        string
	OccurredAt     string
	IdempotencyKey string
	Replayed       bool
}

func AddLabel(options AddLabelOptions) (AddLabelResult, error) {
	options.IdempotencyKey = idempotencyKeyForOptions(options.IdempotencyKey, options.ActorRole)
	request, err := prepareIssueLabelWrite(
		options.WorkspaceRoot,
		options.MigrationsDir,
		options.IssueID,
		options.Label,
		options.Actor,
		options.ActorRole,
		options.IdempotencyKey,
		issueLabelAddArgv,
		"issue label add",
		issueLabelOperationSpec{
			eventType: "issue.label.added",
			effectOp:  "insert",
			direction: "added",
			eventName: "issue-label-added",
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
	return AddLabelResult{IssueID: result.issueID, Label: result.label, IssueRevision: result.revision, EventID: result.eventID, OccurredAt: result.occurredAt, IdempotencyKey: result.idempotencyKey, Replayed: result.replayed}, nil
}

func RemoveLabel(options RemoveLabelOptions) (RemoveLabelResult, error) {
	options.IdempotencyKey = idempotencyKeyForOptions(options.IdempotencyKey, options.ActorRole)
	request, err := prepareIssueLabelWrite(
		options.WorkspaceRoot,
		options.MigrationsDir,
		options.IssueID,
		options.Label,
		options.Actor,
		options.ActorRole,
		options.IdempotencyKey,
		issueLabelRemoveArgv,
		"issue label remove",
		issueLabelOperationSpec{
			eventType: "issue.label.removed",
			effectOp:  "delete",
			direction: "removed",
			eventName: "issue-label-removed",
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
	return RemoveLabelResult{IssueID: result.issueID, Label: result.label, IssueRevision: result.revision, EventID: result.eventID, OccurredAt: result.occurredAt, IdempotencyKey: result.idempotencyKey, Replayed: result.replayed}, nil
}

type issueLabelOperationSpec struct {
	eventType string
	effectOp  string
	direction string
	eventName string
	validate  func(*sql.DB, string, string) error
}

type issueLabelWriteRequest struct {
	workspaceRoot  string
	migrationsDir  string
	issueID        string
	label          string
	actor          standaloneActor
	idempotencyKey string
	operation      string
	commandArgv    []string
	spec           issueLabelOperationSpec
}

type issueLabelWriteResult struct {
	issueID        string
	label          string
	revision       int
	eventID        string
	occurredAt     string
	idempotencyKey string
	replayed       bool
}

func prepareIssueLabelWrite(
	workspaceRoot, migrationsDir, issueID, label, actorValue, legacyRole, idempotencyKey string,
	argvBuilder func(string, string, string, string) []string,
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
	if strings.TrimSpace(label) == "" {
		return issueLabelWriteRequest{}, fmt.Errorf("%s: label is required", opName)
	}
	trimmedLabel, err := validation.NormalizeText("issue label", label, 255)
	if err != nil {
		return issueLabelWriteRequest{}, err
	}
	actor, err := standaloneActorForOptions(actorValue, legacyRole)
	if err != nil {
		return issueLabelWriteRequest{}, err
	}

	return issueLabelWriteRequest{
		workspaceRoot:  workspaceRoot,
		migrationsDir:  migrationsDir,
		issueID:        trimmedIssueID,
		label:          trimmedLabel,
		actor:          actor,
		idempotencyKey: idempotencyKey,
		operation:      strings.ReplaceAll(opName, " ", "."),
		commandArgv:    argvBuilder(trimmedIssueID, trimmedLabel, actorValue, idempotencyKey),
		spec:           spec,
	}, nil
}

func issueLabelWrite(request issueLabelWriteRequest) (issueLabelWriteResult, error) {
	var result issueLabelWriteResult
	err := journal.WithLock(request.workspaceRoot, journal.LockOwner{
		Command:   strings.Join(request.commandArgv, " "),
		ActorRole: request.actor.id,
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
	projectConfig, err := config.ReadProjectConfig(filepath.Join(request.workspaceRoot, config.WorkspaceDirName, config.ConfigFileName))
	if err != nil {
		return err
	}
	digest, err := issueRequestDigest(request.operation, projectConfig.ProjectID, request.actor, map[string]any{"issue_id": request.issueID, "label": request.label})
	if err != nil {
		return err
	}
	idempotencyRequest := journal.IdempotencyRequest{OperationScope: request.operation + ":" + projectConfig.ProjectID, IdempotencyKey: request.idempotencyKey, RequestDigest: digest}
	if outcome, found, err := journal.ReplayIdempotent(conn, idempotencyRequest); err != nil {
		return err
	} else if found {
		decoded, err := issueLabelResultFromJSON(outcome.ResultJSON, outcome.EventID, true)
		if err != nil {
			return err
		}
		*result = decoded
		return nil
	}

	revision, err := issueLabelCurrentRevision(conn, request.issueID)
	if err != nil {
		return err
	}
	if err := request.spec.validate(conn, request.issueID, request.label); err != nil {
		return err
	}

	activityID, err := journal.GenerateID("ACT", func(candidate string) (bool, error) {
		return activityIDExists(conn, candidate)
	})
	if err != nil {
		return err
	}
	sequence, err := journal.NextAggregateSequence(filepath.Join(request.workspaceRoot, config.WorkspaceDirName, config.JournalDirName), journal.Aggregate{ID: request.issueID, Type: "issue"})
	if err != nil {
		return err
	}
	baseEffects, err := issueLabelCoreEffects(request, revision)
	if err != nil {
		return err
	}
	event, err := journal.PrepareEvent(lock, journal.Event{EventSchema: journal.EventSchemaV2, EventType: request.spec.eventType, OccurredAt: issueMutationNow(nil), ActorV2: request.actor.eventActor(), Aggregate: journal.Aggregate{ID: request.issueID, Sequence: sequence, Type: "issue"}, CommandName: request.operation, Effects: baseEffects})
	if err != nil {
		return err
	}
	timestamp := event.OccurredAt.UTC().Format(time.RFC3339)
	event.Effects, err = issueLabelEffects(request, activityID, event.EventID, timestamp, revision)
	if err != nil {
		return err
	}
	seed := issueLabelWriteResult{issueID: request.issueID, label: request.label, revision: revision + 1, occurredAt: event.OccurredAt.UTC().Format(time.RFC3339Nano), idempotencyKey: request.idempotencyKey}
	outcome, err := journal.CommitIdempotent(lock, conn, idempotencyRequest, event, func(journal.Event) (journal.IdempotencyMutation, error) {
		encoded, err := issueLabelResultJSON(seed)
		if err != nil {
			return journal.IdempotencyMutation{}, err
		}
		return journal.IdempotencyMutation{EntityID: request.issueID + ":" + request.label, ResultJSON: encoded}, nil
	})
	if err != nil {
		return err
	}
	decoded, err := issueLabelResultFromJSON(outcome.ResultJSON, outcome.EventID, outcome.Replayed)
	if err != nil {
		return err
	}
	*result = decoded
	return nil
}

func openIssueLabelDB(workspaceRoot string) (*sql.DB, error) {
	return fiddb.Open(filepath.Join(workspaceRoot, config.WorkspaceDirName, config.DatabaseFileName))
}

func issueLabelEffects(request issueLabelWriteRequest, activityID, eventID, timestamp string, revision int) ([]journal.Effect, error) {
	effects, err := issueLabelCoreEffects(request, revision)
	if err != nil {
		return nil, err
	}
	effects = append(effects, issueLabelActivityEffect(
		activityID,
		eventID,
		request.actor.id,
		timestamp,
		request.issueID,
		request.label,
		request.spec.direction,
		request.spec.eventName,
	))
	return effects, nil
}

func issueLabelCoreEffects(request issueLabelWriteRequest, revision int) ([]journal.Effect, error) {
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
	revisionEffects, err := BuildIssueRevisionEffects([]IssueRevision{{IssueID: request.issueID, Observed: revision}})
	if err != nil {
		return nil, err
	}
	effects = append(effects, revisionEffects...)
	return effects, nil
}

func assertIssueLabelAddable(conn *sql.DB, issueID, label string) error {
	if _, found, err := issueKindForID(conn, issueID); err != nil {
		return err
	} else if !found {
		return &NotFoundError{IssueID: issueID}
	}

	if err := issueLabelExists(conn, label); err != nil {
		return err
	}

	labels, err := issueLabels(conn, issueID)
	if err != nil {
		return err
	}
	for _, existing := range labels {
		if existing == label {
			return &DuplicateIssueLabelError{IssueID: issueID, Label: label}
		}
	}
	if len(labels) >= IssueLabelsLimit {
		return &ResultLimitExceededError{Collection: "issue_labels", Limit: IssueLabelsLimit, Observed: len(labels) + 1}
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
	return &LabelNotFoundError{Label: label}
}

func issueLabelCurrentRevision(conn *sql.DB, issueID string) (int, error) {
	var revision int
	err := conn.QueryRow(`SELECT revision FROM issues WHERE id = ?`, issueID).Scan(&revision)
	if err == sql.ErrNoRows {
		return 0, &NotFoundError{IssueID: issueID}
	}
	return revision, err
}

func issueLabelExists(conn *sql.DB, label string) error {
	var name string
	err := conn.QueryRow(`SELECT name FROM labels WHERE name = ?`, label).Scan(&name)
	if err == sql.ErrNoRows {
		return &LabelNotFoundError{Label: label}
	}
	return err
}

func issueLabelAddArgv(issueID, label, actor, key string) []string {
	return []string{"fid", "issue", "label", "add", issueID, label, "--actor", actor, "--idempotency-key", key}
}

func issueLabelRemoveArgv(issueID, label, actor, key string) []string {
	return []string{"fid", "issue", "label", "remove", issueID, label, "--actor", actor, "--idempotency-key", key}
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

func issueLabelResultJSON(result issueLabelWriteResult) (string, error) {
	payload, err := journal.MarshalCanonicalJSON(map[string]any{"event_occurred_at": result.occurredAt, "idempotency_key": result.idempotencyKey, "issue_id": result.issueID, "issue_revision": result.revision, "label": result.label, "replayed": false})
	return string(payload), err
}

func issueLabelResultFromJSON(raw, eventID string, replayed bool) (issueLabelWriteResult, error) {
	var value struct {
		OccurredAt     string `json:"event_occurred_at"`
		IdempotencyKey string `json:"idempotency_key"`
		IssueID        string `json:"issue_id"`
		IssueRevision  int    `json:"issue_revision"`
		Label          string `json:"label"`
	}
	if err := json.Unmarshal([]byte(raw), &value); err != nil {
		return issueLabelWriteResult{}, fmt.Errorf("decode stored issue label result: %w", err)
	}
	return issueLabelWriteResult{issueID: value.IssueID, label: value.Label, revision: value.IssueRevision, eventID: eventID, occurredAt: value.OccurredAt, idempotencyKey: value.IdempotencyKey, replayed: replayed}, nil
}
