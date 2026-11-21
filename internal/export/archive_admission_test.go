package export

import (
	"bytes"
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
	"github.com/mbell-nzl/fid/internal/labels"
)

// fid issue label add and remove write issue.label.added and
// issue.label.removed; a backup holding them must restore.
func TestImportAdmitsEveryFamilyTheWritersProduce(t *testing.T) {
	t.Parallel()
	fixture := newExportSourceFixture(t, func(t *testing.T, root string) string {
		t.Helper()
		if _, err := labels.Add(labels.AddOptions{WorkspaceRoot: root, MigrationsDir: testMigrationsDir(t), Name: "backend", Actor: "human:mark", IdempotencyKey: "label"}); err != nil {
			t.Fatal(err)
		}
		task, err := issues.Create(issues.CreateOptions{WorkspaceRoot: root, MigrationsDir: testMigrationsDir(t), Kind: "task", Title: "Labelled task", Actor: "human:mark", IdempotencyKey: "task"})
		if err != nil {
			t.Fatal(err)
		}
		if _, err := issues.AddLabel(issues.AddLabelOptions{WorkspaceRoot: root, MigrationsDir: testMigrationsDir(t), IssueID: task.ID, Label: "backend", Actor: "human:mark", IdempotencyKey: "add"}); err != nil {
			t.Fatal(err)
		}
		if _, err := issues.RemoveLabel(issues.RemoveLabelOptions{WorkspaceRoot: root, MigrationsDir: testMigrationsDir(t), IssueID: task.ID, Label: "backend", Actor: "human:mark", IdempotencyKey: "remove"}); err != nil {
			t.Fatal(err)
		}
		return task.ID
	})
	if _, err := Validate(ValidateOptions{ArchivePath: fixture.archive, MigrationsDir: testMigrationsDir(t)}); err != nil {
		t.Fatalf("validate = %v", err)
	}
	target := gittest.NewWorkspace(t)
	if _, err := Import(ImportOptions{ArchivePath: fixture.archive, RepositoryRoot: target, MigrationsDir: testMigrationsDir(t), IdempotencyKey: "labels"}); err != nil {
		t.Fatalf("import = %v", err)
	}
	assertImportedJournalEquals(t, fixture.workspace, target)
}

// Import refuses what live append and replay refuse, with nothing changed.
func TestImportRefusesForgedEventsLikeTheJournal(t *testing.T) {
	t.Parallel()
	at := time.Date(2026, 11, 15, 2, 0, 0, 0, time.UTC)
	hold := func() journal.Event {
		return journal.Event{EventSchema: 2, EventType: "hold.taken", EventID: "EVT-20261115-000009", CommandName: "next.claim", OccurredAt: at,
			ActorV2:   journal.EventActor{ActorType: "agent", ActorID: "one", Producer: "fid-standalone", ProvenanceRef: "project-local-attribution"},
			Aggregate: journal.Aggregate{Type: "hold", ID: "HLD-forged", Sequence: 1}, Effects: []journal.Effect{{Op: "insert", Table: "holds", PK: map[string]any{"id": "HLD-forged"}, Values: map[string]any{
				"task_id": "FID-epic", "worktree_id": "sha256:one", "worktree_root": "/w/one", "actor": "agent:one", "since": at.Format(time.RFC3339), "expires_at": at.Add(time.Hour).Format(time.RFC3339), "replaces_hold_id": nil, "ended_at": nil, "end_reason": nil}}}}
	}
	retitle := func() journal.Event {
		stamp := at.Format(time.RFC3339)
		event := journal.Event{EventSchema: 2, EventType: "issue.updated", EventID: "EVT-20261115-000009", CommandName: "issue.update", OccurredAt: at,
			ActorV2:   journal.EventActor{ActorType: "human", ActorID: "mark", Producer: "fid-standalone", ProvenanceRef: "project-local-attribution"},
			Aggregate: journal.Aggregate{Type: "issue", ID: "FID-epic", Sequence: 1}, Effects: []journal.Effect{
				{Op: "update", Table: "issues", PK: map[string]any{"id": "FID-epic"}, Values: map[string]any{"title": "Renamed", "updated_at": stamp}},
				{Op: "update", Table: "issues", PK: map[string]any{"id": "FID-epic"}, Values: map[string]any{"revision": 2}, Preconditions: map[string]any{"revision": 1}},
				{Op: "activity", Table: "activity_log", PK: map[string]any{"id": "ACT-forged"}, Values: map[string]any{"journal_event_id": "EVT-20261115-000009", "ref_type": "issue", "ref_id": "FID-epic", "event": "updated", "details": nil, "actor": "human:mark", "occurred_at": stamp}}}}
		return writerUpdateFixture(t, event, 1)
	}
	unbound := func(e journal.Event) journal.Event {
		e.Effects = append(e.Effects, journal.Effect{Op: "insert", Table: "idempotency_results", PK: map[string]any{"operation_scope": "task.start:project", "idempotency_key": "forged"}, Values: map[string]any{
			"request_digest": "sha256:" + strings.Repeat("b", 64), "request_field_digests_json": "{}", "result_entity_id": e.Aggregate.ID, "result_event_id": e.EventID, "aggregate_sequence": 1, "result_json": `{"replayed":false}`, "created_at": at.Format(time.RFC3339)}})
		return e
	}
	for _, tc := range []struct {
		name, reason string
		event        func() journal.Event
		forge        func(journal.Event) journal.Event
	}{
		{"hold control", "", hold, func(e journal.Event) journal.Event { return e }},
		{"service actor", "actor not admitted", hold, func(e journal.Event) journal.Event { e.ActorV2.ActorType = "service"; return e }},
		{"oversized", "1 MiB line budget", hold, func(e journal.Event) journal.Event {
			e.Effects[0].Values["worktree_root"] = strings.Repeat("<", 200000)
			return e
		}},
		{"hold retry result without its request", "idempotency", hold, unbound},
		// The writers read every ID dated like theirs to choose the next.
		{"dated ID without the EVT prefix", "Event ID", hold, func(e journal.Event) journal.Event { e.EventID = "BAD-20261115-000009"; return e }},
		{"update control", "", retitle, func(e journal.Event) journal.Event { return e }},
		{"update retry result without its request", "idempotency", retitle, unbound},
		{"update keyed by another column", "primary key", retitle, func(e journal.Event) journal.Event {
			e.Effects[0].PK = map[string]any{"title": "Epic"}
			return e
		}},
		{"update without a key", "primary key", retitle, func(e journal.Event) journal.Event {
			e.Effects[1].PK = map[string]any{}
			return e
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			v, _ := epicCloseArchiveFixture(t)
			before := handoffProjectionSnapshot(t, v)
			err := applyLeanEmptyEvent(v, "import", tc.forge(tc.event()))
			if tc.reason == "" {
				if err != nil {
					t.Fatalf("control refused: %v", err)
				}
				return
			}
			if err == nil || !strings.Contains(err.Error(), tc.reason) {
				t.Fatalf("import = %v, want refusal naming %q", err, tc.reason)
			}
			if before != handoffProjectionSnapshot(t, v) {
				t.Fatal("refused import changed the projection")
			}
		})
	}
}

// The writers allocate each dated Event ID as one past its date's highest
// sequence. An archive whose last ID skips ahead, up to the ceiling, would
// leave the restored ledger's writers a sequence admission refuses, so
// import refuses it and changes nothing. The same target then restores from
// the unforged archive and takes a write with an allocated ID.
func TestImportRefusesEventIDsTheWritersWouldNotAllocate(t *testing.T) {
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
	for _, tc := range []struct {
		name     string
		sequence int
	}{{"at the ceiling", 999999999}, {"one below the ceiling", 999999998}, {"a skipped sequence", sequence + 1}} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			forgedID := fmt.Sprintf("%s-%s-%06d", parts[0], parts[1], tc.sequence)
			forgedLines := append(append([][]byte{}, lines[:len(lines)-1]...), bytes.ReplaceAll(lines[len(lines)-1], []byte(last.EventID), []byte(forgedID)))
			var forged journal.Event
			if err := forged.UnmarshalJSON(forgedLines[len(forgedLines)-1]); err != nil {
				t.Fatal(err)
			}
			hash, err := journal.HashEvent(forged)
			if err != nil {
				t.Fatal(err)
			}
			entries := make(map[string][]byte, len(original))
			for name, data := range original {
				entries[name] = bytes.Clone(data)
			}
			entries[path] = append(bytes.Join(forgedLines, []byte{'\n'}), '\n')
			manifest := decodeManifestMap(t, entries["manifest.json"])
			manifest["terminal_event_hash"] = hash
			updateManifestEntry(t, manifest, path, entries[path])
			entries["manifest.json"] = encodeManifestMap(t, manifest)
			archive := filepath.Join(t.TempDir(), "forged.zip")
			writeTestArchive(t, archive, entries)

			_, err = Validate(ValidateOptions{ArchivePath: archive, MigrationsDir: testMigrationsDir(t)})
			assertEventIDCheck(t, err, forgedID)
			target := gittest.NewWorkspace(t)
			before := workspaceTreeDigest(t, target)
			_, err = Import(ImportOptions{ArchivePath: archive, RepositoryRoot: target, MigrationsDir: testMigrationsDir(t), IdempotencyKey: "forged"})
			assertEventIDCheck(t, err, forgedID)
			if workspaceTreeDigest(t, target) != before {
				t.Fatal("refused import changed the target")
			}
			if _, err := Import(ImportOptions{ArchivePath: fixture.archive, RepositoryRoot: target, MigrationsDir: testMigrationsDir(t), IdempotencyKey: "unforged"}); err != nil {
				t.Fatalf("import of the unforged archive = %v", err)
			}
			if _, err := issues.Create(issues.CreateOptions{WorkspaceRoot: target, MigrationsDir: testMigrationsDir(t), Kind: "task", Title: "After the restore", Actor: "human:mark", IdempotencyKey: "after"}); err != nil {
				t.Fatalf("write with an allocated ID after the restore = %v", err)
			}
		})
	}
}

func assertEventIDCheck(t *testing.T, err error, id string) {
	t.Helper()
	var invalid *InvalidArchiveError
	if errors.As(err, &invalid) {
		for _, check := range invalid.Checks {
			if check.Code == "journal-invalid" && strings.Contains(check.Detail, "Event ID \""+id+"\"") {
				return
			}
		}
	}
	t.Fatalf("archive error = %#v, want journal-invalid naming Event ID %s", err, id)
}
