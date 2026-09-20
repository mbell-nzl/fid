package review

import (
	"database/sql"
	"testing"

	_ "modernc.org/sqlite"
)

func TestValidateVerdictFindingAgreement(t *testing.T) {
	tests := []struct {
		name     string
		verdict  string
		findings []FindingInput
		coverage Coverage
		valid    bool
	}{
		{name: "pass without findings", verdict: "pass", coverage: Coverage{Items: []CoverageItem{{Assessed: true}}}, valid: true},
		{name: "pass with blocker", verdict: "pass", coverage: Coverage{Items: []CoverageItem{{Assessed: true}}}, findings: []FindingInput{{Severity: "blocker", Blocking: true}}, valid: false},
		{name: "changes required with blocker", verdict: "changes-required", coverage: Coverage{Items: []CoverageItem{{Assessed: true}}}, findings: []FindingInput{{Severity: "major", Blocking: true}}, valid: true},
		{name: "changes required without blocker", verdict: "changes-required", coverage: Coverage{Items: []CoverageItem{{Assessed: true}}}, findings: []FindingInput{{Severity: "minor", Blocking: false}}, valid: false},
		{name: "blocked with exclusion", verdict: "blocked", coverage: Coverage{Items: []CoverageItem{{Assessed: true}}, Exclusions: []CoverageExclusion{{Reason: "unavailable"}}}, valid: true},
		{name: "blocked without finding or exclusion", verdict: "blocked", coverage: Coverage{Items: []CoverageItem{{Assessed: true}}}, valid: false},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := validateVerdictFindingAgreement(test.verdict, test.findings, test.coverage)
			if (err == nil) != test.valid {
				t.Fatalf("validateVerdictFindingAgreement() error = %v, valid = %v", err, test.valid)
			}
		})
	}
}

func TestValidateCoverageRejectsSourceInItemAndExclusion(t *testing.T) {
	coverage := Coverage{
		Items:      []CoverageItem{{Source: Source{Digest: "sha256:same"}, Assessed: true}},
		Exclusions: []CoverageExclusion{{Source: Source{Digest: "sha256:same"}, Reason: "unavailable"}},
	}
	if err := validateCoverage(coverage); err == nil {
		t.Fatal("validateCoverage() accepted an item and exclusion sharing one Source digest")
	}
}

func TestCheckIndependenceRejectsClaimAndRunHistoryDirections(t *testing.T) {
	conn, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("open test database: %v", err)
	}
	defer func() { _ = conn.Close() }()
	for _, statement := range []string{
		`CREATE TABLE claims (id TEXT, issue_id TEXT, envelope_id TEXT, envelope_revision INTEGER, mode TEXT, actor_id TEXT, acquired_at TEXT)`,
		`CREATE TABLE runs (id TEXT, issue_id TEXT, envelope_id TEXT, envelope_revision INTEGER, claim_id TEXT, actor_id TEXT, started_at TEXT)`,
	} {
		if _, err := conn.Exec(statement); err != nil {
			t.Fatalf("create test table: %v", err)
		}
	}
	if _, err := conn.Exec(`INSERT INTO runs VALUES ('RUN-review', 'ISS-1', 'ENV-1', 1, 'CLM-review', 'actor-reviewer', '2026-09-20T00:00:00Z')`); err != nil {
		t.Fatalf("insert reviewed Run: %v", err)
	}
	if _, err := conn.Exec(`INSERT INTO claims VALUES ('CLM-claim-holder', 'ISS-1', 'ENV-1', 1, 'mutation', 'actor-claim-holder', '2026-09-20T00:00:00Z')`); err != nil {
		t.Fatalf("insert Claim history: %v", err)
	}
	if err := assertIndependenceConflict(conn, "actor-claim-holder", "claim", "CLM-claim-holder"); err != nil {
		t.Fatal(err)
	}
	if _, err := conn.Exec(`INSERT INTO claims VALUES ('CLM-run-holder', 'ISS-1', 'ENV-1', 1, 'mutation', 'actor-claim-owner', '2026-09-20T00:01:00Z')`); err != nil {
		t.Fatalf("insert Run Claim: %v", err)
	}
	if _, err := conn.Exec(`INSERT INTO runs VALUES ('RUN-run-holder', 'ISS-1', 'ENV-1', 1, 'CLM-run-holder', 'actor-run-holder', '2026-09-20T00:01:00Z')`); err != nil {
		t.Fatalf("insert Run history: %v", err)
	}
	if err := assertIndependenceConflict(conn, "actor-run-holder", "run", "RUN-run-holder"); err != nil {
		t.Fatal(err)
	}
}

func assertIndependenceConflict(conn *sql.DB, actorID, subjectType, subjectID string) error {
	err := checkIndependence(conn, actorID, "RUN-review")
	conflict, ok := err.(*NotIndependentError)
	if !ok || conflict.ConflictingType != subjectType || conflict.ConflictingID != subjectID || conflict.ConflictingActorID != actorID {
		return &independenceAssertionError{err: err, actorID: actorID, subjectType: subjectType, subjectID: subjectID}
	}
	return nil
}

type independenceAssertionError struct {
	err                             error
	actorID, subjectType, subjectID string
}

func (err *independenceAssertionError) Error() string {
	return "independence conflict did not name the expected mutating subject"
}
