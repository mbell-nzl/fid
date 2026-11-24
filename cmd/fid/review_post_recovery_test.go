package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"slices"
	"testing"

	"github.com/mbell-nzl/fid/internal/github"
	"github.com/mbell-nzl/fid/internal/reviewpost"
	"github.com/mbell-nzl/fid/internal/reviewwork"
)

func TestReviewReceiptProjectionRecovery(t *testing.T) {
	t.Parallel()
	for _, receiptFailure := range []bool{true, false} {
		name := "transport"
		if receiptFailure {
			name = "receipt projection"
		}
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			f, task := canonicalReviewFixture(t)
			in := reviewInput()
			in["findings"] = []any{reviewFinding("minor")}
			recorded, _ := lifecycleResult(t, recordReviewTask(t, f, in, "receipt-recovery", "human:courier", task), "review.record")
			round := recorded["round"].(map[string]any)
			id := round["id"].(string)
			finding := round["findings"].([]any)[0].(map[string]any)["id"].(string)
			answered, _ := lifecycleResult(t, f.run("--json", "finding", "answer", finding, "--actor", "agent:builder", "--text", "Accepted explanation"), "finding.answer")
			answer := answered["response"].(map[string]any)["id"].(string)
			conn, err := reviewwork.OpenProjection(f.workspace)
			if err != nil {
				t.Fatal(err)
			}
			defer closeCommandDB(t, conn)
			if receiptFailure {
				// Dry-run admission does not write the checkpoint. This fails only
				// after the receipt's durable append, rolling back its projection.
				_, err = conn.Exec(fmt.Sprintf(`CREATE TRIGGER fail_round_receipt BEFORE UPDATE ON projection_meta
WHEN EXISTS (SELECT 1 FROM review_posts WHERE id='%s')
BEGIN SELECT RAISE(ABORT, 'injected receipt checkpoint failure'); END`, id))
				if err != nil {
					t.Fatal(err)
				}
			}
			before := startEventCount(t, f.workspace)
			run, requests, writes := recoveryCommentRunner(t, !receiptFailure)
			first, err := reviewpost.Round(context.Background(), f.workspace, issueMigrationsDir(), id, run)
			if err != nil {
				t.Fatal(err)
			}
			wantPosted, wantPending := []string{answer}, []string{id}
			if receiptFailure {
				wantPosted, wantPending = []string{}, []string{id, answer}
				if _, err := conn.Exec(`DROP TRIGGER fail_round_receipt`); err != nil {
					t.Fatal(err)
				}
			}
			if !slices.Equal(first.Posted, wantPosted) || !slices.Equal(first.Pending, wantPending) {
				t.Errorf("first delivery=%+v; want posted=%v pending=%v", first, wantPosted, wantPending)
			}
			if added := startEventCount(t, f.workspace) - before; added != 1 {
				t.Errorf("first delivery appended %d Events, want one receipt", added)
			}
			if *writes != 1 {
				t.Errorf("first delivery wrote %d comments, want one", *writes)
			}
			for attempt := range 2 {
				retried, retryErr := reviewpost.Round(context.Background(), f.workspace, issueMigrationsDir(), id, run)
				if retryErr != nil || len(retried.Pending) != 0 || !slices.Equal(retried.Posted, []string{id, answer}) {
					t.Fatalf("ordinary retry %d did not recover: %+v / %v", attempt+1, retried, retryErr)
				}
				if attempt == 0 {
					*requests = 0
				} else if *requests != 0 {
					t.Fatal("receipt-backed retry contacted gh")
				}
			}
			if *writes != 2 || startEventCount(t, f.workspace)-before != 2 {
				t.Fatalf("recovery duplicated delivery: comments=%d, receipts=%d", *writes, startEventCount(t, f.workspace)-before)
			}
			var receipts int
			if err := conn.QueryRow(`SELECT count(*) FROM review_posts WHERE id IN (?,?)`, id, answer).Scan(&receipts); err != nil || receipts != 2 {
				t.Fatalf("projected receipts=%d / %v, want two", receipts, err)
			}
		})
	}
}

func recoveryCommentRunner(t *testing.T, failTransport bool) (github.Runner, *int, *int) {
	t.Helper()
	comments := []map[string]string{}
	requests, writes := 0, 0
	run := func(_ context.Context, args, _ []string) github.Result {
		requests++
		if failTransport {
			failTransport = false
			return github.Result{ExitCode: 1}
		}
		if args[0] == "api" {
			raw, err := json.Marshal(comments)
			if err != nil {
				t.Fatal(err)
			}
			return github.Result{Stdout: raw}
		}
		index := slices.Index(args, "--body-file")
		if index < 0 || index+1 >= len(args) {
			t.Fatalf("missing comment body: %v", args)
		}
		body, err := os.ReadFile(args[index+1])
		if err != nil {
			t.Fatal(err)
		}
		writes++
		url := fmt.Sprintf("https://github.com/mbell-nzl/fid/issues/63#issuecomment-%d", writes)
		comments = append(comments, map[string]string{"body": string(body), "html_url": url})
		return github.Result{Stdout: []byte(url)}
	}
	return run, &requests, &writes
}
