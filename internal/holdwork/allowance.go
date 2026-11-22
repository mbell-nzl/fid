package holdwork

import (
	"encoding/json"
	"fmt"

	"github.com/mbell-nzl/fid/internal/journal"
	"github.com/mbell-nzl/fid/internal/persondecision"
)

func AllowRework(options Options) (Result, error) {
	return write(options, "task.rework-allow", func(ctx writeContext) (journal.Event, map[string]any, error) {
		if options.Rounds < 1 {
			return journal.Event{}, nil, fmt.Errorf("rounds must be at least 1")
		}
		var kind, state string
		if err := ctx.conn.QueryRow(`SELECT kind,status FROM issues WHERE id=?`, options.TaskID).Scan(&kind, &state); err != nil {
			return journal.Event{}, nil, err
		}
		if kind != "task" || state == "done" || state == "cancelled" {
			return journal.Event{}, nil, fmt.Errorf("rework allowance requires an open task")
		}
		id, effect, err := persondecision.Effect(ctx.conn, "rework-allowance", options.Actor, "task", options.TaskID, options.Reason, ctx.now)
		if err != nil {
			return journal.Event{}, nil, err
		}
		raw, err := json.Marshal(map[string]int64{"rounds": options.Rounds})
		if err != nil {
			return journal.Event{}, nil, err
		}
		effect.Values["binding_json"] = string(raw)
		return journal.Event{EventType: "rework.allowed", Aggregate: journal.Aggregate{ID: options.TaskID, Type: "task"}, Effects: []journal.Effect{effect}}, map[string]any{"task_id": options.TaskID, "decision_id": id, "rounds": options.Rounds, "reason": options.Reason}, nil
	})
}
