package journal

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/mbell-nzl/fid/internal/config"
)

// The admission probes forge one Event after a writer-shaped prefix: a
// project, two tasks (one created under a retry key) and a hold.
var admissionAt = time.Date(2026, 11, 10, 10, 0, 0, 0, time.UTC)

func admissionStamp(seconds int) string {
	return admissionAt.Add(time.Duration(seconds) * time.Second).Format(time.RFC3339)
}

func admissionEvent(id, family, aggregateType, aggregateID string, sequence int, effects ...Effect) Event {
	return Event{EventID: id, EventSchema: EventSchemaV2, EventType: family, OccurredAt: admissionAt.Add(time.Minute),
		ActorV2:   EventActor{ActorID: "mark", ActorType: "human", Producer: "fid-standalone", ProvenanceRef: "project-local-attribution"},
		Aggregate: Aggregate{ID: aggregateID, Type: aggregateType, Sequence: sequence}, CommandName: "test.forge", Effects: effects}
}

func admissionActivity(eventID, id, issue, name string) Effect {
	return Effect{Op: "activity", Table: "activity_log", PK: map[string]any{"id": id}, Values: map[string]any{"journal_event_id": eventID, "ref_type": "issue", "ref_id": issue, "event": name, "details": nil, "actor": "human:mark", "occurred_at": admissionStamp(0)}}
}

func admissionPrefix() []Event {
	project := admissionEvent("EVT-20261110-000001", "project.initialized", "project", "project-1", 1)
	project.ActorV2 = EventActor{ActorID: "fid-bootstrap", ActorType: "service", Producer: "fid", ProvenanceRef: "fid:init"}
	project.CommandName = "project.init"
	project.Effects = []Effect{}
	events := []Event{project}
	for i, id := range []string{"FID-a", "FID-b"} {
		eventID := []string{"EVT-20261110-000002", "EVT-20261110-000003"}[i]
		created := admissionEvent(eventID, "issue.created", "issue", id, 1,
			Effect{Op: "insert", Table: "issues", PK: map[string]any{"id": id}, Values: map[string]any{"kind": "task", "title": "Task " + id, "workflow_id": "task-default", "status": "todo", "lane": "standard", "created_at": admissionStamp(0), "updated_at": admissionStamp(0)}},
			admissionActivity(eventID, "ACT-"+id, id, "created"))
		created.CommandName = "issue.create"
		if id == "FID-a" {
			request := IdempotencyRequest{OperationScope: "issue.create:project-1", IdempotencyKey: "create-a", RequestDigest: "sha256:" + strings.Repeat("a", 64)}
			created.Idempotency = &EventIdempotency{Key: request.IdempotencyKey, Scope: request.OperationScope, RequestDigest: request.RequestDigest}
			created.Effects = append(created.Effects, idempotencyResultEffect(request, created, id, `{"replayed":false}`))
		}
		events = append(events, created)
	}
	return events
}

// admissionLaneRaise is a valid issue.updated by an agent: it raises FID-a's
// lane and bumps its revision, as the update writer does.
func admissionLaneRaise() Event {
	event := admissionEvent("EVT-20261110-000004", "issue.updated", "issue", "FID-a", 2,
		Effect{Op: "update", Table: "issues", PK: map[string]any{"id": "FID-a"}, Values: map[string]any{"lane": "assured", "updated_at": admissionStamp(1)}, Preconditions: map[string]any{"updated_at": admissionStamp(0)}},
		Effect{Op: "update", Table: "issues", PK: map[string]any{"id": "FID-a"}, Values: map[string]any{"revision": 2}, Preconditions: map[string]any{"revision": 1}},
		admissionActivity("EVT-20261110-000004", "ACT-update", "FID-a", "updated"))
	event.ActorV2 = EventActor{ActorID: "builder", ActorType: "agent", Producer: "fid-standalone", ProvenanceRef: "project-local-attribution"}
	event.CommandName = "issue.update"
	return event
}

func admissionHoldTaken() Event {
	event := holdEvent("EVT-20261110-000004", "HLD-a", "FID-a", "one")
	event.OccurredAt = admissionAt.Add(time.Minute)
	return event
}

func TestForgedEventsRefusedOnEveryJournalPath(t *testing.T) {
	t.Parallel()
	legacy := func(event Event) Event {
		event.EventSchema = EventSchemaV1
		event.Actor = Actor{Role: "developer"}
		event.Command = Command{Argv: []string{"fid", "issue", "update"}, CWD: ".", FIDVersion: "0.1.0"}
		return event
	}
	cases := []struct {
		name   string
		reason string
		forge  func() Event
	}{
		{"table case alias lowers a lane", "lowercase", func() Event {
			e := admissionLaneRaise()
			e.Effects[0].Table, e.Effects[0].Values["lane"] = "ISSUES", "routine"
			return e
		}},
		{"column case alias lowers a lane", "lowercase", func() Event {
			e := admissionLaneRaise()
			delete(e.Effects[0].Values, "lane")
			e.Effects[0].Values["LANE"] = "routine"
			return e
		}},
		{"key case alias", "lowercase", func() Event {
			e := admissionLaneRaise()
			e.Effects[1].PK = map[string]any{"ID": "FID-a"}
			return e
		}},
		{"update moves another issue", "aggregate", func() Event {
			e := admissionLaneRaise()
			e.Effects[0].PK["id"], e.Effects[1].PK["id"] = "FID-b", "FID-b"
			return e
		}},
		{"migration state", "Effect grammar", func() Event {
			e := admissionLaneRaise()
			e.Effects = append(e.Effects[:2:2], Effect{Op: "update", Table: "schema_migrations", PK: map[string]any{"version": 35}, Values: map[string]any{"dirty": true}}, e.Effects[2])
			return e
		}},
		{"projection checkpoint", "Effect grammar", func() Event {
			e := admissionLaneRaise()
			e.Effects = append(e.Effects[:2:2], Effect{Op: "update", Table: "projection_meta", PK: map[string]any{"key": "default"}, Values: map[string]any{"replay_status": "clean"}}, e.Effects[2])
			return e
		}},
		{"deletes another command's retry result", "Effect grammar", func() Event {
			e := admissionLaneRaise()
			// The update family reads result_json from every retry Effect.
			return withEffect(e, Effect{Op: "delete", Table: "idempotency_results", PK: map[string]any{"operation_scope": "issue.create:project-1", "idempotency_key": "create-a"}, Values: map[string]any{"result_json": `{"changes":[]}`}})
		}},
		{"retry result without its request", "idempotency", func() Event {
			e := admissionLaneRaise()
			request := IdempotencyRequest{OperationScope: "task.start:project-1", IdempotencyKey: "forged", RequestDigest: "sha256:" + strings.Repeat("b", 64)}
			e.Effects = append(e.Effects, idempotencyResultEffect(request, e, "FID-a", `{"replayed":false}`))
			return e
		}},
		{"hold taken by a service", "actor", func() Event {
			e := admissionHoldTaken()
			e.ActorV2.ActorType = "service"
			return e
		}},
		{"oversized title", "1 MiB line budget", func() Event {
			e := admissionLaneRaise()
			e.Effects[0].Values["title"] = strings.Repeat("<", 200000)
			return e
		}},
		{"legacy policy snapshot", "Event-v1", func() Event {
			return legacy(admissionEvent("EVT-20261110-000004", "policy.snapshotted", "policy", "policy", 1, Effect{Op: "insert", Table: "policy_snapshots", PK: map[string]any{"sequence": 1}, Values: map[string]any{"source_ref": "built-in defaults", "source_commit": nil, "digest": "sha256:" + strings.Repeat("c", 64), "effective_json": "{}"}}))
		}},
		{"legacy retry result", "Event-v1", func() Event {
			e := admissionLaneRaise()
			request := IdempotencyRequest{OperationScope: "task.start:project-1", IdempotencyKey: "forged", RequestDigest: "sha256:" + strings.Repeat("b", 64)}
			return legacy(admissionEvent("EVT-20261110-000004", "issue.updated", "issue", "FID-a", 2, idempotencyResultEffect(request, e, "FID-a", `{"replayed":false}`)))
		}},
		{"legacy migration state", "Event-v1", func() Event {
			return legacy(admissionEvent("EVT-20261110-000004", "issue.updated", "issue", "FID-a", 2, Effect{Op: "update", Table: "schema_migrations", PK: map[string]any{"version": 35}, Values: map[string]any{"dirty": true}}))
		}},
		{"legacy oversized title", "1 MiB line budget", func() Event {
			e := legacy(admissionLaneRaise())
			e.Effects[0].Values["title"] = strings.Repeat("<", 200000)
			return e
		}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			assertForgedRefusedEveryJournalPath(t, admissionPrefix(), tc.forge(), tc.reason)
		})
	}
}

func withEffect(event Event, effect Effect) Event {
	last := event.Effects[len(event.Effects)-1]
	event.Effects = append(append(event.Effects[:len(event.Effects)-1:len(event.Effects)-1], effect), last)
	return event
}

// The unforged Events pass every path, so each refusal above is the forgery's.
func TestAdmissionControlsAcceptedOnEveryJournalPath(t *testing.T) {
	t.Parallel()
	for name, event := range map[string]Event{"lane raise": admissionLaneRaise(), "hold taken": admissionHoldTaken()} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			for _, mode := range admissionModes(event) {
				t.Run(mode, func(t *testing.T) {
					w := newAdmissionWorkspace(t, admissionPrefix(), mode)
					if err := w.admit(t, mode, event); err != nil {
						t.Fatalf("%s refused valid Event: %v", mode, err)
					}
				})
			}
		})
	}
}

func admissionModes(event Event) []string {
	modes := []string{"dry-run", "project", "replay", "catch-up", "freshen", "rebuild"}
	if event.IsV2() {
		modes = append([]string{"commit"}, modes...)
	}
	return modes
}

type admissionWorkspace struct {
	root, journal, file string
	prefix              []byte
}

// newAdmissionWorkspace writes the prefix to the journal and, except for a
// full replay, projects it first.
func newAdmissionWorkspace(t *testing.T, prefix []Event, mode string) admissionWorkspace {
	t.Helper()
	root := newReplayWorkspace(t)
	dir := filepath.Join(config.LedgerDir(root), "journal")
	file := filepath.Join(dir, "2026-11-10.jsonl")
	writeReplayEvents(t, file, prefix...)
	conn := openReplayDB(t, filepath.Join(config.LedgerDir(root), "fid.sqlite"))
	defer closeProjectDB(t, conn)
	if mode != "replay" {
		if _, err := ReplayJournal(conn, dir, nil); err != nil {
			t.Fatalf("replay prefix: %v", err)
		}
	}
	raw, err := os.ReadFile(file)
	if err != nil {
		t.Fatal(err)
	}
	return admissionWorkspace{root, dir, file, raw}
}

// admit offers event to one admission path. Replay paths read it from the
// journal after the prefix; live paths offer it to the current projection.
func (w admissionWorkspace) admit(t *testing.T, mode string, event Event) error {
	t.Helper()
	conn := openReplayDB(t, filepath.Join(config.LedgerDir(w.root), "fid.sqlite"))
	defer closeProjectDB(t, conn)
	switch mode {
	case "commit":
		var err error
		lockErr := WithLedgerLock(config.LedgerDir(w.root), LockOwner{Command: "test"}, func(lock *WriteLock) error {
			_, err = CommitEvent(lock, conn, event)
			return nil
		})
		if lockErr != nil {
			t.Fatal(lockErr)
		}
		return err
	case "dry-run":
		return DryRunProjectEvent(conn, event)
	case "project":
		hash, err := HashEvent(event)
		if err != nil {
			return err
		}
		return ProjectEvent(conn, ProjectionInput{Event: event, EventHash: hash, File: filepath.Base(w.file), Line: 4, Offset: int64(len(w.prefix)) + 1})
	}
	line, err := MarshalCanonicalJSON(event)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(w.file, append(append(append([]byte{}, w.prefix...), line...), '\n'), 0o644); err != nil {
		t.Fatal(err)
	}
	return w.replayStored(t, mode)
}

// replayStored reads the journal as it stands through one replay path.
func (w admissionWorkspace) replayStored(t *testing.T, mode string) error {
	t.Helper()
	conn := openReplayDB(t, filepath.Join(config.LedgerDir(w.root), "fid.sqlite"))
	defer closeProjectDB(t, conn)
	switch mode {
	case "replay":
		// The pre-scan refuses before any Event of the journal lands.
		_, err := ReplayJournal(conn, w.journal, nil)
		return err
	case "catch-up":
		_, err := ReplayJournal(conn, w.journal, &ReplayOptions{StartFile: filepath.Base(w.file), StartOffset: int64(len(w.prefix))})
		return err
	case "freshen":
		_, err := EnsureProjectionFresh(w.root, testProjectMigrationsDir(t))
		return err
	case "rebuild":
		_, err := RebuildWorkspace(w.root, testProjectMigrationsDir(t))
		return err
	}
	t.Fatalf("unknown admission mode %q", mode)
	return nil
}

func assertForgedRefusedEveryJournalPath(t *testing.T, prefix []Event, forged Event, reason string) {
	t.Helper()
	for _, mode := range admissionModes(forged) {
		t.Run(mode, func(t *testing.T) {
			w := newAdmissionWorkspace(t, prefix, mode)
			conn := openReplayDB(t, filepath.Join(config.LedgerDir(w.root), "fid.sqlite"))
			before := currentProjectionSnapshot(t, conn)
			closeProjectDB(t, conn)
			err := w.admit(t, mode, forged)
			if err == nil || !strings.Contains(err.Error(), reason) {
				t.Fatalf("%s = %v, want refusal naming %q", mode, err, reason)
			}
			conn = openReplayDB(t, filepath.Join(config.LedgerDir(w.root), "fid.sqlite"))
			defer closeProjectDB(t, conn)
			if !reflect.DeepEqual(before, currentProjectionSnapshot(t, conn)) {
				t.Fatalf("%s refusal changed the projection", mode)
			}
			journal, readErr := os.ReadFile(w.file)
			if readErr != nil {
				t.Fatal(readErr)
			}
			if mode == "commit" || mode == "dry-run" || mode == "project" {
				if !reflect.DeepEqual(journal, w.prefix) {
					t.Fatalf("%s refusal changed the journal", mode)
				}
				// The ledger still takes a valid write after the refusal.
				if err := w.admit(t, "commit", admissionLaneRaise()); err != nil {
					t.Fatalf("valid write after %s refusal: %v", mode, err)
				}
			}
		})
	}
}

// The writers' scanners read each stored line, so a line padded past the
// budget would stop every later write even though its Event is small.
func TestReplayRefusesStoredLineOverBudget(t *testing.T) {
	t.Parallel()
	for _, mode := range []string{"replay", "catch-up", "freshen", "rebuild"} {
		t.Run(mode, func(t *testing.T) {
			t.Parallel()
			w := newAdmissionWorkspace(t, admissionPrefix(), mode)
			line, err := MarshalCanonicalJSON(admissionLaneRaise())
			if err != nil {
				t.Fatal(err)
			}
			padded := append(append([]byte("{"), bytes.Repeat([]byte(" "), journalLineBudget)...), line[1:]...)
			if err := os.WriteFile(w.file, append(append(append([]byte{}, w.prefix...), padded...), '\n'), 0o644); err != nil {
				t.Fatal(err)
			}
			conn := openReplayDB(t, filepath.Join(config.LedgerDir(w.root), "fid.sqlite"))
			before := currentProjectionSnapshot(t, conn)
			closeProjectDB(t, conn)
			err = w.replayStored(t, mode)
			var corrupt *CorruptJournalError
			if !errors.As(err, &corrupt) || corrupt.Line != len(admissionPrefix())+1 || !strings.Contains(err.Error(), "1 MiB line budget") {
				t.Fatalf("%s = %v, want the padded line refused", mode, err)
			}
			conn = openReplayDB(t, filepath.Join(config.LedgerDir(w.root), "fid.sqlite"))
			defer closeProjectDB(t, conn)
			if !reflect.DeepEqual(before, currentProjectionSnapshot(t, conn)) {
				t.Fatalf("%s refusal changed the projection", mode)
			}
		})
	}
}

// The writers allocate each Event ID by reading every dated ID of the day.
// One they cannot read, or that leaves no next sequence, would stop every
// write on that date.
func TestForgedEventIdentityRefusedOnEveryJournalPath(t *testing.T) {
	t.Parallel()
	for _, id := range []string{"EVT-20261110-00000x", "EVT-20261110-", "EVT-20261110-000000", "EVT-20261110-9223372036854775807"} {
		t.Run(id, func(t *testing.T) {
			t.Parallel()
			forged := admissionLaneRaise()
			forged.EventID = id
			forged.Effects[2].Values["journal_event_id"] = id
			assertForgedRefusedEveryJournalPath(t, admissionPrefix(), forged, "Event ID")
		})
	}
}

// Export refuses a journal that repeats an Event ID, so replay must too, or
// every later backup would fail. Like aggregate sequences, the check spans
// lines, so it runs as replay reads them: rebuild replays into a scratch
// projection and catch-up stops at the repeat, changing nothing.
func TestReplayRefusesRepeatedEventID(t *testing.T) {
	t.Parallel()
	for _, mode := range []string{"catch-up", "freshen", "rebuild"} {
		t.Run(mode, func(t *testing.T) {
			t.Parallel()
			w := newAdmissionWorkspace(t, admissionPrefix(), mode)
			repeated := admissionLaneRaise()
			repeated.EventID = admissionPrefix()[1].EventID
			repeated.Effects[2].Values["journal_event_id"] = repeated.EventID
			conn := openReplayDB(t, filepath.Join(config.LedgerDir(w.root), "fid.sqlite"))
			before := currentProjectionSnapshot(t, conn)
			closeProjectDB(t, conn)
			err := w.admit(t, mode, repeated)
			var corrupt *CorruptJournalError
			if !errors.As(err, &corrupt) || corrupt.Line != len(admissionPrefix())+1 || !strings.Contains(err.Error(), "repeats Event ID") {
				t.Fatalf("%s = %v, want the repeated ID refused", mode, err)
			}
			conn = openReplayDB(t, filepath.Join(config.LedgerDir(w.root), "fid.sqlite"))
			defer closeProjectDB(t, conn)
			if !reflect.DeepEqual(before, currentProjectionSnapshot(t, conn)) {
				t.Fatalf("%s refusal changed the projection", mode)
			}
		})
	}
}
