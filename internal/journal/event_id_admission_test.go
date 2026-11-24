package journal

import (
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/mbell-nzl/fid/internal/config"
)

// The writers allocate a dated Event ID as one past the highest sequence its
// date holds (nextSequence). An admitted ID that skips ahead would leave them
// a sequence admission refuses, up to the ceiling, and so stop every later
// write on that date. Admission takes only the next sequence of each date.
func TestDatedEventIDTakesItsDatesNextSequenceOnEveryJournalPath(t *testing.T) {
	t.Parallel()
	tomorrow := admissionAt.Add(24 * time.Hour)
	cases := []struct {
		name, id string
		// later is when the writes after the forgery happen: on the date
		// the forged ID belongs to.
		later time.Time
		next  []string
	}{
		{"at the ceiling", "EVT-20261110-999999999", admissionAt, []string{"EVT-20261110-000004", "EVT-20261110-000005"}},
		{"one below the ceiling", "EVT-20261110-999999998", admissionAt, []string{"EVT-20261110-000004", "EVT-20261110-000005"}},
		{"a skipped sequence", "EVT-20261110-000005", admissionAt, []string{"EVT-20261110-000004", "EVT-20261110-000005"}},
		{"a repeated sequence spelled differently", "EVT-20261110-3", admissionAt, []string{"EVT-20261110-000004", "EVT-20261110-000005"}},
		{"a later date at the ceiling", "EVT-20261111-999999999", tomorrow, []string{"EVT-20261111-000001", "EVT-20261111-000002"}},
	}
	for _, tc := range cases {
		for _, mode := range []string{"commit", "replay", "catch-up", "freshen", "rebuild"} {
			t.Run(tc.name+"/"+mode, func(t *testing.T) {
				t.Parallel()
				w := newAdmissionWorkspace(t, admissionPrefix(), mode)
				forged := admissionLaneRaise()
				forged.EventID = tc.id
				forged.Effects[2].Values["journal_event_id"] = tc.id
				conn := openReplayDB(t, filepath.Join(config.LedgerDir(w.root), "fid.sqlite"))
				before := currentProjectionSnapshot(t, conn)
				closeProjectDB(t, conn)
				err := w.admit(t, mode, forged)
				if err == nil || !strings.Contains(err.Error(), "Event ID") {
					t.Fatalf("%s = %v, want %s refused", mode, err, tc.id)
				}
				var corrupt *CorruptJournalError
				if mode != "commit" && (!errors.As(err, &corrupt) || corrupt.Line != len(admissionPrefix())+1) {
					t.Fatalf("%s = %v, want the forged line refused", mode, err)
				}
				conn = openReplayDB(t, filepath.Join(config.LedgerDir(w.root), "fid.sqlite"))
				after := currentProjectionSnapshot(t, conn)
				closeProjectDB(t, conn)
				if !reflect.DeepEqual(before, after) {
					t.Fatalf("%s refusal changed the projection", mode)
				}
				if mode != "commit" {
					return
				}
				journal, readErr := os.ReadFile(w.file)
				if readErr != nil {
					t.Fatal(readErr)
				}
				if !reflect.DeepEqual(journal, w.prefix) {
					t.Fatal("commit refusal changed the journal")
				}
				assertAllocatedWritesFollow(t, w, tc.later, tc.next)
			})
		}
	}
}

// The writers' own next ID passes every path, and the writes after it take
// the sequences after it, across the day boundary too.
func TestDatedEventIDControlLeavesLaterWritesTheirSequences(t *testing.T) {
	t.Parallel()
	for _, mode := range []string{"commit", "replay", "catch-up", "freshen", "rebuild"} {
		t.Run(mode, func(t *testing.T) {
			t.Parallel()
			w := newAdmissionWorkspace(t, admissionPrefix(), mode)
			if err := w.admit(t, mode, admissionLaneRaise()); err != nil {
				t.Fatalf("%s refused the next ID: %v", mode, err)
			}
			assertAllocatedWritesFollow(t, w, admissionAt, []string{"EVT-20261110-000005", "EVT-20261110-000006"})
			assertAllocatedWritesFollow(t, w, admissionAt.Add(24*time.Hour), []string{"EVT-20261111-000001", "EVT-20261111-000002"})
		})
	}
}

// assertAllocatedWritesFollow commits one new task per wanted ID with the
// Event ID left empty, at the given time, and checks the ID each takes.
// A full rebuild then admits the journal the writers left.
func assertAllocatedWritesFollow(t *testing.T, w admissionWorkspace, at time.Time, want []string) {
	t.Helper()
	for _, id := range want {
		issue := "FID-" + strings.ReplaceAll(id, "-", "")
		next := admissionIssueCreated()
		next.EventID, next.OccurredAt, next.Aggregate.ID = "", at.Add(time.Minute), issue
		next.Effects[0].PK["id"], next.Effects[0].Values["title"] = issue, "Task "+issue
		next.Effects[1].PK["id"], next.Effects[1].Values["ref_id"], next.Effects[1].Values["journal_event_id"] = "ACT-"+issue, issue, id
		conn := openReplayDB(t, filepath.Join(config.LedgerDir(w.root), "fid.sqlite"))
		var committed Event
		var err error
		lockErr := WithLedgerLock(config.LedgerDir(w.root), LockOwner{Command: "test"}, func(lock *WriteLock) error {
			committed, err = CommitEvent(lock, conn, next)
			return nil
		})
		closeProjectDB(t, conn)
		if lockErr != nil {
			t.Fatal(lockErr)
		}
		if err != nil || committed.EventID != id {
			t.Fatalf("write with an allocated ID = %q, %v; want %s", committed.EventID, err, id)
		}
	}
	if _, err := RebuildWorkspace(w.root, testProjectMigrationsFS(t)); err != nil {
		t.Fatalf("rebuild after the allocated writes: %v", err)
	}
}
