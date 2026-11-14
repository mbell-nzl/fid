//go:build integration

package integration

import (
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"syscall"
	"testing"

	"github.com/mbell-nzl/fid/internal/config"
)

func TestP22RecoveryPreservesRepositoryOnSuccessAndUnsafeRejection(t *testing.T) {
	t.Parallel()
	cases := []string{"clean", "dirty-in-scope", "staged-restored", "assume-unchanged", "skip-worktree", "unmerged", "missing-head", "unrelated-history", "foreign-worktree", "gitlink", "unreadable", "fifo", "escaping-parent", "protected-adoption", "MERGE_HEAD", "REBASE_HEAD", "CHERRY_PICK_HEAD", "REVERT_HEAD", "BISECT_LOG", "rebase-merge", "rebase-apply"}
	for _, kind := range cases {
		t.Run(kind, func(t *testing.T) {
			a := p22NewFIDRecoveryAdapter(t)
			a.p22Observe(t, "terminated")
			a.p22Expire(t)
			p22Write(t, a.workspace, "base.txt", "stash retained\n")
			p22Git(t, a.workspace, "stash", "push", "-qm", "retained safety fixture")
			p22Git(t, a.workspace, "update-ref", "refs/fq3/retained", "HEAD")
			adopted := []string{}
			wantSuccess := kind == "clean" || kind == "dirty-in-scope" || kind == "staged-restored"
			if kind == "dirty-in-scope" || kind == "staged-restored" {
				p22Write(t, a.workspace, "base.txt", "adopted content\n")
				if kind == "staged-restored" {
					p22Git(t, a.workspace, "add", "base.txt")
					p22Write(t, a.workspace, "base.txt", "base\n")
				}
				adopted = []string{"base.txt"}
				inspection := a.p22Run(t, "run", "inspect", a.runID)
				if inspection.ExitCode != 0 {
					t.Fatalf("inspect dirty: %#v", inspection)
				}
				a.p22Advance(t, "run.inspect", inspection.Stdout)
			} else if kind == "protected-adoption" {
				adopted = []string{".git/fid/config.toml"}
			} else {
				p22MakeRepositoryUnsafe(t, a, kind)
			}
			fixture := p22ConformanceCase{Name: "repository-safety", Operation: "run.recover", Input: map[string]any{"adopted_paths": adopted}}
			args, input := a.p22Args(t, fixture)
			// Place caller input outside source so its creation is not mistaken
			// for a recovery write. Every source path, even ignored ones, is compared.
			args = append(args, "--input", p21Input(t, t.TempDir(), "recover", input), "--idempotency-key", "repository-safety")
			before := p22CompleteRepositorySnapshot(t, a.workspace)
			beforeEvents := p21JournalEventCount(t, a.workspace)
			stdout, stderr, err := closureRunCommand(a.workspace, a.binary, args...)
			response := p22ConformanceResponse{Stdout: stdout, Stderr: stderr}
			if err != nil {
				response.ExitCode = processExitCode(t, err)
			}
			if after := p22CompleteRepositorySnapshot(t, a.workspace); !reflect.DeepEqual(after, before) {
				t.Fatal("recovery changed source bytes/modes, HEAD, complete refs, index, stash, or Git metadata")
			}
			if kind == "unreadable" {
				name := filepath.Join(a.workspace, "base.txt")
				if err := os.Chmod(name, 0o644); err != nil {
					t.Fatal(err)
				}
				contents, err := os.ReadFile(name)
				if err != nil || string(contents) != "unreadable changed bytes" {
					t.Fatalf("unreadable source bytes changed: %q, %v", contents, err)
				}
			}
			if wantSuccess {
				if response.ExitCode != 0 {
					t.Fatalf("valid recovery: %#v", response)
				}
				if p21JournalEventCount(t, a.workspace) != beforeEvents+1 {
					t.Fatal("successful recovery did not append exactly one Event")
				}
			} else {
				p22AssertWireFailure(t, "run.recover", "recovery-state-conflict", response)
				if p21JournalEventCount(t, a.workspace) != beforeEvents {
					t.Fatal("unsafe recovery appended Event/takeover facts")
				}
			}
		})
	}
}

func p22MakeRepositoryUnsafe(t *testing.T, a *p22FIDRecoveryAdapter, kind string) {
	t.Helper()
	root := a.workspace
	switch kind {
	case "clean":
	case "assume-unchanged", "skip-worktree":
		p22Git(t, root, "update-index", "--"+kind, "base.txt")
	case "unmerged":
		blob := p22GitOutput(t, root, "rev-parse", "HEAD:base.txt")
		cmd := exec.Command("git", "-C", root, "update-index", "--index-info")
		cmd.Stdin = strings.NewReader("0 " + strings.Repeat("0", 40) + "\tbase.txt\n100644 " + blob + " 1\tbase.txt\n")
		if output, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("unmerged index: %v %s", err, output)
		}
	case "missing-head":
		p22Git(t, root, "symbolic-ref", "HEAD", "refs/heads/missing")
	case "unrelated-history":
		p22Git(t, root, "checkout", "--orphan", "unrelated")
		p22Git(t, root, "commit", "-qm", "unrelated base")
	case "foreign-worktree":
		moved := filepath.Join(t.TempDir(), "moved")
		if err := os.Rename(root, moved); err != nil {
			t.Fatal(err)
		}
		a.workspace = moved
	case "gitlink":
		p22Git(t, root, "update-index", "--add", "--cacheinfo", "160000,"+p22GitOutput(t, root, "rev-parse", "HEAD")+",module")
	case "unreadable":
		p22Write(t, root, "base.txt", "unreadable changed bytes")
		if err := os.Chmod(filepath.Join(root, "base.txt"), 0); err != nil {
			t.Fatal(err)
		}
		t.Cleanup(func() { _ = os.Chmod(filepath.Join(root, "base.txt"), 0o644) })
	case "fifo":
		if err := os.Remove(filepath.Join(root, "base.txt")); err != nil {
			t.Fatal(err)
		}
		if err := syscall.Mkfifo(filepath.Join(root, "base.txt"), 0o600); err != nil {
			t.Fatal(err)
		}
	case "escaping-parent":
		if err := os.Mkdir(filepath.Join(root, "dir"), 0o755); err != nil {
			t.Fatal(err)
		}
		p22Write(t, root, "dir/child", "inside")
		p22Git(t, root, "add", "dir/child")
		p22Git(t, root, "commit", "-qm", "nested")
		if err := os.RemoveAll(filepath.Join(root, "dir")); err != nil {
			t.Fatal(err)
		}
		outside := t.TempDir()
		p22Write(t, outside, "child", "outside secret")
		if err := os.Symlink(outside, filepath.Join(root, "dir")); err != nil {
			t.Fatal(err)
		}
	default:
		marker := filepath.Join(root, ".git", kind)
		if strings.HasPrefix(kind, "rebase-") {
			if err := os.Mkdir(marker, 0o755); err != nil {
				t.Fatal(err)
			}
		} else if err := os.WriteFile(marker, []byte(p22GitOutput(t, root, "rev-parse", "HEAD")+"\n"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
}

func p22CompleteRepositorySnapshot(t *testing.T, root string) map[string]string {
	t.Helper()
	result := map[string]string{}
	err := filepath.WalkDir(root, func(name string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		relative, err := filepath.Rel(root, name)
		if err != nil {
			return err
		}
		if name == config.LedgerDir(root) {
			return filepath.SkipDir
		}
		info, err := entry.Info()
		if err != nil {
			return err
		}
		value := info.Mode().String()
		switch {
		case info.Mode()&os.ModeSymlink != 0:
			target, err := os.Readlink(name)
			if err != nil {
				return err
			}
			value += ":" + target
		case info.Mode().IsRegular():
			contents, err := os.ReadFile(name)
			if err != nil && !os.IsPermission(err) {
				return err
			}
			value += fmt.Sprintf(":%x:%v", contents, err != nil)
		}
		result[relative] = value
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	return result
}
