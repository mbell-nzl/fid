// Package fingerprint binds reviews and approvals to an aggregate committed
// change. It has no ledger or command adapter; callers supply the hand-off's
// recorded base and head, never a newly computed merge-base.
package fingerprint

import (
	"bytes"
	"crypto/sha256"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
)

const (
	Algorithm    = "patch-id-verbatim-per-file-v1"
	MaxDiffBytes = 10 << 20
)

type Result struct {
	Algorithm, Fingerprint, Base, Head string
}

type FailureKind string

const (
	MissingObject FailureKind = "missing-object"
	EmptyChange   FailureKind = "empty-change"
	OversizedDiff FailureKind = "oversized-diff"
	GitFailure    FailureKind = "git-failure"
)

// Failure is a closed refusal: no Result accompanies it. Err retains git's
// diagnostic, including a producer that emitted a valid prefix before failing.
type Failure struct {
	Kind FailureKind
	Err  error
}

func (f *Failure) Error() string { return fmt.Sprintf("fingerprint %s: %v", f.Kind, f.Err) }
func (f *Failure) Unwrap() error { return f.Err }

var oid = regexp.MustCompile(`^[0-9a-f]{40}$`)
var patchID = regexp.MustCompile(`^([0-9a-f]{40}) 0000000000000000000000000000000000000000\n$`)

// Compute reads only the recorded comparison. The complete diff must fit within
// MaxDiffBytes and its producer must succeed before patch-id receives any bytes.
// This prevents patch-id's successful hashing of a prefix from becoming evidence.
func Compute(root, base, head string) (Result, error) {
	if !oid.MatchString(base) || !oid.MatchString(head) {
		return Result{}, &Failure{Kind: GitFailure, Err: fmt.Errorf("base and head must be full lowercase 40-hex commit IDs")}
	}
	view, err := objectView(root, head)
	if err != nil {
		return Result{}, &Failure{Kind: GitFailure, Err: err}
	}
	defer func() { _ = os.RemoveAll(view) }()
	root = view
	if err := checkCommits(root, base, head); err != nil {
		return Result{}, err
	}
	if err := prepareAttributes(root, head); err != nil {
		return Result{}, &Failure{Kind: GitFailure, Err: err}
	}
	diff, err := runGit(root, nil, MaxDiffBytes, diffArgs(base, head)...)
	if err != nil {
		kind := GitFailure
		if errors.Is(err, errOutputLimit) {
			kind = OversizedDiff
		}
		return Result{}, &Failure{Kind: kind, Err: err}
	}
	if len(diff) == 0 {
		return Result{}, &Failure{Kind: EmptyChange, Err: fmt.Errorf("base and head have no change")}
	}
	// Committed attributes may force binary data into text hunks. patch-id
	// silently truncates such lines at NUL, including in --verbatim mode.
	if bytes.IndexByte(diff, 0) >= 0 {
		return Result{}, &Failure{Kind: GitFailure, Err: fmt.Errorf("diff contains NUL text that patch-id cannot hash completely")}
	}
	id, err := fingerprintDiff(root, diff)
	if err != nil {
		return Result{}, &Failure{Kind: GitFailure, Err: err}
	}
	return Result{Algorithm: Algorithm, Fingerprint: id, Base: base, Head: head}, nil
}

// Git's patch-id can skip the next file header after a binary patch (observed
// with Git 2.53). Give each file its own process so every path is hashed.
func fingerprintDiff(root string, diff []byte) (string, error) {
	if !bytes.HasPrefix(diff, []byte("diff --git ")) {
		return "", fmt.Errorf("diff does not start with a file section")
	}
	var ids []string
	for len(diff) > 0 {
		// Pinned text indicators prefix content with +, - or space; quoted
		// paths escape newlines and base85 binary data contains no spaces.
		// Only a real file header can start a line with this exact marker.
		end := len(diff)
		if next := bytes.Index(diff, []byte("\ndiff --git ")); next >= 0 {
			end = next + 1
		}
		output, err := runGit(root, withoutSimilarity(diff[:end]), 4096, "patch-id", "--verbatim")
		if err != nil {
			return "", err
		}
		match := patchID.FindSubmatch(output)
		if match == nil {
			return "", fmt.Errorf("patch-id did not return exactly one ID for file section %d", len(ids)+1)
		}
		ids = append(ids, string(match[1]))
		diff = diff[end:]
	}
	slices.Sort(ids)
	// Algorithm's SHA-256 input is the decimal section count followed by LF,
	// then every sorted 40-hex ID followed by LF, including the last one.
	input := fmt.Sprintf("%d\n", len(ids)) + strings.Join(ids, "\n") + "\n"
	return fmt.Sprintf("%x", sha256.Sum256([]byte(input))), nil
}

// Similarity scores depend on the size of the upstream file, so a clean rebase
// can change them without changing rename paths or hunks. Algorithm omits these
// redundant extended headers, retaining all content from the first patch line.
func withoutSimilarity(section []byte) []byte {
	var normalized bytes.Buffer
	for len(section) > 0 {
		line, rest, found := bytes.Cut(section, []byte("\n"))
		if bytes.HasPrefix(line, []byte("--- ")) || bytes.HasPrefix(line, []byte("+++ ")) ||
			bytes.HasPrefix(line, []byte("@@")) || bytes.Equal(line, []byte("GIT binary patch")) {
			normalized.Write(section)
			break
		}
		if !bytes.HasPrefix(line, []byte("similarity index ")) && !bytes.HasPrefix(line, []byte("dissimilarity index ")) {
			normalized.Write(line)
			if found {
				normalized.WriteByte('\n')
			}
		}
		section = rest
	}
	return normalized.Bytes()
}

// Materialize only committed attribute files, using a private index. git diff
// does not consistently fall back to index attributes when files are absent.
// This also works before Git's --attr-source option, without using the caller's
// checkout, staged attributes, info/attributes or configured diff drivers.
func prepareAttributes(root, head string) error {
	if _, err := runGit(root, nil, 1024, "read-tree", head); err != nil {
		return err
	}
	paths, err := runGit(root, nil, MaxDiffBytes, "ls-files", "-z", "--", ".gitattributes", "**/.gitattributes")
	if err != nil {
		return err
	}
	if len(paths) == 0 {
		return nil
	}
	_, err = runGit(root, paths, 1024, "checkout-index", "--stdin", "-z")
	return err
}

// Git has no switch to ignore repository config. In particular, named diff
// drivers can change binary classification even with --no-ext-diff/--no-textconv.
// This disposable view borrows only the source object directory (including
// its alternates), with private config, attributes and index. It never writes to
// the source repository and has no remotes from which to fetch.
func objectView(root, head string) (view string, err error) {
	objects, err := runGit(root, nil, 64<<10, "rev-parse", "--path-format=absolute", "--git-path", "objects")
	if err != nil {
		return "", err
	}
	view, err = os.MkdirTemp("", "fid-fingerprint-")
	if err != nil {
		return "", fmt.Errorf("create object view: %w", err)
	}
	defer func() {
		if err != nil {
			_ = os.RemoveAll(view)
		}
	}()
	gitdir := filepath.Join(view, ".git")
	if err = os.MkdirAll(filepath.Join(gitdir, "refs"), 0o700); err != nil {
		return view, fmt.Errorf("create refs: %w", err)
	}
	// Strip exactly the newline git appends: a repository path may end in one.
	if err = os.Symlink(strings.TrimSuffix(string(objects), "\n"), filepath.Join(gitdir, "objects")); err != nil {
		return view, fmt.Errorf("link objects: %w", err)
	}
	if err = os.WriteFile(filepath.Join(gitdir, "HEAD"), []byte(head+"\n"), 0o600); err != nil {
		return view, fmt.Errorf("write private HEAD: %w", err)
	}
	if err = os.WriteFile(filepath.Join(gitdir, "config"), []byte("[core]\n\tbare = false\n\trepositoryformatversion = 0\n"), 0o600); err != nil {
		return view, fmt.Errorf("write private config: %w", err)
	}
	return view, nil
}

func checkCommits(root, base, head string) error {
	output, err := runGit(root, []byte(base+"\n"+head+"\n"), 1024, "cat-file", "--batch-check=%(objectname) %(objecttype)")
	if err != nil {
		return &Failure{Kind: GitFailure, Err: err}
	}
	lines := strings.Split(strings.TrimSuffix(string(output), "\n"), "\n")
	if len(lines) != 2 {
		return &Failure{Kind: GitFailure, Err: fmt.Errorf("cat-file returned incomplete commit types")}
	}
	for i, id := range []string{base, head} {
		if lines[i] == id+" missing" {
			return &Failure{Kind: MissingObject, Err: fmt.Errorf("commit %s is unavailable locally", id)}
		}
		if lines[i] != id+" commit" {
			return &Failure{Kind: GitFailure, Err: fmt.Errorf("object %s is not a commit", id)}
		}
	}
	return nil
}

// These settings are part of Algorithm. Change its identifier if the rendering
// contract changes. --binary keeps full binary identities and contents; --full-index
// alone only pins the index lines. Rename search has no config-dependent limit.
func diffArgs(base, head string) []string {
	return []string{
		"diff", "--patch", "--no-color",
		"--no-ext-diff", "--no-textconv", "--no-relative", "--src-prefix=a/", "--dst-prefix=b/", "--line-prefix=",
		"--unified=3", "--inter-hunk-context=0", "--diff-algorithm=myers", "--indent-heuristic",
		"--find-renames=50%", "-l0", "--full-index", "--binary", "--ignore-submodules=none", "--submodule=short",
		"--output-indicator-new=+", "--output-indicator-old=-", "--output-indicator-context= ",
		"-O" + os.DevNull, base, head, "--",
	}
}

func runGit(root string, input []byte, limit int, args ...string) ([]byte, error) {
	argv := []string{"--no-pager", "--no-replace-objects", "-C", root,
		"-c", "core.quotePath=true", "-c", "core.attributesFile=" + os.DevNull, "-c", "diff.suppressBlankEmpty=false"}
	command := exec.Command("git", append(argv, args...)...)
	command.Env = gitEnvironment()
	command.Stdin = bytes.NewReader(input)
	stdout := &boundedBuffer{limit: limit, stop: true}
	stderr := &boundedBuffer{limit: 8192}
	command.Stdout, command.Stderr = stdout, stderr
	err := command.Run()
	if stdout.overflow {
		return nil, fmt.Errorf("git output exceeds %d bytes: %w", limit, errOutputLimit)
	}
	if err != nil {
		return nil, fmt.Errorf("git %s: %w: %s", args[0], err, strings.TrimSpace(stderr.buffer.String()))
	}
	return stdout.buffer.Bytes(), nil
}

func gitEnvironment() []string {
	var env []string
	for _, entry := range os.Environ() {
		// Keep the test/discovery guard, but no ambient git routing, injected
		// config, diff options, external drivers or lazy-fetch override.
		if !strings.HasPrefix(entry, "GIT_") || strings.HasPrefix(entry, "GIT_CEILING_DIRECTORIES=") {
			env = append(env, entry)
		}
	}
	return append(env, "GIT_CONFIG_NOSYSTEM=1", "GIT_CONFIG_GLOBAL="+os.DevNull,
		"GIT_ATTR_NOSYSTEM=1", "GIT_NO_LAZY_FETCH=1", "GIT_OPTIONAL_LOCKS=0", "GIT_TERMINAL_PROMPT=0", "LC_ALL=C")
}

var errOutputLimit = errors.New("output limit exceeded")

// Git's stdout fails closed at the limit. Its diagnostic is bounded separately
// and drained, so an oversized diagnostic cannot block an otherwise failing child.
type boundedBuffer struct {
	buffer         bytes.Buffer
	limit          int
	stop, overflow bool
}

func (b *boundedBuffer) Write(p []byte) (int, error) {
	n := len(p)
	remaining := b.limit - b.buffer.Len()
	if n > remaining {
		b.overflow = true
		p = p[:remaining]
	}
	_, _ = b.buffer.Write(p)
	if b.overflow && b.stop {
		return len(p), errOutputLimit
	}
	return n, nil
}
