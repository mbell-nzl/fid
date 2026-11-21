package doneeval

import "github.com/mbell-nzl/fid/internal/reviewrecord"

type ReworkLimit struct {
	Count int
	Cap   int64
}

func (r ReworkLimit) Reached() bool { return int64(r.Count) > r.Cap }

// ReadReworkLimit counts both streams over the task's whole history, against
// the policy and raised lane frozen at its latest hand-off. No git is needed.
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
	for _, round := range rounds {
		if round.Verdict == "CHANGES" {
			r.Count++
		}
	}
	return r, nil
}
