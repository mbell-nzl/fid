package defaultconflict

import (
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/mbell-nzl/fid/internal/gittest"
	"github.com/mbell-nzl/fid/internal/policy"
)

func TestMain(m *testing.M) { os.Exit(gittest.Main(m)) }

func TestInspectPreservesGitPathBytes(t *testing.T) {
	t.Parallel()
	root := gittest.NewWorkspace(t)
	git := func(args ...string) string {
		t.Helper()
		cmd := exec.Command("git", append([]string{"-C", root, "-c", "user.name=fid test", "-c", "user.email=fid@example.test"}, args...)...)
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("git %v: %v: %s", args, err, out)
		}
		return strings.TrimSpace(string(out))
	}
	git("commit", "--allow-empty", "-qm", "Base")
	git("branch", "base")
	paths := []string{"line\nname", "x💛", "x\xff"}
	commit := func(content string) {
		t.Helper()
		for _, path := range paths {
			if err := os.WriteFile(filepath.Join(root, path), []byte(content), 0o600); err != nil {
				t.Fatal(err)
			}
		}
		git("add", ".")
		git("commit", "-qm", content)
	}
	commit("Builder\n")
	head := git("rev-parse", "HEAD")
	git("checkout", "--detach", "base")
	commit("Default\n")
	git("update-ref", "refs/remotes/origin/trunk", "HEAD")
	git("symbolic-ref", "refs/remotes/origin/HEAD", "refs/remotes/origin/trunk")
	got, err := Inspect(root, "HOF-abcdef", head)
	if err != nil || got == nil || !reflect.DeepEqual(got.ConflictingPaths, paths) || !reflect.DeepEqual(policy.PathsFromBytes(got.PathBytes), paths) {
		t.Fatalf("conflicting paths=%+v/%v", got, err)
	}
	if git("status", "--porcelain") != "" || git("rev-parse", "HEAD") != got.DefaultCommit {
		t.Fatal("merge test changed checkout")
	}
}
