# Phase Brief - Between Phases After Phase 16

> **Purpose.** Default execution brief for the current project state. Replaced
> at each phase kickoff by the incoming Tech Lead.
>
> **Size target:** under 150 lines.
>
> **State:** Phase 16 is **WRAPPED**. No Developer task is active. Phase 17
> kickoff is the next Tech Lead action.

## Phase 16 Outcome

Phase 16 shipped the v0.1e activity timeline and read-only escape hatches:

- `fid activity list` over projected `activity_log` with `--on`, `--since`, and
  `--event` filters in all read formats.
- CF-011 broad activity reconciliation for the approved issue-scoped writers:
  create, AC add/verify, complete, concern raise/address, unclaim, link add, and
  `dev gate` / `gate_recorded`.
- `fid sql` read-only escape hatch with classifier rejection and read-only
  SQLite backstop for TRD-SEC-002.
- `fid journal inspect <event-id>` with event location, hash parity, effects,
  actor/command metadata, and revert verdict.
- CF-016 `fid issue get --section gates` over raw append-only `gate_runs`.

## Validation Status

Final review: `project-docs/docs/phases/phase-16/qa/final-review.md`

| Range | Outcome | Notes |
|---|---|---|
| VAL-001..012 | PASS | Deterministic QA and phase-level gates passed. |
| VAL-013 | CARRY-FORWARD | Human dogfood outcomes are not recorded in `qa/human-checklist.md`; see CF-020 / DEC-048. |

Human checklist support exists at
`project-docs/docs/phases/phase-16/qa/human-checklist-run-sheet.md`, but the
checklist rows and summary totals remain blank. Do not mark VAL-013 PASS without
file-backed human outcomes.

## Phase 16 Wrap Decisions

- DEC-048 archives Phase 16 as deterministic-QA complete while carrying
  VAL-013 forward as CF-020.
- CF-011 and CF-016 are closed.
- CF-018 and CF-019 capture non-blocking file-size follow-ups for
  `internal/sqlquery/sqlquery.go` and `internal/format/markdown.go`.
- TRG-046 RESUME was audited correct in hindsight. No `ESC-NNN` items are open.

## Next Phase Candidate

Roadmap Phase 17 is "Maintenance + stub guard + v0.1 closure":

- `fid maintenance check`
- `fid maintenance vacuum`
- stub-table CLI guards for governance entities not shipping until later
- v0.1 exit-criteria smoke over the full implemented surface

Phase 17 must account for CF-020 before v0.1 is declared closed. It may also
carry CF-018/CF-019 only if it reopens the relevant packages.

## Next Action

Start Phase 17 kickoff only after reading:

1. `project-docs/specs/10_IMPLEMENTATION_ROADMAP.md` Phase 17 entry.
2. Relevant sections named by that entry in PRD, TRD, schema, design, quality,
   and workflow specs.
3. `project-docs/docs/current/decisions_active.md`, especially CF-020.
4. Phase 16 wrap artifacts: `decisions.md`, `progress.md`,
   `qa/final-review.md`, and `qa/human-checklist.md`.
