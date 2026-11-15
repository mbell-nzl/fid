package export

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/mbell-nzl/fid/internal/config"
	"github.com/mbell-nzl/fid/internal/gittest"
	"github.com/mbell-nzl/fid/internal/holdwork"
	"github.com/mbell-nzl/fid/internal/journal"
	"github.com/mbell-nzl/fid/internal/reviewrecord"
	"github.com/mbell-nzl/fid/internal/reviewwork"
)

func addClosureExportFacts(t *testing.T, root string) string {
	t.Helper()
	task := addReviewExportFacts(t, root)
	conn, err := reviewwork.OpenProjection(root)
	if err != nil {
		t.Fatal(err)
	}
	rounds, err := reviewrecord.Read(conn, task)
	_ = conn.Close()
	if err != nil {
		t.Fatal(err)
	}
	first := rounds[0]
	in := first.Input()
	in.Dispositions = []reviewrecord.Disposition{{FindingID: first.Findings[0].ID, Disposition: "resolved", Reason: "Regression passes"}, {FindingID: first.Findings[1].ID, Disposition: "not-applicable", Reason: "Outside this change"}}
	in.Findings = []reviewrecord.FindingInput{first.Findings[0].FindingInput, first.Findings[0].FindingInput}
	result, err := reviewwork.Record(reviewwork.Options{WorkspaceRoot: root, MigrationsDir: testMigrationsDir(t), TaskID: task, Actor: "human:courier", IdempotencyKey: "close-by-round", Input: in})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := holdwork.CloseFinding(holdwork.Options{WorkspaceRoot: root, MigrationsDir: testMigrationsDir(t), FindingID: result.Round.Findings[0].ID, Actor: "human:mark", Reason: "Accepted risk", IdempotencyKey: "close-by-person"}); err != nil {
		t.Fatal(err)
	}
	return task
}

func closureSnapshot(t *testing.T, root string) string {
	t.Helper()
	conn, err := reviewwork.OpenProjection(root)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = conn.Close() }()
	rows, err := conn.Query(`SELECT finding_id,round_id,decision_id,disposition,reason FROM finding_closures ORDER BY finding_id`)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = rows.Close() }()
	all := []map[string]any{}
	for rows.Next() {
		var finding, disposition, reason string
		var round, decision any
		if err := rows.Scan(&finding, &round, &decision, &disposition, &reason); err != nil {
			t.Fatal(err)
		}
		all = append(all, map[string]any{"finding": finding, "round": round, "decision": decision, "disposition": disposition, "reason": reason})
	}
	if err := rows.Err(); err != nil {
		t.Fatal(err)
	}
	raw, err := json.Marshal(all)
	if err != nil {
		t.Fatal(err)
	}
	return string(raw)
}

func TestFindingClosureExportImportRebuild(t *testing.T) {
	t.Parallel()
	f := newExportSourceFixture(t, addClosureExportFacts)
	before := closureSnapshot(t, f.workspace)
	for _, text := range []string{"resolved", "not-applicable", "closed", "PDN-", "RND-"} {
		if !strings.Contains(before, text) {
			t.Fatalf("fixture missing %s: %s", text, before)
		}
	}
	if _, err := journal.RebuildWorkspace(f.workspace, testMigrationsDir(t)); err != nil {
		t.Fatal(err)
	}
	if closureSnapshot(t, f.workspace) != before {
		t.Fatal("rebuild changed closures")
	}
	target := gittest.NewWorkspace(t)
	if _, err := Import(ImportOptions{ArchivePath: f.archive, RepositoryRoot: target, MigrationsDir: testMigrationsDir(t), IdempotencyKey: "closure-import"}); err != nil {
		t.Fatal(err)
	}
	if closureSnapshot(t, target) != before {
		t.Fatal("import changed closures")
	}
	assertImportedJournalEquals(t, f.workspace, target)
}

func TestFindingClosureImportRefusesWithoutPublication(t *testing.T) {
	t.Parallel()
	f := newExportSourceFixture(t, addClosureExportFacts)
	for _, field := range []string{"actor", "binding", "reason", "closure"} {
		t.Run(field, func(t *testing.T) {
			t.Parallel()
			entries := archiveContents(t, f.archive)
			path := archivePathWithPrefix(t, entries, "journal/")
			data, hash := mutateArchiveEvent(t, entries[path], "finding.closed", func(e *journal.Event) {
				switch field {
				case "actor":
					e.ActorV2.ActorType = "agent"
					e.Effects[0].Values["actor"] = "agent:mark"
				case "binding":
					e.Effects[0].Values["binding_id"] = "FND-missing"
				case "reason":
					e.Effects[0].Values["reason"] = " "
					e.Effects[1].Values["reason"] = " "
				case "closure":
					e.Effects[1].Values["decision_id"] = "PDN-forged"
				}
			})
			entries[path] = data
			manifest := decodeManifestMap(t, entries["manifest.json"])
			manifest["terminal_event_hash"] = hash
			updateManifestEntry(t, manifest, path, data)
			entries["manifest.json"] = encodeManifestMap(t, manifest)
			archive := filepath.Join(t.TempDir(), "forged.zip")
			writeTestArchive(t, archive, entries)
			target := gittest.NewWorkspace(t)
			before := workspaceTreeDigest(t, target)
			_, err := Import(ImportOptions{ArchivePath: archive, RepositoryRoot: target, MigrationsDir: testMigrationsDir(t), IdempotencyKey: "refused"})
			assertArchiveCheck(t, err, "reference-invalid")
			if !reflect.DeepEqual(before, workspaceTreeDigest(t, target)) {
				t.Fatal("refusal changed target")
			}
			if _, err := os.Stat(config.LedgerDir(target)); !os.IsNotExist(err) {
				t.Fatalf("published ledger: %v", err)
			}
		})
	}
}
