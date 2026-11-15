package fingerprint_test

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/mbell-nzl/fid/internal/fingerprint"
	"github.com/mbell-nzl/fid/internal/gittest"
)

func TestMain(m *testing.M) { os.Exit(gittest.Main(m)) }

func git(t *testing.T, root string, args ...string) string {
	t.Helper()
	command := exec.Command("git", append([]string{"-C", root, "-c", "user.name=fid test", "-c", "user.email=fid@example.test", "-c", "commit.gpgsign=false"}, args...)...)
	output, err := command.CombinedOutput()
	if err != nil {
		t.Fatalf("git %v: %v: %s", args, err, output)
	}
	return strings.TrimSpace(string(output))
}

func write(t *testing.T, root, name, content string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(root, name), []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
}

func commit(t *testing.T, root string) string {
	t.Helper()
	git(t, root, "add", "--all")
	git(t, root, "commit", "-q", "--no-verify", "-m", "fixture")
	return git(t, root, "rev-parse", "HEAD")
}

func fixture(t *testing.T) (root, base string) {
	t.Helper()
	root = gittest.NewWorkspace(t)
	write(t, root, "file", "before\n")
	return root, commit(t, root)
}

func compute(t *testing.T, root, base, head string) fingerprint.Result {
	t.Helper()
	got, err := fingerprint.Compute(root, base, head)
	if err != nil {
		t.Fatal(err)
	}
	if got.Algorithm != "patch-id-verbatim-per-file-v1" || !regexp.MustCompile(`^[0-9a-f]{64}$`).MatchString(got.Fingerprint) || got.Base != base || got.Head != head {
		t.Fatalf("incomplete binding: %+v", got)
	}
	return got
}

func refuse(t *testing.T, root, base, head string, kind fingerprint.FailureKind) {
	t.Helper()
	got, err := fingerprint.Compute(root, base, head)
	var failure *fingerprint.Failure
	if !errors.As(err, &failure) || failure.Kind != kind || got != (fingerprint.Result{}) {
		t.Fatalf("got %+v, %v; want no result and %s", got, err, kind)
	}
}

func TestComputeSurvivesCleanRebase(t *testing.T) {
	t.Parallel()
	root := gittest.NewWorkspace(t)
	var lines strings.Builder
	for i := range 40 {
		fmt.Fprintf(&lines, "line %02d\n", i)
	}
	original := lines.String()
	write(t, root, "file", original)
	base := commit(t, root)
	git(t, root, "checkout", "-q", "-b", "task")
	write(t, root, "file", strings.ReplaceAll(original, "line 25", "edited 25"))
	head := commit(t, root)
	want := compute(t, root, base, head)
	git(t, root, "checkout", "-q", "--detach", base)
	write(t, root, "file", "upstream line\n"+original)
	write(t, root, "upstream", "unrelated file\n")
	upstream := commit(t, root)
	git(t, root, "checkout", "-q", "task")
	git(t, root, "rebase", upstream)
	rebased := git(t, root, "rev-parse", "HEAD")
	got := compute(t, root, upstream, rebased)
	if got.Fingerprint != want.Fingerprint || head == rebased {
		t.Fatalf("clean rebase changed fingerprint: %s -> %s", want.Fingerprint, got.Fingerprint)
	}
}

func TestComputeContentEdits(t *testing.T) {
	t.Parallel()
	for _, row := range []struct{ name, content string }{
		{"text", "another\n"}, {"space", "af ter\n"}, {"tab", "af\tter\n"},
		{"trailing space", "after \n"}, {"trailing tab", "after\t\n"}, {"no final newline", "after"}, {"extra newline", "after\n\n"},
	} {
		t.Run(row.name, func(t *testing.T) {
			t.Parallel()
			root, base := fixture(t)
			write(t, root, "file", "after\n")
			want := compute(t, root, base, commit(t, root))
			write(t, root, "file", row.content)
			head := commit(t, root)
			got := compute(t, root, base, head)
			if want.Fingerprint == got.Fingerprint {
				t.Fatalf("edit %q kept fingerprint %s", row.content, got.Fingerprint)
			}
			if repeat := compute(t, root, base, head); repeat != got {
				t.Fatalf("unstable result: %+v != %+v", repeat, got)
			}
		})
	}
}

func TestComputeUsesAggregateChange(t *testing.T) {
	t.Parallel()
	root, base := fixture(t)
	write(t, root, "file", "after\n")
	head := commit(t, root)
	want := compute(t, root, base, head)
	git(t, root, "revert", "--no-edit", head)
	write(t, root, "file", "after\n")
	got := compute(t, root, base, commit(t, root))
	if got.Fingerprint != want.Fingerprint {
		t.Fatalf("intermediate history changed aggregate fingerprint: %+v != %+v", got, want)
	}
	git(t, root, "checkout", "-q", "--detach", base)
	write(t, root, "file", "uncommitted noise\n")
	if repeat := compute(t, root, base, head); repeat != want {
		t.Fatalf("checkout changed recorded comparison: %+v != %+v", repeat, want)
	}
}

func TestComputeSpecialChanges(t *testing.T) {
	t.Parallel()
	seen := map[string]string{}
	for _, row := range []struct {
		name            string
		prepare, change func(*testing.T, string)
	}{
		{"binary add A", nil, func(t *testing.T, r string) { write(t, r, "binary", "\x00binary A") }},
		{"binary add B", nil, func(t *testing.T, r string) { write(t, r, "binary", "\x00binary B") }},
		{"binary change A", binaryBase, func(t *testing.T, r string) { write(t, r, "binary", "\x00binary A") }},
		{"binary change B", binaryBase, func(t *testing.T, r string) { write(t, r, "binary", "\x00binary B") }},
		{"binary delete", binaryBase, func(t *testing.T, r string) { git(t, r, "rm", "binary") }},
		{"rename", nil, func(t *testing.T, r string) { git(t, r, "mv", "file", "renamed") }},
		{"rename with edit", nil, func(t *testing.T, r string) {
			git(t, r, "mv", "file", "renamed")
			write(t, r, "renamed", "before\nafter\n")
		}},
		{"mode", nil, func(t *testing.T, r string) {
			if err := os.Chmod(filepath.Join(r, "file"), 0o755); err != nil {
				t.Fatal(err)
			}
		}},
		{"symlink", func(t *testing.T, r string) { symlink(t, r, "file") }, func(t *testing.T, r string) { git(t, r, "rm", "link"); symlink(t, r, "renamed") }},
		{"gitlink", submoduleBase, func(t *testing.T, r string) {
			sub := filepath.Join(r, "module")
			write(t, sub, "file", "submodule edit\n")
			commit(t, sub)
		}},
	} {
		t.Run(row.name, func(t *testing.T) {
			root, base := fixture(t)
			if row.prepare != nil {
				row.prepare(t, root)
				base = commit(t, root)
			}
			row.change(t, root)
			head := commit(t, root)
			got := compute(t, root, base, head)
			if previous, ok := seen[got.Fingerprint]; ok {
				t.Fatalf("same fingerprint for %s and %s", row.name, previous)
			}
			seen[got.Fingerprint] = row.name
			if repeat := compute(t, root, base, head); repeat != got {
				t.Fatalf("unstable: %+v != %+v", repeat, got)
			}
		})
	}
}

func binaryBase(t *testing.T, root string) { t.Helper(); write(t, root, "binary", "\x00original") }
func symlink(t *testing.T, root, target string) {
	t.Helper()
	if err := os.Symlink(target, filepath.Join(root, "link")); err != nil {
		t.Fatal(err)
	}
}
func submoduleBase(t *testing.T, root string) {
	t.Helper()
	sub := filepath.Join(root, "module")
	if err := os.Mkdir(sub, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := gittest.InitRepository(sub); err != nil {
		t.Fatal(err)
	}
	write(t, sub, "file", "submodule base\n")
	commit(t, sub)
}

func TestComputeRefusals(t *testing.T) {
	t.Parallel()
	root, base := fixture(t)
	for _, row := range []struct {
		name, root, base, head string
		kind                   fingerprint.FailureKind
	}{
		{"empty", root, base, base, fingerprint.EmptyChange},
		{"missing base", root, strings.Repeat("1", 40), base, fingerprint.MissingObject},
		{"missing head", root, base, strings.Repeat("2", 40), fingerprint.MissingObject},
		{"invalid oid", root, "HEAD", base, fingerprint.GitFailure},
		{"blob is not commit", root, git(t, root, "rev-parse", base+":file"), base, fingerprint.GitFailure},
		{"not repository", t.TempDir(), base, base, fingerprint.GitFailure},
	} {
		t.Run(row.name, func(t *testing.T) { refuse(t, row.root, row.base, row.head, row.kind) })
	}
	write(t, root, "file", strings.Repeat("large line\n", fingerprint.MaxDiffBytes/10+1))
	refuse(t, root, base, commit(t, root), fingerprint.OversizedDiff)
}
