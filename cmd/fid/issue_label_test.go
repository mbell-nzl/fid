package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"
	"testing"
	"time"

	fiddb "github.com/mbell-nzl/fid/internal/db"
	"github.com/mbell-nzl/fid/internal/journal"
)

func TestIssueLabelAddPreservesExistingLabelsAndProjectsActivity(t *testing.T) {
	t.Parallel()
	workspace := t.TempDir()
	clonePristineWorkspaceInto(t, workspace)

	runFIDCommand(t, workspace, "label", "add", "api", "--by", "developer")
	runFIDCommand(t, workspace, "label", "add", "core", "--by", "developer")
	runFIDCommand(t, workspace, "label", "add", "legacy", "--by", "developer")

	_, createStderr := runFIDCommand(
		t,
		workspace,
		"issue", "create",
		"--kind", "task",
		"--title", "CLI issue label add",
		"--label", "core",
		"--label", "legacy",
		"--actor", "agent:developer",
	)
	issueID := parseCreatedIssueID(t, createStderr)

	journalPath := mustJournalPath(t, workspace)
	linesBefore := len(splitJournalLines(mustReadCommandFile(t, journalPath)))

	conn, err := fiddb.Open(filepath.Join(workspace, ".fid", "fid.sqlite"))
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	defer closeCommandDB(t, conn)

	activityBefore := countActivityLogRows(t, conn)

	stdout, stderr := runFIDCommand(t, workspace, "issue", "label", "add", issueID, "api", "--by", "developer")
	if stderr != "" {
		t.Fatalf("stderr = %q, want empty", stderr)
	}
	if !strings.Contains(stdout, `added label "api" to `+issueID) {
		t.Fatalf("stdout = %q, want add confirmation", stdout)
	}

	if got := len(splitJournalLines(mustReadCommandFile(t, journalPath))); got != linesBefore+1 {
		t.Fatalf("journal lines after label add = %d, want %d", got, linesBefore+1)
	}
	if got := countActivityLogRows(t, conn); got != activityBefore+1 {
		t.Fatalf("activity_log rows after label add = %d, want %d", got, activityBefore+1)
	}

	getStdout := mustGetIssueJSON(t, workspace, issueID)

	var issue struct {
		Labels []string `json:"labels"`
	}
	if err := json.Unmarshal([]byte(getStdout), &issue); err != nil {
		t.Fatalf("Unmarshal(issue get) error = %v", err)
	}
	if want := []string{"api", "core", "legacy"}; !equalIssueIDs(issue.Labels, want) {
		t.Fatalf("labels = %#v, want %#v", issue.Labels, want)
	}

	event := mustReadLastJournalEvent(t, journalPath)
	if event.EventType != "issue.label.added" {
		t.Fatalf("event type = %q, want %q", event.EventType, "issue.label.added")
	}

	var (
		insertEffect   *journal.Effect
		activityEffect *journal.Effect
	)
	for index := range event.Effects {
		effect := &event.Effects[index]
		switch {
		case effect.Table == "issue_labels" && effect.Op == "insert":
			insertEffect = effect
		case effect.Table == "activity_log" && effect.Op == "activity":
			activityEffect = effect
		}
	}
	if insertEffect == nil {
		t.Fatal("missing issue_labels insert effect")
	}
	if activityEffect == nil {
		t.Fatal("missing activity effect")
	}
	if insertEffect.PK["issue_id"] != issueID {
		t.Fatalf("issue_labels issue_id = %#v, want %q", insertEffect.PK["issue_id"], issueID)
	}
	if insertEffect.PK["label"] != "api" {
		t.Fatalf("issue_labels label = %#v, want %q", insertEffect.PK["label"], "api")
	}
	if activityEffect.Values["journal_event_id"] != event.EventID {
		t.Fatalf("activity journal_event_id = %#v, want %q", activityEffect.Values["journal_event_id"], event.EventID)
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
	if activityEvent != "issue-label-added" {
		t.Fatalf("activity_log event = %q, want %q", activityEvent, "issue-label-added")
	}
	if actor != "developer" {
		t.Fatalf("activity_log actor = %q, want %q", actor, "developer")
	}
	for _, want := range []string{`"label":"api"`, `"direction":"added"`} {
		if !strings.Contains(details.String, want) {
			t.Fatalf("activity_log details = %q, want substring %q", details.String, want)
		}
	}
}

func TestIssueLabelMutationsJSONMatchContract(t *testing.T) {
	workspace := t.TempDir()
	clonePristineWorkspaceInto(t, workspace)
	runFIDCommand(t, workspace, "label", "add", "core", "--by", "developer")
	_, createStderr := runFIDCommand(t, workspace, "issue", "create", "--kind", "task", "--title", "Issue label envelope", "--actor", "agent:developer")
	issueID := parseCreatedIssueID(t, createStderr)

	journalPath := mustJournalPath(t, workspace)
	add := executeCommandForTest(workspace, "--json", "issue", "label", "add", issueID, "core", "--by", "developer")
	assertIssueLabelMutationEnvelope(t, add, journalPath, "issue.label.add", issueID, "core", 2)

	remove := executeCommandForTest(workspace, "--json", "issue", "label", "remove", issueID, "core", "--by", "developer")
	assertIssueLabelMutationEnvelope(t, remove, journalPath, "issue.label.remove", issueID, "core", 3)
}

func TestIssueLabelJSONUsesSpecificAtomicFailureCodes(t *testing.T) {
	workspace := t.TempDir()
	clonePristineWorkspaceInto(t, workspace)
	runFIDCommand(t, workspace, "label", "add", "core", "--by", "developer")
	_, createStderr := runFIDCommand(t, workspace, "issue", "create", "--kind", "task", "--title", "Issue label JSON failures", "--label", "core", "--actor", "agent:developer")
	issueID := parseCreatedIssueID(t, createStderr)

	journalPath := mustJournalPath(t, workspace)
	tests := []struct {
		name      string
		operation string
		args      []string
		code      string
	}{
		{name: "add missing issue", operation: "issue.label.add", args: []string{"--json", "issue", "label", "add", "FID-missing", "core", "--by", "developer"}, code: "issue-not-found"},
		{name: "add missing label", operation: "issue.label.add", args: []string{"--json", "issue", "label", "add", issueID, "missing", "--by", "developer"}, code: "label-not-found"},
		{name: "add duplicate label", operation: "issue.label.add", args: []string{"--json", "issue", "label", "add", issueID, "core", "--by", "developer"}, code: "duplicate-entity"},
		{name: "remove missing issue", operation: "issue.label.remove", args: []string{"--json", "issue", "label", "remove", "FID-missing", "core", "--by", "developer"}, code: "issue-not-found"},
		{name: "remove absent label", operation: "issue.label.remove", args: []string{"--json", "issue", "label", "remove", issueID, "missing", "--by", "developer"}, code: "label-not-found"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			before := len(splitJournalLines(mustReadCommandFile(t, journalPath)))
			result := executeCommandForTest(workspace, test.args...)
			if result.err == nil || result.stdout != "" {
				t.Fatalf("result = %#v, want stderr-only rejection", result)
			}
			assertJSONFailureCode(t, result.stderr, test.operation, test.code)
			if got := len(splitJournalLines(mustReadCommandFile(t, journalPath))); got != before {
				t.Fatalf("journal lines after rejection = %d, want %d", got, before)
			}
		})
	}
}

func TestIssueLabelAddJSONReportsNestedLimitWithoutAppend(t *testing.T) {
	workspace := t.TempDir()
	clonePristineWorkspaceInto(t, workspace)
	_, createStderr := runFIDCommand(t, workspace, "issue", "create", "--kind", "task", "--title", "Issue label JSON limit", "--actor", "agent:developer")
	issueID := parseCreatedIssueID(t, createStderr)
	conn, err := fiddb.Open(filepath.Join(workspace, ".fid", "fid.sqlite"))
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	tx, err := conn.Begin()
	if err != nil {
		t.Fatalf("Begin() error = %v", err)
	}
	for index := 0; index <= 100; index++ {
		label := fmt.Sprintf("limit-%03d", index)
		if _, err := tx.Exec(`INSERT INTO labels (name) VALUES (?)`, label); err != nil {
			t.Fatalf("insert label %d: %v", index, err)
		}
		if index < 100 {
			if _, err := tx.Exec(`INSERT INTO issue_labels (issue_id, label) VALUES (?, ?)`, issueID, label); err != nil {
				t.Fatalf("insert issue label %d: %v", index, err)
			}
		}
	}
	if err := tx.Commit(); err != nil {
		t.Fatalf("Commit() error = %v", err)
	}
	closeCommandDB(t, conn)

	journalPath := mustJournalPath(t, workspace)
	before := len(splitJournalLines(mustReadCommandFile(t, journalPath)))
	result := executeCommandForTest(workspace, "--json", "issue", "label", "add", issueID, "limit-100", "--by", "developer")
	if result.err == nil || result.stdout != "" {
		t.Fatalf("result = %#v, want stderr-only limit rejection", result)
	}
	assertJSONFailureCode(t, result.stderr, "issue.label.add", "result-limit-exceeded")
	if got := len(splitJournalLines(mustReadCommandFile(t, journalPath))); got != before {
		t.Fatalf("journal lines after rejection = %d, want %d", got, before)
	}
}

func assertIssueLabelMutationEnvelope(t *testing.T, result commandExecutionResult, journalPath, operation, issueID, label string, revision int) {
	t.Helper()
	if result.err != nil || result.stderr != "" {
		t.Fatalf("%s result = %#v, want JSON success on stdout only", operation, result)
	}
	var envelope struct {
		ContractVersion string         `json:"contract_version"`
		OK              bool           `json:"ok"`
		Operation       string         `json:"operation"`
		Result          map[string]any `json:"result"`
		Event           map[string]any `json:"event"`
	}
	if err := json.Unmarshal([]byte(result.stdout), &envelope); err != nil {
		t.Fatalf("Unmarshal(%s stdout) error = %v; stdout = %q", operation, err, result.stdout)
	}
	assertEnvelopeKeySet(t, []byte(result.stdout), []string{"contract_version", "event", "ok", "operation", "result"})
	assertExactKeys(t, envelope.Result, []string{"idempotency_key", "issue_id", "issue_revision", "label", "replayed"})
	assertExactKeys(t, envelope.Event, []string{"id", "occurred_at"})
	if envelope.ContractVersion != "fid-cli.v1" || !envelope.OK || envelope.Operation != operation {
		t.Fatalf("envelope = %#v, want fid-cli.v1 %s success", envelope, operation)
	}
	if envelope.Result["issue_id"] != issueID || envelope.Result["label"] != label || envelope.Result["issue_revision"] != float64(revision) {
		t.Fatalf("result = %#v, want issue %q label %q revision %d", envelope.Result, issueID, label, revision)
	}
	last := mustReadLastJournalEvent(t, journalPath)
	if envelope.Event["id"] != last.EventID || envelope.Event["occurred_at"] != last.OccurredAt.UTC().Format(time.RFC3339Nano) {
		t.Fatalf("envelope event = %#v, journal tail = %#v", envelope.Event, last)
	}
}

func TestIssueLabelCommandsRejectBeforeAppendAndHonorPermissions(t *testing.T) {
	t.Parallel()
	t.Run("add rejects missing and duplicate label", func(t *testing.T) {
		workspace := t.TempDir()
		clonePristineWorkspaceInto(t, workspace)

		runFIDCommand(t, workspace, "label", "add", "core", "--by", "developer")
		_, createStderr := runFIDCommand(
			t,
			workspace,
			"issue", "create",
			"--kind", "task",
			"--title", "CLI issue label rejection",
			"--label", "core",
			"--actor", "agent:developer",
		)
		issueID := parseCreatedIssueID(t, createStderr)

		journalPath := mustJournalPath(t, workspace)
		linesBefore := len(splitJournalLines(mustReadCommandFile(t, journalPath)))

		stdout, stderr, err := runFIDCommandExpectError(t, workspace, "issue", "label", "add", issueID, "missing", "--by", "developer")
		if err == nil {
			t.Fatal("Execute(missing label) error = nil, want rejection")
		}
		if stdout != "" {
			t.Fatalf("missing label stdout = %q, want empty", stdout)
		}
		if !strings.Contains(stderr, `error: label "missing" not found`) {
			t.Fatalf("missing label stderr = %q, want missing-label error", stderr)
		}

		stdout, stderr, err = runFIDCommandExpectError(t, workspace, "issue", "label", "add", issueID, "core", "--by", "developer")
		if err == nil {
			t.Fatal("Execute(duplicate label) error = nil, want rejection")
		}
		if stdout != "" {
			t.Fatalf("duplicate label stdout = %q, want empty", stdout)
		}
		if !strings.Contains(stderr, `error: issue "`+issueID+`" already has label "core"`) {
			t.Fatalf("duplicate label stderr = %q, want duplicate-label error", stderr)
		}

		if got := len(splitJournalLines(mustReadCommandFile(t, journalPath))); got != linesBefore {
			t.Fatalf("journal lines after rejected add = %d, want %d", got, linesBefore)
		}
	})

	t.Run("project-local actors are accepted without role policy", func(t *testing.T) {
		workspace := t.TempDir()
		clonePristineWorkspaceInto(t, workspace)

		runFIDCommand(t, workspace, "label", "add", "core", "--by", "developer")
		_, createStderr := runFIDCommand(
			t,
			workspace,
			"issue", "create",
			"--kind", "task",
			"--title", "CLI issue label permissions",
			"--actor", "agent:developer",
		)
		issueID := parseCreatedIssueID(t, createStderr)

		journalPath := mustJournalPath(t, workspace)
		linesBefore := len(splitJournalLines(mustReadCommandFile(t, journalPath)))

		for _, fixture := range []struct{ label, actor string }{{"qa", "agent:qa_reviewer"}, {"intern", "agent:intern"}} {
			runFIDCommand(t, workspace, "label", "add", fixture.label, "--actor", fixture.actor)
			stdout, stderr := runFIDCommand(t, workspace, "issue", "label", "add", issueID, fixture.label, "--actor", fixture.actor)
			if stderr != "" || !strings.Contains(stdout, "added label") {
				t.Fatalf("%s actor issue label = (%q, %q), want stdout confirmation", fixture.actor, stdout, stderr)
			}
		}

		if got := len(splitJournalLines(mustReadCommandFile(t, journalPath))); got != linesBefore+4 {
			t.Fatalf("journal lines after project-local actor checks = %d, want %d", got, linesBefore+4)
		}
	})
}

func TestIssueLabelRemovePreservesRemainingLabelsAndProjectsActivity(t *testing.T) {
	t.Parallel()
	workspace := t.TempDir()
	clonePristineWorkspaceInto(t, workspace)

	runFIDCommand(t, workspace, "label", "add", "api", "--by", "developer")
	runFIDCommand(t, workspace, "label", "add", "core", "--by", "developer")
	runFIDCommand(t, workspace, "label", "add", "legacy", "--by", "developer")

	_, createStderr := runFIDCommand(
		t,
		workspace,
		"issue", "create",
		"--kind", "task",
		"--title", "CLI issue label remove",
		"--label", "api",
		"--label", "core",
		"--label", "legacy",
		"--actor", "agent:developer",
	)
	issueID := parseCreatedIssueID(t, createStderr)

	journalPath := mustJournalPath(t, workspace)
	linesBefore := len(splitJournalLines(mustReadCommandFile(t, journalPath)))

	conn, err := fiddb.Open(filepath.Join(workspace, ".fid", "fid.sqlite"))
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	defer closeCommandDB(t, conn)

	activityBefore := countActivityLogRows(t, conn)

	stdout, stderr := runFIDCommand(t, workspace, "issue", "label", "remove", issueID, "legacy", "--by", "developer")
	if stderr != "" {
		t.Fatalf("stderr = %q, want empty", stderr)
	}
	if !strings.Contains(stdout, `removed label "legacy" from `+issueID) {
		t.Fatalf("stdout = %q, want remove confirmation", stdout)
	}

	if got := len(splitJournalLines(mustReadCommandFile(t, journalPath))); got != linesBefore+1 {
		t.Fatalf("journal lines after label remove = %d, want %d", got, linesBefore+1)
	}
	if got := countActivityLogRows(t, conn); got != activityBefore+1 {
		t.Fatalf("activity_log rows after label remove = %d, want %d", got, activityBefore+1)
	}

	getStdout := mustGetIssueJSON(t, workspace, issueID)

	var issue struct {
		Labels []string `json:"labels"`
	}
	if err := json.Unmarshal([]byte(getStdout), &issue); err != nil {
		t.Fatalf("Unmarshal(issue get) error = %v", err)
	}
	if want := []string{"api", "core"}; !equalIssueIDs(issue.Labels, want) {
		t.Fatalf("labels = %#v, want %#v", issue.Labels, want)
	}

	var labelCount int
	if err := conn.QueryRow(`SELECT COUNT(*) FROM labels WHERE name = ?`, "legacy").Scan(&labelCount); err != nil {
		t.Fatalf("count label row: %v", err)
	}
	if labelCount != 1 {
		t.Fatalf("labels row count for legacy = %d, want 1", labelCount)
	}

	event := mustReadLastJournalEvent(t, journalPath)
	if event.EventType != "issue.label.removed" {
		t.Fatalf("event type = %q, want %q", event.EventType, "issue.label.removed")
	}

	var (
		deleteEffect   *journal.Effect
		activityEffect *journal.Effect
	)
	for index := range event.Effects {
		effect := &event.Effects[index]
		switch {
		case effect.Table == "issue_labels" && effect.Op == "delete":
			deleteEffect = effect
		case effect.Table == "activity_log" && effect.Op == "activity":
			activityEffect = effect
		}
	}
	if deleteEffect == nil {
		t.Fatal("missing issue_labels delete effect")
	}
	if activityEffect == nil {
		t.Fatal("missing activity effect")
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
	if activityEvent != "issue-label-removed" {
		t.Fatalf("activity_log event = %q, want %q", activityEvent, "issue-label-removed")
	}
	if actor != "developer" {
		t.Fatalf("activity_log actor = %q, want %q", actor, "developer")
	}
	for _, want := range []string{`"label":"legacy"`, `"direction":"removed"`} {
		if !strings.Contains(details.String, want) {
			t.Fatalf("activity_log details = %q, want substring %q", details.String, want)
		}
	}
}

func TestIssueLabelRemoveRejectsAbsentLabelBeforeAppend(t *testing.T) {
	t.Parallel()
	workspace := t.TempDir()
	clonePristineWorkspaceInto(t, workspace)

	runFIDCommand(t, workspace, "label", "add", "core", "--by", "developer")
	_, createStderr := runFIDCommand(
		t,
		workspace,
		"issue", "create",
		"--kind", "task",
		"--title", "CLI issue label absent removal",
		"--label", "core",
		"--actor", "agent:developer",
	)
	issueID := parseCreatedIssueID(t, createStderr)

	journalPath := mustJournalPath(t, workspace)
	linesBefore := len(splitJournalLines(mustReadCommandFile(t, journalPath)))

	stdout, stderr, err := runFIDCommandExpectError(t, workspace, "issue", "label", "remove", issueID, "legacy", "--by", "developer")
	if err == nil {
		t.Fatal("Execute(absent removal) error = nil, want rejection")
	}
	if stdout != "" {
		t.Fatalf("stdout = %q, want empty", stdout)
	}
	if !strings.Contains(stderr, `error: label "legacy" not found`) {
		t.Fatalf("stderr = %q, want absent-label error", stderr)
	}
	if got := len(splitJournalLines(mustReadCommandFile(t, journalPath))); got != linesBefore {
		t.Fatalf("journal lines after rejected remove = %d, want %d", got, linesBefore)
	}
}
