package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	fiddb "github.com/mbell-nzl/fid/internal/db"
)

func TestAcceptanceCriterionAddJSONMatchesContract(t *testing.T) {
	workspace := t.TempDir()
	clonePristineWorkspaceInto(t, workspace)

	created := executeCommandForTest(workspace,
		"--json", "issue", "create", "--kind", "task", "--title", "AC envelope target", "--by", "developer",
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
	assertNestedEnvelopeKeySet(t, []byte(result.stdout), "result", []string{"acceptance_criterion", "issue_revision"})
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
	workspace := t.TempDir()
	clonePristineWorkspaceInto(t, workspace)
	_, createStderr := runFIDCommand(t, workspace, "issue", "create", "--kind", "task", "--title", "AC JSON failures", "--by", "developer")
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
	workspace := t.TempDir()
	runInitCommand(t, workspace, "init")
	journalPath := mustJournalPath(t, workspace)

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
	if event.EventType != "issue.ac_added" {
		t.Fatalf("event type = %q, want issue.ac_added", event.EventType)
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
	if activityEvent != "ac_added" {
		t.Fatalf("activity_log event = %q, want ac_added", activityEvent)
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
	workspace := t.TempDir()
	runInitCommand(t, workspace, "init")

	_, createStderr := runFIDCommand(t, workspace, "issue", "create", "--kind", "task", "--title", "CLI AC verify", "--by", "developer")
	issueID := parseCreatedIssueID(t, createStderr)
	addStdout, _ := runFIDCommand(
		t,
		workspace,
		"issue", "ac", "add", issueID,
		"--text", "Verification persists all fields",
		"--ordinal", "1",
		"--by", "developer",
	)
	criterionID := parseAddedAcceptanceCriterionID(t, addStdout)

	stdout, stderr := runFIDCommand(
		t,
		workspace,
		"issue", "ac", "verify", issueID, criterionID,
		"--by", "developer",
		"--status", "pass",
		"--evidence-type", "test",
		"--evidence-ref", "cmd/fid/issue_ac_test.go::TestIssueACVerifyCommandPersistsVerificationAndPrintsConfirmation",
		"--notes", "CLI verification path",
	)
	if stdout != "" {
		t.Fatalf("stdout = %q, want empty", stdout)
	}
	if !strings.Contains(stderr, "verified "+criterionID+" on "+issueID+" by developer: pass") {
		t.Fatalf("stderr = %q, want verification confirmation", stderr)
	}
	if !strings.Contains(stderr, "evidence: test → cmd/fid/issue_ac_test.go::TestIssueACVerifyCommandPersistsVerificationAndPrintsConfirmation") {
		t.Fatalf("stderr = %q, want evidence line", stderr)
	}

	conn, err := fiddb.Open(filepath.Join(workspace, ".fid", "fid.sqlite"))
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	defer closeCommandDB(t, conn)

	var (
		verifierRole string
		status       string
		evidenceType sql.NullString
		evidenceRef  sql.NullString
		notes        sql.NullString
		completionID sql.NullString
	)
	if err := conn.QueryRow(`
		SELECT verifier_role, status, evidence_type, evidence_ref, notes, completion_id
		FROM ac_verifications
		WHERE ac_id = ?
	`, criterionID).Scan(&verifierRole, &status, &evidenceType, &evidenceRef, &notes, &completionID); err != nil {
		t.Fatalf("query ac_verifications: %v", err)
	}
	if verifierRole != "developer" || status != "pass" {
		t.Fatalf("stored verification = (%q, %q), want (developer, pass)", verifierRole, status)
	}
	if !evidenceType.Valid || evidenceType.String != "test" {
		t.Fatalf("stored evidence_type = %#v, want test", evidenceType)
	}
	if !evidenceRef.Valid || evidenceRef.String != "cmd/fid/issue_ac_test.go::TestIssueACVerifyCommandPersistsVerificationAndPrintsConfirmation" {
		t.Fatalf("stored evidence_ref = %#v, want targeted test ref", evidenceRef)
	}
	if !notes.Valid || notes.String != "CLI verification path" {
		t.Fatalf("stored notes = %#v, want CLI verification path", notes)
	}
	if completionID.Valid {
		t.Fatalf("stored completion_id = %#v, want NULL", completionID)
	}

	event := mustReadLastJournalEvent(t, mustJournalPath(t, workspace))
	if event.EventType != "issue.ac_verified" {
		t.Fatalf("event type = %q, want issue.ac_verified", event.EventType)
	}
	if got := countActivityLogRows(t, conn); got != 3 {
		t.Fatalf("activity_log rows after AC verify = %d, want 3", got)
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
	if activityEvent != "ac_verified" {
		t.Fatalf("activity_log event = %q, want ac_verified", activityEvent)
	}
	if actor != "developer" {
		t.Fatalf("activity_log actor = %q, want developer", actor)
	}
	for _, want := range []string{criterionID, `"verifier_role":"developer"`, `"status":"pass"`} {
		if !strings.Contains(details.String, want) {
			t.Fatalf("activity_log details = %q, want substring %q", details.String, want)
		}
	}
}

func TestIssueACVerifyCommandKeepsLatestPerRoleIndependent(t *testing.T) {
	workspace := t.TempDir()
	runInitCommand(t, workspace, "init")

	issueID, criterionID := seedIssueAcceptanceCriterion(t, workspace, "CLI multi-role verify")

	if _, stderr := runFIDCommand(
		t,
		workspace,
		"issue", "ac", "verify", issueID, criterionID,
		"--by", "developer",
		"--status", "pass",
	); !strings.Contains(stderr, "verified "+criterionID+" on "+issueID+" by developer: pass") {
		t.Fatalf("developer verify stderr = %q, want developer pass confirmation", stderr)
	}
	if _, stderr := runFIDCommand(
		t,
		workspace,
		"issue", "ac", "verify", issueID, criterionID,
		"--by", "qa_lead",
		"--status", "fail",
	); !strings.Contains(stderr, "verified "+criterionID+" on "+issueID+" by qa_lead: fail") {
		t.Fatalf("qa_lead verify stderr = %q, want qa_lead fail confirmation", stderr)
	}

	got := mustGetAcceptanceCriteria(t, workspace, issueID)
	if len(got) != 1 {
		t.Fatalf("acceptance criteria len = %d, want 1", len(got))
	}
	if got[0].LatestVerifications["developer"]["status"] != "pass" {
		t.Fatalf("developer latest verification = %#v, want pass", got[0].LatestVerifications["developer"])
	}
	if got[0].LatestVerifications["qa_lead"]["status"] != "fail" {
		t.Fatalf("qa_lead latest verification = %#v, want fail", got[0].LatestVerifications["qa_lead"])
	}
}

func TestIssueACVerifyCommandPreservesFailToPassHistory(t *testing.T) {
	workspace := t.TempDir()
	runInitCommand(t, workspace, "init")

	issueID, criterionID := seedIssueAcceptanceCriterion(t, workspace, "CLI fail then pass")

	if _, stderr := runFIDCommand(
		t,
		workspace,
		"issue", "ac", "verify", issueID, criterionID,
		"--by", "developer",
		"--status", "fail",
	); !strings.Contains(stderr, "verified "+criterionID+" on "+issueID+" by developer: fail") {
		t.Fatalf("developer fail stderr = %q, want fail confirmation", stderr)
	}

	afterFail := mustGetAcceptanceCriteria(t, workspace, issueID)
	if afterFail[0].LatestVerifications["developer"]["status"] != "fail" {
		t.Fatalf("developer latest verification after fail = %#v, want fail", afterFail[0].LatestVerifications["developer"])
	}

	if _, stderr := runFIDCommand(
		t,
		workspace,
		"issue", "ac", "verify", issueID, criterionID,
		"--by", "developer",
		"--status", "pass",
	); !strings.Contains(stderr, "verified "+criterionID+" on "+issueID+" by developer: pass") {
		t.Fatalf("developer pass stderr = %q, want pass confirmation", stderr)
	}

	afterPass := mustGetAcceptanceCriteria(t, workspace, issueID)
	if afterPass[0].LatestVerifications["developer"]["status"] != "pass" {
		t.Fatalf("developer latest verification after pass = %#v, want pass", afterPass[0].LatestVerifications["developer"])
	}

	conn, err := fiddb.Open(filepath.Join(workspace, ".fid", "fid.sqlite"))
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	defer closeCommandDB(t, conn)

	var count int
	if err := conn.QueryRow(`
		SELECT COUNT(*)
		FROM ac_verifications
		WHERE ac_id = ? AND verifier_role = ?
	`, criterionID, "developer").Scan(&count); err != nil {
		t.Fatalf("count ac_verifications: %v", err)
	}
	if count != 2 {
		t.Fatalf("developer ac_verifications rows = %d, want 2", count)
	}
}

func TestIssueACVerifyCommandLatestPerRoleSurvivesRebuild(t *testing.T) {
	workspace := t.TempDir()
	runInitCommand(t, workspace, "init")

	issueID, criterionID := seedIssueAcceptanceCriterion(t, workspace, "CLI rebuild AC parity")

	runFIDCommand(t, workspace, "issue", "ac", "verify", issueID, criterionID, "--by", "developer", "--status", "fail")
	_, developerPassStderr := runFIDCommand(t, workspace, "issue", "ac", "verify", issueID, criterionID, "--by", "developer", "--status", "pass")
	runFIDCommand(t, workspace, "issue", "ac", "verify", issueID, criterionID, "--by", "qa_lead", "--status", "fail")

	before := mustGetAcceptanceCriteria(t, workspace, issueID)
	if before[0].LatestVerifications["developer"]["status"] != "pass" {
		t.Fatalf("developer latest verification before rebuild = %#v, want pass", before[0].LatestVerifications["developer"])
	}
	if before[0].LatestVerifications["qa_lead"]["status"] != "fail" {
		t.Fatalf("qa_lead latest verification before rebuild = %#v, want fail", before[0].LatestVerifications["qa_lead"])
	}

	dbPath := filepath.Join(workspace, ".fid", "fid.sqlite")
	if err := os.Remove(dbPath); err != nil {
		t.Fatalf("Remove(%q) error = %v", dbPath, err)
	}

	rebuildStdout, rebuildStderr := runFIDCommand(t, workspace, "rebuild", "--by", "workflow_runner")
	if !strings.Contains(rebuildStdout, "rebuild complete") || rebuildStderr != "" {
		t.Fatalf("rebuild stdout/stderr = %q/%q, want stdout completion and empty stderr", rebuildStdout, rebuildStderr)
	}

	after := mustGetAcceptanceCriteria(t, workspace, issueID)
	if after[0].LatestVerifications["developer"]["status"] != "pass" {
		t.Fatalf("developer latest verification after rebuild = %#v, want pass", after[0].LatestVerifications["developer"])
	}
	if after[0].LatestVerifications["qa_lead"]["status"] != "fail" {
		t.Fatalf("qa_lead latest verification after rebuild = %#v, want fail", after[0].LatestVerifications["qa_lead"])
	}
	if gotBefore, gotAfter := before[0].LatestVerifications["developer"]["verified_at"], after[0].LatestVerifications["developer"]["verified_at"]; gotBefore != gotAfter {
		t.Fatalf("developer verified_at before/after rebuild = %q / %q, want identical", gotBefore, gotAfter)
	}
	if gotBefore, gotAfter := before[0].LatestVerifications["qa_lead"]["verified_at"], after[0].LatestVerifications["qa_lead"]["verified_at"]; gotBefore != gotAfter {
		t.Fatalf("qa_lead verified_at before/after rebuild = %q / %q, want identical", gotBefore, gotAfter)
	}
	if !strings.Contains(developerPassStderr, "verified "+criterionID+" on "+issueID+" by developer: pass") {
		t.Fatalf("developer pass stderr = %q, want pass confirmation", developerPassStderr)
	}
}

func TestIssueACVerifyCommandRejectsBeforeAppend(t *testing.T) {
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
				issueID, criterionID := seedIssueAcceptanceCriterion(t, workspace, "Missing by")
				return []string{"issue", "ac", "verify", issueID, criterionID, "--status", "pass"}
			},
			wantErr: "error: --by is required",
		},
		{
			name: "missing status flag",
			args: func(t *testing.T, workspace string) []string {
				issueID, criterionID := seedIssueAcceptanceCriterion(t, workspace, "Missing status")
				return []string{"issue", "ac", "verify", issueID, criterionID, "--by", "developer"}
			},
			wantErr: "error: --status is required",
		},
		{
			name: "invalid status",
			args: func(t *testing.T, workspace string) []string {
				issueID, criterionID := seedIssueAcceptanceCriterion(t, workspace, "Invalid status")
				return []string{"issue", "ac", "verify", issueID, criterionID, "--by", "developer", "--status", "bogus"}
			},
			wantErr: `error: verify acceptance criterion: status "bogus" must be one of pass, fail, blocked, not_applicable`,
		},
		{
			name: "invalid evidence type",
			args: func(t *testing.T, workspace string) []string {
				issueID, criterionID := seedIssueAcceptanceCriterion(t, workspace, "Invalid evidence type")
				return []string{"issue", "ac", "verify", issueID, criterionID, "--by", "developer", "--status", "pass", "--evidence-type", "pdf"}
			},
			wantErr: `error: verify acceptance criterion: evidence type "pdf" must be one of test, source, screenshot, manual, command, commit`,
		},
		{
			name: "unknown role",
			args: func(t *testing.T, workspace string) []string {
				issueID, criterionID := seedIssueAcceptanceCriterion(t, workspace, "Unknown role")
				return []string{"issue", "ac", "verify", issueID, criterionID, "--by", "intern", "--status", "pass"}
			},
			wantErr: `error: role "intern" not in project profile`,
		},
		{
			name: "missing acceptance criterion",
			args: func(t *testing.T, workspace string) []string {
				issueID, _ := seedIssueAcceptanceCriterion(t, workspace, "Missing AC")
				return []string{"issue", "ac", "verify", issueID, "AC-missing", "--by", "developer", "--status", "pass"}
			},
			wantErr: `error: verify acceptance criterion: acceptance criterion "AC-missing" not found`,
		},
		{
			name: "acceptance criterion belongs to different issue",
			args: func(t *testing.T, workspace string) []string {
				_, criterionID := seedIssueAcceptanceCriterion(t, workspace, "Criterion issue")
				otherIssueID, _ := seedIssueAcceptanceCriterion(t, workspace, "Other issue")
				return []string{"issue", "ac", "verify", otherIssueID, criterionID, "--by", "developer", "--status", "pass"}
			},
			wantErr: `error: verify acceptance criterion: acceptance criterion "`,
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

func TestIssueACVerifyCommandUsesRoleSpecificPermissionsAndHumanAttribution(t *testing.T) {
	workspace := t.TempDir()
	runInitCommand(t, workspace, "init")

	issueID, criterionID := seedIssueAcceptanceCriterion(t, workspace, "CLI verifier permissions")
	journalPath := mustJournalPath(t, workspace)
	linesBeforeDenied := len(splitJournalLines(mustReadCommandFile(t, journalPath)))

	stdout, stderr, err := runFIDCommandExpectError(
		t,
		workspace,
		"issue", "ac", "verify", issueID, criterionID,
		"--by", "workflow_runner",
		"--status", "pass",
	)
	if err == nil {
		t.Fatal("workflow_runner verify error = nil, want permission denial")
	}
	if stdout != "" {
		t.Fatalf("workflow_runner stdout = %q, want empty", stdout)
	}
	for _, want := range []string{
		`error: role "workflow_runner" not permitted for action "issue.ac.verify.workflow_runner"`,
		"check `fid policy show` for permitted role/action combinations",
	} {
		if !strings.Contains(stderr, want) {
			t.Fatalf("workflow_runner stderr = %q, want substring %q", stderr, want)
		}
	}
	if got := len(splitJournalLines(mustReadCommandFile(t, journalPath))); got != linesBeforeDenied {
		t.Fatalf("journal lines after workflow_runner denial = %d, want %d", got, linesBeforeDenied)
	}

	qaReviewerStdout, qaReviewerStderr := runFIDCommand(
		t,
		workspace,
		"issue", "ac", "verify", issueID, criterionID,
		"--by", "qa_reviewer",
		"--status", "fail",
	)
	if qaReviewerStdout != "" {
		t.Fatalf("qa_reviewer stdout = %q, want empty", qaReviewerStdout)
	}
	if !strings.Contains(qaReviewerStderr, "verified "+criterionID+" on "+issueID+" by qa_reviewer: fail") {
		t.Fatalf("qa_reviewer stderr = %q, want reviewer confirmation", qaReviewerStderr)
	}

	humanStdout, humanStderr := runFIDCommand(
		t,
		workspace,
		"issue", "ac", "verify", issueID, criterionID,
		"--by", "human",
		"--status", "pass",
	)
	if humanStdout != "" {
		t.Fatalf("human stdout = %q, want empty", humanStdout)
	}
	if !strings.Contains(humanStderr, "verified "+criterionID+" on "+issueID+" by human: pass") {
		t.Fatalf("human stderr = %q, want human confirmation", humanStderr)
	}

	got := mustGetAcceptanceCriteria(t, workspace, issueID)
	if len(got) != 1 {
		t.Fatalf("acceptance criteria len = %d, want 1", len(got))
	}
	if got[0].LatestVerifications["qa_reviewer"]["status"] != "fail" {
		t.Fatalf("qa_reviewer latest verification = %#v, want fail", got[0].LatestVerifications["qa_reviewer"])
	}
	if got[0].LatestVerifications["human"]["status"] != "pass" {
		t.Fatalf("human latest verification = %#v, want pass", got[0].LatestVerifications["human"])
	}
}

func TestIssueACListCommandMatchesClosedHumanAndJSONContracts(t *testing.T) {
	workspace := t.TempDir()
	runInitCommand(t, workspace, "init")

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

	runFIDCommand(t, workspace, "issue", "ac", "verify", issueID, criterionID, "--by", "developer", "--status", "fail")
	runFIDCommand(
		t,
		workspace,
		"issue", "ac", "verify", issueID, criterionID,
		"--by", "developer",
		"--status", "pass",
		"--evidence-type", "test",
		"--evidence-ref", "cmd/fid/issue_ac_test.go::TestIssueACListCommandDefaultsToJSONAndMatchesIssueGetAcceptance",
	)
	runFIDCommand(t, workspace, "issue", "ac", "verify", issueID, criterionID, "--by", "qa_lead", "--status", "blocked")

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
	workspace := t.TempDir()
	runInitCommand(t, workspace, "init")

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
		_, emptyIssueStderr := runFIDCommand(t, workspace, "issue", "create", "--kind", "task", "--title", "Empty AC list", "--by", "developer")
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
	workspace := t.TempDir()
	clonePristineWorkspaceInto(t, workspace)
	_, createStderr := runFIDCommand(t, workspace, "issue", "create", "--kind", "task", "--title", "AC list limit", "--by", "developer")
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

	_, createStderr := runFIDCommand(t, workspace, "issue", "create", "--kind", "task", "--title", title, "--by", "developer")
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

type commandAcceptanceCriterion struct {
	ID                  string                    `json:"id"`
	LatestVerifications map[string]map[string]any `json:"latest_verifications"`
}

func mustGetAcceptanceCriteria(t *testing.T, workspace, issueID string) []commandAcceptanceCriterion {
	t.Helper()

	stdout, stderr := runFIDCommand(t, workspace, "--json", "issue", "ac", "list", issueID)
	if stderr != "" {
		t.Fatalf("issue ac list stderr = %q, want empty", stderr)
	}
	stdout = mustCanonicalAcceptanceCriteriaJSON(t, stdout)

	var got []commandAcceptanceCriterion
	if err := json.Unmarshal([]byte(stdout), &got); err != nil {
		t.Fatalf("Unmarshal() error = %v", err)
	}
	conn, err := fiddb.Open(filepath.Join(workspace, ".fid", "fid.sqlite"))
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	defer closeCommandDB(t, conn)
	for index := range got {
		got[index].LatestVerifications = make(map[string]map[string]any)
		rows, err := conn.Query(`
			SELECT verifier_role, status, verified_at
			FROM ac_verifications
			WHERE ac_id = ?
			ORDER BY rowid
		`, got[index].ID)
		if err != nil {
			t.Fatalf("query verifications for %s: %v", got[index].ID, err)
		}
		for rows.Next() {
			var role, status, verifiedAt string
			if err := rows.Scan(&role, &status, &verifiedAt); err != nil {
				closeRows(t, rows)
				t.Fatalf("scan verification for %s: %v", got[index].ID, err)
			}
			got[index].LatestVerifications[role] = map[string]any{
				"status":      status,
				"verified_at": verifiedAt,
			}
		}
		if err := rows.Err(); err != nil {
			closeRows(t, rows)
			t.Fatalf("iterate verifications for %s: %v", got[index].ID, err)
		}
		closeRows(t, rows)
	}
	return got
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
