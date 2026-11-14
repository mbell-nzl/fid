// Package checkrecord defines the lean check's immutable observations, shared
// by execution, replay and archive validation. It has no persistence dependency.
package checkrecord

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"regexp"
	"time"
)

const (
	OutputLimit   = 10 << 20
	ExcerptLimit  = 16 << 10
	ArtifactLimit = 64 << 20
)

type Definition struct {
	Run         []string      `json:"run"`
	Timeout     time.Duration `json:"timeout_ns"`
	Environment []string      `json:"environment"`
	Artifacts   []string      `json:"artifacts"`
}

func Digest(value any) string {
	raw, _ := json.Marshal(value)
	return BytesDigest(raw)
}

func BytesDigest(raw []byte) string { return fmt.Sprintf("sha256:%x", sha256.Sum256(raw)) }

type Tree struct {
	Head        string  `json:"head"`
	Tree        string  `json:"tree"`
	Dirty       bool    `json:"dirty"`
	DirtyDigest *string `json:"dirty_digest"`
}

type Process struct {
	Machine   string `json:"machine"`
	BootID    string `json:"boot_id"`
	PID       int    `json:"pid"`
	StartTime string `json:"start_time"`
}

var bootIDPattern = regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)

func (process Process) Validate() error {
	if process.PID <= 0 || process.Machine == "" || process.StartTime == "" {
		return fmt.Errorf("incomplete check process identity")
	}
	// Other Unix systems use ps start times and have no Linux boot ID.
	if process.BootID != "" && !bootIDPattern.MatchString(process.BootID) {
		return fmt.Errorf("invalid check boot identity")
	}
	return nil
}

type Start struct {
	ID                string     `json:"id"`
	TaskID            string     `json:"task_id"`
	HoldID            string     `json:"hold_id"`
	Actor             string     `json:"actor"`
	Name              string     `json:"name"`
	Definition        Definition `json:"definition"`
	DefinitionDigest  string     `json:"definition_digest"`
	PolicySnapshot    int        `json:"policy_snapshot"`
	EnvironmentDigest string     `json:"environment_digest"`
	StartedAt         string     `json:"started_at"`
	Tree              Tree       `json:"tree"`
	Owner             Process    `json:"owner"`
}

type Toolchain struct {
	Executable string `json:"executable"`
	Digest     string `json:"digest"`
	OS         string `json:"os"`
	Arch       string `json:"arch"`
	FidVersion string `json:"fid_version"`
}

type Artifact struct {
	ID        string `json:"id"`
	Path      string `json:"path"`
	Digest    string `json:"digest"`
	SizeBytes int64  `json:"size_bytes"`
}

type RefusedArtifact struct {
	Path   string `json:"path"`
	Reason string `json:"reason"`
}

type Result struct {
	FinishedAt       string            `json:"finished_at"`
	Outcome          string            `json:"outcome"`
	ExitCode         *int              `json:"exit_code"`
	Signal           *string           `json:"signal"`
	Tree             Tree              `json:"tree"`
	OutputDigest     string            `json:"output_digest"`
	OutputExcerpt    string            `json:"output_excerpt"`
	OutputBytes      int64             `json:"output_bytes"`
	OutputTruncated  bool              `json:"output_truncated"`
	Toolchain        Toolchain         `json:"toolchain"`
	Artifacts        []Artifact        `json:"artifacts"`
	RefusedArtifacts []RefusedArtifact `json:"refused_artifacts"`
	Counts           bool              `json:"counts"`
	Reasons          []string          `json:"reasons"`
}

type Record struct {
	Start   Start    `json:"start"`
	Process *Process `json:"process"`
	Result  *Result  `json:"result"`
}

// Qualification is deliberately independent of the policy's required set.
func Qualification(start, end Tree, outcome string, current, allKept bool) []string {
	reasons := []string{}
	if start.Dirty {
		reasons = append(reasons, "tree-dirty")
	}
	start.Head, end.Head = "", ""
	if Digest(start) != Digest(end) {
		reasons = append(reasons, "tree-changed")
	}
	if !current {
		reasons = append(reasons, "ownership-lost")
	}
	if !allKept {
		reasons = append(reasons, "artifact-not-kept")
	}
	return reasons
}

func EnvironmentNames() []string {
	return []string{"CI", "GOCACHE", "GOMODCACHE", "GOPATH", "HOME", "LANG", "LC_ALL", "PATH", "TMPDIR", "XDG_CACHE_HOME"}
}
