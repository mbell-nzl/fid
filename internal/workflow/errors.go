package workflow

import "fmt"

// NotFoundError identifies a workflow that is absent from the projection.
// It keeps adapter error classification independent of SQLite error text.
type NotFoundError struct{ ID string }

func (e *NotFoundError) Error() string { return fmt.Sprintf("workflow %q does not exist", e.ID) }

// DuplicateError identifies a workflow definition that would reuse a slug.
type DuplicateError struct{ ID string }

func (e *DuplicateError) Error() string { return fmt.Sprintf("workflow slug %q already exists", e.ID) }
