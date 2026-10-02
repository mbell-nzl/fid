package journal

import (
	"database/sql"
	"fmt"
	"regexp"

	fiddb "github.com/mbell-nzl/fid/internal/db"
)

var identifierPattern = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*$`)

type ProjectionInput struct {
	Event      Event
	EventHash  string
	File       string
	Line       int
	Offset     int64
	ReplayMode string
}

type PreconditionConflictError struct {
	EventID   string
	EventType string
	Table     string
	PK        map[string]any
	Expected  map[string]any
	Actual    map[string]any
	File      string
	Line      int
}

func (err *PreconditionConflictError) Error() string {
	if err == nil {
		return "<nil>"
	}

	return fmt.Sprintf(
		"replay precondition conflict at %s:%d for event %s (%s) on %s pk=%v expected=%v actual=%v",
		err.File,
		err.Line,
		err.EventID,
		err.EventType,
		err.Table,
		err.PK,
		err.Expected,
		err.Actual,
	)
}

func DryRunProjectEvent(conn *sql.DB, event Event) error {
	if err := validateProjectionConnection(conn); err != nil {
		return err
	}

	input := ProjectionInput{
		Event: event,
	}
	if input.Event.EventID == "" {
		input.Event.EventID = "EVT-dry-run"
	}

	tx, err := conn.Begin()
	if err != nil {
		return fmt.Errorf("project event %q: begin transaction: %w", input.Event.EventID, err)
	}

	if err := applyEventEffects(tx, input, tx); err != nil {
		_ = tx.Rollback()
		return err
	}
	if err := tx.Rollback(); err != nil {
		return fmt.Errorf("project event %q: rollback transaction: %w", input.Event.EventID, err)
	}

	return nil
}

func ProjectEvent(conn *sql.DB, input ProjectionInput) error {
	if err := validateProjectionConnection(conn); err != nil {
		return err
	}
	if err := validateProjectionInput(input); err != nil {
		return err
	}

	tx, err := conn.Begin()
	if err != nil {
		return fmt.Errorf("project event %q: begin transaction: %w", input.Event.EventID, err)
	}

	return projectEventTransaction(tx, input, tx)
}

func projectEventTransaction(tx *sql.Tx, input ProjectionInput, executor fiddb.ProjectionMetaExecutor) error {
	if err := applyEventEffects(tx, input, executor); err != nil {
		_ = tx.Rollback()
		return err
	}

	replayStatus := input.ReplayMode
	if replayStatus == "" {
		replayStatus = fiddb.ReplayStatusClean
	}
	if err := fiddb.WriteProjectionMetaWithExecutor(executor, fiddb.ProjectionMeta{
		JournalFile:             input.File,
		JournalOffset:           input.Offset,
		LastEventID:             input.Event.EventID,
		LastEventHash:           input.EventHash,
		EventSchema:             input.Event.EventSchema,
		ProjectionSchemaVersion: fiddb.ProjectionSchemaVersion,
		ReplayStatus:            replayStatus,
		UpdatedAt:               input.Event.OccurredAt.UTC(),
	}); err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("project event %q: %w", input.Event.EventID, err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("project event %q: commit transaction: %w", input.Event.EventID, err)
	}

	return nil
}

func validateProjectionConnection(conn *sql.DB) error {
	if conn == nil {
		return fmt.Errorf("project event: connection is nil")
	}
	return nil
}

func validateProjectionInput(input ProjectionInput) error {
	if input.Event.EventID == "" {
		return fmt.Errorf("project event: event_id is required")
	}
	if input.EventHash == "" {
		return fmt.Errorf("project event %q: event hash is required", input.Event.EventID)
	}
	return nil
}
