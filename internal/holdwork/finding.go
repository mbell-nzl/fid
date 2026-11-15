package holdwork

import (
	"fmt"

	"github.com/mbell-nzl/fid/internal/journal"
	"github.com/mbell-nzl/fid/internal/persondecision"
	"github.com/mbell-nzl/fid/internal/reviewrecord"
	"github.com/mbell-nzl/fid/internal/reviewwork"
)

func CloseFinding(o Options) (Result, error) {
	return write(o, "finding.close", func(ctx writeContext) (journal.Event, map[string]any, error) {
		task, err := reviewrecord.FindingTask(ctx.conn, o.FindingID)
		if err != nil {
			return journal.Event{}, nil, err
		}
		c, err := reviewrecord.ReadClosure(ctx.conn, o.FindingID)
		if err != nil {
			return journal.Event{}, nil, err
		}
		if c != nil {
			return journal.Event{}, nil, fmt.Errorf("finding is already closed")
		}
		id, decision, err := persondecision.Effect(ctx.conn, "close-finding", o.Actor, "task", task, o.Reason, ctx.now)
		if err != nil {
			return journal.Event{}, nil, err
		}
		decision.Values["binding_id"] = o.FindingID
		closure := reviewwork.ClosureEffect(o.FindingID, "", id, "closed", o.Reason)
		return journal.Event{EventType: "finding.closed", Aggregate: journal.Aggregate{Type: "finding", ID: o.FindingID}, Effects: []journal.Effect{decision, closure}}, map[string]any{"task_id": task, "finding_id": o.FindingID, "decision_id": id, "disposition": "closed", "reason": o.Reason}, nil
	})
}
