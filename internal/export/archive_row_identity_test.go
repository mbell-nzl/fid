package export

import (
	"bytes"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mbell-nzl/fid/internal/gittest"
	"github.com/mbell-nzl/fid/internal/issues"
	"github.com/mbell-nzl/fid/internal/journal"
	"github.com/mbell-nzl/fid/internal/labels"
	"github.com/mbell-nzl/fid/internal/links"
)

// The same corpus reaches shared admission through import's private projection.
func TestImportRowIdentity(t *testing.T) {
	t.Parallel()
	raw, err := os.ReadFile("../../testdata/journal/row-identity.json")
	if err != nil {
		t.Fatal(err)
	}
	type probe struct {
		Name   string
		Prefix []journal.Event
		Event  journal.Event
	}
	var corpus struct {
		Prefix              []journal.Event
		Controls, Malformed []probe
	}
	if err := json.Unmarshal(raw, &corpus); err != nil {
		t.Fatal(err)
	}
	for _, group := range []struct {
		name    string
		probes  []probe
		refused bool
	}{{"control", corpus.Controls, false}, {"malformed", corpus.Malformed, true}} {
		for _, p := range group.probes {
			if !p.Event.IsV2() {
				continue
			} // Import already refuses all Event-v1 history.
			t.Run(group.name+"/"+p.Name, func(t *testing.T) {
				t.Parallel()
				v, err := newSemanticValidator(testMigrationsDir(t), 35)
				if err != nil {
					t.Fatal(err)
				}
				defer func() {
					if err := v.close(); err != nil {
						t.Error(err)
					}
				}()
				for _, e := range append(append([]journal.Event{}, corpus.Prefix...), p.Prefix...) {
					if err := applyLeanEmptyEvent(v, "import", e); err != nil {
						t.Fatalf("prefix: %v", err)
					}
				}
				before := handoffProjectionSnapshot(t, v)
				err = applyLeanEmptyEvent(v, "import", p.Event)
				if !group.refused {
					if err != nil {
						t.Fatalf("control: %v", err)
					}
					return
				}
				var semantic *semanticValidationError
				if !errors.As(err, &semantic) {
					t.Fatalf("malformed import = %v, want typed semantic failure", err)
				}
				if before != handoffProjectionSnapshot(t, v) {
					t.Fatal("refusal changed projection")
				}
				var valid journal.Event
				encoded, marshalErr := journal.MarshalCanonicalJSON(corpus.Controls[0].Event)
				if marshalErr != nil {
					t.Fatal(marshalErr)
				}
				if err := valid.UnmarshalJSON(encoded); err != nil {
					t.Fatal(err)
				}
				valid.EventID = p.Event.EventID
				valid.Effects[2].Values["journal_event_id"] = valid.EventID
				valid.Effects[2].PK["id"] = "ACT-later-valid"
				if len(p.Prefix) > 0 {
					valid.Effects[1].Values["revision"] = 3
					valid.Effects[1].Preconditions["revision"] = 2
					valid.Effects[1].Prior["revision"] = 2
				}
				if err := applyLeanEmptyEvent(v, "import", valid); err != nil {
					t.Fatalf("valid write after refusal: %v", err)
				}
			})
		}
	}
}

func TestImportFinishedCancellation(t *testing.T) {
	t.Parallel()
	for _, kind := range []string{"task", "epic"} {
		t.Run(kind, func(t *testing.T) {
			t.Parallel()
			v := stateMoveArchiveFixture(t)
			id := "FID-epic"
			if kind == "task" {
				id = "FID-todo"
			}
			// This validator has no workspace for commands. The CLI journey creates
			// both finished subjects through fid done with their required receipts.
			if _, err := v.conn.Exec(`UPDATE issues SET status='done' WHERE id=?`, id); err != nil {
				t.Fatal(err)
			}
			before := handoffProjectionSnapshot(t, v)
			event := stateMoveCancel(id, kind, "done")
			err := applyLeanEmptyEvent(v, "import", event)
			assertImportStateMoveRefused(t, err, id, "done", "cancelled")
			if !strings.Contains(err.Error(), "finished work stays done") {
				t.Fatalf("refusal: %v", err)
			}
			if before != handoffProjectionSnapshot(t, v) {
				t.Fatal("refusal changed projection")
			}
			if err := applyLeanEmptyEvent(v, "import", stateMoveCancel("FID-started", "task", "in_progress")); err != nil {
				t.Fatalf("valid cancellation after refusal: %v", err)
			}
		})
	}
}

func TestImportMalformedRowsLeavesTargetUnchanged(t *testing.T) {
	t.Parallel()
	var first, second string
	fixture := newExportSourceFixture(t, func(t *testing.T, root string) string {
		t.Helper()
		for _, label := range []string{"red", "blue"} {
			if _, err := labels.Add(labels.AddOptions{WorkspaceRoot: root, MigrationsDir: testMigrationsDir(t), Name: label, Actor: "human:mark", IdempotencyKey: label}); err != nil {
				t.Fatal(err)
			}
		}
		for _, title := range []string{"First", "Second"} {
			created, err := issues.Create(issues.CreateOptions{WorkspaceRoot: root, MigrationsDir: testMigrationsDir(t), Kind: "task", Title: title, Actor: "human:mark", IdempotencyKey: title})
			if err != nil {
				t.Fatal(err)
			}
			if first == "" {
				first = created.ID
			} else {
				second = created.ID
			}
		}
		if _, err := issues.AddAcceptanceCriterion(issues.AddAcceptanceCriterionOptions{WorkspaceRoot: root, MigrationsDir: testMigrationsDir(t), IssueID: first, Text: "Works", Ordinal: 1, Actor: "human:mark", IdempotencyKey: "ac"}); err != nil {
			t.Fatal(err)
		}
		if _, err := issues.AddLabel(issues.AddLabelOptions{WorkspaceRoot: root, MigrationsDir: testMigrationsDir(t), IssueID: first, Label: "red", Actor: "human:mark", IdempotencyKey: "add"}); err != nil {
			t.Fatal(err)
		}
		if _, err := links.Add(links.AddOptions{WorkspaceRoot: root, MigrationsDir: testMigrationsDir(t), FromID: first, ToID: second, Relation: "blocks", Actor: "human:mark", IdempotencyKey: "link"}); err != nil {
			t.Fatal(err)
		}
		return first
	})
	original := archiveContents(t, fixture.archive)
	for _, tc := range []struct {
		name, family string
		mutate       func(*journal.Event)
	}{
		{"rename", "issue.created", func(e *journal.Event) { e.Effects[0].Values["id"] = "FID-other" }},
		{"criterion subject", "issue.ac-added", func(e *journal.Event) { e.Effects[0].Values["issue_id"] = second }},
		{"issue label subject", "issue.label.added", func(e *journal.Event) { e.Effects[0].PK["issue_id"] = second }},
		{"label definition", "label.added", func(e *journal.Event) { e.Effects[0].PK["name"] = "green" }},
		{"link arguments", "link.added", func(e *journal.Event) { e.Aggregate.ID = "unrelated" }},
		{"activity subject", "issue.ac-added", func(e *journal.Event) { e.Effects[2].Values["ref_id"] = second }},
		{"repeated stage", "issue.ac-added", func(e *journal.Event) {
			extra := e.Effects[0]
			extra.PK = map[string]any{"id": "AC-extra"}
			e.Effects = append(e.Effects[:1:1], append([]journal.Effect{extra}, e.Effects[1:]...)...)
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			entries := map[string][]byte{}
			for path, data := range original {
				entries[path] = bytes.Clone(data)
			}
			path := archivePathWithPrefix(t, entries, "journal/")
			data, hash := mutateArchiveEvent(t, entries[path], tc.family, tc.mutate)
			entries[path] = data
			manifest := decodeManifestMap(t, entries["manifest.json"])
			manifest["terminal_event_hash"] = hash
			updateManifestEntry(t, manifest, path, data)
			entries["manifest.json"] = encodeManifestMap(t, manifest)
			archive := filepath.Join(t.TempDir(), "malformed.zip")
			writeTestArchive(t, archive, entries)
			target := gittest.NewWorkspace(t)
			before := workspaceTreeDigest(t, target)
			_, err := Import(ImportOptions{ArchivePath: archive, RepositoryRoot: target, MigrationsDir: testMigrationsDir(t), IdempotencyKey: "malformed"})
			var refused *InvalidArchiveError
			if !errors.As(err, &refused) {
				t.Fatalf("import = %v, want typed refusal", err)
			}
			if before != workspaceTreeDigest(t, target) {
				t.Fatal("refused import changed target")
			}
			if _, err := Import(ImportOptions{ArchivePath: fixture.archive, RepositoryRoot: target, MigrationsDir: testMigrationsDir(t), IdempotencyKey: "valid"}); err != nil {
				t.Fatal(err)
			}
			title := "After refusal"
			if _, err := issues.Update(issues.UpdateOptions{WorkspaceRoot: target, MigrationsDir: testMigrationsDir(t), ID: first, Title: &title, Actor: "human:mark", IdempotencyKey: "later"}); err != nil {
				t.Fatalf("later valid write: %v", err)
			}
		})
	}
}
