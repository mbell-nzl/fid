package evidence

import (
	"database/sql"
	"encoding/json"
	"path/filepath"
	"testing"

	fiddb "github.com/mbell-nzl/fid/internal/db"
)

// TestValidateGateRunReferencesAuthorizesAgainstEveryEnvelopeGrant is the
// review-7 I-5 regression (DEC-252). The Envelope carries one Grant per Actor;
// gate.run must be authorized by the Grant that names the running Actor, not
// only by whichever Grant row the store happens to return first.
func TestValidateGateRunReferencesAuthorizesAgainstEveryEnvelopeGrant(t *testing.T) {
	conn := newGateAuthorityTestDB(t)
	repository := task6EvidenceRepository()
	registration := Registration{ID: "REG-1", Version: 1, Scope: GateScope{IssueID: "ISS-1"}}
	gateRun := []string{"gate.run", "run.start"}
	grants := []struct {
		actor        string
		actions      []string
		capabilities []string
		issue        string
	}{
		{"actor-a", gateRun, gateRun, "ISS-1"},
		{"actor-b", gateRun, gateRun, "ISS-1"},
		{"actor-c", []string{"run.start"}, []string{"run.start"}, "ISS-1"},
		{"actor-d", gateRun, gateRun, "ISS-other"},
		{"actor-e", gateRun, []string{"run.start"}, "ISS-1"},
	}
	for ordinal, grant := range grants {
		execGateAuthoritySQL(t, conn, `INSERT INTO envelope_authority_grants VALUES ('ENV-1', 1, ?, ?, ?, ?)`, ordinal, mustJSON(t, map[string]any{"actor_id": grant.actor, "required_capabilities": grant.capabilities}), mustJSON(t, grant.actions), mustJSON(t, map[string]any{"issue_id": grant.issue}))
		execGateAuthoritySQL(t, conn, `INSERT INTO runs VALUES (?, 'ISS-1', 'ENV-1', 1, ?, 'active', ?)`, "RUN-"+grant.actor, grant.actor, mustJSON(t, repository))
	}
	execGateAuthoritySQL(t, conn, `INSERT INTO envelope_obligations VALUES ('OBL-1', 'ENV-1', 1, 'REG-1', 1)`)

	tests := []struct {
		actor   string
		granted bool
	}{
		{"actor-a", true},
		{"actor-b", true},
		{"actor-c", false},
		{"actor-d", false},
		{"actor-e", false},
	}
	for _, test := range tests {
		t.Run(test.actor, func(t *testing.T) {
			err := validateGateRunReferences(conn, registration, test.actor, GateRunRequest{RunID: "RUN-" + test.actor, ObligationID: "OBL-1", Repository: repository})
			if test.granted && err != nil {
				t.Fatalf("granted Actor %s = %v, want gate.run authorized by its own Grant", test.actor, err)
			}
			if !test.granted && (err == nil || err.Error() != "run gate: authority grant does not permit gate.run") {
				t.Fatalf("ungranted Actor %s = %v, want the unchanged authority denial", test.actor, err)
			}
		})
	}
}

func newGateAuthorityTestDB(t *testing.T) *sql.DB {
	t.Helper()
	conn, err := fiddb.Open(filepath.Join(t.TempDir(), "gate-authority.sqlite"))
	if err != nil {
		t.Fatalf("open gate authority test DB: %v", err)
	}
	t.Cleanup(func() { _ = conn.Close() })
	for _, statement := range []string{
		`CREATE TABLE envelope_obligations (id TEXT PRIMARY KEY, envelope_id TEXT, revision INTEGER, registration_id TEXT, registration_version INTEGER)`,
		`CREATE TABLE runs (id TEXT PRIMARY KEY, issue_id TEXT, envelope_id TEXT, envelope_revision INTEGER, actor_id TEXT, status TEXT, start_repository_json TEXT)`,
		`CREATE TABLE envelope_authority_grants (envelope_id TEXT, revision INTEGER, ordinal INTEGER, actor_selector_json TEXT, actions_json TEXT, scope_json TEXT, PRIMARY KEY (envelope_id, revision, ordinal))`,
	} {
		execGateAuthoritySQL(t, conn, statement)
	}
	return conn
}

func execGateAuthoritySQL(t *testing.T, conn *sql.DB, statement string, arguments ...any) {
	t.Helper()
	if _, err := conn.Exec(statement, arguments...); err != nil {
		t.Fatalf("execute gate authority fixture SQL: %v", err)
	}
}

func mustJSON(t *testing.T, value any) string {
	t.Helper()
	encoded, err := json.Marshal(value)
	if err != nil {
		t.Fatalf("marshal fixture JSON: %v", err)
	}
	return string(encoded)
}
