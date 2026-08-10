package claims

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/mbell-nzl/fid/internal/config"
	"github.com/mbell-nzl/fid/internal/journal"
)

func RequireLiveIssueClaimHeld(conn *sql.DB, issueID, actorRole string) error {
	_, err := requireLiveIssueClaimHeldAt(conn, issueID, actorRole, operationTime(time.Time{}))
	return err
}

func requireLiveIssueClaimHeldAt(conn *sql.DB, issueID, actorRole string, now time.Time) (*currentClaim, error) {
	current, err := loadCurrentClaim(conn, issueID)
	if err != nil {
		return nil, err
	}
	if current == nil || !current.isLiveAt(now) || current.ClaimedBy != actorRole {
		return nil, claimNotHeldError(issueID, actorRole, current)
	}
	return current, nil
}

func nextClaimTime(currentClaimedAt, now time.Time) time.Time {
	candidate := now.UTC()
	if candidate.Format(time.RFC3339) <= currentClaimedAt.UTC().Format(time.RFC3339) {
		return currentClaimedAt.UTC().Add(time.Second)
	}
	return candidate
}

func loadCurrentClaim(conn *sql.DB, issueID string) (*currentClaim, error) {
	var (
		claimedBy string
		claimedAt string
		ttl       int
		note      sql.NullString
	)
	err := conn.QueryRow(`
		SELECT claimed_by, claimed_at, ttl_seconds, note
		FROM claims
		WHERE ref_type = 'issue' AND ref_id = ?
	`, issueID).Scan(&claimedBy, &claimedAt, &ttl, &note)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("query claim for issue %q: %w", issueID, err)
	}

	parsedClaimedAt, err := time.Parse(time.RFC3339, claimedAt)
	if err != nil {
		return nil, fmt.Errorf("parse claim time for issue %q: %w", issueID, err)
	}

	return &currentClaim{
		ClaimedBy:  claimedBy,
		ClaimedAt:  parsedClaimedAt.UTC(),
		TTLSeconds: ttl,
		Note:       nullStringPointer(note),
	}, nil
}

func nullStringPointer(value sql.NullString) *string {
	if !value.Valid {
		return nil
	}
	copy := value.String
	return &copy
}

func assertIssueExists(conn *sql.DB, issueID string) error {
	var existing string
	err := conn.QueryRow(`SELECT id FROM issues WHERE id = ? LIMIT 1`, issueID).Scan(&existing)
	if err == nil {
		return nil
	}
	if err == sql.ErrNoRows {
		return fmt.Errorf("acquire claim: issue %q not found", issueID)
	}
	return fmt.Errorf("query issue %q: %w", issueID, err)
}

func normalizeNote(note *string) *string {
	if note == nil {
		return nil
	}
	trimmed := strings.TrimSpace(*note)
	if trimmed == "" {
		return nil
	}
	return &trimmed
}

func projectClaimsEvent(workspaceRoot string, conn *sql.DB, event journal.Event) error {
	eventHash, err := journal.HashEvent(event)
	if err != nil {
		return fmt.Errorf("hash claim event %q: %w", event.EventID, err)
	}

	journalFile := event.OccurredAt.UTC().Format("2006-01-02") + ".jsonl"
	journalPath := filepath.Join(workspaceRoot, config.WorkspaceDirName, config.JournalDirName, journalFile)
	info, err := os.Stat(journalPath)
	if err != nil {
		return fmt.Errorf("stat claim journal %q: %w", journalPath, err)
	}

	return journal.ProjectEvent(conn, journal.ProjectionInput{
		Event:     event,
		EventHash: eventHash,
		File:      journalFile,
		Line:      countJournalLines(journalPath),
		Offset:    info.Size(),
	})
}

func countJournalLines(path string) int {
	payload, err := os.ReadFile(path)
	if err != nil {
		return 0
	}
	trimmed := strings.TrimSpace(string(payload))
	if trimmed == "" {
		return 0
	}
	return len(strings.Split(trimmed, "\n"))
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

func dereferenceString(value *string) any {
	if value == nil {
		return nil
	}
	return *value
}

func takeoverReason(current *currentClaim, force bool, now time.Time) string {
	if force && current.isLiveAt(now) {
		return "force_acquire"
	}
	return "stale_takeover"
}

func claimNotHeldError(issueID, actorRole string, current *currentClaim) error {
	err := &ClaimNotHeldError{
		RefID:       issueID,
		RequestedBy: actorRole,
	}
	if current != nil {
		err.ClaimedBy = current.ClaimedBy
		err.ExpiresAt = current.expiresAt().Format(time.RFC3339)
	}
	return err
}

func closeClaimsConn(conn *sql.DB) {
	_ = conn.Close()
}

func closeClaimsRows(rows *sql.Rows) {
	_ = rows.Close()
}
