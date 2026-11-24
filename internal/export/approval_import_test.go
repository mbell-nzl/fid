package export

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/mbell-nzl/fid/internal/approvalrecord"
	"github.com/mbell-nzl/fid/internal/approvalwork"
	"github.com/mbell-nzl/fid/internal/config"
	"github.com/mbell-nzl/fid/internal/gittest"
	"github.com/mbell-nzl/fid/internal/holdwork"
	"github.com/mbell-nzl/fid/internal/journal"
	"github.com/mbell-nzl/fid/internal/reviewwork"
)

func addApprovalExportFacts(t *testing.T, root string) string {
	t.Helper()
	task := addReviewExportFacts(t, root)
	taken, err := holdwork.Takeover(holdwork.Options{WorkspaceRoot: root, MigrationsFS: testMigrationsFS(t), TaskID: task, Actor: "human:mark", Reason: "Continue the fix"})
	if err != nil {
		t.Fatal(err)
	}
	_, err = holdwork.Handoff(holdwork.Options{WorkspaceRoot: root, MigrationsFS: testMigrationsFS(t), HoldID: taken.Fields["hold_id"].(string), Actor: "human:mark", PR: "https://github.com/owner/repo/pull/63", IdempotencyKey: "approval-handoff"})
	if err != nil {
		t.Fatal(err)
	}
	_, err = approvalwork.Record(approvalwork.Options{WorkspaceRoot: root, MigrationsFS: testMigrationsFS(t), TaskID: task, Actor: "human:mark", Reason: "Reviewed all changes", IdempotencyKey: "approval", Confirm: func(approvalwork.Preview) (bool, error) { return true, nil }})
	if err != nil {
		t.Fatal(err)
	}
	return task
}

func approvalExportSnapshot(t *testing.T, root string) []approvalrecord.Record {
	t.Helper()
	conn, err := reviewwork.OpenProjection(root)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = conn.Close() }()
	var task string
	if err := conn.QueryRow(`SELECT subject_id FROM person_decisions WHERE kind='approve'`).Scan(&task); err != nil {
		t.Fatal(err)
	}
	records, err := approvalrecord.Read(conn, task)
	if err != nil {
		t.Fatal(err)
	}
	return records
}

func TestApprovalExportImportRebuildAndRefusals(t *testing.T) {
	t.Parallel()
	f := newExportSourceFixture(t, addApprovalExportFacts)
	before := approvalExportSnapshot(t, f.workspace)
	if len(before) != 1 {
		t.Fatal(before)
	}
	if _, err := journal.RebuildWorkspace(f.workspace, testMigrationsFS(t)); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(before, approvalExportSnapshot(t, f.workspace)) {
		t.Fatal("rebuild changed approvals")
	}
	target := gittest.NewWorkspace(t)
	if _, err := Import(ImportOptions{ArchivePath: f.archive, RepositoryRoot: target, MigrationsFS: testMigrationsFS(t), IdempotencyKey: "approval-import"}); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(before, approvalExportSnapshot(t, target)) {
		t.Fatal("import changed approvals")
	}
	assertImportedJournalEquals(t, f.workspace, target)
	for _, field := range []string{"agent", "binding", "retry"} {
		t.Run(field, func(t *testing.T) {
			t.Parallel()
			entries := archiveContents(t, f.archive)
			path := archivePathWithPrefix(t, entries, "journal/")
			data, hash := mutateArchiveEvent(t, entries[path], "task.approved", func(e *journal.Event) {
				if field == "retry" {
					e.Effects[1].Values["result_json"] = "{}"
					return
				}
				if field == "agent" {
					e.ActorV2.ActorType = "agent"
					e.Effects[0].Values["actor"] = "agent:mark"
				} else {
					var b approvalrecord.Binding
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
