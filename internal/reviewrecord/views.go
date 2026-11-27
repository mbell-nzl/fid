package reviewrecord

import "github.com/mbell-nzl/fid/internal/fingerprint"

type RoundView struct {
	Round
	Sealed     bool          `json:"sealed,omitempty"`
	PersonCall *PersonCall   `json:"person_call,omitempty"`
	Agrees     *bool         `json:"agrees,omitempty"`
	OtherModel bool          `json:"other_model"`
	Current    bool          `json:"current"`
	Findings   []FindingView `json:"findings"`
}

type FindingView struct {
	Finding
	Open      bool       `json:"open"`
	Responses []Response `json:"responses,omitempty"`
	Closure   *Closure   `json:"closure,omitempty"`
}

// NoOtherModelReview is the shared Assured disclosure for the approve prompt,
// done and read views. Every current primary or separate round can establish
// coverage; stale rounds and rounds for another task cannot.
func NoOtherModelReview(task string, rounds []RoundView) bool {
	for _, r := range rounds {
		if r.TaskID == task && r.Current && (r.Stream == "primary" || r.Stream == "separate") && r.OtherModel {
			return false
		}
	}
	return true
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
			responses, err := Responses(conn, f.ID)
			if err != nil {
				return nil, err
			}
			v.Findings = append(v.Findings, FindingView{Finding: f, Open: c == nil, Closure: c, Responses: responses})
		}
		if r.Stream == "judge" {
			visible, err := VisibleRound(conn, r)
			if err != nil {
				return nil, err
			}
			visible.Current = v.Current
			v = visible
		}
		views = append(views, v)
	}
	return views, nil
}
