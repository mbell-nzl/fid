package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	fiddb "github.com/mbell-nzl/fid/internal/db"
	"github.com/mbell-nzl/fid/internal/journal"
	"gopkg.in/yaml.v3"
)

func TestIssueCreateThenGetReturnsJSONByDefaultAndWithFormatFlag(t *testing.T) {
	workspace := t.TempDir()
	runInitCommand(t, workspace, "init")

	_, stderr := runFIDCommand(t, workspace, "issue", "create", "--kind", "task", "--title", "Implement tracer bullet", "--by", "developer")
	issueID := parseCreatedIssueID(t, stderr)

	defaultStdout, _ := runFIDCommand(t, workspace, "issue", "get", issueID)
	jsonStdout, _ := runFIDCommand(t, workspace, "issue", "get", issueID, "--format", "json")

	if defaultStdout != jsonStdout {
		t.Fatalf("default stdout = %q, want same as --format json %q", defaultStdout, jsonStdout)
	}

	var got struct {
		ID                 string   `json:"id"`
		Kind               string   `json:"kind"`
		Title              string   `json:"title"`
		Status             string   `json:"status"`
		StatusCategory     string   `json:"status_category"`
		WorkflowID         string   `json:"workflow_id"`
		OwnerRole          string   `json:"owner_role"`
		Labels             []string `json:"labels"`
		InInbox            bool     `json:"in_inbox"`
		ConcernsRaised     bool     `json:"concerns_raised"`
		OpenConcernsCount  int      `json:"open_concerns_count"`
		TriageStatus       *string  `json:"triage_status"`
		AcceptanceCriteria []any    `json:"acceptance_criteria"`
	}
	if err := json.Unmarshal([]byte(defaultStdout), &got); err != nil {
		t.Fatalf("Unmarshal() error = %v", err)
	}

	if got.ID != issueID {
		t.Fatalf("id = %q, want %q", got.ID, issueID)
	}
	if got.Kind != "task" {
		t.Fatalf("kind = %q, want %q", got.Kind, "task")
	}
	if got.Title != "Implement tracer bullet" {
		t.Fatalf("title = %q, want %q", got.Title, "Implement tracer bullet")
	}
	if got.Status != "todo" {
		t.Fatalf("status = %q, want %q", got.Status, "todo")
	}
	if got.StatusCategory != "unstarted" {
		t.Fatalf("status_category = %q, want %q", got.StatusCategory, "unstarted")
	}
	if got.WorkflowID != "task-default" {
		t.Fatalf("workflow_id = %q, want %q", got.WorkflowID, "task-default")
	}
	if got.OwnerRole != "developer" {
		t.Fatalf("owner_role = %q, want %q", got.OwnerRole, "developer")
	}
	if got.Labels == nil || len(got.Labels) != 0 {
		t.Fatalf("labels = %#v, want empty array", got.Labels)
	}
	if got.InInbox {
		t.Fatal("in_inbox = true, want false")
	}
	if got.ConcernsRaised {
		t.Fatal("concerns_raised = true, want false")
	}
	if got.OpenConcernsCount != 0 {
		t.Fatalf("open_concerns_count = %d, want 0", got.OpenConcernsCount)
	}
	if got.TriageStatus != nil {
		t.Fatalf("triage_status = %#v, want nil", got.TriageStatus)
	}
	if got.AcceptanceCriteria == nil || len(got.AcceptanceCriteria) != 0 {
		t.Fatalf("acceptance_criteria = %#v, want empty array", got.AcceptanceCriteria)
	}
}

func TestIssueGetReturnsEnrichedJSONAndRebuildParity(t *testing.T) {
	workspace := t.TempDir()
	runInitCommand(t, workspace, "init")

	_, stderr := runFIDCommand(t, workspace, "issue", "create", "--kind", "task", "--title", "Surface enriched issue state", "--by", "developer")
	issueID := parseCreatedIssueID(t, stderr)

	appendCommandIssueEvent(t, workspace, journal.Event{
		EventSchema: journal.EventSchemaV1,
		EventType:   "issue.enriched_state_seeded",
		Actor: journal.Actor{
			Role: "qa_lead",
		},
		Command: journal.Command{
			Argv:       []string{"fid", "maintenance", "seed"},
			CWD:        ".",
			FIDVersion: "0.1.0",
		},
		Effects: []journal.Effect{
			{
				Op:    "insert",
				Table: "issues",
				PK: map[string]any{
					"id": "FID-epiccmd",
				},
				Values: map[string]any{
					"kind":        "epic",
					"title":       "CLI enriched epic",
					"workflow_id": "epic-default",
					"status":      "todo",
					"created_at":  "2026-06-07T14:00:00Z",
					"updated_at":  "2026-06-07T14:00:00Z",
				},
			},
			{
				Op:    "insert",
				Table: "task_completions",
				PK: map[string]any{
					"id": "TC-CMD-001",
				},
				Values: map[string]any{
					"issue_id":           issueID,
					"attempt_number":     1,
					"developer_role":     "developer",
					"model":              "gpt-5",
					"status":             "dev_complete",
					"gate_results":       `["go test ./..."]`,
					"human_intervention": false,
					"completed_at":       "2026-06-07T14:01:00Z",
				},
			},
			{
				Op:    "insert",
				Table: "concerns",
				PK: map[string]any{
					"id": "CON-CMD-001",
				},
				Values: map[string]any{
					"issue_id":                issueID,
					"raised_in_completion_id": "TC-CMD-001",
					"raised_by_role":          "developer",
					"summary":                 "CLI concern",
					"body_md_path":            "project-docs/concerns/CON-CMD-001.md",
					"status":                  "open",
					"created_at":              "2026-06-07T14:01:30Z",
					"updated_at":              "2026-06-07T14:01:30Z",
				},
			},
			{
				Op:    "insert",
				Table: "triage_events",
				PK: map[string]any{
					"id": "TRG-CMD-001",
				},
				Values: map[string]any{
					"epic_id":       "FID-epiccmd",
					"issue_id":      issueID,
					"triager_role":  "qa_lead",
					"source":        "task_completion",
					"source_ref_id": "TC-CMD-001",
					"verdict":       "RESUME",
					"reasoning_md":  "project-docs/triage/TRG-CMD-001.md",
					"created_at":    "2026-06-07T14:02:00Z",
				},
			},
			{
				Op:    "insert",
				Table: "triage_events",
				PK: map[string]any{
					"id": "TRG-CMD-002",
				},
				Values: map[string]any{
					"epic_id":       "FID-epiccmd",
					"issue_id":      issueID,
					"triager_role":  "qa_lead",
					"source":        "task_completion",
					"source_ref_id": "TC-CMD-001",
					"verdict":       "BLOCK",
					"reasoning_md":  "project-docs/triage/TRG-CMD-002.md",
					"created_at":    "2026-06-07T14:03:00Z",
				},
			},
		},
	})

	beforeStdout, _ := runFIDCommand(t, workspace, "issue", "get", issueID, "--format", "json")

	var before struct {
		ConcernsRaised    bool    `json:"concerns_raised"`
		OpenConcernsCount int     `json:"open_concerns_count"`
		TriageStatus      *string `json:"triage_status"`
	}
	if err := json.Unmarshal([]byte(beforeStdout), &before); err != nil {
		t.Fatalf("Unmarshal(before) error = %v", err)
	}
	if !before.ConcernsRaised {
		t.Fatal("before concerns_raised = false, want true")
	}
	if before.OpenConcernsCount != 1 {
		t.Fatalf("before open_concerns_count = %d, want 1", before.OpenConcernsCount)
	}
	if before.TriageStatus == nil || *before.TriageStatus != "blocked" {
		t.Fatalf("before triage_status = %#v, want blocked", before.TriageStatus)
	}

	dbPath := filepath.Join(workspace, ".fid", "fid.sqlite")
	if err := os.Remove(dbPath); err != nil {
		t.Fatalf("Remove(%q) error = %v", dbPath, err)
	}

	runFIDCommand(t, workspace, "rebuild", "--by", "developer")

	afterStdout, _ := runFIDCommand(t, workspace, "issue", "get", issueID, "--format", "json")
	if beforeStdout != afterStdout {
		t.Fatalf("issue JSON after rebuild = %q, want %q", afterStdout, beforeStdout)
	}
}

func TestIssueGetReturnsMarkdownForEmptyAndPopulatedIssues(t *testing.T) {
	t.Run("empty issue omits empty sections", func(t *testing.T) {
		workspace := t.TempDir()
		runInitCommand(t, workspace, "init")

		conn, err := fiddb.Open(filepath.Join(workspace, ".fid", "fid.sqlite"))
		if err != nil {
			t.Fatalf("Open() error = %v", err)
		}
		defer closeCommandDB(t, conn)

		if _, err := conn.Exec(`
			INSERT INTO issues (id, kind, title, status, workflow_id, owner_role, created_at, updated_at)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?)
		`,
			"FID-emptymd",
			"task",
			"Render empty issue",
			"todo",
			"task-default",
			"developer",
			"2026-06-07T14:00:00Z",
			"2026-06-07T14:00:00Z",
		); err != nil {
			t.Fatalf("insert empty markdown issue: %v", err)
		}

		stdout, stderr := runFIDCommand(t, workspace, "issue", "get", "FID-emptymd", "--format", "markdown")
		if stderr != "" {
			t.Fatalf("stderr = %q, want empty", stderr)
		}

		want := "" +
			"# FID-emptymd: Render empty issue\n\n" +
			"**Kind:** task | **Status:** todo | **Priority:** (none)\n" +
			"**Owner:** developer\n\n" +
			"---\n" +
			"Created: 2026-06-07T14:00:00Z | Updated: 2026-06-07T14:00:00Z\n"
		if stdout != want {
			t.Fatalf("markdown stdout = %q, want %q", stdout, want)
		}
		if strings.Contains(stdout, "## Acceptance Criteria") || strings.Contains(stdout, "## Labels") || strings.Contains(stdout, "## Links") {
			t.Fatalf("markdown stdout = %q, want empty sections omitted", stdout)
		}
	})

	t.Run("populated issue matches markdown layout", func(t *testing.T) {
		workspace := t.TempDir()
		runInitCommand(t, workspace, "init")

		conn, err := fiddb.Open(filepath.Join(workspace, ".fid", "fid.sqlite"))
		if err != nil {
			t.Fatalf("Open() error = %v", err)
		}
		defer closeCommandDB(t, conn)

		seedMarkdownPopulatedIssue(t, conn)

		stdout, stderr := runFIDCommand(t, workspace, "issue", "get", "FID-a1b2c3", "--format", "markdown")
		if stderr != "" {
			t.Fatalf("stderr = %q, want empty", stderr)
		}

		want := "" +
			"# FID-a1b2c3: Implement journal append with flush\n\n" +
			"**Kind:** task | **Status:** in_progress | **Priority:** P1\n" +
			"**Epic:** FID-e1f2g3 (Core persistence layer)\n" +
			"**Owner:** developer | **Sequence:** 3\n" +
			"**Spec:** project-docs/docs/phases/phase-1/tasks/FID-a1b2c3.md\n\n" +
			"## Acceptance Criteria\n\n" +
			"- [x] AC-d4e5f6: Journal file is flushed to disk before returning success\n" +
			"  - developer: pass (2026-05-27T10:15:00Z)\n" +
			"- [ ] AC-e5f6g7: Partial write leaves journal intact\n" +
			"  - developer: (not verified)\n\n" +
			"## Labels\n\n" +
			"core, journal\n\n" +
			"## Links\n\n" +
			"- blocks FID-7c3e9f (Schema migration runner)\n" +
			"- implements PRD-ISSUE-001\n\n" +
			"---\n" +
			"Created: 2026-05-25T09:00:00Z | Updated: 2026-05-27T10:20:00Z\n"
		if stdout != want {
			t.Fatalf("markdown stdout = %q, want %q", stdout, want)
		}
	})
}

func TestIssueCreateAppendsOneEventAndProjectsOneIssue(t *testing.T) {
	workspace := t.TempDir()
	runInitCommand(t, workspace, "init")

	journalPath := mustJournalPath(t, workspace)
	beforePayload, err := os.ReadFile(journalPath)
	if err != nil {
		t.Fatalf("ReadFile(%q) error = %v", journalPath, err)
	}
	beforeLines := splitJournalLines(beforePayload)

	_, stderr := runFIDCommand(t, workspace, "issue", "create", "--kind", "task", "--title", "Append one event", "--by", "developer")
	issueID := parseCreatedIssueID(t, stderr)

	afterPayload, err := os.ReadFile(journalPath)
	if err != nil {
		t.Fatalf("ReadFile(%q) after create error = %v", journalPath, err)
	}
	afterLines := splitJournalLines(afterPayload)
	if len(afterLines) != len(beforeLines)+1 {
		t.Fatalf("journal lines = %d, want %d", len(afterLines), len(beforeLines)+1)
	}

	var event journal.Event
	if err := json.Unmarshal([]byte(afterLines[len(afterLines)-1]), &event); err != nil {
		t.Fatalf("Unmarshal() event error = %v", err)
	}
	if event.EventType != "issue.created" {
		t.Fatalf("event_type = %q, want %q", event.EventType, "issue.created")
	}
	if event.EventID == "" || event.EventSchema != journal.EventSchemaV1 || event.Actor.Role != "developer" || len(event.Command.Argv) == 0 || len(event.Effects) != 1 {
		t.Fatalf("event envelope = %+v, want populated v1 envelope with one effect", event)
	}

	conn, err := fiddb.Open(filepath.Join(workspace, ".fid", "fid.sqlite"))
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	defer closeCommandDB(t, conn)

	var projectedTitle string
	if err := conn.QueryRow(`SELECT title FROM issues WHERE id = ?`, issueID).Scan(&projectedTitle); err != nil {
		t.Fatalf("query created issue: %v", err)
	}
	if projectedTitle != "Append one event" {
		t.Fatalf("projected title = %q, want %q", projectedTitle, "Append one event")
	}
}

func TestIssueCreateAndRebuildRoundTripReturnsIdenticalJSON(t *testing.T) {
	workspace := t.TempDir()
	runInitCommand(t, workspace, "init")

	_, stderr := runFIDCommand(t, workspace, "issue", "create", "--kind", "task", "--title", "Rebuild parity", "--by", "developer")
	issueID := parseCreatedIssueID(t, stderr)

	beforeStdout, _ := runFIDCommand(t, workspace, "issue", "get", issueID)

	dbPath := filepath.Join(workspace, ".fid", "fid.sqlite")
	if err := os.Remove(dbPath); err != nil {
		t.Fatalf("Remove(%q) error = %v", dbPath, err)
	}

	runFIDCommand(t, workspace, "rebuild", "--by", "developer")

	afterStdout, _ := runFIDCommand(t, workspace, "issue", "get", issueID)
	if beforeStdout != afterStdout {
		t.Fatalf("issue JSON after rebuild = %q, want %q", afterStdout, beforeStdout)
	}
}

func TestIssueCreateAndRebuildValidateRoleExistence(t *testing.T) {
	workspace := t.TempDir()
	runInitCommand(t, workspace, "init")

	t.Run("issue create missing by", func(t *testing.T) {
		stdout, stderr, err := runFIDCommandExpectError(t, workspace, "issue", "create", "--kind", "task", "--title", "Missing role")
		if err == nil {
			t.Fatal("Execute() error = nil, want missing --by rejection")
		}
		if stdout != "" {
			t.Fatalf("stdout = %q, want empty", stdout)
		}
		if !strings.Contains(stderr, "error: --by is required") {
			t.Fatalf("stderr = %q, want missing --by error", stderr)
		}
	})

	t.Run("issue create unknown role", func(t *testing.T) {
		stdout, stderr, err := runFIDCommandExpectError(t, workspace, "issue", "create", "--kind", "task", "--title", "Unknown role", "--by", "intern")
		if err == nil {
			t.Fatal("Execute() error = nil, want unknown role rejection")
		}
		if stdout != "" {
			t.Fatalf("stdout = %q, want empty", stdout)
		}
		if !strings.Contains(stderr, `error: role "intern" not in project profile`) {
			t.Fatalf("stderr = %q, want unknown role error", stderr)
		}
	})

	t.Run("rebuild missing by", func(t *testing.T) {
		stdout, stderr, err := runFIDCommandExpectError(t, workspace, "rebuild")
		if err == nil {
			t.Fatal("Execute() error = nil, want missing --by rejection")
		}
		if stdout != "" {
			t.Fatalf("stdout = %q, want empty", stdout)
		}
		if !strings.Contains(stderr, "error: --by is required") {
			t.Fatalf("stderr = %q, want missing --by error", stderr)
		}
	})

	t.Run("rebuild unknown role", func(t *testing.T) {
		stdout, stderr, err := runFIDCommandExpectError(t, workspace, "rebuild", "--by", "intern")
		if err == nil {
			t.Fatal("Execute() error = nil, want unknown role rejection")
		}
		if stdout != "" {
			t.Fatalf("stdout = %q, want empty", stdout)
		}
		if !strings.Contains(stderr, `error: role "intern" not in project profile`) {
			t.Fatalf("stderr = %q, want unknown role error", stderr)
		}
	})
}

func TestIssueGetReturnsYAMLForEmptyAndPopulatedIssues(t *testing.T) {
	t.Run("empty issue", func(t *testing.T) {
		workspace := t.TempDir()
		runInitCommand(t, workspace, "init")

		_, stderr := runFIDCommand(t, workspace, "issue", "create", "--kind", "task", "--title", "YAML empty issue", "--by", "developer")
		issueID := parseCreatedIssueID(t, stderr)

		jsonStdout, jsonStderr := runFIDCommand(t, workspace, "issue", "get", issueID, "--format", "json")
		yamlStdout, yamlStderr := runFIDCommand(t, workspace, "issue", "get", issueID, "--format", "yaml")

		if jsonStderr != "" {
			t.Fatalf("json stderr = %q, want empty", jsonStderr)
		}
		if yamlStderr != "" {
			t.Fatalf("yaml stderr = %q, want empty", yamlStderr)
		}

		want := decodeCommandJSONStructure(t, jsonStdout)
		got := decodeCommandYAMLStructure(t, yamlStdout)
		if !commandStructuresEqual(want, got) {
			t.Fatalf("YAML structure = %#v, want %#v", got, want)
		}
	})

	t.Run("populated issue", func(t *testing.T) {
		workspace := t.TempDir()
		runInitCommand(t, workspace, "init")

		conn, err := fiddb.Open(filepath.Join(workspace, ".fid", "fid.sqlite"))
		if err != nil {
			t.Fatalf("Open() error = %v", err)
		}
		defer closeCommandDB(t, conn)

		seedMarkdownPopulatedIssue(t, conn)

		jsonStdout, jsonStderr := runFIDCommand(t, workspace, "issue", "get", "FID-a1b2c3", "--format", "json")
		yamlStdout, yamlStderr := runFIDCommand(t, workspace, "issue", "get", "FID-a1b2c3", "--format", "yaml")

		if jsonStderr != "" {
			t.Fatalf("json stderr = %q, want empty", jsonStderr)
		}
		if yamlStderr != "" {
			t.Fatalf("yaml stderr = %q, want empty", yamlStderr)
		}

		want := decodeCommandJSONStructure(t, jsonStdout)
		got := decodeCommandYAMLStructure(t, yamlStdout)
		if !commandStructuresEqual(want, got) {
			t.Fatalf("YAML structure = %#v, want %#v", got, want)
		}
	})
}

func TestIssueGetDispatchesStructuredFormatsAndRejectsUnknownFormats(t *testing.T) {
	workspace := t.TempDir()
	runInitCommand(t, workspace, "init")

	_, stderr := runFIDCommand(t, workspace, "issue", "create", "--kind", "task", "--title", "Unsupported format", "--by", "developer")
	issueID := parseCreatedIssueID(t, stderr)

	for _, formatName := range []string{"toon", "yaml"} {
		t.Run(formatName, func(t *testing.T) {
			jsonStdout, _ := runFIDCommand(t, workspace, "issue", "get", issueID, "--format", "json")
			formattedStdout, formattedStderr := runFIDCommand(t, workspace, "issue", "get", issueID, "--format", formatName)

			if formattedStderr != "" {
				t.Fatalf("stderr = %q, want empty", formattedStderr)
			}

			want := decodeCommandJSONStructure(t, jsonStdout)
			var got any
			switch formatName {
			case "toon":
				got = decodeCommandTOONStructure(t, formattedStdout)
			case "yaml":
				got = decodeCommandYAMLStructure(t, formattedStdout)
			default:
				t.Fatalf("unexpected format %q", formatName)
			}

			if !commandStructuresEqual(want, got) {
				t.Fatalf("%s structure = %#v, want %#v", formatName, got, want)
			}
		})
	}

	t.Run("xml", func(t *testing.T) {
		stdout, stderr, err := runFIDCommandExpectError(t, workspace, "issue", "get", issueID, "--format", "xml")
		if err == nil {
			t.Fatal("Execute() error = nil, want unsupported format rejection")
		}
		if stdout != "" {
			t.Fatalf("stdout = %q, want empty", stdout)
		}
		if !strings.Contains(stderr, `unsupported format "xml"; known formats: json, toon, markdown, yaml`) {
			t.Fatalf("stderr = %q, want actionable unknown format error", stderr)
		}
	})
}

func seedMarkdownPopulatedIssue(t *testing.T, conn *sql.DB) {
	t.Helper()

	statements := []struct {
		query string
		args  []any
	}{
		{
			query: `
				INSERT INTO issues (id, kind, title, status, workflow_id, owner_role, created_at, updated_at)
				VALUES (?, ?, ?, ?, ?, ?, ?, ?)
			`,
			args: []any{
				"FID-e1f2g3",
				"epic",
				"Core persistence layer",
				"todo",
				"epic-default",
				"developer",
				"2026-05-25T08:00:00Z",
				"2026-05-25T08:00:00Z",
			},
		},
		{
			query: `
				INSERT INTO issues (id, kind, parent_id, title, workflow_id, status, priority, sequence, spec_md_path, owner_role, created_at, updated_at)
				VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
			`,
			args: []any{
				"FID-a1b2c3",
				"task",
				"FID-e1f2g3",
				"Implement journal append with flush",
				"task-default",
				"in_progress",
				1,
				3,
				"project-docs/docs/phases/phase-1/tasks/FID-a1b2c3.md",
				"developer",
				"2026-05-25T09:00:00Z",
				"2026-05-27T10:20:00Z",
			},
		},
		{
			query: `
				INSERT INTO issues (id, kind, title, status, workflow_id, owner_role, created_at, updated_at)
				VALUES (?, ?, ?, ?, ?, ?, ?, ?)
			`,
			args: []any{
				"FID-7c3e9f",
				"task",
				"Schema migration runner",
				"todo",
				"task-default",
				"developer",
				"2026-05-25T09:05:00Z",
				"2026-05-25T09:05:00Z",
			},
		},
		{
			query: `
				INSERT INTO requirements (id, type, summary, source_doc, source_section, status)
				VALUES (?, ?, ?, ?, ?, ?)
			`,
			args: []any{
				"PRD-ISSUE-001",
				"PRD",
				"Implement journal append with flush",
				"01_PRD.md",
				"ISSUE",
				"active",
			},
		},
		{
			query: `INSERT INTO labels (name) VALUES (?), (?)`,
			args:  []any{"core", "journal"},
		},
		{
			query: `INSERT INTO issue_labels (issue_id, label) VALUES (?, ?), (?, ?)`,
			args:  []any{"FID-a1b2c3", "core", "FID-a1b2c3", "journal"},
		},
		{
			query: `
				INSERT INTO acceptance_criteria (id, issue_id, ordinal, text)
				VALUES (?, ?, ?, ?), (?, ?, ?, ?)
			`,
			args: []any{
				"AC-d4e5f6",
				"FID-a1b2c3",
				1,
				"Journal file is flushed to disk before returning success",
				"AC-e5f6g7",
				"FID-a1b2c3",
				2,
				"Partial write leaves journal intact",
			},
		},
		{
			query: `
				INSERT INTO ac_verifications (id, ac_id, verifier_role, status, verified_at)
				VALUES (?, ?, ?, ?, ?)
			`,
			args: []any{
				"ACV-d4e5f6-dev",
				"AC-d4e5f6",
				"developer",
				"pass",
				"2026-05-27T10:15:00Z",
			},
		},
		{
			query: `
				INSERT INTO links (from_type, from_id, relation, to_type, to_id, created_at)
				VALUES (?, ?, ?, ?, ?, ?), (?, ?, ?, ?, ?, ?)
			`,
			args: []any{
				"issue",
				"FID-a1b2c3",
				"blocks",
				"issue",
				"FID-7c3e9f",
				"2026-05-27T10:18:00Z",
				"issue",
				"FID-a1b2c3",
				"implements",
				"requirement",
				"PRD-ISSUE-001",
				"2026-05-27T10:19:00Z",
			},
		},
	}

	for _, statement := range statements {
		if _, err := conn.Exec(statement.query, statement.args...); err != nil {
			t.Fatalf("Exec(%q) error = %v", statement.query, err)
		}
	}
}

func runFIDCommand(t *testing.T, workspace string, args ...string) (string, string) {
	t.Helper()

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	cmd := newRootCommand()
	cmd.SetOut(stdout)
	cmd.SetErr(stderr)
	cmd.SetArgs(args)

	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Getwd() error = %v", err)
	}
	defer mustChdir(t, cwd)
	mustChdir(t, workspace)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute(%v) error = %v", args, err)
	}

	return stdout.String(), stderr.String()
}

func runFIDCommandExpectError(t *testing.T, workspace string, args ...string) (string, string, error) {
	t.Helper()

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	cmd := newRootCommand()
	cmd.SetOut(stdout)
	cmd.SetErr(stderr)
	cmd.SetArgs(args)

	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Getwd() error = %v", err)
	}
	defer mustChdir(t, cwd)
	mustChdir(t, workspace)

	err = cmd.Execute()
	return stdout.String(), stderr.String(), err
}

func parseCreatedIssueID(t *testing.T, stderr string) string {
	t.Helper()

	fields := strings.Fields(stderr)
	for index := range fields {
		if fields[index] == "issue" && index+1 < len(fields) {
			return fields[index+1]
		}
	}
	t.Fatalf("stderr = %q, want created issue line with id", stderr)
	return ""
}

func appendCommandIssueEvent(t *testing.T, workspace string, event journal.Event) {
	t.Helper()

	if err := journal.WithLock(workspace, journal.LockOwner{
		Command:   "fid maintenance seed",
		ActorRole: event.Actor.Role,
	}, func(lock *journal.WriteLock) error {
		_, err := journal.AppendEvent(lock, event)
		return err
	}); err != nil {
		t.Fatalf("AppendEvent() error = %v", err)
	}
}

func decodeCommandJSONStructure(t *testing.T, value string) any {
	t.Helper()

	decoder := json.NewDecoder(strings.NewReader(value))
	decoder.UseNumber()

	var result any
	if err := decoder.Decode(&result); err != nil {
		t.Fatalf("Decode JSON structure error = %v", err)
	}
	return result
}

func decodeCommandTOONStructure(t *testing.T, value string) any {
	t.Helper()

	parser := commandTOONParser{input: strings.TrimSpace(value)}
	result, err := parser.parseValue()
	if err != nil {
		t.Fatalf("parse TOON structure error = %v", err)
	}
	parser.skipWhitespace()
	if parser.position != len(parser.input) {
		t.Fatalf("parse TOON structure stopped at %d of %d", parser.position, len(parser.input))
	}
	return result
}

func decodeCommandYAMLStructure(t *testing.T, value string) any {
	t.Helper()

	var result any
	if err := yaml.Unmarshal([]byte(value), &result); err != nil {
		t.Fatalf("Unmarshal YAML structure error = %v", err)
	}
	return result
}

func commandStructuresEqual(left any, right any) bool {
	switch leftValue := left.(type) {
	case map[string]any:
		rightValue, ok := right.(map[string]any)
		if !ok || len(leftValue) != len(rightValue) {
			return false
		}
		for key, nestedLeft := range leftValue {
			nestedRight, ok := rightValue[key]
			if !ok || !commandStructuresEqual(nestedLeft, nestedRight) {
				return false
			}
		}
		return true
	case []any:
		rightValue, ok := right.([]any)
		if !ok || len(leftValue) != len(rightValue) {
			return false
		}
		for index := range leftValue {
			if !commandStructuresEqual(leftValue[index], rightValue[index]) {
				return false
			}
		}
		return true
	case json.Number:
		return commandNumericStructuresEqual(leftValue, right)
	default:
		return fmt.Sprintf("%v", left) == fmt.Sprintf("%v", right)
	}
}

func commandNumericStructuresEqual(left json.Number, right any) bool {
	switch rightValue := right.(type) {
	case json.Number:
		return left.String() == rightValue.String()
	case int:
		return left.String() == strconv.Itoa(rightValue)
	case int8:
		return left.String() == strconv.FormatInt(int64(rightValue), 10)
	case int16:
		return left.String() == strconv.FormatInt(int64(rightValue), 10)
	case int32:
		return left.String() == strconv.FormatInt(int64(rightValue), 10)
	case int64:
		return left.String() == strconv.FormatInt(rightValue, 10)
	case uint:
		return left.String() == strconv.FormatUint(uint64(rightValue), 10)
	case uint8:
		return left.String() == strconv.FormatUint(uint64(rightValue), 10)
	case uint16:
		return left.String() == strconv.FormatUint(uint64(rightValue), 10)
	case uint32:
		return left.String() == strconv.FormatUint(uint64(rightValue), 10)
	case uint64:
		return left.String() == strconv.FormatUint(rightValue, 10)
	case float32:
		return left.String() == strconv.FormatFloat(float64(rightValue), 'f', -1, 32)
	case float64:
		return left.String() == strconv.FormatFloat(rightValue, 'f', -1, 64)
	default:
		return false
	}
}

type commandTOONParser struct {
	input    string
	position int
}

func (parser *commandTOONParser) parseValue() (any, error) {
	parser.skipWhitespace()
	if parser.position >= len(parser.input) {
		return nil, fmt.Errorf("unexpected end of input")
	}

	switch parser.input[parser.position] {
	case '{':
		return parser.parseObject()
	case '[':
		return parser.parseArray()
	case '"':
		return parser.parseString()
	default:
		return parser.parseScalar()
	}
}

func (parser *commandTOONParser) parseObject() (map[string]any, error) {
	parser.position++
	result := map[string]any{}

	parser.skipWhitespace()
	if parser.consumeIf('}') {
		return result, nil
	}

	for {
		key, err := parser.parseKey()
		if err != nil {
			return nil, err
		}
		if !parser.consumeIf(':') {
			return nil, fmt.Errorf("expected ':' after key %q", key)
		}
		value, err := parser.parseValue()
		if err != nil {
			return nil, err
		}
		result[key] = value

		parser.skipWhitespace()
		if parser.consumeIf('}') {
			return result, nil
		}
		if !parser.consumeIf(',') {
			return nil, fmt.Errorf("expected ',' or '}' after key %q", key)
		}
	}
}

func (parser *commandTOONParser) parseArray() ([]any, error) {
	parser.position++
	var result []any

	parser.skipWhitespace()
	if parser.consumeIf(']') {
		return result, nil
	}

	for {
		value, err := parser.parseValue()
		if err != nil {
			return nil, err
		}
		result = append(result, value)

		parser.skipWhitespace()
		if parser.consumeIf(']') {
			return result, nil
		}
		if !parser.consumeIf(',') {
			return nil, fmt.Errorf("expected ',' or ']' in array")
		}
	}
}

func (parser *commandTOONParser) parseKey() (string, error) {
	parser.skipWhitespace()
	if parser.position >= len(parser.input) {
		return "", fmt.Errorf("unexpected end of input")
	}
	if parser.input[parser.position] == '"' {
		return parser.parseString()
	}

	start := parser.position
	for parser.position < len(parser.input) {
		switch parser.input[parser.position] {
		case ':', ' ', '\n', '\t', '\r':
			return parser.input[start:parser.position], nil
		default:
			parser.position++
		}
	}
	return parser.input[start:parser.position], nil
}

func (parser *commandTOONParser) parseString() (string, error) {
	start := parser.position
	parser.position++
	escaped := false
	for parser.position < len(parser.input) {
		current := parser.input[parser.position]
		if escaped {
			escaped = false
			parser.position++
			continue
		}
		if current == '\\' {
			escaped = true
			parser.position++
			continue
		}
		if current == '"' {
			parser.position++
			var result string
			if err := json.Unmarshal([]byte(parser.input[start:parser.position]), &result); err != nil {
				return "", err
			}
			return result, nil
		}
		parser.position++
	}
	return "", fmt.Errorf("unterminated string")
}

func (parser *commandTOONParser) parseScalar() (any, error) {
	start := parser.position
	for parser.position < len(parser.input) {
		switch parser.input[parser.position] {
		case ',', '}', ']', ' ', '\n', '\t', '\r':
			token := parser.input[start:parser.position]
			return parseCommandTOONScalar(token)
		default:
			parser.position++
		}
	}
	return parseCommandTOONScalar(parser.input[start:parser.position])
}

func (parser *commandTOONParser) skipWhitespace() {
	for parser.position < len(parser.input) {
		switch parser.input[parser.position] {
		case ' ', '\n', '\t', '\r':
			parser.position++
		default:
			return
		}
	}
}

func (parser *commandTOONParser) consumeIf(expected byte) bool {
	parser.skipWhitespace()
	if parser.position < len(parser.input) && parser.input[parser.position] == expected {
		parser.position++
		return true
	}
	return false
}

func parseCommandTOONScalar(token string) (any, error) {
	switch token {
	case "null":
		return nil, nil
	case "true":
		return true, nil
	case "false":
		return false, nil
	}

	if _, err := strconv.ParseInt(token, 10, 64); err == nil {
		return json.Number(token), nil
	}
	if _, err := strconv.ParseFloat(token, 64); err == nil {
		return json.Number(token), nil
	}
	return token, nil
}
