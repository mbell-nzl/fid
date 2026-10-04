package project

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var sourceDigestPattern = regexp.MustCompile(`^sha256:[0-9a-f]{64}$`)

// RecordedRepository is the repository identity attached to journal Events and Sources.
type RecordedRepository struct {
	BaseOID       *string `json:"base_oid"`
	CaptureSource string  `json:"capture_source"`
	Dirty         bool    `json:"dirty"`
	DirtyDigest   *string `json:"dirty_digest"`
	HeadOID       string  `json:"head_oid"`
	ProjectID     string  `json:"project_id"`
	TreeOID       string  `json:"tree_oid"`
	WorktreeID    string  `json:"worktree_id"`
}

// ObserveRecordedRepository captures one stable worktree and its project identity.
func ObserveRecordedRepository(root, projectID string) (RecordedRepository, error) {
	observation, err := ObserveRepository(root)
	if err != nil {
		return RecordedRepository{}, fmt.Errorf("observe repository: %w", err)
	}
	base := observation.HeadOID
	return RecordedRepository{BaseOID: &base, CaptureSource: "git", Dirty: observation.Dirty,
		DirtyDigest: observation.DirtyDigest, HeadOID: observation.HeadOID, ProjectID: projectID,
		TreeOID: observation.TreeOID, WorktreeID: observation.WorktreeID}, nil
}

// SourceInput is the caller-written portion of a Source. Digest is optional
// only for project-file sources; a nil pointer preserves omission in digests.
type SourceInput struct {
	SourceType   string  `json:"source_type"`
	DisplayLabel string  `json:"display_label"`
	Path         *string `json:"path"`
	Section      *string `json:"section"`
	ExternalID   *string `json:"external_id"`
	Digest       *string `json:"digest,omitempty"`
}

type SourceDigestMismatchError struct {
	DisplayLabel   string
	Path           *string
	ExpectedDigest *string
	ObservedDigest *string
}

func (err *SourceDigestMismatchError) Error() string {
	return fmt.Sprintf("Source %q digest does not match the observed project file", err.DisplayLabel)
}

type InvalidSourceInputError struct{ Field, Detail string }

func (err *InvalidSourceInputError) Error() string { return err.Detail }

// ResolveSourceInput checks a project file under root or accepts a supplied
// external digest. It never observes the repository and never mutates input.
func ResolveSourceInput(root string, source SourceInput) (string, error) {
	if source.DisplayLabel == "" || len(source.DisplayLabel) > 128 {
		return "", &InvalidSourceInputError{Field: "display_label", Detail: "Source display label is invalid"}
	}
	if source.Digest != nil && !sourceDigestPattern.MatchString(*source.Digest) {
		return "", &InvalidSourceInputError{Field: "digest", Detail: "Source digest must be a SHA-256 digest"}
	}
	if source.SourceType == "external-reference" {
		if source.Digest == nil || *source.Digest == "" {
			return "", &InvalidSourceInputError{Field: "digest", Detail: "external reference digest is required"}
		}
		if source.Path != nil || source.ExternalID == nil || *source.ExternalID == "" {
			return "", &InvalidSourceInputError{Field: "external_id", Detail: "external reference requires an external ID and no path"}
		}
		return *source.Digest, nil
	}
	if source.SourceType != "project-file" || source.Path == nil || source.ExternalID != nil {
		return "", &InvalidSourceInputError{Field: "path", Detail: "project-file path is required"}
	}
	value := *source.Path
	if value == "" || filepath.IsAbs(value) || strings.Contains(value, "\\") || strings.ContainsRune(value, 0) {
		return "", &InvalidSourceInputError{Field: "path", Detail: "project-file path must be relative"}
	}
	for _, segment := range strings.Split(value, "/") {
		if segment == "." || segment == ".." || segment == "" {
			return "", &InvalidSourceInputError{Field: "path", Detail: "project-file path contains an unsafe segment"}
		}
	}
	filename := filepath.Join(root, filepath.FromSlash(value))
	resolvedRoot, rootErr := filepath.EvalSymlinks(root)
	resolved, err := filepath.EvalSymlinks(filename)
	if rootErr != nil || err != nil {
		return "", &SourceDigestMismatchError{source.DisplayLabel, source.Path, source.Digest, nil}
	}
	relative, err := filepath.Rel(resolvedRoot, resolved)
	if err != nil || relative == ".." || strings.HasPrefix(relative, ".."+string(filepath.Separator)) {
		return "", &SourceDigestMismatchError{source.DisplayLabel, source.Path, source.Digest, nil}
	}
	info, err := os.Stat(resolved)
	if err != nil || !info.Mode().IsRegular() {
		return "", &SourceDigestMismatchError{source.DisplayLabel, source.Path, source.Digest, nil}
	}
	contents, err := os.ReadFile(resolved)
	if err != nil {
		return "", &SourceDigestMismatchError{source.DisplayLabel, source.Path, source.Digest, nil}
	}
	sum := sha256.Sum256(contents)
	observed := "sha256:" + hex.EncodeToString(sum[:])
	if source.Digest != nil && *source.Digest != observed {
		return "", &SourceDigestMismatchError{source.DisplayLabel, source.Path, source.Digest, &observed}
	}
	return observed, nil
}
