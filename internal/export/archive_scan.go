package export

import (
	"archive/zip"
	"bufio"
	"errors"
	"io"
	"path/filepath"
	"strings"
)

var errArchiveResourceLimit = errors.New("archive resource limit exceeded")

const maxArchiveNonArtifactBytes int64 = 64 << 20
const maxArchiveArtifactBytes int64 = 1 << 30
const maxArchiveExpandedBytes int64 = 8 << 30
const maxArchiveEntries = 1000000

type archiveMember struct {
	file     *zip.File
	path     string
	artifact bool
}
type archiveScan struct {
	ExpandedBytes int64
	EntryCount    int
}
type archiveLimits struct {
	nonArtifactBytes int64
	artifactBytes    int64
	expandedBytes    int64
	entries          int
}

func productionArchiveLimits() archiveLimits {
	return archiveLimits{maxArchiveNonArtifactBytes, maxArchiveArtifactBytes, maxArchiveExpandedBytes, maxArchiveEntries}
}

type archiveByteBudget struct {
	used   int64
	limits archiveLimits
}

func newArchiveByteBudget() *archiveByteBudget {
	return &archiveByteBudget{limits: productionArchiveLimits()}
}

func inspectArchiveDirectory(files []*zip.File) ([]archiveMember, []Check) {
	return inspectArchiveDirectoryWithLimit(files, maxArchiveEntries)
}

func inspectArchiveDirectoryWithLimit(files []*zip.File, limit int) ([]archiveMember, []Check) {
	if len(files) > limit {
		return nil, []Check{check("resource-limit", nil, "archive entry count exceeds limit")}
	}
	seen := map[string]bool{}
	members := make([]archiveMember, 0, len(files))
	checks := []Check{}
	for _, file := range files {
		if !safeArchivePath(file.Name) || !file.Mode().IsRegular() {
			checks = append(checks, check("path-unsafe", file.Name, "entry is not a safe regular file"))
			continue
		}
		if seen[file.Name] {
			checks = append(checks, check("duplicate-entry", file.Name, "entry appears more than once"))
			continue
		}
		seen[file.Name] = true
		members = append(members, archiveMember{file: file, path: file.Name, artifact: strings.HasPrefix(file.Name, "artifacts/")})
	}
	return members, checks
}

func scanArchive(reader *zip.Reader, consume func(archiveMember, io.Reader) error) (archiveScan, []Check) {
	return scanArchiveWithLimits(reader, productionArchiveLimits(), consume)
}

func scanArchiveWithLimits(reader *zip.Reader, limits archiveLimits, consume func(archiveMember, io.Reader) error) (archiveScan, []Check) {
	members, checks := inspectArchiveDirectoryWithLimit(reader.File, limits.entries)
	if len(checks) != 0 {
		return archiveScan{}, checks
	}
	budget := &archiveByteBudget{limits: limits}
	scan := archiveScan{EntryCount: len(members)}
	for _, member := range members {
		n, err := streamArchiveMember(member, budget, func(input io.Reader) error { return consume(member, input) })
		scan.ExpandedBytes += n
		if err != nil {
			code := "digest-mismatch"
			if errors.Is(err, errArchiveResourceLimit) {
				code = "resource-limit"
			}
			checks = append(checks, check(code, member.path, err.Error()))
			return scan, checks
		}
	}
	return scan, checks
}

func streamArchiveMember(member archiveMember, budget *archiveByteBudget, consume func(io.Reader) error) (int64, error) {
	input, err := member.file.Open()
	if err != nil {
		return 0, err
	}
	defer func() { _ = input.Close() }()
	limit := budget.limits.nonArtifactBytes
	if member.artifact {
		limit = budget.limits.artifactBytes
	}
	reader := &archiveLimitReader{Reader: input, limit: limit, budget: budget}
	if err := consume(reader); err != nil {
		return reader.n, err
	}
	_, err = io.Copy(io.Discard, reader)
	return reader.n, err
}

type archiveLimitReader struct {
	io.Reader
	limit  int64
	budget *archiveByteBudget
	n      int64
}

func (r *archiveLimitReader) Read(p []byte) (int, error) {
	remaining := r.limit - r.n
	total := r.budget.limits.expandedBytes - r.budget.used
	if total < remaining {
		remaining = total
	}
	if remaining < 0 {
		return 0, errArchiveResourceLimit
	}
	if int64(len(p)) > remaining+1 {
		p = p[:remaining+1]
	}
	n, err := r.Reader.Read(p)
	r.n += int64(n)
	r.budget.used += int64(n)
	if r.n > r.limit || r.budget.used > r.budget.limits.expandedBytes {
		return n, errArchiveResourceLimit
	}
	return n, err
}

func readEventLines(body io.Reader, consume func([]byte) error) error {
	r := bufio.NewReader(body)
	line := make([]byte, 0, 4096)
	for {
		fragment, err := r.ReadSlice('\n')
		line = append(line, fragment...)
		if len(line) > 1<<20+1 {
			return errArchiveResourceLimit
		}
		if err == bufio.ErrBufferFull {
			continue
		}
		if err != nil && err != io.EOF {
			return err
		}
		if len(line) != 0 {
			value := line
			if value[len(value)-1] == '\n' {
				value = value[:len(value)-1]
			}
			if len(value) > 1<<20 {
				return errArchiveResourceLimit
			}
			if err := consume(value); err != nil {
				return err
			}
			line = line[:0]
		}
		if err == io.EOF {
			return nil
		}
	}
}

func safeArchivePath(value string) bool {
	if value == "" || value == "." || strings.ContainsRune(value, 0) || filepath.IsAbs(value) || strings.HasSuffix(value, "/") || strings.Contains(value, "\\") || (len(value) >= 2 && value[1] == ':') {
		return false
	}
	cleaned := filepath.ToSlash(filepath.Clean(value))
	return cleaned == value && cleaned != ".." && !strings.HasPrefix(cleaned, "../")
}
