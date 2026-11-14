package holds

import "time"

// Working is an open hold as next, status and show describe it. Here says
// it belongs to this worktree; Lapsed says it had expired when read.
type Working struct {
	HoldID          string `json:"hold_id"`
	Worktree        string `json:"worktree"`
	Here            bool   `json:"here"`
	Actor           string `json:"actor"`
	Since           string `json:"since"`
	ExpiresAt       string `json:"expires_at"`
	LastHeartbeatAt string `json:"last_heartbeat_at"`
	Lapsed          bool   `json:"lapsed"`
}

// Describe reads a hold for a worktree at now without changing it.
func Describe(hold Hold, worktreeID string, now time.Time) *Working {
	if hold.EndedAt != nil {
		return nil
	}
	return &Working{HoldID: hold.ID, Worktree: hold.Worktree, Here: hold.WorktreeID == worktreeID, Actor: hold.Actor,
		Since: Timestamp(hold.Since), ExpiresAt: Timestamp(hold.ExpiresAt), Lapsed: hold.Lapsed(now), LastHeartbeatAt: Timestamp(hold.LastHeartbeat())}
}

// Kept describes interrupted or waiting work reserved for a worktree with no live hold.
type Kept struct {
	Worktree string  `json:"worktree"`
	Here     bool    `json:"here"`
	Mode     string  `json:"mode"`
	Reason   *string `json:"reason"`
	Since    string  `json:"since"`
	HoldID   string  `json:"hold_id"`
	Actor    string  `json:"actor"`
	Question *string `json:"question,omitempty"`
}

func DescribeKept(hold Hold, worktreeID string, now time.Time) *Kept {
	if hold.Live(now) {
		return nil
	}
	since := Timestamp(hold.LastHeartbeat())
	if hold.EndedAt != nil {
		since = *hold.EndedAt
	}
	kept := &Kept{Worktree: hold.Worktree, Here: hold.WorktreeID == worktreeID, Mode: "interrupted", Reason: hold.ParkReason, Since: since, HoldID: hold.ID, Actor: hold.Actor}
	if hold.Waiting() || (hold.Question != nil && hold.Kept && hold.ParkReason == nil) {
		kept.Mode = "waiting"
		kept.Question = &hold.Question.Text
		kept.Reason = nil
	}
	return kept
}
