package evidence

import (
	"database/sql"
	"errors"
	"path/filepath"
	"testing"

	fiddb "github.com/mbell-nzl/fid/internal/db"
	"github.com/mbell-nzl/fid/internal/workflow"
)

// TestLoadHumanOutcomeEvidenceBindingRejectsClassMismatchAsOutcomeNotPermitted
// closes a class-mismatch gap in VAL-044/VAL-025: contracts.md section 10
// says a Human Outcome "may satisfy only a human-outcome Obligation whose
// outcome_class equals its class... Any other link rejects
// outcome-not-permitted naming the Obligation type." A same-type,
// approved-result Outcome offered to an Obligation with a *different*
// outcome_class is exactly that "any other link" case, so it must reject
// outcome-not-permitted like every other disallowed pairing - not the
// generic Evidence-binding-invalid error the pre-task-10 code path returned
// while Human Outcome linking was still unconditionally blocked upstream.
func TestLoadHumanOutcomeEvidenceBindingRejectsClassMismatchAsOutcomeNotPermitted(t *testing.T) {
	t.Parallel()
	conn := newHumanOutcomeEvidenceTestDB(t)
	execHumanOutcomeEvidenceSQL(t, conn, `INSERT INTO human_outcomes VALUES ('OUT-security','security','issue','ISS-1','approved','ISS-1')`)

	class := "release"
	obligation := workflow.ObligationInput{ID: "OBL-release", Type: "human-outcome", OutcomeClass: &class, AcceptedEvidenceTypes: []string{"human-outcome"}, AcceptedTrustSources: []string{"activated-human"}}
	obligation.Subject.Type, obligation.Subject.ID = "issue", "ISS-1"
	binding := evidenceBinding{EnvelopeID: "ENV-1", Revision: 1, IssueID: "ISS-1", Obligation: obligation}

	outcomeID := "OUT-security"
	request := EvidenceLinkRequest{ObligationID: "OBL-release", EvidenceType: "human-outcome", HumanOutcomeID: &outcomeID, TargetRepository: Repository{ProjectID: "PRJ-1"}, RequiredArtifactIDs: []string{}}
	err := loadHumanOutcomeEvidenceBinding(conn, request, &binding)
	var denied *OutcomeNotPermittedError
	if !errors.As(err, &denied) || denied.ObligationID != "OBL-release" || denied.ObligationType != "human-outcome" || denied.OutcomeID != "OUT-security" || denied.Result != "approved" {
		t.Fatalf("class-mismatched Human Outcome link error = %#v, want exact outcome-not-permitted denial", err)
	}
}

func newHumanOutcomeEvidenceTestDB(t *testing.T) *sql.DB {
	t.Helper()
	conn, err := fiddb.Open(filepath.Join(t.TempDir(), "evidence.sqlite"))
	if err != nil {
		t.Fatalf("open Evidence test DB: %v", err)
	}
	t.Cleanup(func() { _ = conn.Close() })
	execHumanOutcomeEvidenceSQL(t, conn, `CREATE TABLE human_outcomes (id TEXT PRIMARY KEY,class TEXT,subject_type TEXT,subject_id TEXT,result TEXT,scope_issue_id TEXT)`)
	return conn
}

func execHumanOutcomeEvidenceSQL(t *testing.T, conn *sql.DB, statement string, arguments ...any) {
	t.Helper()
	if _, err := conn.Exec(statement, arguments...); err != nil {
		t.Fatalf("execute Human Outcome Evidence fixture SQL: %v", err)
	}
}
