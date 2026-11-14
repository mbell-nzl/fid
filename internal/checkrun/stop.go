package checkrun

import (
	"errors"

	"github.com/mbell-nzl/fid/internal/checkrecord"
	"github.com/mbell-nzl/fid/internal/holdwork"
	"github.com/mbell-nzl/fid/internal/journal"
)

// stop completes a request refused between children, while there is no
// unfinished start. A refusal must remain the same on a same-key repeat.
func (ex *execution) stop(db store, index int, cause error) error {
	ex.response.Complete = true
	ex.response.NotStarted = append([]string{}, ex.options.Names[index:]...)
	var hold *holdwork.Refusal
	var refused *checkrecord.Refusal
	switch {
	case errors.As(cause, &hold):
		ex.response.Failure = &checkrecord.Refusal{Code: hold.Code, TaskID: hold.TaskID, Message: hold.Message, Details: hold.Details}
	case errors.As(cause, &refused):
		ex.response.Failure = refused
	default:
		ex.response.Failure = &checkrecord.Refusal{Code: "internal-error", Message: "check failed: " + cause.Error(), Details: map[string]any{}}
	}
	event, err := db.event("check.stopped", ex.requestID, []journal.Effect{{Op: "update", Table: "check_requests", PK: map[string]any{"id": ex.requestID}, Values: map[string]any{"response_json": encoded(ex.response)}, Preconditions: map[string]any{"response_json": nil}}})
	if err != nil {
		return err
	}
	return ex.commitProgress(db, event)
}
