package main

import (
	"encoding/json"
	"path/filepath"
	"strings"
	"testing"

	fiddb "github.com/mbell-nzl/fid/internal/db"
)

func TestIssueCompleteCommandCreatesAttemptOneAndReturnsJSONCompletionFields(t *testing.T) {
	workspace := t.TempDir()
	runInitCommand(t, workspace, "init")

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

	jsonStdout, getStderr := runFIDCommand(t, workspace, "issue", "get", issueID, "--format", "json")
	if getStderr != "" {
		t.Fatalf("issue get stderr = %q, want empty", getStderr)
	}

	var got struct {
		Completions []struct {
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
		} `json:"completions"`
	}
	if err := json.Unmarshal([]byte(jsonStdout), &got); err != nil {
		t.Fatalf("Unmarshal() error = %v", err)
	}
	if len(got.Completions) != 1 {
		t.Fatalf("completions = %#v, want one completion", got.Completions)
	}
	completion := got.Completions[0]
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
}

func TestIssueCompleteCommandSupersedesPriorAttemptAndStoresInterventionFields(t *testing.T) {
	workspace := t.TempDir()
	runInitCommand(t, workspace, "init")

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

	jsonStdout, getStderr := runFIDCommand(t, workspace, "issue", "get", issueID, "--format", "json")
	if getStderr != "" {
		t.Fatalf("issue get stderr = %q, want empty", getStderr)
	}

	var got struct {
		Completions []struct {
			AttemptNumber         int     `json:"attempt_number"`
			Status                *string `json:"status"`
			ActualTokens          *int    `json:"actual_tokens"`
			HumanIntervention     bool    `json:"human_intervention"`
			HumanInterventionNote *string `json:"human_intervention_note"`
		} `json:"completions"`
	}
	if err := json.Unmarshal([]byte(jsonStdout), &got); err != nil {
		t.Fatalf("Unmarshal() error = %v", err)
	}
	if len(got.Completions) != 2 {
		t.Fatalf("completions len = %d, want 2", len(got.Completions))
	}
	if got.Completions[0].AttemptNumber != 2 {
		t.Fatalf("latest attempt_number = %d, want 2", got.Completions[0].AttemptNumber)
	}
	if got.Completions[0].Status == nil || *got.Completions[0].Status != "dev_complete" {
		t.Fatalf("latest status = %#v, want dev_complete", got.Completions[0].Status)
	}
	if got.Completions[0].ActualTokens == nil || *got.Completions[0].ActualTokens != 5000 {
		t.Fatalf("latest actual_tokens = %#v, want 5000", got.Completions[0].ActualTokens)
	}
	if !got.Completions[0].HumanIntervention {
		t.Fatal("latest human_intervention = false, want true")
	}
	if got.Completions[0].HumanInterventionNote == nil || *got.Completions[0].HumanInterventionNote != "paired on final review" {
		t.Fatalf("latest human_intervention_note = %#v, want paired on final review", got.Completions[0].HumanInterventionNote)
	}
	if got.Completions[1].AttemptNumber != 1 {
		t.Fatalf("prior attempt_number = %d, want 1", got.Completions[1].AttemptNumber)
	}
	if got.Completions[1].Status == nil || *got.Completions[1].Status != "superseded" {
		t.Fatalf("prior status = %#v, want superseded", got.Completions[1].Status)
	}
	if got.Completions[1].HumanIntervention {
		t.Fatal("prior human_intervention = true, want false")
	}
	if got.Completions[1].HumanInterventionNote != nil {
		t.Fatalf("prior human_intervention_note = %#v, want nil", got.Completions[1].HumanInterventionNote)
	}
}

func TestIssueCompleteCommandRejectsEpicBeforeAppend(t *testing.T) {
	workspace := t.TempDir()
	runInitCommand(t, workspace, "init")

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
			runInitCommand(t, workspace, "init")

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
