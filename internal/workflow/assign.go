package workflow

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/mbell-nzl/fid/internal/journal"
	"github.com/mbell-nzl/fid/internal/policy"
)

// AssignOptions supplies an assignment request and its write context.
type AssignOptions struct {
	WorkspaceRoot string
	MigrationsDir string
	ActorRole     string
	Request       AssignmentRequest
}

// AssignResult reports the workflow and scope recorded by one assignment.
type AssignResult struct {
	Name    string
	Kind    string
	Project bool
}

// Assign appends and projects one atomic workflow.assigned event.
func Assign(options AssignOptions) (AssignResult, error) {
	if strings.TrimSpace(options.WorkspaceRoot) == "" {
		return AssignResult{}, fmt.Errorf("assign workflow: workspace root is required")
	}
	if strings.TrimSpace(options.MigrationsDir) == "" {
		return AssignResult{}, fmt.Errorf("assign workflow: migrations dir is required")
	}
	if err := AuthorizeAssignment(options.WorkspaceRoot, options.ActorRole); err != nil {
		return AssignResult{}, err
	}
	conn, err := openWorkflowDB(options.WorkspaceRoot)
	if err != nil {
		return AssignResult{}, err
	}
	err = ValidateAssignmentRequest(conn, options.Request)
	closeWorkflowDB(conn)
	if err != nil {
		return AssignResult{}, err
	}

	result := AssignResult{Name: options.Request.Name, Kind: options.Request.Kind, Project: options.Request.Project}
	err = journal.WithLock(options.WorkspaceRoot, journal.LockOwner{
		Command:   strings.Join(workflowAssignArgv(options), " "),
		ActorRole: options.ActorRole,
	}, func(lock *journal.WriteLock) error {
		return assignWorkflowLocked(lock, options)
	})
	if err != nil {
		return AssignResult{}, err
	}
	return result, nil
}

// AuthorizeAssignment verifies the actor may assign workflows before parsing work begins.
func AuthorizeAssignment(workspaceRoot, actorRole string) error {
	return policy.Authorize(workspaceRoot, actorRole, "workflow.assign")
}

func assignWorkflowLocked(lock *journal.WriteLock, options AssignOptions) error {
	if _, err := journal.EnsureProjectionFreshWithLock(lock, options.MigrationsDir); err != nil {
		return err
	}
	conn, err := openWorkflowDB(options.WorkspaceRoot)
	if err != nil {
		return err
	}
	defer closeWorkflowDB(conn)
	if err := ValidateAssignmentRequest(conn, options.Request); err != nil {
		return err
	}
	replacing, err := assignmentExists(conn, options.Request)
	if err != nil {
		return err
	}
	event, err := journal.PrepareEvent(lock, journal.Event{
		EventSchema: journal.CurrentEventSchema,
		EventType:   "workflow.assigned",
		Actor:       journal.Actor{Role: options.ActorRole},
		Command: journal.Command{
			Argv:       workflowAssignArgv(options),
			CWD:        ".",
			FIDVersion: workflowFIDVersion,
		},
	})
	if err != nil {
		return err
	}
	event.Effects = assignmentEffects(options.Request, replacing, event.OccurredAt.UTC().Format("2006-01-02T15:04:05Z"))
	if err := journal.DryRunProjectEvent(conn, event); err != nil {
		return err
	}
	appended, err := journal.AppendEvent(lock, event)
	if err != nil {
		return err
	}
	return projectWorkflowEvent(options.WorkspaceRoot, conn, appended)
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
		argv = append(argv, "--project")
	} else {
		argv = append(argv, "--kind", options.Request.Kind)
	}
	return append(argv, "--by", options.ActorRole)
}
