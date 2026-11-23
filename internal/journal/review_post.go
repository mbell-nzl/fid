package journal

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/mbell-nzl/fid/internal/github"
	"github.com/mbell-nzl/fid/internal/reviewcomment"
)

func validateReviewPost(conn *sql.Tx, e Event) error {
	if e.CommandName != "review.post" || e.ActorV2.ActorID != "fid-review-post" || e.Aggregate.Sequence != 1 || e.Repository != nil || e.Idempotency != nil || len(e.Effects) != 1 {
		return fmt.Errorf("invalid review post Event")
	}
	p := e.Effects[0]
	if err := validateCurrentEffect(e, p); err != nil {
		return err
	}
	s, err := reviewcomment.ReadSubject(conn, e.Aggregate.ID)
	if err != nil {
		return err
	}
	url, ok := p.Values["url"].(string)
	if !ok || !github.ValidCommentURL(s.PR, url) {
		return fmt.Errorf("comment URL differs from recorded PR")
	}
	want := Effect{Op: "insert", Table: "review_posts", PK: map[string]any{"id": s.ID}, Values: map[string]any{"task_id": s.TaskID, "handoff_id": s.HandoffID, "pr": s.PR, "url": url, "render_version": reviewcomment.RenderVersion, "digest": reviewcomment.Digest(s.Body), "posted_at": e.OccurredAt.UTC().Format(time.RFC3339)}}
	if !currentValuesEqual(p, want) {
		return fmt.Errorf("review post differs from canonical record")
	}
	return nil
}
