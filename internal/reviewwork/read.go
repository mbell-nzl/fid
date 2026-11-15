package reviewwork

import (
	"fmt"

	"github.com/mbell-nzl/fid/internal/fingerprint"
	"github.com/mbell-nzl/fid/internal/handoffrecord"
	"github.com/mbell-nzl/fid/internal/reviewrecord"
)

type RoundView = reviewrecord.RoundView
type FindingView = reviewrecord.FindingView

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
	return reviewrecord.CurrentViews(conn, task, fingerprint.Result{Algorithm: fp.Algorithm, Fingerprint: fp.Fingerprint})
}

// CurrentViews is the shared freshness/closure reader for done. The caller
// supplies the latest hand-off's locally computed fingerprint and projection
// snapshot; no git or network is used here.
func CurrentViews(conn reviewrecord.Queryer, task string, fp fingerprint.Result) ([]RoundView, error) {
	return reviewrecord.CurrentViews(conn, task, fp)
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
