package dev

import (
	"database/sql"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/mbell-nzl/fid/internal/bootstrap"
	"github.com/mbell-nzl/fid/internal/claims"
	fiddb "github.com/mbell-nzl/fid/internal/db"
	"github.com/mbell-nzl/fid/internal/issues"
	"github.com/mbell-nzl/fid/internal/journal"
)

func TestStartClaimsTransitionsAndOpensRunInOneEvent(t *testing.T) {
	t.Parallel()

	workspace := newDevWorkspace(t)

	created, err := issues.Create(issues.CreateOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testDevMigrationsDir(t),
		Kind:          "task",
		Title:         "Start composite target",
		ActorRole:     "developer",
	})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	journalPath := journalPathForDevWorkspace(t, workspace)
	linesBefore := countDevJournalLines(journalPath)

	result, err := Start(StartOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testDevMigrationsDir(t),
		IssueID:       created.ID,
		ActorRole:     "developer",
		Model:         "claude-sonnet-4-6",
	})
	if err != nil {
		t.Fatalf("Start() error = %v", err)
	}

	if result.IssueID != created.ID {
		t.Fatalf("Start().IssueID = %q, want %q", result.IssueID, created.ID)
	}
	if result.ClaimedBy != "developer" {
		t.Fatalf("Start().ClaimedBy = %q, want developer", result.ClaimedBy)
	}
	if result.FromState != "todo" {
		t.Fatalf("Start().FromState = %q, want todo", result.FromState)
	}
	if result.ToState != "in_progress" {
		t.Fatalf("Start().ToState = %q, want in_progress", result.ToState)
	}
	if result.Model != "claude-sonnet-4-6" {
		t.Fatalf("Start().Model = %q, want claude-sonnet-4-6", result.Model)
	}
	if !strings.HasPrefix(result.RunID, "RUN-") {
		t.Fatalf("Start().RunID = %q, want RUN- prefix", result.RunID)
	}
	if result.ExpiresAt == "" {
		t.Fatal("Start().ExpiresAt = empty, want claim expiry")
	}

	if got := countDevJournalLines(journalPath); got != linesBefore+1 {
		t.Fatalf("journal lines after Start() = %d, want %d", got, linesBefore+1)
	}

	stored := readDevStoredEvents(t, journalPath)
	last := stored[len(stored)-1]
	if last.EventType != "dev.started" {
		t.Fatalf("last event type = %q, want dev.started", last.EventType)
	}

	var (
		claimEffect    *journal.Effect
		transition     *journal.Effect
		activityEffect *journal.Effect
		runEffect      *journal.Effect
	)
	for index := range last.Effects {
		effect := &last.Effects[index]
		switch {
		case effect.Table == "claims" && effect.Op == "insert":
			claimEffect = effect
		case effect.Table == "issues" && effect.Op == "update":
			transition = effect
		case effect.Table == "activity_log" && effect.Op == "activity":
			activityEffect = effect
		case effect.Table == "agent_runs" && effect.Op == "insert":
			runEffect = effect
		}
	}
	if claimEffect == nil {
		t.Fatal("missing claim insert effect")
	}
	if transition == nil {
		t.Fatal("missing issue transition effect")
	}
	if activityEffect == nil {
		t.Fatal("missing status_changed activity effect")
	}
	if runEffect == nil {
		t.Fatal("missing agent run insert effect")
	}

	conn := openDevDB(t, filepath.Join(workspace, ".fid", "fid.sqlite"))
	defer closeDevDB(t, conn)

	var (
		status    string
		claimRole string
		runRole   string
		model     string
		endedAt   sql.NullString
	)
	if err := conn.QueryRow(`SELECT status FROM issues WHERE id = ?`, created.ID).Scan(&status); err != nil {
		t.Fatalf("query issues row: %v", err)
	}
	if status != "in_progress" {
		t.Fatalf("issues.status = %q, want in_progress", status)
	}
	if err := conn.QueryRow(`
		SELECT claimed_by
		FROM claims
		WHERE ref_type = 'issue' AND ref_id = ?
	`, created.ID).Scan(&claimRole); err != nil {
		t.Fatalf("query claims row: %v", err)
	}
	if claimRole != "developer" {
		t.Fatalf("claims.claimed_by = %q, want developer", claimRole)
	}
	if err := conn.QueryRow(`
		SELECT role, model, ended_at
		FROM agent_runs
		WHERE id = ?
	`, result.RunID).Scan(&runRole, &model, &endedAt); err != nil {
		t.Fatalf("query agent_runs row: %v", err)
	}
	if runRole != "developer" {
		t.Fatalf("agent_runs.role = %q, want developer", runRole)
	}
	if model != "claude-sonnet-4-6" {
		t.Fatalf("agent_runs.model = %q, want claude-sonnet-4-6", model)
	}
	if endedAt.Valid {
		t.Fatalf("agent_runs.ended_at = %#v, want NULL open run", endedAt)
	}
}

func TestStartRejectsConflictingLiveClaimBeforeAppend(t *testing.T) {
	t.Parallel()

	workspace := newDevWorkspace(t)

	created, err := issues.Create(issues.CreateOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testDevMigrationsDir(t),
		Kind:          "task",
		Title:         "Conflicting claim target",
		ActorRole:     "developer",
	})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if _, err := claims.Acquire(claims.AcquireOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testDevMigrationsDir(t),
		IssueID:       created.ID,
		ActorRole:     "qa_reviewer",
	}); err != nil {
		t.Fatalf("Acquire() error = %v", err)
	}

	journalPath := journalPathForDevWorkspace(t, workspace)
	linesBefore := countDevJournalLines(journalPath)

	_, err = Start(StartOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testDevMigrationsDir(t),
		IssueID:       created.ID,
		ActorRole:     "developer",
		Model:         "claude-sonnet-4-6",
	})
	if err == nil {
		t.Fatal("Start() error = nil, want live claim conflict")
	}

	var conflictErr *claims.LiveClaimConflictError
	if !errors.As(err, &conflictErr) {
		t.Fatalf("Start() error = %T, want *claims.LiveClaimConflictError", err)
	}
	if conflictErr.ClaimedBy != "qa_reviewer" {
		t.Fatalf("LiveClaimConflictError.ClaimedBy = %q, want qa_reviewer", conflictErr.ClaimedBy)
	}

	if got := countDevJournalLines(journalPath); got != linesBefore {
		t.Fatalf("journal lines after rejected Start() = %d, want unchanged %d", got, linesBefore)
	}
}

func TestStartRejectsInvalidWorkflowEdgeBeforeAppend(t *testing.T) {
	t.Parallel()

	workspace := newDevWorkspace(t)

	created, err := issues.Create(issues.CreateOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testDevMigrationsDir(t),
		Kind:          "task",
		Title:         "Already started target",
		ActorRole:     "developer",
	})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if _, err := issues.Transition(issues.TransitionOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testDevMigrationsDir(t),
		ID:            created.ID,
		ToState:       "in_progress",
		ActorRole:     "developer",
	}); err == nil {
		t.Fatal("Transition() error = nil, want missing claim rejection")
	}
	if _, err := claims.Acquire(claims.AcquireOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testDevMigrationsDir(t),
		IssueID:       created.ID,
		ActorRole:     "developer",
	}); err != nil {
		t.Fatalf("Acquire() error = %v", err)
	}
	if _, err := issues.Transition(issues.TransitionOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testDevMigrationsDir(t),
		ID:            created.ID,
		ToState:       "in_progress",
		ActorRole:     "developer",
	}); err != nil {
		t.Fatalf("Transition() error = %v", err)
	}

	journalPath := journalPathForDevWorkspace(t, workspace)
	linesBefore := countDevJournalLines(journalPath)

	_, err = Start(StartOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testDevMigrationsDir(t),
		IssueID:       created.ID,
		ActorRole:     "developer",
		Model:         "claude-sonnet-4-6",
	})
	if err == nil {
		t.Fatal("Start() error = nil, want invalid transition rejection")
	}

	var invalidErr *issues.InvalidTransitionError
	if !errors.As(err, &invalidErr) {
		t.Fatalf("Start() error = %T, want *issues.InvalidTransitionError", err)
	}
	if invalidErr.CurrentState != "in_progress" {
		t.Fatalf("InvalidTransitionError.CurrentState = %q, want in_progress", invalidErr.CurrentState)
	}

	if got := countDevJournalLines(journalPath); got != linesBefore {
		t.Fatalf("journal lines after invalid Start() = %d, want unchanged %d", got, linesBefore)
	}
}

func testDevMigrationsDir(t *testing.T) string {
	t.Helper()

	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller(0) failed")
	}

	return filepath.Clean(filepath.Join(filepath.Dir(currentFile), "..", "..", "migrations"))
}

func newDevWorkspace(t *testing.T) string {
	t.Helper()

	workspace := t.TempDir()
	if _, err := bootstrap.InitializeWorkspace(bootstrap.InitOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testDevMigrationsDir(t),
	}); err != nil {
		t.Fatalf("InitializeWorkspace() error = %v", err)
	}

	return workspace
}

func journalPathForDevWorkspace(t *testing.T, workspace string) string {
	t.Helper()

	entries, err := os.ReadDir(filepath.Join(workspace, ".fid", "journal"))
	if err != nil {
		t.Fatalf("ReadDir(journal) error = %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("journal file count = %d, want 1", len(entries))
	}
	return filepath.Join(workspace, ".fid", "journal", entries[0].Name())
}

func countDevJournalLines(path string) int {
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

func readDevStoredEvents(t *testing.T, journalPath string) []journal.Event {
	t.Helper()

	payload, err := os.ReadFile(journalPath)
	if err != nil {
		t.Fatalf("ReadFile(%q) error = %v", journalPath, err)
	}

	lines := strings.Split(strings.TrimSpace(string(payload)), "\n")
	events := make([]journal.Event, 0, len(lines))
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		var event journal.Event
		if err := json.Unmarshal([]byte(line), &event); err != nil {
			t.Fatalf("Unmarshal(journal line) error = %v", err)
		}
		events = append(events, event)
	}
	return events
}

func openDevDB(t *testing.T, dbPath string) *sql.DB {
	t.Helper()

	conn, err := fiddb.Open(dbPath)
	if err != nil {
		t.Fatalf("Open(%q) error = %v", dbPath, err)
	}
	return conn
}

func closeDevDB(t *testing.T, conn *sql.DB) {
	t.Helper()

	if err := conn.Close(); err != nil {
		t.Fatalf("Close() error = %v", err)
	}
}
