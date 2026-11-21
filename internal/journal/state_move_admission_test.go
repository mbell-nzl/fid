package journal

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/mbell-nzl/fid/internal/config"
)

// The state probes add an epic, FID-e, and a subtask, FID-s under FID-b, to
// the admission prefix. Each probe forges one Event after them, as
// EVT-20261110-000006.
const stateMoveEventID = "EVT-20261110-000006"

func stateMovePrefix() []Event {
	epic := admissionEvent("EVT-20261110-000004", "issue.created", "issue", "FID-e", 1,
		Effect{Op: "insert", Table: "issues", PK: map[string]any{"id": "FID-e"}, Values: map[string]any{"kind": "epic", "title": "Epic FID-e", "workflow_id": "epic-default", "status": "todo", "created_at": admissionStamp(0), "updated_at": admissionStamp(0)}},
		admissionActivity("EVT-20261110-000004", "ACT-FID-e", "FID-e", "created"))
	subtask := admissionEvent("EVT-20261110-000005", "issue.created", "issue", "FID-s", 1,
		Effect{Op: "insert", Table: "issues", PK: map[string]any{"id": "FID-s"}, Values: map[string]any{"kind": "subtask", "parent_id": "FID-b", "title": "Subtask FID-s", "workflow_id": "subtask-default", "status": "todo", "created_at": admissionStamp(0), "updated_at": admissionStamp(0)}},
		admissionActivity("EVT-20261110-000005", "ACT-FID-s", "FID-s", "created"))
	for _, event := range []*Event{&epic, &subtask} {
		event.CommandName = "issue.create"
	}
	return append(admissionPrefix(), epic, subtask)
}

var (
	stateAgent  = EventActor{ActorID: "builder", ActorType: "agent", Producer: "fid-standalone", ProvenanceRef: "project-local-attribution"}
	statePerson = EventActor{ActorID: "mark", ActorType: "human", Producer: "fid-standalone", ProvenanceRef: "project-local-attribution"}
)

var stateWorkflows = map[string]string{"FID-a": "task-default", "FID-e": "epic-default", "FID-s": "subtask-default"}

// stateTransition is issue.transitioned as fid issue transition writes it:
// the state, then the revision, then the activity row.
func stateTransition(issue, to string, actor EventActor) Event {
	event := admissionEvent(stateMoveEventID, "issue.transitioned", "issue", issue, 2,
		Effect{Op: "update", Table: "issues", PK: map[string]any{"id": issue}, Values: map[string]any{"status": to, "updated_at": admissionStamp(60)}, Prior: map[string]any{"status": "todo"}, Preconditions: map[string]any{"status": "todo", "workflow_id": stateWorkflows[issue]}},
		Effect{Op: "update", Table: "issues", PK: map[string]any{"id": issue}, Values: map[string]any{"revision": 2}, Prior: map[string]any{"revision": 1}, Preconditions: map[string]any{"revision": 1}},
		admissionActivity(stateMoveEventID, "ACT-move", issue, "status-changed"))
	event.CommandName, event.ActorV2 = "issue.transition", actor
	return event
}

// stateCancel is a person's cancellation, with the decision fid issue
// transition records.
func stateCancel(issue, kind string) Event {
	event := stateTransition(issue, "cancelled", statePerson)
	return withEffect(event, Effect{Op: "insert", Table: "person_decisions", PK: map[string]any{"id": "PDN-a1b2c3"}, Values: map[string]any{
		"kind": "cancel", "actor": "human:mark", "subject_type": kind, "subject_id": issue, "reason": "Not needed", "binding_id": issue, "decided_at": admissionStamp(60)}})
}

// stateUpdate is issue.updated as fid issue update writes it, with status
// among its values.
func stateUpdate(issue, to string) Event {
	event := admissionEvent(stateMoveEventID, "issue.updated", "issue", issue, 2,
		Effect{Op: "update", Table: "issues", PK: map[string]any{"id": issue}, Values: map[string]any{"status": to, "updated_at": admissionStamp(60)}, Prior: map[string]any{"status": "todo"}, Preconditions: map[string]any{"status": "todo"}},
		Effect{Op: "update", Table: "issues", PK: map[string]any{"id": issue}, Values: map[string]any{"revision": 2}, Prior: map[string]any{"revision": 1}, Preconditions: map[string]any{"revision": 1}},
		admissionActivity(stateMoveEventID, "ACT-update", issue, "updated"))
	event.CommandName, event.ActorV2 = "issue.update", stateAgent
	return event
}

func stateCreated(issue, kind, workflow, status string) Event {
	values := map[string]any{"kind": kind, "title": "Created " + issue, "workflow_id": workflow, "status": status, "created_at": admissionStamp(60), "updated_at": admissionStamp(60)}
	if kind == "task" {
		values["lane"] = "standard"
	}
	event := admissionEvent(stateMoveEventID, "issue.created", "issue", issue, 1,
		Effect{Op: "insert", Table: "issues", PK: map[string]any{"id": issue}, Values: values},
		admissionActivity(stateMoveEventID, "ACT-"+issue, issue, "created"))
	event.CommandName, event.ActorV2 = "issue.create", stateAgent
	return event
}

func stateLegacy(event Event) Event {
	event.EventSchema = EventSchemaV1
	event.Actor = Actor{Role: "developer"}
	event.Command = Command{Argv: []string{"fid", "issue", "transition"}, CWD: ".", FIDVersion: "0.1.0"}
	return event
}

// stateLaneRaise is the valid write each probe makes after a refusal.
func stateLaneRaise() Event {
	event := admissionLaneRaise()
	event.EventID = stateMoveEventID
	event.Effects[2].Values["journal_event_id"] = stateMoveEventID
	return event
}

// Admission takes a task's or an epic's state only from the Event of the
// command that owns the move: fid start, the hand-off, a review round
// asking for changes, fid done, or a person's cancellation. A forged
// issue.transitioned, an issue.updated, another Issue family or Event-v1
// history can't move one into any of those states on any journal path.
func TestForgedStateMovesRefusedOnEveryJournalPath(t *testing.T) {
	t.Parallel()
	type probe struct {
		issue, from, to string
		forge           func() Event
	}
	cases := map[string]probe{}
	for _, to := range []string{"in_progress", "in_review", "done", "cancelled"} {
		cases["task transition to "+to] = probe{"FID-a", "todo", to, func() Event { return stateTransition("FID-a", to, stateAgent) }}
		cases["task update to "+to] = probe{"FID-a", "todo", to, func() Event { return stateUpdate("FID-a", to) }}
		cases["legacy task transition to "+to] = probe{"FID-a", "todo", to, func() Event { return stateLegacy(stateTransition("FID-a", to, stateAgent)) }}
	}
	for _, to := range []string{"in_progress", "done", "cancelled"} {
		cases["epic transition to "+to] = probe{"FID-e", "todo", to, func() Event { return stateTransition("FID-e", to, stateAgent) }}
		cases["epic update to "+to] = probe{"FID-e", "todo", to, func() Event { return stateUpdate("FID-e", to) }}
	}
	cases["task cancelled by a person without the decision"] = probe{"FID-a", "todo", "cancelled", func() Event { return stateTransition("FID-a", "cancelled", statePerson) }}
	cases["epic cancelled by a person without the decision"] = probe{"FID-e", "todo", "cancelled", func() Event { return stateTransition("FID-e", "cancelled", statePerson) }}
	cases["legacy Event in the start family"] = probe{"FID-a", "todo", "in_progress", func() Event {
		event := stateLegacy(stateTransition("FID-a", "in_progress", stateAgent))
		event.EventType, event.CommandName = "task.started", "task.start"
		return event
	}}
	cases["legacy update to in_review"] = probe{"FID-a", "todo", "in_review", func() Event { return stateLegacy(stateUpdate("FID-a", "in_review")) }}
	cases["subtask update to done"] = probe{"FID-s", "todo", "done", func() Event { return stateUpdate("FID-s", "done") }}
	cases["task created done"] = probe{"FID-c", "", "done", func() Event { return stateCreated("FID-c", "task", "task-default", "done") }}
	cases["epic created in progress"] = probe{"FID-f", "", "in_progress", func() Event { return stateCreated("FID-f", "epic", "epic-default", "in_progress") }}
	cases["legacy task created in review"] = probe{"FID-c", "", "in_review", func() Event {
		return stateLegacy(stateCreated("FID-c", "task", "task-default", "in_review"))
	}}
	cases["acceptance criterion moves its task"] = probe{"FID-a", "todo", "done", func() Event {
		event := admissionEvent(stateMoveEventID, "issue.ac-added", "issue", "FID-a", 2,
			Effect{Op: "insert", Table: "acceptance_criteria", PK: map[string]any{"id": "AC-1"}, Values: map[string]any{"issue_id": "FID-a", "ordinal": 1, "text": "Works"}},
			Effect{Op: "update", Table: "issues", PK: map[string]any{"id": "FID-a"}, Values: map[string]any{"revision": 2, "status": "done"}, Prior: map[string]any{"revision": 1}, Preconditions: map[string]any{"revision": 1}},
			admissionActivity(stateMoveEventID, "ACT-ac", "FID-a", "ac-added"))
		event.CommandName, event.ActorV2 = "issue.ac.add", stateAgent
		return event
	}}
	cases["link moves its end"] = probe{"FID-b", "todo", "done", func() Event {
		event := admissionLinkAdded()
		event.EventID = stateMoveEventID
		event.Effects[3].Values["status"] = "done"
		event.Effects[4].Values["journal_event_id"] = stateMoveEventID
		return event
	}}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			assertStateMoveRefused(t, tc.forge(), tc.issue, tc.from, tc.to)
		})
	}
}

// The moves fid's own commands make through issue.transitioned still pass:
// a person's cancellation and a subtask's plain transition.
func TestOwnedStateMovesAcceptedOnEveryJournalPath(t *testing.T) {
	t.Parallel()
	for name, event := range map[string]Event{
		"person cancels a task":   stateCancel("FID-a", "task"),
		"person cancels an epic":  stateCancel("FID-e", "epic"),
		"subtask moves":           stateTransition("FID-s", "in_progress", stateAgent),
		"legacy subtask moves":    stateLegacy(stateTransition("FID-s", "done", stateAgent)),
		"task created to do":      stateCreated("FID-c", "task", "task-default", "todo"),
		"legacy epic created":     stateLegacy(stateCreated("FID-f", "epic", "epic-default", "todo")),
		"subtask created started": stateCreated("FID-t", "subtask", "subtask-default", "in_progress"),
	} {
		if event.Effects[0].Values["kind"] == "subtask" {
			event.Effects[0].Values["parent_id"] = "FID-b"
			delete(event.Effects[0].Values, "lane")
		}
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			for _, mode := range admissionModes(event) {
				t.Run(mode, func(t *testing.T) {
					w := newAdmissionWorkspace(t, stateMovePrefix(), mode)
					if err := w.admit(t, mode, event); err != nil {
						t.Fatalf("%s refused a move fid's commands make: %v", mode, err)
					}
				})
			}
		})
	}
}

// assertStateMoveRefused offers forged on every journal path: each refuses it
// with a StateMoveError naming the move, and changes nothing, and the ledger
// still takes a valid write.
func assertStateMoveRefused(t *testing.T, forged Event, issue, from, to string) {
	t.Helper()
	for _, mode := range admissionModes(forged) {
		t.Run(mode, func(t *testing.T) {
			w := newAdmissionWorkspace(t, stateMovePrefix(), mode)
			// A full replay starts from an empty projection and projects the
			// prefix before it reaches the forged Event; the prefix alone,
			// projected, is what it must leave.
			reference := w
			if mode == "replay" {
				reference = newAdmissionWorkspace(t, stateMovePrefix(), "project")
			}
			before := stateProjection(t, reference.root)
			err := w.admit(t, mode, forged)
			var refused *StateMoveError
			if !errors.As(err, &refused) || refused.IssueID != issue || refused.From != from || refused.To != to || refused.Family != forged.EventType {
				t.Fatalf("%s = %v (%T), want a StateMoveError for %s %s -> %s by %s", mode, err, err, issue, from, to, forged.EventType)
			}
			if !reflect.DeepEqual(before, stateProjection(t, w.root)) {
				t.Fatalf("%s refusal changed the projection", mode)
			}
			stored, readErr := os.ReadFile(w.file)
			if readErr != nil {
				t.Fatal(readErr)
			}
			if mode == "commit" || mode == "dry-run" || mode == "project" {
				if !reflect.DeepEqual(stored, w.prefix) {
					t.Fatalf("%s refusal changed the journal", mode)
				}
			} else if err := os.WriteFile(w.file, w.prefix, 0o644); err != nil {
				t.Fatal(err)
			}
			if err := w.admit(t, "commit", stateLaneRaise()); err != nil {
				t.Fatalf("valid write after %s refusal: %v", mode, err)
			}
		})
	}
}

// stateProjection snapshots every projected table. The checkpoint row keeps
// its position but not its clock, so two workspaces that projected the same
// lines compare equal.
func stateProjection(t *testing.T, root string) map[string][]string {
	t.Helper()
	conn := openReplayDB(t, filepath.Join(config.LedgerDir(root), "fid.sqlite"))
	defer closeProjectDB(t, conn)
	snapshot := currentProjectionSnapshot(t, conn)
	var file, event, hash, status sql.NullString
	var offset int64
	err := conn.QueryRow(`SELECT journal_file, journal_offset, last_event_id, last_event_hash, replay_status FROM projection_meta`).Scan(&file, &offset, &event, &hash, &status)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		t.Fatal(err)
	}
	snapshot["projection_meta"] = []string{fmt.Sprint(file, offset, event, hash, status)}
	return snapshot
}
