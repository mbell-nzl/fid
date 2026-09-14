package workflow

import (
	"database/sql"
	"testing"

	"github.com/mbell-nzl/fid/internal/journal"
)

func TestStopAndAttentionListsFilterLimitAndRemainReadOnly(t *testing.T) {
	workspace, migrationsDir := initializedWorkflowWorkspace(t)
	conn, err := openWorkflowDB(workspace)
	if err != nil {
		t.Fatalf("open workflow database: %v", err)
	}
	defer closeWorkflowDB(conn)
	commitStopReadFacts(t, workspace, conn)
	before := readWorkflowTestFile(t, workflowJournalPath(t, workspace))
	applicable, err := ApplicableStops(conn, "run.recover", Subject{Type: "run", ID: "RUN-1"})
	if err != nil || len(applicable) != 2 || applicable[0].ID != "STOP-1" || applicable[1].ID != "STOP-2" {
		t.Fatalf("ApplicableStops(run.recover, RUN-1) = %#v, %v; want both Stop facts", applicable, err)
	}
	unrelated, err := ApplicableStops(conn, "gate.run", Subject{Type: "run", ID: "RUN-1"})
	if err != nil || len(unrelated) != 0 {
		t.Fatalf("ApplicableStops(gate.run, RUN-1) = %#v, %v; want no unrelated Stop", unrelated, err)
	}

	stops, err := ListStops(StopListOptions{WorkspaceRoot: workspace, MigrationsDir: migrationsDir, Subject: "run:RUN-1", Status: "open", Limit: 1})
	if err != nil || len(stops) != 1 || stops[0].ID != "STOP-1" {
		t.Fatalf("ListStops(filtered limited) = %#v, %v; want STOP-1", stops, err)
	}
	attention, err := ListAttention(StopListOptions{WorkspaceRoot: workspace, MigrationsDir: migrationsDir, Subject: "run:RUN-1", Status: "open", Limit: 10})
	if err != nil || len(attention) != 1 || attention[0].ID != "ATT-1" || len(attention[0].StopIDs) != 2 {
		t.Fatalf("ListAttention(filtered) = %#v, %v; want linked ATT-1", attention, err)
	}
	if _, err := ListStops(StopListOptions{WorkspaceRoot: workspace, MigrationsDir: migrationsDir, Subject: "issue:ISS-1"}); err == nil {
		t.Fatal("ListStops(invalid subject) error = nil, want rejection")
	}
	if after := readWorkflowTestFile(t, workflowJournalPath(t, workspace)); string(after) != string(before) {
		t.Fatal("Stop/Attention reads appended a journal Event")
	}
}

func commitStopReadFacts(t *testing.T, workspace string, conn *sql.DB) {
	t.Helper()
	seedStandaloneStopSubject(t, conn)
	first := stopRaisedEvent("", "STOP-1", "ATT-1", true)
	second := stopRaisedEvent("", "STOP-2", "ATT-1", false)
	if err := journal.WithLock(workspace, journal.LockOwner{Command: "test stop list"}, func(lock *journal.WriteLock) error {
		prepared, err := journal.PrepareEvent(lock, first)
		if err != nil {
			return err
		}
		if _, err := journal.CommitEvent(lock, conn, stopRaisedEvent(prepared.EventID, "STOP-1", "ATT-1", true)); err != nil {
			return err
		}
		prepared, err = journal.PrepareEvent(lock, second)
		if err != nil {
			return err
		}
		_, err = journal.CommitEvent(lock, conn, stopRaisedEvent(prepared.EventID, "STOP-2", "ATT-1", false))
		return err
	}); err != nil {
		t.Fatalf("commit Stop facts: %v", err)
	}
}
