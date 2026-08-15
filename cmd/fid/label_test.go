package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	fiddb "github.com/mbell-nzl/fid/internal/db"
)

func TestLabelAddCreatesOneRowInOneEventWithoutActivityLog(t *testing.T) {
	t.Parallel()
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
	if stderr != "" {
		t.Fatalf("stderr = %q, want empty", stderr)
	}
	if !strings.Contains(stdout, `created label "core"`) {
		t.Fatalf("stdout = %q, want create confirmation", stdout)
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
	if event.Actor.Role != "developer" {
		t.Fatalf("event actor role = %q, want developer", event.Actor.Role)
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

func TestLabelAddJSONMatchesContract(t *testing.T) {
	t.Parallel()
	workspace := t.TempDir()
	clonePristineWorkspaceInto(t, workspace)

	journalPath := mustJournalPath(t, workspace)
	result := executeCommandForTest(workspace, "--json", "label", "add", "core", "--colour", "#4A90D9", "--description", "Core infrastructure", "--by", "developer")
	if result.err != nil || result.stderr != "" {
		t.Fatalf("label add result = %#v, want JSON success on stdout only", result)
	}
	var envelope struct {
		ContractVersion string         `json:"contract_version"`
		OK              bool           `json:"ok"`
		Operation       string         `json:"operation"`
		Result          map[string]any `json:"result"`
		Event           map[string]any `json:"event"`
	}
	if err := json.Unmarshal([]byte(result.stdout), &envelope); err != nil {
		t.Fatalf("Unmarshal(stdout) error = %v; stdout = %q", err, result.stdout)
	}
	assertEnvelopeKeySet(t, []byte(result.stdout), []string{"contract_version", "event", "ok", "operation", "result"})
	assertExactKeys(t, envelope.Result, []string{"label"})
	label := requireObject(t, envelope.Result, "label")
	assertExactKeys(t, label, []string{"name", "colour", "description"})
	assertExactKeys(t, envelope.Event, []string{"id", "occurred_at"})
	if envelope.ContractVersion != "fid-cli.v1" || !envelope.OK || envelope.Operation != "label.add" {
		t.Fatalf("envelope = %#v, want fid-cli.v1 label.add success", envelope)
	}
	if label["name"] != "core" || label["colour"] != "#4A90D9" || label["description"] != "Core infrastructure" {
		t.Fatalf("label = %#v, want exact contract values", label)
	}
	last := mustReadLastJournalEvent(t, journalPath)
	if envelope.Event["id"] != last.EventID || envelope.Event["occurred_at"] != last.OccurredAt.UTC().Format("2006-01-02T15:04:05Z") {
		t.Fatalf("envelope event = %#v, journal tail = %#v", envelope.Event, last)
	}
}

func TestLabelAddJSONDuplicateUsesSpecificAtomicFailureCode(t *testing.T) {
	workspace := t.TempDir()
	clonePristineWorkspaceInto(t, workspace)
	runFIDCommand(t, workspace, "label", "add", "core", "--by", "developer")
	journalPath := mustJournalPath(t, workspace)
	before := len(splitJournalLines(mustReadCommandFile(t, journalPath)))

	result := executeCommandForTest(workspace, "--json", "label", "add", "core", "--by", "developer")
	if result.err == nil || result.stdout != "" {
		t.Fatalf("result = %#v, want stderr-only duplicate rejection", result)
	}
	assertJSONFailureCode(t, result.stderr, "label.add", "duplicate-entity")
	if got := len(splitJournalLines(mustReadCommandFile(t, journalPath))); got != before {
		t.Fatalf("journal lines after duplicate = %d, want %d", got, before)
	}
}

func TestLabelAddRejectsExplicitlyEmptyOptionalFlagsBeforeAppend(t *testing.T) {
	for _, flag := range []string{"colour", "description"} {
		t.Run(flag, func(t *testing.T) {
			workspace := t.TempDir()
			clonePristineWorkspaceInto(t, workspace)
			journalPath := mustJournalPath(t, workspace)
			before := len(splitJournalLines(mustReadCommandFile(t, journalPath)))
			result := executeCommandForTest(workspace, "--json", "label", "add", "core", "--"+flag, " ", "--by", "developer")
			if result.err == nil || result.stdout != "" {
				t.Fatalf("result = %#v, want stderr-only empty-%s rejection", result, flag)
			}
			assertJSONFailureCode(t, result.stderr, "label.add", "invalid-argument")
			if got := len(splitJournalLines(mustReadCommandFile(t, journalPath))); got != before {
				t.Fatalf("journal lines after rejection = %d, want %d", got, before)
			}
		})
	}
}

func TestLabelListMatchesClosedHumanAndJSONContracts(t *testing.T) {
	t.Parallel()
	workspace := t.TempDir()
	clonePristineWorkspaceInto(t, workspace)

	runFIDCommand(t, workspace, "label", "add", "core", "--colour", "#4A90D9", "--description", "Core infrastructure", "--by", "developer")
	runFIDCommand(t, workspace, "label", "add", "tech-debt", "--by", "developer")

	defaultStdout, defaultStderr := runFIDCommand(t, workspace, "label", "list")
	if defaultStderr != "" {
		t.Fatalf("default stderr = %q, want empty", defaultStderr)
	}
	if defaultStdout != "core\ntech-debt\n" {
		t.Fatalf("default stdout = %q, want sorted label names", defaultStdout)
	}

	jsonStdout, jsonStderr := runFIDCommand(t, workspace, "--json", "label", "list")
	if jsonStderr != "" {
		t.Fatalf("json stderr = %q, want empty", jsonStderr)
	}
	envelope := decodeReadEnvelope(t, jsonStdout, "label.list")
	assertExactKeys(t, envelope.Result, []string{"labels", "count"})
	if envelope.Result["count"] != float64(2) {
		t.Fatalf("count = %#v, want 2", envelope.Result["count"])
	}
	rows, ok := envelope.Result["labels"].([]any)
	if !ok || len(rows) != 2 {
		t.Fatalf("labels = %#v, want two rows", envelope.Result["labels"])
	}
	want := []map[string]any{{"name": "core", "colour": "#4A90D9", "description": "Core infrastructure"}, {"name": "tech-debt", "colour": nil, "description": nil}}
	for index, raw := range rows {
		row, ok := raw.(map[string]any)
		if !ok {
			t.Fatalf("labels[%d] = %#v, want object", index, raw)
		}
		assertExactKeys(t, row, []string{"name", "colour", "description"})
		if !commandStructuresEqual(want[index], row) {
			t.Fatalf("labels[%d] = %#v, want %#v", index, row, want[index])
		}
	}
}

func TestLabelListRejectsLegacyFormatFlag(t *testing.T) {
	t.Parallel()
	workspace := t.TempDir()
	clonePristineWorkspaceInto(t, workspace)

	runFIDCommand(t, workspace, "label", "add", "core", "--colour", "#4A90D9", "--description", "Core infrastructure", "--by", "developer")
	runFIDCommand(t, workspace, "label", "add", "tech-debt", "--by", "developer")

	for _, formatName := range []string{"json", "toon", "yaml", "markdown"} {
		t.Run(formatName, func(t *testing.T) {
			stdout, stderr, err := runFIDCommandExpectError(t, workspace, "label", "list", "--format", formatName)
			if err == nil {
				t.Fatal("Execute() error = nil, want removed --format rejection")
			}
			if stdout != "" {
				t.Fatalf("stdout = %q, want empty", stdout)
			}
			if stderr != "" {
				t.Fatalf("stderr = %q, want no rendered legacy-format output", stderr)
			}
			if !strings.Contains(err.Error(), "unknown flag: --format") {
				t.Fatalf("error = %q, want removed --format rejection", err)
			}
		})
	}
}

func TestLabelListJSONReportsOverLimitWithoutOutput(t *testing.T) {
	workspace := t.TempDir()
	clonePristineWorkspaceInto(t, workspace)
	conn, err := fiddb.Open(filepath.Join(workspace, ".fid", "fid.sqlite"))
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	tx, err := conn.Begin()
	if err != nil {
		t.Fatalf("Begin() error = %v", err)
	}
	for index := 0; index <= 10000; index++ {
		if _, err := tx.Exec(`INSERT INTO labels (name) VALUES (?)`, fmt.Sprintf("label-%05d", index)); err != nil {
			t.Fatalf("insert label %d: %v", index, err)
		}
	}
	if err := tx.Commit(); err != nil {
		t.Fatalf("Commit() error = %v", err)
	}
	closeCommandDB(t, conn)

	result := executeCommandForTest(workspace, "--json", "label", "list")
	if result.err == nil || result.stdout != "" {
		t.Fatalf("result = %#v, want stderr-only limit rejection", result)
	}
	assertJSONFailureCode(t, result.stderr, "label.list", "result-limit-exceeded")
}

func TestLabelAddRejectsDuplicateDeniedAndUnknownBeforeAppend(t *testing.T) {
	t.Parallel()
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
	t.Parallel()
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
