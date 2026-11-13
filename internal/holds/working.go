package holds

import "time"

// Working is an open hold as next, status and show describe it. Here says
// it belongs to this worktree; Lapsed says it had expired when read.
type Working struct {
	HoldID    string `json:"hold_id"`
	Worktree  string `json:"worktree"`
	Here      bool   `json:"here"`
	Actor     string `json:"actor"`
	Since     string `json:"since"`
	ExpiresAt string `json:"expires_at"`
	Lapsed    bool   `json:"lapsed"`
}

// Describe reads a hold for a worktree at now without changing it.
func Describe(hold Hold, worktreeID string, now time.Time) *Working {
	return &Working{HoldID: hold.ID, Worktree: hold.Worktree, Here: hold.WorktreeID == worktreeID, Actor: hold.Actor,
		Since: Timestamp(hold.Since), ExpiresAt: Timestamp(hold.ExpiresAt), Lapsed: hold.Lapsed(now)}
}
