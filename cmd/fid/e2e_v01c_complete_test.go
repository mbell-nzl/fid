package main

import (
	"database/sql"
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	fiddb "github.com/mbell-nzl/fid/internal/db"
	"github.com/mbell-nzl/fid/internal/journal"
)

func TestV01CHasCompletionReportGateAndRebuildParity(t *testing.T) {
	workspace := t.TempDir()
	runInitCommand(t, workspace, "init")

	_, createStderr := runFIDCommand(
		t,
		workspace,
		"issue", "create",
		"--kind", "task",
		"--title", "Phase 8 completion gate parity",
		"--by", "developer",
	)
	issueID := parseCreatedIssueID(t, createStderr)
	_, claimStderr := runFIDCommand(t, workspace, "issue", "claim", issueID, "--by", "developer")
	if !strings.Contains(claimStderr, "claimed "+issueID+" by developer (expires ") {
		t.Fatalf("claim stderr = %q, want confirmation", claimStderr)
	}

	_, transitionStderr := runFIDCommand(
		t,
		workspace,
		"issue", "transition", issueID,
		"--to", "in_progress",
		"--by", "developer",
	)
	if !strings.Contains(transitionStderr, "transitioned "+issueID+": todo → in_progress") {
		t.Fatalf("todo -> in_progress stderr = %q, want confirmation", transitionStderr)
	}

	journalPath := mustJournalPath(t, workspace)
	linesBeforeReject := len(splitJournalLines(mustReadCommandFile(t, journalPath)))
	statusBeforeReject := mustGetIssueJSON(t, workspace, issueID)

	stdout, stderr, err := runFIDCommandExpectError(
		t,
		workspace,
		"issue", "transition", issueID,
		"--to", "dev_complete",
		"--by", "developer",
	)
	if err == nil {
		t.Fatal("in_progress -> dev_complete error = nil, want gated rejection")
	}
	if stdout != "" {
		t.Fatalf("reject stdout = %q, want empty", stdout)
	}
	for _, want := range []string{
		"has_completion_report",
		issueID,
	} {
		if !strings.Contains(stderr, want) {
			t.Fatalf("reject stderr = %q, want substring %q", stderr, want)
		}
	}
	if got := len(splitJournalLines(mustReadCommandFile(t, journalPath))); got != linesBeforeReject {
		t.Fatalf("journal lines after reject = %d, want %d", got, linesBeforeReject)
	}
	if got := mustGetIssueJSON(t, workspace, issueID); got != statusBeforeReject {
		t.Fatalf("issue JSON after reject = %q, want %q", got, statusBeforeReject)
	}

	_, completeStderr := runFIDCommand(
		t,
		workspace,
		"issue", "complete", issueID,
		"--by", "developer",
		"--notes-md", "project-docs/docs/phases/phase-8/tasks/task-3-done.md",
		"--tokens", "1600",
	)
	if !strings.Contains(completeStderr, "completed "+issueID+" (attempt 1)") {
		t.Fatalf("attempt 1 complete stderr = %q, want confirmation", completeStderr)
	}

	_, transitionStderr = runFIDCommand(
		t,
		workspace,
		"issue", "transition", issueID,
		"--to", "dev_complete",
		"--by", "developer",
	)
	if !strings.Contains(transitionStderr, "transitioned "+issueID+": in_progress → dev_complete") {
		t.Fatalf("in_progress -> dev_complete stderr = %q, want confirmation", transitionStderr)
	}

	seedTask3GateRuns(t, workspace, issueID, []task3GateRunSeed{
		{
			ID:             "GATE-task3-001",
			AttemptNumber:  2,
			RecordedByRole: "developer",
			Command:        "go test -race ./...",
			ExitCode:       0,
			OutputExcerpt:  task3StringPointer("ok\tgithub.com/mbell-nzl/fid/cmd/fid\t0.321s"),
			RanAt:          "2026-06-20T09:00:00Z",
		},
		{
			ID:             "GATE-task3-002",
			AttemptNumber:  2,
			RecordedByRole: "developer",
			Command:        "golangci-lint run ./...",
			ExitCode:       1,
			OutputExcerpt:  task3StringPointer("cmd/fid/e2e_v01c_complete_test.go: line exceeds budget"),
			RanAt:          "2026-06-20T09:01:00Z",
		},
	})

	_, secondCompleteStderr := runFIDCommand(
		t,
		workspace,
		"issue", "complete", issueID,
		"--by", "developer",
		"--notes-md", "project-docs/docs/phases/phase-8/tasks/task-3-done.md",
		"--tokens", "1700",
		"--human-intervention",
		"--intervention-note", "paired with qa on rebuild parity",
	)
	if !strings.Contains(secondCompleteStderr, "completed "+issueID+" (attempt 2)") {
		t.Fatalf("attempt 2 complete stderr = %q, want confirmation", secondCompleteStderr)
	}

	beforeCompletionsJSON := mustGetTask3CompletionsJSON(t, workspace, issueID)
	assertTask3CompletionSnapshot(t, beforeCompletionsJSON)
	beforeGateRuns := snapshotTask3GateRuns(t, workspace, issueID)

	dbPath := filepath.Join(workspace, ".fid", "fid.sqlite")
	if err := os.Remove(dbPath); err != nil {
		t.Fatalf("Remove(%q) error = %v", dbPath, err)
	}

	rebuildStdout, rebuildStderr := runFIDCommand(t, workspace, "rebuild")
	if !strings.Contains(rebuildStdout, "rebuild complete") || rebuildStderr != "" {
		t.Fatalf("rebuild stdout/stderr = %q/%q, want stdout completion and empty stderr", rebuildStdout, rebuildStderr)
	}

	if got := mustGetTask3CompletionsJSON(t, workspace, issueID); got != beforeCompletionsJSON {
		t.Fatalf("completions JSON after rebuild = %q, want %q", got, beforeCompletionsJSON)
	}
	if got := snapshotTask3GateRuns(t, workspace, issueID); !reflect.DeepEqual(got, beforeGateRuns) {
		t.Fatalf("gate_runs after rebuild = %#v, want %#v", got, beforeGateRuns)
	}
}

type task3CompletionView struct {
	ID                    string                `json:"id"`
	AttemptNumber         int                   `json:"attempt_number"`
	Status                *string               `json:"status"`
	GatesPassed           int                   `json:"gates_passed"`
	GatesFailed           int                   `json:"gates_failed"`
	GateResults           []task3GateResultView `json:"gate_results"`
	HumanIntervention     bool                  `json:"human_intervention"`
	HumanInterventionNote *string               `json:"human_intervention_note"`
}

type task3GateResultView struct {
	Cmd           string  `json:"cmd"`
	ExitCode      int     `json:"exit_code"`
	RanAt         string  `json:"ran_at"`
	OutputExcerpt *string `json:"output_excerpt"`
}

type task3GateRunSeed struct {
	ID             string
	AttemptNumber  int
	RecordedByRole string
	Command        string
	ExitCode       int
	OutputExcerpt  *string
	RanAt          string
}

type task3GateRunRow struct {
	ID            string
	IssueID       string
	AttemptNumber int
	Command       string
	ExitCode      int
	OutputExcerpt string
	RanAt         string
}

func mustGetTask3CompletionsJSON(t *testing.T, workspace, issueID string) string {
	t.Helper()

	rows := mustReadV01ECompletionsForIssue(t, workspace, issueID)
	views := make([]task3CompletionView, 0, len(rows))
	for _, row := range rows {
		var gateResults []task3GateResultView
		if err := json.Unmarshal([]byte(row.GateResults), &gateResults); err != nil {
			t.Fatalf("Unmarshal(gate_results) error = %v", err)
		}
		var status *string
		if row.Status != "" {
			status = &row.Status
		}
		var interventionNote *string
		if row.HumanInterventionNote != "" {
			interventionNote = &row.HumanInterventionNote
		}
		views = append(views, task3CompletionView{
			ID:                    row.ID,
			AttemptNumber:         row.AttemptNumber,
			Status:                status,
			GatesPassed:           row.GatesPassed,
			GatesFailed:           row.GatesFailed,
			GateResults:           gateResults,
			HumanIntervention:     row.HumanIntervention,
			HumanInterventionNote: interventionNote,
		})
	}
	payload, err := json.Marshal(views)
	if err != nil {
		t.Fatalf("Marshal(completions) error = %v", err)
	}
	return string(payload)
}

func assertTask3CompletionSnapshot(t *testing.T, completionsJSON string) {
	t.Helper()

	var got []task3CompletionView
	if err := json.Unmarshal([]byte(completionsJSON), &got); err != nil {
		t.Fatalf("Unmarshal(completions) error = %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("completions len = %d, want 2", len(got))
	}
	if got[0].AttemptNumber != 2 {
		t.Fatalf("latest attempt_number = %d, want 2", got[0].AttemptNumber)
	}
	if got[0].Status == nil || *got[0].Status != "dev_complete" {
		t.Fatalf("latest status = %#v, want dev_complete", got[0].Status)
	}
	if !got[0].HumanIntervention {
		t.Fatal("latest human_intervention = false, want true")
	}
	if got[0].HumanInterventionNote == nil || *got[0].HumanInterventionNote != "paired with qa on rebuild parity" {
		t.Fatalf("latest human_intervention_note = %#v, want note", got[0].HumanInterventionNote)
	}
	if got[0].GatesPassed != 1 {
		t.Fatalf("latest gates_passed = %d, want 1", got[0].GatesPassed)
	}
	if got[0].GatesFailed != 1 {
		t.Fatalf("latest gates_failed = %d, want 1", got[0].GatesFailed)
	}
	wantGateResults := []task3GateResultView{
		{
			Cmd:           "go test -race ./...",
			ExitCode:      0,
			RanAt:         "2026-06-20T09:00:00Z",
			OutputExcerpt: task3StringPointer("ok\tgithub.com/mbell-nzl/fid/cmd/fid\t0.321s"),
		},
		{
			Cmd:           "golangci-lint run ./...",
			ExitCode:      1,
			RanAt:         "2026-06-20T09:01:00Z",
			OutputExcerpt: task3StringPointer("cmd/fid/e2e_v01c_complete_test.go: line exceeds budget"),
		},
	}
	if !reflect.DeepEqual(got[0].GateResults, wantGateResults) {
		t.Fatalf("latest gate_results = %#v, want %#v", got[0].GateResults, wantGateResults)
	}
	if got[1].AttemptNumber != 1 {
		t.Fatalf("prior attempt_number = %d, want 1", got[1].AttemptNumber)
	}
	if got[1].Status == nil || *got[1].Status != "superseded" {
		t.Fatalf("prior status = %#v, want superseded", got[1].Status)
	}
	if got[1].HumanIntervention {
		t.Fatal("prior human_intervention = true, want false")
	}
	if got[1].HumanInterventionNote != nil {
		t.Fatalf("prior human_intervention_note = %#v, want nil", got[1].HumanInterventionNote)
	}
	if got[1].GatesPassed != 0 || got[1].GatesFailed != 0 {
		t.Fatalf("prior gate counts = %d/%d, want 0/0", got[1].GatesPassed, got[1].GatesFailed)
	}
	if len(got[1].GateResults) != 0 {
		t.Fatalf("prior gate_results = %#v, want empty", got[1].GateResults)
	}
}

func seedTask3GateRuns(t *testing.T, workspace, issueID string, runs []task3GateRunSeed) {
	t.Helper()

	effects := make([]journal.Effect, 0, len(runs))
	for _, run := range runs {
		values := map[string]any{
			"issue_id":         issueID,
			"attempt_number":   run.AttemptNumber,
			"recorded_by_role": run.RecordedByRole,
			"command":          run.Command,
			"exit_code":        run.ExitCode,
			"ran_at":           run.RanAt,
		}
		if run.OutputExcerpt != nil {
			values["output_excerpt"] = *run.OutputExcerpt
		}
		effects = append(effects, journal.Effect{
			Op:     "insert",
			Table:  "gate_runs",
			PK:     map[string]any{"id": run.ID},
			Values: values,
		})
	}

	appendCommandIssueEvent(t, workspace, journal.Event{
		EventSchema: journal.EventSchemaV1,
		EventType:   "gate.seeded",
		Actor: journal.Actor{
			Role: "developer",
		},
		Command: journal.Command{
			Argv:       []string{"fid", "maintenance", "seed"},
			CWD:        ".",
			FIDVersion: "0.1.0",
		},
		Effects: effects,
	})
}

func snapshotTask3GateRuns(t *testing.T, workspace, issueID string) []task3GateRunRow {
	t.Helper()

	conn, err := fiddb.Open(filepath.Join(workspace, ".fid", "fid.sqlite"))
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	defer closeCommandDB(t, conn)

	rows, err := conn.Query(`
		SELECT id, issue_id, attempt_number, command, exit_code, COALESCE(output_excerpt, ''), ran_at
		FROM gate_runs
		WHERE issue_id = ?
		ORDER BY attempt_number, ran_at, id
	`, issueID)
	if err != nil {
		t.Fatalf("query gate_runs: %v", err)
	}
	defer closeTask3Rows(t, rows)

	var snapshot []task3GateRunRow
	for rows.Next() {
		var row task3GateRunRow
		if err := rows.Scan(
			&row.ID,
			&row.IssueID,
			&row.AttemptNumber,
			&row.Command,
			&row.ExitCode,
			&row.OutputExcerpt,
			&row.RanAt,
		); err != nil {
			t.Fatalf("scan gate_runs: %v", err)
		}
		snapshot = append(snapshot, row)
	}
	if err := rows.Err(); err != nil {
		t.Fatalf("iterate gate_runs: %v", err)
	}
	return snapshot
}

func closeTask3Rows(t *testing.T, rows *sql.Rows) {
	t.Helper()

	if err := rows.Close(); err != nil {
		t.Fatalf("Rows.Close() error = %v", err)
	}
}

func task3StringPointer(value string) *string {
	return &value
}
