package status

import (
	"testing"
	"time"

	"github.com/mbell-nzl/fid/internal/backup"
	"github.com/mbell-nzl/fid/internal/display"
	"github.com/mbell-nzl/fid/internal/holds"
	"github.com/mbell-nzl/fid/internal/issues"
	"github.com/mbell-nzl/fid/internal/next"
)

// Mark's approved step 6 mock-up pins the layout independently of the
// ledger fixtures, including its illustrative progress counts.
func TestReadyApprovedMock(t *testing.T) {
	t.Parallel()
	task := func(id, title string) issues.OverviewTask { return issues.OverviewTask{ID: id, Title: title} }
	export := task("FID-71a0c2", "Export validation")
	export.Status = "in_progress"
	contract := task("FID-1f2e3d", "Contract fields and exit codes")
	contract.Status = "in_progress"
	retry := task("FID-c2d4e5", "Retry a rejected backup push")
	login := task("FID-9f13aa", "Harden login rate limiting")
	backupAge := task("FID-0d4e8b", "Document backup age")
	golden := task("FID-2f2e3d", "Golden tests for the contract")
	rework := task("FID-3f2e3d", "Scenario test for rework")
	version := task("FID-4f2e3d", "Version the contract")
	ops := textEpic("FID-e3a1c0", issues.EpicInProgress, 1, 4, export)
	ops.Title = "Operating experience"
	loop := textEpic("FID-8a1b2c", issues.EpicInProgress, 0, 4, contract)
	loop.Title = "Loop contract"
	opsReady, loopReady := ops, loop
	opsReady.Tasks, loopReady.Tasks = []issues.OverviewTask{retry, login, backupAge}, []issues.OverviewTask{golden, rework, version}
	blocked, after := next.ReasonBlocked, next.ReasonNotNextInSequence
	result := Result{GeneratedAt: time.Date(2026, 11, 11, 9, 10, 0, 0, nzdt), Backup: Backup{State: string(backup.HealthNeedsYou)},
		InProgress: issues.OverviewSection{Epics: []issues.OverviewEpic{ops, loop}},
		Ready:      issues.OverviewSection{Epics: []issues.OverviewEpic{opsReady, loopReady}, Tasks: []issues.OverviewTask{task("FID-77aa01", "Tidy help text")}},
		Next:       &login.ID, entries: map[string]next.Entry{
			contract.ID:  {Working: &holds.Working{Worktree: "/home/mark/projects/fid-wt2", Lapsed: true}},
			retry.ID:     {OverviewTask: issues.OverviewTask{BlockedBy: []string{export.ID}}, Reason: &blocked},
			backupAge.ID: {Reason: stringPtr(next.ReasonHeldElsewhere), Working: &holds.Working{Worktree: "/home/mark/projects/fid-wt3"}},
			golden.ID:    {Reason: &after, After: &contract.ID}, rework.ID: {Reason: &after, After: &contract.ID}, version.ID: {Reason: &after, After: &contract.ID},
		}}
	const want = `╭ fid status ────────────────────────────────────────────── 2026-11-11 ╮

  IN PROGRESS
  ◆ FID-e3a1c0  Operating experience              ██░░░░░░░░  1/4 done
      ● FID-71a0c2  Export validation                      in progress
  ◆ FID-8a1b2c  Loop contract                     ░░░░░░░░░░  0/4 done
      ● FID-1f2e3d  Contract fields and exit …  interrupted in fid-wt2

  READY
  ◆ FID-e3a1c0  Operating experience                           3 to do
      ○ FID-c2d4e5  Retry a rejected backup pu…  blocked by FID-71a0c2
      ○ FID-9f13aa  Harden login rate limiting                    next
      ○ FID-0d4e8b  Document backup age             working in fid-wt3
  ◆ FID-8a1b2c  Loop contract                                  3 to do
      ○ FID-2f2e3d  Golden tests for the contract     after FID-1f2e3d
      ○ FID-3f2e3d  Scenario test for rework          after FID-1f2e3d
      ○ FID-4f2e3d  Version the contract              after FID-1f2e3d
  ○ FID-77aa01  Tidy help text
╰──────────────────────────────────────────────────────────────────────╯
`
	text := writeText(t, result, display.Renderer{}, false)
	if text != want {
		t.Fatalf("approved mock differs:\n%s\nwant:\n%s", text, want)
	}
	assertNoKernelWords(t, text)
}

func stringPtr(value string) *string { return &value }
