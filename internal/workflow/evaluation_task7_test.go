package workflow

import (
	"encoding/json"
	"path/filepath"
	"strings"
	"testing"

	fiddb "github.com/mbell-nzl/fid/internal/db"
)

func TestHumanOutcomeSupersessionFollowsRecordedEventOrderBothDirections(t *testing.T) {
	conn, err := fiddb.Open(filepath.Join(t.TempDir(), "outcomes.sqlite"))
	if err != nil {
		t.Fatalf("open Outcome fixture: %v", err)
	}
	defer closeWorkflowDB(conn)
	if _, err := conn.Exec(`CREATE TABLE human_outcomes (id TEXT PRIMARY KEY,class TEXT,scope_issue_id TEXT,result TEXT,created_at TEXT,created_event_id TEXT UNIQUE)`); err != nil {
		t.Fatalf("create Outcome fixture: %v", err)
	}
	for _, statement := range []string{
		`INSERT INTO human_outcomes VALUES ('OUT-approved-first','release','ISS-1','approved','2026-09-20T23:00:00Z','EVT-20260920-000002')`,
		`INSERT INTO human_outcomes VALUES ('OUT-rejected-later','release','ISS-1','rejected','2026-09-20T01:00:00Z','EVT-20260920-000003')`,
		`INSERT INTO human_outcomes VALUES ('OUT-rejected-first','security','ISS-1','rejected','2026-09-20T23:00:00Z','EVT-20260920-000004')`,
		`INSERT INTO human_outcomes VALUES ('OUT-approved-later','security','ISS-1','approved','2026-09-20T01:00:00Z','EVT-20260920-000005')`,
	} {
		if _, err := conn.Exec(statement); err != nil {
			t.Fatalf("seed Outcome fixture: %v", err)
		}
	}

	for _, test := range []struct {
		class      string
		wantID     string
		wantResult string
	}{
		{class: "release", wantID: "OUT-rejected-later", wantResult: "rejected"},
		{class: "security", wantID: "OUT-approved-later", wantResult: "approved"},
	} {
		got, err := ResolveGoverningHumanOutcome(conn, test.class, "ISS-1")
		if err != nil || got.ID != test.wantID || got.Result != test.wantResult {
			t.Fatalf("resolve governing %s Outcome = %#v, %v; want %s/%s by Event order", test.class, got, err, test.wantID, test.wantResult)
		}
	}
}

func TestWorkflowEvaluationUsesHumanOutcomeFreshnessAndSupersession(t *testing.T) {
	conn, err := fiddb.Open(filepath.Join(t.TempDir(), "workflow-outcomes.sqlite"))
	if err != nil {
		t.Fatalf("open workflow Outcome fixture: %v", err)
	}
	defer closeWorkflowDB(conn)
	for _, statement := range []string{
		`CREATE TABLE evidence_records (id TEXT PRIMARY KEY,obligation_id TEXT,evidence_type TEXT,receipt_id TEXT,review_id TEXT,human_outcome_id TEXT,validity TEXT,target_repository_json TEXT)`,
		`CREATE TABLE execution_receipts (id TEXT PRIMARY KEY,terminal_status TEXT)`,
		`CREATE TABLE review_records (id TEXT PRIMARY KEY,verdict TEXT,issue_id TEXT,envelope_id TEXT,envelope_revision INTEGER)`,
		`CREATE TABLE human_outcomes (id TEXT PRIMARY KEY,class TEXT,scope_issue_id TEXT,result TEXT,created_at TEXT,created_event_id TEXT UNIQUE)`,
		`CREATE TABLE evidence_artifact_links (evidence_id TEXT,artifact_id TEXT,required_material INTEGER)`,
		`CREATE TABLE referenced_artifacts (id TEXT PRIMARY KEY,availability TEXT)`,
	} {
		if _, err := conn.Exec(statement); err != nil {
			t.Fatalf("create workflow Outcome fixture: %v", err)
		}
	}
	bound := Repository{CaptureSource: "git", HeadOID: strings.Repeat("a", 40), ProjectID: "PRJ-1", TreeOID: strings.Repeat("b", 40), WorktreeID: "WT-1"}
	changed := bound
	changed.TreeOID = strings.Repeat("c", 40)
	payload, err := json.Marshal(bound)
	if err != nil {
		t.Fatalf("marshal workflow Outcome repository: %v", err)
	}
	for _, statement := range []string{
		`INSERT INTO human_outcomes VALUES ('OUT-approved','release','ISS-1','approved','2026-09-20T23:00:00Z','EVT-20260920-000020')`,
		`INSERT INTO human_outcomes VALUES ('OUT-rejected','release','ISS-1','rejected','2026-09-20T01:00:00Z','EVT-20260920-000021')`,
	} {
		if _, err := conn.Exec(statement); err != nil {
			t.Fatalf("seed workflow Outcome fixture: %v", err)
		}
	}
	if _, err := conn.Exec(`INSERT INTO evidence_records VALUES ('EVD-approved','OBL-human','human-outcome',NULL,NULL,'OUT-approved','valid',?)`, string(payload)); err != nil {
		t.Fatalf("seed workflow Evidence fixture: %v", err)
	}
	states, err := loadObligationEvidenceStates(conn, "OBL-human", changed)
	if err != nil {
		t.Fatalf("load workflow Outcome state: %v", err)
	}
	class := "release"
	obligation := Obligation{ObligationInput: ObligationInput{ID: "OBL-human", Type: "human-outcome", OutcomeClass: &class, FreshnessTarget: "none", ArtifactPolicy: "metadata-sufficient"}}
	got := EvaluateEnvelopeObligations([]Obligation{obligation}, map[string][]ObligationEvidenceState{"OBL-human": states})[0]
	if got.Status != "unsatisfied" || got.Validity != "valid" || got.Freshness != "fresh" || len(got.Reasons) != 1 || got.Reasons[0].Code != "outcome-superseded" {
		t.Fatalf("workflow Outcome evaluation = %#v, want valid/fresh superseded result", got)
	}
}

func TestHumanOutcomeEvaluationUsesLaterDifferingResult(t *testing.T) {
	class := "release"
	obligation := Obligation{ObligationInput: ObligationInput{
		ID: "OBL-human", Type: "human-outcome", OutcomeClass: &class,
		FreshnessTarget: "none", ArtifactPolicy: "metadata-sufficient",
	}}
	tests := []struct {
		name       string
		state      ObligationEvidenceState
		wantStatus string
		wantReason string
	}{
		{
			name: "approved then rejected",
			state: ObligationEvidenceState{
				ObligationID: "OBL-human", EvidenceID: "EVD-approved", EvidenceType: "human-outcome",
				HumanOutcomeID: "OUT-approved", HumanOutcomeResult: "approved", OutcomeSuperseded: true,
				Validity: "valid", Freshness: "fresh", ArtifactAvailability: "not-applicable",
			},
			wantStatus: "unsatisfied", wantReason: "outcome-superseded",
		},
		{
			name: "rejected then approved",
			state: ObligationEvidenceState{
				ObligationID: "OBL-human", EvidenceID: "EVD-approved-later", EvidenceType: "human-outcome",
				HumanOutcomeID: "OUT-approved-later", HumanOutcomeResult: "approved",
				Validity: "valid", Freshness: "fresh", ArtifactAvailability: "not-applicable",
			},
			wantStatus: "satisfied",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := EvaluateEnvelopeObligations([]Obligation{obligation}, map[string][]ObligationEvidenceState{"OBL-human": {test.state}})[0]
			if got.Status != test.wantStatus || got.Validity != "valid" || got.Freshness != "fresh" || got.ArtifactAvailability != "not-applicable" {
				t.Fatalf("Human Outcome evaluation = %#v, want %s with independent dimensions", got, test.wantStatus)
			}
			if test.wantReason == "" {
				if len(got.Reasons) != 0 {
					t.Fatalf("Human Outcome reasons = %#v, want none", got.Reasons)
				}
			} else if len(got.Reasons) != 1 || got.Reasons[0].Code != test.wantReason || got.Reasons[0].SubjectID != test.state.HumanOutcomeID || got.Reasons[0].SubjectType != "human-outcome" {
				t.Fatalf("Human Outcome reasons = %#v, want %s for displaced Outcome", got.Reasons, test.wantReason)
			}
		})
	}
}

func TestAbsentHumanOutcomeFreshnessRemainsFresh(t *testing.T) {
	class := "release"
	obligation := Obligation{ObligationInput: ObligationInput{ID: "OBL-human", Type: "human-outcome", OutcomeClass: &class, FreshnessTarget: "none", ArtifactPolicy: "metadata-sufficient"}}
	got := EvaluateEnvelopeObligations([]Obligation{obligation}, nil)[0]
	if got.Status != "unsatisfied" || got.Validity != "invalid" || got.Freshness != "fresh" || got.ArtifactAvailability != "not-applicable" || len(got.Reasons) != 1 || got.Reasons[0].Code != "human-outcome-absent" {
		t.Fatalf("absent Human Outcome evaluation = %#v, want invalid/fresh/not-applicable with human-outcome-absent", got)
	}
}
