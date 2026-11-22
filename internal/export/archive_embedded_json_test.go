package export

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mbell-nzl/fid/internal/config"
	fiddb "github.com/mbell-nzl/fid/internal/db"
	"github.com/mbell-nzl/fid/internal/gittest"
	"github.com/mbell-nzl/fid/internal/issues"
	"github.com/mbell-nzl/fid/internal/journal"
)

// An actual criterion writer supplies both a valid retry and its prefix. Check
// the reported placements through live commit, then the public backup/restore
// boundary: every accepted variant must validate and import.
func TestEmbeddedJSONLiveArchiveRoundTrip(t *testing.T) {
	t.Parallel()
	fixture := newExportSourceFixture(t, func(t *testing.T, root string) string {
		t.Helper()
		issue, err := issues.Create(issues.CreateOptions{WorkspaceRoot: root, MigrationsDir: testMigrationsDir(t), Kind: "task", Title: "Round trip", Actor: "human:mark", IdempotencyKey: "create"})
		if err != nil {
			t.Fatal(err)
		}
		if _, err := issues.AddAcceptanceCriterion(issues.AddAcceptanceCriterionOptions{WorkspaceRoot: root, MigrationsDir: testMigrationsDir(t), IssueID: issue.ID, Text: "Works", Ordinal: 1, Actor: "human:mark", IdempotencyKey: "ac"}); err != nil {
			t.Fatal(err)
		}
		return issue.ID
	})
	for _, placement := range []string{"values", "prior", "preconditions"} {
		for _, variant := range []string{"valid", "exact", "overflow", "depth"} {
			malformed := variant != "valid"
			reason := "duplicate key"
			switch variant {
			case "overflow":
				reason = "cannot unmarshal number"
			case "depth":
				reason = "exceeded max depth"
			}
			name := placement + "/valid"
			if malformed {
				name = placement + "/" + variant
			}
			t.Run(name, func(t *testing.T) {
				t.Parallel()
				root, _ := fixture.clone(t)
				files, err := filepath.Glob(filepath.Join(config.LedgerDir(root), "journal", "*.jsonl"))
				if err != nil || len(files) != 1 {
					t.Fatalf("journal files: %v %v", files, err)
				}
				data, err := os.ReadFile(files[0])
				if err != nil {
					t.Fatal(err)
				}
				lines := bytes.Split(bytes.TrimSpace(data), []byte{'\n'})
				var e journal.Event
				if err := e.UnmarshalJSON(lines[len(lines)-1]); err != nil {
					t.Fatal(err)
				}
				if e.EventType != "issue.ac-added" {
					t.Fatal("missing writer event")
				}
				prefix := append(bytes.Join(lines[:len(lines)-1], []byte{'\n'}), '\n')
				if err := os.WriteFile(files[0], prefix, 0o600); err != nil {
					t.Fatal(err)
				}
				if _, err := journal.RebuildWorkspace(root, testMigrationsDir(t)); err != nil {
					t.Fatal(err)
				}
				row := &e.Effects[len(e.Effects)-1]
				raw := row.Values["result_json"].(string)
				raw = strings.Replace(raw, `"replayed":false`, `"replayed":false,"extra":{"build":1,"Build":2}`, 1)
				switch variant {
				case "exact":
					raw = strings.Replace(raw, `"replayed":false`, `"replayed":true,"replayed":false`, 1)
				case "overflow":
					raw = strings.Replace(raw, `"replayed":false`, `"replayed":false,"overflow":1e999`, 1)
				case "depth":
					raw = strings.Replace(raw, `"replayed":false`, `"replayed":false,"deep":`+strings.Repeat("[", 10_001)+"0"+strings.Repeat("]", 10_001), 1)
				}
				switch placement {
				case "values":
					row.Values["result_json"] = raw
				case "prior":
					row.Prior = map[string]any{"result_json": raw}
				case "preconditions":
					row.Preconditions = map[string]any{"result_json": raw}
				}
				conn, err := fiddb.Open(config.WorkspacePathsForRoot(root).DatabasePath)
				if err != nil {
					t.Fatal(err)
				}
				defer func() { _ = conn.Close() }()
				err = journal.WithLedgerLock(config.LedgerDir(root), journal.LockOwner{Command: "test"}, func(lock *journal.WriteLock) error {
					_, err := journal.CommitEvent(lock, conn, e)
					return err
				})
				if malformed {
					if err == nil || !strings.Contains(err.Error(), reason) {
						t.Fatalf("commit = %v, want %s refusal", err, reason)
					}
					// Independently exercise public import with the malformed writer Event.
					original := archiveContents(t, fixture.archive)
					path := archivePathWithPrefix(t, original, "journal/")
					changed, hash := mutateArchiveEvent(t, original[path], "issue.ac-added", func(event *journal.Event) { *event = e })
					original[path] = changed
					manifest := decodeManifestMap(t, original["manifest.json"])
					manifest["terminal_event_hash"] = hash
					updateManifestEntry(t, manifest, path, changed)
					original["manifest.json"] = encodeManifestMap(t, manifest)
					archive := filepath.Join(t.TempDir(), "malformed.zip")
					writeTestArchive(t, archive, original)
					target := gittest.NewWorkspace(t)
					before := workspaceTreeDigest(t, target)
					_, err := Import(ImportOptions{ArchivePath: archive, RepositoryRoot: target, MigrationsDir: testMigrationsDir(t), IdempotencyKey: "malformed"})
					var refused *InvalidArchiveError
					if !errors.As(err, &refused) || len(refused.Checks) != 1 || !strings.Contains(refused.Checks[0].Detail, reason) {
						t.Fatalf("import = %v, checks %#v", err, refused)
					}
					if workspaceTreeDigest(t, target) != before {
						t.Fatal("refusal changed target")
					}
					return
				}
				if err != nil {
					t.Fatalf("valid JSON refused: %v", err)
				}
				archive := filepath.Join(t.TempDir(), "accepted.zip")
				if _, err := Create(CreateOptions{WorkspaceRoot: root, ArchivePath: archive}); err != nil {
					t.Fatal(err)
				}
				if _, err := Validate(ValidateOptions{ArchivePath: archive, MigrationsDir: testMigrationsDir(t)}); err != nil {
					t.Fatal(err)
				}
				target := gittest.NewWorkspace(t)
				if _, err := Import(ImportOptions{ArchivePath: archive, RepositoryRoot: target, MigrationsDir: testMigrationsDir(t), IdempotencyKey: "accepted"}); err != nil {
					t.Fatal(err)
				}
				if _, err := issues.Create(issues.CreateOptions{WorkspaceRoot: target, MigrationsDir: testMigrationsDir(t), Kind: "task", Title: "After restore", Actor: "human:mark", IdempotencyKey: "later"}); err != nil {
					t.Fatal(err)
				}
			})
		}
	}
}
