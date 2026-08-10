package main

import (
	"database/sql"
	"encoding/json"
	"path/filepath"
	"strings"
	"testing"

	fiddb "github.com/mbell-nzl/fid/internal/db"
)

func TestLabelAddCreatesOneRowInOneEventWithoutActivityLog(t *testing.T) {
	workspace := t.TempDir()
	clonePristineWorkspaceInto(t, workspace)

	journalPath := mustJournalPath(t, workspace)
	linesBefore := len(splitJournalLines(mustReadCommandFile(t, journalPath)))

	conn, err := fiddb.Open(filepath.Join(workspace, ".fid", "fid.sqlite"))
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	defer closeCommandDB(t, conn)

	activityBefore := countActivityLogRows(t, conn)

	stdout, stderr := runFIDCommand(t, workspace, "label", "add", "core", "--by", "developer")
	if stdout != "" {
		t.Fatalf("stdout = %q, want empty", stdout)
	}
	if !strings.Contains(stderr, `created label "core"`) {
		t.Fatalf("stderr = %q, want create confirmation", stderr)
	}

	linesAfter := len(splitJournalLines(mustReadCommandFile(t, journalPath)))
	if linesAfter != linesBefore+1 {
		t.Fatalf("journal lines after label add = %d, want %d", linesAfter, linesBefore+1)
	}

	var (
		name        string
		colour      sql.NullString
		description sql.NullString
	)
	if err := conn.QueryRow(`SELECT name, colour, description FROM labels WHERE name = ?`, "core").Scan(&name, &colour, &description); err != nil {
		t.Fatalf("query label row: %v", err)
	}
	if name != "core" {
		t.Fatalf("label name = %q, want %q", name, "core")
	}
	if colour.Valid {
		t.Fatalf("label colour = %q, want NULL", colour.String)
	}
	if description.Valid {
		t.Fatalf("label description = %q, want NULL", description.String)
	}

	if got := countActivityLogRows(t, conn); got != activityBefore {
		t.Fatalf("activity_log rows after label add = %d, want unchanged %d", got, activityBefore)
	}

	event := mustReadLastJournalEvent(t, journalPath)
	if event.EventType != "label.added" {
		t.Fatalf("event type = %q, want %q", event.EventType, "label.added")
	}
	if len(event.Effects) != 1 {
		t.Fatalf("effect count = %d, want 1", len(event.Effects))
	}
	effect := event.Effects[0]
	if effect.Op != "insert" {
		t.Fatalf("effect op = %q, want insert", effect.Op)
	}
	if effect.Table != "labels" {
		t.Fatalf("effect table = %q, want labels", effect.Table)
	}
}

func TestLabelListDefaultsToJSON(t *testing.T) {
	workspace := t.TempDir()
	clonePristineWorkspaceInto(t, workspace)

	runFIDCommand(t, workspace, "label", "add", "core", "--colour", "#4A90D9", "--description", "Core infrastructure", "--by", "developer")
	runFIDCommand(t, workspace, "label", "add", "tech-debt", "--by", "developer")

	defaultStdout, defaultStderr := runFIDCommand(t, workspace, "label", "list")
	if defaultStderr != "" {
		t.Fatalf("default stderr = %q, want empty", defaultStderr)
	}

	jsonStdout, jsonStderr := runFIDCommand(t, workspace, "label", "list", "--format", "json")
	if jsonStderr != "" {
		t.Fatalf("json stderr = %q, want empty", jsonStderr)
	}
	if defaultStdout != jsonStdout {
		t.Fatalf("default stdout = %q, want json stdout %q", defaultStdout, jsonStdout)
	}

	var rows []map[string]any
	if err := json.Unmarshal([]byte(jsonStdout), &rows); err != nil {
		t.Fatalf("Unmarshal(stdout) error = %v", err)
	}

	want := []map[string]any{
		{"name": "core", "colour": "#4A90D9", "description": "Core infrastructure"},
		{"name": "tech-debt", "colour": nil, "description": nil},
	}
	if !commandStructuresEqual(want, rows) {
		t.Fatalf("rows = %#v, want %#v", rows, want)
	}
}

func TestLabelListFormatsMatchJSONReference(t *testing.T) {
	workspace := t.TempDir()
	clonePristineWorkspaceInto(t, workspace)

	runFIDCommand(t, workspace, "label", "add", "core", "--colour", "#4A90D9", "--description", "Core infrastructure", "--by", "developer")
	runFIDCommand(t, workspace, "label", "add", "tech-debt", "--by", "developer")

	jsonStdout, jsonStderr := runFIDCommand(t, workspace, "label", "list", "--format", "json")
	if jsonStderr != "" {
		t.Fatalf("json stderr = %q, want empty", jsonStderr)
	}
	want := decodeCommandJSONStructure(t, jsonStdout)

	for _, formatName := range []string{"toon", "yaml", "markdown"} {
		t.Run(formatName, func(t *testing.T) {
			stdout, stderr := runFIDCommand(t, workspace, "label", "list", "--format", formatName)
			if stderr != "" {
				t.Fatalf("%s stderr = %q, want empty", formatName, stderr)
			}

			var got any
			switch formatName {
			case "toon":
				got = decodeCommandTOONStructure(t, stdout)
			case "yaml":
				got = decodeCommandYAMLStructure(t, stdout)
			case "markdown":
				got = decodeCommandLabelMarkdownRows(t, stdout)
			}

			if !commandStructuresEqual(want, got) {
				t.Fatalf("%s structure = %#v, want %#v", formatName, got, want)
			}
		})
	}
}

func TestLabelAddRejectsDuplicateDeniedAndUnknownBeforeAppend(t *testing.T) {
	t.Run("duplicate add", func(t *testing.T) {
		workspace := t.TempDir()
		clonePristineWorkspaceInto(t, workspace)
		runFIDCommand(t, workspace, "label", "add", "core", "--by", "developer")

		journalPath := mustJournalPath(t, workspace)
		linesBefore := len(splitJournalLines(mustReadCommandFile(t, journalPath)))

		stdout, stderr, err := runFIDCommandExpectError(t, workspace, "label", "add", "core", "--by", "developer")
		if err == nil {
			t.Fatal("duplicate Execute() error = nil, want rejection")
		}
		if stdout != "" {
			t.Fatalf("duplicate stdout = %q, want empty", stdout)
		}
		if !strings.Contains(stderr, `error: label add: label "core" already exists`) {
			t.Fatalf("duplicate stderr = %q, want duplicate-name error", stderr)
		}
		if got := len(splitJournalLines(mustReadCommandFile(t, journalPath))); got != linesBefore {
			t.Fatalf("journal lines after duplicate add = %d, want %d", got, linesBefore)
		}
	})

	t.Run("denied and unknown reject before append", func(t *testing.T) {
		workspace := t.TempDir()
		clonePristineWorkspaceInto(t, workspace)

		journalPath := mustJournalPath(t, workspace)
		linesBefore := len(splitJournalLines(mustReadCommandFile(t, journalPath)))

		stdout, stderr, err := runFIDCommandExpectError(t, workspace, "label", "add", "core", "--by", "qa_lead")
		if err == nil {
			t.Fatal("qa_lead Execute() error = nil, want denial")
		}
		if stdout != "" {
			t.Fatalf("qa_lead stdout = %q, want empty", stdout)
		}
		if !strings.Contains(stderr, `error: role "qa_lead" not permitted for action "label.add"`) {
			t.Fatalf("qa_lead stderr = %q, want permission denial", stderr)
		}
		if got := len(splitJournalLines(mustReadCommandFile(t, journalPath))); got != linesBefore {
			t.Fatalf("journal lines after denied role = %d, want %d", got, linesBefore)
		}

		stdout, stderr, err = runFIDCommandExpectError(t, workspace, "label", "add", "core", "--by", "intern")
		if err == nil {
			t.Fatal("unknown-role Execute() error = nil, want rejection")
		}
		if stdout != "" {
			t.Fatalf("unknown-role stdout = %q, want empty", stdout)
		}
		if !strings.Contains(stderr, `error: role "intern" not in project profile`) {
			t.Fatalf("unknown-role stderr = %q, want unknown-role error", stderr)
		}
		for _, want := range []string{
			"known roles:",
			"developer",
			"qa_reviewer",
			"qa_lead",
			"tl_lite",
			"tl",
			"workflow_runner",
			"human",
		} {
			if !strings.Contains(stderr, want) {
				t.Fatalf("unknown-role stderr = %q, want substring %q", stderr, want)
			}
		}
		if got := len(splitJournalLines(mustReadCommandFile(t, journalPath))); got != linesBefore {
			t.Fatalf("journal lines after unknown role = %d, want %d", got, linesBefore)
		}
	})
}

func TestLabelAddPolicyShowIncludesMatrixRow(t *testing.T) {
	workspace := t.TempDir()
	clonePristineWorkspaceInto(t, workspace)

	stdout, stderr := runFIDCommand(t, workspace, "policy", "show", "--format", "json")
	if stderr != "" {
		t.Fatalf("policy show stderr = %q, want empty", stderr)
	}
	rows := decodeCommandJSONStructure(t, stdout).([]any)
	if !commandPolicyCellAllowed(t, rows, "label.add", "developer") {
		t.Fatal("policy show developer/label.add = false, want true")
	}
	if !commandPolicyCellAllowed(t, rows, "label.add", "tl") {
		t.Fatal("policy show tl/label.add = false, want true")
	}
	if !commandPolicyCellAllowed(t, rows, "label.add", "workflow_runner") {
		t.Fatal("policy show workflow_runner/label.add = false, want true")
	}
	if !commandPolicyCellAllowed(t, rows, "label.add", "human") {
		t.Fatal("policy show human/label.add = false, want true")
	}
	if commandPolicyCellAllowed(t, rows, "label.add", "qa_lead") {
		t.Fatal("policy show qa_lead/label.add = true, want false")
	}
}

func decodeCommandLabelMarkdownRows(t *testing.T, value string) any {
	t.Helper()

	lines := strings.Split(strings.TrimSpace(value), "\n")
	if len(lines) == 0 || lines[0] != "# Labels" {
		t.Fatalf("markdown output = %q, want # Labels header", value)
	}
	if len(lines) >= 3 && strings.TrimSpace(lines[2]) == "(no labels)" {
		return []any{}
	}
	if len(lines) < 5 {
		t.Fatalf("markdown output = %q, want labels table", value)
	}

	rows := []any{}
	for _, line := range lines[4:] {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}
		cells := splitMarkdownTableRow(trimmed)
		if len(cells) != 3 {
			t.Fatalf("markdown row = %q, want 3 cells", line)
		}
		row := map[string]any{
			"name":        cells[0],
			"colour":      cells[1],
			"description": cells[2],
		}
		if row["colour"] == "" {
			row["colour"] = nil
		}
		if row["description"] == "" {
			row["description"] = nil
		}
		rows = append(rows, row)
	}
	return rows
}
