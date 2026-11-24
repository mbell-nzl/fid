package main

import (
	"context"
	"encoding/json"
	"fmt"
	"slices"
	"strings"
	"testing"

	"github.com/mbell-nzl/fid/internal/github"
	"github.com/mbell-nzl/fid/internal/journal"
	"github.com/mbell-nzl/fid/internal/reviewcomment"
	"github.com/mbell-nzl/fid/internal/reviewpost"
	"github.com/mbell-nzl/fid/internal/reviewwork"
	fidmigrations "github.com/mbell-nzl/fid/migrations"
)

func TestReviewPostRecoversAcrossRenderUpgrade(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name    string
		version int64
		edited  int // One-based subject index; zero means all comments are canonical.
	}{
		{"legacy receipt lost", 1, 0},
		{"current receipt lost", 2, 0},
		{"edited legacy round", 1, 1},
		{"edited legacy answer", 1, 2},
		{"edited legacy deferral", 1, 3},
		{"edited current round", 2, 1},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			f, task := canonicalReviewFixture(t)
			in := reviewInput()
			in["findings"] = []any{reviewFinding("minor")}
			recorded, _ := lifecycleResult(t, recordReviewTask(t, f, in, "upgrade-recovery", "human:courier", task), "review.record")
			round := recorded["round"].(map[string]any)["id"].(string)
			answer, _ := lifecycleResult(t, f.run("--json", "finding", "answer", round+"-1", "--actor", "agent:builder", "--text", "Accepted explanation"), "finding.answer")
			deferred, _ := lifecycleResult(t, f.run("--json", "finding", "defer", round+"-1", "--actor", "agent:builder", "--title", "Follow up"), "finding.defer")
			ids := []string{round, answer["response"].(map[string]any)["id"].(string), deferred["response"].(map[string]any)["id"].(string)}
			conn, err := reviewwork.OpenProjection(f.workspace)
			if err != nil {
				t.Fatal(err)
			}
			defer closeCommandDB(t, conn)
			comments := []map[string]string{}
			digests, urls := map[string]string{}, map[string]string{}
			for i, id := range ids {
				s, err := reviewcomment.ReadSubjectVersion(conn, id, tc.version)
				if err != nil {
					t.Fatal(err)
				}
				digests[id] = reviewcomment.Digest(s.Body)
				urls[id] = fmt.Sprintf("https://github.com/mbell-nzl/fid/issues/63#issuecomment-%d", i+100)
				body := s.Body + "\n" + reviewcomment.Marker(id, s.Body) + "\n"
				if i+1 == tc.edited {
					body = "Edited " + body // Retain the original marker.
				}
				comments = append(comments, map[string]string{"body": body, "html_url": urls[id]})
			}
			requests, writes := 0, 0
			run := func(_ context.Context, args, _ []string) github.Result {
				requests++
				if args[0] != "api" {
					writes++
					return github.Result{Stdout: []byte(fmt.Sprintf("https://github.com/mbell-nzl/fid/issues/63#issuecomment-%d", writes+1000))}
				}
				page := comments
				if strings.Contains(strings.Join(args, " "), "&page=1") {
					page = make([]map[string]string, 100) // Recovery must reach page two.
				}
				raw, err := json.Marshal(page)
				if err != nil {
					t.Fatal(err)
				}
				return github.Result{Stdout: raw}
			}
			before := startEventCount(t, f.workspace)
			got, err := reviewpost.Round(context.Background(), f.workspace, fidmigrations.FS, round, run)
			wantPosted, wantPending := slices.Clone(ids), []string{}
			if tc.edited != 0 {
				wantPending = []string{ids[tc.edited-1]}
				wantPosted = slices.Delete(wantPosted, tc.edited-1, tc.edited)
			}
			if err != nil || !slices.Equal(got.Posted, wantPosted) || !slices.Equal(got.Pending, wantPending) || writes != 0 {
				t.Fatalf("recovery=%+v / %v, writes=%d; want posted=%v pending=%v, no writes", got, err, writes, wantPosted, wantPending)
			}
			if added := startEventCount(t, f.workspace) - before; added != len(wantPosted) {
				t.Fatalf("recovery appended %d Events, want %d receipts", added, len(wantPosted))
			}
			if _, err := journal.RebuildWorkspace(f.workspace, fidmigrations.FS); err != nil {
				t.Fatal("rebuild recovered receipts", err)
			}
			rebuilt, err := reviewwork.OpenProjection(f.workspace)
			if err != nil {
				t.Fatal(err)
			}
			defer closeCommandDB(t, rebuilt)
			for _, id := range wantPosted {
				var version int64
				var digest, url string
				if err := rebuilt.QueryRow(`SELECT render_version,digest,url FROM review_posts WHERE id=?`, id).Scan(&version, &digest, &url); err != nil || version != tc.version || digest != digests[id] || url != urls[id] {
					t.Fatalf("receipt %s: version=%d digest=%s url=%s / %v; want %d %s %s", id, version, digest, url, err, tc.version, digests[id], urls[id])
				}
			}
			for _, id := range wantPending {
				var count int
				if err := rebuilt.QueryRow(`SELECT count(*) FROM review_posts WHERE id=?`, id).Scan(&count); err != nil || count != 0 {
					t.Fatalf("edited comment acquired %d receipts / %v", count, err)
				}
			}
			if tc.edited == 0 {
				requests = 0
				retried, err := reviewpost.Round(context.Background(), f.workspace, fidmigrations.FS, round, run)
				if err != nil || !slices.Equal(retried.Posted, ids) || len(retried.Pending) != 0 || requests != 0 || startEventCount(t, f.workspace)-before != len(ids) {
					t.Fatalf("receipt-backed retry=%+v / %v, requests=%d", retried, err, requests)
				}
			}
		})
	}
}
