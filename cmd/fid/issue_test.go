package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	fiddb "github.com/mbell-nzl/fid/internal/db"
	"github.com/mbell-nzl/fid/internal/journal"
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

func TestIssueGetRejectsUnsupportedFormats(t *testing.T) {
	workspace := t.TempDir()
	runInitCommand(t, workspace, "init")

	_, stderr := runFIDCommand(t, workspace, "issue", "create", "--kind", "task", "--title", "Unsupported format", "--by", "developer")
	issueID := parseCreatedIssueID(t, stderr)

	t.Run("toon", func(t *testing.T) {
		jsonStdout, _ := runFIDCommand(t, workspace, "issue", "get", issueID, "--format", "json")
		toonStdout, toonStderr := runFIDCommand(t, workspace, "issue", "get", issueID, "--format", "toon")

		if toonStderr != "" {
			t.Fatalf("stderr = %q, want empty", toonStderr)
		}

		want := decodeCommandJSONStructure(t, jsonStdout)
		got := decodeCommandTOONStructure(t, toonStdout)
		if !commandStructuresEqual(want, got) {
			t.Fatalf("TOON structure = %#v, want %#v", got, want)
		}
	})

	for _, formatName := range []string{"markdown", "yaml"} {
		t.Run(formatName, func(t *testing.T) {
			stdout, stderr, err := runFIDCommandExpectError(t, workspace, "issue", "get", issueID, "--format", formatName)
			if err == nil {
				t.Fatal("Execute() error = nil, want unsupported format rejection")
			}
			if stdout != "" {
				t.Fatalf("stdout = %q, want empty", stdout)
			}
			if !strings.Contains(stderr, `not yet supported`) {
				t.Fatalf("stderr = %q, want not yet supported error", stderr)
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
		rightValue, ok := right.(json.Number)
		return ok && leftValue.String() == rightValue.String()
	default:
		return fmt.Sprintf("%v", left) == fmt.Sprintf("%v", right)
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
