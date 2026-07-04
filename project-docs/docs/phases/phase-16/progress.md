# Phase 16 Progress

> Phase: 16 - Activity Timeline + Read-Only Escape Hatches (roadmap v0.1e,
> slice 3 of 4)
> Status: READY FOR PHASE RUNNER. Heavy Mode, validation contract, QA plan, and
> six task specs are ready. Human confirmed DC-A, DC-B, and Heavy Mode on
> 2026-07-05 (workspace timezone).

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
| 1 | `fid activity list` projected timeline read | Pending |
| 2 | CF-011 backfill part 1 - create / ac add / ac verify / complete | Pending |
| 3 | CF-011 backfill part 2 + timeline-consistency e2e | Pending |
| 4 | `fid sql` read-only escape hatch | Pending |
| 5 | `fid journal inspect <event-id>` | Pending |
| 6 | CF-016 `fid issue get --section gates` | Pending |

## Validation

Validation contract at `qa/validation-contract.md` includes VAL-001..013 and is
READY. QA plan at `qa/qa-plan.md` is READY. Heavy Mode cadence plans mid-phase
reviews after Tasks 2 and 4, then final review after Task 6.

## Progress History

- Kickoff (2026-07-05): human confirmed DC-A (CF-011 scope/vocabulary), DC-B
  (include CF-016 as Task 6), and Heavy Mode. Recorded as DEC-040..042. Phase
  Runner may start with Task 1; Tasks 4 and 5 are also independent of CF-011.

## Carry-Forward

- CF-008 re-deferred to v0.2 (resolution-reason field).
- CF-011 scheduled for Phase 16 Tasks 2-3; remains open until VAL-005/006/007
  pass.
- CF-016 scheduled for Phase 16 Task 6; remains open until VAL-012 passes.
- CF-017 re-deferred to v0.2c (`internal/dev` function-size cleanup).

## Escalations

- No open escalations at kickoff.
