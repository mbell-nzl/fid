package policy

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"slices"
	"strings"

	"github.com/mbell-nzl/fid/internal/checkrecord"
	"github.com/mbell-nzl/fid/internal/pathglob"
)

// Comparison describes committed changes against the current local default
// branch. Recompute it for each check and hand-off; start retains no base.
type Comparison struct {
	Ref, Base, Head string
	Paths           []string
}

// HeadTree reads only the committed tree, independently of worktree dirt or
// merge state. A missing promisor object must not cause a lazy fetch.
func HeadTree(root string) (string, error) {
	tree, err := gitPolicy(root, "rev-parse", "--verify", "HEAD^{tree}")
	return strings.TrimSpace(string(tree)), err
}

// Compare never fetches. Disabling rename detection includes both names as a
// deletion and an addition, and -z preserves literal repository path bytes.
func Compare(root, task string) (Comparison, error) {
	comparison := Comparison{}
	fail := func(problem string) (Comparison, error) {
		var ref any
		if comparison.Ref != "" {
			ref = comparison.Ref
		}
		return Comparison{}, &checkrecord.Refusal{Code: "comparison-unavailable", TaskID: task, Message: "cannot compare this task's committed changes", Details: map[string]any{"ref": ref, "problem": problem}}
	}
	target, err := gitPolicy(root, "symbolic-ref", "--quiet", "refs/remotes/origin/HEAD")
	ref := strings.TrimSpace(string(target))
	if err != nil || !strings.HasPrefix(ref, "refs/remotes/origin/") || ref == "refs/remotes/origin/HEAD" {
		return fail("origin's default branch is unknown")
	}
	comparison.Ref = strings.TrimPrefix(ref, "refs/remotes/")
	head, err := gitPolicy(root, "rev-parse", "--verify", "HEAD^{commit}")
	if err != nil {
		return fail(err.Error())
	}
	comparison.Head = strings.TrimSpace(string(head))
	base, err := gitPolicy(root, "merge-base", comparison.Head, ref)
	if err != nil {
		var exit *exec.ExitError
		if errors.As(err, &exit) && exit.ExitCode() == 1 {
			return fail("HEAD and " + comparison.Ref + " have no common ancestor in the available history")
		}
		return fail(strings.TrimSpace(err.Error()))
	}
	comparison.Base = strings.TrimSpace(string(base))
	// This is evidence, so local display settings must not filter, rewrite or
	// replace paths. In particular, include committed gitlinks even when git
	// config or .gitmodules asks to hide them. Ignore optional order files too.
	paths, err := gitPolicy(root, "diff", "--no-relative", "--no-ext-diff", "--no-textconv",
		"--no-renames", "--ignore-submodules=none", "-O"+os.DevNull, "--name-only", "-z", comparison.Base, comparison.Head, "--")
	if err != nil {
		return fail(err.Error())
	}
	comparison.Paths = []string{}
	if len(paths) > 0 {
		comparison.Paths = strings.Split(strings.TrimSuffix(string(paths), "\x00"), "\x00")
	}
	return comparison, nil
}

// Required selects the lane's configured checks against a freshly computed
// comparison. An omitted or empty when_changed always requires the check.
// Paths outside the matcher grammar conservatively require every conditional
// check: unsupported must never be mistaken for an evaluated nonmatch.
func Required(effective Effective, lane, task string, paths []string) ([]string, error) {
	configured := effective.Lanes[lane]
	if !configured.ChecksSet {
		return nil, &checkrecord.Refusal{Code: "checks-not-configured", TaskID: task, Message: "no checks set for " + lane + " in .fid.toml", Details: map[string]any{"lane": lane}, Next: "add checks to [lanes." + lane + "], then merge .fid.toml and git fetch"}
	}
	names := []string{}
	for _, name := range configured.Checks {
		check := effective.Checks[name]
		required := len(check.WhenChanged) == 0
		for _, pattern := range check.WhenChanged {
			glob, err := pathglob.Compile(pattern)
			if err != nil {
				return nil, fmt.Errorf("required check %s: %w", name, err)
			}
			for _, path := range paths {
				required = required || !pathglob.ValidPath(path) || glob.Match(path)
			}
		}
		if required {
			names = append(names, name)
		}
	}
	slices.Sort(names)
	return names, nil
}
