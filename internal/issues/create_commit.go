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
	"github.com/mbell-nzl/fid/internal/workflow"
)

type createRequest struct {
	kind, scheduling string
	actor            standaloneActor
	parentID         *string
	specMDPath       *string
	labels           []string
}

func normalizeCreateRequest(options CreateOptions) (createRequest, error) {
	if strings.TrimSpace(options.WorkspaceRoot) == "" {
		return createRequest{}, fmt.Errorf("create issue: workspace root is required")
	}
	if strings.TrimSpace(options.MigrationsDir) == "" {
		return createRequest{}, fmt.Errorf("create issue: migrations dir is required")
	}
	if strings.TrimSpace(options.Title) == "" {
		return createRequest{}, fmt.Errorf("create issue: title is required")
	}
	request := createRequest{kind: strings.TrimSpace(options.Kind), scheduling: strings.TrimSpace(options.Scheduling), parentID: trimmedOptionalString(options.ParentID), specMDPath: trimmedOptionalString(options.SpecMDPath)}
	if !supportedIssueKind(request.kind) {
		return createRequest{}, &UnsupportedKindError{Kind: request.kind}
	}
	var err error
	if request.actor, err = standaloneActorForOptions(options.Actor); err != nil {
		return createRequest{}, err
	}
	if request.labels, err = normalizedLabels(options.Labels); err != nil {
		return createRequest{}, err
	}
	return request, validateCreateOptions(request.kind, request.parentID, options.Priority, options.Sequence, request.scheduling)
}

func ensureIssueProjection(options CreateOptions) (config.ProjectConfig, error) {
	if _, err := journal.EnsureProjectionFresh(options.WorkspaceRoot, options.MigrationsDir); err != nil {
		return config.ProjectConfig{}, err
	}
	return config.ReadProjectConfig(filepath.Join(options.WorkspaceRoot, config.WorkspaceDirName, config.ConfigFileName))
}

func commitCreate(options CreateOptions, request createRequest) (CreateResult, error) {
	projectConfig, err := ensureIssueProjection(options)
	if err != nil {
		return CreateResult{}, err
	}
	argv := issueCreateArgv(request.kind, options.Title, options.Actor, options.IdempotencyKey, request.parentID, options.Priority, options.Sequence, request.specMDPath, request.labels, request.scheduling)
	var result CreateResult
	err = journal.WithLock(options.WorkspaceRoot, journal.LockOwner{Command: strings.Join(argv, " ")}, func(lock *journal.WriteLock) error {
		conn, err := fiddb.Open(filepath.Join(options.WorkspaceRoot, config.WorkspaceDirName, config.DatabaseFileName))
		if err != nil {
			return err
		}
		defer closeIssueDB(conn)
		result, err = createUnderLock(conn, lock, options, request, projectConfig)
		return err
	})
	return result, err
}

func createUnderLock(conn *sql.DB, lock *journal.WriteLock, options CreateOptions, request createRequest, project config.ProjectConfig) (CreateResult, error) {
	issueID, activityID, resolved, err := createIdentifiers(conn, project.IssuePrefix, request)
	if err != nil {
		return CreateResult{}, err
	}
	values, effects := createIssueEffects(options, request, issueID, resolved)
	event, err := prepareCreateEvent(lock, options, request, issueID, activityID, values, effects)
	if err != nil {
		return CreateResult{}, err
	}
	result := CreateResult{ID: issueID, Kind: request.kind, Title: options.Title, Status: resolved.InitialState, Revision: 1, OccurredAt: event.OccurredAt.UTC().Format(time.RFC3339Nano), IdempotencyKey: options.IdempotencyKey}
	return commitCreateEvent(conn, lock, project.ProjectID, request, options, event, result)
}

func createIdentifiers(conn *sql.DB, prefix string, request createRequest) (string, string, workflow.ResolvedIssueWorkflow, error) {
	issueID, err := journal.GenerateID(prefix, func(candidate string) (bool, error) { return issueIDExists(conn, candidate) })
	if err != nil {
		return "", "", workflow.ResolvedIssueWorkflow{}, err
	}
	if err := validateParentConstraint(conn, request.kind, request.parentID); err != nil {
		return "", "", workflow.ResolvedIssueWorkflow{}, err
	}
	if err := requireExistingLabels(conn, request.labels); err != nil {
		return "", "", workflow.ResolvedIssueWorkflow{}, err
	}
	activityID, err := nextIssueActivityID(conn)
	if err != nil {
		return "", "", workflow.ResolvedIssueWorkflow{}, err
	}
	resolved, err := workflow.ResolveIssueWorkflow(conn, request.kind)
	return issueID, activityID, resolved, err
}

func createIssueEffects(options CreateOptions, request createRequest, issueID string, resolved workflow.ResolvedIssueWorkflow) (map[string]any, []journal.Effect) {
	timestamp := issueMutationNow(options.now).Format(time.RFC3339)
	values := map[string]any{"kind": request.kind, "title": options.Title, "workflow_id": resolved.ID, "status": resolved.InitialState, "revision": 1, "owner_role": defaultIssueOwnerRole(request.kind), "created_at": timestamp, "updated_at": timestamp}
	if request.parentID != nil {
		values["parent_id"] = *request.parentID
	}
	if options.Priority != nil {
		values["priority"] = *options.Priority
	}
	if options.Sequence != nil {
		values["sequence"] = *options.Sequence
	}
	if request.specMDPath != nil {
		values["spec_md_path"] = *request.specMDPath
	}
	effects := []journal.Effect{{Op: "insert", Table: "issues", PK: map[string]any{"id": issueID}, Values: values}}
	if request.kind == "epic" {
		effects = append(effects, journal.Effect{Op: "insert", Table: "epic_metadata", PK: map[string]any{"issue_id": issueID}, Values: map[string]any{"scheduling": request.scheduling}})
	}
	for _, label := range request.labels {
		effects = append(effects, journal.Effect{Op: "insert", Table: "issue_labels", PK: map[string]any{"issue_id": issueID, "label": label}, Values: map[string]any{}})
	}
	return values, effects
}

func prepareCreateEvent(lock *journal.WriteLock, options CreateOptions, request createRequest, issueID, activityID string, values map[string]any, effects []journal.Effect) (journal.Event, error) {
	event := journal.Event{EventSchema: journal.EventSchemaV2, EventType: "issue.created", OccurredAt: issueMutationNow(options.now), ActorV2: request.actor.eventActor(), Aggregate: journal.Aggregate{ID: issueID, Sequence: 1, Type: "issue"}, CommandName: "issue.create", Effects: effects}
	prepared, err := journal.PrepareEvent(lock, event)
	if err != nil {
		return journal.Event{}, err
	}
	timestamp := prepared.OccurredAt.UTC().Format(time.RFC3339)
	values["created_at"], values["updated_at"] = timestamp, timestamp
	event.Effects = append(effects, issueActivityEffect(activityID, prepared.EventID, request.actor.id, timestamp, issueID, "issue-created", map[string]any{"kind": request.kind, "title": options.Title, "parent_id": parentIDValue(request.parentID)}))
	return event, nil
}

func commitCreateEvent(conn *sql.DB, lock *journal.WriteLock, projectID string, request createRequest, options CreateOptions, event journal.Event, result CreateResult) (CreateResult, error) {
	resultJSON, err := createResultJSON(result)
	if err != nil {
		return CreateResult{}, err
	}
	digest, err := issueRequestDigest("issue.create", projectID, request.actor, map[string]any{"kind": request.kind, "labels": request.labels, "parent_id": parentIDValue(request.parentID), "priority": options.Priority, "scheduling": request.scheduling, "sequence": options.Sequence, "spec_md_path": request.specMDPath, "title": options.Title})
	if err != nil {
		return CreateResult{}, err
	}
	outcome, err := journal.CommitIdempotent(lock, conn, journal.IdempotencyRequest{OperationScope: "issue.create:" + projectID, IdempotencyKey: options.IdempotencyKey, RequestDigest: digest}, event, func(journal.Event) (journal.IdempotencyMutation, error) {
		return journal.IdempotencyMutation{EntityID: result.ID, ResultJSON: resultJSON}, nil
	})
	if err != nil {
		return CreateResult{}, err
	}
	return createResultFromJSON(outcome.ResultJSON, outcome.EventID, outcome.Replayed)
}
