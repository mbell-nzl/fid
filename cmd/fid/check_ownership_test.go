package main

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/mbell-nzl/fid/internal/checkrecord"
	"github.com/mbell-nzl/fid/internal/checkrun"
	"github.com/mbell-nzl/fid/internal/config"
	fiddb "github.com/mbell-nzl/fid/internal/db"
	"github.com/mbell-nzl/fid/internal/holdwork"
)

func leanOptions(f *statusFixture, key string) checkrun.Options {
	return checkrun.Options{WorkspaceRoot: f.workspace, MigrationsDir: issueMigrationsDir(), HoldID: "HLD-aaaaaa", Actor: "agent:claude-builder", IdempotencyKey: key, Names: []string{"lean"}, Clock: func() time.Time { return statusNow }, Version: "test", RenewInterval: 20 * time.Millisecond}
}

func TestLeanCheckIncompleteKeyAndLiveEarlierStart(t *testing.T) {
	t.Parallel()
	f := leanJourneyFixture(t, leanHelperDefinition(t, "pass", "1m", nil))
	options := leanOptions(f, "interrupted")
	cut := errors.New("cut off before spawning")
	options.Before = func(string) error { return cut }
	if _, err := checkrun.Run(context.Background(), options); !errors.Is(err, cut) {
		t.Fatalf("cut off=%v", err)
	}
	for _, row := range []struct{ key, code string }{{"interrupted", "check-incomplete"}, {"different", "check-still-running"}} {
		got := f.run("--json", "check", "lean", "--hold", "HLD-aaaaaa", "--actor", "agent:claude-builder", "--idempotency-key", row.key)
		assertJSONFailureCode(t, got.stderr, "check.execute", row.code)
	}
	if record := latestLeanRecord(t, f); record.Result != nil || record.Process != nil {
		t.Fatalf("invented result/process=%#v", record)
	}
}

func TestLeanCheckCancellationDoesNotStartLaterNames(t *testing.T) {
	t.Parallel()
	f := leanJourneyFixture(t, leanHelperDefinition(t, "wait", "1m", nil)+"\n[checks.zpass]\nrun=['true']\ntimeout='1m'\n")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	options := leanOptions(f, "cancel")
	options.Names = []string{"lean", "zpass"}
	options.Before = func(string) error { cancel(); return nil }
	response, err := checkrun.Run(ctx, options)
	var refusal *checkrecord.Refusal
	if !errors.As(err, &refusal) || refusal.Code != "check-unsuccessful" || len(response.Checks) != 1 || response.Checks[0].Outcome != "cancelled" || len(response.NotStarted) != 1 || response.NotStarted[0] != "zpass" {
		t.Fatalf("cancel=%#v/%v", response, err)
	}
}

func TestLeanCheckOwnershipLostBetweenChecksReplaysTheOriginalRefusal(t *testing.T) {
	t.Parallel()
	f := leanJourneyFixture(t, leanHelperDefinition(t, "pass", "1m", nil)+"\n[checks.zpass]\nrun=['true']\ntimeout='1m'\n")
	options := leanOptions(f, "between")
	options.Names = []string{"lean", "zpass"}
	options.After = func(checkrecord.Record) error {
		_, err := holdwork.Takeover(holdwork.Options{WorkspaceRoot: f.workspace, MigrationsDir: issueMigrationsDir(), TaskID: "FID-9f13aa", Actor: "human:mark", Reason: "Take control", Clock: func() time.Time { return statusNow }})
		return err
	}
	response, err := checkrun.Run(context.Background(), options)
	if err == nil || len(response.Checks) != 1 {
		t.Fatalf("between=%#v/%v", response, err)
	}
	got := f.run("--json", "check", "lean", "zpass", "--hold", "HLD-aaaaaa", "--actor", "agent:claude-builder", "--idempotency-key", "between")
	assertJSONFailureCode(t, got.stderr, "check.execute", "hold-not-current")
}

func TestLeanCheckRenewsHoldAndStopsOnTakeover(t *testing.T) {
	t.Parallel()
	for _, takeover := range []bool{false, true} {
		t.Run(map[bool]string{false: "renewal", true: "takeover"}[takeover], func(t *testing.T) {
			t.Parallel()
			mode := "pause"
			if takeover {
				mode = "wait"
			}
			f := leanJourneyFixture(t, leanHelperDefinition(t, mode, "1m", nil)+"\n[checks.zpass]\nrun=['true']\ntimeout='1m'\n")
			options := leanOptions(f, "check")
			options.Names = []string{"lean", "zpass"}
			var taken chan error
			if takeover {
				taken = make(chan error, 1)
				go func() {
					deadline := time.Now().Add(10 * time.Second)
					for time.Now().Before(deadline) {
						conn, err := fiddb.Open(filepath.Join(config.LedgerDir(f.workspace), config.DatabaseFileName))
						if err != nil {
							taken <- err
							return
						}
						var count int
						err = conn.QueryRow(`SELECT count(*) FROM check_processes`).Scan(&count)
						_ = conn.Close()
						if err != nil {
							taken <- err
							return
						}
						if count > 0 {
							_, err := holdwork.Takeover(holdwork.Options{WorkspaceRoot: f.workspace, MigrationsDir: issueMigrationsDir(), TaskID: "FID-9f13aa", Actor: "human:mark", Reason: "Take control", Clock: func() time.Time { return statusNow }})
							taken <- err
							return
						}
						time.Sleep(5 * time.Millisecond)
					}
					taken <- errors.New("check process did not start")
				}()
			}
			response, err := checkrun.Run(context.Background(), options)
			if takeover {
				if takeErr := <-taken; takeErr != nil {
					t.Fatal(takeErr)
				}
				var refused *checkrecord.Refusal
				if !errors.As(err, &refused) || refused.Code != "hold-not-current" || len(response.NotStarted) != 1 || response.NotStarted[0] != "zpass" {
					t.Fatalf("takeover=%#v, %v", response, err)
				}
				record := latestLeanRecord(t, f)
				if record.Result == nil || !strings.Contains(strings.Join(record.Result.Reasons, ","), "ownership-lost") || record.Process == nil {
					t.Fatalf("record=%#v", record)
				}
				if err := syscall.Kill(record.Process.PID, 0); err == nil {
					t.Fatal("child survived takeover")
				}
			} else {
				if err != nil {
					t.Fatal(err)
				}
				entries, err := os.ReadDir(filepath.Join(config.LedgerDir(f.workspace), "journal"))
				if err != nil {
					t.Fatal(err)
				}
				renewals := 0
				for _, entry := range entries {
					raw, err := os.ReadFile(filepath.Join(config.LedgerDir(f.workspace), "journal", entry.Name()))
					if err != nil {
						t.Fatal(err)
					}
					renewals += strings.Count(string(raw), `"event_type":"check.hold-renewed"`)
				}
				if renewals < 2 {
					t.Fatalf("renewals=%d, want multiple accepted renewals", renewals)
				}
			}
		})
	}
}
