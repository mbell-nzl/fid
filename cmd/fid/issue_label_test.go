package main

import (
	"database/sql"
	"encoding/json"
	"path/filepath"
	"strings"
	"testing"

	fiddb "github.com/mbell-nzl/fid/internal/db"
	"github.com/mbell-nzl/fid/internal/journal"
)

func TestIssueLabelAddPreservesExistingLabelsAndProjectsActivity(t *testing.T) {
	workspace := t.TempDir()
	runInitCommand(t, workspace, "init")

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
		"--by", "developer",
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
	if stdout != "" {
		t.Fatalf("stdout = %q, want empty", stdout)
	}
	if !strings.Contains(stderr, `added label "api" to `+issueID) {
		t.Fatalf("stderr = %q, want add confirmation", stderr)
	}

	if got := len(splitJournalLines(mustReadCommandFile(t, journalPath))); got != linesBefore+1 {
		t.Fatalf("journal lines after label add = %d, want %d", got, linesBefore+1)
	}
	if got := countActivityLogRows(t, conn); got != activityBefore+1 {
		t.Fatalf("activity_log rows after label add = %d, want %d", got, activityBefore+1)
	}

	getStdout, getStderr := runFIDCommand(t, workspace, "issue", "get", issueID, "--format", "json")
	if getStderr != "" {
		t.Fatalf("issue get stderr = %q, want empty", getStderr)
	}

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
	if activityEvent != "issue_label_added" {
		t.Fatalf("activity_log event = %q, want %q", activityEvent, "issue_label_added")
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

func TestIssueLabelCommandsRejectBeforeAppendAndHonorPermissions(t *testing.T) {
	t.Run("add rejects missing and duplicate label", func(t *testing.T) {
		workspace := t.TempDir()
		runInitCommand(t, workspace, "init")

		runFIDCommand(t, workspace, "label", "add", "core", "--by", "developer")
		_, createStderr := runFIDCommand(
			t,
			workspace,
			"issue", "create",
			"--kind", "task",
			"--title", "CLI issue label rejection",
			"--label", "core",
			"--by", "developer",
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
		if !strings.Contains(stderr, `error: issue label add: label "missing" does not exist`) {
			t.Fatalf("missing label stderr = %q, want missing-label error", stderr)
		}

		stdout, stderr, err = runFIDCommandExpectError(t, workspace, "issue", "label", "add", issueID, "core", "--by", "developer")
		if err == nil {
			t.Fatal("Execute(duplicate label) error = nil, want rejection")
		}
		if stdout != "" {
			t.Fatalf("duplicate label stdout = %q, want empty", stdout)
		}
		if !strings.Contains(stderr, `error: issue label add: issue "`+issueID+`" already has label "core"`) {
			t.Fatalf("duplicate label stderr = %q, want duplicate-label error", stderr)
		}

		if got := len(splitJournalLines(mustReadCommandFile(t, journalPath))); got != linesBefore {
			t.Fatalf("journal lines after rejected add = %d, want %d", got, linesBefore)
		}
	})

	t.Run("permissions reject denied and unknown roles", func(t *testing.T) {
		workspace := t.TempDir()
		runInitCommand(t, workspace, "init")

		runFIDCommand(t, workspace, "label", "add", "core", "--by", "developer")
		_, createStderr := runFIDCommand(
			t,
			workspace,
			"issue", "create",
			"--kind", "task",
			"--title", "CLI issue label permissions",
			"--by", "developer",
		)
		issueID := parseCreatedIssueID(t, createStderr)

		journalPath := mustJournalPath(t, workspace)
		linesBefore := len(splitJournalLines(mustReadCommandFile(t, journalPath)))

		stdout, stderr, err := runFIDCommandExpectError(t, workspace, "issue", "label", "add", issueID, "core", "--by", "qa_reviewer")
		if err == nil {
			t.Fatal("Execute(denied add) error = nil, want rejection")
		}
		if stdout != "" {
			t.Fatalf("denied add stdout = %q, want empty", stdout)
		}
		if !strings.Contains(stderr, `error: role "qa_reviewer" not permitted for action "issue.label.add"`) {
			t.Fatalf("denied add stderr = %q, want permission error", stderr)
		}

		stdout, stderr, err = runFIDCommandExpectError(t, workspace, "issue", "label", "add", issueID, "core", "--by", "intern")
		if err == nil {
			t.Fatal("Execute(unknown add role) error = nil, want rejection")
		}
		if stdout != "" {
			t.Fatalf("unknown add role stdout = %q, want empty", stdout)
		}
		if !strings.Contains(stderr, `error: role "intern" not in project profile`) {
			t.Fatalf("unknown add role stderr = %q, want unknown-role error", stderr)
		}

		if got := len(splitJournalLines(mustReadCommandFile(t, journalPath))); got != linesBefore {
			t.Fatalf("journal lines after rejected permission checks = %d, want %d", got, linesBefore)
		}
	})
}

func TestIssueLabelRemovePreservesRemainingLabelsAndProjectsActivity(t *testing.T) {
	workspace := t.TempDir()
	runInitCommand(t, workspace, "init")

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
		"--by", "developer",
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
	if stdout != "" {
		t.Fatalf("stdout = %q, want empty", stdout)
	}
	if !strings.Contains(stderr, `removed label "legacy" from `+issueID) {
		t.Fatalf("stderr = %q, want remove confirmation", stderr)
	}

	if got := len(splitJournalLines(mustReadCommandFile(t, journalPath))); got != linesBefore+1 {
		t.Fatalf("journal lines after label remove = %d, want %d", got, linesBefore+1)
	}
	if got := countActivityLogRows(t, conn); got != activityBefore+1 {
		t.Fatalf("activity_log rows after label remove = %d, want %d", got, activityBefore+1)
	}

	getStdout, getStderr := runFIDCommand(t, workspace, "issue", "get", issueID, "--format", "json")
	if getStderr != "" {
		t.Fatalf("issue get stderr = %q, want empty", getStderr)
	}

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
	if activityEvent != "issue_label_removed" {
		t.Fatalf("activity_log event = %q, want %q", activityEvent, "issue_label_removed")
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
	workspace := t.TempDir()
	runInitCommand(t, workspace, "init")

	runFIDCommand(t, workspace, "label", "add", "core", "--by", "developer")
	_, createStderr := runFIDCommand(
		t,
		workspace,
		"issue", "create",
		"--kind", "task",
		"--title", "CLI issue label absent removal",
		"--label", "core",
		"--by", "developer",
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
	if !strings.Contains(stderr, `error: issue label remove: issue "`+issueID+`" does not have label "legacy"`) {
		t.Fatalf("stderr = %q, want absent-label error", stderr)
	}
	if got := len(splitJournalLines(mustReadCommandFile(t, journalPath))); got != linesBefore {
		t.Fatalf("journal lines after rejected remove = %d, want %d", got, linesBefore)
	}
}

func TestIssueLabelCommandsPolicyShowIncludesMatrixRows(t *testing.T) {
	workspace := t.TempDir()
	runInitCommand(t, workspace, "init")

	stdout, stderr := runFIDCommand(t, workspace, "policy", "show", "--format", "json")
	if stderr != "" {
		t.Fatalf("policy show stderr = %q, want empty", stderr)
	}
	rows := decodeCommandJSONStructure(t, stdout).([]any)

	for _, action := range []string{"issue.label.add", "issue.label.remove"} {
		for _, role := range []string{"developer", "tl", "workflow_runner", "human"} {
			if !commandPolicyCellAllowed(t, rows, action, role) {
				t.Fatalf("policy show %s/%s = false, want true", role, action)
			}
		}
		for _, role := range []string{"qa_reviewer", "qa_lead", "tl_lite"} {
			if commandPolicyCellAllowed(t, rows, action, role) {
				t.Fatalf("policy show %s/%s = true, want false", role, action)
			}
		}
	}
}
