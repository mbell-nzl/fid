package policy

import (
	"fmt"
	"os"
	"os/exec"
	"slices"
	"strings"

	"github.com/mbell-nzl/fid/internal/config"
)

type Source struct {
	Ref    string `json:"ref"`
	Commit string `json:"commit,omitempty"`
}

type Loaded struct {
	Source    Source
	Effective Effective
}

// SourceError is a policy refusal, including the source ref/commit when known.
type SourceError struct {
	Problem string
	Ref     string
	Commit  string
}

const (
	DefaultBranchFix = "git remote set-head origin <branch>"
	PolicyFileFix    = "fix .fid.toml on the default branch, merge it and git fetch"
)

func (err *SourceError) NextAction() string {
	if err.Ref == "" {
		return DefaultBranchFix
	}
	return PolicyFileFix
}

func (err *SourceError) Error() string { return err.Problem }

// Load uses only local git plumbing. It never consults the worktree or index,
// nor runs fetch. An origin with no known default is refused rather than guessed.
func Load(root string) (Loaded, error) {
	remotes, err := gitPolicy(root, "remote")
	if err != nil {
		return Loaded{}, sourceError(Source{}, "read policy remotes: %v", err)
	}
	if !slices.Contains(strings.Fields(string(remotes)), "origin") {
		return builtInDefaults()
	}
	target, err := gitPolicy(root, "symbolic-ref", "--quiet", "refs/remotes/origin/HEAD")
	ref := strings.TrimSpace(string(target))
	if err != nil || !strings.HasPrefix(ref, "refs/remotes/origin/") || ref == "refs/remotes/origin/HEAD" {
		return Loaded{}, sourceError(Source{}, "origin's default branch is unknown")
	}
	commitBytes, err := gitPolicy(root, "rev-parse", "--verify", ref+"^{commit}")
	if err != nil {
		return Loaded{}, sourceError(Source{Ref: strings.TrimPrefix(ref, "refs/remotes/")}, "read policy: %v", err)
	}
	commit := strings.TrimSpace(string(commitBytes))
	source := Source{Ref: strings.TrimPrefix(ref, "refs/remotes/"), Commit: commit}
	entry, err := gitPolicy(root, "ls-tree", commit, "--", config.ProjectIdentityFileName)
	if err != nil {
		return Loaded{}, sourceError(source, "read .fid.toml: %v", err)
	}
	if len(entry) == 0 {
		return builtInDefaults()
	}
	// Refuse a symlink or directory at the policy path, as identity readers do.
	if !strings.HasPrefix(string(entry), "100644 blob ") && !strings.HasPrefix(string(entry), "100755 blob ") {
		return Loaded{}, sourceError(source, "invalid .fid.toml: expected a regular file")
	}
	data, err := gitPolicy(root, "cat-file", "blob", commit+":"+config.ProjectIdentityFileName)
	if err != nil {
		return Loaded{}, sourceError(source, "read .fid.toml: %v", err)
	}
	file, err := config.DecodeProject(data)
	if err != nil {
		return Loaded{}, sourceError(source, "invalid .fid.toml: %s", strings.TrimPrefix(err.Error(), "parse .fid.toml: "))
	}
	return Loaded{Source: source, Effective: Resolve(file)}, nil
}

func builtInDefaults() (Loaded, error) {
	file, err := config.DecodeProject(nil)
	if err != nil {
		return Loaded{}, err
	}
	return Loaded{Source: Source{Ref: "built-in defaults"}, Effective: Resolve(file)}, nil
}

func sourceError(source Source, format string, args ...any) error {
	return &SourceError{Problem: fmt.Sprintf(format, args...), Ref: source.Ref, Commit: source.Commit}
}

func gitPolicy(root string, args ...string) ([]byte, error) {
	command := exec.Command("git", append([]string{"-C", root}, args...)...)
	// Reading a missing promisor object must not cause git to fetch it for us.
	command.Env = append(os.Environ(), "GIT_NO_LAZY_FETCH=1")
	output, err := command.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("git %s: %w: %s", strings.Join(args, " "), err, strings.TrimSpace(string(output)))
	}
	return output, nil
}
