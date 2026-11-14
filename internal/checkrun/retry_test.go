package checkrun

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/mbell-nzl/fid/internal/bootstrap"
	"github.com/mbell-nzl/fid/internal/checkrecord"
	"github.com/mbell-nzl/fid/internal/gittest"
	"github.com/mbell-nzl/fid/internal/holdwork"
	"github.com/mbell-nzl/fid/internal/issues"
)

func retryFixture(t *testing.T) Options {
	t.Helper()
	root := gittest.NewWorkspace(t)
	migrations := filepath.Join(gittest.ModuleRoot(), "migrations")
	if _, err := bootstrap.InitializeWorkspace(bootstrap.InitOptions{WorkspaceRoot: root, MigrationsDir: migrations, Prefix: "FID"}); err != nil {
		t.Fatal(err)
	}
	file, err := os.OpenFile(filepath.Join(root, ".fid.toml"), os.O_APPEND|os.O_WRONLY, 0)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := file.WriteString("\n[checks.test]\nrun=['true']\ntimeout='1m'\n[checks.ztest]\nrun=['true']\ntimeout='1m'\n[lanes.standard]\nchecks=[]\n"); err != nil {
		t.Fatal(err)
	}
	if err := file.Close(); err != nil {
		t.Fatal(err)
	}
	for _, args := range [][]string{
		{"add", "."},
		{"-c", "user.name=fid test", "-c", "user.email=fid@example.test", "commit", "-qm", "Check policy"},
		{"remote", "add", "origin", root},
		{"update-ref", "refs/remotes/origin/trunk", "HEAD"},
		{"symbolic-ref", "refs/remotes/origin/HEAD", "refs/remotes/origin/trunk"},
	} {
		if output, err := exec.Command("git", append([]string{"-C", root}, args...)...).CombinedOutput(); err != nil {
			t.Fatalf("git %v: %v: %s", args, err, output)
		}
	}
	lane := "standard"
	task, err := issues.Create(issues.CreateOptions{WorkspaceRoot: root, MigrationsDir: migrations, Kind: "task", Lane: &lane, Title: "Concurrent retries", Actor: "human:mark", IdempotencyKey: "task"})
	if err != nil {
		t.Fatal(err)
	}
	held, err := holdwork.Takeover(holdwork.Options{WorkspaceRoot: root, MigrationsDir: migrations, TaskID: task.ID, Actor: "human:mark", Reason: "Build fixture"})
	if err != nil {
		t.Fatal(err)
	}
	options := Options{WorkspaceRoot: root, MigrationsDir: migrations, HoldID: held.Fields["hold_id"].(string), Actor: "human:mark", Names: []string{"test"}, IdempotencyKey: "check", Version: "test"}
	if _, err := holdwork.Start(holdwork.Options{WorkspaceRoot: root, MigrationsDir: migrations, HoldID: options.HoldID, Actor: options.Actor, IdempotencyKey: "start"}); err != nil {
		t.Fatal(err)
	}
	return options
}

func incomplete(err error) bool {
	var refused *checkrecord.Refusal
	return errors.As(err, &refused) && refused.Code == "check-incomplete"
}

func TestSameKeyCompetingAdmission(t *testing.T) {
	t.Parallel()
	for _, row := range []struct {
		name      string
		before    int
		interrupt bool
	}{
		{"complete before first lock", 1, false},
		{"incomplete before first lock", 1, true},
		{"complete between locks", 2, false},
		{"incomplete between locks", 2, true},
	} {
		t.Run(row.name, func(t *testing.T) {
			t.Parallel()
			options := retryFixture(t)
			other := options
			cut := errors.New("interrupted before spawning")
			if row.interrupt {
				other.Before = func(string) error { return cut }
			}
			var competing checkrecord.Response
			var competingErr error
			locks := 0
			options.beforeWrite = func() {
				locks++
				if locks == row.before {
					competing, competingErr = Run(context.Background(), other)
				}
			}
			response, err := Run(context.Background(), options)
			switch {
			case competingErr == nil:
				competing.Replayed = true
				if err != nil || !reflect.DeepEqual(response, competing) {
					t.Fatalf("completed competitor=%#v; retry=%#v/%v", competing, response, err)
				}
			case errors.Is(competingErr, cut):
				if !incomplete(err) || !reflect.DeepEqual(response.CheckIDs, competing.CheckIDs) {
					t.Fatalf("interrupted competitor=%#v; retry=%#v/%v", competing, response, err)
				}
			case incomplete(competingErr):
				if err != nil || !response.Complete || !reflect.DeepEqual(response.CheckIDs, competing.CheckIDs) {
					t.Fatalf("reserved request=%#v/%v; competitor=%#v/%v", response, err, competing, competingErr)
				}
			default:
				t.Fatalf("unexpected competitor error: %v", competingErr)
			}
			if len(response.CheckIDs) != 1 {
				t.Fatalf("retry invented starts: %#v", response)
			}
			if err := other.write(func(db store) error {
				var starts int
				if err := db.conn.QueryRow(`SELECT count(*) FROM check_starts`).Scan(&starts); err != nil {
					return err
				}
				if starts != 1 {
					t.Fatalf("durable starts=%d, want 1", starts)
				}
				return nil
			}); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestSameKeyRetryAtResultAndBetweenChecksStop(t *testing.T) {
	t.Parallel()
	for _, stopping := range []bool{false, true} {
		t.Run(map[bool]string{false: "result", true: "stop"}[stopping], func(t *testing.T) {
			t.Parallel()
			options := retryFixture(t)
			if stopping {
				options.Names = []string{"test", "ztest"}
			}
			retry := options
			afterFirst := false
			options.After = func(record checkrecord.Record) error {
				if !stopping {
					return nil
				}
				_, err := holdwork.Takeover(holdwork.Options{WorkspaceRoot: options.WorkspaceRoot, MigrationsDir: options.MigrationsDir, TaskID: record.Start.TaskID, Actor: "human:operator", Reason: "Stop between checks"})
				afterFirst = true
				return err
			}
			var pending checkrecord.Response
			var pendingErr error
			probed := false
			options.beforeWrite = func() {
				if probed {
					return
				}
				ready := afterFirst
				if !stopping {
					if err := retry.write(func(db store) error {
						var processes, results int
						if err := db.conn.QueryRow(`SELECT (SELECT count(*) FROM check_processes),(SELECT count(*) FROM check_results)`).Scan(&processes, &results); err != nil {
							return err
						}
						ready = processes == 1 && results == 0
						return nil
					}); err != nil {
						t.Fatal(err)
					}
				}
				if ready {
					probed = true
					pending, pendingErr = Run(context.Background(), retry)
				}
			}
			response, err := Run(context.Background(), options)
			if !probed || !incomplete(pendingErr) || !reflect.DeepEqual(pending.CheckIDs, response.CheckIDs) {
				t.Fatalf("pending=%#v/%v; final=%#v/%v", pending, pendingErr, response, err)
			}
			if stopping {
				var refused *checkrecord.Refusal
				if !errors.As(err, &refused) || refused.Code != "hold-not-current" || !reflect.DeepEqual(response.NotStarted, []string{"ztest"}) {
					t.Fatalf("stop=%#v/%v", response, err)
				}
			} else if err != nil {
				t.Fatal(err)
			}
			replayed, replayErr := Run(context.Background(), retry)
			response.Replayed = true
			if !reflect.DeepEqual(replayed, response) || (replayErr != nil) != stopping {
				t.Fatalf("replayed=%#v/%v, want %#v", replayed, replayErr, response)
			}
		})
	}
}
