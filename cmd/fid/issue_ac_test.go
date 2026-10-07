package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"
	"testing"
	"time"

	fiddb "github.com/mbell-nzl/fid/internal/db"
)

func TestAcceptanceCriterionAddJSONMatchesContract(t *testing.T) {
	t.Parallel()
	workspace := t.TempDir()
	clonePristineWorkspaceInto(t, workspace)

	created := executeCommandForTest(workspace,
		"--json", "issue", "create", "--kind", "task", "--title", "AC envelope target", "--actor", "agent:developer",
	)
	if created.err != nil || created.stderr != "" {
		t.Fatalf("issue create = %#v, want stdout-only success", created)
	}
	var createEnvelope struct {
		Result struct {
			Issue struct {
				ID string `json:"id"`
			} `json:"issue"`
		} `json:"result"`
	}
	if err := json.Unmarshal([]byte(created.stdout), &createEnvelope); err != nil {
		t.Fatalf("Unmarshal(create stdout) error = %v", err)
	}

	journalPath := mustJournalPath(t, workspace)
	before := len(splitJournalLines(mustReadCommandFile(t, journalPath)))
	result := executeCommandForTest(workspace,
		"--json", "issue", "ac", "add", createEnvelope.Result.Issue.ID,
		"--text", "Response envelope has the exact AC result", "--ordinal", "1", "--by", "developer",
	)
	if result.err != nil {
		t.Fatalf("issue ac add error = %v, stderr = %q", result.err, result.stderr)
	}
	if result.stderr != "" {
		t.Fatalf("stderr = %q, want empty", result.stderr)
	}

	var envelope struct {
		ContractVersion string `json:"contract_version"`
		OK              bool   `json:"ok"`
		Operation       string `json:"operation"`
		Result          struct {
			AcceptanceCriterion struct {
				ID      string `json:"id"`
				IssueID string `json:"issue_id"`
				Ordinal int    `json:"ordinal"`
				Text    string `json:"text"`
			} `json:"acceptance_criterion"`
			IssueRevision int `json:"issue_revision"`
		} `json:"result"`
		Event struct {
			ID         string `json:"id"`
			OccurredAt string `json:"occurred_at"`
		} `json:"event"`
	}
	if err := json.Unmarshal([]byte(result.stdout), &envelope); err != nil {
		t.Fatalf("Unmarshal(stdout) error = %v; stdout = %q", err, result.stdout)
	}
	assertEnvelopeKeySet(t, []byte(result.stdout), []string{"contract_version", "event", "ok", "operation", "result"})
	assertNestedEnvelopeKeySet(t, []byte(result.stdout), "result", []string{"acceptance_criterion", "idempotency_key", "issue_revision", "replayed"})
	assertNestedEnvelopeKeySet(t, []byte(result.stdout), "result.acceptance_criterion", []string{"id", "issue_id", "ordinal", "text"})
	assertNestedEnvelopeKeySet(t, []byte(result.stdout), "event", []string{"id", "occurred_at"})
	if envelope.ContractVersion != "fid-cli.v1" || !envelope.OK || envelope.Operation != "issue.ac.add" {
		t.Fatalf("envelope = %+v, want fid-cli.v1 issue.ac.add success", envelope)
	}
	if envelope.Result.AcceptanceCriterion.IssueID != createEnvelope.Result.Issue.ID || envelope.Result.AcceptanceCriterion.Ordinal != 1 || envelope.Result.AcceptanceCriterion.Text != "Response envelope has the exact AC result" || !strings.HasPrefix(envelope.Result.AcceptanceCriterion.ID, "AC-") {
		t.Fatalf("acceptance criterion = %+v, want exact contract values", envelope.Result.AcceptanceCriterion)
	}
	if envelope.Result.IssueRevision != 2 {
		t.Fatalf("issue_revision = %d, want 2", envelope.Result.IssueRevision)
	}
	if envelope.Event.ID == "" || envelope.Event.OccurredAt == "" {
		t.Fatalf("event = %+v, want identity", envelope.Event)
	}
	if got := len(splitJournalLines(mustReadCommandFile(t, journalPath))); got != before+1 {
		t.Fatalf("journal lines = %d, want %d", got, before+1)
	}
	last := mustReadLastJournalEvent(t, journalPath)
	if last.EventID != envelope.Event.ID || last.OccurredAt.UTC().Format(time.RFC3339Nano) != envelope.Event.OccurredAt {
		t.Fatalf("envelope event = %+v, journal tail = %+v", envelope.Event, last)
	}
}

func TestAcceptanceCriterionAddJSONUsesSpecificAtomicFailureCodes(t *testing.T) {
	t.Parallel()
	workspace := t.TempDir()
	clonePristineWorkspaceInto(t, workspace)
	_, createStderr := runFIDCommand(t, workspace, "issue", "create", "--kind", "task", "--title", "AC JSON failures", "--actor", "agent:developer")
	issueID := parseCreatedIssueID(t, createStderr)
	runFIDCommand(t, workspace, "issue", "ac", "add", issueID, "--text", "existing", "--ordinal", "1", "--by", "developer")

	journalPath := mustJournalPath(t, workspace)
	tests := []struct {
		name string
		args []string
		code string
	}{
		{name: "duplicate ordinal", args: []string{"--json", "issue", "ac", "add", issueID, "--text", "duplicate", "--ordinal", "1", "--by", "developer"}, code: "duplicate-entity"},
		{name: "missing issue", args: []string{"--json", "issue", "ac", "add", "FID-missing", "--text", "missing", "--ordinal", "1", "--by", "developer"}, code: "issue-not-found"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			before := len(splitJournalLines(mustReadCommandFile(t, journalPath)))
			result := executeCommandForTest(workspace, test.args...)
			if result.err == nil || result.stdout != "" {
				t.Fatalf("result = %#v, want stderr-only rejection", result)
			}
			assertJSONFailureCode(t, result.stderr, "issue.ac.add", test.code)
			if got := len(splitJournalLines(mustReadCommandFile(t, journalPath))); got != before {
				t.Fatalf("journal lines after rejection = %d, want %d", got, before)
			}
		})
	}
}

func TestIssueACAddCommandAddsAcceptanceCriterionAndPrintsConfirmation(t *testing.T) {
	t.Parallel()
	workspace := clonePristineWorkspace(t)
	journalPath := mustJournalPath(t, workspace)

	_, createStderr := runFIDCommand(t, workspace, "issue", "create", "--kind", "task", "--title", "CLI AC add", "--actor", "agent:developer")
	issueID := parseCreatedIssueID(t, createStderr)

	stdout, stderr := runFIDCommand(
		t,
		workspace,
		"issue", "ac", "add", issueID,
		"--text", "Journal file is flushed to disk before returning success",
		"--ordinal", "1",
		"--by", "developer",
	)
	if stderr != "" {
		t.Fatalf("stderr = %q, want empty", stderr)
	}
	if !strings.Contains(stdout, "added AC AC-") {
		t.Fatalf("stdout = %q, want AC confirmation prefix", stdout)
	}
	if !strings.Contains(stdout, "on "+issueID+" (ordinal 1)") {
		t.Fatalf("stdout = %q, want issue id and ordinal", stdout)
	}
	if !strings.Contains(stdout, `"Journal file is flushed to disk before returning success"`) {
		t.Fatalf("stdout = %q, want quoted criterion text", stdout)
	}

	getStdout, getStderr := runFIDCommand(t, workspace, "--json", "issue", "ac", "list", issueID)
	if getStderr != "" {
		t.Fatalf("issue ac list stderr = %q, want empty", getStderr)
	}
	getStdout = mustCanonicalAcceptanceCriteriaJSON(t, getStdout)

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

	conn, err := fiddb.Open(filepath.Join(workspace, ".fid", "fid.sqlite"))
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	defer closeCommandDB(t, conn)

	event := mustReadLastJournalEvent(t, journalPath)
	if event.EventType != "issue.ac-added" {
		t.Fatalf("event type = %q, want issue.ac-added", event.EventType)
	}
	if got := countActivityLogRows(t, conn); got != 2 {
		t.Fatalf("activity_log rows after AC add = %d, want 2", got)
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
	if activityEvent != "ac-added" {
		t.Fatalf("activity_log event = %q, want ac-added", activityEvent)
	}
	if actor != "developer" {
		t.Fatalf("activity_log actor = %q, want developer", actor)
	}
	for _, want := range []string{got[0].ID, `"ordinal":1`} {
		if !strings.Contains(details.String, want) {
			t.Fatalf("activity_log details = %q, want substring %q", details.String, want)
		}
	}
}

func TestIssueACAddCommandRejectsBeforeAppend(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		args    func(t *testing.T, workspace string) []string
		wantErr string
		verify  func(t *testing.T, workspace string)
	}{
		{
			name: "missing by flag",
			args: func(t *testing.T, workspace string) []string {
				_, createStderr := runFIDCommand(t, workspace, "issue", "create", "--kind", "task", "--title", "Missing by", "--actor", "agent:developer")
				issueID := parseCreatedIssueID(t, createStderr)
				return []string{"issue", "ac", "add", issueID, "--text", "Missing role", "--ordinal", "1"}
			},
			wantErr: "error: --actor is required",
		},
		{
			name: "missing text flag",
			args: func(t *testing.T, workspace string) []string {
				_, createStderr := runFIDCommand(t, workspace, "issue", "create", "--kind", "task", "--title", "Missing text", "--actor", "agent:developer")
				issueID := parseCreatedIssueID(t, createStderr)
				return []string{"issue", "ac", "add", issueID, "--ordinal", "1", "--by", "developer"}
			},
			wantErr: "error: --text is required",
		},
		{
			name: "missing ordinal flag",
			args: func(t *testing.T, workspace string) []string {
				_, createStderr := runFIDCommand(t, workspace, "issue", "create", "--kind", "task", "--title", "Missing ordinal", "--actor", "agent:developer")
				issueID := parseCreatedIssueID(t, createStderr)
				return []string{"issue", "ac", "add", issueID, "--text", "Missing ordinal", "--by", "developer"}
			},
			wantErr: "error: --ordinal is required",
		},
		{
			name: "invalid ordinal",
			args: func(t *testing.T, workspace string) []string {
				_, createStderr := runFIDCommand(t, workspace, "issue", "create", "--kind", "task", "--title", "Bad ordinal", "--actor", "agent:developer")
				issueID := parseCreatedIssueID(t, createStderr)
				return []string{"issue", "ac", "add", issueID, "--text", "Bad ordinal", "--ordinal", "0", "--by", "developer"}
			},
			wantErr: `error: invalid argument "0" for "--ordinal" flag`,
		},
		{
			name: "invalid actor",
			args: func(t *testing.T, workspace string) []string {
				_, createStderr := runFIDCommand(t, workspace, "issue", "create", "--kind", "task", "--title", "Unknown role", "--actor", "agent:developer")
				issueID := parseCreatedIssueID(t, createStderr)
				return []string{"issue", "ac", "add", issueID, "--text", "Invalid actor", "--ordinal", "1", "--actor", "intern"}
			},
			wantErr: `error: invalid argument "intern" for "--actor" flag`,
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
				_, createStderr := runFIDCommand(t, workspace, "issue", "create", "--kind", "epic", "--title", "Epic AC reject", "--scheduling", "prioritised", "--actor", "agent:tl")
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

func TestIssueACVerifyCommandPersistsVerificationAndPrintsConfirmation(t *testing.T) {
	t.Parallel()
	workspace := clonePristineWorkspace(t)

	issueID, criterionID := seedIssueAcceptanceCriterion(t, workspace, "CLI AC verify")
	assertIssueACVerifyUnavailable(t, workspace, issueID, criterionID, "agent:developer")
}

func TestIssueACVerifyCommandKeepsLatestPerRoleIndependent(t *testing.T) {
	t.Parallel()
	workspace := clonePristineWorkspace(t)

	issueID, criterionID := seedIssueAcceptanceCriterion(t, workspace, "CLI multi-role verify")
	assertIssueACVerifyUnavailable(t, workspace, issueID, criterionID, "agent:developer")
	assertIssueACVerifyUnavailable(t, workspace, issueID, criterionID, "agent:qa_lead")
}

func TestIssueACVerifyCommandPreservesFailToPassHistory(t *testing.T) {
	t.Parallel()
	workspace := clonePristineWorkspace(t)

	issueID, criterionID := seedIssueAcceptanceCriterion(t, workspace, "CLI fail then pass")
	assertIssueACVerifyUnavailable(t, workspace, issueID, criterionID, "agent:developer")
	assertIssueACVerifyUnavailable(t, workspace, issueID, criterionID, "agent:developer")
}

func TestIssueACVerifyCommandLatestPerRoleSurvivesRebuild(t *testing.T) {
	t.Parallel()
	workspace := clonePristineWorkspace(t)

	issueID, criterionID := seedIssueAcceptanceCriterion(t, workspace, "CLI rebuild AC parity")
	assertIssueACVerifyUnavailable(t, workspace, issueID, criterionID, "agent:developer")
	assertIssueACVerifyUnavailable(t, workspace, issueID, criterionID, "agent:qa_lead")
}

func TestIssueACVerifyCommandRejectsBeforeAppend(t *testing.T) {
	t.Parallel()

	for _, actor := range []string{"agent:developer", "agent:qa_lead", "human:reviewer"} {
		t.Run(actor, func(t *testing.T) {
			workspace := t.TempDir()
			clonePristineWorkspaceInto(t, workspace)
			issueID, criterionID := seedIssueAcceptanceCriterion(t, workspace, "Fail-closed AC verification")
			assertIssueACVerifyUnavailable(t, workspace, issueID, criterionID, actor)
		})
	}
}

func TestIssueACVerifyCommandUsesRoleSpecificPermissionsAndHumanAttribution(t *testing.T) {
	t.Parallel()
	workspace := clonePristineWorkspace(t)

	issueID, criterionID := seedIssueAcceptanceCriterion(t, workspace, "CLI verifier permissions")
	assertIssueACVerifyUnavailable(t, workspace, issueID, criterionID, "agent:workflow_runner")
	assertIssueACVerifyUnavailable(t, workspace, issueID, criterionID, "agent:qa_reviewer")
	assertIssueACVerifyUnavailable(t, workspace, issueID, criterionID, "human:reviewer")
}

func assertIssueACVerifyUnavailable(t *testing.T, workspace, issueID, criterionID, actor string) {
	t.Helper()
	journalPath := mustJournalPath(t, workspace)
	before := mustReadCommandFile(t, journalPath)
	result := executeCommandForTest(workspace, "issue", "ac", "verify")
	if result.err == nil || result.stdout != "" || !strings.Contains(result.err.Error(), `unknown command "verify"`) {
		t.Fatalf("removed issue ac verify for %s/%s/%s = %#v, want unknown-command rejection", issueID, criterionID, actor, result)
	}
	if after := mustReadCommandFile(t, journalPath); !bytes.Equal(after, before) {
		t.Fatalf("issue ac verify actor %s changed journal bytes", actor)
	}
	conn, err := fiddb.Open(filepath.Join(workspace, ".fid", "fid.sqlite"))
	if err != nil {
		t.Fatalf("open AC verification projection: %v", err)
	}
	defer closeCommandDB(t, conn)
	var count int
	if err := conn.QueryRow(`SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name='ac_verifications'`).Scan(&count); err != nil {
		t.Fatalf("count removed AC verification tables: %v", err)
	}
	if count != 0 {
		t.Fatalf("AC verification rows = %d, want no writer rows", count)
	}
}

func TestIssueACListCommandMatchesClosedHumanAndJSONContracts(t *testing.T) {
	t.Parallel()
	workspace := clonePristineWorkspace(t)

	issueID, criterionID := seedIssueAcceptanceCriterion(t, workspace, "CLI AC list")
	secondCriterionStdout, secondCriterionStderr := runFIDCommand(
		t,
		workspace,
		"issue", "ac", "add", issueID,
		"--text", "Second criterion for CLI AC list",
		"--ordinal", "2",
		"--by", "developer",
	)
	if secondCriterionStderr != "" {
		t.Fatalf("second add stderr = %q, want empty", secondCriterionStderr)
	}
	secondCriterionID := parseAddedAcceptanceCriterionID(t, secondCriterionStdout)

	assertIssueACVerifyUnavailable(t, workspace, issueID, criterionID, "agent:developer")
	assertIssueACVerifyUnavailable(t, workspace, issueID, criterionID, "agent:developer")
	assertIssueACVerifyUnavailable(t, workspace, issueID, criterionID, "agent:qa_lead")

	defaultStdout, defaultStderr := runFIDCommand(t, workspace, "issue", "ac", "list", issueID)
	if defaultStderr != "" {
		t.Fatalf("default stderr = %q, want empty", defaultStderr)
	}
	for _, want := range []string{"1. Criterion for CLI AC list", "2. Second criterion for CLI AC list"} {
		if !strings.Contains(defaultStdout, want) {
			t.Fatalf("default stdout = %q, want line %q", defaultStdout, want)
		}
	}
	if strings.Contains(defaultStdout, criterionID) || strings.Contains(defaultStdout, secondCriterionID) {
		t.Fatalf("default stdout = %q, want concise human output without IDs", defaultStdout)
	}

	jsonStdout, jsonStderr := runFIDCommand(t, workspace, "--json", "issue", "ac", "list", issueID)
	if jsonStderr != "" {
		t.Fatalf("json stderr = %q, want empty", jsonStderr)
	}
	envelope := decodeReadEnvelope(t, jsonStdout, "issue.ac.list")
	assertExactKeys(t, envelope.Result, []string{"issue_id", "acceptance_criteria"})
	if envelope.Result["issue_id"] != issueID {
		t.Fatalf("issue_id = %#v, want %q", envelope.Result["issue_id"], issueID)
	}
	rows, ok := envelope.Result["acceptance_criteria"].([]any)
	if !ok {
		t.Fatalf("acceptance_criteria = %T, want []any", envelope.Result["acceptance_criteria"])
	}
	if len(rows) != 2 {
		t.Fatalf("ac list rows len = %d, want 2", len(rows))
	}
	firstRow := rows[0].(map[string]any)
	assertExactKeys(t, firstRow, []string{"id", "ordinal", "text"})
	if firstRow["id"] != criterionID {
		t.Fatalf("first row id = %#v, want %q", firstRow["id"], criterionID)
	}
	secondRow := rows[1].(map[string]any)
	assertExactKeys(t, secondRow, []string{"id", "ordinal", "text"})
	if secondRow["id"] != secondCriterionID {
		t.Fatalf("second row id = %#v, want %q", secondRow["id"], secondCriterionID)
	}
}

func TestIssueACListCommandRejectsLegacyFormatsAndHandlesEmptyAndMissingIssues(t *testing.T) {
	t.Parallel()
	workspace := clonePristineWorkspace(t)

	issueID, _ := seedIssueAcceptanceCriterion(t, workspace, "CLI AC transport replacement")

	for _, formatName := range []string{"json", "toon", "yaml", "markdown"} {
		t.Run(formatName, func(t *testing.T) {
			stdout, stderr, err := runFIDCommandExpectError(t, workspace, "issue", "ac", "list", issueID, "--format", formatName)
			if err == nil {
				t.Fatal("Execute() error = nil, want removed --format rejection")
			}
			if stdout != "" {
				t.Fatalf("stdout = %q, want empty", stdout)
			}
			if stderr != "error: unknown flag: --format\n" {
				t.Fatalf("stderr = %q, want concise removed-flag error", stderr)
			}
			if !strings.Contains(err.Error(), "unknown flag: --format") {
				t.Fatalf("error = %q, want removed --format rejection", err)
			}
		})
	}

	t.Run("empty list", func(t *testing.T) {
		_, emptyIssueStderr := runFIDCommand(t, workspace, "issue", "create", "--kind", "task", "--title", "Empty AC list", "--actor", "agent:developer")
		emptyIssueID := parseCreatedIssueID(t, emptyIssueStderr)

		stdout, stderr := runFIDCommand(t, workspace, "--json", "issue", "ac", "list", emptyIssueID)
		if stderr != "" {
			t.Fatalf("stderr = %q, want empty", stderr)
		}
		envelope := decodeReadEnvelope(t, stdout, "issue.ac.list")
		assertExactKeys(t, envelope.Result, []string{"issue_id", "acceptance_criteria"})
		if got := envelope.Result["acceptance_criteria"]; !commandStructuresEqual([]any{}, got) {
			t.Fatalf("empty acceptance_criteria = %#v, want named empty array", got)
		}
	})

	t.Run("missing issue", func(t *testing.T) {
		stdout, stderr, err := runFIDCommandExpectError(t, workspace, "--json", "issue", "ac", "list", "FID-missing")
		if err == nil {
			t.Fatal("Execute() error = nil, want missing issue rejection")
		}
		if stdout != "" {
			t.Fatalf("stdout = %q, want empty", stdout)
		}
		assertJSONFailureCode(t, stderr, "issue.ac.list", "issue-not-found")
	})
}

func TestIssueACListJSONReportsOverLimitWithoutOutput(t *testing.T) {
	t.Parallel()
	workspace := t.TempDir()
	clonePristineWorkspaceInto(t, workspace)
	_, createStderr := runFIDCommand(t, workspace, "issue", "create", "--kind", "task", "--title", "AC list limit", "--actor", "agent:developer")
	issueID := parseCreatedIssueID(t, createStderr)
	conn, err := fiddb.Open(filepath.Join(workspace, ".fid", "fid.sqlite"))
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	tx, err := conn.Begin()
	if err != nil {
		t.Fatalf("Begin() error = %v", err)
	}
	for ordinal := 1; ordinal <= 1001; ordinal++ {
		if _, err := tx.Exec(`INSERT INTO acceptance_criteria (id, issue_id, ordinal, text) VALUES (?, ?, ?, ?)`, fmt.Sprintf("AC-cli-limit-%04d", ordinal), issueID, ordinal, "criterion"); err != nil {
			t.Fatalf("insert criterion %d: %v", ordinal, err)
		}
	}
	if err := tx.Commit(); err != nil {
		t.Fatalf("Commit() error = %v", err)
	}
	closeCommandDB(t, conn)

	result := executeCommandForTest(workspace, "--json", "issue", "ac", "list", issueID)
	if result.err == nil || result.stdout != "" {
		t.Fatalf("result = %#v, want stderr-only limit rejection", result)
	}
	assertJSONFailureCode(t, result.stderr, "issue.ac.list", "result-limit-exceeded")
}

func parseAddedAcceptanceCriterionID(t *testing.T, output string) string {
	t.Helper()

	fields := strings.Fields(output)
	for index := range fields {
		if fields[index] == "AC" && index+1 < len(fields) {
			return fields[index+1]
		}
	}
	t.Fatalf("output = %q, want added AC line with id", output)
	return ""
}

func seedIssueAcceptanceCriterion(t *testing.T, workspace string, title string) (string, string) {
	t.Helper()

	_, createStderr := runFIDCommand(t, workspace, "issue", "create", "--kind", "task", "--title", title, "--actor", "agent:developer")
	issueID := parseCreatedIssueID(t, createStderr)
	addStdout, _ := runFIDCommand(
		t,
		workspace,
		"issue", "ac", "add", issueID,
		"--text", "Criterion for "+title,
		"--ordinal", "1",
		"--by", "developer",
	)
	return issueID, parseAddedAcceptanceCriterionID(t, addStdout)
}

func mustCanonicalAcceptanceCriteriaJSON(t *testing.T, output string) string {
	t.Helper()
	result := decodeCanonicalResult(t, output)
	criteria, ok := result["acceptance_criteria"]
	if !ok {
		t.Fatalf("issue ac list result = %#v, missing acceptance_criteria", result)
	}
	payload, err := json.Marshal(criteria)
	if err != nil {
		t.Fatalf("Marshal(acceptance criteria) error = %v", err)
	}
	return string(payload)
}
