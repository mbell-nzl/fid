package decisionview

import (
	"github.com/mbell-nzl/fid/internal/approvalrecord"
	"github.com/mbell-nzl/fid/internal/fingerprint"
	"github.com/mbell-nzl/fid/internal/handoffrecord"
	"github.com/mbell-nzl/fid/internal/reviewrecord"
)

type Approval struct {
	Waivers            []Waiver `json:"waivers,omitempty"`
	Decision           *View    `json:"approval"`
	NoOtherModelReview *bool    `json:"no_other_model_review"`
}

// ReadApproval uses the latest hand-off's lane and recorded endpoints, never
// today's policy or merge-base. A missing fingerprint leaves unknown fields
// null; an absence of approvals or reviews is still known without git objects.
// The caller owns the projection snapshot and refreshes it before reading.
func ReadApproval(root string, conn reviewrecord.Queryer, h handoffrecord.Record) (Approval, error) {
	waivers, err := ReadWaivers(root, conn, h)
	result := Approval{Waivers: waivers}
	if err != nil {
		return result, err
	}
	if waiver := LatestWaiver(waivers); waiver != nil && waiver.Current {
		result.Decision = &View{State: "not-required"}
		return result, nil
	}
	if h.Lane != "assured" {
		return result, nil
	}
	approvals, err := approvalrecord.Read(conn, h.TaskID)
	if err != nil {
		return result, err
	}
	rounds, err := reviewrecord.Read(conn, h.TaskID)
	if err != nil {
		return result, err
	}
	if len(approvals) == 0 {
		result.Decision = &View{State: "missing"}
	}
	if len(rounds) == 0 {
		flag := true
		result.NoOtherModelReview = &flag
	}
	fp := fingerprint.Result{}
	if len(approvals) > 0 || len(rounds) > 0 {
		fp, err = fingerprint.Compute(root, h.Base, h.Head)
		if err != nil {
			return result, err
		}
	}
	view, err := approvalrecord.Current(conn, h.TaskID, fp)
	if err != nil {
		return result, err
	}
	result.Decision = &View{State: view.State}
	if d := view.Decision; d != nil {
		result.Decision = &View{State: view.State, DecisionID: d.ID, Actor: d.Actor, At: d.At, Fingerprint: d.Binding.Fingerprint}
	}
	views, err := reviewrecord.CurrentViews(conn, h.TaskID, fp)
	if err != nil {
		return result, err
	}
	flag := reviewrecord.NoOtherModelReview(h.TaskID, views)
	result.NoOtherModelReview = &flag
	return result, nil
}
