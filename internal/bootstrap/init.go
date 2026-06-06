package bootstrap

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/mbell-nzl/fid/internal/config"
	fiddb "github.com/mbell-nzl/fid/internal/db"
	"github.com/mbell-nzl/fid/internal/journal"
)

const fidVersion = "0.1.0"

type RoleSeed struct {
	Name        string
	Description string
	IsHuman     bool
}

type InitOptions struct {
	WorkspaceRoot string
	MigrationsDir string
	Prefix        string
}

type InitResult struct {
	Prefix      string
	JournalPath string
	EventID     string
	SeededRoles []string
}

type AlreadyInitializedError struct {
	Path string
}

func (err *AlreadyInitializedError) Error() string {
	if err == nil {
		return "<nil>"
	}
	return fmt.Sprintf("%s already exists", err.Path)
}

func InitializeWorkspace(options InitOptions) (InitResult, error) {
	if strings.TrimSpace(options.WorkspaceRoot) == "" {
		return InitResult{}, fmt.Errorf("initialize workspace: workspace root is required")
	}
	if strings.TrimSpace(options.MigrationsDir) == "" {
		return InitResult{}, fmt.Errorf("initialize workspace: migrations dir is required")
	}

	paths := config.WorkspacePathsForRoot(options.WorkspaceRoot)
	if _, err := os.Stat(paths.WorkspaceDir); err == nil {
		return InitResult{}, &AlreadyInitializedError{Path: config.WorkspaceDirName + "/"}
	} else if !os.IsNotExist(err) {
		return InitResult{}, fmt.Errorf("stat workspace dir %q: %w", paths.WorkspaceDir, err)
	}

	prefix := options.Prefix
	explicitPrefix := strings.TrimSpace(options.Prefix) != ""
	if strings.TrimSpace(prefix) == "" {
		prefix = DefaultPrefix(options.WorkspaceRoot)
	}

	if _, err := config.WriteBootstrapFiles(options.WorkspaceRoot, config.DefaultProjectConfig(prefix), config.DefaultPricingConfig()); err != nil {
		return InitResult{}, err
	}

	success := false
	defer func() {
		if !success {
			_ = os.RemoveAll(paths.WorkspaceDir)
		}
	}()

	if err := fiddb.MigrateUp(paths.DatabasePath, options.MigrationsDir); err != nil {
		return InitResult{}, err
	}

	result, err := initializeJournaledState(options.WorkspaceRoot, paths.DatabasePath, prefix, explicitPrefix)
	if err != nil {
		return InitResult{}, err
	}

	success = true
	return result, nil
}

func DefaultPrefix(workspaceRoot string) string {
	return strings.ToUpper(filepath.Base(workspaceRoot))
}

func DefaultRoles() []RoleSeed {
	return []RoleSeed{
		{Name: "developer", Description: "Developer"},
		{Name: "qa_reviewer", Description: "QA Reviewer"},
		{Name: "qa_lead", Description: "QA Lead"},
		{Name: "tl_lite", Description: "TL-Lite"},
		{Name: "tl", Description: "TL"},
		{Name: "workflow_runner", Description: "Workflow Runner"},
		{Name: "human", Description: "Human", IsHuman: true},
	}
}

func initializeJournaledState(workspaceRoot, dbPath, prefix string, explicitPrefix bool) (InitResult, error) {
	var result InitResult
	err := journal.WithLock(workspaceRoot, journal.LockOwner{
		Command:   initCommand(prefix),
		ActorRole: "human",
	}, func(lock *journal.WriteLock) error {
		event := BuildProjectInitializedEvent(prefix, explicitPrefix)
		appended, err := journal.AppendEvent(lock, event)
		if err != nil {
			return err
		}

		eventHash, err := journal.HashEvent(appended)
		if err != nil {
			return fmt.Errorf("hash bootstrap event %q: %w", appended.EventID, err)
		}

		journalFile := appended.OccurredAt.UTC().Format("2006-01-02") + ".jsonl"
		journalPath := filepath.Join(workspaceRoot, config.WorkspaceDirName, config.JournalDirName, journalFile)
		info, err := os.Stat(journalPath)
		if err != nil {
			return fmt.Errorf("stat bootstrap journal %q: %w", journalPath, err)
		}

		conn, err := fiddb.Open(dbPath)
		if err != nil {
			return err
		}
		defer closeBootstrapDB(conn)

		if err := journal.ProjectEvent(conn, journal.ProjectionInput{
			Event:     appended,
			EventHash: eventHash,
			File:      journalFile,
			Line:      1,
			Offset:    info.Size(),
		}); err != nil {
			return err
		}

		result = InitResult{
			Prefix:      prefix,
			JournalPath: filepath.ToSlash(filepath.Join(config.WorkspaceDirName, config.JournalDirName, journalFile)),
			EventID:     appended.EventID,
			SeededRoles: SeededRoleNames(),
		}
		return nil
	})
	if err != nil {
		return InitResult{}, err
	}

	return result, nil
}

func BuildProjectInitializedEvent(prefix string, explicitPrefix bool) journal.Event {
	roles := DefaultRoles()
	effects := make([]journal.Effect, 0, len(roles))
	for _, role := range roles {
		effects = append(effects, journal.Effect{
			Op:    "insert",
			Table: "roles",
			PK: map[string]any{
				"name": role.Name,
			},
			Values: map[string]any{
				"description": role.Description,
				"is_human":    role.IsHuman,
				"active":      true,
			},
		})
	}

	argv := []string{"fid", "init"}
	if explicitPrefix {
		argv = append(argv, "--prefix", prefix)
	}

	return journal.Event{
		EventSchema: journal.CurrentEventSchema,
		EventType:   "project.initialized",
		Actor: journal.Actor{
			Role: "human",
		},
		Command: journal.Command{
			Argv:       argv,
			CWD:        ".",
			FIDVersion: fidVersion,
		},
		Effects: effects,
	}
}

func SeededRoleNames() []string {
	roles := DefaultRoles()
	names := make([]string, 0, len(roles))
	for _, role := range roles {
		names = append(names, role.Name)
	}
	return names
}

func initCommand(prefix string) string {
	if strings.TrimSpace(prefix) == "" {
		return "fid init"
	}
	return "fid init --prefix " + prefix
}

func closeBootstrapDB(conn *sql.DB) {
	_ = conn.Close()
}

func IsAlreadyInitialized(err error) bool {
	var target *AlreadyInitializedError
	return errors.As(err, &target)
}
