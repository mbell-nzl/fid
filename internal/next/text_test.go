package next

import (
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/mbell-nzl/fid/internal/display"
	"github.com/mbell-nzl/fid/internal/issues"
)

// kernelWords are the kernel nouns and IDs that human output never names
// (RESET 1.7).
var kernelWords = regexp.MustCompile(`(?i)\b(envelopes?|activations?|obligations?|attempts?|attention items?|runs?|claims?)\b|ENV-|STP-|ATT-`)

// mockTask is a task as the approved mock-ups show it.
func mockTask(id, title, category, epic string, priority int, reason, after string, held *Working, blockedBy ...string) Entry {
	entry := Entry{OverviewTask: issues.OverviewTask{ID: id, Title: title, StatusCategory: category, Status: map[string]string{"unstarted": "todo", "started": "in_progress"}[category],
		EpicID: &epic, BlockedBy: append([]string{}, blockedBy...)}, Working: held}
	if priority > 0 {
		entry.Priority = &priority
	}
	if reason != "" {
		entry.Reason = &reason
	}
	if after != "" {
		entry.After = &after
	}
	return entry
}

func mockHold(worktree string, here bool, since, until string, lapsed bool) *Working {
	return &Working{HoldID: "HLD-77b2c0", Worktree: worktree, Here: here, Actor: "agent:claude-builder", Since: since, ExpiresAt: until, Lapsed: lapsed}
}

// mockQueue is the --explain mock-up's order; the served task is second.
func mockQueue(served *Working) []Entry {
	elsewhere := mockHold("/home/mark/projects/fid-wt2", false, "2026-11-10T08:40:12Z", "2026-11-10T09:40:12Z", false)
	return []Entry{
		mockTask("FID-c2d4e5", "Retry a rejected backup push", "unstarted", "FID-e3a1c0", 0, ReasonBlocked, "", nil, "FID-71a0c2"),
		mockTask("FID-9f13aa", "Harden login rate limiting", "unstarted", "FID-e3a1c0", 1, "", "", served),
		mockTask("FID-0d4e8b", "Document backup age", "unstarted", "FID-e3a1c0", 0, "", "", nil),
		mockTask("FID-1f2e3d", "Contract fields and exit codes", "started", "FID-8a1b2c", 0, ReasonHeldElsewhere, "", elsewhere),
		mockTask("FID-2f2e3d", "Golden tests for the contract", "unstarted", "FID-8a1b2c", 0, ReasonNotNextInSequence, "FID-1f2e3d", nil, "FID-1f2e3d"),
		mockTask("FID-3f2e3d", "Scenario test for rework", "unstarted", "FID-8a1b2c", 0, ReasonNotNextInSequence, "FID-1f2e3d", nil),
		mockTask("FID-77aa01", "Tidy help text", "unstarted", "", 2, "", "", nil),
	}
}

var (
	operating = &issues.WorkEpic{ID: "FID-e3a1c0", Title: "Operating experience", Progress: issues.Progress{Done: 0, Total: 3}}
	loop      = &issues.WorkEpic{ID: "FID-8a1b2c", Title: "Loop contract", Progress: issues.Progress{Done: 0, Total: 3}, Sequential: true}
)

// The text renders each approved mock-up exactly (RESET 1.4a brief, step 4).
func TestWriteTextMatchesTheApprovedMockUps(t *testing.T) {
	t.Parallel()
	at := func(clock string) time.Time {
		parsed, err := time.Parse(time.RFC3339, "2026-11-10T"+clock+"Z")
		if err != nil {
			t.Fatal(err)
		}
		return parsed
	}
	taken := mockHold("/home/mark/projects/fid", true, "2026-11-10T09:02:11Z", "2026-11-10T10:02:11Z", false)
	live := mockHold("/home/mark/projects/fid", true, "2026-11-10T08:40:12Z", "2026-11-10T09:40:12Z", false)
	resumed := mockHold("/home/mark/projects/fid", true, "2026-11-10T11:15:00Z", "2026-11-10T12:15:00Z", false)
	interrupted := mockHold("/home/mark/projects/fid", true, "2026-11-10T09:58:00Z", "2026-11-10T10:58:00Z", true)
	resume := func(held *Working) []Entry {
		return []Entry{mockTask("FID-1f2e3d", "Contract fields and exit codes", "started", "FID-8a1b2c", 0, "", "", held)}
	}
	stuck := []Entry{
		mockTask("FID-c2d4e5", "Retry a rejected backup push", "unstarted", "FID-e3a1c0", 0, ReasonBlocked, "", nil, "FID-71a0c2"),
		mockQueue(nil)[3], mockQueue(nil)[4], mockQueue(nil)[5],
	}
	for _, test := range []struct {
		name    string
		result  Result
		explain bool
		want    string
	}{
		{"fid next", Result{Kind: KindNew, Queue: mockQueue(nil), Served: 1, Epic: operating, Now: at("09:02:11")}, false, `
╭ fid next ──────────────────────────────────────────────── 2026-11-10 ╮

  ○ FID-9f13aa  Harden login rate limiting                       to do

  Epic        ◆ FID-e3a1c0  Operating experience              0/3 done
  Priority    P1

  Passed over 1 task ahead of this one
      → fid next --explain
╰──────────────────────────────────────────────────────────────────────╯
`},
		{"taken", Result{Kind: KindNew, Queue: mockQueue(taken), Served: 1, Epic: operating, Now: at("09:02:11"), Claimed: true}, false, `
╭ fid next ──────────────────────────────────────────────── 2026-11-10 ╮

  ○ FID-9f13aa  Harden login rate limiting                       to do

  Epic        ◆ FID-e3a1c0  Operating experience              0/3 done
  Priority    P1
  Working     here, since 09:02, until 10:02
      → fid show FID-9f13aa
╰──────────────────────────────────────────────────────────────────────╯
`},
		{"live hold", Result{Kind: KindResume, Queue: resume(live), Served: 0, Epic: loop, Now: at("09:00:00"), Claimed: true}, false, `
╭ fid next ──────────────────────────────────────────────── 2026-11-10 ╮

  ● FID-1f2e3d  Contract fields and exit codes             in progress

  Epic        ◆ FID-8a1b2c  Loop contract                     0/3 done
  Working     here, since 08:40, until 09:40

  Already yours in this worktree: carry on with it
      → fid show FID-1f2e3d
╰──────────────────────────────────────────────────────────────────────╯
`},
		{"interrupted", Result{Kind: KindResume, Queue: resume(resumed), Served: 0, Epic: loop, Now: at("11:15:00"), Claimed: true, Interrupted: interrupted}, false, `
╭ fid next ──────────────────────────────────────────────── 2026-11-10 ╮

  ● FID-1f2e3d  Contract fields and exit codes             in progress

  Epic        ◆ FID-8a1b2c  Loop contract                     0/3 done
  Working     here, since 11:15, until 12:15

  ⚠ Interrupted here: work stopped at 10:58 without a park.
    Carry on from what's in this worktree.
      → fid show FID-1f2e3d
╰──────────────────────────────────────────────────────────────────────╯
`},
		{"nothing ready", Result{Kind: KindNone, Queue: stuck, Served: -1, Now: at("09:02:11")}, false, `
╭ fid next ──────────────────────────────────────────────── 2026-11-10 ╮

  ○ Nothing ready here: 1 blocked, 1 being worked in fid-wt2,
    2 not next in their epic's sequence
      → fid next --explain
╰──────────────────────────────────────────────────────────────────────╯
`},
		{"explain", Result{Kind: KindNew, Queue: mockQueue(nil), Served: 1, Epic: operating, Now: at("09:02:11")}, true, `
╭ fid next ──────────────────────────────────────────────── 2026-11-10 ╮

  ○ FID-9f13aa  Harden login rate limiting                       to do

  Epic        ◆ FID-e3a1c0  Operating experience              0/3 done
  Priority    P1

  IN ORDER
  ○ FID-c2d4e5  Retry a rejected backup push     blocked by FID-71a0c2
  ○ FID-9f13aa  Harden login rate limiting                        next
  ○ FID-0d4e8b  Document backup age                              ready
  ● FID-1f2e3d  Contract fields and exit codes      working in fid-wt2
  ○ FID-2f2e3d  Golden tests for the contract         after FID-1f2e3d
  ○ FID-3f2e3d  Scenario test for rework              after FID-1f2e3d
  ○ FID-77aa01  Tidy help text                                   ready
╰──────────────────────────────────────────────────────────────────────╯
`},
	} {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			var text strings.Builder
			if err := WriteText(&text, test.result, display.Renderer{}, test.explain); err != nil {
				t.Fatal(err)
			}
			if got := text.String(); got != strings.TrimPrefix(test.want, "\n") {
				t.Fatalf("text =\n%s\nwant\n%s", got, strings.TrimPrefix(test.want, "\n"))
			}
			assertNoKernelWords(t, text.String())
		})
	}
}

func assertNoKernelWords(t *testing.T, text string) {
	t.Helper()
	if found := kernelWords.FindAllString(text, -1); len(found) > 0 {
		t.Fatalf("human output names kernel words %q:\n%s", found, text)
	}
}
