package reviewpost

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io/fs"
	"strings"
	"time"

	"github.com/mbell-nzl/fid/internal/github"
	"github.com/mbell-nzl/fid/internal/journal"
	"github.com/mbell-nzl/fid/internal/reviewcomment"
	"github.com/mbell-nzl/fid/internal/reviewrecord"
	"github.com/mbell-nzl/fid/internal/reviewwork"
)

type Result struct {
	Posted  []string `json:"posted"`
	Pending []string `json:"pending"`
}

var errReceiptCommit = errors.New("commit review post receipt")

// Round also retries answers and deferrals to that round's findings. Posting
// failures leave the canonical records intact and their derived queue pending.
func Round(ctx context.Context, root string, migrations fs.FS, round string, run github.Runner) (Result, error) {
	if !strings.HasPrefix(round, "RND-") {
		return Result{}, fmt.Errorf("review post requires a round")
	}
	result := Result{Posted: []string{}, Pending: []string{}}
	err := journal.WithLock(root, journal.LockOwner{Command: "fid review post"}, func(lock *journal.WriteLock) error {
		if _, err := journal.EnsureProjectionFreshWithLock(lock, migrations); err != nil {
			return err
		}
		conn, err := reviewwork.OpenProjection(root)
		if err != nil {
			return err
		}
		defer func() { _ = conn.Close() }()
		if _, err := reviewcomment.ReadSubject(conn, round); err != nil {
			if errors.Is(err, reviewcomment.ErrSealed) {
				return nil
			}
			return err
		}
		ids := []string{round}
		rows, err := conn.Query(`SELECT a.id FROM finding_responses a JOIN round_findings f ON f.id=a.finding_id WHERE f.round_id=? ORDER BY a.rowid`, round)
		if err != nil {
			return err
		}
		for rows.Next() {
			var id string
			if err := rows.Scan(&id); err != nil {
				_ = rows.Close()
				return err
			}
			ids = append(ids, id)
		}
		if err := rows.Err(); err != nil {
			_ = rows.Close()
			return err
		}
		_ = rows.Close()
		for i, id := range ids {
			if err := post(ctx, lock, conn, id, run); err != nil {
				if errors.Is(err, errReceiptCommit) {
					// A receipt may already be appended but not projected. Leave
					// the checkpoint behind it so the next invocation catches up.
					result.Pending = append(result.Pending, ids[i:]...)
					break
				}
				result.Pending = append(result.Pending, id)
			} else {
				result.Posted = append(result.Posted, id)
			}
		}
		return nil
	})
	return result, err
}

func post(ctx context.Context, lock *journal.WriteLock, conn *sql.DB, id string, run github.Runner) error {
	var count int
	if err := conn.QueryRow(`SELECT count(*) FROM review_posts WHERE id=?`, id).Scan(&count); err != nil {
		return err
	}
	if count != 0 {
		return nil
	}
	s, err := reviewcomment.ReadSubject(conn, id)
	if err != nil {
		return err
	}
	legacy, err := reviewcomment.ReadSubjectVersion(conn, id, 1)
	if err != nil {
		return err
	}
	p, err := github.ParsePRURL(s.PR)
	if err != nil {
		return err
	}
	client, err := github.NewCommenter(p.Repository, run)
	if err != nil {
		return err
	}
	url, variant, err := client.PostCommentVariants(ctx, s.PR,
		github.CommentBody{Marker: reviewcomment.Marker(s.ID, s.Body), Body: s.Body},
		github.CommentBody{Marker: reviewcomment.Marker(legacy.ID, legacy.Body), Body: legacy.Body})
	if err != nil {
		return err
	}
	version := reviewcomment.RenderVersion
	if variant == 1 {
		s, version = legacy, 1
	}
	now := time.Now().UTC().Truncate(time.Second)
	e := journal.Event{EventSchema: 2, EventType: "review.posted", CommandName: "review.post", OccurredAt: now, ActorV2: journal.EventActor{ActorID: "fid-review-post", ActorType: "service", Producer: "fid", ProvenanceRef: "fid:review-post"}, Aggregate: journal.Aggregate{ID: id, Type: "review-post", Sequence: 1}, Effects: []journal.Effect{{Op: "insert", Table: "review_posts", PK: map[string]any{"id": id}, Values: map[string]any{"task_id": s.TaskID, "handoff_id": s.HandoffID, "pr": s.PR, "url": url, "render_version": version, "digest": reviewcomment.Digest(s.Body), "posted_at": now.Format(time.RFC3339)}}}}
	_, err = journal.CommitEvent(lock, conn, e)
	if err != nil {
		return fmt.Errorf("%w: %w", errReceiptCommit, err)
	}
	return nil
}

func Pending(conn reviewrecord.Queryer) (map[string][]string, error) {
	result := map[string][]string{}
	rows, err := conn.Query(`SELECT r.id,r.id FROM review_rounds r LEFT JOIN review_posts p ON p.id=r.id WHERE p.id IS NULL UNION ALL SELECT f.round_id,a.id FROM finding_responses a JOIN round_findings f ON f.id=a.finding_id LEFT JOIN review_posts p ON p.id=a.id WHERE p.id IS NULL ORDER BY 1,2`)
	if err != nil {
		return nil, fmt.Errorf("read pending reviews: %w", err)
	}
	defer func() { _ = rows.Close() }()
	type pendingItem struct{ round, id string }
	items := []pendingItem{}
	for rows.Next() {
		var round, id string
		if err := rows.Scan(&round, &id); err != nil {
			return nil, err
		}
		items = append(items, pendingItem{round, id})
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	_ = rows.Close()
	for _, item := range items {
		round, id := item.round, item.id
		if sealed, err := reviewrecord.RoundSealed(conn, round); err == nil && sealed {
			continue
		} else if err != nil {
			return nil, err
		}
		result[round] = append(result[round], id)
	}
	return result, rows.Err()
}
