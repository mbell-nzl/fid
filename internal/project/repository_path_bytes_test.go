package project

import "testing"

func TestObserveRepositoryWithNonUTF8CommittedPath(t *testing.T) {
	t.Parallel()
	for _, row := range []struct {
		name, reason string
	}{
		{"clean", ""},
		{"assume-unchanged", "tracked-flags"},
		{"skip-worktree", "tracked-flags"},
		{"dirty", "unsafe-path"},
	} {
		t.Run(row.name, func(t *testing.T) {
			t.Parallel()
			const path = "committed-\xff.go"
			root := initRepositoryWithFiles(t, map[string]string{path: "base\n"})
			switch row.name {
			case "assume-unchanged", "skip-worktree":
				runGit(t, root, "update-index", "--"+row.name, "--", path)
			case "dirty":
				writeFixtureFile(t, root, path, "edited\n")
			}
			if row.reason != "" {
				assertRepositoryReason(t, root, row.reason)
				return
			}
			got, err := ObserveRepository(root)
			if err != nil || got.Dirty || len(got.Entries) != 0 {
				t.Fatalf("clean committed byte path observation=%+v/%v", got, err)
			}
		})
	}
}
