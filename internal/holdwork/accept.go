package holdwork

import (
	"database/sql"
	"fmt"

	"github.com/mbell-nzl/fid/internal/journal"
	"github.com/mbell-nzl/fid/internal/persondecision"
)

// Accept keeps one task-bound decision for each acceptance. Repeating it is
// harmless, including accepting a task originally created by a person.
func Accept(options Options) (Result, error) {
	return write(options, "task.accept", func(ctx writeContext) (journal.Event, map[string]any, error) {
		var kind, category string
		if err := ctx.conn.QueryRow(`SELECT i.kind,s.category FROM issues i JOIN workflow_states s ON s.workflow_id=i.workflow_id AND s.state=i.status WHERE i.id=?`, options.TaskID).Scan(&kind, &category); err != nil {
			return journal.Event{}, nil, fmt.Errorf("accept requires an existing task or subtask: %w", err)
		}
		if (kind != "task" && kind != "subtask") || category == "completed" || category == "cancelled" {
			return journal.Event{}, nil, fmt.Errorf("accept requires an open task or subtask")
		}
		fields := map[string]any{"issue_id": options.TaskID, "accepted": true}
		var id string
		err := ctx.conn.QueryRow(`SELECT id FROM person_decisions WHERE kind='accept' AND subject_id=? ORDER BY rowid LIMIT 1`, options.TaskID).Scan(&id)
		if err == nil {
			fields["decision_id"] = id
			return journal.Event{}, fields, nil
		}
		if err != sql.ErrNoRows {
			return journal.Event{}, nil, err
		}
		reason := options.Reason
		if reason == "" {
			reason = "Accepted for work"
		}
		id, effect, err := persondecision.Effect(ctx.conn, "accept", options.Actor, kind, options.TaskID, reason, ctx.now)
		if err != nil {
			return journal.Event{}, nil, err
		}
		fields["decision_id"], fields["reason"] = id, reason
		return journal.Event{EventType: "work.accepted", Aggregate: journal.Aggregate{ID: options.TaskID, Type: "issue"}, Effects: []journal.Effect{effect}}, fields, nil
	})
}
