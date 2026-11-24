package github

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestCommentReconcilesAcceptedPost(t *testing.T) {
	t.Parallel()
	for _, changed := range []bool{false, true} {
		t.Run(map[bool]string{false: "receipt lost", true: "edited comment"}[changed], func(t *testing.T) {
			t.Parallel()
			body := "Canonical review\n"
			full := body + "\n<!-- fid:key -->\n"
			if changed {
				full = "Edited " + full
			}
			run := func(_ context.Context, args, _ []string) Result {
				if args[0] != "api" {
					t.Fatal("recovery posted a second comment")
				}
				raw, _ := json.Marshal([]map[string]string{{"body": full, "html_url": "https://github.com/owner/repo/issues/63#issuecomment-123"}})
				return Result{Stdout: raw}
			}
			client, err := New(Repository{Owner: "owner", Name: "repo"}, "main", run)
			if err != nil {
				t.Fatal(err)
			}
			url, err := client.PostComment(context.Background(), "https://github.com/owner/repo/pull/63", "<!-- fid:key -->", body)
			if changed {
				if err == nil {
					t.Fatal("edited comment accepted")
				}
			} else if err != nil || url != "https://github.com/owner/repo/issues/63#issuecomment-123" {
				t.Fatalf("recovery=%s/%v", url, err)
			}
		})
	}
}

func TestCommentListingGuards(t *testing.T) {
	t.Parallel()
	for _, mode := range []string{"second page", "page limit", "size limit", "malformed", "wrong write URL"} {
		t.Run(mode, func(t *testing.T) {
			t.Parallel()
			calls, writes := 0, 0
			var deadline time.Time
			run := func(ctx context.Context, args, _ []string) Result {
				until, ok := ctx.Deadline()
				if !ok || !deadline.IsZero() && until != deadline {
					t.Fatal("requests must share one bounded lifetime")
				}
				deadline = until
				if args[0] != "api" {
					writes++
					return Result{Stdout: []byte("https://github.com/owner/repo/issues/64#issuecomment-1")}
				}
				calls++
				if !strings.Contains(strings.Join(args, " "), fmt.Sprintf("repos/owner/repo/issues/63/comments?per_page=100&page=%d", calls)) {
					t.Fatalf("wrong destination/page: %v", args)
				}
				if mode == "malformed" {
					return Result{Stdout: []byte("null")}
				}
				comments := []map[string]string{}
				if mode == "second page" && calls == 2 {
					comments = append(comments, map[string]string{"body": "Review\n<!-- key -->\n", "html_url": "https://github.com/owner/repo/issues/63#issuecomment-1"})
				} else if mode != "wrong write URL" {
					body := "Other comment"
					if mode == "size limit" {
						body = strings.Repeat("x", DefaultOutputLimit/150)
					}
					for range pageSize {
						comments = append(comments, map[string]string{"body": body})
					}
				}
				raw, _ := json.Marshal(comments)
				return Result{Stdout: raw}
			}
			client, err := NewCommenter(Repository{Owner: "owner", Name: "repo"}, run)
			if err != nil {
				t.Fatal(err)
			}
			url, err := client.PostComment(context.Background(), "https://github.com/owner/repo/pull/63", "<!-- key -->", "Review")
			if mode == "second page" {
				if err != nil || calls != 2 || writes != 0 || url != "https://github.com/owner/repo/issues/63#issuecomment-1" {
					t.Fatalf("paginated recovery=%s/%v, calls=%d writes=%d", url, err, calls, writes)
				}
			} else if err == nil || !strings.Contains(err.Error(), map[string]string{"page limit": "page limit", "size limit": "size limit", "malformed": "invalid comment listing", "wrong write URL": "invalid posted comment URL"}[mode]) || writes != map[bool]int{true: 1}[mode == "wrong write URL"] || calls > maxCIPages {
				t.Fatalf("listing guard=%s/%v, calls=%d writes=%d", url, err, calls, writes)
			}
		})
	}
}

func TestCommentURLs(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		url   string
		valid bool
	}{
		{"https://github.com/owner/repo/pull/63#issuecomment-123", true},
		{"https://github.com/owner/repo/issues/63#issuecomment-123", true},
		{"https://github.com/owner/repo/issues/64#issuecomment-123", false},
		{"https://github.com/owner/other/issues/63#issuecomment-123", false},
		{"https://github.com/owner/repo/issues/63", false},
		{"https://github.com/owner/repo/issues/63#issuecomment-0", false},
		{"https://elsewhere.example/owner/repo/issues/63#issuecomment-123", false},
	} {
		t.Run(tc.url, func(t *testing.T) {
			t.Parallel()
			if got := ValidCommentURL("https://github.com/owner/repo/pull/63", tc.url); got != tc.valid {
				t.Fatalf("valid=%t", got)
			}
		})
	}
}
