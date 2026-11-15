package reviewwork

import (
	"fmt"

	"github.com/mbell-nzl/fid/internal/fingerprint"
	"github.com/mbell-nzl/fid/internal/handoffrecord"
	"github.com/mbell-nzl/fid/internal/reviewrecord"
)

type RoundView struct {
	reviewrecord.Round
	OtherModel bool          `json:"other_model"`
	Current    bool          `json:"current"`
	Findings   []FindingView `json:"findings"`
}

type FindingView struct {
	reviewrecord.Finding
	Open    bool                  `json:"open"`
	Closure *reviewrecord.Closure `json:"closure,omitempty"`
}

// List derives freshness from the latest hand-off, without rewriting history.
// The caller refreshes the projection before reading. All git is local.
func List(root, task string) ([]RoundView, error) {
	conn, err := OpenProjection(root)
	if err != nil {
		return nil, err
	}
	defer func() { _ = conn.Close() }()
	var kind string
	if err := conn.QueryRow(`SELECT kind FROM issues WHERE id=?`, task).Scan(&kind); err != nil {
		return nil, err
	}
	if kind != "task" {
		return nil, fmt.Errorf("review list requires a task")
	}
	rounds, err := reviewrecord.Read(conn, task)
	if err != nil {
		return nil, err
	}
	views := []RoundView{}
	if len(rounds) == 0 {
		return views, nil
	}
	h, err := reviewrecord.LatestHandoff(conn, task)
	if err != nil {
		return nil, err
	}
	fp, err := latestBinding(root, h, rounds)
	if err != nil {
		return nil, err
	}
	for _, r := range rounds {
		other, err := reviewrecord.OtherModel(conn, r.TaskID, r.Reviewer.Model)
		if err != nil {
			return nil, err
		}
		v := RoundView{Round: r, OtherModel: other, Current: r.Binding.Algorithm == fp.Algorithm && r.Binding.Fingerprint == fp.Fingerprint, Findings: []FindingView{}}
		for _, f := range r.Findings {
			c, err := reviewrecord.ReadClosure(conn, f.ID)
			if err != nil {
				return nil, err
			}
			v.Findings = append(v.Findings, FindingView{f, c == nil, c})
		}
		views = append(views, v)
	}
	return views, nil
}

// A recorded round already computed this exact pair of endpoints. Reuse that
// evidence for restored ledgers whose repository objects are not present yet.
// A later hand-off with new endpoints still needs a local computation.
func latestBinding(root string, h handoffrecord.Record, rounds []reviewrecord.Round) (reviewrecord.Binding, error) {
	for _, r := range rounds {
		if r.Binding.Base == h.Base && r.Binding.Head == h.Head {
			return r.Binding, nil
		}
	}
	result, err := fingerprint.Compute(root, h.Base, h.Head)
	return reviewrecord.Bind(result), err
}
