package maintenance

import (
	"database/sql"
	"fmt"
	"time"
)

func checkInfo(db *sql.DB, now time.Time) (Info, error) {
	issueCount, err := countIssues(db)
	if err != nil {
		return Info{}, err
	}
	claims, err := loadClaims(db)
	if err != nil {
		return Info{}, err
	}
	liveClaims, err := countLiveClaims(claims, now)
	if err != nil {
		return Info{}, err
	}
	recentActivity, err := countRecentActivity(db, now)
	if err != nil {
		return Info{}, err
	}
	return Info{IssueCount: issueCount, ClaimCountLive: liveClaims, RecentActivity: recentActivity}, nil
}

func countIssues(db *sql.DB) (int, error) {
	var count int
	if err := db.QueryRow(`SELECT COUNT(*) FROM issues`).Scan(&count); err != nil {
		return 0, fmt.Errorf("check info: count issues: %w", err)
	}
	return count, nil
}

func countLiveClaims(claims []claimRecord, now time.Time) (int, error) {
	count := 0
	for _, claim := range claims {
		live, _, err := claim.liveAt(now)
		if err != nil {
			return 0, err
		}
		if live {
			count++
		}
	}
	return count, nil
}

func countRecentActivity(db *sql.DB, now time.Time) (int, error) {
	rows, err := db.Query(`SELECT occurred_at FROM activity_log ORDER BY rowid`)
	if err != nil {
		return 0, fmt.Errorf("check info: query activity: %w", err)
	}
	defer func() { _ = rows.Close() }()

	cutoff := now.Add(-7 * 24 * time.Hour)
	count := 0
	for rows.Next() {
		var occurredAt sql.NullString
		if err := rows.Scan(&occurredAt); err != nil {
			return 0, fmt.Errorf("check info: scan activity: %w", err)
		}
		if !occurredAt.Valid {
			continue
		}
		at, err := time.Parse(time.RFC3339, occurredAt.String)
		if err != nil {
			return 0, fmt.Errorf("check info activity time: %w", err)
		}
		if !at.UTC().Before(cutoff) && !at.UTC().After(now) {
			count++
		}
	}
	if err := rows.Err(); err != nil {
		return 0, fmt.Errorf("check info: iterate activity: %w", err)
	}
	return count, nil
}
