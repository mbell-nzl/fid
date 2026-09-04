package issues

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/mbell-nzl/fid/internal/journal"
)

func issueMutationNow(now func() time.Time) time.Time {
	if now == nil {
		return time.Now().UTC()
	}
	return now().UTC()
}

func activityIDExists(conn *sql.DB, activityID string) (bool, error) {
	var existing string
	err := conn.QueryRow(`SELECT id FROM activity_log WHERE id = ? LIMIT 1`, activityID).Scan(&existing)
	if err == nil {
		return true, nil
	}
	if err == sql.ErrNoRows {
		return false, nil
	}

	return false, fmt.Errorf("query activity id %q: %w", activityID, err)
}

func nextIssueActivityID(conn *sql.DB) (string, error) {
	return journal.GenerateID("ACT", func(candidate string) (bool, error) {
		return activityIDExists(conn, candidate)
	})
}

func issueActivityEffect(
	activityID, eventID, actorRole, timestamp, refID, eventName string,
	details map[string]any,
) journal.Effect {
	return journal.Effect{
		Op:    "activity",
		Table: "activity_log",
		PK:    map[string]any{"id": activityID},
		Values: map[string]any{
			"journal_event_id": eventID,
			"ref_type":         "issue",
			"ref_id":           refID,
			"event":            eventName,
			"details":          details,
			"actor":            actorRole,
			"occurred_at":      timestamp,
		},
	}
}
