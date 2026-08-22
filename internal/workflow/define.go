package workflow

import (
	"database/sql"
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

const workflowFIDVersion = "0.1.0"

// DefineOptions supplies a validated workflow declaration and its write context.
type DefineOptions struct {
	WorkspaceRoot string
	MigrationsDir string
	ActorRole     string
	Request       DefinitionRequest
}

// DefineResult reports the projected rows written for one workflow declaration.
type DefineResult struct {
	Name         string
	AppliesTo    string
	States       int
	Transitions  int
	Requirements int
	EventID      string
	OccurredAt   time.Time
}

// Define appends and projects one atomic workflow.defined event.
func Define(options DefineOptions) (DefineResult, error) {
	if strings.TrimSpace(options.WorkspaceRoot) == "" {
		return DefineResult{}, fmt.Errorf("define workflow: workspace root is required")
	}
	if strings.TrimSpace(options.MigrationsDir) == "" {
		return DefineResult{}, fmt.Errorf("define workflow: migrations dir is required")
	}
	if err := AuthorizeDefinition(options.WorkspaceRoot, options.ActorRole); err != nil {
		return DefineResult{}, err
	}
	conn, err := openWorkflowDB(options.WorkspaceRoot)
	if err != nil {
		return DefineResult{}, err
	}
	_, err = BuildDefinitionEffects(conn, options.Request)
	closeWorkflowDB(conn)
	if err != nil {
		return DefineResult{}, err
	}

	result := DefineResult{Name: options.Request.Name, AppliesTo: options.Request.AppliesTo}
	err = journal.WithLock(options.WorkspaceRoot, journal.LockOwner{
		Command:   strings.Join(workflowDefineArgv(options), " "),
		ActorRole: options.ActorRole,
	}, func(lock *journal.WriteLock) error {
		return defineWorkflowLocked(lock, options, &result)
	})
	if err != nil {
		return DefineResult{}, err
	}
	return result, nil
}

// AuthorizeDefinition verifies the actor may define workflows before parsing work begins.
func AuthorizeDefinition(workspaceRoot, actorRole string) error {
	return policy.Authorize(workspaceRoot, actorRole, "workflow.define")
}

func defineWorkflowLocked(lock *journal.WriteLock, options DefineOptions, result *DefineResult) error {
	if _, err := journal.EnsureProjectionFreshWithLock(lock, options.MigrationsDir); err != nil {
		return err
	}
	conn, err := openWorkflowDB(options.WorkspaceRoot)
	if err != nil {
		return err
	}
	defer closeWorkflowDB(conn)
	effects, err := BuildDefinitionEffects(conn, options.Request)
	if err != nil {
		return err
	}
	event, err := journal.PrepareEvent(lock, journal.Event{
		EventSchema: journal.CurrentEventSchema,
		EventType:   "workflow.defined",
		Actor:       journal.Actor{Role: options.ActorRole},
		Command: journal.Command{
			Argv:       workflowDefineArgv(options),
			CWD:        ".",
			FIDVersion: workflowFIDVersion,
		},
		Effects: effects,
	})
	if err != nil {
		return err
	}
	if err := journal.DryRunProjectEvent(conn, event); err != nil {
		return err
	}
	appended, err := journal.AppendEvent(lock, event)
	if err != nil {
		return err
	}
	if err := projectWorkflowEvent(options.WorkspaceRoot, conn, appended); err != nil {
		return err
	}
	if err := loadDefinitionCounts(conn, options.Request.Name, result); err != nil {
		return err
	}
	result.EventID = appended.EventID
	result.OccurredAt = appended.OccurredAt.UTC()
	return nil
}

func openWorkflowDB(workspaceRoot string) (*sql.DB, error) {
	return fiddb.Open(filepath.Join(workspaceRoot, config.WorkspaceDirName, config.DatabaseFileName))
}

func closeWorkflowDB(conn *sql.DB) {
	_ = conn.Close()
}

func projectWorkflowEvent(workspaceRoot string, conn *sql.DB, event journal.Event) error {
	eventHash, err := journal.HashEvent(event)
	if err != nil {
		return fmt.Errorf("hash workflow event %q: %w", event.EventID, err)
	}
	journalFile := event.OccurredAt.UTC().Format("2006-01-02") + ".jsonl"
	journalPath := filepath.Join(workspaceRoot, config.WorkspaceDirName, config.JournalDirName, journalFile)
	info, err := os.Stat(journalPath)
	if err != nil {
		return fmt.Errorf("stat workflow journal %q: %w", journalPath, err)
	}
	return journal.ProjectEvent(conn, journal.ProjectionInput{
		Event: event, EventHash: eventHash, File: journalFile,
		Line: workflowJournalLineCount(journalPath), Offset: info.Size(),
	})
}

func workflowJournalLineCount(path string) int {
	payload, err := os.ReadFile(path)
	if err != nil {
		return 0
	}
	trimmed := strings.TrimSpace(string(payload))
	if trimmed == "" {
		return 0
	}
	return len(strings.Split(trimmed, "\n"))
}

func loadDefinitionCounts(conn *sql.DB, name string, result *DefineResult) error {
	if err := conn.QueryRow(`SELECT COUNT(*) FROM workflow_states WHERE workflow_id = ?`, name).Scan(&result.States); err != nil {
		return fmt.Errorf("count workflow %q states: %w", name, err)
	}
	if err := conn.QueryRow(`SELECT COUNT(*) FROM workflow_transitions WHERE workflow_id = ?`, name).Scan(&result.Transitions); err != nil {
		return fmt.Errorf("count workflow %q transitions: %w", name, err)
	}
	if err := conn.QueryRow(`SELECT COUNT(*) FROM workflow_transition_requires WHERE workflow_id = ?`, name).Scan(&result.Requirements); err != nil {
		return fmt.Errorf("count workflow %q requirements: %w", name, err)
	}
	return nil
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
	return append(argv, "--by", options.ActorRole)
}
