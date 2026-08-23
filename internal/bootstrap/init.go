package bootstrap

import (
	"bytes"
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
)

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
	ProjectID   string
	Prefix      string
	JournalPath string
	EventID     string
	OccurredAt  time.Time
	SeededRoles []string
}

type initStage string

const (
	initStageValidated         initStage = "validated"
	initStageStaged            initStage = "staged"
	initStageSchemaInitialized initStage = "schema-initialized"
	initStageGenesisAppended   initStage = "genesis-appended"
	initStageRebuilt           initStage = "rebuilt"
	initStageCompared          initStage = "compared"
	initStageSynced            initStage = "synced"
	initStageActivated         initStage = "activated"
)

type initHooks struct {
	before func(initStage) error
}

type initPlan struct {
	workspaceRoot string
	targetPaths   config.WorkspacePaths
	projectConfig config.ProjectConfig
	prefix        string
	explicit      bool
	migrationsDir string
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
	return initializeWorkspaceWithHooks(options, initHooks{})
}

func initializeWorkspaceWithHooks(options InitOptions, hooks initHooks) (InitResult, error) {
	plan, err := prepareInitialization(options)
	if err != nil {
		return InitResult{}, err
	}
	if err := runInitHook(hooks, initStageValidated); err != nil {
		return InitResult{}, err
	}

	stagingRoot, err := os.MkdirTemp(plan.workspaceRoot, ".fid-staging-")
	if err != nil {
		return InitResult{}, fmt.Errorf("create staging directory: %w", err)
	}
	activated := false
	defer func() {
		if !activated {
			_ = os.RemoveAll(stagingRoot)
		}
	}()

	stagePaths, err := stageWorkspace(stagingRoot, plan.projectConfig, hooks)
	if err != nil {
		return InitResult{}, err
	}
	result, err := initializeStagedWorkspace(stagingRoot, stagePaths.DatabasePath, plan, hooks)
	if err != nil {
		return InitResult{}, err
	}
	if err := verifyStagedRebuild(stagingRoot, stagePaths.DatabasePath, plan.migrationsDir, hooks); err != nil {
		return InitResult{}, err
	}
	if err := syncStagedWorkspace(stagePaths.WorkspaceDir); err != nil {
		return InitResult{}, err
	}
	if err := runInitHook(hooks, initStageSynced); err != nil {
		return InitResult{}, err
	}
	if err := runInitHook(hooks, initStageActivated); err != nil {
		return InitResult{}, err
	}
	if err := activateStagedWorkspace(stagePaths.WorkspaceDir, plan.targetPaths.WorkspaceDir); err != nil {
		return InitResult{}, err
	}
	activated = true
	if err := finishActivation(stagingRoot, plan.workspaceRoot); err != nil {
		return InitResult{}, err
	}
	result.ProjectID = plan.projectConfig.ProjectID
	return result, nil
}

func prepareInitialization(options InitOptions) (initPlan, error) {
	if strings.TrimSpace(options.WorkspaceRoot) == "" {
		return initPlan{}, fmt.Errorf("initialize workspace: workspace root is required")
	}
	if strings.TrimSpace(options.MigrationsDir) == "" {
		return initPlan{}, fmt.Errorf("initialize workspace: migrations dir is required")
	}
	workspaceRoot, err := filepath.Abs(options.WorkspaceRoot)
	if err != nil {
		return initPlan{}, fmt.Errorf("resolve workspace root: %w", err)
	}
	if info, err := os.Stat(workspaceRoot); err != nil || !info.IsDir() {
		if err != nil {
			return initPlan{}, fmt.Errorf("stat workspace root %q: %w", workspaceRoot, err)
		}
		return initPlan{}, fmt.Errorf("workspace root %q is not a directory", workspaceRoot)
	}
	targetPaths := config.WorkspacePathsForRoot(workspaceRoot)
	if _, err := os.Stat(targetPaths.WorkspaceDir); err == nil {
		return initPlan{}, &AlreadyInitializedError{Path: targetPaths.WorkspaceDir}
	} else if !os.IsNotExist(err) {
		return initPlan{}, fmt.Errorf("stat workspace dir %q: %w", targetPaths.WorkspaceDir, err)
	}
	prefix := options.Prefix
	explicit := strings.TrimSpace(options.Prefix) != ""
	if strings.TrimSpace(prefix) == "" {
		prefix = "FID"
	}
	projectConfig, err := config.NewDefaultProjectConfig(prefix)
	if err != nil {
		return initPlan{}, err
	}
	return initPlan{workspaceRoot, targetPaths, projectConfig, prefix, explicit, options.MigrationsDir}, nil
}

func stageWorkspace(stagingRoot string, projectConfig config.ProjectConfig, hooks initHooks) (config.WorkspacePaths, error) {
	stagePaths, err := config.WriteBootstrapFiles(stagingRoot, projectConfig, config.DefaultPricingConfig())
	if err != nil {
		return config.WorkspacePaths{}, err
	}
	if err := os.Mkdir(filepath.Join(stagePaths.WorkspaceDir, projectConfig.ArtifactRoot), 0o755); err != nil {
		return config.WorkspacePaths{}, fmt.Errorf("create artifacts directory: %w", err)
	}
	if err := runInitHook(hooks, initStageStaged); err != nil {
		return config.WorkspacePaths{}, err
	}
	return stagePaths, nil
}

func initializeStagedWorkspace(stagingRoot, databasePath string, plan initPlan, hooks initHooks) (InitResult, error) {
	if err := fiddb.MigrateUp(databasePath, plan.migrationsDir); err != nil {
		return InitResult{}, err
	}
	if err := runInitHook(hooks, initStageSchemaInitialized); err != nil {
		return InitResult{}, err
	}
	result, err := initializeJournaledState(stagingRoot, databasePath, plan.projectConfig.ProjectID, plan.prefix)
	if err != nil {
		return InitResult{}, err
	}
	if err := runInitHook(hooks, initStageGenesisAppended); err != nil {
		return InitResult{}, err
	}
	return result, nil
}

func verifyStagedRebuild(stagingRoot, databasePath, migrationsDir string, hooks initHooks) error {
	beforeRebuild, err := os.ReadFile(databasePath)
	if err != nil {
		return fmt.Errorf("read staged projection before rebuild: %w", err)
	}
	if _, err := journal.RebuildWorkspace(stagingRoot, migrationsDir); err != nil {
		return fmt.Errorf("rebuild staged projection: %w", err)
	}
	if err := runInitHook(hooks, initStageRebuilt); err != nil {
		return err
	}
	afterRebuild, err := os.ReadFile(databasePath)
	if err != nil {
		return fmt.Errorf("read staged projection after rebuild: %w", err)
	}
	if !bytes.Equal(beforeRebuild, afterRebuild) {
		return fmt.Errorf("verify staged rebuild: projection bytes differ")
	}
	return runInitHook(hooks, initStageCompared)
}

func finishActivation(stagingRoot, workspaceRoot string) error {
	if err := os.Remove(stagingRoot); err != nil {
		return fmt.Errorf("remove empty staging directory: %w", err)
	}
	if err := syncPath(workspaceRoot); err != nil {
		return fmt.Errorf("sync workspace root: %w", err)
	}
	return nil
}

func runInitHook(hooks initHooks, stage initStage) error {
	if hooks.before == nil {
		return nil
	}
	if err := hooks.before(stage); err != nil {
		return fmt.Errorf("initialize stage %s: %w", stage, err)
	}
	return nil
}

func isStagingDirectoryName(name string) bool {
	return strings.HasPrefix(name, ".fid-staging-")
}

func activateStagedWorkspace(stagedWorkspace, target string) error {
	if _, err := os.Lstat(target); err == nil {
		return &AlreadyInitializedError{Path: target}
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("inspect activation target %q: %w", target, err)
	}
	if err := os.Rename(stagedWorkspace, target); err != nil {
		return fmt.Errorf("activate staged workspace: %w", err)
	}
	return nil
}

func syncStagedWorkspace(workspaceDir string) error {
	var directories []string
	err := filepath.WalkDir(workspaceDir, func(path string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() {
			directories = append(directories, path)
			return nil
		}
		return syncPath(path)
	})
	if err != nil {
		return fmt.Errorf("sync staged workspace: %w", err)
	}
	for index := len(directories) - 1; index >= 0; index-- {
		if err := syncPath(directories[index]); err != nil {
			return fmt.Errorf("sync staged workspace: %w", err)
		}
	}
	return nil
}

func syncPath(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("open for sync %q: %w", path, err)
	}
	if err := file.Sync(); err != nil {
		_ = file.Close()
		return fmt.Errorf("sync %q: %w", path, err)
	}
	if err := file.Close(); err != nil {
		return fmt.Errorf("close synced path %q: %w", path, err)
	}
	return nil
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

func initializeJournaledState(workspaceRoot, dbPath, projectID, prefix string) (InitResult, error) {
	var result InitResult
	err := journal.WithLock(workspaceRoot, journal.LockOwner{
		Command:   "fid init",
		ActorRole: "human",
	}, func(lock *journal.WriteLock) error {
		event := BuildProjectInitializedEvent(projectID)
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
			OccurredAt:  appended.OccurredAt.UTC(),
		}
		return nil
	})
	if err != nil {
		return InitResult{}, err
	}

	return result, nil
}

func BuildProjectInitializedEvent(projectID string) journal.Event {
	return journal.Event{
		EventSchema: journal.EventSchemaV2,
		EventType:   "project.initialized",
		ActorV2: journal.EventActor{
			ActorID:       "fid-bootstrap",
			ActorType:     "service",
			Producer:      "fid",
			ProvenanceRef: "fid:init",
		},
		Aggregate:   journal.Aggregate{ID: projectID, Sequence: 1, Type: "project"},
		CommandName: "project.init",
		Effects:     []journal.Effect{},
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

func closeBootstrapDB(conn *sql.DB) {
	_ = conn.Close()
}

func IsAlreadyInitialized(err error) bool {
	var target *AlreadyInitializedError
	return errors.As(err, &target)
}
