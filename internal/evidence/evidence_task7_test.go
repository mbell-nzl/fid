package evidence

import (
	"database/sql"
	"encoding/json"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	fiddb "github.com/mbell-nzl/fid/internal/db"
	"github.com/mbell-nzl/fid/internal/workflow"
)

func TestReviewEvidenceStalesOnTargetTreeChange(t *testing.T) {
	t.Parallel()
	bound := task6EvidenceRepository()
	changed := bound
	changed.TreeOID = strings.Repeat("f", 40)

	if got := freshnessForTarget("repository-tree", bound, changed); got != "stale" {
		t.Fatalf("Review freshness = %q, want stale after target tree change", got)
	}
	if got := freshnessForTarget("none", bound, changed); got != "fresh" {
		t.Fatalf("Human Outcome freshness = %q, want fresh after target tree change", got)
	}
}

func TestHumanOutcomeEvaluationKeepsFreshnessAndUsesSupersession(t *testing.T) {
	t.Parallel()
	conn := newTask7EvaluationDB(t)
	bound := task6EvidenceRepository()
	target := bound
	target.TreeOID = strings.Repeat("e", 40)
	repositoryJSON, err := json.Marshal(bound)
	if err != nil {
		t.Fatalf("marshal Outcome repository: %v", err)
	}
	class := "release"
	obligation := workflow.Obligation{ObligationInput: workflow.ObligationInput{
		ID: "OBL-human", Type: "human-outcome", OutcomeClass: &class, Subject: workflow.Subject{Type: "issue", ID: "ISS-1"},
		FreshnessTarget: "none", ArtifactPolicy: "metadata-sufficient",
	}}
	execTask7SQL(t, conn, `INSERT INTO envelope_obligations VALUES ('OBL-human','ENV-1',1)`)
	execTask7SQL(t, conn, `INSERT INTO human_outcomes VALUES ('OUT-approved','release','ISS-1','approved','2026-09-20T23:00:00Z','EVT-20260920-000010')`)
	execTask7SQL(t, conn, `INSERT INTO evidence_records VALUES ('EVD-approved','OBL-human','human-outcome',NULL,NULL,'OUT-approved','valid',?)`, string(repositoryJSON))

	before := task7OutcomeRows(t, conn)
	states, err := loadEvaluationStates(conn, obligation, target)
	if err != nil {
		t.Fatalf("load approved Outcome state: %v", err)
	}
	assertTask7OutcomeEvaluation(t, obligation, states, "satisfied", "")

	// The later Event governs even though its caller-authored timestamp sorts first.
	execTask7SQL(t, conn, `INSERT INTO human_outcomes VALUES ('OUT-rejected','release','ISS-1','rejected','2026-09-20T01:00:00Z','EVT-20260920-000011')`)
	afterInsert := task7OutcomeRows(t, conn)
	states, err = loadEvaluationStates(conn, obligation, target)
	if err != nil {
		t.Fatalf("load superseded Outcome state: %v", err)
	}
	assertTask7OutcomeEvaluation(t, obligation, states, "unsatisfied", "outcome-superseded")
	if after := task7OutcomeRows(t, conn); !reflect.DeepEqual(after, afterInsert) || !reflect.DeepEqual(after[:1], before) {
		t.Fatalf("Outcome evaluation changed readable records: before=%#v inserted=%#v after=%#v", before, afterInsert, after)
	}
}

func assertTask7OutcomeEvaluation(t *testing.T, obligation workflow.Obligation, states []workflow.ObligationEvidenceState, wantStatus, wantReason string) {
	t.Helper()
	got := workflow.EvaluateEnvelopeObligations([]workflow.Obligation{obligation}, map[string][]workflow.ObligationEvidenceState{obligation.ID: states})[0]
	if got.Status != wantStatus || got.Validity != "valid" || got.Freshness != "fresh" || got.ArtifactAvailability != "not-applicable" {
		t.Fatalf("Outcome evaluation = %#v, want %s with valid/fresh/not-applicable dimensions", got, wantStatus)
	}
	if wantReason == "" && len(got.Reasons) != 0 {
		t.Fatalf("Outcome reasons = %#v, want none", got.Reasons)
	}
	if wantReason != "" && (len(got.Reasons) != 1 || got.Reasons[0].Code != wantReason) {
		t.Fatalf("Outcome reasons = %#v, want %s", got.Reasons, wantReason)
	}
}

func newTask7EvaluationDB(t *testing.T) *sql.DB {
	t.Helper()
	conn, err := fiddb.Open(filepath.Join(t.TempDir(), "task7.sqlite"))
	if err != nil {
		t.Fatalf("open Task 7 fixture: %v", err)
	}
	t.Cleanup(func() { _ = conn.Close() })
	for _, statement := range []string{
		`CREATE TABLE envelope_obligations (id TEXT PRIMARY KEY,envelope_id TEXT,revision INTEGER)`,
		`CREATE TABLE evidence_records (id TEXT PRIMARY KEY,obligation_id TEXT,evidence_type TEXT,receipt_id TEXT,review_id TEXT,human_outcome_id TEXT,validity TEXT,target_repository_json TEXT)`,
		`CREATE TABLE execution_receipts (id TEXT PRIMARY KEY,terminal_status TEXT)`,
		`CREATE TABLE review_records (id TEXT PRIMARY KEY,verdict TEXT)`,
		`CREATE TABLE human_outcomes (id TEXT PRIMARY KEY,class TEXT,scope_issue_id TEXT,result TEXT,created_at TEXT,created_event_id TEXT UNIQUE)`,
		`CREATE TABLE evidence_artifact_links (evidence_id TEXT,artifact_id TEXT,required_material INTEGER)`,
		`CREATE TABLE referenced_artifacts (id TEXT PRIMARY KEY,availability TEXT)`,
	} {
		execTask7SQL(t, conn, statement)
	}
	return conn
}

func execTask7SQL(t *testing.T, conn *sql.DB, statement string, arguments ...any) {
	t.Helper()
	if _, err := conn.Exec(statement, arguments...); err != nil {
		t.Fatalf("execute Task 7 fixture SQL: %v", err)
	}
}

func task7OutcomeRows(t *testing.T, conn *sql.DB) [][]string {
	t.Helper()
	rows, err := conn.Query(`SELECT id,class,scope_issue_id,result,created_at,created_event_id FROM human_outcomes ORDER BY created_event_id`)
	if err != nil {
		t.Fatalf("read Outcome fixture: %v", err)
	}
	defer func() { _ = rows.Close() }()
	result := [][]string{}
	for rows.Next() {
		values := make([]string, 6)
		if err := rows.Scan(&values[0], &values[1], &values[2], &values[3], &values[4], &values[5]); err != nil {
			t.Fatalf("scan Outcome fixture: %v", err)
		}
		result = append(result, values)
	}
	return result
}
