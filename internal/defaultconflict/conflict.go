// Package defaultconflict tests a hand-off against the locally known default
// branch. It never fetches or changes the worktree, index or refs.
package defaultconflict

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"slices"
	"strings"

	"github.com/mbell-nzl/fid/internal/fingerprint"
	"github.com/mbell-nzl/fid/internal/policy"
)

// Evidence is the optional reason next returns approved work to its builder.
// PathBytes preserves non-UTF-8 Git names alongside the readable paths.
type Evidence struct {
	Cause            string   `json:"cause"`
	DefaultRef       string   `json:"default_ref"`
	DefaultCommit    string   `json:"default_commit"`
	HandoffID        string   `json:"handoff_id"`
	Head             string   `json:"head"`
	ConflictingPaths []string `json:"conflicting_paths"`
	PathBytes        [][]byte `json:"conflicting_path_bytes,omitempty"`
}

// Record keeps the fingerprint observation needed for admission without git.
type Record struct {
	Evidence
	Binding fingerprint.Result `json:"binding"`
}

// Inspect uses the same origin/HEAD resolution as policy. A missing ref or
// object is unavailable evidence, never proof of a conflict.
func Inspect(root, handoff, head string) (*Evidence, error) {
	source, err := policy.DefaultSource(root)
	if err != nil {
		return nil, err
	}
	cmd := exec.Command("git", "-C", root, "merge-tree", "--write-tree", "--name-only", "-z", "--no-messages", source.Commit, head)
	cmd.Env = append(os.Environ(), "GIT_NO_LAZY_FETCH=1", "GIT_NO_REPLACE_OBJECTS=1")
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	out, err := cmd.Output()
	if err == nil {
		return nil, nil
	}
	var exit *exec.ExitError
	if !errors.As(err, &exit) || exit.ExitCode() != 1 {
		return nil, fmt.Errorf("local merge unavailable: %w: %s", err, strings.TrimSpace(stderr.String()))
	}
	fields := bytes.Split(out, []byte{0})
	paths := []string{}
	for _, field := range fields[1:] {
		if len(field) == 0 {
			break
		}
		paths = append(paths, string(field))
	}
	if len(paths) == 0 {
		return nil, fmt.Errorf("local merge unavailable: git reported no conflicting paths")
	}
	slices.Sort(paths)
	e := &Evidence{Cause: "default-conflict", DefaultRef: source.Ref, DefaultCommit: source.Commit, HandoffID: handoff, Head: head, ConflictingPaths: paths}
	for _, path := range paths {
		if strings.ToValidUTF8(path, "�") != path {
			e.PathBytes = policy.PathBytes(paths)
			break
		}
	}
	return e, nil
}
