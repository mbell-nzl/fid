package main

import (
	"database/sql"
	"encoding/json"
	"path/filepath"
	"strings"
	"testing"

	fiddb "github.com/mbell-nzl/fid/internal/db"
)

func TestIssueCompleteCommandCreatesAttemptOneAndReturnsJSONCompletionFields(t *testing.T) {
	t.Parallel()
	workspace := t.TempDir()
	clonePristineWorkspaceInto(t, workspace)

	_, createStderr := runFIDCommand(t, workspace, "issue", "create", "--kind", "task", "--title", "CLI completion", "--by", "developer")
	issueID := parseCreatedIssueID(t, createStderr)

	stdout, stderr := runFIDCommand(
		t,
		workspace,
		"issue", "complete", issueID,
		"--by", "developer",
		"--notes-md", "project-docs/docs/phases/phase-8/tasks/task-1-done.md",
		"--tokens", "14823",
		"--commits", "abc123, def456",
		"--files", "internal/issues/complete.go",
		"--files", "cmd/fid/issue.go",
	)
	if stdout != "" {
		t.Fatalf("stdout = %q, want empty", stdout)
	}
	if !strings.Contains(stderr, "completed "+issueID+" (attempt 1)") {
		t.Fatalf("stderr = %q, want completion confirmation", stderr)
	}
	if !strings.Contains(stderr, "  completion: TC-") {
		t.Fatalf("stderr = %q, want completion id", stderr)
	}
	if !strings.Contains(stderr, "  tokens: 14823") {
		t.Fatalf("stderr = %q, want tokens line", stderr)
	}
	if !strings.Contains(stderr, "  commits: abc123, def456") {
		t.Fatalf("stderr = %q, want commits line", stderr)
	}

	completions := mustGetIssueCompletionViews(t, workspace, issueID)
	if len(completions) != 1 {
		t.Fatalf("completions = %#v, want one completion", completions)
	}
	completion := completions[0]
	if !strings.HasPrefix(completion.ID, "TC-") {
		t.Fatalf("completion id = %q, want TC- prefix", completion.ID)
	}
	if completion.AttemptNumber != 1 {
		t.Fatalf("completion attempt_number = %d, want 1", completion.AttemptNumber)
	}
	if completion.DeveloperRole == nil || *completion.DeveloperRole != "developer" {
		t.Fatalf("completion developer_role = %#v, want developer", completion.DeveloperRole)
	}
	if completion.Status == nil || *completion.Status != "dev_complete" {
		t.Fatalf("completion status = %#v, want dev_complete", completion.Status)
	}
	if completion.NotesMDPath == nil || *completion.NotesMDPath != "project-docs/docs/phases/phase-8/tasks/task-1-done.md" {
		t.Fatalf("completion notes_md_path = %#v, want task-1 report path", completion.NotesMDPath)
	}
	if completion.ActualTokens == nil || *completion.ActualTokens != 14823 {
		t.Fatalf("completion actual_tokens = %#v, want 14823", completion.ActualTokens)
	}
	if strings.Join(completion.CommitHashes, ",") != "abc123,def456" {
		t.Fatalf("completion commit_hashes = %#v, want [abc123 def456]", completion.CommitHashes)
	}
	if strings.Join(completion.FilesTouched, ",") != "internal/issues/complete.go,cmd/fid/issue.go" {
		t.Fatalf("completion files_touched = %#v, want stored files", completion.FilesTouched)
	}
	if completion.HumanIntervention {
		t.Fatal("completion human_intervention = true, want false")
	}
	if completion.HumanInterventionNote != nil {
		t.Fatalf("completion human_intervention_note = %#v, want nil", completion.HumanInterventionNote)
	}

	conn, err := fiddb.Open(filepath.Join(workspace, ".fid", "fid.sqlite"))
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	defer closeCommandDB(t, conn)

	event := mustReadLastJournalEvent(t, mustJournalPath(t, workspace))
	if event.EventType != "issue.completed" {
		t.Fatalf("event type = %q, want issue.completed", event.EventType)
	}
	if got := countActivityLogRows(t, conn); got != 2 {
		t.Fatalf("activity_log rows after completion = %d, want 2", got)
	}

	var (
		activityCount int
		activityEvent string
		actor         string
		details       sql.NullString
	)
	if err := conn.QueryRow(`
		SELECT COUNT(*), event, COALESCE(actor, ''), details
		FROM activity_log
		WHERE journal_event_id = ?
		GROUP BY event, actor, details
	`, event.EventID).Scan(&activityCount, &activityEvent, &actor, &details); err != nil {
		t.Fatalf("query activity_log row: %v", err)
	}
	if activityCount != 1 {
		t.Fatalf("activity_log grouped row count = %d, want 1", activityCount)
	}
	if activityEvent != "task_completed" {
		t.Fatalf("activity_log event = %q, want task_completed", activityEvent)
	}
	if actor != "developer" {
		t.Fatalf("activity_log actor = %q, want developer", actor)
	}
	for _, want := range []string{completion.ID, `"status":"dev_complete"`, `"attempt_number":1`} {
		if !strings.Contains(details.String, want) {
			t.Fatalf("activity_log details = %q, want substring %q", details.String, want)
		}
	}
}

func TestIssueCompleteCommandSupersedesPriorAttemptAndStoresInterventionFields(t *testing.T) {
	t.Parallel()
	workspace := t.TempDir()
	clonePristineWorkspaceInto(t, workspace)

	_, createStderr := runFIDCommand(t, workspace, "issue", "create", "--kind", "task", "--title", "CLI supersession", "--by", "developer")
	issueID := parseCreatedIssueID(t, createStderr)

	runFIDCommand(t, workspace, "issue", "complete", issueID, "--by", "developer", "--tokens", "4000")
	_, secondStderr := runFIDCommand(
		t,
		workspace,
		"issue", "complete", issueID,
		"--by", "developer",
		"--tokens", "5000",
		"--human-intervention",
		"--intervention-note", "paired on final review",
	)
	if !strings.Contains(secondStderr, "completed "+issueID+" (attempt 2)") {
		t.Fatalf("second stderr = %q, want attempt 2 confirmation", secondStderr)
	}

	completions := mustGetIssueCompletionViews(t, workspace, issueID)
	if len(completions) != 2 {
		t.Fatalf("completions len = %d, want 2", len(completions))
	}
	if completions[0].AttemptNumber != 2 {
		t.Fatalf("latest attempt_number = %d, want 2", completions[0].AttemptNumber)
	}
	if completions[0].Status == nil || *completions[0].Status != "dev_complete" {
		t.Fatalf("latest status = %#v, want dev_complete", completions[0].Status)
	}
	if completions[0].ActualTokens == nil || *completions[0].ActualTokens != 5000 {
		t.Fatalf("latest actual_tokens = %#v, want 5000", completions[0].ActualTokens)
	}
	if !completions[0].HumanIntervention {
		t.Fatal("latest human_intervention = false, want true")
	}
	if completions[0].HumanInterventionNote == nil || *completions[0].HumanInterventionNote != "paired on final review" {
		t.Fatalf("latest human_intervention_note = %#v, want paired on final review", completions[0].HumanInterventionNote)
	}
	if completions[1].AttemptNumber != 1 {
		t.Fatalf("prior attempt_number = %d, want 1", completions[1].AttemptNumber)
	}
	if completions[1].Status == nil || *completions[1].Status != "superseded" {
		t.Fatalf("prior status = %#v, want superseded", completions[1].Status)
	}
	if completions[1].HumanIntervention {
		t.Fatal("prior human_intervention = true, want false")
	}
	if completions[1].HumanInterventionNote != nil {
		t.Fatalf("prior human_intervention_note = %#v, want nil", completions[1].HumanInterventionNote)
	}
}

type commandCompletionView struct {
	ID                    string   `json:"id"`
	AttemptNumber         int      `json:"attempt_number"`
	DeveloperRole         *string  `json:"developer_role"`
	Status                *string  `json:"status"`
	NotesMDPath           *string  `json:"notes_md_path"`
	ActualTokens          *int     `json:"actual_tokens"`
	CommitHashes          []string `json:"commit_hashes"`
	FilesTouched          []string `json:"files_touched"`
	HumanIntervention     bool     `json:"human_intervention"`
	HumanInterventionNote *string  `json:"human_intervention_note"`
}

func mustGetIssueCompletionViews(t *testing.T, workspace, issueID string) []commandCompletionView {
	t.Helper()

	rows := mustReadV01ECompletionsForIssue(t, workspace, issueID)
	views := make([]commandCompletionView, 0, len(rows))
	for _, row := range rows {
		var commitHashes, filesTouched []string
		if row.CommitHashes != "" {
			if err := json.Unmarshal([]byte(row.CommitHashes), &commitHashes); err != nil {
				t.Fatalf("Unmarshal(commit_hashes) error = %v", err)
			}
		}
		if row.FilesTouched != "" {
			if err := json.Unmarshal([]byte(row.FilesTouched), &filesTouched); err != nil {
				t.Fatalf("Unmarshal(files_touched) error = %v", err)
			}
		}
		var developerRole, status, notesMDPath, interventionNote *string
		if row.DeveloperRole != "" {
			developerRole = &row.DeveloperRole
		}
		if row.Status != "" {
			status = &row.Status
		}
		if row.NotesMDPath != "" {
			notesMDPath = &row.NotesMDPath
		}
		if row.HumanInterventionNote != "" {
			interventionNote = &row.HumanInterventionNote
		}
		var actualTokens *int
		if row.ActualTokens >= 0 {
			actualTokens = &row.ActualTokens
		}
		views = append(views, commandCompletionView{
			ID:                    row.ID,
			AttemptNumber:         row.AttemptNumber,
			DeveloperRole:         developerRole,
			Status:                status,
			NotesMDPath:           notesMDPath,
			ActualTokens:          actualTokens,
			CommitHashes:          commitHashes,
			FilesTouched:          filesTouched,
			HumanIntervention:     row.HumanIntervention,
			HumanInterventionNote: interventionNote,
		})
	}
	return views
}

func TestIssueCompleteCommandRejectsEpicBeforeAppend(t *testing.T) {
	t.Parallel()
	workspace := t.TempDir()
	clonePristineWorkspaceInto(t, workspace)

	_, createStderr := runFIDCommand(t, workspace, "issue", "create", "--kind", "epic", "--title", "Epic completion reject", "--scheduling", "sequential", "--by", "tl")
	issueID := parseCreatedIssueID(t, createStderr)

	journalPath := mustJournalPath(t, workspace)
	linesBefore := splitJournalLines(mustReadCommandFile(t, journalPath))

	stdout, stderr, err := runFIDCommandExpectError(t, workspace, "issue", "complete", issueID, "--by", "developer")
	if err == nil {
		t.Fatal("Execute() error = nil, want epic rejection")
	}
	if stdout != "" {
		t.Fatalf("stdout = %q, want empty", stdout)
	}
	if !strings.Contains(stderr, `error: complete issue: issue "`) || !strings.Contains(stderr, `task completions require task or subtask`) {
		t.Fatalf("stderr = %q, want epic kind rejection", stderr)
	}
	linesAfter := splitJournalLines(mustReadCommandFile(t, journalPath))
	if len(linesAfter) != len(linesBefore) {
		t.Fatalf("journal lines after epic rejection = %d, want unchanged %d", len(linesAfter), len(linesBefore))
	}

	conn, err := fiddb.Open(filepath.Join(workspace, ".fid", "fid.sqlite"))
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	defer closeCommandDB(t, conn)

	var count int
	if err := conn.QueryRow(`SELECT COUNT(*) FROM task_completions`).Scan(&count); err != nil {
		t.Fatalf("count task_completions: %v", err)
	}
	if count != 0 {
		t.Fatalf("task_completions rows = %d, want 0", count)
	}
}

func TestIssueCompleteCommandRejectsUnknownRoleAndMissingIssueBeforeAppend(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		args    func(t *testing.T, workspace string) []string
		wantErr string
		verify  func(t *testing.T, workspace string)
	}{
		{
			name: "unknown role",
			args: func(t *testing.T, workspace string) []string {
				_, createStderr := runFIDCommand(t, workspace, "issue", "create", "--kind", "task", "--title", "Unknown role reject", "--by", "developer")
				issueID := parseCreatedIssueID(t, createStderr)
				return []string{"issue", "complete", issueID, "--by", "intern"}
			},
			wantErr: `error: role "intern" not in project profile`,
			verify: func(t *testing.T, workspace string) {
				conn, err := fiddb.Open(filepath.Join(workspace, ".fid", "fid.sqlite"))
				if err != nil {
					t.Fatalf("Open() error = %v", err)
				}
				defer closeCommandDB(t, conn)

				var count int
				if err := conn.QueryRow(`SELECT COUNT(*) FROM task_completions`).Scan(&count); err != nil {
					t.Fatalf("count task_completions: %v", err)
				}
				if count != 0 {
					t.Fatalf("task_completions rows = %d, want 0", count)
				}
			},
		},
		{
			name: "missing issue",
			args: func(t *testing.T, workspace string) []string {
				return []string{"issue", "complete", "FID-missing", "--by", "developer"}
			},
			wantErr: `error: issue "FID-missing" not found`,
			verify: func(t *testing.T, workspace string) {
				conn, err := fiddb.Open(filepath.Join(workspace, ".fid", "fid.sqlite"))
				if err != nil {
					t.Fatalf("Open() error = %v", err)
				}
				defer closeCommandDB(t, conn)

				var count int
				if err := conn.QueryRow(`SELECT COUNT(*) FROM task_completions`).Scan(&count); err != nil {
					t.Fatalf("count task_completions: %v", err)
				}
				if count != 0 {
					t.Fatalf("task_completions rows = %d, want 0", count)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			workspace := t.TempDir()
			clonePristineWorkspaceInto(t, workspace)

			args := tt.args(t, workspace)
			journalPath := mustJournalPath(t, workspace)
			linesBefore := splitJournalLines(mustReadCommandFile(t, journalPath))

			stdout, stderr, err := runFIDCommandExpectError(t, workspace, args...)
			if err == nil {
				t.Fatal("Execute() error = nil, want rejection")
			}
			if stdout != "" {
				t.Fatalf("stdout = %q, want empty", stdout)
			}
			if !strings.Contains(stderr, tt.wantErr) {
				t.Fatalf("stderr = %q, want substring %q", stderr, tt.wantErr)
			}

			linesAfter := splitJournalLines(mustReadCommandFile(t, journalPath))
			if len(linesAfter) != len(linesBefore) {
				t.Fatalf("journal lines after rejection = %d, want unchanged %d", len(linesAfter), len(linesBefore))
			}

			if tt.verify != nil {
				tt.verify(t, workspace)
			}
		})
	}
}
