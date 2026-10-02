package review

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/mbell-nzl/fid/internal/journal"
	"github.com/mbell-nzl/fid/internal/workflow"
	_ "modernc.org/sqlite"
)

func TestChangesRequiredReviewChargesExactlyOneReviewCycle(t *testing.T) {
	t.Parallel()
	conn, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("open review-cycle database: %v", err)
	}
	defer func() { _ = conn.Close() }()
	for _, statement := range []string{
		`CREATE TABLE envelope_budgets (envelope_id TEXT, revision INTEGER, name TEXT, budget_json TEXT)`,
		`CREATE TABLE budget_samples (id TEXT, envelope_id TEXT, envelope_revision INTEGER, budget_name TEXT, run_id TEXT, actor_id TEXT, source TEXT, source_sample_id TEXT, sample_kind TEXT, previous_sample_id TEXT, value INTEGER, measurement_ref_json TEXT, observed_at TEXT, created_event_id TEXT)`,
	} {
		if _, err := conn.Exec(statement); err != nil {
			t.Fatalf("prepare review-cycle database: %v", err)
		}
	}
	policy, err := json.Marshal(workflow.Budget{Name: "reviews", Unit: "review-cycles", Maximum: 2, Scope: "envelope", ExhaustionAction: "review-rework-exhausted"})
	if err != nil {
		t.Fatalf("marshal budget policy: %v", err)
	}
	if _, err := conn.Exec(`INSERT INTO envelope_budgets VALUES ('ENV-1',1,'reviews',?)`, string(policy)); err != nil {
		t.Fatalf("insert review-cycle budget: %v", err)
	}

	effects, budget, stops, attention, err := reviewCycleEffects(conn, recordBinding{EnvelopeID: "ENV-1", EnvelopeRevision: 1}, "RUN-review", "actor-reviewer", "REV-current", "changes-required", "2026-09-20T00:00:00Z")
	if err != nil {
		t.Fatalf("reviewCycleEffects() error = %v", err)
	}
	value, _ := effects[0].Values["value"].(*int64)
	if len(effects) != 1 || effects[0].Table != "budget_samples" || effects[0].Values["source_sample_id"] != "REV-current" || value == nil || *value != 1 {
		t.Fatalf("reviewCycleEffects() effects = %#v, want one Review-keyed charge", effects)
	}
	if budget == nil || budget.Consumed == nil || *budget.Consumed != 1 || len(stops) != 0 || len(attention) != 0 {
		t.Fatalf("reviewCycleEffects() budget/stops/attention = %#v/%#v/%#v", budget, stops, attention)
	}
	for _, verdict := range []string{"pass", "blocked"} {
		effects, budget, stops, attention, err := reviewCycleEffects(conn, recordBinding{EnvelopeID: "ENV-1", EnvelopeRevision: 1}, "RUN-"+verdict, "actor-reviewer", "REV-"+verdict, verdict, "2026-09-20T00:00:00Z")
		if err != nil || len(effects) != 0 || budget != nil || len(stops) != 0 || len(attention) != 0 {
			t.Fatalf("reviewCycleEffects(%s) = effects %#v budget %#v stops %#v attention %#v error %v, want no charge", verdict, effects, budget, stops, attention, err)
		}
	}
}

func TestTerminalFindingDispositionAuthorityMatrix(t *testing.T) {
	t.Parallel()
	reason := "no longer applies"
	for _, action := range []string{"resolved", "inapplicable"} {
		action := action
		for _, test := range []struct {
			name      string
			actor     string
			verdict   string
			creating  bool
			mutation  bool
			terminal  string
			wantError any
		}{
			{name: "creating Review", actor: "actor-new", verdict: "pass", creating: true, wantError: &FindingTargetInvalidError{}},
			{name: "mutating Claim holder", actor: "actor-builder", verdict: "pass", mutation: true, wantError: &NotIndependentError{}},
			{name: "blocked verdict", actor: "actor-new", verdict: "blocked", wantError: &VerdictInvalidError{}},
			{name: "different independent Actor", actor: "actor-new", verdict: "pass"},
			{name: "same independent Reviewer later", actor: "actor-creator", verdict: "pass"},
			{name: "second terminal", actor: "actor-new", verdict: "pass", terminal: action, wantError: &FindingAlreadyDisposedError{}},
		} {
			t.Run(action+"/"+test.name, func(t *testing.T) {
				conn := terminalMatrixDB(t, test.actor, test.mutation, test.terminal)
				defer func() { _ = conn.Close() }()
				request := terminalMatrixRequest(action, test.verdict, reason)
				if test.creating {
					request.Findings = []FindingInput{{ID: "FND-target"}}
				}
				err := checkIndependence(conn, test.actor, "RUN-current")
				if err == nil {
					err = validateTerminalFindingVerdict(request)
				}
				var effects []journal.Effect
				if err == nil {
					effects, err = resolveFindingEffects(conn, request, recordBinding{IssueID: "ISS-1", EnvelopeID: "ENV-1", EnvelopeRevision: 1}, test.actor, "REV-current")
				}
				if test.wantError == nil {
					if err != nil || len(effects) != 1 || effects[0].Values["action"] != action || effects[0].Values["actor_id"] != test.actor || effects[0].Values["qualifying_review_id"] != "REV-current" {
						t.Fatalf("terminal admission error/effects = %v/%#v", err, effects)
					}
					return
				}
				switch test.wantError.(type) {
				case *FindingTargetInvalidError:
					var target *FindingTargetInvalidError
					if !errors.As(err, &target) || target.ID != "FND-target" {
						t.Fatalf("terminal admission error = %T %v, want target error", err, err)
					}
				case *NotIndependentError:
					var independent *NotIndependentError
					if !errors.As(err, &independent) || independent.ConflictingType != "claim" {
						t.Fatalf("terminal admission error = %T %v, want independence error", err, err)
					}
				case *VerdictInvalidError:
					var verdict *VerdictInvalidError
					if !errors.As(err, &verdict) {
						t.Fatalf("terminal admission error = %T %v, want verdict error", err, err)
					}
				case *FindingAlreadyDisposedError:
					var disposed *FindingAlreadyDisposedError
					if !errors.As(err, &disposed) || disposed.State != action {
						t.Fatalf("terminal admission error = %T %v, want disposed %s", err, err, action)
					}
				}
			})
		}
	}
}

func TestValidateTerminalFindingInputBoundaries(t *testing.T) {
	t.Parallel()
	reason := "reason"
	resolved64 := make([]string, 64)
	for index := range resolved64 {
		resolved64[index] = fmt.Sprintf("FND-%02d", index)
	}
	if err := validateTerminalFindingInputs(resolved64, nil); err != nil {
		t.Fatalf("64 terminal Finding IDs rejected: %v", err)
	}
	conn, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("open 64-Finding database: %v", err)
	}
	defer func() { _ = conn.Close() }()
	for _, statement := range []string{
		`CREATE TABLE review_records (id TEXT, issue_id TEXT, envelope_id TEXT, envelope_revision INTEGER)`,
		`CREATE TABLE review_findings (id TEXT, review_id TEXT)`,
		`CREATE TABLE finding_dispositions (id TEXT, finding_id TEXT, action TEXT)`,
		`INSERT INTO review_records VALUES ('REV-prior','ISS-1','ENV-1',1)`,
	} {
		if _, err := conn.Exec(statement); err != nil {
			t.Fatalf("prepare 64-Finding database: %v", err)
		}
	}
	for _, id := range resolved64 {
		if _, err := conn.Exec(`INSERT INTO review_findings VALUES (?,'REV-prior')`, id); err != nil {
			t.Fatalf("insert Finding %s: %v", id, err)
		}
	}
	effects, err := resolveFindingEffects(conn, RecordRequest{ResolvesFindingIDs: resolved64}, recordBinding{IssueID: "ISS-1", EnvelopeID: "ENV-1", EnvelopeRevision: 1}, "actor-reviewer", "REV-current")
	if err != nil || len(effects) != 64 {
		t.Fatalf("resolveFindingEffects(64) = %d effects, error %v", len(effects), err)
	}
	for _, test := range []struct {
		name         string
		resolved     []string
		inapplicable []InapplicableFindingInput
	}{
		{name: "65 combined", resolved: append(append([]string{}, resolved64...), "FND-64")},
		{name: "duplicate resolved", resolved: []string{"FND-1", "FND-1"}},
		{name: "overlap", resolved: []string{"FND-1"}, inapplicable: []InapplicableFindingInput{{FindingID: "FND-1", Reason: &reason}}},
		{name: "missing reason", inapplicable: []InapplicableFindingInput{{FindingID: "FND-1"}}},
		{name: "blank reason", inapplicable: []InapplicableFindingInput{{FindingID: "FND-1", Reason: stringPointer(" \t")}}},
		{name: "oversize reason", inapplicable: []InapplicableFindingInput{{FindingID: "FND-1", Reason: stringPointer(string(make([]byte, 4097)))}}},
	} {
		t.Run(test.name, func(t *testing.T) {
			if err := validateTerminalFindingInputs(test.resolved, test.inapplicable); err == nil {
				t.Fatal("invalid terminal Finding input accepted")
			}
		})
	}
}

func terminalMatrixRequest(action, verdict, reason string) RecordRequest {
	request := RecordRequest{Verdict: verdict}
	if action == "resolved" {
		request.ResolvesFindingIDs = []string{"FND-target"}
	} else {
		request.InapplicableFindings = []InapplicableFindingInput{{FindingID: "FND-target", Reason: &reason}}
	}
	return request
}

func terminalMatrixDB(t *testing.T, actor string, mutation bool, terminal string) *sql.DB {
	t.Helper()
	conn, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("open terminal matrix database: %v", err)
	}
	for _, statement := range []string{
		`CREATE TABLE claims (id TEXT, issue_id TEXT, envelope_id TEXT, envelope_revision INTEGER, mode TEXT, actor_id TEXT, acquired_at TEXT)`,
		`CREATE TABLE runs (id TEXT, issue_id TEXT, envelope_id TEXT, envelope_revision INTEGER, claim_id TEXT, actor_id TEXT, started_at TEXT)`,
		`CREATE TABLE review_records (id TEXT, issue_id TEXT, envelope_id TEXT, envelope_revision INTEGER, actor_id TEXT)`,
		`CREATE TABLE review_findings (id TEXT, review_id TEXT)`,
		`CREATE TABLE finding_dispositions (id TEXT, finding_id TEXT, action TEXT)`,
		`INSERT INTO claims VALUES ('CLM-current','ISS-1','ENV-1',1,'review','actor-current','2026-09-20T00:00:00Z')`,
		`INSERT INTO runs VALUES ('RUN-current','ISS-1','ENV-1',1,'CLM-current','actor-current','2026-09-20T00:00:00Z')`,
		`INSERT INTO review_records VALUES ('REV-creator','ISS-1','ENV-1',1,'actor-creator')`,
		`INSERT INTO review_findings VALUES ('FND-target','REV-creator')`,
	} {
		if _, err := conn.Exec(statement); err != nil {
			t.Fatalf("prepare terminal matrix database: %v", err)
		}
	}
	if mutation {
		if _, err := conn.Exec(`INSERT INTO claims VALUES ('CLM-mutation','ISS-1','ENV-1',1,'mutation',?,'2026-09-20T00:00:01Z')`, actor); err != nil {
			t.Fatalf("insert mutating history: %v", err)
		}
	}
	if terminal != "" {
		if _, err := conn.Exec(`INSERT INTO finding_dispositions VALUES ('FDP-terminal','FND-target',?)`, terminal); err != nil {
			t.Fatalf("insert terminal disposition: %v", err)
		}
	}
	return conn
}

func stringPointer(value string) *string { return &value }

func TestValidateVerdictFindingAgreement(t *testing.T) {
	t.Parallel()
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
			if !test.valid {
				if _, ok := err.(*VerdictInvalidError); !ok {
					t.Fatalf("validateVerdictFindingAgreement() error = %T, want *VerdictInvalidError", err)
				}
			}
		})
	}
}

func TestValidateCoverageRejectsSourceInItemAndExclusion(t *testing.T) {
	t.Parallel()
	coverage := Coverage{
		Items:      []CoverageItem{{Source: Source{Digest: "sha256:same"}, Assessed: true}},
		Exclusions: []CoverageExclusion{{Source: Source{Digest: "sha256:same"}, Reason: "unavailable"}},
	}
	if err := validateCoverage(coverage); err == nil {
		t.Fatal("validateCoverage() accepted an item and exclusion sharing one Source digest")
	}
}

func TestCheckIndependenceRejectsClaimAndRunHistoryDirections(t *testing.T) {
	t.Parallel()
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
