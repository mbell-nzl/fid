package workflow

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/mbell-nzl/fid/internal/config"
	"github.com/mbell-nzl/fid/internal/journal"
)

// AssignOptions supplies an assignment request and its write context.
type AssignOptions struct {
	WorkspaceRoot  string
	MigrationsDir  string
	ActorRole      string
	Actor          string
	IdempotencyKey string
	Request        AssignmentRequest
}

// AssignResult reports the workflow and scope recorded by one assignment.
type AssignResult struct {
	Name           string
	Kind           string
	Project        bool
	EventID        string
	OccurredAt     time.Time
	IdempotencyKey string
	Replayed       bool
}

// Assign appends and projects one atomic workflow.assigned event.
func Assign(options AssignOptions) (AssignResult, error) {
	options.IdempotencyKey = workflowIdempotencyKey(options.IdempotencyKey, options.ActorRole)
	if strings.TrimSpace(options.WorkspaceRoot) == "" {
		return AssignResult{}, fmt.Errorf("assign workflow: workspace root is required")
	}
	if strings.TrimSpace(options.MigrationsDir) == "" {
		return AssignResult{}, fmt.Errorf("assign workflow: migrations dir is required")
	}
	actor, err := standaloneActorForOptions(options.Actor, options.ActorRole)
	if err != nil {
		return AssignResult{}, err
	}

	result := AssignResult{Name: options.Request.Name, Kind: options.Request.Kind, Project: options.Request.Project}
	err = journal.WithLock(options.WorkspaceRoot, journal.LockOwner{
		Command:   strings.Join(workflowAssignArgv(options), " "),
		ActorRole: actor.id,
	}, func(lock *journal.WriteLock) error {
		return assignWorkflowLocked(lock, options, actor, &result)
	})
	if err != nil {
		return AssignResult{}, err
	}
	return result, nil
}

// AuthorizeAssignment remains for legacy in-process callers; standalone
// attribution deliberately grants no workflow authority.
func AuthorizeAssignment(workspaceRoot, actorRole string) error {
	return nil
}

func assignWorkflowLocked(lock *journal.WriteLock, options AssignOptions, actor standaloneActor, result *AssignResult) error {
	if _, err := journal.EnsureProjectionFreshWithLock(lock, options.MigrationsDir); err != nil {
		return err
	}
	conn, err := openWorkflowDB(options.WorkspaceRoot)
	if err != nil {
		return err
	}
	defer closeWorkflowDB(conn)
	projectConfig, err := config.ReadProjectConfig(filepath.Join(options.WorkspaceRoot, config.WorkspaceDirName, config.ConfigFileName))
	if err != nil {
		return err
	}
	digest, err := workflowRequestDigest("workflow.assign", projectConfig.ProjectID, actor, options.Request)
	if err != nil {
		return err
	}
	idempotencyRequest := journal.IdempotencyRequest{OperationScope: "workflow.assign:" + projectConfig.ProjectID, IdempotencyKey: options.IdempotencyKey, RequestDigest: digest}
	if outcome, found, err := journal.ReplayIdempotent(conn, idempotencyRequest); err != nil {
		return err
	} else if found {
		decoded, err := assignResultFromJSON(outcome.ResultJSON, outcome.EventID, true)
		if err != nil {
			return err
		}
		*result = decoded
		return nil
	}
	if err := ValidateAssignmentRequest(conn, options.Request); err != nil {
		return err
	}
	replacing, err := assignmentExists(conn, options.Request)
	if err != nil {
		return err
	}
	aggregateID := options.Request.Name + ":project"
	if !options.Request.Project {
		aggregateID = options.Request.Name + ":kind:" + options.Request.Kind
	}
	sequence, err := journal.NextAggregateSequence(filepath.Join(options.WorkspaceRoot, config.WorkspaceDirName, config.JournalDirName), journal.Aggregate{ID: aggregateID, Type: "workflow-assignment"})
	if err != nil {
		return err
	}
	event := journal.Event{
		EventSchema: journal.EventSchemaV2,
		EventType:   "workflow.assigned",
		OccurredAt:  time.Now().UTC(),
		ActorV2:     actor.eventActor(),
		Aggregate:   journal.Aggregate{ID: aggregateID, Sequence: sequence, Type: "workflow-assignment"},
		CommandName: "workflow.assign",
		Effects:     assignmentEffects(options.Request, replacing, time.Now().UTC().Format(time.RFC3339)),
	}
	result.OccurredAt = event.OccurredAt.UTC()
	result.IdempotencyKey = options.IdempotencyKey
	resultJSON, err := assignResultJSON(*result)
	if err != nil {
		return err
	}
	outcome, err := journal.CommitIdempotent(lock, conn, idempotencyRequest, event, func(journal.Event) (journal.IdempotencyMutation, error) {
		return journal.IdempotencyMutation{EntityID: aggregateID, ResultJSON: resultJSON}, nil
	})
	if err != nil {
		return err
	}
	decoded, err := assignResultFromJSON(outcome.ResultJSON, outcome.EventID, outcome.Replayed)
	if err != nil {
		return err
	}
	*result = decoded
	return nil
}

func assignResultJSON(result AssignResult) (string, error) {
	payload, err := journal.MarshalCanonicalJSON(map[string]any{"event_occurred_at": result.OccurredAt.Format(time.RFC3339Nano), "idempotency_key": result.IdempotencyKey, "kind": result.Kind, "name": result.Name, "project": result.Project, "replayed": false})
	return string(payload), err
}

func assignResultFromJSON(raw, eventID string, replayed bool) (AssignResult, error) {
	var result struct {
		OccurredAt     string `json:"event_occurred_at"`
		IdempotencyKey string `json:"idempotency_key"`
		Kind           string `json:"kind"`
		Name           string `json:"name"`
		Project        bool   `json:"project"`
	}
	if err := json.Unmarshal([]byte(raw), &result); err != nil {
		return AssignResult{}, fmt.Errorf("decode stored workflow assign result: %w", err)
	}
	occurredAt, err := time.Parse(time.RFC3339Nano, result.OccurredAt)
	if err != nil {
		return AssignResult{}, fmt.Errorf("decode stored workflow assign occurred_at: %w", err)
	}
	return AssignResult{Name: result.Name, Kind: result.Kind, Project: result.Project, EventID: eventID, OccurredAt: occurredAt, IdempotencyKey: result.IdempotencyKey, Replayed: replayed}, nil
}

func assignmentExists(conn *sql.DB, request AssignmentRequest) (bool, error) {
	query := `SELECT 1 FROM workflow_assignments WHERE scope = 'project' LIMIT 1`
	args := []any(nil)
	if !request.Project {
		query = `SELECT 1 FROM workflow_assignments WHERE scope = 'kind' AND kind = ? LIMIT 1`
		args = []any{request.Kind}
	}
	var found int
	err := conn.QueryRow(query, args...).Scan(&found)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("read workflow assignment: %w", err)
	}
	return true, nil
}

func assignmentEffects(request AssignmentRequest, replacing bool, assignedAt string) []journal.Effect {
	pk := map[string]any{"scope": "project"}
	if !request.Project {
		pk = map[string]any{"scope": "kind", "kind": request.Kind}
	}
	effects := make([]journal.Effect, 0, 2)
	if replacing {
		effects = append(effects, journal.Effect{Op: "delete", Table: "workflow_assignments", PK: pk, Values: map[string]any{}})
	}
	effects = append(effects, journal.Effect{
		Op: "insert", Table: "workflow_assignments", PK: pk,
		Values: map[string]any{"workflow_id": request.Name, "assigned_at": assignedAt},
	})
	return effects
}

func workflowAssignArgv(options AssignOptions) []string {
	argv := []string{"fid", "workflow", "assign", "--name", options.Request.Name}
	if options.Request.Project {
		argv = append(argv, "--project-wide")
	} else {
		argv = append(argv, "--kind", options.Request.Kind)
	}
	return append(argv, "--actor", options.Actor, "--idempotency-key", options.IdempotencyKey)
}
