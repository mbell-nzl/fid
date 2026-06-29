package dev

import (
	"database/sql"
	"errors"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mbell-nzl/fid/internal/claims"
	"github.com/mbell-nzl/fid/internal/issues"
)

func TestGateRecordsReportedResultOnActiveAttempt(t *testing.T) {
	t.Parallel()

	workspace := newDevWorkspace(t)

	created, err := issues.Create(issues.CreateOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testDevMigrationsDir(t),
		Kind:          "task",
		Title:         "Gate target",
		ActorRole:     "developer",
	})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	started, err := Start(StartOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testDevMigrationsDir(t),
		IssueID:       created.ID,
		ActorRole:     "developer",
		Model:         "claude-sonnet-4-6",
	})
	if err != nil {
		t.Fatalf("Start() error = %v", err)
	}

	journalPath := journalPathForDevWorkspace(t, workspace)
	linesBefore := countDevJournalLines(journalPath)

	result, err := Gate(GateOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testDevMigrationsDir(t),
		IssueID:       created.ID,
		ActorRole:     "developer",
		Command:       "go test ./...",
		ExitCode:      0,
		OutputExcerpt: stringPointer("ok"),
	})
	if err != nil {
		t.Fatalf("Gate() error = %v", err)
	}

	if result.IssueID != created.ID {
		t.Fatalf("Gate().IssueID = %q, want %q", result.IssueID, created.ID)
	}
	if result.AttemptNumber != 1 {
		t.Fatalf("Gate().AttemptNumber = %d, want 1", result.AttemptNumber)
	}
	if result.Command != "go test ./..." {
		t.Fatalf("Gate().Command = %q, want go test ./...", result.Command)
	}
	if result.ExitCode != 0 {
		t.Fatalf("Gate().ExitCode = %d, want 0", result.ExitCode)
	}
	if result.GatesPassed != 1 || result.GatesFailed != 0 {
		t.Fatalf("Gate() counts = (%d, %d), want (1, 0)", result.GatesPassed, result.GatesFailed)
	}
	if !strings.HasPrefix(result.GateID, "GATE-") {
		t.Fatalf("Gate().GateID = %q, want GATE- prefix", result.GateID)
	}
	if result.AgentRunID == nil || *result.AgentRunID != started.RunID {
		t.Fatalf("Gate().AgentRunID = %#v, want %q", result.AgentRunID, started.RunID)
	}

	if got := countDevJournalLines(journalPath); got != linesBefore+1 {
		t.Fatalf("journal lines after Gate() = %d, want %d", got, linesBefore+1)
	}

	stored := readDevStoredEvents(t, journalPath)
	last := stored[len(stored)-1]
	if last.EventType != "dev.gate.recorded" {
		t.Fatalf("last event type = %q, want dev.gate.recorded", last.EventType)
	}
	if len(last.Effects) != 1 {
		t.Fatalf("len(last.Effects) = %d, want 1", len(last.Effects))
	}
	effect := last.Effects[0]
	if effect.Table != "gate_runs" || effect.Op != "insert" {
		t.Fatalf("gate effect = (%s, %s), want (gate_runs, insert)", effect.Table, effect.Op)
	}

	conn := openDevDB(t, filepath.Join(workspace, ".fid", "fid.sqlite"))
	defer closeDevDB(t, conn)

	var (
		attemptNumber int
		agentRunID    sql.NullString
		recordedBy    string
		command       string
		exitCode      int
		outputExcerpt sql.NullString
		ranAt         string
	)
	if err := conn.QueryRow(`
		SELECT attempt_number, agent_run_id, recorded_by_role, command, exit_code, output_excerpt, ran_at
		FROM gate_runs
		WHERE id = ?
	`, result.GateID).Scan(
		&attemptNumber,
		&agentRunID,
		&recordedBy,
		&command,
		&exitCode,
		&outputExcerpt,
		&ranAt,
	); err != nil {
		t.Fatalf("query gate_runs row: %v", err)
	}
	if attemptNumber != 1 {
		t.Fatalf("gate_runs.attempt_number = %d, want 1", attemptNumber)
	}
	if !agentRunID.Valid || agentRunID.String != started.RunID {
		t.Fatalf("gate_runs.agent_run_id = %#v, want %q", agentRunID, started.RunID)
	}
	if recordedBy != "developer" {
		t.Fatalf("gate_runs.recorded_by_role = %q, want developer", recordedBy)
	}
	if command != "go test ./..." {
		t.Fatalf("gate_runs.command = %q, want go test ./...", command)
	}
	if exitCode != 0 {
		t.Fatalf("gate_runs.exit_code = %d, want 0", exitCode)
	}
	if !outputExcerpt.Valid || outputExcerpt.String != "ok" {
		t.Fatalf("gate_runs.output_excerpt = %#v, want ok", outputExcerpt)
	}
	if ranAt == "" {
		t.Fatal("gate_runs.ran_at = empty, want timestamp")
	}
}

func TestGateUsesNextAttemptNumberAfterCompletion(t *testing.T) {
	t.Parallel()

	workspace := newDevWorkspace(t)

	created, err := issues.Create(issues.CreateOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testDevMigrationsDir(t),
		Kind:          "task",
		Title:         "Gate attempt target",
		ActorRole:     "developer",
	})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if _, err := Start(StartOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testDevMigrationsDir(t),
		IssueID:       created.ID,
		ActorRole:     "developer",
		Model:         "claude-sonnet-4-6",
	}); err != nil {
		t.Fatalf("Start() error = %v", err)
	}
	if _, err := issues.Complete(issues.CompleteOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testDevMigrationsDir(t),
		IssueID:       created.ID,
		ActorRole:     "developer",
	}); err != nil {
		t.Fatalf("Complete() error = %v", err)
	}

	result, err := Gate(GateOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testDevMigrationsDir(t),
		IssueID:       created.ID,
		ActorRole:     "developer",
		Command:       "golangci-lint run ./...",
		ExitCode:      1,
	})
	if err != nil {
		t.Fatalf("Gate() error = %v", err)
	}

	if result.AttemptNumber != 2 {
		t.Fatalf("Gate().AttemptNumber = %d, want 2", result.AttemptNumber)
	}
	if result.GatesPassed != 0 || result.GatesFailed != 1 {
		t.Fatalf("Gate() counts = (%d, %d), want (0, 1)", result.GatesPassed, result.GatesFailed)
	}
}

func TestGateRejectsWithoutStartedStateBeforeAppend(t *testing.T) {
	t.Parallel()

	workspace := newDevWorkspace(t)

	created, err := issues.Create(issues.CreateOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testDevMigrationsDir(t),
		Kind:          "task",
		Title:         "Not started gate target",
		ActorRole:     "developer",
	})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	journalPath := journalPathForDevWorkspace(t, workspace)
	linesBefore := countDevJournalLines(journalPath)

	_, err = Gate(GateOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testDevMigrationsDir(t),
		IssueID:       created.ID,
		ActorRole:     "developer",
		Command:       "go test ./...",
		ExitCode:      0,
	})
	if err == nil {
		t.Fatal("Gate() error = nil, want started-category rejection")
	}
	if !strings.Contains(err.Error(), "started-category") {
		t.Fatalf("Gate() error = %q, want started-category message", err.Error())
	}

	if got := countDevJournalLines(journalPath); got != linesBefore {
		t.Fatalf("journal lines after rejected Gate() = %d, want unchanged %d", got, linesBefore)
	}
}

func TestGateRejectsWithoutLiveClaimBeforeAppend(t *testing.T) {
	t.Parallel()

	workspace := newDevWorkspace(t)

	created, err := issues.Create(issues.CreateOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testDevMigrationsDir(t),
		Kind:          "task",
		Title:         "No claim gate target",
		ActorRole:     "developer",
	})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	if _, err := Start(StartOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testDevMigrationsDir(t),
		IssueID:       created.ID,
		ActorRole:     "developer",
		Model:         "claude-sonnet-4-6",
	}); err != nil {
		t.Fatalf("Start() error = %v", err)
	}
	if _, err := claims.Release(claims.ReleaseOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testDevMigrationsDir(t),
		IssueID:       created.ID,
		ActorRole:     "developer",
	}); err != nil {
		t.Fatalf("Release() error = %v", err)
	}

	journalPath := journalPathForDevWorkspace(t, workspace)
	linesBefore := countDevJournalLines(journalPath)

	_, err = Gate(GateOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testDevMigrationsDir(t),
		IssueID:       created.ID,
		ActorRole:     "developer",
		Command:       "go test ./...",
		ExitCode:      0,
	})
	if err == nil {
		t.Fatal("Gate() error = nil, want missing-claim rejection")
	}

	var claimErr *claims.ClaimNotHeldError
	if !errors.As(err, &claimErr) {
		t.Fatalf("Gate() error = %T, want *claims.ClaimNotHeldError", err)
	}

	if got := countDevJournalLines(journalPath); got != linesBefore {
		t.Fatalf("journal lines after rejected Gate() = %d, want unchanged %d", got, linesBefore)
	}
}

func stringPointer(value string) *string {
	return &value
}
