package approvalwork

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/mbell-nzl/fid/internal/approvalrecord"
	"github.com/mbell-nzl/fid/internal/checkrecord"
	"github.com/mbell-nzl/fid/internal/cli"
	"github.com/mbell-nzl/fid/internal/fingerprint"
	"github.com/mbell-nzl/fid/internal/handoffrecord"
	"github.com/mbell-nzl/fid/internal/reviewrecord"
)

type Preview struct {
	TaskID       string
	Handoff      handoffrecord.Record
	Fingerprint  fingerprint.Result
	Diffstat     string
	Checks       []string
	Streams      map[string]string
	Majors       []reviewrecord.Finding
	Earlier      []approvalrecord.Record
	NoOtherModel bool
}

func previewFor(conn reviewrecord.Queryer, task string, h handoffrecord.Record, fp fingerprint.Result, stat string) (Preview, error) {
	p := Preview{TaskID: task, Handoff: h, Fingerprint: fp, Diffstat: stat, Checks: []string{}, Streams: map[string]string{}, NoOtherModel: true}
	for i, name := range h.RequiredChecks {
		var raw string
		if err := conn.QueryRow(`SELECT result_json FROM check_results WHERE check_id=?`, h.ResultIDs[i]).Scan(&raw); err != nil {
			return p, err
		}
		var result checkrecord.Result
		if err := json.Unmarshal([]byte(raw), &result); err != nil {
			return p, err
		}
		p.Checks = append(p.Checks, fmt.Sprintf("%s: %s · counts=%t · %s", name, result.Outcome, result.Counts, h.ResultIDs[i]))
	}
	rounds, err := reviewrecord.Read(conn, task)
	if err != nil {
		return p, err
	}
	for _, r := range rounds {
		current := r.Binding.Algorithm == fp.Algorithm && r.Binding.Fingerprint == fp.Fingerprint
		freshness := "stale"
		if current {
			freshness = "current"
			other, err := reviewrecord.OtherModel(conn, task, r.Reviewer.Model)
			if err != nil {
				return p, err
			}
			p.NoOtherModel = p.NoOtherModel && !other
		}
		model := r.Reviewer.Model
		if model == "" {
			model = "unknown model"
		}
		p.Streams[r.Stream] = fmt.Sprintf("%s · %s · %s · %s · %s", r.Verdict, r.Reviewer.Actor, model, freshness, r.ID)
	}
	p.Majors, err = reviewrecord.OpenMajors(conn, task)
	if err != nil {
		return p, err
	}
	p.Earlier, err = approvalrecord.Read(conn, task)
	return p, err
}

// WritePrompt writes all evidence before the confirmation, escaping ledger text
// so a reason, reviewer or path cannot control the person's terminal.
func WritePrompt(w io.Writer, p Preview) error {
	var b strings.Builder
	line := func(format string, args ...any) { fmt.Fprintf(&b, format+"\n", args...) }
	esc := cli.EscapeControlText
	line("Approve %s · %s", esc(p.TaskID), esc(p.Handoff.Lane))
	line("Hand-off %s · %s..%s", esc(p.Handoff.ID), p.Handoff.Base, p.Handoff.Head)
	line("Fingerprint: %s · %s", p.Fingerprint.Algorithm, p.Fingerprint.Fingerprint)
	line("Diffstat:")
	for _, s := range strings.Split(strings.TrimSuffix(p.Diffstat, "\n"), "\n") {
		line("%s", esc(s))
	}
	if len(p.Checks) == 0 {
		line("Checks: none required")
	} else {
		line("Checks:")
		for _, s := range p.Checks {
			line("  %s", esc(s))
		}
	}
	line("Reviews:")
	streams := []string{"primary"}
	if p.Handoff.Lane == "assured" {
		streams = append(streams, "separate")
	}
	for _, stream := range streams {
		value := p.Streams[stream]
		if value == "" {
			value = "missing"
		}
		line("  %s: %s", stream, esc(value))
	}
	if len(p.Majors) == 0 {
		line("Open majors: none")
	} else {
		line("Open majors:")
		for _, f := range p.Majors {
			line("  %s: %s", esc(f.ID), esc(f.Summary))
		}
	}
	if p.NoOtherModel {
		line("no other-model review")
	}
	if len(p.Earlier) == 0 {
		line("Earlier approvals: none")
	} else {
		line("Earlier approvals:")
		for _, r := range p.Earlier {
			state := "stale"
			if r.Binding.Algorithm == p.Fingerprint.Algorithm && r.Binding.Fingerprint == p.Fingerprint.Fingerprint {
				state = "current"
			}
			line("  %s · %s · %s · %s · %s", esc(r.ID), esc(r.Actor), state, esc(r.At), esc(r.Reason))
		}
	}
	line("Approval records your decision; done still requires every gate.")
	fmt.Fprintf(&b, "Type %s to approve: ", esc(p.TaskID))
	_, err := io.WriteString(w, b.String())
	return err
}
