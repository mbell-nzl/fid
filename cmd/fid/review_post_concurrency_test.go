package main

import (
	"context"
	"fmt"
	"sync/atomic"
	"testing"
	"time"

	"github.com/mbell-nzl/fid/internal/github"
	"github.com/mbell-nzl/fid/internal/reviewpost"
)

func TestReviewConcurrentPosting(t *testing.T) {
	t.Parallel()
	f, task := canonicalReviewFixture(t)
	fields, _ := lifecycleResult(t, recordReviewTask(t, f, reviewInput(), "concurrent", "human:courier", task), "review.record")
	id := fields["round"].(map[string]any)["id"].(string)
	var queries, writes atomic.Int32
	secondQuery := make(chan struct{})
	run := func(_ context.Context, args, _ []string) github.Result {
		if args[0] == "api" {
			switch queries.Add(1) {
			case 1:
				select {
				case <-secondQuery:
				case <-time.After(100 * time.Millisecond):
				}
			case 2:
				close(secondQuery)
			}
			return github.Result{Stdout: []byte("[]")}
		}
		n := writes.Add(1)
		return github.Result{Stdout: []byte(fmt.Sprintf("https://github.com/mbell-nzl/fid/issues/63#issuecomment-%d", n))}
	}
	type outcome struct {
		result reviewpost.Result
		err    error
	}
	results := make(chan outcome, 2)
	for range 2 {
		go func() {
			r, err := reviewpost.Round(context.Background(), f.workspace, issueMigrationsDir(), id, run)
			results <- outcome{r, err}
		}()
	}
	for range 2 {
		o := <-results
		if o.err != nil || len(o.result.Pending) != 0 {
			t.Fatalf("post=%v/%v", o.result, o.err)
		}
	}
	if queries.Load() != 1 || writes.Load() != 1 {
		t.Fatalf("concurrent requests=%d, posts=%d; want one each", queries.Load(), writes.Load())
	}
}
