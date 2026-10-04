package workflow

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
)

// DefineOptions supplies a validated workflow declaration and its write context.
type DefineOptions struct {
	WorkspaceRoot  string
	MigrationsDir  string
	Actor          string
	IdempotencyKey string
	Request        DefinitionRequest
	RequestInput   map[string]any
}

// DefineResult reports the projected rows written for one workflow declaration.
type DefineResult struct {
	Name           string
	AppliesTo      string
	States         int
	Transitions    int
	Requirements   int
	EventID        string
	OccurredAt     time.Time
	IdempotencyKey string
	Replayed       bool
}

// Define appends and projects one atomic workflow.defined event.
func Define(options DefineOptions) (DefineResult, error) {
	if strings.TrimSpace(options.WorkspaceRoot) == "" {
		return DefineResult{}, fmt.Errorf("define workflow: workspace root is required")
	}
	if strings.TrimSpace(options.MigrationsDir) == "" {
		return DefineResult{}, fmt.Errorf("define workflow: migrations dir is required")
	}
	actor, err := standaloneActorForOptions(options.Actor)
	if err != nil {
		return DefineResult{}, err
	}

	result := DefineResult{Name: options.Request.Name, AppliesTo: options.Request.AppliesTo}
	err = journal.WithLock(options.WorkspaceRoot, journal.LockOwner{
		Command: strings.Join(workflowDefineArgv(options), " "),
	}, func(lock *journal.WriteLock) error {
		return defineWorkflowLocked(lock, options, actor, &result)
	})
	if err != nil {
		return DefineResult{}, err
	}
	return result, nil
}

// AuthorizeDefinition remains for legacy in-process callers; standalone attribution
// deliberately grants no workflow authority.
func AuthorizeDefinition(workspaceRoot, actorRole string) error {
	return nil
}

func defineWorkflowLocked(lock *journal.WriteLock, options DefineOptions, actor standaloneActor, result *DefineResult) error {
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
	requestInput := any(map[string]any{"Name": options.Request.Name, "AppliesTo": options.Request.AppliesTo, "States": options.Request.States, "Categories": options.Request.Categories, "Transitions": options.Request.Transitions})
	aliases := map[string]string{"Name": "--name", "AppliesTo": "--applies-to", "States": "--states", "Categories": "--categories", "Transitions": "--transitions"}
	if options.RequestInput != nil {
		requestInput, aliases = options.RequestInput, nil
	}
	digest, fields, err := journal.ComputeRequestDigest(map[string]any{"actor": actor, "operation": "workflow.define", "project_id": projectConfig.ProjectID, "request": requestInput}, options.Actor, aliases)
	if err != nil {
		return err
	}
	idempotencyRequest := journal.IdempotencyRequest{OperationScope: "workflow.define:" + projectConfig.ProjectID, IdempotencyKey: options.IdempotencyKey, RequestDigest: digest, RequestFieldDigests: fields}
	if outcome, found, err := journal.ReplayIdempotent(conn, idempotencyRequest); err != nil {
		return err
	} else if found {
		decoded, err := defineResultFromJSON(outcome.ResultJSON, outcome.EventID, true)
		if err != nil {
			return err
		}
		*result = decoded
		return nil
	}
	effects, err := BuildDefinitionEffects(conn, options.Request)
	if err != nil {
		return err
	}
	sequence, err := journal.NextAggregateSequence(filepath.Join(options.WorkspaceRoot, config.WorkspaceDirName, config.JournalDirName), journal.Aggregate{ID: options.Request.Name, Type: "workflow"})
	if err != nil {
		return err
	}
	event := journal.Event{
		EventSchema: journal.EventSchemaV2,
		EventType:   "workflow.defined",
		OccurredAt:  time.Now().UTC(),
		ActorV2:     actor.eventActor(),
		Aggregate:   journal.Aggregate{ID: options.Request.Name, Sequence: sequence, Type: "workflow"},
		CommandName: "workflow.define",
		Effects:     effects,
	}
	result.Name, result.AppliesTo, result.OccurredAt = options.Request.Name, options.Request.AppliesTo, event.OccurredAt.UTC()
	if err := loadDefinitionCountsAfterEffects(conn, options.Request, result); err != nil {
		return err
	}
	result.IdempotencyKey = options.IdempotencyKey
	resultJSON, err := defineResultJSON(*result)
	if err != nil {
		return err
	}
	outcome, err := journal.CommitIdempotent(lock, conn, idempotencyRequest, event, func(journal.Event) (journal.IdempotencyMutation, error) {
		return journal.IdempotencyMutation{EntityID: result.Name, ResultJSON: resultJSON}, nil
	})
	if err != nil {
		return err
	}
	decoded, err := defineResultFromJSON(outcome.ResultJSON, outcome.EventID, outcome.Replayed)
	if err != nil {
		return err
	}
	*result = decoded
	return nil
}

func loadDefinitionCountsAfterEffects(conn *sql.DB, request DefinitionRequest, result *DefineResult) error {
	result.States = len(request.States)
	result.Transitions = len(materializedTransitions(request))
	result.Requirements = len(request.Requirements)
	return nil
}

func defineResultJSON(result DefineResult) (string, error) {
	payload, err := journal.MarshalCanonicalJSON(map[string]any{"applies_to": result.AppliesTo, "event_occurred_at": result.OccurredAt.Format(time.RFC3339Nano), "idempotency_key": result.IdempotencyKey, "name": result.Name, "requirements": result.Requirements, "replayed": false, "state_count": result.States, "transition_count": result.Transitions})
	return string(payload), err
}

func defineResultFromJSON(raw, eventID string, replayed bool) (DefineResult, error) {
	var result struct {
		AppliesTo      string `json:"applies_to"`
		OccurredAt     string `json:"event_occurred_at"`
		IdempotencyKey string `json:"idempotency_key"`
		Name           string `json:"name"`
		Requirements   int    `json:"requirements"`
		States         int    `json:"state_count"`
		Transitions    int    `json:"transition_count"`
	}
	if err := json.Unmarshal([]byte(raw), &result); err != nil {
		return DefineResult{}, fmt.Errorf("decode stored workflow define result: %w", err)
	}
	occurredAt, err := time.Parse(time.RFC3339Nano, result.OccurredAt)
	if err != nil {
		return DefineResult{}, fmt.Errorf("decode stored workflow define occurred_at: %w", err)
	}
	return DefineResult{Name: result.Name, AppliesTo: result.AppliesTo, States: result.States, Transitions: result.Transitions, Requirements: result.Requirements, EventID: eventID, OccurredAt: occurredAt, IdempotencyKey: result.IdempotencyKey, Replayed: replayed}, nil
}

func openWorkflowDB(workspaceRoot string) (*sql.DB, error) {
	return fiddb.Open(filepath.Join(workspaceRoot, config.WorkspaceDirName, config.DatabaseFileName))
}

func closeWorkflowDB(conn *sql.DB) {
	_ = conn.Close()
}

func workflowDefineArgv(options DefineOptions) []string {
	request := options.Request
	argv := []string{
		"fid", "workflow", "define", "--name", request.Name,
		"--applies-to", request.AppliesTo,
		"--states", strings.Join(request.States, ","),
		"--categories", strings.Join(request.Categories, ","),
	}
	transitions := make([]string, 0, len(request.Transitions))
	for _, transition := range request.Transitions {
		transitions = append(transitions, transition.From+"->"+transition.To)
	}
	argv = append(argv, "--transitions", strings.Join(transitions, ","))
	for _, requirement := range request.Requirements {
		argv = append(argv, "--require", requirement.From+"->"+requirement.To+":"+requirement.Rule+requirement.ParamsJSON)
	}
	return append(argv, "--actor", options.Actor, "--idempotency-key", options.IdempotencyKey)
}
