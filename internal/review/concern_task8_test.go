package review

import (
	"database/sql"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/mbell-nzl/fid/internal/journal"
	_ "modernc.org/sqlite"
)

func TestConcernDisposeCarriedForwardBindsExistingIssue(t *testing.T) {
	t.Parallel()
	conn, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("open target database: %v", err)
	}
	defer func() { _ = conn.Close() }()
	for _, statement := range []string{
		`CREATE TABLE issues (id TEXT PRIMARY KEY)`,
		`CREATE TABLE concerns (id TEXT PRIMARY KEY)`,
		`CREATE TABLE review_findings (id TEXT PRIMARY KEY)`,
		`INSERT INTO issues VALUES ('ISS-future-work')`,
	} {
		if _, err := conn.Exec(statement); err != nil {
			t.Fatalf("prepare target database: %v", err)
		}
	}

	request := ConcernDisposeRequest{
		ConcernID: "CON-source",
		Action:    "carried-forward",
		IssueID:   stringPointer("ISS-future-work"),
		Reason:    "tracked as durable work",
	}
	if err := validateConcernDispositionTargets(conn, request); err != nil {
		t.Fatalf("validateConcernDispositionTargets() error = %v", err)
	}
}

func TestConcernSubjectAndDispositionMatrix(t *testing.T) {
	t.Parallel()
	options := JudgmentOptions{WorkspaceRoot: "root", MigrationsDir: "migrations", ActorID: "actor-builder", IdempotencyKey: "key"}
	for _, subjectType := range []string{"issue", "envelope", "run", "review", "finding", "requirement"} {
		request := ConcernRecordRequest{ID: "CON-" + subjectType, Subject: Subject{Type: subjectType, ID: "unresolved-" + subjectType}, Scope: Scope{Paths: []string{}}, Severity: "blocker", Summary: "observation", Reasoning: "reason", References: []Source{}}
		if err := validateConcernRecord(options, request); err != nil {
			t.Fatalf("subject type %q rejected: %v", subjectType, err)
		}
	}

	conn := concernTargetMatrixDB(t)
	defer func() { _ = conn.Close() }()
	issueID, concernID, findingID := "ISS-existing", "CON-existing", "FND-existing"
	valid := []ConcernDisposeRequest{
		{ConcernID: "CON-source", Action: "carried-forward", IssueID: &issueID, Reason: "future work"},
		{ConcernID: "CON-source", Action: "rejected", Reason: "not a problem"},
		{ConcernID: "CON-source", Action: "superseded", SupersedesConcernID: &concernID, Reason: "new concern"},
		{ConcernID: "CON-source", Action: "superseded", SupersedesFindingID: &findingID, Reason: "review finding"},
	}
	for _, request := range valid {
		if err := validateConcernDispositionTargets(conn, request); err != nil {
			t.Fatalf("valid %s disposition rejected: %v", request.Action, err)
		}
	}

	missingIssue, missingConcern, missingFinding := "ISS-missing", "CON-missing", "FND-missing"
	invalid := []struct {
		name    string
		request ConcernDisposeRequest
		kind    any
	}{
		{name: "unknown Issue", request: ConcernDisposeRequest{ConcernID: "CON-source", Action: "carried-forward", IssueID: &missingIssue}, kind: &ConcernIssueNotFoundError{}},
		{name: "unknown Concern", request: ConcernDisposeRequest{ConcernID: "CON-source", Action: "superseded", SupersedesConcernID: &missingConcern}, kind: &ConcernNotFoundError{}},
		{name: "unknown Finding", request: ConcernDisposeRequest{ConcernID: "CON-source", Action: "superseded", SupersedesFindingID: &missingFinding}, kind: &FindingNotFoundError{}},
		{name: "self supersession", request: ConcernDisposeRequest{ConcernID: "CON-source", Action: "superseded", SupersedesConcernID: stringPointer("CON-source")}},
		{name: "outside closed action", request: ConcernDisposeRequest{ConcernID: "CON-source", Action: "addressed"}},
	}
	for _, test := range invalid {
		err := validateConcernDispositionTargets(conn, test.request)
		if err == nil {
			t.Fatalf("%s accepted", test.name)
		}
		if test.kind != nil && !errorHasType(err, test.kind) {
			t.Fatalf("%s error = %T %v, want %T", test.name, err, err, test.kind)
		}
	}
}

func TestConcernSelfDisposalKeepsRaisingAndDisposingActors(t *testing.T) {
	t.Parallel()
	value := map[string]any{"id": "CON-self", "actor_id": "actor-builder", "status": "open"}
	event := journal.Event{EventID: "EVT-dispose", OccurredAt: time.Date(2026, 9, 21, 1, 2, 3, 0, time.UTC)}
	got := disposedConcernValue(value, ConcernDisposeRequest{ConcernID: "CON-self", Action: "rejected", Reason: "cleared"}, "actor-builder", event)
	if got["actor_id"] != "actor-builder" || got["disposition"].(map[string]any)["actor_id"] != "actor-builder" {
		t.Fatalf("self disposal Actor fields = %#v", got)
	}
	if !reflect.DeepEqual(value, map[string]any{"id": "CON-self", "actor_id": "actor-builder", "status": "open"}) {
		t.Fatalf("source Concern mutated: %#v", value)
	}
}

func TestDeviationAcceptsExplicitPredecessorReference(t *testing.T) {
	t.Parallel()
	options := JudgmentOptions{WorkspaceRoot: "root", MigrationsDir: "migrations", ActorID: "actor-builder", IdempotencyKey: "key"}
	request := DeviationRecordRequest{ID: "DEV-next", Subject: Subject{Type: "requirement", ID: "REQ-unresolved"}, Expected: "accepted contract", Observed: "departure", Authority: DeviationAuthority{ActorID: "human-reviewer"}, Disposition: "accepted", Reasoning: "approved exception", References: []Source{}, SupersedesDeviationID: stringPointer("DEV-prior")}
	if err := validateDeviationRecord(options, request); err != nil {
		t.Fatalf("valid Deviation rejected: %v", err)
	}
}

func concernTargetMatrixDB(t *testing.T) *sql.DB {
	t.Helper()
	conn, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	for _, statement := range []string{
		`CREATE TABLE issues (id TEXT PRIMARY KEY)`,
		`CREATE TABLE concerns (id TEXT PRIMARY KEY)`,
		`CREATE TABLE review_findings (id TEXT PRIMARY KEY)`,
		`INSERT INTO issues VALUES ('ISS-existing')`,
		`INSERT INTO concerns VALUES ('CON-existing')`,
		`INSERT INTO review_findings VALUES ('FND-existing')`,
	} {
		if _, err := conn.Exec(statement); err != nil {
			t.Fatal(err)
		}
	}
	return conn
}

func errorHasType(err error, target any) bool {
	switch target.(type) {
	case *ConcernIssueNotFoundError:
		var value *ConcernIssueNotFoundError
		return errors.As(err, &value)
	case *ConcernNotFoundError:
		var value *ConcernNotFoundError
		return errors.As(err, &value)
	case *FindingNotFoundError:
		var value *FindingNotFoundError
		return errors.As(err, &value)
	}
	return false
}
