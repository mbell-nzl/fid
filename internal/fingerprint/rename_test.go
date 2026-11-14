package fingerprint_test

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mbell-nzl/fid/internal/gittest"
)

//nolint:paralleltest // Supplies producer output through the git shim.
func TestComputeIgnoresSimilarityExtendedHeaders(t *testing.T) {
	root, base := fixture(t)
	content := "similarity index 97%\ndissimilarity index 3%\n"
	write(t, root, "file", content)
	head := commit(t, root)
	want := compute(t, root, base, head)
	patch := referenceFilePatch(t, root, base, head, "file")
	// Header-looking added content must remain significant.
	for _, line := range []string{"similarity index 97%", "dissimilarity index 3%"} {
		write(t, root, "file", strings.ReplaceAll(content, line, line+" edited"))
		if got := compute(t, root, base, commit(t, root)); got.Fingerprint == want.Fingerprint {
			t.Fatalf("ignored hunk content %q", line)
		}
	}
	installGitShim(t)
	dir := t.TempDir()
	t.Setenv("FID_TEST_DIFF", filepath.Join(dir, "diff"))
	t.Setenv("FID_TEST_GIT_MODE", "fixture-diff")
	for _, header := range []string{"similarity index 97%\n", "dissimilarity index 3%\n", "similarity index 99%\ndissimilarity index 1%\n"} {
		t.Run(strings.TrimSpace(header), func(t *testing.T) {
			firstLine, rest, _ := strings.Cut(patch, "\n")
			write(t, dir, "diff", firstLine+"\n"+header+rest)
			if got := compute(t, root, base, head); got.Fingerprint != want.Fingerprint {
				t.Fatalf("similarity metadata changed fingerprint: %s != %s", got.Fingerprint, want.Fingerprint)
			}
		})
	}
}

func TestComputeRenamesSurviveCleanRebaseAndBindPaths(t *testing.T) {
	t.Parallel()
	seen := map[string]string{}
	for _, row := range []struct {
		name, beforeSimilarity, afterSimilarity string
		edit                                    bool
	}{
		{"pure", "100%", "100%", false},
		{"edited", "97%", "99%", true},
	} {
		for _, target := range []string{"renamed", "other-target"} {
			t.Run(row.name+"/"+target, func(t *testing.T) {
				root := gittest.NewWorkspace(t)
				var lines strings.Builder
				for i := range 40 {
					fmt.Fprintf(&lines, "line %02d\n", i)
				}
				original := lines.String()
				write(t, root, "file", original)
				base := commit(t, root)
				git(t, root, "checkout", "-q", "-b", "task")
				git(t, root, "mv", "file", target)
				if row.edit {
					write(t, root, target, strings.ReplaceAll(original, "line 25", "edit 25"))
				}
				head := commit(t, root)
				beforePatch := referenceFilePatch(t, root, base, head, ".")
				assertContains(t, beforePatch, "similarity index "+row.beforeSimilarity+"\n")
				want := compute(t, root, base, head)
				if previous, ok := seen[want.Fingerprint]; ok {
					t.Fatalf("same fingerprint for %s and %s", previous, t.Name())
				}
				seen[want.Fingerprint] = t.Name()
				if repeat := compute(t, root, base, head); repeat != want {
					t.Fatalf("unstable rename: %+v != %+v", repeat, want)
				}
				git(t, root, "checkout", "-q", "--detach", base)
				write(t, root, "file", strings.Repeat("upstream line\n", 40)+original)
				upstream := commit(t, root)
				git(t, root, "checkout", "-q", "task")
				git(t, root, "rebase", upstream)
				rebased := git(t, root, "rev-parse", "HEAD")
				afterPatch := referenceFilePatch(t, root, upstream, rebased, ".")
				assertContains(t, afterPatch, "similarity index "+row.afterSimilarity+"\n")
				if row.edit {
					// The text index changes with upstream blob contents too.
					if patchLine(beforePatch, "index ") == patchLine(afterPatch, "index ") {
						t.Fatal("rebase did not change the text index header")
					}
				}
				got := compute(t, root, upstream, rebased)
				if got.Fingerprint != want.Fingerprint || head == rebased {
					t.Fatalf("clean rename rebase changed fingerprint: %s -> %s", want.Fingerprint, got.Fingerprint)
				}
				if repeat := compute(t, root, upstream, rebased); repeat != got {
					t.Fatalf("unstable rebased rename: %+v != %+v", repeat, got)
				}
			})
		}
	}
}

func assertContains(t *testing.T, patch, line string) {
	t.Helper()
	if !strings.Contains(patch, line) {
		t.Fatalf("patch missing %q:\n%s", line, patch)
	}
}

func patchLine(patch, prefix string) string {
	for line := range strings.SplitSeq(patch, "\n") {
		if strings.HasPrefix(line, prefix) {
			return line
		}
	}
	return ""
}
