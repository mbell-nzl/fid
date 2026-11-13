package holdwork

import (
	"fmt"

	"github.com/mbell-nzl/fid/internal/holds"
	"github.com/mbell-nzl/fid/internal/journal"
	"github.com/mbell-nzl/fid/internal/persondecision"
	"github.com/mbell-nzl/fid/internal/questions"
)

func (ctx writeContext) parkWaiting(hold holds.Hold) (journal.Event, map[string]any, error) {
	id, err := journal.GenerateID("QST", func(id string) (bool, error) {
		var n int
		err := ctx.conn.QueryRow(`SELECT count(*) FROM questions WHERE id=?`, id).Scan(&n)
		return n > 0, err
	})
	if err != nil {
		return journal.Event{}, nil, err
	}
	effect := holds.End(hold.ID, "waiting", ctx.now)
	effect.Values["kept"] = true
	question := journal.Effect{Op: "insert", Table: "questions", PK: map[string]any{"id": id}, Values: map[string]any{
		"task_id": hold.TaskID, "hold_id": hold.ID, "text": ctx.options.Question, "asker": ctx.options.Actor, "asked_at": holds.Timestamp(ctx.now), "resolved_at": nil, "resolution": nil, "decision_id": nil,
	}}
	return journal.Event{EventType: "hold.waiting", Aggregate: journal.Aggregate{ID: hold.ID, Type: "hold"}, Effects: []journal.Effect{effect, question}}, map[string]any{"task_id": hold.TaskID, "hold_id": hold.ID, "mode": "waiting", "question_id": id, "question": ctx.options.Question, "kept_for_worktree": true}, nil
}

func Answer(options Options) (Result, error) {
	return write(options, "task.answer", func(ctx writeContext) (journal.Event, map[string]any, error) {
		history, err := questions.Read(ctx.conn, options.TaskID)
		if err != nil {
			return journal.Event{}, nil, err
		}
		for _, q := range history {
			if q.ResolvedAt != nil {
				continue
			}
			id, effect, err := persondecision.Effect(ctx.conn, "answer", options.Actor, "task", options.TaskID, options.Answer, ctx.now)
			if err != nil {
				return journal.Event{}, nil, err
			}
			effect.Values["binding_id"] = q.ID
			return journal.Event{EventType: "task.answered", Aggregate: journal.Aggregate{ID: options.TaskID, Type: "task"}, Effects: []journal.Effect{effect, q.Resolve(holds.Timestamp(ctx.now), "answered", id)}}, map[string]any{"task_id": options.TaskID, "question_id": q.ID, "decision_id": id, "answer": options.Answer}, nil
		}
		return journal.Event{}, nil, &Refusal{Code: "not-waiting", TaskID: options.TaskID, Message: fmt.Sprintf("%s isn't waiting for an answer", options.TaskID), Details: map[string]any{"task_id": options.TaskID}}
	})
}
