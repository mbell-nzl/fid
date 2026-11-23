package export

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mbell-nzl/fid/internal/config"
	"github.com/mbell-nzl/fid/internal/gittest"
	"github.com/mbell-nzl/fid/internal/journal"
)

// Representative semantic refusals use the real append and journal scanner;
// the full malformed matrix runs once at their common projection boundary.
func TestReviewDeliveryJournalPaths(t *testing.T) {
	t.Parallel()
	for _, family := range []string{"answer", "defer", "post"} {
		for _, mode := range []string{"append", "replay"} {
			t.Run(family+"/"+mode, func(t *testing.T) {
				t.Parallel()
				v, e := responseArchiveFixture(t, family == "defer")
				if family == "post" {
					v, e = postArchiveFixture(t)
				}
				e.EventID = "EVT-review-delivery"
				goodRaw, err := journal.MarshalCanonicalJSON(e)
				if err != nil {
					t.Fatal(err)
				}
				if family == "post" {
					e.Effects[0].Values["task_id"] = "FID-other"
				} else {
					e.Effects[len(e.Effects)-1].Values["finding_id"] = "FND-1"
				}
				root := gittest.NewWorkspace(t)
				dir := filepath.Join(config.LedgerDir(root), "journal")
				if err := os.MkdirAll(dir, 0o700); err != nil {
					t.Fatal(err)
				}
				file := filepath.Join(dir, "2026-11-10.jsonl")
				before := handoffProjectionSnapshot(t, v)
				var good journal.Event
				if err := good.UnmarshalJSON(goodRaw); err != nil {
					t.Fatal(err)
				}
				if mode == "append" {
					err = journal.WithLedgerLock(config.LedgerDir(root), journal.LockOwner{Command: "test"}, func(lock *journal.WriteLock) error {
						if _, err := journal.CommitEvent(lock, v.conn, e); err == nil || !strings.Contains(err.Error(), "canonical record") {
							t.Fatalf("append refusal=%v", err)
						}
						if handoffProjectionSnapshot(t, v) != before {
							t.Fatal("append refusal changed state")
						}
						if raw, err := os.ReadFile(file); !os.IsNotExist(err) && (err != nil || len(raw) != 0) {
							t.Fatalf("append refusal changed bytes: %s/%v", raw, err)
						}
						_, err := journal.CommitEvent(lock, v.conn, good)
						return err
					})
				} else {
					raw, marshalErr := journal.MarshalCanonicalJSON(e)
					if marshalErr != nil {
						t.Fatal(marshalErr)
					}
					raw = append(raw, '\n')
					if err := os.WriteFile(file, raw, 0o600); err != nil {
						t.Fatal(err)
					}
					if _, err := journal.ReplayJournal(v.conn, dir, nil); err == nil || !strings.Contains(err.Error(), "canonical record") {
						t.Fatalf("replay refusal=%v", err)
					}
					if handoffProjectionSnapshot(t, v) != before {
						t.Fatal("replay refusal changed state or checkpoint")
					}
					stored, err := os.ReadFile(file)
					if err != nil || !bytes.Equal(stored, raw) {
						t.Fatal("replay refusal changed journal")
					}
					if err := os.WriteFile(file, append(goodRaw, '\n'), 0o600); err != nil {
						t.Fatal(err)
					}
					_, err = journal.ReplayJournal(v.conn, dir, nil)
				}
				if err != nil {
					t.Fatal("valid write after refusal", err)
				}
			})
		}
	}
}
