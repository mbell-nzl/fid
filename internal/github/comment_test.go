package github

import (
	"context"
	"encoding/json"
	"testing"
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
