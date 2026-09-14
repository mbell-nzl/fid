package workflow

import "fmt"

type StopNotFoundError struct{ ID string }

func (err *StopNotFoundError) Error() string { return fmt.Sprintf("Stop %q was not found", err.ID) }

type AttentionNotFoundError struct{ ID string }

func (err *AttentionNotFoundError) Error() string {
	return fmt.Sprintf("Attention %q was not found", err.ID)
}

type AttentionVersionConflictError struct {
	AttentionID     string
	ExpectedVersion int
	ActualVersion   int
}

func (err *AttentionVersionConflictError) Error() string { return "Attention version conflict" }

type ResolutionInvalidError struct {
	Subject      Subject
	FactEventIDs []string
	Reasons      []EvaluationReason
}

func (err *ResolutionInvalidError) Error() string { return "Stop resolution is invalid" }
