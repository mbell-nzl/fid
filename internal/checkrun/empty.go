package checkrun

import (
	"time"

	"github.com/mbell-nzl/fid/internal/journal"
)

// An empty selection is a completed request, with no invented start or result.
func (ex *execution) empty(db store) error {
	id, err := journal.GenerateID("CHK", func(candidate string) (bool, error) {
		var n int
		err := db.conn.QueryRow(`SELECT count(*) FROM check_requests WHERE id=?`, candidate).Scan(&n)
		return n > 0, err
	})
	if err != nil {
		return err
	}
	ex.response.Complete = true
	event, err := db.event("check.completed", id, []journal.Effect{{Op: "insert", Table: "check_requests", PK: map[string]any{"id": id}, Values: map[string]any{"task_id": ex.response.TaskID, "hold_id": ex.options.HoldID, "actor": ex.options.Actor, "names_json": "[]", "response_json": nil}}})
	if err != nil {
		return err
	}
	_, err = journal.CommitIdempotent(db.lock, db.conn, ex.request, event, func(prepared journal.Event) (journal.IdempotencyMutation, error) {
		ex.response.EventID = prepared.EventID
		ex.response.OccurredAt = prepared.OccurredAt.UTC().Format(time.RFC3339Nano)
		prepared.Effects[0].Values["response_json"] = encoded(ex.response)
		return journal.IdempotencyMutation{EntityID: id, ResultJSON: encoded(ex.response)}, nil
	})
	return err
}
