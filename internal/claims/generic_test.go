package claims

import (
	"database/sql"
	"encoding/json"
	"path/filepath"
	"testing"
	"time"

	"github.com/mbell-nzl/fid/internal/workflow"
	_ "modernc.org/sqlite"
)

func TestClaimAcquireRejectsOpenReviewReworkStop(t *testing.T) {
	conn, err := sql.Open("sqlite", filepath.Join(t.TempDir(), "claim.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = conn.Close() }()
	for _, statement := range []string{
		`CREATE TABLE stops (id TEXT,subject_type TEXT,subject_id TEXT,code TEXT,blocked_actions_json TEXT,details_json TEXT,resolution_condition_json TEXT,resolver_capability TEXT,actor_id TEXT,raised_at TEXT)`,
		`CREATE TABLE stop_resolutions (stop_id TEXT,actor_id TEXT,fact_event_id TEXT,resolved_at TEXT,created_event_id TEXT)`,
		`CREATE TABLE attention_stop_links (attention_id TEXT,stop_id TEXT)`,
		`CREATE TABLE envelope_budgets (envelope_id TEXT,revision INTEGER,name TEXT,budget_json TEXT)`,
		`CREATE TABLE budget_samples (id TEXT,envelope_id TEXT,envelope_revision INTEGER,budget_name TEXT,run_id TEXT,actor_id TEXT,source TEXT,source_sample_id TEXT,sample_kind TEXT,previous_sample_id TEXT,value INTEGER,measurement_ref_json TEXT,observed_at TEXT,created_event_id TEXT)`,
	} {
		if _, err := conn.Exec(statement); err != nil {
			t.Fatal(err)
		}
	}
	policy, _ := json.Marshal(workflow.Budget{Name: "reviews", Scope: "envelope", Unit: "review-cycles", Maximum: 1, ExhaustionAction: "review-rework-exhausted"})
	for _, statement := range []string{
		`INSERT INTO stops VALUES ('STP-review','envelope','ENV-1','review-rework-exhausted','["claim.acquire","run.start"]','{"envelope_id":"ENV-1","budget_name":"reviews"}','{"type":"budget-observed","subject":{"type":"envelope","id":"ENV-1"},"budget_name":"reviews","capability_id":null}','stop.resolve','actor-reviewer','2026-09-20T00:00:00Z')`,
		`INSERT INTO attention_stop_links VALUES ('ATT-review','STP-review')`,
	} {
		if _, err := conn.Exec(statement); err != nil {
			t.Fatal(err)
		}
	}
	if _, err := conn.Exec(`INSERT INTO envelope_budgets VALUES ('ENV-1',1,'reviews',?)`, string(policy)); err != nil {
		t.Fatal(err)
	}
	if _, err := conn.Exec(`INSERT INTO budget_samples VALUES ('BSM-review','ENV-1',1,'reviews','RUN-review','actor-reviewer','fid','REV-review','delta',NULL,1,NULL,'2026-09-20T00:00:00Z','EVT-review')`); err != nil {
		t.Fatal(err)
	}

	blocked, err := claimAcquireStopBlock(conn, AcquireRequest{EnvelopeID: "ENV-1"})
	if err != nil || blocked == nil || blocked.Stops[0].ID != "STP-review" || blocked.AttentionIDs[0] != "ATT-review" || blocked.Budget == nil || blocked.Budget.State != "exhausted" {
		t.Fatalf("claimAcquireStopBlock() = %#v, %v", blocked, err)
	}
}

func TestGenericClaimExpiresAtEqualityAndNeverBefore(t *testing.T) {
	expiresAt := time.Date(2026, time.September, 5, 12, 0, 0, 0, time.UTC)
	claim := GenericClaim{ExpiresAt: expiresAt.Format(time.RFC3339Nano)}
	for _, testCase := range []struct {
		name string
		now  time.Time
		want bool
	}{
		{name: "before expiry", now: expiresAt.Add(-time.Nanosecond), want: false},
		{name: "at expiry", now: expiresAt, want: true},
		{name: "after expiry", now: expiresAt.Add(time.Nanosecond), want: true},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			if got := claim.expiredAt(testCase.now); got != testCase.want {
				t.Fatalf("expiredAt(%s) = %t, want %t", testCase.now.Format(time.RFC3339Nano), got, testCase.want)
			}
		})
	}
}
