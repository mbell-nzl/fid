package export

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/mbell-nzl/fid/internal/gittest"
	"github.com/mbell-nzl/fid/internal/issues"
	"github.com/mbell-nzl/fid/internal/journal"
)

var stateMoveAt = time.Date(2026, 11, 15, 2, 0, 0, 0, time.UTC)

// stateMoveArchiveFixture is the epic-close fixture's two epics in todo plus
// three tasks, one in each state a command moves a task out of.
func stateMoveArchiveFixture(t *testing.T) *semanticValidator {
	t.Helper()
	v, _ := epicCloseArchiveFixture(t)
	for id, state := range map[string]string{"FID-todo": "todo", "FID-started": "in_progress", "FID-review": "in_review"} {
		if _, err := v.conn.Exec(`INSERT INTO issues(id,kind,title,workflow_id,status,lane,owner_role,created_at,updated_at) VALUES(?,'task','Task','task-default',?,'standard','developer','2026-11-15T00:00:00Z','2026-11-15T00:00:00Z')`, id, state); err != nil {
			t.Fatal(err)
		}
	}
	return v
}

var stateMoveWorkflows = map[string]string{"FID-todo": "task-default", "FID-started": "task-default", "FID-review": "task-default", "FID-epic": "epic-default", "FID-sub": "subtask-default"}

// stateMoveEvent writes to as fid issue transition or fid issue update
// would: the state, the revision, then the activity row.
func stateMoveEvent(family, issue, from, to, actor string) journal.Event {
	stamp := stateMoveAt.Format(time.RFC3339)
	preconditions := map[string]any{"status": from}
	if family == "issue.transitioned" {
		preconditions["workflow_id"] = stateMoveWorkflows[issue]
	}
	kind, id, _ := strings.Cut(actor, ":")
	return journal.Event{EventSchema: 2, EventType: family, EventID: "EVT-20261115-000009", CommandName: map[string]string{"issue.transitioned": "issue.transition", "issue.updated": "issue.update"}[family], OccurredAt: stateMoveAt,
		ActorV2:   journal.EventActor{ActorType: kind, ActorID: id, Producer: "fid-standalone", ProvenanceRef: "project-local-attribution"},
		Aggregate: journal.Aggregate{Type: "issue", ID: issue, Sequence: 2}, Effects: []journal.Effect{
			{Op: "update", Table: "issues", PK: map[string]any{"id": issue}, Values: map[string]any{"status": to, "updated_at": stamp}, Prior: map[string]any{"status": from}, Preconditions: preconditions},
			{Op: "update", Table: "issues", PK: map[string]any{"id": issue}, Values: map[string]any{"revision": 2}, Prior: map[string]any{"revision": 1}, Preconditions: map[string]any{"revision": 1}},
			{Op: "activity", Table: "activity_log", PK: map[string]any{"id": "ACT-move"}, Values: map[string]any{"journal_event_id": "EVT-20261115-000009", "ref_type": "issue", "ref_id": issue, "event": "status-changed", "details": nil, "actor": actor, "occurred_at": stamp}}}}
}

// stateMoveCancel is a person's cancellation with its decision.
func stateMoveCancel(issue, kind, from string) journal.Event {
	event := stateMoveEvent("issue.transitioned", issue, from, "cancelled", "human:mark")
	decision := journal.Effect{Op: "insert", Table: "person_decisions", PK: map[string]any{"id": "PDN-a1b2c3"}, Values: map[string]any{
		"kind": "cancel", "actor": "human:mark", "subject_type": kind, "subject_id": issue, "reason": "Not needed", "binding_id": issue, "decided_at": stateMoveAt.Format(time.RFC3339)}}
	event.Effects = append(event.Effects[:2:2], decision, event.Effects[2])
	return event
}

// Import admits a task's or an epic's state only from the Event of the
// command that owns the move, as live append and replay do: a forged
// transition or update is refused around each of those commands, and the
// validator still takes a valid Event afterwards.
func TestImportRefusesForgedStateMoves(t *testing.T) {
	t.Parallel()
	type probe struct{ issue, from, to, actor string }
	cases := map[string]probe{
		"around fid done":                    {"FID-review", "in_review", "done", "agent:builder"},
		"around the hand-off":                {"FID-started", "in_progress", "in_review", "agent:builder"},
		"around a review round":              {"FID-review", "in_review", "in_progress", "agent:builder"},
		"around fid start":                   {"FID-todo", "todo", "in_progress", "agent:builder"},
		"around a person's cancel":           {"FID-started", "in_progress", "cancelled", "agent:builder"},
		"cancel without its decision":        {"FID-started", "in_progress", "cancelled", "human:mark"},
		"epic around fid start":              {"FID-epic", "todo", "in_progress", "agent:builder"},
		"epic around fid done":               {"FID-epic", "todo", "done", "human:mark"},
		"epic around a person's cancel":      {"FID-epic", "todo", "cancelled", "agent:builder"},
		"task moved back to its first state": {"FID-started", "in_progress", "todo", "human:mark"},
	}
	for name, tc := range cases {
		for _, family := range []string{"issue.transitioned", "issue.updated"} {
			t.Run(family+"/"+name, func(t *testing.T) {
				t.Parallel()
				v := stateMoveArchiveFixture(t)
				before := handoffProjectionSnapshot(t, v)
				err := applyLeanEmptyEvent(v, "import", stateMoveEvent(family, tc.issue, tc.from, tc.to, tc.actor))
				assertImportStateMoveRefused(t, err, tc.issue, tc.from, tc.to)
				if before != handoffProjectionSnapshot(t, v) {
					t.Fatal("refused import changed the projection")
				}
				if err := applyLeanEmptyEvent(v, "import", stateMoveCancel("FID-todo", "task", "todo")); err != nil {
					t.Fatalf("valid Event after the refusal = %v", err)
				}
			})
		}
	}
	// A cancellation takes exactly one person's decision, not two.
	for _, tc := range []struct{ issue, kind, from string }{{"FID-started", "task", "in_progress"}, {"FID-epic", "epic", "todo"}} {
		t.Run("cancel with two decisions/"+tc.kind, func(t *testing.T) {
			t.Parallel()
			v := stateMoveArchiveFixture(t)
			event := stateMoveCancel(tc.issue, tc.kind, tc.from)
			second := event.Effects[2]
			second.PK = map[string]any{"id": "PDN-d4e5f6"}
			event.Effects = append(event.Effects[:3:3], second, event.Effects[3])
			before := handoffProjectionSnapshot(t, v)
			assertImportStateMoveRefused(t, applyLeanEmptyEvent(v, "import", event), tc.issue, tc.from, "cancelled")
			if before != handoffProjectionSnapshot(t, v) {
				t.Fatal("refused import changed the projection")
			}
		})
	}
	t.Run("task created done", func(t *testing.T) {
		t.Parallel()
		v := stateMoveArchiveFixture(t)
		created := stateMoveEvent("issue.created", "FID-new", "", "done", "agent:builder")
		created.CommandName, created.Aggregate.Sequence = "issue.create", 1
		created.Effects = []journal.Effect{{Op: "insert", Table: "issues", PK: map[string]any{"id": "FID-new"}, Values: map[string]any{"kind": "task", "title": "New", "workflow_id": "task-default", "status": "done", "lane": "standard", "created_at": stateMoveAt.Format(time.RFC3339), "updated_at": stateMoveAt.Format(time.RFC3339)}}, created.Effects[2]}
		created.Effects[1].Values["ref_id"], created.Effects[1].Values["event"] = "FID-new", "created"
		before := handoffProjectionSnapshot(t, v)
		assertImportStateMoveRefused(t, applyLeanEmptyEvent(v, "import", created), "FID-new", "", "done")
		if before != handoffProjectionSnapshot(t, v) {
			t.Fatal("refused import changed the projection")
		}
	})
}

// Import keys a state change to one Issue once per Event, as live append and
// replay do. A subtask's moves are free, so nothing else refuses two in one
// Event; the same Event with one move then imports.
func TestImportRefusesTwoStateChangesInOneEvent(t *testing.T) {
	t.Parallel()
	v := stateMoveArchiveFixture(t)
	if _, err := v.conn.Exec(`INSERT INTO issues(id,kind,parent_id,title,workflow_id,status,owner_role,created_at,updated_at) VALUES('FID-sub','subtask','FID-todo','Subtask','subtask-default','todo','developer','2026-11-15T00:00:00Z','2026-11-15T00:00:00Z')`); err != nil {
		t.Fatal(err)
	}
	twice := stateMoveEvent("issue.transitioned", "FID-sub", "todo", "in_progress", "agent:builder")
	again := journal.Effect{Op: "update", Table: "issues", PK: map[string]any{"id": "FID-sub"}, Values: map[string]any{"status": "done", "updated_at": stateMoveAt.Format(time.RFC3339)}, Prior: map[string]any{"status": "in_progress"}, Preconditions: map[string]any{"status": "in_progress"}}
	twice.Effects = append([]journal.Effect{twice.Effects[0], again}, twice.Effects[1:]...)
	before := handoffProjectionSnapshot(t, v)
	err := applyLeanEmptyEvent(v, "import", twice)
	var semantic *semanticValidationError
	if reason := "a state change keys one Issue by id alone, once per Event"; !errors.As(err, &semantic) || semantic.Code != "reference-invalid" || !strings.Contains(err.Error(), reason) {
		t.Fatalf("import = %v (%T), want reference-invalid naming %q", err, err, reason)
	}
	if before != handoffProjectionSnapshot(t, v) {
		t.Fatal("refused import changed the projection")
	}
	if err := applyLeanEmptyEvent(v, "import", stateMoveEvent("issue.transitioned", "FID-sub", "todo", "in_progress", "agent:builder")); err != nil {
		t.Fatalf("one move of the subtask after the refusal = %v", err)
	}
}

// A person's cancellation still imports, for a task in each open state and
// for an epic.
func TestImportAdmitsPersonCancellation(t *testing.T) {
	t.Parallel()
	for issue, from := range map[string]string{"FID-todo": "todo", "FID-started": "in_progress", "FID-review": "in_review", "FID-epic": "todo"} {
		t.Run(issue, func(t *testing.T) {
			t.Parallel()
			v := stateMoveArchiveFixture(t)
			kind := "task"
			if issue == "FID-epic" {
				kind = "epic"
			}
			if err := applyLeanEmptyEvent(v, "import", stateMoveCancel(issue, kind, from)); err != nil {
				t.Fatalf("person's cancellation = %v", err)
			}
		})
	}
}

func assertImportStateMoveRefused(t *testing.T, err error, issue, from, to string) {
	t.Helper()
	var semantic *semanticValidationError
	var refused *journal.StateMoveError
	if !errors.As(err, &semantic) || semantic.Code != "reference-invalid" || !errors.As(err, &refused) || refused.IssueID != issue || refused.From != from || refused.To != to {
		t.Fatalf("import = %v (%T), want reference-invalid with a StateMoveError for %s %s -> %s", err, err, issue, from, to)
	}
}

// A backup whose journal ends in a forged move is refused by validation and
// by import, and the target is left as it was; the unforged backup then
// restores and takes a write.
func TestImportRefusesArchiveWithForgedStateMove(t *testing.T) {
	t.Parallel()
	fixture := newExportSourceFixture(t, func(t *testing.T, root string) string {
		t.Helper()
		task, err := issues.Create(issues.CreateOptions{WorkspaceRoot: root, MigrationsDir: testMigrationsDir(t), Kind: "task", Title: "Before the backup", Actor: "human:mark", IdempotencyKey: "task"})
		if err != nil {
			t.Fatal(err)
		}
		return task.ID
	})
	original := archiveContents(t, fixture.archive)
	path := archivePathWithPrefix(t, original, "journal/")
	lines := bytes.Split(bytes.TrimSuffix(original[path], []byte{'\n'}), []byte{'\n'})
	var last journal.Event
	if err := last.UnmarshalJSON(lines[len(lines)-1]); err != nil {
		t.Fatal(err)
	}
	parts := strings.Split(last.EventID, "-")
	sequence, err := strconv.Atoi(parts[2])
	if err != nil {
		t.Fatal(err)
	}
	forged := stateMoveEvent("issue.transitioned", last.Aggregate.ID, "todo", "done", "agent:builder")
	forged.EventID = fmt.Sprintf("%s-%s-%06d", parts[0], parts[1], sequence+1)
	forged.OccurredAt = last.OccurredAt.Add(time.Second)
	forged.Effects[0].Values["updated_at"] = forged.OccurredAt.UTC().Format(time.RFC3339)
	forged.Effects[0].Preconditions["workflow_id"] = "task-default"
	forged.Effects[2].Values["journal_event_id"], forged.Effects[2].Values["occurred_at"] = forged.EventID, forged.OccurredAt.UTC().Format(time.RFC3339)
	encoded, err := journal.MarshalCanonicalJSON(forged)
	if err != nil {
		t.Fatal(err)
	}
	entries := make(map[string][]byte, len(original))
	for name, data := range original {
		entries[name] = bytes.Clone(data)
	}
	entries[path] = append(entries[path], append(encoded, '\n')...)
	manifest := decodeManifestMap(t, entries["manifest.json"])
	if manifest["terminal_event_hash"], err = journal.HashEvent(forged); err != nil {
		t.Fatal(err)
	}
	manifest["event_count"] = json.Number(fmt.Sprint(len(lines) + 1))
	updateManifestEntry(t, manifest, path, entries[path])
	entries["manifest.json"] = encodeManifestMap(t, manifest)
	archive := filepath.Join(t.TempDir(), "forged.zip")
	writeTestArchive(t, archive, entries)

	_, err = Validate(ValidateOptions{ArchivePath: archive, MigrationsDir: testMigrationsDir(t)})
	assertStateMoveCheck(t, err, last.Aggregate.ID)
	target := gittest.NewWorkspace(t)
	before := workspaceTreeDigest(t, target)
	_, err = Import(ImportOptions{ArchivePath: archive, RepositoryRoot: target, MigrationsDir: testMigrationsDir(t), IdempotencyKey: "forged"})
	assertStateMoveCheck(t, err, last.Aggregate.ID)
	if workspaceTreeDigest(t, target) != before {
		t.Fatal("refused import changed the target")
	}
	if _, err := Import(ImportOptions{ArchivePath: fixture.archive, RepositoryRoot: target, MigrationsDir: testMigrationsDir(t), IdempotencyKey: "unforged"}); err != nil {
		t.Fatalf("import of the unforged archive = %v", err)
	}
	if _, err := issues.Create(issues.CreateOptions{WorkspaceRoot: target, MigrationsDir: testMigrationsDir(t), Kind: "task", Title: "After the restore", Actor: "human:mark", IdempotencyKey: "after"}); err != nil {
		t.Fatalf("write after the restore = %v", err)
	}
}

func assertStateMoveCheck(t *testing.T, err error, issue string) {
	t.Helper()
	var invalid *InvalidArchiveError
	if errors.As(err, &invalid) {
		for _, check := range invalid.Checks {
			if check.Code == "reference-invalid" && strings.Contains(check.Detail, "can't move task "+issue+" from todo to done") {
				return
			}
		}
	}
	t.Fatalf("archive error = %#v, want reference-invalid naming the move of %s", err, issue)
}
