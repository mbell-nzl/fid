package export

import (
	"bytes"
	"errors"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mbell-nzl/fid/internal/gittest"
	"github.com/mbell-nzl/fid/internal/issues"
)

func TestImportMalformedEventDuplicateJSONKeys(t *testing.T) {
	t.Parallel()
	var issueID string
	fixture := newExportSourceFixture(t, func(t *testing.T, root string) string {
		t.Helper()
		created, err := issues.Create(issues.CreateOptions{WorkspaceRoot: root, MigrationsDir: testMigrationsDir(t), Kind: "task", Title: "Task", Actor: "human:mark", IdempotencyKey: "create"})
		if err != nil {
			t.Fatal(err)
		}
		issueID = created.ID
		return issueID
	})
	original := archiveContents(t, fixture.archive)
	for _, tc := range []struct{ name, old, replacement string }{
		{"unicode header", `"event_schema":2`, `"event_schema":1,"event_ſchema":2`},
		{"unicode effects", `"effects":[`, `"effects":[],"effectſ":[`},
		{"unicode Effect", `"values":{`, `"values":{},"valueſ":{`},
		{"unicode aggregate", `"sequence":1`, `"sequence":0,"ſequence":1`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			entries := map[string][]byte{}
			for path, data := range original {
				entries[path] = bytes.Clone(data)
			}
			path := archivePathWithPrefix(t, entries, "journal/")
			data := bytes.Replace(entries[path], []byte(tc.old), []byte(tc.replacement), 1)
			if bytes.Equal(data, entries[path]) {
				t.Fatal("probe did not change the Event JSON")
			}
			entries[path] = data
			manifest := decodeManifestMap(t, entries["manifest.json"])
			// Each alias overwrites its earlier value with the original value,
			// so the decoded terminal Event's canonical hash is unchanged.
			updateManifestEntry(t, manifest, path, data)
			entries["manifest.json"] = encodeManifestMap(t, manifest)
			archive := filepath.Join(t.TempDir(), "malformed.zip")
			writeTestArchive(t, archive, entries)
			target := gittest.NewWorkspace(t)
			before := workspaceTreeDigest(t, target)
			_, err := Import(ImportOptions{ArchivePath: archive, RepositoryRoot: target, MigrationsDir: testMigrationsDir(t), IdempotencyKey: "malformed"})
			var refused *InvalidArchiveError
			if !errors.As(err, &refused) {
				t.Fatalf("import = %v, want typed archive refusal", err)
			}
			if len(refused.Checks) != 1 || !strings.Contains(refused.Checks[0].Detail, "duplicate key") {
				t.Errorf("import checks = %#v, want duplicate-key refusal", refused.Checks)
			}
			if before != workspaceTreeDigest(t, target) {
				t.Fatal("refused import changed target")
			}
			if _, err := Import(ImportOptions{ArchivePath: fixture.archive, RepositoryRoot: target, MigrationsDir: testMigrationsDir(t), IdempotencyKey: "valid"}); err != nil {
				t.Fatal(err)
			}
			title := "After refusal"
			if _, err := issues.Update(issues.UpdateOptions{WorkspaceRoot: target, MigrationsDir: testMigrationsDir(t), ID: issueID, Title: &title, Actor: "human:mark", IdempotencyKey: "later"}); err != nil {
				t.Fatalf("later valid write: %v", err)
			}
		})
	}
}
