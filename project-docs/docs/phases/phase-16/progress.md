# Phase 16 Progress

> Phase: 16 - Activity Timeline + Read-Only Escape Hatches (roadmap v0.1e,
> slice 3 of 4)
> Status: RUNNING. Tasks 1-4 and Reviews 1-2 complete. Task 5 implementation is
> blocked from close until DEC-046 internal-activity fixture remediation lets the
> required repo-wide gate pass. DEC-047 requires TL-Lite triage before resume.

## Plan

Phase 16 ships the human-facing activity timeline and read-only escape hatches:

- `fid activity list`: projected `activity_log` read with entity, event, and
  date filters.
- CF-011 reconciliation: backfill unlogged issue-state writers so the timeline
  has no silent gaps and rebuilds deterministically.
- `fid sql`: raw read-only SQL escape hatch with classifier rejection and a
  read-only SQLite connection.
- `fid journal inspect <event-id>`: journal event explanation with canonical
  hash and revert verdict.
- CF-016: `fid issue get --section gates`, a first-class issue-scoped read of
  raw `gate_runs` history.

No schema migration is planned. CF-011 is additive: new events carry extra
`activity` effects; old events replay unchanged.

## Task Index

| Task | Title | Status |
|---|---|---|
| 1 | `fid activity list` projected timeline read | Complete |
| 2 | CF-011 backfill part 1 - create / ac add / ac verify / complete | Complete |
| 3 | CF-011 backfill part 2 + timeline-consistency e2e | Complete |
| 4 | `fid sql` read-only escape hatch | Complete |
| 5 | `fid journal inspect <event-id>` | Blocked; DEC-046 fixture remediation required before close |
| 6 | CF-016 `fid issue get --section gates` | Pending |

## Validation

Validation contract at `qa/validation-contract.md` includes VAL-001..013 and is
READY. QA plan at `qa/qa-plan.md` is READY. Heavy Mode cadence plans mid-phase
reviews after Tasks 2 and 4, then final review after Task 6.

## Progress History

- Kickoff (2026-07-05): human confirmed DC-A (CF-011 scope/vocabulary), DC-B
  (include CF-016 as Task 6), and Heavy Mode. Recorded as DEC-040..042. Phase
  Runner may start with Task 1; Tasks 4 and 5 are also independent of CF-011.
- Task 1 complete (2026-07-05): `fid activity list` implemented; report
  recorded at `tasks/task-1-done.md`, commit `977974c`.
- Task 2 complete (2026-07-05): CF-011 create/ac add/ac verify/complete
  backfill implemented; report recorded at `tasks/task-2-done.md`, commit
  `931d83a`.
- Task 3 complete (2026-07-05): CF-011 concern/unclaim/link add/dev gate and
  VAL-007 lifecycle/rebuild proof complete; report recorded at
  `tasks/task-3-done.md`, commits `fd98a02` and docs finalization `f10897f`.
- Task 4 blocked (2026-07-05): `fid sql` implementation and tests are present,
  but required global gates fail in `cmd/fid/activity_test.go` because the test
  seeds a `2026-07-04` journal event after `fid init` created a `2026-07-05`
  journal file. DEC-044 authorizes narrow do-now fixture remediation before
  Task 4 finalization.
- Workflow correction (2026-07-05): The Phase 16 QA plan required Review 1 after
  Tasks 1-2, but the ledger has no REVIEW entry and no `qa/review-1.md` exists.
  DEC-045 requires the runner to run overdue Review 1 before more Developer
  work, then run Review 2 after Task 4 finalizes and before Task 5.
- Reviews 1-2 and Task 4 complete (2026-07-05): Review 1 covered Tasks 1-2;
  Task 4 closed with commit `9cae372`; Review 2 covered Tasks 3-4.
- Task 5 blocked (2026-07-05): `fid journal inspect` implementation and tests
  are present, but required repo-wide gates fail in
  `internal/activity/activity_test.go` due to the same backdated seed fixture
  family as DEC-044. DEC-046 authorizes narrow fixture remediation.
- Workflow correction (2026-07-05): Ledger review found no TRIAGE entries and no
  `triage-log.md` despite multiple task-level STOPs. DEC-047 restores TL-Lite as
  first-step triager for task blockers before STOP/resume.

## Carry-Forward

- CF-008 re-deferred to v0.2 (resolution-reason field).
- CF-011 scheduled for Phase 16 Tasks 2-3; remains open until VAL-005/006/007
  pass.
- CF-016 scheduled for Phase 16 Task 6; remains open until VAL-012 passes.
- CF-017 re-deferred to v0.2c (`internal/dev` function-size cleanup).

## Escalations

- No open escalations at kickoff.
