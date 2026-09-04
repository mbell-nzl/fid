package cli

import (
	"errors"
	"fmt"
	"io"
)

const CategoryGateUnsuccessful Category = "gate-unsuccessful"

// ExitCodeForCategory returns the process status reserved by fid-cli.v1.
func ExitCodeForCategory(category Category) int {
	switch category {
	case CategoryInvalidInput:
		return 2
	case CategoryNotFound:
		return 3
	case CategoryConflict, CategoryPrecondition:
		return 4
	case CategoryAuthority:
		return 5
	case CategoryIntegrity:
		return 6
	case CategoryGateUnsuccessful:
		return 8
	default:
		return 7
	}
}

// FailureError carries the emitted failure through Cobra to the process entrypoint.
type FailureError struct {
	Code     Code
	Category Category
	cause    error
}

func (err *FailureError) Error() string {
	if err.cause != nil {
		return err.cause.Error()
	}
	return fmt.Sprintf("fid failure: %s", err.Code)
}

func NewFailureError(failure Failure, cause ...error) *FailureError {
	err := &FailureError{Code: failure.Code, Category: failure.Category}
	if len(cause) > 0 {
		err.cause = cause[0]
	}
	return err
}

// WriteFailureAndReturn emits a closed envelope then preserves its exit class.
func WriteFailureAndReturn(writer io.Writer, operation string, failure Failure, cause ...error) error {
	if err := WriteFailure(writer, operation, failure); err != nil {
		return err
	}
	return NewFailureError(failure, cause...)
}

func ExitCode(err error) int {
	var failure *FailureError
	if errors.As(err, &failure) {
		return ExitCodeForCategory(failure.Category)
	}
	return ExitCodeForCategory(CategoryInternal)
}
