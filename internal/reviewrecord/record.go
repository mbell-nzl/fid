// Package reviewrecord holds immutable, attributed review rounds and findings.
// It validates recorded evidence without consulting git or the network.
package reviewrecord

import (
	"fmt"
	"regexp"
	"slices"
	"strings"
	"time"

	"github.com/mbell-nzl/fid/internal/fingerprint"
	"github.com/mbell-nzl/fid/internal/modelid"
	"github.com/mbell-nzl/fid/internal/policy"
)

type Reviewer struct {
	Actor     string `json:"actor"`
	Model     string `json:"model,omitempty"`
	Effort    string `json:"effort,omitempty"`
	SessionID string `json:"session_id"`
}

type Location struct {
	File string `json:"file"`
	Line int    `json:"line"`
}

type FindingInput struct {
	Profile  string    `json:"profile"`
	Severity string    `json:"severity"`
	Summary  string    `json:"summary"`
	Location *Location `json:"location"`
	Evidence string    `json:"evidence"`
	Fix      string    `json:"fix"`
}

type Provenance struct {
	Prompt  string `json:"prompt,omitempty"`
	Command string `json:"command,omitempty"`
	Output  string `json:"output,omitempty"`
}

type Disposition struct {
	FindingID   string `json:"finding_id"`
	Disposition string `json:"disposition"`
	Reason      string `json:"reason"`
}

func (d Disposition) Closes() bool {
	return d.Disposition == "resolved" || d.Disposition == "not-applicable"
}

type LaneRaise struct {
	Lane   string `json:"lane"`
	Reason string `json:"reason"`
}

type Input struct {
	Reason       string         `json:"reason,omitempty"`
	LaneRaise    *LaneRaise     `json:"lane_raise,omitempty"`
	Reviewer     Reviewer       `json:"reviewer"`
	Stream       string         `json:"stream"`
	Profiles     []string       `json:"profiles"`
	Verdict      string         `json:"verdict"`
	Purpose      string         `json:"purpose"`
	Findings     []FindingInput `json:"findings"`
	Dispositions []Disposition  `json:"dispositions,omitempty"`
	Verification string         `json:"verification"`
	Provenance   *Provenance    `json:"provenance,omitempty"`
}

type Binding struct {
	Algorithm   string `json:"algorithm"`
	Fingerprint string `json:"fingerprint"`
	Base        string `json:"base"`
	Head        string `json:"head"`
}

func Bind(f fingerprint.Result) Binding { return Binding{f.Algorithm, f.Fingerprint, f.Base, f.Head} }

type Finding struct {
	ID string `json:"id"`
	FindingInput
}

type Round struct {
	Reason       string        `json:"reason,omitempty"`
	LaneRaise    *LaneRaise    `json:"lane_raise,omitempty"`
	ID           string        `json:"id"`
	TaskID       string        `json:"task_id"`
	HandoffID    string        `json:"handoff_id"`
	Reviewer     Reviewer      `json:"reviewer"`
	Stream       string        `json:"stream"`
	Profiles     []string      `json:"profiles"`
	Verdict      string        `json:"verdict"`
	Purpose      string        `json:"purpose"`
	Findings     []Finding     `json:"findings"`
	Dispositions []Disposition `json:"dispositions,omitempty"`
	Verification string        `json:"verification"`
	Provenance   *Provenance   `json:"provenance,omitempty"`
	Binding      Binding       `json:"binding"`
	RecordedBy   string        `json:"recorded_by"`
	At           string        `json:"at"`
}

var token = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._-]{0,127}$`)
var actor = regexp.MustCompile(policy.ActorPattern)
var hexID = regexp.MustCompile(`^[0-9a-f]{40}$`)

// fingerprintHex is the form of fingerprint.Algorithm's value: a SHA-256 over
// the per-file patch-ids.
var fingerprintHex = regexp.MustCompile(`^[0-9a-f]{64}$`)

func (in Input) Validate() error {
	if !actor.MatchString(in.Reviewer.Actor) || !policy.Allows(in.Reviewer.Actor, policy.WorkSet) || !token.MatchString(in.Reviewer.SessionID) {
		return fmt.Errorf("reviewer requires a work actor and session ID")
	}
	if in.Reviewer.Model != "" && !modelid.Valid(in.Reviewer.Model) {
		return fmt.Errorf("reviewer model must be <family>/<model>")
	}
	if in.Reviewer.Effort != "" && !token.MatchString(in.Reviewer.Effort) {
		return fmt.Errorf("invalid reviewer effort")
	}
	if in.Stream == "judge" {
		if (in.Verdict != "MERGE" && in.Verdict != "HOLD") || !slices.Equal(in.Profiles, []string{"merge"}) || !text(in.Reason) || len(in.Findings) != 0 || len(in.Dispositions) != 0 {
			return fmt.Errorf("judge requires merge profile, MERGE or HOLD, reason, and no findings or dispositions")
		}
		if in.LaneRaise != nil && (!text(in.LaneRaise.Reason) || (in.LaneRaise.Lane != "standard" && in.LaneRaise.Lane != "assured")) {
			return fmt.Errorf("judge lane raise requires a higher lane and reason")
		}
	} else if in.Reason != "" || in.LaneRaise != nil {
		return fmt.Errorf("judge fields require judge stream")
	}
	if in.Stream != "primary" && in.Stream != "separate" && in.Stream != "judge" {
		return fmt.Errorf("review stream must be primary, separate or judge")
	}
	if in.Stream != "judge" && in.Verdict != "APPROVE" && in.Verdict != "CHANGES" {
		return fmt.Errorf("review verdict must be APPROVE or CHANGES")
	}
	if in.Profiles == nil || in.Findings == nil || !text(in.Purpose) || !text(in.Verification) {
		return fmt.Errorf("review requires profiles, findings, purpose and verification")
	}
	seen := map[string]bool{}
	for _, p := range in.Profiles {
		if !token.MatchString(p) || seen[p] {
			return fmt.Errorf("invalid or duplicate review profile")
		}
		seen[p] = true
	}
	for _, f := range in.Findings {
		if !seen[f.Profile] || (f.Severity != "major" && f.Severity != "minor" && f.Severity != "nit") || !text(f.Summary) || !text(f.Evidence) || !text(f.Fix) {
			return fmt.Errorf("finding requires a reviewed profile, severity, summary, evidence and fix")
		}
		if f.Location != nil && (!text(f.Location.File) || f.Location.Line < 1) {
			return fmt.Errorf("finding location requires a file and positive line")
		}
	}
	seenFindings := map[string]bool{}
	for _, d := range in.Dispositions {
		if d.FindingID == "" || seenFindings[d.FindingID] || (d.Disposition != "open" && !d.Closes()) || !text(d.Reason) {
			return fmt.Errorf("disposition requires a unique finding, resolved, not-applicable or open, and a reason")
		}
		seenFindings[d.FindingID] = true
	}
	return nil
}

func text(s string) bool { return strings.TrimSpace(s) != "" && len(s) <= 1<<20 }

func (r Round) Input() Input {
	findings := make([]FindingInput, len(r.Findings))
	for i, f := range r.Findings {
		findings[i] = f.FindingInput
	}
	return Input{Reviewer: r.Reviewer, Stream: r.Stream, Profiles: r.Profiles, Verdict: r.Verdict, Purpose: r.Purpose, Findings: findings, Dispositions: r.Dispositions, Verification: r.Verification, Provenance: r.Provenance, Reason: r.Reason, LaneRaise: r.LaneRaise}
}

func (r Round) Validate() error {
	if !strings.HasPrefix(r.ID, "RND-") || r.TaskID == "" || r.HandoffID == "" || !actor.MatchString(r.RecordedBy) || !policy.Allows(r.RecordedBy, policy.WorkSet) || r.Findings == nil {
		return fmt.Errorf("incomplete review round identity")
	}
	if err := r.Input().Validate(); err != nil {
		return err
	}
	if at, err := time.Parse(time.RFC3339, r.At); err != nil || at.IsZero() || !strings.HasSuffix(r.At, "Z") {
		return fmt.Errorf("invalid review time")
	}
	b := r.Binding
	if b.Algorithm != fingerprint.Algorithm || !fingerprintHex.MatchString(b.Fingerprint) || !hexID.MatchString(b.Base) || !hexID.MatchString(b.Head) {
		return fmt.Errorf("invalid review fingerprint")
	}
	ids := []string{}
	for _, f := range r.Findings {
		if f.ID == "" || slices.Contains(ids, f.ID) {
			return fmt.Errorf("missing or duplicate finding ID")
		}
		ids = append(ids, f.ID)
	}
	return nil
}
