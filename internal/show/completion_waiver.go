package show

import (
	"fmt"

	"github.com/mbell-nzl/fid/internal/checkread"
	"github.com/mbell-nzl/fid/internal/checkrecord"
	"github.com/mbell-nzl/fid/internal/decisionview"
	"github.com/mbell-nzl/fid/internal/fingerprint"
	"github.com/mbell-nzl/fid/internal/policy"
	"github.com/mbell-nzl/fid/internal/reviewrecord"
)

func (result *Result) readCompletionWaiver(root string, conn reviewrecord.Queryer) error {
	c := result.Completion
	fp := fingerprint.Result{Algorithm: c.Fingerprint.Algorithm, Fingerprint: c.Fingerprint.Fingerprint}
	waivers, err := decisionview.ReadRecordedWaivers(conn, c.TaskID, fp)
	if err != nil {
		return err
	}
	result.Approval.Waivers = waivers
	if w := decisionview.LatestWaiver(waivers); w != nil && w.Current {
		c.Waiver = w
		result.Approval.Decision = &decisionview.View{State: "not-required"}
	}
	return c.readCheckRequirements(root, conn, result.IssueView.Checks)
}

func (c *Completion) readCheckRequirements(root string, conn reviewrecord.Queryer, records []checkrecord.Record) error {
	h := c.Handoff
	if h == nil {
		return nil
	}
	c.requiredChecks, c.resultIDs = h.RequiredChecks, map[string]string{}
	for i, name := range h.RequiredChecks {
		if i < len(h.ResultIDs) {
			c.resultIDs[name] = h.ResultIDs[i]
		}
	}
	if c.Waiver == nil {
		return nil
	}
	c.requiredChecks = nil
	p, err := completionRow[policy.Effective](conn, c, "policy snapshot", `SELECT effective_json FROM policy_snapshots WHERE sequence=?`, fmt.Sprint(h.PolicySnapshot))
	if err != nil || p == nil {
		return err
	}
	c.requiredChecks, err = policy.Required(*p, c.Waiver.Binding.Lane, c.TaskID, h.Paths())
	if err != nil {
		return fmt.Errorf("read completed waived checks: %w", err)
	}
	// Keep the hand-off's named results. A lower lane can also require a check
	// absent from the higher lane; select that pass on the recorded tree and
	// definitions, just as done did. Only results recorded before completion
	// were available to done, even if a later finish reports an earlier time.
	records, err = c.checksBeforeCompletion(root, records)
	if err != nil {
		return err
	}
	checks := checkread.Evaluate(c.TaskID, h.Tree, c.requiredChecks, p.Checks, records)
	for name, id := range checks.LatestResultIDs {
		if _, named := c.resultIDs[name]; !named {
			c.resultIDs[name] = id
		}
	}
	return nil
}
