package decisionview

import (
	"github.com/mbell-nzl/fid/internal/fingerprint"
	"github.com/mbell-nzl/fid/internal/handoffrecord"
	"github.com/mbell-nzl/fid/internal/reviewrecord"
	"github.com/mbell-nzl/fid/internal/waiverrecord"
)

type Waiver struct {
	waiverrecord.Record
	Current bool `json:"current"`
}

func (w Waiver) View() View {
	state := "stale"
	if w.Current {
		state = "current"
	}
	return View{State: state, DecisionID: w.ID, Actor: w.Actor, At: w.At, Fingerprint: w.Binding.Fingerprint}
}

func ReadWaivers(root string, conn reviewrecord.Queryer, h handoffrecord.Record) ([]Waiver, error) {
	records, err := waiverrecord.Read(conn, h.TaskID)
	if err != nil || len(records) == 0 {
		return nil, err
	}
	fp, err := fingerprint.Compute(root, h.Base, h.Head)
	if err != nil {
		return nil, err
	}
	return waiverViews(records, fp), nil
}

// ReadRecordedWaivers presents currency against an immutable completion's
// fingerprint, without needing its git objects or today's comparison.
func ReadRecordedWaivers(conn reviewrecord.Queryer, task string, fp fingerprint.Result) ([]Waiver, error) {
	records, err := waiverrecord.Read(conn, task)
	if err != nil || len(records) == 0 {
		return nil, err
	}
	return waiverViews(records, fp), nil
}

func waiverViews(records []waiverrecord.Record, fp fingerprint.Result) []Waiver {
	result := make([]Waiver, 0, len(records))
	for _, r := range records {
		result = append(result, Waiver{Record: r, Current: r.Matches(fp)})
	}
	return result
}

// LatestWaiver prefers the most recent current decision, as done does.
func LatestWaiver(records []Waiver) *Waiver {
	for i := len(records) - 1; i >= 0; i-- {
		if records[i].Current {
			return &records[i]
		}
	}
	if len(records) > 0 {
		return &records[len(records)-1]
	}
	return nil
}
