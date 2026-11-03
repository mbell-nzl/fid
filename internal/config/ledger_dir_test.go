package config

import (
	"path/filepath"
	"testing"
)

// Every other test finds ledger files through LedgerDir, so this test is what
// pins where the ledger is. Until RESET task 1.1d moves it to the Git common
// directory, it is the .fid directory under the project root.
func TestLedgerDirIsTheFidDirectoryUnderTheRoot(t *testing.T) {
	t.Parallel()
	for _, test := range []struct{ root, want string }{
		{root: "/work/project", want: "/work/project/.fid"},
		{root: "relative/project", want: "relative/project/.fid"},
		{root: ".", want: ".fid"},
	} {
		if got := LedgerDir(filepath.FromSlash(test.root)); got != filepath.FromSlash(test.want) {
			t.Errorf("LedgerDir(%q) = %q, want %q", test.root, got, test.want)
		}
	}
}
