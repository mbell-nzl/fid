package reviewrecord

import "github.com/mbell-nzl/fid/internal/fingerprint"

type RoundView struct {
	Round
	OtherModel bool          `json:"other_model"`
	Current    bool          `json:"current"`
	Findings   []FindingView `json:"findings"`
}

type FindingView struct {
	Finding
	Open    bool     `json:"open"`
	Closure *Closure `json:"closure,omitempty"`
}

func CurrentViews(conn Queryer, task string, fp fingerprint.Result) ([]RoundView, error) {
	rounds, err := Read(conn, task)
	if err != nil {
		return nil, err
	}
	return currentViews(conn, rounds, Bind(fp))
}

func currentViews(conn Queryer, rounds []Round, fp Binding) ([]RoundView, error) {
	views := []RoundView{}
	for _, r := range rounds {
		other, err := OtherModel(conn, r.TaskID, r.Reviewer.Model)
		if err != nil {
			return nil, err
		}
		v := RoundView{Round: r, OtherModel: other, Current: r.Binding.Algorithm == fp.Algorithm && r.Binding.Fingerprint == fp.Fingerprint, Findings: []FindingView{}}
		for _, f := range r.Findings {
			c, err := ReadClosure(conn, f.ID)
			if err != nil {
				return nil, err
			}
			v.Findings = append(v.Findings, FindingView{f, c == nil, c})
		}
		views = append(views, v)
	}
	return views, nil
}
