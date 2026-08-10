package main

import (
	"database/sql"
	"encoding/json"
	"path/filepath"
	"strings"
	"testing"

	fiddb "github.com/mbell-nzl/fid/internal/db"
)

func TestIssueGetSectionsAndDueDateReadBack(t *testing.T) {
	t.Run("json read-back exposes populated and omitted due dates", func(t *testing.T) {
		workspace := t.TempDir()
		clonePristineWorkspaceInto(t, workspace)

		_, withDueDateStderr := runFIDCommand(
			t,
			workspace,
			"issue", "create",
			"--kind", "task",
			"--title", "Due date read-back",
			"--due-date", "2026-06-30",
			"--by", "developer",
		)
		withDueDateID := parseCreatedIssueID(t, withDueDateStderr)

		_, withoutDueDateStderr := runFIDCommand(
			t,
			workspace,
			"issue", "create",
			"--kind", "task",
			"--title", "No due date read-back",
			"--by", "developer",
		)
		withoutDueDateID := parseCreatedIssueID(t, withoutDueDateStderr)

		type issueJSON struct {
			DueDate *string `json:"due_date"`
		}

		withDueDateJSON := mustGetIssueJSON(t, workspace, withDueDateID)
		var withDueDate issueJSON
		if err := json.Unmarshal([]byte(withDueDateJSON), &withDueDate); err != nil {
			t.Fatalf("Unmarshal(with due date) error = %v", err)
		}
		if withDueDate.DueDate == nil || *withDueDate.DueDate != "2026-06-30" {
			t.Fatalf("with due date JSON = %#v, want due_date 2026-06-30", withDueDate)
		}

		withoutDueDateJSON := mustGetIssueJSON(t, workspace, withoutDueDateID)
		var withoutDueDate issueJSON
		if err := json.Unmarshal([]byte(withoutDueDateJSON), &withoutDueDate); err != nil {
			t.Fatalf("Unmarshal(without due date) error = %v", err)
		}
		if withoutDueDate.DueDate != nil {
			t.Fatalf("without due date JSON = %#v, want due_date null", withoutDueDate)
		}
	})

	t.Run("sections return only the requested structured data and stable empty states", func(t *testing.T) {
		workspace := t.TempDir()
		clonePristineWorkspaceInto(t, workspace)

		conn, err := fiddb.Open(filepath.Join(workspace, ".fid", "fid.sqlite"))
		if err != nil {
			t.Fatalf("Open() error = %v", err)
		}
		defer closeCommandDB(t, conn)

		seedIssueGetSectionFixture(t, conn)

		acceptanceStdout, acceptanceStderr := runFIDCommand(t, workspace, "issue", "get", "FID-a1b2c3", "--section", "acceptance", "--format", "json")
		if acceptanceStderr != "" {
			t.Fatalf("acceptance stderr = %q, want empty", acceptanceStderr)
		}
		var acceptance []struct {
			ID string `json:"id"`
		}
		if err := json.Unmarshal([]byte(acceptanceStdout), &acceptance); err != nil {
			t.Fatalf("Unmarshal(acceptance) error = %v", err)
		}
		if len(acceptance) != 2 || acceptance[0].ID != "AC-d4e5f6" {
			t.Fatalf("acceptance = %#v, want two acceptance criteria", acceptance)
		}

		concernsStdout, concernsStderr := runFIDCommand(t, workspace, "issue", "get", "FID-a1b2c3", "--section", "concerns", "--format", "json")
		if concernsStderr != "" {
			t.Fatalf("concerns stderr = %q, want empty", concernsStderr)
		}
		var concerns []struct {
			ID     string `json:"id"`
			Status string `json:"status"`
		}
		if err := json.Unmarshal([]byte(concernsStdout), &concerns); err != nil {
			t.Fatalf("Unmarshal(concerns) error = %v", err)
		}
		if len(concerns) != 1 || concerns[0].ID != "CON-000401" || concerns[0].Status != "open" {
			t.Fatalf("concerns = %#v, want one open concern", concerns)
		}

		completionsStdout, completionsStderr := runFIDCommand(t, workspace, "issue", "get", "FID-a1b2c3", "--section", "completions", "--format", "json")
		if completionsStderr != "" {
			t.Fatalf("completions stderr = %q, want empty", completionsStderr)
		}
		var completions []struct {
			ID            string `json:"id"`
			AttemptNumber int    `json:"attempt_number"`
		}
		if err := json.Unmarshal([]byte(completionsStdout), &completions); err != nil {
			t.Fatalf("Unmarshal(completions) error = %v", err)
		}
		if len(completions) != 2 || completions[0].ID != "TC-000402" || completions[0].AttemptNumber != 2 {
			t.Fatalf("completions = %#v, want latest completion first", completions)
		}

		linksStdout, linksStderr := runFIDCommand(t, workspace, "issue", "get", "FID-a1b2c3", "--section", "links", "--format", "json")
		if linksStderr != "" {
			t.Fatalf("links stderr = %q, want empty", linksStderr)
		}
		var links []struct {
			Relation    string  `json:"relation"`
			TargetType  string  `json:"target_type"`
			TargetID    string  `json:"target_id"`
			TargetTitle *string `json:"target_title"`
		}
		if err := json.Unmarshal([]byte(linksStdout), &links); err != nil {
			t.Fatalf("Unmarshal(links) error = %v", err)
		}
		if len(links) != 2 || links[0].Relation != "blocks" || links[0].TargetType != "issue" || links[0].TargetTitle == nil || *links[0].TargetTitle != "Schema migration runner" {
			t.Fatalf("links = %#v, want issue + requirement links", links)
		}

		claimsStdout, claimsStderr := runFIDCommand(t, workspace, "issue", "get", "FID-a1b2c3", "--section", "claims", "--format", "json")
		if claimsStderr != "" {
			t.Fatalf("claims stderr = %q, want empty", claimsStderr)
		}
		var claims []struct {
			ClaimedBy string  `json:"claimed_by"`
			ExpiresAt *string `json:"expires_at"`
		}
		if err := json.Unmarshal([]byte(claimsStdout), &claims); err != nil {
			t.Fatalf("Unmarshal(claims) error = %v", err)
		}
		if len(claims) != 1 || claims[0].ClaimedBy != "developer" || claims[0].ExpiresAt == nil {
			t.Fatalf("claims = %#v, want one developer claim with expires_at", claims)
		}

		wantGates := []any{
			map[string]any{
				"id":               "GATE-000401",
				"attempt_number":   1,
				"agent_run_id":     "RUN-000401",
				"recorded_by_role": "developer",
				"command":          "go test ./...",
				"exit_code":        0,
				"output_excerpt":   "ok",
				"ran_at":           "2026-05-27T09:20:00Z",
			},
			map[string]any{
				"id":               "GATE-000402",
				"attempt_number":   2,
				"agent_run_id":     nil,
				"recorded_by_role": "developer",
				"command":          "golangci-lint run ./...",
				"exit_code":        1,
				"output_excerpt":   nil,
				"ran_at":           "2026-05-27T10:25:00Z",
			},
		}

		for _, formatName := range []string{"json", "toon", "yaml"} {
			t.Run("gates_"+formatName, func(t *testing.T) {
				stdout, stderr := runFIDCommand(t, workspace, "issue", "get", "FID-a1b2c3", "--section", "gates", "--format", formatName)
				if stderr != "" {
					t.Fatalf("%s stderr = %q, want empty", formatName, stderr)
				}

				var got any
				switch formatName {
				case "json":
					got = decodeCommandJSONStructure(t, stdout)
				case "toon":
					got = decodeCommandTOONStructure(t, stdout)
				case "yaml":
					got = decodeCommandYAMLStructure(t, stdout)
				}
				if !commandStructuresEqual(wantGates, got) {
					t.Fatalf("%s gates structure = %#v, want %#v", formatName, got, wantGates)
				}
			})
		}

		markdownStdout, markdownStderr := runFIDCommand(t, workspace, "issue", "get", "FID-a1b2c3", "--section", "gates", "--format", "markdown")
		if markdownStderr != "" {
			t.Fatalf("gates markdown stderr = %q, want empty", markdownStderr)
		}
		for _, wantSubstring := range []string{"# Gates", "GATE-000401", "go test ./...", "golangci-lint run ./...", "RUN-000401"} {
			if !strings.Contains(markdownStdout, wantSubstring) {
				t.Fatalf("gates markdown stdout = %q, want substring %q", markdownStdout, wantSubstring)
			}
		}

		_, emptyIssueStderr := runFIDCommand(t, workspace, "issue", "create", "--kind", "task", "--title", "Empty sections", "--by", "developer")
		emptyIssueID := parseCreatedIssueID(t, emptyIssueStderr)

		for _, section := range []string{"acceptance", "concerns", "completions", "links", "claims", "gates"} {
			section := section
			t.Run("empty "+section, func(t *testing.T) {
				stdout, stderr := runFIDCommand(t, workspace, "issue", "get", emptyIssueID, "--section", section, "--format", "json")
				if stderr != "" {
					t.Fatalf("%s stderr = %q, want empty", section, stderr)
				}
				if strings.TrimSpace(stdout) != "[]" {
					t.Fatalf("%s stdout = %q, want []", section, stdout)
				}
			})
		}

		stdout, stderr, err := runFIDCommandExpectError(t, workspace, "issue", "get", "FID-a1b2c3", "--section", "unknown", "--format", "json")
		if err == nil {
			t.Fatal("Execute() error = nil, want unknown section rejection")
		}
		if stdout != "" {
			t.Fatalf("stdout = %q, want empty", stdout)
		}
		if !strings.Contains(stderr, `supported sections: acceptance, concerns, completions, links, claims, gates`) {
			t.Fatalf("stderr = %q, want supported section list", stderr)
		}
	})
}

func seedIssueGetSectionFixture(t *testing.T, conn *sql.DB) {
	t.Helper()

	seedMarkdownPopulatedIssue(t, conn)

	statements := []struct {
		query string
		args  []any
	}{
		{
			query: `
				INSERT INTO task_completions (
					id, issue_id, attempt_number, developer_role, model, status, actual_tokens, gate_results, human_intervention, completed_at
				) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?), (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
			`,
			args: []any{
				"TC-000401",
				"FID-a1b2c3",
				1,
				"developer",
				"gpt-5",
				"dev_complete",
				1440,
				`["go test ./..."]`,
				false,
				"2026-05-27T09:30:00Z",
				"TC-000402",
				"FID-a1b2c3",
				2,
				"developer",
				"gpt-5",
				"qa_pass",
				1550,
				`[{"cmd":"go test ./...","exit_code":0,"ran_at":"2026-05-27T10:10:00Z","output_excerpt":"ok"},{"cmd":"go test -race ./...","exit_code":0,"ran_at":"2026-05-27T10:12:00Z","output_excerpt":"race ok"}]`,
				false,
				"2026-05-27T10:20:00Z",
			},
		},
		{
			query: `
				INSERT INTO agent_runs (id, issue_id, role, model, started_at)
				VALUES (?, ?, ?, ?, ?)
			`,
			args: []any{
				"RUN-000401",
				"FID-a1b2c3",
				"developer",
				"gpt-5",
				"2026-05-27T09:00:00Z",
			},
		},
		{
			query: `
				INSERT INTO gate_runs (
					id, issue_id, attempt_number, agent_run_id, recorded_by_role, command, exit_code, output_excerpt, ran_at
				) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?), (?, ?, ?, ?, ?, ?, ?, ?, ?)
			`,
			args: []any{
				"GATE-000401",
				"FID-a1b2c3",
				1,
				"RUN-000401",
				"developer",
				"go test ./...",
				0,
				"ok",
				"2026-05-27T09:20:00Z",
				"GATE-000402",
				"FID-a1b2c3",
				2,
				nil,
				"developer",
				"golangci-lint run ./...",
				1,
				nil,
				"2026-05-27T10:25:00Z",
			},
		},
		{
			query: `
				INSERT INTO concerns (
					id, issue_id, raised_in_completion_id, raised_by_role, summary, body_md_path, status, created_at, updated_at
				) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
			`,
			args: []any{
				"CON-000401",
				"FID-a1b2c3",
				"TC-000402",
				"developer",
				"Need triage before closing",
				"project-docs/concerns/CON-000401.md",
				"open",
				"2026-05-27T10:18:00Z",
				"2026-05-27T10:18:00Z",
			},
		},
		{
			query: `
				INSERT INTO triage_events (
					id, epic_id, issue_id, triager_role, source, source_ref_id, verdict, reasoning_md, created_at
				) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
			`,
			args: []any{
				"TRG-000401",
				"FID-e1f2g3",
				"FID-a1b2c3",
				"qa_lead",
				"task_completion",
				"TC-000402",
				"BLOCK",
				"project-docs/triage/TRG-000401.md",
				"2026-05-27T10:19:00Z",
			},
		},
		{
			query: `
				INSERT INTO claims (ref_type, ref_id, claimed_by, claimed_at, ttl_seconds, note)
				VALUES (?, ?, ?, ?, ?, ?)
			`,
			args: []any{
				"issue",
				"FID-a1b2c3",
				"developer",
				"2026-05-27T09:50:00Z",
				1800,
				"Working the flush path",
			},
		},
	}

	for _, statement := range statements {
		if _, err := conn.Exec(statement.query, statement.args...); err != nil {
			t.Fatalf("Exec(%q) error = %v", statement.query, err)
		}
	}
}
