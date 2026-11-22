package reviewrecord

import (
	"database/sql"
	"fmt"
)

// Claim is a review session reservation, separate from the builder's kept hold.
// It never contributes to builder independence or model attribution.
type Claim struct {
	ID, TaskID, HandoffID, Stream, WorktreeID, Worktree, Actor, Model string
	Since, LastHeartbeatAt, ExpiresAt                                 string
	EndedAt                                                           *string
}

func Claims(conn Queryer) ([]Claim, error) {
	rows, err := conn.Query(`SELECT id,task_id,handoff_id,stream,worktree_id,worktree_root,actor,model,since,last_heartbeat_at,expires_at,ended_at FROM review_claims ORDER BY rowid`)
	if err != nil {
		return nil, fmt.Errorf("read review reservations: %w", err)
	}
	defer func() { _ = rows.Close() }()
	claims := []Claim{}
	for rows.Next() {
		var c Claim
		var end sql.NullString
		if err := rows.Scan(&c.ID, &c.TaskID, &c.HandoffID, &c.Stream, &c.WorktreeID, &c.Worktree, &c.Actor, &c.Model, &c.Since, &c.LastHeartbeatAt, &c.ExpiresAt, &end); err != nil {
			return nil, err
		}
		if end.Valid {
			c.EndedAt = &end.String
		}
		claims = append(claims, c)
	}
	return claims, rows.Err()
}

func BuiltBy(conn Queryer, task, actor string) (bool, error) {
	var count int
	err := conn.QueryRow(`SELECT count(*) FROM holds WHERE task_id=? AND actor=?`, task, actor).Scan(&count)
	return count > 0, err
}
