package export

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/mbell-nzl/fid/internal/config"
	"github.com/mbell-nzl/fid/internal/fingerprint"
	"github.com/mbell-nzl/fid/internal/gittest"
	"github.com/mbell-nzl/fid/internal/holdwork"
	"github.com/mbell-nzl/fid/internal/journal"
	"github.com/mbell-nzl/fid/internal/reviewwork"
	"github.com/mbell-nzl/fid/internal/waiverrecord"
	"github.com/mbell-nzl/fid/internal/waiverwork"
)

func addWaiverExportFacts(t *testing.T, root string) string {
	t.Helper()
	task := addReviewExportFacts(t, root)
	taken, err := holdwork.Takeover(holdwork.Options{WorkspaceRoot: root, MigrationsFS: testMigrationsFS(t), TaskID: task, Actor: "human:mark", Reason: "Continue the fix"})
	if err != nil {
		t.Fatal(err)
	}
	_, err = holdwork.Handoff(holdwork.Options{WorkspaceRoot: root, MigrationsFS: testMigrationsFS(t), HoldID: taken.Fields["hold_id"].(string), Actor: "human:mark", PR: "https://github.com/owner/repo/pull/63", IdempotencyKey: "waiver-handoff"})
	if err != nil {
		t.Fatal(err)
	}
	_, err = waiverwork.Record(waiverwork.Options{WorkspaceRoot: root, MigrationsFS: testMigrationsFS(t), TaskID: task, Actor: "human:mark", Reason: "Reviewed all changes", Lane: "routine", IdempotencyKey: "waiver", Confirm: func(waiverwork.Preview) (bool, error) { return true, nil }})
	if err != nil {
		t.Fatal(err)
	}
	return task
}

func waiverExportSnapshot(t *testing.T, root string) []waiverrecord.Record {
	t.Helper()
	conn, err := reviewwork.OpenProjection(root)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = conn.Close() }()
	var task string
	if err := conn.QueryRow(`SELECT subject_id FROM person_decisions WHERE kind='waive'`).Scan(&task); err != nil {
		t.Fatal(err)
	}
	records, err := waiverrecord.Read(conn, task)
	if err != nil {
		t.Fatal(err)
	}
	return records
}

func TestWaiverExportImportRebuildAndRefusals(t *testing.T) {
	t.Parallel()
	f := newExportSourceFixture(t, addWaiverExportFacts)
	before := waiverExportSnapshot(t, f.workspace)
	if len(before) != 1 {
		t.Fatal(before)
	}
	if _, err := journal.RebuildWorkspace(f.workspace, testMigrationsFS(t)); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(before, waiverExportSnapshot(t, f.workspace)) {
		t.Fatal("rebuild changed waivers")
	}
	target := gittest.NewWorkspace(t)
	if _, err := Import(ImportOptions{ArchivePath: f.archive, RepositoryRoot: target, MigrationsFS: testMigrationsFS(t), IdempotencyKey: "waiver-import"}); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(before, waiverExportSnapshot(t, target)) {
		t.Fatal("import changed waivers")
	}
	assertImportedJournalEquals(t, f.workspace, target)
	for _, field := range []string{"agent", "binding", "retry"} {
		t.Run(field, func(t *testing.T) {
			t.Parallel()
			entries := archiveContents(t, f.archive)
			path := archivePathWithPrefix(t, entries, "journal/")
			data, hash := mutateArchiveEvent(t, entries[path], "task.waived", func(e *journal.Event) {
				if field == "retry" {
					e.Effects[1].Values["result_json"] = "{}"
					return
				}
				if field == "agent" {
					e.ActorV2.ActorType = "agent"
					e.Effects[0].Values["actor"] = "agent:mark"
				} else {
					var b waiverrecord.Binding
					if err := json.Unmarshal([]byte(e.Effects[0].Values["binding_json"].(string)), &b); err != nil {
						t.Fatal(err)
					}
					b.Head = strings.Repeat("e", 40)
					e.Effects[0].Values["binding_json"] = encodedLean(b)
				}
			})
			entries[path] = data
			manifest := decodeManifestMap(t, entries["manifest.json"])
			manifest["terminal_event_hash"] = hash
			updateManifestEntry(t, manifest, path, data)
			entries["manifest.json"] = encodeManifestMap(t, manifest)
			archive := filepath.Join(t.TempDir(), "forged.zip")
			writeTestArchive(t, archive, entries)
			refused := gittest.NewWorkspace(t)
			state := workspaceTreeDigest(t, refused)
			_, err := Import(ImportOptions{ArchivePath: archive, RepositoryRoot: refused, MigrationsFS: testMigrationsFS(t), IdempotencyKey: "refused"})
			assertArchiveCheck(t, err, "reference-invalid")
			if !reflect.DeepEqual(state, workspaceTreeDigest(t, refused)) {
				t.Fatal("refused import changed target")
			}
			if _, err := os.Stat(config.LedgerDir(refused)); !os.IsNotExist(err) {
				t.Fatalf("refused import published ledger: %v", err)
			}
		})
	}
}

func TestWaiverCurrencyMatchesOnlyAlgorithmAndFingerprint(t *testing.T) {
	t.Parallel()
	root := newExportSourceFixture(t, addWaiverExportFacts).workspace
	records := waiverExportSnapshot(t, root)
	b := records[0].Binding
	conn, err := reviewwork.OpenProjection(root)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = conn.Close() }()
	for _, tc := range []struct{ name, algorithm, fp, want string }{
		{"same change", b.Algorithm, b.Fingerprint, "current"},
		{"another algorithm", "other-algorithm", b.Fingerprint, "stale"},
		{"edited change", b.Algorithm, strings.Repeat("e", 64), "stale"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := waiverrecord.Current(conn, records[0].TaskID, fingerprint.Result{Algorithm: tc.algorithm, Fingerprint: tc.fp, Base: strings.Repeat("e", 40), Head: strings.Repeat("f", 40)})
			if err != nil || got.State != tc.want || got.Decision.ID != records[0].ID {
				t.Fatalf("currency=%v/%v", got, err)
			}
		})
	}
}
