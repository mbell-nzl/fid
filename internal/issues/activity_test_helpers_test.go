package issues

import "database/sql"

func countIssueActivityRows(t interface {
	Helper()
	Fatalf(string, ...any)
}, conn *sql.DB) int {
	t.Helper()

	var count int
	if err := conn.QueryRow(`SELECT COUNT(*) FROM activity_log`).Scan(&count); err != nil {
		t.Fatalf("count activity_log rows: %v", err)
	}
	return count
}
