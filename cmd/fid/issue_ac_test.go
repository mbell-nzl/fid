package main

import (
	"encoding/json"
	"path/filepath"
	"strings"
	"testing"

	fiddb "github.com/mbell-nzl/fid/internal/db"
)

func TestIssueACAddCommandAddsAcceptanceCriterionAndPrintsConfirmation(t *testing.T) {
	workspace := t.TempDir()
	runInitCommand(t, workspace, "init")

	_, createStderr := runFIDCommand(t, workspace, "issue", "create", "--kind", "task", "--title", "CLI AC add", "--by", "developer")
	issueID := parseCreatedIssueID(t, createStderr)

	stdout, stderr := runFIDCommand(
		t,
		workspace,
		"issue", "ac", "add", issueID,
		"--text", "Journal file is flushed to disk before returning success",
		"--ordinal", "1",
		"--by", "developer",
	)
	if stdout != "" {
		t.Fatalf("stdout = %q, want empty", stdout)
	}
	if !strings.Contains(stderr, "added AC AC-") {
		t.Fatalf("stderr = %q, want AC confirmation prefix", stderr)
	}
	if !strings.Contains(stderr, "on "+issueID+" (ordinal 1)") {
		t.Fatalf("stderr = %q, want issue id and ordinal", stderr)
	}
	if !strings.Contains(stderr, `"Journal file is flushed to disk before returning success"`) {
		t.Fatalf("stderr = %q, want quoted criterion text", stderr)
	}

	getStdout, getStderr := runFIDCommand(t, workspace, "issue", "get", issueID, "--section", "acceptance", "--format", "json")
	if getStderr != "" {
		t.Fatalf("issue get stderr = %q, want empty", getStderr)
	}

	var got []struct {
		ID      string `json:"id"`
		Ordinal int    `json:"ordinal"`
		Text    string `json:"text"`
	}
	if err := json.Unmarshal([]byte(getStdout), &got); err != nil {
		t.Fatalf("Unmarshal() error = %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("acceptance criteria len = %d, want 1", len(got))
	}
	if !strings.HasPrefix(got[0].ID, "AC-") || got[0].Ordinal != 1 || got[0].Text != "Journal file is flushed to disk before returning success" {
		t.Fatalf("acceptance criterion = %#v, want id prefix/text/ordinal", got[0])
	}
}

func TestIssueACAddCommandRejectsBeforeAppend(t *testing.T) {
	tests := []struct {
		name    string
		args    func(t *testing.T, workspace string) []string
		wantErr string
		verify  func(t *testing.T, workspace string)
	}{
		{
			name: "missing by flag",
			args: func(t *testing.T, workspace string) []string {
				_, createStderr := runFIDCommand(t, workspace, "issue", "create", "--kind", "task", "--title", "Missing by", "--by", "developer")
				issueID := parseCreatedIssueID(t, createStderr)
				return []string{"issue", "ac", "add", issueID, "--text", "Missing role", "--ordinal", "1"}
			},
			wantErr: "error: --by is required",
		},
		{
			name: "missing text flag",
			args: func(t *testing.T, workspace string) []string {
				_, createStderr := runFIDCommand(t, workspace, "issue", "create", "--kind", "task", "--title", "Missing text", "--by", "developer")
				issueID := parseCreatedIssueID(t, createStderr)
				return []string{"issue", "ac", "add", issueID, "--ordinal", "1", "--by", "developer"}
			},
			wantErr: "error: --text is required",
		},
		{
			name: "missing ordinal flag",
			args: func(t *testing.T, workspace string) []string {
				_, createStderr := runFIDCommand(t, workspace, "issue", "create", "--kind", "task", "--title", "Missing ordinal", "--by", "developer")
				issueID := parseCreatedIssueID(t, createStderr)
				return []string{"issue", "ac", "add", issueID, "--text", "Missing ordinal", "--by", "developer"}
			},
			wantErr: "error: --ordinal is required",
		},
		{
			name: "invalid ordinal",
			args: func(t *testing.T, workspace string) []string {
				_, createStderr := runFIDCommand(t, workspace, "issue", "create", "--kind", "task", "--title", "Bad ordinal", "--by", "developer")
				issueID := parseCreatedIssueID(t, createStderr)
				return []string{"issue", "ac", "add", issueID, "--text", "Bad ordinal", "--ordinal", "0", "--by", "developer"}
			},
			wantErr: "error: add acceptance criterion: ordinal must be a positive integer",
		},
		{
			name: "unknown role",
			args: func(t *testing.T, workspace string) []string {
				_, createStderr := runFIDCommand(t, workspace, "issue", "create", "--kind", "task", "--title", "Unknown role", "--by", "developer")
				issueID := parseCreatedIssueID(t, createStderr)
				return []string{"issue", "ac", "add", issueID, "--text", "Unknown role", "--ordinal", "1", "--by", "intern"}
			},
			wantErr: `error: role "intern" not in project profile`,
		},
		{
			name: "missing issue",
			args: func(t *testing.T, workspace string) []string {
				return []string{"issue", "ac", "add", "FID-missing", "--text", "Missing issue", "--ordinal", "1", "--by", "developer"}
			},
			wantErr: `error: issue "FID-missing" not found`,
		},
		{
			name: "epic rejection before append",
			args: func(t *testing.T, workspace string) []string {
				_, createStderr := runFIDCommand(t, workspace, "issue", "create", "--kind", "epic", "--title", "Epic AC reject", "--scheduling", "prioritised", "--by", "tl")
				issueID := parseCreatedIssueID(t, createStderr)
				return []string{"issue", "ac", "add", issueID, "--text", "Epic AC", "--ordinal", "1", "--by", "developer"}
			},
			wantErr: `error: add acceptance criterion: issue "`,
			verify: func(t *testing.T, workspace string) {
				conn, err := fiddb.Open(filepath.Join(workspace, ".fid", "fid.sqlite"))
				if err != nil {
					t.Fatalf("Open() error = %v", err)
				}
				defer closeCommandDB(t, conn)

				var count int
				if err := conn.QueryRow(`SELECT COUNT(*) FROM acceptance_criteria`).Scan(&count); err != nil {
					t.Fatalf("count acceptance_criteria: %v", err)
				}
				if count != 0 {
					t.Fatalf("acceptance_criteria rows = %d, want 0", count)
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
