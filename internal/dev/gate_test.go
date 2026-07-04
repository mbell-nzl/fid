package dev

import (
	"database/sql"
	"errors"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mbell-nzl/fid/internal/claims"
	"github.com/mbell-nzl/fid/internal/issues"
	"github.com/mbell-nzl/fid/internal/journal"
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
	conn := openDevDB(t, filepath.Join(workspace, ".fid", "fid.sqlite"))
	activityBefore := countDevActivityRows(t, conn)
	closeDevDB(t, conn)

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

	var (
		gateEffect     *journal.Effect
		activityEffect *journal.Effect
	)
	for index := range last.Effects {
		effect := &last.Effects[index]
		switch {
		case effect.Table == "gate_runs" && effect.Op == "insert":
			gateEffect = effect
		case effect.Table == "activity_log" && effect.Op == "activity":
			activityEffect = effect
		}
	}
	if gateEffect == nil {
		t.Fatal("missing gate_runs insert effect")
	}
	if activityEffect == nil {
		t.Fatal("missing gate_recorded activity effect")
	}
	if got := activityEffect.Values["event"]; got != "gate_recorded" {
		t.Fatalf("activity event = %#v, want gate_recorded", got)
	}
	if got := activityEffect.Values["ref_type"]; got != "issue" {
		t.Fatalf("activity ref_type = %#v, want issue", got)
	}
	if got := activityEffect.Values["ref_id"]; got != created.ID {
		t.Fatalf("activity ref_id = %#v, want %q", got, created.ID)
	}
	if got := activityEffect.Values["actor"]; got != "developer" {
		t.Fatalf("activity actor = %#v, want developer", got)
	}
	details, ok := activityEffect.Values["details"].(map[string]any)
	if !ok {
		t.Fatalf("activity details = %#v, want map", activityEffect.Values["details"])
	}
	for key, want := range map[string]string{
		"gate_run_id":  result.GateID,
		"command":      "go test ./...",
		"agent_run_id": started.RunID,
	} {
		if got := details[key]; got != want {
			t.Fatalf("activity details[%q] = %#v, want %#v", key, got, want)
		}
	}
	if got := details["attempt_number"]; got != float64(1) {
		t.Fatalf("activity details[attempt_number] = %#v, want 1", got)
	}
	if got := details["exit_code"]; got != float64(0) {
		t.Fatalf("activity details[exit_code] = %#v, want 0", got)
	}

	conn = openDevDB(t, filepath.Join(workspace, ".fid", "fid.sqlite"))
	defer closeDevDB(t, conn)

	if got := countDevActivityRows(t, conn); got != activityBefore+1 {
		t.Fatalf("activity_log rows after Gate() = %d, want %d", got, activityBefore+1)
	}

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

	var (
		activityCount int
		eventName     string
		actor         string
		activityJSON  string
	)
	if err := conn.QueryRow(`
		SELECT COUNT(*), event, COALESCE(actor, ''), COALESCE(details, '')
		FROM activity_log
		WHERE journal_event_id = ?
		GROUP BY event, actor, details
	`, last.EventID).Scan(&activityCount, &eventName, &actor, &activityJSON); err != nil {
		t.Fatalf("query gate activity row: %v", err)
	}
	if activityCount != 1 {
		t.Fatalf("activity_log grouped row count = %d, want 1", activityCount)
	}
	if eventName != "gate_recorded" {
		t.Fatalf("activity_log event = %q, want gate_recorded", eventName)
	}
	if actor != "developer" {
		t.Fatalf("activity_log actor = %q, want developer", actor)
	}
	for _, want := range []string{
		`"gate_run_id":"` + result.GateID + `"`,
		`"attempt_number":1`,
		`"command":"go test ./..."`,
		`"exit_code":0`,
		`"agent_run_id":"` + started.RunID + `"`,
	} {
		if !strings.Contains(activityJSON, want) {
			t.Fatalf("activity_log details = %q, want substring %q", activityJSON, want)
		}
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
	conn := openDevDB(t, filepath.Join(workspace, ".fid", "fid.sqlite"))
	activityBefore := countDevActivityRows(t, conn)
	closeDevDB(t, conn)

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

	conn = openDevDB(t, filepath.Join(workspace, ".fid", "fid.sqlite"))
	defer closeDevDB(t, conn)
	if got := countDevActivityRows(t, conn); got != activityBefore {
		t.Fatalf("activity_log rows after rejected Gate() = %d, want unchanged %d", got, activityBefore)
	}
}

func stringPointer(value string) *string {
	return &value
}

func countDevActivityRows(t *testing.T, conn *sql.DB) int {
	t.Helper()

	var count int
	if err := conn.QueryRow(`SELECT COUNT(*) FROM activity_log`).Scan(&count); err != nil {
		t.Fatalf("count activity_log rows: %v", err)
	}
	return count
}
