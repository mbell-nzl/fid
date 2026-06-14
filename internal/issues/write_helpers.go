package issues

import (
	"database/sql"
	"fmt"
	"time"
)

func nextIssueUpdateTime(currentUpdatedAt string, now time.Time) time.Time {
	parsed, err := time.Parse(time.RFC3339, currentUpdatedAt)
	if err != nil {
		return now.UTC()
	}

	candidate := now.UTC()
	if candidate.Format(time.RFC3339) <= parsed.UTC().Format(time.RFC3339) {
		return parsed.UTC().Add(time.Second)
	}

	return candidate
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
