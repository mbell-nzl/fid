# Phase 9 Escalations (archived)

> Phase: 9 - Concerns + v0.1c closure. Status: CLOSED 2026-06-21.

## ESC-006 - Phase 9 final review blockers: benchmark gate failure and concern address read-back gap

Status: **RESOLVED** at Phase 9 wrap (2026-06-21).
Source: Phase 9 final review / TL-Lite TRG-002.
Raised: 2026-06-20. Direction issued: Tech Lead 2026-06-20. Verified: QA Review 2.

Summary:
- The required phase benchmark gate `go test -bench . -benchmem ./...` failed in
  `cmd/fid` at `TestConcernListFiltersByIssueAndStatus/all_concerns`.
- VAL-004 / QA-12 failed public read-back: after `fid concern address`, read-back
  showed `status="addressed"` but neither `fid concern list` nor
  `fid issue get --section concerns` exposed `addressed_in_completion_id`.

Tech Lead diagnosis:
- Read-back gap: the write path (`concern_address.go`) persisted
  `addressed_in_completion_id`, but `ConcernListItem` (`concern_list.go`) and
  `IssueConcern` (`detail.go`) never selected/exposed it.
- Benchmark flake: both concern read queries `ORDER BY created_at DESC, id DESC`;
  second-resolution `created_at` ties fell back to the content-hash `id`, ordering
  rows nondeterministically vs creation order (same class as DEC-021).

Resolution:
- **DEC-023** — expose `addressed_in_completion_id` (nullable) on both concern
  read models.
- **DEC-024** — order concern reads by journal append order (`rowid DESC`,
  rebuild-stable; mirrors DEC-021).
- **Task 4** (`tasks/task-4.md`, commit `c927de0`) implemented both under TDD; the
  required benchmark gate exited 0 and VAL-004 PASSED while VAL-008/009 stayed
  PASS. QA Review 2 (`qa/review-2.md`) verified the fix (CONTINUE, 0 blocking).
- VAL-011 (human dogfood) was tracked on the human checklist and resolved 5/5 PASS
  (HUMAN-011a..e), not a Task 4 deliverable.

Outcome: all eleven Phase 9 VAL assertions stand PASS; ESC-006 closed.

## Other

No other escalations were raised during Phase 9. TRG-001 (Task 3 expanded write
set) resolved RESUME without escalation; see `triage-log.md`.
