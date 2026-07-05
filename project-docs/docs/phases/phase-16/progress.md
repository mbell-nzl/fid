# Phase 16 Progress

> Phase: 16 - Activity Timeline + Read-Only Escape Hatches (roadmap v0.1e,
> slice 3 of 4)
> Status: WRAPPED. Tasks 1-6, Reviews 1-2, final review, and Tech Lead wrap are
> complete. Deterministic validation is PASS for VAL-001..012. VAL-013 remains
> CARRY-FORWARD as CF-020 because the human checklist outcomes are not recorded.

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
| 5 | `fid journal inspect <event-id>` | Complete |
| 6 | CF-016 `fid issue get --section gates` | Complete |

## Validation

Validation contract at `qa/validation-contract.md` includes VAL-001..013 and is
terminal except for the human-only carry-forward item. Final review records
VAL-001..012 PASS and VAL-013 CARRY-FORWARD. Tech Lead wrap records the
VAL-013 disposition as CF-020 instead of converting it to an evidence-free PASS.

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
- TRG-046 and Task 5 complete (2026-07-05): TL-Lite returned RESUME because
  DEC-046 already authorized the fixture remediation. Task 5 closed with commits
  `12c5107` and `91ee35b`.
- Task 6 complete (2026-07-05): `fid issue get --section gates` shipped with
  VAL-012 evidence and commits `31fb3dd` / `56ab49f`.
- Final review complete (2026-07-05): PASS WITH IMPROVEMENTS. VAL-001..012
  PASS; VAL-013 CARRY-FORWARD to human checklist. Phase-level format, lint,
  coverage, integration, benchmark, and build gates passed in QA.
- Human run sheet prepared (2026-07-05): commit `89e014d` added
  `qa/human-checklist-run-sheet.md` and tester name, but checklist outcomes and
  summary totals remain blank.
- Tech Lead wrap (2026-07-05): DEC-048 records deterministic wrap, closes
  CF-011/CF-016, opens CF-018/CF-019/CF-020, records spec-sync verdicts, and
  audits TRG-046.

## Carry-Forward

- CF-008 re-deferred to v0.2 (resolution-reason field).
- CF-011 closed at Phase 16 wrap (VAL-005/006/007 PASS).
- CF-016 closed at Phase 16 wrap (VAL-012 PASS).
- CF-017 re-deferred to v0.2c (`internal/dev` function-size cleanup).
- CF-018 opened at Phase 16 wrap: split `internal/sqlquery/sqlquery.go` when the
  SQL package next reopens.
- CF-019 opened at Phase 16 wrap: split `internal/format/markdown.go` when
  formatter-heavy work next reopens.
- CF-020 opened at Phase 16 wrap: complete or replace Phase 16 human dogfood
  evidence for VAL-013 before v0.1 closure.

## Escalations

- No ESC-NNN escalations were raised. TRG-046 RESUME was audited correct at wrap.
