package config

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mbell-nzl/fid/internal/gittest"
)

// Every other test finds ledger files through LedgerDir, so these tests are
// what pin where the ledger is: the fid directory in the git common
// directory, which every worktree of a repository shares (ADR-0013).

func TestLedgerDirIsInTheGitCommonDirectoryOfEveryWorktree(t *testing.T) {
	t.Parallel()
	repository := gittest.NewWorkspace(t)
	worktree := gittest.AddWorktree(t, repository)
	output, err := exec.Command("git", "-C", worktree, "rev-parse", "--path-format=absolute", "--git-common-dir").Output()
	if err != nil {
		t.Fatalf("git rev-parse --git-common-dir: %v", err)
	}
	want := filepath.Join(strings.TrimSpace(string(output)), "fid")
	for _, root := range []string{gittest.ResolvedPath(t, repository), gittest.ResolvedPath(t, worktree)} {
		if got := LedgerDir(root); got != want {
			t.Errorf("LedgerDir(%q) = %q, want %q", root, got, want)
		}
	}
}

// LedgerDir reads the worktree's .git entry the way git does, without running
// git: a directory is the common directory, unless it names another in a
// commondir file; a file names the git directory, relative to the worktree.
func TestLedgerDirReadsTheWorktreesGitEntry(t *testing.T) {
	t.Parallel()
	for _, test := range []struct {
		name  string
		files map[string]string
		dirs  []string
		root  string
		want  string
	}{
		{
			name: "a main worktree's git directory",
			dirs: []string{"work/.git"},
			root: "work",
			want: "work/.git/fid",
		},
		{
			name: "a linked worktree's absolute gitdir and relative commondir",
			files: map[string]string{
				"linked/.git":                          "gitdir: $TMP/main/.git/worktrees/linked\n",
				"main/.git/worktrees/linked/commondir": "../..\n",
			},
			root: "linked",
			want: "main/.git/fid",
		},
		{
			name: "a relative gitdir, read from the worktree, with CRLF line ends",
			files: map[string]string{
				"linked/.git":                          "gitdir: ../main/.git/worktrees/linked\r\n",
				"main/.git/worktrees/linked/commondir": "../..\r\n",
			},
			root: "linked",
			want: "main/.git/fid",
		},
		{
			name: "an absolute commondir",
			files: map[string]string{
				"linked/.git":                          "gitdir: $TMP/main/.git/worktrees/linked\n",
				"main/.git/worktrees/linked/commondir": "$TMP/elsewhere/common.git\n",
			},
			root: "linked",
			want: "elsewhere/common.git/fid",
		},
		{
			name:  "a gitdir with no commondir, as a submodule or a separate git directory has",
			files: map[string]string{"sub/.git": "gitdir: ../modules/sub\n"},
			dirs:  []string{"modules/sub"},
			root:  "sub",
			want:  "modules/sub/fid",
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			temp := t.TempDir()
			for _, dir := range test.dirs {
				if err := os.MkdirAll(filepath.Join(temp, filepath.FromSlash(dir)), 0o755); err != nil {
					t.Fatalf("MkdirAll(%q) error = %v", dir, err)
				}
			}
			for name, contents := range test.files {
				path := filepath.Join(temp, filepath.FromSlash(name))
				if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
					t.Fatalf("MkdirAll(%q) error = %v", filepath.Dir(path), err)
				}
				if err := os.WriteFile(path, []byte(strings.ReplaceAll(contents, "$TMP", temp)), 0o644); err != nil {
					t.Fatalf("WriteFile(%q) error = %v", path, err)
				}
			}
			want := filepath.Join(temp, filepath.FromSlash(test.want))
			if got := LedgerDir(filepath.Join(temp, test.root)); got != want {
				t.Fatalf("LedgerDir() = %q, want %q", got, want)
			}
		})
	}
}

// LedgerDir returns no error. Where root has no .git entry it can read, it
// names the ledger a main worktree would have, so a missing or broken entry
// never redirects fid somewhere else.
func TestLedgerDirFallsBackToTheMainWorktreesLocation(t *testing.T) {
	t.Parallel()
	for _, test := range []struct{ root, want string }{
		{root: "/work/project", want: "/work/project/.git/fid"},
		{root: "relative/project", want: "relative/project/.git/fid"},
		{root: ".", want: ".git/fid"},
	} {
		if got := LedgerDir(filepath.FromSlash(test.root)); got != filepath.FromSlash(test.want) {
			t.Errorf("LedgerDir(%q) = %q, want %q", test.root, got, filepath.FromSlash(test.want))
		}
	}
	for name, contents := range map[string]string{
		"an empty .git file":               "",
		"a .git file that names no gitdir": "not a gitdir line\n",
		"a .git file with an empty gitdir": "gitdir: \n",
	} {
		root := t.TempDir()
		if err := os.WriteFile(filepath.Join(root, ".git"), []byte(contents), 0o644); err != nil {
			t.Fatalf("%s: WriteFile error = %v", name, err)
		}
		if got, want := LedgerDir(root), filepath.Join(root, ".git", "fid"); got != want {
			t.Errorf("%s: LedgerDir() = %q, want %q", name, got, want)
		}
	}
}
