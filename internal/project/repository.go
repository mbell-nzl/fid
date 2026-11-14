package project

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"unicode/utf8"

	"github.com/mbell-nzl/fid/internal/journal"
)

const maxRepositoryPaths = 10000
const maxRepositoryOutputBytes = 10 << 20

var ErrUnsafeRepository = errors.New("repository observation is unsafe")

// RepositoryUnavailableError classifies a failed capture without relying on
// Git's platform-specific message text.
type RepositoryUnavailableError struct {
	Reason string
	Detail string
}

func (err *RepositoryUnavailableError) Error() string {
	return fmt.Sprintf("%s: %s: %s", ErrUnsafeRepository, err.Reason, err.Detail)
}

func (err *RepositoryUnavailableError) Is(target error) bool { return target == ErrUnsafeRepository }

type Image struct {
	Mode   string `json:"mode"`
	Digest string `json:"digest"`
}

type Entry struct {
	Path     string `json:"path"`
	Head     *Image `json:"head"`
	Index    *Image `json:"index"`
	Worktree *Image `json:"worktree"`
}

type Observation struct {
	Root        string   `json:"-"`
	WorktreeID  string   `json:"worktree_id"`
	HeadOID     string   `json:"head_oid"`
	TreeOID     string   `json:"tree_oid"`
	Dirty       bool     `json:"dirty"`
	DirtyDigest *string  `json:"dirty_digest"`
	Entries     []Entry  `json:"entries"`
	Operation   string   `json:"operation"`
	Unmerged    []string `json:"unmerged_paths"`
}

// ObserveRepository reads a complete Git worktree snapshot without modifying it.
func ObserveRepository(root string) (Observation, error) {
	return captureRepository(func() (Observation, error) { return observeRepository(root) })
}

func captureRepository(read func() (Observation, error)) (Observation, error) {
	first, err := read()
	if err != nil {
		return Observation{}, err
	}
	second, err := read()
	if err != nil {
		return Observation{}, err
	}
	if !reflect.DeepEqual(first, second) {
		return Observation{}, unsafeRepository("repository changed during capture")
	}
	return first, nil
}

func observeRepository(root string) (Observation, error) {
	root, err := filepath.Abs(root)
	if err != nil {
		return Observation{}, fmt.Errorf("resolve repository root: %w", err)
	}
	root, err = gitText(root, "rev-parse", "--show-toplevel")
	if err != nil {
		return Observation{}, unsafeRepositoryReason("not-a-repository", "repository root: %v", err)
	}
	head, err := gitText(root, "rev-parse", "HEAD")
	if err != nil {
		return Observation{}, unsafeRepositoryReason("not-a-repository", "missing base: %v", err)
	}
	tree, err := gitText(root, "rev-parse", "HEAD^{tree}")
	if err != nil {
		return Observation{}, unsafeRepositoryReason("not-a-repository", "missing tree: %v", err)
	}
	if unsafe, err := trackedFlags(root); err != nil {
		return Observation{}, err
	} else if unsafe {
		return Observation{}, unsafeRepositoryReason("tracked-flags", "tracked assume-unchanged or skip-worktree entry")
	}
	unmerged, err := unmergedPaths(root)
	if err != nil {
		return Observation{}, err
	}
	if len(unmerged) != 0 {
		return Observation{}, unsafeRepositoryReason("unmerged-paths", "repository has unmerged paths")
	}
	operation, err := operationMarker(root)
	if err != nil {
		return Observation{}, err
	}
	paths, err := changedPaths(root)
	if err != nil {
		return Observation{}, err
	}
	entries := make([]Entry, 0, len(paths))
	for _, changed := range paths {
		entry, err := observeEntry(root, changed)
		if err != nil {
			return Observation{}, err
		}
		entries = append(entries, entry)
	}
	observation := Observation{Root: root, WorktreeID: WorktreeID(root), HeadOID: head, TreeOID: tree, Entries: entries, Operation: operation, Unmerged: unmerged}
	if len(entries) == 0 {
		return observation, nil
	}
	payload, err := journal.MarshalCanonicalJSON(map[string]any{"schema": "fid.repository-dirty.v1", "entries": entries})
	if err != nil {
		return Observation{}, fmt.Errorf("canonical dirty digest: %w", err)
	}
	digest := sha256Digest(payload)
	observation.Dirty, observation.DirtyDigest = true, &digest
	return observation, nil
}

// IsAncestor reports whether earlier is an ancestor of later without changing refs.
func IsAncestor(root, earlier, later string) (bool, error) {
	command := exec.Command("git", "-C", root, "merge-base", "--is-ancestor", earlier, later)
	err := command.Run()
	if err == nil {
		return true, nil
	}
	var exitError *exec.ExitError
	if errors.As(err, &exitError) && exitError.ExitCode() == 1 {
		return false, nil
	}
	return false, unsafeRepository("check ancestry: %v", err)
}

func changedPaths(root string) ([]string, error) {
	staged, err := gitNUL(root, "diff", "--cached", "--name-only", "-z", "--no-renames", "HEAD", "--")
	if err != nil {
		return nil, unsafeRepository("enumerate index changes: %v", err)
	}
	tracked, err := gitNUL(root, "diff", "--name-only", "-z", "--no-renames", "--")
	if err != nil {
		return nil, unsafeRepository("enumerate tracked changes: %v", err)
	}
	untracked, err := gitNUL(root, "ls-files", "--others", "--exclude-standard", "-z")
	if err != nil {
		return nil, unsafeRepository("enumerate untracked changes: %v", err)
	}
	paths := map[string]bool{}
	for _, raw := range append(append(staged, tracked...), untracked...) {
		if err := validPath(raw); err != nil {
			return nil, err
		}
		paths[raw] = true
	}
	if len(paths) > maxRepositoryPaths {
		return nil, unsafeRepositoryReason("unsafe-path", "path count exceeds %d", maxRepositoryPaths)
	}
	result := make([]string, 0, len(paths))
	for changed := range paths {
		result = append(result, changed)
	}
	sort.Slice(result, func(i, j int) bool { return bytes.Compare([]byte(result[i]), []byte(result[j])) < 0 })
	return result, nil
}

func observeEntry(root, changed string) (Entry, error) {
	head, err := gitImage(root, "ls-tree", "-z", "HEAD", "--", changed)
	if err != nil {
		return Entry{}, err
	}
	index, err := gitImage(root, "ls-files", "-s", "-z", "--", changed)
	if err != nil {
		return Entry{}, err
	}
	worktree, err := worktreeImage(root, changed)
	if err != nil {
		return Entry{}, err
	}
	return Entry{Path: changed, Head: head, Index: index, Worktree: worktree}, nil
}

func gitImage(root string, arguments ...string) (*Image, error) {
	output, err := gitOutput(root, arguments...)
	if err != nil {
		return nil, unsafeRepository("read Git image: %v", err)
	}
	if len(output) == 0 {
		return nil, nil
	}
	record := strings.TrimSuffix(string(output), "\x00")
	parts := strings.SplitN(record, "\t", 2)
	fields := strings.Fields(parts[0])
	if len(parts) != 2 || len(fields) != 3 {
		return nil, unsafeRepositoryReason("unsafe-path", "unsupported Git image for %q", arguments[len(arguments)-1])
	}
	mode, objectID := fields[0], fields[1]
	if arguments[0] == "ls-tree" {
		if fields[1] != "blob" {
			return nil, unsafeRepositoryReason("unsafe-path", "unsupported Git object type for %q", arguments[len(arguments)-1])
		}
		objectID = fields[2]
	} else if fields[2] != "0" {
		return nil, unsafeRepositoryReason("unmerged-paths", "unmerged Git image for %q", arguments[len(arguments)-1])
	}
	if mode != "100644" && mode != "100755" && mode != "120000" {
		return nil, unsafeRepositoryReason("unsafe-path", "unsupported Git mode %q", mode)
	}
	contents, err := gitOutput(root, "cat-file", "blob", objectID)
	if err != nil {
		return nil, unsafeRepository("read Git blob: %v", err)
	}
	return &Image{Mode: mode, Digest: sha256Digest(contents)}, nil
}

func worktreeImage(root, changed string) (*Image, error) {
	filename := filepath.Join(root, filepath.FromSlash(changed))
	info, err := os.Lstat(filename)
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, unsafeRepository("read worktree image %q: %v", changed, err)
	}
	mode := "100644"
	var contents []byte
	switch {
	case info.Mode()&os.ModeSymlink != 0:
		mode = "120000"
		target, err := os.Readlink(filename)
		if err != nil {
			return nil, unsafeRepository("read symlink %q: %v", changed, err)
		}
		contents = []byte(target)
	case info.Mode().IsRegular():
		if info.Mode().Perm()&0o111 != 0 {
			mode = "100755"
		}
		contents, err = os.ReadFile(filename)
		if err != nil {
			return nil, unsafeRepository("read worktree file %q: %v", changed, err)
		}
	default:
		return nil, unsafeRepositoryReason("unsafe-path", "unsupported worktree file type %q", changed)
	}
	return &Image{Mode: mode, Digest: sha256Digest(contents)}, nil
}

func unmergedPaths(root string) ([]string, error) {
	records, err := gitNUL(root, "ls-files", "-u", "-z")
	if err != nil {
		return nil, unsafeRepository("list unmerged entries: %v", err)
	}
	return parseUnmergedPaths(records)
}

func parseUnmergedPaths(records []string) ([]string, error) {
	paths := map[string]bool{}
	for _, record := range records {
		parts := strings.SplitN(record, "\t", 2)
		if len(parts) != 2 || validPath(parts[1]) != nil {
			return nil, unsafeRepositoryReason("unmerged-paths", "invalid unmerged path")
		}
		paths[parts[1]] = true
	}
	result := make([]string, 0, len(paths))
	for changed := range paths {
		result = append(result, changed)
	}
	sort.Strings(result)
	return result, nil
}

func operationMarker(root string) (string, error) {
	gitDir, err := gitText(root, "rev-parse", "--git-dir")
	if err != nil {
		return "", unsafeRepository("resolve Git directory: %v", err)
	}
	if !filepath.IsAbs(gitDir) {
		gitDir = filepath.Join(root, gitDir)
	}
	for marker, operation := range map[string]string{"MERGE_HEAD": "merge", "REBASE_HEAD": "rebase", "CHERRY_PICK_HEAD": "cherry-pick", "REVERT_HEAD": "revert", "BISECT_LOG": "bisect", "rebase-merge": "rebase", "rebase-apply": "rebase"} {
		if _, err := os.Lstat(filepath.Join(gitDir, marker)); err == nil {
			return operation, nil
		} else if !errors.Is(err, os.ErrNotExist) {
			return "", unsafeRepository("inspect Git operation marker: %v", err)
		}
	}
	return "none", nil
}

func trackedFlags(root string) (bool, error) {
	// Only the flag byte is needed here. Clean committed paths are not encoded
	// in an observation, so their arbitrary Git filename bytes are safe to keep.
	records, err := gitNULBytes(root, "ls-files", "-v", "-z")
	if err != nil {
		return false, unsafeRepository("list tracked flags: %v", err)
	}
	for _, record := range records {
		if len(record) > 0 && (record[0] == 'S' || (record[0] >= 'a' && record[0] <= 'z')) {
			return true, nil
		}
	}
	return false, nil
}

func gitNUL(root string, arguments ...string) ([]string, error) {
	parts, err := gitNULBytes(root, arguments...)
	if err != nil {
		return nil, err
	}
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		if !utf8.Valid(part) {
			return nil, unsafeRepositoryReason("unsafe-path", "Git path is not valid UTF-8")
		}
		result = append(result, string(part))
	}
	return result, nil
}

func gitNULBytes(root string, arguments ...string) ([][]byte, error) {
	output, err := gitOutput(root, arguments...)
	if err != nil {
		return nil, err
	}
	if len(output) == 0 {
		return [][]byte{}, nil
	}
	if output[len(output)-1] != 0 {
		return nil, unsafeRepository("Git NUL output was incomplete")
	}
	parts := bytes.Split(output[:len(output)-1], []byte{0})
	if len(parts) > maxRepositoryPaths {
		return nil, unsafeRepositoryReason("unsafe-path", "path count exceeds %d", maxRepositoryPaths)
	}
	return parts, nil
}

func gitText(root string, arguments ...string) (string, error) {
	output, err := gitOutput(root, arguments...)
	if err != nil {
		return "", err
	}
	return strings.TrimSuffix(strings.TrimSuffix(string(output), "\n"), "\r"), nil
}

func gitOutput(root string, arguments ...string) ([]byte, error) {
	command := exec.Command("git", append([]string{"-C", root}, arguments...)...)
	output, err := command.Output()
	if err != nil {
		return nil, err
	}
	if err := checkRepositoryOutputSize(len(output)); err != nil {
		return nil, err
	}
	return output, nil
}

func checkRepositoryOutputSize(size int) error {
	if size > maxRepositoryOutputBytes {
		return unsafeRepositoryReason("unsafe-path", "Git output exceeds %d bytes", maxRepositoryOutputBytes)
	}
	return nil
}

func validPath(value string) error {
	if value == "" || !utf8.ValidString(value) || strings.ContainsRune(value, 0) || strings.HasPrefix(value, "/") || path.Clean(value) != value || value == "." || strings.HasPrefix(value, "../") {
		return unsafeRepositoryReason("unsafe-path", "invalid repository path")
	}
	return nil
}

// WorktreeID identifies the worktree whose top level is root: sha256 of the
// path, so a moved worktree is another worktree.
func WorktreeID(root string) string {
	return sha256Digest([]byte(root))
}

func sha256Digest(contents []byte) string {
	sum := sha256.Sum256(contents)
	return "sha256:" + hex.EncodeToString(sum[:])
}

func unsafeRepository(format string, arguments ...any) error {
	reason := "changed-during-capture"
	for _, argument := range arguments {
		if err, ok := argument.(error); ok {
			var unavailable *RepositoryUnavailableError
			if errors.As(err, &unavailable) {
				reason = unavailable.Reason
			}
		}
	}
	return unsafeRepositoryReason(reason, format, arguments...)
}

func unsafeRepositoryReason(reason, format string, arguments ...any) error {
	return &RepositoryUnavailableError{Reason: reason, Detail: fmt.Sprintf(format, arguments...)}
}
