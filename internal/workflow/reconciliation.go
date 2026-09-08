package workflow

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"

	"github.com/mbell-nzl/fid/internal/journal"
	"github.com/mbell-nzl/fid/internal/project"
)

type PathChange struct {
	Path           string  `json:"path"`
	BeforeDigest   *string `json:"before_digest"`
	AfterDigest    *string `json:"after_digest"`
	Classification string  `json:"classification"`
}

type Reconciliation struct {
	Repository    Repository   `json:"repository"`
	CheckpointID  *string      `json:"checkpoint_id"`
	Paths         []PathChange `json:"paths"`
	UnmergedPaths []string     `json:"unmerged_paths"`
	Operation     string       `json:"operation"`
	Ancestry      string       `json:"ancestry"`
	Digest        string       `json:"digest"`
}

type RecoveryStateConflictError struct {
	Paths   []string
	Reasons []string
}

func (err *RecoveryStateConflictError) Error() string { return "recovery state conflict" }

// ReconcileRepository converts a read-only Git observation into a bound recovery fact.
func ReconcileRepository(start Repository, current project.Observation, checkpointID *string, approvedPaths []string) (Reconciliation, error) {
	return reconcileRepository(start, current, checkpointID, approvedPaths, false)
}

// ReconcileImportedDestination applies the ordinary reconciliation checks after
// a verified local import has bound the historical project to this worktree.
func ReconcileImportedDestination(start Repository, current project.Observation, checkpointID *string, approvedPaths []string) (Reconciliation, error) {
	return reconcileRepository(start, current, checkpointID, approvedPaths, true)
}

func reconcileRepository(start Repository, current project.Observation, checkpointID *string, approvedPaths []string, allowDestinationWorktree bool) (Reconciliation, error) {
	if start.BaseOID == nil || *start.BaseOID == "" || start.HeadOID == "" || start.WorktreeID == "" || current.Root == "" || current.WorktreeID == "" {
		return Reconciliation{}, recoveryConflict(nil, "missing repository base or identity")
	}
	if start.WorktreeID != current.WorktreeID && !allowDestinationWorktree {
		return Reconciliation{}, recoveryConflict(nil, "repository identity changed")
	}
	if exists, err := project.IsAncestor(current.Root, *start.BaseOID, *start.BaseOID); err != nil || !exists {
		return Reconciliation{}, recoveryConflict(nil, "repository base is unavailable")
	}
	if len(current.Unmerged) != 0 {
		return Reconciliation{}, recoveryConflict(current.Unmerged, "repository has unmerged paths")
	}
	if current.Operation != "none" {
		return Reconciliation{}, recoveryConflict(nil, "repository operation is in progress")
	}
	ancestry := "unchanged"
	if current.HeadOID != start.HeadOID {
		ancestor, err := project.IsAncestor(current.Root, start.HeadOID, current.HeadOID)
		if err != nil || !ancestor {
			return Reconciliation{}, recoveryConflict(nil, "repository history is unattributable")
		}
		ancestry = "descendant"
	}
	approved := make(map[string]bool, len(approvedPaths))
	for _, value := range approvedPaths {
		approved[value] = true
	}
	paths := make([]PathChange, 0, len(current.Entries))
	for _, entry := range current.Entries {
		classification := "outside-scope"
		if entry.Path == ".fid" || strings.HasPrefix(entry.Path, ".fid/") || entry.Path == ".git" || strings.HasPrefix(entry.Path, ".git/") {
			classification = "protected"
		} else if approved[entry.Path] {
			classification = "in-scope"
		}
		paths = append(paths, PathChange{Path: entry.Path, BeforeDigest: imageDigest(entry.Head), AfterDigest: changedImageDigest(entry), Classification: classification})
	}
	sort.Slice(paths, func(i, j int) bool { return paths[i].Path < paths[j].Path })
	repository := start
	if allowDestinationWorktree {
		repository.WorktreeID = current.WorktreeID
	}
	repository.HeadOID, repository.TreeOID, repository.Dirty, repository.DirtyDigest = current.HeadOID, current.TreeOID, current.Dirty, current.DirtyDigest
	reconciliation := Reconciliation{Repository: repository, CheckpointID: checkpointID, Paths: paths, UnmergedPaths: []string{}, Operation: current.Operation, Ancestry: ancestry}
	payload, err := journal.MarshalCanonicalJSON(map[string]any{"repository": reconciliation.Repository, "checkpoint_id": reconciliation.CheckpointID, "paths": reconciliation.Paths, "unmerged_paths": reconciliation.UnmergedPaths, "operation": reconciliation.Operation, "ancestry": reconciliation.Ancestry})
	if err != nil {
		return Reconciliation{}, fmt.Errorf("canonical reconciliation digest: %w", err)
	}
	sum := sha256.Sum256(payload)
	reconciliation.Digest = "sha256:" + hex.EncodeToString(sum[:])
	return reconciliation, nil
}

// ValidateReconciliationAdoption proves the exact current observed dirty set was adopted.
func ValidateReconciliationAdoption(reconciliation Reconciliation, adoptedPaths []string, digest string) error {
	if reconciliation.Digest == "" || digest != reconciliation.Digest {
		return recoveryConflict(nil, "reconciliation digest is stale or changed")
	}
	if reconciliation.Operation != "none" || len(reconciliation.UnmergedPaths) != 0 || reconciliation.Ancestry == "conflict" || reconciliation.Ancestry == "unknown" {
		return recoveryConflict(reconciliation.UnmergedPaths, "reconciliation is unsafe")
	}
	want := make([]string, 0, len(reconciliation.Paths))
	for _, change := range reconciliation.Paths {
		if change.Classification != "in-scope" {
			return recoveryConflict([]string{change.Path}, "repository contains non-adoptable change")
		}
		want = append(want, change.Path)
	}
	sort.Strings(want)
	got := append([]string(nil), adoptedPaths...)
	sort.Strings(got)
	if len(want) != len(got) {
		return recoveryConflict(want, "adopted paths are incomplete")
	}
	for index := range want {
		if want[index] != got[index] {
			return recoveryConflict(want, "adopted paths differ from the observation")
		}
	}
	return nil
}

func imageDigest(image *project.Image) *string {
	if image == nil {
		return nil
	}
	digest := image.Digest
	return &digest
}

func changedImageDigest(entry project.Entry) *string {
	if entry.Index != nil && (entry.Head == nil || entry.Index.Mode != entry.Head.Mode || entry.Index.Digest != entry.Head.Digest) {
		return imageDigest(entry.Index)
	}
	return imageDigest(entry.Worktree)
}

func recoveryConflict(paths []string, reason string) error {
	return &RecoveryStateConflictError{Paths: append([]string(nil), paths...), Reasons: []string{reason}}
}
