package doneeval

import (
	"github.com/mbell-nzl/fid/internal/reviewrecord"
	"github.com/mbell-nzl/fid/internal/workrecord"
)

type ReworkLimit struct {
	Count     int
	Cap       int64
	Allowance *workrecord.Allowance
}

func (r ReworkLimit) Reached() bool { return int64(r.Count) > r.Cap }

// ReadReworkLimit counts both streams over the task's whole history, against
// the policy and raised lane frozen at its latest hand-off, plus every rework
// allowance. No git is needed.
func ReadReworkLimit(conn reviewrecord.Queryer, task string) (*ReworkLimit, error) {
	s, err := reviewBase(conn, task)
	if err != nil || s.Handoff == nil || s.Policy == nil || s.Task.State != "in_progress" {
		return nil, err
	}
	rounds, err := reviewrecord.Read(conn, task)
	if err != nil {
		return nil, err
	}
	r := &ReworkLimit{Cap: s.Policy.Lanes[s.Handoff.Lane].ReworkCap}
	r.Cap, r.Allowance, err = workrecord.Allowances(conn, task, r.Cap)
	if err != nil {
		return nil, err
	}
	for _, round := range rounds {
		if round.Verdict == "CHANGES" {
			r.Count++
		}
	}
	return r, nil
}
