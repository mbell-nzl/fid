// Package reviewcomment renders and delivers canonical review records.
package reviewcomment

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/mbell-nzl/fid/internal/reviewrecord"
)

const RenderVersion = 2

var ErrSealed = errors.New("judge call recorded and sealed")

type Subject struct{ ID, TaskID, HandoffID, PR, Body string }

func ReadSubject(conn reviewrecord.Queryer, id string) (Subject, error) {
	return ReadSubjectVersion(conn, id, RenderVersion)
}

// ReadSubjectVersion uses the receipt's layout so older digests remain verifiable.
func ReadSubjectVersion(conn reviewrecord.Queryer, id string, version int64) (Subject, error) {
	if version != 1 && version != RenderVersion {
		return Subject{}, fmt.Errorf("unknown review render version %d", version)
	}
	var raw, task, handoff string
	var body string
	if strings.HasPrefix(id, "RND-") {
		if err := conn.QueryRow(`SELECT record_json,task_id,handoff_id FROM review_rounds WHERE id=?`, id).Scan(&raw, &task, &handoff); err != nil {
			return Subject{}, err
		}
		var r reviewrecord.Round
		if err := json.Unmarshal([]byte(raw), &r); err != nil {
			return Subject{}, err
		}
		if r.Stream == "judge" {
			v, err := reviewrecord.VisibleRound(conn, r)
			if err != nil {
				return Subject{}, err
			}
			if v.Sealed {
				return Subject{}, ErrSealed
			}
		}
		var number int
		if err := conn.QueryRow(`SELECT count(*) FROM review_rounds WHERE task_id=? AND rowid<=(SELECT rowid FROM review_rounds WHERE id=?)`, task, id).Scan(&number); err != nil {
			return Subject{}, err
		}
		body = RenderRoundV1(r, number)
		if version == 2 {
			body = renderRoundV2(r, number)
		}
		if r.Stream == "judge" {
			v, err := reviewrecord.VisibleRound(conn, r)
			if err != nil {
				return Subject{}, err
			}
			body = renderJudge(v, body)
		}
	} else {
		if err := conn.QueryRow(`SELECT record_json FROM finding_responses WHERE id=?`, id).Scan(&raw); err != nil {
			return Subject{}, err
		}
		var r reviewrecord.Response
		if err := json.Unmarshal([]byte(raw), &r); err != nil {
			return Subject{}, err
		}
		task, handoff = r.TaskID, r.HandoffID
		body = fmt.Sprintf("## Finding answer: %s\n\nBy %s\n\n%s\n", r.FindingID, r.Actor, r.Text)
		if d := r.DeferredTask; d != nil {
			body = fmt.Sprintf("## Deferred finding: %s\n\nBy %s · Discovered task %s: %s\n\n%s\n", r.FindingID, r.Actor, d.ID, d.Title, r.Text)
		}
		if version == 2 {
			body = fmt.Sprintf("## Finding answer: %s\n\n**By** %s\n\n%s\n", codeSpan(r.FindingID), codeSpan(r.Actor), r.Text)
			if d := r.DeferredTask; d != nil {
				body = fmt.Sprintf("## Deferred finding: %s\n\n**By** %s · **Discovered task** %s: %s\n\n%s\n", codeSpan(r.FindingID), codeSpan(r.Actor), codeSpan(d.ID), d.Title, r.Text)
			}
		}
	}
	var h struct {
		PR   string `json:"pr"`
		Lane string `json:"lane"`
	}
	if err := conn.QueryRow(`SELECT record_json FROM handoffs WHERE id=? AND task_id=?`, handoff, task).Scan(&raw); err != nil {
		return Subject{}, err
	}
	if err := json.Unmarshal([]byte(raw), &h); err != nil {
		return Subject{}, err
	}
	if strings.HasPrefix(id, "RND-") {
		if version == 1 {
			body = strings.Replace(body, "\nClaim:", "\nLane: "+h.Lane+"  \nClaim:", 1)
		} else {
			body = strings.Replace(body, "\n**Claim**", "\n**Lane** "+codeSpan(h.Lane)+"  \n**Claim**", 1)
		}
	}
	return Subject{id, task, handoff, h.PR, body}, nil
}

func Digest(body string) string {
	sum := sha256.Sum256([]byte(body))
	return "sha256:" + hex.EncodeToString(sum[:])
}
func Marker(id, body string) string { return "<!-- fid-review:" + id + ":" + Digest(body) + " -->" }

// RenderRoundV1 is retained as the receipt digest format for version 1.
func RenderRoundV1(r reviewrecord.Round, number int) string {
	var b strings.Builder
	fmt.Fprintf(&b, "## Review: %s\n\nReviewer: %s (%s)", r.Verdict, r.Reviewer.Actor, r.Reviewer.Model)
	if r.Reviewer.Effort != "" {
		fmt.Fprintf(&b, ", effort %s", r.Reviewer.Effort)
	}
	fmt.Fprintf(&b, " · Profiles: %s · Round %d  \nSession: %s · Stream: %s · %s  \nClaim: %s\n\n### Findings\n\n", strings.Join(r.Profiles, ", "), number, r.Reviewer.SessionID, r.Stream, r.ID, r.Purpose)
	if len(r.Findings) == 0 {
		b.WriteString("None.\n")
	}
	for _, f := range r.Findings {
		action := map[string]string{"major": "fix before merge", "minor": "fix or answer", "nit": "optional"}[f.Severity]
		fmt.Fprintf(&b, "- **%s · %s · %s** — %s\n", f.ID, f.Severity, action, f.Summary)
		if f.Location != nil {
			fmt.Fprintf(&b, "  Location: %s:%d\n", f.Location.File, f.Location.Line)
		}
		fmt.Fprintf(&b, "  Evidence: %s\n  Fix: %s\n", f.Evidence, f.Fix)
	}
	if len(r.Dispositions) > 0 {
		b.WriteString("\n### Earlier findings\n\n")
		for _, d := range r.Dispositions {
			fmt.Fprintf(&b, "- %s · %s: %s\n", d.FindingID, d.Disposition, d.Reason)
		}
	}
	fmt.Fprintf(&b, "\n### Verification\n\n%s\n", r.Verification)
	return b.String()
}

func renderJudge(v reviewrecord.RoundView, body string) string {
	r, c := v.Round, v.PersonCall
	agreement := "disagree"
	if v.Agrees != nil && *v.Agrees {
		agreement = "agree"
	}

	body = strings.Replace(body, "## Review:", "## Shadow merge judge:", 1)
	body += fmt.Sprintf("\n### Shadow comparison\n\n**Judge's reason** %s\n\n**Person's first call** %s by %s: %s\n\n**Comparison** Calls %s.\n", r.Reason, codeSpan(c.Verdict), codeSpan(c.Actor), c.Reason, agreement)

	if r.LaneRaise != nil {
		body += fmt.Sprintf("\nRecommended lane: %s (shadow only): %s\n", codeSpan(r.LaneRaise.Lane), r.LaneRaise.Reason)
	}
	return body
}
