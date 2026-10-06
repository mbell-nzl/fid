package evidence

import (
	"database/sql"
	"encoding/json"
	"errors"
	"path/filepath"
	"reflect"
	"testing"

	fiddb "github.com/mbell-nzl/fid/internal/db"
	"github.com/mbell-nzl/fid/internal/workflow"
)

// TestEvidenceAdmissionTask6PureProductionClassifierMatrix is intermediate
// DEC-243 evidence for Human Outcome cells; it does not claim a durable link.
func TestEvidenceAdmissionTask6PureProductionClassifierMatrix(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name           string
		obligationType string
		acceptedType   string
		acceptedTrust  string
		evidenceType   string
		trustSource    string
		outcomeID      string
		outcomeResult  string
		wantCode       string
	}{
		{name: "receipt to gate evidence", obligationType: "gate-evidence", acceptedType: "receipt", acceptedTrust: "fid-executed", evidenceType: "receipt", trustSource: "fid-executed"},
		{name: "receipt to review", obligationType: "review", acceptedType: "review", acceptedTrust: "activated-actor", evidenceType: "receipt", trustSource: "fid-executed", wantCode: "evidence-type-not-accepted"},
		{name: "receipt to human outcome", obligationType: "human-outcome", acceptedType: "human-outcome", acceptedTrust: "activated-human", evidenceType: "receipt", trustSource: "fid-executed", wantCode: "evidence-type-not-accepted"},
		{name: "review to gate evidence", obligationType: "gate-evidence", acceptedType: "receipt", acceptedTrust: "fid-executed", evidenceType: "review", trustSource: "activated-actor", wantCode: "evidence-type-not-accepted"},
		{name: "review to review", obligationType: "review", acceptedType: "review", acceptedTrust: "activated-actor", evidenceType: "review", trustSource: "activated-actor"},
		{name: "review to human outcome", obligationType: "human-outcome", acceptedType: "human-outcome", acceptedTrust: "activated-human", evidenceType: "review", trustSource: "activated-actor", wantCode: "evidence-type-not-accepted"},
		{name: "human outcome to gate evidence", obligationType: "gate-evidence", acceptedType: "receipt", acceptedTrust: "fid-executed", evidenceType: "human-outcome", trustSource: "activated-human", outcomeID: "OUT-gate", outcomeResult: "approved", wantCode: "outcome-not-permitted"},
		{name: "human outcome to review", obligationType: "review", acceptedType: "review", acceptedTrust: "activated-actor", evidenceType: "human-outcome", trustSource: "activated-human", outcomeID: "OUT-review", outcomeResult: "approved", wantCode: "outcome-not-permitted"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := classifyEvidenceAdmission(evidenceAdmission{
				ObligationID:          "OBL-matrix",
				ObligationType:        test.obligationType,
				AcceptedEvidenceTypes: []string{test.acceptedType},
				AcceptedTrustSources:  []string{test.acceptedTrust},
				EvidenceType:          test.evidenceType,
				TrustSource:           test.trustSource,
				HumanOutcomeID:        test.outcomeID,
				HumanOutcomeResult:    test.outcomeResult,
			})
			switch test.wantCode {
			case "":
				if err != nil {
					t.Fatalf("classifyEvidenceAdmission() error = %v, want admission", err)
				}
			case "evidence-type-not-accepted":
				var mismatch *EvidenceTypeNotAcceptedError
				if !errors.As(err, &mismatch) || mismatch.ObligationID != "OBL-matrix" || mismatch.EvidenceType != test.evidenceType || !reflect.DeepEqual(mismatch.AcceptedEvidenceTypes, []string{test.acceptedType}) {
					t.Fatalf("classifyEvidenceAdmission() error = %#v, want exact type mismatch", err)
				}
			case "outcome-not-permitted":
				var denied *OutcomeNotPermittedError
				if !errors.As(err, &denied) || denied.ObligationType != test.obligationType || denied.OutcomeID != test.outcomeID || denied.Result != test.outcomeResult {
					t.Fatalf("classifyEvidenceAdmission() error = %#v, want exact Human Outcome denial", err)
				}
			}
		})
	}
}

func TestEvidenceAdmissionRejectsTrustMismatchWithAcceptedTypeDetails(t *testing.T) {
	t.Parallel()
	err := classifyEvidenceAdmission(evidenceAdmission{
		ObligationID:          "OBL-review",
		ObligationType:        "review",
		AcceptedEvidenceTypes: []string{"review"},
		AcceptedTrustSources:  []string{"different-trust"},
		EvidenceType:          "review",
		TrustSource:           "activated-actor",
	})
	var mismatch *EvidenceTypeNotAcceptedError
	if !errors.As(err, &mismatch) || !reflect.DeepEqual(mismatch.AcceptedEvidenceTypes, []string{"review"}) {
		t.Fatalf("trust mismatch = %#v, want evidence-type-not-accepted with accepted type list", err)
	}
}

func TestEvidenceWireCarriesExactlyOneNullableSubject(t *testing.T) {
	t.Parallel()
	reviewID := "REV-pass"
	record := Evidence{ID: "EVD-review", ObligationID: "OBL-review", EvidenceType: "review", ReviewID: &reviewID}
	payload, err := json.Marshal(record)
	if err != nil {
		t.Fatalf("marshal Evidence: %v", err)
	}
	var fields map[string]any
	if err := json.Unmarshal(payload, &fields); err != nil {
		t.Fatalf("decode Evidence: %v", err)
	}
	if fields["receipt_id"] != nil || fields["review_id"] != reviewID || fields["human_outcome_id"] != nil {
		t.Fatalf("Evidence subject fields = %#v, want null/review/null", fields)
	}
}

func TestEvidenceLinkInputRequiresOneSubjectMatchingType(t *testing.T) {
	t.Parallel()
	receiptID, reviewID := "RCP-1", "REV-1"
	base := EvidenceLinkOptions{WorkspaceRoot: "/workspace", MigrationsDir: "/migrations", ActorID: "ACT-1", IdempotencyKey: "key", Request: EvidenceLinkRequest{ObligationID: "OBL-1", EvidenceType: "receipt", ReceiptID: &receiptID, RequiredArtifactIDs: []string{}}}
	if err := validateEvidenceLinkOptions(base); err != nil {
		t.Fatalf("valid typed Evidence input = %v", err)
	}
	for _, test := range []struct {
		name   string
		mutate func(*EvidenceLinkRequest)
	}{
		{name: "missing type", mutate: func(request *EvidenceLinkRequest) { request.EvidenceType = "" }},
		{name: "wrong subject", mutate: func(request *EvidenceLinkRequest) { request.EvidenceType = "review" }},
		{name: "two subjects", mutate: func(request *EvidenceLinkRequest) { request.ReviewID = &reviewID }},
		{name: "blank subject", mutate: func(request *EvidenceLinkRequest) { blank := ""; request.ReceiptID = &blank }},
	} {
		t.Run(test.name, func(t *testing.T) {
			candidate := base
			test.mutate(&candidate.Request)
			if err := validateEvidenceLinkOptions(candidate); err == nil {
				t.Fatal("validateEvidenceLinkOptions() error = nil, want closed subject rejection")
			}
		})
	}
}

func TestReviewAdmissionReturnsBlockingFindingIDs(t *testing.T) {
	t.Parallel()
	for _, test := range []struct {
		verdict string
		ids     []string
	}{{verdict: "changes-required", ids: []string{"FND-1", "FND-2"}}, {verdict: "blocked", ids: []string{}}} {
		err := validateReviewEvidenceAdmission("OBL-review", "REV-"+test.verdict, test.verdict, test.ids)
		var changes *ReviewChangesRequiredError
		if !errors.As(err, &changes) || !reflect.DeepEqual(changes.BlockingFindingIDs, test.ids) {
			t.Fatalf("%s Review admission error = %#v, want ordered blocking Finding IDs %#v", test.verdict, err, test.ids)
		}
	}
	if err := validateReviewEvidenceAdmission("OBL-review", "REV-pass", "pass", nil); err != nil {
		t.Fatalf("pass Review admission error = %v", err)
	}
}

func TestLoadReviewEvidenceBindingEnforcesVerdictTypeAndLinkTimeIndependence(t *testing.T) {
	t.Parallel()
	conn := newReviewEvidenceTestDB(t)
	repository := task6EvidenceRepository()
	repositoryJSON, err := json.Marshal(repository)
	if err != nil {
		t.Fatalf("marshal repository: %v", err)
	}
	execReviewEvidenceSQL(t, conn, `INSERT INTO review_records VALUES ('REV-pass','RUN-review','ISS-1','ENV-1',1,'ACT-reviewer','standard',?,'pass')`, string(repositoryJSON))

	reviewID := "REV-pass"
	request := EvidenceLinkRequest{ObligationID: "OBL-review", EvidenceType: "review", ReviewID: &reviewID, RequiredArtifactIDs: []string{}}
	binding := reviewEvidenceBindingForTest("review")
	if err := loadReviewEvidenceBinding(conn, request, &binding); err != nil {
		t.Fatalf("load pass Review binding: %v", err)
	}

	mismatchBinding := reviewEvidenceBindingForTest("gate-evidence")
	request.ObligationID = "OBL-gate"
	err = loadReviewEvidenceBinding(conn, request, &mismatchBinding)
	var mismatch *EvidenceTypeNotAcceptedError
	if !errors.As(err, &mismatch) {
		t.Fatalf("review/gate mismatch = %#v, want evidence-type-not-accepted", err)
	}

	request.ObligationID = "OBL-review"
	execReviewEvidenceSQL(t, conn, `UPDATE review_records SET verdict='changes-required' WHERE id='REV-pass'`)
	execReviewEvidenceSQL(t, conn, `INSERT INTO review_findings VALUES ('FND-block','REV-pass',1)`)
	binding = reviewEvidenceBindingForTest("review")
	err = loadReviewEvidenceBinding(conn, request, &binding)
	var changes *ReviewChangesRequiredError
	if !errors.As(err, &changes) || !reflect.DeepEqual(changes.BlockingFindingIDs, []string{"FND-block"}) {
		t.Fatalf("changes-required Review = %#v, want exact blocker", err)
	}

	execReviewEvidenceSQL(t, conn, `UPDATE review_records SET verdict='pass' WHERE id='REV-pass'`)
	execReviewEvidenceSQL(t, conn, `DELETE FROM review_findings`)
	execReviewEvidenceSQL(t, conn, `INSERT INTO claims VALUES ('CLM-later','ISS-1','ENV-1',1,'mutation','ACT-reviewer','2026-09-20T00:00:00Z')`)
	binding = reviewEvidenceBindingForTest("review")
	err = loadReviewEvidenceBinding(conn, request, &binding)
	var independent *ReviewNotIndependentError
	if !errors.As(err, &independent) || independent.ConflictingType != "claim" || independent.ConflictingID != "CLM-later" || independent.ReviewActorID != "ACT-reviewer" {
		t.Fatalf("later mutating Claim = %#v, want exact link-time independence conflict", err)
	}
}

func newReviewEvidenceTestDB(t *testing.T) *sql.DB {
	t.Helper()
	conn, err := fiddb.Open(filepath.Join(t.TempDir(), "evidence.sqlite"))
	if err != nil {
		t.Fatalf("open Evidence test DB: %v", err)
	}
	t.Cleanup(func() { _ = conn.Close() })
	statements := []string{
		`CREATE TABLE review_records (id TEXT PRIMARY KEY,run_id TEXT,issue_id TEXT,envelope_id TEXT,envelope_revision INTEGER,actor_id TEXT,profile TEXT,repository_json TEXT,verdict TEXT)`,
		`CREATE TABLE review_findings (id TEXT PRIMARY KEY,review_id TEXT,blocking INTEGER)`,
		`CREATE TABLE finding_dispositions (finding_id TEXT,action TEXT)`,
		`CREATE TABLE claims (id TEXT PRIMARY KEY,issue_id TEXT,envelope_id TEXT,envelope_revision INTEGER,mode TEXT,actor_id TEXT,acquired_at TEXT)`,
		`CREATE TABLE runs (id TEXT PRIMARY KEY,claim_id TEXT,issue_id TEXT,envelope_id TEXT,envelope_revision INTEGER,actor_id TEXT,started_at TEXT)`,
	}
	for _, statement := range statements {
		execReviewEvidenceSQL(t, conn, statement)
	}
	return conn
}

func execReviewEvidenceSQL(t *testing.T, conn *sql.DB, statement string, arguments ...any) {
	t.Helper()
	if _, err := conn.Exec(statement, arguments...); err != nil {
		t.Fatalf("execute Evidence fixture SQL: %v", err)
	}
}

func reviewEvidenceBindingForTest(obligationType string) evidenceBinding {
	obligation := workflow.ObligationInput{ID: "OBL-review", Type: obligationType}
	obligation.Subject.Type, obligation.Subject.ID = "issue", "ISS-1"
	if obligationType == "review" {
		profile := "standard"
		obligation.ReviewProfile = &profile
		obligation.AcceptedEvidenceTypes = []string{"review"}
		obligation.AcceptedTrustSources = []string{"activated-actor"}
	} else {
		obligation.AcceptedEvidenceTypes = []string{"receipt"}
		obligation.AcceptedTrustSources = []string{"fid-executed"}
	}
	return evidenceBinding{EnvelopeID: "ENV-1", Revision: 1, IssueID: "ISS-1", Obligation: obligation}
}

func task6EvidenceRepository() Repository {
	return Repository{CaptureSource: "git", HeadOID: "0123456789abcdef0123456789abcdef01234567", ProjectID: "PRJ-1", TreeOID: "0123456789abcdef0123456789abcdef01234567", WorktreeID: "WT-1"}
}
