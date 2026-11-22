package export

import (
	"errors"

	"github.com/mbell-nzl/fid/internal/journal"
)

func semanticCheckCode(err error) string {
	var unsupported *journal.SchemaUnsupportedError
	if errors.As(err, &unsupported) {
		return "schema-unsupported"
	}
	var semanticErr *semanticValidationError
	if errors.As(err, &semanticErr) {
		return semanticErr.Code
	}
	return ""
}
