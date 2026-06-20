package main

import (
	"database/sql"
	"encoding/json"
	"os"
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

func TestIssueACVerifyCommandPersistsVerificationAndPrintsConfirmation(t *testing.T) {
	workspace := t.TempDir()
	runInitCommand(t, workspace, "init")

	_, createStderr := runFIDCommand(t, workspace, "issue", "create", "--kind", "task", "--title", "CLI AC verify", "--by", "developer")
	issueID := parseCreatedIssueID(t, createStderr)
	_, addStderr := runFIDCommand(
		t,
		workspace,
		"issue", "ac", "add", issueID,
		"--text", "Verification persists all fields",
		"--ordinal", "1",
		"--by", "developer",
	)
	criterionID := parseAddedAcceptanceCriterionID(t, addStderr)

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

	_, rebuildStderr := runFIDCommand(t, workspace, "rebuild", "--by", "workflow_runner")
	if !strings.Contains(rebuildStderr, "rebuild complete") {
		t.Fatalf("rebuild stderr = %q, want completion output", rebuildStderr)
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

func TestIssueACListCommandDefaultsToJSONAndMatchesIssueGetAcceptance(t *testing.T) {
	workspace := t.TempDir()
	runInitCommand(t, workspace, "init")

	issueID, criterionID := seedIssueAcceptanceCriterion(t, workspace, "CLI AC list")
	_, secondCriterionStderr := runFIDCommand(
		t,
		workspace,
		"issue", "ac", "add", issueID,
		"--text", "Second criterion for CLI AC list",
		"--ordinal", "2",
		"--by", "developer",
	)
	secondCriterionID := parseAddedAcceptanceCriterionID(t, secondCriterionStderr)

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

	jsonStdout, jsonStderr := runFIDCommand(t, workspace, "issue", "ac", "list", issueID, "--format", "json")
	if jsonStderr != "" {
		t.Fatalf("json stderr = %q, want empty", jsonStderr)
	}
	if defaultStdout != jsonStdout {
		t.Fatalf("default stdout = %q, want json stdout %q", defaultStdout, jsonStdout)
	}

	getStdout, getStderr := runFIDCommand(t, workspace, "issue", "get", issueID, "--section", "acceptance", "--format", "json")
	if getStderr != "" {
		t.Fatalf("issue get stderr = %q, want empty", getStderr)
	}

	listStructure := decodeCommandJSONStructure(t, jsonStdout)
	getStructure := decodeCommandJSONStructure(t, getStdout)
	if !commandStructuresEqual(listStructure, getStructure) {
		t.Fatalf("ac list structure = %#v, want issue get acceptance structure %#v", listStructure, getStructure)
	}

	rows, ok := listStructure.([]any)
	if !ok {
		t.Fatalf("ac list structure = %T, want []any", listStructure)
	}
	if len(rows) != 2 {
		t.Fatalf("ac list rows len = %d, want 2", len(rows))
	}
	firstRow := rows[0].(map[string]any)
	if firstRow["id"] != criterionID {
		t.Fatalf("first row id = %#v, want %q", firstRow["id"], criterionID)
	}
	secondRow := rows[1].(map[string]any)
	if secondRow["id"] != secondCriterionID {
		t.Fatalf("second row id = %#v, want %q", secondRow["id"], secondCriterionID)
	}
}

func TestIssueACListCommandDispatchesFormatsAndHandlesEmptyAndMissingIssues(t *testing.T) {
	workspace := t.TempDir()
	runInitCommand(t, workspace, "init")

	issueID, criterionID := seedIssueAcceptanceCriterion(t, workspace, "CLI AC format dispatch")
	runFIDCommand(t, workspace, "issue", "ac", "verify", issueID, criterionID, "--by", "developer", "--status", "pass")

	jsonStdout, jsonStderr := runFIDCommand(t, workspace, "issue", "ac", "list", issueID, "--format", "json")
	if jsonStderr != "" {
		t.Fatalf("json stderr = %q, want empty", jsonStderr)
	}
	want := decodeCommandJSONStructure(t, jsonStdout)

	for _, formatName := range []string{"toon", "yaml"} {
		t.Run(formatName, func(t *testing.T) {
			stdout, stderr := runFIDCommand(t, workspace, "issue", "ac", "list", issueID, "--format", formatName)
			if stderr != "" {
				t.Fatalf("stderr = %q, want empty", stderr)
			}

			var got any
			switch formatName {
			case "toon":
				got = decodeCommandTOONStructure(t, stdout)
			case "yaml":
				got = decodeCommandYAMLStructure(t, stdout)
			}
			if !commandStructuresEqual(want, got) {
				t.Fatalf("%s structure = %#v, want %#v", formatName, got, want)
			}
		})
	}

	t.Run("markdown", func(t *testing.T) {
		stdout, stderr := runFIDCommand(t, workspace, "issue", "ac", "list", issueID, "--format", "markdown")
		if stderr != "" {
			t.Fatalf("stderr = %q, want empty", stderr)
		}
		if !strings.Contains(stdout, "# Acceptance Criteria") {
			t.Fatalf("markdown stdout = %q, want heading", stdout)
		}
		if !strings.Contains(stdout, criterionID) {
			t.Fatalf("markdown stdout = %q, want criterion id", stdout)
		}
		if !strings.Contains(stdout, "developer: pass") {
			t.Fatalf("markdown stdout = %q, want latest verification status", stdout)
		}
	})

	t.Run("empty list", func(t *testing.T) {
		_, emptyIssueStderr := runFIDCommand(t, workspace, "issue", "create", "--kind", "task", "--title", "Empty AC list", "--by", "developer")
		emptyIssueID := parseCreatedIssueID(t, emptyIssueStderr)

		stdout, stderr := runFIDCommand(t, workspace, "issue", "ac", "list", emptyIssueID, "--format", "json")
		if stderr != "" {
			t.Fatalf("stderr = %q, want empty", stderr)
		}
		if got := decodeCommandJSONStructure(t, stdout); !commandStructuresEqual([]any{}, got) {
			t.Fatalf("empty json structure = %#v, want empty list", got)
		}
	})

	t.Run("missing issue", func(t *testing.T) {
		stdout, stderr, err := runFIDCommandExpectError(t, workspace, "issue", "ac", "list", "FID-missing", "--format", "json")
		if err == nil {
			t.Fatal("Execute() error = nil, want missing issue rejection")
		}
		if stdout != "" {
			t.Fatalf("stdout = %q, want empty", stdout)
		}
		if !strings.Contains(stderr, `error: issue "FID-missing" not found`) {
			t.Fatalf("stderr = %q, want missing issue error", stderr)
		}
	})
}

func parseAddedAcceptanceCriterionID(t *testing.T, stderr string) string {
	t.Helper()

	fields := strings.Fields(stderr)
	for index := range fields {
		if fields[index] == "AC" && index+1 < len(fields) {
			return fields[index+1]
		}
	}
	t.Fatalf("stderr = %q, want added AC line with id", stderr)
	return ""
}

func seedIssueAcceptanceCriterion(t *testing.T, workspace string, title string) (string, string) {
	t.Helper()

	_, createStderr := runFIDCommand(t, workspace, "issue", "create", "--kind", "task", "--title", title, "--by", "developer")
	issueID := parseCreatedIssueID(t, createStderr)
	_, addStderr := runFIDCommand(
		t,
		workspace,
		"issue", "ac", "add", issueID,
		"--text", "Criterion for "+title,
		"--ordinal", "1",
		"--by", "developer",
	)
	return issueID, parseAddedAcceptanceCriterionID(t, addStderr)
}

type commandAcceptanceCriterion struct {
	ID                  string                    `json:"id"`
	LatestVerifications map[string]map[string]any `json:"latest_verifications"`
}

func mustGetAcceptanceCriteria(t *testing.T, workspace, issueID string) []commandAcceptanceCriterion {
	t.Helper()

	stdout, stderr := runFIDCommand(t, workspace, "issue", "get", issueID, "--section", "acceptance", "--format", "json")
	if stderr != "" {
		t.Fatalf("issue get acceptance stderr = %q, want empty", stderr)
	}

	var got []commandAcceptanceCriterion
	if err := json.Unmarshal([]byte(stdout), &got); err != nil {
		t.Fatalf("Unmarshal() error = %v", err)
	}
	return got
}
