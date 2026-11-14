package journal

import (
	"database/sql"
	"strings"
)

func validateCurrentAppliedResults(tx *sql.Tx, event Event) error {
	if strings.HasPrefix(event.EventType, "check.") {
		return validateCheckResponse(tx, event)
	}
	if event.EventType == "artifact.availability-observed" {
		return validateArtifactResult(tx, event)
	}
	return nil
}
