package status

import (
	"fmt"
	"slices"

	"github.com/mbell-nzl/fid/internal/reviewpost"
	"github.com/mbell-nzl/fid/internal/reviewrecord"
)

func (r *Result) readPendingPosts(conn reviewrecord.Queryer) error {
	pending, err := reviewpost.Pending(conn)
	if err != nil {
		return err
	}
	rounds := []string{}
	for round := range pending {
		rounds = append(rounds, round)
	}
	slices.Sort(rounds)
	for _, round := range rounds {
		r.Warnings = append(r.Warnings, NewWarning("review-post-pending", fmt.Sprintf("PR comments pending for review %s", round), "fid review post "+round, map[string]any{"round_id": round, "pending": pending[round]}))
	}
	return nil
}
